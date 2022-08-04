use serde::{Deserialize, Serialize};

#[derive(Serialize, Deserialize, Debug, Clone)]
#[serde(rename_all = "camelCase")]
pub struct SourceInfo {
	pub id: String,
	pub name: String,
	pub lang: String,
	pub version: usize,
	pub nsfw: usize,

	#[serde(skip_serializing_if = "Option::is_none")]
	pub min_app_version: Option<String>,

	#[serde(skip_serializing_if = "Option::is_none")]
	pub max_app_version: Option<String>,
}

#[derive(Serialize, Deserialize, Debug, Clone, Default)]
#[serde(rename_all = "camelCase")]
pub struct ExternalSourceInfo {
	pub id: String,
	pub name: String,
	pub file: String,
	pub icon: String,
	pub lang: String,
	pub version: usize,
	pub nsfw: usize,

	#[serde(skip_serializing_if = "Option::is_none")]
	pub min_app_version: Option<String>,

	#[serde(skip_serializing_if = "Option::is_none")]
	pub max_app_version: Option<String>,
}

impl From<SourceInfo> for ExternalSourceInfo {
	fn from(src: SourceInfo) -> Self {
		Self {
			file: format!("{}-v{}.aix", src.id, src.version),
			icon: format!("{}-v{}.png", src.id, src.version),
			id: src.id,
			name: src.name,
			lang: src.lang,
			version: src.version,
			nsfw: src.nsfw,
			min_app_version: src.min_app_version,
			max_app_version: src.max_app_version,
		}
	}
}

#[derive(Serialize, Deserialize, Debug, Clone)]
pub struct Source {
	pub info: SourceInfo,

	#[serde(skip_serializing_if = "Option::is_none")]
	pub languages: Option<Vec<Language>>,

	#[serde(skip_serializing_if = "Option::is_none")]
	pub listings: Option<Vec<Listing>>,
}

#[derive(Serialize, Deserialize, Debug, Clone)]
pub struct Language {
	pub code: String,

	#[serde(skip_serializing_if = "Option::is_none")]
	pub value: Option<String>,

	#[serde(skip_serializing_if = "Option::is_none")]
	pub default: Option<bool>,
}

#[derive(Serialize, Deserialize, Debug, Clone)]
pub struct Listing {
	pub name: String,
}
