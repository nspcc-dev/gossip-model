FROM golang:1.11-alpine3.8 as builder

RUN apk add --no-cache git

COPY . /gossipmodel

WORKDIR /gossipmodel

# https://github.com/golang/go/wiki/Modules#how-do-i-use-vendoring-with-modules-is-vendoring-going-away
# go build -mod=vendor
RUN set -x \
    && export CGO_ENABLED=0 \
    && go build -mod=vendor -o /go/bin/gossipmodel ./main.go

# Executable image
FROM alpine:3.8

COPY --from=builder /go/bin/gossipmodel /usr/local/sbin/gossipmodel