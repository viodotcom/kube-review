package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"path"
	"strconv"
	"time"

	"github.com/gobwas/glob"
	cli "github.com/urfave/cli/v2"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

// GHClient is a client for the github API
type GHClient struct {
	APIEndpoint string
	UserName    string
	APIToken    string
	httpClient  *http.Client
}

// NewGHClient creates a new GHClient
func NewGHClient(apiEndpoint, userName, apiToken string) *GHClient {
	httpClient := http.Client{}
	return &GHClient{
		APIEndpoint: apiEndpoint,
		UserName:    userName,
		APIToken:    apiToken,
		httpClient:  &httpClient,
	}
}

// IsMerged returns true if a PR is open false otherwise
// if PR number is not present use the branch name, in that case
// checks if branch exists
func (gh *GHClient) IsMerged(owner, repo, number, branch string) (bool, error) {
	if number != "" {
		return gh.IsPRMerged(owner, repo, number)
	}

	if branch != "" {
		return gh.IsBranchMerged(owner, repo, branch)
	}

	return false, fmt.Errorf("Both PR number and branch name annotations are missing")
}

// IsBranchMerged returns true if the branch is merged false otherwise
func (gh *GHClient) IsBranchMerged(owner, repo, branch string) (bool, error) {
	endpoint, err := url.Parse(gh.APIEndpoint)
	if err != nil {
		return false, err
	}

	// GET /repos/{owner}/{repo}/branches/{branch}
	method := fmt.Sprintf("repos/%s/%s/branches/%s", owner, repo, branch)
	endpoint.Path = path.Join(endpoint.Path, method)
	req, err := http.NewRequest("GET", endpoint.String(), nil)
	if err != nil {
		return false, err
	}

	req.SetBasicAuth(gh.UserName, gh.APIToken)
	resp, err := gh.httpClient.Do(req)
	if err != nil {
		return false, err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return false, err
	}

	switch resp.StatusCode {
	case 200:
		return false, nil
	case 404:
		return true, nil
	default:
		return false, fmt.Errorf("Got status code %d with body: %s", resp.Status, body)
	}

	return false, nil
}

// IsPRMerged returns true if a PR is merged false otherwise
func (gh *GHClient) IsPRMerged(owner, repo, number string) (bool, error) {
	endpoint, err := url.Parse(gh.APIEndpoint)
	if err != nil {
		return false, err
	}

	// GET /repos/:owner/:repo/pulls/:pull_number/merge
	method := fmt.Sprintf("repos/%s/%s/pulls/%s/merge", owner, repo, number)
	endpoint.Path = path.Join(endpoint.Path, method)
	req, err := http.NewRequest("GET", endpoint.String(), nil)
	if err != nil {
		return false, err
	}

	req.SetBasicAuth(gh.UserName, gh.APIToken)
	resp, err := gh.httpClient.Do(req)
	if err != nil {
		return false, err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return false, err
	}

	switch resp.StatusCode {
	case 204:
		return false, nil
	case 404:
		return true, nil
	default:
		return false, fmt.Errorf("Got status code %d with body: %s", resp.Status, body)
	}

	return false, nil
}

// K8sClient for kubernetes cluter
type K8sClient struct {
	ClientSet *kubernetes.Clientset
}

// K8sNamespace is a kubernete namespace
type K8sNamespace struct {
	Name              string
	Instance          string
	UpdatedAt         *time.Time
	PullRequestNumber string
	BranchName        string
	RepositoryName    string
	RepositoryOwner   string
	IsEphemeral       bool
}

func buildConfigFromFlags(contextName, kubeconfigPath string) (*rest.Config, error) {
	return clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
		&clientcmd.ClientConfigLoadingRules{ExplicitPath: kubeconfigPath},
		&clientcmd.ConfigOverrides{
			CurrentContext: contextName,
		}).ClientConfig()
}

// NewK8sClient return a k8s client
func NewK8sClient(contextName, kubeconfig string) (K8sClient, error) {
	config, err := buildConfigFromFlags(contextName, kubeconfig)
	if err != nil {
		return K8sClient{}, err
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return K8sClient{}, err
	}

	return K8sClient{clientset}, nil
}

// DeleteNamespace deletes a namespace
func (k8s *K8sClient) DeleteNamespace(name string) error {
	// The name of the deployment and the namespace are the same
	namespacesClient := k8s.ClientSet.CoreV1().Namespaces()
	deletePolicy := metav1.DeletePropagationForeground
	if err := namespacesClient.Delete(name, &metav1.DeleteOptions{
		PropagationPolicy: &deletePolicy,
	}); err != nil {
		return err
	}
	return nil
}

// NamespaceList deletes a namespace
func (k8s *K8sClient) NamespaceList() ([]K8sNamespace, error) {
	namespacesClient := k8s.ClientSet.CoreV1().Namespaces()
	namespaceList, err := namespacesClient.List(metav1.ListOptions{
		LabelSelector: "app.kubernetes.io/name in (cf-review-env, kube-review)",
	})
	if err != nil {
		return nil, err
	}

	nk8sNamespacesList := make([]K8sNamespace, 0)
	for _, namespace := range namespaceList.Items {
		// Backward compatibility layer. This can be removed once all
		// envs using labels for storing data are purged.
		data := namespace.Annotations
		if _, ok := data["app.kubernetes.io/instance"]; !ok {
			data = namespace.Labels
		}

		var updatedAt time.Time
		updatedAtTimestamp := data["app.kubernetes.io/updated_at"]
		if updatedAtTimestamp != "" {
			updatedAtTimestamp, err := strconv.ParseInt(updatedAtTimestamp, 10, 64)
			if err != nil {
				log.Printf("Got bad data at updated_at annotation for namespace: %s", namespace.Name)
				continue
			}

			updatedAt = time.Unix(updatedAtTimestamp, 0)
		}

		isEphemeral := true
		s, err := strconv.ParseBool(data["app.kubernetes.io/is_ephemeral"])
		if err != nil {
			log.Printf("Got bad data at is_ephemeral annotation for namespace: %s", namespace.Name)
		} else {
			isEphemeral = s
		}

		nk8sNamespacesList = append(nk8sNamespacesList, K8sNamespace{
			UpdatedAt:         &updatedAt,
			Name:              namespace.Name,
			Instance:          data["app.kubernetes.io/instance"],
			BranchName:        data["app.kubernetes.io/branch_name"],
			PullRequestNumber: data["app.kubernetes.io/pull_request_number"],
			RepositoryName:    data["app.kubernetes.io/repository_name"],
			RepositoryOwner:   data["app.kubernetes.io/repository_owner"],
			IsEphemeral:       isEphemeral,
		})
	}
	return nk8sNamespacesList, nil
}

// run will list all namespaces that belong to kube-review
// if the name of the namespace matches and it's expired,
// the namespace will be deleted.
func run(c *cli.Context) error {
	ghClient := NewGHClient(c.String("ghEndpoint"), c.String("ghUserName"), c.String("ghToken"))
	k8sClient, err := NewK8sClient(c.String("k8sContextName"), c.String("k8sKubeconfig"))
	if err != nil {
		log.Fatalf("error creating k8s client: %s", err)
	}

	namespaces, err := k8sClient.NamespaceList()
	if err != nil {
		log.Fatalf("error listing namespaces: %s", err)
	}

	g := glob.MustCompile(c.String("name"))
	for _, namespace := range namespaces {
		matches := g.Match(namespace.Name)
		if !matches {
			continue
		}

		if namespace.UpdatedAt == nil {
			log.Printf("Namespace %s has no updated_at label", namespace.Name)
			continue
		}

		merged := false
		merged, err := ghClient.IsMerged(namespace.RepositoryOwner, namespace.RepositoryName,
			namespace.PullRequestNumber, namespace.BranchName)
		if err != nil {
			log.Printf("error checking env: %s error: %s", namespace.Name, err)
			continue
		} else {
			merged = merged
		}

		expired := time.Since(*namespace.UpdatedAt) >= time.Duration(c.Int("expiration"))*time.Hour
		log.Printf("Name: %s Duration: %s, Expired: %t, Merged: %t, Ephemeral: %t",
			namespace.Name, time.Since(*namespace.UpdatedAt), expired, merged, namespace.IsEphemeral)
		// If an env is not ephemeral we dont need to check
		// We never touch non ephemeral environments
		if !namespace.IsEphemeral {
			continue
		}

		if expired || merged {
			if !c.Bool("dryRun") {
				if err := k8sClient.DeleteNamespace(namespace.Name); err != nil {
					log.Printf("error deleting k8s namespace: %s", err)
					continue
				}
			}

			log.Printf("K8s namespace deleted: %s", namespace.Name)
		}
	}
	return nil
}

func main() {
	app := cli.NewApp()
	app.Flags = []cli.Flag{
		&cli.StringFlag{
			Name:     "name",
			Usage:    "environment name to filter, accepts glob expressions",
			Required: false,
			Value:    "*",
		},
		&cli.IntFlag{
			Name:  "expiration",
			Usage: "how many hous to consider an environment stale",
			Value: 120,
		},
		&cli.BoolFlag{
			Name:  "dryRun",
			Usage: "only show logs but don'r perform deletess",
			Value: false,
		},
		&cli.StringFlag{
			Name:     "k8sKubeconfig",
			Usage:    "absolute path to the kubeconfig file",
			Required: true,
		},
		&cli.StringFlag{
			Name:     "k8sContextName",
			Usage:    "the k8s context name to operate on",
			Required: true,
		},
		&cli.StringFlag{
			Name:  "ghEndpoint",
			Usage: "github api endpoint",
			Value: "https://api.github.com",
		},
		&cli.StringFlag{
			Name:     "ghToken",
			Usage:    "github api token",
			Required: true,
		},
		&cli.StringFlag{
			Name:     "ghUserName",
			Usage:    "github username to use for auth",
			Required: true,
		},
	}

	app.Action = func(c *cli.Context) error {
		return run(c)
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
