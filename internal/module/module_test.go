// Copyright 2021-2023 Zenauth Ltd.
// SPDX-License-Identifier: Apache-2.0

package module_test

import (
	"bytes"
	"os"
	"testing"

	pgs "github.com/lyft/protoc-gen-star/v2"
	"github.com/stretchr/testify/require"

	"github.com/cerbos/protoc-gen-jsonschema/internal/common"
	"github.com/cerbos/protoc-gen-jsonschema/internal/module"
	"github.com/cerbos/protoc-gen-jsonschema/internal/test"
)

const requestName = "code_generator_request.pb.bin"

func TestModule(t *testing.T) {
	reqFile, err := os.Open(test.PathToDir(t, requestName))
	require.NoError(t, err)
	resBytes := &bytes.Buffer{}
	pgs.Init(
		pgs.DebugEnv(common.DebugEnv),
		pgs.ProtocInput(reqFile),
		pgs.ProtocOutput(resBytes),
	).RegisterModule(module.New()).Render()
}
