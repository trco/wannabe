# ./Dockerfile

FROM golang:1.22-alpine AS builder

WORKDIR /build

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -ldflags="-s -w" -o wannabe .

FROM alpine

WORKDIR /usr/src/app

COPY --from=builder /build/wannabe /usr/src/app/wannabe
COPY certs/wannabe.crt /usr/src/app/certs/wannabe.crt
COPY certs/wannabe.key /usr/src/app/certs/wannabe.key

RUN apk add --no-cache bash curl

ENV CERT_PATH=/usr/src/app/certs/wannabe.crt
ENV CERT_KEY_PATH=/usr/src/app/certs/wannabe.key

EXPOSE 6789
EXPOSE 6790

ENTRYPOINT ["/usr/src/app/wannabe"]