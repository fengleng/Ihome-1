
GOPATH:=$(shell go env GOPATH)


.PHONY: proto
proto:
	protoc --proto_path=${GOPATH}/src:. --micro_out=. --go_out=. proto/DeleteSession/DeleteSession.proto

.PHONY: build
build: proto

	go build -o DeleteSession-srv *.go

.PHONY: test
test:
	go test -v ./... -cover

.PHONY: docker
docker:
	docker build . -t DeleteSession-srv:latest
