.PHONY: help test
.DEFAULT_GOAL := help
DOCKER_TAG_DATE=$(shell date -u +%F.%H%M)

help: ## Displays this help message.
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)

test: ## Runs the tests and vetting.
	staticcheck $$(go list ./... | grep -v /vendor/)
	go test -cover -race -count=1 ./...
	go vet ./...
