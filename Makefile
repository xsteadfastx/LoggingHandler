GO := go

GORELEASER := $(GO) run github.com/goreleaser/goreleaser@v1.3.1
GOLANGCI_LINT := $(GO) run github.com/golangci/golangci-lint/cmd/golangci-lint@v1.44.2

.PHONY: test
test:
	$(GO) test -count=1 -v ./...

.PHONY: lint
lint:
	$(GOLANGCI_LINT) run \
		--enable-all \
		--disable=godox,varnamelen \
		--timeout 10m

.PHONY: tidy
tidy:
	$(GO) mod tidy
	$(GO) mod vendor

.PHONY: build
build:
	$(GORELEASER) build --rm-dist --snapshot

.PHONY: release
release:
	$(GORELEASER) release --rm-dist

.PHONY: release-snapshot
release-snapshot:
	$(GORELEASER) release --rm-dist --snapshot --skip-publish
