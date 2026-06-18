# Chapter 12 — Working with JavaScript

> *"TypeScript's relationship with JavaScript is its greatest strength. You can add types gradually, use any JavaScript library, and interop with the entire npm ecosystem — all with full type safety."*

---

## 12.1 TypeScript and JavaScript Together

TypeScript is designed to work alongside JavaScript. You can:
1. Use any JavaScript library (with `@types` packages or your own declarations)
2. Gradually migrate JavaScript files to TypeScript
3. Include `.js` files in a TypeScript project

### Project Configuration for Mixed Files

```json
// tsconfig.json — configuring JavaScript/TypeScript interop
{
  "compilerOptions": {
    "allowJs": true,          // Allow .js files in the project
    "checkJs": true,          // Type-check .js files (optional but recommended)
    "maxNodeModuleJsDepth": 1, // How deep to check JS in node_modules
    "strict": true,
    "outDir": "./dist",
    "rootDir": "./src",
    "esModuleInterop": true
  },
  "include": ["src/**/*"]  // Picks up both .ts and .js files
}
```

```
src/
├── main.ts          ← TypeScript file
├── utils.js         ← JavaScript file — allowJs lets TypeScript include it
├── legacy.js        ← JavaScript file — checkJs adds type checking
└── types.d.ts       ← Declaration file — pure types
```

---

## 12.2 Using JavaScript Libraries

When you `import` from a JavaScript library, TypeScript needs type information. There are three situations:

### Situation 1: Library Ships Its Own Types

Modern libraries (axios, zod, prisma, etc.) include `.d.ts` files in their npm package:

```typescript
import axios from "axios";
// axios ships its own types — full autocomplete and type safety immediately

const response = await axios.get<{ id: number; name: string }[]>("/api/users");
const users = response.data;  // TypeScript knows: { id: number; name: string }[]
```

### Situation 2: @types Package Available

Older JavaScript libraries often have community-maintained types on DefinitelyTyped:

```bash
npm install express
npm install --save-dev @types/express
```

```typescript
import express, { Request, Response } from "express";
// @types/express provides Request, Response, etc.

const app = express();

app.get("/", (req: Request, res: Response) => {
  res.json({ status: "ok" });
});
```

```bash
# Common @types packages
npm install --save-dev @types/node    # Node.js built-ins (fs, path, http, etc.)
npm install --save-dev @types/express # Express framework
npm install --save-dev @types/lodash  # Lodash utility library
npm install --save-dev @types/jest    # Jest test framework
```

### Situation 3: No Types Available — Write Your Own

```typescript
// You: import "old-library" has no type information

// Option A: quick and dirty — cast to any (avoid this)
const lib = require("old-library") as any;
lib.doSomething("value");  // no type safety

// Option B: write a declaration file
// old-library.d.ts
declare module "old-library" {
  export interface Config {
    timeout: number;
    retries: number;
  }

  export function doSomething(value: string): Promise<{ result: string }>;
  export function initialize(config: Config): void;

  export class Client {
    constructor(config: Config);
    send(data: unknown): Promise<unknown>;
    close(): void;
  }
}
```

---

## 12.3 Declaration Files — Writing .d.ts

A `.d.ts` file is a TypeScript type declaration file — it contains only type information, no runtime code. The JavaScript code exists separately.

### Anatomy of a Declaration File

```typescript
// src/types/api.d.ts

// Declare exported values from a module
declare module "my-api-sdk" {

  // Interface declarations
  export interface User {
    id: string;
    name: string;
    email: string;
    createdAt: string;  // ISO date string from API
  }

  export interface ApiOptions {
    baseUrl: string;
    apiKey: string;
    timeout?: number;
  }

  // Function declarations
  export function createClient(options: ApiOptions): ApiClient;
  export function parseUser(raw: unknown): User;

  // Class declaration
  export class ApiClient {
    constructor(options: ApiOptions);
    getUser(id: string): Promise<User>;
    listUsers(page?: number): Promise<{ users: User[]; total: number }>;
    createUser(input: { name: string; email: string }): Promise<User>;
    deleteUser(id: string): Promise<void>;
  }

  // Namespace (for grouping related types)
  export namespace Errors {
    export class ApiError extends Error {
      constructor(message: string, public statusCode: number);
    }
    export class NetworkError extends Error {}
    export class AuthError extends Error {}
  }
}
```

### Global Declaration Files

```typescript
// src/types/globals.d.ts — augments global types

// Extend the global Window interface
interface Window {
  analytics: {
    track(event: string, properties?: Record<string, unknown>): void;
    identify(userId: string, traits?: Record<string, unknown>): void;
  };
  __DEV__: boolean;
}

// Declare global variables (e.g., from webpack DefinePlugin)
declare const __APP_VERSION__: string;
declare const __BUILD_DATE__: string;
declare const __COMMIT_HASH__: string;

// Augment process.env types (Node.js)
declare namespace NodeJS {
  interface ProcessEnv {
    readonly NODE_ENV: "development" | "production" | "test";
    readonly PORT?: string;
    readonly DATABASE_URL: string;
    readonly JWT_SECRET: string;
    readonly REDIS_URL?: string;
  }
}

// Declare a CSS module type
declare module "*.css" {
  const styles: Record<string, string>;
  export default styles;
}

// Declare image imports
declare module "*.png" {
  const url: string;
  export default url;
}
```

### Module Augmentation — Adding to Existing Types

```typescript
// src/types/express-extension.d.ts

// Augment Express Request to add our custom properties
// (added by our authentication middleware)
declare module "express" {
  interface Request {
    user?: {
      id: string;
      role: "admin" | "user" | "guest";
      permissions: string[];
    };
    requestId: string;
    startTime: number;
  }
}

// Now in our route handlers:
import { Request, Response } from "express";

function getProfile(req: Request, res: Response): void {
  // TypeScript knows req.user exists (it was added above)
  if (!req.user) {
    res.status(401).json({ error: "Unauthorized" });
    return;
  }
  // req.user: { id: string; role: ...; permissions: string[] }
  res.json({ userId: req.user.id });
}
```

---

## 12.4 checkJs — Type-Checking JavaScript Files

With `checkJs: true`, TypeScript checks your JavaScript files using JSDoc annotations:

```javascript
// utils.js — plain JavaScript, but type-checked via JSDoc

/**
 * @param {string} name
 * @param {number} [age]
 * @returns {string}
 */
function greet(name, age) {
  if (age !== undefined) {
    return `Hello, ${name}! You are ${age} years old.`;
  }
  return `Hello, ${name}!`;
}

/**
 * @typedef {Object} Config
 * @property {string} host
 * @property {number} port
 * @property {boolean} [debug]
 */

/**
 * @param {Config} config
 * @returns {void}
 */
function startServer(config) {
  console.log(`Starting server at ${config.host}:${config.port}`);
}

// TypeScript (via checkJs) will type-check calls to these functions
greet("Alice", 30);  // OK
greet(42);           // ERROR: 42 is not assignable to string
```

---

## 12.5 Gradual Migration: JavaScript to TypeScript

The practical approach to migrating a JavaScript project:

### Phase 1: Add TypeScript Infrastructure

```bash
npm install --save-dev typescript @types/node
npx tsc --init
```

```json
// tsconfig.json — permissive to start
{
  "compilerOptions": {
    "target": "ES2020",
    "module": "commonjs",
    "allowJs": true,          // keep .js files working
    "checkJs": false,         // don't check JS files yet
    "outDir": "./dist",
    "rootDir": "./src",
    "strict": false,          // start without strict mode
    "skipLibCheck": true
  },
  "include": ["src/**/*"]
}
```

### Phase 2: Convert Files One at a Time

Rename `.js` to `.ts` and fix errors:

```javascript
// Before: utils.js
function formatDate(date) {
  return date.toISOString().split("T")[0];
}
module.exports = { formatDate };
```

```typescript
// After: utils.ts
export function formatDate(date: Date): string {
  return date.toISOString().split("T")[0];
}
```

### Phase 3: Enable Strict Mode Gradually

```json
// tsconfig.json — enable strict checks one by one
{
  "compilerOptions": {
    "strictNullChecks": true,    // Phase 3a
    "noImplicitAny": true,       // Phase 3b
    "strictFunctionTypes": true, // Phase 3c
    // ...eventually: "strict": true
  }
}
```

### Phase 4: Clean Up any Types

```typescript
// Search for 'any' in your codebase and replace with real types
// Before:
function process(data: any): any {
  return data.value;
}

// After:
interface DataWithValue {
  value: string;
}
function process(data: DataWithValue): string {
  return data.value;
}
```

---

## 12.6 Type Assertions vs Type Guards for Unknown Data

When receiving data from external sources (APIs, JSON, localStorage), use type guards rather than assertions:

```typescript
// BAD: type assertion — TypeScript trusts you, no validation
const user = JSON.parse(localStorage.getItem("user")!) as User;
// If the stored JSON doesn't match User, runtime errors occur

// GOOD: type guard — validates before using
function isUser(value: unknown): value is User {
  return (
    typeof value === "object" &&
    value !== null &&
    typeof (value as Record<string, unknown>).id === "number" &&
    typeof (value as Record<string, unknown>).name === "string"
  );
}

const raw = JSON.parse(localStorage.getItem("user") ?? "null");
if (isUser(raw)) {
  // raw: User — validated
  console.log(raw.name);
} else {
  console.error("Invalid user data in localStorage");
}

// A generic validator approach
type Schema<T> = {
  [K in keyof T]: (value: unknown) => value is T[K];
};

function createValidator<T>(schema: Schema<T>) {
  return (value: unknown): value is T => {
    if (typeof value !== "object" || value === null) return false;
    const obj = value as Record<string, unknown>;
    return (Object.keys(schema) as Array<keyof T>).every((key) =>
      schema[key](obj[key as string])
    );
  };
}

const isUserValid = createValidator<User>({
  id: (v): v is number => typeof v === "number",
  name: (v): v is string => typeof v === "string",
  email: (v): v is string => typeof v === "string" && v.includes("@"),
});
```

---

## 12.7 Working with JSON

```typescript
// JSON.parse returns 'any' — annotate carefully
const raw: unknown = JSON.parse(jsonString);  // cast to unknown for safety

// Safe JSON parsing with error handling
function parseJSON<T>(json: string, validator: (v: unknown) => v is T): T | null {
  try {
    const parsed: unknown = JSON.parse(json);
    return validator(parsed) ? parsed : null;
  } catch {
    return null;
  }
}

// JSON serialization — JSON.stringify handles most types
const user: User = { id: 1, name: "Alice", email: "alice@ex.com" };
const json = JSON.stringify(user);  // string

// Classes with toJSON
class Money {
  constructor(
    public readonly amount: number,
    public readonly currency: string
  ) {}

  toJSON(): { amount: number; currency: string } {
    return { amount: this.amount, currency: this.currency };
  }
}

const price = new Money(9.99, "USD");
console.log(JSON.stringify(price));  // {"amount":9.99,"currency":"USD"}
```

---

## 12.8 require() and CommonJS

When working in Node.js with CommonJS modules, TypeScript handles the `require` function:

```typescript
// CommonJS-style import (TypeScript compiles these to require() with "module": "commonjs")
import fs from "fs";
import path from "path";

// The esModuleInterop option enables this default import syntax
// Without it, you'd need:
// import * as fs from "fs";
// import * as path from "path";

// Reading a file
const content: string = fs.readFileSync(
  path.join(__dirname, "config.json"),
  "utf-8"
);

// Dynamic require — loses type safety
const moduleName = "fs";
const dynamicModule = require(moduleName);  // any

// Safe dynamic import (ES modules)
const dynamicImport = await import("./config");  // fully typed!
```

---

## Complete Example: Typing a Third-Party JavaScript SDK

Suppose you're working with a payment SDK that has no TypeScript support:

```typescript
// src/types/payment-sdk.d.ts

declare module "payment-sdk" {

  type Currency = "USD" | "EUR" | "GBP" | "JPY";
  type PaymentStatus = "pending" | "completed" | "failed" | "refunded";
  type CardBrand = "visa" | "mastercard" | "amex" | "discover";

  export interface PaymentIntent {
    id: string;
    amount: number;  // in cents
    currency: Currency;
    status: PaymentStatus;
    createdAt: number;  // unix timestamp
    metadata: Record<string, string>;
  }

  export interface Card {
    id: string;
    last4: string;
    brand: CardBrand;
    expMonth: number;
    expYear: number;
  }

  export interface CreatePaymentIntentOptions {
    amount: number;
    currency: Currency;
    description?: string;
    metadata?: Record<string, string>;
    customerId?: string;
  }

  export interface PaymentSDKConfig {
    apiKey: string;
    environment: "sandbox" | "production";
  }

  export class PaymentClient {
    constructor(config: PaymentSDKConfig);

    createPaymentIntent(options: CreatePaymentIntentOptions): Promise<PaymentIntent>;
    getPaymentIntent(id: string): Promise<PaymentIntent>;
    confirmPayment(intentId: string, cardId: string): Promise<PaymentIntent>;
    cancelPayment(intentId: string): Promise<PaymentIntent>;
    refundPayment(intentId: string, amount?: number): Promise<PaymentIntent>;

    listCards(customerId: string): Promise<Card[]>;
    addCard(customerId: string, cardToken: string): Promise<Card>;
    removeCard(cardId: string): Promise<void>;
  }

  export function formatAmount(amount: number, currency: Currency): string;
  export function isPaymentSuccessful(intent: PaymentIntent): boolean;
}

// src/payment.service.ts — using the typed SDK
import { PaymentClient, type PaymentIntent, type Currency } from "payment-sdk";

export class PaymentService {
  private client: PaymentClient;

  constructor(apiKey: string) {
    this.client = new PaymentClient({
      apiKey,
      environment: process.env.NODE_ENV === "production" ? "production" : "sandbox",
    });
  }

  async charge(
    amountInCents: number,
    currency: Currency,
    customerId: string,
    cardId: string,
    description?: string
  ): Promise<PaymentIntent> {
    const intent = await this.client.createPaymentIntent({
      amount: amountInCents,
      currency,
      description,
      customerId,
    });

    return this.client.confirmPayment(intent.id, cardId);
  }

  async refund(intentId: string, amountInCents?: number): Promise<PaymentIntent> {
    return this.client.refundPayment(intentId, amountInCents);
  }
}
```

---

## Summary

TypeScript's relationship with JavaScript is its superpower. The `allowJs` and `checkJs` options let you include and check JavaScript files. `@types` packages provide types for thousands of JavaScript libraries. When types don't exist, you write `.d.ts` declaration files. Module augmentation extends existing type definitions. Gradual migration allows large JavaScript codebases to adopt TypeScript incrementally. Type assertions bypass the type system — prefer type guards for data from external sources.

---

## Key Takeaways

- **`allowJs: true`** includes `.js` files; **`checkJs: true`** type-checks them with JSDoc annotations
- **`@types/*` packages** are the primary way to type JavaScript libraries
- **`.d.ts` files** contain only type declarations — no runtime code
- **Module augmentation** (`declare module "..."`) extends existing type definitions, including third-party types
- **Gradual migration**: rename `.js` to `.ts` one file at a time; start with `strict: false`, tighten over time
- **`unknown` instead of `any`** for external data — forces validation before use
- **`import type`** is erased completely at compile time — use for type-only imports

---

## Practice Questions

1. What is the difference between `allowJs` and `checkJs` in tsconfig?
2. How do you add types to a JavaScript library that has no TypeScript support?
3. What is module augmentation, and when would you use it?
4. Why is using `unknown` safer than `any` for parsed JSON?
5. How do you type `process.env` variables to avoid runtime surprises?

---

## Exercises

**Exercise 1**: Write a complete `.d.ts` declaration file for a hypothetical `logger.js` module that exports: `log(level, message, meta?)`, `createLogger(name, options?)`, and a `Logger` class.

**Exercise 2**: Augment the `Express.Request` interface to add `session: { userId?: string; cart: string[] }` and write a middleware that populates it.

**Exercise 3**: Write a `safeParseJSON<T>(json: string, validator: (v: unknown) => v is T): Result<T, string>` function where `Result` is a discriminated union.

**Exercise 4**: Take an existing small JavaScript module (e.g., a utility file you've written), add JSDoc type annotations, enable `checkJs`, and fix any type errors TypeScript finds.

---

*Next: [Chapter 13 — Async Programming](13-async.md)*
