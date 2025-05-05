// Copyright 2021-2025 Zenauth Ltd.
// SPDX-License-Identifier: Apache-2.0

package main

import (
	pgs "github.com/lyft/protoc-gen-star/v2"

	"github.com/cerbos/protoc-gen-jsonschema/internal/common"
	"github.com/cerbos/protoc-gen-jsonschema/internal/module"
)

func main() {
	pgs.Init(pgs.DebugEnv(common.DebugEnv)).RegisterModule(module.New()).Render()
}
