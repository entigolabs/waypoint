# Waypoint

Waypoint example application.

## ArgoCD application examples

Database needs to be deployed and synced before the Waypoint application.

### Waypoint database application

```yaml
apiVersion: argoproj.io/v1alpha1
kind: Application
metadata:
  name: waypoint-helm-db
spec:
  destination:
    namespace: helm-test
    server: https://kubernetes.default.svc
  source:
    path: ''
    repoURL: ghcr.io/entigolabs
    targetRevision: '*.*.*'
    chart: waypoint-db-helm
    helm:
      parameters:
        - name: deletionProtection
          value: 'true'
  sources: []
  project: default
```

### Waypoint application

```yaml
apiVersion: argoproj.io/v1alpha1
kind: Application
metadata:
  name: waypoint-helm
spec:
  destination:
    namespace: helm-test
    server: https://kubernetes.default.svc
  source:
    path: ''
    repoURL: ghcr.io/entigolabs
    targetRevision: '*.*.*'
    chart: waypoint-helm
    helm:
      parameters:
        - name: backend.config.ALLOWED_ORIGINS
          value: https://waypoint.example.com
        - name: ingress.host
          value: waypoint.example.com
        - name: database.releaseName
          value: waypoint-helm-db
  sources: []
  project: default
```