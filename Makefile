# Setup name variables for the package/tool
NAME := func

export GO111MODULE := on

# Set an output prefix, which is the local directory if not specified
PREFIX?=$(shell pwd)
# Set the build dir, where built cross-compiled binaries will be output
BUILDDIR := $(PREFIX)/dist

# Binary dependencies for this Makefile
BIN_DIR := $(GOPATH)/bin
LINTER := $(BIN_DIR)/golint
PACKR2 := $(BIN_DIR)/packr2
STATIC_CHECK := $(BIN_DIR)/staticcheck

all: help

.PHONY: ci ## Runs all tests and static code analysis
ci: build fmt lint test staticcheck vet

$(PACKR2):
	go get -u github.com/gobuffalo/packr/v2/packr2

.PHONY: build
build: $(PACKR2) ## Builds a static executable
	@echo "+ $@"
	@packr2
	@CGO_ENABLED=0 go build -o $(BUILDDIR)/$(NAME) .
	@packr2 clean
	@go mod tidy

.PHONY: fmt
fmt: ## Verifies all files have men `gofmt`ed
	@echo "+ $@"
	@gofmt -s -l . | grep -v vendor | tee /dev/stderr

$(LINTER):
	go get -u golang.org/x/lint/golint

.PHONY: lint
lint: $(LINTER) ## Verifies `golint` passes
	@echo "+ $@"
	@golint ./... | grep -v vendor | tee /dev/stderr

.PHONY: vet
vet: ## Verifies `go vet` passes
	@echo "+ $@"
	@go vet $(shell go list ./... | grep -v vendor) | tee /dev/stderr

$(STATIC_CHECK):
	go get -u honnef.co/go/tools/cmd/staticcheck

.PHONY: staticcheck
staticcheck: $(STATIC_CHECK) ## Verifies `staticcheck` passes
	@echo "+ $@"
	@staticcheck $(shell go list ./... | grep -v vendor)  | tee /dev/stderr

.PHONY: test
test: ## Runs all go tests
	@echo "+ $@"
	@go test -v -race $(shell go list ./... | grep -v vendor)

.PHONY: cover
cover: ## Runs all go tests (including integration tests) with coverage
	@echo "" > coverage.txt
	@for d in $(shell go list ./... | grep -v vendor); do \
		go test -race -coverprofile=profile.out -covermode=atomic "$$d"; \
		if [ -f profile.out ]; then \
			cat profile.out >> coverage.txt; \
			rm profile.out; \
		fi; \
	done;

.PHONY: clean
clean: $(PACKR2) ## Cleanup any build binaries or packages
	@echo "+ $@"
	@$(RM) -r $(BUILDDIR)
	@packr2 clean
	@go mod tidy

.PHONY: help
help: ## Display this help screen
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'