# Build the manager binary
FROM golang:1.17 as builder

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
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 GO111MODULE=on go build -o whosdriving-be

# Use distroless as minimal base image to package the manager binary
# Refer to https://github.com/GoogleContainerTools/distroless for more details
FROM gcr.io/distroless/static:nonroot
# FROM gcr.io/distroless/static-debian11
WORKDIR /app

# copy the assets
COPY assets/ assets/

COPY --from=builder /app .
# Do we create minimal stucture here ?
# RUN mkdir -p /app/data

USER nonroot:nonroot

EXPOSE 9000
ENTRYPOINT ["/app/whosdriving-be"]