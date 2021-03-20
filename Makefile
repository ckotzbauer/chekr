BINDIR      := $(CURDIR)/bin
DIST_DIRS   := find * -type d -exec
TARGETS     := darwin/amd64 linux/amd64 linux/386 linux/arm linux/arm64 windows/amd64
TARGET_OBJS ?= darwin-amd64.tar.gz linux-amd64.tar.gz linux-386.tar.gz linux-arm.tar.gz linux-arm64.tar.gz windows-amd64.zip
BINNAME     ?= chekr

GOPATH        = $(shell go env GOPATH)
GOX           = $(GOPATH)/bin/gox
ARCH          = $(shell uname -p)

# go option
TAGS       :=
LDFLAGS    := -w -s
GOFLAGS    :=
SRC        := $(shell find . -type f -name '*.go' -print)

# Required for globs to work correctly
SHELL      = /usr/bin/env bash

.PHONY: all
all: build

# ------------------------------------------------------------------------------
#  build

.PHONY: build
build: $(BINDIR)/$(BINNAME)

$(BINDIR)/$(BINNAME): $(SRC)
	go build $(GOFLAGS) -tags '$(TAGS)' -ldflags '$(LDFLAGS)' -o $(BINDIR)/$(BINNAME)

# ------------------------------------------------------------------------------
#  dependencies

# If go get is run from inside the project directory it will add the dependencies
# to the go.mod file. To avoid that we change to a directory without a go.mod file
# when downloading the following dependencies

$(GOX):
	(cd /; go get -u github.com/mitchellh/gox)

# ------------------------------------------------------------------------------
#  release

.PHONY: build-cross
build-cross: LDFLAGS += -extldflags "-static"
build-cross: $(GOX)
	CGO_ENABLED=0 $(GOX) -parallel=3 -output="bin/{{.OS}}-{{.Arch}}/$(BINNAME)" -osarch='$(TARGETS)' $(GOFLAGS) -tags '$(TAGS)' -ldflags '$(LDFLAGS)'

.PHONY: dist
dist:
	( \
		cd bin && \
		$(DIST_DIRS) cp ../LICENSE {} \; && \
		$(DIST_DIRS) cp ../README.md {} \; && \
		$(DIST_DIRS) tar -zcf chekr-${VERSION}-{}.tar.gz {} \; && \
		$(DIST_DIRS) zip -r chekr-${VERSION}-{}.zip {} \; \
	)

.PHONY: clean
clean:
	@rm -rf $(BINDIR) ./bin
