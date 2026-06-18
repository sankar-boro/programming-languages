# Chapter 15 — Modules, Crates, and Cargo

---

## 15.1 The Module System

Rust's module system organizes code into namespaced units. Modules control visibility, enable code splitting, and prevent naming conflicts.

### Declaring Modules

```rust
// src/main.rs or src/lib.rs

mod garden {
    pub mod vegetables {
        pub struct Asparagus;
    }
}

use crate::garden::vegetables::Asparagus;

fn main() {
    let plant = Asparagus;
}
```

### Module Hierarchy

```
crate (root)
├── front_of_house
│   ├── hosting
│   │   ├── add_to_waitlist
│   │   └── seat_at_table
│   └── serving
│       ├── take_order
│       ├── serve_order
│       └── take_payment
└── back_of_house
    ├── fix_incorrect_order
    └── cook_order
```

```rust
mod front_of_house {
    pub mod hosting {
        pub fn add_to_waitlist() {
            println!("Added to waitlist");
        }

        fn seat_at_table() {}  // private
    }

    mod serving {
        fn take_order() {}
        fn serve_order() {}
    }
}

pub fn eat_at_restaurant() {
    // Absolute path
    crate::front_of_house::hosting::add_to_waitlist();

    // Relative path
    front_of_house::hosting::add_to_waitlist();
}
```

---

## 15.2 Visibility — pub

Everything in Rust is **private by default**. Use `pub` to make items public.

```rust
mod outer {
    pub fn public_fn() {}
    fn private_fn() {}  // only accessible within this module and children

    pub mod inner {
        pub fn inner_public() {
            // Can access parent's private items
            super::private_fn();
        }

        fn inner_private() {}
    }

    pub struct PublicStruct {
        pub public_field: i32,
        private_field: String,  // fields are private by default
    }

    impl PublicStruct {
        pub fn new(value: i32) -> Self {
            PublicStruct {
                public_field: value,
                private_field: String::from("secret"),
            }
        }
    }

    pub enum PublicEnum {
        Variant1,  // enum variants are public if enum is public
        Variant2(i32),
    }
}

fn main() {
    outer::public_fn();
    // outer::private_fn();  // ERROR — private

    let s = outer::PublicStruct::new(5);
    println!("{}", s.public_field);
    // println!("{}", s.private_field);  // ERROR — private

    let e = outer::PublicEnum::Variant1;  // OK — variants are public
}
```

### Visibility Levels

```rust
pub fn visible_everywhere() {}
pub(crate) fn visible_in_crate() {}      // only within this crate
pub(super) fn visible_in_parent() {}     // only in parent module
pub(in crate::some::path) fn specific() {}  // only in specific path
fn private() {}                           // only in this module
```

---

## 15.3 use — Bringing Paths Into Scope

```rust
use std::collections::HashMap;
use std::io::{self, Write};        // import io and io::Write
use std::io::*;                    // glob import (use sparingly)

// Renaming with as
use std::collections::HashMap as Map;
use std::fmt::Result as FmtResult;

fn main() {
    let mut m: HashMap<&str, i32> = HashMap::new();
    m.insert("key", 42);

    // io and io::Write both available
    let stdout = io::stdout();
    let _ = stdout.lock();
}
```

### Re-exporting with pub use

```rust
// lib.rs — create a clean public API
mod internal {
    pub struct Implementation;
    impl Implementation {
        pub fn do_work(&self) {}
    }
}

// Re-export — users see lib::Implementation, not lib::internal::Implementation
pub use internal::Implementation;
```

---

## 15.4 Module Files — Multi-File Projects

For larger projects, modules live in separate files:

### File Structure

```
src/
├── main.rs
├── lib.rs
├── config.rs          ← mod config; in main.rs
├── network/
│   ├── mod.rs         ← mod network; in main.rs (old style)
│   ├── server.rs      ← mod server; in network/mod.rs
│   └── client.rs
└── database.rs
```

**Modern style (Rust 2018+):**

```
src/
├── main.rs
├── config.rs          ← module file
├── network.rs         ← module file (replaces network/mod.rs)
├── network/
│   ├── server.rs      ← submodule
│   └── client.rs
└── database.rs
```

```rust
// src/main.rs
mod config;     // looks for src/config.rs
mod network;    // looks for src/network.rs
mod database;

use config::Config;
use network::server::Server;

fn main() {
    let config = Config::new();
    let server = Server::start(&config);
}
```

```rust
// src/network.rs
pub mod server;   // looks for src/network/server.rs
pub mod client;   // looks for src/network/client.rs

pub use server::Server;  // re-export for convenience
```

```rust
// src/network/server.rs
use crate::config::Config;  // absolute path from crate root

pub struct Server {
    port: u16,
}

impl Server {
    pub fn start(config: &Config) -> Self {
        Server { port: config.port }
    }
}
```

---

## 15.5 Crates — Packages and Libraries

**Package**: a Cargo project (one `Cargo.toml`). Can contain multiple crates.  
**Crate**: a compilation unit (one library or binary). Named by the package name.  
**Library crate**: `src/lib.rs` — the core of most libraries  
**Binary crate**: `src/main.rs` — an executable

### Adding Dependencies

```toml
# Cargo.toml
[package]
name = "my_app"
version = "0.1.0"
edition = "2021"

[dependencies]
serde = { version = "1.0", features = ["derive"] }
tokio = { version = "1", features = ["full"] }
reqwest = "0.12"
anyhow = "1"
clap = { version = "4", features = ["derive"] }

[dev-dependencies]       # only for tests and benches
criterion = "0.5"
pretty_assertions = "1"

[build-dependencies]     # for build scripts (build.rs)
cc = "1"
```

### Version Specification

```toml
dep = "1.2.3"      # exactly 1.2.3 (rare)
dep = "^1.2.3"     # >=1.2.3, <2.0.0 (default when you write "1.2.3")
dep = "~1.2.3"     # >=1.2.3, <1.3.0
dep = ">=1.2, <2"  # explicit range
dep = "*"          # any version (dangerous)
```

### Common Ecosystem Crates

```toml
[dependencies]
# Serialization
serde = { version = "1", features = ["derive"] }
serde_json = "1"

# Async runtime
tokio = { version = "1", features = ["full"] }

# HTTP client
reqwest = "0.12"

# Error handling
anyhow = "1"
thiserror = "1"

# CLI
clap = { version = "4", features = ["derive"] }

# Logging
tracing = "0.1"
tracing-subscriber = "0.3"

# Random numbers
rand = "0.8"

# Date/time
chrono = "0.4"

# Regex
regex = "1"

# Uuid
uuid = { version = "1", features = ["v4"] }
```

---

## 15.6 Cargo Features

Features are optional functionality you can enable in dependencies:

```toml
[features]
default = ["networking"]     # enabled by default
networking = ["dep:reqwest"] # enable reqwest only when networking is enabled
async-support = ["dep:tokio"]

[dependencies]
reqwest = { version = "0.12", optional = true }
tokio = { version = "1", optional = true }
```

```bash
# Enable features
cargo build --features networking,async-support
cargo build --all-features
cargo build --no-default-features
```

---

## 15.7 Cargo Workspaces

Workspaces allow multiple related crates to share a `Cargo.lock` and build output:

```
my-workspace/
├── Cargo.toml          ← workspace manifest
├── Cargo.lock          ← shared lock file
├── core/               ← library crate
│   ├── Cargo.toml
│   └── src/lib.rs
├── cli/                ← binary that uses core
│   ├── Cargo.toml
│   └── src/main.rs
└── server/             ← another binary
    ├── Cargo.toml
    └── src/main.rs
```

```toml
# Root Cargo.toml
[workspace]
members = ["core", "cli", "server"]
resolver = "2"
```

```toml
# cli/Cargo.toml
[dependencies]
core = { path = "../core" }  # reference sibling crate by path
```

```bash
# Build everything
cargo build

# Run a specific binary
cargo run -p cli

# Test a specific crate
cargo test -p core
```

---

## 15.8 Publishing to crates.io

```bash
# 1. Create an account at crates.io and get a token
cargo login <your-token>

# 2. Ensure Cargo.toml has required fields
# [package]
# name = "my-awesome-crate"
# version = "0.1.0"
# edition = "2021"
# description = "Does awesome things"
# license = "MIT OR Apache-2.0"
# homepage = "https://example.com"
# repository = "https://github.com/you/my-awesome-crate"
# documentation = "https://docs.rs/my-awesome-crate"

# 3. Dry run — check without publishing
cargo publish --dry-run

# 4. Publish
cargo publish
```

---

## 15.9 Documentation

Rust has a built-in documentation system:

```rust
//! Module-level documentation (inner doc comment)
//! This module provides utilities for parsing configuration files.

/// Parses a configuration string.
///
/// # Arguments
///
/// * `input` - The configuration string in KEY=VALUE format
///
/// # Returns
///
/// A `HashMap` mapping keys to values.
///
/// # Errors
///
/// Returns `ParseError` if the input is malformed.
///
/// # Examples
///
/// ```
/// let config = my_crate::parse_config("HOST=localhost\nPORT=8080").unwrap();
/// assert_eq!(config.get("HOST"), Some(&"localhost".to_string()));
/// ```
pub fn parse_config(input: &str) -> Result<HashMap<String, String>, ParseError> {
    todo!()
}
```

```bash
cargo doc         # generate docs
cargo doc --open  # generate and open in browser
cargo test --doc  # run doc tests (code in /// ``` blocks)
```

---

## Complete Example: Organizing a Library

```
mylib/
├── Cargo.toml
└── src/
    ├── lib.rs
    ├── parser.rs
    ├── validator.rs
    └── error.rs
```

```rust
// src/error.rs
use std::fmt;

#[derive(Debug)]
pub enum Error {
    ParseError(String),
    ValidationError(String),
}

impl fmt::Display for Error {
    fn fmt(&self, f: &mut fmt::Formatter<'_>) -> fmt::Result {
        match self {
            Error::ParseError(s) => write!(f, "Parse error: {}", s),
            Error::ValidationError(s) => write!(f, "Validation error: {}", s),
        }
    }
}

impl std::error::Error for Error {}
```

```rust
// src/lib.rs
mod error;
mod parser;
mod validator;

// Clean public API
pub use error::Error;
pub use parser::parse;
pub use validator::validate;

pub type Result<T> = std::result::Result<T, Error>;

pub fn process(input: &str) -> Result<String> {
    let parsed = parse(input)?;
    validate(&parsed)?;
    Ok(parsed.to_uppercase())
}
```

---

## Summary

Rust's module system organizes code into namespaced, visibility-controlled units. Items are private by default — `pub` makes them public. `use` brings paths into scope; `pub use` re-exports for clean APIs. For multi-file projects, `mod foo;` looks for `src/foo.rs` or `src/foo/mod.rs`. Crates are compilation units; packages are Cargo projects. Cargo manages dependencies in `Cargo.toml` with semantic versioning. Workspaces unite related crates under a shared build environment. Documentation is built-in with `///` doc comments and `cargo doc`.

---

## Key Takeaways

- Everything private by default — explicitly `pub` what should be public
- `use crate::...` for absolute paths; `use super::...` for relative parent paths
- Module files: `src/foo.rs` OR `src/foo/mod.rs` — not both
- `pub use` re-exports flatten your module hierarchy for users
- `cargo add <crate>` is the modern way to add dependencies
- Cargo.lock: commit for binaries; gitignore for libraries
- `///` doc comments appear in `cargo doc`; doc tests in ` ``` ` blocks run with `cargo test --doc`

---

## Exercises

**Exercise 1:** Create a library with three modules: `shapes`, `colors`, and `canvas`. Each module has its own file. Export a clean public API from `lib.rs`.

**Exercise 2:** Add `serde` with the `derive` feature to a project. Create a struct with `#[derive(Serialize, Deserialize)]` and serialize it to JSON.

**Exercise 3:** Create a workspace with two crates: `math-core` (library with arithmetic functions) and `math-cli` (binary that uses the library).

**Exercise 4:** Write documentation for a function including a `# Examples` section with runnable code. Verify the doc test passes with `cargo test --doc`.

**Exercise 5:** Use `pub(crate)` to make an implementation detail accessible within the crate but not from external code. Write a test that verifies external code cannot access it.

---

*Next: [Chapter 16 — Rust Internals](16-internals.md)*
