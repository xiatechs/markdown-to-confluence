.PHONY: help
help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

.PHONY: test
test: ## Test the Project with coverage and race detector.
	go test -v -cover -race ./...

.PHONY: lint
lint: ## Lint the Project using golangci-lint. If not installed this will install golangci-lint.
	@command -v golangci-lint > /dev/null || \
		GO111MODULE=off go get -u github.com/golangci/golangci-lint/cmd/golangci-lint

	golangci-lint run

.PHONY: update
update: ## Update the Project's dependencies.
	@go get -u ./...
	go mod tidy

.PHONY: run-docker
run-docker: ## Run the Project in the docker container used for it's action.
	docker build . -t mdtc
	docker run -v $(shell pwd):/github/workspace mdtc 