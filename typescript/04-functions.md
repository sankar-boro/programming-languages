# Chapter 4 — Functions

> *"A well-typed function is a contract. It tells you exactly what it needs and exactly what it provides."*

---

## 4.1 Function Type Syntax

TypeScript allows you to type every aspect of a function: its parameters, their types, and its return type.

### Basic Function Declaration

```typescript
// Syntax: function name(param: Type, ...): ReturnType
function add(a: number, b: number): number {
  return a + b;
}

// TypeScript infers return type if you omit it
function multiply(a: number, b: number) {
  return a * b;  // TypeScript infers: number
}

// Explicit return type is better for public APIs — it's a contract
function divide(a: number, b: number): number {
  if (b === 0) throw new Error("Division by zero");
  return a / b;
}
```

### Function Type Expressions

Functions are first-class values in JavaScript. TypeScript has syntax to type them:

```typescript
// Function type expression
type Adder = (a: number, b: number) => number;

// Assign a function to a typed variable
const add: Adder = (a, b) => a + b;

// Inline function type in function parameters
function applyOperation(
  x: number,
  y: number,
  operation: (a: number, b: number) => number
): number {
  return operation(x, y);
}

applyOperation(5, 3, (a, b) => a + b);  // 8
applyOperation(5, 3, (a, b) => a * b);  // 15

// Type an object property as a function
interface Calculator {
  add: (a: number, b: number) => number;
  subtract: (a: number, b: number) => number;
}

const calc: Calculator = {
  add: (a, b) => a + b,
  subtract: (a, b) => a - b,
};
```

### Call Signatures

For more complex function types (with properties), use call signatures:

```typescript
// A function with properties
interface LogFunction {
  (message: string, level: "info" | "warn" | "error"): void;
  prefix: string;
  callCount: number;
}

function createLogger(): LogFunction {
  const log = function (message: string, level: "info" | "warn" | "error") {
    log.callCount++;
    console.log(`[${log.prefix}] [${level.toUpperCase()}] ${message}`);
  };
  log.prefix = "APP";
  log.callCount = 0;
  return log;
}

const logger = createLogger();
logger("Server started", "info");
logger("Low memory", "warn");
console.log(`Called ${logger.callCount} times`);
```

---

## 4.2 Optional Parameters

Mark a parameter optional with `?`. Optional parameters must come after required ones.

```typescript
// Optional parameter: string | undefined
function greet(name: string, greeting?: string): string {
  // greeting is string | undefined here
  const g = greeting ?? "Hello";
  return `${g}, ${name}!`;
}

greet("Alice");              // "Hello, Alice!"
greet("Alice", "Hi");       // "Hi, Alice!"
greet("Alice", undefined);  // "Hello, Alice!" — same as omitting

// TypeScript ensures optional params aren't accidentally used without checking
function process(value: string, transformer?: (s: string) => string): string {
  if (transformer) {
    return transformer(value);  // safe — checked above
  }
  return value;
}

// Optional in interfaces
interface RequestOptions {
  timeout?: number;      // may or may not be present
  headers?: Record<string, string>;
  retries?: number;
}

function makeRequest(url: string, options?: RequestOptions): Promise<Response> {
  const timeout = options?.timeout ?? 5000;  // safe chaining
  return fetch(url, { signal: AbortSignal.timeout(timeout) });
}
```

---

## 4.3 Default Parameters

Default parameters provide a fallback value when no argument is passed:

```typescript
// Default parameter — not optional, just has a default
function createUser(
  name: string,
  role: "admin" | "user" | "guest" = "user",
  active: boolean = true
): object {
  return { name, role, active };
}

createUser("Alice");                  // { name: "Alice", role: "user", active: true }
createUser("Bob", "admin");           // { name: "Bob", role: "admin", active: true }
createUser("Charlie", "guest", false); // { name: "Charlie", role: "guest", active: false }

// Key difference from optional: default params have a definite type
function format(value: number, precision: number = 2): string {
  return value.toFixed(precision);
  // precision is `number` here, not `number | undefined`
  // because it has a default value
}

// Defaults can be expressions — evaluated at call time
let defaultCount = 0;
function increment(step: number = ++defaultCount): number {
  return step;
}
```

---

## 4.4 Rest Parameters

Rest parameters collect remaining arguments into a typed array:

```typescript
// Rest parameter: must be last, always an array type
function sum(...numbers: number[]): number {
  return numbers.reduce((acc, n) => acc + n, 0);
}

sum(1, 2, 3);           // 6
sum(1, 2, 3, 4, 5);     // 15
sum();                   // 0 — empty array

// Rest with typed tuples (TypeScript 4.0+)
function concat(separator: string, ...parts: string[]): string {
  return parts.join(separator);
}

concat("-", "2024", "01", "15");  // "2024-01-15"

// Spread into rest parameters
const nums = [1, 2, 3, 4, 5];
console.log(sum(...nums));  // 15

// Rest and required params together
function log(level: "info" | "warn" | "error", ...messages: string[]): void {
  messages.forEach((msg) => console.log(`[${level}] ${msg}`));
}

log("info", "Server started", "Listening on port 3000");

// Typed tuple rest parameters (advanced)
function first<T>(...args: [T, ...unknown[]]): T {
  return args[0];
}

console.log(first(42, "extra", true));  // 42, typed as number
```

---

## 4.5 Arrow Functions and Their Types

Arrow functions (`=>`) are concise function expressions. TypeScript types them identically to regular functions.

```typescript
// Regular function
function double(n: number): number {
  return n * 2;
}

// Arrow function — equivalent
const double2 = (n: number): number => n * 2;

// Explicit type annotation
const double3: (n: number) => number = (n) => n * 2;

// Multi-line arrow function
const processUser = (user: { name: string; age: number }): string => {
  const normalized = user.name.toLowerCase();
  return `${normalized}-${user.age}`;
};

// Array methods with arrow functions
const numbers = [1, 2, 3, 4, 5];

const doubled: number[] = numbers.map((n) => n * 2);
const evens: number[] = numbers.filter((n) => n % 2 === 0);
const total: number = numbers.reduce((acc, n) => acc + n, 0);
const found: number | undefined = numbers.find((n) => n > 3);

// TypeScript infers correctly from the callback
numbers.forEach((n) => {
  // TypeScript knows n is number here
  console.log(n.toFixed(2));
});
```

### Key Difference: this Binding

```typescript
// Regular functions have their own `this`
// Arrow functions capture `this` from the enclosing scope

class Timer {
  private seconds: number = 0;

  start(): void {
    // BAD: regular function — `this` is undefined in strict mode
    // setInterval(function() {
    //   this.seconds++;  // ERROR at runtime
    // }, 1000);

    // GOOD: arrow function — `this` is the Timer instance
    setInterval(() => {
      this.seconds++;  // `this` is correctly the Timer instance
      console.log(this.seconds);
    }, 1000);
  }
}
```

---

## 4.6 Function Overloading

TypeScript allows multiple function signatures for the same function — useful when a function behaves differently based on input types.

### The Problem Without Overloads

```typescript
// Without overloads — imprecise typing
function padLeft(value: string, padding: string | number): string {
  if (typeof padding === "number") {
    return " ".repeat(padding) + value;
  }
  return padding + value;
}

// The return type is always string, but we can be more precise
// about what inputs produce what output
```

### Overload Signatures

```typescript
// Overload signatures: define the variants
function createElement(tag: "div"): HTMLDivElement;
function createElement(tag: "span"): HTMLSpanElement;
function createElement(tag: "input"): HTMLInputElement;
// Implementation signature: not visible to callers
function createElement(tag: string): HTMLElement {
  return document.createElement(tag);
}

// Usage — TypeScript knows the exact return type
const div: HTMLDivElement = createElement("div");
const span: HTMLSpanElement = createElement("span");
const input: HTMLInputElement = createElement("input");
// const p: HTMLParagraphElement = createElement("p");  // ERROR — no overload for "p"
```

### A Practical Overload Example

```typescript
// formatDate: different behavior based on input type
function formatDate(date: Date): string;
function formatDate(timestamp: number): string;
function formatDate(dateString: string): string;
function formatDate(input: Date | number | string): string {
  let date: Date;

  if (input instanceof Date) {
    date = input;
  } else if (typeof input === "number") {
    date = new Date(input);
  } else {
    date = new Date(input);
  }

  return date.toISOString().split("T")[0];
}

formatDate(new Date());           // "2024-01-15"
formatDate(1705276800000);        // "2024-01-15"
formatDate("2024-01-15");         // "2024-01-15"
// formatDate(true);              // ERROR — no matching overload
```

### When to Use Overloads vs Unions

```typescript
// Prefer overloads when return type depends on input type:
function parse(value: string): string[];
function parse(value: number): number[];
function parse(value: string | number): string[] | number[] {
  if (typeof value === "string") {
    return value.split("");
  }
  return [value];
}

// Caller gets precise type based on what they passed:
const strings: string[] = parse("hello");
const numbers: number[] = parse(42);

// Without overloads, return type is always string[] | number[]
// even when caller knows they passed a string
```

---

## 4.7 void and never as Return Types

### void

`void` means the function doesn't return a meaningful value. It's the return type of functions called for their side effects.

```typescript
function log(message: string): void {
  console.log(message);
  // implicit return undefined — that's fine for void
}

// void functions can technically return undefined
function setup(): void {
  initDatabase();
  // no return needed
  return;           // OK
  // return undefined;  // OK
  // return 5;          // ERROR — can't return a value from void function
}

// void in callbacks — you don't care what the callback returns
type Callback = () => void;

function runCallback(cb: Callback): void {
  cb();  // call it, ignore return value
}

// Interesting: void callbacks can return a value — it's just ignored
const nums = [1, 2, 3];
nums.forEach((n): void => {
  // TypeScript allows returning from a void callback
  // because the return value is ignored anyway
  return;
});
```

### void vs undefined

```typescript
// void: "I don't return anything useful"
function sideEffect(): void {
  console.log("side effect");
}

// undefined: "I explicitly return undefined"
function explicitUndefined(): undefined {
  return undefined;  // must return undefined explicitly
}

// The difference matters for callbacks
type Returns = () => undefined;
type VoidCallback = () => void;

const r: Returns = () => undefined;  // must return undefined
const v: VoidCallback = () => 42;    // fine — return value ignored
```

---

## 4.8 this in Functions — Typing the Context

TypeScript can type the `this` parameter to prevent incorrect usage.

```typescript
// TypeScript's fake 'this' parameter — not a real parameter
interface User {
  name: string;
  greet(this: User): string;  // 'this' must be a User
}

const user: User = {
  name: "Alice",
  greet() {
    return `Hello, I'm ${this.name}`;  // TypeScript knows this is User
  },
};

// Detaching a method causes TypeScript to complain
const greet = user.greet;
// greet();  // ERROR: The 'this' context is not of type 'User'
greet.call(user);  // OK — passing the correct `this`

// Explicit this typing in functions
function formatUser(this: { name: string; age: number }): string {
  return `${this.name} (${this.age})`;
}

formatUser.call({ name: "Alice", age: 30 });  // "Alice (30)"
```

---

## 4.9 Higher-Order Functions

Functions that take functions as arguments or return functions.

```typescript
// Map, filter, reduce — standard higher-order functions
function mapArray<T, U>(arr: T[], fn: (item: T, index: number) => U): U[] {
  return arr.map(fn);
}

const doubled = mapArray([1, 2, 3], (n) => n * 2);     // number[]
const asStrings = mapArray([1, 2, 3], (n) => `${n}`);  // string[]

// Function factories
function makeMultiplier(factor: number): (n: number) => number {
  return (n) => n * factor;
}

const double = makeMultiplier(2);
const triple = makeMultiplier(3);
console.log(double(5));  // 10
console.log(triple(5));  // 15

// Currying
function curry<A, B, C>(
  fn: (a: A, b: B) => C
): (a: A) => (b: B) => C {
  return (a) => (b) => fn(a, b);
}

const add = (a: number, b: number) => a + b;
const curriedAdd = curry(add);
const addFive = curriedAdd(5);
console.log(addFive(3));   // 8
console.log(addFive(10));  // 15

// Function composition
function compose<A, B, C>(
  f: (b: B) => C,
  g: (a: A) => B
): (a: A) => C {
  return (a) => f(g(a));
}

const trim = (s: string) => s.trim();
const upper = (s: string) => s.toUpperCase();
const trimAndUpper = compose(upper, trim);

console.log(trimAndUpper("  hello  "));  // "HELLO"

// Memoization — cache function results
function memoize<T extends (...args: unknown[]) => unknown>(fn: T): T {
  const cache = new Map<string, unknown>();
  return ((...args: unknown[]) => {
    const key = JSON.stringify(args);
    if (cache.has(key)) return cache.get(key);
    const result = fn(...args);
    cache.set(key, result);
    return result;
  }) as T;
}

const expensiveCalc = memoize((n: number): number => {
  console.log(`Computing for ${n}...`);
  return n * n;
});

console.log(expensiveCalc(5));  // Computing for 5... → 25
console.log(expensiveCalc(5));  // (cache hit) → 25
```

### Partial Application

```typescript
// Make some arguments of a function fixed
function partial<T extends unknown[], U extends unknown[], R>(
  fn: (...args: [...T, ...U]) => R,
  ...presetArgs: T
): (...remainingArgs: U) => R {
  return (...remainingArgs: U) => fn(...presetArgs, ...remainingArgs);
}

function fetchWithAuth(
  baseUrl: string,
  token: string,
  endpoint: string
): Promise<Response> {
  return fetch(`${baseUrl}${endpoint}`, {
    headers: { Authorization: `Bearer ${token}` },
  });
}

// Fix baseUrl and token, leaving endpoint flexible
const apiGet = partial(fetchWithAuth, "https://api.example.com", "abc123");
apiGet("/users");    // fetches https://api.example.com/users
apiGet("/products"); // fetches https://api.example.com/products
```

---

## Function Types — Complete Summary

```typescript
// 1. Named function
function add(a: number, b: number): number { return a + b; }

// 2. Function expression
const add2 = function(a: number, b: number): number { return a + b; };

// 3. Arrow function
const add3 = (a: number, b: number): number => a + b;

// 4. Type alias for function type
type AddFn = (a: number, b: number) => number;
const add4: AddFn = (a, b) => a + b;

// 5. Method in interface (two syntaxes)
interface Math {
  add(a: number, b: number): number;     // method syntax
  subtract: (a: number, b: number) => number;  // property syntax
}

// Difference: method syntax is bivariant for this, property syntax is stricter
```

---

## Summary

TypeScript functions are typed along three dimensions: parameter types, return type, and optionality. Optional parameters use `?`, defaults use `=`, and rest parameters use `...`. Arrow functions capture `this` from the enclosing scope — important for class methods used as callbacks. Function overloading allows precise return types based on input types. `void` means "returns nothing useful"; `never` means "never returns". Higher-order functions — functions that take or return other functions — work naturally with TypeScript's type inference.

---

## Key Takeaways

- Always type function **parameters** — they cannot be inferred
- Return types are optional (inference works) but recommended for exported functions
- `optional?: T` means the parameter may be absent; `default = value` means it has a fallback
- Rest parameters (`...args: T[]`) must be the last parameter
- Arrow functions capture `this` from surrounding scope — prefer for callbacks
- **Overloads** provide precise type mapping from input to output when a function behaves differently based on input type
- `void` = "doesn't return usefully"; `never` = "never returns at all"

---

## Practice Questions

1. What is the difference between an optional parameter and a parameter with a default value?
2. Can a function declared with return type `void` return a value? What happens?
3. When would you use function overloading instead of union types for parameters?
4. What is the difference between `this` binding in regular functions vs arrow functions?
5. What is a call signature, and when would you use it instead of a function type expression?

---

## Exercises

**Exercise 1**: Write a `pipeline` function that takes a value and an array of transformation functions, applying each in sequence. Type it with generics so TypeScript knows the output type.

```typescript
// pipeline(5, [(n: number) => n * 2, (n: number) => n + 1])  // should return 11
```

**Exercise 2**: Write an overloaded `stringify` function that:
- Takes a `number` and returns a `string` with 2 decimal places
- Takes a `boolean` and returns `"yes"` or `"no"`
- Takes a `Date` and returns an ISO date string

**Exercise 3**: Implement `debounce<T extends (...args: unknown[]) => unknown>(fn: T, delay: number): T` that delays function execution until `delay` ms have passed since the last call.

**Exercise 4**: Write a type-safe `once<T extends (...args: unknown[]) => unknown>(fn: T): T` that ensures a function is only called once, returning the cached result on subsequent calls.

---

*Next: [Chapter 5 — Objects and Interfaces](05-objects-interfaces.md)*
