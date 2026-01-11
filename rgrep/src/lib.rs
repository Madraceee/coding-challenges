use clap::Parser;
use regex::Regex;
use std::error::Error;
use std::fs;

#[derive(Parser, Debug)]
#[command(version, about, long_about = None)]
pub struct Args {
    #[arg(short, long, default_value_t = false)]
    ignore_case: bool,

    #[arg(short = 'v', long, default_value_t = false)]
    invert_search: bool,

    query: String,

    file: String,
}

impl Args {
    pub fn build() -> Result<Args, Box<dyn Error>> {
        Ok(Args::parse())
    }
}

// If match is found return true,
// else return false
pub fn run(args: Args) -> Result<bool, Box<dyn Error>> {
    let contents = fs::read_to_string(args.file)?;

    let results = if args.ignore_case {
        search_insensitive(&args.query, &contents, args.invert_search)?
    } else {
        search(&args.query, &contents, args.invert_search)?
    };

    let count = results.len();

    for line in results {
        println!("{line}");
    }

    if count > 0 {
        return Ok(true);
    }

    return Ok(false);
}

fn search<'a>(
    query: &str,
    contents: &'a str,
    invert_search: bool,
) -> Result<Vec<&'a str>, Box<dyn Error>> {
    let regex = Regex::new(query)?;
    Ok(contents
        .lines()
        .filter(|line| regex.is_match(line) ^ invert_search)
        .collect())
}

fn search_insensitive<'a>(
    query: &str,
    contents: &'a str,
    invert_search: bool,
) -> Result<Vec<&'a str>, Box<dyn Error>> {
    let regex = Regex::new(&query.to_lowercase())?;
    Ok(contents
        .lines()
        .filter(|line| regex.is_match(&line.to_lowercase()) ^ invert_search)
        .collect())
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn case_sensitive() {
        let query = "duct";
        let contents = "\
Rust:
safe, fast, productive.
Pick three.
Duct tape.";

        assert_eq!(
            vec!["safe, fast, productive."],
            search(query, contents, false).unwrap()
        );
    }

    #[test]
    fn case_insensitive() {
        let query = "rUsT";
        let contents = "\
Rust:
safe, fast, productive.
Pick three.
Trust me.";

        assert_eq!(
            vec!["Rust:", "Trust me."],
            search_insensitive(query, contents, false).unwrap()
        );
    }

    #[test]
    fn search_revert() {
        let query = "duct";
        let contents = "\
Rust:
safe, fast, productive.
Pick three.
Duct tape.";

        assert_eq!(
            vec!["Rust:", "Pick three.", "Duct tape."],
            search(query, contents, true).unwrap()
        )
    }

    #[test]
    fn search_revert_case_insensitive() {
        let query = "duct";
        let contents = "\
Rust:
safe, fast, productive.
Pick three.
Duct tape.";

        assert_eq!(
            vec!["Rust:", "Pick three."],
            search_insensitive(query, contents, true).unwrap()
        )
    }

    #[test]
    fn search_digits() {
        let query = "\\d";
        let contents = "\
1
2
three
4
five
six";
        assert_eq!(vec!["1", "2", "4"], search(query, contents, false).unwrap())
    }

    #[test]
    fn search_words() {
        let query = "\\w";
        let contents = "\
!
@
three
$
five
six";
        assert_eq!(
            vec!["three", "five", "six"],
            search(query, contents, false).unwrap()
        )
    }

    #[test]
    fn search_starting() {
        let query = "^a";
        let contents = "\
apple
man
ends with a";
        assert_eq!(vec!["apple"], search(query, contents, false).unwrap())
    }

    #[test]
    fn search_ending() {
        let query = "a$";
        let contents = "\
apple
man
ends with a";
        assert_eq!(vec!["ends with a"], search(query, contents, false).unwrap())
    }
}
