BINARY  := lite-switch
MODULE  := github.com/nlink-jp/lite-switch
VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
LDFLAGS := -ldflags "-X main.version=$(VERSION)"

DIST_DIR := dist

# macOS Developer ID signing / notarization (see nlink-jp/.github
# CONVENTIONS.md §Code Signing). Defaults match any Developer ID
# Application cert in the keychain and the org-standard notary
# profile. Builds without these fall back to ad-hoc / un-notarized
# with a one-line warning — see scripts/codesign-darwin.sh.
CODESIGN_IDENTITY ?= Developer ID Application
NOTARY_PROFILE    ?= nlink-jp-notary

PLATFORMS := \
	linux/amd64 \
	linux/arm64 \
	darwin/amd64 \
	darwin/arm64 \
	windows/amd64

.PHONY: all build test vet lint check setup build-all package clean

## all: default target — build the binary
all: build

## build: compile the binary for the current platform
build:
	@mkdir -p $(DIST_DIR)
	go build $(LDFLAGS) -o $(DIST_DIR)/$(BINARY) .
	@scripts/codesign-darwin.sh $(DIST_DIR)/$(BINARY) "$(CODESIGN_IDENTITY)"

## test: run all unit tests
test:
	go test ./...

## vet: run go vet
vet:
	go vet ./...

## lint: run golangci-lint (must be installed)
lint:
	golangci-lint run ./...

## check: full quality gate — vet + lint + test + build + security scan
check: vet lint test build
	@echo "--- security scan ---"
	govulncheck ./...
	@echo "--- all checks passed ---"

## setup: install git hooks
setup:
	@echo "Installing git hooks..."
	@cp scripts/hooks/pre-commit .git/hooks/pre-commit
	@cp scripts/hooks/pre-push .git/hooks/pre-push
	@chmod +x .git/hooks/pre-commit .git/hooks/pre-push
	@echo "Git hooks installed."

## build-all: cross-compile for all target platforms
build-all:
	@mkdir -p $(DIST_DIR)
	$(foreach PLATFORM,$(PLATFORMS), \
		$(eval GOOS=$(word 1,$(subst /, ,$(PLATFORM)))) \
		$(eval GOARCH=$(word 2,$(subst /, ,$(PLATFORM)))) \
		$(eval EXT=$(if $(filter windows,$(GOOS)),.exe,)) \
		GOOS=$(GOOS) GOARCH=$(GOARCH) go build $(LDFLAGS) \
			-o $(DIST_DIR)/$(BINARY)-$(GOOS)-$(GOARCH)$(EXT) . ; \
	)
	@scripts/codesign-darwin.sh $(DIST_DIR)/$(BINARY)-darwin-amd64 "$(CODESIGN_IDENTITY)"
	@scripts/codesign-darwin.sh $(DIST_DIR)/$(BINARY)-darwin-arm64 "$(CODESIGN_IDENTITY)"
	@echo "Cross-compiled binaries in $(DIST_DIR)/"

## package: Build all platforms, zip with versioned naming + README, notarize darwin → dist/
package: build-all
	@cd $(DIST_DIR) && for f in $(BINARY)-*; do \
		case "$$f" in *.zip) continue ;; esac; \
		suffix=$${f#$(BINARY)-}; \
		suffix=$${suffix%%.exe}; \
		cp ../README.md .; \
		zip -j "$(BINARY)-$(VERSION)-$${suffix}.zip" "$$f" README.md; \
		rm -f README.md; \
	done
	@scripts/notarize-darwin.sh $(DIST_DIR)/$(BINARY)-$(VERSION)-darwin-amd64.zip "$(NOTARY_PROFILE)"
	@scripts/notarize-darwin.sh $(DIST_DIR)/$(BINARY)-$(VERSION)-darwin-arm64.zip "$(NOTARY_PROFILE)"

## clean: remove build artifacts
clean:
	@rm -rf $(DIST_DIR)
