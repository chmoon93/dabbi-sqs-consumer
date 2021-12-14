APPNAME := dabbi-sqs-consumer
PKG := github.com/chmoon93/$(APPNAME)
OS := linux
PLATFORM := amd64
DATE := $(shell date +%F-%T%z)

VERSION := $(strip $(shell cat VERSION.txt))
GITCOMMIT := $(shell git rev-parse --short --verify HEAD)
GITUNTRACKEDCHANGES := $(shell git status --porcelain --untracked-files=no) -X main.buildDate=${DATE}
ifneq ($(GITUNTRACKEDCHANGES),)
	GITCOMMIT := $(GITCOMMIT)-dirty
endif
CTIMEVAR=-X main.appName=${APPNAME} -X main.buildVersion=${VERSION} -X main.buildCommit=${GITCOMMIT} -X main.buildDate=${DATE}
GO_LDFLAGS=-ldflags "-w -s $(CTIMEVAR)"

#---

# define build-target
# env	GOOS=$(OS) GOARCH=$(PLATFORM) go build -v $(GO_LDFLAGS) -o ./$(APPNAME) ./cmd/$(APPNAME).go
# endef

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
all: build-app # build-docker-dev clean

.PHONY: clean
clean:
	rm -f ./$(APPNAME)