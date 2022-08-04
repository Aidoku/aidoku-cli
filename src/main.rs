#![deny(missing_debug_implementations)]

mod build;
mod logcat;

use clap::{CommandFactory, Parser, Subcommand};
use log::error;
use std::{
	net::{IpAddr, Ipv4Addr, SocketAddr},
	path::PathBuf,
};
use warp::Filter;

pub static PLACEHOLDER_PNG: &[u8] = include_bytes!("./res/1x1-000000ff.png");

#[derive(Debug, Clone, Parser)]
struct Cli {
	#[clap(flatten)]
	verbose: clap_verbosity_flag::Verbosity,

	#[clap(subcommand)]
	command: Commands,
}

#[derive(Subcommand, Debug, Clone)]
enum Commands {
	/// Build a source list from the given packages
	Build {
		#[clap(value_parser)]
		files: Vec<String>,

		/// Output folder
		#[clap(short, long, value_parser, default_value_t = String::from("public"))]
		output: String,
	},

	/// Generate completion script to stdout
	Completion {
		#[clap(arg_enum, value_parser)]
		shell: clap_complete::Shell,
	},

	/// Generate the boilerplate for an Aidoku source
	Init {
		#[clap(value_parser)]
		path: String,

		#[clap(short, long, value_parser)]
		homepage: Option<String>,

		#[clap(short, long, value_parser)]
		language: Option<String>,

		#[clap(short, long, value_parser)]
		name: Option<String>,

		#[clap(short, long, value_parser)]
		nsfw: Option<usize>,
	},

	/// Build a source list and serve it on the local network
	Serve {
		#[clap(value_parser)]
		files: Vec<String>,

		/// Output folder
		#[clap(short, long, value_parser, default_value_t = String::from("public"))]
		output: String,

		/// Address to serve the source list on
		#[clap(short, long, value_parser, default_value_t = IpAddr::V4(Ipv4Addr::new(0, 0, 0, 0)))]
		address: IpAddr,

		/// Port to serve the source list on
		#[clap(short, long, value_parser, default_value_t = 8080)]
		port: u16,

		#[clap(short, long, value_parser)]
		watch: bool,
	},

	/// Log streaming
	Logcat {
		/// Address to broadcast log server on
		#[clap(short, long, value_parser, default_value_t = IpAddr::V4(Ipv4Addr::new(0, 0, 0, 0)))]
		address: IpAddr,

		/// The port to broadcast the log server on
		#[clap(short, long, value_parser, default_value_t = 9000)]
		port: u16,
	},

	/// Test Aidoku packages to see if they are ready for publishing
	Verify {
		#[clap(value_parser)]
		files: Vec<String>,
	},
}

#[tokio::main]
async fn main() {
	let cli = Cli::parse_from(wild::args_os());

	env_logger::builder()
		.filter_level(cli.verbose.log_level_filter())
		.format_timestamp(None)
		.init();

	match &cli.command {
		Commands::Build { files, output } => {
			let files = files.iter().map(|p| p.into()).collect::<Vec<PathBuf>>();
			let output = output.into();
			match build::build(files, output) {
				Ok(_) => (),
				Err(e) => {
					error!("{}", e);
					std::process::exit(1);
				}
			}
		}
		Commands::Completion { shell } => {
			clap_complete::generate(
				shell.to_owned(),
				&mut Cli::command(),
				"aidoku",
				&mut std::io::stdout(),
			);
		}
		Commands::Logcat { address, port } => {
			let path = warp::post()
				.and(warp::body::bytes())
				.map(logcat::logcat_handler);
			println!("Listening on {}:{}", address, port);
			warp::serve(path)
				.run(SocketAddr::new(*address, *port))
				.await;
		}
		Commands::Serve {
			files,
			output,
			address,
			port,
			watch,
		} => {
			let files = files.iter().map(|p| p.into()).collect::<Vec<PathBuf>>();
			let output: PathBuf = output.into();
			match build::build(files, output.clone()) {
				Ok(_) => (),
				Err(e) => {
					error!("{}", e);
					std::process::exit(1);
				}
			}
			let route = warp::any()
				.and(warp::fs::dir(output))
				.with(warp::log("aidoku_cli"));
			println!("Listening on {}:{}", address, port);
			warp::serve(route)
				.run(SocketAddr::new(*address, *port))
				.await;
		}
		Commands::Verify { files } => {
			todo!()
		}
		Commands::Init {
			path,
			homepage,
			language,
			name,
			nsfw,
		} => {
			todo!()
		}
	}
}
