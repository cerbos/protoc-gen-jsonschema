set dotenv-load := true

tools_mod_dir := join(justfile_directory(), "tools")

export TOOLS_BIN_DIR := join(env_var_or_default("XDG_CACHE_HOME", join(env_var("HOME"), ".cache")), "cerbos/protoc-gen-jsonschema/bin")

default:
    @ just --list

compile:
	@ go build -o bin/protoc-gen-jsonschema cmd/protoc-gen-jsonschema/main.go

deps:
    @ go mod tidy

# Run after testproto package is modified to generate new testdata
generate-testdata: _buf
	@ rm -rf $(TESTDATA_DIR)/code_generator_request.pb.bin
	@ cd $(TEST_DIR) && $(BUF) generate .

lint: _golangcilint _buf
    @ "${TOOLS_BIN_DIR}/golangci-lint" run --config=.golangci.yaml --fix
    @ "${TOOLS_BIN_DIR}/buf" format -w

tests PKG='./...' TEST='.*': _gotestsum
    @ "${TOOLS_BIN_DIR}/gotestsum" --format=dots-v2 --format-hide-empty-pkg -- -tags=tests,integration -failfast -count=1 -run='{{ TEST }}' '{{ PKG }}'

install:
	@ go install ./cmd/protoc-gen-jsonschema

_buf: (_install "buf" "github.com/bufbuild/buf" "cmd/buf")

_golangcilint: (_install "golangci-lint" "github.com/golangci/golangci-lint/v2" "cmd/golangci-lint")

_gotestsum: (_install "gotestsum" "gotest.tools/gotestsum")

_install EXECUTABLE MODULE CMD_PKG="":
    #!/usr/bin/env bash
    set -euo pipefail
    cd {{ tools_mod_dir }}
    TMP_VERSION=$(GOWORK=off go list -m -f "{{{{.Version}}" "{{ MODULE }}")
    VERSION="${TMP_VERSION#v}"
    BINARY="${TOOLS_BIN_DIR}/{{ EXECUTABLE }}"
    SYMLINK="${BINARY}-${VERSION}"
    if [[ ! -e "$SYMLINK" ]]; then
      echo "Installing $SYMLINK" 1>&2
      mkdir -p "$TOOLS_BIN_DIR"
      find "${TOOLS_BIN_DIR}" -lname "$BINARY" -delete
      if [[ "{{ EXECUTABLE }}" == "golangci-lint" ]]; then
        curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b "$TOOLS_BIN_DIR" "v${VERSION}"
      else
        export CGO_ENABLED={{ if EXECUTABLE =~ "(^sql|^tbls)" { "1" } else { "0" } }}
        GOWORK=off GOBIN="$TOOLS_BIN_DIR" go install {{ if CMD_PKG != "" { MODULE + "/" + CMD_PKG } else { MODULE } }}@v${VERSION}
      fi
      ln -s "$BINARY" "$SYMLINK"
    fi