.PHONY: cover start test test-integration build-lambda deploy build

cover:
	go tool cover -html=cover.out

start:
	go run cmd/server/*.go

build-lambda:
	GOOS=linux go build -o bin/main -ldflags '-s -w' cmd/server/*.go 
	zip bin_archived/main_archived.zip bin/main

build:
	go build -o bin/main -ldflags '-s -w' cmd/server/*.go 

deploy:
	cd terraform/ && terraform apply

test:
	go test -coverprofile=cover.out -short ./...

test-integration:
	go test -coverprofile=cover.out -p 1 ./...
