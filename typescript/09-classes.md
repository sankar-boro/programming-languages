# Chapter 9 — Classes

> *"Classes in TypeScript are full-featured: they compile to JavaScript classes but add access modifiers, readonly, abstract, and the ability to implement interfaces — turning classes into a powerful tool for encapsulation and polymorphism."*

---

## 9.1 Class Basics

TypeScript classes are JavaScript classes plus types. Everything valid in ES2020 classes is valid in TypeScript, plus additional features.

```typescript
class Person {
  // Properties must be declared before use
  name: string;
  age: number;

  constructor(name: string, age: number) {
    this.name = name;
    this.age = age;
  }

  greet(): string {
    return `Hello, I'm ${this.name}, ${this.age} years old.`;
  }

  birthday(): void {
    this.age++;
  }
}

const alice = new Person("Alice", 30);
console.log(alice.greet());  // "Hello, I'm Alice, 30 years old."
alice.birthday();
console.log(alice.age);      // 31

// TypeScript knows the shape of the instance
function describe(person: Person): string {
  return `${person.name} is ${person.age}`;
}
```

### Shorthand Property Declaration

TypeScript offers parameter properties — a shorthand to declare and initialize in one step:

```typescript
class Point {
  constructor(
    public x: number,   // public shorthand
    public y: number    // declares and assigns this.x, this.y
  ) {}

  distanceTo(other: Point): number {
    const dx = this.x - other.x;
    const dy = this.y - other.y;
    return Math.sqrt(dx * dx + dy * dy);
  }

  toString(): string {
    return `(${this.x}, ${this.y})`;
  }
}

// Equivalent to:
// class Point {
//   public x: number;
//   public y: number;
//   constructor(x: number, y: number) { this.x = x; this.y = y; }
// }

const p1 = new Point(0, 0);
const p2 = new Point(3, 4);
console.log(p1.distanceTo(p2));  // 5
```

---

## 9.2 Access Modifiers

TypeScript adds `public`, `private`, and `protected` to control visibility.

### public (default)

```typescript
class Car {
  public make: string;     // accessible from anywhere
  public model: string;
  public year: number;

  constructor(make: string, model: string, year: number) {
    this.make = make;
    this.model = model;
    this.year = year;
  }
}

const car = new Car("Toyota", "Camry", 2024);
car.make;   // accessible from outside class
```

### private — TypeScript + JavaScript Private

```typescript
class BankAccount {
  private balance: number = 0;   // TypeScript's private — compile-time only
  #pin: string;                  // JavaScript's # private — runtime enforced

  constructor(initialBalance: number, pin: string) {
    this.balance = initialBalance;
    this.#pin = pin;
  }

  deposit(amount: number): void {
    if (amount <= 0) throw new Error("Amount must be positive");
    this.balance += amount;
  }

  withdraw(amount: number, pin: string): boolean {
    if (pin !== this.#pin) {
      console.log("Invalid PIN");
      return false;
    }
    if (amount > this.balance) {
      console.log("Insufficient funds");
      return false;
    }
    this.balance -= amount;
    return true;
  }

  getBalance(pin: string): number | null {
    if (pin !== this.#pin) return null;
    return this.balance;
  }
}

const account = new BankAccount(1000, "1234");
// account.balance;  // TypeScript ERROR — private
// account.#pin;     // TypeError — JavaScript runtime error
account.deposit(500);
account.withdraw(200, "1234");
console.log(account.getBalance("1234")); // 1300
```

**Key difference**: TypeScript's `private` is erased at compile time — at runtime the property is public. JavaScript's `#` is a runtime privacy mechanism.

### protected — Accessible in Subclasses

```typescript
class Animal {
  protected name: string;
  private health: number;

  constructor(name: string, health: number) {
    this.name = name;
    this.health = health;
  }

  protected takeDamage(amount: number): void {
    this.health = Math.max(0, this.health - amount);
  }

  isAlive(): boolean {
    return this.health > 0;
  }
}

class Dog extends Animal {
  private tricks: string[] = [];

  constructor(name: string) {
    super(name, 100);  // must call super() before accessing this
  }

  speak(): string {
    return `${this.name} says: Woof!`;  // OK — protected
    // this.health  // ERROR — private to Animal
  }

  learnTrick(trick: string): void {
    this.tricks.push(trick);
  }

  getDamaged(amount: number): void {
    this.takeDamage(amount);  // OK — protected method
  }
}

const dog = new Dog("Rex");
console.log(dog.speak());  // "Rex says: Woof!"
// dog.name;               // ERROR — protected, not accessible outside class hierarchy
```

---

## 9.3 readonly Properties

`readonly` properties can only be set during declaration or in the constructor:

```typescript
class Config {
  readonly host: string;
  readonly port: number;
  mutable: string = "changeable";

  constructor(host: string, port: number) {
    this.host = host;  // OK — setting in constructor
    this.port = port;
  }

  update(): void {
    // this.host = "newhost";  // ERROR — cannot reassign readonly
    this.mutable = "new value";  // OK
  }
}

const config = new Config("localhost", 3000);
// config.host = "other";  // ERROR — readonly
console.log(config.host);   // "localhost" — can read

// Combine with parameter properties
class ImmutablePoint {
  constructor(
    public readonly x: number,
    public readonly y: number
  ) {}
}

const p = new ImmutablePoint(3, 4);
// p.x = 5;  // ERROR — readonly
```

---

## 9.4 Getters and Setters

```typescript
class Temperature {
  private _celsius: number;

  constructor(celsius: number) {
    this._celsius = celsius;
  }

  // Getter
  get celsius(): number {
    return this._celsius;
  }

  // Setter with validation
  set celsius(value: number) {
    if (value < -273.15) {
      throw new RangeError("Temperature below absolute zero");
    }
    this._celsius = value;
  }

  // Derived computed property
  get fahrenheit(): number {
    return this._celsius * 9 / 5 + 32;
  }

  set fahrenheit(value: number) {
    this.celsius = (value - 32) * 5 / 9;  // triggers celsius setter validation
  }

  get kelvin(): number {
    return this._celsius + 273.15;
  }
}

const temp = new Temperature(100);
console.log(temp.celsius);     // 100
console.log(temp.fahrenheit);  // 212
console.log(temp.kelvin);      // 373.15

temp.fahrenheit = 32;
console.log(temp.celsius);  // 0
// temp.celsius = -300;     // RangeError!

// Getter-only (no setter) — derived computed property
class Circle {
  constructor(public readonly radius: number) {}

  get area(): number {
    return Math.PI * this.radius ** 2;
  }

  get circumference(): number {
    return 2 * Math.PI * this.radius;
  }
}
```

---

## 9.5 Static Members

Static members belong to the class itself, not to instances:

```typescript
class MathUtils {
  static readonly PI = 3.14159265358979;

  static square(n: number): number {
    return n * n;
  }

  static cube(n: number): number {
    return n * n * n;
  }

  static hypotenuse(a: number, b: number): number {
    return Math.sqrt(MathUtils.square(a) + MathUtils.square(b));
  }
}

// Call on the class, not on instances
MathUtils.square(5);         // 25
MathUtils.hypotenuse(3, 4);  // 5
// const m = new MathUtils();
// m.square(5);  // Not idiomatic — square is static

// Static factory methods — the classic pattern
class User {
  private constructor(  // private constructor — must use factory
    public readonly id: string,
    public readonly name: string,
    public readonly email: string,
    public readonly createdAt: Date
  ) {}

  static create(name: string, email: string): User {
    return new User(
      Math.random().toString(36).slice(2),
      name,
      email,
      new Date()
    );
  }

  static fromJson(json: { id: string; name: string; email: string; createdAt: string }): User {
    return new User(json.id, json.name, json.email, new Date(json.createdAt));
  }
}

const alice = User.create("Alice", "alice@example.com");
// const bad = new User(...);  // ERROR — constructor is private

// Static counter / singleton
class RequestCounter {
  private static count = 0;

  static increment(): number {
    return ++RequestCounter.count;
  }

  static getCount(): number {
    return RequestCounter.count;
  }

  static reset(): void {
    RequestCounter.count = 0;
  }
}
```

---

## 9.6 Inheritance

TypeScript classes support single inheritance with `extends`:

```typescript
class Shape {
  constructor(public readonly color: string) {}

  area(): number {
    return 0;  // default — subclasses should override
  }

  toString(): string {
    return `${this.constructor.name}(color=${this.color}, area=${this.area().toFixed(2)})`;
  }
}

class Circle extends Shape {
  constructor(color: string, public readonly radius: number) {
    super(color);  // MUST call super() before using 'this'
  }

  override area(): number {  // 'override' keyword catches accidental shadowing
    return Math.PI * this.radius ** 2;
  }
}

class Rectangle extends Shape {
  constructor(
    color: string,
    public readonly width: number,
    public readonly height: number
  ) {
    super(color);
  }

  override area(): number {
    return this.width * this.height;
  }
}

class Square extends Rectangle {
  constructor(color: string, side: number) {
    super(color, side, side);
  }
}

const shapes: Shape[] = [
  new Circle("red", 5),
  new Rectangle("blue", 4, 6),
  new Square("green", 3),
];

shapes.forEach((s) => console.log(s.toString()));
// Circle(color=red, area=78.54)
// Rectangle(color=blue, area=24.00)
// Square(color=green, area=9.00)

// instanceof works correctly with inheritance
console.log(shapes[2] instanceof Square);     // true
console.log(shapes[2] instanceof Rectangle);  // true — Square extends Rectangle
console.log(shapes[2] instanceof Shape);      // true
```

### The override Keyword

```typescript
class Base {
  method(): void {
    console.log("base");
  }
}

class Derived extends Base {
  // override keyword: TypeScript errors if 'method' doesn't exist in Base
  override method(): void {
    console.log("derived");
  }

  // method2(): void {}  // OK — a new method
  // override method3(): void {}  // ERROR — method3 doesn't exist in Base
}
```

---

## 9.7 Abstract Classes

Abstract classes cannot be instantiated directly — they serve as blueprints for subclasses. Abstract methods must be implemented by subclasses.

```typescript
abstract class DataSource<T> {
  // Abstract method: must be implemented by subclasses
  abstract connect(): Promise<void>;
  abstract findAll(): Promise<T[]>;
  abstract findById(id: string): Promise<T | null>;
  abstract insert(item: T): Promise<T>;
  abstract update(id: string, item: Partial<T>): Promise<T | null>;
  abstract delete(id: string): Promise<boolean>;
  abstract disconnect(): Promise<void>;

  // Concrete method: inherited by subclasses as-is
  async exists(id: string): Promise<boolean> {
    return (await this.findById(id)) !== null;
  }

  async insertMany(items: T[]): Promise<T[]> {
    return Promise.all(items.map((item) => this.insert(item)));
  }
}

interface User {
  id: string;
  name: string;
  email: string;
}

// Concrete implementation
class InMemoryUserSource extends DataSource<User> {
  private store = new Map<string, User>();

  async connect(): Promise<void> {
    console.log("InMemory: connected");
  }

  async findAll(): Promise<User[]> {
    return [...this.store.values()];
  }

  async findById(id: string): Promise<User | null> {
    return this.store.get(id) ?? null;
  }

  async insert(user: User): Promise<User> {
    this.store.set(user.id, user);
    return user;
  }

  async update(id: string, data: Partial<User>): Promise<User | null> {
    const user = this.store.get(id);
    if (!user) return null;
    const updated = { ...user, ...data };
    this.store.set(id, updated);
    return updated;
  }

  async delete(id: string): Promise<boolean> {
    return this.store.delete(id);
  }

  async disconnect(): Promise<void> {
    this.store.clear();
    console.log("InMemory: disconnected");
  }
}

// const ds = new DataSource<User>();  // ERROR: Cannot create instance of abstract class
const ds = new InMemoryUserSource();   // OK

// The abstract class type is valid as a type reference
async function runMigration(source: DataSource<User>): Promise<void> {
  await source.connect();
  const users = await source.findAll();
  console.log(`Found ${users.length} users`);
  await source.disconnect();
}
```

---

## 9.8 Implementing Interfaces

A class can implement one or more interfaces. This is a compile-time contract.

```typescript
interface Printable {
  print(): void;
  toPrettyString(): string;
}

interface Serializable {
  serialize(): string;
  deserialize(data: string): void;
}

interface Comparable<T> {
  compareTo(other: T): number;  // negative: less, 0: equal, positive: greater
}

class Product implements Printable, Serializable, Comparable<Product> {
  constructor(
    public id: string,
    public name: string,
    public price: number
  ) {}

  print(): void {
    console.log(this.toPrettyString());
  }

  toPrettyString(): string {
    return `[${this.id}] ${this.name} — $${this.price.toFixed(2)}`;
  }

  serialize(): string {
    return JSON.stringify({ id: this.id, name: this.name, price: this.price });
  }

  deserialize(data: string): void {
    const parsed = JSON.parse(data);
    // Note: 'this' properties would need to be mutable for this to work
    Object.assign(this, parsed);
  }

  compareTo(other: Product): number {
    return this.price - other.price;
  }
}

const products = [
  new Product("p3", "Widget", 19.99),
  new Product("p1", "Gadget", 49.99),
  new Product("p2", "Doohickey", 9.99),
];

// Sort using Comparable
products.sort((a, b) => a.compareTo(b));
products.forEach((p) => p.print());
// [p2] Doohickey — $9.99
// [p3] Widget — $19.99
// [p1] Gadget — $49.99

// Type-safe: treat Product as any of its interfaces
function printAll(items: Printable[]): void {
  items.forEach((item) => item.print());
}
printAll(products);  // OK — Product implements Printable
```

---

## 9.9 The Singleton Pattern

```typescript
class DatabaseConnection {
  private static instance: DatabaseConnection | null = null;
  private isConnected = false;

  private constructor(
    private readonly host: string,
    private readonly port: number
  ) {}

  static getInstance(): DatabaseConnection {
    if (!DatabaseConnection.instance) {
      DatabaseConnection.instance = new DatabaseConnection("localhost", 5432);
    }
    return DatabaseConnection.instance;
  }

  async connect(): Promise<void> {
    if (this.isConnected) return;
    console.log(`Connecting to ${this.host}:${this.port}...`);
    this.isConnected = true;
  }

  async query(sql: string): Promise<unknown[]> {
    if (!this.isConnected) throw new Error("Not connected");
    console.log(`Executing: ${sql}`);
    return [];
  }
}

// Both return the same instance
const db1 = DatabaseConnection.getInstance();
const db2 = DatabaseConnection.getInstance();
console.log(db1 === db2);  // true
```

---

## 9.10 Mixins — Composing Behavior

TypeScript supports a mixin pattern for composing behaviors without deep inheritance:

```typescript
// Mixin: a function that takes a class and returns a new class
type Constructor<T = {}> = new (...args: unknown[]) => T;

// Timestamped mixin
function Timestamped<TBase extends Constructor>(Base: TBase) {
  return class extends Base {
    createdAt = new Date();
    updatedAt = new Date();

    touch(): void {
      this.updatedAt = new Date();
    }
  };
}

// Activatable mixin
function Activatable<TBase extends Constructor>(Base: TBase) {
  return class extends Base {
    isActive = false;

    activate(): void {
      this.isActive = true;
    }

    deactivate(): void {
      this.isActive = false;
    }
  };
}

// Base class
class Entity {
  constructor(public id: string) {}
}

// Compose mixins
class User extends Timestamped(Activatable(Entity)) {
  constructor(id: string, public name: string) {
    super(id);
  }
}

const user = new User("u1", "Alice");
console.log(user.createdAt);  // Date
console.log(user.isActive);   // false
user.activate();
console.log(user.isActive);   // true
user.touch();
console.log(user.updatedAt);  // updated Date
```

---

## Summary

TypeScript classes extend JavaScript classes with access modifiers (`public`, `private`, `protected`), `readonly` properties, and the `override` keyword for safe inheritance. Parameter properties (`constructor(public x: T)`) reduce boilerplate. Getters and setters allow computed and validated properties. Static members belong to the class constructor. Abstract classes define contracts without implementation — they enforce that subclasses provide specific methods. Classes can implement one or more interfaces, declaring they satisfy a contract. The `private constructor` pattern enables factories and singletons.

---

## Key Takeaways

- **Access modifiers** are a compile-time feature — TypeScript's `private` is not runtime-enforced; use `#` for true runtime privacy
- **`readonly`** prevents reassignment after construction
- **`override`** catches accidental shadowing — use it always when overriding
- **Abstract classes** = blueprint + partial implementation; cannot be instantiated
- **Implementing interfaces** declares a compile-time contract
- **Static factory methods** with `private constructor` are more flexible than public constructors
- **Mixins** compose behaviors without multiple inheritance

---

## Practice Questions

1. What is the difference between TypeScript's `private` and JavaScript's `#`?
2. What does `super()` do, and when must you call it?
3. What is an abstract class, and why can't it be instantiated?
4. What is the difference between `extends` and `implements`?
5. When would you use a `private constructor`?
6. What does the `override` keyword do? Why should you always use it?

---

## Exercises

**Exercise 1**: Implement a `LinkedList<T>` class with `push`, `pop`, `peek`, `isEmpty`, and iteration support (implement `[Symbol.iterator]`).

**Exercise 2**: Create an abstract `Logger` class with abstract methods `log(level: string, message: string): void` and `clear(): void`. Implement `ConsoleLogger` and `MemoryLogger` (stores logs in an array).

**Exercise 3**: Build a `EventEmitter` class using private `Map<string, Function[]>` internals with `on`, `off`, `once`, and `emit` methods. Make it generic: `EventEmitter<T extends Record<string, unknown[]>>`.

**Exercise 4**: Implement the Observer design pattern: an abstract `Subject<T>` that notifies observers when it changes, and an `Observer<T>` interface.

---

*Next: [Chapter 10 — Modules and Namespaces](10-modules.md)*
