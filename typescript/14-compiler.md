# Chapter 14 — The TypeScript Compiler

> *"Understanding tsconfig.json is understanding TypeScript. The compiler options shape what TypeScript accepts, what JavaScript it emits, and how strictly it checks your code."*

---

## 14.1 How the TypeScript Compiler Works

```
TypeScript Source (.ts)
        │
        ▼
  ┌───────────┐
  │  Lexer    │ — Tokenizes source into tokens
  └─────┬─────┘
        │
        ▼
  ┌───────────┐
  │  Parser   │ — Tokens → Abstract Syntax Tree (AST)
  └─────┬─────┘
        │
        ▼
  ┌────────────────┐
  │  Binder        │ — Creates symbols, links AST nodes
  └────────┬───────┘
           │
           ▼
  ┌────────────────┐
  │  Type Checker  │ — Checks types, reports errors
  └────────┬───────┘
           │
           ▼
  ┌────────────────┐
  │  Emitter       │ — Generates .js, .d.ts, source maps
  └────────────────┘
```

Key insight: **Types are erased**. The emitter generates JavaScript by stripping all type annotations. The type checker is a separate phase that runs before emission.

```bash
# The two main things tsc does:
tsc                    # type-check AND emit .js files
tsc --noEmit           # type-check ONLY — no file generation (used in CI)
tsc --emitDeclarationOnly  # emit .d.ts only (not .js)
```

---

## 14.2 tsconfig.json — Deep Dive

### Locating tsconfig

TypeScript looks for `tsconfig.json` by searching up from the current directory. The `--project` flag specifies an explicit path:

```bash
tsc                           # finds tsconfig.json upward from cwd
tsc --project tsconfig.prod.json  # explicit config file
tsc src/main.ts               # compile single file (ignores tsconfig.json)
```

### Extending tsconfigs

```json
// tsconfig.base.json — shared settings
{
  "compilerOptions": {
    "target": "ES2020",
    "strict": true,
    "esModuleInterop": true,
    "skipLibCheck": true,
    "forceConsistentCasingInFileNames": true
  }
}
```

```json
// tsconfig.json — extends the base
{
  "extends": "./tsconfig.base.json",
  "compilerOptions": {
    "module": "commonjs",
    "outDir": "./dist",
    "rootDir": "./src",
    "declaration": true
  },
  "include": ["src/**/*"],
  "exclude": ["node_modules", "dist", "**/*.test.ts"]
}
```

```json
// tsconfig.test.json — for test files
{
  "extends": "./tsconfig.json",
  "compilerOptions": {
    "strict": false,
    "types": ["jest"]
  },
  "include": ["src/**/*", "tests/**/*"]
}
```

---

## 14.3 Target and Lib — What JavaScript to Emit

### target

Controls the JavaScript version of the emitted code:

```typescript
// TypeScript source
class Animal {
  #name: string;  // JavaScript private field

  constructor(name: string) {
    this.#name = name;
  }

  get name(): string { return this.#name; }
}

const fn = async () => {
  const nums = [1, 2, 3];
  return nums.at(-1) ?? 0;
};
```

With `"target": "ES5"`:
```javascript
var Animal = /** @class */ (function () {
  function Animal(name) {
    this._name = name;  // no # support in ES5
  }
  // ... complex getter polyfill
})();
```

With `"target": "ES2022"`:
```javascript
class Animal {
  #name;
  constructor(name) {
    this.#name = name;
  }
  get name() { return this.#name; }
}
const fn = async () => {
  const nums = [1, 2, 3];
  return nums.at(-1) ?? 0;
};
```

Common targets:
- `ES5` — maximum compat (IE11, very old Node.js)
- `ES2015` / `ES6` — modern browsers, Node.js 6+
- `ES2017` — async/await native
- `ES2020` — optional chaining, nullish coalescing, BigInt
- `ES2022` — class fields, top-level await
- `ESNext` — latest features (use with care)

### lib

`lib` specifies which built-in API type definitions to include. TypeScript uses this to know what globals are available:

```json
{
  "compilerOptions": {
    "target": "ES2020",
    "lib": ["ES2020", "DOM", "DOM.Iterable"]
  }
}
```

Common lib values:
- `ES5`, `ES2015`, ..., `ES2023`, `ESNext` — ECMAScript standard APIs
- `DOM` — browser globals (window, document, fetch, etc.)
- `DOM.Iterable` — iteration on DOM collections
- `WebWorker` — Web Worker globals
- `Node` — not a built-in; use `@types/node` package

```typescript
// If you don't include DOM in lib:
document.querySelector("div");  // ERROR: 'document' not found
fetch("/api");                  // ERROR: 'fetch' not found

// If you don't include ES2021 in lib:
[1, 2, 3].at(-1);              // ERROR: 'at' not found on Array
"hello".replaceAll("l", "L");  // ERROR: 'replaceAll' not found
```

---

## 14.4 Module and moduleResolution

### module

Controls the module format of the emitted JavaScript:

```typescript
// TypeScript source
import { readFileSync } from "fs";
export function readConfig(): string {
  return readFileSync("./config.json", "utf-8");
}
```

With `"module": "commonjs"`:
```javascript
"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.readConfig = void 0;
const fs_1 = require("fs");
function readConfig() {
    return (0, fs_1.readFileSync)("./config.json", "utf-8");
}
exports.readConfig = readConfig;
```

With `"module": "es2020"`:
```javascript
import { readFileSync } from "fs";
export function readConfig() {
    return readFileSync("./config.json", "utf-8");
}
```

Common module options:
- `commonjs` — Node.js default (require/module.exports)
- `es2015`/`es2020`/`esnext` — ES modules (import/export)
- `node16`/`nodenext` — Modern Node.js with dual CJS/ESM support
- `none` — No module system (globals only)

### moduleResolution

Controls how TypeScript finds imported modules:

```json
{
  "compilerOptions": {
    "module": "node16",
    "moduleResolution": "node16"  // must match module
  }
}
```

| Setting | When to use |
|---------|------------|
| `node` | Classic Node.js, most common |
| `node16` | Modern Node.js 16+ with ESM |
| `nodenext` | Latest Node.js resolution |
| `bundler` | Vite, webpack, esbuild — allows extensionless imports |

---

## 14.5 The strict Flag — What It Enables

`"strict": true` is a single flag that enables a set of checks. Understanding each sub-flag helps when migrating legacy code:

### strictNullChecks

```typescript
// strictNullChecks: false (dangerous!)
let name: string = null;  // allowed — but crashes when used
let element = document.getElementById("app");
element.textContent = "hello";  // crashes at runtime if element is null

// strictNullChecks: true (safe!)
let name: string = null;  // ERROR: null not assignable to string
let name2: string | null = null;  // OK — explicitly nullable

let element = document.getElementById("app");
// element: HTMLElement | null — TypeScript forces you to handle null
if (element) {
  element.textContent = "hello";  // OK
}
element!.textContent = "hello";  // non-null assertion — risky
```

### noImplicitAny

```typescript
// noImplicitAny: false
function process(data) {  // data is implicitly any — no error
  return data.value;  // might crash
}

// noImplicitAny: true
function process(data) {  // ERROR: 'data' implicitly has an 'any' type
  return data.value;
}

function process(data: { value: string }): string {  // OK — explicit
  return data.value;
}
```

### strictFunctionTypes

```typescript
// strictFunctionTypes: false — bivariant function parameters
type Logger = (message: string) => void;
const handler: (message: string | number) => void = (msg) => console.log(msg);
const logger: Logger = handler;  // allowed with bivariant (unsound!)
logger("hello");  // calls handler("hello") — ok
logger(42);       // TypeScript says this is ERROR... but handler handles it?

// strictFunctionTypes: true — contravariant function parameters
// Correctly rejects assignments that could cause runtime errors
```

### strictPropertyInitialization

```typescript
class User {
  name: string;     // ERROR with strict: class property 'name' has no initializer

  constructor() {
    // forgot to set this.name!
  }
}

// Fix: initialize in declaration or constructor
class User2 {
  name: string = "";  // OK — initialized in declaration

  // Or:
  name2: string;
  constructor() {
    this.name2 = "default";  // OK — set in constructor
  }

  // Or: use ! to opt out (be careful)
  name3!: string;  // tells TypeScript "I'll set this later"
}
```

### noImplicitThis

```typescript
// noImplicitThis: true
class Timer {
  seconds = 0;

  start() {
    // BAD: regular function — 'this' is implicitly any
    setInterval(function() {
      this.seconds++;  // ERROR: 'this' is typed 'any'
    }, 1000);

    // GOOD: arrow function captures 'this'
    setInterval(() => {
      this.seconds++;  // OK — 'this' is Timer
    }, 1000);
  }
}
```

### Full strict Flag Matrix

```json
// What "strict": true enables:
{
  "strictNullChecks": true,
  "noImplicitAny": true,
  "strictFunctionTypes": true,
  "strictBindCallApply": true,
  "strictPropertyInitialization": true,
  "noImplicitThis": true,
  "alwaysStrict": true,
  "useUnknownInCatchVariables": true,  // catch (e) where e: unknown
  "exactOptionalPropertyTypes": false  // not included in strict — opt-in separately
}
```

---

## 14.6 Additional Useful Flags

### noUnusedLocals and noUnusedParameters

```typescript
// noUnusedLocals: true
function process(): void {
  const unusedVar = 42;  // ERROR: 'unusedVar' is declared but never read
  console.log("done");
}

// noUnusedParameters: true
function greet(name: string, age: number): string {  // ERROR: 'age' declared but never read
  return `Hello, ${name}!`;
}

// Fix: prefix with _ to suppress the error
function greet(name: string, _age: number): string {  // OK
  return `Hello, ${name}!`;
}
```

### noImplicitReturns

```typescript
// noImplicitReturns: true — every code path must return
function getStatus(code: number): string {  // ERROR: function lacks ending return
  if (code === 200) return "OK";
  if (code === 404) return "Not Found";
  // missing return for other codes!
}

function getStatus2(code: number): string {
  if (code === 200) return "OK";
  if (code === 404) return "Not Found";
  return "Unknown";  // OK — all paths return
}
```

### exactOptionalPropertyTypes

```typescript
// exactOptionalPropertyTypes: true
interface Config {
  debug?: boolean;
}

const config: Config = {};

// With exactOptionalPropertyTypes: true
config.debug = undefined;  // ERROR! Optional means absent, not undefined

// Correctly absent:
const config2: Config = {};  // debug is absent

// If you need to explicitly store undefined:
interface Config2 {
  debug?: boolean | undefined;  // explicit undefined in the type
}
```

### declaration and declarationMap

```json
{
  "declaration": true,           // emit .d.ts files
  "declarationMap": true,        // emit .d.ts.map (source maps for declarations)
  "declarationDir": "./types"    // separate output directory for .d.ts files
}
```

### Strict Template Literal Type Checking

```typescript
// noUncheckedIndexedAccess: true (not in strict, but very useful)
const arr = [1, 2, 3];
const first = arr[0];  // number | undefined (not just number!)
// Forces you to check:
if (first !== undefined) {
  console.log(first.toFixed(2));
}
```

---

## 14.7 Project References — Multi-Package Builds

For monorepos or large projects with multiple packages, project references enable incremental builds:

```
packages/
├── core/
│   ├── tsconfig.json
│   └── src/
├── api/
│   ├── tsconfig.json  (references core)
│   └── src/
└── cli/
    ├── tsconfig.json  (references core and api)
    └── src/
```

```json
// packages/api/tsconfig.json
{
  "compilerOptions": {
    "composite": true,  // required for referenced projects
    "outDir": "./dist",
    "rootDir": "./src"
  },
  "references": [
    { "path": "../core" }  // depends on core
  ]
}

// packages/cli/tsconfig.json
{
  "compilerOptions": { "composite": true },
  "references": [
    { "path": "../core" },
    { "path": "../api" }
  ]
}
```

```bash
# Build all projects in dependency order
tsc --build          # or tsc -b

# Build a specific project
tsc --build packages/cli

# Clean build artifacts
tsc --build --clean

# Only rebuild what changed (incremental)
tsc --build --incremental
```

---

## 14.8 Performance Options

### skipLibCheck

```json
{
  "skipLibCheck": true  // Skip type-checking .d.ts files — massive speedup
}
```

`skipLibCheck` avoids re-checking declaration files for libraries. Safe because library authors test their types. Without it, conflicting `@types` versions can cause build failures.

### incremental

```json
{
  "incremental": true,
  "tsBuildInfoFile": "./dist/.tsbuildinfo"  // where to store build cache
}
```

Stores information about the last compilation. On subsequent runs, only recompiles changed files.

### isolatedModules

```json
{
  "isolatedModules": true
}
```

Ensures each file can be safely transpiled in isolation (required by esbuild, babel). Flags patterns that require cross-file type information:

```typescript
// With isolatedModules: true

// ERROR: const enums require cross-file knowledge
const enum Direction { Up, Down }  // ERROR

// ERROR: type-only re-exports need 'export type'
export { type SomeType } from "./types";  // must use 'export type'

// Use regular enums instead:
enum Direction { Up = "UP", Down = "DOWN" }  // OK
```

---

## 14.9 Practical tsconfig Recipes

### Node.js API Server

```json
{
  "compilerOptions": {
    "target": "ES2022",
    "module": "commonjs",
    "lib": ["ES2022"],
    "outDir": "./dist",
    "rootDir": "./src",
    "strict": true,
    "noUnusedLocals": true,
    "noUnusedParameters": true,
    "noImplicitReturns": true,
    "noFallthroughCasesInSwitch": true,
    "exactOptionalPropertyTypes": true,
    "declaration": true,
    "declarationMap": true,
    "sourceMap": true,
    "esModuleInterop": true,
    "skipLibCheck": true,
    "forceConsistentCasingInFileNames": true,
    "resolveJsonModule": true,
    "moduleResolution": "node"
  },
  "include": ["src/**/*"],
  "exclude": ["node_modules", "dist", "**/*.test.ts", "**/*.spec.ts"]
}
```

### Modern Node.js (ESM)

```json
{
  "compilerOptions": {
    "target": "ES2022",
    "module": "node16",
    "moduleResolution": "node16",
    "lib": ["ES2022"],
    "outDir": "./dist",
    "rootDir": "./src",
    "strict": true,
    "declaration": true,
    "sourceMap": true,
    "skipLibCheck": true
  },
  "include": ["src/**/*"]
}
```

### Library (published to npm)

```json
{
  "compilerOptions": {
    "target": "ES2020",
    "module": "esnext",
    "moduleResolution": "bundler",
    "lib": ["ES2020"],
    "outDir": "./dist",
    "rootDir": "./src",
    "strict": true,
    "declaration": true,
    "declarationMap": true,
    "sourceMap": true,
    "noUnusedLocals": true,
    "noUnusedParameters": true,
    "exactOptionalPropertyTypes": true,
    "skipLibCheck": true,
    "verbatimModuleSyntax": true  // ensure import type for type-only imports
  },
  "include": ["src/**/*"],
  "exclude": ["node_modules", "dist", "**/*.test.ts"]
}
```

---

## Summary

The TypeScript compiler does two separate jobs: type-checking and code emission. `tsconfig.json` controls both. `target` determines what JavaScript version is emitted; `lib` determines what built-in APIs are available to the type system. `module` and `moduleResolution` control import/export syntax and how TypeScript finds files. `strict: true` enables a suite of checks — `strictNullChecks` and `noImplicitAny` are the most impactful. Project references enable efficient builds in monorepos. `skipLibCheck` and `incremental` are key performance options.

---

## Key Takeaways

- **`tsc --noEmit`** — type check only, no output; use in CI
- **`strict: true`** enables 8 sub-flags — never ship production code without it
- **`target`** = what JavaScript is emitted; **`lib`** = what built-ins TypeScript knows about
- **`module: "node16"` + `moduleResolution: "node16"`** — modern Node.js ESM support
- **`skipLibCheck: true`** — safe performance win; skip checking third-party `.d.ts` files
- **`incremental: true`** — dramatically speeds up rebuild times in watch mode
- **Project references** enable parallelizable, incremental builds in monorepos

---

## Practice Questions

1. What are the two main things the TypeScript compiler does?
2. What sub-flags does `strict: true` enable?
3. What is the difference between `target` and `lib`?
4. What does `noEmit` do and when would you use it?
5. What is `isolatedModules` and why does esbuild/babel require it?
6. What is `skipLibCheck` and is it safe to use?

---

## Exercises

**Exercise 1**: Create three tsconfig files: `tsconfig.json` (base), `tsconfig.build.json` (production, extends base, excludes tests), and `tsconfig.test.json` (test build, extends base, includes tests). Verify they all work.

**Exercise 2**: Gradually enable strict flags one at a time on an existing TypeScript codebase. Start with `noImplicitAny`, then add `strictNullChecks`, tracking how many errors each adds.

**Exercise 3**: Set up a monorepo with two packages: `core` (a utility library) and `app` (imports from core). Configure project references and verify `tsc --build` works correctly.

**Exercise 4**: Compare the JavaScript output of a class with private fields (`#field`) between `target: "ES5"`, `target: "ES2015"`, and `target: "ES2022"`. Note the differences.

---

*Next: [Chapter 15 — TypeScript Internals](15-internals.md)*
