WORKDIR      := $(shell pwd)
NATIVEOS	 := $(shell go version | awk -F '[ /]' '{print $$4}')
NATIVEARCH	 := $(shell go version | awk -F '[ /]' '{print $$5}')
INTEGRATION  := mongodb
BINARY_NAME   = nri-$(INTEGRATION)
GO_PKGS      := $(shell go list ./... | grep -v "/vendor/")
GO_FILES     := ./src/
GOFLAGS			= -mod=readonly


all: build

build: clean test compile

clean:
	@echo "=== $(INTEGRATION) === [ clean ]: Removing binaries..."
	@rm -rfv bin

compile:
	@echo "=== $(INTEGRATION) === [ compile ]: Building $(BINARY_NAME)..."
	@go build -o bin/$(BINARY_NAME) $(GO_FILES)

test:
	@echo "=== $(INTEGRATION) === [ test ]: Running unit tests..."
	@go test -race ./... -count=1

integration-test:
	@echo "=== $(INTEGRATION) === [ test ]: running integration tests..."
	@docker compose -f tests/integration/docker-compose.yml up -d --build
	@go test -v -tags=integration ./tests/integration/. ; (ret=$$?; docker compose -f tests/integration/docker-compose.yml down && exit $$ret)

# rt-update-changelog runs the release-toolkit run.sh script by piping it into bash to update the CHANGELOG.md.
# It also passes down to the script all the flags added to the make target. To check all the accepted flags,
# see: https://github.com/newrelic/release-toolkit/blob/main/contrib/ohi-release-notes/run.sh
#  e.g. `make rt-update-changelog -- -v`
rt-update-changelog:
	curl "https://raw.githubusercontent.com/newrelic/release-toolkit/v1/contrib/ohi-release-notes/run.sh" | bash -s -- $(filter-out $@,$(MAKECMDGOALS))

# Include thematic Makefiles
include $(CURDIR)/build/ci.mk
include $(CURDIR)/build/release.mk

.PHONY: all build clean compile test rt-update-changelog
