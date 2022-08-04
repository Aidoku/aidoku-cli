use anyhow::{anyhow, Result};
use log::{debug, warn};
use rayon::prelude::{IntoParallelRefIterator, ParallelIterator};
use std::{
	collections::HashMap,
	fs::File,
	io::{BufReader, Write},
	path::PathBuf,
	sync::{Arc, RwLock},
};

mod structs;
use structs::*;

static PLACEHOLDER_PNG: &[u8] = include_bytes!("./1x1-000000ff.png");

pub fn build(files: Vec<PathBuf>, output: PathBuf) -> Result<()> {
	if output.exists() {
		std::fs::remove_dir_all(&output)?;
	}
	std::fs::create_dir_all(&output)?;
	std::fs::create_dir_all(output.join("icons"))?;
	std::fs::create_dir_all(output.join("sources"))?;

	let source_ids: Arc<RwLock<HashMap<String, PathBuf>>> = Arc::new(RwLock::new(HashMap::new()));
	let (external_sources, errors): (Vec<_>, Vec<_>) = files.par_iter().map(|file| {
        let mut archive = zip::ZipArchive::new(
            BufReader::new(
                File::open(file)?
            )
        )?;
        debug!("Opened archive {}", file.display());

        let source: Source = serde_json::from_reader(archive.by_name("Payload/source.json")?)?;
        debug!("Loaded source info: {:?}", source);

        let ids = source_ids.read().expect("RwLock poisoned");
        if ids.contains_key(&source.info.id) {
            return Err(
                anyhow!(
                    "Duplicate source id: {}, first found in file {}", 
                    source.info.id,
                    ids[&source.info.id].display()
                )
            );
        }
        drop(ids);
        source_ids.write().expect("RwLock poisoned").insert(source.info.id.clone(), file.to_path_buf());

        let external_source_info: ExternalSourceInfo = source.info.into();

        if let Ok(mut img) = archive.by_name("Payload/Icon.png")
           && let mut icon = File::create(output.join("icons").join(&external_source_info.icon))?
           && std::io::copy(&mut img, &mut icon).is_ok() {
            debug!("Extracted icon to {}", output.join("icons").join(&external_source_info.icon).display());
        } else {
            warn!("Failed to extract icon from {}, generating a placeholder", file.display());
            File::create(
                output.join("icons").join(&external_source_info.icon)
            )?.write_all(PLACEHOLDER_PNG)?;
        }

        std::fs::copy(file, output.join("sources").join(&external_source_info.file))?;

        Ok(external_source_info)
    })
    .partition(Result::is_ok);

	let external_sources: Vec<ExternalSourceInfo> = external_sources
		.into_iter()
		.map(Result::unwrap_or_default)
		.collect();

	let index_json = output.join("index.json");
	let mut index_file = File::create(&index_json)?;
	serde_json::to_writer_pretty(&mut index_file, &external_sources)?;

	let index_min_json = output.join("index.min.json");
	let mut index_min_file = File::create(&index_min_json)?;
	serde_json::to_writer(&mut index_min_file, &external_sources)?;

	if errors.is_empty() {
		Ok(())
	} else {
		Err(anyhow!(
			"Errors: \n- {}",
			errors
				.into_iter()
				.map(|v| v.unwrap_err().to_string())
				.collect::<Vec<_>>()
				.join("\n- ")
		))
	}
}
