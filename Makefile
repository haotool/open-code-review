.PHONY: build build-full test clean run help fmt vet fmt-check doc-whitespace-check check coverage install-skill \
	plugin-check delegate-surface-check host-agent-smoke \
	build-all dist sha256sum version-info \
	build-linux-amd64 build-linux-arm64 build-darwin-amd64 build-darwin-arm64 \
	build-windows-amd64 build-windows-arm64

DELEGATE_BINARY := ocr-delegate
FULL_BINARY     := opencodereview
GO              := go
DIST_DIR        := ./dist

# Version info — use git tag if available, fallback to short commit hash
GIT_TAG     := $(shell git describe --tags --abbrev=0 2>/dev/null || echo "")
GIT_COMMIT  := $(shell git rev-parse --short HEAD)
BUILD_DATE  := $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")

VERSION     ?= $(if $(GIT_TAG),$(GIT_TAG),v0.0.0-$(GIT_COMMIT))

LD_FLAGS    := \
	-X main.Version=$(VERSION) \
	-X main.GitCommit=$(GIT_COMMIT) \
	-X main.BuildDate=$(BUILD_DATE)

RELEASE_LD_FLAGS := -s -w $(LD_FLAGS)

define BUILD_PLATFORM
	GOOS=$(1) GOARCH=$(2) CGO_ENABLED=0 $(GO) build -trimpath -ldflags "$(RELEASE_LD_FLAGS)" \
		-o $(DIST_DIR)/$(DELEGATE_BINARY)-$(1)-$(2)$(3) \
		./cmd/ocr-delegate
endef

# ── Development targets ──────────────────────────────────────────────────────
build:
	CGO_ENABLED=0 $(GO) build -trimpath -ldflags "$(RELEASE_LD_FLAGS)" \
		-o $(DIST_DIR)/$(DELEGATE_BINARY) ./cmd/ocr-delegate

build-full:
	$(GO) build -ldflags "$(LD_FLAGS)" -o $(DIST_DIR)/$(FULL_BINARY) ./cmd/opencodereview

SKILL_DIR ?= $(HOME)/.cursor/skills
SKILL_SRC := skills/open-code-review
SKILL_DELEGATE_SRC := skills/open-code-review-delegate

install-skill:
	mkdir -p $(SKILL_DIR)/open-code-review $(SKILL_DIR)/open-code-review-delegate
	cp -R $(SKILL_SRC)/. $(SKILL_DIR)/open-code-review/
	cp -R $(SKILL_DELEGATE_SRC)/. $(SKILL_DIR)/open-code-review-delegate/

plugin-check:
	bash test/ci/check-plugin-skills.sh

delegate-surface-check:
	bash test/ci/check-delegate-surface.sh

host-agent-smoke:
	bash test/e2e/host-agent-smoke.sh

PACKAGES := $(shell $(GO) list ./... | grep -v /extensions/)

test:
	LC_ALL=C $(GO) test -v -race -count=1 $(PACKAGES)

COVERAGE_THRESHOLD := 80

coverage:
	LC_ALL=C $(GO) test -count=1 -coverprofile=coverage.out $(PACKAGES)
	$(GO) tool cover -func=coverage.out | grep total:
	@COVERAGE=$$($(GO) tool cover -func=coverage.out | grep total: | awk '{print $$3}' | sed 's/%//'); \
	if awk "BEGIN {exit !($$COVERAGE < $(COVERAGE_THRESHOLD))}"; then \
		echo "FAIL: Coverage $${COVERAGE}% is below $(COVERAGE_THRESHOLD)% threshold"; \
		exit 1; \
	fi; \
	echo "PASS: Coverage $${COVERAGE}% meets $(COVERAGE_THRESHOLD)% threshold"

clean:
	rm -rf $(DIST_DIR) coverage.out

run: build-full
	$(DIST_DIR)/$(FULL_BINARY) --staged

help: build
	$(DIST_DIR)/$(DELEGATE_BINARY) -h

fmt:
	gofmt -s -w .

vet:
	LC_ALL=C $(GO) vet $(PACKAGES)

fmt-check:
	@unformatted=$$(gofmt -s -l $$( $(GO) list -f '{{.Dir}}' $(PACKAGES) )); \
	if [ -n "$$unformatted" ]; then \
		echo "FAIL: gofmt needed on:"; \
		echo "$$unformatted"; \
		exit 1; \
	fi; \
	echo "fmt-check passed"

doc-whitespace-check:
	@bash test/ci/check-doc-whitespace.sh

check:
	$(GO) mod tidy
	gofmt -s -w .
	LC_ALL=C $(GO) vet $(PACKAGES)
	@echo "check passed"

# ── Cross-platform targets (local reproduction only; no release upload) ───────
build-linux-amd64:
	$(call BUILD_PLATFORM,linux,amd64)

build-linux-arm64:
	$(call BUILD_PLATFORM,linux,arm64)

build-darwin-amd64:
	$(call BUILD_PLATFORM,darwin,amd64)

build-darwin-arm64:
	$(call BUILD_PLATFORM,darwin,arm64)

build-windows-amd64:
	$(call BUILD_PLATFORM,windows,amd64,.exe)

build-windows-arm64:
	$(call BUILD_PLATFORM,windows,arm64,.exe)

build-all: build-linux-amd64 build-linux-arm64 build-darwin-amd64 build-darwin-arm64 build-windows-amd64 build-windows-arm64

sha256sum: build-all
	cd $(DIST_DIR) && shasum -a 256 $(DELEGATE_BINARY)-* | sort > sha256sum.txt

dist: clean build-all sha256sum
	@echo $(VERSION) > $(DIST_DIR)/VERSION

version-info:
	@echo "Version:   $(VERSION)"
	@echo "GitCommit: $(GIT_COMMIT)"
	@echo "BuildDate: $(BUILD_DATE)"
	@echo "LD_FLAGS:  $(LD_FLAGS)"
