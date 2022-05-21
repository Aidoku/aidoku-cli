package templates

import (
	"os"
)

func RustGenerator(output string, source Source) error {
	err := GenerateCommon(output, source)
	if err != nil {
		return err
	}

	files := map[string]func() []byte{
		"/src/lib.rs": templateFactory(box, "rust/src/lib.rs.tmpl"),
		"/Cargo.toml": templateFactory(box, "rust/Cargo.toml.tmpl"),
	}
	if len(source.TemplateName) == 0 {
		os.MkdirAll(output+"/.cargo", os.FileMode(0754))
		files["/.cargo/config"] = templateFactory(box, "rust/.cargo/config.tmpl")
		files["/build.sh"] = templateFactory(box, "rust/build.sh.tmpl")
		files["/build.ps1"] = templateFactory(box, "rust/build.ps1.tmpl")
		files["/src/helper.rs"] = templateFactory(box, "rust/src/helper.rs.tmpl")
	}
	return GenerateFilesFromMap(output, source, files)
}
