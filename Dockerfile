# syntax=docker/dockerfile:1
FROM golang:1.25.7-bookworm AS builder

RUN apt-get update -qq \
    && apt-get install -y --no-install-recommends \
         git ca-certificates \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /go/src

COPY go.mod go.sum ./
RUN go mod download

COPY cmd ./cmd
COPY internal ./internal

ENV CGO_ENABLED=0
ENV GOOS=linux
RUN go build -o /go/bin/app ./cmd

RUN update-ca-certificates

FROM scratch AS app

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /go/bin/app /app

ENV HTTP_SERVER_LISTEN_ADDR=0.0.0.0:8080
EXPOSE 8080

ENTRYPOINT ["/app"]
