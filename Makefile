build:
	go build -v -o bin/api ./cmd/api

run:
	./bin/api

run-rebuild: build run

docker/build:
	docker buildx build -t bghji/pthd-notifications . --platform=linux/amd64

docker/push:
	docker push bghji/pthd-notifications

test:
	go test ./... -v

fmt:
	go fmt ./...
