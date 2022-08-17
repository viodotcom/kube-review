- [Customization](#customization)
  * [Usage](#usage)

# Customization

The deploy component was designed to accept the most common options that one needs when deploying a simple review environment, like container port and the image to be deployed. However, it's very important to allow more advanced use cases through customization. Users need to be able to customize any part of a review env if necessary.

To implement that, the deploy command is based on **Kustomize**. **Kustomize** allow us to have a common base that can be customized through the creation of overlays on top of that base. By using customize we allow our users to fully customize any part of a review environment.

## Usage

To customize any part of a review environment, one needs to create a folder with at least a [kustomization.yml](https://kubectl.docs.kubernetes.io/references/kustomize/kustomization/), while adding the base overlay in the [resources](https://kubectl.docs.kubernetes.io/references/kustomize/resource/) list. After that one can patch by either adding [patches](https://kubectl.docs.kubernetes.io/references/kustomize/patches/) or by adding new resources under `resources`: This is an example from the *Kube-Review* repo:

```yaml
apiVersion: kustomize.config.k8s.io/v1beta1
kind: Kustomization

resources:
  - ../base

secretGenerator:
- name: secret
  envs:
  - secrets.txt
  type: Opaque
  options:
    disableNameSuffixHash: true
  behavior: merge

patches:
  - path: deployment.yml
  - path: patches/deployment.patch.json
    target:
      kind: Deployment
      name: deployment
```

Note that as secret ingestion is also done through **kustomize**, you also have to add a [secretGenerator](https://kubectl.docs.kubernetes.io/references/kustomize/secretgenerator/) and a `secrets.txt` file, even if a secret is not necessary.

Patches can change any resource created from base. That can be done using either a simple [patch](https://kubectl.docs.kubernetes.io/references/kustomize/patches/) file or a [json](https://kubectl.docs.kubernetes.io/references/kustomize/patchesjson6902/) file. This is how the `deployment.yml` looks like:

```yaml
metadata:
  name: deployment
spec:
  template:
    spec:
      containers:
        - name: kube-review
          startupProbe:
            periodSeconds: 10
        - name: redis
          image: "redis:alpine"
          ports:
            - name: redis
              containerPort: 6379
              protocol: TCP
```

If one needs to inject dynamic variables in the resources, that can be done by using the json patch file. The deploy component will replace any env variables present in the file and in the environment. This is how the `deployment.patch.json` file looks like:

```json
[
    {
        "op": "replace",
        "path": "/spec/template/spec/containers/1/image",
        "value": "redis:$LABEL"
    }
]
```

Finally, in order to use the overlay one has to specify the `KR_OVERLAY_PATH` variable and the env vars to be injected in the resources:

```
KR_ID=nginx \
KR_IMAGE=nginx:latest \
KR_DOMAIN="my-domain.io" \
KR_CONTAINER_PORT="80" \
KR_OVERLAY_PATH=src/deploy/resources/example \
LABEL=6.2.1 \
src/deploy/deploy
```
