# syntax=docker/dockerfile:1

# Build
FROM golang:1.23-bookworm AS build-stage

WORKDIR /go/src/app

COPY go.mod go.sum ./
RUN go mod download

COPY *.go ./
COPY internal/ internal/

RUN CGO_ENABLED=0 go build -o /go/bin/app

# Tests
FROM build-stage AS run-test-stage
RUN go test -v ./...

## Deploy
FROM gcr.io/distroless/static-debian12 AS build-release-stage

COPY --from=build-stage /go/bin/app /

EXPOSE 8080

USER nonroot:nonroot

ENTRYPOINT ["/app"]