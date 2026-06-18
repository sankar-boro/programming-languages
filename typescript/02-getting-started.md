# Chapter 2 — Getting Started

---

## 2.1 Installing Node.js

TypeScript runs on Node.js for development, and your compiled JavaScript also runs on Node.js for server-side code. Even for browser-targeting projects, Node.js is required for the TypeScript compiler.

### Installing Node.js

**Recommended: Use a version manager**

```bash
# Linux/macOS: install nvm (Node Version Manager)
curl -o- https://raw.githubusercontent.com/nvm-sh/nvm/v0.39.7/install.sh | bash

# Restart your terminal, then:
nvm install --lts        # install the Long-Term Support version
nvm use --lts            # switch to it
node --version           # verify: v20.x.x or later
npm --version            # verify: 10.x.x or later
```

**Windows: Use fnm**

```powershell
# Using winget
winget install Schniz.fnm

# Then in a new terminal:
fnm install --lts
fnm use lts-latest
node --version
```

**Direct installer**: Download from nodejs.org if you prefer.

### Verify Installation

```bash
node --version   # v20.11.0 or similar
npm --version    # 10.2.4 or similar
```

---

## 2.2 Installing TypeScript

TypeScript is distributed as an npm package. The compiler is `tsc` (TypeScript Compiler).

### Global Installation

```bash
npm install -g typescript
tsc --version   # Version 5.x.x
```

### Project-Local Installation (Recommended)

For consistent versions across teams:

```bash
mkdir my-ts-project
cd my-ts-project
npm init -y                          # create package.json
npm install --save-dev typescript    # install TypeScript locally
npx tsc --version                    # use local version
```

**Why local over global?** Different projects may require different TypeScript versions. Local installation ensures everyone on the team uses the same version. Always prefer `npx tsc` over `tsc` in scripts.

### ts-node — Running TypeScript Directly

`ts-node` compiles and runs TypeScript on-the-fly without generating files:

```bash
npm install --save-dev ts-node @types/node

# Run a TypeScript file directly
npx ts-node src/main.ts
```

For development, `ts-node` is convenient. For production, always compile with `tsc` first.

---

## 2.3 Your First TypeScript File

```bash
# Create the project structure
mkdir my-ts-project && cd my-ts-project
npm init -y
npm install --save-dev typescript ts-node @types/node

# Create the source file
mkdir src
```

```typescript
// src/main.ts

// TypeScript adds types to JavaScript
const message: string = "Hello, TypeScript!";
const year: number = 2026;
const isAwesome: boolean = true;

function greet(name: string, times: number = 1): string {
  return `${"Hello, ".repeat(times)}${name}!`;
}

// TypeScript catches errors at compile time
// greet(42);  // ERROR: Argument of type 'number' is not assignable to parameter of type 'string'

console.log(greet("World"));           // Hello, World!
console.log(greet("TypeScript", 3));   // Hello, Hello, Hello, TypeScript!
console.log(message, year, isAwesome);
```

### Compiling and Running

```bash
# Compile TypeScript to JavaScript
npx tsc src/main.ts

# This creates src/main.js
node src/main.js

# Or run directly with ts-node
npx ts-node src/main.ts
```

### Examining the Compiled Output

The compiled `src/main.js` will look like:

```javascript
"use strict";
// Types are completely gone
const message = "Hello, TypeScript!";
const year = 2026;
const isAwesome = true;

function greet(name, times = 1) {
  return `${"Hello, ".repeat(times)}${name}!`;
}

console.log(greet("World"));
console.log(greet("TypeScript", 3));
console.log(message, year, isAwesome);
```

Notice: `: string`, `: number`, `: boolean` — all type annotations — have been stripped. The output is clean JavaScript.

---

## 2.4 tsconfig.json — The TypeScript Configuration File

Instead of passing flags to `tsc` every time, you define your project configuration in `tsconfig.json`. This file controls everything about how TypeScript compiles your code.

### Generating a Default tsconfig.json

```bash
npx tsc --init
```

This creates a `tsconfig.json` with many options (most commented out). Let's build one from scratch:

### A Production-Ready tsconfig.json

```json
{
  "compilerOptions": {
    /* ── Target Output ──────────────────────────────── */
    "target": "ES2020",          // What JS version to emit
    "module": "commonjs",        // Module system for output
    "lib": ["ES2020"],           // Type definitions included
    
    /* ── Output ─────────────────────────────────────── */
    "outDir": "./dist",          // Where to put compiled .js files
    "rootDir": "./src",          // Where your .ts source files are
    "declaration": true,         // Generate .d.ts files alongside .js
    "declarationMap": true,      // Source maps for .d.ts files
    "sourceMap": true,           // Source maps for debugging
    
    /* ── Strict Type Checking ───────────────────────── */
    "strict": true,              // Enable ALL strict checks (highly recommended)
    // strict is shorthand for all of these:
    // "strictNullChecks": true,        // null and undefined are not every type
    // "strictFunctionTypes": true,     // stricter function type checking
    // "strictBindCallApply": true,     // correct types for bind, call, apply
    // "strictPropertyInitialization": true, // class properties must be initialized
    // "noImplicitAny": true,           // error when type is implicitly any
    // "noImplicitThis": true,          // error when this is implicitly any
    // "alwaysStrict": true,            // emit "use strict" in every file
    
    /* ── Additional Checks ──────────────────────────── */
    "noUnusedLocals": true,       // Error on unused local variables
    "noUnusedParameters": true,   // Error on unused function parameters
    "noImplicitReturns": true,    // Error if function doesn't always return
    "noFallthroughCasesInSwitch": true, // Error on switch case fallthrough
    
    /* ── Module Resolution ──────────────────────────── */
    "moduleResolution": "node",   // How to resolve imports
    "esModuleInterop": true,      // Better CommonJS/ESM interop
    "allowSyntheticDefaultImports": true, // Allow default imports from CJS modules
    "resolveJsonModule": true,    // Allow importing .json files
    
    /* ── Paths ──────────────────────────────────────── */
    "baseUrl": ".",              // Base for path aliases
    "paths": {
      "@/*": ["src/*"]           // @ alias for src/
    },
    
    /* ── Misc ──────────────────────────────────────── */
    "forceConsistentCasingInFileNames": true, // Case-sensitive imports
    "skipLibCheck": true          // Skip type-checking declaration files
  },
  "include": ["src/**/*"],      // Files to compile
  "exclude": [
    "node_modules",
    "dist",
    "**/*.test.ts"              // Don't include test files in main build
  ]
}
```

### Essential Options Explained

#### `target`

Controls what JavaScript version the output uses:

```typescript
// TypeScript source
const nums = [1, 2, 3];
const doubled = nums.map(n => n * 2);
```

With `"target": "ES5"` (for old browsers):
```javascript
var nums = [1, 2, 3];
var doubled = nums.map(function(n) { return n * 2; });
```

With `"target": "ES2020"` (modern environments):
```javascript
const nums = [1, 2, 3];
const doubled = nums.map(n => n * 2);  // arrow functions preserved
```

Common targets:
- `ES5`: Maximum compatibility (IE11)
- `ES2015`/`ES6`: Modern browsers (no IE)
- `ES2020`: Node.js 14+, modern browsers
- `ESNext`: Latest features (risky for compatibility)

#### `strict`

The most important option. Always enable it. `strict: true` enables a set of checks that catch real bugs:

```typescript
// With strict: false (BAD)
function greet(name) {           // name is implicitly `any` — no error
  console.log(name.toUpperCase());  // might crash at runtime if name isn't a string
}

// With strict: true (GOOD)
// ERROR: Parameter 'name' implicitly has an 'any' type.
function greet(name) { ... }

// Must be explicit:
function greet(name: string) { ... }  // OK
```

#### `strictNullChecks`

Part of `strict`. Without it, `null` and `undefined` can be assigned to any type:

```typescript
// strictNullChecks: false (dangerous)
let name: string = null;  // allowed — but will crash when you use it

// strictNullChecks: true (safe)
let name: string = null;  // ERROR: Type 'null' is not assignable to type 'string'
let safeName: string | null = null;  // OK — explicitly nullable
```

#### `outDir` and `rootDir`

```
Project structure:
├── tsconfig.json
├── src/
│   ├── main.ts
│   └── utils/
│       └── helper.ts
└── dist/         ← compiled output mirrors src/ structure
    ├── main.js
    └── utils/
        └── helper.js
```

---

## 2.5 Running TypeScript

### Option 1: Compile Once

```bash
# Compile all files described by tsconfig.json
npx tsc

# Output appears in dist/ (or whatever outDir specifies)
node dist/main.js
```

### Option 2: Watch Mode

```bash
# Recompile automatically whenever files change
npx tsc --watch
# or
npx tsc -w
```

### Option 3: ts-node (Development)

```bash
# Run TypeScript directly — no compilation step
npx ts-node src/main.ts

# With ts-node-dev (auto-restart on file changes, like nodemon)
npm install --save-dev ts-node-dev
npx ts-node-dev --respawn src/main.ts
```

### Option 4: Faster Alternatives

For large projects, `tsc` can be slow (it does full type-checking). These tools use `tsc` for type-checking but faster transpilers for running:

```bash
# esbuild — extremely fast bundler/transpiler
npm install --save-dev esbuild

# tsx — modern ts-node alternative using esbuild
npm install --save-dev tsx
npx tsx src/main.ts
```

**Important**: `tsx` and `esbuild` do NOT type-check — they just strip types. Always run `tsc --noEmit` separately for type checking.

### Adding npm Scripts

```json
// package.json
{
  "scripts": {
    "build": "tsc",
    "build:watch": "tsc --watch",
    "start": "node dist/main.js",
    "dev": "ts-node-dev --respawn src/main.ts",
    "typecheck": "tsc --noEmit",
    "clean": "rm -rf dist"
  }
}
```

```bash
npm run build      # compile TypeScript
npm run dev        # run in development mode with auto-restart
npm run typecheck  # check types without emitting files
```

---

## 2.6 Project Structure Best Practices

### Small Project

```
my-project/
├── package.json
├── tsconfig.json
├── .gitignore
├── src/
│   ├── index.ts          ← entry point
│   ├── types.ts          ← shared type definitions
│   └── utils.ts          ← utility functions
└── dist/                 ← git-ignored, compiled output
```

### Medium Project

```
my-project/
├── package.json
├── tsconfig.json
├── tsconfig.build.json   ← production build config (excludes tests)
├── .gitignore
├── src/
│   ├── index.ts
│   ├── types/
│   │   ├── index.ts      ← re-exports all types
│   │   ├── user.ts
│   │   └── api.ts
│   ├── services/
│   │   ├── user.service.ts
│   │   └── auth.service.ts
│   ├── utils/
│   │   ├── validation.ts
│   │   └── formatting.ts
│   └── config.ts
├── tests/
│   ├── user.service.test.ts
│   └── validation.test.ts
└── dist/
```

### .gitignore

```gitignore
node_modules/
dist/
*.js.map
.env
.env.local
```

### tsconfig for Tests

```json
// tsconfig.build.json — excludes test files
{
  "extends": "./tsconfig.json",
  "exclude": [
    "node_modules",
    "dist",
    "tests/**/*",
    "**/*.test.ts",
    "**/*.spec.ts"
  ]
}
```

---

## 2.7 IDE Setup — VS Code

VS Code has built-in TypeScript support (it's written in TypeScript). No extensions are strictly required, but these enhance the experience:

### Essential Settings (`.vscode/settings.json`)

```json
{
  "typescript.preferences.importModuleSpecifier": "relative",
  "typescript.updateImportsOnFileMove.enabled": "always",
  "editor.formatOnSave": true,
  "editor.defaultFormatter": "esbenp.prettier-vscode",
  "[typescript]": {
    "editor.defaultFormatter": "esbenp.prettier-vscode"
  },
  "typescript.tsdk": "node_modules/typescript/lib"
}
```

The last setting (`typescript.tsdk`) tells VS Code to use the project's local TypeScript version, not the one bundled with VS Code. This is important for consistency.

### Useful VS Code Commands

- **Go to Definition**: `F12` — jump to where a type is defined
- **Peek Definition**: `Alt+F12` — view definition inline
- **Find All References**: `Shift+F12` — every usage of a symbol
- **Rename Symbol**: `F2` — rename across all files
- **Quick Fix**: `Ctrl+.` — TypeScript's suggested fixes
- **Hover**: hover over any identifier to see its type

---

## Complete Getting Started Example

Let's build a complete mini-project to verify everything works:

```bash
mkdir hello-typescript && cd hello-typescript
npm init -y
npm install --save-dev typescript ts-node @types/node
npx tsc --init
```

Edit `tsconfig.json`:

```json
{
  "compilerOptions": {
    "target": "ES2020",
    "module": "commonjs",
    "strict": true,
    "outDir": "./dist",
    "rootDir": "./src",
    "esModuleInterop": true,
    "skipLibCheck": true
  },
  "include": ["src/**/*"]
}
```

```typescript
// src/types.ts
export interface Todo {
  id: number;
  title: string;
  completed: boolean;
  createdAt: Date;
}
```

```typescript
// src/todo-service.ts
import { Todo } from "./types";

let nextId = 1;
const todos: Todo[] = [];

export function createTodo(title: string): Todo {
  const todo: Todo = {
    id: nextId++,
    title,
    completed: false,
    createdAt: new Date(),
  };
  todos.push(todo);
  return todo;
}

export function completeTodo(id: number): Todo | undefined {
  const todo = todos.find((t) => t.id === id);
  if (todo) {
    todo.completed = true;
  }
  return todo;
}

export function getTodos(completedOnly?: boolean): Todo[] {
  if (completedOnly === undefined) return todos;
  return todos.filter((t) => t.completed === completedOnly);
}
```

```typescript
// src/main.ts
import { createTodo, completeTodo, getTodos } from "./todo-service";

const t1 = createTodo("Learn TypeScript");
const t2 = createTodo("Build a project");
const t3 = createTodo("Write tests");

completeTodo(t1.id);

const allTodos = getTodos();
const completedTodos = getTodos(true);

console.log("All todos:");
allTodos.forEach((t) => {
  console.log(`  [${t.completed ? "x" : " "}] ${t.title}`);
});

console.log(`\nCompleted: ${completedTodos.length}/${allTodos.length}`);
```

```bash
# Run it
npx ts-node src/main.ts

# Or compile and run
npx tsc
node dist/main.js
```

Output:
```
All todos:
  [x] Learn TypeScript
  [ ] Build a project
  [ ] Write tests

Completed: 1/3
```

---

## Summary

TypeScript requires Node.js and the `tsc` compiler, both easily installed via npm. Configuration lives in `tsconfig.json` — the most important setting is `strict: true`. You can run TypeScript directly with `ts-node` during development, or compile to JavaScript with `tsc` for production. VS Code provides excellent TypeScript support out of the box. A well-structured TypeScript project separates source files (`src/`) from compiled output (`dist/`), with types defined in dedicated files.

---

## Key Takeaways

- Install TypeScript locally per-project for consistent versions: `npm install --save-dev typescript`
- Always use `npx tsc` or `npm run` scripts — not global `tsc`
- `strict: true` in `tsconfig.json` enables all important checks — always use it
- `ts-node` for development, compiled `tsc` output for production
- VS Code uses `node_modules/typescript/lib` for TypeScript — set `typescript.tsdk` accordingly
- Types are erased at compile time — the output is plain JavaScript

---

## Practice Questions

1. What is the difference between `tsc` and `ts-node`?
2. What does the `strict` option in tsconfig.json enable?
3. Why should TypeScript be installed locally per project rather than globally?
4. What does `outDir` in tsconfig.json control?
5. How do you run the TypeScript compiler in watch mode?

---

## Exercises

**Exercise 1**: Set up a new TypeScript project from scratch with `strict: true`, `outDir: "./dist"`, and `rootDir: "./src"`. Verify it compiles successfully.

**Exercise 2**: Add a `build` and `dev` script to `package.json`. Run both to verify they work.

**Exercise 3**: In your tsconfig.json, set `noUnusedLocals: true`. Create a TypeScript file with an unused variable and verify the compiler reports an error.

**Exercise 4**: Compare the TypeScript source of the todo example above with the compiled JavaScript in `dist/`. List every TypeScript-specific syntax that was removed.

---

*Next: [Chapter 3 — Basic Types and Variables](03-basics.md)*
