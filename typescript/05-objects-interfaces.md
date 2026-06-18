# Chapter 5 — Objects and Interfaces

> *"Interfaces are TypeScript's primary way of naming object shapes. They are the vocabulary through which your codebase communicates what data looks like."*

---

## 5.1 Object Type Literals

The simplest way to type an object in TypeScript is with an **object type literal** — a description of the object's shape written inline.

```typescript
// Object type literal: { property: Type; property: Type }
let user: { id: number; name: string; active: boolean } = {
  id: 1,
  name: "Alice",
  active: true,
};

// TypeScript ensures you provide all required properties
let invalid: { id: number; name: string } = {
  id: 1,
  // ERROR: Property 'name' is missing in type '{ id: number; }' but required in type '{ id: number; name: string; }'
};

// Object type literals in function parameters
function displayUser(user: { id: number; name: string }): void {
  console.log(`User #${user.id}: ${user.name}`);
}

// TypeScript knows the shape — full autocomplete
displayUser({ id: 1, name: "Alice" });
// displayUser({ id: "one", name: "Alice" });  // ERROR — id must be number
```

Object type literals are fine for simple, one-off types. But when you need to reuse a type or give it a meaningful name, use interfaces or type aliases.

---

## 5.2 Interfaces — Naming Object Shapes

An `interface` declares a named type for an object shape. It's one of TypeScript's most-used features.

### Declaring an Interface

```typescript
interface User {
  id: number;
  name: string;
  email: string;
  active: boolean;
  createdAt: Date;
}

// Using the interface
const alice: User = {
  id: 1,
  name: "Alice",
  email: "alice@example.com",
  active: true,
  createdAt: new Date(),
};

// Interface as function parameter
function sendWelcomeEmail(user: User): void {
  console.log(`Sending welcome to ${user.email}`);
}

// Interface as return type
function createUser(name: string, email: string): User {
  return {
    id: Math.floor(Math.random() * 10000),
    name,
    email,
    active: true,
    createdAt: new Date(),
  };
}
```

### Methods in Interfaces

```typescript
interface Calculator {
  // Method signature (two equivalent syntaxes)
  add(a: number, b: number): number;
  subtract: (a: number, b: number) => number;  // property with function type

  // Optional method
  multiply?(a: number, b: number): number;
}

const calc: Calculator = {
  add(a, b) { return a + b; },
  subtract: (a, b) => a - b,
  // multiply is optional — OK to omit
};

const calcWithMultiply: Calculator = {
  add: (a, b) => a + b,
  subtract: (a, b) => a - b,
  multiply: (a, b) => a * b,
};
```

### Nested Interfaces

```typescript
interface Address {
  street: string;
  city: string;
  country: string;
  postalCode: string;
}

interface UserProfile {
  id: number;
  name: string;
  address: Address;         // nested interface
  addresses: Address[];     // array of interfaces
  primaryAddress?: Address; // optional nested interface
}

const profile: UserProfile = {
  id: 1,
  name: "Alice",
  address: {
    street: "123 Main St",
    city: "Springfield",
    country: "US",
    postalCode: "12345",
  },
  addresses: [
    { street: "123 Main St", city: "Springfield", country: "US", postalCode: "12345" },
    { street: "456 Oak Ave", city: "Shelbyville", country: "US", postalCode: "67890" },
  ],
};
```

---

## 5.3 Type Aliases

A `type` alias gives a name to any type — not just object shapes, but unions, intersections, primitives, and more.

```typescript
// Type alias for an object shape
type Point = {
  x: number;
  y: number;
};

// Type alias for primitives
type ID = string | number;
type Email = string;

// Type alias for union
type Status = "active" | "inactive" | "pending";

// Type alias for function type
type Transformer<T> = (value: T) => T;

// Type alias for complex generic type
type ApiResponse<T> = {
  data: T;
  status: number;
  message: string;
  timestamp: Date;
};

// Usage
const point: Point = { x: 3, y: 4 };
const userId: ID = "user-123";
const userIdNumber: ID = 42;
const status: Status = "active";
const double: Transformer<number> = (n) => n * 2;
const response: ApiResponse<User> = {
  data: { id: 1, name: "Alice", email: "alice@example.com", active: true, createdAt: new Date() },
  status: 200,
  message: "Success",
  timestamp: new Date(),
};
```

---

## 5.4 Interface vs Type Alias — The Real Differences

This is one of the most asked questions in TypeScript. They are more similar than different, but have key distinctions.

### Syntactic Differences

```typescript
// interface: only for object shapes (and functions with call signatures)
interface Point {
  x: number;
  y: number;
}

// type: for any type
type Point2 = {
  x: number;
  y: number;
};

// type can do things interface can't:
type ID = string | number;         // union — interface can't do this
type Tuple = [string, number];     // tuple — interface can't do this
type Fn = (x: number) => string;   // simple function — interface needs call signature
```

### Declaration Merging — Interfaces Only

```typescript
// Interfaces can be declared multiple times — they merge!
interface User {
  id: number;
  name: string;
}

interface User {
  email: string;  // added to the existing User interface
}

// User is now { id: number; name: string; email: string }
const user: User = { id: 1, name: "Alice", email: "alice@example.com" };

// Types CANNOT be merged:
type Product = { id: number };
// type Product = { name: string };  // ERROR: Duplicate identifier 'Product'
```

Declaration merging is used by `@types` packages to augment existing types (like adding properties to Express's `Request` object).

### Extending

```typescript
// Interface extending interface
interface Animal {
  name: string;
  sound(): string;
}

interface Dog extends Animal {
  breed: string;
}

// Interface extending type alias
type Serializable = { serialize(): string };
interface Config extends Serializable {
  version: number;
}

// Type alias extending via intersection
type Animal2 = { name: string; sound(): string };
type Dog2 = Animal2 & { breed: string };

// Multiple inheritance — only in interfaces
interface A { a: string }
interface B { b: number }
interface C extends A, B {
  c: boolean;
}
```

### Error Messages

```typescript
// Interfaces tend to produce cleaner error messages
// because TypeScript can refer to the named type

interface UserShape {
  id: number;
  name: string;
}

// type UserShape = { id: number; name: string };

const bad: UserShape = { id: "one", name: "Alice" };
// With interface: Type 'string' is not assignable to type 'number' in 'UserShape.id'
// With type: Type 'string' is not assignable to type 'number' in '{ id: string; name: string }.id'
// Interface error is cleaner
```

### When to Use Which

```typescript
// Use interface when:
// 1. Defining a public API for a library (allows declaration merging for users)
// 2. Defining class implementations
// 3. It's an object/class shape

// Use type when:
// 1. You need unions, intersections, or tuples
// 2. You need conditional or mapped types
// 3. You prefer type aliases for everything (valid preference — stay consistent)

// The TypeScript team's guidance: either works; be consistent
// Popular convention: interfaces for object shapes, types for everything else
```

---

## 5.5 Optional Properties

Properties marked with `?` may be present or absent:

```typescript
interface UserPreferences {
  theme: "light" | "dark";         // required
  language: string;                 // required
  timezone?: string;                // optional — may be absent
  notifications?: {                  // optional nested object
    email: boolean;
    push: boolean;
  };
  fontSize?: number;                // optional
}

const prefs: UserPreferences = {
  theme: "dark",
  language: "en",
  // timezone, notifications, fontSize are all optional
};

// Accessing optional properties
function getTimezone(prefs: UserPreferences): string {
  // prefs.timezone is string | undefined
  return prefs.timezone ?? "UTC";  // default if absent
}

// Optional chaining for nested optionals
function getEmailPref(prefs: UserPreferences): boolean {
  return prefs.notifications?.email ?? true;
}
```

---

## 5.6 Readonly Properties

`readonly` prevents property reassignment after object creation:

```typescript
interface ImmutablePoint {
  readonly x: number;
  readonly y: number;
}

const point: ImmutablePoint = { x: 3, y: 4 };
// point.x = 5;  // ERROR: Cannot assign to 'x' because it is a read-only property

// Readonly arrays
interface Config {
  readonly allowedMethods: readonly string[];
  version: string;
}

const config: Config = {
  allowedMethods: ["GET", "POST", "PUT"],
  version: "1.0",
};

// config.allowedMethods = [];          // ERROR — readonly property
// config.allowedMethods.push("DELETE"); // ERROR — readonly array
config.version = "2.0";                // OK — not readonly

// Readonly<T> utility type — makes all properties readonly
type ReadonlyUser = Readonly<User>;
// { readonly id: number; readonly name: string; ... }

// Deep readonly (TypeScript doesn't have this built-in — need utility)
type DeepReadonly<T> = {
  readonly [K in keyof T]: T[K] extends object ? DeepReadonly<T[K]> : T[K];
};
```

---

## 5.7 Index Signatures

Index signatures allow an interface to describe objects with dynamic keys:

```typescript
// Index signature: any key of type string maps to a value of type number
interface NumberMap {
  [key: string]: number;
}

const scores: NumberMap = {
  alice: 95,
  bob: 87,
  charlie: 92,
};

scores.dave = 88;  // OK — dynamic key
const aliceScore: number = scores.alice;  // TypeScript knows: number

// Combining with specific properties
interface UserRecord {
  id: number;
  name: string;
  [key: string]: string | number;  // additional dynamic properties
  // Note: id and name must match the index signature's value type!
}

// Numeric index signature
interface NumberedList {
  [index: number]: string;
  length: number;
}

const list: NumberedList = {
  0: "first",
  1: "second",
  2: "third",
  length: 3,
};

// Record type — cleaner for dynamic key-value pairs
type WordCount = Record<string, number>;
const counts: WordCount = { hello: 3, world: 5 };

// Typed dictionary
type ColorMap = Record<"red" | "green" | "blue", string>;
const colors: ColorMap = {
  red: "#FF0000",
  green: "#00FF00",
  blue: "#0000FF",
};
```

---

## 5.8 Excess Property Checking

TypeScript performs **excess property checking** when you create an object literal and assign it directly to a typed variable. This catches typos in property names.

```typescript
interface User {
  id: number;
  name: string;
}

// Direct assignment — excess property checking applies
const user: User = {
  id: 1,
  name: "Alice",
  // emai: "alice@example.com",  // ERROR: Object literal may only specify known properties,
  //                              // and 'emai' does not exist in type 'User'
  //                              // (you probably meant 'email')
};

// But through an intermediate variable — no excess property check!
const userData = {
  id: 1,
  name: "Alice",
  emai: "alice@example.com",  // Typo! But TypeScript doesn't catch it here
};
const user2: User = userData;  // OK — structural compatibility check, not excess check

// This is a key structural typing behavior:
// An object with MORE properties than required is assignable to a type requiring fewer
function printUser(user: { name: string }): void {
  console.log(user.name);
}

const detailedUser = { name: "Alice", age: 30, email: "alice@example.com" };
printUser(detailedUser);  // OK — has 'name', extra properties don't matter
```

---

## 5.9 Extending Interfaces

Interfaces can extend other interfaces to build more specific types:

```typescript
interface Entity {
  id: number;
  createdAt: Date;
  updatedAt: Date;
}

interface User extends Entity {
  name: string;
  email: string;
}

interface Admin extends User {
  permissions: string[];
  isSuperAdmin: boolean;
}

// Admin has: id, createdAt, updatedAt, name, email, permissions, isSuperAdmin
const admin: Admin = {
  id: 1,
  createdAt: new Date(),
  updatedAt: new Date(),
  name: "Alice",
  email: "alice@admin.com",
  permissions: ["read", "write", "delete"],
  isSuperAdmin: true,
};

// Multiple extension
interface Serializable {
  serialize(): string;
  deserialize(data: string): void;
}

interface Cacheable {
  cacheKey: string;
  ttl: number;
}

interface CacheableUser extends User, Serializable, Cacheable {
  displayName: string;
}
```

---

## 5.10 Structural Typing — Duck Typing in TypeScript

This is the foundation of TypeScript's type compatibility: two types are compatible if they have the same structure, regardless of their names.

### The Structural Compatibility Rule

```typescript
interface Point2D {
  x: number;
  y: number;
}

interface Coordinate {
  x: number;
  y: number;
}

// These are different interface names but the same structure
function plotPoint(p: Point2D): void {
  console.log(`(${p.x}, ${p.y})`);
}

const coord: Coordinate = { x: 3, y: 4 };
plotPoint(coord);  // OK! Coordinate is structurally compatible with Point2D
```

### Subtype Compatibility

```typescript
// A type with MORE properties is a subtype of a type with fewer
interface Named {
  name: string;
}

interface NamedAndAged {
  name: string;
  age: number;
}

function greet(entity: Named): void {
  console.log(`Hello, ${entity.name}!`);
}

const person: NamedAndAged = { name: "Alice", age: 30 };
greet(person);  // OK — NamedAndAged has everything Named requires, and more
```

### Structural Typing Surprises

```typescript
// Same structure = same type, even if semantically different
type Dollars = { amount: number; currency: string };
type Euros = { amount: number; currency: string };

function deposit(money: Dollars): void {
  // assumes money is in dollars
}

const euros: Euros = { amount: 100, currency: "EUR" };
deposit(euros);  // TypeScript: OK! Both have same shape.
// This is a valid TypeScript complaint — structural typing doesn't prevent semantic errors

// Solution: use branded types (covered in advanced types)
type Dollars2 = { amount: number; currency: "USD"; readonly __brand: "Dollars" };
type Euros2 = { amount: number; currency: "EUR"; readonly __brand: "Euros" };
// Now they're structurally different!
```

---

## Complete Example: A Type-Safe Data Model

```typescript
// types.ts — a complete data model using interfaces

interface Entity {
  readonly id: string;
  readonly createdAt: Date;
  updatedAt: Date;
}

interface User extends Entity {
  name: string;
  email: string;
  role: "admin" | "editor" | "viewer";
  preferences: UserPreferences;
}

interface UserPreferences {
  theme: "light" | "dark";
  language: string;
  notifications: NotificationSettings;
}

interface NotificationSettings {
  email: boolean;
  push: boolean;
  frequency: "immediate" | "daily" | "weekly";
}

interface Post extends Entity {
  title: string;
  content: string;
  author: User;
  tags: readonly string[];
  status: "draft" | "published" | "archived";
  publishedAt?: Date;
  viewCount: number;
}

// service.ts — functions using the type model
type PostFilter = Partial<Pick<Post, "status" | "tags">>;
type CreatePostInput = Omit<Post, "id" | "createdAt" | "updatedAt" | "viewCount" | "author"> & {
  authorId: string;
};

function createPost(input: CreatePostInput, author: User): Post {
  const now = new Date();
  return {
    id: generateId(),
    createdAt: now,
    updatedAt: now,
    title: input.title,
    content: input.content,
    author,
    tags: input.tags,
    status: input.status,
    publishedAt: input.status === "published" ? now : input.publishedAt,
    viewCount: 0,
  };
}

function filterPosts(posts: Post[], filter: PostFilter): Post[] {
  return posts.filter((post) => {
    if (filter.status && post.status !== filter.status) return false;
    if (filter.tags && !filter.tags.some((tag) => post.tags.includes(tag))) return false;
    return true;
  });
}

function generateId(): string {
  return Math.random().toString(36).slice(2);
}
```

---

## Summary

Objects in TypeScript are described by their *shape* — the names and types of their properties. This shape can be expressed as an inline object type literal, named with an `interface`, or aliased with `type`. Interfaces support optional properties (`?`), readonly properties (`readonly`), index signatures for dynamic keys, and extension with `extends`. Type aliases support unions, intersections, and all type expressions. Structural typing means compatibility is based on shape, not name — enabling flexible, ergonomic code at the cost of some semantic safety.

---

## Key Takeaways

- **Interface vs type alias**: interfaces support declaration merging and are preferred for object shapes; type aliases handle unions, tuples, and any type
- **`readonly`** prevents reassignment after creation — useful for configuration objects and data models
- **Optional properties** (`?`) represent data that may be absent — always check before using
- **Index signatures** describe dynamic key-value objects — use `Record<K, V>` for simpler cases
- **Excess property checking** only applies to object literals assigned directly — not intermediate variables
- **Structural typing**: TypeScript checks *shape*, not *name* — a compatible shape is compatible, regardless of the interface name

---

## Practice Questions

1. What is the difference between `interface` and `type` for object shapes?
2. Can an interface extend a type alias? Can a type alias extend an interface?
3. What is declaration merging, and which construct supports it?
4. What is excess property checking, and when does it NOT apply?
5. Why is TypeScript's structural typing sometimes called "duck typing"?
6. What is the difference between `readonly` on a property and `const` for a variable?

---

## Exercises

**Exercise 1**: Design a complete type system for a shopping cart application. Include interfaces for: `Product`, `CartItem`, `Cart`, `Discount`, and `Order`. Use optional properties, readonly properties, and proper relationships.

**Exercise 2**: Create an interface `Repository<T extends Entity>` with methods: `findById`, `findAll`, `create`, `update`, and `delete`. Then implement it for `UserRepository`.

**Exercise 3**: Write a function `mergeDefaults<T>(partial: Partial<T>, defaults: T): T` that fills in missing properties from defaults. Make it fully type-safe.

**Exercise 4**: Demonstrate the difference between excess property checking for direct literals vs. intermediate variables. Write a case where a bug (typo in property name) passes type checking due to structural compatibility.

---

*Next: [Chapter 6 — Advanced Type System](06-advanced-types.md)*
