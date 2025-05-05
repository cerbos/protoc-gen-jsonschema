// Copyright 2021-2025 Zenauth Ltd.
// SPDX-License-Identifier: Apache-2.0

package module

import (
	"fmt"
	"strings"

	pgs "github.com/lyft/protoc-gen-star/v2"

	"github.com/cerbos/protoc-gen-jsonschema/internal/jsonschema"
)

type namedEntity interface {
	FullyQualifiedName() string
}

func (m *Module) pushMessage(message pgs.Message) {
	m.Push(fmt.Sprintf("message:%s", message.Name()))
	m.Debug("pushMessage")

	if m.nestedUnderMessage == nil {
		m.nestedUnderMessage = message
		m.definitions = make(map[string]jsonschema.Schema)
	}
}

func (m *Module) popMessage(message pgs.Message, schema jsonschema.NonTrivialSchema) {
	m.Debug("popMessage")
	if m.nestedUnder(message) {
		schema.Define(m.definitions)
		m.definitions = nil
		m.nestedUnderMessage = nil
	}

	m.Pop()
}

func (m *Module) ref(entity namedEntity, schema func() jsonschema.Schema) *jsonschema.GenericSchema {
	m.Debug("ref")
	if m.nestedUnder(entity) {
		return jsonschema.Ref("#")
	}

	key := strings.TrimPrefix(entity.FullyQualifiedName(), ".")

	if _, ok := m.definitions[key]; !ok {
		m.definitions[key] = jsonschema.True // avoid cycles
		m.definitions[key] = schema()
	}

	return jsonschema.Ref("#/definitions/" + key)
}

func (m *Module) nestedUnder(entity namedEntity) bool {
	return entity.FullyQualifiedName() == m.nestedUnderMessage.FullyQualifiedName()
}
