package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"path"
	"time"

	"github.com/gobwas/glob"
	cli "github.com/urfave/cli/v2"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

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

// NewCFCLient creates a new CFClient
func NewCFCLient(apiEndpoint string, apiToken string) *CFClient {
	httpClient := http.Client{}
	return &CFClient{
		APIEndpoint: apiEndpoint,
		APIToken:    apiToken,
		httpClient:  &httpClient,
	}
}

// EnvironmentsList lists all environments using environments-v2 endpoint
func (cf *CFClient) EnvironmentsList() ([]CFEnvironment, error) {
	endpoint, err := url.Parse(cf.APIEndpoint)
	if err != nil {
		return nil, err
	}

	endpoint.Path = path.Join(endpoint.Path, "api/environments-v2")
	req, err := http.NewRequest("GET", endpoint.String(), nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Authorization", cf.APIToken)
	resp, err := cf.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("Got status code %d with body: %s", resp.Status, string(body))
	}

	cfResponse := CFResponse{}
	if err := json.Unmarshal(body, &cfResponse); err != nil {
		return nil, err
	}

	environments := make([]CFEnvironment, 0)
	for _, environment := range cfResponse.Doc {
		t, err := time.Parse("2006-01-02T15:04:05.000Z", environment.Metadata.UpdatedAt)
		if err != nil {
			return nil, err
		}

		environments = append(environments, CFEnvironment{
			ID:        environment.Metadata.ID,
			Name:      environment.Metadata.Name,
			UpdatedAt: t,
		})
	}
	return environments, nil
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

func run(c *cli.Context) error {
	cfClient := NewCFCLient(c.String("cfEndpoint"), c.String("cfToken"))
	environments, err := cfClient.EnvironmentsList()
	if err != nil {
		log.Fatalf("error listing environments: %s", err)
	}

	k8sClient, err := NewK8sClient(c.String("k8sContextName"), c.String("k8sKubeconfig"))
	if err != nil {
		log.Fatalf("error creating k8s client: %s", err)
	}

	g := glob.MustCompile(c.String("name"))
	for _, environment := range environments {
		matches := g.Match(environment.Name)
		expired := time.Since(environment.UpdatedAt) >= time.Duration(c.Int("expiration"))*time.Hour
		log.Printf("Name: %s Duration: %s, Expired: %t, Matches: %t, Forced: %t",
			environment.Name, time.Since(environment.UpdatedAt), expired, matches, c.Bool("force"))
		if matches && (expired || c.Bool("force")) {
			if !c.Bool("dryRun") {
				if err := cfClient.DeleteEnvironment(environment.Name); err != nil {
					log.Printf("error deleting environment: %s", err)
					continue
				}
			}
			log.Printf("Environment deleted: %s", environment.Name)
			if !c.Bool("dryRun") {
				if err := k8sClient.DeleteNamespace(environment.Name); err != nil {
					log.Printf("error deleting k8s namespace: %s", err)
					continue
				}
			}

			log.Printf("K8s namespace deleted: %s", environment.Name)
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
			Required: true,
		},
		&cli.BoolFlag{
			Name:  "force",
			Usage: "force will ignore expiration check",
			Value: false,
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
	}

	app.Action = func(c *cli.Context) error {
		return run(c)
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
