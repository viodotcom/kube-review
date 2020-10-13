# TODO

## Use startupProbe on the deployments

We need to use the `startupProbe` in our deployments, but this new feature will be available only on `EKS version 1.18`, so now we have it added in our [deployment.yaml](https://github.com/FindHotel/cf-review-env/blob/master/charts/cf-review-env/templates/deployment.yaml#L62) file but we can't use it.
**Note:** According to AWS the `EKS 1.18 version` will be available on `October, 2020`.


Issue: https://github.com/aws/containers-roadmap/issues/947
AWS Roadmap: https://docs.aws.amazon.com/eks/latest/userguide/kubernetes-versions.html#kubernetes-release-calendar
