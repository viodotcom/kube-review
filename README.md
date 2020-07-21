cf-review-env
=============

Review environment support for CodeFresh using Kubernetes.

## Requirements

You will need `brew` installed and a CodeFresh API Token.

## Setup

To instal requied depencies run:

    CFTOKEN=[YOUR-TOKEN-HERE] make setup

## Create and Update

In order to create or update the pipelines you need a project already crated in CodeFresh.
We recommend using the `cf-review-env-[STAGE]` as name of the project.

To create the pipelines run:

    STAGE=stg make create

After that, you can update the pipelines by running:

    STAGE=stg make update
