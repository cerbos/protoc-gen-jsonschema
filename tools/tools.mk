XDG_CACHE_HOME ?= $(HOME)/.cache
TOOLS_BIN_DIR := $(abspath $(XDG_CACHE_HOME)/cerbos/protoc-gen-jsonschema/bin)
TOOLS_MOD := tools/go.mod

BUF := $(TOOLS_BIN_DIR)/buf
GOLANGCI_LINT := $(TOOLS_BIN_DIR)/golangci-lint
GOTESTSUM := $(TOOLS_BIN_DIR)/gotestsum

$(TOOLS_BIN_DIR):
	@ mkdir -p $(TOOLS_BIN_DIR)

$(BUF): $(TOOLS_BIN_DIR)
	@ GOWORK=off GOBIN=$(TOOLS_BIN_DIR) go install -modfile=$(TOOLS_MOD) github.com/bufbuild/buf/cmd/buf

$(GOLANGCI_LINT): $(TOOLS_BIN_DIR)
	@ curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(TOOLS_BIN_DIR)

$(GOTESTSUM): $(TOOLS_BIN_DIR)
	@ GOWORK=off GOBIN=$(TOOLS_BIN_DIR) go install -modfile=$(TOOLS_MOD) gotest.tools/gotestsum
