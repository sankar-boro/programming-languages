# Chapter 13 — Advanced Features: Closures, Smart Pointers, Unsafe

---

## 13.1 Closures — Deep Dive

Closures capture their environment. Understanding *how* they capture determines which traits they implement.

### Three Closure Traits

```rust
// FnOnce — can be called once (might move captured values out)
// FnMut  — can be called multiple times, may mutate captured values
// Fn     — can be called multiple times, doesn't mutate captured values
//
// Fn ⊂ FnMut ⊂ FnOnce
// Every Fn implements FnMut, every FnMut implements FnOnce
```

### How Closures Capture

```rust
fn main() {
    let s = String::from("hello");

    // Capture by reference — Fn
    let borrow_closure = || println!("{}", s);  // borrows s
    borrow_closure();
    borrow_closure();  // can call multiple times
    println!("{}", s);  // s still valid

    // Capture by mutable reference — FnMut
    let mut count = 0;
    let mut increment = || { count += 1; };  // borrows count mutably
    increment();
    increment();
    // println!("{}", count);  // ERROR — count is mutably borrowed by closure

    // Capture by move — move keyword
    let s = String::from("hello");
    let moved = move || println!("{}", s);  // s is MOVED into the closure
    moved();
    // println!("{}", s);  // ERROR — s was moved

    // FnOnce — closure that consumes what it captures
    let s = String::from("hello");
    let consume = move || {
        let s2 = s;  // s moved out of closure — can only call once
        println!("{}", s2);
    };
    consume();
    // consume();  // ERROR — s was moved in first call
}
```

### move Closures — When to Use

```rust
use std::thread;

fn main() {
    let s = String::from("hello");

    // Threads need ownership — can't borrow across thread boundaries
    let handle = thread::spawn(move || {
        println!("Thread: {}", s);  // s moved into thread
    });

    handle.join().unwrap();
    // println!("{}", s);  // ERROR — moved into thread
}
```

### Returning Closures

```rust
// Must use Box<dyn Fn> — closures have different sizes
fn make_adder(x: i32) -> Box<dyn Fn(i32) -> i32> {
    Box::new(move |y| x + y)
}

// With impl Fn (simpler when possible)
fn make_adder_v2(x: i32) -> impl Fn(i32) -> i32 {
    move |y| x + y
}

fn make_multiplier(factor: i32) -> impl FnMut(i32) -> i32 {
    let mut call_count = 0;
    move |x| {
        call_count += 1;
        println!("Called {} times", call_count);
        x * factor
    }
}

fn main() {
    let add5 = make_adder(5);
    println!("{}", add5(10));   // 15
    println!("{}", add5(20));   // 25

    let mut triple = make_multiplier(3);
    triple(10);  // Called 1 times
    triple(10);  // Called 2 times
}
```

---

## 13.2 Iterators — Advanced

### Iterator Adapters Chain

```rust
fn main() {
    // Complex iterator pipelines
    let words = vec!["hello", "world", "rust", "is", "awesome"];

    let result: Vec<String> = words.iter()
        .filter(|w| w.len() > 3)       // keep words longer than 3 chars
        .map(|w| w.to_uppercase())      // uppercase each
        .enumerate()                    // add index
        .map(|(i, w)| format!("{}:{}", i, w))  // format with index
        .take(3)                        // take first 3
        .collect();

    println!("{:?}", result);
}
```

### scan — Stateful Map

```rust
fn main() {
    // scan is like fold but emits intermediate values
    let v = vec![1, 2, 3, 4, 5];

    // Running sum
    let running_sum: Vec<i32> = v.iter()
        .scan(0, |acc, &x| {
            *acc += x;
            Some(*acc)
        })
        .collect();
    println!("{:?}", running_sum);  // [1, 3, 6, 10, 15]
}
```

### Implementing IntoIterator

```rust
struct Grid {
    rows: Vec<Vec<i32>>,
}

struct GridIterator<'a> {
    grid: &'a Grid,
    row: usize,
    col: usize,
}

impl<'a> Iterator for GridIterator<'a> {
    type Item = &'a i32;

    fn next(&mut self) -> Option<&'a i32> {
        if self.row >= self.grid.rows.len() {
            return None;
        }
        let item = &self.grid.rows[self.row][self.col];
        self.col += 1;
        if self.col >= self.grid.rows[self.row].len() {
            self.col = 0;
            self.row += 1;
        }
        Some(item)
    }
}

impl<'a> IntoIterator for &'a Grid {
    type Item = &'a i32;
    type IntoIter = GridIterator<'a>;

    fn into_iter(self) -> GridIterator<'a> {
        GridIterator { grid: self, row: 0, col: 0 }
    }
}

fn main() {
    let grid = Grid { rows: vec![vec![1, 2, 3], vec![4, 5, 6]] };
    let sum: i32 = grid.into_iter().sum();
    println!("{}", sum);  // 21
}
```

---

## 13.3 Box<T> — Heap Allocation

`Box<T>` allocates `T` on the heap. Use it when:
1. You have a type whose size isn't known at compile time (recursive types, trait objects)
2. You want to transfer ownership of a large amount of data

```rust
fn main() {
    // Box for a heap-allocated i32 (rarely needed, just for illustration)
    let b = Box::new(5);
    println!("{}", b);   // auto-dereferences: prints 5
    println!("{}", *b);  // explicit deref

    // Box is dropped when it goes out of scope — heap freed automatically
}
```

### Recursive Types with Box

Without `Box`, recursive types have infinite size:

```rust
// ERROR: recursive type `List` has infinite size
// enum List {
//     Cons(i32, List),  // List contains List, which contains List...
//     Nil,
// }

// CORRECT: Box breaks the infinite size
#[derive(Debug)]
enum List {
    Cons(i32, Box<List>),  // heap pointer — fixed 8 bytes
    Nil,
}

fn main() {
    let list = List::Cons(1,
        Box::new(List::Cons(2,
            Box::new(List::Cons(3,
                Box::new(List::Nil))))));
    println!("{:?}", list);
}
```

### Box as a Trait Object

```rust
trait Animal {
    fn name(&self) -> &str;
    fn sound(&self) -> &str;
}

struct Dog;
struct Cat;
struct Cow;

impl Animal for Dog { fn name(&self) -> &str { "Dog" } fn sound(&self) -> &str { "Woof" } }
impl Animal for Cat { fn name(&self) -> &str { "Cat" } fn sound(&self) -> &str { "Meow" } }
impl Animal for Cow { fn name(&self) -> &str { "Cow" } fn sound(&self) -> &str { "Moo" }  }

fn make_animals() -> Vec<Box<dyn Animal>> {
    vec![Box::new(Dog), Box::new(Cat), Box::new(Cow)]
}

fn main() {
    for animal in make_animals() {
        println!("{} says {}", animal.name(), animal.sound());
    }
}
```

---

## 13.4 Rc<T> — Reference Counting

`Rc<T>` (Reference Counted) allows **multiple owners** of the same data. Uses reference counting to track when to free. **Single-threaded only.**

```rust
use std::rc::Rc;

fn main() {
    let a = Rc::new(5);
    println!("count: {}", Rc::strong_count(&a));  // 1

    let b = Rc::clone(&a);  // increments reference count (fast — no deep copy)
    println!("count: {}", Rc::strong_count(&a));  // 2

    {
        let c = Rc::clone(&a);
        println!("count: {}", Rc::strong_count(&a));  // 3
    }  // c dropped — count decrements

    println!("count: {}", Rc::strong_count(&a));  // 2
    // When count reaches 0, data is freed
}
```

### Rc with Shared Graph/Tree Nodes

```rust
use std::rc::Rc;

#[derive(Debug)]
struct Node {
    value: i32,
    children: Vec<Rc<Node>>,
}

fn main() {
    let leaf = Rc::new(Node { value: 3, children: vec![] });

    // Two parents sharing the same leaf node
    let parent1 = Rc::new(Node {
        value: 1,
        children: vec![Rc::clone(&leaf)],
    });
    let parent2 = Rc::new(Node {
        value: 2,
        children: vec![Rc::clone(&leaf)],
    });

    println!("leaf ref count: {}", Rc::strong_count(&leaf));  // 3
    println!("parent1 has child with value {}", parent1.children[0].value);
    println!("parent2 has child with value {}", parent2.children[0].value);
}
```

---

## 13.5 RefCell<T> — Interior Mutability

`RefCell<T>` allows mutating data through a shared reference, with borrow checking enforced at **runtime** instead of compile time. **Single-threaded only.**

```rust
use std::cell::RefCell;

fn main() {
    // RefCell: can have immutable outer binding but mutable inner data
    let data = RefCell::new(vec![1, 2, 3]);

    // borrow() — immutable borrow (runtime check)
    let b1 = data.borrow();
    println!("{:?}", *b1);
    drop(b1);  // release borrow

    // borrow_mut() — mutable borrow (runtime check)
    data.borrow_mut().push(4);
    println!("{:?}", data.borrow());

    // Runtime panic if borrow rules violated:
    // let b1 = data.borrow();
    // let b2 = data.borrow_mut();  // PANIC at runtime: already borrowed
}
```

### Rc<RefCell<T>> — Shared Mutable State

The classic pattern for shared mutable ownership in single-threaded code:

```rust
use std::rc::Rc;
use std::cell::RefCell;

fn main() {
    let shared = Rc::new(RefCell::new(0));

    let a = Rc::clone(&shared);
    let b = Rc::clone(&shared);

    *a.borrow_mut() += 10;
    *b.borrow_mut() += 20;

    println!("{}", shared.borrow());  // 30
}
```

---

## 13.6 Arc<T> and Mutex<T> — Thread-Safe Sharing

`Arc<T>` (Atomic Reference Counted) is like `Rc<T>` but safe to share across threads. `Mutex<T>` provides mutual exclusion:

```rust
use std::sync::{Arc, Mutex};
use std::thread;

fn main() {
    let counter = Arc::new(Mutex::new(0));
    let mut handles = vec![];

    for _ in 0..10 {
        let counter = Arc::clone(&counter);
        let handle = thread::spawn(move || {
            let mut num = counter.lock().unwrap();  // acquire lock
            *num += 1;
        });  // lock released when `num` is dropped
        handles.push(handle);
    }

    for handle in handles {
        handle.join().unwrap();
    }

    println!("Result: {}", *counter.lock().unwrap());  // 10
}
```

---

## 13.7 Cell<T> — Simple Interior Mutability

For `Copy` types, `Cell<T>` is simpler than `RefCell`:

```rust
use std::cell::Cell;

struct Counter {
    value: Cell<u32>,
    name: String,
}

impl Counter {
    fn new(name: &str) -> Self {
        Counter { value: Cell::new(0), name: name.to_string() }
    }

    fn increment(&self) {  // takes &self (not &mut self)!
        self.value.set(self.value.get() + 1);
    }

    fn count(&self) -> u32 {
        self.value.get()
    }
}

fn main() {
    let c = Counter::new("visits");
    c.increment();
    c.increment();
    c.increment();
    println!("{}: {}", c.name, c.count());  // visits: 3
}
```

---

## 13.8 Deref and DerefMut — Smart Pointer Ergonomics

The `Deref` trait allows smart pointers to behave like references. This is why `*box_value` works and why `Box<String>` automatically coerces to `&str`:

```rust
use std::ops::Deref;

struct MyBox<T>(T);

impl<T> MyBox<T> {
    fn new(x: T) -> MyBox<T> {
        MyBox(x)
    }
}

impl<T> Deref for MyBox<T> {
    type Target = T;
    fn deref(&self) -> &T {
        &self.0
    }
}

fn hello(name: &str) {
    println!("Hello, {}!", name);
}

fn main() {
    let s = MyBox::new(String::from("Rust"));

    // Deref coercion chain: MyBox<String> → String → str → &str
    hello(&s);  // automatically: &s → &String (Deref) → &str (Deref)
}
```

---

## 13.9 Unsafe Rust

`unsafe` unlocks capabilities that the borrow checker normally prevents. Required for FFI, raw pointers, and certain low-level operations.

```rust
fn main() {
    // Raw pointers
    let mut num = 5;
    let r1 = &num as *const i32;      // raw immutable pointer
    let r2 = &mut num as *mut i32;    // raw mutable pointer

    unsafe {
        println!("r1: {}", *r1);  // deref raw pointer — unsafe
        *r2 = 10;
        println!("num: {}", num);  // 10
    }

    // Calling unsafe functions
    unsafe fn dangerous() {
        println!("I'm dangerous!");
    }

    unsafe { dangerous(); }

    // FFI — calling C code
    extern "C" {
        fn abs(x: i32) -> i32;
    }

    unsafe {
        println!("{}", abs(-5));  // 5
    }
}
```

### Safe Abstractions Over Unsafe Code

The pattern: write `unsafe` internally, expose a safe API:

```rust
fn split_at_mut(slice: &mut [i32], mid: usize) -> (&mut [i32], &mut [i32]) {
    let len = slice.len();
    assert!(mid <= len);

    let ptr = slice.as_mut_ptr();

    unsafe {
        (
            std::slice::from_raw_parts_mut(ptr, mid),
            std::slice::from_raw_parts_mut(ptr.add(mid), len - mid),
        )
        // Safe: we checked mid <= len; the two slices don't overlap
    }
}

fn main() {
    let mut v = vec![1, 2, 3, 4, 5];
    let (left, right) = split_at_mut(&mut v, 3);
    left[0] = 10;
    right[0] = 40;
    println!("{:?}", v);  // [10, 2, 3, 40, 5]
}
```

---

## Summary

Closures capture their environment in three ways — by reference (Fn), by mutable reference (FnMut), or by moving (FnOnce). Use `move` when the closure must own its captures (threads, returned closures). `Box<T>` heap-allocates and gives unique ownership — use for recursive types and trait objects. `Rc<T>` allows multiple owners with reference counting (single-threaded). `RefCell<T>` moves borrow checking to runtime, enabling interior mutability. `Arc<Mutex<T>>` is the thread-safe combination for shared mutable state. `unsafe` enables raw pointers, FFI, and certain low-level patterns — wrap unsafe code in safe APIs.

---

## Key Takeaways

- Closure traits: `Fn` (multiple immutable calls), `FnMut` (multiple mutable calls), `FnOnce` (one call)
- `move` forces capture by value — required for threads and returned closures
- Use `Box<T>` for recursive types and dynamic dispatch (`Box<dyn Trait>`)
- `Rc<T>` = multiple owners, single thread. `Arc<T>` = multiple owners, multi-thread
- `RefCell<T>` = runtime borrow checking, single thread. `Mutex<T>` = runtime mutual exclusion, multi-thread
- `Arc<Mutex<T>>` is the standard pattern for shared mutable state across threads
- All `unsafe` code should be wrapped in a safe abstraction — document invariants

---

## Exercises

**Exercise 1:** Write a function `memoize<A: Eq + Hash + Copy, B: Copy>(f: impl Fn(A) -> B) -> impl FnMut(A) -> B` that caches results.

**Exercise 2:** Implement a `Tree<T>` using `Rc<RefCell<Node<T>>>` where each node can have multiple children and a parent reference.

**Exercise 3:** Create a `Logger` struct that uses `Cell<u32>` to count calls while keeping `log(&self, msg: &str)` signature (not `&mut self`).

**Exercise 4:** Write a safe wrapper around `split_at_mut` that panics with a custom message if `mid > len`.

**Exercise 5:** Implement a thread-safe counter using `Arc<Mutex<i32>>` that can be incremented from 5 threads simultaneously, printing the final count.

---

*Next: [Chapter 14 — Concurrency](14-concurrency.md)*
