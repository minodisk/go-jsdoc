package jsdoc

import (
	"net/url"

	"github.com/lestrrat/go-jshschema"
	"github.com/lestrrat/go-jsschema"
)

type JSDoc struct {
	Title       string
	Description string
	URL         *url.URL
	Properties  map[string]*Schema
	Links       LinkList
}

type Schema struct {
	*schema.Schema
	LiteralExample string
	Properties     map[string]*Schema
	Items          *ItemSpec
	RefLink        string
}

type ItemSpec struct {
	*schema.ItemSpec
	Schemas []*Schema
}

type LinkList []*Link

type Link struct {
	hschema.Link
	Schema       *Schema
	TargetSchema *Schema
	Description  string
	Request      Request
	Response     Response
}

type Request struct {
	ContentType string
	Body        string
}

type Response struct {
	StatusCode   int
	ReasonPhrase string
	ContentType  string
	Body         string
}

type Builder struct {
	Schema *hschema.HyperSchema
}
