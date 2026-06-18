# Chapter 8 — Lifetimes: The Borrow Checker's Logic

> *"Lifetime annotations don't change how long references live — they describe relationships between reference lifetimes so the borrow checker can validate your logic."*

This is the most conceptually challenging chapter in the Rust book. Read it slowly. Run every example. The investment pays off: once you understand lifetimes, the borrow checker becomes a tool, not an obstacle.

---

## 8.1 The Problem Lifetimes Solve

Consider this function — it takes two string slices and returns the longer one:

```rust
fn longest(x: &str, y: &str) -> &str {  // COMPILE ERROR
    if x.len() > y.len() { x } else { y }
}
```

The compiler rejects this with:
```
error[E0106]: missing lifetime specifier
 --> src/main.rs:1:33
  |
1 | fn longest(x: &str, y: &str) -> &str {
  |               ----     ----     ^ expected named lifetime parameter
```

**Why?** The compiler looks at the return type `&str` and asks: "This returned reference borrows from... where? From `x`? From `y`? How long is it valid?"

The answer depends on runtime input — but the compiler needs to verify correctness at compile time. It cannot know which branch will be taken. So it needs you to tell it: "the returned reference lives as long as the shorter-lived of x and y."

That's what a lifetime annotation expresses.

---

## 8.2 Lifetime Annotation Syntax

Lifetime annotations are written with a tick and a short name: `'a`, `'b`, `'static`. They go after the `&`:

```
&'a T        — immutable reference with lifetime 'a
&'a mut T    — mutable reference with lifetime 'a
```

They describe **relationships** between lifetimes, not absolute durations. They don't change how long anything lives — they're just labels so the compiler can verify consistency.

---

## 8.3 Lifetime Annotations in Functions

### Fixing the longest Function

```rust
fn longest<'a>(x: &'a str, y: &'a str) -> &'a str {
    if x.len() > y.len() { x } else { y }
}
```

Read this as: "There exists some lifetime `'a`. `x` lives at least as long as `'a`. `y` lives at least as long as `'a`. The returned reference lives at most as long as `'a`."

In practice: `'a` will be the **intersection** (shorter) of x's and y's actual lifetimes.

```rust
fn main() {
    let s1 = String::from("long string");
    let result;

    {
        let s2 = String::from("xy");       // s2 has a shorter lifetime
        result = longest(s1.as_str(), s2.as_str());
        println!("Longest: {}", result);   // OK — result used within s2's scope
    }

    // println!("{}", result);  // COMPILE ERROR: result might point to s2,
                                // which is already dropped!
}
```

The lifetime annotation told the compiler: "the result lives only as long as the shorter input." The compiler uses this to reject the code above when `result` is used after `s2` is dropped.

### Single-Input Lifetime

When there's only one possible source for the output:

```rust
// No ambiguity — the result must come from s
fn first_word(s: &str) -> &str {
    &s[..s.find(' ').unwrap_or(s.len())]
}
// lifetime elision fills in: fn first_word<'a>(s: &'a str) -> &'a str
```

The compiler handles this automatically (lifetime elision rules, section 8.6).

---

## 8.4 Lifetime Annotations in Structs

When a struct holds references, those references must not outlive the struct:

```rust
struct ImportantExcerpt<'a> {
    part: &'a str,  // this struct cannot outlive the string 'part' points to
}

impl<'a> ImportantExcerpt<'a> {
    fn level(&self) -> i32 {
        3
    }

    fn announce_and_return_part(&self, announcement: &str) -> &str {
        println!("Attention: {}!", announcement);
        self.part  // returns a reference with the lifetime of 'self
    }
}

fn main() {
    let novel = String::from("Call me Ishmael. Some years ago...");

    let excerpt = {
        let first_sentence = novel
            .split('.')
            .next()
            .expect("Could not find a '.'");
        ImportantExcerpt { part: first_sentence }
        // ImportantExcerpt borrows from `novel` via `first_sentence`
    };

    // excerpt is valid because `novel` is still alive
    println!("{}", excerpt.part);

    // If novel were dropped before excerpt, the compiler would reject it:
    // drop(novel);              // novel freed
    // println!("{}", excerpt.part);  // COMPILE ERROR — excerpt contains
                                       // a reference to freed data
}
```

---

## 8.5 The Borrow Checker — How It Works

The borrow checker is a dataflow analysis that tracks:
1. The scope of every value (where it was created, where it's dropped)
2. The scope of every reference (where it's created, where it's last used)
3. Whether any reference's scope exceeds its referent's scope

```rust
fn main() {
    let r;                    // r is declared — no value yet

    {
        let x = 5;            // x is created — lifetime begins
        r = &x;              // r borrows x
    }                         // x is dropped — lifetime ends

    println!("{}", r);        // COMPILE ERROR: r's lifetime extends beyond x's
}
```

The borrow checker visualizes this as:

```
let r;         |-----------|  r's lifetime
{              |            |
    let x = 5; |   |--------|  x's lifetime
    r = &x;    |   |        |
}              |   |  x dropped here
println!(r);   |            |  r used here — but x is gone!
```

Since `r`'s lifetime (`r`) extends beyond `x`'s lifetime (`x`), the reference is dangling. Rejected.

Fix:

```rust
fn main() {
    let x = 5;        // x's lifetime is the entire main function
    let r = &x;       // r borrows x — r's lifetime fits inside x's
    println!("{}", r); // OK
}
```

---

## 8.6 Lifetime Elision Rules

In many common patterns, the compiler can infer lifetimes automatically. These are the **elision rules** — you don't need to write annotations when these rules apply.

**Rule 1:** Each reference parameter gets its own lifetime parameter.
```rust
fn foo(x: &str) → &str            // becomes:
fn foo<'a>(x: &'a str) -> &'a str // (after applying rule 3)
```

**Rule 2:** If there is exactly one input lifetime, it's assigned to all outputs.
```rust
fn first_word(s: &str) -> &str    // becomes:
fn first_word<'a>(s: &'a str) -> &'a str  // output gets input's lifetime
```

**Rule 3:** If one of the inputs is `&self` or `&mut self`, its lifetime is assigned to all outputs.
```rust
impl Excerpt {
    fn part(&self) -> &str    // becomes:
    fn part<'a>(&'a self) -> &'a str  // output gets self's lifetime
}
```

If all three rules are exhausted and there's still ambiguity, the compiler requires you to annotate explicitly:

```rust
// Two references, not self — rules don't resolve the output lifetime
fn longest(x: &str, y: &str) -> &str  // ERROR — ambiguous
// Must annotate:
fn longest<'a>(x: &'a str, y: &'a str) -> &'a str  // OK
```

---

## 8.7 The 'static Lifetime

`'static` is the lifetime that lasts the entire program duration:

```rust
// String literals have 'static lifetime — they're stored in the binary
let s: &'static str = "I live forever";

// Functions returning 'static references never dangle
fn static_greeting() -> &'static str {
    "Hello, World!"  // stored in binary — always valid
}

// Constants have 'static lifetime
static PI: f64 = 3.14159;
const MAX: u32 = 100;
```

You'll sometimes see `'static` in trait bounds:

```rust
fn print_and_store<T: std::fmt::Display + 'static>(value: T) {
    // 'static here means T contains no non-static references
    // i.e., T owns all its data, or its references are 'static
    println!("{}", value);
}
```

---

## 8.8 Lifetime Annotations in impl Blocks

```rust
struct Config<'a> {
    host: &'a str,
    port: u16,
}

// Lifetime is part of the type — must be specified in impl
impl<'a> Config<'a> {
    fn new(host: &'a str, port: u16) -> Self {
        Config { host, port }
    }

    fn host(&self) -> &str {
        self.host  // lifetime elision: output lifetime = 'self lifetime
    }

    fn with_port(&self, port: u16) -> Config<'_> {
        // '_ means "infer the lifetime" — shorthand for 'a
        Config { host: self.host, port }
    }
}

fn main() {
    let hostname = String::from("localhost");
    let config = Config::new(&hostname, 8080);
    println!("{}:{}", config.host(), config.port);
}
```

---

## 8.9 Multiple Lifetime Parameters

Sometimes you need to express different lifetime constraints:

```rust
// x and y have different lifetimes; the result is tied to x, not y
fn first_if_longer<'a, 'b>(x: &'a str, y: &'b str) -> &'a str {
    if x.len() > y.len() { x } else { x }  // always return x
}

// Structure with multiple reference lifetimes
struct Pair<'a, 'b> {
    first: &'a str,
    second: &'b str,
}

impl<'a, 'b> Pair<'a, 'b> {
    fn longer(&self) -> &str {
        if self.first.len() >= self.second.len() {
            self.first    // output lifetime = 'a (from self.first)
        } else {
            self.second   // wait — this is 'b, not 'a...
        }
    }
    // Actually this is ambiguous — the compiler would require annotations
}
```

When you have multiple lifetime parameters, think carefully about which references the output could come from, and annotate accordingly.

---

## 8.10 Combining Lifetimes with Generics and Traits

```rust
use std::fmt::Display;

fn longest_with_announcement<'a, T>(
    x: &'a str,
    y: &'a str,
    ann: T,
) -> &'a str
where
    T: Display,
{
    println!("Announcement: {}", ann);
    if x.len() > y.len() { x } else { y }
}

fn main() {
    let s1 = String::from("long string");
    let s2 = String::from("xy");
    let result = longest_with_announcement(
        &s1,
        &s2,
        "Today is June 2026",
    );
    println!("Longest: {}", result);
}
```

---

## 8.11 Common Lifetime Patterns

### Owned Data in Structs (No Lifetime Needed)

If your struct owns its data (no references), no lifetime annotations needed:

```rust
// Owned — no lifetime
struct Config {
    host: String,   // owned String
    port: u16,
}

// Borrowed — needs lifetime
struct ConfigRef<'a> {
    host: &'a str,  // reference — must not outlive source
    port: u16,
}
```

In general, prefer owned data in structs unless you have a specific performance reason to use references. Lifetime-free structs are simpler to use.

### Returning References to Input Data

```rust
struct Cache {
    data: Vec<String>,
}

impl Cache {
    fn get(&self, index: usize) -> Option<&str> {
        // returns a reference into self.data
        // lifetime elision: output lifetime = 'self lifetime
        self.data.get(index).map(|s| s.as_str())
    }

    fn longest_entry(&self) -> Option<&str> {
        self.data.iter().max_by_key(|s| s.len()).map(|s| s.as_str())
    }
}

fn main() {
    let mut cache = Cache { data: Vec::new() };
    cache.data.push(String::from("hello"));
    cache.data.push(String::from("world!"));

    if let Some(s) = cache.longest_entry() {
        println!("Longest: {}", s);  // s borrows from cache
    }
    // cache is still valid — we only borrowed from it
}
```

---

## 8.12 Lifetime Variance

Lifetimes have a subtyping relationship. A `'long` lifetime is a subtype of a `'short` lifetime (because something that lives longer can be used wherever something shorter is expected):

```rust
fn takes_short_lived<'a>(x: &'a str) -> &'a str { x }

fn main() {
    let long_lived = String::from("hello");  // lives for entire main
    let result;

    {
        let s: &str = "static string";  // 'static lifetime
        // 'static can be coerced to any shorter lifetime
        result = takes_short_lived(s);
        println!("{}", result);
    }
}
```

This is called **covariance** — longer lifetimes can coerce to shorter ones when needed.

---

## The Lifetime Mental Model

When you're confused about lifetimes, ask these questions:

1. **What is the function doing?** Transforming? Selecting? Creating?
2. **If it returns a reference, where does that reference point?** Into the first argument? The second? A global?
3. **Label the "where it points to" with the same lifetime as the return.** If it could be either input, they must share a lifetime.
4. **For structs holding references**: the struct cannot outlive the reference source.

```rust
// Q: "Where does the return reference come from?"
// A: "From whichever of x or y is longer — so both."
fn longest<'a>(x: &'a str, y: &'a str) -> &'a str { ... }

// Q: "Where does the return reference come from?"
// A: "Always from x."
fn always_first<'a, 'b>(x: &'a str, _y: &'b str) -> &'a str { x }

// Q: "The struct holds a &str — where does that str come from?"
// A: "From 'a — wherever the caller points it."
struct Excerpt<'a> { part: &'a str }
```

---

## Summary

Lifetimes are compile-time annotations that describe how long references are valid relative to each other. They don't change program behavior — they give the borrow checker the information it needs to verify correctness. Most lifetimes are inferred via elision rules. You need explicit annotations when the compiler can't determine output lifetime from the inputs alone — primarily: functions returning a reference that could come from multiple inputs, and structs holding references. `'static` lifetime lasts the entire program. Lifetimes prevent dangling references entirely at compile time.

---

## Key Takeaways

- Lifetimes are labels describing relationships between reference scopes
- `'a` in `fn foo<'a>(x: &'a T) -> &'a T` means: "output lives as long as input"
- The borrow checker rejects code where a reference outlives its source
- Elision rules handle common patterns automatically — you rarely need explicit annotations
- Structs holding references need lifetime annotations: `struct Foo<'a> { s: &'a str }`
- `'static` = lives for the whole program (string literals, constants)
- Prefer owned data in structs to avoid lifetime complexity

---

## Exercises

**Exercise 1:** Without running it, explain exactly why this fails and how to fix it:
```rust
fn main() {
    let result;
    {
        let s = String::from("hello");
        result = &s;
    }
    println!("{}", result);
}
```

**Exercise 2:** Write `fn longest_starting_with<'a>(text: &'a str, prefix: &str) -> &'a str` that returns the longest word in `text` starting with `prefix`, or an empty slice if none found.

**Exercise 3:** Create a struct `StrSplit<'a>` that holds a `&'a str` (the string to split) and a `&'a str` (the delimiter). Implement `fn next(&mut self) -> Option<&'a str>` returning the next split segment.

**Exercise 4:** What lifetime relationship does `fn first_or_default<'a>(slice: &'a [i32], default: &'a i32) -> &'a i32` express? When would the compiler reject calls to this function?

**Exercise 5:** Rewrite this to compile correctly by switching from references to owned types in the struct, and explain the trade-off:
```rust
struct User<'a> {
    name: &'a str,
    email: &'a str,
}
```

---

*Next: [Chapter 9 — Structs and Enums](09-structs-enums.md)*
