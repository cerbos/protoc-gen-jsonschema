// Copyright 2021-2023 Zenauth Ltd.
// SPDX-License-Identifier: Apache-2.0

package module

import (
	"fmt"

	pgs "github.com/lyft/protoc-gen-star/v2"

	"github.com/cerbos/protoc-gen-jsonschema/gen/pb/buf/validate"
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

	constraints := &validate.FieldConstraints{}
	_, err := field.Extension(validate.E_Field, constraints)
	m.CheckErr(err, "unable to read validation constraints from field")

	var schema jsonschema.Schema
	var required bool
	switch {
	case field.Type().IsEmbed():
		schema, required = m.schemaForEmbed(field.Type().Embed(), constraints)
	case field.Type().IsEnum():
		schema, required = m.schemaForEnum(field.Type().Enum(), constraints.GetEnum())
	case field.Type().IsMap():
		schema = m.schemaForMap(field.Type().Element(), constraints.GetMap())
	case field.Type().IsRepeated():
		schema = m.schemaForRepeated(field.Type().Element(), constraints.GetRepeated())
		required = constraints.Required
	default:
		schema, required = m.schemaForScalar(field.Type().ProtoType(), constraints)
	}

	return schema, required && !field.InOneOf()
}

func (m *Module) schemaForEmbed(embed pgs.Message, constraints *validate.FieldConstraints) (jsonschema.Schema, bool) {
	m.Debug("schemaForEmbed")
	if embed.IsWellKnown() {
		return m.schemaForWellKnownType(embed.WellKnownType(), constraints)
	}

	return m.schemaForMessage(embed), false
}

func (m *Module) schemaForMessage(message pgs.Message) jsonschema.Schema {
	m.Debug("schemaForMessage")
	return m.messageRef(message)
}

func (m *Module) schemaForOneOf(oneOf pgs.OneOf) jsonschema.NonTrivialSchema {
	m.Debug("schemaForOneOf")
	constraint := validate.OneofConstraints{}
	_, err := oneOf.Extension(validate.E_Oneof, &constraint)
	m.CheckErr(err, "unable to read oneOf option")

	if constraint.Required == nil || !*constraint.Required {
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
