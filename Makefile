.PHONY: help
## help: prints this help message
help:
	@echo "Usage: \n"
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' |  sed -e 's/^/ /'

BIN_PATH = bin
BIN_NAME = alarm-name

.PHONY: build
## build: builds app locally
build:
	@mkdir -p $(BIN_PATH)
	@go build -o $(BIN_PATH)/$(BIN_NAME) -v

.PHONY: test
## test: runs `go test`
test:
	@go test ./...

.PHONY: lint
## lint: runs `golangci-lint`
lint:
	@golangci-lint run ./...

.PHONY: run-local
## run-local: runs app locally
# example:
# make run-local ARGS="*/15 0 1,15 * 1-5 /usr/bin/find"
run-local:
	@go run -v main.go ${ARGS}

.PHONY: run-docker
## run: runs app in docker
# example:
# make run-docker ARGS="'*/15 0 1,15 * 1-5 /usr/bin/find'"
run-docker:
	docker build -t ${BIN_NAME} .
	docker run ${BIN_NAME} ./${BIN_NAME} ${ARGS}