all: kube-onion

push: docker
	docker push kragniz/kube-onion:latest

docker: kube-onion
	docker build . -t kragniz/kube-onion:latest

kube-onion: Makefile main.go
	CGO_ENABLED=0 GOOS=linux go build -a -ldflags '-extldflags "-static"' -o kube-onion

vendor: Gopkg.toml
	dep ensure

