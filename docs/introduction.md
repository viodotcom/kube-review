# Introduction

## Review Environments

A review environment, or preview environment, is a step commonly applied during the development phase of a software product, usually before deploying to staging or production.

By using review environments, product development teams can deploy their applications during the development phase, updating the environment on each commit made to a branch.

With review environments, development teams can get quick feedback on their work by sharing a public url of the environment with it's peers.

On top of that, environments can also be used for other purposes, like automated or manual testing.

Environments can be deployed manually or through a CI/CD pipeline. Connecting the deployment to a CI/CD tool, makes possible to update the environment every time that a commit is pushed, which will speed up the feedback cycle.

When deployed across an organization, review environment will be deployed in the same across all teams. This removes the confusion and the work of maintaining adhoc infrastructure for each product.

### Infrascture

The *kube-review* project is based on Kubernetes. Every environment deployed runs in its own namespace.

Every environment deployed can expose a service port, which will be publicly exposed using **Nginx Ingress**. The ingress itself will be load balanced through a public load balancer. At this point we only support **AWS ELB**.

Therefore, the Load Balancer will send requests to the **Nginx Ingress**, which will route these requests to the correct container port based on the domain requested.

In order to support HTTPS, Let's Encrypt is used as ACME to automatically handle certificate issuing and renewing.

To use a well defined domain, users can choose any product for DNS resolution, we recommend **AWS Route 53**. Therefore, Route 53 will point the DNS entry to AWS ELB.

As we want to support a high number of review environments with HTTPS, we use WildCard domain certificates, in order to avoid hitting **Let's Encrypt** rate limiting.

To manage the autoscaling and Nodes on your Kubernetes cluster in a cost effective way we recommend the use of spot instances. You can use pure and simple **AWS Spot Instances** but we recommend the **Spotinst** product.

## Deploy Component

The Deploy component is the one that is deployed and executed on Kubernetes, running a Container Image built by the user.

In order to deploy the necessary resources and configurations for the app, we use a **Kustomize** based process that takes care of doing all the necessary lifting while allowing customization.

In order to make installation even simpler, the app is deployed through a shell script. The script is baked inside the *kube-review* public container image.

Using this image is the easiest way to deploy a review environment, specially from inside a CI/CD pipeline, as the image already contains all the necessary requirements.

Finally, the id of the review env is automatically generated based on the prefix *re*, the branch and a hash of the whole name. This way, we avoid issues with long branch names and DNS limitations, at same time that we achieve unique and yet readable names.

## Prune Component

The prune component is a Kubernetes cron job responsible by purging expired environments or environment on which the attached pull request or branch was already merged or removed.

The prune component will run every hour. Once it starts, it will scan all environments belonging to a *kube-review* environment. This is done through scanning namespaces that have the *kube-review* annotation.

Once it finds an environment that is deployed using *kube-review*, prune will get the metadata from the annotation and check if the PR is already merged or if the branch was deleted.

To check for merged or removed PRs and branches, prune will contact the source code repository API. Right now we only support **GitHub**. If the environment is ephemeral, it will be removed if the expired time has passed, 5 days by default, or if the branch or PR was deleted/merged, whatever happens first.

If the environment is non-ephemeral, it will never be deleted. So, in order to destroy a non-ephemeral environment it first has to be changed to ephemeral.