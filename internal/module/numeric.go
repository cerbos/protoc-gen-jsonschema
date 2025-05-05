// Copyright 2021-2025 Zenauth Ltd.
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

func (m *Module) schemaForNumericScalar(numeric pgs.ProtoType, rules *validate.FieldRules) jsonschema.Schema {
	m.Debug("schemaForNumericScalar")
	value := m.valueSchemaForNumericScalar(numeric)
	schemas := []jsonschema.NonTrivialSchema{m.stringValueSchemaForNumericScalar(numeric, value)}
	r := m.numericRules(numeric, rules)

	//nolint:nestif
	if r != nil {
		if r.Const != nil {
			value.Const = r.Const
		}

		if r.GreaterThan.Gt != nil {
			value.ExclusiveMinimum = r.GreaterThan.Gt
		}

		if r.GreaterThan.Gte != nil {
			value.Minimum = r.GreaterThan.Gte
		}

		if len(r.In) > 0 {
			value.Enum = r.In
		}

		if r.LessThan.Lt != nil {
			value.ExclusiveMaximum = r.LessThan.Lt
		}

		if r.LessThan.Lte != nil {
			value.Maximum = r.LessThan.Lte
		}

		if len(r.NotIn) > 0 {
			in := jsonschema.NewNumberSchema()
			in.Enum = r.NotIn

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

func (m *Module) numericRules(numeric pgs.ProtoType, rules *validate.FieldRules) *numericRules {
	m.Debug("numericRules")
	var source proto.Message

	switch numeric {
	case pgs.DoubleT:
		source = rules.GetDouble()

	case pgs.Fixed32T:
		source = rules.GetFixed32()

	case pgs.Fixed64T:
		source = rules.GetFixed64()

	case pgs.FloatT:
		source = rules.GetFloat()

	case pgs.Int32T:
		source = rules.GetInt32()

	case pgs.Int64T:
		source = rules.GetInt64()

	case pgs.SFixed32:
		source = rules.GetSfixed32()

	case pgs.SFixed64:
		source = rules.GetSfixed64()

	case pgs.SInt32:
		source = rules.GetSint32()

	case pgs.SInt64:
		source = rules.GetSint64()

	case pgs.StringT:
		source = rules.GetString()

	case pgs.UInt32T:
		source = rules.GetUint32()

	case pgs.UInt64T:
		source = rules.GetUint64()

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
