# Chapter 6 — Ownership: Rust's Superpower

> *"The ownership system is the most innovative aspect of Rust. Once you understand it, everything clicks."*

This chapter is the heart of Rust. Take your time. Read slowly. Run every example. Ownership is unlike anything in other languages, but once it clicks, you'll understand exactly what Rust's compiler is doing and why.

---

## 6.1 Why Ownership Exists

Every program must manage memory. The computer has a finite amount of RAM. When you allocate memory to store data, you need to eventually release it. Failing to do so is a **memory leak**. Releasing memory too early — and then accessing it — is a **use-after-free** bug, one of the most dangerous security vulnerabilities.

Three approaches exist:

**1. Garbage Collection (Java, Go, Python)**
```
Your code → allocates memory freely
GC → periodically scans for unreachable memory → frees it
```
- Pro: Easy for the programmer
- Con: Unpredictable pauses, runtime overhead, not usable for systems/embedded

**2. Manual Management (C, C++)**
```
malloc(size) → allocate
free(ptr)    → deallocate (programmer's job)
```
- Pro: Maximum control and performance
- Con: Use-after-free, double-free, memory leaks — catastrophic bugs

**3. Ownership (Rust)**
```
The compiler statically tracks every allocation.
When the owner goes out of scope → memory is automatically freed.
No GC. No manual free(). Decided at compile time.
```
- Pro: Zero runtime overhead, memory safe
- Con: Requires learning a new programming model

---

## 6.2 The Three Ownership Rules

Every piece of heap-allocated data in Rust follows exactly three rules, always:

```
Rule 1: Each value in Rust has an owner.
Rule 2: There can only be one owner at a time.
Rule 3: When the owner goes out of scope, the value will be dropped.
```

These three rules, enforced by the compiler, eliminate the entire class of memory safety bugs.

---

## 6.3 Move Semantics

### Copying Integers (Stack Data)

```rust
fn main() {
    let x = 5;
    let y = x;  // x is COPIED — both x and y are 5
    println!("x={}, y={}", x, y);  // OK — both still valid
}
```

This works because `i32` is a **Copy** type. Stack data with known fixed size is cheap to copy. A full duplicate is made.

### Moving Strings (Heap Data)

```rust
fn main() {
    let s1 = String::from("hello");
    let s2 = s1;  // s1 is MOVED to s2

    // println!("{}", s1);  // COMPILE ERROR: value borrowed here after move
    println!("{}", s2);  // OK
}
```

Why doesn't s1 work? Let's visualize what happens in memory:

```
Before move:
Stack:          Heap:
s1:
┌──────────┐   ┌───────────────┐
│ ptr ─────┼──►│ h e l l o    │
│ len: 5   │   └───────────────┘
│ cap: 5   │
└──────────┘

After `let s2 = s1;` (MOVE — not copy):
Stack:          Heap:
s1: (invalidated — compiler marks it as moved)
s2:
┌──────────┐   ┌───────────────┐
│ ptr ─────┼──►│ h e l l o    │
│ len: 5   │   └───────────────┘
│ cap: 5   │
└──────────┘
```

**Why move instead of copy?**

Option A: Copy the entire heap data — expensive for large strings.  
Option B: Both s1 and s2 point to the same heap data — when one is dropped, the other has a dangling pointer (double-free bug).  
Option C: **Move** — transfer ownership. s1 is now invalid. Only s2 owns the data. When s2 is dropped, the heap is freed exactly once.

Rust chose Option C. The compiler marks s1 as "moved" — accessing it after the move is a compile error.

### Move in Assignment

```rust
fn main() {
    let s1 = String::from("hello");
    let s2 = String::from("world");

    let s3 = s1;  // s1 moved to s3
    // s1 is now invalid

    // s2 is still valid — it wasn't moved
    let s4 = s2;  // s2 moved to s4
    // s2 is now invalid

    println!("{}", s3);  // OK
    println!("{}", s4);  // OK
}
```

---

## 6.4 Clone and Copy

### Clone — Explicit Deep Copy

When you genuinely need two independent copies of heap data:

```rust
fn main() {
    let s1 = String::from("hello");
    let s2 = s1.clone();  // deep copy — new heap allocation, new data

    println!("s1={}, s2={}", s1, s2);  // both valid
}
```

After `clone()`:
```
Stack:          Heap:
s1:
┌──────────┐   ┌───────────────┐
│ ptr ─────┼──►│ h e l l o    │
│ len: 5   │   └───────────────┘
│ cap: 5   │
└──────────┘
s2:
┌──────────┐   ┌───────────────┐
│ ptr ─────┼──►│ h e l l o    │  (new allocation)
│ len: 5   │   └───────────────┘
│ cap: 5   │
└──────────┘
```

`clone()` is always explicit in Rust — unlike Java where `=` can hide expensive copies. Seeing `.clone()` in code tells you: "this is doing a heap allocation."

### The Copy Trait

Types that implement the `Copy` trait are **automatically copied** on assignment instead of moved. These are always stack-allocated types with known fixed size:

```rust
// Copy types — automatically copied, not moved
let x: i32 = 5;
let y = x;          // COPY — x still valid
println!("{}", x);  // OK

let a: bool = true;
let b = a;          // COPY
println!("{}", a);  // OK

let c: char = 'z';
let d = c;          // COPY
println!("{}", c);  // OK

let tup = (1, 2.0, 'x');
let tup2 = tup;    // COPY — all elements are Copy types
println!("{:?}", tup);  // OK

// Non-Copy types — moved, not copied
let s = String::from("hello");
let s2 = s;         // MOVE
// println!("{}", s);  // ERROR — moved
```

**Types that implement Copy:**
- All integer types: `i8` through `i128`, `u8` through `u128`
- `f32`, `f64`
- `bool`, `char`
- Tuples — only if ALL elements are `Copy`
- Arrays — only if element type is `Copy`
- Raw pointers `*const T`, `*mut T`
- References `&T` (shared references, not mutable)

**Types that do NOT implement Copy:**
- `String` (heap data)
- `Vec<T>` (heap data)
- `Box<T>` (heap data)
- Any type containing non-Copy types

### Implementing Copy for Your Types

```rust
#[derive(Debug, Clone, Copy)]  // derive both Clone and Copy
struct Point {
    x: f64,
    y: f64,  // f64 is Copy, so Point can be Copy
}

fn main() {
    let p1 = Point { x: 1.0, y: 2.0 };
    let p2 = p1;  // COPY — p1 still valid
    println!("{:?}", p1);  // OK
    println!("{:?}", p2);  // OK
}

// This would NOT work — String is not Copy
// #[derive(Clone, Copy)]
// struct NamedPoint {
//     name: String,  // String is not Copy — compile error
//     x: f64,
// }
```

---

## 6.5 The drop Function

When a value goes out of scope, Rust automatically calls its `drop` function, which frees the associated resources:

```rust
fn main() {
    let s1 = String::from("first");   // s1 allocated

    {
        let s2 = String::from("second");  // s2 allocated
        println!("{}", s2);
    }  // ← s2 goes out of scope here: drop(s2) called, "second" freed

    println!("{}", s1);

}  // ← s1 goes out of scope here: drop(s1) called, "first" freed
```

This is called **RAII** (Resource Acquisition Is Initialization) — the same pattern used in C++. Resources are always tied to the lifetime of an object.

### Implementing Drop for Your Types

```rust
struct ResourceHandle {
    name: String,
}

impl Drop for ResourceHandle {
    fn drop(&mut self) {
        println!("Releasing resource: {}", self.name);
        // cleanup code — close file handles, network connections, etc.
    }
}

fn main() {
    let h1 = ResourceHandle { name: String::from("Database connection") };
    let h2 = ResourceHandle { name: String::from("File handle") };

    println!("Resources acquired");
    // h2 dropped first (LIFO — last acquired, first released)
    // h1 dropped second
}
// Output:
// Resources acquired
// Releasing resource: File handle
// Releasing resource: Database connection
```

### Explicitly Dropping Early

```rust
fn main() {
    let h = ResourceHandle { name: String::from("Lock") };
    println!("Lock acquired");

    drop(h);  // explicitly drop before end of scope
    println!("Lock released early");

    // h cannot be used here — it's been dropped
    // println!("{}", h.name);  // COMPILE ERROR: use of moved value
}
```

---

## 6.6 Ownership and Functions

### Passing Values to Functions

```rust
fn takes_ownership(s: String) {  // s comes into scope
    println!("{}", s);
}  // s goes out of scope — drop called, heap freed

fn makes_copy(x: i32) {   // x comes into scope (COPY)
    println!("{}", x);
}  // x goes out of scope — but i32 is on stack, nothing special happens

fn main() {
    let s = String::from("hello");
    takes_ownership(s);          // s is MOVED into the function
    // println!("{}", s);        // COMPILE ERROR — s was moved

    let x = 5;
    makes_copy(x);               // x is COPIED — x is still valid
    println!("{}", x);           // OK
}
```

### Returning Values Transfers Ownership

```rust
fn gives_ownership() -> String {
    let s = String::from("hello");
    s  // s is MOVED out of the function — caller becomes owner
}

fn takes_and_gives_back(s: String) -> String {
    s  // takes ownership, then returns it
}

fn main() {
    let s1 = gives_ownership();       // s1 receives ownership
    let s2 = String::from("world");
    let s3 = takes_and_gives_back(s2); // s2 moved to fn, then moved to s3
    // s2 invalid here. s1 and s3 valid.
    println!("{}", s1);
    println!("{}", s3);
}
```

### The Problem: Ownership Dance

```rust
fn calculate_length(s: String) -> (String, usize) {
    let length = s.len();
    (s, length)  // return both to give ownership back
}

fn main() {
    let s1 = String::from("hello");
    let (s2, len) = calculate_length(s1);
    println!("'{}' has length {}", s2, len);
}
```

This is clunky. You must return the String back just to use it again. The solution is **borrowing** — covered in the next chapter. But first, understand WHY this problem exists: once you move a value into a function, you've transferred ownership. The function owns it. If you want it back, you have to return it.

---

## Complete Example: Understanding Ownership Through Memory

```rust
fn main() {
    // Stack variable — Copy type, trivial
    let n: i32 = 42;
    let m = n;  // copy — both n and m are 42
    println!("n={}, m={}", n, m);

    // Heap variable — move semantics
    let s1 = String::from("ownership");

    // Move: s1 → s2
    let s2 = s1;

    // Use s2, not s1
    println!("s2 = {}", s2);
    println!("length = {}", s2.len());

    // Clone: make an independent copy
    let s3 = s2.clone();
    println!("s2={}, s3={}", s2, s3);  // both valid

    // Pass to function — s3 is moved
    print_string(s3);
    // println!("{}", s3);  // ERROR — moved into print_string

    // Pass by reference — no move (next chapter)
    print_length(&s2);
    println!("{}", s2);  // still valid!

}  // s2 is dropped here — heap freed

fn print_string(s: String) {
    println!("Function received: {}", s);
}  // s is dropped here

fn print_length(s: &String) {
    println!("Length is: {}", s.len());
    // s is a reference — it doesn't own the String
}  // s (the reference) goes out of scope — but the String is NOT freed
```

---

## Summary

Ownership is Rust's compile-time mechanism for managing heap memory without a garbage collector. Every heap value has exactly one owner. When the owner goes out of scope, the value is dropped. Assignment moves heap data — the original variable becomes invalid. Stack data (Copy types: integers, booleans, chars, fixed tuples/arrays of Copy types) is automatically copied on assignment. `clone()` creates an explicit, expensive deep copy. Functions take ownership of their parameters unless you pass references (borrowing, next chapter).

---

## Key Takeaways

- **Rule 1**: Each value has exactly one owner
- **Rule 2**: Only one owner at a time
- **Rule 3**: Owner goes out of scope → value is dropped
- Heap data (String, Vec): assignment **moves**, original invalidated
- Stack data (i32, bool, etc.): assignment **copies**, original still valid
- `clone()` performs an explicit deep copy — avoid when borrowing works
- `drop()` is called automatically when owner goes out of scope — RAII
- Passing to a function = moving — caller loses ownership unless returned

---

## Exercises

**Exercise 1:** Without running it, determine which lines cause compile errors and why:
```rust
let s1 = String::from("hello");
let s2 = s1;
let s3 = s1.clone();
println!("{} {} {}", s1, s2, s3);
```

**Exercise 2:** Write a function `first_word(s: String) -> String` that returns the first word. What's the problem with this signature? (Hint: think about ownership of the input.)

**Exercise 3:** Why can you do `let y = x` for `i32` but not for `String` without moving? What trait controls this?

**Exercise 4:** Implement `Drop` for a struct `DatabaseConnection { url: String }` that prints "Closing connection to {url}" when dropped. Verify it's called correctly with a nested scope.

**Exercise 5:** Draw a memory diagram (stack + heap) for this code at each step:
```rust
let s1 = String::from("abc");
let s2 = s1.clone();
let s3 = s2;
```

---

*Next: [Chapter 7 — Borrowing](07-borrowing.md)*
