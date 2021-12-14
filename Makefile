APPNAME := dabbi-sqs-consumer
OS := linux
PLATFORM := amd64
VERSION := $(strip $(shell cat VERSION.txt))
COMMIT := $(shell git rev-parse --short --verify HEAD)
DATE := $(shell date +%F-%T%z)

GO_LDFLAGS += -X main.appName=${APPNAME}
GO_LDFLAGS += -X main.buildVersion=${VERSION}
GO_LDFLAGS += -X main.buildCommit=${COMMIT}
GO_LDFLAGS += -X main.buildDate=${DATE}
GO_LDFLAGS := -ldflags="$(GO_LDFLAGS) -s -w"

# BUILD
.PHONY: build-app
build-app:
	@echo "# $@ : $(APPNAME)"
	GOOS=$(OS) GOARCH=$(PLATFORM) go build -v $(GO_LDFLAGS) -o ./$(APPNAME) ./cmd/$(APPNAME).go

# DOCKERIZE
define build-docker
docker build -t au.icr.io/$(strip $(1))/$(strip $(2)):$(strip $(3)) .
endef

.PHONY: build-docker-dev build-docker-prod
build-docker-dev:
	@echo "# $@: $(APPNAME):$(VERSION)"
	@$(call build-docker,"zmondev",$(APPNAME),$(VERSION))

build-docker-prod:
	@echo "# $@: $(APPNAME):$(VERSION)"
	@$(call build-docker,"zmon",$(APPNAME),$(VERSION))

.PHONY: all
all: build-app build-docker-dev clean

.PHONY: clean
clean:
	rm -f ./$(APPNAME)