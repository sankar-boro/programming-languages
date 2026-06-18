# Chapter 13 — Async Programming

> *"TypeScript makes async programming type-safe. Promises have types, async functions have typed return values, and the compiler ensures you handle errors and nullability at every await point."*

---

## 13.1 The Promise Type

`Promise<T>` is the fundamental type for async operations. `T` is the type of the resolved value.

```typescript
// Promise<T> — represents an eventual value of type T
const immediate: Promise<number> = Promise.resolve(42);
const delayed: Promise<string> = new Promise((resolve) => {
  setTimeout(() => resolve("hello"), 1000);
});
const failing: Promise<never> = Promise.reject(new Error("failed"));

// TypeScript tracks the resolved type
immediate.then((value) => {
  // value: number — TypeScript knows
  console.log(value.toFixed(2));
});

delayed.then((s) => {
  // s: string — TypeScript knows
  console.log(s.toUpperCase());
});

// Creating typed Promises
function delay(ms: number): Promise<void> {
  return new Promise((resolve) => setTimeout(resolve, ms));
}

function fetchNumber(url: string): Promise<number> {
  return fetch(url).then((r) => r.json() as Promise<number>);
}

// Chaining Promises with preserved types
function processUser(id: number): Promise<string> {
  return fetch(`/api/users/${id}`)
    .then((r): Promise<{ name: string; email: string }> => r.json())
    .then((user) => user.name);  // user is { name: string; email: string }
}
```

### Promise Constructor — Typing resolve and reject

```typescript
function readFile(path: string): Promise<string> {
  return new Promise((resolve, reject) => {
    // resolve: (value: string) => void
    // reject: (reason?: unknown) => void

    try {
      const content = require("fs").readFileSync(path, "utf-8") as string;
      resolve(content);  // must pass string — TypeScript enforces it
    } catch (err) {
      reject(err);
    }
  });
}
```

---

## 13.2 async/await — Typed Async Functions

`async` functions always return a `Promise`. TypeScript infers the promise type from the return value.

### Basic async/await

```typescript
// Explicit return type: Promise<string>
async function fetchUserName(id: number): Promise<string> {
  const response = await fetch(`/api/users/${id}`);

  if (!response.ok) {
    throw new Error(`HTTP error: ${response.status}`);
  }

  const user: { id: number; name: string } = await response.json();
  return user.name;  // TypeScript knows this becomes the resolved value
}

// TypeScript infers Promise<string> from the return value
async function greet() {
  return "Hello, World!";  // returns Promise<string>
}

// Await unwraps the Promise type
async function main() {
  const name = await fetchUserName(1);
  // name: string — await removes the Promise wrapper

  const greeting = await greet();
  // greeting: string

  console.log(`${greeting} ${name}`);
}
```

### async/await vs .then()

```typescript
interface User {
  id: number;
  name: string;
  posts: Post[];
}

interface Post {
  id: number;
  title: string;
}

// Promise chain style
function loadUserWithPosts_chain(userId: number): Promise<User> {
  return fetch(`/api/users/${userId}`)
    .then((r) => r.json() as Promise<User>)
    .then((user) =>
      fetch(`/api/users/${userId}/posts`)
        .then((r) => r.json() as Promise<Post[]>)
        .then((posts) => ({ ...user, posts }))
    );
}

// async/await style — same logic, much more readable
async function loadUserWithPosts(userId: number): Promise<User> {
  const userRes = await fetch(`/api/users/${userId}`);
  const user: User = await userRes.json();

  const postsRes = await fetch(`/api/users/${userId}/posts`);
  const posts: Post[] = await postsRes.json();

  return { ...user, posts };
}
```

---

## 13.3 Error Handling in Async Functions

TypeScript's `catch` blocks receive `unknown` by default (with `useUnknownInCatchVariables: true`, part of `strict`).

```typescript
// Basic try/catch
async function fetchData(url: string): Promise<string> {
  try {
    const response = await fetch(url);

    if (!response.ok) {
      throw new Error(`HTTP ${response.status}: ${response.statusText}`);
    }

    return await response.text();
  } catch (error: unknown) {
    // error is unknown — must narrow before using
    if (error instanceof Error) {
      console.error("Network error:", error.message);
      throw error;  // rethrow for caller to handle
    }
    throw new Error("Unknown error occurred");
  }
}

// Typed error hierarchy
class ApiError extends Error {
  constructor(
    message: string,
    public readonly statusCode: number,
    public readonly endpoint: string
  ) {
    super(message);
    this.name = "ApiError";
  }
}

class NetworkError extends Error {
  constructor(message: string, public readonly cause: unknown) {
    super(message);
    this.name = "NetworkError";
  }
}

async function apiRequest<T>(endpoint: string): Promise<T> {
  let response: Response;

  try {
    response = await fetch(endpoint);
  } catch (err) {
    throw new NetworkError("Failed to connect", err);
  }

  if (!response.ok) {
    throw new ApiError(
      `API request failed: ${response.statusText}`,
      response.status,
      endpoint
    );
  }

  return response.json() as Promise<T>;
}

// Handling specific error types
async function safeApiRequest<T>(endpoint: string): Promise<T | null> {
  try {
    return await apiRequest<T>(endpoint);
  } catch (error) {
    if (error instanceof ApiError) {
      if (error.statusCode === 404) return null;
      if (error.statusCode === 401) {
        // redirect to login
        return null;
      }
    }
    if (error instanceof NetworkError) {
      console.error("Network issue:", error.message);
      return null;
    }
    throw error;  // unexpected error — rethrow
  }
}

// Result type for explicit error handling (no exceptions)
type AsyncResult<T, E = Error> =
  | { success: true; data: T }
  | { success: false; error: E };

async function tryCatch<T>(
  fn: () => Promise<T>
): Promise<AsyncResult<T>> {
  try {
    const data = await fn();
    return { success: true, data };
  } catch (error) {
    return {
      success: false,
      error: error instanceof Error ? error : new Error(String(error)),
    };
  }
}

// Usage — no try/catch in caller
const result = await tryCatch(() => apiRequest<User[]>("/api/users"));
if (result.success) {
  console.log(result.data.length);  // data: User[]
} else {
  console.error(result.error.message);  // error: Error
}
```

---

## 13.4 Promise.all, race, allSettled, any

TypeScript types all Promise combinators precisely.

### Promise.all — All Must Succeed

```typescript
// TypeScript infers the tuple type of Promise.all results
async function loadDashboard(userId: string) {
  const [user, posts, notifications] = await Promise.all([
    fetch(`/api/users/${userId}`).then((r) => r.json() as Promise<User>),
    fetch(`/api/posts?userId=${userId}`).then((r) => r.json() as Promise<Post[]>),
    fetch(`/api/notifications?userId=${userId}`).then((r) => r.json() as Promise<Notification[]>),
  ]);

  // user: User
  // posts: Post[]
  // notifications: Notification[]
  // TypeScript knows each element's type from the tuple

  return { user, posts, notifications };
}

// Sequential vs parallel
async function sequential(ids: number[]): Promise<User[]> {
  const users: User[] = [];
  for (const id of ids) {
    users.push(await fetchUser(id));  // one at a time
  }
  return users;
}

async function parallel(ids: number[]): Promise<User[]> {
  return Promise.all(ids.map((id) => fetchUser(id)));  // all at once
}
```

### Promise.allSettled — Collect All Results

```typescript
// allSettled never rejects — each result is either fulfilled or rejected
async function loadAllUsers(ids: number[]): Promise<User[]> {
  const results = await Promise.allSettled(ids.map((id) => fetchUser(id)));

  // results: PromiseSettledResult<User>[]
  // Each element is either:
  //   { status: "fulfilled"; value: User }
  //   { status: "rejected"; reason: unknown }

  const users: User[] = [];
  const errors: unknown[] = [];

  for (const result of results) {
    if (result.status === "fulfilled") {
      users.push(result.value);  // value: User
    } else {
      errors.push(result.reason);  // reason: unknown
    }
  }

  if (errors.length > 0) {
    console.warn(`${errors.length} users failed to load`);
  }

  return users;
}
```

### Promise.race — First One Wins

```typescript
// race returns Promise<User | Timeout>
async function fetchWithTimeout<T>(
  promise: Promise<T>,
  timeoutMs: number
): Promise<T> {
  const timeout = new Promise<never>((_, reject) =>
    setTimeout(() => reject(new Error(`Timeout after ${timeoutMs}ms`)), timeoutMs)
  );

  // race: the first to resolve or reject wins
  return Promise.race([promise, timeout]);
}

const user = await fetchWithTimeout(fetchUser(1), 5000);
```

### Promise.any — First Success Wins

```typescript
// any resolves with the first fulfilled promise
// only rejects if ALL promises reject (with AggregateError)
async function fetchFromNearestServer(userId: string): Promise<User> {
  const endpoints = [
    `https://server-us.api.com/users/${userId}`,
    `https://server-eu.api.com/users/${userId}`,
    `https://server-ap.api.com/users/${userId}`,
  ];

  try {
    return await Promise.any(
      endpoints.map((url) => fetch(url).then((r) => r.json() as Promise<User>))
    );
  } catch (error) {
    if (error instanceof AggregateError) {
      // All servers failed
      throw new Error(`All servers failed: ${error.errors.map(String).join(", ")}`);
    }
    throw error;
  }
}
```

---

## 13.5 Async Generators and Iterators

Async generators produce values over time using `yield`:

```typescript
// Async generator function — yields values asynchronously
async function* paginate<T>(
  fetchPage: (page: number) => Promise<{ data: T[]; hasNext: boolean }>,
  startPage = 1
): AsyncGenerator<T[], void, unknown> {
  let page = startPage;

  while (true) {
    const { data, hasNext } = await fetchPage(page);
    yield data;  // yields T[] each iteration

    if (!hasNext) break;
    page++;
  }
}

// Consuming an async generator
async function getAllUsers(): Promise<User[]> {
  const allUsers: User[] = [];

  const pages = paginate(
    (page) =>
      fetch(`/api/users?page=${page}&limit=100`)
        .then((r) => r.json() as Promise<{ data: User[]; hasNext: boolean }>),
    1
  );

  for await (const pageUsers of pages) {
    // pageUsers: User[] — each page
    allUsers.push(...pageUsers);
  }

  return allUsers;
}

// Async event stream
async function* streamEvents(url: string): AsyncGenerator<{ type: string; data: unknown }> {
  const response = await fetch(url);
  const reader = response.body!.getReader();
  const decoder = new TextDecoder();

  while (true) {
    const { done, value } = await reader.read();
    if (done) break;

    const text = decoder.decode(value);
    const lines = text.split("\n").filter((l) => l.startsWith("data:"));

    for (const line of lines) {
      const json = line.slice(5).trim();
      if (json) yield JSON.parse(json);
    }
  }
}
```

---

## 13.6 Async Patterns — Practical Examples

### Retry with Exponential Backoff

```typescript
interface RetryOptions {
  maxAttempts: number;
  initialDelayMs: number;
  maxDelayMs: number;
  backoffFactor: number;
}

async function withRetry<T>(
  fn: () => Promise<T>,
  options: RetryOptions = {
    maxAttempts: 3,
    initialDelayMs: 100,
    maxDelayMs: 5000,
    backoffFactor: 2,
  }
): Promise<T> {
  let delay = options.initialDelayMs;

  for (let attempt = 1; attempt <= options.maxAttempts; attempt++) {
    try {
      return await fn();
    } catch (error) {
      if (attempt === options.maxAttempts) {
        throw error;  // last attempt — rethrow
      }

      console.warn(`Attempt ${attempt} failed. Retrying in ${delay}ms...`);
      await new Promise((resolve) => setTimeout(resolve, delay));

      delay = Math.min(delay * options.backoffFactor, options.maxDelayMs);
    }
  }

  // TypeScript needs this even though it's unreachable
  throw new Error("Should not reach here");
}

// Usage
const user = await withRetry(() => fetchUser(1), {
  maxAttempts: 5,
  initialDelayMs: 200,
  maxDelayMs: 10000,
  backoffFactor: 2,
});
```

### Async Queue — Limiting Concurrency

```typescript
class AsyncQueue {
  private queue: Array<() => Promise<unknown>> = [];
  private running = 0;

  constructor(private readonly concurrency: number = 3) {}

  async add<T>(fn: () => Promise<T>): Promise<T> {
    return new Promise((resolve, reject) => {
      this.queue.push(async () => {
        try {
          resolve(await fn());
        } catch (e) {
          reject(e);
        }
      });
      this.process();
    });
  }

  private async process(): Promise<void> {
    if (this.running >= this.concurrency || this.queue.length === 0) return;

    this.running++;
    const task = this.queue.shift()!;

    try {
      await task();
    } finally {
      this.running--;
      this.process();
    }
  }
}

// Fetch with max 3 concurrent requests
const queue = new AsyncQueue(3);

const results = await Promise.all(
  userIds.map((id) => queue.add(() => fetchUser(id)))
);
```

### Debounce — Typed

```typescript
function debounce<T extends (...args: unknown[]) => unknown>(
  fn: T,
  delay: number
): (...args: Parameters<T>) => void {
  let timer: ReturnType<typeof setTimeout> | undefined;

  return (...args: Parameters<T>) => {
    clearTimeout(timer);
    timer = setTimeout(() => fn(...args), delay);
  };
}

// Async debounce — returns the result
function debounceAsync<T extends (...args: unknown[]) => Promise<unknown>>(
  fn: T,
  delay: number
): (...args: Parameters<T>) => Promise<Awaited<ReturnType<T>>> {
  let timer: ReturnType<typeof setTimeout> | undefined;
  let resolve: ((value: Awaited<ReturnType<T>>) => void) | undefined;
  let reject: ((reason: unknown) => void) | undefined;

  return (...args: Parameters<T>): Promise<Awaited<ReturnType<T>>> => {
    clearTimeout(timer);

    return new Promise((res, rej) => {
      resolve = res;
      reject = rej;

      timer = setTimeout(async () => {
        try {
          const result = await fn(...args);
          resolve?.(result as Awaited<ReturnType<T>>);
        } catch (e) {
          reject?.(e);
        }
      }, delay);
    });
  };
}

// Use debounce for search
const searchUsers = debounce(async (query: string) => {
  const users = await fetch(`/api/users?q=${query}`).then((r) => r.json());
  displayUsers(users);
}, 300);

function displayUsers(users: unknown): void {}
```

---

## 13.7 Typing Promise-Based Callbacks

```typescript
// Promisify a callback-based function
function promisify<T>(
  fn: (callback: (err: Error | null, result?: T) => void) => void
): () => Promise<T>;

function promisify<A, T>(
  fn: (arg1: A, callback: (err: Error | null, result?: T) => void) => void
): (arg1: A) => Promise<T>;

function promisify(fn: Function): Function {
  return (...args: unknown[]) =>
    new Promise((resolve, reject) => {
      fn(...args, (err: Error | null, result?: unknown) => {
        if (err) reject(err);
        else resolve(result);
      });
    });
}

// Node.js built-in: util.promisify
import { promisify as nodePromisify } from "util";
import fs from "fs";

const readFile = nodePromisify(fs.readFile);
const content = await readFile("./config.json", "utf-8");
// content: string | Buffer — TypeScript correctly infers based on encoding
```

---

## Summary

TypeScript's async support is built on the `Promise<T>` type, which carries the resolved value's type through the chain. `async` functions return `Promise<T>` where `T` is the return type. `await` unwraps `Promise<T>` to `T`. Error handling uses `try/catch` where `catch` receives `unknown` — requiring explicit narrowing. `Promise.all`, `allSettled`, `race`, and `any` are all precisely typed. Async generators (`async function*`) yield typed values asynchronously. Practical patterns like retry, concurrency limiting, and debounce are straightforwardly typed using generics.

---

## Key Takeaways

- **`Promise<T>`** is the core async type — `T` is what you get when the promise resolves
- **`async` functions** always return `Promise<T>` — TypeScript infers `T` from the return value
- **`await`** unwraps `Promise<T>` to `T` — inside an async function
- **`catch` receives `unknown`** with strict mode — always narrow error types before using them
- **`Promise.all`** types the result as a tuple matching its input array
- **`Promise.allSettled`** never rejects — each result is `{ status: "fulfilled", value: T }` or `{ status: "rejected", reason: unknown }`
- **`Awaited<T>`** extracts the resolved type: `Awaited<Promise<string>>` = `string`
- **Async generators** (`async function*`) enable lazy async sequences

---

## Practice Questions

1. What does TypeScript infer as the return type of `async function foo() { return 42; }`?
2. What type does `catch (error)` receive with `strict: true`? How do you narrow it?
3. What is the difference between `Promise.all` and `Promise.allSettled`?
4. What is `Awaited<T>` and when do you need it?
5. How do you type a function that returns `Promise<never>`?

---

## Exercises

**Exercise 1**: Implement a typed `timeout<T>(promise: Promise<T>, ms: number): Promise<T>` that rejects with a `TimeoutError` if the promise doesn't resolve within `ms` milliseconds.

**Exercise 2**: Write a generic `cache<T>(fn: (...args: unknown[]) => Promise<T>, ttlMs: number)` that memoizes the result of an async function, expiring after `ttlMs` milliseconds.

**Exercise 3**: Implement `mapAsync<T, U>(items: T[], fn: (item: T) => Promise<U>, concurrency: number): Promise<U[]>` that maps an async function over an array with limited concurrency.

**Exercise 4**: Build a `createEventStream<T>(url: string, parser: (raw: string) => T): AsyncGenerator<T>` that reads Server-Sent Events from a URL and yields parsed events.

---

*Next: [Chapter 14 — The TypeScript Compiler](14-compiler.md)*
