# syntax=docker/dockerfile:1

# Build the application from source
FROM golang:1.21 AS build-stage

WORKDIR /queue

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN GOOS=linux go build cmd/main.go

# Run the tests in the container
FROM build-stage AS run-test-stage
RUN go test -v ./...

# Deploy the application binary into a lean image

WORKDIR /

ENTRYPOINT ["/queue/main"]