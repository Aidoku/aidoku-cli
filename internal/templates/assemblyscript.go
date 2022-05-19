package templates

import "fmt"

func asConfigJson() []byte {
	return []byte(`{
	"targets": {
		"debug": {
			"binaryFile": "build/untouched.wasm",
			"textFile": "build/untouched.wat",
			"sourceMap": true,
			"debug": true
		},
		"release": {
			"binaryFile": "build/optimized.wasm",
			"textFile": "build/optimized.wat",
			"sourceMap": true,
			"optimizeLevel": 3,
			"shrinkLevel": 0,
			"converge": false,
			"noAssert": false
		}
	},
	"options": {}
}
`)
}

func packageJson() []byte {
	return []byte(`{
	"name": "{{ .Name | ToLower }}",
	"version": "0.1.0",
	"ascMain": "src/index.ts",
	"dependencies": {
		"aidoku-as": "github:Aidoku/aidoku-as"
	},
	"devDependencies": {
		"assemblyscript": "^0.19.23"
	},
	"scripts": {
		"build": "rm -rf build/package.aix build/Payload\nnpm run asbuild:optimized\nmkdir -p build/Payload\ncp build/optimized.wasm build/Payload/main.wasm\ncp res/* build/Payload/\ncd build/ ; zip -r package.aix Payload",
		"asbuild:untouched": "asc src/index.ts --target debug",
		"asbuild:optimized": "asc src/index.ts --target release",
		"asbuild": "npm run asbuild:untouched && npm run asbuild:optimized"
	}
}
`)
}

func tsConfigJson() []byte {
	return []byte(`{
	"extends": "assemblyscript/std/assembly.json",
	"include": [
		"./**/*.ts"
	]
}
`)
}

func indexTs() []byte {
	return []byte(`import {
	ArrayRef,
	Filter,
	Listing,
	Request,
	ValueRef,
	DeepLink,
} from "aidoku-as/src/index";

import { {{ .Name }} as Source } from "./{{ .Name }}";

let source = new Source();

export function get_manga_list(filter_list_descriptor: i32, page: i32): i32 {
	let filters: Filter[] = [];
	let objects = new ValueRef(filter_list_descriptor).asArray().toArray();
	for (let i = 0; i < objects.length; i++) {
		filters.push(new Filter(objects[i].asObject()));
	}
	let result = source.getMangaList(filters, page);
	return result.value;
}

export function get_manga_listing(listing: i32, page: i32): i32 {
	return source.getMangaListing(new Listing(listing), page).value;
}

export function get_manga_details(manga_descriptor: i32): i32 {
	let id = new ValueRef(manga_descriptor).asObject().get("id").toString();
	return source.getMangaDetails(id).value;
}

export function get_chapter_list(manga_descriptor: i32): i32 {
	let id = new ValueRef(manga_descriptor).asObject().get("id").toString();
	let array = ArrayRef.new();
	let result = source.getChapterList(id);
	for (let i = 0; i < result.length; i++) {
		array.push(new ValueRef(result[i].value));
	}
	return array.value.rid;
}

export function get_page_list(chapter_descriptor: i32): i32 {
	let id = new ValueRef(chapter_descriptor).asObject().get("id").toString();
	let array = ArrayRef.new();
	let result = source.getPageList(id);
	for (let i = 0; i < result.length; i++) {
		array.push(new ValueRef(result[i].value));
	}
	return array.value.rid;
}

export function modify_image_request(req: i32): void {
	let request = new Request(req);
	source.modifyImageRequest(request);
}

export function handle_url(url: i32): i32 {
	let result = source.handleUrl(new ValueRef(url).toString());
	if (result == null) return -1;
	return (result as DeepLink).value;
}
`)
}

func sourceTs() []byte {
	return []byte(`import {
	Chapter,
	DeepLink,
	Filter,
	Listing,
	Manga,
	MangaPageResult,
	Page,
	Request,
	Source,
} from "aidoku-as/src/index";


export class {{ .Name }} extends Source {
	constructor() {
		super();
		// TODO
	}

	modifyImageRequest(request: Request): void {
		// TODO
	}

	getMangaList(filters: Filter[], page: i32): MangaPageResult {
		// TODO
	}

	getMangaListing(listing: Listing, page: i32): MangaPageResult {
		// TODO
		
	}

	getMangaDetails(mangaId: string): Manga {
		// TODO
	}

	getChapterList(mangaId: string): Chapter[] {
		// TODO
	}

	getPageList(chapterId: string): Page[] {
		// TODO
	}

	private getMangaDetailsFromChapterPage(chapterId: string): Manga {
		// TODO
	}

	handleUrl(url: string): DeepLink | null {
		// TODO
	}
}
`)
}

func AscGenerator(output string, source Source) error {
	err := GenerateCommon(output, source)
	if err != nil {
		return err
	}

	files := map[string]func() []byte{
		"/tsconfig.json":                       tsConfigJson,
		"/asconfig.json":                       asConfigJson,
		"/package.json":                        packageJson,
		"/src/index.ts":                        indexTs,
		fmt.Sprintf("/src/%s.ts", source.Name): sourceTs,
	}
	return GenerateFilesFromMap(output, source, files)
}
