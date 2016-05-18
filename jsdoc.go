package jsdoc

import "github.com/lestrrat/go-jsschema"

func New(length int) *JSDoc {
	return &JSDoc{
		Links: make(LinkList, length),
	}
}

func (d *JSDoc) Host() string {
	if d.URL == nil {
		return "example.com"
	}
	return d.URL.Host
}

func (s *Schema) IsArray() bool {
	return s.Type.Contains(schema.ArrayType)
}

func (s *Schema) IsObject() bool {
	return s.Type.Contains(schema.ObjectType)
}
