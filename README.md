# kube-review

*kube-review* is a Kubernetes based platform to deploy and manage review environments.

*kube-review* is a battle tested platform based on Kubernetes to deploy and managed review environments at scale. *kube-review* has been in use at FindHotel for more than a year and it's considered resilient and production ready.

Based on Kubernetes, *kube-review* is very lean on resources as it is optimized to require few infra resources and to scale to high number of concurrent environments.

In order to control the life cycle of a environment, *kube-review* requires integration with your source code repository. At this point we only support **GitHub**, but there is no reason why others couldn't be added.

**WARNING**: Although, *kube-review* is resilient and scalable, it's not meant to run live or customer facing workloads.

### Features

These are some of the features supported by *kube-review*:

- Simple and universal deployment through any CI/CD tool;
- Public accessible URL with HTTPS support;
- Automatic removal of expired environments by time on when the branch is merged;
- Complex environments with side car containers support;
- Secrets and environment variables;
- Ephemerals and Non-Ephemerals environments;
- Custom environment names;
- Kubernetes resource parametrization through Helm values file.
- Customization through pre and post install hooks;
- Scalable and lean infra using Nginx Ingress and Let's Encrypt WildCards domains;
- Environment isolation through namespaces;
- Connection test after deployment.

## Documentation

**WARNING**: Docs are WIP and not complete yet.

- [Introduction](docs/introductions.md)
- [Getting Started](docs/getting-started.md)
- [Reference](docs/reference.md)
- [Developing](docs/developing.md)

## Licensing

Apache License 2.0