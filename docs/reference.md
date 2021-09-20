# Reference

The `kube-review` is composed of two main components, `deploy` and `prune`.

## Deploy Component

The deploy component is a `bash` script that can be used to deploy a review environment. At it's core it's a simple script that will do all the necessary magic to get a review env running. Therefore, the `deploy` component is a command that is called every-time that an environment needs to created or updated.

The basic flow of the components is, configure the correct kube context, create the resources, and test the installation. Depending on the use case, some other features might be used, like running pre/post scripts.

### Options

The `deploy` components contains many options which can be passed as environment variables. This is the list of all options:

| Name | Description | Default Value | Required |
| - | - | - | - |
| KR_ID | A unique identifier for the review environment. It's recommended this to be the branch name. | - | true |
| KR_IMAGE | The full url of the container image, including the tag that the app should run. e.g: `KR_IMAGE=123456789012.dkr.ecr.eu-west-1.amazonaws.com/ecr-repo-name@sha256:2de5d881318951e48c69987df478b076f2b2ac9971d12f1b342347ed51ff7f4b` | - | true |
| KR_DOMAIN | The domain on which the app should be available. e.g: `foo.com` | - | true |
| KR_CONTAINER_PORT | The port exposed by the main container. | - | true |
| KR_KUBE_CONTEXT | The kube context from the kube config file that should be used. | Default to current context. | false |
| KR_PREFIX | A prefix to be added to the name of the environment. | re | false |
| KR_IS_EPHEMERAL | If the environment is ephemeral or not. Non ephemeral environments will never be expired. | true | false |
| KR_KUBE_CONFIG_FILE | The kube config file used for connecting to Kuberneres. The file has to be accessible on the local file system during execution of the script. | $HOME/.kube/config | false |
| KR_PULL_REQUEST_NUMBER | The pull request number that is getting deployed. This will be saved as annotation into the namespace so that the prune command can check the expiration of the branch/pr with the source code service. | - | false |
| KR_BRANCH_NAME | The branch that is getting deployed. This will be saved as annotation into the namespace so that the prune command can check the expiration of the branch/pr with the source code service. | - | false |
| KR_REPO_NAME | The repo name source code in question. This will be saved as annotation into the namespace so that the prune command can check the expiration of the branch/pr with the source code service. | - | false |
| KR_REPO_OWNER | The repository owner of the source code in question. This will be saved as annotation into the namespace so that the prune command can check the expiration of the branch/pr with the source code service.  | - | false |
| KR_TEST_CONNECTION | Enable/disable testing the url of the environment once the deployment is done. If the connection fails the deployment will also fails. | true | false |
| KR_TEST_CONNECTION_URL_PATH | A custom url path to be appended to the final url when running a connection test. | / | false |
| KR_PRE_HOOK | A shell command to be executed before the deployment starts. | - | false |
| KR_POST_HOOK | A shell command to be executed after the deployment is finished. | - | false |
| KR_BASE_OVERLAY_PATH | The path containing the base kustomize overlay to be used. | src/deploy/resources/base | false |
| KR_OVERLAY_PATH | The path containing a kustomize overlay to be used. | - | false |
| KR_VERBOSE | Prints verbose or debug messages. Should not be used in production. | false | false |

## Prune Component

The `prune` component is a simple go application that is responsible to cleanup expired environments. The component will scan all namespaces deployed by the `deploy` component, if the environment is expired it will delete the namespace. Thus, the `prune` component is deployed as a `Cron Job` and only needs to be installed once and doesn't need to be executed manually.

### Options

Although the component is not executed manually but through a `Cron Job` it's still worth to have its options documented here:

| Name | Description | Default Value | Required |
| - | - | - | - |
| name | environment name to filter, accepts glob expressions | *  | false |
| expiration | how many hous to consider an environment stale | 120 | false |
| dryRun | only show logs but don'r perform deletes | false | false |
| k8sKubeconfig | absolute path to the kubeconfig file, if not informed will use in cluster config | - | false |
| k8sContextName | the k8s context name to operate on | - | false |
| ghEndpoint | github api endpoint | https://api.github.com | false |
| ghToken | github api token | - | true |
| ghUserName | github username to use for auth | - | true |
