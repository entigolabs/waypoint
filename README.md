# Waypoint

Waypoint is a demo application that showcases automated application provisioning into a Kubernetes cluster using [Entigo Infralib](https://github.com/entigolabs/entigo-infralib). It demonstrates data gathering from an external public API, storing the results in a database, and serving them through a web interface. Capable of running as multiple replicas for redundancy.

The application is API-first: both components generate their client and server code from a shared [OpenAPI specification](./openapi.yaml).

## Components

Backend is a Go application that features:
* Database sql schema and migrations management.
* Collector that collects data from [Andmete Teabevärav](https://andmed.eesti.ee/) and stores it in a PostgreSQL database. Collected data includes categories, EMS categories and themes. Supports multiple instances by using table locks to ensure that only one instance is collecting data at a time.
* REST API server that serves the collected data from the database from `/api` endpoints.
* Prometheus metrics endpoint for observability.
* Technical documentation in [README.md](./backend/README.md).

Frontend is a React application that features:
* Client for retrieving the collected data from the backend API.
* Ant Design based components for displaying the data.
* Accessibility features.
* Technical documentation in [README.md](./frontend/README.md).

## Local development

To run the full application locally:
1. Start the database: `docker compose up -d` in the `backend/` directory.
2. Run database migrations, described in [backend/README.md](./backend/README.md#database).
3. Start the [backend](./backend/README.md#running-locally).
4. Start the [frontend](./frontend/README.md#development).

## CI/CD

Pull requests and pushes to `main` run linting and tests as a quality gate for both backend and frontend, as well as Helm chart linting.

The release pipeline is triggered by a version tag (`vX.Y.Z`) and:
1. Runs all quality gates.
2. Builds and publishes Docker images for the backend, database migrations, and frontend to GHCR with build provenance attestations.
3. Packages and publishes the `waypoint-helm` Helm chart to GHCR.

The database Helm chart (`waypoint-db-helm`) is versioned and published independently on every merge to `main` that modifies `backend/db` files.

## Attestations

To view the build provenance attestations for the published images, you can use the `docker buildx imagetools`. For example:

```bash
docker buildx imagetools inspect ghcr.io/entigolabs/waypoint:latest --format '{{json .SBOM}}'
docker buildx imagetools inspect ghcr.io/entigolabs/waypoint-front:latest --format '{{json .SBOM}}'
docker buildx imagetools inspect ghcr.io/entigolabs/waypoint-db:latest --format '{{json .SBOM}}'
```

## Deployment

Prerequisites: a Kubernetes cluster with ArgoCD.

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