# Chapter 90 — Best Practices and Idiomatic Rust

---

## 1. Use &str Instead of &String in Function Parameters

```rust
// BAD — only accepts String
fn greet(name: &String) { println!("Hello, {}", name); }

// GOOD — accepts String, &str, and more via deref coercion
fn greet(name: &str) { println!("Hello, {}", name); }

fn main() {
    greet("literal");                        // &str
    greet(&String::from("owned"));           // String via deref
    greet(&String::from("owned")[0..3]);     // slice
}
```

## 2. Use &[T] Instead of &Vec<T> in Function Parameters

```rust
// BAD
fn sum(v: &Vec<i32>) -> i32 { v.iter().sum() }

// GOOD — more general
fn sum(v: &[i32]) -> i32 { v.iter().sum() }

fn main() {
    sum(&[1, 2, 3]);              // array
    sum(&vec![1, 2, 3]);          // Vec
    sum(&vec![1, 2, 3, 4, 5][1..3]); // slice
}
```

## 3. Prefer Iterators Over Index Loops

```rust
let v = vec![1, 2, 3, 4, 5];

// BAD — manual indexing, panic risk
for i in 0..v.len() {
    println!("{}", v[i]);
}

// GOOD — idiomatic, safe, expressive
for x in &v {
    println!("{}", x);
}

// GOOD — with index
for (i, x) in v.iter().enumerate() {
    println!("[{}] = {}", i, x);
}
```

## 4. Use the Entry API for HashMap

```rust
use std::collections::HashMap;

let mut map: HashMap<String, Vec<i32>> = HashMap::new();

// BAD — double lookup
if !map.contains_key("key") {
    map.insert("key".to_string(), Vec::new());
}
map.get_mut("key").unwrap().push(1);

// GOOD — single lookup, idiomatic
map.entry("key".to_string()).or_insert_with(Vec::new).push(1);
```

## 5. Return impl Trait Instead of Box<dyn Trait> When Possible

```rust
// BAD — heap allocation, slower
fn make_greeting(name: &str) -> Box<dyn Fn() -> String> {
    let name = name.to_string();
    Box::new(move || format!("Hello, {}!", name))
}

// GOOD — zero overhead when type is known at compile time
fn make_greeting(name: &str) -> impl Fn() -> String {
    let name = name.to_string();
    move || format!("Hello, {}!", name)
}
```

## 6. Use Option and Result Methods Instead of Explicit match

```rust
let v: Option<i32> = Some(42);

// BAD — verbose
let doubled = match v {
    Some(x) => Some(x * 2),
    None => None,
};

// GOOD — idiomatic
let doubled = v.map(|x| x * 2);

// GOOD chaining
let result: Option<String> = v
    .filter(|&x| x > 0)
    .map(|x| x.to_string());
```

## 7. Use ? Early and Often

```rust
// BAD — nested match hell
fn process(s: &str) -> Result<i32, String> {
    match s.parse::<i32>() {
        Ok(n) => match check_range(n) {
            Ok(_) => Ok(n * 2),
            Err(e) => Err(e),
        },
        Err(e) => Err(e.to_string()),
    }
}

// GOOD — flat and readable
fn process(s: &str) -> Result<i32, String> {
    let n: i32 = s.parse().map_err(|e: std::num::ParseIntError| e.to_string())?;
    check_range(n)?;
    Ok(n * 2)
}
```

## 8. Derive Common Traits

```rust
// Always derive Debug for any type you create
// Derive Clone, Copy, PartialEq, Eq, Hash, Default when it makes sense

#[derive(Debug, Clone, PartialEq, Eq, Hash)]
struct UserId(u32);

#[derive(Debug, Clone, Default)]
struct Config {
    host: String,
    port: u16,
}
```

## 9. Use Newtype Pattern for Type Safety

```rust
// BAD — easy to mix up meters and seconds
fn travel_time(distance: f64, speed: f64) -> f64 {
    distance / speed
}

// GOOD — distinct types prevent mix-ups
struct Meters(f64);
struct MetersPerSecond(f64);
struct Seconds(f64);

fn travel_time(distance: Meters, speed: MetersPerSecond) -> Seconds {
    Seconds(distance.0 / speed.0)
}
```

## 10. Prefer Enums Over Bool Parameters

```rust
// BAD — what does `true` mean?
fn connect(host: &str, secure: bool) {}
connect("example.com", true);

// GOOD — self-documenting
enum Security { Secure, Insecure }
fn connect(host: &str, security: Security) {}
connect("example.com", Security::Secure);
```

## 11. Use Clippy and rustfmt

```bash
# Run linter — catches many anti-patterns
cargo clippy

# Format code
cargo fmt

# Pedantic mode — stricter suggestions
cargo clippy -- -W clippy::pedantic

# Add to CI
cargo clippy -- -D warnings  # fail on any warning
```

## 12. Avoid Unwrap in Production Code

```rust
// BAD — panics if None (in library/production code)
let val = some_option.unwrap();
let result = some_result.unwrap();

// GOOD — handle the error
let val = some_option?;
let val = some_option.unwrap_or_default();
let val = some_option.ok_or(MyError::Missing)?;

// OK in tests and prototypes
#[test]
fn test_something() {
    let val = some_option.unwrap();  // acceptable in tests
}
```

## 13. Write Small, Focused Functions

```rust
// BAD — does everything in one function
fn process_data(raw: &str) -> String {
    // 100 lines of parse + validate + transform + format
}

// GOOD — each function does one thing
fn parse(raw: &str) -> Result<Data, ParseError> { ... }
fn validate(data: &Data) -> Result<(), ValidationError> { ... }
fn transform(data: Data) -> Output { ... }
fn format(output: &Output) -> String { ... }

fn process_data(raw: &str) -> Result<String, AppError> {
    let data = parse(raw)?;
    validate(&data)?;
    let output = transform(data);
    Ok(format(&output))
}
```

## 14. Use Constants for Magic Numbers

```rust
// BAD
fn is_retirement_age(age: u32) -> bool { age >= 65 }
let timeout = 30_000;

// GOOD
const RETIREMENT_AGE: u32 = 65;
const TIMEOUT_MS: u64 = 30_000;

fn is_retirement_age(age: u32) -> bool { age >= RETIREMENT_AGE }
```

## 15. Prefer Owned Data in Structs (Unless You Have a Reason)

```rust
// Usually BAD for structs — lifetime annotations are contagious
struct User<'a> {
    name: &'a str,  // now User has a lifetime parameter
}

// Usually GOOD — simpler, more flexible
struct User {
    name: String,  // owns its data, no lifetime needed
}
```

## 16. Write Tests

```rust
fn add(a: i32, b: i32) -> i32 { a + b }

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_add_positive() {
        assert_eq!(add(2, 3), 5);
    }

    #[test]
    fn test_add_negative() {
        assert_eq!(add(-1, -1), -2);
    }

    #[test]
    #[should_panic(expected = "overflow")]
    fn test_overflow() {
        let _ = i32::MAX + 1;
    }

    #[test]
    fn test_with_result() -> Result<(), Box<dyn std::error::Error>> {
        let n: i32 = "42".parse()?;
        assert_eq!(n, 42);
        Ok(())
    }
}
```

## 17. Use cargo check Before cargo build

```bash
# cargo check is 10x faster than cargo build
# Use it for rapid iteration
cargo check

# Only build when you need to run or need the artifact
cargo build
```

## 18. Profile Before Optimizing

```rust
// WRONG approach: guess and micro-optimize

// RIGHT approach:
// 1. Write clean, correct code first
// 2. Measure with a profiler (cargo flamegraph, perf)
// 3. Identify the actual bottleneck
// 4. Optimize specifically that part
// 5. Verify the optimization helped
```

## 19. Use the Type System to Enforce Invariants

```rust
// BAD — invalid states are representable
struct Connection {
    host: String,
    is_connected: bool,
    socket: Option<Socket>,  // must be Some when is_connected, None otherwise
}

// GOOD — invalid states are impossible
enum Connection {
    Disconnected { host: String },
    Connected { host: String, socket: Socket },
}
```

## 20. Read the Error Messages — They're Excellent

```
error[E0502]: cannot borrow `s` as mutable because it is also borrowed as immutable
  --> src/main.rs:5:5
   |
4  |     let r = &s;
   |             -- immutable borrow occurs here
5  |     s.push_str(" world");
   |     ^^^^^^^^^^^^^^^^^^^^ mutable borrow occurs here
6  |     println!("{}", r);
   |                    - immutable borrow later used here
```

Rust error messages tell you **exactly** what went wrong, **where**, and often **how to fix it**. Always read the full error, especially the `help:` and `note:` sections.

---

## Quick Reference Card

| Pattern | Prefer | Avoid |
|---------|--------|-------|
| String params | `&str` | `&String` |
| Slice params | `&[T]` | `&Vec<T>` |
| Error propagation | `?` | nested `match` |
| Optional values | `Option` methods | `if let` chains |
| Loops | iterators | index loops |
| HashMap upsert | `.entry().or_insert()` | `.contains_key()` + `.insert()` |
| Error handling | `Result` + `?` | `.unwrap()` in library code |
| Type disambiguation | newtype pattern | raw primitives |
| Decision params | enums | booleans |
| Production code | `expect("message")` or `?` | `.unwrap()` |

---

*Next: [Chapter 91 — Common Pitfalls](91-common-pitfalls.md)*
