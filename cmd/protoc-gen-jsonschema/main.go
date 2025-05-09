// Copyright 2021-2025 Zenauth Ltd.
// SPDX-License-Identifier: Apache-2.0

package main

import (
	pgs "github.com/lyft/protoc-gen-star/v2"
	"google.golang.org/protobuf/types/pluginpb"

	"github.com/cerbos/protoc-gen-jsonschema/internal/common"
	"github.com/cerbos/protoc-gen-jsonschema/internal/module"
)

func main() {
	supportedFeatures := uint64(pluginpb.CodeGeneratorResponse_FEATURE_PROTO3_OPTIONAL | pluginpb.CodeGeneratorResponse_FEATURE_SUPPORTS_EDITIONS)
	pgs.Init(pgs.SupportedFeatures(&supportedFeatures), pgs.DebugEnv(common.DebugEnv)).RegisterModule(module.New()).Render()
}
