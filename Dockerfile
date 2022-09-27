# Build the manager binary
FROM golang:1.18-alpine as builder
RUN apk add --no-cache gcc g++

WORKDIR /app
# Copy the Go Modules manifests
COPY go.mod go.mod
COPY go.sum go.sum
# cache deps before building and copying source so that we don't need to re-download as much
# and so that source changes don't invalidate our downloaded layer
RUN go mod download

# Copy the go source
COPY assets/ assets/
COPY data_interface/ data_interface/
COPY graph/ graph/
COPY tools.go tools.go
COPY server.go server.go
COPY server_test.go server_test.go

# Build
RUN GOOS=linux GOARCH=amd64 GO111MODULE=on CGO_ENABLED=1 go build -ldflags="-w -s" -o whosdriving-be


FROM alpine:latest
RUN apk add --no-cache musl-dev

WORKDIR /app

RUN adduser -D -g 'nonroot' nonroot
USER nonroot:nonroot

# copy the assets
COPY --from=builder --chown=nonroot:nonroot /app/assets/ assets/
COPY --from=builder --chown=nonroot:nonroot /app .

EXPOSE 8080
ENTRYPOINT ["/app/whosdriving-be"]