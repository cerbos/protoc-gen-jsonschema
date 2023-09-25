include tools/tools.mk

PROTOVALIDATE_MODULE := "buf.build/bufbuild/protovalidate"
PROTOS_DIR := "protos"
GENERATED_DIR := "gen"
TEST_DIR := "internal/test"

.PHONY: build
build: deps generate lint test compile install

.PHONY: compile
compile:
	@ go build -o bin/protoc-gen-jsonschema cmd/protoc-gen-jsonschema/main.go

.PHONY: deps
deps:
	@ go mod tidy --compat=1.21

.PHONY: generate
generate: $(BUF) generate-buf generate-testdata

.PHONY: install
install:
	@ go install ./cmd/protoc-gen-jsonschema

.PHONY: lint
lint: $(BUF) $(GOLANGCI_LINT)
	@ $(GOLANGCI_LINT) run --config=.golangci.yaml --fix
	@ $(BUF) format -w

.PHONY: test
test: $(GOTESTSUM)
	@ CGO_ENABLED=0 $(GOTESTSUM) -- -tags=tests ./...

.PHONY: generate-buf
generate-buf: $(BUF)
	@ rm -rf $(GENERATED_DIR) $(PROTOS_DIR)
	@ $(BUF) export $(PROTOVALIDATE_MODULE) -o $(PROTOS_DIR)
	@ $(BUF) generate $(PROTOS_DIR)
	@ rm -rf $(PROTOS_DIR)

.PHONY: generate-testdata
generate-testdata: $(BUF)
	@ cd $(TEST_DIR) && $(BUF) generate .
