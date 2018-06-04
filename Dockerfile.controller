FROM golang:1.9.3 as builder

# Copy in the go src
WORKDIR /go/src/github.com/kragniz/tor-controller
COPY pkg/    pkg/
COPY cmd/    cmd/
COPY vendor/ vendor/

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -o controller-manager ./cmd/controller-manager/main.go

# Copy the controller-manager into a thin image
FROM ubuntu:latest
# RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /go/src/github.com/kragniz/tor-controller/controller-manager .
ENTRYPOINT ["./controller-manager"]
CMD ["--install-crds=false"]
