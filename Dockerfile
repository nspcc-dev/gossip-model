FROM golang:1.14 as builder

COPY . /gossipmodel

WORKDIR /gossipmodel
RUN set -x \
    && export CGO_ENABLED=0 \
    && go build -o /go/bin/gossipmodel ./main.go

# Executable image
FROM alpine:3.11

COPY --from=builder /go/bin/gossipmodel /usr/local/sbin/gossipmodel