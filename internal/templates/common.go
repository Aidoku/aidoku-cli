package templates

import (
	"os"
	"strings"
	"sync"
	"text/template"

	"github.com/beerpiss/aidoku-cli/internal/common"
	"github.com/fatih/color"
)

type Source struct {
	Language, Name, Homepage, TemplateName string
	Nsfw                                   int
}

func commonSourceJson() []byte {
	return []byte(`{
	"info": {
		"id": "{{ .Language }}.{{ .Name | ToLower }}",
		"lang": "{{ .Language }}",
		"name": "{{ .Name }}",
		"version": 1,
		"url": "{{ .Homepage }}",
		"nsfw": {{ .Nsfw }}
	},
	"languages": [

	],
	"listings": [

	]
}
`)
}

func commonSettingsJson() []byte {
	return []byte(`[
	{
		"type": "group",
		"title": "Settings",
		"footer": "You can have footers for a setting group",
		"items": [
			{
				"type": "select",
				"key": "something",
				"title": "A select option",
				"values": ["your", "values", "here"],
				"titles": ["your", "titles", "here"],
				"default": "here"
			},
			{
				"type": "switch",
				"key": "something2",
				"title": "A switch",
				"subtitle": "A subtext to describe this option",
				"default": false
			},
			{
				"type": "text",
				"key": "something3",
				"placeholder": "A text box"
			}
		]
	}
]`)
}

func commonFilterJson() []byte {
	return []byte(`[
	{
		"type": "title"
	},
	{
		"type": "author"
	},
	{
		"type": "group",
		"name": "A group filter",
		"filters": [
			{
				"type": "check",
				"name": "A checkbox",
				"default": true
			}
		]
	},
	{
		"type": "group",
		"name": "Another group filter",
		"filters": [
			{
				"type": "genre",
				"name": "A genre",
				"canExclude": true
			}
		]
	},
	{
		"type": "select",
		"name": "A select filter",
		"options": [

		]
	},
	{
		"type": "sort",
		"name": "A sorter",
		"canAscend": true,
		"options": [

		],
		"default": {
			"index": 0,
			"ascending": false
		}
	}
]
`)
}

func GenerateFilesFromMap(output string, source Source, files map[string]func() []byte) error {
	funcMap := template.FuncMap{
		"ToLower": strings.ToLower,
	}
	var wg sync.WaitGroup
	errc := make(chan error, 10)
	for key, value := range files {
		wg.Add(1)
		go func(key string, value func() []byte) {
			defer wg.Done()
			file, err := os.Create(output + key)
			if err != nil {
				color.Red("error: could not create %s: %s", key, err.Error())
				errc <- err
				return
			}
			defer file.Close()

			fileTemplate := template.Must(template.New(key).Funcs(funcMap).Parse(string(value())))
			err = fileTemplate.Execute(file, source)
			if err != nil {
				color.Red("error: could not generate %s from template: %s", key, err.Error())
				errc <- err
				return
			}
		}(key, value)
	}
	go func() {
		wg.Wait()
		close(errc)
	}()

	for err := range errc {
		if err != nil {
			return err
		}
	}
	return nil
}

func GenerateCommon(output string, source Source) error {
	// Create source directory
	os.MkdirAll(output+"/src", os.FileMode(0754))
	os.MkdirAll(output+"/res", os.FileMode(0754))

	// generate placeholder Icon.png
	err := common.GeneratePng(output + "/res/Icon.png")
	if err != nil {
		color.Red("error: could not generate placeholder icon")
		return err
	}

	files := map[string]func() []byte{
		"/res/source.json":   commonSourceJson,
		"/res/filters.json":  commonFilterJson,
		"/res/settings.json": commonSettingsJson,
	}
	return GenerateFilesFromMap(output, source, files)
}
