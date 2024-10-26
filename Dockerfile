FROM golang:1.23-alpine AS builder
MAINTAINER Ugo Landini <ugo@confluent.io>

ARG VERSION=2.0.0
ARG GOVERSION=$(go version)
ARG USER=$(id -u -n)
ARG TIME=$(date)

RUN apk update \
    && apk add --no-cache git ca-certificates \
    && apk add --update gcc musl-dev libssl3 librdkafka-dev pkgconf \
    && update-ca-certificates

RUN adduser \
    --disabled-password \
    --gecos "" \
    --home "/home/jr" \
    --shell "/bin/sh" \
    --uid "100001" \
    "jr-user"

WORKDIR /go/src/github.com/jrnd-io/jr
COPY cmd cmd
COPY pkg pkg
COPY config config
COPY templates templates
COPY go.mod go.mod
COPY go.sum go.sum

RUN go mod tidy
RUN CGO_ENABLED=1 GOOS=linux go build \
     -tags musl -v \
     -ldflags="-X 'github.com/jrnd-io/jrv2/cmd.Version=${VERSION}' -X 'github.com/jrnd-io/jrv2/pkg/cmd.GoVersion=${GOVERSION}' -X 'github.com/jrnd-io/jrv2/cmd.BuildUser=${USER}' -X 'github.com/jrnd-io/jrv2/cmd.BuildTime=${TIME}' -linkmode external -w -s -extldflags '-static'" \
     -a -o build/jr github.com/jrnd-io/jrv2/cmd/jr

FROM scratch
COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /etc/group /etc/group
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /go/src/github.com/jrnd-io/jr/templates/ /home/jr/.jr/templates/
COPY --from=builder /go/src/github.com/jrnd-io/jr/config/ /home/jr/.jr/
COPY --from=builder /go/src/github.com/jrnd-io/jr/build/jr /bin/jr

USER jr-user:jr-user

ENV JR_SYSTEM_DIR=/home/jr/.jr
