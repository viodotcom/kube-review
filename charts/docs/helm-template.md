## How to debug the templates files

Ref.: https://helm.sh/docs/chart_template_guide/debugging/

When will be necessary to add new changes on the templates file it's possible to do a debug and then to check if the new change is correct or not, basically you need to run this command:

`helm template test cf-review-env -f values.yaml --debug`

**Note:** It's necessary to use a true `values.yaml` file.
