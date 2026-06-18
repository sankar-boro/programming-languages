# Chapter 15 — TypeScript Internals

> *"Understanding how TypeScript thinks about types — erasure, assignability, structural subtyping, variance — makes you a better TypeScript programmer. These aren't academic concepts; they explain why TypeScript accepts and rejects code in ways that surprise you."*

---

## 15.1 Type Erasure — Types Exist Only at Compile Time

TypeScript's most fundamental property: **types do not exist at runtime**. The compiler strips all type annotations before producing JavaScript.

```typescript
// TypeScript source
interface User {
  id: number;
  name: string;
}

function greet(user: User): string {
  return `Hello, ${user.name}!`;
}

const alice: User = { id: 1, name: "Alice" };
const message: string = greet(alice);
```

```javascript
// Compiled JavaScript — notice what's gone:
function greet(user) {    // : User removed
  return `Hello, ${user.name}!`;
}

const alice = { id: 1, name: "Alice" };  // : User removed
const message = greet(alice);             // : string removed
// The 'interface User' declaration is gone entirely — no trace of it
```

### What This Means In Practice

```typescript
// This works at compile time but fails at runtime:
function isUser(value: unknown): value is User {
  // At runtime, 'User' doesn't exist — you can't check it
  // return value instanceof User;  // ERROR at runtime — User has no runtime representation

  // You must check the shape manually:
  return (
    typeof value === "object" &&
    value !== null &&
    typeof (value as Record<string, unknown>).id === "number" &&
    typeof (value as Record<string, unknown>).name === "string"
  );
}

// Type assertions have zero runtime cost — they're compile-time only
const raw: unknown = fetchData();
const user = raw as User;  // no runtime check! Just tells TypeScript to trust you
// If raw isn't a User, you'll get a runtime error when you use user.name

// typeof only knows JavaScript types — not TypeScript interfaces
type MyType = { x: number };
const val: MyType = { x: 1 };
console.log(typeof val);  // "object" — not "MyType"

// Enums DO have runtime representation — a special case
enum Direction { Up = "UP", Down = "DOWN" }
console.log(Direction.Up);   // "UP" — exists at runtime
console.log(typeof Direction);  // "object" — it's a real JS object

// const enums are erased:
const enum Color { Red, Green, Blue }
// Color.Red gets inlined as 0 — the enum object itself disappears
```

---

## 15.2 The Type System — Compile-Time vs Runtime

TypeScript exists in two parallel worlds:

```
┌─────────────────────────────────────────────┐
│           COMPILE TIME (TypeScript)         │
│  Types, interfaces, generics, type checks   │
│  Exists only in .ts files and type checker  │
└────────────────────────┬────────────────────┘
                         │ tsc compiles
                         ▼
┌─────────────────────────────────────────────┐
│            RUNTIME (JavaScript)             │
│  Values, objects, functions, classes        │
│  typeof, instanceof, JSON.parse work here   │
└─────────────────────────────────────────────┘
```

```typescript
// Type space vs value space — the same name can mean different things

interface Animal { name: string }  // type space — erased
class Animal { name: string = ""; }  // value space — real JS class

// typeof in type position vs value position
type StringType = typeof "hello";    // type position: string
const str = typeof "hello";          // value position: "object" (JavaScript's typeof)

function process(x: typeof "hello"): void {}  // type-position typeof
const t = typeof x;                           // value-position typeof

// import type — purely type space
import type { User } from "./types";  // erased — zero runtime bytes
import { createUser } from "./service";  // value — real function
```

---

## 15.3 Structural Typing — How Type Compatibility Works

TypeScript uses **structural typing**: two types are compatible if they have the same structure (shape), regardless of their names. This is sometimes called "duck typing."

### The Assignability Rule

Type `A` is assignable to type `B` if `A` has all the properties that `B` requires (and possibly more):

```typescript
interface Animal {
  name: string;
}

interface Dog {
  name: string;
  breed: string;
}

// Dog has ALL of Animal's properties (and more)
// So Dog is assignable to Animal

function greet(animal: Animal): void {
  console.log(`Hello, ${animal.name}!`);
}

const dog: Dog = { name: "Rex", breed: "Labrador" };
greet(dog);  // OK — Dog is assignable to Animal (Dog ⊆ Animal in type terms)
// The extra 'breed' property is fine — Dog satisfies the Animal contract

// Reversed: Animal is NOT assignable to Dog
function showBreed(dog: Dog): void {
  console.log(dog.breed);
}

const animal: Animal = { name: "Unknown" };
// showBreed(animal);  // ERROR — animal doesn't have 'breed'
```

### Nominal vs Structural — Why This Matters

```typescript
// In Java/C# (nominal typing):
// class A {} class B {}  — A and B are incompatible even if they have the same members

// In TypeScript (structural typing):
class A { x: number = 0; }
class B { x: number = 0; }

function takeA(a: A): void {}
const b = new B();
takeA(b);  // OK! B is structurally identical to A

// This is why branded types work:
type USD = number & { readonly __brand: "USD" };
type EUR = number & { readonly __brand: "EUR" };

const price: USD = 9.99 as USD;
const tax: EUR = 2.50 as EUR;

function totalCost(price: USD, tax: USD): USD {
  return (price + tax) as USD;
}

// totalCost(price, tax);  // ERROR — EUR is not structurally compatible with USD
// (because __brand has different literal types)
```

---

## 15.4 Assignability — The Full Rules

TypeScript's assignability algorithm determines when one type can be used where another is expected.

### Primitive Assignability

```typescript
// Literal types are subtypes of their base types
type A = "hello";
type B = string;

const a: A = "hello";
const b: B = a;  // OK — "hello" is assignable to string (subtype)

const b2: B = "anything";
// const a2: A = b2;  // ERROR — string is not assignable to "hello"

// Numeric literal
type Forty = 40;
const n: Forty = 40;
const num: number = n;  // OK — 40 is subtype of number
// const m: Forty = num;  // ERROR
```

### Object Assignability

```typescript
// An object with more properties is assignable to one with fewer (not the reverse)
type Minimal = { id: number };
type Full = { id: number; name: string; email: string };

const full: Full = { id: 1, name: "Alice", email: "alice@ex.com" };
const minimal: Minimal = full;  // OK — Full ⊆ Minimal (Full is a subtype)
// const full2: Full = minimal;  // ERROR — Minimal missing name, email

// This is called "width subtyping" or "structural subtyping"
// More properties = more specific type = subtype
```

### Union Assignability

```typescript
// Narrower union is assignable to wider union
type Narrow = "a" | "b";
type Wide = "a" | "b" | "c";

const narrow: Narrow = "a";
const wide: Wide = narrow;  // OK — {"a","b"} ⊆ {"a","b","c"}
// const narrow2: Narrow = wide;  // ERROR — wide might be "c"

// never is assignable to everything (bottom type)
type N = never;
const n: N = null as never;
const anyValue: string = n;  // OK — never is assignable to anything
```

### Function Assignability — Covariance and Contravariance

This is where TypeScript gets nuanced:

```typescript
// Return types: COVARIANT — more specific return type is OK
type FnReturnsAnimal = () => Animal;
type FnReturnsDog = () => Dog;

// A function returning Dog can be used where Animal is expected
// (Dog has all of Animal's properties — it satisfies the Animal contract)
const returnsAnimal: FnReturnsAnimal = (() => ({ name: "Rex", breed: "Lab" } as Dog));
// OK: Dog is assignable to Animal, so () => Dog is assignable to () => Animal

// Parameter types: CONTRAVARIANT — less specific param is OK (with strictFunctionTypes)
type HandlerAnimal = (animal: Animal) => void;
type HandlerDog = (dog: Dog) => void;

// A handler for Animal can be used where a Dog handler is expected
// (If you can handle any Animal, you can handle a Dog)
const handleDog: HandlerDog = ((animal: Animal) => console.log(animal.name));
// OK: HandlerAnimal is assignable to HandlerDog (contravariant in parameter)

// A handler for Dog CANNOT be used where Animal handler is expected
// (The Dog handler might access .breed — which Animal doesn't have)
// const handleAnimal: HandlerAnimal = ((dog: Dog) => console.log(dog.breed));
// ERROR: HandlerDog is NOT assignable to HandlerAnimal
```

```
Summary of variance:
  Return types:    COVARIANT     — Dog ≤ Animal, so () => Dog ≤ () => Animal
  Parameter types: CONTRAVARIANT — Animal ≤ Dog → (Dog→void) ≤ (Animal→void)
  Generic types:   depends on usage (TypeScript approximates)
```

---

## 15.5 Type Widening — How TypeScript Infers Types

TypeScript "widens" literal types to their base type in many situations:

```typescript
// Let — TypeScript widens to the base type
let x = "hello";  // x: string (widened from "hello")
x = "world";      // OK — x is string

// Const — TypeScript keeps the literal type
const y = "hello";  // y: "hello" (literal type preserved)
// y = "world";      // ERROR — y is "hello", can't reassign

// In objects — TypeScript widens
const obj = { name: "Alice" };  // obj.name: string (not "Alice")
obj.name = "Bob";               // OK — name is string

// as const — prevents widening
const obj2 = { name: "Alice" } as const;  // obj2.name: "Alice" (literal)
// obj2.name = "Bob";  // ERROR — readonly

// Array widening
const arr = ["a", "b", "c"];  // string[] (widened)
const tuple = ["a", "b", "c"] as const;  // readonly ["a", "b", "c"] (literal tuple)

// Function return type widening
function getDirection() {
  return "left";  // returns string (widened)
}

function getDirection2(): "left" | "right" {
  return "left";  // returns "left" | "right"
}
```

---

## 15.6 Control Flow and Type Narrowing Internals

TypeScript's type narrowing is powered by **control flow analysis** — a dataflow analysis that tracks types through every code path.

```typescript
// How TypeScript tracks types through branches
function process(value: string | null | undefined): string {
  // Entry: value is string | null | undefined

  if (value === null) {
    // Branch 1: value is null — return early
    return "null";
  }
  // After branch 1: value is string | undefined (null eliminated)

  if (value === undefined) {
    // Branch 2: value is undefined — return early
    return "undefined";
  }
  // After branch 2: value is string (null and undefined eliminated)

  return value.toUpperCase();  // TypeScript knows: string
}

// Narrowing through assignments
function example(x: string | number): void {
  // x: string | number

  let y: string | number = x;
  // y: string | number

  if (typeof x === "string") {
    // x: string
    y = x.toUpperCase();  // y: string (narrowed from assignment)
  }

  // y: string | number (back to the union after if)
  // (TypeScript doesn't know if we took the if branch)
}
```

### Why Some Narrowing Fails

```typescript
// TypeScript can't narrow through function calls (in general)
function isString(x: unknown): x is string {
  return typeof x === "string";
}

// Works: type guard function
function process(x: unknown): void {
  if (isString(x)) {
    x.toUpperCase();  // OK — isString is a type guard
  }
}

// Fails: arbitrary boolean conditions don't narrow
let userValidator = (x: unknown): boolean => typeof x === "string";

function process2(x: unknown): void {
  if (userValidator(x)) {
    // x.toUpperCase();  // ERROR — userValidator doesn't narrow (returns boolean, not x is string)
  }
}

// TypeScript also can't narrow across async boundaries
async function asyncNarrow(value: string | null): Promise<void> {
  if (value !== null) {
    // value: string — narrowed
    await someAsyncOperation();  // value might have changed?
    value.toUpperCase();  // still string — TypeScript is optimistic here
  }
}
```

---

## 15.7 The never Type — The Bottom Type

`never` is the bottom type — it's a subtype of every type, and no type is a subtype of `never` (except `never` itself).

```typescript
// Functions that never return have type 'never'
function throwError(msg: string): never {
  throw new Error(msg);
}

function infiniteLoop(): never {
  while (true) {}
}

// 'never' is useful for exhaustiveness checking
type Shape = "circle" | "square" | "triangle";

function area(shape: Shape): number {
  switch (shape) {
    case "circle": return Math.PI * 5 ** 2;
    case "square": return 25;
    case "triangle": return 12.5;
    default:
      // If you add "hexagon" to Shape without handling it here,
      // TypeScript errors because shape is "hexagon", not never
      const exhaustive: never = shape;
      throw new Error(`Unknown shape: ${exhaustive}`);
  }
}

// Filtering with 'never'
type NonNullable2<T> = T extends null | undefined ? never : T;
// "a" | null | undefined → "a" | never | never → "a"

// Conditional types with 'never'
type ExcludeFromUnion<T, U> = T extends U ? never : T;
// Distributes over union, replacing matching members with never
// never is absorbed in a union: string | never = string
```

---

## 15.8 Declaration Spaces — Type Space and Value Space

TypeScript has two declaration spaces: **types** and **values**. Understanding this explains many confusing errors.

```typescript
// A name can exist in type space, value space, or both

// interface — type space only (erased)
interface IFoo { x: number }

// class — both type space AND value space
class Foo {
  x: number = 0;
}

// Using Foo in type position (type space):
const foo: Foo = new Foo();  // Foo as a type

// Using Foo in value position (value space):
const instance = new Foo();    // Foo as a class constructor
const isFoo = instance instanceof Foo;  // Foo as a value

// enum — both spaces (with caveats)
enum Direction { Up = "UP", Down = "DOWN" }
const d: Direction = Direction.Up;  // type and value
console.log(Direction.Up);          // value at runtime

// type alias — type space only
type Bar = { y: string };
const bar: Bar = { y: "hello" };
// new Bar();  // ERROR — Bar is not a value

// namespace — value space (and type space for inner types)
namespace Utils {
  export function greet(): string { return "hello"; }
  export interface Config { debug: boolean }
}

Utils.greet();      // value — works at runtime
const c: Utils.Config = { debug: true };  // type — erased
```

---

## 15.9 Variance in Generic Types

Generic types have variance — how they relate to their type parameter's subtype relationship:

```typescript
// Covariant: if A extends B, then Container<A> extends Container<B>
// Appropriate for read-only containers
interface ReadonlyBox<T> {
  readonly value: T;  // only read (returns T, doesn't accept T)
}

// ReadonlyBox<Dog> is assignable to ReadonlyBox<Animal>
// because Dog extends Animal, and reading gives you Animal (or better)
const dogBox: ReadonlyBox<Dog> = { value: { name: "Rex", breed: "Lab" } };
const animalBox: ReadonlyBox<Animal> = dogBox;  // OK — covariant

// Contravariant: if A extends B, then Container<B> extends Container<A>
// Appropriate for write-only containers (consumers)
type Consumer<T> = (value: T) => void;

// Consumer<Animal> is assignable to Consumer<Dog>
// A function that accepts any Animal certainly handles a Dog
const animalConsumer: Consumer<Animal> = (a) => console.log(a.name);
const dogConsumer: Consumer<Dog> = animalConsumer;  // OK — contravariant

// Invariant: no subtype relationship in either direction
// Mutable containers are invariant (both read and write)
interface MutableBox<T> {
  get(): T;    // reads: covariant
  set(v: T): void;  // writes: contravariant
  // together: invariant
}

// TypeScript approximates with bivariance for method parameters (historical)
// but is stricter with function property parameters (strictFunctionTypes)
```

---

## 15.10 Type Compatibility Deep Dive

```typescript
// Excess property checking vs structural compatibility
interface Options {
  timeout: number;
  retries?: number;
}

// Object literal assignment — excess property check (stricter)
// const opts: Options = { timeout: 5000, extraProp: true };  // ERROR

// Variable assignment — structural check (looser)
const config = { timeout: 5000, extraProp: true };
const opts: Options = config;  // OK — structural compatibility

// Why the difference?
// Object literals are likely bugs (you typed the wrong name)
// Variables might legitimately have extra properties that are fine to ignore

// Discriminated union compatibility
type Result<T> =
  | { success: true; data: T }
  | { success: false; error: Error };

// A { success: true; data: string; metadata: string } is assignable to Result<string>
// because the discriminant 'success' matches and the required fields are present
const res: Result<string> = {
  success: true,
  data: "hello",
  // metadata: "extra"  // excess property check catches this for literals
};

// But through an intermediate variable:
const withExtra = { success: true as const, data: "hello", metadata: "extra" };
const res2: Result<string> = withExtra;  // OK — structural check passes
```

---

## Summary

TypeScript's types are purely compile-time constructs — they're erased completely from the generated JavaScript. The type system uses structural typing (duck typing): two types are compatible based on their shape, not their names. Assignability follows the rule: "A is assignable to B if A has everything B requires." Type widening converts literal types to base types in mutable contexts; `as const` prevents this. Control flow analysis tracks types through every code path — narrowing removes impossible types. `never` is the bottom type: a subtype of everything, used for exhaustiveness and impossible cases. Understanding variance (covariance/contravariance) explains why some generic types are or aren't compatible.

---

## Key Takeaways

- **Types are erased at runtime** — no type information survives compilation
- **Structural typing**: compatibility is about shape, not name — two identically-shaped types are compatible
- **Assignability**: A is assignable to B if A has all of B's required properties (and possibly more)
- **Type widening**: `let` widens literals to base types; `const` and `as const` preserve literals
- **Control flow analysis**: TypeScript tracks types through every if/else/switch branch
- **`never`**: the bottom type — subtype of everything; used for exhaustiveness and impossible states
- **Variance**: return types are covariant; function parameters are contravariant (with `strictFunctionTypes`)

---

## Practice Questions

1. What does "type erasure" mean? What runtime representation do interfaces have?
2. Why is TypeScript's type system "structural" rather than "nominal"?
3. What is the assignability rule for object types?
4. What is the difference between type widening for `let` vs `const`?
5. What is `never`, and why is it useful for exhaustiveness checking?
6. Explain covariance and contravariance in the context of function types.

---

## Exercises

**Exercise 1**: Demonstrate type erasure: write a TypeScript file with interfaces and type aliases, compile it to JavaScript, and annotate what was removed.

**Exercise 2**: Create an example showing structural typing: two unrelated interfaces with the same shape, and a function that accepts one but is given the other. Explain why TypeScript accepts this.

**Exercise 3**: Write a generic `Invariant<T>` type that has both a getter and setter for `T`. Demonstrate that `Invariant<Dog>` is not assignable to `Invariant<Animal>` even though `Dog extends Animal`.

**Exercise 4**: Build a complete exhaustiveness checker: a function `assertNever(value: never): never` and a union with 5 variants. Show that adding a 6th variant without handling it causes a compile error.

---

*Next: [Chapter 90 — Best Practices](90-best-practices.md)*
