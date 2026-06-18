# Chapter 11 — Advanced Types

> *"Mapped types, conditional types, and template literal types are the meta-programming layer of TypeScript's type system. They let you derive new types from existing ones programmatically."*

---

## 11.1 Mapped Types — Transforming Every Property

A mapped type creates a new type by iterating over the keys of an existing type and applying a transformation to each property.

### Basic Syntax

```typescript
// Syntax: { [K in keyof T]: ... }
// K iterates over every key in T

type Stringify<T> = {
  [K in keyof T]: string;
};

interface User {
  id: number;
  name: string;
  active: boolean;
}

type StringifiedUser = Stringify<User>;
// { id: string; name: string; active: string }

// Making all properties optional (like Partial<T>)
type MyPartial<T> = {
  [K in keyof T]?: T[K];  // T[K] = the original value type
};

// Making all properties readonly (like Readonly<T>)
type MyReadonly<T> = {
  readonly [K in keyof T]: T[K];
};

// Making all properties nullable
type Nullable<T> = {
  [K in keyof T]: T[K] | null;
};

type NullableUser = Nullable<User>;
// { id: number | null; name: string | null; active: boolean | null }
```

### Modifiers in Mapped Types

Use `+`/`-` to add or remove `readonly` and `?`:

```typescript
// Remove readonly from all properties
type Mutable<T> = {
  -readonly [K in keyof T]: T[K];  // '-' removes 'readonly'
};

type ReadonlyPoint = { readonly x: number; readonly y: number };
type MutablePoint = Mutable<ReadonlyPoint>;
// { x: number; y: number }

// Remove optional from all properties (like Required<T>)
type MyRequired<T> = {
  [K in keyof T]-?: T[K];  // '-?' removes optionality
};

interface Options {
  timeout?: number;
  retries?: number;
  debug?: boolean;
}

type RequiredOptions = MyRequired<Options>;
// { timeout: number; retries: number; debug: boolean }

// Add optional (-? removed) and readonly together
type ImmutablePartial<T> = {
  readonly [K in keyof T]?: T[K];
};
```

### Remapping Keys with as

TypeScript 4.1+ allows renaming keys in mapped types using `as`:

```typescript
// Rename all keys to their camelCase getter form
type Getters<T> = {
  [K in keyof T as `get${Capitalize<string & K>}`]: () => T[K];
};

interface Person {
  name: string;
  age: number;
}

type PersonGetters = Getters<Person>;
// { getName: () => string; getAge: () => number }

// Filter properties by type using 'as never' to exclude
type OnlyStrings<T> = {
  [K in keyof T as T[K] extends string ? K : never]: T[K];
};

interface Mixed {
  id: number;
  name: string;
  email: string;
  active: boolean;
}

type StringFields = OnlyStrings<Mixed>;
// { name: string; email: string }

// Build event handler types
type EventHandlers<T> = {
  [K in keyof T as `on${Capitalize<string & K>}`]: (value: T[K]) => void;
};

interface FormFields {
  username: string;
  password: string;
  rememberMe: boolean;
}

type FormHandlers = EventHandlers<FormFields>;
// {
//   onUsername: (value: string) => void;
//   onPassword: (value: string) => void;
//   onRememberMe: (value: boolean) => void;
// }
```

---

## 11.2 Conditional Types — Type-Level if/else

Conditional types choose between two types based on whether a type satisfies a constraint.

### Basic Syntax

```typescript
// T extends U ? TypeIfTrue : TypeIfFalse
type IsString<T> = T extends string ? "yes" : "no";

type A = IsString<string>;    // "yes"
type B = IsString<number>;    // "no"
type C = IsString<"hello">;   // "yes" — "hello" extends string

// Practical: IsArray
type IsArray<T> = T extends unknown[] ? true : false;

type D = IsArray<string[]>;  // true
type E = IsArray<string>;    // false

// Practical: Unwrap a Promise
type Awaited2<T> = T extends Promise<infer U> ? U : T;

type F = Awaited2<Promise<string>>;           // string
type G = Awaited2<Promise<Promise<number>>>;  // Promise<number> (one level)
type H = Awaited2<boolean>;                   // boolean (not a Promise)
```

### Distributive Conditional Types

When the checked type is a naked type parameter, the conditional type distributes over union members:

```typescript
// IsString distributes over union
type IsString<T> = T extends string ? "yes" : "no";

type I = IsString<string | number>;
// "yes" | "no" — distributed! Applied to string ("yes") and number ("no") separately

// This is how Exclude and Extract work:
type MyExclude<T, U> = T extends U ? never : T;
// For T = "a" | "b" | "c", U = "b":
// "a" extends "b" ? never : "a" → "a"
// "b" extends "b" ? never : "b" → never
// "c" extends "b" ? never : "c" → "c"
// Result: "a" | "c"

type MyExtract<T, U> = T extends U ? T : never;

type J = MyExclude<"a" | "b" | "c", "b">;  // "a" | "c"
type K = MyExtract<"a" | "b" | "c", "a" | "c">;  // "a" | "c"

// Prevent distribution by wrapping in a tuple
type IsStringNonDistrib<T> = [T] extends [string] ? "yes" : "no";
type L = IsStringNonDistrib<string | number>;  // "no" — not distributed
```

---

## 11.3 The infer Keyword — Pattern Matching on Types

`infer` captures a type from within a conditional type. It's type-level pattern matching.

```typescript
// Infer the return type of a function
type ReturnType2<T> = T extends (...args: unknown[]) => infer R ? R : never;

function greet(name: string): string { return `Hello, ${name}`; }

type GreetReturn = ReturnType2<typeof greet>;  // string

// Infer the first element of a tuple
type Head<T extends unknown[]> = T extends [infer First, ...unknown[]] ? First : never;

type H1 = Head<[string, number, boolean]>;  // string
type H2 = Head<[number]>;                  // number
type H3 = Head<[]>;                        // never

// Infer the rest of a tuple
type Tail<T extends unknown[]> = T extends [unknown, ...infer Rest] ? Rest : never;

type T1 = Tail<[string, number, boolean]>;  // [number, boolean]
type T2 = Tail<[string]>;                  // []

// Infer the element type of an array or Promise
type ElementType<T> = T extends Array<infer E> ? E : never;
type Awaited3<T> = T extends Promise<infer U> ? Awaited3<U> : T;  // recursive!

type E1 = ElementType<string[]>;              // string
type E2 = ElementType<Array<{ id: number }>>;  // { id: number }
type E3 = Awaited3<Promise<Promise<string>>>;  // string — recursive unwrap!

// Infer parameter types
type Parameters2<T> = T extends (...args: infer P) => unknown ? P : never;

function createUser(name: string, age: number, role: "admin" | "user"): void {}
type CreateUserParams = Parameters2<typeof createUser>;
// [name: string, age: number, role: "admin" | "user"]

// Infer the last element of a union
type LastOf<T> = (
  T extends unknown ? (x: () => T) => void : never
) extends (x: infer L) => void ? ReturnType<L> : never;
```

---

## 11.4 Template Literal Types — String Manipulation at the Type Level

Template literal types use backtick syntax to construct string literal types from existing types.

### Basic Usage

```typescript
type Greeting = `Hello, ${string}`;
// Matches any string starting with "Hello, "

const g1: Greeting = "Hello, Alice";   // OK
const g2: Greeting = "Hello, World";   // OK
// const g3: Greeting = "Hi, Alice";   // ERROR

// Combining union types
type Direction = "left" | "right" | "up" | "down";
type CSSProperty = `margin-${Direction}` | `padding-${Direction}`;
// "margin-left" | "margin-right" | "margin-up" | "margin-down" |
// "padding-left" | "padding-right" | "padding-up" | "padding-down"

// Two unions = cartesian product
type Size = "sm" | "md" | "lg";
type Color = "red" | "blue" | "green";
type ColoredButton = `${Color}-${Size}`;
// "red-sm" | "red-md" | "red-lg" | "blue-sm" | ... (9 combinations)
```

### Template Literals with Utility String Types

TypeScript provides built-in string manipulation types:
- `Uppercase<S>` — `"hello"` → `"HELLO"`
- `Lowercase<S>` — `"HELLO"` → `"hello"`
- `Capitalize<S>` — `"hello"` → `"Hello"`
- `Uncapitalize<S>` — `"Hello"` → `"hello"`

```typescript
type Shout<T extends string> = Uppercase<T>;
type Whisper<T extends string> = Lowercase<T>;

type Shouted = Shout<"hello world">;  // "HELLO WORLD"
type Whispered = Whisper<"LOUD TEXT">;  // "loud text"

// Building event types
type EventName<T extends string> = `on${Capitalize<T>}`;

type ClickEvent = EventName<"click">;   // "onClick"
type HoverEvent = EventName<"hover">;   // "onHover"
type SubmitEvent = EventName<"submit">; // "onSubmit"

// Deriving setters from property names
type Setter<T extends string> = `set${Capitalize<T>}`;

interface FormState {
  username: string;
  email: string;
  age: number;
}

type Setters = {
  [K in keyof FormState as Setter<string & K>]: (value: FormState[K]) => void;
};
// {
//   setUsername: (value: string) => void;
//   setEmail: (value: string) => void;
//   setAge: (value: number) => void;
// }
```

### Parsing String Literal Types

```typescript
// Extract route parameters from a URL pattern string
type ExtractRouteParams<S extends string> =
  S extends `${string}:${infer Param}/${infer Rest}`
    ? { [K in Param | keyof ExtractRouteParams<`/${Rest}`>]: string }
    : S extends `${string}:${infer Param}`
    ? { [K in Param]: string }
    : {};

type Params = ExtractRouteParams<"/users/:userId/posts/:postId">;
// { userId: string; postId: string }

type Params2 = ExtractRouteParams<"/users/:id">;
// { id: string }

type Params3 = ExtractRouteParams<"/health">;
// {} — no params

// SQL query typing (simplified)
type SqlSelect<
  Table extends string,
  Columns extends string
> = `SELECT ${Columns} FROM ${Table}`;

type Query = SqlSelect<"users", "id, name, email">;
// "SELECT id, name, email FROM users"
```

---

## 11.5 Recursive Types

Types can reference themselves to describe recursive data structures.

```typescript
// JSON value — recursive
type JSONValue =
  | string
  | number
  | boolean
  | null
  | JSONValue[]
  | { [key: string]: JSONValue };

const json: JSONValue = {
  name: "Alice",
  age: 30,
  tags: ["admin", "user"],
  address: {
    city: "Springfield",
    coords: [40.7128, -74.0060],
  },
};

// Recursive tree structure
interface TreeNode<T> {
  value: T;
  children: TreeNode<T>[];
}

const tree: TreeNode<string> = {
  value: "root",
  children: [
    {
      value: "child1",
      children: [
        { value: "grandchild1", children: [] },
        { value: "grandchild2", children: [] },
      ],
    },
    {
      value: "child2",
      children: [],
    },
  ],
};

// Recursive function over recursive type
function mapTree<T, U>(node: TreeNode<T>, fn: (value: T) => U): TreeNode<U> {
  return {
    value: fn(node.value),
    children: node.children.map((child) => mapTree(child, fn)),
  };
}

const numbered = mapTree(tree, (s) => s.length);

// Deep recursive type
type DeepPartial<T> = T extends object
  ? { [K in keyof T]?: DeepPartial<T[K]> }
  : T;

interface Config {
  server: {
    host: string;
    port: number;
    ssl: { cert: string; key: string };
  };
  database: {
    url: string;
    poolSize: number;
  };
}

type PartialConfig = DeepPartial<Config>;
// All nested properties become optional
const config: PartialConfig = {
  server: { host: "localhost" },  // port and ssl are optional
};
```

---

## 11.6 Variadic Tuple Types

TypeScript 4.0+ supports spreading tuple types, enabling powerful generic tuple manipulation.

```typescript
// Spread in tuples
type Prepend<T, Tuple extends unknown[]> = [T, ...Tuple];
type Append<Tuple extends unknown[], T> = [...Tuple, T];
type Concat<A extends unknown[], B extends unknown[]> = [...A, ...B];

type P = Prepend<string, [number, boolean]>;  // [string, number, boolean]
type A = Append<[string, number], boolean>;   // [string, number, boolean]
type C = Concat<[string, number], [boolean, null]>;  // [string, number, boolean, null]

// Function parameters with variadic tuples
function call<F extends (...args: unknown[]) => unknown>(
  fn: F,
  ...args: Parameters<F>
): ReturnType<F> {
  return fn(...args) as ReturnType<F>;
}

function add(a: number, b: number): number { return a + b; }

const result = call(add, 2, 3);  // TypeScript knows: number
// call(add, "x", 3);            // ERROR — "x" is not number

// Typed pipeline
type Pipe<Fns extends ((arg: unknown) => unknown)[]> =
  Fns extends []
    ? never
    : Fns extends [infer F]
    ? F
    : Fns extends [infer F, ...infer Rest extends ((arg: unknown) => unknown)[]]
    ? (arg: Parameters<F & Function>[0]) => ReturnType<Pipe<Rest> & Function>
    : never;
```

---

## 11.7 Putting It All Together — Advanced Patterns

### A Type-Safe Object Builder

```typescript
// Builder pattern using conditional types
class Builder<T extends Record<string, unknown>> {
  private data: Partial<T> = {};

  set<K extends keyof T>(key: K, value: T[K]): Builder<T> {
    this.data[key] = value;
    return this;
  }

  build(): T {
    return this.data as T;
  }
}

// Type-safe: only keys of T can be set
const user = new Builder<{ name: string; age: number; role: string }>()
  .set("name", "Alice")
  .set("age", 30)
  .set("role", "admin")
  // .set("phone", "555-1234")  // ERROR: "phone" not in type
  .build();
```

### Generating Type-Safe Validators

```typescript
// Runtime validator that mirrors compile-time types
type Validator<T> = {
  [K in keyof T]: (value: unknown) => value is T[K];
};

interface User {
  id: number;
  name: string;
  email: string;
}

const userValidator: Validator<User> = {
  id: (v): v is number => typeof v === "number",
  name: (v): v is string => typeof v === "string",
  email: (v): v is string => typeof v === "string" && v.includes("@"),
};

function validate<T>(data: Record<string, unknown>, validator: Validator<T>): data is T {
  return (Object.keys(validator) as Array<keyof T>).every((key) => {
    return validator[key](data[key as string]);
  });
}

const rawData = { id: 1, name: "Alice", email: "alice@ex.com" };
if (validate(rawData, userValidator)) {
  // rawData: User — narrowed by the validator
  console.log(rawData.name);  // "Alice" — TypeScript knows: string
}
```

### Type-Safe CSS-in-JS

```typescript
// A mini type-safe style object builder
type CSSUnit = "px" | "em" | "rem" | "%" | "vh" | "vw";
type CSSValue = number | `${number}${CSSUnit}` | "auto" | "inherit";

type StyleProperty =
  | "width" | "height" | "margin" | "padding"
  | "marginTop" | "marginBottom" | "marginLeft" | "marginRight"
  | "paddingTop" | "paddingBottom" | "paddingLeft" | "paddingRight"
  | "fontSize" | "lineHeight"
  | "color" | "backgroundColor"
  | "display" | "flexDirection" | "alignItems" | "justifyContent";

type StyleValue<P extends StyleProperty> =
  P extends "display" ? "block" | "inline" | "flex" | "grid" | "none" :
  P extends "flexDirection" ? "row" | "column" | "row-reverse" | "column-reverse" :
  P extends "alignItems" | "justifyContent" ?
    "flex-start" | "flex-end" | "center" | "space-between" | "space-around" :
  P extends "color" | "backgroundColor" ? string :
  CSSValue;

type StyleObject = {
  [P in StyleProperty]?: StyleValue<P>;
};

function style(styles: StyleObject): StyleObject {
  return styles;
}

const buttonStyle = style({
  display: "flex",
  alignItems: "center",
  padding: "8px",
  backgroundColor: "#007bff",
  color: "white",
  // display: "block",  // allowed — it's valid
  // flexDirection: "diagonal",  // ERROR — not a valid value
});
```

---

## Summary

Mapped types transform every property of a type — they're how TypeScript's built-in utility types like `Partial`, `Required`, and `Readonly` are implemented. Conditional types express type-level if/else logic, and distribute over unions when the checked type is a bare type parameter. The `infer` keyword pattern-matches on types, extracting sub-types. Template literal types construct string types via interpolation, enabling typed event names, CSS properties, and route parameters. Recursive types model recursive data structures like JSON and trees. Together, these features form TypeScript's type meta-programming layer.

---

## Key Takeaways

- **Mapped types** iterate over keys: `{ [K in keyof T]: ... }` — use `-readonly` and `-?` to remove modifiers
- **`as` in mapped types** allows key renaming and filtering
- **Conditional types** are type-level `if/else`: `T extends U ? A : B`
- **Distributive conditional types** automatically apply over each union member when `T` is a naked type parameter
- **`infer`** captures a type during pattern matching — it's how `ReturnType`, `Parameters`, and `Awaited` work
- **Template literal types** compose string types from string, union, and built-in string-manipulation types
- **Recursive types** model recursive data — but TypeScript limits recursion depth

---

## Practice Questions

1. What is a mapped type? Write the definition of `Partial<T>` using only mapped type syntax.
2. What is distributivity in conditional types? How do you prevent it?
3. What does `infer` do in a conditional type?
4. What are template literal types? Give a real example where they're useful.
5. What is the result of `"a" | "b" extends "a" ? true : false`? Explain why.

---

## Exercises

**Exercise 1**: Implement `DeepRequired<T>` — makes all properties required recursively, including nested objects.

**Exercise 2**: Write `FlattenObject<T>` that transforms `{ a: { b: string; c: number }; d: boolean }` into `{ "a.b": string; "a.c": number; d: boolean }` using template literal types.

**Exercise 3**: Implement `Zip<A extends unknown[], B extends unknown[]>` — a type-level zip that produces `[[A[0], B[0]], [A[1], B[1]], ...]` from two tuples.

**Exercise 4**: Create a `TypedRouter<Routes>` type that takes a map of route patterns (like `/users/:id`) and produces typed handler functions where the parameter types are inferred from the route string.

---

*Next: [Chapter 12 — Working with JavaScript](12-interop.md)*
