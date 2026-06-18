# Chapter 91 — Common Pitfalls and How to Fix Them

---

## Pitfall 1: Moving Into a Loop

**Problem:** Moving a value that should be used in every iteration.

```rust
// WRONG — String moved on first iteration
let s = String::from("hello");
for i in 0..5 {
    println!("{} {}", s, i);  // ERROR: use of moved value
    consume_string(s);        // s moved here
}

fn consume_string(s: String) {}
```

**Fix:** Clone, borrow, or restructure:

```rust
// Fix 1: borrow instead of move
fn use_string(s: &str) {}
let s = String::from("hello");
for i in 0..5 {
    use_string(&s);  // borrows, doesn't move
}

// Fix 2: clone if you need ownership each iteration
for i in 0..5 {
    consume_string(s.clone());  // explicit clone each iteration
}
println!("{}", s);  // still valid
```

---

## Pitfall 2: Using a Value After Move

**Problem:** Trying to use a value that was moved.

```rust
let v1 = vec![1, 2, 3];
let v2 = v1;           // v1 moved to v2
println!("{:?}", v1);  // ERROR: value borrowed here after move
```

**Fix:** Use `clone()` when you need two independent copies, or reorganize code:

```rust
let v1 = vec![1, 2, 3];
let v2 = v1.clone();   // explicit deep copy
println!("{:?}", v1);  // still valid
println!("{:?}", v2);
```

---

## Pitfall 3: Mutating a Vec While Iterating

**Problem:** Cannot hold a reference into a Vec while modifying it (the Vec might reallocate).

```rust
let mut v = vec![1, 2, 3, 4, 5];
let first = &v[0];  // immutable borrow
v.push(6);          // ERROR: cannot borrow `v` as mutable because it's borrowed as immutable
println!("{}", first);
```

**Fix:** Either finish using the reference first, or use indices:

```rust
// Fix 1: don't hold reference across mutation
let mut v = vec![1, 2, 3];
let first_val = v[0];  // copy the value (i32 is Copy)
v.push(4);
println!("{}", first_val);  // OK — we have a copy, not a reference

// Fix 2: separate the operations
let mut v = vec![1, 2, 3, 4, 5];
{
    let first = &v[0];
    println!("{}", first);
}  // reference dropped here
v.push(6);  // now OK
```

---

## Pitfall 4: The Lifetime of Temporaries

**Problem:** Temporary created in a chain drops before the binding is used.

```rust
// WRONG — the String returned by to_string() is a temporary
// that is dropped immediately
let s: &str = String::from("hello").as_str();  // COMPILER ERROR
// (Actually Rust catches this — it won't compile)

// Less obvious case:
fn get_string() -> String { String::from("hello") }
let r: &str = &get_string();  // ERROR in some contexts
```

**Fix:** Bind the owned value first:

```rust
let owned = get_string();   // bind the String
let r: &str = &owned;       // then take a reference
println!("{}", r);
```

---

## Pitfall 5: Shadowing vs Mutation Confusion

**Problem:** Thinking shadowing mutates when it creates a new binding.

```rust
let x = 5;
let x = x + 1;  // shadowing — creates NEW x = 6
let x = x * 2;  // shadowing — creates NEW x = 12

// The old x values are gone. This is NOT the same as:
let mut x = 5;
x = x + 1;  // mutation
x = x * 2;  // mutation
```

**When shadowing causes bugs:**

```rust
fn main() {
    let result = compute();
    let result = result.unwrap();  // OK — shadowing changes type from Option to i32

    // But if you shadow a value you needed later:
    let data = expensive_operation();
    let data = process(&data);   // data is shadowed — expensive_operation's data gone
    println!("{:?}", data);      // only has processed result, not original
}
```

---

## Pitfall 6: String Indexing Panic

**Problem:** Assuming strings can be indexed by character position.

```rust
let s = "hello";
// let ch = s[0];    // ERROR: cannot index into &str

let s = "héllo";  // 'é' is a 2-byte UTF-8 character
let slice = &s[0..2];  // You might expect "hé" — you get "h" and first byte of 'é'
// If the slice cuts a multi-byte char: PANIC at runtime
```

**Fix:** Use char-aware methods:

```rust
let s = "héllo";

// Get Nth character safely
let third: Option<char> = s.chars().nth(2);
println!("{:?}", third);  // Some('l')

// Iterate over characters
for c in s.chars() {
    println!("{}", c);
}

// Byte slicing only works if you know the byte boundaries
```

---

## Pitfall 7: Forgetting to Handle Errors (unwrap in Library Code)

**Problem:** Using `.unwrap()` that panics on errors in code others depend on.

```rust
// NEVER in library code
pub fn read_config(path: &str) -> Config {
    let content = std::fs::read_to_string(path).unwrap();  // panic!
    parse_config(&content).unwrap()                         // panic!
}
```

**Fix:** Return Result and let the caller decide:

```rust
pub fn read_config(path: &str) -> Result<Config, ConfigError> {
    let content = std::fs::read_to_string(path)
        .map_err(|e| ConfigError::Io(e))?;
    parse_config(&content)
}
```

---

## Pitfall 8: Rc Cycles — Memory Leaks

**Problem:** `Rc<T>` cycles prevent reference count from reaching zero — memory leak.

```rust
use std::rc::Rc;
use std::cell::RefCell;

#[derive(Debug)]
struct Node {
    value: i32,
    next: Option<Rc<RefCell<Node>>>,
}

fn main() {
    let a = Rc::new(RefCell::new(Node { value: 1, next: None }));
    let b = Rc::new(RefCell::new(Node { value: 2, next: None }));

    // Create a cycle: a → b → a
    a.borrow_mut().next = Some(Rc::clone(&b));
    b.borrow_mut().next = Some(Rc::clone(&a));

    // Neither a nor b will ever be freed — reference counts never reach 0
}
```

**Fix:** Use `Weak<T>` for back-references:

```rust
use std::rc::{Rc, Weak};
use std::cell::RefCell;

struct Node {
    value: i32,
    parent: Option<Weak<RefCell<Node>>>,  // weak reference — doesn't own
    children: Vec<Rc<RefCell<Node>>>,     // strong reference — owns
}
```

---

## Pitfall 9: Integer Overflow in Debug vs Release

**Problem:** Overflow panics in debug but wraps silently in release.

```rust
fn main() {
    let x: u8 = 255;
    let y = x + 1;  // PANIC in debug; wraps to 0 in release
}
```

**Fix:** Use explicit overflow methods:

```rust
fn main() {
    let x: u8 = 255;

    let wrapped = x.wrapping_add(1);       // 0 — explicit wrapping
    let saturated = x.saturating_add(1);   // 255 — saturates at max
    let checked = x.checked_add(1);        // None — returns Option

    println!("{} {} {:?}", wrapped, saturated, checked);
}
```

---

## Pitfall 10: Deadlock With Mutex

**Problem:** Holding a MutexGuard while trying to acquire another lock.

```rust
use std::sync::Mutex;

fn main() {
    let m = Mutex::new(5);

    let guard = m.lock().unwrap();

    // DEADLOCK — guard holds the lock, this tries to acquire it again
    let guard2 = m.lock().unwrap();  // blocks forever
}
```

**Fix:** Drop the guard before re-acquiring, or use a different structure:

```rust
fn main() {
    let m = Mutex::new(5);

    {
        let guard = m.lock().unwrap();
        println!("{}", *guard);
    }  // guard dropped here — lock released

    let guard2 = m.lock().unwrap();  // OK
    println!("{}", *guard2);
}
```

---

## Pitfall 11: Confusing iter(), iter_mut(), into_iter()

**Problem:** Picking the wrong iteration method causes ownership issues.

```rust
let v = vec![String::from("a"), String::from("b")];

// into_iter() — MOVES v, yields String
for s in v.into_iter() { println!("{}", s); }
// println!("{:?}", v);  // ERROR: v moved

// v is gone. If you needed v after, use iter() instead.
let v = vec![String::from("a"), String::from("b")];
for s in v.iter() { println!("{}", s); }  // borrows v, yields &String
println!("{:?}", v);  // OK
```

**Quick Reference:**
```rust
v.iter()        → yields &T       (borrows v)
v.iter_mut()    → yields &mut T   (mutably borrows v)
v.into_iter()   → yields T        (moves v — v is gone)
&v              → auto into_iter() → yields &T
&mut v          → auto into_iter() → yields &mut T
```

---

## Pitfall 12: Holding a Lock Across an Await

**Problem:** Holding a `Mutex` lock across an `.await` point blocks other tasks.

```rust
// BAD (in async code)
use tokio::sync::Mutex;

async fn bad(mutex: &Mutex<i32>) {
    let guard = mutex.lock().await;
    some_async_operation().await;  // holding lock while suspended!
    println!("{}", *guard);
}

// GOOD — release lock before awaiting
async fn good(mutex: &Mutex<i32>) {
    let value = {
        let guard = mutex.lock().await;
        *guard  // copy the value (if Copy)
    };  // guard dropped — lock released
    some_async_operation().await;  // no lock held
    println!("{}", value);
}
```

---

## Pitfall 13: Not Using clippy

Many of the above pitfalls are caught by `cargo clippy`:

```bash
cargo clippy
# Catches: unnecessary clones, wrong iteration methods, 
# redundant patterns, unwrap in production code, etc.

# Run clippy in CI to prevent regressions:
cargo clippy -- -D warnings
```

---

## Pitfall 14: Forgetting that String != &str

```rust
// This looks like it should work:
fn greet(name: String) {}

fn main() {
    greet("hello");  // ERROR: expected String, found &str
}

// Fix: accept &str (more flexible) or convert explicitly
fn greet(name: &str) {}   // accepts &str, &String, ...
fn greet(name: String) {}  // requires: greet("hello".to_string())
```

---

## Pitfall 15: Implementing From Without Implementing Into

```rust
// BAD — implementing both manually (redundant)
impl From<i32> for MyType { fn from(x: i32) -> Self { MyType(x) } }
impl Into<MyType> for i32 { fn into(self) -> MyType { MyType(self) } }

// GOOD — Into is automatically derived from From
impl From<i32> for MyType { fn from(x: i32) -> Self { MyType(x) } }
// MyType::from(5) works
// 5i32.into() also works automatically
```

---

## Summary

The most common Rust pitfalls fall into three categories:

**Ownership mistakes**: moving instead of borrowing, using after move, moving in loops  
**Borrowing errors**: references outliving data, simultaneous mutable + immutable borrows, holding references across mutations  
**Lifetime confusion**: temporaries dropping too early, returning references to locals

All of these are caught at compile time — Rust won't let you ship these bugs. The error messages guide you to the fix. `cargo clippy` catches even more issues before you encounter them at runtime.

---

*Next: [Chapter 92 — Interview Prep](92-interview-prep.md)*
