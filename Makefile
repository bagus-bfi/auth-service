.PHONY: vendor
vendor:
	go mod tidy && go mod vendor

# copy .env.example so we can adjust without affecting git file change
# actual server run does not required .env file to be exist.
.PHONY: env
env:
	cp .env.example .env

.PHONY: server/start
server/start:
	docker-compose -f docker-compose.yml up -d

.PHONY: server/restart
server/restart:
	docker-compose -f docker-compose.yml restart db-schema grpc

.PHONY: server/stop
server/stop:
	docker-compose -f docker-compose.yml down

GOLANGCI_VERSION=1.50.1
GOLANGCI_CHECK := $(shell golangci-lint -v 2> /dev/null)

.PHONY: lint
lint:
# if golangci-lint failed on MacOS Ventura, try: brew install diffutils
ifndef GOLANGCI_CHECK
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin v$(GOLANGCI_VERSION)
endif
	golangci-lint run -c .golangci.yml ./...

COVERAGE_OUTPUT=coverage.out
COVERAGE_OUTPUT_HTML=coverage.html
TEST_LIST := $(shell go list ./... | grep -v mock | grep -v /vendor/)

.PHONY: test
test:
	go test $(TEST_LIST) -race -coverprofile=$(COVERAGE_OUTPUT) -coverpkg=./...
	go tool cover -html=$(COVERAGE_OUTPUT) -o $(COVERAGE_OUTPUT_HTML)

################
# BUILD BINARY
################

GOOS ?= linux
GOARCH ?= amd64
CGO_ENABLED ?= 0
GIT_URL := $(shell git config --get remote.origin.url)
GIT_COMMIT_HASH := $(shell git rev-parse --short HEAD)
GIT_BRANCH := $(shell git branch --show-current)
GIT_TAG := $(shell git tag --points-at HEAD)
GIT_TAG := $(or $(GIT_TAG),$(GIT_BRANCH))
BUILD_OS := $(shell uname -rns | sed -e 's/ /_/g') # replace space with _
BUILD_TIME := $(shell date -u +%FT%T%Z)
# GO_MOD_NAME=$(shell go list -m)
GO_APP_INFO_MOD_NAME=github.com/bfi-finance/bfi-go-pkg
# default to static build
GO_BUILD_FLAGS=-trimpath -a -tags "osusergo,netgo" -ldflags '-extldflags=-static -w -s -v' \
-ldflags "-X '$(GO_APP_INFO_MOD_NAME)/appinfo.GitURL=$(GIT_URL)' \
-X '$(GO_APP_INFO_MOD_NAME)/appinfo.GitCommitHash=$(GIT_COMMIT_HASH)' \
-X '$(GO_APP_INFO_MOD_NAME)/appinfo.GitTag=$(GIT_TAG)' \
-X '$(GO_APP_INFO_MOD_NAME)/appinfo.BuildTime=$(BUILD_TIME)' \
-X '$(GO_APP_INFO_MOD_NAME)/appinfo.BuildOS=$(BUILD_OS)'"

.PHONY: build/grpc
build/grpc:
	CGO_ENABLED=$(CGO_ENABLED) GOOS=$(GOOS) GOARCH=$(GOARCH) go build $(GO_BUILD_FLAGS) -o build/grpc cmd/grpc/main.go

################
# BUILD DOCKER IMAGE
################

DOCKER_IMAGE_NAME=$(or $(APP_NAME),service-name)

.PHONY: docker/db-schema
docker/db-schema:
	docker build --no-cache --rm -f Dockerfile.db-schema -t $(DOCKER_IMAGE_NAME)-db-schema:$(GIT_TAG) .

.PHONY: docker/grpc
docker/grpc:
	docker build --no-cache --rm -f Dockerfile -t $(DOCKER_IMAGE_NAME)-grpc:$(GIT_TAG) .
