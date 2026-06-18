# Chapter 16 — Rust Internals: Compilation, Memory, and Zero-Cost Abstractions

---

## 16.1 The Compilation Pipeline

```
Source (.rs files)
       │
       ▼  [Lexing + Parsing]
    AST (Abstract Syntax Tree)
       │
       ▼  [Name Resolution + Macro Expansion]
    HIR (High-level Intermediate Representation)
       │
       ▼  [Type Checking + Trait Resolution + Borrow Checking]
    MIR (Mid-level Intermediate Representation)
       │
       ▼  [Optimizations: inlining, dead code elimination, etc.]
    LLVM IR
       │
       ▼  [LLVM backend: optimization + code generation]
    Native machine code
```

### The MIR — Where Most Magic Happens

MIR is Rust's internal representation where borrow checking occurs. It's a simplified, explicit CFG (Control Flow Graph):

```bash
# See MIR for any function
cargo rustc -- --emit=mir
# Output in target/debug/deps/*.mir
```

### What the Borrow Checker Actually Checks

The borrow checker (implemented as a dataflow analysis on MIR) verifies:
1. Every reference has a valid source (no dangling refs)
2. No aliasing between `&T` and `&mut T` to the same location
3. Moves and copies are consistent (moved values not used after move)
4. All paths through the function obey the rules

---

## 16.2 Monomorphization

Rust generics are **monomorphized** at compile time — the compiler generates separate code for each concrete type:

```rust
fn largest<T: PartialOrd>(list: &[T]) -> &T {
    let mut largest = &list[0];
    for item in list {
        if item > largest { largest = item; }
    }
    largest
}

fn main() {
    largest(&[1, 2, 3]);    // Compiler generates: largest_i32
    largest(&[1.0, 2.0]);   // Compiler generates: largest_f64
}
```

The compiled binary contains two separate functions — as if you'd written:

```rust
fn largest_i32(list: &[i32]) -> &i32 { ... }
fn largest_f64(list: &[f64]) -> &f64 { ... }
```

**Benefits:**
- Zero runtime overhead — no vtable, no boxing
- Each specialization can be fully optimized for its type

**Trade-offs:**
- Larger binary (code is duplicated per type)
- Longer compile times for heavily generic code

---

## 16.3 Trait Objects and vtables

Dynamic dispatch via `dyn Trait` uses a **vtable** — a pointer to a table of function pointers:

```
Box<dyn Draw>:
┌─────────────────────────┐
│ data pointer (8 bytes)  ├──► actual object on heap
│ vtable pointer (8 bytes)├──► vtable: [draw fn ptr, drop fn ptr, size, align]
└─────────────────────────┘
```

```rust
trait Draw { fn draw(&self); }

struct Circle { radius: f64 }
impl Draw for Circle { fn draw(&self) { println!("circle"); } }

fn main() {
    let d: Box<dyn Draw> = Box::new(Circle { radius: 5.0 });
    // d is 16 bytes: pointer to Circle + pointer to Circle's vtable
    println!("{}", std::mem::size_of_val(&d));

    d.draw();  // indirect call through vtable — one pointer indirection
}
```

**Static dispatch**: call is resolved at compile time — zero overhead  
**Dynamic dispatch**: call goes through vtable — one pointer indirection + cache miss risk

---

## 16.4 Memory Layout

### Stack Layout

```rust
fn foo() {
    let x: i32 = 5;       // 4 bytes on stack
    let y: f64 = 3.14;    // 8 bytes on stack
    let z: bool = true;   // 1 byte on stack (but aligned to 4 or 8)
}
// Stack frame: ~20+ bytes, automatically freed on return
```

### Struct Layout and Alignment

Rust aligns struct fields to their natural alignment (size of the field). It may add padding:

```rust
use std::mem::{size_of, align_of};

struct A {
    x: u8,   // 1 byte
    y: u32,  // 4 bytes (needs 4-byte alignment)
    z: u8,   // 1 byte
}
// Actual layout: [x: 1 byte][pad: 3 bytes][y: 4 bytes][z: 1 byte][pad: 3 bytes]
// Total: 12 bytes (not 6!)

struct B {
    y: u32,  // 4 bytes
    x: u8,   // 1 byte
    z: u8,   // 1 byte
}
// Actual layout: [y: 4 bytes][x: 1 byte][z: 1 byte][pad: 2 bytes]
// Total: 8 bytes — better field ordering!

fn main() {
    println!("A: {} bytes, align {}", size_of::<A>(), align_of::<A>());  // 12, 4
    println!("B: {} bytes, align {}", size_of::<B>(), align_of::<B>());  // 8, 4
}
```

### repr Attributes — Control Layout

```rust
// Default Rust layout (may reorder fields for optimization)
struct Normal { a: u8, b: u32 }

// C-compatible layout (same as C struct, fields in order)
#[repr(C)]
struct CCompatible { a: u8, b: u32 }

// Packed — no padding (unaligned access — dangerous!)
#[repr(packed)]
struct Packed { a: u8, b: u32 }

// Specific alignment
#[repr(align(64))]  // 64-byte aligned (useful for cache line alignment)
struct CacheAligned { data: [u8; 64] }

fn main() {
    println!("{}", std::mem::size_of::<Normal>());     // 8 (Rust reorders)
    println!("{}", std::mem::size_of::<CCompatible>()); // 8 (C layout)
    println!("{}", std::mem::size_of::<Packed>());     // 5 (no padding!)
}
```

### Enum Layout

Rust enums are discriminated unions:

```rust
enum Option<T> {
    Some(T),
    None,
}

// Size of Option<i32>:
// discriminant (tag): usually 1-4 bytes (optimized away when possible)
// variant data: size of largest variant

// Null Pointer Optimization (NPO):
// Option<&T> is the same size as &T!
// None is represented as a null pointer — no extra byte needed

fn main() {
    use std::mem::size_of;
    println!("{}", size_of::<Option<i32>>());  // 8 (4 data + 4 tag, aligned)
    println!("{}", size_of::<Option<&i32>>());  // 8 (same as &i32 — NPO!)
    println!("{}", size_of::<&i32>());          // 8

    // Box<T> also benefits from NPO
    println!("{}", size_of::<Option<Box<i32>>>());  // 8 (same as Box<i32>)
}
```

---

## 16.5 Zero-Cost Abstractions

Rust's abstractions compile away to nothing:

### Iterators

```rust
// High-level iterator code
fn sum_of_squares_filtered(v: &[i32]) -> i32 {
    v.iter()
        .filter(|&&x| x > 0)
        .map(|&x| x * x)
        .sum()
}

// Compiles to exactly the same machine code as:
fn sum_of_squares_manual(v: &[i32]) -> i32 {
    let mut result = 0;
    for &x in v {
        if x > 0 {
            result += x * x;
        }
    }
    result
}
// No overhead. No allocation. No indirection.
```

### Closures

```rust
let factor = 2;
let double = |x| x * factor;

// Compiled to:
struct DoubleClosure { factor: i32 }
impl DoubleClosure {
    fn call(&self, x: i32) -> i32 { x * self.factor }
}
// When called through `fn` pointer or generic: inlined completely
```

### String Formatting

```rust
println!("{}", x);
// Compiled to a direct write syscall — no intermediate String allocation
// format!() does allocate — use sparingly in hot paths
```

---

## 16.6 RAII and Drop Order

When multiple values go out of scope, they're dropped in **reverse declaration order** (LIFO):

```rust
struct Droppable { name: &'static str }
impl Drop for Droppable {
    fn drop(&mut self) { println!("Dropping {}", self.name); }
}

fn main() {
    let _a = Droppable { name: "a" };
    let _b = Droppable { name: "b" };
    let _c = Droppable { name: "c" };
}
// Output:
// Dropping c
// Dropping b
// Dropping a
```

This is critical for resource management: a file buffer is flushed before the file handle is closed because the buffer is declared after the file.

---

## 16.7 The Allocator

Rust uses the system allocator (jemalloc was removed in favor of the OS allocator). You can replace it:

```rust
use std::alloc::{GlobalAlloc, System, Layout};

// Use the system allocator explicitly
#[global_allocator]
static GLOBAL: System = System;

// Or use a custom allocator:
// extern crate jemallocator;
// #[global_allocator]
// static ALLOC: jemallocator::Jemalloc = jemallocator::Jemalloc;
```

---

## 16.8 Unsafe and Raw Pointers

How raw pointers work internally:

```rust
fn main() {
    let x = 42;
    let r: &i32 = &x;                  // safe reference: borrow-checked
    let ptr: *const i32 = &x;          // raw pointer: no checks
    let ptr2: *const i32 = r as *const i32;  // reference to raw pointer

    unsafe {
        println!("{}", *ptr);           // dereference raw pointer
        println!("{}", *ptr2);
        println!("addr: {:p}", ptr);    // print memory address
    }

    // Null raw pointer (cannot create null reference!)
    let null: *const i32 = std::ptr::null();
    println!("is null: {}", null.is_null());  // true
}
```

---

## 16.9 Inline and LTO

```rust
// Force inlining
#[inline(always)]
fn hot_path(x: i32) -> i32 { x * 2 }

// Prevent inlining (for code size reduction)
#[inline(never)]
fn cold_path(x: i32) -> i32 { expensive_computation(x) }

// Hint compiler this is unlikely
fn parse_input(input: &str) -> i32 {
    input.parse().unwrap_or_else(|_| {
        // #[cold] hint on the function would help here
        eprintln!("Invalid input");
        0
    })
}
```

```toml
# Cargo.toml — Link-Time Optimization
[profile.release]
opt-level = 3
lto = true          # enables LTO — better optimization across crates
codegen-units = 1   # single codegen unit — better optimization, slower compile
strip = true        # strip debug symbols from binary
```

---

## 16.10 Compile-Time Computation

`const fn` runs at compile time:

```rust
const fn fibonacci(n: u64) -> u64 {
    match n {
        0 => 0,
        1 => 1,
        _ => fibonacci(n - 1) + fibonacci(n - 2),
    }
}

// Evaluated at compile time — zero runtime cost
const FIB_30: u64 = fibonacci(30);  // 832040

fn main() {
    println!("{}", FIB_30);  // already computed — just loads a constant
}
```

---

## Summary

Rust compiles through AST → HIR → MIR → LLVM IR → machine code. Borrow checking happens on MIR. Generics are monomorphized: one specialized copy per type, zero runtime overhead. `dyn Trait` uses vtables for dynamic dispatch: two pointers (data + vtable), one indirect call. Struct layout may include padding for alignment; use `#[repr(C)]` for C compatibility. Rust's abstractions — iterators, closures, RAII — compile to the same code as handwritten loops and manual memory management. `const fn` moves computation to compile time.

---

## Key Takeaways

- Borrow checker works on MIR (control flow graph), not the raw source
- Monomorphization = generics expand to concrete types — fast, larger binary
- Dynamic dispatch (`dyn Trait`) = 16-byte fat pointer, one vtable call
- Null Pointer Optimization: `Option<&T>` is the same size as `&T`
- Iterators and closures are zero-cost — no allocation, inlined away
- Drop order is LIFO (reverse declaration) within a scope
- `#[repr(C)]` for FFI; `#[repr(packed)]` for space, `#[repr(align(N))]` for alignment
- `const fn` enables compile-time computation

---

## Exercises

**Exercise 1:** Use `std::mem::size_of` to measure the size of: `bool`, `i32`, `(i32, bool)`, `[i32; 3]`, `&i32`, `Box<i32>`, `Option<&i32>`, `Option<i32>`. Explain any surprising values.

**Exercise 2:** Write two versions of the same algorithm — one using iterator chains, one using explicit loops. Use `cargo build --release` and compare the assembly output with `cargo rustc --release -- --emit=asm`.

**Exercise 3:** Create a struct with fields in poor ordering and measure its size. Reorder the fields optimally. Verify the size difference.

**Exercise 4:** Write a `const fn` that computes prime numbers up to N at compile time, storing them in a const array.

**Exercise 5:** Implement a custom smart pointer `Logged<T>` that prints "drop" when the inner value is dropped. Verify the drop order when multiple `Logged` values go out of scope together.

---

*Next: [Chapter 90 — Best Practices](90-best-practices.md)*
