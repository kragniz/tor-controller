all: onion-controller

push: docker
	docker push kragniz/kube-onions:latest

docker: onion-controller
	docker build . -t kragniz/kube-onions:latest

onion-controller: Makefile main.go
	CGO_ENABLED=0 GOOS=linux go build -a -ldflags '-extldflags "-static"' -o onion-controller

vendor: Gopkg.toml
	dep ensure

generate:
	./hack/update-codegen.sh
