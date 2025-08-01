// Copyright 2021-2025 Zenauth Ltd.
// SPDX-License-Identifier: Apache-2.0

package module

import (
	"fmt"

	"buf.build/gen/go/bufbuild/protovalidate/protocolbuffers/go/buf/validate"
	pgs "github.com/lyft/protoc-gen-star/v2"

	"github.com/cerbos/protoc-gen-jsonschema/internal/jsonschema"
)

func (m *Module) defineMessage(message pgs.Message) jsonschema.NonTrivialSchema {
	m.pushMessage(message)
	m.Debug("defineMessage")

	schema := jsonschema.NewObjectSchema()
	schema.AdditionalProperties = jsonschema.False
	schemas := []jsonschema.NonTrivialSchema{schema}

	for _, field := range message.Fields() {
		name := m.propertyName(field)
		valueSchema, required := m.schemaForField(field)
		schema.Properties[name] = valueSchema
		if required {
			schema.Required = append(schema.Required, name)
		}
	}

	for _, oneOf := range message.OneOfs() {
		oneOfSchema := m.schemaForOneOf(oneOf)
		if oneOfSchema != nil {
			schemas = append(schemas, oneOfSchema)
		}
	}

	result := jsonschema.AllOf(schemas...)
	m.popMessage(message, result)
	return result
}

func (m *Module) propertyName(field pgs.Field) string {
	return field.Descriptor().GetJsonName()
}

func (m *Module) schemaForField(field pgs.Field) (jsonschema.Schema, bool) {
	m.Push(fmt.Sprintf("field:%s", field.Name()))
	defer m.Pop()
	m.Debug("schemaForField")

	rules := &validate.FieldRules{}
	_, err := field.Extension(validate.E_Field, rules)
	m.CheckErr(err, "unable to read validation rules from field")

	required := rules.GetRequired()
	if rules.GetIgnore() == validate.Ignore_IGNORE_IF_ZERO_VALUE {
		required = false
	}

	if field.HasOptionalKeyword() {
		required = false
	}

	var schema jsonschema.Schema
	switch {
	case field.Type().IsEmbed():
		schema = m.schemaForEmbed(field.Type().Embed(), rules)
	case field.Type().IsEnum():
		schema = m.schemaForEnum(field.Type().Enum(), rules.GetEnum())
	case field.Type().IsMap():
		schema = m.schemaForMap(field.Type().Element(), rules.GetMap())
	case field.Type().IsRepeated():
		schema = m.schemaForRepeated(field.Type().Element(), rules.GetRepeated())
	default:
		schema = m.schemaForScalar(field.Type().ProtoType(), rules)
	}

	return schema, required && !field.InOneOf()
}

func (m *Module) schemaForEmbed(embed pgs.Message, rules *validate.FieldRules) jsonschema.Schema {
	m.Debug("schemaForEmbed")
	if embed.IsWellKnown() {
		return m.schemaForWellKnownType(embed.WellKnownType(), rules)
	}

	return m.schemaForMessage(embed)
}

func (m *Module) schemaForMessage(message pgs.Message) jsonschema.Schema {
	m.Debug("schemaForMessage")
	return m.messageRef(message)
}

func (m *Module) schemaForOneOf(oneOf pgs.OneOf) jsonschema.NonTrivialSchema {
	m.Debug("schemaForOneOf")
	rules := validate.OneofRules{}
	_, err := oneOf.Extension(validate.E_Oneof, &rules)
	m.CheckErr(err, "unable to read oneOf option")

	if rules.Required == nil || !*rules.Required {
		return nil
	}

	schemas := make([]jsonschema.NonTrivialSchema, len(oneOf.Fields()))
	for i, field := range oneOf.Fields() {
		schema := jsonschema.NewObjectSchema()
		schema.Required = []string{m.propertyName(field)}
		schemas[i] = schema
	}

	return jsonschema.OneOf(schemas...)
}

func (m *Module) messageRef(message pgs.Message) jsonschema.Schema {
	m.Debug("messageRef")
	return m.ref(message, func() jsonschema.Schema {
		return m.defineMessage(message)
	})
}
