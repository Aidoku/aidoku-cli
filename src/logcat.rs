use std::io::Write;
use termcolor::{BufferWriter, Color, ColorChoice, ColorSpec, WriteColor};
use warp::hyper::body::Bytes;

pub fn logcat_handler(data: Bytes) -> impl warp::Reply {
	let log = String::from_utf8_lossy(&data);

	let start_date = log.find('[').unwrap_or(0) + 1;
	let end_date = log[start_date..]
		.find(']')
		.map(|v| v + start_date)
		.unwrap_or(19);
	let date =
		chrono::NaiveDateTime::parse_from_str(&log[start_date..end_date], "%m/%d %H:%M:%S.%L")
			.unwrap_or_else(|_| chrono::Utc::now().naive_utc())
			.format("%Y-%m-%dT%H:%M:%S");

	let loglevel = if log.chars().nth(end_date + 2).unwrap() == '[' {
		let end = log[end_date + 3..]
			.find(']')
			.map(|v| v + end_date + 3)
			.unwrap_or(log.len());
		&log[end_date + 3..end]
	} else {
		"INFO"
	};

	let log_content = log[log.rfind(']').unwrap_or(0) + 2..].to_string();

	if atty::is(atty::Stream::Stdout) {
		let bufwtr = BufferWriter::stdout(ColorChoice::Auto);
		let mut buffer = bufwtr.buffer();

		buffer
			.set_color(
				ColorSpec::new()
					.set_fg(Some(Color::Black))
					.set_intense(true),
			)
			.unwrap();
		write!(&mut buffer, "[").unwrap();

		buffer.reset().unwrap();
		write!(&mut buffer, "{}", date).unwrap();

		buffer
			.set_color(
				ColorSpec::new()
					.set_fg(Some(match loglevel {
						"DEBUG" => Color::Blue,
						"WARN" => Color::Yellow,
						"ERROR" => Color::Red,
						_ => Color::Green,
					}))
					.set_bold(loglevel == "ERROR"),
			)
			.unwrap();
		write!(&mut buffer, " {:<5}", loglevel).unwrap();

		buffer
			.set_color(
				ColorSpec::new()
					.set_fg(Some(Color::Black))
					.set_intense(true),
			)
			.unwrap();
		write!(&mut buffer, "]").unwrap();

		buffer.reset().unwrap();
		writeln!(&mut buffer, " {}", log_content).unwrap();

		bufwtr.print(&buffer).ok();
	} else {
		println!("[{} {:<5}] {}", date, loglevel, log_content);
	}

	warp::reply()
}
