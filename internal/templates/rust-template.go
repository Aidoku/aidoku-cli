package templates

import (
	"os"
	"os/exec"
	"path"
	"strings"
)

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
	// Make the build script executable
	err = GenerateFilesFromMap(output, source, files)
	if err != nil {
		return err
	}
	os.Chmod(path.Join(output, "build.sh"), os.FileMode(0755))
	git, err := exec.LookPath("git")
	if err == nil {
		cmd := exec.Command(git, "rev-parse", "--is-inside-work-tree")
		stdout, _ := cmd.Output()
		if strings.Contains(string(stdout), "true") {
			exec.Command(git, "update-index", "--chmod=+x", "build.sh")
		}
		return nil
	} else {
		return err
	}
}
