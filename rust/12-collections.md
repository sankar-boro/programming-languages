# Chapter 12 — Collections, Strings, and Iterators

---

## 12.1 Vec<T> — Dynamic Arrays

`Vec<T>` is the most commonly used collection in Rust. It's a heap-allocated, dynamically-sized list.

### Creating Vecs

```rust
fn main() {
    // Empty vec — type must be specified
    let v1: Vec<i32> = Vec::new();
    let v2: Vec<i32> = Vec::with_capacity(10);  // pre-allocate capacity

    // With elements — type inferred
    let v3 = vec![1, 2, 3, 4, 5];
    let v4 = vec![0; 10];  // ten zeros

    // From an iterator
    let v5: Vec<i32> = (1..=5).collect();
    println!("{:?}", v5);  // [1, 2, 3, 4, 5]
}
```

### Modifying Vecs

```rust
fn main() {
    let mut v = Vec::new();

    v.push(1);
    v.push(2);
    v.push(3);

    // Access
    let third = &v[2];  // panics if out of bounds
    println!("{}", third);

    let third_safe = v.get(2);  // returns Option<&i32>
    println!("{:?}", third_safe);

    // Pop (removes and returns last element)
    while let Some(top) = v.pop() {
        print!("{} ", top);
    }
    println!();

    // Extend
    let mut a = vec![1, 2, 3];
    let b = vec![4, 5, 6];
    a.extend(b.iter());  // borrow b
    a.extend([7, 8, 9]); // extend with array
    println!("{:?}", a);

    // Insert and remove at position
    let mut v = vec![1, 2, 4, 5];
    v.insert(2, 3);   // insert 3 at index 2
    v.remove(4);      // remove element at index 4
    println!("{:?}", v);  // [1, 2, 3, 4]

    // Retain — keep only elements satisfying predicate
    let mut v = vec![1, 2, 3, 4, 5, 6];
    v.retain(|&x| x % 2 == 0);
    println!("{:?}", v);  // [2, 4, 6]

    // Sort
    let mut v = vec![3, 1, 4, 1, 5, 9, 2, 6];
    v.sort();
    println!("{:?}", v);  // [1, 1, 2, 3, 4, 5, 6, 9]

    v.sort_by(|a, b| b.cmp(a));  // reverse sort
    v.sort_by_key(|k| std::cmp::Reverse(*k));  // reverse using key
    println!("{:?}", v);  // [9, 6, 5, 4, 3, 2, 1, 1]

    // Dedup (removes consecutive duplicates — sort first!)
    let mut v = vec![1, 1, 2, 3, 3, 3, 4];
    v.dedup();
    println!("{:?}", v);  // [1, 2, 3, 4]
}
```

### Slicing and Iterating

```rust
fn main() {
    let v = vec![1, 2, 3, 4, 5];

    // Iterate (borrow — v still valid)
    for x in &v { print!("{} ", x); }
    println!();

    // Iterate mutably
    let mut v = vec![1, 2, 3, 4, 5];
    for x in &mut v { *x *= 2; }
    println!("{:?}", v);  // [2, 4, 6, 8, 10]

    // Consuming iteration (v is moved)
    let v = vec![1, 2, 3];
    for x in v { println!("{}", x); }
    // v is gone

    // Slices
    let v = vec![10, 20, 30, 40, 50];
    let slice = &v[1..4];  // [20, 30, 40]
    println!("{:?}", slice);

    // Common operations
    println!("{}", v.len());
    println!("{}", v.is_empty());
    println!("{:?}", v.first());   // Some(10)
    println!("{:?}", v.last());    // Some(50)
    println!("{:?}", v.contains(&30));  // true
    println!("{:?}", v.iter().position(|&x| x == 30));  // Some(2)
}
```

---

## 12.2 HashMap<K, V>

```rust
use std::collections::HashMap;

fn main() {
    // Create
    let mut scores: HashMap<String, i32> = HashMap::new();

    // Insert
    scores.insert(String::from("Alice"), 100);
    scores.insert(String::from("Bob"), 90);
    scores.insert(String::from("Charlie"), 85);

    // Access
    let alice = scores.get("Alice");  // Option<&i32>
    println!("{:?}", alice);  // Some(100)

    // Access with default
    let dave_score = scores.get("Dave").copied().unwrap_or(0);
    println!("{}", dave_score);  // 0

    // Check containment
    println!("{}", scores.contains_key("Bob"));

    // Entry API — insert only if not present
    scores.entry(String::from("Dave")).or_insert(75);
    scores.entry(String::from("Alice")).or_insert(0);  // doesn't change Alice
    println!("{:?}", scores.get("Dave"));    // Some(75)
    println!("{:?}", scores.get("Alice"));   // Some(100)

    // Entry — update existing value
    let alice_score = scores.entry(String::from("Alice")).or_insert(0);
    *alice_score += 10;  // Alice now has 110

    // Iterate
    for (name, score) in &scores {
        println!("{}: {}", name, score);
    }

    // Remove
    scores.remove("Charlie");

    // Collect from pairs
    let pairs = vec![("x", 1), ("y", 2), ("z", 3)];
    let map: HashMap<&str, i32> = pairs.into_iter().collect();
    println!("{:?}", map);

    // Word frequency counter
    let text = "hello world hello rust world hello";
    let mut freq: HashMap<&str, usize> = HashMap::new();
    for word in text.split_whitespace() {
        *freq.entry(word).or_insert(0) += 1;
    }
    println!("{:?}", freq);
}
```

---

## 12.3 HashSet<T>

```rust
use std::collections::HashSet;

fn main() {
    let mut set: HashSet<i32> = HashSet::new();

    set.insert(1);
    set.insert(2);
    set.insert(3);
    set.insert(2);  // duplicate — ignored
    println!("{:?}", set);  // {1, 2, 3} (order not guaranteed)

    // Set operations
    let set_a: HashSet<i32> = [1, 2, 3, 4].iter().cloned().collect();
    let set_b: HashSet<i32> = [3, 4, 5, 6].iter().cloned().collect();

    // Union
    let union: HashSet<&i32> = set_a.union(&set_b).collect();
    println!("Union: {:?}", union);

    // Intersection
    let intersection: HashSet<&i32> = set_a.intersection(&set_b).collect();
    println!("Intersection: {:?}", intersection);

    // Difference
    let difference: HashSet<&i32> = set_a.difference(&set_b).collect();
    println!("a - b: {:?}", difference);

    // Subset/superset
    let small: HashSet<i32> = [1, 2].iter().cloned().collect();
    println!("{}", small.is_subset(&set_a));  // true
}
```

---

## 12.4 Strings — Two Types

Rust has two string types:
- `String` — owned, heap-allocated, mutable
- `&str` — borrowed string slice, immutable view into string data

```rust
fn main() {
    // &str — string literal (stored in binary)
    let s1: &str = "hello";

    // String — heap allocated, owned
    let s2: String = String::from("hello");
    let s3: String = "hello".to_string();
    let s4: String = "hello".to_owned();

    // Conversion
    let slice: &str = &s2;    // String → &str (deref coercion)
    let owned: String = s1.to_string();  // &str → String

    // Concatenation
    let s1 = String::from("Hello, ");
    let s2 = String::from("world!");
    let s3 = s1 + &s2;  // s1 is MOVED here (it's the operator's `self`)
    // println!("{}", s1);  // ERROR — moved

    // Use format! to avoid moves
    let s1 = String::from("Hello");
    let s2 = String::from("world");
    let s3 = format!("{}, {}!", s1, s2);  // both s1 and s2 still valid
    println!("{}", s3);

    // String capacity and length
    let s = String::with_capacity(50);
    println!("len={}, cap={}", s.len(), s.capacity());
}
```

### String Methods

```rust
fn main() {
    let s = String::from("  Hello, World!  ");

    // Trim
    println!("{}", s.trim());         // "Hello, World!"
    println!("{}", s.trim_start());   // "Hello, World!  "
    println!("{}", s.trim_end());     // "  Hello, World!"

    // Case
    println!("{}", s.to_uppercase());
    println!("{}", s.to_lowercase());

    // Contains, starts_with, ends_with
    println!("{}", s.contains("World"));
    println!("{}", s.starts_with("  Hello"));
    println!("{}", s.ends_with("  "));

    // Replace
    let replaced = s.replace("World", "Rust");
    println!("{}", replaced);

    // Split
    let csv = "a,b,c,d,e";
    let parts: Vec<&str> = csv.split(',').collect();
    println!("{:?}", parts);

    // Split and count
    let word_count = "hello world foo bar".split_whitespace().count();
    println!("{}", word_count);  // 4

    // Find
    println!("{:?}", s.find("World"));  // Some(9)

    // Chars — iterate over Unicode characters
    let s = "hello 🦀";
    for c in s.chars() {
        print!("{} ", c);
    }
    println!();
    println!("char count: {}", s.chars().count());  // 7
    println!("byte count: {}", s.len());             // 10 (crab emoji = 4 bytes)

    // Bytes
    for b in "hello".bytes() {
        print!("{} ", b);  // 104 101 108 108 111
    }
    println!();

    // Parse
    let n: i32 = "42".parse().unwrap();
    println!("{}", n);
}
```

### String Slicing and Unicode

```rust
fn main() {
    let s = String::from("hello world");

    // String slicing uses byte indices, not char indices
    let hello = &s[0..5];   // "hello"
    let world = &s[6..11];  // "world"
    println!("{} {}", hello, world);

    // DANGER: slicing in the middle of a multi-byte character panics!
    let emoji_str = "🦀 rust";
    // &emoji_str[0..1]  // PANIC — 🦀 is 4 bytes; byte 1 is mid-char

    let crab = &emoji_str[0..4];  // OK — full crab emoji
    println!("{}", crab);

    // Safe: use char_indices for Unicode-safe slicing
    let s = "αβγδ";
    for (i, c) in s.char_indices() {
        println!("byte {}: '{}'", i, c);
    }
}
```

---

## 12.5 Iterators

Rust's iterator system is one of its most powerful features. Iterators are lazy — they produce values on demand, not all at once.

### The Iterator Trait

```rust
pub trait Iterator {
    type Item;
    fn next(&mut self) -> Option<Self::Item>;
    // + 70 default methods built on next()
}
```

### Creating Iterators

```rust
fn main() {
    let v = vec![1, 2, 3, 4, 5];

    // iter() — yields &T (borrows)
    let iter = v.iter();

    // iter_mut() — yields &mut T (mutable borrows)
    let mut v_mut = vec![1, 2, 3];
    for x in v_mut.iter_mut() { *x *= 2; }

    // into_iter() — yields T (moves)
    let v2 = vec![1, 2, 3];
    for x in v2.into_iter() { println!("{}", x); }
    // v2 is moved — cannot use after

    // Ranges are iterators
    for i in 0..10 { print!("{} ", i); }
    println!();

    // Manual iteration
    let mut iter = vec![1, 2, 3].into_iter();
    println!("{:?}", iter.next());  // Some(1)
    println!("{:?}", iter.next());  // Some(2)
    println!("{:?}", iter.next());  // Some(3)
    println!("{:?}", iter.next());  // None
}
```

### Iterator Adapters (Lazy)

These transform iterators without computing anything yet:

```rust
fn main() {
    let v = vec![1, 2, 3, 4, 5, 6, 7, 8, 9, 10];

    // map — transform each element
    let doubled: Vec<i32> = v.iter().map(|&x| x * 2).collect();
    println!("{:?}", doubled);

    // filter — keep elements matching predicate
    let evens: Vec<&i32> = v.iter().filter(|&&x| x % 2 == 0).collect();
    println!("{:?}", evens);

    // filter_map — map + filter in one step
    let doubled_evens: Vec<i32> = v.iter()
        .filter_map(|&x| if x % 2 == 0 { Some(x * 2) } else { None })
        .collect();
    println!("{:?}", doubled_evens);

    // take — first N elements
    let first_three: Vec<&i32> = v.iter().take(3).collect();
    println!("{:?}", first_three);  // [1, 2, 3]

    // skip — skip N elements
    let after_three: Vec<&i32> = v.iter().skip(3).collect();
    println!("{:?}", after_three);  // [4, 5, 6, 7, 8, 9, 10]

    // chain — concatenate iterators
    let a = vec![1, 2, 3];
    let b = vec![4, 5, 6];
    let combined: Vec<&i32> = a.iter().chain(b.iter()).collect();
    println!("{:?}", combined);

    // zip — pair elements from two iterators
    let names = vec!["Alice", "Bob", "Charlie"];
    let scores = vec![100, 90, 85];
    let pairs: Vec<(&&str, &i32)> = names.iter().zip(scores.iter()).collect();
    for (name, score) in &pairs {
        println!("{}: {}", name, score);
    }

    // enumerate — add index
    for (i, val) in v.iter().enumerate() {
        print!("[{}]={} ", i, val);
    }
    println!();

    // flat_map — map + flatten
    let words = vec!["hello world", "foo bar"];
    let chars: Vec<&str> = words.iter().flat_map(|s| s.split_whitespace()).collect();
    println!("{:?}", chars);

    // peekable — peek at next without consuming
    let mut iter = v.iter().peekable();
    println!("{:?}", iter.peek());  // Some(1) — not consumed
    println!("{:?}", iter.next());  // Some(1) — consumed
}
```

### Consumer Methods (Eager)

These consume the iterator and produce a result:

```rust
fn main() {
    let v = vec![1, 2, 3, 4, 5];

    // collect — gather into a collection
    let doubled: Vec<i32> = v.iter().map(|&x| x * 2).collect();

    // sum, product
    let sum: i32 = v.iter().sum();
    let product: i32 = v.iter().product();
    println!("{} {}", sum, product);  // 15 120

    // count
    let count = v.iter().filter(|&&x| x > 3).count();
    println!("{}", count);  // 2

    // any, all
    println!("{}", v.iter().any(|&x| x > 4));   // true
    println!("{}", v.iter().all(|&x| x > 0));   // true

    // find, position
    println!("{:?}", v.iter().find(|&&x| x > 3));       // Some(4)
    println!("{:?}", v.iter().position(|&x| x > 3));     // Some(3)

    // min, max
    println!("{:?}", v.iter().min());  // Some(1)
    println!("{:?}", v.iter().max());  // Some(5)

    // min_by_key, max_by_key
    let words = vec!["alpha", "beta", "gamma", "delta"];
    println!("{:?}", words.iter().max_by_key(|s| s.len()));  // Some("gamma" or "delta")

    // reduce — combine elements
    let sum = v.iter().copied().reduce(|acc, x| acc + x);
    println!("{:?}", sum);  // Some(15)

    // fold — reduce with initial value
    let sum = v.iter().fold(0, |acc, &x| acc + x);
    println!("{}", sum);  // 15

    // for_each
    v.iter().for_each(|x| print!("{} ", x));
    println!();

    // Unzip pairs into two collections
    let pairs = vec![(1, 'a'), (2, 'b'), (3, 'c')];
    let (nums, chars): (Vec<i32>, Vec<char>) = pairs.into_iter().unzip();
    println!("{:?} {:?}", nums, chars);
}
```

### Custom Iterators

```rust
struct Fibonacci {
    curr: u64,
    next: u64,
}

impl Fibonacci {
    fn new() -> Self {
        Fibonacci { curr: 0, next: 1 }
    }
}

impl Iterator for Fibonacci {
    type Item = u64;

    fn next(&mut self) -> Option<u64> {
        let result = self.curr;
        let new_next = self.curr + self.next;
        self.curr = self.next;
        self.next = new_next;
        Some(result)  // infinite iterator — always returns Some
    }
}

fn main() {
    // Take first 10 Fibonacci numbers
    let fibs: Vec<u64> = Fibonacci::new().take(10).collect();
    println!("{:?}", fibs);

    // Sum of Fibonacci numbers under 1000
    let sum: u64 = Fibonacci::new()
        .take_while(|&x| x < 1000)
        .sum();
    println!("Sum: {}", sum);

    // First Fibonacci number over 1000
    let first_over_1000 = Fibonacci::new().find(|&x| x > 1000);
    println!("{:?}", first_over_1000);  // Some(1597)
}
```

---

## Other Collections

```rust
use std::collections::{BTreeMap, BTreeSet, VecDeque, LinkedList, BinaryHeap};

fn main() {
    // BTreeMap — sorted HashMap (useful when order matters)
    let mut btree: BTreeMap<String, i32> = BTreeMap::new();
    btree.insert("b".to_string(), 2);
    btree.insert("a".to_string(), 1);
    btree.insert("c".to_string(), 3);
    for (k, v) in &btree { println!("{}: {}", k, v); }  // sorted by key

    // VecDeque — double-ended queue
    let mut deque: VecDeque<i32> = VecDeque::new();
    deque.push_back(1);
    deque.push_back(2);
    deque.push_front(0);
    println!("{:?}", deque);  // [0, 1, 2]
    deque.pop_front();
    deque.pop_back();

    // BinaryHeap — max-heap priority queue
    let mut heap = BinaryHeap::new();
    heap.push(3);
    heap.push(1);
    heap.push(4);
    heap.push(1);
    heap.push(5);
    while let Some(top) = heap.pop() {
        print!("{} ", top);  // 5 4 3 1 1 (largest first)
    }
    println!();
}
```

---

## Summary

`Vec<T>` is the go-to dynamic collection — use it for ordered, indexed data. `HashMap<K, V>` maps keys to values with O(1) average access. `HashSet<T>` is a set with O(1) membership testing. Rust has two string types: owned `String` and borrowed `&str` — prefer `&str` for function parameters. Iterators are lazy and composable — chain adapters (`map`, `filter`, `take`) and consume with collectors (`collect`, `sum`, `fold`). Implementing the `Iterator` trait requires only `next()` — all other methods come for free.

---

## Key Takeaways

- Prefer `Vec` for ordered data, `HashMap` for key-value, `HashSet` for uniqueness
- `iter()` borrows, `iter_mut()` mutably borrows, `into_iter()` consumes
- Iterator adapters are lazy — nothing runs until a consumer is called
- `collect()` is powerful: it can produce `Vec`, `HashMap`, `HashSet`, `String`, etc.
- `entry().or_insert()` is the idiomatic way to insert-or-update in a HashMap
- `&str` for function parameters, `String` for owned/stored data
- String indexing uses bytes — use `.chars()` for Unicode-correct character operations

---

## Exercises

**Exercise 1:** Write a function `group_by<T, K, F>(items: Vec<T>, key_fn: F) -> HashMap<K, Vec<T>> where K: Eq + Hash, F: Fn(&T) -> K` that groups items by a key function.

**Exercise 2:** Implement a word frequency counter that reads a string and returns the top-N most frequent words as a `Vec<(String, usize)>` sorted by frequency descending.

**Exercise 3:** Write a function `flatten<T: Clone>(nested: &[Vec<T>]) -> Vec<T>` using `flat_map`.

**Exercise 4:** Implement `fn running_average(numbers: &[f64]) -> Vec<f64>` that returns a vec where each element is the average of all preceding elements (inclusive).

**Exercise 5:** Create a `Primes` iterator struct that yields prime numbers indefinitely. Use it with `.take(20).collect::<Vec<_>>()`.

---

*Next: [Chapter 13 — Advanced Features](13-advanced-features.md)*
