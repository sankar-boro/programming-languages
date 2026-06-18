# Chapter 11 — Generics and Traits

---

## 11.1 Generics

Generics allow you to write code that works with multiple types. Rust implements generics through **monomorphization** — at compile time, generic code is expanded into concrete versions for each type it's used with. Zero runtime cost.

### Generic Functions

```rust
// Without generics — must write one for each type
fn largest_i32(list: &[i32]) -> &i32 {
    let mut largest = &list[0];
    for item in list {
        if item > largest { largest = item; }
    }
    largest
}

fn largest_char(list: &[char]) -> &char {
    let mut largest = &list[0];
    for item in list {
        if item > largest { largest = item; }
    }
    largest
}

// With generics — one function for any comparable type
fn largest<T: PartialOrd>(list: &[T]) -> &T {
    let mut largest = &list[0];
    for item in list {
        if item > largest { largest = item; }
    }
    largest
}

fn main() {
    let numbers = vec![34, 50, 25, 100, 65];
    println!("Largest: {}", largest(&numbers));

    let chars = vec!['y', 'm', 'a', 'q'];
    println!("Largest: {}", largest(&chars));
}
```

### Generic Structs

```rust
#[derive(Debug)]
struct Point<T> {
    x: T,
    y: T,
}

impl<T> Point<T> {
    fn new(x: T, y: T) -> Self {
        Point { x, y }
    }

    fn x(&self) -> &T {
        &self.x
    }
}

// Methods only available for specific types
impl Point<f64> {
    fn distance_from_origin(&self) -> f64 {
        (self.x.powi(2) + self.y.powi(2)).sqrt()
    }
}

// Struct with multiple type parameters
#[derive(Debug)]
struct Pair<T, U> {
    first: T,
    second: U,
}

fn main() {
    let int_point = Point::new(5, 10);
    let float_point = Point::new(1.0, 4.0);

    println!("{:?}", int_point);
    println!("x = {}", float_point.x());
    println!("distance = {}", float_point.distance_from_origin());

    let mixed = Pair { first: 5, second: "hello" };
    println!("{:?}", mixed);
}
```

### Generic Enums

You've already seen the most important generic enums:

```rust
enum Option<T> {
    Some(T),
    None,
}

enum Result<T, E> {
    Ok(T),
    Err(E),
}
```

You can define your own:

```rust
#[derive(Debug)]
enum Either<L, R> {
    Left(L),
    Right(R),
}

fn main() {
    let x: Either<i32, &str> = Either::Left(42);
    let y: Either<i32, &str> = Either::Right("hello");

    match &x {
        Either::Left(n) => println!("Left: {}", n),
        Either::Right(s) => println!("Right: {}", s),
    }
}
```

---

## 11.2 Traits — Defining Shared Behavior

A trait defines an interface — a set of methods that a type must implement. Like interfaces in Java/Go, but more powerful.

### Defining Traits

```rust
trait Summary {
    fn summarize(&self) -> String;  // method signature — implementors must provide this

    // Default implementation — implementors can override
    fn author(&self) -> String {
        String::from("(anonymous)")
    }

    fn full_summary(&self) -> String {
        format!("{} — {}", self.summarize(), self.author())
    }
}
```

### Implementing Traits

```rust
struct Article {
    title: String,
    content: String,
    author: String,
}

struct Tweet {
    username: String,
    content: String,
}

impl Summary for Article {
    fn summarize(&self) -> String {
        format!("{}: {}...", self.title, &self.content[..50.min(self.content.len())])
    }

    fn author(&self) -> String {
        format!("@{}", self.author)
    }
}

impl Summary for Tweet {
    fn summarize(&self) -> String {
        format!("{}: {}", self.username, self.content)
    }
    // Uses default author() and full_summary()
}

fn main() {
    let article = Article {
        title: String::from("Rust is Amazing"),
        content: String::from("Rust combines safety and performance..."),
        author: String::from("rustacean"),
    };

    let tweet = Tweet {
        username: String::from("rustlang"),
        content: String::from("exciting new Rust release!"),
    };

    println!("{}", article.full_summary());
    println!("{}", tweet.summarize());
    println!("{}", tweet.author());  // uses default
}
```

---

## 11.3 Trait Bounds

Trait bounds specify what a generic type must implement:

### impl Trait Syntax (Simple Cases)

```rust
// "notify accepts any type that implements Summary"
fn notify(item: &impl Summary) {
    println!("Breaking news! {}", item.summarize());
}

// Return impl Trait
fn create_summary() -> impl Summary {
    Tweet {
        username: String::from("bot"),
        content: String::from("Hello!"),
    }
}
```

### where Clauses (Complex Cases)

```rust
// Trait bound syntax — equivalent to impl Trait but more flexible
fn notify_generic<T: Summary>(item: &T) {
    println!("Breaking news! {}", item.summarize());
}

// Multiple bounds
fn notify_multi<T: Summary + std::fmt::Display>(item: &T) {
    println!("{}: {}", item, item.summarize());
}

// where clause — cleaner for complex bounds
fn some_function<T, U>(t: &T, u: &U) -> String
where
    T: std::fmt::Display + Clone,
    U: std::fmt::Debug + Summary,
{
    format!("{:?} {}", u, t)
}
```

### Conditional Method Implementation

```rust
use std::fmt::Display;

struct Wrapper<T>(T);

impl<T: Display + PartialOrd> Wrapper<T> {
    fn display_if_positive(&self) where T: Default + PartialOrd {
        if self.0 > T::default() {
            println!("{}", self.0);
        }
    }
}

// Blanket implementations — implement trait for all types satisfying bounds
impl<T: Display> ToString for Wrapper<T> {
    fn to_string(&self) -> String {
        format!("Wrapper({})", self.0)
    }
}
```

---

## 11.4 Common Standard Library Traits

### Display and Debug

```rust
use std::fmt;

#[derive(Debug)]  // auto-implement Debug
struct Color {
    r: u8,
    g: u8,
    b: u8,
}

// Manual Display implementation
impl fmt::Display for Color {
    fn fmt(&self, f: &mut fmt::Formatter<'_>) -> fmt::Result {
        write!(f, "#{:02X}{:02X}{:02X}", self.r, self.g, self.b)
    }
}

fn main() {
    let c = Color { r: 255, g: 128, b: 0 };
    println!("{}", c);    // Display: #FF8000
    println!("{:?}", c);  // Debug: Color { r: 255, g: 128, b: 0 }
    println!("{:#?}", c); // Pretty Debug
}
```

### Clone and Copy

```rust
#[derive(Debug, Clone)]  // Clone: must opt in
struct Config {
    host: String,
    port: u16,
}

#[derive(Debug, Clone, Copy)]  // Copy: auto-copy on assignment
struct Point {
    x: f64,
    y: f64,
}

fn main() {
    let c1 = Config { host: String::from("localhost"), port: 8080 };
    let c2 = c1.clone();  // explicit deep copy
    println!("{:?}", c1);  // still valid

    let p1 = Point { x: 1.0, y: 2.0 };
    let p2 = p1;           // automatic copy
    println!("{:?}", p1);  // still valid — Copy semantics
}
```

### PartialEq, Eq, PartialOrd, Ord

```rust
#[derive(Debug, PartialEq, Eq, PartialOrd, Ord)]
struct Version {
    major: u32,
    minor: u32,
    patch: u32,
}

fn main() {
    let v1 = Version { major: 1, minor: 0, patch: 0 };
    let v2 = Version { major: 2, minor: 0, patch: 0 };
    let v3 = Version { major: 1, minor: 0, patch: 0 };

    println!("{}", v1 == v3);   // true
    println!("{}", v1 < v2);    // true
    println!("{:?}", v1.cmp(&v2));  // Less

    let mut versions = vec![v2, v3, v1];
    versions.sort();
    println!("{:?}", versions);  // sorted by major, then minor, then patch
}
```

### Default

```rust
#[derive(Debug, Default)]
struct Config {
    host: String,      // Default: empty string
    port: u16,         // Default: 0
    debug: bool,       // Default: false
    timeout: Option<u32>,  // Default: None
}

fn main() {
    let config = Config::default();
    println!("{:?}", config);

    // Struct update with defaults
    let custom = Config {
        host: String::from("localhost"),
        port: 8080,
        ..Default::default()
    };
    println!("{:?}", custom);
}
```

### From and Into

```rust
#[derive(Debug)]
struct Celsius(f64);

#[derive(Debug)]
struct Fahrenheit(f64);

impl From<Celsius> for Fahrenheit {
    fn from(c: Celsius) -> Self {
        Fahrenheit(c.0 * 9.0 / 5.0 + 32.0)
    }
}

fn main() {
    let boiling = Celsius(100.0);
    let f: Fahrenheit = boiling.into();  // Into is automatically implemented if From exists
    println!("{:?}", f);  // Fahrenheit(212.0)

    let freezing = Fahrenheit::from(Celsius(0.0));
    println!("{:?}", freezing);  // Fahrenheit(32.0)
}
```

---

## 11.5 Trait Objects — Dynamic Dispatch

Sometimes you need a collection of different types that all implement the same trait. Use `dyn Trait`:

```rust
trait Draw {
    fn draw(&self);
}

struct Circle { radius: f64 }
struct Rectangle { width: f64, height: f64 }
struct Triangle { base: f64, height: f64 }

impl Draw for Circle {
    fn draw(&self) { println!("Drawing circle with radius {}", self.radius); }
}

impl Draw for Rectangle {
    fn draw(&self) { println!("Drawing {}x{} rectangle", self.width, self.height); }
}

impl Draw for Triangle {
    fn draw(&self) { println!("Drawing triangle"); }
}

fn main() {
    // Vec of trait objects — different types, unified interface
    let shapes: Vec<Box<dyn Draw>> = vec![
        Box::new(Circle { radius: 5.0 }),
        Box::new(Rectangle { width: 4.0, height: 3.0 }),
        Box::new(Triangle { base: 6.0, height: 4.0 }),
    ];

    for shape in &shapes {
        shape.draw();  // dynamic dispatch — resolved at runtime
    }
}
```

### Static vs Dynamic Dispatch

```rust
// STATIC DISPATCH (monomorphization) — zero cost, compiler generates specific code
fn draw_static<T: Draw>(shape: &T) {
    shape.draw();
}

// DYNAMIC DISPATCH (vtable) — small overhead, flexible at runtime
fn draw_dynamic(shape: &dyn Draw) {
    shape.draw();
}

// When to use each:
// - impl Trait / generics → static dispatch → prefer this when possible
// - dyn Trait → dynamic dispatch → when type varies at runtime (heterogeneous collections)
```

### Object Safety

Not all traits can be made into trait objects. A trait is object-safe if:
1. It has no methods that return `Self`
2. It has no generic methods

```rust
// NOT object-safe — Clone requires knowing the concrete type size
// let _: Box<dyn Clone> = Box::new(5);  // ERROR

// Object-safe — Draw doesn't return Self or use generics
let _: Box<dyn Draw> = Box::new(Circle { radius: 1.0 });  // OK
```

---

## 11.6 Deriving Traits

Many common traits can be automatically derived with `#[derive(...)]`:

```rust
#[derive(Debug, Clone, PartialEq, Eq, PartialOrd, Ord, Hash, Default)]
struct Key {
    name: String,
    value: i32,
}

// Using derived traits
fn main() {
    let k1 = Key { name: String::from("a"), value: 1 };
    let k2 = k1.clone();
    println!("{}", k1 == k2);  // true — PartialEq

    let mut map = std::collections::HashMap::new();
    map.insert(k1, "value");  // Key implements Hash + Eq — can be a HashMap key

    let default_key = Key::default();  // Default
    println!("{:?}", default_key);
}
```

---

## Complete Example: A Generic Data Store

```rust
use std::collections::HashMap;

trait Identifiable {
    fn id(&self) -> u32;
}

trait Describable {
    fn describe(&self) -> String;
}

struct Store<T: Identifiable + Clone> {
    items: HashMap<u32, T>,
}

impl<T: Identifiable + Clone> Store<T> {
    fn new() -> Self {
        Store { items: HashMap::new() }
    }

    fn insert(&mut self, item: T) {
        self.items.insert(item.id(), item);
    }

    fn get(&self, id: u32) -> Option<&T> {
        self.items.get(&id)
    }

    fn remove(&mut self, id: u32) -> Option<T> {
        self.items.remove(&id)
    }

    fn count(&self) -> usize {
        self.items.len()
    }

    fn all(&self) -> Vec<&T> {
        self.items.values().collect()
    }
}

impl<T: Identifiable + Clone + Describable> Store<T> {
    fn describe_all(&self) {
        for item in self.all() {
            println!("[{}] {}", item.id(), item.describe());
        }
    }
}

// Concrete types
#[derive(Clone)]
struct Product {
    id: u32,
    name: String,
    price: f64,
}

impl Identifiable for Product {
    fn id(&self) -> u32 { self.id }
}

impl Describable for Product {
    fn describe(&self) -> String {
        format!("{} (${:.2})", self.name, self.price)
    }
}

fn main() {
    let mut store: Store<Product> = Store::new();

    store.insert(Product { id: 1, name: "Widget".into(), price: 9.99 });
    store.insert(Product { id: 2, name: "Gadget".into(), price: 24.99 });
    store.insert(Product { id: 3, name: "Doohickey".into(), price: 4.99 });

    println!("Products: {}", store.count());
    store.describe_all();

    if let Some(p) = store.get(2) {
        println!("Found: {}", p.describe());
    }
}
```

---

## Summary

Generics let you write code once and use it with any type. The compiler monomorphizes generic code — expands it into concrete versions — so there's zero runtime cost. Traits define shared behavior (like interfaces). Trait bounds on generics constrain which types can be used. `impl Trait` in function signatures is shorthand for simple generic bounds. `dyn Trait` enables dynamic dispatch for heterogeneous collections. Standard library traits (`Debug`, `Clone`, `Display`, `From`, `Into`, `Default`, `PartialEq`, `Ord`) can often be derived automatically.

---

## Key Takeaways

- Generics are zero-cost — compiled to concrete code per type (monomorphization)
- Trait bounds restrict generic types: `T: Display + Clone`
- `impl Trait` = static dispatch (generics); `dyn Trait` = dynamic dispatch (vtable)
- Use static dispatch by default; switch to `dyn Trait` for heterogeneous collections
- `From<T>` automatically provides `Into<T>` — implement `From`, get `Into` for free
- `#[derive(...)]` auto-implements common traits when all fields implement them
- Traits can have default method implementations — implementors can override selectively

---

## Exercises

**Exercise 1:** Write a generic `Stack<T>` with `push`, `pop`, `peek`, `is_empty`, and implement `Display` for `Stack<T> where T: Display`.

**Exercise 2:** Define a trait `Area` with method `area(&self) -> f64`. Implement it for `Circle`, `Rectangle`, and `Triangle`. Write `largest_area(shapes: &[&dyn Area]) -> f64`.

**Exercise 3:** Implement `From<(f64, f64)>` for `Point` and `From<Point>` for `(f64, f64)`. Verify that `Into` works automatically.

**Exercise 4:** Write a generic function `zip_map<T, U, V, F>(a: &[T], b: &[U], f: F) -> Vec<V> where F: Fn(&T, &U) -> V` that combines two slices element-by-element using a function.

**Exercise 5:** Create a `Cache<K, V>` struct using `HashMap`. Add methods `get_or_insert(&mut self, key: K, compute: impl Fn() -> V) -> &V` that returns the cached value or computes and stores it.

---

*Next: [Chapter 12 — Collections](12-collections.md)*
