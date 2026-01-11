use std::process;

use rgrep::Args;

fn main() {
    if let Ok(args) = Args::build() {
        match rgrep::run(args) {
            Ok(is_match) => {
                if is_match {
                    process::exit(0);
                } else {
                    process::exit(1);
                }
            }
            Err(e) => {
                eprintln!("Application error:{e}");
                process::exit(1);
            }
        }
    } else {
        eprintln!("Error while parsing arguments");
        process::exit(1);
    }
}
