.ONESHELL:
.PHONY: test

test:
	go test -v -cover -race ./...

lint:
	@command -v golangci-lint > /dev/null || \
		GO111MODULE=off go get -u github.com/golangci/golangci-lint/cmd/golangci-lint

	golangci-lint run

update:
	@go get -u ./...
	go mod tidy