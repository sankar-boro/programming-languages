# Chapter 92 — Interview Preparation

> *"TypeScript interviews test two things: understanding of the type system, and practical TypeScript fluency. This chapter covers both."*

---

## Part 1 — Conceptual Questions

### Q1: What is TypeScript? How does it relate to JavaScript?

**Answer**: TypeScript is a statically typed superset of JavaScript developed by Microsoft. Every valid JavaScript program is valid TypeScript. TypeScript adds a type system on top of JavaScript — the types are checked at compile time by `tsc` and then **completely erased** from the output. The emitted code is plain JavaScript that runs anywhere JavaScript runs. TypeScript provides compile-time type safety, better tooling (autocomplete, refactoring), and documentation of function contracts through types.

---

### Q2: What is the difference between any and unknown?

**Answer**: Both accept any value, but they differ in what you can do with them:

```typescript
// any: disables type checking — TypeScript gives up
let a: any = "hello";
a.toUpperCase();  // OK
a.nonExistentMethod();  // OK (TypeScript doesn't check)
a = 42;  // OK

// unknown: type-safe alternative — must narrow before use
let u: unknown = "hello";
// u.toUpperCase();  // ERROR — must narrow first
if (typeof u === "string") {
  u.toUpperCase();  // OK — narrowed to string
}
```

Use `unknown` for values of truly unknown type (API responses, JSON.parse). Use `any` only as a last resort (legacy code, type assertions you can't avoid).

---

### Q3: Explain structural typing vs nominal typing.

**Answer**: TypeScript uses **structural typing**: two types are compatible if they have the same shape (same properties), regardless of their names. In contrast, Java/C# use **nominal typing**: types are compatible only if they have the same name or explicit inheritance.

```typescript
interface Point2D { x: number; y: number }
interface Coordinate { x: number; y: number }  // different name, same shape

function plot(p: Point2D): void {}

const coord: Coordinate = { x: 3, y: 4 };
plot(coord);  // OK! Structurally compatible — TypeScript accepts it
// In Java, this would fail — different types, no inheritance
```

This makes TypeScript flexible and enables duck typing, but can cause semantic bugs when types that should be distinct happen to have the same structure.

---

### Q4: What is the difference between interface and type?

**Answer**:

| Feature | `interface` | `type` |
|---------|------------|--------|
| Object shapes | ✓ | ✓ |
| Union/intersection | ✗ | ✓ |
| Primitive aliases | ✗ | ✓ |
| Declaration merging | ✓ | ✗ |
| `extends` keyword | ✓ | Use `&` |
| Conditional/mapped types | ✗ | ✓ |

```typescript
// interface — merges with same-name interfaces (declaration merging)
interface User { id: number }
interface User { name: string }  // OK — merges into { id: number; name: string }

// type — cannot be redeclared
type Product = { id: number };
// type Product = { name: string };  // ERROR: Duplicate identifier
```

Use `interface` for object/class shapes. Use `type` for unions, intersections, primitives, and complex type operations.

---

### Q5: What is a discriminated union? When would you use it?

**Answer**: A discriminated union is a union of types that share a common literal property (the "discriminant"). It enables TypeScript's control flow analysis to narrow types precisely in switch/if statements.

```typescript
type Shape =
  | { kind: "circle"; radius: number }
  | { kind: "rectangle"; width: number; height: number }
  | { kind: "triangle"; base: number; height: number };

function area(shape: Shape): number {
  switch (shape.kind) {
    case "circle": return Math.PI * shape.radius ** 2;
    case "rectangle": return shape.width * shape.height;
    case "triangle": return (shape.base * shape.height) / 2;
    // TypeScript errors if a case is unhandled (when combined with never check)
  }
}
```

Use discriminated unions for: API response states (loading/success/error), domain models with variant behavior, Redux/Flux action types, FSM states.

---

### Q6: What is type narrowing? Name all narrowing techniques.

**Answer**: Type narrowing is TypeScript's control flow analysis that refines a broader type to a more specific one within a code block. TypeScript tracks the type through every branch.

Narrowing techniques:
1. `typeof` — primitives: `typeof x === "string"`
2. `instanceof` — classes: `x instanceof Error`
3. `in` — property presence: `"name" in obj`
4. Equality — `x === null`, `x === "active"`
5. Truthiness — `if (x)` (removes null/undefined/false/0/"")
6. User-defined type guards — `function isUser(x): x is User`
7. Assertion functions — `function assertUser(x): asserts x is User`
8. Discriminant check — `if (shape.kind === "circle")`

---

### Q7: Explain generics. Why are they better than any?

**Answer**: Generics are type parameters that are filled in at call time. They allow writing reusable code that's still type-safe — unlike `any`, the type information is preserved.

```typescript
// With any: reusable but loses type information
function first(arr: any[]): any { return arr[0]; }
const n = first([1, 2, 3]);
n.toFixed(2);  // OK, but also n.toUpperCase() — TypeScript doesn't catch the error

// With generics: reusable AND type-safe
function first<T>(arr: T[]): T { return arr[0]; }
const m = first([1, 2, 3]);  // T = number, returns number
// m.toUpperCase();  // ERROR — TypeScript knows m is number
```

Generics are better than `any` because they preserve the caller's type information, enabling proper type checking downstream.

---

### Q8: What is the difference between Partial, Pick, Omit, and Record?

| Utility | What it does | Example |
|---------|-------------|---------|
| `Partial<T>` | All properties optional | `Partial<User>` for update inputs |
| `Required<T>` | All properties required | `Required<Options>` after validation |
| `Readonly<T>` | All properties readonly | `Readonly<Config>` for constants |
| `Pick<T, K>` | Keep only properties K | `Pick<User, "id" \| "name">` for public data |
| `Omit<T, K>` | Remove properties K | `Omit<User, "password">` for safe response |
| `Record<K, V>` | Object with keys K and values V | `Record<string, number>` for word counts |
| `Exclude<T, U>` | Remove union members assignable to U | `Exclude<Status, "error">` |
| `Extract<T, U>` | Keep union members assignable to U | `Extract<Input, string>` |
| `NonNullable<T>` | Remove null and undefined | `NonNullable<string \| null>` |
| `ReturnType<T>` | Return type of function type T | `ReturnType<typeof createUser>` |
| `Parameters<T>` | Parameter types as tuple | `Parameters<typeof fn>` |

---

### Q9: What is type erasure, and why does it matter?

**Answer**: TypeScript's types are completely removed from the generated JavaScript. No type information survives compilation — not interfaces, not type aliases, not generic type parameters.

This matters because:
1. **Runtime checks on types don't work**: `value instanceof User` fails for interfaces (no runtime representation)
2. **Type assertions are free but dangerous**: `x as User` costs nothing at runtime but provides no validation
3. **You must validate external data**: JSON.parse returns `any` — you can't check if it's a `User` without explicit property checks
4. **Types don't affect performance**: complex conditional types or mapped types have zero runtime cost

The correct response to external data is type guards (runtime validation) rather than type assertions.

---

### Q10: When would you use never?

**Answer**: `never` is the bottom type — a value of type `never` can never exist. It's used for:

1. **Functions that never return** (throw or infinite loop)
2. **Exhaustiveness checking** in switch statements
3. **Filtering types** in conditional types

```typescript
// Exhaustiveness check — if you add "hexagon" to Shape without handling it,
// TypeScript errors because shape becomes "hexagon", not never
function area(shape: Shape): number {
  switch (shape.kind) {
    case "circle": return Math.PI * shape.radius ** 2;
    // ...
    default:
      const _exhaustive: never = shape;  // ERROR if Shape has unhandled variants
      throw new Error(`Unknown shape: ${_exhaustive}`);
  }
}

// Filtering in conditional types
type Exclude<T, U> = T extends U ? never : T;
// "a" | "b" | "c" excludes "b" → "a" | never | "c" → "a" | "c"
// never is absorbed in unions
```

---

## Part 2 — Tricky Code Questions

### Q11: What is the output?

```typescript
type IsString<T> = T extends string ? "yes" : "no";

type A = IsString<string | number>;
```

**Answer**: `"yes" | "no"`

Because `T` is a naked type parameter in a conditional type, it **distributes** over the union:
- `string extends string ? "yes" : "no"` → `"yes"`
- `number extends string ? "yes" : "no"` → `"no"`
- Result: `"yes" | "no"`

To prevent distribution: `type IsString<T> = [T] extends [string] ? "yes" : "no"` → `"no"` for `string | number`.

---

### Q12: Will this compile? Why or why not?

```typescript
interface Animal { name: string }
interface Dog extends Animal { breed: string }

function takeDog(dog: Dog): void {}

const animal: Animal = { name: "Rex" };
takeDog(animal);
```

**Answer**: No, it will not compile. `Animal` is not assignable to `Dog` because `Dog` requires `breed` (which `Animal` doesn't have). The reverse works: `Dog` is assignable to `Animal` (a Dog has everything an Animal has, plus more).

However, note that excess property checking applies to object literals only. Through a variable (as shown here), it's a structural compatibility check.

---

### Q13: What type does TypeScript infer for `result`?

```typescript
const pair = [1, "hello"] as const;
const [first, second] = pair;
type First = typeof first;
type Second = typeof second;
```

**Answer**:
- `first`: `1` (number literal, because `as const` preserves literals)
- `second`: `"hello"` (string literal)
- `pair`: `readonly [1, "hello"]` (readonly tuple with literal types)

Without `as const`, `pair` would be `(string | number)[]` and `first`/`second` would both be `string | number`.

---

### Q14: Explain what happens here:

```typescript
function merge<T extends object, U extends object>(a: T, b: U): T & U {
  return { ...a, ...b };
}

const result = merge({ x: 1, y: 2 }, { y: "hello", z: true });
```

**Answer**:
- TypeScript infers `T = { x: number; y: number }` and `U = { y: string; z: boolean }`
- Return type: `T & U = { x: number; y: number } & { y: string; z: boolean }`
- For `y`: the intersection gives `number & string = never` — TypeScript's intersection of incompatible types is `never`
- At runtime: `y` will be `"hello"` (the second spread wins), but TypeScript types `y` as `never`

This demonstrates that intersection types can produce `never` for properties that conflict.

---

### Q15: What's wrong with this code?

```typescript
async function loadUsers(): Promise<User[]> {
  const users = [];
  for (const id of [1, 2, 3]) {
    users.push(await fetchUser(id));
  }
  return users;
}
```

**Answer**: TypeScript infers `users` as `any[]` because:
1. It's initialized as an empty array `[]`
2. TypeScript can't infer what type will be pushed into it
3. So it widens to `any[]`

The return type annotation `Promise<User[]>` doesn't help narrow what goes into `users`.

Fix:
```typescript
const users: User[] = [];  // explicit annotation
// Or:
return Promise.all([1, 2, 3].map(fetchUser));  // TypeScript infers correctly
```

---

## Part 3 — Coding Problems

### Problem 1: Implement a Type-Safe EventEmitter

```typescript
// Implement this:
class TypedEmitter<Events extends Record<string, unknown[]>> {
  on<K extends keyof Events>(event: K, listener: (...args: Events[K]) => void): this;
  off<K extends keyof Events>(event: K, listener: (...args: Events[K]) => void): this;
  emit<K extends keyof Events>(event: K, ...args: Events[K]): boolean;
}

// Usage:
type AppEvents = {
  login: [userId: string, timestamp: number];
  logout: [userId: string];
  error: [error: Error];
};

const emitter = new TypedEmitter<AppEvents>();
emitter.on("login", (userId, timestamp) => {
  // userId: string, timestamp: number — correct types!
});
emitter.emit("login", "user-123", Date.now());  // type-safe arguments
```

**Solution**:

```typescript
class TypedEmitter<Events extends Record<string, unknown[]>> {
  private listeners = new Map<keyof Events, Array<(...args: unknown[]) => void>>();

  on<K extends keyof Events>(event: K, listener: (...args: Events[K]) => void): this {
    if (!this.listeners.has(event)) this.listeners.set(event, []);
    this.listeners.get(event)!.push(listener as (...args: unknown[]) => void);
    return this;
  }

  off<K extends keyof Events>(event: K, listener: (...args: Events[K]) => void): this {
    const list = this.listeners.get(event);
    if (list) {
      const idx = list.indexOf(listener as (...args: unknown[]) => void);
      if (idx > -1) list.splice(idx, 1);
    }
    return this;
  }

  emit<K extends keyof Events>(event: K, ...args: Events[K]): boolean {
    const list = this.listeners.get(event);
    if (!list || list.length === 0) return false;
    list.forEach((fn) => fn(...args));
    return true;
  }
}
```

---

### Problem 2: Type-Safe Object Path Access

```typescript
// Implement DeepGet that safely accesses nested object paths
type DeepGet<T, Path extends string> =
  Path extends `${infer Key}.${infer Rest}`
    ? Key extends keyof T
      ? DeepGet<T[Key], Rest>
      : never
    : Path extends keyof T
    ? T[Path]
    : never;

// Test
type Config = {
  server: {
    host: string;
    port: number;
    ssl: { cert: string; key: string };
  };
  database: { url: string };
};

type A = DeepGet<Config, "server.host">;       // string
type B = DeepGet<Config, "server.ssl.cert">;   // string
type C = DeepGet<Config, "database.url">;      // string
type D = DeepGet<Config, "server.port">;       // number
type E = DeepGet<Config, "nonexistent">;       // never

// Runtime implementation
function deepGet<T, P extends string>(obj: T, path: P): DeepGet<T, P> {
  return path.split(".").reduce(
    (current: unknown, key: string) =>
      current !== null && typeof current === "object"
        ? (current as Record<string, unknown>)[key]
        : undefined,
    obj
  ) as DeepGet<T, P>;
}

const config: Config = {
  server: { host: "localhost", port: 3000, ssl: { cert: "cert.pem", key: "key.pem" } },
  database: { url: "postgres://localhost/db" },
};

const host = deepGet(config, "server.host");  // string
const cert = deepGet(config, "server.ssl.cert");  // string
```

---

### Problem 3: Builder Pattern with Required Fields

```typescript
// Implement a Builder that tracks which required fields have been set

type BuilderState<Required extends string, T> = {
  [K in Required]: K extends keyof T ? T[K] | undefined : never;
};

class UserBuilder<HasName extends boolean = false, HasEmail extends boolean = false> {
  private data: Partial<User> = {};

  setName(name: string): UserBuilder<true, HasEmail> {
    this.data.name = name;
    return this as unknown as UserBuilder<true, HasEmail>;
  }

  setEmail(email: string): UserBuilder<HasName, true> {
    this.data.email = email;
    return this as unknown as UserBuilder<HasName, true>;
  }

  setAge(age: number): this {
    this.data.age = age;
    return this;
  }
}

// build() is only available when both required fields are set
interface UserBuilder<HasName extends boolean, HasEmail extends boolean> {
  build: HasName extends true
    ? HasEmail extends true
      ? () => User
      : never
    : never;
}

// Usage:
// const builder = new UserBuilder();
// builder.build();  // ERROR — TypeScript knows not all required fields are set
// builder.setName("Alice").setEmail("alice@ex.com").build();  // OK
```

---

### Problem 4: Implement a Type-Safe Store

```typescript
// A type-safe state store (simplified Redux)

type Reducer<State, Action> = (state: State, action: Action) => State;
type Selector<State, Result> = (state: State) => Result;
type Listener = () => void;

class Store<State, Action> {
  private state: State;
  private listeners: Listener[] = [];

  constructor(
    private reducer: Reducer<State, Action>,
    initialState: State
  ) {
    this.state = initialState;
  }

  getState(): Readonly<State> {
    return this.state;
  }

  dispatch(action: Action): void {
    this.state = this.reducer(this.state, action);
    this.listeners.forEach((l) => l());
  }

  select<Result>(selector: Selector<State, Result>): Result {
    return selector(this.state);
  }

  subscribe(listener: Listener): () => void {
    this.listeners.push(listener);
    return () => {
      this.listeners = this.listeners.filter((l) => l !== listener);
    };
  }
}

// Usage
interface CounterState { count: number; history: number[] }
type CounterAction =
  | { type: "INCREMENT"; by: number }
  | { type: "DECREMENT"; by: number }
  | { type: "RESET" };

const counterReducer: Reducer<CounterState, CounterAction> = (state, action) => {
  switch (action.type) {
    case "INCREMENT":
      return { count: state.count + action.by, history: [...state.history, state.count] };
    case "DECREMENT":
      return { count: state.count - action.by, history: [...state.history, state.count] };
    case "RESET":
      return { count: 0, history: [] };
  }
};

const store = new Store(counterReducer, { count: 0, history: [] });

const unsubscribe = store.subscribe(() => {
  console.log("State changed:", store.getState());
});

store.dispatch({ type: "INCREMENT", by: 5 });  // TypeScript: only valid action shapes
// store.dispatch({ type: "ADD", value: 1 });   // ERROR — not a valid action

const count = store.select((state) => state.count);  // number
```

---

## Quick Flashcards

**Q**: What does `T extends never` evaluate to?
**A**: `never` — because `never` is the empty set, and distributing over it yields nothing.

**Q**: What's `keyof never`?
**A**: `string | number | symbol` — the widest possible key union.

**Q**: What's `keyof any`?
**A**: `string | number | symbol` — same as `keyof never`.

**Q**: Can you use `never` as a function parameter?
**A**: Yes, but calling such a function is impossible without a `never`-typed argument.

**Q**: What happens when you `await` a non-Promise value?
**A**: It returns the value wrapped in a resolved Promise — `await 42` works, returns `42`.

**Q**: What is `infer` and where can you use it?
**A**: `infer` captures a type inside a conditional type. It can only be used in the `extends` clause of a conditional type.

**Q**: What is the `satisfies` operator?
**A**: `satisfies` checks that a value matches a type without widening the inferred type. Preserves literal types while still type-checking.

**Q**: What is declaration merging?
**A**: The ability to declare an interface (or namespace) multiple times — TypeScript merges them into a single definition. Only interfaces support it (not type aliases).

---

*Next: [Chapter 99 — Final Project](99-final-project.md)*
