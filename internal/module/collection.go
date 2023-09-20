// Copyright 2021-2023 Zenauth Ltd.
// SPDX-License-Identifier: Apache-2.0

package module

import (
	pgs "github.com/lyft/protoc-gen-star/v2"

	"github.com/cerbos/protoc-gen-jsonschema/gen/pb/buf/validate"
	"github.com/cerbos/protoc-gen-jsonschema/internal/jsonschema"
)

func (m *Module) schemaForMap(value pgs.FieldTypeElem, rules *validate.MapRules) jsonschema.Schema {
	schema := jsonschema.NewObjectSchema()
	schema.AdditionalProperties, _ = m.schemaForElement(value, rules.GetValues())

	if rules != nil {
		if rules.GetKeys().GetString_() != nil {
			schema.PropertyNames, _ = m.schemaForString(rules.GetKeys().GetString_(), false) // TODO
		}

		if rules.MaxPairs != nil {
			schema.MaxProperties = jsonschema.Size(rules.GetMaxPairs())
		}

		if rules.MinPairs != nil {
			schema.MinProperties = jsonschema.Size(rules.GetMinPairs())
		}
	}

	return schema
}

func (m *Module) schemaForRepeated(item pgs.FieldTypeElem, rules *validate.RepeatedRules) jsonschema.Schema {
	schema := jsonschema.NewArraySchema()
	schema.Items, _ = m.schemaForElement(item, rules.GetItems())

	if rules != nil {
		if rules.MaxItems != nil {
			schema.MaxItems = jsonschema.Size(rules.GetMaxItems())
		}

		if rules.MinItems != nil {
			schema.MinItems = jsonschema.Size(rules.GetMinItems())
		}

		if rules.Unique != nil {
			schema.UniqueItems = rules.GetUnique()
		}
	}

	return schema
}

func (m *Module) schemaForElement(element pgs.FieldTypeElem, constraints *validate.FieldConstraints) (jsonschema.Schema, bool) {
	if element.IsEmbed() {
		return m.schemaForEmbed(element.Embed(), constraints)
	}

	if element.IsEnum() {
		return m.schemaForEnum(element.Enum(), constraints.GetEnum())
	}

	return m.schemaForScalar(element.ProtoType(), constraints)
}