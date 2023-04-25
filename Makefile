UNAME := $(shell uname)
GO_PATH:=$(shell go env GOPATH)

LINTER_ARGS = run -c .golangci.yml --timeout 5m
GODOG_FORMAT = pretty
GOLANGCI_LINT_VERSION = v1.52.2

.PHONY: help
help:	## Show a list of available commands
	@fgrep -h "##" $(MAKEFILE_LIST) | fgrep -v fgrep | sed -e 's/\\$$//' | sed -e 's/##//'

.PHONY: test-acceptance
test-acceptance:	## Run acceptance tests
	go test ./test/acceptance -v --godog.format=$(GODOG_FORMAT)

.PHONY: test-run-tag
test-run-tag:	## Run feature from tag using variable TAG='<@tag_name>'
	go test ./test/acceptance -v --godog.tags=$(TAG) --godog.format=$(GODOG_FORMAT)

.PHONY: download-tools
download-tools:	## Download all required tools to validate and generate documentation, code analysis...
	@echo "Installing tools on $(GO_PATH)/bin"
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@$(GOLANGCI_LINT_VERSION)
	@echo "Tools installed"

.PHONY: lint
lint:	## Run static linting of source files. See .golangci.yml for options
	golangci-lint $(LINTER_ARGS)
