# Chapter 6 — Advanced Type System

> *"TypeScript's type narrowing is control flow analysis. The compiler understands your code's logic and adjusts types accordingly."*

---

## 6.1 Union Types — A or B

A union type `A | B` means a value can be of type A **or** type B. It's the TypeScript equivalent of "this OR that".

```typescript
// Basic union
let id: string | number;
id = "user-123";  // OK
id = 42;          // OK
// id = true;     // ERROR: boolean is not in the union

// Union in function parameters
function formatId(id: string | number): string {
  if (typeof id === "string") {
    return id.toUpperCase();
  }
  return id.toString();  // TypeScript knows id is number here
}

// Wider unions
type Status = "pending" | "active" | "completed" | "failed";
type HttpStatus = 200 | 201 | 400 | 401 | 403 | 404 | 500;

function handleResponse(status: HttpStatus): string {
  if (status === 200 || status === 201) return "Success";
  if (status >= 400 && status < 500) return "Client error";
  return "Server error";
}

// Nullable via union
type MaybeUser = User | null;
type OptionalName = string | undefined;
```

### Properties Available on Unions

With a union, you can only access properties that exist on ALL members:

```typescript
interface Cat {
  name: string;
  purr(): void;
  meow(): void;
}

interface Dog {
  name: string;
  bark(): void;
  fetch(): void;
}

type Pet = Cat | Dog;

function greetPet(pet: Pet): void {
  console.log(pet.name);  // OK — name exists on both Cat and Dog
  // pet.purr();  // ERROR — not on Dog
  // pet.bark();  // ERROR — not on Cat
}
```

---

## 6.2 Intersection Types — A and B

An intersection type `A & B` means a value must satisfy both A **and** B — it combines types.

```typescript
interface Serializable {
  serialize(): string;
}

interface Loggable {
  log(message: string): void;
}

// Intersection: must have both
type SerializableAndLoggable = Serializable & Loggable;

function processEntity(entity: SerializableAndLoggable): void {
  entity.log("Processing...");
  const data = entity.serialize();
  entity.log(`Serialized: ${data}`);
}

// Merging object types with intersection
type UserBase = {
  id: number;
  name: string;
};

type UserWithEmail = UserBase & {
  email: string;
};

type AdminUser = UserWithEmail & {
  role: "admin";
  permissions: string[];
};

const admin: AdminUser = {
  id: 1,
  name: "Alice",
  email: "alice@admin.com",
  role: "admin",
  permissions: ["read", "write", "delete"],
};
```

### Impossible Intersections

```typescript
// Intersecting incompatible primitives produces never
type NumberAndString = number & string;  // never — impossible

function getValueNever(): NumberAndString {
  // Can never return a value that is both number and string
  throw new Error("Impossible");
}
```

### Union vs Intersection

```typescript
// Union (|): either one OR the other — less specific, more flexible
// Intersection (&): BOTH simultaneously — more specific, more constrained

// Think of it mathematically:
// Union: A ∪ B — elements in A OR B (bigger set of types, smaller set of properties)
// Intersection: A ∩ B — elements in BOTH A AND B (smaller set of types, bigger set of properties)

interface HasId { id: number }
interface HasName { name: string }

// Union: objects that have EITHER id OR name (or both)
// type Either = HasId | HasName;  // can access neither id nor name without narrowing

// Intersection: objects that have BOTH id AND name
type Both = HasId & HasName;  // can access both id and name freely
```

---

## 6.3 Type Narrowing — Drilling Down to Specific Types

**Type narrowing** is the process by which TypeScript refines a broader type to a more specific one within a block of code, based on the code's logic. This is one of TypeScript's most sophisticated and useful features.

### The Core Idea

```typescript
function process(value: string | number): string {
  // Here: value is string | number
  
  if (typeof value === "string") {
    // Here: TypeScript NARROWS value to string
    return value.toUpperCase();  // string methods available
  }
  
  // Here: TypeScript NARROWS value to number (string was eliminated)
  return value.toFixed(2);  // number methods available
}
```

TypeScript's control flow analysis tracks the type through every branch.

---

## 6.4 typeof Guards

```typescript
function processInput(input: string | number | boolean | null | undefined): string {
  if (typeof input === "string") {
    // input: string
    return input.toUpperCase();
  }
  
  if (typeof input === "number") {
    // input: number
    return input.toFixed(2);
  }
  
  if (typeof input === "boolean") {
    // input: boolean
    return input ? "yes" : "no";
  }
  
  // input: null | undefined
  return "nothing";
}

// typeof narrows to: "string" | "number" | "boolean" | "bigint" | "symbol" | "undefined" | "function" | "object"
// Note: typeof null === "object" — JavaScript legacy!

function handleNullable(value: string | null | undefined): string {
  if (typeof value === "string") {
    return value;  // string — both null and undefined eliminated
  }
  return "";  // value is null | undefined here
}

// Equality narrowing
function processStatus(status: "active" | "inactive" | "pending"): string {
  if (status === "active") {
    // status: "active"
    return "Running";
  }
  if (status === "inactive") {
    // status: "inactive"
    return "Stopped";
  }
  // status: "pending"
  return "Waiting";
}
```

---

## 6.5 instanceof Guards

`instanceof` narrows based on class/constructor:

```typescript
class Circle {
  constructor(public radius: number) {}
  area(): number { return Math.PI * this.radius ** 2; }
}

class Rectangle {
  constructor(public width: number, public height: number) {}
  area(): number { return this.width * this.height; }
}

type Shape = Circle | Rectangle;

function describeShape(shape: Shape): string {
  if (shape instanceof Circle) {
    // shape: Circle — TypeScript knows
    return `Circle with radius ${shape.radius}, area ${shape.area().toFixed(2)}`;
  }
  // shape: Rectangle
  return `Rectangle ${shape.width}×${shape.height}, area ${shape.area()}`;
}

// instanceof with error handling
function handleError(error: unknown): string {
  if (error instanceof Error) {
    // error: Error — has message, stack, etc.
    return error.message;
  }
  if (typeof error === "string") {
    return error;
  }
  return "Unknown error";
}
```

---

## 6.6 in Operator Narrowing

The `in` operator checks if a property exists on an object, narrowing the type:

```typescript
interface Bird {
  fly(): void;
  wingspan: number;
}

interface Fish {
  swim(): void;
  depth: number;
}

type Animal = Bird | Fish;

function moveAnimal(animal: Animal): void {
  if ("fly" in animal) {
    // animal: Bird — has 'fly' property
    animal.fly();
    console.log(`Wingspan: ${animal.wingspan}`);
  } else {
    // animal: Fish — 'fly' not present
    animal.swim();
    console.log(`Depth: ${animal.depth}`);
  }
}

// in narrowing with optional properties
interface Config {
  mode: "development" | "production";
  debug?: boolean;
  verbose?: boolean;
}

function applyDebug(config: Config): void {
  if ("debug" in config && config.debug) {
    console.log("Debug mode enabled");
  }
}
```

---

## 6.7 Discriminated Unions — Tagged Unions

A **discriminated union** is a union of types that each have a common "tag" property — a literal type that uniquely identifies which variant you're dealing with. This is the most robust narrowing technique.

```typescript
// Each variant has a 'kind' property with a unique literal type
interface Circle {
  kind: "circle";  // discriminant
  radius: number;
}

interface Rectangle {
  kind: "rectangle";  // discriminant
  width: number;
  height: number;
}

interface Triangle {
  kind: "triangle";  // discriminant
  base: number;
  height: number;
}

type Shape = Circle | Rectangle | Triangle;

function calculateArea(shape: Shape): number {
  switch (shape.kind) {
    case "circle":
      // shape: Circle — TypeScript knows
      return Math.PI * shape.radius ** 2;
    
    case "rectangle":
      // shape: Rectangle — TypeScript knows
      return shape.width * shape.height;
    
    case "triangle":
      // shape: Triangle — TypeScript knows
      return (shape.base * shape.height) / 2;
    
    default:
      // Exhaustiveness check — if you add a new Shape, this errors
      const _exhaustive: never = shape;
      throw new Error(`Unknown shape: ${JSON.stringify(_exhaustive)}`);
  }
}

// Usage
const shapes: Shape[] = [
  { kind: "circle", radius: 5 },
  { kind: "rectangle", width: 4, height: 6 },
  { kind: "triangle", base: 3, height: 4 },
];

shapes.forEach((s) => console.log(`Area: ${calculateArea(s).toFixed(2)}`));
```

### Discriminated Unions for API Responses

```typescript
// The canonical discriminated union pattern for results
type ApiResult<T> =
  | { status: "success"; data: T; statusCode: 200 | 201 }
  | { status: "error"; error: string; statusCode: 400 | 401 | 403 | 404 | 500 }
  | { status: "loading" }
  | { status: "idle" };

function renderResult<T>(result: ApiResult<T>): string {
  switch (result.status) {
    case "success":
      return `Data: ${JSON.stringify(result.data)}`;  // result.data available
    case "error":
      return `Error ${result.statusCode}: ${result.error}`;  // result.error available
    case "loading":
      return "Loading...";
    case "idle":
      return "Nothing to show";
  }
}

// Redux-style action types
type Action =
  | { type: "INCREMENT"; by: number }
  | { type: "DECREMENT"; by: number }
  | { type: "RESET" }
  | { type: "SET_VALUE"; value: number };

interface State {
  count: number;
}

function reducer(state: State, action: Action): State {
  switch (action.type) {
    case "INCREMENT":
      return { count: state.count + action.by };
    case "DECREMENT":
      return { count: state.count - action.by };
    case "RESET":
      return { count: 0 };
    case "SET_VALUE":
      return { count: action.value };
  }
}
```

---

## 6.8 User-Defined Type Guards

Sometimes TypeScript can't automatically narrow a type. You write a **type guard function** — a function that returns a type predicate.

```typescript
// Syntax: parameter is Type
function isString(value: unknown): value is string {
  return typeof value === "string";
}

function isNumber(value: unknown): value is number {
  return typeof value === "number" && !isNaN(value);
}

// Using the type guard
function process(value: unknown): void {
  if (isString(value)) {
    // value: string — TypeScript narrows based on the type guard
    console.log(value.toUpperCase());
  } else if (isNumber(value)) {
    // value: number
    console.log(value.toFixed(2));
  }
}

// Type guard for interfaces
interface User {
  id: number;
  name: string;
  email: string;
}

function isUser(value: unknown): value is User {
  return (
    typeof value === "object" &&
    value !== null &&
    typeof (value as Record<string, unknown>).id === "number" &&
    typeof (value as Record<string, unknown>).name === "string" &&
    typeof (value as Record<string, unknown>).email === "string"
  );
}

// Now you can safely work with API responses
async function fetchUser(id: number): Promise<User> {
  const response = await fetch(`/api/users/${id}`);
  const data: unknown = await response.json();
  
  if (!isUser(data)) {
    throw new Error("Invalid user data received from API");
  }
  
  return data;  // TypeScript knows: User
}

// Array type guard
function isStringArray(value: unknown): value is string[] {
  return Array.isArray(value) && value.every(isString);
}

// Discriminated union type guard
function isCircle(shape: Shape): shape is Circle {
  return shape.kind === "circle";
}
```

---

## 6.9 Assertion Functions

Assertion functions throw if a condition is false, narrowing the type for code that follows:

```typescript
// Assertion function signature
function assert(condition: boolean, message: string): asserts condition {
  if (!condition) throw new Error(message);
}

// After calling assert, TypeScript narrows based on the condition
function processUser(user: User | null): void {
  assert(user !== null, "User must not be null");
  // user: User — narrowed after the assertion
  console.log(user.name);
}

// Asserts parameter is type
function assertIsString(val: unknown): asserts val is string {
  if (typeof val !== "string") {
    throw new TypeError(`Expected string, got ${typeof val}`);
  }
}

function processInput(input: unknown): void {
  assertIsString(input);
  // input: string — narrowed after the assertion
  console.log(input.toUpperCase());
}

// Non-null assertion function
function assertDefined<T>(val: T | null | undefined, name: string): asserts val is T {
  if (val === null || val === undefined) {
    throw new Error(`Expected ${name} to be defined`);
  }
}

function render(element: HTMLElement | null): void {
  assertDefined(element, "element");
  // element: HTMLElement
  element.textContent = "Hello!";
}
```

---

## 6.10 Control Flow Analysis

TypeScript's type narrowing is powered by **control flow analysis** — the compiler traces every code path and tracks what types are possible at each point.

```typescript
function analyze(value: string | number | null | undefined): void {
  // value: string | number | null | undefined
  
  if (value === null || value === undefined) {
    return;  // exits function
  }
  // value: string | number — null and undefined eliminated by the early return
  
  if (typeof value === "string") {
    // value: string
    value.toUpperCase();
    return;
  }
  
  // value: number — string was handled above, and we're past the return
  value.toFixed(2);
}

// Narrowing across variable assignments
function processData(data: string | null): void {
  let result: string;
  
  if (data !== null) {
    result = data.toUpperCase();  // TypeScript knows data is string
  } else {
    result = "default";
  }
  
  // result: string — TypeScript knows all paths assign a string
  console.log(result);
}

// Control flow with switch
type Color = "red" | "green" | "blue" | "unknown";

function rgb(color: Color): [number, number, number] {
  switch (color) {
    case "red":    return [255, 0, 0];
    case "green":  return [0, 255, 0];
    case "blue":   return [0, 0, 255];
    case "unknown": return [128, 128, 128];
    // TypeScript ensures all cases are handled
    // If you remove "unknown" from Color, TypeScript warns about exhaustiveness
  }
}

// Narrowing with loops
function findFirstString(items: (string | number)[]): string | undefined {
  for (const item of items) {
    if (typeof item === "string") {
      return item;  // item: string — narrowed
    }
  }
  return undefined;
}

// Discriminant narrowing in complex scenarios
function processEvents(events: Array<Action>): State {
  let state: State = { count: 0 };
  
  for (const event of events) {
    // TypeScript narrows event.type in each case
    state = reducer(state, event);
  }
  
  return state;
}
```

---

## Complete Example: A Type-Safe Event System

```typescript
// Event system using discriminated unions and type guards

type EventPayload = {
  "user:login": { userId: string; timestamp: Date; ip: string };
  "user:logout": { userId: string; timestamp: Date };
  "post:created": { postId: string; authorId: string; title: string };
  "post:deleted": { postId: string; deletedBy: string };
  "error:occurred": { code: string; message: string; stack?: string };
};

type EventType = keyof EventPayload;

// Create a discriminated union from the EventPayload map
type AppEvent = {
  [K in EventType]: { type: K; payload: EventPayload[K] };
}[EventType];

// Type-safe event handler
type EventHandler<T extends EventType> = (payload: EventPayload[T]) => void;

class EventBus {
  private handlers: Map<EventType, EventHandler<any>[]> = new Map();

  on<T extends EventType>(event: T, handler: EventHandler<T>): void {
    if (!this.handlers.has(event)) {
      this.handlers.set(event, []);
    }
    this.handlers.get(event)!.push(handler);
  }

  emit<T extends EventType>(event: T, payload: EventPayload[T]): void {
    this.handlers.get(event)?.forEach((handler) => handler(payload));
  }
}

const bus = new EventBus();

// Fully type-safe — TypeScript knows what payload each event carries
bus.on("user:login", (payload) => {
  // payload: { userId: string; timestamp: Date; ip: string }
  console.log(`User ${payload.userId} logged in from ${payload.ip}`);
});

bus.on("post:created", (payload) => {
  // payload: { postId: string; authorId: string; title: string }
  console.log(`New post: "${payload.title}" by ${payload.authorId}`);
});

bus.emit("user:login", {
  userId: "user-123",
  timestamp: new Date(),
  ip: "192.168.1.1",
});

// TYPE ERROR: wrong payload shape
// bus.emit("user:login", { userId: "user-123" });  // missing required fields
```

---

## Summary

TypeScript's advanced type system enables highly expressive type modeling. Union types (`|`) model "A or B" relationships; intersection types (`&`) combine types. Type narrowing is the compiler's ability to refine types within conditional blocks through control flow analysis. `typeof`, `instanceof`, `in`, and equality checks all perform narrowing. **Discriminated unions** — unions with a shared literal "tag" property — are the most powerful and safe narrowing pattern. User-defined type guards (`x is T`) extend narrowing to custom validation logic. Assertion functions (`asserts condition`) narrow types for code that follows them.

---

## Key Takeaways

- **Union (`|`)**: only common properties are available without narrowing
- **Intersection (`&`)**: all properties of all types are available
- TypeScript narrows types through **control flow analysis** — tracking every code path
- **`typeof`** narrows primitives; **`instanceof`** narrows class instances; **`in`** narrows by property presence
- **Discriminated unions** are the gold standard: a `kind`/`type` literal property makes narrowing precise and exhaustive
- **User-defined type guards** (`value is Type`) allow custom validation to inform the type system
- **Exhaustiveness checking** with `never` ensures all union cases are handled at compile time

---

## Practice Questions

1. What is the difference between a union type and an intersection type?
2. What is a discriminated union, and why is it more reliable than a plain union?
3. What is a type guard function? How does its return type work?
4. What does TypeScript's control flow analysis mean in practice?
5. When would you use an assertion function vs. a type guard function?
6. Why is `typeof null === "object"` in JavaScript, and how does TypeScript handle this?

---

## Exercises

**Exercise 1**: Model a payment system using discriminated unions. A `Payment` can be a credit card payment (with card number, expiry, cvv), a bank transfer (with account number, routing number), or a crypto payment (with wallet address, currency). Write a `processPayment(payment: Payment): string` function that handles all cases.

**Exercise 2**: Write a `safeGet<T>(obj: unknown, key: string): T | undefined` function that safely accesses a property of an unknown object, with proper type narrowing.

**Exercise 3**: Implement a `narrow<T>(value: unknown, guard: (v: unknown) => v is T): T` function that applies a type guard and throws if the guard fails.

**Exercise 4**: Write a complete type guard for a `Config` object that validates all required fields have the correct types, including nested objects.

---

*Next: [Chapter 7 — Generics](07-generics.md)*
