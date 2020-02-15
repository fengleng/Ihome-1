
GOPATH:=$(shell go env GOPATH)


.PHONY: build
build:

	go build -o IhomeWeb-web *.go

.PHONY: test
test:
	go test -v ./... -cover

.PHONY: docker
docker:
	docker build . -t IhomeWeb-web:latest
