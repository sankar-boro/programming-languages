# Chapter 10 — Error Handling

> *"Rust forces you to acknowledge that errors can happen. The Result type is not optional — it's the language's way of making you write robust code."*

---

## 10.1 Two Categories of Errors

Rust distinguishes between two fundamentally different kinds of failures:

**Unrecoverable errors** (`panic!`)  
- Programming bugs — things that should never happen
- Out-of-bounds array access, arithmetic overflow, assertion failures
- The program crashes immediately with a stack trace
- Not for user-facing error handling

**Recoverable errors** (`Result<T, E>`)  
- Expected failure modes — file not found, invalid input, network timeout
- The function returns an Err variant, caller decides how to handle
- No crash — program continues with error handling logic

---

## 10.2 panic! — Unrecoverable Errors

```rust
fn main() {
    // Explicit panic
    panic!("Something went terribly wrong");

    // Automatic panics from the standard library:
    let v = vec![1, 2, 3];
    // v[99]  // panics: index out of bounds
    
    // None.unwrap()  // panics: called unwrap() on None
    
    // let n: i32 = "abc".parse().unwrap();  // panics: invalid digit

    // assertion macros
    assert!(1 == 1);                    // panics if false
    assert_eq!(2 + 2, 4);              // panics if not equal
    assert_ne!("hello", "world");      // panics if equal

    // unreachable! — marks code that should never execute
    let x = 5;
    match x {
        1..=10 => println!("in range"),
        _ => unreachable!("x is always 1-10"),
    }

    // todo! — marks unimplemented code
    fn not_yet_implemented() -> i32 {
        todo!()  // panics when called
    }
}
```

### RUST_BACKTRACE

When a panic occurs, set `RUST_BACKTRACE=1` to see the full stack trace:

```bash
RUST_BACKTRACE=1 cargo run
```

---

## 10.3 Result<T, E> — Recoverable Errors

```rust
enum Result<T, E> {
    Ok(T),   // success
    Err(E),  // failure
}
```

Functions that can fail return `Result`. The caller is forced to handle both cases.

```rust
use std::fs;
use std::io;

fn read_file(path: &str) -> Result<String, io::Error> {
    fs::read_to_string(path)  // returns Result<String, io::Error>
}

fn main() {
    // Method 1: match
    match read_file("hello.txt") {
        Ok(contents) => println!("File: {}", contents),
        Err(e) => println!("Error: {}", e),
    }

    // Method 2: if let
    if let Ok(contents) = read_file("hello.txt") {
        println!("File: {}", contents);
    }

    // Method 3: unwrap (panics on Err — use only in tests/prototypes)
    let contents = read_file("hello.txt").unwrap();

    // Method 4: expect (panics with a message)
    let contents = read_file("hello.txt").expect("Failed to read hello.txt");

    // Method 5: unwrap_or (default value)
    let contents = read_file("hello.txt").unwrap_or_default();

    // Method 6: unwrap_or_else (compute default)
    let contents = read_file("hello.txt").unwrap_or_else(|_| String::from("default"));
}
```

---

## 10.4 The ? Operator — Propagating Errors

The `?` operator is Rust's primary tool for error propagation. It's the most important operator in error handling:

```rust
use std::fs;
use std::io;

// WITHOUT ? operator — verbose
fn read_username_from_file_verbose() -> Result<String, io::Error> {
    let f = fs::File::open("hello.txt");
    let mut f = match f {
        Ok(file) => file,
        Err(e) => return Err(e),  // propagate error up
    };

    let mut s = String::new();
    match f.read_to_string(&mut s) {
        Ok(_) => Ok(s),
        Err(e) => Err(e),  // propagate error up
    }
}

// WITH ? operator — concise
fn read_username_from_file() -> Result<String, io::Error> {
    let mut s = String::new();
    fs::File::open("hello.txt")?.read_to_string(&mut s)?;
    Ok(s)
}

// EVEN SHORTER — using the stdlib shortcut
fn read_username_from_file_short() -> Result<String, io::Error> {
    fs::read_to_string("hello.txt")
}
```

### How ? Works

`?` placed after a `Result` expression:
1. If `Ok(value)`: unwraps to `value`, execution continues
2. If `Err(e)`: converts the error (via `From` trait) and **immediately returns** `Err(converted_e)` from the current function

```rust
let x = some_result?;
// Expands roughly to:
let x = match some_result {
    Ok(v) => v,
    Err(e) => return Err(e.into()),  // .into() allows error type conversion
};
```

### ? Also Works with Option

```rust
fn get_first_char(s: &str) -> Option<char> {
    let first_word = s.split_whitespace().next()?;  // None if empty string
    first_word.chars().next()  // None if word is empty (shouldn't happen)
}

fn main() {
    println!("{:?}", get_first_char("hello world"));  // Some('h')
    println!("{:?}", get_first_char(""));              // None
}
```

### Chaining with ?

```rust
use std::num::ParseIntError;

fn double_first(vec: Vec<&str>) -> Result<i32, ParseIntError> {
    let first = vec.first().ok_or("".parse::<i32>().unwrap_err())?;
    let parsed = first.parse::<i32>()?;
    Ok(2 * parsed)
}

// More idiomatic with explicit error types:
fn parse_and_double(s: &str) -> Result<i32, ParseIntError> {
    let n = s.trim().parse::<i32>()?;
    Ok(n * 2)
}

fn main() {
    println!("{:?}", parse_and_double("21"));    // Ok(42)
    println!("{:?}", parse_and_double("abc"));   // Err(...)
    println!("{:?}", parse_and_double("  10 ")); // Ok(20) — trim handles whitespace
}
```

---

## 10.5 Custom Error Types

For real applications, define your own error types:

### Simple String Errors

```rust
fn find_user(id: u32) -> Result<String, String> {
    if id == 0 {
        Err(String::from("ID cannot be zero"))
    } else if id > 1000 {
        Err(format!("ID {} is out of range", id))
    } else {
        Ok(format!("User_{}", id))
    }
}
```

### Proper Error Enum

```rust
use std::fmt;
use std::num::ParseIntError;

#[derive(Debug)]
enum AppError {
    NotFound(String),
    ParseError(ParseIntError),
    InvalidInput { field: String, reason: String },
    IoError(std::io::Error),
}

// Display for human-readable error messages
impl fmt::Display for AppError {
    fn fmt(&self, f: &mut fmt::Formatter<'_>) -> fmt::Result {
        match self {
            AppError::NotFound(msg) => write!(f, "Not found: {}", msg),
            AppError::ParseError(e) => write!(f, "Parse error: {}", e),
            AppError::InvalidInput { field, reason } => {
                write!(f, "Invalid {}: {}", field, reason)
            }
            AppError::IoError(e) => write!(f, "IO error: {}", e),
        }
    }
}

// std::error::Error trait — enables interoperability
impl std::error::Error for AppError {
    fn source(&self) -> Option<&(dyn std::error::Error + 'static)> {
        match self {
            AppError::ParseError(e) => Some(e),
            AppError::IoError(e) => Some(e),
            _ => None,
        }
    }
}

// From implementations — enable ? operator to convert errors automatically
impl From<ParseIntError> for AppError {
    fn from(e: ParseIntError) -> Self {
        AppError::ParseError(e)
    }
}

impl From<std::io::Error> for AppError {
    fn from(e: std::io::Error) -> Self {
        AppError::IoError(e)
    }
}

fn process_id(s: &str) -> Result<String, AppError> {
    let id: u32 = s.trim().parse()?;  // ParseIntError → AppError via From
    if id == 0 {
        return Err(AppError::InvalidInput {
            field: "id".to_string(),
            reason: "must be non-zero".to_string(),
        });
    }
    Ok(format!("user_{}", id))
}

fn main() {
    println!("{:?}", process_id("42"));   // Ok("user_42")
    println!("{:?}", process_id("abc"));  // Err(ParseError(...))
    println!("{:?}", process_id("0"));    // Err(InvalidInput {...})

    // Display
    match process_id("abc") {
        Ok(u) => println!("{}", u),
        Err(e) => eprintln!("Error: {}", e),
    }
}
```

---

## 10.6 Using Box<dyn Error>

For quick programs where you don't want to define custom error types, use trait objects:

```rust
use std::error::Error;
use std::fs;

fn main() -> Result<(), Box<dyn Error>> {
    let content = fs::read_to_string("config.toml")?;   // io::Error
    let n: i32 = content.trim().parse()?;                // ParseIntError
    println!("Config value: {}", n);
    Ok(())
}
```

`Box<dyn Error>` can hold any error type — convenient but loses static type information. Fine for `main()` and CLI tools.

---

## 10.7 The anyhow and thiserror Crates

In real Rust projects, two crates simplify error handling:

### thiserror — Custom Error Types Without Boilerplate

```rust
// In Cargo.toml: thiserror = "1"

use thiserror::Error;

#[derive(Debug, Error)]
enum AppError {
    #[error("Not found: {0}")]
    NotFound(String),

    #[error("Parse error")]
    Parse(#[from] std::num::ParseIntError),

    #[error("IO error")]
    Io(#[from] std::io::Error),

    #[error("Invalid {field}: {reason}")]
    InvalidInput { field: String, reason: String },
}
```

### anyhow — Easy Error Propagation in Applications

```rust
// In Cargo.toml: anyhow = "1"

use anyhow::{Context, Result};

fn read_config() -> Result<String> {
    let content = std::fs::read_to_string("config.txt")
        .context("Failed to read config file")?;
    Ok(content)
}

fn main() -> Result<()> {
    let config = read_config()?;
    println!("{}", config);
    Ok(())
}
```

---

## 10.8 Error Handling Patterns

### Early Return Pattern

```rust
fn process(input: &str) -> Result<i32, AppError> {
    if input.is_empty() {
        return Err(AppError::InvalidInput {
            field: "input".to_string(),
            reason: "cannot be empty".to_string(),
        });
    }

    let n: i32 = input.parse()?;

    if n < 0 {
        return Err(AppError::InvalidInput {
            field: "n".to_string(),
            reason: "must be non-negative".to_string(),
        });
    }

    Ok(n * 2)
}
```

### Collecting Results

```rust
fn main() {
    let strings = vec!["1", "2", "three", "4"];

    // Collect stops at first error
    let results: Result<Vec<i32>, _> = strings.iter()
        .map(|s| s.parse::<i32>())
        .collect();
    println!("{:?}", results);  // Err(ParseIntError { ... })

    // Collect all, keep only successes
    let numbers: Vec<i32> = strings.iter()
        .filter_map(|s| s.parse::<i32>().ok())
        .collect();
    println!("{:?}", numbers);  // [1, 2, 4]

    // Collect all, keep all Results
    let all: Vec<Result<i32, _>> = strings.iter()
        .map(|s| s.parse::<i32>())
        .collect();
    for (s, r) in strings.iter().zip(all.iter()) {
        println!("{} → {:?}", s, r);
    }
}
```

### Converting Option to Result

```rust
fn get_env_var(name: &str) -> Result<String, String> {
    std::env::var(name).ok()                              // Option<String>
        .ok_or_else(|| format!("{} is not set", name))   // Result<String, String>
}

fn find_config(name: &str) -> Result<String, AppError> {
    vec!["config.json", "config.toml"].into_iter()
        .find(|&f| f.starts_with(name))
        .map(|s| s.to_string())
        .ok_or_else(|| AppError::NotFound(name.to_string()))
}
```

---

## Complete Example: A Robust Configuration Parser

```rust
use std::collections::HashMap;
use std::fmt;

#[derive(Debug)]
enum ConfigError {
    MissingKey(String),
    InvalidValue { key: String, expected: String, got: String },
    ParseError(String),
}

impl fmt::Display for ConfigError {
    fn fmt(&self, f: &mut fmt::Formatter<'_>) -> fmt::Result {
        match self {
            ConfigError::MissingKey(k) => write!(f, "Missing required key: {}", k),
            ConfigError::InvalidValue { key, expected, got } => {
                write!(f, "Key '{}': expected {}, got '{}'", key, expected, got)
            }
            ConfigError::ParseError(s) => write!(f, "Parse error: {}", s),
        }
    }
}

impl std::error::Error for ConfigError {}

fn parse_config(input: &str) -> Result<HashMap<String, String>, ConfigError> {
    let mut map = HashMap::new();
    for line in input.lines() {
        let line = line.trim();
        if line.is_empty() || line.starts_with('#') {
            continue;
        }
        let parts: Vec<&str> = line.splitn(2, '=').collect();
        if parts.len() != 2 {
            return Err(ConfigError::ParseError(format!("Invalid line: {}", line)));
        }
        map.insert(parts[0].trim().to_string(), parts[1].trim().to_string());
    }
    Ok(map)
}

fn get_port(config: &HashMap<String, String>) -> Result<u16, ConfigError> {
    let s = config.get("port").ok_or_else(|| ConfigError::MissingKey("port".to_string()))?;
    s.parse::<u16>().map_err(|_| ConfigError::InvalidValue {
        key: "port".to_string(),
        expected: "0-65535".to_string(),
        got: s.clone(),
    })
}

fn main() {
    let config_str = "
        # Server configuration
        host = localhost
        port = 8080
        debug = true
    ";

    match parse_config(config_str) {
        Ok(config) => {
            println!("Config loaded: {} keys", config.len());
            match get_port(&config) {
                Ok(port) => println!("Port: {}", port),
                Err(e) => eprintln!("Port error: {}", e),
            }
        }
        Err(e) => eprintln!("Config error: {}", e),
    }
}
```

---

## Summary

Rust has two error kinds: `panic!` for unrecoverable programming bugs and `Result<T, E>` for recoverable failures. The `?` operator propagates errors up the call stack automatically. Custom error types should implement `Display`, `std::error::Error`, and `From` for automatic conversions. `Box<dyn Error>` is a quick escape hatch when you don't need specific error types. The `thiserror` crate reduces custom error boilerplate; `anyhow` simplifies error propagation in applications.

---

## Key Takeaways

- `panic!` = bug. `Result` = expected failure. Never use `panic!` for user-facing errors
- `?` is the primary tool: propagates `Err` up, unwraps `Ok` — cleaner than `match`
- `?` works on both `Result` and `Option` (in functions returning `Option`)
- Custom error types need: `Debug`, `Display`, `std::error::Error`, and `From<X>` impls
- `From` impls enable automatic `?` conversion between error types
- `collect::<Result<Vec<_>, _>>()` collects iterator results, stopping at first error
- Use `ok()` to convert `Result` → `Option`, `ok_or(e)` to convert `Option` → `Result`

---

## Exercises

**Exercise 1:** Write a `parse_date(s: &str) -> Result<(u32, u32, u32), String>` function that parses "YYYY-MM-DD" format and returns `Err` with a helpful message for invalid inputs.

**Exercise 2:** Create a `VendingMachine` struct. Its `purchase(item: &str, payment: f64) -> Result<f64, VendingError>` method should return `Ok(change)` or `Err(reason)`. Define a `VendingError` enum.

**Exercise 3:** Use the `?` operator to write a function that reads a file containing a number, doubles it, and writes it back. Handle all three operations with proper error propagation.

**Exercise 4:** Implement `fn try_all(ops: Vec<fn() -> Result<i32, String>>) -> Result<Vec<i32>, String>` that runs all operations and returns `Ok` only if all succeed.

**Exercise 5:** Write a small CSV parser that returns a custom `ParseError` enum with variants for empty input, missing columns, and invalid data types.

---

*Next: [Chapter 11 — Generics and Traits](11-generics-traits.md)*
