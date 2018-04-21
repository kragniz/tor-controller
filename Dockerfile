# build stage
FROM golang:alpine AS build-env
RUN apk update && apk add make
ADD . /go/src/github.com/kragniz/kube-onions
RUN cd /go/src/github.com/kragniz/kube-onions && make

# final stage
FROM scratch
WORKDIR /app
VOLUME /tmp
COPY --from=build-env /go/src/github.com/kragniz/kube-onions/onion-controller /app/
ENTRYPOINT ["./onion-controller"]
