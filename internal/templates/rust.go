package templates

import (
	"os"
	"os/exec"
	"path"
	"strings"
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
