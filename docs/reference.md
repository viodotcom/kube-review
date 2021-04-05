# Reference

## App Component

The recommended way of installing/deploying an app is through the deploy shell script.

### Shell Options

The script contains many options which can be passed as environment variables. This is the list of all options:

| Name | Description | Default Value | Required |
| - | - | - | - |
| KR_ID | A unique identifier for the review environment. It's recommended this to be the branch name. | - | true |
| KR_IMAGE_URL | The url of the container image that the app should run. | - | true |
| KR_IMAGE_TAG | The tag of the container image that the app should run. | - | true |
| KR_KUBE_CONTEXT | The kube context from the kube config file that should be used. | - | true |
| KR_DOMAIN | The domain on which the app should be available. e.g: `foo.com` | - | true |
| KR_PREFIX | A prefix to be added to the name of the environment. | re | false |
| KR_IS_EPHEMERAL | If the environment is ephemeral or not. Non ephemeral environments will never be expired. | true | false |
| KR_CHART_VERSION | The version of the `kube-review` app chart to be used. | latest | false |
| KR_KUBE_CONFIG_FILE | The kube config file used for connecting to Kuberneres. The file has to be accessible on the local file system during execution of the script. | - | false |
| KR_POST_INSTALL_MSG | A message to be printed once the deployment is done. | - | false |
| KR_SECRETS_FILE | The secrets file from which secrets will be loaded and inject as environment variable secrets on Kubernetes. The file has to be accessible on the local file system during execution of the script. | - | false |
| KR_PULL_REQUEST_NUMBER | The pull request number that is getting deployed. This will be saved as annotation into the namespace so that the purge command can check the expiration of the branch/pr with the source code service. | - | false |
| KR_BRANCH_NAME | The branch that is getting deployed. This will be saved as annotation into the namespace so that the purge command can check the expiration of the branch/pr with the source code service. | - | false |
| KR_REPO_NAME | The repo name source code in question. This will be saved as annotation into the namespace so that the purge command can check the expiration of the branch/pr with the source code service. | - | false |
| KR_REPO_OWNER | The repository owner of the source code in question. This will be saved as annotation into the namespace so that the purge command can check the expiration of the branch/pr with the source code service.  | - | false |
| KR_TEST_CONNECTION | Enable/disable testing the url of the environment once the deployment is done. If the connection fails the deployment will also fails. | true | false |
| KR_PRE_HOOK | A shell command to be executed before the deployment starts. | - | false |
| KR_POST_HOOK | A shell command to be executed after the deployment is finished. | - | false |
| KR_HELM_REPO_URL | The helm repo url from which to get the app chart from. | cm://h.cfcr.io/findhotel/default/ | false |
| KR_VALUES_FILE | The helm values file to be used to customize the Kubernetes resources. | - | false |

### Values File Options

To allow customization of the resources deployed to Kubernetes, it's possible to pass a yaml values file, which will be override the default values from the chart.

#### Options

TODO

## Purge Component

TODO
