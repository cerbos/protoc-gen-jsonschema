// Copyright 2021-2023 Zenauth Ltd.
// SPDX-License-Identifier: Apache-2.0

//go:build toolsx
// +build toolsx

package tools

import (
	_ "github.com/bufbuild/buf/cmd/buf"
	_ "github.com/lyft/protoc-gen-star"
	_ "gotest.tools/gotestsum"
)
