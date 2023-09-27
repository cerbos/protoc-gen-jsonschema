// Copyright 2021-2023 Zenauth Ltd.
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"

	pgs "github.com/lyft/protoc-gen-star/v2"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/pluginpb"

	"github.com/cerbos/protoc-gen-jsonschema/internal/common"
	"github.com/cerbos/protoc-gen-jsonschema/internal/module"
)

const (
	debugPrintRequest = "PGJS_DEBUG_REQUEST"
	requestPath       = "internal/test/testdata/code_generator_request.pb.bin"
)

func main() {
	if _, exists := os.LookupEnv(debugPrintRequest); exists {
		if err := printRequest(requestPath); err != nil {
			log.Fatalf("failed to print request: %s", err.Error())
		}
	}

	reqFile, err := os.Open(requestPath)
	if err != nil {
		log.Fatalf("failed to open code generator request file: %s", err.Error())
	}

	resBytes := &bytes.Buffer{}
	pgs.Init(
		pgs.DebugEnv(common.DebugEnv),
		pgs.ProtocInput(reqFile),
		pgs.ProtocOutput(resBytes),
	).RegisterModule(module.New()).Render()

	res := &pluginpb.CodeGeneratorResponse{}
	if err := proto.Unmarshal(resBytes.Bytes(), res); err != nil {
		log.Fatalf("failed to unmarshal code generator response: %s", err.Error())
	}

	for _, f := range res.File {
		log.Printf("%s\n", *f.Name)
		log.Printf("\n%s", *f.Content)
		log.Printf("--------------------------------------\n")
	}
}

func printRequest(path string) error {
	f, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("failed to open code generator request file: %w", err)
	}

	reqBytes, err := io.ReadAll(f)
	if err != nil {
		return fmt.Errorf("failed to read code generator request bytes: %w", err)
	}

	req := &pluginpb.CodeGeneratorRequest{}
	if err := proto.Unmarshal(reqBytes, req); err != nil {
		return fmt.Errorf("failed to unmarshal code generator response: %w", err)
	}

	log.Printf("%s", protojson.Format(req))
	log.Printf("--------------------------------------\n")
	return nil
}
