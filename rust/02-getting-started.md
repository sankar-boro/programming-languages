# Chapter 2 — The Rust Toolchain

---

## 2.1 Installing Rust with rustup

`rustup` is the official Rust version manager — think of it like `nvm` for Node or `pyenv` for Python. It manages Rust compiler versions, targets, and components.

### Installation

```bash
# Linux / macOS
curl --proto '=https' --tlsv1.2 -sSf https://sh.rustup.rs | sh

# Follow the on-screen prompts, then:
source ~/.cargo/env

# Windows — download and run rustup-init.exe from rustup.rs
```

### Verify Installation

```bash
rustc --version
# rustc 1.78.0 (9b00956e5 2024-04-29)

cargo --version
# cargo 1.78.0 (54d8815d0 2024-03-26)

rustup --version
# rustup 1.27.0 (bbb9276d2 2024-02-21)
```

### Managing Versions

```bash
rustup update              # update to latest stable
rustup show                # show installed toolchains
rustup toolchain list      # list all toolchains

# Install specific versions
rustup toolchain install 1.70.0
rustup default 1.70.0      # set as default

# Channels
rustup toolchain install nightly  # bleeding edge
rustup toolchain install beta     # pre-release testing
```

### Components

```bash
# Add useful components
rustup component add rustfmt      # code formatter
rustup component add clippy       # linter
rustup component add rust-analyzer # language server for IDEs

# Add cross-compilation targets
rustup target add wasm32-unknown-unknown  # WebAssembly
rustup target add aarch64-apple-darwin   # Apple Silicon
```

---

## 2.2 cargo — The Build Tool

`cargo` is Rust's package manager and build tool. It is significantly better than most language build tools — no Makefile, no Maven, no Gradle. One tool does everything.

### Creating a Project

```bash
cargo new hello_world        # creates a binary project
cargo new mylib --lib        # creates a library project

# Project structure created:
# hello_world/
# ├── Cargo.toml     ← manifest (like package.json)
# └── src/
#     └── main.rs    ← entry point
```

### Cargo.toml — The Manifest

```toml
[package]
name = "hello_world"
version = "0.1.0"
edition = "2021"          # Rust edition (2015, 2018, 2021)
authors = ["Your Name <you@example.com>"]
description = "A hello world program"

[dependencies]
serde = { version = "1.0", features = ["derive"] }
tokio = "1.0"

[dev-dependencies]         # only for tests and benchmarks
pretty_assertions = "1.0"

[profile.release]          # optimization settings for release builds
opt-level = 3
lto = true
```

### Common cargo Commands

```bash
cargo build              # compile in debug mode
cargo build --release    # compile with optimizations
cargo run                # build and run
cargo run -- arg1 arg2   # build and run with arguments
cargo test               # run all tests
cargo check              # check for errors WITHOUT compiling (fast)
cargo clippy             # run linter
cargo fmt                # format code
cargo doc                # generate documentation
cargo doc --open         # generate and open in browser
cargo clean              # remove build artifacts

# Dependencies
cargo add serde          # add a dependency
cargo add serde@1.0.100  # specific version
cargo remove serde       # remove a dependency
cargo update             # update all dependencies
cargo tree               # show dependency tree
```

### The Cargo.lock File

```toml
# Cargo.lock — auto-generated, exact dependency versions
# Commit this for binaries (reproducible builds)
# Do NOT commit for libraries (let users pick compatible versions)

[[package]]
name = "hello_world"
version = "0.1.0"
```

### Workspaces — Multi-Package Projects

```toml
# Root Cargo.toml
[workspace]
members = [
    "core",
    "cli",
    "server",
]
```

---

## 2.3 rustc — The Compiler

`rustc` is the Rust compiler. In practice, you almost never call it directly — `cargo` calls it for you. But knowing it exists helps:

```bash
# Compile a single file
rustc main.rs

# With optimizations
rustc -O main.rs

# Target a specific platform
rustc --target wasm32-unknown-unknown main.rs

# Show what the compiler produces
rustc --emit=asm main.rs     # assembly output
rustc --emit=mir main.rs     # mid-level IR
rustc --emit=llvm-ir main.rs # LLVM IR
```

### The Compilation Pipeline

```
[.rs source files]
      │
      ▼
  [Parsing]  →  Abstract Syntax Tree (AST)
      │
      ▼
  [Name Resolution + Type Checking]
      │
      ▼
  [Borrow Checking]  ← this is where ownership rules are enforced
      │
      ▼
  [MIR]  (Mid-level Intermediate Representation)
      │
      ▼
  [LLVM IR]  (handed to LLVM backend)
      │
      ▼
  [Native machine code / WASM / etc.]
```

---

## 2.4 Project Structure

### Binary Project

```
my_app/
├── Cargo.toml
├── Cargo.lock
├── src/
│   ├── main.rs          ← binary entry point
│   ├── lib.rs           ← library root (optional)
│   └── modules/
│       ├── mod.rs       ← or just module.rs in newer Rust
│       └── helper.rs
├── tests/
│   └── integration_test.rs  ← integration tests
├── benches/
│   └── benchmark.rs     ← benchmarks
├── examples/
│   └── example.rs       ← runnable examples
└── target/              ← build artifacts (gitignored)
```

### Library Project

```
my_lib/
├── Cargo.toml
└── src/
    └── lib.rs           ← library root (no main function)
```

---

## 2.5 Hello, World!

### The Program

```rust
// src/main.rs
fn main() {
    println!("Hello, World!");
}
```

### Run It

```bash
cargo run
# Compiling hello_world v0.1.0
# Finished dev [unoptimized + debuginfo] target(s) in 0.42s
# Running `target/debug/hello_world`
# Hello, World!
```

### What Is println!?

The `!` after `println` means it's a **macro**, not a regular function. Macros in Rust are powerful metaprogramming tools. `println!` specifically is a macro because it handles format strings at compile time — it verifies the format string matches the arguments you provide.

```rust
fn main() {
    // Basic output
    println!("Hello, World!");

    // Formatted output
    let name = "Rust";
    let year = 2015;
    println!("Hello from {} (since {})!", name, year);

    // Debug format with {:?}
    let numbers = vec![1, 2, 3];
    println!("Numbers: {:?}", numbers);

    // Pretty debug with {:#?}
    println!("Numbers pretty:\n{:#?}", numbers);

    // Positional arguments
    println!("{0} is {1} and {0} is great", "Rust", "fast");

    // Named arguments
    println!("{lang} was created in {year}", lang = "Rust", year = 2006);

    // print! without newline
    print!("no newline here");
    print!(" — continuing on same line\n");

    // eprintln! for stderr
    eprintln!("This goes to stderr");
}
```

### A More Complete First Program

```rust
use std::io;
use std::io::Write;

fn main() {
    // Print without newline, then flush to ensure it shows
    print!("Enter your name: ");
    io::stdout().flush().unwrap();

    // Read a line from stdin
    let mut name = String::new();
    io::stdin()
        .read_line(&mut name)
        .expect("Failed to read line");

    // trim() removes the newline at the end
    let name = name.trim();

    println!("Hello, {}! Welcome to Rust.", name);

    // Basic arithmetic
    let x: i32 = 10;
    let y: i32 = 3;
    println!("{} + {} = {}", x, y, x + y);
    println!("{} / {} = {}", x, y, x / y);       // integer division
    println!("{} % {} = {}", x, y, x % y);       // remainder

    // Floating point
    let a: f64 = 10.0;
    let b: f64 = 3.0;
    println!("{:.4}", a / b);                     // 3.3333
}
```

### Running Tests

```rust
// src/main.rs or src/lib.rs
fn add(a: i32, b: i32) -> i32 {
    a + b
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_add() {
        assert_eq!(add(2, 3), 5);
    }

    #[test]
    fn test_add_negative() {
        assert_eq!(add(-1, 1), 0);
    }

    #[test]
    #[should_panic]
    fn test_panic() {
        panic!("this test expects a panic");
    }
}
```

```bash
cargo test
# running 3 tests
# test tests::test_add ... ok
# test tests::test_add_negative ... ok
# test tests::test_panic ... ok
# test result: ok. 3 passed; 0 failed
```

---

## Summary

Rust's toolchain is centered on three tools: `rustup` (version manager), `cargo` (build tool + package manager), and `rustc` (compiler). `cargo new` sets up a complete project structure. `cargo run` compiles and executes. `cargo test` runs the test suite. `cargo check` quickly validates code without producing a binary. The Rust toolchain is one of the best in any language ecosystem.

---

## Key Takeaways

- Install Rust via `rustup` — never download rustc directly
- `cargo` does everything: build, run, test, format, lint, manage dependencies
- `cargo check` is faster than `cargo build` — use it for rapid feedback
- `println!` is a macro (note the `!`), not a function
- `Cargo.toml` is the project manifest; `Cargo.lock` pins exact versions
- Commit `Cargo.lock` for binaries; don't commit it for libraries
- `cargo clippy` and `cargo fmt` should be part of every workflow

---

## Exercises

**Exercise 1:** Create a new Rust project called `calculator`. Add a function `multiply(a: i32, b: i32) -> i32` and write three tests for it.

**Exercise 2:** Modify `main.rs` to read two numbers from stdin and print their sum, difference, product, and quotient.

**Exercise 3:** Run `cargo build --release` and compare the binary size in `target/debug/` vs `target/release/`. What do you notice?

**Exercise 4:** Run `cargo doc --open` on any project. Explore the auto-generated documentation.

---

*Next: [Chapter 3 — Variables, Types, and Expressions](03-basic-syntax.md)*
