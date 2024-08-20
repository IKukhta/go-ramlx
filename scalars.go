package raml

import (
	"gopkg.in/yaml.v3"
)

type EnumFacets struct {
	Enum []*Node
}

func MakeEnum(v *yaml.Node, location string) ([]*Node, error) {
	if v.Kind != yaml.SequenceNode {
		return nil, NewError("enum must be sequence node", location, WithNodePosition(v))
	}
	var enums []*Node = make([]*Node, len(v.Content))
	for i, v := range v.Content {
		n, err := MakeNode(v, location)
		if err != nil {
			return nil, NewWrappedError("make node enum", err, location, WithNodePosition(v))
		}
		enums[i] = n
	}
	return enums, nil
}

type FormatFacets struct {
	Format *string
}

type IntegerFacets struct {
	Minimum *any
	Maximum *any
}

type IntegerShape struct {
	BaseShape

	EnumFacets
	FormatFacets
	IntegerFacets
}

func (s *IntegerShape) Base() *BaseShape {
	return &s.BaseShape
}

func (s *IntegerShape) Clone() Shape {
	c := *s
	return &c
}

// func (s *IntegerShape) Validate(v interface{}) error {
// 	i, ok := v.(int64)
// 	if !ok {
// 		return fmt.Errorf("invalid value")
// 	}

// 	if t.Minimum != nil && *t.Minimum < i {
// 		return fmt.Errorf("value must be in range")
// 	}
// 	if t.Maximum != nil && i > *t.Maximum {
// 		return fmt.Errorf("value must be in range")
// 	}

// 	return nil
// }

func (s *IntegerShape) UnmarshalYAMLNodes(v []*yaml.Node) error {
	for i := 0; i != len(v); i += 2 {
		node := v[i]
		valueNode := v[i+1]
		if node.Value == "minimum" {
			if err := valueNode.Decode(&s.Minimum); err != nil {
				return NewWrappedError("decode minimum", err, s.Location, WithNodePosition(valueNode))
			}
		} else if node.Value == "maximum" {
			if err := valueNode.Decode(&s.Maximum); err != nil {
				return NewWrappedError("decode maximum", err, s.Location, WithNodePosition(valueNode))
			}
		} else if node.Value == "format" {
			if err := valueNode.Decode(&s.Format); err != nil {
				return NewWrappedError("decode format", err, s.Location, WithNodePosition(valueNode))
			}
		} else if node.Value == "enum" {
			enums, err := MakeEnum(valueNode, s.Location)
			if err != nil {
				return NewWrappedError("make enum", err, s.Location, WithNodePosition(valueNode))
			}
			s.Enum = enums
		} else {
			n, err := MakeNode(valueNode, s.Location)
			if err != nil {
				return NewWrappedError("make node", err, s.Location, WithNodePosition(valueNode))
			}
			s.CustomShapeFacets[node.Value] = n
		}
	}
	return nil
}

type NumberFacets struct {
	IntegerFacets
	MultipleOf *float64
}

type NumberShape struct {
	BaseShape

	EnumFacets
	FormatFacets
	NumberFacets
}

func (s *NumberShape) Base() *BaseShape {
	return &s.BaseShape
}

func (s *NumberShape) Clone() Shape {
	c := *s
	return &c
}

func (s *NumberShape) UnmarshalYAMLNodes(v []*yaml.Node) error {
	for i := 0; i != len(v); i += 2 {
		node := v[i]
		valueNode := v[i+1]
		if node.Value == "minimum" {
			if err := valueNode.Decode(&s.Minimum); err != nil {
				return NewWrappedError("decode minimum", err, s.Location, WithNodePosition(valueNode))
			}
		} else if node.Value == "maximum" {
			if err := valueNode.Decode(&s.Maximum); err != nil {
				return NewWrappedError("decode maximum", err, s.Location, WithNodePosition(valueNode))
			}
		} else if node.Value == "format" {
			if err := valueNode.Decode(&s.Format); err != nil {
				return NewWrappedError("decode format", err, s.Location, WithNodePosition(valueNode))
			}
		} else if node.Value == "enum" {
			enums, err := MakeEnum(valueNode, s.Location)
			if err != nil {
				return NewWrappedError("make enum", err, s.Location, WithNodePosition(valueNode))
			}
			s.Enum = enums
		} else if node.Value == "multipleOf" {
			if err := valueNode.Decode(&s.MultipleOf); err != nil {
				return NewWrappedError("decode multipleOf", err, s.Location, WithNodePosition(valueNode))
			}
		} else {
			n, err := MakeNode(valueNode, s.Location)
			if err != nil {
				return NewWrappedError("make node", err, s.Location, WithNodePosition(valueNode))
			}
			s.CustomShapeFacets[node.Value] = n
		}
	}
	return nil
}

type ByteLengthFacets struct {
	MaxLength *int64
	MinLength *int64
}

type StringFacets struct {
	ByteLengthFacets
	Pattern *string
}

type StringShape struct {
	BaseShape

	EnumFacets
	StringFacets
}

func (s *StringShape) Base() *BaseShape {
	return &s.BaseShape
}

func (s *StringShape) Clone() Shape {
	c := *s
	return &c
}

func (s *StringShape) UnmarshalYAMLNodes(v []*yaml.Node) error {
	for i := 0; i != len(v); i += 2 {
		node := v[i]
		valueNode := v[i+1]

		if node.Value == "minLength" {
			if err := valueNode.Decode(&s.MinLength); err != nil {
				return NewWrappedError("decode minLength", err, s.Location, WithNodePosition(valueNode))
			}
		} else if node.Value == "maxLength" {
			if err := valueNode.Decode(&s.MaxLength); err != nil {
				return NewWrappedError("decode maxLength", err, s.Location, WithNodePosition(valueNode))
			}
		} else if node.Value == "pattern" {
			if err := valueNode.Decode(&s.Pattern); err != nil {
				return NewWrappedError("decode pattern", err, s.Location, WithNodePosition(valueNode))
			}
		} else if node.Value == "enum" {
			enums, err := MakeEnum(valueNode, s.Location)
			if err != nil {
				return NewWrappedError("make enum", err, s.Location, WithNodePosition(valueNode))
			}
			s.Enum = enums
		} else {
			n, err := MakeNode(valueNode, s.Location)
			if err != nil {
				return NewWrappedError("make node", err, s.Location, WithNodePosition(valueNode))
			}
			s.CustomShapeFacets[node.Value] = n
		}
	}
	return nil
}

type FileFacets struct {
	FileTypes []Node
}

type FileShape struct {
	BaseShape

	ByteLengthFacets
	FileFacets
}

func (s *FileShape) Base() *BaseShape {
	return &s.BaseShape
}

func (s *FileShape) Clone() Shape {
	c := *s
	return &c
}

func (s *FileShape) UnmarshalYAMLNodes(v []*yaml.Node) error {
	for i := 0; i != len(v); i += 2 {
		node := v[i]
		valueNode := v[i+1]

		if node.Value == "minLength" {
			if err := valueNode.Decode(&s.MinLength); err != nil {
				return NewWrappedError("decode minLength", err, s.Location, WithNodePosition(valueNode))
			}
		} else if node.Value == "maxLength" {
			if err := valueNode.Decode(&s.MaxLength); err != nil {
				return NewWrappedError("decode maxLength", err, s.Location, WithNodePosition(valueNode))
			}
		} else if node.Value == "fileTypes" {
			if err := valueNode.Decode(&s.FileTypes); err != nil {
				return NewWrappedError("decode fileTypes", err, s.Location, WithNodePosition(valueNode))
			}
		} else {
			n, err := MakeNode(valueNode, s.Location)
			if err != nil {
				return NewWrappedError("make node", err, s.Location, WithNodePosition(valueNode))
			}
			s.CustomShapeFacets[node.Value] = n
		}
	}
	return nil
}

type BooleanShape struct {
	BaseShape

	EnumFacets
}

func (s *BooleanShape) Clone() Shape {
	c := *s
	return &c
}

func (s *BooleanShape) Base() *BaseShape {
	return &s.BaseShape
}

func (s *BooleanShape) UnmarshalYAMLNodes(v []*yaml.Node) error {
	for i := 0; i != len(v); i += 2 {
		node := v[i]
		valueNode := v[i+1]

		if node.Value == "enum" {
			enums, err := MakeEnum(valueNode, s.Location)
			if err != nil {
				return NewWrappedError("make enum", err, s.Location, WithNodePosition(valueNode))
			}
			s.Enum = enums
		} else {
			n, err := MakeNode(valueNode, s.Location)
			if err != nil {
				return NewWrappedError("make node", err, s.Location, WithNodePosition(valueNode))
			}
			s.CustomShapeFacets[node.Value] = n
		}
	}

	return nil
}

type DateTimeShape struct {
	BaseShape

	FormatFacets
}

func (s *DateTimeShape) Base() *BaseShape {
	return &s.BaseShape
}

func (s *DateTimeShape) Clone() Shape {
	c := *s
	return &c
}

func (s *DateTimeShape) UnmarshalYAMLNodes(v []*yaml.Node) error {
	for i := 0; i != len(v); i += 2 {
		node := v[i]
		valueNode := v[i+1]
		if node.Value == "format" {
			if err := valueNode.Decode(&s.Format); err != nil {
				return NewWrappedError("decode format", err, s.Location, WithNodePosition(valueNode))
			}
		} else {
			n, err := MakeNode(valueNode, s.Location)
			if err != nil {
				return NewWrappedError("make node", err, s.Location, WithNodePosition(valueNode))
			}
			s.CustomShapeFacets[node.Value] = n
		}
	}
	return nil
}

type DateTimeOnlyShape struct {
	BaseShape
}

func (s *DateTimeOnlyShape) Base() *BaseShape {
	return &s.BaseShape
}

func (s *DateTimeOnlyShape) Clone() Shape {
	c := *s
	return &c
}

func (s *DateTimeOnlyShape) UnmarshalYAMLNodes(v []*yaml.Node) error {
	return nil
}

type DateOnlyShape struct {
	BaseShape
}

func (s *DateOnlyShape) Base() *BaseShape {
	return &s.BaseShape
}

func (s *DateOnlyShape) Clone() Shape {
	c := *s
	return &c
}

func (s *DateOnlyShape) UnmarshalYAMLNodes(v []*yaml.Node) error {
	return nil
}

type TimeOnlyShape struct {
	BaseShape
}

func (s *TimeOnlyShape) Base() *BaseShape {
	return &s.BaseShape
}

func (s *TimeOnlyShape) Clone() Shape {
	c := *s
	return &c
}

func (s *TimeOnlyShape) UnmarshalYAMLNodes(v []*yaml.Node) error {
	return nil
}

type AnyShape struct {
	BaseShape
}

func (s *AnyShape) Base() *BaseShape {
	return &s.BaseShape
}

func (s *AnyShape) Clone() Shape {
	c := *s
	return &c
}

func (s *AnyShape) UnmarshalYAMLNodes(v []*yaml.Node) error {
	return nil
}

type NilShape struct {
	BaseShape
}

func (s *NilShape) Base() *BaseShape {
	return &s.BaseShape
}

func (s *NilShape) Clone() Shape {
	c := *s
	return &c
}

func (s *NilShape) UnmarshalYAMLNodes(v []*yaml.Node) error {
	return nil
}
