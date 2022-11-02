.PHONY: build
build:
	go build -v -o bin/api ./cmd/api

.PHONY: run
run:
	./bin/api

.PHONY: run-rebuild
run-rebuild: build run

.PHONY: docker/build
docker/build:
	docker buildx build -t bghji/pthd-notifications . --platform=linux/amd64

.PHONY: docker/push
docker/push:
	docker push bghji/pthd-notifications

.PHONY: test
test:
	go test ./... -v

.PHONY: fmt
fmt:
	go fmt ./...
