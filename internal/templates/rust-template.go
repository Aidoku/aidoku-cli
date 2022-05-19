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
name = "{{ .Name | ToLower }}_template"
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

func rustTemplateHelpers() []byte {
	return []byte(`use aidoku::{
	std::String, std::ArrayRef, std::Vec, MangaStatus, prelude::format,
};

pub fn extract_f32_from_string(title: String, text: String) -> f32 {  
    text.replace(&title, "")
        .chars()
        .filter(|a| (*a >= '0' && *a <= '9') || *a == ' ' || *a == '.')
        .collect::<String>()
        .split(" ")
        .collect::<Vec<&str>>().into_iter()
        .map(|a| a.parse::<f32>().unwrap_or(0.0))
        .find(|a| *a > 0.0)
        .unwrap_or(0.0)
}

pub fn append_protocol(url: String) -> String {
    if !url.starts_with("http") {
        return format!("{}{}", "https:", url);
    } else {
        return url;
    }
}

pub fn https_upgrade(url: String) -> String {
    return url.replacen("http://", "https://", 1);
}

pub fn urlencode(string: String) -> String {
    let mut result: Vec<u8> = Vec::with_capacity(string.len() * 3);
    let hex = "0123456789abcdef".as_bytes();
    let bytes = string.as_bytes();
    
    for byte in bytes {
        let curr = *byte;
        if (b'a' <= curr && curr <= b'z')
            || (b'A' <= curr && curr <= b'Z')
            || (b'0' <= curr && curr <= b'9') {
                result.push(curr);
        } else {
            result.push(b'%');
            result.push(hex[curr as usize >> 4]);
            result.push(hex[curr as usize & 15]);
        }
    }

    String::from_utf8(result).unwrap_or(String::new())
}

pub fn i32_to_string(mut integer: i32) -> String {
    if integer == 0 {
        return String::from("0");
    }
    let mut string = String::with_capacity(11);
    let pos = if integer < 0 {
        string.insert(0, '-');
        1
    } else {
        0
    };
    while integer != 0 {
        let mut digit = integer % 10;
        if pos == 1 {
            digit *= -1;
        }
        string.insert(pos, char::from_u32((digit as u32) + ('0' as u32)).unwrap());
        integer /= 10;
    }
    return string;
}

pub fn join_string_array(array: ArrayRef, delimeter: String) -> String {
    let mut string = String::new();
    let mut at = 0;
    for item in array {
        if at != 0 {
            string.push_str(&delimeter);
        }
        string.push_str(item.as_node().text().read().as_str());
        at += 1;
    }
    return string;
}

pub fn status_from_string(status: String) -> MangaStatus {
    if status == "Ongoing" {
        return MangaStatus::Ongoing;
    } else if status == "Completed" {
        return MangaStatus::Completed;
    } else if status == "Hiatus" {
        return MangaStatus::Hiatus;
    } else if status == "Cancelled" {
        return MangaStatus::Cancelled;
    } else {
        return MangaStatus::Unknown;
    }
}

pub fn is_numeric_char(c: char) -> bool {
    return (c >= '0' && c <= '9') || c == '.';
}

pub fn get_chapter_number(id: String) -> f32 {
    let mut number_string = String::new();
    let mut i = id.len() - 1;
    for c in id.chars().rev() {
        if !is_numeric_char(c) {
            number_string = String::from(&id[i + 1..]);
            break;
        }
        i -= 1;
    }
    if number_string.len() == 0 {
        return 0.0;
    }
    return number_string.parse::<f32>().unwrap_or(0.0);
}

pub fn string_replace(string: String, search: String, replace: String) -> String {
    let mut result = String::new();
    let mut at = 0;
    for c in string.chars() {
        if c == search.chars().next().unwrap() {
            if string[at..].starts_with(&search) {
                result.push_str(&replace);
                at += search.len();
            } else {
                result.push(c);
            }
        } else {
            result.push(c);
        }
        at += 1;
    }
    return result;
}

pub fn stupidencode(string: String) -> String {
    let mut result = String::new();
    for c in string.chars() {
        if c.is_alphanumeric() {
            result.push(c);
        } else if c == ' ' {
            result.push('_');
        }
    }
    return result;
}
`)
}

func rustTemplateTemplate() []byte {
	return []byte(`use aidoku::{
	error::Result, std::String, std::Vec, std::net::Request, std::net::HttpMethod,
	Listing, Manga, MangaPageResult, Page, MangaStatus, MangaContentRating, MangaViewer, Chapter, DeepLink
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
		"/template/src/helper.rs":   rustTemplateHelpers,
		"/template/src/template.rs": rustTemplateTemplate,
	}
	return GenerateFilesFromMap(output, source, files)
}
