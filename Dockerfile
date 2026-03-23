FROM --platform=$BUILDPLATFORM golang:1.26-alpine-slim AS build
WORKDIR /src

COPY go.* ./
RUN go mod download
COPY . .

ARG TARGETARCH
ARG TARGETOS
ARG GITHUB_SHA=main
ARG VERSION=latest

RUN --mount=type=cache,target=/root/.cache/go-build \
    --mount=type=cache,target=/go/pkg \
    GOOS="$TARGETOS" GOARCH="$TARGETARCH" CGO_ENABLED=0 go build -ldflags \
    "-X github.com/entigolabs/waypoint/internal/version.version=${VERSION} \
     -X github.com/entigolabs/waypoint/internal/version.buildDate=$(date -u +'%Y-%m-%dT%H:%M:%SZ') \
     -X github.com/entigolabs/waypoint/internal/version.gitCommit=${GITHUB_SHA} \
     -extldflags -static -s -w" -o /bin/server . \

FROM gcr.io/distroless/static-debian13:nonroot
COPY --from=build /bin/server /bin/
EXPOSE 8081
ENTRYPOINT [ "/bin/server" ]
