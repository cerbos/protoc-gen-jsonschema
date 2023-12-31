// Copyright 2021-2023 Zenauth Ltd.
// SPDX-License-Identifier: Apache-2.0

package module

import (
	"encoding/json"

	"buf.build/gen/go/bufbuild/protovalidate/protocolbuffers/go/buf/validate"
	pgs "github.com/lyft/protoc-gen-star/v2"
	"google.golang.org/protobuf/proto"

	"github.com/cerbos/protoc-gen-jsonschema/internal/jsonschema"
)

const (
	signedDecimalString   = `^-?(?:0|[1-9]\d*)(?:\.\d+)?(?:[eE][+-]?\d+)?$`
	unsignedDecimalString = `^(?:0|[1-9]\d*)(?:\.\d+)?(?:[eE][+-]?\d+)?$`
)

//nolint:tagliatelle
type numericRules struct {
	Const       jsonschema.Number `json:"const,omitempty"`
	GreaterThan struct {
		Gt  jsonschema.Number `json:"Gt,omitempty"`
		Gte jsonschema.Number `json:"Gte,omitempty"`
	} `json:"GreaterThan,omitempty"`
	LessThan struct {
		Lt  jsonschema.Number `json:"Lt,omitempty"`
		Lte jsonschema.Number `json:"Lte,omitempty"`
	} `json:"LessThan,omitempty"`
	In    []jsonschema.Number `json:"in,omitempty"`
	NotIn []jsonschema.Number `json:"not_in,omitempty"`
}

func (m *Module) schemaForNumericScalar(numeric pgs.ProtoType, constraints *validate.FieldConstraints) jsonschema.Schema {
	m.Debug("schemaForNumericScalar")
	value := m.valueSchemaForNumericScalar(numeric)
	schemas := []jsonschema.NonTrivialSchema{m.stringValueSchemaForNumericScalar(numeric, value)}
	rules := m.numericRules(numeric, constraints)

	//nolint:nestif
	if rules != nil {
		if rules.Const != nil {
			value.Const = rules.Const
		}

		if rules.GreaterThan.Gt != nil {
			value.ExclusiveMinimum = rules.GreaterThan.Gt
		}

		if rules.GreaterThan.Gte != nil {
			value.Minimum = rules.GreaterThan.Gte
		}

		if len(rules.In) > 0 {
			value.Enum = rules.In
		}

		if rules.LessThan.Lt != nil {
			value.ExclusiveMaximum = rules.LessThan.Lt
		}

		if rules.LessThan.Lte != nil {
			value.Maximum = rules.LessThan.Lte
		}

		if len(rules.NotIn) > 0 {
			in := jsonschema.NewNumberSchema()
			in.Enum = rules.NotIn

			schemas = append(schemas, jsonschema.Not(in))
		}
	}

	return jsonschema.AllOf(schemas...)
}

func (m *Module) valueSchemaForNumericScalar(numeric pgs.ProtoType) *jsonschema.NumberSchema {
	m.Debug("valueSchemaForNumericScalar")
	switch numeric {
	case pgs.Fixed32T, pgs.UInt32T, pgs.Fixed64T, pgs.UInt64T:
		schema := jsonschema.NewIntegerSchema()
		schema.Minimum = jsonschema.Number("0")
		return schema

	case pgs.Int32T, pgs.SFixed32, pgs.SInt32, pgs.Int64T, pgs.SFixed64, pgs.SInt64:
		return jsonschema.NewIntegerSchema()

	case pgs.DoubleT, pgs.FloatT:
		return jsonschema.NewNumberSchema()

	default:
		m.Failf("unknown numeric scalar type %q", numeric)
		return nil
	}
}

func (m *Module) stringValueSchemaForNumericScalar(numeric pgs.ProtoType, value *jsonschema.NumberSchema) jsonschema.NonTrivialSchema {
	m.Debug("stringValueSchemaForNumericScalar")
	var pattern string

	switch numeric {
	case pgs.Fixed64T, pgs.UInt64T:
		pattern = unsignedDecimalString

	case pgs.Int64T, pgs.SFixed64, pgs.SInt64:
		pattern = signedDecimalString

	default:
		return value
	}

	stringValue := jsonschema.NewStringSchema()
	stringValue.Pattern = pattern

	return jsonschema.OneOf(value, stringValue)
}

func (m *Module) numericRules(numeric pgs.ProtoType, constraints *validate.FieldConstraints) *numericRules {
	m.Debug("numericRules")
	var source proto.Message

	switch numeric {
	case pgs.DoubleT:
		source = constraints.GetDouble()

	case pgs.Fixed32T:
		source = constraints.GetFixed32()

	case pgs.Fixed64T:
		source = constraints.GetFixed64()

	case pgs.FloatT:
		source = constraints.GetFloat()

	case pgs.Int32T:
		source = constraints.GetInt32()

	case pgs.Int64T:
		source = constraints.GetInt64()

	case pgs.SFixed32:
		source = constraints.GetSfixed32()

	case pgs.SFixed64:
		source = constraints.GetSfixed64()

	case pgs.SInt32:
		source = constraints.GetSint32()

	case pgs.SInt64:
		source = constraints.GetSint64()

	case pgs.StringT:
		source = constraints.GetString_()

	case pgs.UInt32T:
		source = constraints.GetUint32()

	case pgs.UInt64T:
		source = constraints.GetUint64()

	default:
		m.Failf("unknown numeric scalar type %q", numeric)
		return nil
	}

	if source == nil {
		return nil
	}

	data, err := json.Marshal(source)
	m.CheckErr(err, "failed to marshal numeric validation rules to JSON")

	target := &numericRules{}
	err = json.Unmarshal(data, target)
	m.CheckErr(err, "failed to unmarshal numeric validation rules from JSON")

	return target
}
