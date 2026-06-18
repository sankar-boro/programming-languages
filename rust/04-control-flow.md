# Chapter 4 — Control Flow

---

## 4.1 if Expressions

In Rust, `if` is an **expression** — it returns a value. There is no ternary operator; `if` replaces it.

```rust
fn main() {
    let number = 7;

    // Basic if-else
    if number < 5 {
        println!("small");
    } else if number < 10 {
        println!("medium");
    } else {
        println!("large");
    }

    // if as an expression — returns a value
    let description = if number % 2 == 0 { "even" } else { "odd" };
    println!("{} is {}", number, description);

    // Both arms must return the same type
    let value = if true { 42 } else { 0 };  // i32
    // let bad = if true { 42 } else { "hello" };  // COMPILE ERROR: type mismatch

    // In let bindings
    let abs = if number < 0 { -number } else { number };
}
```

### Conditions Must Be bool

Unlike C or JavaScript, Rust does **not** coerce non-boolean types to bool:

```rust
fn main() {
    let x = 5;
    // if x { }         // COMPILE ERROR: expected bool, found integer
    // if x != 0 { }   // CORRECT

    let s = "hello";
    // if s { }         // COMPILE ERROR
    // if !s.is_empty() { }   // CORRECT
}
```

---

## 4.2 match — Deep Dive

`match` is Rust's most powerful control flow construct. It is an expression that pattern-matches a value against a series of patterns, executing the code for the first match.

### Basic match

```rust
fn main() {
    let number = 3;

    match number {
        1 => println!("one"),
        2 => println!("two"),
        3 => println!("three"),
        4 | 5 => println!("four or five"),       // multiple patterns with |
        6..=9 => println!("six through nine"),    // inclusive range
        _ => println!("something else"),          // wildcard — catches everything
    }
}
```

### match Is Exhaustive

The compiler requires every possible value to be handled:

```rust
fn main() {
    let x: bool = true;

    match x {
        true => println!("yes"),
        // COMPILE ERROR if false arm is missing:
        // error[E0004]: non-exhaustive patterns: `false` not covered
        false => println!("no"),
    }
}
```

### match as an Expression

```rust
fn main() {
    let dice_roll = 6;

    let message = match dice_roll {
        1 => "Critical failure!",
        2..=5 => "Normal roll",
        6 => "Critical success!",
        _ => unreachable!("Dice only goes 1-6"),
    };

    println!("{}", message);
}
```

### Matching Tuples

```rust
fn main() {
    let point = (0, -2);

    match point {
        (0, 0) => println!("Origin"),
        (x, 0) | (0, x) => println!("On an axis at {}", x),
        (x, y) if x == y => println!("On the diagonal at {}", x),
        (x, y) => println!("At ({}, {})", x, y),
    }
}
```

### match Guards — Adding if Conditions

```rust
fn main() {
    let num = Some(4);

    match num {
        Some(x) if x < 0 => println!("Got negative: {}", x),
        Some(x) if x == 0 => println!("Got zero"),
        Some(x) => println!("Got positive: {}", x),
        None => println!("Got nothing"),
    }
}
```

### Binding with @ (at operator)

```rust
fn main() {
    let num = 15;

    match num {
        // n @ pattern — binds the value to `n` while matching the pattern
        n @ 1..=12 => println!("Got {} (small)", n),
        n @ 13..=19 => println!("Got {} (teen)", n),
        n => println!("Got {} (large)", n),
    }
}
```

### Destructuring in match

```rust
struct Point { x: i32, y: i32 }
enum Shape { Circle { radius: f64 }, Rectangle { width: f64, height: f64 } }

fn main() {
    // Struct destructuring
    let p = Point { x: 0, y: 7 };
    match p {
        Point { x: 0, y } => println!("On Y axis at {}", y),
        Point { x, y: 0 } => println!("On X axis at {}", x),
        Point { x, y } => println!("At ({}, {})", x, y),
    }

    // Enum destructuring
    let shape = Shape::Circle { radius: 5.0 };
    let area = match shape {
        Shape::Circle { radius } => std::f64::consts::PI * radius * radius,
        Shape::Rectangle { width, height } => width * height,
    };
    println!("Area: {:.2}", area);
}
```

### Ignoring Values with ..

```rust
struct Point3D { x: i32, y: i32, z: i32 }

fn main() {
    let p = Point3D { x: 0, y: 5, z: 10 };

    match p {
        Point3D { x: 0, .. } => println!("X is zero, don't care about y and z"),
        Point3D { x, .. } => println!("X is {}", x),
    }

    let numbers = (2, 4, 8, 16, 32);
    match numbers {
        (first, .., last) => println!("first={}, last={}", first, last),
    }
}
```

---

## 4.3 loop, while, for

### loop — Infinite Loop with break Value

```rust
fn main() {
    // Basic infinite loop
    let mut count = 0;
    loop {
        count += 1;
        if count == 5 {
            break;
        }
    }
    println!("count = {}", count);

    // loop returns a value via break
    let result = loop {
        count += 1;
        if count == 10 {
            break count * 2;  // returns 20
        }
    };
    println!("result = {}", result);  // 20

    // Labeled loops — break from an outer loop
    'outer: loop {
        let mut x = 0;
        loop {
            x += 1;
            if x == 3 {
                break 'outer;  // breaks the outer loop
            }
        }
        println!("This never prints");
    }
    println!("Exited outer loop");
}
```

### while — Condition-Based Loop

```rust
fn main() {
    let mut n = 1;

    while n < 100 {
        n *= 2;
    }
    println!("First power of 2 >= 100: {}", n);  // 128

    // while with complex condition
    let mut x = 0;
    let mut y = 10;
    while x < 5 && y > 0 {
        x += 1;
        y -= 2;
    }
    println!("x={}, y={}", x, y);

    // while let — loop while a pattern matches
    let mut stack = vec![1, 2, 3];
    while let Some(top) = stack.pop() {
        println!("popped: {}", top);
    }
    // 3, 2, 1
}
```

### for — Iterating Over Collections

```rust
fn main() {
    // Over a range
    for i in 0..5 {
        print!("{} ", i);  // 0 1 2 3 4
    }
    println!();

    // Inclusive range
    for i in 0..=5 {
        print!("{} ", i);  // 0 1 2 3 4 5
    }
    println!();

    // Over an array
    let arr = [10, 20, 30, 40, 50];
    for element in &arr {
        print!("{} ", element);  // 10 20 30 40 50
    }
    println!();

    // With index using enumerate()
    for (i, value) in arr.iter().enumerate() {
        println!("arr[{}] = {}", i, value);
    }

    // Over a Vec
    let words = vec!["hello", "world", "rust"];
    for word in &words {
        println!("{}", word);
    }

    // Consuming iteration (no &)
    let numbers = vec![1, 2, 3, 4, 5];
    for n in numbers {      // numbers is moved here — cannot use it after
        println!("{}", n);
    }
    // println!("{:?}", numbers);  // COMPILE ERROR: moved

    // Mutable iteration
    let mut data = vec![1, 2, 3, 4, 5];
    for x in &mut data {
        *x *= 2;  // dereference to modify
    }
    println!("{:?}", data);  // [2, 4, 6, 8, 10]
}
```

### continue and break

```rust
fn main() {
    // continue — skip to next iteration
    for i in 0..10 {
        if i % 2 == 0 {
            continue;
        }
        print!("{} ", i);  // 1 3 5 7 9
    }
    println!();

    // break — exit the loop
    for i in 0..100 {
        if i * i > 50 {
            println!("First number whose square > 50: {}", i);
            break;
        }
    }
}
```

---

## 4.4 Pattern Matching — The Full Picture

Pattern matching is used in `match`, `if let`, `while let`, `for`, `let`, and function parameters.

### Patterns in let

```rust
fn main() {
    // Tuple destructuring in let
    let (x, y, z) = (1, 2, 3);
    println!("{} {} {}", x, y, z);

    // Struct destructuring in let
    struct Point { x: i32, y: i32 }
    let Point { x, y } = Point { x: 5, y: 10 };
    println!("{} {}", x, y);

    // Ignored values with _
    let (a, _, c) = (1, 2, 3);
    println!("{} {}", a, c);

    // Ignored remainder with ..
    let (first, ..) = (1, 2, 3, 4, 5);
    println!("{}", first);
}
```

### if let — Single Pattern Match

`if let` is syntactic sugar for a `match` with one arm:

```rust
fn main() {
    let some_value: Option<i32> = Some(42);

    // Long form with match:
    match some_value {
        Some(n) => println!("Got: {}", n),
        None => {}
    }

    // Short form with if let:
    if let Some(n) = some_value {
        println!("Got: {}", n);
    }

    // With else:
    if let Some(n) = some_value {
        println!("Some: {}", n);
    } else {
        println!("None");
    }

    // Chaining if let
    let favorite_color: Option<&str> = None;
    let is_tuesday = false;
    let age: Result<u8, _> = "34".parse();

    if let Some(color) = favorite_color {
        println!("Using your color: {}", color);
    } else if is_tuesday {
        println!("Tuesday is green day!");
    } else if let Ok(age) = age {
        if age > 30 {
            println!("Using purple as the background color");
        } else {
            println!("Using orange as the background color");
        }
    } else {
        println!("Using blue as the background color");
    }
}
```

### Patterns in Function Parameters

```rust
fn print_point(&(x, y): &(i32, i32)) {
    println!("({}, {})", x, y);
}

fn process_pair((first, second): (i32, i32)) -> i32 {
    first + second
}

fn main() {
    let point = (3, 5);
    print_point(&point);
    println!("{}", process_pair((10, 20)));
}
```

---

## 4.5 Ranges

```rust
fn main() {
    // Exclusive range (start..end)
    let r1 = 1..5;   // 1, 2, 3, 4
    for i in r1 { print!("{} ", i); }
    println!();

    // Inclusive range (start..=end)
    let r2 = 1..=5;  // 1, 2, 3, 4, 5
    for i in r2 { print!("{} ", i); }
    println!();

    // Ranges in match
    let x = 42;
    match x {
        0 => println!("zero"),
        1..=9 => println!("single digit"),
        10..=99 => println!("two digits"),
        100..=999 => println!("three digits"),
        _ => println!("very large"),
    }

    // Ranges have methods
    let range = 1..=100;
    println!("{}", range.contains(&50));  // true

    // Step through a range (not a range method, use step_by on iter)
    for i in (0..20).step_by(3) {
        print!("{} ", i);  // 0 3 6 9 12 15 18
    }
    println!();

    // Reverse
    for i in (0..5).rev() {
        print!("{} ", i);  // 4 3 2 1 0
    }
}
```

---

## Summary

Rust's `if` is an expression that returns a value. `match` is exhaustive, pattern-based, and also an expression — it is far more powerful than switch statements in other languages. The `loop` construct creates infinite loops with `break` values; `while` uses conditions; `for` iterates over any iterable. Pattern matching appears everywhere: `match`, `if let`, `while let`, `let`, and function parameters. Patterns can destructure, bind, guard, and ignore values.

---

## Key Takeaways

- `if` is an expression — both arms must return the same type
- `match` is exhaustive — the compiler rejects non-exhaustive patterns
- `match` arms can: match ranges, multiple values (`|`), guard with `if`, bind with `@`
- `loop` returns a value via `break value`
- `for` works with any type that implements `IntoIterator`
- `if let` is concise pattern matching for single-pattern cases
- Patterns can appear in `let`, function args, `for` — not just `match`

---

## Exercises

**Exercise 1:** Write a function `grade(score: u32) -> &'static str` using `match` that returns "A" for 90+, "B" for 80+, etc.

**Exercise 2:** Use a labeled `'outer` loop to find the first pair (i, j) where `i * j > 100`, with i in 1..=10 and j in 1..=10.

**Exercise 3:** Write a `fizzbuzz(n: u32)` function using `match` with a tuple `(n % 3 == 0, n % 5 == 0)`.

**Exercise 4:** Use `while let` to drain a `Vec<Option<i32>>`, printing each `Some` value and stopping at the first `None`.

**Exercise 5:** Write a function that takes a `(i32, i32)` point and returns which quadrant it's in (I, II, III, IV) or "on axis" using `match`.

---

*Next: [Chapter 5 — Functions](05-functions.md)*
