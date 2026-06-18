# Chapter 7 — Borrowing and References

> *"The borrow checker is not your enemy. It's teaching you how to write correct concurrent code."*

Borrowing solves the problem from the last chapter: how do you use a value without taking ownership of it? The answer is **references**.

---

## 7.1 References — Borrowing Without Owning

A reference lets you **refer to a value without owning it**. When you create a reference, you're "borrowing" the value — like borrowing a book from a library. You can read it (or write in it, with permission), but you don't own it.

```rust
fn calculate_length(s: &String) -> usize {  // s is a reference to String
    s.len()
}  // s goes out of scope, but it does NOT drop the String — it doesn't own it

fn main() {
    let s1 = String::from("hello");

    let len = calculate_length(&s1);  // &s1 creates a reference to s1

    // s1 is still valid — we only lent it to the function
    println!("The length of '{}' is {}.", s1, len);
}
```

The `&` operator creates a reference. The `&s1` syntax "borrows" s1 — the function gets to use the String without taking ownership.

### What a Reference Looks Like in Memory

```
s1 (owns the String):          String on heap:
┌──────────┐                  ┌───────────────┐
│ ptr ─────┼─────────────────►│ h e l l o    │
│ len: 5   │                  └───────────────┘
│ cap: 5   │
└──────────┘
     ▲
     │
ref (just a pointer to s1):
┌──────────┐
│ ptr ─────┘
└──────────┘
```

The reference is a pointer to the owner (or directly to the data). It doesn't own anything.

---

## 7.2 The Borrowing Rules

Rust enforces these rules at compile time, with no exceptions:

```
Rule 1: At any given time, you can have EITHER:
        — any number of immutable references (&T), OR
        — exactly ONE mutable reference (&mut T)
        
Rule 2: References must always be valid (no dangling references)
```

These rules prevent **data races** at compile time — the same guarantee that mutexes provide at runtime, but with zero overhead.

---

## 7.3 Immutable References (&T)

You can have as many immutable (shared) references as you want, simultaneously:

```rust
fn main() {
    let s = String::from("hello");

    let r1 = &s;  // immutable reference
    let r2 = &s;  // another immutable reference — OK!
    let r3 = &s;  // yet another — still OK!

    println!("{} {} {}", r1, r2, r3);
    // All valid — reading the same data simultaneously is always safe
}
```

Multiple readers, no writers — this is the classic reader/writer pattern. Reading is always safe because no one is mutating the data.

---

## 7.4 Mutable References (&mut T)

```rust
fn change(s: &mut String) {
    s.push_str(" world");
}

fn main() {
    let mut s = String::from("hello");  // must be mut to create &mut
    change(&mut s);                     // create mutable reference
    println!("{}", s);                  // "hello world"
}
```

**Only one mutable reference at a time:**

```rust
fn main() {
    let mut s = String::from("hello");

    let r1 = &mut s;
    // let r2 = &mut s;  // COMPILE ERROR: cannot borrow `s` as mutable more than once

    println!("{}", r1);
}
```

**Cannot mix mutable and immutable:**

```rust
fn main() {
    let mut s = String::from("hello");

    let r1 = &s;      // immutable — OK
    let r2 = &s;      // another immutable — OK
    // let r3 = &mut s;  // COMPILE ERROR: cannot borrow as mutable because it is
                         // also borrowed as immutable

    println!("{} {}", r1, r2);
}
```

### Why These Rules Exist

```rust
// Without these rules, this would be a data race:
let mut s = String::from("hello");
let r1 = &s;           // reads s
let r2 = &mut s;       // modifies s — but r1 is still "looking at" s!
println!("{}", r1);    // What does r1 see? Undefined behavior.
```

In C++, this is undefined behavior — your program might crash, return garbage, or appear to work but have subtle corruption. Rust prevents it at compile time.

---

## 7.5 Non-Lexical Lifetimes (NLL)

The borrow checker is smart about when a reference's lifetime actually ends. It ends at the **last use**, not the end of the scope:

```rust
fn main() {
    let mut s = String::from("hello");

    let r1 = &s;      // immutable borrow begins
    let r2 = &s;      // another immutable borrow
    println!("{} {}", r1, r2);
    // r1 and r2 are no longer used after this point — their borrows END here

    // Now a mutable reference is OK:
    let r3 = &mut s;  // mutable borrow begins
    r3.push_str(" world");
    println!("{}", r3);
}
```

This is called **Non-Lexical Lifetimes** (NLL) — the borrow checker uses data flow analysis to determine where borrows actually end. Earlier Rust versions were more conservative and would reject this code.

---

## 7.6 Dangling References

Rust prevents dangling references — pointers to memory that has been freed:

```rust
// COMPILE ERROR: this would be a dangling reference
fn dangle() -> &String {  // trying to return a reference to a local String
    let s = String::from("hello");
    &s  // return a reference to s
}  // s goes out of scope here — its memory is freed
   // but we're returning a reference to it!
   // In C: this returns a pointer to freed memory — crash or security bug
   // In Rust: COMPILE ERROR

fn main() {
    let ref_to_nothing = dangle();  // this line causes the error above
}
```

The solution is to return the owned String directly (transfer ownership):

```rust
fn no_dangle() -> String {
    let s = String::from("hello");
    s  // move s out — no dangling reference
}

fn main() {
    let s = no_dangle();  // s owns the String
    println!("{}", s);
}
```

---

## 7.7 Slices — References to Parts of Collections

Slices are references to a contiguous sequence of elements in a collection. They don't own their data.

### String Slices (&str)

```rust
fn main() {
    let s = String::from("hello world");

    // Create slices with range indices
    let hello = &s[0..5];    // "hello"
    let world = &s[6..11];   // "world"

    println!("{} {}", hello, world);

    // Shorthand ranges
    let hello2 = &s[..5];    // same as [0..5]
    let world2 = &s[6..];    // same as [6..11]
    let all = &s[..];        // the whole string

    // String literals ARE string slices
    let literal: &str = "hello world";
    // "hello world" is stored in the program binary
    // literal is a &str pointing into that binary data

    println!("{:?}", std::mem::size_of::<&str>());  // 16 bytes: (ptr, len) fat pointer
}
```

### Using &str Instead of &String

```rust
// WORSE: only accepts String
fn first_word_v1(s: &String) -> &str {
    let bytes = s.as_bytes();
    for (i, &byte) in bytes.iter().enumerate() {
        if byte == b' ' {
            return &s[..i];
        }
    }
    &s[..]
}

// BETTER: accepts both String and string literals
fn first_word(s: &str) -> &str {
    let bytes = s.as_bytes();
    for (i, &byte) in bytes.iter().enumerate() {
        if byte == b' ' {
            return &s[..i];
        }
    }
    &s[..]
}

fn main() {
    let s = String::from("hello world");
    let word = first_word(&s);          // works with String
    let word2 = first_word("hello world"); // works with &str literal

    println!("First word: {}", word);
    println!("First word: {}", word2);

    // The lifetime connection: word is a slice of s
    // If s changes, word becomes invalid:
    // s.clear();          // would drop the string data
    // println!("{}", word); // COMPILE ERROR — s was mutated while word borrows it
}
```

### Array Slices

```rust
fn sum(slice: &[i32]) -> i32 {
    slice.iter().sum()
}

fn main() {
    let arr = [1, 2, 3, 4, 5];
    let vec = vec![10, 20, 30];

    // Both arrays and vecs can be sliced
    println!("{}", sum(&arr));           // 15
    println!("{}", sum(&arr[1..4]));     // 9 (2+3+4)
    println!("{}", sum(&vec));           // 60

    // Slice type: &[i32]
    let slice: &[i32] = &arr[1..3];  // [2, 3]
    println!("{:?}", slice);
    println!("len={}", slice.len());

    // Fat pointer: (pointer to data, length)
    println!("{}", std::mem::size_of::<&[i32]>());  // 16 bytes
}
```

---

## 7.8 The Relationship Between References and Lifetimes

Every reference has a **lifetime** — the scope for which that reference is valid. The compiler tracks lifetimes to ensure references never outlive the data they point to.

Most of the time, the compiler infers lifetimes automatically (lifetime elision). You only need explicit lifetime annotations when the compiler can't figure it out — covered in depth in Chapter 8.

```rust
// The compiler automatically knows:
// - r1 borrows s and is valid as long as s is valid
// - word borrows s and is valid until s is modified or dropped

fn main() {
    let s = String::from("hello world");
    let word = first_word(&s);  // word borrows from s

    // As long as we use word, s cannot be mutably borrowed
    println!("{}", word);  // last use of word

    // Now s can be modified freely
    // s.clear();  // would be OK after the last use of word
}

fn first_word(s: &str) -> &str {
    match s.find(' ') {
        Some(i) => &s[..i],
        None => s,
    }
}
```

---

## Complete Example: Text Analysis

```rust
fn count_words(text: &str) -> usize {
    text.split_whitespace().count()
}

fn longest_word<'a>(text: &'a str) -> &'a str {
    text.split_whitespace()
        .max_by_key(|w| w.len())
        .unwrap_or("")
}

fn contains_word(text: &str, word: &str) -> bool {
    text.split_whitespace().any(|w| w == word)
}

fn word_frequencies(text: &str) -> std::collections::HashMap<&str, usize> {
    let mut map = std::collections::HashMap::new();
    for word in text.split_whitespace() {
        *map.entry(word).or_insert(0) += 1;
    }
    map
}

fn main() {
    let text = "the quick brown fox jumps over the lazy dog the";

    println!("Words: {}", count_words(text));               // 9
    println!("Longest: {}", longest_word(text));            // "jumps" or "quick"
    println!("Has 'fox': {}", contains_word(text, "fox"));  // true

    let frequencies = word_frequencies(text);
    println!("'the' appears {} times", frequencies.get("the").unwrap_or(&0)); // 3

    // All these functions borrow text — no cloning needed
    // text is still valid throughout
    println!("Original text still valid: {}", &text[..3]);  // "the"
}
```

---

## Summary

References let you borrow values without taking ownership. Shared references (`&T`) allow multiple concurrent readers. Mutable references (`&mut T`) give exclusive write access — no other references can exist while it's active. These rules prevent data races at compile time. Slices (`&[T]` and `&str`) are fat pointers (pointer + length) to contiguous data. The compiler automatically tracks reference lifetimes and rejects code where references outlive their source data.

---

## Key Takeaways

- `&T` is a shared (immutable) reference — many can coexist
- `&mut T` is an exclusive (mutable) reference — only one can exist at a time
- Cannot have both `&T` and `&mut T` to the same data simultaneously
- No dangling references — the compiler rejects them
- Borrowing ends at last use (NLL), not end of lexical scope
- `&str` is a string slice — prefer it over `&String` in function parameters
- Slices are fat pointers: 16 bytes (pointer + length)

---

## Exercises

**Exercise 1:** Write a function `largest<'a>(list: &'a [i32]) -> &'a i32` that returns a reference to the largest element. Why does it need lifetime annotations?

**Exercise 2:** What does this code do? Why is it a compile error?
```rust
let mut v = vec![1, 2, 3];
let first = &v[0];
v.push(4);
println!("{}", first);
```

**Exercise 3:** Write a function `split_at_first(s: &str, ch: char) -> (&str, &str)` that splits a string slice at the first occurrence of a character, returning two slices.

**Exercise 4:** Explain why `&String` is less flexible than `&str` for function parameters, and how Rust handles the conversion automatically.

**Exercise 5:** Write a function that takes `&[f64]` and returns `(f64, f64)` — the minimum and maximum values — without cloning.

---

*Next: [Chapter 8 — Lifetimes](08-lifetimes.md)*
