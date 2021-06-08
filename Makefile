BLACK     := $(shell tput setab 7)$(shell tput -Txterm setaf 0)
RED       := $(shell tput setab 0)$(shell tput -Txterm setaf 1)
GREEN     := $(shell tput setab 0)$(shell tput -Txterm setaf 2)
YELLOW    := $(shell tput setab 0)$(shell tput -Txterm setaf 3)
BLUE      := $(shell tput setab 4)$(shell tput setaf 7)
PURPLE    := $(shell tput setab 0)$(shell tput -Txterm setaf 5)
LIGHTBLUE := $(shell tput setab 0)$(shell tput -Txterm setaf 6)
WHITE     := $(shell tput setab 0)$(shell tput -Txterm setaf 7)

RESET := $(shell tput -Txterm sgr0)
CLEAR := -en "\033c"

HAS_GIT := $(shell ls .git 2> /dev/null)
BRANCHENAME=
ifneq ("$(HAS_GIT)","")
	BRANCHENAME=$(shell git branch --show-current)
endif
PROJECTNAME=$(shell basename "$(PWD)")

GOTEST_PRESENT := $(shell command -v gotest 2> /dev/null)
GOTEST=go test
ifneq ("$(GOTEST_PRESENT)","")
 GOTEST=~/go/bin/gotest
endif
READY := "$(GREEN) > Ready: $(BLUE) $(PROJECTNAME) $(YELLOW) $(BRANCHENAME)$(RESET)"


.PHONY: help
help: ## Print this help message
	@echo ${CLEAR}
	@echo "$(LIGHTBLUE) > help $(BLUE) $(PROJECTNAME) $(YELLOW) $(BRANCHENAME)$(RESET)"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
	@echo ${READY}

## check: run golangci-lint
check:
	@echo ${CLEAR}
	@echo "$(LIGHTBLUE) > Checking with golangci-lint $(BLUE) $(PROJECTNAME) $(YELLOW) $(BRANCHENAME)$(RESET)"
	@golangci-lint run --fix
	@echo ${READY}

.PHONY: test
test: ## Test the Project with coverage and race detector.
	@echo ${CLEAR}
	@echo "$(LIGHTBLUE) > Testing $(BLUE) $(PROJECTNAME) $(YELLOW) $(BRANCHENAME)$(RESET)"
	@GOPATH=$(GOPATH) GOBIN=$(GOBIN) $(GOTEST) $(VENDOR) -count=1 -race $(TESTTIMEOUT) ./...
	@echo ${READY}

## test-c: run all tests with codecoverage
test-c:
	@echo ${CLEAR}
	@echo "$(LIGHTBLUE) > Testing with code codecoverage $(BLUE) $(PROJECTNAME) $(YELLOW) $(BRANCHENAME)$(RESET)"
	@GOPATH=$(GOPATH) GOBIN=$(GOBIN) $(GOTEST) $(VENDOR) -count=1 $(TESTTIMEOUT) -coverprofile coverage.out ./...
	@echo "Coverage overall:"
	@go tool cover -func=coverage.out | tail -n 1
	@GOPATH=$(GOPATH) GOBIN=$(GOBIN) go tool cover -html coverage.out
	#@gocheckcov check init --profile-file coverage.out
	#@gocheckcov check --minimum-coverage 80 --profile-file coverage.out --print-src
	@echo ${READY}

.PHONY: test-v
test-v: ## Test the Project verbose and race detector.
	@echo ${CLEAR}
	@echo "$(LIGHTBLUE) > Testing verbose $(BLUE) $(PROJECTNAME) $(YELLOW) $(BRANCHENAME)$(RESET)"
	@GOPATH=$(GOPATH) GOBIN=$(GOBIN) $(GOTEST) $(VENDOR) -v -count=1 -race $(TESTTIMEOUT) ./...
	@echo ${READY}

.PHONY: lint
lint: ## Lint the Project using golangci-lint. If not installed this will install golangci-lint.
	@echo ${CLEAR}
	@echo "$(LIGHTBLUE) > Checking with golangci-lint $(BLUE) $(PROJECTNAME) $(YELLOW) $(BRANCHENAME)$(RESET)"
	@command -v golangci-lint > /dev/null || \
		GO111MODULE=off go get -u github.com/golangci/golangci-lint/cmd/golangci-lint
	@golangci-lint run --fix
	@echo ${READY}

.PHONY: update
update: ## Update the Project's dependencies.
	@echo ${CLEAR}
	@echo "$(LIGHTBLUE) > Updating mod-file $(BLUE) $(PROJECTNAME) $(YELLOW) $(BRANCHENAME)$(RESET)"
	@rm -f go.sum

	@GOPATH=$(GOPATH) GOBIN=$(GOBIN) go mod vendor
	@GOPATH=$(GOPATH) GOBIN=$(GOBIN) go get -u ./...
	@GOPATH=$(GOPATH) GOBIN=$(GOBIN) go mod tidy
	@GOPATH=$(GOPATH) GOBIN=$(GOBIN) go mod vendor
	@echo ${READY}

.PHONY: run-docker
run-docker: ## Run the Project in the docker container used for it's action.
	@echo ${CLEAR}
	@echo "$(LIGHTBLUE) > Creating and run Docker.. $(BLUE) $(PROJECTNAME) $(YELLOW) $(BRANCHENAME)$(RESET)"
	@docker build . -t mdtc
	@docker run -v $(shell pwd):/github/workspace mdtc 
	@echo ${READY}
