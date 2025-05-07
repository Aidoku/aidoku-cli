package templates

import (
	"embed"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"text/template"
	"unicode"

	"github.com/Aidoku/aidoku-cli/internal/common"
	"github.com/fatih/color"
	"github.com/iancoleman/strcase"
	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
)

type Source struct {
	Language, Name, Homepage, TemplateName string
	Nsfw                                   int
}

var (
	//go:embed all:resources/*
	resources embed.FS
)

func templateFactory(box embed.FS, path string) func() []byte {
	return func() []byte {
		if bytes, err := box.ReadFile("resources/" + path); err != nil {
			panic(err)
		} else {
			return bytes
		}
	}
}

type ToPascalCase struct {
}

func (t ToPascalCase) Transform(dst, src []byte, atEOF bool) (nDst int, nSrc int, err error) {
	result := []byte(strcase.ToCamel(string(src)))
	nDst = copy(dst, result)
	nSrc = len(src)
	if nDst < nSrc {
		err = transform.ErrShortDst
	}
	return
}

func (t ToPascalCase) Reset() {

}

func slugifyFactory(whitespaceReplacer string, t transform.Transformer) func(string) string {
	return func(val string) string {
		temp := strings.ReplaceAll(strings.TrimSpace(val), " ", whitespaceReplacer)
		ret, _, _ := transform.String(t, temp)
		return ret
	}
}

func GenerateFilesFromMap(output string, source Source, files map[string]func() []byte) error {
	funcMap := template.FuncMap{
		"ToLower":      strings.ToLower,
		"SlugifyRust":  slugifyFactory("_", transform.Chain(norm.NFD, runes.Remove(runes.In(unicode.Mn)), norm.NFC)),
		"SlugifyAs":    slugifyFactory("-", transform.Chain(norm.NFD, runes.Remove(runes.In(unicode.Mn)), norm.NFC)),
		"SlugifyClass": slugifyFactory(" ", transform.Chain(norm.NFD, runes.Remove(runes.In(unicode.Mn)), norm.NFC, ToPascalCase{})),
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
			if filepath.Ext(key) == "sh" {
				os.Chmod(output+key, os.FileMode(0755))
			}

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
		"/res/source.json":   templateFactory(resources, "common/res/source.json.tmpl"),
		"/res/filters.json":  templateFactory(resources, "common/res/filters.json.tmpl"),
		"/res/settings.json": templateFactory(resources, "common/res/settings.json.tmpl"),
	}
	return GenerateFilesFromMap(output, source, files)
}
