# Chapter 14 — Concurrency

> *"Fearless concurrency — Rust's ownership system catches data races at compile time, the same type of bug that brings down production servers in C++."*

---

## 14.1 Threads

Rust threads map 1:1 to OS threads:

```rust
use std::thread;
use std::time::Duration;

fn main() {
    // Spawn a thread
    let handle = thread::spawn(|| {
        for i in 1..=10 {
            println!("spawned: {}", i);
            thread::sleep(Duration::from_millis(1));
        }
    });

    for i in 1..=5 {
        println!("main: {}", i);
        thread::sleep(Duration::from_millis(1));
    }

    handle.join().unwrap();  // wait for thread to finish
}
```

### Moving Data Into Threads

Threads must own their data — you can't borrow across thread boundaries:

```rust
use std::thread;

fn main() {
    let v = vec![1, 2, 3];

    // ERROR: cannot borrow `v` — v might be dropped before thread finishes
    // let handle = thread::spawn(|| println!("{:?}", v));

    // CORRECT: move v into the thread
    let handle = thread::spawn(move || {
        println!("{:?}", v);
    });

    handle.join().unwrap();
}
```

### Thread Builder — Named Threads, Stack Size

```rust
use std::thread;

fn main() {
    let builder = thread::Builder::new()
        .name("worker-1".to_string())
        .stack_size(4 * 1024 * 1024);  // 4MB stack

    let handle = builder.spawn(|| {
        println!("Thread name: {:?}", thread::current().name());
    }).unwrap();

    handle.join().unwrap();
}
```

---

## 14.2 Channels — Message Passing

Channels allow threads to communicate by sending messages. Rust uses **multi-producer, single-consumer (mpsc)** channels:

```rust
use std::sync::mpsc;
use std::thread;

fn main() {
    let (tx, rx) = mpsc::channel();

    // Producer thread
    thread::spawn(move || {
        let values = vec!["hello", "from", "the", "thread"];
        for val in values {
            tx.send(val).unwrap();
            thread::sleep(std::time::Duration::from_millis(100));
        }
        // tx is dropped here — receiver knows channel is closed
    });

    // Consumer (main thread)
    for received in rx {  // iterates until sender drops
        println!("Got: {}", received);
    }
}
```

### Multiple Producers

```rust
use std::sync::mpsc;
use std::thread;

fn main() {
    let (tx, rx) = mpsc::channel();

    // Clone the sender for multiple producers
    let tx1 = tx.clone();
    let tx2 = tx.clone();

    thread::spawn(move || {
        for i in 0..5 {
            tx1.send(format!("from thread 1: {}", i)).unwrap();
        }
    });

    thread::spawn(move || {
        for i in 0..5 {
            tx2.send(format!("from thread 2: {}", i)).unwrap();
        }
    });

    // tx is dropped here — but tx1 and tx2 are still alive
    drop(tx);  // explicit drop to close that end

    for msg in rx {
        println!("{}", msg);
    }
}
```

### Sync Channels — Bounded Channels

```rust
use std::sync::mpsc;
use std::thread;

fn main() {
    // Bounded: sender blocks when buffer is full (backpressure)
    let (tx, rx) = mpsc::sync_channel(2);  // buffer size 2

    let handle = thread::spawn(move || {
        println!("sending 1...");
        tx.send(1).unwrap();  // ok — buffer empty
        println!("sending 2...");
        tx.send(2).unwrap();  // ok — buffer has 1 item
        println!("sending 3...");
        tx.send(3).unwrap();  // BLOCKS until receiver reads something
        println!("all sent");
    });

    thread::sleep(std::time::Duration::from_millis(100));
    println!("receiving...");
    for val in rx {
        println!("got: {}", val);
        thread::sleep(std::time::Duration::from_millis(100));
    }

    handle.join().unwrap();
}
```

---

## 14.3 Mutex<T> — Mutual Exclusion

When threads need to share and mutate data, use `Mutex<T>`:

```rust
use std::sync::{Arc, Mutex};
use std::thread;

fn main() {
    // Mutex protects the data — only one thread can access at a time
    let counter = Arc::new(Mutex::new(0));
    let mut handles = vec![];

    for id in 0..10 {
        let counter = Arc::clone(&counter);
        let handle = thread::spawn(move || {
            let mut num = counter.lock().unwrap();  // acquire lock
            *num += 1;
            println!("Thread {} incremented to {}", id, *num);
        });  // MutexGuard is dropped here — lock released automatically
        handles.push(handle);
    }

    for handle in handles {
        handle.join().unwrap();
    }

    println!("Final: {}", *counter.lock().unwrap());  // 10
}
```

### Deadlock — What to Avoid

```rust
use std::sync::{Arc, Mutex};
use std::thread;

fn main() {
    let lock_a = Arc::new(Mutex::new(1));
    let lock_b = Arc::new(Mutex::new(2));

    let la = Arc::clone(&lock_a);
    let lb = Arc::clone(&lock_b);

    // Thread 1: acquires A then tries B
    let h1 = thread::spawn(move || {
        let _a = la.lock().unwrap();
        thread::sleep(std::time::Duration::from_millis(10));
        let _b = lb.lock().unwrap();  // may deadlock if thread 2 has B
    });

    // Thread 2: acquires B then tries A
    let h2 = thread::spawn(move || {
        let _b = lock_b.lock().unwrap();
        thread::sleep(std::time::Duration::from_millis(10));
        let _a = lock_a.lock().unwrap();  // may deadlock if thread 1 has A
    });

    // Prevention: always acquire locks in the same order
    // Or use a single lock protecting both values
}
```

---

## 14.4 RwLock — Multiple Readers, One Writer

```rust
use std::sync::{Arc, RwLock};
use std::thread;

fn main() {
    let data = Arc::new(RwLock::new(vec![1, 2, 3]));
    let mut handles = vec![];

    // Multiple simultaneous readers
    for i in 0..3 {
        let data = Arc::clone(&data);
        handles.push(thread::spawn(move || {
            let r = data.read().unwrap();  // shared read lock
            println!("Reader {}: {:?}", i, *r);
        }));
    }

    // One exclusive writer
    let data = Arc::clone(&data);
    handles.push(thread::spawn(move || {
        let mut w = data.write().unwrap();  // exclusive write lock
        w.push(4);
        println!("Writer added 4");
    }));

    for h in handles {
        h.join().unwrap();
    }
}
```

---

## 14.5 Send and Sync — The Thread Safety Markers

Two marker traits control concurrency safety:

**`Send`** — safe to transfer ownership across thread boundaries  
**`Sync`** — safe to share a reference across thread boundaries (`&T` is `Send` if `T: Sync`)

```rust
// Types that implement Send:
// - Most primitives: i32, String, Vec<T> (if T: Send)
// - Arc<T> (if T: Send + Sync)
// - Mutex<T> (if T: Send)

// Types that do NOT implement Send:
// - Rc<T> — reference counting is not atomic
// - *mut T — raw pointer — could be used unsafely from multiple threads
// - RefCell<T> — not thread-safe borrow checking

// This is a compile error:
// fn send_rc() {
//     let rc = std::rc::Rc::new(5);
//     thread::spawn(move || println!("{}", rc));  // ERROR: Rc is not Send
// }

// Use Arc instead
fn main() {
    let arc = std::sync::Arc::new(5);
    thread::spawn(move || println!("{}", arc)).join().unwrap();
}
```

---

## 14.6 Thread Pools — Practical Concurrency

For many short tasks, spawning a thread per task is expensive. Use a thread pool:

```rust
use std::sync::{Arc, Mutex};
use std::thread;

struct ThreadPool {
    workers: Vec<thread::JoinHandle<()>>,
    sender: std::sync::mpsc::Sender<Box<dyn FnOnce() + Send + 'static>>,
}

impl ThreadPool {
    fn new(size: usize) -> Self {
        let (tx, rx) = std::sync::mpsc::channel::<Box<dyn FnOnce() + Send>>();
        let rx = Arc::new(Mutex::new(rx));

        let mut workers = vec![];
        for _ in 0..size {
            let rx = Arc::clone(&rx);
            let handle = thread::spawn(move || {
                loop {
                    let job = rx.lock().unwrap().recv();
                    match job {
                        Ok(job) => job(),
                        Err(_) => break,  // sender dropped — shut down
                    }
                }
            });
            workers.push(handle);
        }

        ThreadPool { workers, sender: tx }
    }

    fn execute<F: FnOnce() + Send + 'static>(&self, job: F) {
        self.sender.send(Box::new(job)).unwrap();
    }
}

fn main() {
    let pool = ThreadPool::new(4);

    for i in 0..8 {
        pool.execute(move || {
            println!("Task {} running on {:?}", i, thread::current().id());
            thread::sleep(std::time::Duration::from_millis(100));
        });
    }

    // Wait for completion
    thread::sleep(std::time::Duration::from_millis(500));
}
```

---

## 14.7 Concurrent Data Structures

### Atomic Types — Lock-Free Primitives

```rust
use std::sync::atomic::{AtomicUsize, Ordering};
use std::sync::Arc;
use std::thread;

fn main() {
    let counter = Arc::new(AtomicUsize::new(0));
    let mut handles = vec![];

    for _ in 0..10 {
        let counter = Arc::clone(&counter);
        handles.push(thread::spawn(move || {
            for _ in 0..1000 {
                counter.fetch_add(1, Ordering::SeqCst);
            }
        }));
    }

    for h in handles { h.join().unwrap(); }
    println!("Count: {}", counter.load(Ordering::SeqCst));  // 10000
}
```

### Memory Ordering

```rust
use std::sync::atomic::Ordering;

// SeqCst (Sequentially Consistent) — strongest, most expensive
// Acquire — load; pairs with Release
// Release — store; pairs with Acquire
// AcqRel — both Acquire and Release (for read-modify-write)
// Relaxed — weakest, fastest — no ordering guarantees

// For simple counters: Relaxed is sufficient
// For producer/consumer flags: Release (store) + Acquire (load)
// When in doubt: SeqCst (correct but slower)
```

---

## Complete Example: Parallel Map

```rust
use std::sync::{Arc, Mutex};
use std::thread;

fn parallel_map<T, R, F>(items: Vec<T>, workers: usize, f: F) -> Vec<R>
where
    T: Send + 'static,
    R: Send + 'static + Default + Clone,
    F: Fn(T) -> R + Send + Sync + 'static,
{
    let items = Arc::new(Mutex::new(items.into_iter().enumerate()));
    let results = Arc::new(Mutex::new(vec![]));
    let f = Arc::new(f);
    let mut handles = vec![];

    for _ in 0..workers {
        let items = Arc::clone(&items);
        let results = Arc::clone(&results);
        let f = Arc::clone(&f);

        handles.push(thread::spawn(move || {
            loop {
                let item = items.lock().unwrap().next();
                match item {
                    Some((i, item)) => {
                        let result = f(item);
                        results.lock().unwrap().push((i, result));
                    }
                    None => break,
                }
            }
        }));
    }

    for h in handles { h.join().unwrap(); }

    let mut results = Arc::try_unwrap(results).unwrap().into_inner().unwrap();
    results.sort_by_key(|(i, _)| *i);
    results.into_iter().map(|(_, r)| r).collect()
}

fn main() {
    let numbers = vec![1, 2, 3, 4, 5, 6, 7, 8, 9, 10];
    let results = parallel_map(numbers, 4, |n| {
        thread::sleep(std::time::Duration::from_millis(10));
        n * n
    });
    println!("{:?}", results);  // [1, 4, 9, 16, 25, 36, 49, 64, 81, 100]
}
```

---

## Summary

Rust threads are OS threads — lightweight to use, ownership-safe by design. Move closures transfer data into threads. Channels (`mpsc`) are for message passing between threads. `Mutex<T>` with `Arc<T>` provides shared mutable state. `RwLock<T>` allows multiple readers or one writer. `Send` and `Sync` marker traits are automatically implemented and enforced by the compiler — using `Rc<T>` across threads is a compile error. Atomic types enable lock-free operations for simple counters and flags.

---

## Key Takeaways

- `thread::spawn` + `move` = transfer ownership to thread
- `Arc::clone` = cheap shared ownership across threads (atomic reference counting)
- `Arc<Mutex<T>>` = the standard pattern for shared mutable state
- `mpsc::channel()` = multiple producers, one consumer; `sync_channel(n)` = bounded
- `Send` = ownership can cross thread boundaries; `Sync` = reference can cross boundaries
- `Rc<T>` is NOT `Send` — use `Arc<T>`. `RefCell<T>` is NOT `Sync` — use `Mutex<T>`
- Atomics (`AtomicUsize`, etc.) = lock-free, fine-grained shared state

---

## Exercises

**Exercise 1:** Write a program that spawns 5 threads. Each thread generates 1000 random numbers and sends their sum to the main thread via a channel. The main thread prints the grand total.

**Exercise 2:** Implement a `Cache<K, V>` that uses `RwLock<HashMap<K, V>>` for concurrent reads and exclusive writes.

**Exercise 3:** Build a simple work queue: producer thread generates 100 items, 4 consumer threads process them. Use `Arc<Mutex<VecDeque<Item>>>`.

**Exercise 4:** Implement a barrier synchronization: N threads each do work, then wait until all have finished before proceeding. (Hint: look at `std::sync::Barrier`.)

**Exercise 5:** Create a `SharedCounter` with methods `increment()` and `get()` that use `AtomicI32` instead of `Mutex<i32>`. Compare the performance.

---

*Next: [Chapter 15 — Modules and Crates](15-modules-crates.md)*
