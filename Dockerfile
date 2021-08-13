FROM golang:1.16 AS build
WORKDIR /app
COPY go.mod ./
COPY go.sum ./
COPY *.go ./
COPY sls_store/*.go ./sls_store/
RUN go mod download && CGO_ENABLED=0 GOARCH=amd64 GOOS=linux go build -o /jaeger-sls-plugin

FROM jaegertracing/all-in-one:1.25.0
ENV ACCESS_KEY_SECRET="" \
    ACCESS_KEY_ID="" \
    PROJECT="" \
    ENDPOINT="" \
    INSTANCE="" \
    GRPC_STORAGE_PLUGIN_BINARY="/jaeger-sls-plugin" \
    SPAN_STORAGE_TYPE=grpc-plugin \
    JAEGER_DISABLED=true  \
    GRPC_STORAGE_PLUGIN_LOG_LEVEL=DEBUG
COPY --from=build /jaeger-sls-plugin /jaeger-sls-plugin

