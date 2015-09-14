package main

import (
	"encoding/json"
	"os"
	"regexp"
	"strings"
	"text/template"

	"github.com/akiomik/schematic"
	"github.com/codegangsta/cli"
	log "github.com/golang/glog"
)

func main() {
	app := cli.NewApp()
	app.Name = "subschema"
	app.Version = Version
	app.Usage = "An utility for JSON Schema and API documentations"
	app.Author = "Akiomi Kamakura"
	app.Email = "akiomik@gmail.com"
	app.Commands = []cli.Command{
		cli.Command{
			Name:      "convert",
			ShortName: "c",
			Usage:     "Convert a JSON Schema into API documentations",
			Action:    doConvert,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "format",
					Value: "apib", // TODO
					Usage: "output format for the conversion",
				},
			},
		},
	}
	app.Run(os.Args)
}

func exampleObject(ps map[string]*schematic.Schema) map[string]interface{} {
	example := map[string]interface{}{}
	for i, p := range ps {
		if p.Type == "array" {
			example[i] = exampleArray(p.Items)
		} else if p.Type == "object" {
			example[i] = exampleObject(p.Definitions)
		} else {
			example[i] = p.Example
		}
	}

	return example
}

func exampleArray(s *schematic.Schema) []interface{} {
	var example []interface{}

	if s.Type == "array" {
		example = []interface{}{exampleArray(s.Items)}
	} else if s.Type == "object" {
		example = []interface{}{exampleObject(s.Definitions)}
	} else {
		example = []interface{}{s.Example}
	}

	return example
}

func doConvert(c *cli.Context) {
	// read args
	jsonSchema := c.Args().Get(0)
	if jsonSchema == "" {
		log.Error("JSON Schema is not provided")
		os.Exit(1)
	}
	tmplName := c.String("format") + ".tmpl"
	if tmplName == "" {
	}

	// open json schema file
	reader, err1 := os.Open(jsonSchema)
	if err1 != nil {
		panic(err1)
	}
	defer reader.Close()

	// define template
	funcMap := template.FuncMap{
		"title": func(s string) string {
			return strings.Title(s)
		},
		"canFormat": func(s string) bool {
			matched, _ := regexp.MatchString("%v", s)
			return matched
		},
		"statusCode": func(l *schematic.Link) int {
			if l.Method == "POST" {
				return 201
			} else if l.MediaType == "null" {
				return 204
			} else {
				return 200
			}
		},
		"body": func(l *schematic.Link, d *schematic.Schema) string {
			var props map[string]*schematic.Schema

			if l.Schema == nil {
				props = d.Properties
			} else {
				props = l.Schema.Properties
			}

			bodyJson, _ := json.Marshal(exampleObject(props))
			return string(bodyJson)
		},
	}
	tmpl := template.Must(template.New(tmplName).Funcs(funcMap).ParseFiles("templates/" + tmplName))

	// decode json schema
	var schema schematic.Schema
	decoder := json.NewDecoder(reader)
	decoder.Decode(&schema)

	// convert json schema into api documentation
	err := tmpl.Execute(os.Stdout, schema.Resolve(&schema))
	if err != nil {
		panic(err)
	}
}
