# Chapter 3 — Variables, Types, and Expressions

> *"In Rust, every value has an explicit type. There is no implicit coercion. What you see is what you get."*

---

## 3.1 let and mut

### Immutable by Default

Variables in Rust are **immutable by default**. This is the opposite of most languages. You must explicitly opt into mutability.

```rust
fn main() {
    let x = 5;
    // x = 6;  // COMPILE ERROR: cannot assign twice to immutable variable `x`

    let mut y = 5;  // mut makes it mutable
    y = 6;          // OK
    println!("y = {}", y);
}
```

**Why immutable by default?** Because mutation is a common source of bugs. Making immutability the default forces you to consciously decide when state needs to change.

### Type Annotations

Rust infers types in most cases, but you can annotate explicitly:

```rust
fn main() {
    let x = 5;          // inferred: i32
    let y: i32 = 5;     // explicit
    let z: f64 = 5.0;   // explicit float
    let s: &str = "hello";  // explicit string slice

    // Type annotation is required when inference can't determine the type
    let mut v: Vec<i32> = Vec::new();  // must specify element type
    v.push(1);
}
```

---

## 3.2 Shadowing

Shadowing lets you re-declare a variable with the same name, even changing its type. This is different from mutation.

```rust
fn main() {
    let x = 5;
    let x = x + 1;      // shadows the previous x
    let x = x * 2;      // shadows again
    println!("x = {}", x);  // 12

    // Shadowing can change the type
    let spaces = "   ";         // &str
    let spaces = spaces.len();  // usize — completely different type
    println!("spaces = {}", spaces);  // 3

    // With mut, you CANNOT change the type:
    let mut spaces = "   ";
    // spaces = spaces.len();  // COMPILE ERROR: mismatched types
}
```

### Shadowing vs Mutation

```rust
// Mutation — same variable, same type, same memory location
let mut x = 5;
x = x + 1;  // x is still the same binding

// Shadowing — new binding, potentially new type, new memory
let x = 5;
let x = x + 1;  // x is a completely new binding that happens to have the same name
                 // the old x is no longer accessible
```

---

## 3.3 Primitive Types

### Integer Types

```rust
fn main() {
    // Signed integers
    let a: i8  = 127;
    let b: i16 = 32_767;
    let c: i32 = 2_147_483_647;   // default integer type
    let d: i64 = 9_223_372_036_854_775_807;
    let e: i128 = 170_141_183_460_469_231_731_687_303_715_884_105_727;
    let f: isize = 42;  // pointer-sized (64-bit on 64-bit systems)

    // Unsigned integers
    let g: u8  = 255;
    let h: u16 = 65_535;
    let i: u32 = 4_294_967_295;
    let j: u64 = 18_446_744_073_709_551_615;
    let k: u128 = 340_282_366_920_938_463_463_374_607_431_768_211_455;
    let l: usize = 0;  // used for indexing

    // Integer literals
    let decimal     = 1_000_000;  // underscores for readability
    let hex         = 0xFF;
    let octal       = 0o77;
    let binary      = 0b1111_0000;
    let byte: u8    = b'A';       // byte literal (u8 only)

    println!("{} {} {} {}", decimal, hex, octal, binary);
}
```

### Integer Overflow

```rust
fn main() {
    let x: u8 = 255;

    // In debug mode: panics with overflow
    // In release mode: wraps around (255 + 1 = 0 for u8)

    // Explicit wrapping:
    let wrapped = x.wrapping_add(1);  // 0
    println!("wrapped: {}", wrapped);

    // Checked (returns Option):
    let checked = x.checked_add(1);  // None
    println!("checked: {:?}", checked);

    // Saturating (stays at max):
    let saturated = x.saturating_add(1);  // 255
    println!("saturated: {}", saturated);
}
```

### Floating-Point Types

```rust
fn main() {
    let x: f32 = 3.14;    // 32-bit float
    let y: f64 = 3.14;    // 64-bit float — default
    let z = 3.14;          // inferred as f64

    // Floating point operations
    println!("{}", 2.0_f64.sqrt());     // 1.4142135623730951
    println!("{}", f64::MAX);           // 1.7976931348623157e308
    println!("{}", f64::INFINITY);      // inf
    println!("{}", f64::NAN);           // NaN
    println!("{}", f64::NAN.is_nan());  // true

    // Precision matters
    println!("{:.10}", 0.1 + 0.2);  // 0.3000000000 (IEEE 754 float)
}
```

### Booleans

```rust
fn main() {
    let t: bool = true;
    let f: bool = false;

    println!("{}", t && f);  // false
    println!("{}", t || f);  // true
    println!("{}", !t);      // false

    // Booleans are 1 byte in Rust
    println!("{}", std::mem::size_of::<bool>());  // 1
}
```

### Characters

```rust
fn main() {
    let c: char = 'z';
    let emoji: char = '🦀';         // Rust supports full Unicode
    let chinese: char = '中';

    println!("{} {} {}", c, emoji, chinese);
    println!("{}", c as u32);        // Unicode scalar value: 122
    println!("{}", '🦀' as u32);    // 129408

    // char is 4 bytes (Unicode scalar value)
    println!("{}", std::mem::size_of::<char>());  // 4
}
```

---

## 3.4 Tuples and Arrays

### Tuples — Fixed-Length, Mixed Types

```rust
fn main() {
    // Tuple with mixed types
    let tup: (i32, f64, bool, char) = (42, 3.14, true, 'z');

    // Destructuring
    let (x, y, z, w) = tup;
    println!("{} {} {} {}", x, y, z, w);

    // Index access with .0, .1, etc.
    println!("{}", tup.0);  // 42
    println!("{}", tup.1);  // 3.14

    // Nested tuples
    let nested = ((1, 2), (3, 4));
    println!("{}", nested.0.1);  // 2

    // Unit tuple — empty tuple, used for "no value"
    let unit: () = ();
    // Functions that return nothing actually return ()
}
```

### Arrays — Fixed-Length, Same Type

```rust
fn main() {
    // Array syntax: [Type; length]
    let arr: [i32; 5] = [1, 2, 3, 4, 5];

    // Repeat syntax: [value; count]
    let zeros = [0; 10];  // [0, 0, 0, 0, 0, 0, 0, 0, 0, 0]

    // Indexing
    println!("{}", arr[0]);  // 1
    println!("{}", arr[4]);  // 5
    println!("{}", arr.len());  // 5

    // Arrays are stack-allocated
    println!("{}", std::mem::size_of::<[i32; 5]>());  // 20 bytes (5 * 4)

    // Out-of-bounds access panics at runtime (not UB like C)
    // println!("{}", arr[10]);  // thread 'main' panicked: index out of bounds

    // Slices — reference to part of an array
    let slice: &[i32] = &arr[1..4];  // [2, 3, 4]
    println!("{:?}", slice);

    // Iterating
    for element in &arr {
        print!("{} ", element);
    }
    println!();

    for (i, element) in arr.iter().enumerate() {
        println!("arr[{}] = {}", i, element);
    }
}
```

### Arrays vs Vectors

```rust
fn main() {
    // Array — fixed size, stack allocated, known at compile time
    let arr = [1, 2, 3, 4, 5];  // [i32; 5]

    // Vec — dynamic size, heap allocated
    let vec = vec![1, 2, 3, 4, 5];  // Vec<i32>

    // Use arrays when size is fixed and known at compile time
    // Use Vec when size can change or is unknown at compile time
}
```

---

## 3.5 Constants and Statics

### const — Compile-Time Constants

```rust
// Constants: must have a type, evaluated at compile time
const MAX_POINTS: u32 = 100_000;
const PI: f64 = 3.14159265358979;

fn main() {
    println!("Max: {}", MAX_POINTS);
    println!("Pi: {}", PI);

    // const in a function — limited to that scope
    const LOCAL_MAX: i32 = 42;
    println!("{}", LOCAL_MAX);
}
```

### static — Global Variables

```rust
// Static: lives for the entire program duration
static GREETING: &str = "Hello, World!";
static mut COUNTER: i32 = 0;  // mutable static is unsafe

fn main() {
    println!("{}", GREETING);

    // Accessing mutable static requires unsafe
    unsafe {
        COUNTER += 1;
        println!("Counter: {}", COUNTER);
    }
}
```

### const vs static

| Aspect | `const` | `static` |
|--------|---------|---------|
| Evaluation | Compile time | Compile time |
| Memory | Inlined at use sites | Single memory location |
| Mutability | Never mutable | Can be `static mut` (unsafe) |
| Lifetime | N/A (no address) | `'static` (whole program) |

---

## 3.6 Stack vs Heap (Critical Foundation)

This is the most important section in this chapter. Without understanding stack vs. heap, ownership is impossible to understand.

### The Stack

The **stack** is a region of memory that operates like a stack of plates:
- You put data on top (push)
- You take data from the top (pop)
- **LIFO**: Last In, First Out
- Data must have a **known, fixed size at compile time**
- Allocation and deallocation are **instantaneous** (just move a pointer)
- **Very fast**

```
Stack (grows down):
┌─────────────────┐  ← stack pointer (top)
│  local var: y   │  (just pushed)
│  local var: x   │  
│  return address │
│  (previous fn)  │
└─────────────────┘
```

```rust
fn main() {
    let x: i32 = 5;    // pushed onto stack — 4 bytes
    let y: bool = true; // pushed onto stack — 1 byte
    // When main() returns, x and y are automatically popped
}
```

### The Heap

The **heap** is a large, less organized region of memory:
- You **request** a block of memory from the OS/allocator
- The allocator finds a free block of the right size and returns a **pointer**
- The pointer (address) is fixed size — it goes on the stack
- The actual data lives on the heap
- You must **explicitly free** the memory when done (or it leaks)

```
Stack:               Heap:
┌─────────┐         ┌────────────────────┐
│  ptr ──────────►  │ "hello"            │
│  len: 5 │         │ h e l l o          │
│  cap: 5 │         └────────────────────┘
└─────────┘
(String metadata   (actual string data —
 on the stack)      heap-allocated)
```

```rust
fn main() {
    // Stack-allocated: i32, bool, char, arrays, tuples (of stack types)
    let x: i32 = 5;         // 4 bytes on the stack, that's it

    // Heap-allocated: String, Vec, Box, etc.
    let s = String::from("hello");
    // Stack: pointer to heap data + length + capacity (24 bytes total on 64-bit)
    // Heap: the actual bytes "hello" (5 bytes)

    println!("{}", x);
    println!("{}", s);
}
// When s goes out of scope here, Rust automatically frees the heap memory
// This is the DROP mechanism — the heart of Rust's memory management
```

### Why This Matters for Ownership

The fundamental question Rust's ownership system answers is: **when should heap memory be freed?**

- Stack memory is automatically freed when a function returns — no decision needed
- Heap memory must be tracked — who is responsible for freeing it?

In C: programmer's problem → use-after-free bugs, double-free bugs, memory leaks  
In Java: garbage collector's problem → GC pauses, runtime overhead  
In Rust: **ownership system's problem** → compile-time tracking, zero runtime cost

```rust
fn main() {
    let s1 = String::from("hello");  // s1 owns the heap data

    {
        let s2 = String::from("world");  // s2 owns this heap data
        println!("{}", s2);
    }  // s2 goes out of scope — Rust automatically calls drop(s2), freeing heap memory

    println!("{}", s1);
}  // s1 goes out of scope — Rust automatically calls drop(s1)
   // No manual free() needed. No garbage collector. Zero overhead.
```

This is the key insight: **Rust ties the lifetime of heap data to the variable that owns it**. When the variable goes out of scope, the data is freed. The compiler tracks this statically.

---

## Type Sizes at a Glance

```rust
fn main() {
    use std::mem::size_of;

    // Primitives
    println!("bool:  {} bytes", size_of::<bool>());   // 1
    println!("i8:    {} bytes", size_of::<i8>());     // 1
    println!("i16:   {} bytes", size_of::<i16>());    // 2
    println!("i32:   {} bytes", size_of::<i32>());    // 4
    println!("i64:   {} bytes", size_of::<i64>());    // 8
    println!("i128:  {} bytes", size_of::<i128>());   // 16
    println!("f32:   {} bytes", size_of::<f32>());    // 4
    println!("f64:   {} bytes", size_of::<f64>());    // 8
    println!("char:  {} bytes", size_of::<char>());   // 4
    println!("usize: {} bytes", size_of::<usize>());  // 8 on 64-bit

    // Compound
    println!("(i32, bool): {} bytes", size_of::<(i32, bool)>());  // 8 (padding!)
    println!("[i32; 5]:   {} bytes", size_of::<[i32; 5]>());      // 20

    // Pointer-sized types
    println!("&i32:   {} bytes", size_of::<&i32>());   // 8 (pointer)
    println!("String: {} bytes", size_of::<String>()); // 24 (ptr + len + cap)
    println!("Vec<i32>: {} bytes", size_of::<Vec<i32>>());  // 24
}
```

---

## Summary

Rust variables are immutable by default — use `mut` to opt into mutation. Shadowing creates a new binding with the same name, allowing type changes. Rust has rich integer types (i8–i128, u8–u128), two float types (f32, f64), bool, and char (4-byte Unicode). Tuples hold mixed types; arrays hold same-type fixed-length sequences. Constants (`const`) are compile-time evaluated; statics (`static`) live for the program's entire duration. The stack/heap distinction is critical for understanding ownership — stack data is fixed-size and automatically managed, heap data requires explicit ownership tracking.

---

## Key Takeaways

- Variables are immutable by default — this prevents accidental mutation
- Shadowing ≠ mutation — shadowing creates a new binding; mutation updates the existing one
- Default integer type is `i32`; default float type is `f64`
- Arrays have a fixed size known at compile time; `Vec` is a dynamic heap-allocated list
- The stack is fast, fixed-size, LIFO; the heap is flexible but requires management
- Rust ties heap memory lifetime to the owning variable — the basis of ownership

---

## Exercises

**Exercise 1:** What are the min and max values of `i8`, `u8`, `i32`, and `u32`? Verify with Rust: `println!("{}", i8::MAX)`.

**Exercise 2:** Create a tuple `(name, age, height)` where name is `&str`, age is `u32`, and height is `f64`. Destructure it and print each field.

**Exercise 3:** Create an array of 5 temperatures in Celsius. Write a loop that converts each to Fahrenheit (`F = C * 9/5 + 32`) and prints both.

**Exercise 4:** Using shadowing, take a string `" 42 "`, shadow it to its trimmed version, then shadow it again to a parsed `i32`.

**Exercise 5:** Look up `std::mem::size_of_val` and use it to print the size of a `String` vs the size of the string's content on the heap.

---

*Next: [Chapter 4 — Control Flow](04-control-flow.md)*
