package templates

func makefile() []byte {
	return []byte(`CC=clang
TARGET=wasm32
OPT=-O3
CC_FLAGS=$(OPT) --target=$(TARGET)
LD_FLAGS=-nostdlib -Wl,--no-entry -Wl,--export-dynamic #-W -Wall
SRCS=$(wildcard src/*.c) $(wildcard src/lib/*.c)
INCLUDE=-Isrc/include

build: clean package.aix

.build/main.wasm:
	@mkdir -p .build
	$(CC) $(CC_FLAGS) $(LD_FLAGS) $(INCLUDE) $(SRCS) -o $@

main.wat: .build/main.wasm
	wasm2wat $^ -o $@

package.aix: .build/main.wasm
	@mkdir -p .build/Payload
	@cp .build/main.wasm .build/Payload/main.wasm
	@cp res/* .build/Payload/
	@cd .build/ ; zip -r package.aix Payload > /dev/null
	@mv .build/package.aix package.aix

clean:
	@rm -rf main.wat package.aix .build
`)
}

func mainC() []byte {
	return []byte(`#include <stdbool.h>
#include <walloc.h>
#include <aidoku.h>
#include <std.h>
#include <net.h>
#include <html.h>

WASM_EXPORT
std_obj_t get_manga_list(std_obj_t filter_list_obj, int page) {
	// TODO
}

WASM_EXPORT
std_obj_t get_manga_details(std_obj_t manga_obj) {
	// TODO
}

WASM_EXPORT
std_obj_t get_chapter_list(std_obj_t manga_obj) {
	// TODO
}

WASM_EXPORT
std_obj_t get_page_list(std_obj_t chapter_obj) {
	// TODO
}
`)
}

func ClangGenerator(output string, source Source) error {
	err := GenerateCommon(output, source)
	if err != nil {
		return err
	}

	files := map[string]func() []byte{
		"/Makefile":   makefile,
		"/src/main.c": mainC,
	}
	return GenerateFilesFromMap(output, source, files)
}
