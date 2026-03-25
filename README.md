# Waypoint

Waypoint example application.

## ArgoCD application examples

Database needs to be deployed and synced before the Waypoint application.

### Waypoint database application

```yaml
apiVersion: argoproj.io/v1alpha1
kind: Application
metadata:
  name: waypoint-db
spec:
  destination:
    namespace: default
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
  name: waypoint
spec:
  destination:
    namespace: default
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
          value: waypoint-db
  sources: []
  project: default
```