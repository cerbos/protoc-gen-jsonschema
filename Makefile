include tools/tools.mk

PROTOVALIDATE_MODULE := "buf.build/bufbuild/protovalidate"
PROTOS_DIR := "protos"
GENERATED_DIR := "gen"
TEST_DIR := "internal/test"
TESTDATA_DIR := "$(TEST_DIR)/testdata"

.PHONY: build
build: deps lint test compile install

.PHONY: compile
compile:
	@ go build -o bin/protoc-gen-jsonschema cmd/protoc-gen-jsonschema/main.go

.PHONY: deps
deps:
	@ go mod tidy --compat=1.21

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

# Run after testproto package is modified to generate new testdata
.PHONY: generate-testdata
generate-testdata: $(BUF)
	@ rm -rf $(TESTDATA_DIR)/code_generator_request.pb.bin
	@ cd $(TEST_DIR) && $(BUF) generate .
