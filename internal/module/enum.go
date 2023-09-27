// Copyright 2021-2023 Zenauth Ltd.
// SPDX-License-Identifier: Apache-2.0

package module

import (
	pgs "github.com/lyft/protoc-gen-star/v2"

	"github.com/cerbos/protoc-gen-jsonschema/gen/pb/buf/validate"
	"github.com/cerbos/protoc-gen-jsonschema/internal/jsonschema"
)

func (m *Module) defineEnum(enum pgs.Enum) *jsonschema.StringSchema {
	schema := jsonschema.NewStringSchema()

	for _, value := range enum.Values() {
		schema.Enum = append(schema.Enum, value.Name().String())
	}

	return schema
}

func (m *Module) schemaForEnum(enum pgs.Enum, rules *validate.EnumRules) jsonschema.Schema {
	m.Debug("schemaForEnum")
	if rules != nil {
		switch {
		case rules.Const != nil:
			return m.schemaForEnumConst(enum, rules.GetConst())
		case len(rules.In) > 0:
			return m.schemaForEnumIn(enum, rules.In)
		case len(rules.NotIn) > 0:
			return m.schemaForEnumNotIn(enum, rules.NotIn)
		}
	}

	return m.enumRef(enum)
}

func (m *Module) schemaForEnumConst(enum pgs.Enum, value int32) *jsonschema.StringSchema {
	m.Debug("schemaForEnumConst")
	schema := jsonschema.NewStringSchema()
	schema.Const = jsonschema.String(m.lookUpEnumName(enum, value))

	return schema
}

func (m *Module) schemaForEnumIn(enum pgs.Enum, values []int32) *jsonschema.StringSchema {
	m.Debug("schemaForEnumIn")
	schema := jsonschema.NewStringSchema()
	for _, value := range values {
		schema.Enum = append(schema.Enum, m.lookUpEnumName(enum, value))
	}

	return schema
}

func (m *Module) schemaForEnumNotIn(enum pgs.Enum, values []int32) *jsonschema.StringSchema {
	m.Debug("schemaForEnumNotIn")
	exclude := make(map[int32]struct{}, len(values))
	for _, v := range values {
		exclude[v] = struct{}{}
	}

	schema := jsonschema.NewStringSchema()
	for _, v := range enum.Values() {
		value := v.Value()
		if _, ok := exclude[value]; !ok {
			schema.Enum = append(schema.Enum, v.Name().String())
		}
	}

	return schema
}

func (m *Module) lookUpEnumName(enum pgs.Enum, value int32) string {
	m.Debug("lookUpEnumName")
	for _, enumValue := range enum.Values() {
		if enumValue.Value() == value {
			return enumValue.Name().String()
		}
	}

	m.Failf("unknown enum value %d", value)
	return ""
}

func (m *Module) enumRef(enum pgs.Enum) *jsonschema.GenericSchema {
	m.Debug("enumRef")
	return m.ref(enum, func() jsonschema.Schema {
		return m.defineEnum(enum)
	})
}
