# Developing

## CI/CD

We use GitHub Actions to power the `kube-review` CI/CD flow. The pipelines are quite simple. Every commit on a `non master` branch will trigger a CI pipeline, which will build and publish the container and image and also build and publish the helm charts.

## Releases

### Container image

We use a very simple versioning model. The Docker image is always tagged with `latest` for the stable release (CD pipeline). During development the Docker image is tagged with the `commit hash` and `branch name` (CI pipeline).

Therefore, if you want to test changes affecting the docker image, you must use the image with the tag corresponding to the commit you want to target.

The docker image is hosted on [Docker Hub](https://hub.docker.com/r/findhotelamsterdam/kube-review).

### Helm Chart

The Helm Chart version uses `SemVer` for stable releases, but during development the version is always `0.0.1:{COMMIT_HASH}`.

When the new chart version is created by the CI pipeline, it will push the new chart to the lab repo, using dev version schema, for example `0.0.1+87d6164`. To use this version you must pass it as the `KR_CHART_VERSION` variable.

The Helm Chart is stored for stable releases is stored on the [main repo](cm://h.cfcr.io/findhotel/default/) and development releases at the [lab repo](cm://h.cfcr.io/findhotel/lab/).

Because the chart for development is stored in a different repo, you also need to override the helm repo when using the `deploy` script. You can do that by either using a development container image, which is already pointing to the `lab` repo or just override the `KR_HELM_REPO_URL`.

But if you are installing the chart directly with the Helm cli, you just specify the options in the cli itself. This is usually useful when testing the chart, but specially when working with the `kube-review-prune` chart.

**Note**: You should always use a `staging` environment when it is necessary to test a new `helm chart version`.