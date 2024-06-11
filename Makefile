# 
# Makefile for catalogue
# 
GOCMD=go
GOBUILD=$(GOCMD) build
GORUN=$(GOCMD) run
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOTOOL=$(GOCMD) tool
GOGET=$(GOCMD) get
GOPATH?=`$(GOCMD) env GOPATH`

BINARY=catalogue
BINARY_DARWIN=$(BINARY)_darwin
BINARY_LINUX=$(BINARY)_linux
BINARY_WINDOWS=$(BINARY).exe

TESTS=./...
COVERAGE_FILE=coverage.out

SOURCES_PATH=github.com/creativeprojects/catalogue
BUILD=build/
GO_VERSION=1.14
DOCKER_TAG=creativeprojects/catalogue

.PHONY: all test test-short build build-mac build-linux build-windows build-all coverage clean test-docker build-docker ramdisk run-docker

all: test build

build:
		$(GOBUILD) -o $(BINARY) -v

build-mac:
		GOOS="darwin" GOARCH="amd64" $(GOBUILD) -o $(BINARY_DARWIN)_x86 -v
		GOOS="darwin" GOARCH="arm64" $(GOBUILD) -o $(BINARY_DARWIN)_arm64 -v

build-linux:
		GOOS="linux" GOARCH="amd64" $(GOBUILD) -o $(BINARY_LINUX)_x86 -v
		GOOS="linux" GOARCH="arm64" $(GOBUILD) -o $(BINARY_LINUX)_arm64 -v

build-windows:
		GOOS="windows" GOARCH="amd64" $(GOBUILD) -o $(BINARY_WINDOWS) -v

build-all: build-mac build-linux build-windows

test: ramdisk
		DB_TEST_PATH=/Volumes/RAMDisk $(GOTEST) -v $(TESTS)

test-short:
		$(GOTEST) -v -short $(TESTS)

coverage: ramdisk
		DB_TEST_PATH=/Volumes/RAMDisk $(GOTEST) -coverprofile=$(COVERAGE_FILE) $(TESTS)
		$(GOTOOL) cover -html=$(COVERAGE_FILE)

clean:
		$(GOCLEAN)
		rm -f $(BINARY) $(BINARY_DARWIN) $(BINARY_LINUX) $(BINARY_WINDOWS) $(COVERAGE_FILE) ${BUILD}$(BINARY)*

test-docker:
		docker run --rm -v "${GOPATH}":/go -w /go/src/${SOURCES_PATH} golang:${GO_VERSION} $(GOTEST) -v $(TESTS)

build-docker: clean
		CGO_ENABLED=0 GOARCH=amd64 GOOS=linux $(GOBUILD) -v -o ${BUILD}$(BINARY) .
		cd ${BUILD}; docker build --pull --tag ${DOCKER_TAG} .

ramdisk: /Volumes/RAMDisk

/Volumes/RAMDisk:
		diskutil erasevolume HFS+ RAMDisk `hdiutil attach -nomount ram://4194304`

run-docker:
		docker run --rm -v "${GOPATH}":/go -v "${PWD}/.cache":/root/.cache -w /go/src/${SOURCES_PATH} golang:${GO_VERSION} ${GORUN} -v . ${PARAMS}
