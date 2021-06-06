.DEFAULT_GOAL := build-darwin

#swagger installed as per https://goswagger.io/install.html
check-swagger:
	which swagger || (GO111MODULE=off go get -u github.com/go-swagger/go-swagger/cmd/swagger)

swagger: check-swagger
	swagger generate spec -o ./swagger.yaml --scan-models

swagger-ui: check-swagger
	swagger serve -F=swagger swagger.yaml

fmt:
	go fmt ./...
.PHONY:fmt

lint: fmt
	golint ./...
.PHONY:lint

vet: lint
	go vet ./...
.PHONY:vet

test: vet
	go test ./...
.PHONY:test

build-linux: test
	GOOS=linux GOARCH=amd64 go build -o nursery-linux-amd64 main.go
.PHONY:build-linux

build-darwin: test
	GOOS=darwin GOARCH=amd64 go build -o nusery-darwin-amd64 main.go
.PHONY:builddarwin