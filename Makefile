EXECUTABLE ?= nexus-cli
GO := CGO_ENABLED=0 go
DATE := $(shell date -u '+%FT%T%z')

LDFLAGS += -X main.Version=$(CI_COMMIT_TAG)
LDFLAGS += -extldflags '-static'

PACKAGES = $(shell go list ./... | grep -v /vendor/)

.PHONY: all
all: build

.PHONY: clean
clean:
	$(GO) clean -i ./...
	rm -rf dist/

.PHONY: fmt
fmt:
	$(GO) fmt $(PACKAGES)

.PHONY: vet
vet:
	$(GO) vet $(PACKAGES)

.PHONY: lint
lint:
	@which golint > /dev/null; if [ $$? -ne 0 ]; then \
		$(GO) get -u golang.org/x/lint/golint; \
	fi
	for PKG in $(PACKAGES); do golint -set_exit_status $$PKG || exit 1; done;

.PHONY: dep
dep:
	$(GO) mod vendor

.PHONY: build
build:
	$(GO) build -v -ldflags '-w $(LDFLAGS)' -o dist/$(EXECUTABLE) ./cmd/nexus-cli

.PHONY: release
release:
	@which gox > /dev/null; if [ $$? -ne 0 ]; then \
		$(GO) install github.com/mitchellh/gox@v1.0.1; \
	fi
	CGO_ENABLED=0 gox -arch="amd64" -verbose -ldflags '-w $(LDFLAGS)' -output="dist/$(EXECUTABLE)-{{.OS}}-{{.Arch}}" ./cmd/nexus-cli
