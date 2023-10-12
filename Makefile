GOPATH := $(shell go env GOPATH)
GOBIN := $(GOPATH)/bin

default: help

.PHONY: build
build: ## Build the shockrays binary in the build sub-directory of the project's top level directory
ifeq ($(OS),Windows_NT)
	@echo "building shockrays..."
	@go build -mod=readonly $(BUILD_FLAGS) -o build/shockrays.exe main.go
else
	@echo "building shockrays..."
	@go build $(BUILD_FLAGS) -o build/shockrays main.go
endif

.PHONY: install
install: ## Install the shockrays binary to the $GOBIN directory
ifeq ($(OS),Windows_NT)
	@echo "installing shockrays..."
	@go build -mod=readonly $(BUILD_FLAGS) -o $(GOBIN)/shockrays.exe main.go
else
	@echo "installing shockrays..."
	@go build -mod=readonly $(BUILD_FLAGS) -o $(GOBIN)/shockrays main.go
endif

.PHONY: help
help: ## Prints the available Make targets, this will not work on Windows unless using WSL or installing sort, grep & awk
	@echo "Available make commands:"; grep -h -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
