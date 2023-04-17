- [Tutorial](#tutorial)
  * [Requirements](#requirements)
  * [EKS Cluster](#eks-cluster)
    + [Create](#create)
    + [Configure](#configure)
  * [Domain](#domain)
  * [certmanager](#certmanager)
    + [Install](#install)
    + [Configure](#configure-1)
  * [Nginx Ingress](#nginx-ingress)
    + [Install](#install-1)
    + [Configure](#configure-2)
    + [Vertical Pod autoscaling - VPA](#vertical-pod-autoscaling---vpa)
  * [Scaling From or To zero with Keda](#scaling-from-or-to-zero-with-keda)
  * [Deploying an environment](#deploying-an-environment)
  * [Kudos](#kudos)

# Tutorial

This tutorial follows the installation and use of *Kube-Review* in a clean **EKS Cluster**. Therefore, this document assumes the use of **AWS** as a cloud provider. That said, the same instructions can easily be adapted to other cloud providers. 

Finally, the tutorial assumes commands are running on a modern MacOS computer, and the commands are executed from the root of the repo.

Before we start, we should define this variable that control the domain name to be used along then tutorial.

```shell
export MY_DOMAIN=my-domain.io
```

## Requirements

To follow this guide you will need an **AWS Account** and the following software installed. Besides that, you will need **aws cli**, **eksctl**, **kubectl**, **helm** and **gettext**:

```shel
brew tap weaveworks/tap
brew install weaveworks/tap/eksctl awscli kubectl helm gettext rhash
```

## EKS Cluster

### Create

Now that we have satisfied all the requirements, we can start by creating the EKS Cluster. If you already have a cluster configured you can skip this step, just to make sure to use Kubernetes version `>= 1.16`.

The following command creates a EKS Cluster and node group named `tutorial`:

```shell
eksctl create cluster \
--name tutorial \
--version 1.24 \
--with-oidc \
--nodegroup-name tutorial \
--node-type m5.large \
--managed \
--region eu-west-1
```

### Configure

With the cluster ready we have to update the kube config file:

```shell
aws eks --region eu-west-1 update-kubeconfig --name tutorial
```

## Domain

With cluster running, we can go ahead and create the domain that will use, you can skip this if you already have a domain:

```shell
aws route53 create-hosted-zone --name ${MY_DOMAIN} --caller-reference ${MY_DOMAIN}
```

## certmanager

On this section we will install and configure **certmanager** to issue and manage SSL certificates for our domain. 

### Install

With cluster properly configured, the first thin we need is to install **certmanager**:

```shell
helm repo add jetstack https://charts.jetstack.io
helm repo update
helm install \
  cert-manager jetstack/cert-manager \
  --create-namespace
  --namespace cert-manager \
  --version v1.9.1 \
  --set installCRDs=true \
  --set 'extraArgs={--dns01-recursive-nameservers-only}'
```

You can check that the installation went well by running:

```shell
kubectl get pods --namespace cert-manager
```

### Configure

Before we configure **certmanager** itself, we first have to create the **AWS Policy** that will be used to create the validation record on **AWS Route53**:

```shell
aws iam create-policy \
--policy-name AmazonRoute53Domains-cert-manager \
--description "Policy required by cert-manager to be able to modify Route 53 when generating wildcard certificates using Lets Encrypt" \
--policy-document file://./docs/files/route_53_change_policy.json
```

With the policy in-place, we can create the user that will be used by **certmanager** and attach the above policy:

```shell
aws iam create-user --user-name eks-cert-manager-route53
POLICY_ARN=$(aws iam list-policies --query "Policies[?PolicyName==\`AmazonRoute53Domains-cert-manager\`].{ARN:Arn}" --output text)
aws iam attach-user-policy --user-name "eks-cert-manager-route53" --policy-arn $POLICY_ARN
aws iam create-access-key --user-name eks-cert-manager-route53 > $HOME/.aws/eks-cert-manager-route53
export EKS_CERT_MANAGER_ROUTE53_AWS_ACCESS_KEY_ID=$(awk -F\" "/AccessKeyId/ { print \$4 }" $HOME/.aws/eks-cert-manager-route53)
export EKS_CERT_MANAGER_ROUTE53_AWS_SECRET_ACCESS_KEY=$(awk -F\" "/SecretAccessKey/ { print \$4 }" $HOME/.aws/eks-cert-manager-route53)
```

Now that we have a user with the proper policy, we can create the **ClusterIssuer** and the **Certificate**. In order to be able to issue a valid Certificate, we need prove to let's encrypt that we own the domain. To do so, we will configure **certmanager** to perform a DNS validation:

```shell
export EKS_CERT_MANAGER_ROUTE53_AWS_SECRET_ACCESS_KEY_BASE64=$(echo -n "$EKS_CERT_MANAGER_ROUTE53_AWS_SECRET_ACCESS_KEY" | base64)
envsubst < ./docs/files/cert-manager-letsencrypt-aws-route53-clusterissuer.yaml | kubectl apply -f -
envsubst < ./docs/files/cert-manager-certificate.yaml  | kubectl apply -f -
```

You can check the status by running these commands:

```shell
kubectl describe clusterissuer --namespace cert-manager
kubectl describe certificate --namespace cert-manager
```

If everything went well, you should see something like:

```shell
   ...
    Status:
    Conditions:
        Last Transition Time:  2021-06-08T13:46:24Z
        Message:               Certificate is up to date and has not expired
        Observed Generation:   1
        Reason:                Ready
        Status:                True
        Type:                  Ready
    Not After:               2021-09-05T14:00:55Z
    Not Before:              2021-06-07T14:00:55Z
    Renewal Time:            2021-08-06T14:00:55Z
    Events:                    <none>
    ...
```

For more info about this, check the certmanager's [official docs](https://cert-manager.io/docs/configuration/acme/).

## Nginx Ingress

Now it's time to install and configure **nginx ingress**:

### Install

Here we install ingress and point it to use the ssl certificate generated by **certmanager**:

```shell
helm repo add ingress-nginx https://kubernetes.github.io/ingress-nginx
helm repo update
helm install \
  ingress-nginx ingress-nginx/ingress-nginx \
  --create-namespace
  --wait \
  --namespace nginx-ingress \
  --version 4.4.0 \
  --set rbac.create=true \
  --set controller.extraArgs.default-ssl-certificate=cert-manager/tutorial-prd
```

Check that everything went well:

```shell
kubectl get service -n nginx-ingress
```

### Configure

Now that we have **nginx ingress** running, we just have to create the wild card DNS record pointing to the **AWS NLB**:

```shell
export LOADBALANCER_HOSTNAME=$(kubectl get svc ingress-nginx-controller -n nginx-ingress -o jsonpath="{.status.loadBalancer.ingress[0].hostname}")
export CANONICAL_HOSTED_ZONE_NAME_ID=$(aws elb describe-load-balancers --query "LoadBalancerDescriptions[?DNSName==\`$LOADBALANCER_HOSTNAME\`].CanonicalHostedZoneNameID" --output text)
export HOSTED_ZONE_ID=$(aws route53 list-hosted-zones --query "HostedZones[?Name==\`${MY_DOMAIN}.\`].Id" --output text)
envsubst < ./docs/files/aws_route53-dns_change.json | aws route53 change-resource-record-sets --hosted-zone-id ${HOSTED_ZONE_ID} --change-batch=file:///dev/stdin
```

### Vertical Pod autoscaling - VPA

The VPA service is a mandatory component to install and for that you can follow these steps below:

```shell
git clone https://github.com/kubernetes/autoscaler.git
cd autoscaler/vertical-pod-autoscaler/
./hack/vpa-up.sh
```

For more information, you can also check this [guide](https://github.com/kubernetes/autoscaler/tree/master/vertical-pod-autoscaler).

We are using the VPA service with `updateMode: "Off"` by default for all containers, including the `kube-review` and `sidecar`. To change these settings, we recommend that you use the [customization](customization.md) page. You can see the file created for this project [here](../src/deploy/resources/base/vpa.yml).

## Scaling From or To zero with Keda

**NOTE: KEDA requires Kubernetes cluster version 1.24 and higher**

We implemented the [Keda - Kubernetes-based Event Driven Autoscaling](https://github.com/kedacore/keda) project as a component in the Kube Review project because the review environments are a temporal environment running for a few days, we considered that saving money is an essential decision in that case. Implementing Keda help us with the possibility of scaling from/to zero the environment through HTTP requests with the [HTTP Add-On](https://github.com/kedacore/http-add-on) project.

To install both components Keda and HTTP Add-On, you can follow their guides below:

**NOTE: We tested the following versions: Keda 2.10.1 and HTTP Add-on 0.4.1**

- [Keda](https://keda.sh/docs/2.10/deploy/)
- [HTTP Add-on](https://github.com/kedacore/http-add-on/blob/main/docs/install.md)

Using the Keda project in the review environments, we don't need to take care of an ingress setting per each review environment (namespace). We should move it to the Keda namespace where we will have a `wildcard` covering all review environments URLs, then per each review environment (namespace), we just need to have an `HTTPScaledObject` created and used by Keda to collect metrics to scaling up/down the environment checking the HTTP requests.

```yaml
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: kube-review-ingress
  namespace: keda
  annotations:
    kubernetes.io/ingress.class: nginx
spec:
  tls:
    - hosts:
        - '*.example.com'
  rules:
    - host: '*.example.com'
      http:
        paths:
          - path: /
            pathType: Prefix
            backend:
              service:
                name: keda-add-ons-http-interceptor-proxy
                port:
                  number: 8080
```

## Deploying an environment

Now that everything is installed and working, we just need to call the deploy script to actually deploy a review env.

Note that this example is using kustomize overlays to add a redis as a sidecar container. This example also demonstrates how one can use dynamic variables with overlays. The `LABEL` variable will be used to replace the version of redis to be used.

With this command we will deploy a container running Nginx as a review env:

```shell 
KR_ID=nginx \
KR_IMAGE=nginx:latest \
KR_DOMAIN="${MY_DOMAIN}" \
KR_CONTAINER_PORT="80" \
KR_OVERLAY_PATH=src/deploy/resources/example \
KR_OVERLAY_TARGET_DIR=example \
LABEL=6.2.1 \
src/deploy/deploy
```

If everything goes well you should see something like this at end:

```
Connection test executed successfully
Environment deployed with url: https://re-xxxx.my-domain.io
```

You can use that url to test the deployed environment.

## Kudos

This tutorial was inspired by [K8s and Harbor setup docs](https://ruzickap.github.io/k8s-harbor/part-03/#install-cert-manager).
