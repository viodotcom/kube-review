# How to use the Helm Chart

These are the basic commands to install the helm chart from the local folder.

## How to Install

REF.: https://helm.sh/docs/helm/helm_install/

This is an example of `how to install`.

```
helm install CHART_NAME ./kube-review \
--install \
--reset-values \
--namespace NAMESPACE_NAME \
--values values.yaml \
--set envFrom.secretRef.name=SECRET_NAME \
--set image.repository=IMAGE_REPOSITORY_URL/IMAGE_REPOSITORY_NAME \
--set image.tag=IMAGE_TAG \
--set imagePullSecrets=IMAGE_PULL_SECRETS_NAME \
--set "ingress.hosts[0].host=HOST_NAME.shared-prod.fih.io" \
--set "ingress.hosts[0].paths[0]=/" \
--set "ingress.tls[0].hosts[0]=HOST_NAME.shared-prod.fih.io" \
--wait
```

## How to Upgrade

REF: https://helm.sh/docs/helm/helm_upgrade/#helm

This is an example of `how to upgrade`.

```
helm upgrade CHART_NAME ./kube-review \
--install \
--reset-values \
--namespace NAMESPACE_NAME \
--values values.yaml \
--set envFrom.secretRef.name=SECRET_NAME \
--set image.repository=IMAGE_REPOSITORY_URL/IMAGE_REPOSITORY_NAME \
--set image.tag=IMAGE_TAG \
--set imagePullSecrets=IMAGE_PULL_SECRETS_NAME \
--set "ingress.hosts[0].host=HOST_NAME.shared-prod.fih.io" \
--set "ingress.hosts[0].paths[0]=/" \
--set "ingress.tls[0].hosts[0]=HOST_NAME.shared-prod.fih.io" \
--wait
```

The same for the `kube-review-prune` chart:

```
helm upgrade CHART_NAME ./kube-review-prune \
--install \
--reset-values \
--namespace NAMESPACE_NAME \
--wait  \
--set image.tag=IMAGE_TAG \
--set github.ghToken=GITHUB_TOKEN \
--set github.ghUserName=GITHUB_USERNAME
```

## How to Uninstall

REF.: https://helm.sh/docs/helm/helm_uninstall/

This is an example of `how to uninstall`.

`helm uninstall CHART_NAME --namespace NAMESPACE_NAME`
