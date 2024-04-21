# ./Dockerfile

FROM golang:1.21-alpine AS builder

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

FROM scratch

# Copy binary and config files from /build 
# to root folder of scratch container.
COPY --from=builder /build/wannabe /wannabe

# Export necessary port.
EXPOSE 1234

# Command to run when starting the container.
ENTRYPOINT ["/wannabe"]