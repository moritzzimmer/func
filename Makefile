# Setup name variables for the package/tool
NAME := func

# Set an output prefix, which is the local directory if not specified
PREFIX?=$(shell pwd)
BUILDDIR := $(PREFIX)/build

# App version
VERSION := $(shell cat VERSION.txt)

all: help

.PHONY: ci ## Runs all tests and static code analysis
ci: build fmt lint test staticcheck vet

.PHONY: build
build: ## Builds a static executable
	@echo "+ $@"
	@packr2
	@CGO_ENABLED=0 go build -o $(BUILDDIR)/$(NAME) .
	@packr2 clean

.PHONY: install
install:
	@echo "+ $@"
	@GO111MODULE=on packr2 install
	@packr2 clean
	
.PHONY: init-releaser
init-releaser: ## Initializes goreleaser for GitHub actions
	@echo "+ $@"
	@go mod tidy
	@packr2

.PHONY: fmt
fmt: ## Verifies all files have men `gofmt`ed
	@echo "+ $@"
	@gofmt -s -l . | grep -v vendor | tee /dev/stderr

.PHONY: lint
lint: ## Verifies `golint` passes
	@echo "+ $@"
	@golint ./... | grep -v vendor | tee /dev/stderr

.PHONY: vet
vet: ## Verifies `go vet` passes
	@echo "+ $@"
	@go vet $(shell go list ./... | grep -v vendor) | tee /dev/stderr

.PHONY: staticcheck
staticcheck: ## Verifies `staticcheck` passes
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

$(SEMBUMP):
	GO111MODULE=off go get -u github.com/jessfraz/junk/sembump

.PHONY: bump-version
BUMP := patch
bump-version: $(SEMBUMP) ## Bump the version in the version file. Set BUMP to [ patch | major | minor ].
	$(eval NEW_VERSION = $(shell $(BINDIR)/sembump --kind $(BUMP) $(VERSION)))
	@echo "Bumping VERSION.txt from $(VERSION) to $(NEW_VERSION)"
	echo $(NEW_VERSION) > VERSION.txt
	@echo "Updating links in README.md"
	sed -i '' s/$(subst v,,$(VERSION))/$(subst v,,$(NEW_VERSION))/g README.md
	git add VERSION.txt README.md
	git commit -vsam "Bump version to $(NEW_VERSION)"
	@echo "Run make tag to create and push the tag for new version $(NEW_VERSION)"

.PHONY: tag
tag: ## Create a new git tag to prepare to build a release
	git tag -a $(VERSION) -m "$(VERSION)"
	@echo "Run git push origin $(VERSION) to push your new tag to GitHub and trigger a build."

.PHONY: clean
clean: ## Cleanup any build binaries or packages
	@echo "+ $@"
	@$(RM) -r $(BUILDDIR)
	@packr2 clean

.PHONY: help
help: ## Display this help screen
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'