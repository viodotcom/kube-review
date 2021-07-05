# Tutorial

This tutorial follows the installation and use of *kube-review* in a clean **EKS Cluster**. Therefore, this document assumes the use of **AWS** as a cloud provider. That said, the same instructions can easily be adapted to other cloud providers. 

Finally, the tutorial assumes commands are running on a modern MacOS computer, and the commands are executed from the root of the repo.

Before we start, we should define this variable that control the domain name to be used along then tutorial.

    export MY_DOMAIN=my-domain.io

## Requirements

To follow this guide you will need an **AWS Account** and the following software installed. Besides that, you will need **aws cli**, **eksctl**, **kubectl**, **helm** and **gettext**:

    brew tap weaveworks/tap
    brew install weaveworks/tap/eksctl awscli kubectl helm gettext

## EKS Cluster

### Create

Now that we have satisfied all the requirements, we can start by creating the EKS Cluster. If you already have a cluster configured you can skip this step, just to make sure to use Kubernetes version `>= 1.16`.

The following command creates a EKS Cluster and node group named `tutorial`:

    eksctl create cluster \
    --name tutorial \
    --version 1.19 \
    --with-oidc \
    --nodegroup-name tutorial \
    --node-type m5.large \
    --managed \
    --region eu-west-1

### Configure

With the cluster ready we have to update the kube config file:

    aws eks --region eu-west-1 update-kubeconfig --name tutorial

## Domain

With cluster running, we can go ahead and create the domain that will use, you can skip this if you already have a domain:

    aws route53 create-hosted-zone --name ${MY_DOMAIN} --caller-reference ${MY_DOMAIN}

## certmanager

On this section we will install and configure **certmanager** to issue and manage SSL certificates for our domain. 

### Install

With cluster properly configured, the first thin we need is to install **certmanager**:

    helm repo add jetstack https://charts.jetstack.io
    helm repo update
    kubectl create namespace cert-manager
    helm install \
        cert-manager jetstack/cert-manager \
        --namespace cert-manager \
        --version v1.3.1 \
        --set installCRDs=true \
        --set 'extraArgs={--dns01-recursive-nameservers-only}'

You can check that the installation went well by running:

    kubectl get pods --namespace cert-manager

### Configure

Before we configure **certmanager** itself, we first have to create the **AWS Policy** that will be used to create the validation record on **AWS Route53**:

    aws iam create-policy \
    --policy-name AmazonRoute53Domains-cert-manager \
    --description "Policy required by cert-manager to be able to modify Route 53 when generating wildcard certificates using Lets Encrypt" \
    --policy-document file://./docs/files/route_53_change_policy.json

With the policy in-place, we can create the user that will be used by **certmanager** and attach the above policy:

    aws iam create-user --user-name eks-cert-manager-route53
    POLICY_ARN=$(aws iam list-policies --query "Policies[?PolicyName==\`AmazonRoute53Domains-cert-manager\`].{ARN:Arn}" --output text)
    aws iam attach-user-policy --user-name "eks-cert-manager-route53" --policy-arn $POLICY_ARN
    aws iam create-access-key --user-name eks-cert-manager-route53 > $HOME/.aws/eks-cert-manager-route53
    export EKS_CERT_MANAGER_ROUTE53_AWS_ACCESS_KEY_ID=$(awk -F\" "/AccessKeyId/ { print \$4 }" $HOME/.aws/eks-cert-manager-route53)
    export EKS_CERT_MANAGER_ROUTE53_AWS_SECRET_ACCESS_KEY=$(awk -F\" "/SecretAccessKey/ { print \$4 }" $HOME/.aws/eks-cert-manager-route53)

Now that we have a user with the proper policy, we can create the **ClusterIssuer** and the **Certificate**. In order to be able to issue a valid Certificate, we need prove to let's encrypt that we own the domain. To do so, we will configure **certmanager** to perform a DNS validation:

    export EKS_CERT_MANAGER_ROUTE53_AWS_SECRET_ACCESS_KEY_BASE64=$(echo -n "$EKS_CERT_MANAGER_ROUTE53_AWS_SECRET_ACCESS_KEY" | base64)
    envsubst < ./docs/files/cert-manager-letsencrypt-aws-route53-clusterissuer.yaml | kubectl apply -f -
    envsubst < ./docs/files/cert-manager-certificate.yaml  | kubectl apply -f -

You can check the status by running these commands:

    kubectl describe clusterissuer --namespace cert-manager
    kubectl describe certificate --namespace cert-manager

If everything went well, you should see something like:

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

For more info about this, check the certmanager's [official docs](https://cert-manager.io/docs/configuration/acme/).

## Nginx Ingress

Now it's time to install and configure **nginx ingress**:

### Install

Here we install ingress and point it to use the ssl certificate generated by **certmanager**:

    helm repo add ingress-nginx https://kubernetes.github.io/ingress-nginx
    helm repo update
    kubectl create namespace nginx-ingress
    helm install \
        ingress-nginx ingress-nginx/ingress-nginx  \
        --wait \
        --namespace nginx-ingress \
        --version 3.8.0 \
        --set rbac.create=true \
        --set controller.extraArgs.default-ssl-certificate=cert-manager/tutorial-prd

Check that everything went well:
    
    kubectl get service -n nginx-ingress

### Configure

Now that we have **nginx ingress** running, we just have to create the wild card DNS record pointing to the **AWS NLB**:

    export LOADBALANCER_HOSTNAME=$(kubectl get svc ingress-nginx-controller -n nginx-ingress -o jsonpath="{.status.loadBalancer.ingress[0].hostname}")
    export CANONICAL_HOSTED_ZONE_NAME_ID=$(aws elb describe-load-balancers --query "LoadBalancerDescriptions[?DNSName==\`$LOADBALANCER_HOSTNAME\`].CanonicalHostedZoneNameID" --output text)
    export HOSTED_ZONE_ID=$(aws route53 list-hosted-zones --query "HostedZones[?Name==\`${MY_DOMAIN}.\`].Id" --output text)

    envsubst < ./docs/files/aws_route53-dns_change.json | aws route53 change-resource-record-sets --hosted-zone-id ${HOSTED_ZONE_ID} --change-batch=file:///dev/stdin

## Deploying an environment

Now that everything is installed and working, we just need to call the deploy script to actually deploy a review env.

Note that this example is using kustomize overlays to add a redis as a sidecar container. This example also demonstrates how one can use dynamic variables with overlays. The `LABEL` variable will be used to replace the version of redis to be used.

With this command we will deploy a container running Nginx as a review env:
    
    KR_ID=nginx \
    KR_IMAGE_URL=nginx \
    KR_IMAGE_TAG=latest \
    KR_DOMAIN="${MY_DOMAIN}" \
    KR_CONTAINER_PORT="80" \
    KR_OVERLAY_PATH=src/deploy/resources/example \
    LABEL=6.2.1 \
    src/deploy/deploy

If everything goes well you should see something like this at end:

    Connection test executed successfully
    Environment deployed with url: https://re-xxxx.my-domain.io

You can use that url to test the deployed environment.

## Kudos

This tutorial was inspired by [K8s and Harbor setup docs](https://ruzickap.github.io/k8s-harbor/part-03/#install-cert-manager).