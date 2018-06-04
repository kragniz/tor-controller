FROM golang:1.9.3 as builder

WORKDIR /go/src/github.com/kragniz/tor-controller
COPY pkg/    pkg/
COPY cmd/    cmd/
COPY vendor/ vendor/

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -o tor-local-manager ./cmd/tor-local-manager/main.go

FROM alpine:edge
RUN apk update \
  && apk add tor --update-cache --repository http://dl-3.alpinelinux.org/alpine/edge/community/ \
  && rm -rf /var/cache/apk/* \
  && mkdir -p /run/tor/service
WORKDIR /root/
COPY --from=builder /go/src/github.com/kragniz/tor-controller/tor-local-manager .
ENTRYPOINT ["./tor-local-manager"]
