# Chapter 10 — Modules and Namespaces

> *"TypeScript modules are ES modules with types. Understanding module resolution is the key to eliminating mysterious 'cannot find module' errors and structuring large codebases."*

---

## 10.1 Why Modules?

Before modules, all JavaScript code shared a single global scope. Every variable was global, name collisions were inevitable, and organizing code was a nightmare.

Modules solve this by giving each file its own scope. A file is a module if it contains any `import` or `export` statement.

```typescript
// Without modules — everything is global (bad):
var userName = "Alice";  // global — any file can overwrite this

// With modules — each file has its own scope:
// user.ts
export const userName = "Alice";  // only accessible via import

// main.ts
import { userName } from "./user";
```

In TypeScript, module behavior is controlled by `"module"` in tsconfig.json:
- `"commonjs"` — for Node.js (`require`/`module.exports`)
- `"es2020"` or `"esnext"` — native ES modules (`import`/`export`)
- `"node16"` or `"nodenext"` — modern Node.js with ESM support

---

## 10.2 Named Exports and Imports

The primary way to share code between files.

### Exporting

```typescript
// math.ts — multiple named exports

export const PI = 3.14159265358979;

export function add(a: number, b: number): number {
  return a + b;
}

export function multiply(a: number, b: number): number {
  return a * b;
}

// Export a type
export interface Vector2D {
  x: number;
  y: number;
}

export class Matrix {
  constructor(private data: number[][]) {}

  get rows(): number { return this.data.length; }
  get cols(): number { return this.data[0]?.length ?? 0; }
}

// Export at the bottom — also valid
const GOLDEN_RATIO = 1.61803398875;
function fibonacci(n: number): number {
  if (n <= 1) return n;
  return fibonacci(n - 1) + fibonacci(n - 2);
}

export { GOLDEN_RATIO, fibonacci };

// Re-export with rename
export { fibonacci as fib };
```

### Importing

```typescript
// main.ts

// Import specific names
import { add, multiply, PI } from "./math";

// Import with rename (alias)
import { multiply as mul } from "./math";

// Import type (only imports the type — no runtime code)
import type { Vector2D } from "./math";

// Import all into a namespace
import * as MathLib from "./math";

console.log(MathLib.PI);           // 3.14159...
console.log(MathLib.add(2, 3));    // 5
console.log(MathLib.GOLDEN_RATIO); // 1.618...

// Use imported type (erased at runtime)
const v: Vector2D = { x: 1, y: 2 };
```

### import type — Type-Only Imports

```typescript
// Explicitly import only the type — zero runtime cost
import type { User } from "./types";
import type { Config } from "./config";

// Mixing value and type imports
import { createUser, type User } from "./user-service";  // inline type import

function processUser(user: User): void {
  // User is only a type — erased at runtime
  console.log(user.name);
}
```

`import type` is important for:
1. Performance — bundlers can confidently eliminate it
2. Clarity — makes it obvious the import is type-only
3. Required in some configurations (verbatimModuleSyntax)

---

## 10.3 Default Exports and Imports

A module can have at most one `default` export. Default exports are often used for the "main thing" a module provides.

```typescript
// user-service.ts — default export

export interface User {
  id: number;
  name: string;
  email: string;
}

// The default export — the main value of this module
export default class UserService {
  private users: User[] = [];

  create(name: string, email: string): User {
    const user: User = { id: Date.now(), name, email };
    this.users.push(user);
    return user;
  }

  findById(id: number): User | undefined {
    return this.users.find((u) => u.id === id);
  }

  getAll(): User[] {
    return [...this.users];
  }
}
```

```typescript
// main.ts — importing default export

// Default import: no curly braces, any name you choose
import UserService from "./user-service";
import { type User } from "./user-service";  // named import alongside default

const service = new UserService();
const alice: User = service.create("Alice", "alice@example.com");
```

### Default Export Considerations

```typescript
// Default exports can be expressions, not just class/function declarations
export default 42;
export default { key: "value" };
export default function() { return "unnamed function"; }  // anonymous — bad for debugging
export default function namedFn() { return "named"; }  // better

// Common pattern: named declaration + export default separately
function processData(data: unknown[]): string {
  return JSON.stringify(data);
}
export default processData;

// Many teams prefer named exports for everything:
// - Better refactoring support (rename works across files)
// - Better tree-shaking
// - Clearer at the import site
```

---

## 10.4 Re-Exports — Building a Public API

Re-exports let you create a public interface for a module or a whole folder.

```typescript
// src/utils/string.ts
export function capitalize(s: string): string {
  return s.charAt(0).toUpperCase() + s.slice(1);
}
export function trim(s: string): string {
  return s.trim();
}

// src/utils/array.ts
export function unique<T>(arr: T[]): T[] {
  return [...new Set(arr)];
}
export function flatten<T>(arr: T[][]): T[] {
  return arr.flat();
}

// src/utils/number.ts
export function clamp(n: number, min: number, max: number): number {
  return Math.min(Math.max(n, min), max);
}

// src/utils/index.ts — barrel file: re-exports everything
export { capitalize, trim } from "./string";
export { unique, flatten } from "./array";
export { clamp } from "./number";

// Re-export with rename
export { capitalize as cap } from "./string";

// Re-export all from a module
export * from "./string";
export * from "./array";

// Re-export a type
export type { SomeType } from "./types";

// src/main.ts — clean import from the barrel
import { capitalize, unique, clamp } from "./utils";
// Instead of:
// import { capitalize } from "./utils/string";
// import { unique } from "./utils/array";
// import { clamp } from "./utils/number";
```

---

## 10.5 Module Resolution

TypeScript needs to find the file for every import. The resolution algorithm depends on your `moduleResolution` setting.

### `node` Resolution (Classic)

```typescript
import { foo } from "./relative/path";
// 1. ./relative/path.ts
// 2. ./relative/path.tsx
// 3. ./relative/path.d.ts
// 4. ./relative/path/index.ts
// 5. ./relative/path/index.tsx
// 6. ./relative/path/index.d.ts

import { bar } from "some-package";
// Looks in node_modules/some-package:
// 1. package.json "main" field
// 2. index.js -> index.d.ts
// 3. node_modules/@types/some-package
```

### `node16` / `nodenext` Resolution (Modern)

Matches Node.js's native ESM resolution. Requires file extensions in relative imports:

```typescript
// tsconfig.json
// "moduleResolution": "node16"

// Must include .js extension (TypeScript resolves .ts but you write .js)
import { foo } from "./foo.js";      // resolves to foo.ts
import { bar } from "./bar/index.js"; // resolves to bar/index.ts
```

### Path Aliases

```json
// tsconfig.json
{
  "compilerOptions": {
    "baseUrl": ".",
    "paths": {
      "@/*": ["src/*"],
      "@types/*": ["src/types/*"],
      "@utils/*": ["src/utils/*"]
    }
  }
}
```

```typescript
// Instead of: import { User } from "../../types/user"
import type { User } from "@types/user";

// Instead of: import { capitalize } from "../../../utils/string"
import { capitalize } from "@utils/string";
```

**Important**: `paths` in tsconfig is only for TypeScript. For runtime, you also need `tsconfig-paths` or a bundler config (webpack alias, Vite's `resolve.alias`).

---

## 10.6 Declaration Files (.d.ts)

Declaration files describe the types of JavaScript code — they contain only type information, no runtime code.

### What They Look Like

```typescript
// express/index.d.ts (simplified)
declare module "express" {
  interface Request {
    body: unknown;
    params: Record<string, string>;
    query: Record<string, string | string[]>;
  }

  interface Response {
    json(body: unknown): this;
    send(body: string): this;
    status(code: number): this;
  }

  type RequestHandler = (req: Request, res: Response, next: () => void) => void;

  interface Application {
    get(path: string, ...handlers: RequestHandler[]): this;
    post(path: string, ...handlers: RequestHandler[]): this;
    listen(port: number, callback?: () => void): void;
  }

  function express(): Application;
  export = express;
}
```

### Writing Your Own .d.ts

```typescript
// legacy-lib.d.ts — typing a plain JavaScript file

// Declare a module's types
declare module "legacy-lib" {
  export function doSomething(value: string): number;
  export const VERSION: string;

  export interface Config {
    timeout: number;
    retries: number;
  }

  export class LegacyClient {
    constructor(config: Config);
    connect(): void;
    send(data: unknown): Promise<unknown>;
  }
}

// Declare a global variable (set by a script tag, for example)
declare const __APP_VERSION__: string;
declare const __DEV__: boolean;

// Augment an existing module's types (declaration merging)
declare module "express" {
  interface Request {
    user?: { id: string; role: string };  // added by our auth middleware
  }
}
```

### @types Packages

The DefinitelyTyped repository provides `.d.ts` files for thousands of JavaScript packages:

```bash
npm install --save-dev @types/node
npm install --save-dev @types/lodash
npm install --save-dev @types/express
```

---

## 10.7 Ambient Declarations

`declare` tells TypeScript about things that exist at runtime but aren't in your TypeScript files.

```typescript
// globals.d.ts

// Declare global variables
declare const process: {
  env: Record<string, string | undefined>;
  argv: string[];
  exit(code?: number): never;
};

// Declare a global function
declare function fetch(url: string, init?: RequestInit): Promise<Response>;

// Declare a global class
declare class EventEmitter {
  on(event: string, listener: (...args: unknown[]) => void): this;
  emit(event: string, ...args: unknown[]): boolean;
  off(event: string, listener: (...args: unknown[]) => void): this;
}

// Declare a namespace
declare namespace NodeJS {
  interface ProcessEnv {
    NODE_ENV: "development" | "production" | "test";
    PORT?: string;
    DATABASE_URL?: string;
  }
}

// After this declaration:
const port = process.env.PORT;  // string | undefined — typed!
const env = process.env.NODE_ENV;  // "development" | "production" | "test"
```

---

## 10.8 Namespaces (Legacy)

Namespaces are TypeScript's older module system, predating ES modules. Still used in declaration files but rarely in new application code.

```typescript
// validation.ts — namespace organization
namespace Validation {
  export interface StringValidator {
    isAcceptable(s: string): boolean;
  }

  export class LettersOnlyValidator implements StringValidator {
    isAcceptable(s: string): boolean {
      return /^[A-Za-z]+$/.test(s);
    }
  }

  export class ZipCodeValidator implements StringValidator {
    isAcceptable(s: string): boolean {
      return s.length === 5 && /\d{5}/.test(s);
    }
  }
}

// Usage
const validators: { [key: string]: Validation.StringValidator } = {
  letters: new Validation.LettersOnlyValidator(),
  zip: new Validation.ZipCodeValidator(),
};

// Nested namespaces
namespace App {
  export namespace Models {
    export interface User { id: number; name: string; }
  }
  export namespace Services {
    export class UserService {
      getUser(id: number): Models.User {
        return { id, name: "Alice" };
      }
    }
  }
}

const service = new App.Services.UserService();
const user: App.Models.User = service.getUser(1);
```

**When to use namespaces**: Almost never in new application code. Use ES modules instead. Namespaces are appropriate in `.d.ts` files to organize type declarations for large global APIs.

---

## 10.9 Structuring a TypeScript Project

### Barrel Files (index.ts)

```
src/
├── index.ts              ← main entry point
├── types/
│   ├── index.ts          ← re-exports all types
│   ├── user.ts
│   ├── post.ts
│   └── api.ts
├── services/
│   ├── index.ts          ← re-exports all services
│   ├── user.service.ts
│   └── post.service.ts
├── utils/
│   ├── index.ts          ← re-exports all utils
│   ├── string.ts
│   ├── date.ts
│   └── validation.ts
└── config.ts
```

```typescript
// src/types/index.ts
export type { User, CreateUserInput } from "./user";
export type { Post, CreatePostInput } from "./post";
export type { ApiResponse, PaginatedResponse } from "./api";

// src/services/index.ts
export { UserService } from "./user.service";
export { PostService } from "./post.service";

// src/index.ts — public API of the entire app/library
export * from "./types";
export * from "./services";
export { default as config } from "./config";
```

### Circular Imports — How to Avoid

```typescript
// BAD: user.ts imports from post.ts, post.ts imports from user.ts
// user.ts
import { Post } from "./post";  // Post imports User → circular!
export interface User {
  posts: Post[];
}

// post.ts
import { User } from "./user";  // User imports Post → circular!
export interface Post {
  author: User;
}

// GOOD: Extract shared types to a third file
// types.ts — no imports from the other modules
export interface User {
  posts: PostRef[];  // Use IDs or refs instead of full objects
}
export interface Post {
  authorId: string;  // reference, not full User object
}

// Or: extract to a shared types file
// shared-types.ts
export interface UserRef { id: string; name: string; }
export interface PostRef { id: string; title: string; }
```

---

## Complete Example: Organizing a REST API Client

```typescript
// src/api/types.ts
export interface ApiConfig {
  baseUrl: string;
  token: string;
  timeout?: number;
}

export type HttpMethod = "GET" | "POST" | "PUT" | "PATCH" | "DELETE";

export interface RequestOptions {
  method?: HttpMethod;
  body?: unknown;
  headers?: Record<string, string>;
}

export interface ApiError {
  status: number;
  message: string;
  code: string;
}

// src/api/client.ts
import type { ApiConfig, RequestOptions, ApiError } from "./types";

export class ApiClient {
  constructor(private readonly config: ApiConfig) {}

  async request<T>(endpoint: string, options: RequestOptions = {}): Promise<T> {
    const { method = "GET", body, headers = {} } = options;
    const url = `${this.config.baseUrl}${endpoint}`;
    const timeout = this.config.timeout ?? 10000;

    const controller = new AbortController();
    const timeoutId = setTimeout(() => controller.abort(), timeout);

    try {
      const response = await fetch(url, {
        method,
        headers: {
          "Content-Type": "application/json",
          Authorization: `Bearer ${this.config.token}`,
          ...headers,
        },
        body: body !== undefined ? JSON.stringify(body) : undefined,
        signal: controller.signal,
      });

      if (!response.ok) {
        const error: ApiError = await response.json();
        throw new Error(`${error.code}: ${error.message}`);
      }

      return response.json() as Promise<T>;
    } finally {
      clearTimeout(timeoutId);
    }
  }

  get<T>(endpoint: string, headers?: Record<string, string>): Promise<T> {
    return this.request<T>(endpoint, { method: "GET", headers });
  }

  post<T>(endpoint: string, body: unknown): Promise<T> {
    return this.request<T>(endpoint, { method: "POST", body });
  }

  put<T>(endpoint: string, body: unknown): Promise<T> {
    return this.request<T>(endpoint, { method: "PUT", body });
  }

  delete<T>(endpoint: string): Promise<T> {
    return this.request<T>(endpoint, { method: "DELETE" });
  }
}

// src/api/resources/users.ts
import type { ApiClient } from "../client";

export interface User {
  id: string;
  name: string;
  email: string;
}

export interface CreateUserInput {
  name: string;
  email: string;
  password: string;
}

export class UsersResource {
  constructor(private client: ApiClient) {}

  getAll(): Promise<User[]> {
    return this.client.get<User[]>("/users");
  }

  getById(id: string): Promise<User> {
    return this.client.get<User>(`/users/${id}`);
  }

  create(input: CreateUserInput): Promise<User> {
    return this.client.post<User>("/users", input);
  }

  update(id: string, data: Partial<CreateUserInput>): Promise<User> {
    return this.client.put<User>(`/users/${id}`, data);
  }

  delete(id: string): Promise<void> {
    return this.client.delete<void>(`/users/${id}`);
  }
}

// src/api/index.ts
import { ApiClient } from "./client";
import { UsersResource } from "./resources/users";
import type { ApiConfig } from "./types";

export type { ApiConfig, ApiError } from "./types";
export type { User, CreateUserInput } from "./resources/users";

export class ApiSdk {
  public readonly users: UsersResource;

  constructor(config: ApiConfig) {
    const client = new ApiClient(config);
    this.users = new UsersResource(client);
  }
}

// src/main.ts
import { ApiSdk, type User } from "./api";

const api = new ApiSdk({
  baseUrl: "https://api.example.com",
  token: "my-token",
  timeout: 5000,
});

async function main() {
  const users: User[] = await api.users.getAll();
  console.log(users);

  const alice: User = await api.users.create({
    name: "Alice",
    email: "alice@example.com",
    password: "secret",
  });
  console.log(alice.id);
}

main();
```

---

## Summary

TypeScript modules are ES modules extended with type-only imports, declaration files, and ambient declarations. Named exports provide explicit, refactoring-friendly interfaces; default exports are best for single-responsibility modules. Barrel files (`index.ts`) create clean public APIs for folders. Module resolution controls how TypeScript finds imports — `node16`/`nodenext` is the modern choice. Declaration files (`.d.ts`) and `@types` packages add types to JavaScript libraries. Namespaces are a legacy feature — prefer ES modules for new code.

---

## Key Takeaways

- A file is a module if it contains `import` or `export`; otherwise it's a script (global scope)
- **Named exports** are preferred over defaults for better tooling support and refactoring
- **`import type`** ensures type-only imports are erased cleanly at compile time
- **Barrel files** (`index.ts`) aggregate re-exports to create clean public APIs
- **Module resolution** (`node16`/`nodenext`) must match your runtime expectations
- **`.d.ts` files** add types to JavaScript without changing the JavaScript
- **`@types/*` packages** are the standard way to add types to third-party JavaScript

---

## Practice Questions

1. What makes a TypeScript file a "module" vs. a "script"?
2. What is the difference between `import type` and `import`?
3. When would you use a default export vs. named exports?
4. What is a barrel file, and what problem does it solve?
5. What is a `.d.ts` file, and when would you write one?
6. What is declaration merging, and how is it useful with `@types` packages?

---

## Exercises

**Exercise 1**: Organize a hypothetical e-commerce codebase into modules. Define the folder structure and create barrel `index.ts` files for: `types`, `services`, `utils`, and `api`.

**Exercise 2**: Write a `.d.ts` declaration file for a hypothetical legacy JavaScript module `string-utils` that exports: `capitalize(s)`, `slugify(s)`, `truncate(s, maxLength)`, and a constant `VERSION`.

**Exercise 3**: Create a path alias setup in `tsconfig.json` for `@services`, `@types`, and `@utils`. Then rewrite some imports to use these aliases.

**Exercise 4**: Implement a circular import situation (two files importing each other), diagnose why it's a problem, and refactor it by extracting shared types to a third file.

---

*Next: [Chapter 11 — Advanced Types](11-advanced-types.md)*
