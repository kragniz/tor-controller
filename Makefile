.PHONY: tor-daemon

build:
	go build cmd/controller-manager/main.go

tor-daemon:
	docker build . -f tor-daemon/Dockerfile -t kragniz/tor-daemon:latest
