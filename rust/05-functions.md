# Chapter 5 — Functions in Rust

---

## 5.1 Function Syntax

```rust
// Basic function
fn function_name(param1: Type1, param2: Type2) -> ReturnType {
    // body
}

fn add(x: i32, y: i32) -> i32 {
    x + y  // no semicolon — this is a return expression
}

fn greet(name: &str) {
    println!("Hello, {}!", name);
    // returns () implicitly
}

fn main() {
    println!("{}", add(3, 4));  // 7
    greet("Rust");              // Hello, Rust!
}
```

### Parameter Types Are Always Required

Unlike local variables, function parameters cannot use type inference — types must always be annotated:

```rust
// COMPILE ERROR: missing type for function parameter
// fn add(x, y) -> i32 { x + y }

// CORRECT
fn add(x: i32, y: i32) -> i32 { x + y }
```

### Functions Can Be Defined Anywhere

Unlike C, Rust doesn't require forward declarations. Functions can be defined after they're called:

```rust
fn main() {
    println!("{}", helper());  // works fine
}

fn helper() -> i32 {
    42
}
```

---

## 5.2 Expressions vs Statements

This is one of the most important conceptual distinctions in Rust.

**Statement** — performs an action, does NOT return a value:
- `let x = 5;` — variable binding
- `5;` — expression terminated with `;` becomes a statement
- `fn foo() {}` — function declaration

**Expression** — evaluates to a value:
- `5` — literal
- `x + y` — arithmetic
- `if condition { a } else { b }` — if expression
- `{ let x = 3; x + 1 }` — block expression
- `add(5, 6)` — function call

```rust
fn main() {
    // Statement — no value
    let x = 5;

    // Expression — has a value (the block evaluates to 6)
    let y = {
        let x = 3;
        x + 1  // no semicolon — this is the block's value
    };
    println!("y = {}", y);  // 6

    // Adding a semicolon turns an expression into a statement:
    let z = {
        let x = 3;
        x + 1;  // semicolon! This is now a statement — block returns ()
    };
    // z is () — the unit type
    println!("{:?}", z);  // ()
}
```

### The Semicolon Rule

```rust
fn returns_five() -> i32 {
    5  // expression — the function returns 5
}

fn returns_nothing() {
    5;  // statement — the function returns ()
}

// COMPILE ERROR: expected i32, found ()
// fn wrong() -> i32 {
//     5;  // this returns ()
// }
```

---

## 5.3 Return Values

### Implicit Return (Last Expression)

The last expression in a function body is its return value:

```rust
fn double(x: i32) -> i32 {
    x * 2  // no 'return' keyword needed
}

fn max(a: i32, b: i32) -> i32 {
    if a > b { a } else { b }  // if expression as return value
}
```

### Explicit Return

`return` is used for early returns:

```rust
fn find_first_negative(numbers: &[i32]) -> Option<i32> {
    for &n in numbers {
        if n < 0 {
            return Some(n);  // early return
        }
    }
    None  // implicit return at end
}

fn divide(a: f64, b: f64) -> Option<f64> {
    if b == 0.0 {
        return None;  // early return to avoid division by zero
    }
    Some(a / b)  // implicit return
}
```

### Multiple Return Values via Tuples

Rust doesn't have multiple return values directly, but tuples serve the same purpose:

```rust
fn min_max(numbers: &[i32]) -> (i32, i32) {
    let mut min = numbers[0];
    let mut max = numbers[0];

    for &n in numbers {
        if n < min { min = n; }
        if n > max { max = n; }
    }

    (min, max)  // return tuple
}

fn main() {
    let nums = [3, 1, 4, 1, 5, 9, 2, 6];
    let (min, max) = min_max(&nums);
    println!("min={}, max={}", min, max);
}
```

### The Never Type (!)

Some functions never return — they diverge. Their return type is `!` (the "never" type):

```rust
fn crash(msg: &str) -> ! {
    panic!("{}", msg);  // panic! never returns — the thread crashes
}

fn infinite() -> ! {
    loop {
        // never exits
    }
}

fn main() {
    // ! can coerce to any type — useful in match arms:
    let x: i32 = match "5".parse() {
        Ok(n) => n,
        Err(_) => panic!("not a number"),  // panic! returns !, coerces to i32
    };
    println!("{}", x);
}
```

---

## 5.4 Closures — Introduction

Closures are anonymous functions that can capture their environment. Full coverage is in Chapter 13, but here's the foundation:

```rust
fn main() {
    // Basic closure syntax
    let double = |x| x * 2;             // inferred types
    let add = |x: i32, y: i32| x + y;  // explicit types
    let greet = |name: &str| {           // multi-line closure
        println!("Hello, {}!", name);
    };

    println!("{}", double(5));           // 10
    println!("{}", add(3, 4));           // 7
    greet("Rust");                       // Hello, Rust!

    // Closures capture their environment
    let base = 10;
    let add_to_base = |x| x + base;  // captures 'base' from surrounding scope
    println!("{}", add_to_base(5));   // 15

    // Functions vs closures
    fn fn_double(x: i32) -> i32 { x * 2 }  // function: can't capture env
    let cl_double = |x| x * 2;              // closure: can capture env

    // Both can be called the same way
    println!("{}", fn_double(5));
    println!("{}", cl_double(5));

    // Closures as arguments (higher-order functions)
    let numbers = vec![1, 2, 3, 4, 5];
    let doubled: Vec<i32> = numbers.iter().map(|&x| x * 2).collect();
    let evens: Vec<&i32> = numbers.iter().filter(|&&x| x % 2 == 0).collect();

    println!("{:?}", doubled);  // [2, 4, 6, 8, 10]
    println!("{:?}", evens);    // [2, 4]
}
```

### Passing Functions as Arguments

```rust
fn apply(x: i32, f: fn(i32) -> i32) -> i32 {
    f(x)
}

fn square(x: i32) -> i32 { x * x }

fn main() {
    println!("{}", apply(5, square));       // 25
    println!("{}", apply(5, |x| x + 1));   // 6 (but closure needs impl Fn)
}

// For closures, use generic bounds:
fn apply_generic<F: Fn(i32) -> i32>(x: i32, f: F) -> i32 {
    f(x)
}

fn main2() {
    let offset = 10;
    println!("{}", apply_generic(5, |x| x + offset));  // 15
}
```

---

## Function Pointers and the fn Type

```rust
fn add(a: i32, b: i32) -> i32 { a + b }
fn subtract(a: i32, b: i32) -> i32 { a - b }

fn apply_operation(a: i32, b: i32, op: fn(i32, i32) -> i32) -> i32 {
    op(a, b)
}

fn main() {
    let op: fn(i32, i32) -> i32 = add;
    println!("{}", op(3, 4));  // 7

    println!("{}", apply_operation(10, 3, add));       // 13
    println!("{}", apply_operation(10, 3, subtract));  // 7

    // Function pointers in arrays/vectors
    let ops: Vec<fn(i32, i32) -> i32> = vec![add, subtract];
    for op in &ops {
        println!("{}", op(10, 5));  // 15, then 5
    }
}
```

---

## Complete Example: A Mini Calculator

```rust
fn parse_number(s: &str) -> Result<f64, String> {
    s.trim().parse::<f64>().map_err(|e| format!("Parse error: {}", e))
}

fn calculate(a: f64, op: char, b: f64) -> Result<f64, String> {
    match op {
        '+' => Ok(a + b),
        '-' => Ok(a - b),
        '*' => Ok(a * b),
        '/' => {
            if b == 0.0 {
                Err("Division by zero".to_string())
            } else {
                Ok(a / b)
            }
        }
        _ => Err(format!("Unknown operator: {}", op)),
    }
}

fn format_result(result: f64) -> String {
    if result.fract() == 0.0 {
        format!("{}", result as i64)  // show as integer if whole number
    } else {
        format!("{:.4}", result)
    }
}

fn main() {
    let expressions = vec![
        (10.0, '+', 5.0),
        (10.0, '-', 3.0),
        (4.0, '*', 7.0),
        (15.0, '/', 4.0),
        (10.0, '/', 0.0),
        (5.0, '^', 2.0),
    ];

    for (a, op, b) in expressions {
        match calculate(a, op, b) {
            Ok(result) => println!("{} {} {} = {}", a, op, b, format_result(result)),
            Err(e) => println!("{} {} {} = ERROR: {}", a, op, b, e),
        }
    }
}
```

---

## Summary

Rust functions require explicit parameter types. The last expression in a function body is its return value — no `return` keyword needed unless exiting early. Statements do not return values; expressions do. Adding a semicolon to an expression converts it to a statement that returns `()`. Closures are anonymous functions that can capture their environment, making them more powerful than named functions. The `!` (never) type indicates a function that never returns. Functions can be stored in variables and passed as arguments.

---

## Key Takeaways

- Function parameters always require type annotations
- Last expression = return value; semicolon converts expression to statement
- `return` is for early exits only — don't use it at the end of functions
- Closures capture their environment; functions cannot
- Use `fn(T) -> R` for function pointers, `impl Fn(T) -> R` for closures in generics
- `!` is the "never" type — functions that crash or loop forever
- Rust functions are zero-cost — they compile to the same machine code as C functions

---

## Exercises

**Exercise 1:** Write a function `is_prime(n: u64) -> bool` using a `for` loop and early `return false`.

**Exercise 2:** Write a higher-order function `compose(f: impl Fn(i32) -> i32, g: impl Fn(i32) -> i32) -> impl Fn(i32) -> i32` that returns a closure applying f then g.

**Exercise 3:** Write a function `apply_n(f: impl Fn(i32) -> i32, n: u32, x: i32) -> i32` that applies `f` to `x` exactly `n` times.

**Exercise 4:** Explain why this code is a compile error and fix it:
```rust
fn add_five(x: i32) -> i32 {
    x + 5;
}
```

**Exercise 5:** Write a function `count_where<T, F>(slice: &[T], pred: F) -> usize where F: Fn(&T) -> bool` that counts elements satisfying a predicate.

---

*Next: [Chapter 6 — Ownership](06-ownership.md)*
