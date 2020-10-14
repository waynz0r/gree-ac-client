.DEFAULT_GOAL := help
.PHONY: help build

OS := $(shell uname -s)
GOFILES_NOVENDOR = $(shell find . -type f -name '*.go' -not -path "./vendor/*" -not -path "./client/*")
SYMLINKS=$(shell find -L ./vendor -type l)

PKGS=$(shell go list ./... | grep -v /client)

VERSION = 0.1.0
GITREV = $(shell git rev-parse --short HEAD)

build: ## Builds binary package
	go build  -ldflags "-X main.Version=$(VERSION) -X main.GitRev=$(GITREV)" -o gree-ac-client cmd/cli/main.go

build-ci:
	CGO_ENABLED=0 GOOS=linux go build .

clean:
	rm -f pipeline

help: ## Generates this help message
	@grep -E '^[0-9a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

list:
	@$(MAKE) -pRrn : -f $(MAKEFILE_LIST) 2>/dev/null | awk -v RS= -F: '/^# File/,/^# Finished Make data base/ {if ($$1 !~ "^[#.]") {print $$1}}' | egrep -v -e '^[^[:alnum:]]' -e '^$@$$' | sort

fmt:
	@gofmt -w ${GOFILES_NOVENDOR}

check-fmt:
	PKGS="${GOFILES_NOVENDOR}" GOFMT="gofmt" ./scripts/fmt-check.sh

check-misspell: install-misspell
	PKGS="${GOFILES_NOVENDOR}" MISSPELL="misspell" ./scripts/misspell-check.sh

misspell: install-misspell
	misspell -w ${GOFILES_NOVENDOR}

vet:
	@go vet -composites=false ./...

lint: install-golint
	golint -min_confidence 0.9 -set_exit_status $(PKGS)

install-golint:
	GOLINT_CMD=$(shell command -v golint 2> /dev/null)
ifndef GOLINT_CMD
	go get github.com/golang/lint/golint
endif

install-misspell:
	MISSPELL_CMD=$(shell command -v misspell 2> /dev/null)
ifndef MISSPELL_CMD
	go get -u github.com/client9/misspell/cmd/misspell
endif

clean-vendor:
	find -L ./vendor -type l | xargs rm -rf

ineffassign: install-ineffassign
	ineffassign ${GOFILES_NOVENDOR}

gocyclo: install-gocyclo
	gocyclo -over 15 ${GOFILES_NOVENDOR}

install-ineffassign:
	INEFFASSIGN_CMD=$(shell command -v ineffassign 2> /dev/null)
ifndef INEFFASSIGN_CMD
	go get -u github.com/gordonklaus/ineffassign
endif

install-gocyclo:
	GOCYCLO_CMD=$(shell command -v gocyclo 2> /dev/null)
ifndef GOCYCLO_CMD
	go get -u github.com/fzipp/gocyclo
endif
