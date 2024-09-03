package raml

import "encoding/json"

// Adapted from https://github.com/invopop/jsonschema/blob/main/schema.go

// Version is the JSON Schema version.
var Version = "http://json-schema.org/draft-07/schema"

// Schema represents a JSON Schema object type.
//
// https://json-schema.org/draft-07/draft-handrews-json-schema-00.pdf
type JSONSchema struct {
	Version     string      `json:"$schema,omitempty"`
	ID          string      `json:"$id,omitempty"`
	Anchor      string      `json:"$anchor,omitempty"`
	Ref         string      `json:"$ref,omitempty"`
	Definitions Definitions `json:"definitions,omitempty"`
	Comment     string      `json:"$comment,omitempty"`

	AllOf []*JSONSchema `json:"allOf,omitempty"`
	AnyOf []*JSONSchema `json:"anyOf,omitempty"`
	OneOf []*JSONSchema `json:"oneOf,omitempty"`
	Not   *JSONSchema   `json:"not,omitempty"`

	If   *JSONSchema `json:"if,omitempty"`
	Then *JSONSchema `json:"then,omitempty"`
	Else *JSONSchema `json:"else,omitempty"`

	Items *JSONSchema `json:"items,omitempty"`

	Properties           map[string]*JSONSchema `json:"properties,omitempty"`
	PatternProperties    map[string]*JSONSchema `json:"patternProperties,omitempty"`
	AdditionalProperties *bool                  `json:"additionalProperties,omitempty"`
	PropertyNames        *JSONSchema            `json:"propertyNames,omitempty"`

	Type             string      `json:"type,omitempty"`
	Enum             []any       `json:"enum,omitempty"`
	Const            any         `json:"const,omitempty"`
	MultipleOf       json.Number `json:"multipleOf,omitempty"`
	Maximum          json.Number `json:"maximum,omitempty"`
	Minimum          json.Number `json:"minimum,omitempty"`
	MaxLength        *uint64     `json:"maxLength,omitempty"`
	MinLength        *uint64     `json:"minLength,omitempty"`
	Pattern          string      `json:"pattern,omitempty"`
	MaxItems         *uint64     `json:"maxItems,omitempty"`
	MinItems         *uint64     `json:"minItems,omitempty"`
	UniqueItems      bool        `json:"uniqueItems,omitempty"`
	MaxContains      *uint64     `json:"maxContains,omitempty"`
	MinContains      *uint64     `json:"minContains,omitempty"`
	MaxProperties    *uint64     `json:"maxProperties,omitempty"`
	MinProperties    *uint64     `json:"minProperties,omitempty"`
	Required         []string    `json:"required,omitempty"`
	ContentEncoding  string      `json:"contentEncoding,omitempty"`
	ContentMediaType string      `json:"contentMediaType,omitempty"`

	Format string `json:"format,omitempty"`

	Title       string `json:"title,omitempty"`
	Description string `json:"description,omitempty"`
	Default     any    `json:"default,omitempty"`
	Examples    []any  `json:"examples,omitempty"`

	// TODO: There's no better way to serialize custom properties on the same level in Go.
	Extras map[string]any `json:"x-custom,omitempty"`

	// Special boolean representation of the Schema
	boolean *bool
}

var (
	// TrueSchema defines a schema with a true value
	TrueSchema = &JSONSchema{boolean: &[]bool{true}[0]}
	// FalseSchema defines a schema with a false value
	FalseSchema = &JSONSchema{boolean: &[]bool{false}[0]}
)

// Definitions hold schema definitions.
// http://json-schema.org/latest/json-schema-validation.html#rfc.section.5.26
// RFC draft-wright-json-schema-validation-00, section 5.26
type Definitions map[string]*JSONSchema

func (schema *JSONSchema) WithRamlData(base *BaseShape) {
	if base.DisplayName != nil {
		schema.Title = *base.DisplayName
	}
	if base.Description != nil {
		schema.Description = *base.Description
	}
	if base.Default != nil {
		schema.Default = base.Default.Value
	}
	if base.Examples != nil {
		for _, ex := range base.Examples.Map {
			schema.Examples = append(schema.Examples, ex.Data.Value)
		}
	}
	if base.Example != nil {
		schema.Examples = []any{base.Example.Data.Value}
	}
	for k, v := range base.CustomDomainProperties {
		schema.Extras["x-domainExt-"+k] = v.Extension.Value
	}
	for k, v := range base.CustomShapeFacetDefinitions {
		m := schema.Extras["x-shapeExt-definitions"]
		if m == nil {
			m = make(map[string]interface{})
			schema.Extras["x-shapeExt-definitions"] = m
		}
		shapeExtDefs := m.(map[string]interface{})
		shapeExtDefs[k] = (*v.Shape).ToJSONSchema()
	}
	for k, v := range base.CustomShapeFacets {
		schema.Extras["x-shapeExt-data-"+k] = v.Value
	}
}
