# Chapter 92 — Interview Preparation

---

## Conceptual Questions and Answers

### Q1: What is ownership in Rust?

**Answer:** Ownership is Rust's compile-time mechanism for managing memory without garbage collection or manual deallocation. The rules are:
1. Every value has exactly one owner
2. Only one owner at a time
3. When the owner goes out of scope, the value is dropped

Ownership applies to heap-allocated data. Stack data (integers, booleans) is `Copy` — it's trivially duplicated on assignment. Heap data is **moved** on assignment — the original binding becomes invalid, preventing double-free and use-after-free bugs.

---

### Q2: What's the difference between borrowing and ownership?

**Answer:** Borrowing lets you use a value without taking ownership. A borrow is a reference (`&T` or `&mut T`) to data owned by someone else. Key rules:
- Any number of immutable borrows (`&T`) can coexist
- At most one mutable borrow (`&mut T`) can exist at a time
- Immutable and mutable borrows cannot coexist
- References must not outlive their source

This prevents data races at compile time — the same guarantee mutexes provide at runtime.

---

### Q3: What are lifetimes, and when do you need explicit annotations?

**Answer:** Lifetimes are compile-time metadata that track how long references are valid. They prevent dangling references. The compiler automatically infers lifetimes via **elision rules** in most cases.

Explicit annotations are needed when:
1. A function returns a reference that could come from multiple inputs
2. A struct holds references (the struct gets a lifetime parameter)

Example: `fn longest<'a>(x: &'a str, y: &'a str) -> &'a str` tells the compiler: the returned reference is valid for the shorter of x's and y's lifetimes.

---

### Q4: What is the difference between String and &str?

**Answer:**
- `String`: owned, heap-allocated, mutable, variable-length. Think of it as `Vec<u8>` with UTF-8 guarantees. You use `String` when you need to own or mutate text.
- `&str`: borrowed string slice — a fat pointer (pointer + length) to UTF-8 data stored somewhere. The data might be in the binary (string literals), in a `String` on the heap, or anywhere else. Immutable. You use `&str` for function parameters and read-only views.

Rule of thumb: function parameters take `&str`; fields/return types use `String`.

---

### Q5: Explain the difference between Rc and Arc.

**Answer:** Both are reference-counted smart pointers that allow multiple owners:
- `Rc<T>`: single-threaded. Uses non-atomic reference counting. Cannot be sent across threads (`!Send`). Cheaper than Arc.
- `Arc<T>`: thread-safe. Uses atomic reference counting. Implements `Send` and `Sync`. Slightly more expensive due to atomic operations.

Use `Rc` when you know data stays on one thread. Use `Arc` when ownership must cross thread boundaries.

---

### Q6: What is interior mutability?

**Answer:** Interior mutability allows mutating data through a shared reference (`&T`), bypassing Rust's normal borrowing rules. It moves borrow checking from compile time to runtime.

- `Cell<T>`: for `Copy` types. Get/set the value without runtime cost.
- `RefCell<T>`: for any type. Tracks borrows at runtime; panics on violation.
- `Mutex<T>`: thread-safe version — blocks instead of panicking.

Common pattern: `Rc<RefCell<T>>` for shared mutable data on a single thread; `Arc<Mutex<T>>` for multiple threads.

---

### Q7: What are zero-cost abstractions?

**Answer:** Rust's abstractions compile to the same machine code as equivalent hand-written low-level code. Examples:
- **Iterators**: `v.iter().filter(...).map(...).sum()` compiles to a single tight loop — no heap allocation, no dynamic dispatch
- **Closures**: inlined by the compiler when the type is known
- **Generics**: monomorphized — one specialized copy per type
- **RAII/Drop**: resource cleanup happens exactly when needed, at the compiler-determined drop point

"Zero-cost" means you don't pay for what you don't use, and you can't hand-write it faster.

---

### Q8: How does Rust prevent data races?

**Answer:** Through the type system — specifically the `Send` and `Sync` marker traits:
- A type is `Send` if it's safe to transfer ownership across threads
- A type is `Sync` if it's safe to share a reference across threads (`&T: Send if T: Sync`)
- `Rc<T>` is not `Send` (non-atomic ref count) — the compiler prevents sending it to another thread
- `RefCell<T>` is not `Sync` — prevents sharing a `&RefCell<T>` across threads
- `Mutex<T>` implements `Sync` — it's designed for concurrent access

These are not runtime checks — they're compile-time type constraints that reject data races before the code runs.

---

### Q9: What is the difference between panic! and Result?

**Answer:**
- `panic!`: for unrecoverable errors — programming bugs, invariant violations, index out of bounds. Crashes the thread with a stack trace. Not for expected failures.
- `Result<T, E>`: for recoverable errors — file not found, invalid input, network timeout. The caller must handle both `Ok` and `Err` cases. Propagate with `?`.

Rule: `panic!` for bugs (you made a mistake in the program), `Result` for errors (the world didn't cooperate).

---

### Q10: What is monomorphization, and what are its trade-offs?

**Answer:** Monomorphization is the process by which the compiler generates concrete implementations of generic code for each type it's used with:

```rust
fn largest<T: PartialOrd>(list: &[T]) -> &T { ... }
// Becomes:
fn largest_i32(list: &[i32]) -> &i32 { ... }
fn largest_f64(list: &[f64]) -> &f64 { ... }
```

**Benefits**: Zero runtime overhead, each version can be fully optimized, no virtual dispatch  
**Trade-offs**: Larger binary (code duplicated per type), longer compile times for heavily generic code

Alternative: `dyn Trait` uses dynamic dispatch — single code path, small overhead per call, smaller binary.

---

## Tricky Code Challenges

### Challenge 1: What does this print?

```rust
fn main() {
    let mut v = vec![1, 2, 3];
    let first = &v[0];
    v.push(4);
    println!("{}", first);
}
```

**Answer:** COMPILE ERROR. `first` holds an immutable reference into `v`. `v.push(4)` requires a mutable reference. These cannot coexist. Potential issue: Vec reallocation would invalidate `first`.

---

### Challenge 2: What's wrong?

```rust
fn main() {
    let s;
    {
        let temp = String::from("hello");
        s = &temp;
    }
    println!("{}", s);
}
```

**Answer:** COMPILE ERROR. `temp` is dropped at the end of the inner scope, but `s` holds a reference to it. The reference would dangle. The borrow checker rejects this.

---

### Challenge 3: Will this compile?

```rust
fn first_word(s: &str) -> &str {
    &s[..s.find(' ').unwrap_or(s.len())]
}

fn main() {
    let s = String::from("hello world");
    let word = first_word(&s);
    s.clear();
    println!("{}", word);
}
```

**Answer:** COMPILE ERROR. `word` borrows from `s`. `s.clear()` requires `&mut s`. A mutable borrow cannot coexist with `word`'s immutable borrow. Rust prevents using `word` after modifying the String it points into.

---

### Challenge 4: Fix the compilation error

```rust
fn double(v: Vec<i32>) -> Vec<i32> {
    v.iter().map(|x| x * 2).collect()
}

fn main() {
    let numbers = vec![1, 2, 3];
    let doubled = double(numbers);
    println!("{:?}", numbers);  // ERROR
    println!("{:?}", doubled);
}
```

**Answer:** `numbers` is moved into `double`. Fix options:
```rust
// Option 1: Change function to take a slice (better API)
fn double(v: &[i32]) -> Vec<i32> {
    v.iter().map(|&x| x * 2).collect()
}
// Now numbers is borrowed, not moved.

// Option 2: Clone before passing
let doubled = double(numbers.clone());
```

---

### Coding Problem 1: Reverse Words

Write a function that takes a sentence and returns it with words in reverse order, without allocating a new string per word:

```rust
fn reverse_words(s: &str) -> String {
    s.split_whitespace()
        .rev()
        .collect::<Vec<&str>>()
        .join(" ")
}

fn main() {
    println!("{}", reverse_words("hello world rust"));
    // Output: rust world hello
}
```

---

### Coding Problem 2: Implement a Stack

```rust
struct Stack<T> {
    data: Vec<T>,
}

impl<T> Stack<T> {
    fn new() -> Self {
        Stack { data: Vec::new() }
    }

    fn push(&mut self, item: T) {
        self.data.push(item);
    }

    fn pop(&mut self) -> Option<T> {
        self.data.pop()
    }

    fn peek(&self) -> Option<&T> {
        self.data.last()
    }

    fn is_empty(&self) -> bool {
        self.data.is_empty()
    }

    fn size(&self) -> usize {
        self.data.len()
    }
}

fn main() {
    let mut stack = Stack::new();
    stack.push(1);
    stack.push(2);
    stack.push(3);
    println!("{:?}", stack.peek());  // Some(3)
    println!("{:?}", stack.pop());   // Some(3)
    println!("{}", stack.size());    // 2
}
```

---

### Coding Problem 3: Count Word Frequencies

```rust
use std::collections::HashMap;

fn word_frequency(text: &str) -> Vec<(String, usize)> {
    let mut freq: HashMap<&str, usize> = HashMap::new();
    for word in text.split_whitespace() {
        *freq.entry(word).or_insert(0) += 1;
    }
    let mut result: Vec<(String, usize)> = freq
        .into_iter()
        .map(|(w, c)| (w.to_string(), c))
        .collect();
    result.sort_by(|a, b| b.1.cmp(&a.1).then(a.0.cmp(&b.0)));
    result
}

fn main() {
    let text = "the quick brown fox jumps over the lazy dog the fox";
    for (word, count) in word_frequency(text).iter().take(5) {
        println!("{}: {}", word, count);
    }
}
```

---

### Coding Problem 4: Binary Search

```rust
fn binary_search<T: Ord>(slice: &[T], target: &T) -> Option<usize> {
    let mut low = 0;
    let mut high = slice.len();

    while low < high {
        let mid = low + (high - low) / 2;
        match slice[mid].cmp(target) {
            std::cmp::Ordering::Equal => return Some(mid),
            std::cmp::Ordering::Less => low = mid + 1,
            std::cmp::Ordering::Greater => high = mid,
        }
    }
    None
}

fn main() {
    let nums = vec![1, 3, 5, 7, 9, 11, 13];
    println!("{:?}", binary_search(&nums, &7));   // Some(3)
    println!("{:?}", binary_search(&nums, &6));   // None
}
```

---

### Coding Problem 5: Thread-Safe Counter

```rust
use std::sync::{Arc, Mutex};
use std::thread;

fn parallel_count(values: Vec<i32>, threshold: i32) -> usize {
    let count = Arc::new(Mutex::new(0usize));
    let values = Arc::new(values);
    let mut handles = vec![];

    let chunk_size = (values.len() + 3) / 4;
    for chunk_start in (0..values.len()).step_by(chunk_size.max(1)) {
        let values = Arc::clone(&values);
        let count = Arc::clone(&count);
        handles.push(thread::spawn(move || {
            let local_count = values[chunk_start..]
                .iter()
                .take(chunk_size)
                .filter(|&&x| x > threshold)
                .count();
            *count.lock().unwrap() += local_count;
        }));
    }

    for h in handles { h.join().unwrap(); }
    *count.lock().unwrap()
}

fn main() {
    let values = (1..=100).collect::<Vec<i32>>();
    println!("{}", parallel_count(values, 50));  // 50
}
```

---

## Key Facts to Memorize

- **Ownership rules**: one owner, one at a time, dropped when out of scope
- **Borrowing rules**: any number of `&T` OR exactly one `&mut T`
- **Copy types**: i8..i128, u8..u128, f32, f64, bool, char, tuples/arrays of Copy types
- **Move types**: String, Vec, Box, anything heap-allocated
- **Lifetime elision**: 1 input → output gets same; `&self` method → output gets self's lifetime
- **Send**: ownership can cross threads; **Sync**: reference can cross threads
- **Rc**: single-thread multiple ownership; **Arc**: multi-thread; **RefCell**: runtime borrow checking; **Mutex**: thread-safe RefCell
- **? operator**: propagates Err/None up the call stack, converts via From
- **dyn Trait**: 16-byte fat pointer (data + vtable); **impl Trait**: static dispatch, monomorphized
- **String**: owned heap UTF-8; **&str**: borrowed slice; deref coercion converts String → &str

---

*Next: [Chapter 99 — HTTP Server Project](99-http-server-project.md)*
