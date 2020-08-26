Template Helm Chart
===================

This is a `template helm chart` created to be used in the CI pipeline for all projects.

Ref.: https://helm.sh/docs/helm/helm_create/

How to Install
--------------

This is an example of `how to install`.

```
helm install CHART_NAME cf-review-env \
--install \
--reset-values \
--repo cm://h.cfcr.io/findhotel/default/ \
--version 0.2.0 \
--namespace NAMESPACE_NAME \
--values values.yaml \
--set envFrom.secretRef.name=SECRET_NAME \
--set image.repository=265409992602.dkr.ecr.eu-west-1.amazonaws.com/daedalus-server-main \
--set image.tag=IMAGE_TAG \
--set imagePullSecrets=IMAGE_PULL_SECRETS_NAME \
--set "ingress.hosts[0].host=HOST_NAME" \
--set "ingress.hosts[0].paths[0]=/" \
--set "ingress.tls[0].hosts[0]=HOST_NAME" \
--set "ingress.tls[0].secretName=SECRET_TLS_NAME" \
--wait
```

How to Upgrade
--------------

This is an example of `how to upgrade`.

```
helm upgrade CHART_NAME cf-review-env \
--install \
--reset-values \
--repo cm://h.cfcr.io/findhotel/default/ \
--version 0.2.0 \
--namespace NAMESPACE_NAME \
--values values.yaml \
--set envFrom.secretRef.name=SECRET_NAME \
--set image.repository=265409992602.dkr.ecr.eu-west-1.amazonaws.com/daedalus-server-main \
--set image.tag=IMAGE_TAG \
--set imagePullSecrets=IMAGE_PULL_SECRETS_NAME \
--set "ingress.hosts[0].host=HOST_NAME" \
--set "ingress.hosts[0].paths[0]=/" \
--set "ingress.tls[0].hosts[0]=HOST_NAME" \
--set "ingress.tls[0].secretName=SECRET_TLS_NAME" \
--wait
```

How to Delete
-------------

This is an example of `how to delete`.

`helm delete CHART_NAME --namespace NAMESPACE_NAME`
