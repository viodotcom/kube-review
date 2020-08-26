cf-review-env
=============

Review environment support for CodeFresh using Kubernetes.

## Releases

We use a very simple versioning model. The Docker image is always tagged with `latest` for the
stable release (CD pipeline). During development the Docker image is tagged with `short commit hash`
and `branch name` (CI pipeline).

The Helm Chart version uses SemVer for stable releases, but during development the version
is always `0.0.1:{SHORT_COMMIT_HASH}`.