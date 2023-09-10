#######
# Build
#######

.PHONY: build
build:
	go build -v -o bin/api ./cmd/api


##################
# Code style tools
##################

.PHONY: fmt
fmt:
	go fmt ./...

.PHONY: lint/fmt
lint/fmt:
	gofmt -l ./pkg ./cmd


.PHONY: lint/vet
lint/vet:
	go vet ./...

.PHONY:
lint: lint/fmt lint/vet

#########
# Testing
#########

.PHONY: test
test:
	go test ./... -v

#############
# Entrypoints
#############

.PHONY: run/api
run/api:
	./bin/api

.PHONY: run-rebuild
run-rebuild: build run/api

########
# Docker
########

.PHONY: docker/build
docker/build:
	docker buildx build -t bghji/pthd-notifications . --platform=linux/amd64


.PHONY: docker/push
docker/push:
	docker push bghji/pthd-notifications

.PHONY: local-deploy/infrastructure
local-deploy/infrastructure:
	docker-compose -f docker-compose.yml up -d

.PHONY: local-deploy/application
local-deploy/application:
	docker-compose -f docker-compose.yml --profile application up -d

