package templates

import "os"

func RustTemplateGenerator(output string, source Source) error {
	err := GenerateCommon(output, source)
	if err != nil {
		return err
	}
	os.MkdirAll(output+"/.cargo", os.FileMode(0754))
	os.MkdirAll(output+"/sources", os.FileMode(0754))
	os.MkdirAll(output+"/template/src", os.FileMode(0754))
	os.Remove(output + "/src")
	os.RemoveAll(output + "/res")

	files := map[string]func() []byte{
		"/Cargo.toml":               templateFactory(box, "rust-template/Cargo.toml.tmpl"),
		"/build.sh":                 templateFactory(box, "rust-template/build.sh.tmpl"),
		"/build.ps1":                templateFactory(box, "rust-template/build.ps1.tmpl"),
		"/.cargo/config":            templateFactory(box, "rust/.cargo/config.tmpl"),
		"/template/Cargo.toml":      templateFactory(box, "rust-template/template/Cargo.toml.tmpl"),
		"/template/src/lib.rs":      templateFactory(box, "rust-template/template/src/lib.rs.tmpl"),
		"/template/src/helper.rs":   templateFactory(box, "rust/src/helper.rs.tmpl"),
		"/template/src/template.rs": templateFactory(box, "rust-template/template/src/template.rs.tmpl"),
	}
	return GenerateFilesFromMap(output, source, files)
}
