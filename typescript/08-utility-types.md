# Chapter 8 — Utility Types

> *"Utility types are TypeScript's standard library for type transformation. They let you create new types by manipulating existing ones — without rewriting type definitions."*

---

## 8.1 What Are Utility Types?

Utility types are generic types built into TypeScript that transform existing types into new ones. They are the "standard library" of the type system — reusable transformations you apply to your own types.

```typescript
interface User {
  id: number;
  name: string;
  email: string;
  password: string;
  role: "admin" | "user" | "guest";
  createdAt: Date;
}

// Instead of duplicating User with some fields optional:
interface PartialUser {
  id?: number;
  name?: string;
  // ... tedious and error-prone
}

// Use a utility type:
type PartialUser = Partial<User>;  // All properties become optional
type RequiredUser = Required<User>; // All properties become required
type ReadonlyUser = Readonly<User>; // All properties become readonly
```

All utility types discussed in this chapter are available globally — no import needed.

---

## 8.2 Partial\<T\> — Make All Properties Optional

`Partial<T>` constructs a type with all properties of `T` set to optional.

```typescript
interface Config {
  host: string;
  port: number;
  debug: boolean;
  timeout: number;
  maxConnections: number;
}

type PartialConfig = Partial<Config>;
// {
//   host?: string;
//   port?: number;
//   debug?: boolean;
//   timeout?: number;
//   maxConnections?: number;
// }

// The classic use case: updating an object
const defaultConfig: Config = {
  host: "localhost",
  port: 3000,
  debug: false,
  timeout: 5000,
  maxConnections: 100,
};

function updateConfig(current: Config, updates: Partial<Config>): Config {
  return { ...current, ...updates };
}

// Only provide what changes — everything else stays
const devConfig = updateConfig(defaultConfig, { debug: true, port: 3001 });
const prodConfig = updateConfig(defaultConfig, { host: "prod.server.com", port: 80 });

// Another use: form state
interface FormState<T> {
  values: T;
  errors: Partial<Record<keyof T, string>>;
  touched: Partial<Record<keyof T, boolean>>;
}
```

---

## 8.3 Required\<T\> — Make All Properties Required

`Required<T>` is the opposite of `Partial<T>` — it removes optionality from all properties.

```typescript
interface UserPreferences {
  theme?: "light" | "dark";
  language?: string;
  timezone?: string;
  fontSize?: number;
}

type RequiredPreferences = Required<UserPreferences>;
// {
//   theme: "light" | "dark";
//   language: string;
//   timezone: string;
//   fontSize: number;
// }

// Use case: after validation, you know all fields exist
function validatePreferences(prefs: UserPreferences): Required<UserPreferences> {
  if (!prefs.theme) throw new Error("theme required");
  if (!prefs.language) throw new Error("language required");
  if (!prefs.timezone) throw new Error("timezone required");
  if (!prefs.fontSize) throw new Error("fontSize required");
  // TypeScript still thinks the return is UserPreferences — cast needed
  return prefs as Required<UserPreferences>;
}

// After calling validatePreferences, you can use the result without optional checks
const validated = validatePreferences({ theme: "dark", language: "en", timezone: "UTC", fontSize: 14 });
console.log(validated.theme.toUpperCase());  // OK — theme is string, not string | undefined
```

---

## 8.4 Readonly\<T\> — Make All Properties Readonly

`Readonly<T>` constructs a type where all properties are `readonly`.

```typescript
interface Point {
  x: number;
  y: number;
}

type ReadonlyPoint = Readonly<Point>;
// { readonly x: number; readonly y: number }

const point: ReadonlyPoint = { x: 3, y: 4 };
// point.x = 5;  // ERROR: Cannot assign to 'x' because it is a read-only property

// Use case: function parameters that shouldn't be mutated
function translate(point: Readonly<Point>, dx: number, dy: number): Point {
  // point.x += dx;  // ERROR — can't mutate
  return { x: point.x + dx, y: point.y + dy };  // must create a new point
}

// Use case: immutable state
function reducer(state: Readonly<AppState>, action: Action): AppState {
  // Cannot mutate state directly
  switch (action.type) {
    case "INCREMENT":
      return { ...state, count: state.count + 1 };  // new object
    default:
      return state;
  }
}

// Note: Readonly is shallow — nested objects can still be mutated
interface Config {
  server: { host: string; port: number };
}
const config: Readonly<Config> = { server: { host: "localhost", port: 3000 } };
// config.server = {};           // ERROR — readonly
config.server.port = 8080;       // OK — Readonly is shallow!
```

---

## 8.5 Pick\<T, K\> — Select Specific Properties

`Pick<T, K>` constructs a type with only the properties `K` from `T`.

```typescript
interface User {
  id: number;
  name: string;
  email: string;
  password: string;
  role: "admin" | "user";
  createdAt: Date;
  lastLogin: Date;
}

// Public user info — no password
type PublicUser = Pick<User, "id" | "name" | "email" | "role">;
// { id: number; name: string; email: string; role: "admin" | "user" }

// For display in a list
type UserListItem = Pick<User, "id" | "name" | "role">;

// Type error: "phone" is not a key of User
// type BadPick = Pick<User, "id" | "phone">;  // ERROR

// Function returning only safe fields
function getPublicUser(user: User): PublicUser {
  return {
    id: user.id,
    name: user.name,
    email: user.email,
    role: user.role,
  };
  // TypeScript ensures password is not included
}

// Pick for partial update forms
type UserProfileForm = Pick<User, "name" | "email">;

function updateProfile(userId: number, form: UserProfileForm): Promise<PublicUser> {
  // Only name and email can be updated through this function
  return fetch(`/api/users/${userId}`, {
    method: "PATCH",
    body: JSON.stringify(form),
  }).then((r) => r.json());
}
```

---

## 8.6 Omit\<T, K\> — Exclude Specific Properties

`Omit<T, K>` constructs a type with all properties of `T` except those in `K`. It's the inverse of `Pick`.

```typescript
interface User {
  id: number;
  name: string;
  email: string;
  password: string;
  role: "admin" | "user";
  createdAt: Date;
}

// Omit sensitive field
type SafeUser = Omit<User, "password">;
// { id: number; name: string; email: string; role: "admin" | "user"; createdAt: Date }

// Omit auto-generated fields for creation input
type CreateUserInput = Omit<User, "id" | "createdAt">;
// { name: string; email: string; password: string; role: "admin" | "user" }

// Omit multiple fields
type UserSummary = Omit<User, "password" | "createdAt">;

// Use case: base entity with auto fields
interface Entity {
  id: string;
  createdAt: Date;
  updatedAt: Date;
}

interface Post extends Entity {
  title: string;
  content: string;
  authorId: string;
}

// For creation, omit all auto-generated fields
type CreatePost = Omit<Post, "id" | "createdAt" | "updatedAt">;

async function createPost(input: CreatePost): Promise<Post> {
  const now = new Date();
  const post: Post = {
    id: generateId(),
    createdAt: now,
    updatedAt: now,
    ...input,
  };
  return post;
}

function generateId(): string {
  return Math.random().toString(36).slice(2);
}
```

### Pick vs Omit — When to Use Which

```typescript
// Use Pick when: specifying a small subset of properties to KEEP
// Use Omit when: specifying a small subset of properties to EXCLUDE

// If you want 2 out of 10 properties: Pick (name the 2)
type Preview = Pick<LargeType, "id" | "title">;

// If you want 8 out of 10 properties: Omit (name the 2 to remove)
type CreateInput = Omit<LargeType, "id" | "createdAt">;

interface LargeType {
  id: string;
  title: string;
  content: string;
  author: string;
  tags: string[];
  status: string;
  views: number;
  likes: number;
  createdAt: Date;
  updatedAt: Date;
}
```

---

## 8.7 Record\<K, V\> — Build a Map Type

`Record<K, V>` constructs a type whose keys are type `K` and values are type `V`.

```typescript
// Record<string, number> — any string key, number value
type WordCount = Record<string, number>;
const counts: WordCount = { hello: 3, world: 5, typescript: 10 };

// Record with literal union keys — ensures all keys are present
type RolePermissions = Record<"admin" | "user" | "guest", string[]>;
const permissions: RolePermissions = {
  admin: ["read", "write", "delete"],
  user: ["read", "write"],
  guest: ["read"],
  // TypeScript errors if any key is missing
};

// HTTP status code mapping
type StatusMessages = Record<200 | 201 | 400 | 401 | 404 | 500, string>;
const statusText: StatusMessages = {
  200: "OK",
  201: "Created",
  400: "Bad Request",
  401: "Unauthorized",
  404: "Not Found",
  500: "Internal Server Error",
};

// Cache with any string key
type Cache<T> = Record<string, T | undefined>;  // undefined for cache misses

function createCache<T>(): Cache<T> {
  return {};
}

const userCache: Cache<User> = createCache();
userCache["user-123"] = { id: 1, name: "Alice", email: "a@ex.com" };
const cached = userCache["user-123"];  // User | undefined

// Building Record dynamically
function indexBy<T, K extends string | number>(
  items: T[],
  keyFn: (item: T) => K
): Record<K, T> {
  return items.reduce(
    (acc, item) => ({ ...acc, [keyFn(item)]: item }),
    {} as Record<K, T>
  );
}

const users = [
  { id: 1, name: "Alice" },
  { id: 2, name: "Bob" },
];

const userById = indexBy(users, (u) => u.id);
// { 1: { id: 1, name: "Alice" }, 2: { id: 2, name: "Bob" } }
```

---

## 8.8 Exclude\<T, U\> and Extract\<T, U\> — Filter Union Members

`Exclude<T, U>` removes from union `T` all types assignable to `U`.
`Extract<T, U>` keeps from union `T` only types assignable to `U`.

```typescript
type AllEvents = "click" | "scroll" | "mousemove" | "keydown" | "keyup" | "focus" | "blur";

// Exclude mouse events — keep keyboard and focus events
type NonMouseEvents = Exclude<AllEvents, "click" | "scroll" | "mousemove">;
// "keydown" | "keyup" | "focus" | "blur"

// Extract only keyboard events
type KeyboardEvents = Extract<AllEvents, "keydown" | "keyup">;
// "keydown" | "keyup"

// Practical: filter a union by pattern
type StringOrNumber = string | number | boolean | null | undefined;

type StringTypes = Extract<StringOrNumber, string>;      // string
type NullishTypes = Extract<StringOrNumber, null | undefined>;  // null | undefined
type NonNullTypes = Exclude<StringOrNumber, null | undefined>; // string | number | boolean

// Removing specific types from a complex union
type InputType = string | number | string[] | number[] | null;
type OnlyArrays = Extract<InputType, unknown[]>;  // string[] | number[]
type NoArrays = Exclude<InputType, unknown[]>;    // string | number | null
```

---

## 8.9 NonNullable\<T\> — Remove null and undefined

`NonNullable<T>` removes `null` and `undefined` from type `T`.

```typescript
type MaybeString = string | null | undefined;
type DefinitelyString = NonNullable<MaybeString>;  // string

type MaybeUser = User | null | undefined;
type DefinitelyUser = NonNullable<MaybeUser>;  // User

// Use case: after validation
function assertNotNull<T>(value: T | null | undefined, name: string): NonNullable<T> {
  if (value === null || value === undefined) {
    throw new Error(`${name} must not be null or undefined`);
  }
  return value as NonNullable<T>;
}

const user: User | null = findUser(1);
const validUser = assertNotNull(user, "user");
// validUser: User — null removed

// With filter — TypeScript doesn't narrow through filter automatically
const items = ["hello", null, "world", undefined, "!"];
const strings: string[] = items.filter((item): item is string => item !== null && item !== undefined);
// Need type guard in filter callback for TypeScript to understand
```

---

## 8.10 ReturnType\<T\> and Parameters\<T\> — Introspect Functions

These utility types extract information from function types.

### ReturnType\<T\>

```typescript
function getUser(): { id: number; name: string } {
  return { id: 1, name: "Alice" };
}

type UserReturnType = ReturnType<typeof getUser>;
// { id: number; name: string }

// Use case: you don't own the function type — derive it
function processResponse(response: ReturnType<typeof getUser>): void {
  console.log(response.name);
}

// With async functions
async function fetchUsers(): Promise<User[]> {
  return [];
}

type FetchUsersReturn = ReturnType<typeof fetchUsers>;  // Promise<User[]>
type AwaitedUsers = Awaited<ReturnType<typeof fetchUsers>>;  // User[]

// Use case: reuse inferred types from third-party functions
import fs from "fs";
type ReadFileResult = ReturnType<typeof fs.readFileSync>;
// Buffer — you don't need to import or repeat the type
```

### Parameters\<T\>

```typescript
function createServer(host: string, port: number, debug: boolean): void {
  // ...
}

type ServerParams = Parameters<typeof createServer>;
// [host: string, port: number, debug: boolean]  (a tuple type)

// Extract individual parameter types
type FirstParam = Parameters<typeof createServer>[0];  // string (host)
type SecondParam = Parameters<typeof createServer>[1]; // number (port)

// Use case: wrap a function with logging
function withLogging<T extends (...args: unknown[]) => unknown>(
  fn: T,
  name: string
): (...args: Parameters<T>) => ReturnType<T> {
  return (...args: Parameters<T>): ReturnType<T> => {
    console.log(`Calling ${name} with`, args);
    const result = fn(...args) as ReturnType<T>;
    console.log(`${name} returned`, result);
    return result;
  };
}

function add(a: number, b: number): number {
  return a + b;
}

const loggedAdd = withLogging(add, "add");
loggedAdd(2, 3);  // Calling add with [2, 3]; add returned 5
// loggedAdd("x", 3);  // ERROR — parameters match add's signature
```

### ConstructorParameters\<T\> and InstanceType\<T\>

```typescript
class Database {
  constructor(
    public host: string,
    public port: number,
    public dbName: string
  ) {}

  query(sql: string): Promise<unknown[]> {
    return Promise.resolve([]);
  }
}

type DbConstructorParams = ConstructorParameters<typeof Database>;
// [host: string, port: number, dbName: string]

type DbInstance = InstanceType<typeof Database>;
// Database

// Use case: factory function that mirrors a class constructor
function createDatabase(...args: ConstructorParameters<typeof Database>): Database {
  return new Database(...args);
}

createDatabase("localhost", 5432, "mydb");  // fully type-safe
```

---

## 8.11 Awaited\<T\> — Unwrap Promise Types

`Awaited<T>` recursively unwraps Promise types:

```typescript
type A = Awaited<Promise<string>>;                  // string
type B = Awaited<Promise<Promise<number>>>;          // number
type C = Awaited<string | Promise<number>>;          // string | number

// Use case: get the resolved type of an async function
async function loadUser(): Promise<User> {
  return { id: 1, name: "Alice", email: "alice@ex.com" };
}

type LoadedUser = Awaited<ReturnType<typeof loadUser>>;  // User (not Promise<User>)

// Use case: typing Promise.all results
async function main() {
  const results = await Promise.all([
    fetch("/api/users").then((r) => r.json() as Promise<User[]>),
    fetch("/api/posts").then((r) => r.json() as Promise<Post[]>),
  ]);

  type Results = Awaited<typeof results>;  // [User[], Post[]]
}
```

---

## 8.12 Combining Utility Types

Utility types are most powerful when combined:

```typescript
interface User {
  id: number;
  name: string;
  email: string;
  password: string;
  role: "admin" | "user";
  preferences: {
    theme: "light" | "dark";
    language: string;
  };
  createdAt: Date;
  updatedAt: Date;
}

// Safe public user (no password)
type PublicUser = Readonly<Omit<User, "password">>;

// Input for creating a user
type CreateUserInput = Required<Omit<User, "id" | "createdAt" | "updatedAt">>;

// Partial update (all optional, no id or timestamps)
type UpdateUserInput = Partial<Omit<User, "id" | "createdAt" | "updatedAt" | "password">>;

// Admin-only fields
type AdminOnlyFields = Pick<User, "role">;

// User search/filter
type UserFilter = Partial<Pick<User, "name" | "email" | "role">>;

// Form field errors
type UserFormErrors = Partial<Record<keyof Omit<User, "id" | "createdAt" | "updatedAt">, string>>;

// Apply these in service functions
class UserService {
  async create(input: CreateUserInput): Promise<PublicUser> {
    // implementation
    return {} as PublicUser;
  }

  async update(id: number, input: UpdateUserInput): Promise<PublicUser> {
    // Only the fields in UpdateUserInput can be changed
    return {} as PublicUser;
  }

  async search(filter: UserFilter): Promise<PublicUser[]> {
    return [];
  }
}
```

---

## Summary

Utility types are TypeScript's toolkit for type transformation. They operate on existing types to produce new ones — without duplicating type definitions. `Partial`, `Required`, and `Readonly` adjust property mutability and optionality. `Pick` and `Omit` select or exclude specific properties. `Record` builds indexed types. `Exclude` and `Extract` filter union members. `NonNullable` removes nullish types. `ReturnType`, `Parameters`, `ConstructorParameters`, and `InstanceType` introspect function and class types. `Awaited` unwraps Promise chains. The real power comes from **combining** them.

---

## Key Takeaways

- **`Partial<T>`**: all props optional — great for update functions and form state
- **`Required<T>`**: removes optionality — use after validation
- **`Readonly<T>`**: prevents mutation — use for state and config
- **`Pick<T, K>`**: keeps only specified props — use when selecting a few from many
- **`Omit<T, K>`**: removes specified props — use to exclude auto-generated fields from inputs
- **`Record<K, V>`**: typed dictionary — prefer over `{ [key: string]: T }` when keys are known
- **`ReturnType<T>` / `Parameters<T>`**: introspect existing functions — avoid repeating type definitions
- **`Awaited<T>`**: unwrap Promise — get the resolved type of async functions

---

## Practice Questions

1. What is the difference between `Partial<T>` and making all properties optional manually?
2. When would you use `Pick` vs `Omit`?
3. What is the difference between `Exclude<T, U>` and `Omit<T, K>`?
4. How does `ReturnType<typeof fn>` help when you don't control a function's type?
5. Why is `Awaited<ReturnType<typeof asyncFn>>` useful?

---

## Exercises

**Exercise 1**: Given a `Product` interface with many fields, create: `ProductPreview` (only id, name, price), `CreateProduct` (everything except id), `UpdateProduct` (all fields optional except id), and `ProductSummary` (readonly version of ProductPreview).

**Exercise 2**: Build a `typeSafeOmit<T, K extends keyof T>(obj: T, ...keys: K[]): Omit<T, K>` function that omits properties from an object and returns a properly-typed result.

**Exercise 3**: Write a function `mapObject<T extends Record<string, unknown>, U>(obj: T, fn: (value: T[keyof T], key: keyof T) => U): Record<keyof T, U>` that maps over an object's values while preserving its key type.

**Exercise 4**: Implement a `deepPartial<T>` type that makes all properties optional recursively (including nested objects).

---

*Next: [Chapter 9 — Classes](09-classes.md)*
