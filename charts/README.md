# Template Helm Chart

This is a `template helm chart` created to be used in the CI pipeline for all projects.

Ref.: https://helm.sh/docs/helm/helm_create/

## How to Install

This is an example of `how to install`.

```
helm install CHART_NAME cf-review-env \
--install \
--reset-values \
--repo REPO_NAME \
--version CHART_VERSION \
--namespace NAMESPACE_NAME \
--values values.yaml \
--set envFrom.secretRef.name=SECRET_NAME \
--set image.repository=IMAGE_REPOSITORY_NAME \
--set image.tag=IMAGE_TAG \
--set imagePullSecrets=IMAGE_PULL_SECRETS_NAME \
--set "ingress.hosts[0].host=HOST_NAME" \
--set "ingress.hosts[0].paths[0]=/" \
--set "ingress.tls[0].hosts[0]=HOST_NAME" \
--set "ingress.tls[0].secretName=SECRET_TLS_NAME" \
--wait
```

## How to Upgrade

This is an example of `how to upgrade`.

```
helm upgrade CHART_NAME cf-review-env \
--install \
--reset-values \
--repo REPO_NAME \
--version 0.2.0 \
--namespace NAMESPACE_NAME \
--values values.yaml \
--set envFrom.secretRef.name=SECRET_NAME \
--set image.repository=IMAGE_REPOSITORY_NAME \
--set image.tag=IMAGE_TAG \
--set imagePullSecrets=IMAGE_PULL_SECRETS_NAME \
--set "ingress.hosts[0].host=HOST_NAME" \
--set "ingress.hosts[0].paths[0]=/" \
--set "ingress.tls[0].hosts[0]=HOST_NAME" \
--set "ingress.tls[0].secretName=SECRET_TLS_NAME" \
--wait
```

## How to Delete

This is an example of `how to delete`.

`helm delete CHART_NAME --namespace NAMESPACE_NAME`

## How to debut the templates files

Ref.: https://helm.sh/docs/chart_template_guide/debugging/

When will be necessary to add new changes on the templates file it's possible to do a debug and then to check if the new change is correct or not, basically you need to run this command:

`helm template test cf-review-env -f values.yaml --debug`

**Note:** It's necessary to use a true `values.yaml` file.
