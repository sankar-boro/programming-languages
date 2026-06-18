# Chapter 90 — Best Practices

> *"Good TypeScript is not about using the most advanced features. It's about writing types that are honest, precise, and maintainable — types that help your team rather than fight them."*

---

## 1. Always Use strict: true

```json
// tsconfig.json
{
  "compilerOptions": {
    "strict": true  // non-negotiable
  }
}
```

`strict: true` catches real bugs at compile time. The pain of fixing strict errors upfront is far less than debugging production runtime errors. Never ship without it.

---

## 2. Prefer unknown Over any for External Data

```typescript
// BAD
async function fetchUser(id: number): Promise<any> {
  const response = await fetch(`/api/users/${id}`);
  return response.json();  // any — type safety gone
}

// GOOD
async function fetchUser(id: number): Promise<User> {
  const response = await fetch(`/api/users/${id}`);
  const data: unknown = await response.json();

  if (!isUser(data)) throw new Error("Invalid user data");
  return data;  // User — validated
}
```

`unknown` forces you to validate. `any` silently bypasses the type system.

---

## 3. Use Type Inference — Don't Over-Annotate

```typescript
// BAD — redundant annotations
const name: string = "Alice";
const nums: number[] = [1, 2, 3];
const fn: (x: number) => number = (x) => x * 2;

// GOOD — let TypeScript infer
const name = "Alice";         // string
const nums = [1, 2, 3];       // number[]
const fn = (x: number) => x * 2;  // (x: number) => number

// ANNOTATE explicit types for:
// 1. Function parameters (cannot be inferred)
// 2. Public API return types (documentation + early error detection)
// 3. Complex generic types that inference gets wrong
export function processUsers(users: User[]): ProcessedUser[] {
  return users.map(transform);
}
```

---

## 4. Prefer Interfaces for Object Shapes, Types for Unions

```typescript
// Object shapes — use interface
interface User {
  id: number;
  name: string;
}

// Unions, intersections, primitives — use type alias
type ID = string | number;
type Status = "active" | "inactive" | "pending";
type UserOrAdmin = User | Admin;

// Consistency matters more than which you choose
// — pick one convention and stick to it
```

---

## 5. Use Discriminated Unions Instead of Boolean Flags

```typescript
// BAD — boolean flags lead to invalid state
interface Request {
  isLoading: boolean;
  isError: boolean;
  data?: User;
  error?: string;
  // What if isLoading AND isError are both true? Invalid!
}

// GOOD — discriminated union: only valid states representable
type RequestState<T> =
  | { status: "idle" }
  | { status: "loading" }
  | { status: "success"; data: T }
  | { status: "error"; error: string };
```

---

## 6. Never Use Enums — Use Literal Unions or as const

```typescript
// BAD — regular enum has surprising runtime behavior
enum Direction {
  Up,    // 0
  Down,  // 1
  Left,  // 2
  Right, // 3
}
// Direction[0] === "Up" — reverse mapping! Unexpected.

// GOOD — string literal union
type Direction = "up" | "down" | "left" | "right";

// GOOD — const object (has runtime access AND type safety)
const Direction = {
  Up: "up",
  Down: "down",
  Left: "left",
  Right: "right",
} as const;
type Direction = typeof Direction[keyof typeof Direction];
// type Direction = "up" | "down" | "left" | "right"
```

---

## 7. Keep Types DRY — Derive from Existing Types

```typescript
// BAD — duplicate type definition
interface User {
  id: string;
  name: string;
  email: string;
  role: "admin" | "user";
  createdAt: Date;
}

// CreateUserInput mirrors User — maintained separately (diverges over time)
interface CreateUserInput {
  name: string;
  email: string;
  role: "admin" | "user";
}

// GOOD — derive from the source of truth
type CreateUserInput = Omit<User, "id" | "createdAt">;
type PublicUser = Omit<User, "password">;
type UserSummary = Pick<User, "id" | "name" | "role">;
type UpdateUserInput = Partial<Omit<User, "id" | "createdAt">>;
```

---

## 8. Use Readonly for Data That Shouldn't Change

```typescript
// BAD — no protection against mutation
interface Config {
  host: string;
  port: number;
}

function applyConfig(config: Config): void {
  config.host = "changed";  // silent bug!
}

// GOOD — readonly prevents mutation
interface Config {
  readonly host: string;
  readonly port: number;
}

// Or use Readonly<T> at usage sites
function applyConfig(config: Readonly<Config>): void {
  // config.host = "changed";  // ERROR
}

// Readonly arrays
function first(arr: readonly number[]): number | undefined {
  // arr.push(1);  // ERROR — readonly
  return arr[0];
}
```

---

## 9. Write Type Guards Instead of Assertions

```typescript
// BAD — type assertion bypasses type checking
function process(data: unknown): void {
  const user = data as User;  // no validation — runtime bomb
  console.log(user.name);     // might crash
}

// GOOD — type guard validates AND narrows
function isUser(value: unknown): value is User {
  return (
    typeof value === "object" &&
    value !== null &&
    typeof (value as Record<string, unknown>).id === "number" &&
    typeof (value as Record<string, unknown>).name === "string"
  );
}

function process(data: unknown): void {
  if (!isUser(data)) throw new Error("Expected User");
  console.log(data.name);  // safe — validated
}
```

Reserve `as` for cases where you genuinely have more information than TypeScript:

```typescript
// Acceptable uses of 'as':
const el = document.getElementById("app") as HTMLInputElement;  // you know it's an input
const initial = {} as Config;  // building the object incrementally
const narrowed = value as "admin" | "user";  // narrowing within a known union
```

---

## 10. Use import type for Type-Only Imports

```typescript
// BAD — imports runtime values just for types
import { User, Config, ApiError } from "./types";  // might import values

// GOOD — explicit type-only imports
import type { User } from "./types";
import type { Config } from "./config";

// Mixed — import type inline
import { createUser, type User } from "./user-service";
```

---

## 11. Prefer Composition Over Deep Inheritance

```typescript
// BAD — deep class hierarchy
class Animal { }
class Mammal extends Animal { }
class Domestic extends Mammal { }
class Dog extends Domestic { }
class TrainedDog extends Dog { }

// GOOD — compose behaviors
interface Loggable { log(): void }
interface Persistable { save(): Promise<void> }
interface Trackable { track(event: string): void }

class User implements Loggable, Persistable {
  log(): void { console.log(this); }
  async save(): Promise<void> { /* ... */ }
}
```

---

## 12. Model Domain Errors with Discriminated Unions

```typescript
// BAD — throw Error with string message
async function getUser(id: string): Promise<User> {
  if (!id) throw new Error("ID required");
  const user = await db.find(id);
  if (!user) throw new Error("Not found");
  return user;
}
// Caller has no idea what errors to expect — must read implementation

// GOOD — explicit error types
type GetUserError =
  | { code: "INVALID_ID"; id: string }
  | { code: "NOT_FOUND"; id: string }
  | { code: "DB_ERROR"; cause: Error };

type GetUserResult = { success: true; user: User } | { success: false; error: GetUserError };

async function getUser(id: string): Promise<GetUserResult> {
  if (!id) return { success: false, error: { code: "INVALID_ID", id } };
  try {
    const user = await db.find(id);
    if (!user) return { success: false, error: { code: "NOT_FOUND", id } };
    return { success: true, user };
  } catch (cause) {
    return { success: false, error: { code: "DB_ERROR", cause: cause as Error } };
  }
}
```

---

## 13. Don't Suppress Errors with // @ts-ignore — Use // @ts-expect-error with Comments

```typescript
// BAD — silent suppression, no documentation
// @ts-ignore
const result = legacyFunction(arg);

// BETTER — documents why, fails if error disappears
// @ts-expect-error LegacyAPI returns string but we know it's actually number in practice
const result: number = legacyFunction(arg);

// BEST — fix the root cause
const result = legacyFunction(arg) as unknown as number;  // at least visible
// Or write a type declaration for legacyFunction
```

---

## 14. Use satisfies for Type-Checked Literals That Preserve Narrowing

```typescript
// Problem: type annotation widens literal types
const config: Record<string, string> = {
  host: "localhost",
  port: "3000",
};
config.host;  // string — narrowed away

// satisfies: validates without widening
const config2 = {
  host: "localhost",
  port: "3000",
} satisfies Record<string, string>;

config2.host;  // "localhost" — literal type preserved!
// But TypeScript also checked: { host: string; port: string } ✓

// Another example
type Colors = "red" | "green" | "blue";
const palette = {
  red: [255, 0, 0],
  green: "#00ff00",
  blue: [0, 0, 255],
} satisfies Record<Colors, string | number[]>;

palette.red;    // [number, number, number] — array type preserved
palette.green;  // string — string type preserved
```

---

## 15. Write Explicit Return Types for Public Functions

```typescript
// Private implementation — inference is fine
function _helper(x: number) {
  return x * 2 + 1;
}

// Exported/public function — explicit return type
// 1. Documents the contract
// 2. Errors at the function definition, not the call site
// 3. Prevents accidental return type changes
export function processData(items: DataItem[]): ProcessedResult[] {
  return items.map(transform);
}

// Especially important for async functions
export async function fetchUser(id: string): Promise<User | null> {
  const raw = await db.findById(id);
  return raw ? mapToUser(raw) : null;
}
```

---

## 16. Validate at System Boundaries

```typescript
// System boundaries: where unknown data enters your system
// - HTTP request bodies
// - JSON.parse
// - localStorage / sessionStorage
// - External API responses
// - Environment variables

// Validate ONCE at the boundary, then trust the type inside
function parseRequestBody(raw: unknown): CreateUserInput {
  if (!isCreateUserInput(raw)) {
    throw new BadRequestError("Invalid request body");
  }
  return raw;  // trusted
}

// Inside the system: trust your types
function handleCreateUser(input: CreateUserInput): void {
  // No need to re-validate — it was validated at the boundary
  userService.create(input.name, input.email);
}
```

---

## 17. Prefer const Assertions for Configuration Objects

```typescript
// BAD — properties are mutable, types are widened
const ROUTES = {
  home: "/",
  about: "/about",
  users: "/users",
};
ROUTES.home;  // string — not "/"

// GOOD — immutable, literal types preserved
const ROUTES = {
  home: "/",
  about: "/about",
  users: "/users",
} as const;
ROUTES.home;  // "/"

type Route = typeof ROUTES[keyof typeof ROUTES];
// "/" | "/about" | "/users"
```

---

## 18. Enable noUncheckedIndexedAccess for Safer Array Access

```json
// tsconfig.json
{
  "compilerOptions": {
    "noUncheckedIndexedAccess": true
  }
}
```

```typescript
// Without noUncheckedIndexedAccess (default):
const arr = [1, 2, 3];
const first: number = arr[0];  // TypeScript trusts you — might be undefined at runtime

// With noUncheckedIndexedAccess:
const arr2 = [1, 2, 3];
const first2: number | undefined = arr2[0];  // TypeScript forces you to handle undefined
if (first2 !== undefined) {
  console.log(first2.toFixed(2));  // safe
}
```

---

## 19. Co-locate Types with Their Usage

```typescript
// BAD — types in a separate global types.ts file (often becomes a dumping ground)
// src/types.ts: 500 lines of everything

// GOOD — types live with the code that uses them
// user.service.ts
interface User { id: string; name: string; }
interface CreateUserInput { name: string; email: string; }
export function createUser(input: CreateUserInput): Promise<User> { /* ... */ }

// Only move types to a shared file when genuinely used in multiple places
```

---

## 20. Use Exhaustive Switch Statements

```typescript
type Shape = "circle" | "rectangle" | "triangle";

function area(shape: Shape, ...dims: number[]): number {
  switch (shape) {
    case "circle": return Math.PI * dims[0] ** 2;
    case "rectangle": return dims[0] * dims[1];
    case "triangle": return (dims[0] * dims[1]) / 2;
    default:
      // If you add "hexagon" to Shape, TypeScript errors here
      const exhaustive: never = shape;
      throw new Error(`Unhandled shape: ${exhaustive}`);
  }
}
```

---

## Quick Reference — The Rules

| Rule | Prefer | Avoid |
|------|--------|-------|
| Null safety | `unknown` + type guards | `any`, `!` operator |
| Annotation | Annotate params, infer the rest | Annotate everything |
| Object shapes | `interface` (supports merging) | Inline object types |
| Union types | `type` alias | `interface` (can't express unions) |
| State | Discriminated union | Boolean flags |
| Enums | Literal union or `as const` | Regular `enum` |
| External data | Validate at boundary | `as SomeType` |
| Import types | `import type` | Plain `import` for types |
| Configuration | `as const` | Mutable objects |
| Errors | Typed discriminated union | `throw new Error(string)` |

---

*Next: [Chapter 91 — Common Pitfalls](91-common-pitfalls.md)*
