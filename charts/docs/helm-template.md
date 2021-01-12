## How to debug the templates files

Ref.: https://helm.sh/docs/chart_template_guide/debugging/

When will be necessary to add new changes on the templates file it's possible to do a debug and then to check if the new change is correct or not, basically you need to run this command:

`helm template test cf-review-env -f cf-review-env/values.yaml --debug \
--set envFrom.secretRef.name=SECRET_NAME \
--set image.repository=IMAGE_REPOSITORY_URL/IMAGE_REPOSITORY_NAME \
--set image.tag=IMAGE_TAG \
--set imagePullSecrets=IMAGE_PULL_SECRETS_NAME \
--set "ingress.hosts[0].host=HOST_NAME.shared-prod.fih.io" \
--set "ingress.hosts[0].paths[0]=/" \
--set "ingress.tls[0].hosts[0]=HOST_NAME.shared-prod.fih.io"`

**Note:** It's necessary to use a true `values.yaml` file.
