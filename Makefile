.PHONY: cover start test test-integration build-lambda deploy-lambda build

export image := `aws lightsail get-container-images --service-name canvas | jq -r '.containerImages[0].image'`

cover:
	go tool cover -html=cover.out

deploy:
	aws lightsail push-container-image --service-name canvas --label app --image canvas
	aws lightsail create-container-service-deployment --service-name canvas \
		--containers '{"app":{"image":"'$(image)'","environment":{"HOST":"","PORT":"8080","LOG_ENV":"production"},"ports":{"8080":"HTTP"}}}' \
		--public-endpoint '{"containerName":"app","containerPort":8080,"healthCheck":{"path":"/health"}}'

start:
	go run cmd/server/*.go

build-lambda:
	GOOS=linux go build -o bin/main -ldflags '-s -w' cmd/server/*.go 
	zip bin_archived/main_archived.zip bin/main

build:
	go build -o bin/main -ldflags '-s -w' cmd/server/*.go 

deploy-lambda:
	cd terraform/ && terraform apply

test:
	go test -coverprofile=cover.out -short ./...

test-integration:
	go test -coverprofile=cover.out -p 1 ./...
