# Chapter 91 — Common Pitfalls

> *"TypeScript's type system is powerful, but it has gaps — places where it trusts you too much, where runtime behavior differs from compile-time expectations, or where common patterns produce subtle bugs."*

---

## Pitfall 1: any Spreads Like a Disease

```typescript
// The problem: any contaminates everything it touches
function parseConfig(json: string): any {  // any return type
  return JSON.parse(json);
}

const config = parseConfig('{"host":"localhost"}');
// config: any — TypeScript gives up on type checking here

const host = config.host;       // host: any — still infected
const port = config.port + 1;  // port: any — no error, even if port is undefined
const upper = host.toUpperCase(); // any — no autocomplete, might crash

// The fix: use unknown and validate
function parseConfig2(json: string): unknown {
  return JSON.parse(json);
}

const config2 = parseConfig2('{"host":"localhost"}');
// config2: unknown — must narrow before use
if (
  typeof config2 === "object" &&
  config2 !== null &&
  "host" in config2 &&
  typeof (config2 as { host: unknown }).host === "string"
) {
  const host = (config2 as { host: string }).host;  // string — safe
}
```

---

## Pitfall 2: Mutating Readonly At Runtime

```typescript
// TypeScript's 'readonly' is compile-time only — no runtime enforcement
interface Config {
  readonly host: string;
}

const config: Config = { host: "localhost" };
// config.host = "changed";  // TypeScript ERROR

// But...
const mutable = config as { host: string };  // cast away readonly
mutable.host = "changed";  // works at runtime!

// Or via Object.assign, spread:
Object.assign(config, { host: "changed" });  // TypeScript: no error! Any-typed operation

// Or via JSON parse:
const fresh: Config = JSON.parse(JSON.stringify(config));  // fresh: Config — but mutable

// The lesson: readonly is a TypeScript convention — it doesn't survive casts
// For true immutability, use Object.freeze()
const frozenConfig = Object.freeze({ host: "localhost" });
// frozenConfig.host = "changed";  // runtime TypeError in strict mode
```

---

## Pitfall 3: Type Assertions Hide Runtime Errors

```typescript
// BAD: using 'as' to silence TypeScript
async function getUser(id: string): Promise<User> {
  const response = await fetch(`/api/users/${id}`);
  return response.json() as User;  // TypeScript: OK. Runtime: ???
  // If the API returns { error: "not found" }, TypeScript thinks it's a User
}

// Later...
const user = await getUser("missing-id");
console.log(user.name.toUpperCase());  // TypeError: Cannot read property of undefined

// FIX: validate the response
async function getUser2(id: string): Promise<User> {
  const response = await fetch(`/api/users/${id}`);
  const data: unknown = await response.json();
  if (!isUser(data)) throw new Error(`Invalid user response for id ${id}`);
  return data;  // TypeScript: User. Runtime: validated.
}

// When is 'as' safe?
// 1. You have more information than TypeScript (e.g., DOM APIs)
const input = document.getElementById("email") as HTMLInputElement;
// You know this element is an input — TypeScript doesn't

// 2. Narrowing within a known union
type Fruit = "apple" | "banana" | "cherry";
const fruits: string[] = ["apple", "banana", "cherry"];
const fruit = fruits[0] as Fruit;  // You know the strings are valid Fruit values
```

---

## Pitfall 4: The typeof null === "object" Trap

```typescript
// JavaScript legacy: typeof null === "object"
function isObject(value: unknown): value is object {
  return typeof value === "object";  // WRONG: passes null through!
}

const result = isObject(null);  // true — but null is not a useful object

// Correct check:
function isObject2(value: unknown): value is object {
  return typeof value === "object" && value !== null;
}

// TypeScript narrows correctly with a null check:
function process(value: unknown): void {
  if (typeof value === "object") {
    // value: object | null — TypeScript knows null is still possible!
    value.toString();  // ERROR — value might be null

    // Must additionally check:
    if (value !== null) {
      value.toString();  // OK
    }
  }
}
```

---

## Pitfall 5: Optional Properties vs undefined Values

```typescript
// Optional property: may be absent
interface A {
  x?: number;  // x might not be present at all
}

// Property with undefined type: present, but value might be undefined
interface B {
  x: number | undefined;  // x is always present, but value might be undefined
}

// With exactOptionalPropertyTypes: true, these are different!
const a: A = {};          // OK — x is absent
const a2: A = { x: undefined };  // ERROR with exactOptionalPropertyTypes
const b: B = { x: undefined };   // OK — x is present but undefined
// const b2: B = {};              // ERROR — x must be present

// The bug this prevents:
function merge<T extends object>(target: T, source: Partial<T>): T {
  return { ...target, ...source };  // spreads 'undefined' over defined values!
}

const merged = merge({ x: 1, y: 2 }, { x: undefined });
// result: { x: undefined, y: 2 } — x was overwritten with undefined!
// exactOptionalPropertyTypes would catch this
```

---

## Pitfall 6: Structurally Identical Types Are Interchangeable

```typescript
// TypeScript's structural typing: same shape = same type
// This can cause semantic bugs

type Meters = number;
type Seconds = number;
type Speed = number;

function calculateSpeed(distance: Meters, time: Seconds): Speed {
  return distance / time;
}

const distance: Meters = 100;
const time: Seconds = 10;
const wrongOrder = calculateSpeed(time, distance);  // TypeScript: OK!
// Runtime: 0.1 instead of 10 — wrong order, no error

// Fix: branded types prevent this
type Meters2 = number & { readonly __brand: "Meters" };
type Seconds2 = number & { readonly __brand: "Seconds" };

function meters(n: number): Meters2 { return n as Meters2; }
function seconds(n: number): Seconds2 { return n as Seconds2; }

function calculateSpeed2(distance: Meters2, time: Seconds2): number {
  return distance / time;
}

const d = meters(100);
const t = seconds(10);
// calculateSpeed2(t, d);  // ERROR — Seconds2 is not Meters2
calculateSpeed2(d, t);    // OK
```

---

## Pitfall 7: Non-Null Assertion (!) Without Safety

```typescript
// The ! operator: "trust me, this is not null/undefined"
// It's a promise to TypeScript — if wrong, runtime crash

function getElement(): HTMLElement {
  return document.getElementById("app")!;  // ! says: I know it's not null
}

// What happens when it IS null?
function render(): void {
  const app = document.getElementById("app")!;
  app.innerHTML = "<h1>Hello</h1>";  // TypeError if "app" doesn't exist
}

// Common misuse: in loops where the element might not exist
const inputs = document.querySelectorAll("input");
inputs.forEach((input) => {
  const value = (input as HTMLInputElement).value!;  // ! on a string — useless!
  // string never has undefined via ! — but it can be ""
});

// Better approach:
function getElement2(): HTMLElement {
  const el = document.getElementById("app");
  if (!el) throw new Error("Element #app not found");
  return el;  // TypeScript: HTMLElement
}

// Only use ! when:
// 1. You've checked elsewhere (just outside TypeScript's view)
// 2. The element is guaranteed by the framework/context
// 3. You're in test code with known data
```

---

## Pitfall 8: Async Error Handling with unknown

```typescript
// With strictNullChecks, catch error is 'unknown'
// BAD — treating error as Error without checking
async function fetch_bad(url: string): Promise<string> {
  try {
    const response = await fetch(url);
    return await response.text();
  } catch (error) {
    // @ts-expect-error or error is treated as unknown
    console.error(error.message);  // ERROR: error is unknown
    throw error;
  }
}

// GOOD — narrow the error type
async function fetch_good(url: string): Promise<string> {
  try {
    const response = await fetch(url);
    return await response.text();
  } catch (error: unknown) {
    if (error instanceof Error) {
      console.error(`Failed: ${error.message}`);
      throw error;
    }
    throw new Error(`Unknown error: ${String(error)}`);
  }
}

// Helper to extract error message safely
function getErrorMessage(error: unknown): string {
  if (error instanceof Error) return error.message;
  if (typeof error === "string") return error;
  return "An unknown error occurred";
}
```

---

## Pitfall 9: Mutating Objects Through Union Types

```typescript
// TypeScript narrows in place, but mutation can break it
interface Square { kind: "square"; size: number }
interface Circle { kind: "circle"; radius: number }
type Shape = Square | Circle;

function process(shape: Shape): void {
  if (shape.kind === "square") {
    // shape: Square
    (shape as unknown as Circle).kind = "circle";  // mutating! (bypasses TypeScript)
    // Now shape.kind is "circle" but TypeScript thinks it's a Square
    console.log(shape.size);  // TypeScript: OK. Runtime: undefined crash
  }
}

// The lesson: TypeScript's narrowing is based on the code path, not live values
// Mutating a discriminant property after narrowing causes type/value divergence
// Solution: never mutate discriminant properties; keep shape objects immutable
```

---

## Pitfall 10: Function Parameter Bivariance (Method Shorthand)

```typescript
// Interface with method shorthand — bivariant (less safe)
interface Handler1 {
  handle(event: MouseEvent): void;  // method shorthand
}

// Interface with function property — contravariant (safer with strictFunctionTypes)
interface Handler2 {
  handle: (event: MouseEvent) => void;  // function property
}

// The difference:
type GenericHandler = (event: Event) => void;

// Method shorthand — TypeScript accepts both covariant and contravariant
const h1: Handler1 = { handle: (event: Event) => {} };  // OK (less specific param)

// Function property — stricter (contravariant with strictFunctionTypes)
// const h2: Handler2 = { handle: (event: Event) => {} };  // ERROR with strictFunctionTypes
// A handler for Event is NOT assignable to a handler expecting MouseEvent
// because MouseEvent has properties that Event might not

// Use function property syntax in interfaces for stricter type checking
```

---

## Pitfall 11: Class Properties vs Interface Properties in Constructor

```typescript
// BAD: forgetting to declare class property before use
class User {
  // name: string;  // missing!

  constructor(name: string) {
    this.name = name;  // TypeScript ERROR with strict: property not declared
  }
}

// TypeScript's strictPropertyInitialization catches this
class User2 {
  name: string;  // must declare

  constructor(name: string) {
    this.name = name;  // OK
  }
}

// Using definite assignment assertion when you initialize elsewhere
class Database {
  connection!: Connection;  // ! = "I'll set this before use, trust me"

  async init(): Promise<void> {
    this.connection = await createConnection();
  }
}
// If you call db.query() before db.init(), you'll get a runtime error
// Use ! only when you're certain of the initialization order
```

---

## Pitfall 12: Widening in Array Literals

```typescript
// Array literals are typed as mutable arrays — not tuples
const pair = [1, "hello"];  // (number | string)[] — NOT [number, string]

function swap([a, b]: [number, string]): [string, number] {
  return [b, a];
}

// swap(pair);  // ERROR — (number | string)[] is not [number, string]

// Fix 1: explicit tuple type
const pair2: [number, string] = [1, "hello"];
swap(pair2);  // OK

// Fix 2: as const
const pair3 = [1, "hello"] as const;  // readonly [1, "hello"]
// swap(pair3);  // Still ERROR — readonly vs mutable
// swap([...pair3]);  // OK — spread creates new (non-readonly) tuple

// Fix 3: annotation
function makePair(n: number, s: string): [number, string] {
  return [n, s];  // TypeScript infers [number, string] from return type
}
swap(makePair(1, "hello"));  // OK
```

---

## Pitfall 13: Promise Not Awaited

```typescript
// TypeScript doesn't always warn about un-awaited Promises
async function saveUser(user: User): Promise<void> {
  await db.save(user);
}

function processAndSave(user: User): void {
  saveUser(user);  // TypeScript: no error! But the save is not awaited
  // If saveUser throws, the error is silently swallowed
  console.log("saved");  // runs before saveUser completes!
}

// Enable: "no-floating-promises" ESLint rule (not a TypeScript flag)
// OR make the function async and await:
async function processAndSave2(user: User): Promise<void> {
  await saveUser(user);  // proper await
  console.log("saved");
}

// TypeScript does warn with void operator for explicit discard:
void saveUser(user);  // intentionally fire-and-forget
```

---

## Pitfall 14: Circular Type Aliases (Infinite Types)

```typescript
// TypeScript allows recursive types up to a limit
type JSON = string | number | boolean | null | JSON[] | { [key: string]: JSON };

// But circular non-recursive types cause TypeScript to error
// type A = B;   // ERROR if this creates a circular reference
// type B = A;

// Interfaces handle recursion better than type aliases in some cases
interface TreeNode {
  value: string;
  children: TreeNode[];  // OK — interfaces can reference themselves
}

// Use interfaces for recursive data structures
// Use type aliases for recursive mapped/conditional types with care
```

---

## Pitfall 15: Losing Types Through JSON.parse

```typescript
// JSON.parse returns 'any' — types are lost
const json = '{"id":1,"name":"Alice"}';
const parsed = JSON.parse(json);  // any — dangerous!

// parsed.id, parsed.name, parsed.anything — all 'any' with no errors

// FIX 1: annotate with 'unknown' and validate
const parsed2: unknown = JSON.parse(json);
// Must narrow before use

// FIX 2: use a typed wrapper
function typedParse<T>(json: string): unknown {
  return JSON.parse(json);
}
// Still unknown — but at least you get a reminder to validate

// FIX 3: use a validation library (zod, yup, io-ts)
// import { z } from "zod";
// const UserSchema = z.object({ id: z.number(), name: z.string() });
// const user = UserSchema.parse(JSON.parse(json));
// user is typed correctly AND validated at runtime
```

---

## Pitfall Summary — Quick Reference

| Pitfall | Root Cause | Fix |
|---------|-----------|-----|
| `any` spreads | `any` bypasses type checking | Use `unknown`, validate boundaries |
| `readonly` not runtime | `as` casts remove readonly | Use `Object.freeze()` |
| Type assertions crash | `as T` is a promise, not a check | Validate with type guards |
| `typeof null === "object"` | JavaScript legacy | Always add `&& value !== null` |
| Optional vs undefined | Subtle difference in `exactOptionalPropertyTypes` | Enable `exactOptionalPropertyTypes` |
| Structural type bugs | Same-shape types are interchangeable | Use branded types |
| `!` operator crashes | Promise without validation | Use type guards + if checks |
| async catch `unknown` | `catch` error is `unknown` | `instanceof Error` check |
| Un-awaited Promises | TypeScript doesn't always warn | ESLint `no-floating-promises` |
| Array widening | `[a,b]` is array, not tuple | `as const` or explicit tuple type |

---

*Next: [Chapter 92 — Interview Preparation](92-interview-prep.md)*
