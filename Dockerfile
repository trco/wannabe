# ./Dockerfile

# FROM golang:1.21-alpine AS builder
FROM golang:1.21-alpine AS builder

# RUN apk --no-cache add ca-certificates bash

# Move to working directory (/build).
WORKDIR /build

# Copy and download dependency using go mod.
COPY go.mod go.sum ./
RUN go mod download

# Copy the code into the container.
COPY . .

# Set necessary environment variables needed for our image 
# and build the API server.
RUN go build -ldflags="-s -w" -o wannabe .

FROM alpine

# Copy binary and config files from /build 
# to root folder of scratch container.
COPY --from=builder /build/wannabe /wannabe

RUN apk --no-cache add ca-certificates bash curl

# COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# Export necessary port.
EXPOSE 1234

# Command to run when starting the container.
ENTRYPOINT ["/wannabe"]