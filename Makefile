all: onion-controller

push: docker
	docker push kragniz/kube-onions:latest

docker: onion-controller
	docker build . -t kragniz/kube-onions:latest

.PHONY: onion-controller
onion-controller:
	CGO_ENABLED=0 GOOS=linux go build -ldflags '-extldflags "-static"' -o onion-controller

vendor: Gopkg.toml
	dep ensure

generate:
	./hack/update-codegen.sh

kube-tor-daemon:
	docker build . -f tor-daemon/Dockerfile -t kragniz/kube-tor-daemon:latest
