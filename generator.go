package jsdoc

import (
	"io"
	"io/ioutil"
	"strings"
	"text/template"

	"github.com/lestrrat/go-jsschema"
)

type Generator struct{}

func NewGenerator() *Generator {
	return &Generator{}
}

func (g *Generator) Process(out io.Writer, doc *JSDoc, tmpl io.Reader) error {
	b, err := ioutil.ReadAll(tmpl)
	if err != nil {
		return err
	}
	t := template.Must(template.New("doc").Funcs(map[string]interface{}{
		"joinTypes": func(ts schema.PrimitiveTypes, sep string) string {
			var strs []string
			for _, t := range ts {
				strs = append(strs, t.String())
			}
			return strings.Join(strs, sep)
		},
	}).Parse(string(b)))
	if err := t.Execute(out, doc); err != nil {
		return err
	}
	return nil
}
