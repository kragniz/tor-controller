.PHONY: tor-daemon

build:
	go build cmd/controller-manager/main.go

tor-daemon_docker:
	docker build . -f tor-daemon/Dockerfile -t kragniz/tor-daemon:latest

tor-controller_docker:
	docker build . -f Dockerfile.controller -t kragniz/tor-controller-manager:latest

images: tor-daemon_docker tor-controller_docker

install.yaml:
	kubebuilder create config --name=tor-controller --controller-image=kragniz/tor-controller-manager:latest --output=install.yaml
