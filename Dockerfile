# build stage
FROM golang:alpine AS build-env
RUN apk update && apk add make
ADD . /go/src/github.com/kragniz/kube-onion
RUN cd /go/src/github.com/kragniz/kube-onion && make

# final stage
FROM scratch
WORKDIR /app
COPY --from=build-env /go/src/github.com/kragniz/kube-onion/kube-onion /app/
ENTRYPOINT ["./kube-onion"]
