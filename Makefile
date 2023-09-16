#############
# Environment
#############
.PHONY: install/golangci-lint
install/golangci-lint:
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.54.2

.PHONY: install-ci-deps
install: install/golangci-lint

#######
# Build
#######

.PHONY: build/api
build/api:
	go build -v -o bin/api ./cmd/api

.PHONY: build/async-api
build/async-api:
	go build -v -o bin/async-api ./cmd/async-api

.PHONY: build
build: build/api build/async-api

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

.PHONY: lint/golangci-lint
lint/golangci-lint:
	golangci-lint run ./...

.PHONY:
lint: lint/fmt lint/vet lint/golangci-lint

##########
# Generate
##########

.PHONY: generate
generate:
	go generate ./...

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

.PHONY: run/async-api
run/async-api:
	./bin/async-api

.PHONY: run/rebuild/api
run/rebuild/api: build run/api

.PHONY: run/rebuild/async-api
run/rebuild/async-api: build run/async-api

########
# Docker
########

.PHONY: docker/build
docker/build:
	docker buildx build -t bghji/pthd-notifications . --platform=linux/amd64

.PHONY: docker/push
docker/push:
	docker push bghji/pthd-notifications

.PHONY: local-deploy/infrastructure/up
local-deploy/infrastructure/up:
	docker-compose -f docker-compose.yml up -d

.PHONY: local-deploy/infrastructure/stop
local-deploy/infrastructure/stop:
	docker-compose -f docker-compose.yml stop

.PHONY: local-deploy/application/up
local-deploy/application/up:
	docker-compose -f docker-compose.yml --profile application up -d

.PHONY: local-deploy/application/stop
local-deploy/application/stop:
	docker-compose -f docker-compose.yml --profile application stop
