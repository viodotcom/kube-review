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

	log.Printf(branch)
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
		return true, nil
	case 404:
		return false, nil
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
		return true, nil
	case 404:
		return false, nil
	default:
		return false, fmt.Errorf("Got status code %d with body: %s", resp.Status, body)
	}

	return false, nil
}

// CFEnvironment holds information about a codefresh environment
type CFEnvironment struct {
	Name      string
	ID        string
	UpdatedAt time.Time
}

// CFResponse is a codefresh API response
type CFResponse struct {
	Doc []struct {
		Metadata struct {
			ID        string `json:"id"`
			Name      string `json:"name"`
			UpdatedAt string `json:"updated_at"`
		} `json:"metadata"`
	} `json:"docs"`
}

// CFClient is a client for the codefresh API
type CFClient struct {
	APIEndpoint string
	APIToken    string
	httpClient  *http.Client
}

// NewCFClient creates a new CFClient
func NewCFClient(apiEndpoint string, apiToken string) *CFClient {
	httpClient := http.Client{}
	return &CFClient{
		APIEndpoint: apiEndpoint,
		APIToken:    apiToken,
		httpClient:  &httpClient,
	}
}

// DeleteEnvironment deletes an environment using environments-v2 endpoint
func (cf *CFClient) DeleteEnvironment(name string) error {
	endpoint, err := url.Parse(cf.APIEndpoint)
	if err != nil {
		return err
	}

	endpoint.Path = path.Join(endpoint.Path, "api/environments-v2", name)
	req, err := http.NewRequest("DELETE", endpoint.String(), nil)
	if err != nil {
		return err
	}

	req.Header.Add("Authorization", cf.APIToken)
	resp, err := cf.httpClient.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if resp.StatusCode != 200 {
		return fmt.Errorf("Got status code %d with body: %s", resp.Status, body)
	}

	return nil
}

// K8sClient for kubernetes cluter
type K8sClient struct {
	ClientSet *kubernetes.Clientset
}

// K8sNamespace is a kubernete namespace
type K8sNamespace struct {
	Name              string
	UpdatedAt         *time.Time
	PullRequestNumber string
	BranchName        string
	RepositoryName    string
	RepositoryOwner   string
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
		LabelSelector: "app.kubernetes.io/name=cf-review-env",
	})
	if err != nil {
		return nil, err
	}

	nk8sNamespacesList := make([]K8sNamespace, 0)
	for _, namespace := range namespaceList.Items {
		var updatedAt time.Time
		updatedAtTimestamp := namespace.Labels["app.kubernetes.io/updated_at"]
		if updatedAtTimestamp != "" {
			updatedAtTimestamp, err := strconv.ParseInt(updatedAtTimestamp, 10, 64)
			if err != nil {
				log.Printf("Got bad data at updated_at label for namespace: %s", namespace.Name)
				continue
			}

			updatedAt = time.Unix(updatedAtTimestamp, 0)
		}

		nk8sNamespacesList = append(nk8sNamespacesList, K8sNamespace{
			UpdatedAt:         &updatedAt,
			Name:              namespace.Labels["app.kubernetes.io/instance"],
			BranchName:        namespace.Labels["app.kubernetes.io/branch_name"],
			PullRequestNumber: namespace.Labels["app.kubernetes.io/pull_request_number"],
			RepositoryName:    namespace.Labels["app.kubernetes.io/repository_name"],
			RepositoryOwner:   namespace.Labels["app.kubernetes.io/repository_owner"],
		})
	}
	return nk8sNamespacesList, nil
}

// run will list all namespaces that belong to cf-review-env
// if the name of the namespace matches and it's expired, both
// cf review environment and the namespace will be deleted.
func run(c *cli.Context) error {
	ghClient := NewGHClient(c.String("ghEndpoint"), c.String("ghUserName"), c.String("ghToken"))
	cfClient := NewCFClient(c.String("cfEndpoint"), c.String("cfToken"))
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
			log.Printf("error checking env status: %s", err)
		} else {
			merged = merged
		}

		expired := time.Since(*namespace.UpdatedAt) >= time.Duration(c.Int("expiration"))*time.Hour
		log.Printf("Name: %s Duration: %s, Expired: %t, Merged: %t",
			namespace.Name, time.Since(*namespace.UpdatedAt), expired, merged)
		if expired || merged {
			if !c.Bool("dryRun") {
				if err := cfClient.DeleteEnvironment(namespace.Name); err != nil {
					log.Printf("warn could not delete environment: %s", err)
				}
			}
			log.Printf("Environment deleted: %s", namespace.Name)
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
			Name:     "cfToken",
			Usage:    "codefresh api token",
			Required: true,
		},
		&cli.StringFlag{
			Name:  "cfEndpoint",
			Usage: "codefresh api endpoint",
			Value: "https://g.codefresh.io",
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
