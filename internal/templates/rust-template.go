package templates

import "os"

func rustTemplateMainCargo() []byte {
	return []byte(`[workspace]
members = ["template", "sources/*"]

[profile.dev]
panic = "abort"

[profile.release]
panic = "abort"
opt-level = "s"
strip = true
lto = true
`)
}

func rustTemplateTemplateCargo() []byte {
	return []byte(`[package]
name = "{{ .Name | ToLower | SlugifyRust }}_template"
version = "0.1.0"
edition = "2021"
publish = false

[dependencies]
aidoku = { git = "https://github.com/Aidoku/aidoku-rs/" }
`)
}

func rustTemplateLib() []byte {
	return []byte(`#![no_std]
pub mod helper;
pub mod template;
`)
}

func rustTemplateTemplate() []byte {
	return []byte(`use aidoku::{
	error::Result, 
	prelude::*,
	std::{String, Vec, net::Request},
	Manga, MangaPageResult, Page, Chapter, DeepLink
};

pub fn get_manga_list(todo!()) -> Result<MangaPageResult> {
	todo!()
}

pub fn get_manga_listing(todo!()) -> Result<MangaPageResult> {
	todo!()
}

pub fn get_manga_details(todo!()) -> Result<Manga> {
	todo!()
}

pub fn get_chapter_list(todo!()) -> Result<Vec<Chapter>> {
	todo!()
}

pub fn get_page_list(todo!()) -> Result<Vec<Page>> {
	todo!()
}

pub fn modify_image_request(todo!(), request: Request) {
	todo!()
}

pub fn handle_url(todo!()) -> Result<DeepLink> {
	todo!()
}
`)
}

func rustTemplatePOSIXBuildScript() []byte {
	return []byte(`# template source build script
# usage: ./build.sh [source_name/-a]

if [ "$1" != "-a" ]; then
	# compile specified source
	cargo +nightly build --release
	
	echo "packaging $1";
	mkdir -p target/wasm32-unknown-unknown/release/Payload
	cp res/* target/wasm32-unknown-unknown/release/Payload
	cp sources/$1/res/* target/wasm32-unknown-unknown/release/Payload
	cd target/wasm32-unknown-unknown/release
	cp $1.wasm Payload/main.wasm
	zip -r $1.aix Payload
	mv $1.aix ../../../$1.aix
	rm -rf Payload
else
	# compile all sources
	cargo +nightly build --release

	for dir in sources/*/
	do
		dir=${dir%*/}
		dir=${dir##*/}
		echo "packaging $dir";

		mkdir -p target/wasm32-unknown-unknown/release/Payload
		cp res/* target/wasm32-unknown-unknown/release/Payload
		cp sources/$dir/res/* target/wasm32-unknown-unknown/release/Payload
		cd target/wasm32-unknown-unknown/release
		cp $dir.wasm Payload/main.wasm
		zip -r $dir.aix Payload >> /dev/null
		mv $dir.aix ../../../$dir.aix
		rm -rf Payload
		cd ../../../
	done
fi
`)
}

func rustTemplatePS1BuildScript() []byte {
	return []byte(`<#
.SYNOPSIS
	Template source build script for Windows
#>
#requires -version 5
[cmdletbinding()]
param (
	[Parameter(ParameterSetName="help", Mandatory)]
	[alias('h')]
	[switch]$help,

	[Parameter(ParameterSetName="all", Mandatory)]
	[alias('a')]
	[switch]$all,

	[Parameter(Position=0, ParameterSetName="some", Mandatory)]
	[alias('s')]
	[string[]]$sources
)

function Package-Source {
	param (
		[Parameter(Mandatory = $true, Position = 0)]
		[String[]]$Name,

		[switch]$Build
	)
	$Name | ForEach-Object  {
		$source = $_
		if ($Build) {
			Write-Output "building $source"
			Set-Location ./sources/$source
			cargo +nightly build --release
			Set-Location ../..
		}

		Write-Output "packaging $source"
		New-Item -ItemType Directory -Path target/wasm32-unknown-unknown/release/Payload -Force | Out-Null
		Copy-Item res/* target/wasm32-unknown-unknown/release/Payload -ErrorAction SilentlyContinue
		Copy-Item sources/$source/res/* target/wasm32-unknown-unknown/release/Payload -ErrorAction SilentlyContinue
		Set-Location target/wasm32-unknown-unknown/release
		Copy-Item "$source.wasm" Payload/main.wasm
		Compress-Archive -Force -Path Payload -DestinationPath "../../../$source.aix"
		Remove-Item -Recurse -Force Payload/
		Set-Location ../../..
	}
}

if ($help -or ($null -eq $PSBoundParameters.Keys)) {
	Get-Help $MyInvocation.MyCommand.Path -Detailed
	break
}

if ($all) {
	cargo +nightly build --release
	Get-ChildItem ./sources | ForEach-Object {
		$source = (Split-Path -Leaf $_)
		Package-Source $source
	}
} else {
	$sources | ForEach-Object {
		Package-Source $_ -Build
	}
}
	`)
}

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
		"/Cargo.toml":               rustTemplateMainCargo,
		"/build.sh":                 rustTemplatePOSIXBuildScript,
		"/build.ps1":                rustTemplatePS1BuildScript,
		"/.cargo/config":            rustCargoConfigTemplate,
		"/template/Cargo.toml":      rustTemplateTemplateCargo,
		"/template/src/lib.rs":      rustTemplateLib,
		"/template/src/helper.rs":   rustHelpers,
		"/template/src/template.rs": rustTemplateTemplate,
	}
	return GenerateFilesFromMap(output, source, files)
}
