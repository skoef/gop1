# Dockerfile References: https://docs.docker.com/engine/reference/builder/

# Start from the latest golang base image
FROM golang:latest AS builder

# Add Maintainer Info
LABEL maintainer="Reinier Schoof <reinier@skoef.nl>"

# Copy go mod file
COPY go.mod /

# Download all dependencies. Dependencies will be cached if the go.mod file is not changed
RUN go mod download

# Copy the source from the current directory to the Working Directory inside the container
COPY . /

# Build the Go app
WORKDIR /
RUN CGO_ENABLED=0 go build -o /gop1

FROM alpine:latest

WORKDIR /
COPY --from=builder /gop1 /gop1

CMD ["/gop1"]
