package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"path"
	"strings"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
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

// NewK8sClient return a k8s client
func NewK8sClient(kubeconfig string) (K8sClient, error) {
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		return K8sClient{}, err
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return K8sClient{}, err
	}

	return K8sClient{clientset}, nil
}

// DeleteDeployment deletes a deployment
func (k8s *K8sClient) DeleteDeployment(nameSpace string, name string) error {
	deploymentsClient := k8s.ClientSet.AppsV1().Deployments(nameSpace)
	deletePolicy := metav1.DeletePropagationForeground
	if err := deploymentsClient.Delete(name, &metav1.DeleteOptions{
		PropagationPolicy: &deletePolicy,
	}); err != nil {
		return err
	}
	return nil
}

// DeleteService deletes a service
func (k8s *K8sClient) DeleteService(nameSpace string, name string) error {
	servicesClient := k8s.ClientSet.CoreV1().Services(nameSpace)
	deletePolicy := metav1.DeletePropagationForeground
	if err := servicesClient.Delete(name, &metav1.DeleteOptions{
		PropagationPolicy: &deletePolicy,
	}); err != nil {
		return err
	}
	return nil
}

// DeleteIngress deletes an ingress
func (k8s *K8sClient) DeleteIngress(nameSpace string, name string) error {
	ingressesClient := k8s.ClientSet.ExtensionsV1beta1().Ingresses(nameSpace)
	deletePolicy := metav1.DeletePropagationForeground
	if err := ingressesClient.Delete(name, &metav1.DeleteOptions{
		PropagationPolicy: &deletePolicy,
	}); err != nil {
		return err
	}
	return nil
}

// DeleteEnvironment deletes an environment
func (k8s *K8sClient) DeleteEnvironment(nameSpace string, name string) error {
	/*err := k8s.DeleteDeployment(nameSpace, name)
	if err != nil {
		return err
	}*/

	err := k8s.DeleteService(nameSpace, name)
	if err != nil {
		return err
	}

	err = k8s.DeleteIngress(nameSpace, name)
	if err != nil {
		return err
	}

	return nil
}

func main() {
	// TODO figured how to select the cluster and implement filtering by cluster
	var prefix = flag.String("prefix", "", "prefix string to filter environments to check")
	var expiration = flag.Int("expirarion", 120, "how many hous to consider and environment stale")
	var dryRun = flag.Bool("dryRun", false, "only show logs but don'r perform deletess")
	var cfToken = flag.String("cfToken", "", "codefresh api token")
	var cfEndpoint = flag.String("cfEndpoint", "https://g.codefresh.io", "codefresh api endpoint")
	var kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	var k8sNamespace = flag.String("k8sNamespace", "", "the k8s namespace to operate on")
	flag.Parse()

	cfClient := NewCFCLient(*cfEndpoint, *cfToken)
	environments, err := cfClient.EnvironmentsList()
	if err != nil {
		log.Fatalf("error listing environments: %s", err)
	}

	k8sClient, err := NewK8sClient(*kubeconfig)
	if err != nil {
		log.Fatalf("error creating k8s client: %s", err)
	}

	for _, environment := range environments {
		expired := time.Since(environment.UpdatedAt) >= time.Duration(*expiration)*time.Hour
		match := strings.HasPrefix(environment.Name, *prefix)
		log.Printf("Name: %s Duration: %s, Expired: %t, Matches: %t",
			environment.Name, time.Since(environment.UpdatedAt), expired, match)
		if !expired || !match {
			continue
		}

		if !*dryRun {
			if err := cfClient.DeleteEnvironment(environment.Name); err != nil {
				log.Printf("error deleting environment: %s", err)
				continue
			}
		}
		log.Printf("Environment deleted: %s", environment.Name)
		if !*dryRun {
			if err := k8sClient.DeleteDeployment(*k8sNamespace, environment.Name); err != nil {
				log.Printf("error deleting k8s deployment: %s", err)
				continue
			}
		}

		log.Printf("K8s deployment deleted: %s", environment.Name)
	}
}
