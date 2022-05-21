package templates

import (
	"fmt"
	"unicode"

	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
)

func AscGenerator(output string, source Source) error {
	err := GenerateCommon(output, source)
	if err != nil {
		return err
	}
	className := slugifyFactory(" ", transform.Chain(norm.NFD, runes.Remove(runes.In(unicode.Mn)), norm.NFC, ToPascalCase{}))(source.Name)

	files := map[string]func() []byte{
		"/tsconfig.json":                     templateFactory(box, "assemblyscript/tsconfig.json.tmpl"),
		"/asconfig.json":                     templateFactory(box, "assemblyscript/asconfig.json.tmpl"),
		"/package.json":                      templateFactory(box, "assemblyscript/package.json.tmpl"),
		"/src/index.ts":                      templateFactory(box, "assemblyscript/src/index.ts.tmpl"),
		fmt.Sprintf("/src/%s.ts", className): templateFactory(box, "assemblyscript/src/Source.ts.tmpl"),
	}
	return GenerateFilesFromMap(output, source, files)
}
