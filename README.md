# cf-review-env

Review environment support for [CodeFresh](http://codefresh.io/) using [Kubernetes](https://kubernetes.io/).

## Releases

We use a very simple versioning model. The Docker image is always tagged with `latest` for the stable release (CD pipeline). During development the Docker image is tagged with `short commit hash` and `branch name` (CI pipeline).

The Helm Chart version uses SemVer for stable releases, but during development the version is always `0.0.1:{SHORT_COMMIT_HASH}`.

## How to use the development version

**Note**: We should always use the `staging` environment when it is necessary to test a new `helm chart version`.

When the new `chart version` is created the CI pipeline run and push the new chart to the registry with the new version for example `0.0.1+87d6164` and we can use this new version following these steps:

1 - Choose the project development pipeline that you would like to test (e.g.: [geolocation-service-lab](https://g.codefresh.io/projects/geolocation-service-lab/edit/pipelines/?projectId=5fbf87e2b4b6c926b5fe6ebc))

2 - Run a new deployment but first you need to add these variables in `BUILD VARIABLES`:

| Variable  | Value |
|----- |-------|
| KUBE_CONTEXT | k8s-context |
| APP_DOMAIN | k8s-domain |
| CHART_VERSION | 0.0.1+87d6164 |
| CF_REVIEW_ENV_IMAGE_TAG | your-branch-name |

**Note**: The value of the `CHART_VERSION` variable, you must add the new version of the chart created for testing.
