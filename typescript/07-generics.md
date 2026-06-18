# Chapter 7 — Generics

> *"Generics let you write code that works with any type, while still being fully type-safe. They are the mechanism by which TypeScript achieves reusability without sacrificing correctness."*

---

## 7.1 The Problem Generics Solve

Without generics, you face a choice: be specific (lose reusability) or use `any` (lose type safety).

```typescript
// Specific but not reusable
function firstNumber(arr: number[]): number {
  return arr[0];
}

function firstString(arr: string[]): string {
  return arr[0];
}

// Reusable but loses type safety
function firstAny(arr: any[]): any {
  return arr[0];
}

const n = firstAny([1, 2, 3]);
n.toUpperCase();  // No error here! But crashes at runtime.

// With generics: reusable AND type-safe
function first<T>(arr: T[]): T {
  return arr[0];
}

const num = first([1, 2, 3]);      // TypeScript infers: T = number
num.toFixed(2);                    // OK — TypeScript knows it's a number
// num.toUpperCase();              // ERROR — number has no toUpperCase

const str = first(["a", "b"]);    // TypeScript infers: T = string
str.toUpperCase();                 // OK — TypeScript knows it's a string
```

---

## 7.2 Generic Functions

### Basic Syntax

```typescript
// <T> declares a type parameter named T
// T is a placeholder that gets filled in at call time
function identity<T>(value: T): T {
  return value;
}

// TypeScript infers T from the argument
const n = identity(42);            // T = number, returns number
const s = identity("hello");       // T = string, returns string
const b = identity(true);          // T = boolean, returns boolean

// You can also provide T explicitly (rarely needed)
const explicit = identity<string>("hello");
```

### Multiple Type Parameters

```typescript
// Multiple type parameters
function pair<A, B>(first: A, second: B): [A, B] {
  return [first, second];
}

const p1 = pair(1, "hello");       // [number, string]
const p2 = pair(true, { x: 1 });   // [boolean, { x: number }]

// Swap: two parameters, swapped return
function swap<A, B>(pair: [A, B]): [B, A] {
  return [pair[1], pair[0]];
}

const swapped = swap([1, "hello"]);  // [string, number]
console.log(swapped[0].toUpperCase());  // TypeScript knows index 0 is string
```

### Generic Array Operations

```typescript
// A fully typed map — like Array.prototype.map but standalone
function map<T, U>(arr: T[], fn: (item: T, index: number) => U): U[] {
  return arr.map(fn);
}

const doubled = map([1, 2, 3], (n) => n * 2);        // number[]
const asStrings = map([1, 2, 3], (n) => `${n}`);     // string[]
const lengths = map(["hi", "hello"], (s) => s.length); // number[]

// filter
function filter<T>(arr: T[], predicate: (item: T) => boolean): T[] {
  return arr.filter(predicate);
}

const evens = filter([1, 2, 3, 4], (n) => n % 2 === 0);  // number[]

// reduce — input and output can be different types
function reduce<T, U>(
  arr: T[],
  fn: (accumulator: U, item: T, index: number) => U,
  initial: U
): U {
  return arr.reduce(fn, initial);
}

const sum = reduce([1, 2, 3], (acc, n) => acc + n, 0);           // number
const joined = reduce(["a", "b", "c"], (acc, s) => acc + s, ""); // string
const wordLengths = reduce(
  ["hello", "world"],
  (acc, s) => ({ ...acc, [s]: s.length }),
  {} as Record<string, number>
);  // { hello: 5, world: 5 }
```

---

## 7.3 Generic Interfaces

```typescript
// Generic interface
interface Container<T> {
  value: T;
  map<U>(fn: (val: T) => U): Container<U>;
  filter(predicate: (val: T) => boolean): Container<T | undefined>;
}

// A simple Box
interface Box<T> {
  readonly value: T;
  unwrap(): T;
}

function makeBox<T>(value: T): Box<T> {
  return {
    value,
    unwrap() { return this.value; },
  };
}

const numBox: Box<number> = makeBox(42);
const strBox: Box<string> = makeBox("hello");
console.log(numBox.unwrap().toFixed(2));   // "42.00"
console.log(strBox.unwrap().toUpperCase()); // "HELLO"

// Generic interface for a key-value store
interface KeyValueStore<K extends string | number, V> {
  get(key: K): V | undefined;
  set(key: K, value: V): void;
  has(key: K): boolean;
  delete(key: K): boolean;
  keys(): K[];
  values(): V[];
  entries(): [K, V][];
}

class MapStore<K extends string | number, V> implements KeyValueStore<K, V> {
  private store = new Map<K, V>();

  get(key: K): V | undefined { return this.store.get(key); }
  set(key: K, value: V): void { this.store.set(key, value); }
  has(key: K): boolean { return this.store.has(key); }
  delete(key: K): boolean { return this.store.delete(key); }
  keys(): K[] { return [...this.store.keys()]; }
  values(): V[] { return [...this.store.values()]; }
  entries(): [K, V][] { return [...this.store.entries()]; }
}

const store = new MapStore<string, number>();
store.set("a", 1);
store.set("b", 2);
console.log(store.get("a"));  // 1 — typed as number | undefined
```

---

## 7.4 Generic Type Aliases

```typescript
// Generic type alias
type Maybe<T> = T | null | undefined;
type Result<T, E = Error> = { success: true; data: T } | { success: false; error: E };
type Nullable<T> = T | null;
type Optional<T> = T | undefined;

// Using Maybe
function findUser(id: number): Maybe<User> {
  return id === 1 ? { id: 1, name: "Alice", email: "alice@ex.com" } : null;
}

// Using Result
function parseNumber(s: string): Result<number, string> {
  const n = Number(s);
  if (isNaN(n)) {
    return { success: false, error: `"${s}" is not a valid number` };
  }
  return { success: true, data: n };
}

const result = parseNumber("42");
if (result.success) {
  console.log(result.data.toFixed(2));   // data: number
} else {
  console.error(result.error.toUpperCase()); // error: string
}

// Generic pair/tuple aliases
type Pair<A, B = A> = [A, B];  // default B = A
type Triple<T> = [T, T, T];

const coords: Pair<number> = [3, 4];
const rgb: Triple<number> = [255, 128, 0];
const mixed: Pair<string, number> = ["age", 30];

// Generic function type
type Transformer<T, U = T> = (value: T) => U;
type Predicate<T> = (value: T) => boolean;
type Comparator<T> = (a: T, b: T) => number;

const isEven: Predicate<number> = (n) => n % 2 === 0;
const compareNumbers: Comparator<number> = (a, b) => a - b;
```

---

## 7.5 Generic Classes

```typescript
// Generic stack
class Stack<T> {
  private items: T[] = [];

  push(item: T): void {
    this.items.push(item);
  }

  pop(): T | undefined {
    return this.items.pop();
  }

  peek(): T | undefined {
    return this.items[this.items.length - 1];
  }

  isEmpty(): boolean {
    return this.items.length === 0;
  }

  get size(): number {
    return this.items.length;
  }
}

const numStack = new Stack<number>();
numStack.push(1);
numStack.push(2);
numStack.push(3);
console.log(numStack.pop());   // 3, typed as number | undefined
console.log(numStack.peek());  // 2

const strStack = new Stack<string>();
strStack.push("hello");
// strStack.push(42);  // ERROR — Stack<string> only accepts strings

// Generic queue
class Queue<T> {
  private items: T[] = [];

  enqueue(item: T): void {
    this.items.push(item);
  }

  dequeue(): T | undefined {
    return this.items.shift();
  }

  peek(): T | undefined {
    return this.items[0];
  }

  get size(): number {
    return this.items.length;
  }
}

// Generic event emitter
class TypedEventEmitter<Events extends Record<string, unknown[]>> {
  private listeners: {
    [K in keyof Events]?: Array<(...args: Events[K]) => void>;
  } = {};

  on<K extends keyof Events>(
    event: K,
    listener: (...args: Events[K]) => void
  ): this {
    if (!this.listeners[event]) {
      this.listeners[event] = [];
    }
    this.listeners[event]!.push(listener);
    return this;
  }

  emit<K extends keyof Events>(event: K, ...args: Events[K]): void {
    this.listeners[event]?.forEach((listener) => listener(...args));
  }
}

// Define the event map
type AppEvents = {
  connect: [userId: string];
  disconnect: [userId: string, reason: string];
  message: [from: string, to: string, content: string];
};

const emitter = new TypedEventEmitter<AppEvents>();

emitter.on("connect", (userId) => {
  // userId: string — TypeScript knows
  console.log(`${userId} connected`);
});

emitter.on("message", (from, to, content) => {
  // All parameters are correctly typed strings
  console.log(`${from} → ${to}: ${content}`);
});

emitter.emit("connect", "user-123");
emitter.emit("message", "alice", "bob", "Hello!");
// emitter.emit("connect", 123);  // ERROR — userId must be string
```

---

## 7.6 Generic Constraints with extends

Without constraints, TypeScript doesn't know what operations are valid on a type parameter. `extends` adds constraints:

```typescript
// Without constraint: T could be anything — can't access any properties
function getLength<T>(value: T): number {
  // return value.length;  // ERROR: T might not have 'length'
  return 0;  // not useful
}

// With constraint: T must have a 'length' property
function getLength2<T extends { length: number }>(value: T): number {
  return value.length;  // OK — constraint guarantees .length exists
}

getLength2("hello");   // 5 — string has length
getLength2([1, 2, 3]); // 3 — array has length
getLength2({ length: 10, name: "test" });  // 10 — any object with length

// Real use: generic function that only works on objects with certain properties
interface HasId {
  id: number;
}

function findById<T extends HasId>(items: T[], id: number): T | undefined {
  return items.find((item) => item.id === id);
}

const users = [
  { id: 1, name: "Alice", role: "admin" },
  { id: 2, name: "Bob", role: "user" },
];

const alice = findById(users, 1);
// alice: { id: number; name: string; role: string } | undefined
// TypeScript preserves the full type, not just HasId!
console.log(alice?.name);  // "Alice"
console.log(alice?.role);  // "admin"

// Extending a union
function formatValue<T extends string | number>(value: T): string {
  if (typeof value === "string") return value.toUpperCase();
  return value.toFixed(2);
}

formatValue("hello");  // "HELLO"
formatValue(3.14);     // "3.14"
// formatValue(true);  // ERROR — boolean not in string | number
```

---

## 7.7 keyof and Generic Constraints

`keyof` extracts the keys of a type. Combined with generics, it enables type-safe property access:

```typescript
// keyof T: the union of keys of type T
type UserKeys = keyof { name: string; age: number; email: string };
// "name" | "age" | "email"

// Generic function constrained by keyof
function getProperty<T, K extends keyof T>(obj: T, key: K): T[K] {
  return obj[key];
}

const user = { name: "Alice", age: 30, email: "alice@example.com" };

const name = getProperty(user, "name");   // type: string
const age = getProperty(user, "age");     // type: number
// getProperty(user, "phone");            // ERROR: "phone" not in keyof user

// setProperty
function setProperty<T, K extends keyof T>(obj: T, key: K, value: T[K]): void {
  obj[key] = value;
}

setProperty(user, "age", 31);        // OK — 31 is number, matches T[K]
// setProperty(user, "age", "old");  // ERROR — "old" is string, not number

// Picking properties from an object
function pick<T, K extends keyof T>(obj: T, keys: K[]): Pick<T, K> {
  const result = {} as Pick<T, K>;
  keys.forEach((key) => {
    result[key] = obj[key];
  });
  return result;
}

const partial = pick(user, ["name", "email"]);
// partial: { name: string; email: string }
console.log(partial.name);   // "Alice"
console.log(partial.email);  // "alice@example.com"
// console.log(partial.age); // ERROR — age was not picked
```

---

## 7.8 Default Type Parameters

Type parameters can have defaults:

```typescript
// Default type parameter
interface ApiResponse<T = unknown> {
  data: T;
  status: number;
  message: string;
}

// Without specifying T — uses default (unknown)
function fetchData(url: string): Promise<ApiResponse> {
  return fetch(url).then((r) => r.json());
}

// With explicit T
function fetchUsers(url: string): Promise<ApiResponse<User[]>> {
  return fetch(url).then((r) => r.json());
}

// Generic with multiple defaults
type PaginatedResponse<T, Meta = { page: number; total: number }> = {
  items: T[];
  meta: Meta;
};

// Uses default Meta
const response1: PaginatedResponse<User> = {
  items: [],
  meta: { page: 1, total: 100 },
};

// Custom Meta
const response2: PaginatedResponse<User, { cursor: string; hasMore: boolean }> = {
  items: [],
  meta: { cursor: "abc", hasMore: true },
};
```

---

## 7.9 Conditional Types with Generics (Introduction)

Generics become even more powerful when combined with conditional types (covered in depth in Chapter 11):

```typescript
// A function that returns different types based on input
type Unwrapped<T> = T extends Promise<infer U> ? U : T;

// If T is a Promise<U>, return U; otherwise return T
type A = Unwrapped<Promise<string>>;  // string
type B = Unwrapped<number>;           // number
type C = Unwrapped<Promise<User[]>>;  // User[]

// NonNullable — removes null and undefined
type NonNull<T> = T extends null | undefined ? never : T;
type D = NonNull<string | null | undefined>;  // string
```

---

## 7.10 Real-World Generic Patterns

### A Generic Repository Pattern

```typescript
interface Entity {
  id: string;
}

interface Repository<T extends Entity> {
  findById(id: string): Promise<T | null>;
  findAll(filter?: Partial<T>): Promise<T[]>;
  create(data: Omit<T, "id">): Promise<T>;
  update(id: string, data: Partial<Omit<T, "id">>): Promise<T | null>;
  delete(id: string): Promise<boolean>;
}

// A concrete User type
interface User extends Entity {
  name: string;
  email: string;
  role: "admin" | "user";
}

// In-memory implementation
class InMemoryRepository<T extends Entity> implements Repository<T> {
  protected items: T[] = [];
  private nextId = 1;

  async findById(id: string): Promise<T | null> {
    return this.items.find((item) => item.id === id) ?? null;
  }

  async findAll(filter?: Partial<T>): Promise<T[]> {
    if (!filter) return this.items;
    return this.items.filter((item) =>
      Object.entries(filter).every(
        ([key, value]) => item[key as keyof T] === value
      )
    );
  }

  async create(data: Omit<T, "id">): Promise<T> {
    const item = { ...data, id: String(this.nextId++) } as T;
    this.items.push(item);
    return item;
  }

  async update(id: string, data: Partial<Omit<T, "id">>): Promise<T | null> {
    const index = this.items.findIndex((item) => item.id === id);
    if (index === -1) return null;
    this.items[index] = { ...this.items[index], ...data };
    return this.items[index];
  }

  async delete(id: string): Promise<boolean> {
    const index = this.items.findIndex((item) => item.id === id);
    if (index === -1) return false;
    this.items.splice(index, 1);
    return true;
  }
}

// Type-safe usage
const userRepo = new InMemoryRepository<User>();

async function main() {
  const alice = await userRepo.create({ name: "Alice", email: "alice@ex.com", role: "admin" });
  const bob = await userRepo.create({ name: "Bob", email: "bob@ex.com", role: "user" });

  const admins = await userRepo.findAll({ role: "admin" });
  console.log(admins[0].name);  // TypeScript knows: name is string

  await userRepo.update(alice.id, { name: "Alice Smith" });
  // await userRepo.update(alice.id, { role: "superuser" });  // ERROR — not a valid role
}
```

### A Generic Observable

```typescript
type Observer<T> = {
  next: (value: T) => void;
  error?: (err: Error) => void;
  complete?: () => void;
};

class Observable<T> {
  constructor(
    private producer: (observer: Observer<T>) => (() => void) | void
  ) {}

  subscribe(observer: Observer<T>): () => void {
    const cleanup = this.producer(observer);
    return cleanup ?? (() => {});
  }

  map<U>(fn: (value: T) => U): Observable<U> {
    return new Observable<U>((observer) => {
      return this.subscribe({
        next: (value) => observer.next(fn(value)),
        error: observer.error,
        complete: observer.complete,
      });
    });
  }

  filter(predicate: (value: T) => boolean): Observable<T> {
    return new Observable<T>((observer) => {
      return this.subscribe({
        next: (value) => { if (predicate(value)) observer.next(value); },
        error: observer.error,
        complete: observer.complete,
      });
    });
  }
}

// Usage
const numbers$ = new Observable<number>((observer) => {
  [1, 2, 3, 4, 5].forEach((n) => observer.next(n));
  observer.complete?.();
});

const evenDoubles$ = numbers$
  .filter((n) => n % 2 === 0)   // Observable<number>
  .map((n) => n * 2);           // Observable<number>

evenDoubles$.subscribe({
  next: (n) => console.log(n),  // n: number — correctly typed
  complete: () => console.log("done"),
});
// 4, 8, done
```

---

## Summary

Generics enable **type-safe reusability**: write a function, interface, or class once and use it with any type. Type parameters (like `T`) are placeholders filled in at use time — either inferred by TypeScript or explicitly provided. Constraints with `extends` restrict what types are valid for a parameter, enabling access to specific properties. `keyof` combined with generics enables type-safe property access. Default type parameters (`T = SomeType`) make generic code ergonomic to use without always specifying types explicitly.

---

## Key Takeaways

- **Generics** are preferable to `any` — they reuse code without losing type information
- **Type inference** usually makes explicit generic arguments unnecessary (`first([1,2,3])` not `first<number>([1,2,3])`)
- **Multiple type parameters** (`<T, U>`) enable rich mappings between input and output types
- **`extends` constraints** restrict valid types for a parameter — access only what the constraint guarantees
- **`keyof T`** produces a union of the type's keys — combine with generics for type-safe property access
- **Default type parameters** make generic APIs ergonomic while remaining flexible
- Generic patterns (Repository, Observable, Result) are common in TypeScript codebases

---

## Practice Questions

1. What is the difference between `Array<T>` and `T[]`? Are they the same thing?
2. Why is a generic function with `T extends { length: number }` more useful than one with `T`?
3. What does `keyof T` produce? When would you use it with generics?
4. What is the difference between `<T = string>` and `<T extends string>`?
5. In `function getProperty<T, K extends keyof T>(obj: T, key: K): T[K]`, what does `T[K]` mean?

---

## Exercises

**Exercise 1**: Implement a generic `groupBy<T, K extends string>(items: T[], keyFn: (item: T) => K): Record<K, T[]>` function that groups an array of items by a key.

**Exercise 2**: Write a generic `zip<A, B>(a: A[], b: B[]): [A, B][]` function that pairs elements from two arrays.

**Exercise 3**: Implement a generic `LRUCache<K, V>` class with `get(key: K): V | undefined`, `set(key: K, value: V): void`, and a `capacity` constructor argument.

**Exercise 4**: Write a generic `retry<T>(fn: () => Promise<T>, attempts: number, delay: number): Promise<T>` that retries a failing async function up to `attempts` times.

---

*Next: [Chapter 8 — Utility Types](08-utility-types.md)*
