# Introduction

## Motivation

In modern agile development workflows, getting quick feedback on one's work is quite important. Review environments, also known as preview environments, are a way to easily deploy and quick updated an usually lean version of an app.

Every-time a developer pushes to a branch, a review environment is created/updated with latest changes. This allows developers to quickly get feedback or their work.

Once the branch is merged or developer has not made any changes in some time, the environment is considered expired and is purged. However, in some cases it might be useful to keep environments around for longer.

Companies usually implement such environments in a bespoke manner for each different system. This creates a lot of duplication and complexity.

Having only one way of deploying review environments across a organization helps reduce the operations burden on individual teams and enables a streamlined process.

However, exposing direct Kubernetes to developer, for such environments, add too much complexity.

## Our Take

KReview tries to make deployment of review environments easy, by providing a opinionated, yet customizable, package that abstracts kubernetes details away from developers.

Developers only need to specify details about the app that they want to deploy, things like container image, port number, secrets and so on.

Kubernetes resources are managed by the App helm chart, which is templated to expose through values.yml file possible customizations.

To make things even easier for CI/CD environments, a docker image takes care of doing all the work of installing the chart, making sure it's running and so on.

To make sure that expired environments are purge, the Manager helm char takes care of installing a cron job that will periodically cleanup expired environments.
