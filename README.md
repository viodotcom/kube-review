- [Kube-Review](#kube-review)
    + [Features](#features)
  * [Documentation](#documentation)
  * [Licensing](#licensing)

# Kube-Review

*Kube-Review* is a Kubernetes based platform to deploy and manage review environments.

*Kube-Review* is a battle tested platform based on Kubernetes to deploy and managed review environments at scale. *Kube-Review* has been in use at FindHotel for more than a year and it's considered resilient and production ready.

Based on Kubernetes, *Kube-Review* is very lean on resources as it is optimized to require few infra resources and to scale to high number of concurrent environments.

In order to control the life cycle of a environment, *Kube-Review* requires integration with your source code repository. At this point we only support **GitHub**, but there is no reason why others couldn't be added.

**WARNING**: Although, *Kube-Review* is resilient and scalable, it's not meant to run live or customer facing workloads.

### Features

These are some of the features supported by *Kube-Review*:

- Simple and universal deployment through any CI/CD tool;
- Public accessible URL with HTTPS support;
- Automatic removal of expired environments by time on when the branch is merged;
- Complex environments with side car containers support;
- Secrets and environment variables;
- Ephemerals and Non-Ephemerals environments;
- Custom environment names;
- Open to full customization through kustomize overlays;
- Extension through pre and post install hooks;
- Scalable and lean infra using Nginx Ingress and Let's Encrypt WildCards domains;
- Environment isolation through namespaces;
- Connection test after deployment;
- [Vertical Pod Autoscaling](https://cloud.google.com/kubernetes-engine/docs/concepts/verticalpodautoscaler);
- [Scaling From/To zero with Keda HTTP Add-On](https://github.com/kedacore/charts/tree/main/http-add-on)

## Documentation

**WARNING**: Docs are WIP and not complete yet.

- [Introduction](docs/introduction.md)
- [Tutorial](docs/tutorial.md)
- [Reference](docs/reference.md)
- [Customization](docs/customization.md)

## Licensing

Apache License 2.0
