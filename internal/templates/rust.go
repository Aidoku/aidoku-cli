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
name = "{{ .Name | ToLower | SlugifyRust }}"
version = "0.1.0"
edition = "2021"

[lib]
crate-type = ["cdylib"]
{{ if eq (len .TemplateName) 0 }}
[profile.dev]
panic = "abort"

[profile.release]
panic = "abort"
opt-level = "s"
strip = true
lto = true
{{ end }}
[dependencies]
aidoku = { git = "https://github.com/Aidoku/aidoku-rs" }{{ if not (eq (len .TemplateName) 0) }}
{{ .TemplateName }} = { path = "../../template" }
{{ end }}
`)
}

func rustCargoConfigTemplate() []byte {
	return []byte(`[build]
target = "wasm32-unknown-unknown"
`)
}

func rustHelpers() []byte {
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
    let hex = "0123456789ABCDEF".as_bytes();
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
    return match status.as_str() {
        "Ongoing" => MangaStatus::Ongoing,
        "Completed" => MangaStatus::Completed,
        "Hiatus" => MangaStatus::Hiatus,
        "Cancelled" => MangaStatus::Cancelled,
        _ => MangaStatus::Unknown,
    };
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

func rustLibTemplate() []byte {
	return []byte(`#![no_std]
mod helper;
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
	if len(source.TemplateName) == 0 {
		os.MkdirAll(output+"/.cargo", os.FileMode(0754))
		files["/.cargo/config"] = rustCargoConfigTemplate
		files["/build.sh"] = rustPOSIXBuildScript
		files["/build.ps1"] = rustPS1BuildScript
		files["/src/helper.rs"] = rustHelpers
	}
	return GenerateFilesFromMap(output, source, files)
}
