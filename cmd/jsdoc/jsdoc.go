package main

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/ghodss/yaml"
	"github.com/jessevdk/go-flags"
	"github.com/lestrrat/go-jshschema"
	"github.com/minodisk/go-jsdoc"
)

func main() {
	os.Exit(_main())
}

type options struct {
	Schema   string   `short:"s" long:"schema" description:"the source JSON Schema file"`
	OutFile  string   `short:"o" long:"outfile" description:"output file to generate"`
	Pointer  []string `short:"p" long:"ptr" description:"JSON pointer(s) within the document to create validators with"`
	Template string   `short:"t" log:"tmpl" description:"template file to generate document"`
}

func _main() int {
	var opts options
	if _, err := flags.Parse(&opts); err != nil {
		log.Printf("fail to parse flags: %s", err)
		return 1
	}

	f, err := os.Open(opts.Schema)
	if err != nil {
		log.Printf("fail to open the source JSON Schema file: %s", err)
		return 1
	}
	defer f.Close()

	var m map[string]interface{}
	switch ext := filepath.Ext(opts.Schema); ext {
	case ".json":
		if err := json.NewDecoder(f).Decode(&m); err != nil {
			log.Printf("fail to decode JSON: %s", err)
			return 1
		}
	case ".yml", ".yaml":
		b, err := ioutil.ReadFile(opts.Schema)
		if err != nil {
			log.Printf("fail to read the source JSON Schema file: %s", err)
			return 1
		}
		if err := yaml.Unmarshal(b, &m); err != nil {
			log.Printf("fail to unmarshal YAML: %s", err)
			return 1
		}
	default:
		log.Printf("undefined extension: %s", ext)
		return 1
	}

	s := hschema.New()
	if err := s.Extract(m); err != nil {
		log.Printf("fail to extract JSON Hyper Schema: %s", err)
		return 1
	}

	b := jsdoc.NewBuilder()
	d, err := b.Build(s)
	if err != nil {
		log.Printf("fail to build JSDoc: %s", err)
		return 1
	}

	var out io.Writer

	out = os.Stdout
	if fn := opts.OutFile; fn != "" {
		f, err := os.Create(fn)
		if err != nil {
			log.Printf("%s", err)
			return 1
		}
		defer f.Close()

		out = f
	}

	t, err := os.Open(opts.Template)
	if err != nil {
		log.Printf("fail to open the template file: %s", err)
		return 1
	}
	defer f.Close()

	g := jsdoc.NewGenerator()
	if err := g.Process(out, d, t); err != nil {
		log.Printf("%s", err)
		return 1
	}

	// var schemas []*schema.Schema
	// ptrs := opts.Pointer
	// if len(ptrs) == 0 {
	// 	s, err := schema.ReadFile(opts.Schema)
	// 	if err != nil {
	// 		log.Printf("fail to parse schema: %s", err)
	// 		return 1
	// 	}
	// 	schemas = []*schema.Schema{s}
	// } else {
	// 	for _, ptr := range ptrs {
	// 		log.Printf("Resolving pointer '%s'", ptr)
	// 		resolver, err := jspointer.New(ptr)
	// 		if err != nil {
	// 			log.Println("fail to resolve pointer: %s", err)
	// 			return 1
	// 		}
	//
	// 		resolved, err := resolver.Get(m)
	// 		if err != nil {
	// 			log.Println("%s", err)
	// 			return 1
	// 		}
	//
	// 		m2, ok := resolved.(map[string]interface{})
	// 		if !ok {
	// 			log.Printf("Expected map")
	// 			return 1
	// 		}
	//
	// 		s := schema.New()
	// 		if err := s.Extract(m2); err != nil {
	// 			log.Printf("%s", err)
	// 			return 1
	// 		}
	// 		schemas = append(schemas, s)
	// 	}
	// }
	//
	// log.Printf("%+v", schemas)

	// b := builder.New()

	// docs := make([]*jsdoc.JSDoc, len(schemas))
	// for i, s := range schemas {
	// 	d, err := b.BuildWithCtx(s, m)
	// 	if err != nil {
	// 		log.Printf("%s", err)
	// 		return 1
	// 	}
	// 	docs[i] = d
	// }
	//
	// var out io.Writer
	//
	// out = os.Stdout
	// if fn := opts.OutFile; fn != "" {
	// 	f, err := os.Create(fn)
	// 	if err != nil {
	// 		log.Printf("%s", err)
	// 		return 1
	// 	}
	// 	defer f.Close()
	//
	// 	out = f
	// }
	//
	// g := jsdoc.NewGenerator()
	// if err := g.Process(out, docs...); err != nil {
	// 	log.Printf("%s", err)
	// 	return 1
	// }

	return 0
}
