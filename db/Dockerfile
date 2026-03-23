FROM --platform=$BUILDPLATFORM golang:1.26-alpine AS build

ARG TARGETOS
ARG TARGETARCH

WORKDIR /temp-build

RUN go mod init temp-build
RUN go get github.com/jackc/tern/v2@v2.3.5

RUN --mount=type=cache,target=/root/.cache/go-build \
    --mount=type=cache,target=/go/pkg/mod \
    CGO_ENABLED=0 \
    GOOS=${TARGETOS} \
    GOARCH=${TARGETARCH} \
    go build \
    -ldflags="-s -w -extldflags '-static'" \
    -o /out/tern \
    github.com/jackc/tern/v2

FROM alpine:3 AS final

ARG UID=10001
RUN adduser \
    --disabled-password \
    --gecos "" \
    --home "/nonexistent" \
    --shell "/sbin/nologin" \
    --no-create-home \
    --uid "${UID}" \
    appuser
USER appuser

COPY --from=build /out/tern /bin/tern

COPY db.sh /db.sh
COPY core/ /core/
COPY migrations/ /migrations/

ENTRYPOINT ["/db.sh"]