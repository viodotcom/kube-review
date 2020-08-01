cf-review-env
=============

Review environment support for CodeFresh using Kubernetes.

## Requirements

You will need `brew` installed and a CodeFresh API Token.

## Setup

To instal requied depencies run:

    CFTOKEN=[YOUR-TOKEN-HERE] make setup

## Create and Update

In order to create or update the pipelines you need a project already created in CodeFresh.
We recommend using `cf-review-env-[STAGE]` as name of the project.

To create the pipelines run:

    STAGE=stg make create

After that, you can update the pipelines by running:

    STAGE=stg make update

## Releases

We use a very simple versioning model. The Docker image is always tagged with `latest` for the
stable release (CD pipeline). During development the Docker image is tagged with `short commit hash`
and `branch name` (CI pipeline).

The Helm Chart version uses SemVer for stable releases, but during development the version
is always `0.0.1:{SHORT_COMMIT_HASH}`.