# Chapter 9 — Structs, Enums, and Pattern Matching

---

## 9.1 Structs

Structs are named groups of related data — Rust's equivalent of classes (without inheritance).

### Defining and Instantiating Structs

```rust
struct User {
    username: String,
    email: String,
    sign_in_count: u64,
    active: bool,
}

fn main() {
    // Create an instance — all fields must be initialized
    let user1 = User {
        email: String::from("alice@example.com"),
        username: String::from("alice"),
        active: true,
        sign_in_count: 1,
    };

    // Access fields with dot notation
    println!("{}", user1.username);

    // To mutate, the entire instance must be mut
    let mut user2 = User {
        email: String::from("bob@example.com"),
        username: String::from("bob"),
        active: true,
        sign_in_count: 0,
    };
    user2.email = String::from("newemail@example.com");

    // Struct update syntax — copy remaining fields from another instance
    let user3 = User {
        email: String::from("charlie@example.com"),
        ..user2  // remaining fields from user2
        // Note: user2.username is MOVED into user3 here
    };
    println!("{}", user3.username);
    // println!("{}", user2.username);  // ERROR — moved
}
```

### Field Init Shorthand

```rust
fn build_user(email: String, username: String) -> User {
    User {
        email,         // shorthand when field name == variable name
        username,
        active: true,
        sign_in_count: 0,
    }
}
```

### Tuple Structs

Structs without named fields — useful for wrapping types:

```rust
struct Color(i32, i32, i32);    // RGB
struct Point(f64, f64, f64);    // XYZ

fn main() {
    let black = Color(0, 0, 0);
    let origin = Point(0.0, 0.0, 0.0);

    // Access with .0, .1, .2
    println!("{} {} {}", black.0, black.1, black.2);

    // Destructuring
    let Color(r, g, b) = black;
    println!("r={}, g={}, b={}", r, g, b);

    // Note: Color and Point are different types even though they have the same structure
    // fn takes_color(c: Color) is different from fn takes_point(p: Point)
}
```

### Unit Structs

Structs with no fields — used as marker types:

```rust
struct AlwaysEqual;

fn main() {
    let _s = AlwaysEqual;
}
```

---

## 9.2 Methods

Methods are functions defined within the context of a struct (or enum or trait):

```rust
#[derive(Debug)]
struct Rectangle {
    width: f64,
    height: f64,
}

impl Rectangle {
    // Constructor (by convention called new)
    fn new(width: f64, height: f64) -> Rectangle {
        Rectangle { width, height }
    }

    // Instance method — &self borrows the instance immutably
    fn area(&self) -> f64 {
        self.width * self.height
    }

    fn perimeter(&self) -> f64 {
        2.0 * (self.width + self.height)
    }

    fn is_square(&self) -> bool {
        self.width == self.height
    }

    // Mutable method — &mut self borrows mutably
    fn scale(&mut self, factor: f64) {
        self.width *= factor;
        self.height *= factor;
    }

    // Takes ownership of self (rare)
    fn into_string(self) -> String {
        format!("{}x{}", self.width, self.height)
    }

    // Associated function (no self) — like a static method
    fn square(size: f64) -> Rectangle {
        Rectangle { width: size, height: size }
    }
}

fn main() {
    let mut rect = Rectangle::new(10.0, 5.0);

    println!("Area: {}", rect.area());
    println!("Perimeter: {}", rect.perimeter());
    println!("Is square: {}", rect.is_square());

    rect.scale(2.0);
    println!("After scaling: {:?}", rect);

    // Associated function called with ::
    let sq = Rectangle::square(5.0);
    println!("Square: {:?}", sq);

    // Chaining
    println!("{}", rect.into_string());
    // println!("{:?}", rect);  // ERROR — rect was moved into into_string
}
```

### Multiple impl Blocks

```rust
impl Rectangle {
    fn area(&self) -> f64 { self.width * self.height }
}

impl Rectangle {
    fn perimeter(&self) -> f64 { 2.0 * (self.width + self.height) }
}
// Multiple impl blocks are allowed — useful for organizing code
```

---

## 9.3 Enums

Enums define a type by enumerating its possible variants. Rust enums are **algebraic data types** — variants can hold data.

### Basic Enums

```rust
enum Direction {
    North,
    South,
    East,
    West,
}

fn move_player(dir: Direction) {
    match dir {
        Direction::North => println!("Moving north"),
        Direction::South => println!("Moving south"),
        Direction::East => println!("Moving east"),
        Direction::West => println!("Moving west"),
    }
}

fn main() {
    move_player(Direction::North);
    move_player(Direction::East);
}
```

### Enums with Data

This is where Rust enums shine — each variant can hold different types of data:

```rust
enum Message {
    Quit,                          // no data
    Move { x: i32, y: i32 },      // named fields (struct variant)
    Write(String),                 // single string
    ChangeColor(u8, u8, u8),       // three u8s (tuple variant)
}

impl Message {
    fn call(&self) {
        match self {
            Message::Quit => println!("Quit"),
            Message::Move { x, y } => println!("Move to ({}, {})", x, y),
            Message::Write(s) => println!("Write: {}", s),
            Message::ChangeColor(r, g, b) => println!("Color: ({}, {}, {})", r, g, b),
        }
    }
}

fn main() {
    let messages = vec![
        Message::Quit,
        Message::Move { x: 10, y: 20 },
        Message::Write(String::from("hello")),
        Message::ChangeColor(255, 0, 128),
    ];

    for msg in &messages {
        msg.call();
    }
}
```

### Real-World Enum: Shape

```rust
use std::f64::consts::PI;

#[derive(Debug)]
enum Shape {
    Circle { radius: f64 },
    Rectangle { width: f64, height: f64 },
    Triangle { base: f64, height: f64 },
    RegularPolygon { sides: u32, side_length: f64 },
}

impl Shape {
    fn area(&self) -> f64 {
        match self {
            Shape::Circle { radius } => PI * radius * radius,
            Shape::Rectangle { width, height } => width * height,
            Shape::Triangle { base, height } => 0.5 * base * height,
            Shape::RegularPolygon { sides, side_length } => {
                let n = *sides as f64;
                (n * side_length * side_length) / (4.0 * (PI / n).tan())
            }
        }
    }

    fn name(&self) -> &str {
        match self {
            Shape::Circle { .. } => "Circle",
            Shape::Rectangle { .. } => "Rectangle",
            Shape::Triangle { .. } => "Triangle",
            Shape::RegularPolygon { .. } => "Polygon",
        }
    }
}

fn main() {
    let shapes = vec![
        Shape::Circle { radius: 5.0 },
        Shape::Rectangle { width: 4.0, height: 6.0 },
        Shape::Triangle { base: 3.0, height: 4.0 },
    ];

    for shape in &shapes {
        println!("{}: area = {:.2}", shape.name(), shape.area());
    }
}
```

---

## 9.4 Option<T> — Null Safety

Rust has no `null`. Instead, it has `Option<T>`:

```rust
enum Option<T> {
    Some(T),  // contains a value
    None,     // no value
}
```

This forces you to explicitly handle the "no value" case at compile time. No null pointer exceptions ever.

```rust
fn divide(a: f64, b: f64) -> Option<f64> {
    if b == 0.0 { None } else { Some(a / b) }
}

fn main() {
    // Must handle both cases
    match divide(10.0, 3.0) {
        Some(result) => println!("Result: {:.4}", result),
        None => println!("Division by zero"),
    }

    // Common Option methods
    let some_val: Option<i32> = Some(42);
    let no_val: Option<i32> = None;

    // unwrap() — panics if None (use sparingly)
    println!("{}", some_val.unwrap());

    // unwrap_or() — provide a default
    println!("{}", no_val.unwrap_or(0));

    // unwrap_or_else() — compute default lazily
    println!("{}", no_val.unwrap_or_else(|| 2 + 2));

    // map() — transform the inner value if Some
    let doubled = some_val.map(|x| x * 2);  // Some(84)
    println!("{:?}", doubled);

    // and_then() — chain operations that might fail (flatMap)
    let result = some_val
        .and_then(|x| if x > 0 { Some(x) } else { None })
        .and_then(|x| Some(x.to_string()));
    println!("{:?}", result);

    // is_some() and is_none()
    println!("{} {}", some_val.is_some(), no_val.is_none());

    // ? operator — propagate None early (in functions returning Option)
    fn get_first_doubled(v: &[i32]) -> Option<i32> {
        let first = v.first()?;  // return None if empty
        Some(first * 2)
    }

    println!("{:?}", get_first_doubled(&[1, 2, 3]));  // Some(2)
    println!("{:?}", get_first_doubled(&[]));          // None
}
```

### if let with Option

```rust
fn main() {
    let config_value: Option<i32> = Some(42);

    // Concise pattern for "do something only if Some"
    if let Some(v) = config_value {
        println!("Config: {}", v);
    }

    // while let — loop while Some
    let mut stack = vec![1, 2, 3];
    while let Some(top) = stack.pop() {
        println!("{}", top);
    }
}
```

---

## 9.5 Result<T, E> — Error Handling Preview

`Result` is Option's sibling for operations that can fail with an error:

```rust
enum Result<T, E> {
    Ok(T),   // success with value
    Err(E),  // failure with error
}
```

```rust
fn parse_positive(s: &str) -> Result<u32, String> {
    let n: i64 = s.parse().map_err(|_| format!("'{}' is not a number", s))?;
    if n < 0 {
        Err(format!("{} is negative", n))
    } else {
        Ok(n as u32)
    }
}

fn main() {
    println!("{:?}", parse_positive("42"));    // Ok(42)
    println!("{:?}", parse_positive("-5"));    // Err("−5 is negative")
    println!("{:?}", parse_positive("abc"));  // Err("'abc' is not a number")
}
```

Full error handling coverage is in Chapter 10.

---

## 9.6 Patterns in Enums — Advanced Match

```rust
#[derive(Debug)]
enum Coin {
    Penny,
    Nickel,
    Dime,
    Quarter(String),  // Quarter has a state name
}

fn value_in_cents(coin: &Coin) -> u32 {
    match coin {
        Coin::Penny => {
            println!("Lucky penny!");
            1
        }
        Coin::Nickel => 5,
        Coin::Dime => 10,
        Coin::Quarter(state) => {
            println!("State quarter from {}!", state);
            25
        }
    }
}

fn main() {
    let coins = vec![
        Coin::Penny,
        Coin::Quarter(String::from("Alaska")),
        Coin::Dime,
    ];

    let total: u32 = coins.iter().map(value_in_cents).sum();
    println!("Total: {} cents", total);
}
```

### Destructuring Nested Enums

```rust
enum Color {
    Rgb(u8, u8, u8),
    Hsl(f64, f64, f64),
    Named(String),
}

enum Theme {
    Light { background: Color, foreground: Color },
    Dark { background: Color, foreground: Color },
    Custom(String),
}

fn describe_theme(theme: &Theme) {
    match theme {
        Theme::Light { background: Color::Named(name), .. } => {
            println!("Light theme with {} background", name);
        }
        Theme::Dark { background: Color::Rgb(r, g, b), .. } => {
            println!("Dark theme, bg=rgb({},{},{})", r, g, b);
        }
        Theme::Custom(name) => println!("Custom: {}", name),
        _ => println!("Other theme"),
    }
}

fn main() {
    describe_theme(&Theme::Light {
        background: Color::Named(String::from("white")),
        foreground: Color::Named(String::from("black")),
    });
}
```

---

## Complete Example: A Simple Interpreter

```rust
#[derive(Debug)]
enum Token {
    Number(f64),
    Plus,
    Minus,
    Star,
    Slash,
}

#[derive(Debug)]
enum Expr {
    Literal(f64),
    BinaryOp {
        op: char,
        left: Box<Expr>,
        right: Box<Expr>,
    },
}

impl Expr {
    fn eval(&self) -> f64 {
        match self {
            Expr::Literal(n) => *n,
            Expr::BinaryOp { op, left, right } => {
                let l = left.eval();
                let r = right.eval();
                match op {
                    '+' => l + r,
                    '-' => l - r,
                    '*' => l * r,
                    '/' => l / r,
                    _ => panic!("Unknown op"),
                }
            }
        }
    }
}

fn main() {
    // Represents: (2 + 3) * 4
    let expr = Expr::BinaryOp {
        op: '*',
        left: Box::new(Expr::BinaryOp {
            op: '+',
            left: Box::new(Expr::Literal(2.0)),
            right: Box::new(Expr::Literal(3.0)),
        }),
        right: Box::new(Expr::Literal(4.0)),
    };

    println!("Result: {}", expr.eval());  // 20.0
}
```

---

## Summary

Structs group related data with named fields. Methods are defined in `impl` blocks with `&self`, `&mut self`, or by taking `self` (consuming). Associated functions (no self) serve as constructors and utilities. Enums are algebraic data types — each variant can hold different data, and `match` handles all variants exhaustively. `Option<T>` replaces null: `Some(T)` for a value, `None` for absence. `Result<T, E>` represents operations that can succeed or fail. Both `Option` and `Result` are just enums with a rich set of methods (`map`, `and_then`, `unwrap_or`, etc.).

---

## Key Takeaways

- Structs organize data; `impl` blocks add behavior; both are separated intentionally
- The whole instance must be `mut` to mutate any field — not individual fields
- Struct update syntax `..other` moves non-Copy fields from `other`
- Rust enums are algebraic — variants carry data, enabling expressive type modeling
- `Option<T>` = nullability without null — the compiler forces handling
- `match` on enums is exhaustive — the compiler rejects unhandled variants
- `Box<T>` is needed for recursive types (like tree nodes) to give them a known size

---

## Exercises

**Exercise 1:** Create a `Point3D` struct with `x`, `y`, `z: f64`. Implement `distance_from_origin(&self) -> f64` and `translate(&mut self, dx, dy, dz: f64)`.

**Exercise 2:** Create an enum `JsonValue` with variants: `Null`, `Bool(bool)`, `Number(f64)`, `Str(String)`, `Array(Vec<JsonValue>)`, `Object(HashMap<String, JsonValue>)`. Implement `is_truthy(&self) -> bool`.

**Exercise 3:** Implement a `Stack<T>` struct wrapping a `Vec<T>` with `push`, `pop`, `peek`, `is_empty`, and `size` methods.

**Exercise 4:** Create a `Direction` enum (North/South/East/West) with a method `opposite(&self) -> Direction` and `to_vector(&self) -> (i32, i32)`.

**Exercise 5:** Write a function `flatten(nested: Vec<Option<i32>>) -> Vec<i32>` that extracts all `Some` values, discarding `None`s. Use `filter_map`.

---

*Next: [Chapter 10 — Error Handling](10-error-handling.md)*
