GOOS=macos
GOARCH=arm64

.PHONY: build

GIT_COMMIT := $(shell git rev-list -1 HEAD)
GIT_BRANCH := $(shell git rev-parse --abbrev-ref HEAD)
GIT_TAG := $(shell git describe --tags --abbrev=0)
GIT_DIRTY := $(shell test -n "`git status --porcelain`" && echo "+CHANGES" || true)
BUILD_DATE := $(shell date -u '+%Y-%m-%d_%I:%M:%S%p')
BUILD_VERSION := $(GIT_TAG)$(GIT_DIRTY)
BUILD_HASH := $(GIT_COMMIT)

LDFLAGS:=-X main.GitCommit=${GIT_COMMIT}
LDFLAGS+=-X main.GitBranch=${GIT_BRANCH}
LDFLAGS+=-X main.GitTag=${GIT_TAG}
LDFLAGS+=-X main.GitDirty=${GIT_DIRTY}
LDFLAGS+=-X main.BuildDate=${BUILD_DATE}
LDFLAGS+=-X main.BuildVersion=${BUILD_VERSION}
LDFLAGS+=-X main.BuildHash=${BUILD_HASH}

build:
	rm -rf bin/
	@echo ${LDFLAGS}
	go build -ldflags="${LDFLAGS}" -o bin/cl main.go

test:
	go test -v ./...