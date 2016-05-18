package jsdoc

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strings"

	"github.com/lestrrat/go-jshschema"
	"github.com/lestrrat/go-jsschema"
	"github.com/lestrrat/go-pdebug"
)

var (
	statusCodes = map[string]int{
		"GET":    200,
		"POST":   201,
		"PUT":    204,
		"DELETE": 204,
	}
)

func NewBuilder() *Builder {
	return &Builder{}
}

func (b *Builder) Build(hs *hschema.HyperSchema) (d *JSDoc, err error) {
	if pdebug.Enabled {
		g := pdebug.IPrintf("START Builder.Build")
		defer func() {
			if err == nil {
				g.IRelease("END Builder.Build (OK)")
			} else {
				g.IRelease("END Builder.Build (FAIL): %s", err)
			}
		}()
	}

	d = New(len(hs.Links))
	d.Title = hs.Schema.Title
	d.Description = hs.Schema.Description
	str := getString(hs.Schema.Extras, "href")
	d.URL, err = url.Parse(str)
	if err != nil {
		return nil, err
	}

	d.Properties = map[string]*Schema{}
	for n, prop := range hs.Properties {
		s, err := resolve(prop, hs.Schema)
		if err != nil {
			return nil, err
		}
		d.Properties[n] = s
	}

	for i, l := range hs.Links {
		var (
			s, ts *Schema
		)

		req := Request{}
		if l.Schema != nil {
			s, err := resolve(l.Schema, hs.Schema)
			if err != nil {
				return nil, err
			}
			req.Body, err = encodeExample(s)
			if err != nil {
				return nil, err
			}
			if req.Body != "" {
				req.ContentType = "application/json"
			}
		}

		res := Response{
			StatusCode: statusCodes[l.Method],
		}
		if l.TargetSchema != nil {
			ts, err := resolve(l.TargetSchema, hs.Schema)
			if err != nil {
				return nil, err
			}
			res.Body, err = encodeExample(ts)
			if err != nil {
				return nil, err
			}
			if res.Body != "" {
				res.ContentType = "application/json"
			}
		}

		d.Links[i] = &Link{
			Link:         *l,
			Description:  getString(l.Extras, "description"),
			Schema:       s,
			TargetSchema: ts,
			Request:      req,
			Response:     res,
		}
	}

	return
}

func resolve(src, ctx *schema.Schema) (*Schema, error) {
	rs, err := src.Resolve(ctx)
	if err != nil {
		return nil, err
	}

	e := rs.Extras["example"]
	le := ""
	if e != nil {
		b, err := json.Marshal(e)
		if err == nil {
			le = string(b)
		}
	}
	dest := &Schema{
		Schema:         rs,
		LiteralExample: le,
		Properties:     map[string]*Schema{},
		RefLink:        fmt.Sprintf("#%s", strings.TrimLeft(src.Reference, "#/definitions/")),
	}

	if rs.Items != nil {
		dest.Items = &ItemSpec{
			ItemSpec: rs.Items,
			Schemas:  make([]*Schema, len(rs.Items.Schemas)),
		}
		for i, prop := range rs.Items.Schemas {
			p, err := resolve(prop, ctx)
			if err != nil {
				return nil, err
			}
			dest.Items.Schemas[i] = p
		}
	}

	for name, prop := range rs.Properties {
		if prop != nil {
			p, err := resolve(prop, ctx)
			if err != nil {
				return nil, err
			}
			dest.Properties[name] = p
		}
	}

	return dest, nil
}

func getString(extras map[string]interface{}, key string) string {
	v, _ := extras[key].(string)
	return v
}

func encodeExample(s *Schema) (string, error) {
	i, err := encodeProperty(s)
	if err != nil {
		return "", err
	}

	b, err := json.MarshalIndent(i, "", "  ")
	if err != nil {
		return "", err
	}

	return string(b), nil
}

func encodeItems(items *ItemSpec) ([]interface{}, error) {
	if items == nil || len(items.Schemas) == 0 {
		return []interface{}{}, nil
	}
	prop, err := encodeProperty(items.Schemas[0])
	if err != nil {
		return nil, err
	}
	return []interface{}{prop}, nil
}

func encodeProperties(props map[string]*Schema) (map[string]interface{}, error) {
	obj := map[string]interface{}{}
	for name, prop := range props {
		e, err := encodeProperty(prop)
		if err != nil {
			return nil, err
		}
		obj[name] = e
	}
	return obj, nil
}

func encodeProperty(prop *Schema) (interface{}, error) {
	if prop.Type.Contains(schema.ArrayType) {
		return encodeItems(prop.Items)
	}
	if prop.Type.Contains(schema.ObjectType) {
		return encodeProperties(prop.Properties)
	}

	i := prop.Extras["example"]
	if i != nil {
		return i, nil
	}

	if prop.Type.Contains(schema.NumberType) || prop.Type.Contains(schema.IntegerType) {
		return 0, nil
	}
	if prop.Type.Contains(schema.BooleanType) {
		return false, nil
	}
	if prop.Type.Contains(schema.StringType) {
		return "", nil
	}
	if prop.Type.Contains(schema.NullType) {
		return nil, nil
	}
	return nil, fmt.Errorf("no example for %s", prop.ID)
}
