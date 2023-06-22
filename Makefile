.PHONY: cover start test test-integration build-lambda

cover:
	go tool cover -html=cover.out

start:
	go run cmd/server/*.go

build-lambda:
	GOOS=linux go build -ldflags '-s -w' cmd/server/*.go 
	zip main_archived_cloud.zip main

test:
	go test -coverprofile=cover.out -short ./...

test-integration:
	go test -coverprofile=cover.out -p 1 ./...
