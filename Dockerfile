# Build the scraper binary
FROM golang:1.17 as builder
WORKDIR /workspace
# Copy the Go Modules manifests
COPY go.mod go.mod
COPY go.sum go.sum
# cache deps before building and copying source so that we don't need to re-download as much
# and so that source changes don't invalidate our downloaded layer
RUN go mod download
# Copy the go source
COPY main.go main.go
COPY datastore/ datastore/
COPY scraper/ scraper/
# Build
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -o scraper-api main.go

FROM alpine:latest as runner
WORKDIR /
COPY --from=builder /workspace/scraper-api .
ENTRYPOINT ["/scraper-api"]