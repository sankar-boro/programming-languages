# Chapter 3 — Basic Types and Variables

> *"The type system is a lens through which you see your data."*

---

## 3.1 let, const, var — and What TypeScript Changes

TypeScript uses the same variable declarations as JavaScript — `let`, `const`, and `var` — with the same scoping semantics. TypeScript adds type annotations and tightens usage rules.

### var — Function Scoped (Avoid)

```typescript
// var is function-scoped, not block-scoped
function demonstrateVar(): void {
  if (true) {
    var x = 10;  // scoped to the function, not the if block
  }
  console.log(x);  // 10 — accessible outside the if block
}

// TypeScript allows var but you should avoid it
// use let or const instead
```

### let — Block Scoped, Mutable

```typescript
let count: number = 0;
count = 1;    // OK — let is mutable

let name: string = "Alice";
// name = 42;  // ERROR: Type 'number' is not assignable to type 'string'

// Block scope
{
  let blockScoped = "only here";
  console.log(blockScoped);  // OK
}
// console.log(blockScoped);  // ERROR: Cannot find name 'blockScoped'
```

### const — Block Scoped, Immutable Binding

```typescript
const PI: number = 3.14159;
// PI = 3;  // ERROR: Cannot assign to 'PI' because it is a constant

// Important: const prevents reassignment, NOT mutation
const arr: number[] = [1, 2, 3];
arr.push(4);   // OK — we're mutating the array, not reassigning `arr`
// arr = [];   // ERROR — reassignment of `arr` is forbidden

const obj = { name: "Alice" };
obj.name = "Bob";  // OK — mutating the object's property
// obj = {};       // ERROR — reassignment forbidden
```

### Type Annotation Syntax

```typescript
// Explicit annotation
let variableName: TypeName = value;

// Examples
let age: number = 25;
let username: string = "alice";
let isLoggedIn: boolean = false;
let data: unknown = null;

// Usually, TypeScript infers the type — you don't need to write it
let inferredAge = 25;          // TypeScript infers: number
let inferredName = "alice";    // TypeScript infers: string
let inferredFlag = false;      // TypeScript infers: boolean
```

---

## 3.2 Primitive Types: string, number, boolean

### number

TypeScript (like JavaScript) has a single `number` type for all numeric values — integers and floats alike.

```typescript
let integer: number = 42;
let float: number = 3.14;
let negative: number = -100;
let bigNumber: number = 1_000_000;  // underscores for readability
let hex: number = 0xFF;             // 255
let binary: number = 0b1010;        // 10
let octal: number = 0o17;           // 15

// Special values (all type: number)
let infinity: number = Infinity;
let negInfinity: number = -Infinity;
let notANumber: number = NaN;

// Number methods
let n = 3.14159;
console.log(n.toFixed(2));        // "3.14"
console.log(n.toString());        // "3.14159"
console.log(Math.floor(n));       // 3
console.log(Math.ceil(n));        // 4
console.log(Number.isInteger(n)); // false
console.log(Number.isFinite(n));  // true
```

### string

```typescript
let firstName: string = "Alice";
let lastName: string = 'Smith';     // single quotes also OK
let greeting: string = `Hello, ${firstName} ${lastName}!`;  // template literal

// String methods — TypeScript knows them all
let message = "Hello, World!";
console.log(message.length);          // 13
console.log(message.toUpperCase());   // "HELLO, WORLD!"
console.log(message.includes("World")); // true
console.log(message.split(", "));     // ["Hello", "World!"]
console.log(message.slice(7, 12));    // "World"
console.log(message.trim());          // removes leading/trailing whitespace

// Multi-line strings
const multiLine: string = `
  Line 1
  Line 2
  Line 3
`.trim();
```

### boolean

```typescript
let isActive: boolean = true;
let isComplete: boolean = false;

// TypeScript correctly types conditional expressions
function checkAge(age: number): boolean {
  return age >= 18;
}

const canVote: boolean = checkAge(20);  // true
```

### Type Checking at Work

```typescript
function formatName(first: string, last: string): string {
  return `${first} ${last}`;
}

formatName("Alice", "Smith");    // OK
formatName("Alice", 42);         // ERROR: Argument of type 'number' is not assignable to parameter of type 'string'
formatName(true, "Smith");       // ERROR: Argument of type 'boolean' is not assignable to parameter of type 'string'
```

---

## 3.3 bigint and symbol

### bigint

`bigint` is for integers larger than `Number.MAX_SAFE_INTEGER` (2^53 - 1):

```typescript
// bigint literals end with 'n'
const bigNumber: bigint = 9007199254740993n;
const huge: bigint = 123456789012345678901234567890n;

// Arithmetic with bigint
const a: bigint = 10n;
const b: bigint = 3n;
console.log(a + b);   // 13n
console.log(a * b);   // 30n
console.log(a / b);   // 3n (integer division — no fractions)
console.log(a % b);   // 1n

// Cannot mix bigint and number without explicit conversion
// const mixed = 10n + 5;  // ERROR: Cannot mix BigInt and other types
const converted: number = Number(10n);  // explicit conversion OK

// Use case: cryptography, financial calculations, precise large integers
const factorial = (n: bigint): bigint => {
  if (n <= 1n) return 1n;
  return n * factorial(n - 1n);
};
console.log(factorial(20n));  // 2432902008176640000n
```

**Note**: `bigint` requires `"target": "ES2020"` or later in `tsconfig.json`.

### symbol

`symbol` creates unique, immutable values. Useful for property keys that can't accidentally collide:

```typescript
// Every Symbol() call creates a unique value
const sym1: symbol = Symbol("description");
const sym2: symbol = Symbol("description");

console.log(sym1 === sym2);  // false — every symbol is unique
console.log(sym1.toString()); // "Symbol(description)"
console.log(sym1.description); // "description"

// Symbol as object keys — guaranteed to not collide
const USER_ID = Symbol("userId");
const ADMIN_FLAG = Symbol("adminFlag");

const user = {
  name: "Alice",
  [USER_ID]: 42,          // Symbol as property key
  [ADMIN_FLAG]: true,
};

console.log(user[USER_ID]);   // 42
console.log(user[ADMIN_FLAG]); // true

// Symbols don't appear in JSON.stringify or for...in
console.log(JSON.stringify(user));  // {"name":"Alice"} — symbols excluded
for (const key in user) {
  console.log(key);  // only "name" — symbol keys skipped
}

// Well-known symbols (built into JS)
class Collection {
  private items: number[] = [1, 2, 3];

  [Symbol.iterator]() {
    let index = 0;
    return {
      next: () => ({
        value: this.items[index++],
        done: index > this.items.length,
      }),
    };
  }
}

const col = new Collection();
for (const item of col) {  // uses [Symbol.iterator]
  console.log(item);  // 1, 2, 3
}
```

---

## 3.4 null and undefined — Two Kinds of Nothing

JavaScript has two "nothing" values. TypeScript makes them explicit with `strictNullChecks`.

### Without strictNullChecks (Dangerous — Never Use)

```typescript
// strictNullChecks: false
let name: string = null;       // allowed — bad idea
let count: number = undefined; // allowed — bad idea
```

### With strictNullChecks: true (Always Use)

```typescript
// null and undefined are their own types
let nothing: null = null;
let missing: undefined = undefined;

// They cannot be assigned to other types
let name: string = null;       // ERROR with strictNullChecks
let count: number = undefined; // ERROR with strictNullChecks

// Explicit nullable types
let nullableName: string | null = null;     // OK — explicitly nullable
let optionalAge: number | undefined = undefined; // OK

// The optional chaining operator (?.) 
interface User {
  name: string;
  address?: {       // optional — might be undefined
    city: string;
    country: string;
  };
}

function getCity(user: User): string | undefined {
  return user.address?.city;  // safe — returns undefined if address is undefined
}

// Nullish coalescing (??)
function displayName(name: string | null | undefined): string {
  return name ?? "Anonymous";  // use "Anonymous" if name is null or undefined
}

console.log(displayName("Alice"));    // "Alice"
console.log(displayName(null));       // "Anonymous"
console.log(displayName(undefined));  // "Anonymous"
console.log(displayName(""));         // "" — empty string is NOT nullish
```

### null vs undefined — When to Use Which

```typescript
// undefined: absence of a value — something was never set
// null: explicit emptiness — something was intentionally set to "nothing"

interface ApiResponse {
  data: User | null;  // null = server confirmed: no user exists
  error?: string;     // undefined = no error (never set)
}

// Function return: undefined means "not found", null means "explicitly nothing"
function findUser(id: number): User | undefined {
  return database.get(id);  // undefined if not found
}

function clearCurrentUser(): null {
  return null;  // explicitly clearing
}
```

---

## 3.5 Type Inference — When You Don't Need to Write Types

TypeScript is extremely good at inferring types from context. Most of the time, you don't need to write explicit annotations.

### Variable Inference

```typescript
// TypeScript infers all these types:
const name = "Alice";       // string
const age = 25;             // number
const active = true;        // boolean
const ratio = 3.14;         // number
const items = [1, 2, 3];   // number[]
const nothing = null;       // null
```

### Return Type Inference

```typescript
// TypeScript infers the return type
function add(a: number, b: number) {
  return a + b;  // TypeScript infers: number
}

function getUser() {
  return { id: 1, name: "Alice" };  // TypeScript infers: { id: number; name: string }
}

// You can still add explicit return types for documentation/safety
function add2(a: number, b: number): number {  // explicit: better for public APIs
  return a + b;
}
```

### Array and Object Inference

```typescript
const numbers = [1, 2, 3];          // number[]
const strings = ["a", "b", "c"];    // string[]
const mixed = [1, "two", true];     // (string | number | boolean)[]

const user = {
  id: 1,
  name: "Alice",
  active: true,
};
// Inferred type:
// { id: number; name: string; active: boolean; }

// TypeScript knows the type of each property
user.name.toUpperCase();   // OK — TypeScript knows name is string
// user.name.toFixed(2);   // ERROR — toFixed is not on string
```

### When to Write Explicit Types

```typescript
// 1. Function parameters — always annotate
function greet(name: string): void {  // parameters can't be inferred
  console.log(`Hello, ${name}!`);
}

// 2. When inference gives something too broad
const items = [];       // inferred: never[] — too restrictive
const items2: string[] = [];  // explicit: string[]

// 3. Public API / exported functions — for documentation
export function calculate(a: number, b: number): number {
  return a + b;
}

// 4. When dealing with complex types
const handlers: Record<string, (data: unknown) => void> = {};

// 5. When TypeScript infers a wider type than you want
let status = "active";  // inferred: string (any string)
// Better: 
let status2: "active" | "inactive" | "pending" = "active";  // literal type
```

---

## 3.6 any — The Escape Hatch (and Why to Avoid It)

`any` is TypeScript's escape hatch. A variable of type `any` bypasses all type checking.

### What any Does

```typescript
let x: any = 5;
x = "hello";      // OK — no type checking
x = true;         // OK
x = [1, 2, 3];   // OK
x = null;         // OK

// You can access any property on any
x.foo.bar.baz;    // TypeScript: fine. Runtime: crash.
x();              // TypeScript: fine. Runtime: crash if x isn't callable.

// any infects everything it touches
function process(data: any) {
  return data.value;  // inferred return type: any
}

const result = process({ value: 42 });
result.anything.you.want;  // TypeScript: fine. Runtime: maybe crash.
```

### The Problem With any

```typescript
// any spreads "type unsafety" through your codebase
function parseData(json: any): any {
  return JSON.parse(json);
}

const data = parseData('{"name": "Alice"}');
// TypeScript thinks data.name exists, but has no idea what type it is
// data.nme — typo! TypeScript won't catch this. Runtime returns undefined.
```

### When any Is Acceptable

```typescript
// 1. Migrating from JavaScript — temporary
const legacyData: any = getLegacyApiData();

// 2. Highly dynamic code that's hard to type precisely
function mixin(target: any, ...sources: any[]): any {
  return Object.assign(target, ...sources);
}

// 3. In tests, sometimes for convenience
// (but prefer unknown + type guards)
```

### Escaping any Safely

```typescript
// Instead of working with any, validate and narrow
function processInput(input: any): string {
  // Validate before trusting the any value
  if (typeof input !== "string") {
    throw new TypeError(`Expected string, got ${typeof input}`);
  }
  return input.toUpperCase();  // TypeScript now knows input is string
}
```

---

## 3.7 unknown — Type-Safe any

`unknown` was added in TypeScript 3.0 as the type-safe alternative to `any`. Like `any`, `unknown` can hold any value. Unlike `any`, you must narrow it before using it.

### unknown vs any

```typescript
let a: any = "hello";
a.toUpperCase();      // TypeScript: OK (trusting you blindly)
a.nonExistent();      // TypeScript: OK (still trusting you)

let u: unknown = "hello";
// u.toUpperCase();   // ERROR: Object is of type 'unknown'
// u.nonExistent();   // ERROR: Object is of type 'unknown'

// Must narrow first
if (typeof u === "string") {
  u.toUpperCase();    // OK — TypeScript knows u is string here
}
```

### Using unknown Safely

```typescript
function processData(data: unknown): string {
  // Must handle all cases before using
  if (typeof data === "string") {
    return data.toUpperCase();
  }
  if (typeof data === "number") {
    return data.toFixed(2);
  }
  if (Array.isArray(data)) {
    return data.join(", ");
  }
  return String(data);  // fallback
}

// unknown for API responses — the right approach
async function fetchUser(id: number): Promise<unknown> {
  const response = await fetch(`/api/users/${id}`);
  return response.json();  // JSON can be anything — unknown is honest
}

// Type guard to validate the shape
interface User {
  id: number;
  name: string;
}

function isUser(value: unknown): value is User {
  return (
    typeof value === "object" &&
    value !== null &&
    "id" in value &&
    "name" in value &&
    typeof (value as any).id === "number" &&
    typeof (value as any).name === "string"
  );
}

async function getUser(): Promise<User> {
  const data = await fetchUser(1);
  if (!isUser(data)) {
    throw new Error("Invalid user data from API");
  }
  return data;  // TypeScript knows this is User
}
```

### The Key Difference

```typescript
// any: TypeScript trusts you — gives up type checking
// unknown: TypeScript forces you to prove what you know about the value

// Use any when: you genuinely don't care about types (rare)
// Use unknown when: you don't know the type yet but will validate it
```

---

## 3.8 never — The Bottom Type

`never` represents the type of values that **never occur**. It is the "bottom" type — a subtype of every type, but no value is a `never`.

### When never Appears

```typescript
// 1. Functions that never return
function throwError(message: string): never {
  throw new Error(message);  // always throws — never returns
}

function infiniteLoop(): never {
  while (true) {}  // never returns
}

// 2. Exhaustiveness checking — the most important use case
type Shape = "circle" | "square" | "triangle";

function getArea(shape: Shape, size: number): number {
  switch (shape) {
    case "circle":
      return Math.PI * size * size;
    case "square":
      return size * size;
    case "triangle":
      return (size * size) / 2;
    default:
      // If we've handled all cases, shape is `never` here
      // If we add a new Shape and forget to handle it, TypeScript errors here
      const _exhaustiveCheck: never = shape;
      throw new Error(`Unknown shape: ${shape}`);
  }
}

// 3. Impossible type intersections
type NumberAndString = number & string;  // impossible — this is never
```

### Exhaustiveness with never

```typescript
// This is the real power of never
type Status = "pending" | "active" | "completed" | "failed";

function describeStatus(status: Status): string {
  switch (status) {
    case "pending":   return "Waiting to start";
    case "active":    return "Currently running";
    case "completed": return "Done successfully";
    case "failed":    return "Failed with error";
    default:
      // TypeScript infers status as never here
      // because all cases are handled
      // If you add "cancelled" to Status, TypeScript errors here
      const impossible: never = status;
      throw new Error(`Unhandled status: ${impossible}`);
  }
}
```

---

## 3.9 Literal Types — Specific Values as Types

TypeScript can use specific string, number, or boolean values as types. Instead of "any string", you can say "exactly this string".

### String Literal Types

```typescript
// A type that is exactly the string "north"
type Direction = "north" | "south" | "east" | "west";

let heading: Direction = "north";   // OK
// let heading2: Direction = "up";  // ERROR: Type '"up"' is not assignable to type 'Direction'

function move(direction: Direction, distance: number): void {
  console.log(`Moving ${distance} units ${direction}`);
}

move("south", 10);    // OK
// move("sideways", 5);  // ERROR
```

### Number Literal Types

```typescript
type DiceValue = 1 | 2 | 3 | 4 | 5 | 6;

function rollDice(): DiceValue {
  return (Math.floor(Math.random() * 6) + 1) as DiceValue;
}

type HttpStatus = 200 | 201 | 400 | 401 | 403 | 404 | 500;
```

### Boolean Literal Types

```typescript
// Rarely used for booleans alone, but useful in combinations
type AlwaysTrue = true;
type Flag = true | false;  // same as boolean
```

### Literal Types in Objects

```typescript
interface ApiResponse<T> {
  status: "success" | "error";  // literal union
  data: T | null;
  message: string;
}

const successResponse: ApiResponse<User> = {
  status: "success",   // must be exactly "success" or "error"
  data: { id: 1, name: "Alice" },
  message: "User retrieved",
};

// Discriminated union using literals
type Result<T> =
  | { status: "success"; data: T }
  | { status: "error"; error: string; code: number };

function handleResult<T>(result: Result<T>): void {
  if (result.status === "success") {
    console.log(result.data);     // TypeScript knows data exists
  } else {
    console.log(result.error);    // TypeScript knows error and code exist
    console.log(result.code);
  }
}
```

### const Assertion — Widening Prevention

```typescript
// Without const assertion — TypeScript widens the type
let config = {
  host: "localhost",  // inferred: string (any string)
  port: 3000,         // inferred: number (any number)
};

// With const assertion — all values become literal types
const config2 = {
  host: "localhost",  // inferred: "localhost" (exactly this string)
  port: 3000,         // inferred: 3000 (exactly this number)
} as const;

type Config2 = typeof config2;
// { readonly host: "localhost"; readonly port: 3000; }

// Useful for arrays
const directions = ["north", "south", "east", "west"] as const;
type Direction = typeof directions[number];  // "north" | "south" | "east" | "west"
```

---

## 3.10 Type Assertions and Type Casting

Sometimes you know more about a type than TypeScript does. Type assertions let you override TypeScript's inference.

### The as Keyword

```typescript
// Tells TypeScript: "trust me, I know this is a string"
const input = document.getElementById("username") as HTMLInputElement;
console.log(input.value);  // TypeScript now allows .value

// Without assertion:
const input2 = document.getElementById("username");
// input2.value  // ERROR: Property 'value' does not exist on type 'HTMLElement'
```

### When Type Assertions are Acceptable

```typescript
// 1. DOM elements — you know the specific type
const canvas = document.querySelector("#canvas") as HTMLCanvasElement;
const ctx = canvas.getContext("2d")!;  // ! is non-null assertion

// 2. JSON parsing — you know the shape
const data = JSON.parse(json) as User;

// 3. Type narrowing that TypeScript can't verify
const value: string | number = getValue();
const str = value as string;  // assertion

// Better alternative to the above: use a type guard
function isString(v: string | number): v is string {
  return typeof v === "string";
}
```

### The Non-Null Assertion Operator (!)

```typescript
// ! tells TypeScript: "this value is definitely not null or undefined"
function getElement(id: string): HTMLElement {
  const el = document.getElementById(id);
  // el is HTMLElement | null
  return el!;  // asserts non-null — you're responsible for ensuring this is true
}

// Be careful — this crashes at runtime if el is actually null
const el = document.getElementById("nonexistent")!;
el.textContent = "hello";  // TypeError: Cannot set properties of null
```

### Double Assertions (Force Casting)

```typescript
// Sometimes TypeScript rejects an assertion as "too different"
const n: number = 5;
// const s = n as string;  // ERROR: Conversion of type 'number' to type 'string' may be a mistake

// Force it with a double assertion (red flag — usually a design error)
const s = n as unknown as string;  // Force assertion via unknown

// When you see this in code, it's a warning sign
```

### Type Assertion vs Type Annotation

```typescript
// Assertion: you override TypeScript's conclusion
const x = value as string;       // "I'm telling TypeScript it's a string"

// Annotation: you tell TypeScript from the start
const y: string = someString;    // "TypeScript, this is a string from the beginning"

// Annotations are checked — assertions are trusted
const z: string = 42;            // ERROR: 42 is not a string
const w = 42 as unknown as string;  // No error — you're overriding TypeScript
```

---

## Type System Summary Diagram

```
TypeScript Types
│
├── Primitive Types (have JavaScript equivalents)
│   ├── string
│   ├── number
│   ├── boolean
│   ├── bigint
│   ├── symbol
│   ├── null
│   └── undefined
│
├── Special Types
│   ├── any      — escape hatch, no type checking
│   ├── unknown  — safe any, must narrow before use
│   ├── never    — impossible type, bottom of the type hierarchy
│   └── void     — function returns nothing meaningful
│
├── Literal Types
│   ├── "hello"         — exact string value
│   ├── 42              — exact number value
│   └── true / false    — exact boolean value
│
└── Compound Types (covered in later chapters)
    ├── string[]         — array of strings
    ├── string | number  — union type
    ├── { id: number }   — object type
    └── (a: T) => R      — function type
```

---

## Summary

TypeScript's basic type system maps closely to JavaScript's primitive types: `string`, `number`, `boolean`, `bigint`, `symbol`, `null`, and `undefined`. TypeScript adds three special types: `any` (unsafe escape hatch), `unknown` (safe escape hatch — must narrow before use), and `never` (impossible type — used for exhaustiveness checking). Literal types allow specific values to be types. Type inference reduces the need to write explicit annotations. Type assertions (`as`) override TypeScript's judgment but shift responsibility to the developer.

---

## Key Takeaways

- TypeScript's primitive types mirror JavaScript's, but `null` and `undefined` are separate types (with `strictNullChecks`)
- **Prefer `const` over `let`** — and always add types to function parameters
- **`any` disables type checking** — use `unknown` as the safe alternative
- **`never` is for exhaustiveness** — catch unhandled cases in switch statements at compile time
- **Literal types** restrict values to specific constants — more precise than `string` or `number`
- **Type inference is powerful** — only add annotations where inference fails or for documentation
- **`as` assertions bypass safety** — use only when you genuinely know more than the compiler

---

## Practice Questions

1. What is the difference between `any` and `unknown`? Give an example where you'd use each.
2. What does `never` represent? How is it used in exhaustiveness checking?
3. What is a literal type? How does it differ from `string`?
4. What does `strictNullChecks` change about `null` and `undefined`?
5. When should you use a type annotation vs. relying on inference?
6. What is the `as const` assertion and why is it useful?

---

## Exercises

**Exercise 1**: Without using `any`, write a function `safeParseJSON(input: string): unknown` that parses JSON. Then write a type guard `isUser(data: unknown): data is User` that validates the returned data has the shape `{ id: number; name: string; email: string }`.

**Exercise 2**: Create a type `HttpMethod = "GET" | "POST" | "PUT" | "DELETE" | "PATCH"`. Write a function `createRequest(method: HttpMethod, url: string, body?: unknown): string` that formats an HTTP request description.

**Exercise 3**: Write a function `processValue(value: string | number | boolean | null): string` that handles each case and returns a string representation. Use `never` in a default case to ensure exhaustiveness.

**Exercise 4**: Create an object `COLORS` with five color names as keys and hex strings as values, using `as const`. Then derive a type `ColorName` from its keys and `HexColor` from its values.

---

*Next: [Chapter 4 — Functions](04-functions.md)*
