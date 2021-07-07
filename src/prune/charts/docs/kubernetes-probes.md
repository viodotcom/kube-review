# How to work Kubernetes Probes (livenessProbe, readinessProbe and startupProbe)

REF.: https://kubernetes.io/docs/tasks/configure-pod-container/configure-liveness-readiness-startup-probes/

- `initialDelaySeconds`: number of seconds to wait before initiating liveness or readiness probes
- `periodSeconds`: how often to check the probe
- `timeoutSeconds`: number of seconds before marking the probe as timing out (failing the health check)
- `successThreshold`: minimum number of consecutive successful checks for the probe to pass
- `failureThreshold`: number of retries before marking the probe as failed. For liveness probes, this will lead to the pod restarting. For readiness probes, this will mark the pod as unready.
