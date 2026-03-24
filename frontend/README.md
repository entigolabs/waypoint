# Waypoint Frontend

A React dashboard for browsing Waypoint API data, built with Vite and Ant Design.

* [Requirements](#requirements)
* [Setup](#setup)
* [Environment variables](#environment-variables)
* [Development](#development)
* [Build](#build)
* [Testing](#testing)
* [Code generation](#code-generation)

## Requirements

- Node.js
- npm

## Setup

```bash
cp .env.example .env
# Edit .env as needed
npm install
```

## Environment variables

| Variable | Description | Default |
|---|---|---|
| `VITE_API_ENDPOINT` | Base URL of the Waypoint API | Same origin as the frontend |
| `VITE_APP_NAME` | Display name shown in the header | `Entigo Portal minimal` |

Leave `VITE_API_ENDPOINT` empty to have the frontend call the API at the same origin (useful when served behind a reverse proxy).

## Development
For development, with hot reload:

```bash
npm run dev
```

## Build
To build for production preview:

```bash
npm run build
```

The output is written to `dist/`.

To serve the build locally for preview:

```bash
npm run serve
```

## Testing
To test the frontend for accessibility and other issues:

```bash
npm test
```

## Code generation

The API client in `src/client/` is generated from the OpenAPI spec at `../openapi/openapi.yaml`. To regenerate after the spec changes:

```bash
npm run generate
```
