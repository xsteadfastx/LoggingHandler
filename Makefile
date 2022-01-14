GO := go
GORELEASER := $(GO) run github.com/goreleaser/goreleaser
GOLANGCI_LINT := $(GO) run github.com/golangci/golangci-lint/cmd/golangci-lint

.PHONY: test
test:
	$(GO) test -v ./...

.PHONY: lint
lint:
	$(GOLANGCI_LINT) run --enable-all --disable=godox --timeout 10m

.PHONY: tidy
tidy:
	$(GO) mod tidy
	$(GO) mod vendor

.PHONY: install-tools
install-tools:
	$(GO) list -f '{{range .Imports}}{{.}} {{end}}' third_party/tools/tools.go | xargs go install -v

.PHONY: build
build:
	$(GORELEASER) build --rm-dist --snapshot

.PHONY: release
release:
	$(GORELEASER) release --rm-dist --snapshot --skip-publish
