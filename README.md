# KReview

KReview is a [Kubernetes](https://kubernetes.io/) native tool that allows easy deployment and management of review environments.

The most common use case is to easily deploy a container app, for each push made to a branch, from the a CI/CD pipeline. At the same time that it ensures that the environment is removed once it's no longer needed.

KReview is a helm chart and a controller job. The helm chart will be installed for each review environment and the controller will make sure that expired environments are removed.

## Features

The following is non exhaustive list of the main features:

- Support for ephemeral and non ephemeral deployments;
- Easy integration with any CI/CD tool through Docker;
- Self contained deployments;
- SSL with let's encrypt;
- Support for wild card domains;
- Automatically purge of expired environments;
- Flexible configuration through Helm values file;
- Support for pre/post deployment hooks.

## Start using

To start using, check our [docs](docs/index.md).
