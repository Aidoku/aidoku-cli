package templates

import (
	"os"
)

func rustPS1BuildScript() []byte {
	return []byte(`function Package-Source {
	param (
		[Parameter(Mandatory = $true, Position = 0)]
		[String[]]$Name,

		[switch]$Build
	)
	$Name | ForEach-Object	{
		$source = $_
		if ($Build) {
			Write-Output "building $source"
			cargo +nightly build --release
		}

		Write-Output "packaging $source"
		New-Item -ItemType Directory -Path target/wasm32-unknown-unknown/release/Payload -Force | Out-Null
		Copy-Item res/* target/wasm32-unknown-unknown/release/Payload -ErrorAction SilentlyContinue
		Copy-Item sources/$source/res/* target/wasm32-unknown-unknown/release/Payload -ErrorAction SilentlyContinue
		Set-Location target/wasm32-unknown-unknown/release
		Copy-Item "$source.wasm" Payload/main.wasm
		Compress-Archive -Force -DestinationPath "../../../$source.aix" -Path Payload
		Remove-Item -Recurse -Force Payload/
		Set-Location ../../..
	}
}

Package-Source {{ .Name | ToLower }} -Build
`)
}

func rustPOSIXBuildScript() []byte {
	return []byte(`cargo +nightly build --release
mkdir -p target/wasm32-unknown-unknown/release/Payload
cp res/* target/wasm32-unknown-unknown/release/Payload
cp target/wasm32-unknown-unknown/release/*.wasm target/wasm32-unknown-unknown/release/Payload/main.wasm
cd target/wasm32-unknown-unknown/release ; zip -r package.aix Payload
mv package.aix ../../../package.aix
`)
}

func rustCargoTemplate() []byte {
	return []byte(`[package]
name = "{{ .Name | ToLower }}"
version = "0.1.0"
edition = "2021"

[lib]
crate-type = ["cdylib"]
{{ if not .ChildTemplate }}
[profile.dev]
panic = "abort"

[profile.release]
panic = "abort"
opt-level = "s"
strip = true
lto = true
{{ end }}
[dependencies]
aidoku = { git = "https://github.com/Aidoku/aidoku-rs" }{{ if .ChildTemplate }}
template = { path = "../../template" }
{{ end }}
`)
}

func rustCargoConfigTemplate() []byte {
	return []byte(`[build]
target = "wasm32-unknown-unknown"
`)
}

func rustLibTemplate() []byte {
	return []byte(`#![no_std]
use aidoku::{
	error::Result,
	prelude::*,
	std::{String, Vec},
	Chapter, Filter, Listing, Manga, MangaPageResult, Page, DeepLink
};

#[get_manga_list]
fn get_manga_list(_: Vec<Filter>, _: i32) -> Result<MangaPageResult> {
	todo!()
}

#[get_manga_listing]
fn get_manga_listing(_: Listing, _: i32) -> Result<MangaPageResult> {
	todo!()
}

#[get_manga_details]
fn get_manga_details(_: String) -> Result<Manga> {
	todo!()
}

#[get_chapter_list]
fn get_chapter_list(_: String) -> Result<Vec<Chapter>> {
	todo!()
}

#[get_page_list]
fn get_page_list(_: String) -> Result<Vec<Page>> {
	todo!()
}

#[modify_image_request]
fn modify_image_request(_: Request) {
	todo!()
}

#[handle_url]
fn handle_url(_: String) -> Result<DeepLink> {
	todo!()
}
`)
}

func RustGenerator(output string, source Source) error {
	err := GenerateCommon(output, source)
	if err != nil {
		return err
	}

	files := map[string]func() []byte{
		"/src/lib.rs": rustLibTemplate,
		"/Cargo.toml": rustCargoTemplate,
	}
	if !source.ChildTemplate {
		os.MkdirAll(output+"/.cargo", os.FileMode(0754))
		files["/.cargo/config"] = rustCargoConfigTemplate
		files["/build.sh"] = rustPOSIXBuildScript
		files["/build.ps1"] = rustPS1BuildScript
	}
	return GenerateFilesFromMap(output, source, files)
}
