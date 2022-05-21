package templates

func ClangGenerator(output string, source Source) error {
	err := GenerateCommon(output, source)
	if err != nil {
		return err
	}

	files := map[string]func() []byte{
		"/Makefile":   templateFactory(box, "clang/Makefile.tmpl"),
		"/src/main.c": templateFactory(box, "clang/src/main.c.tmpl"),
	}
	return GenerateFilesFromMap(output, source, files)
}
