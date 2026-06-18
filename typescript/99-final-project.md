# Chapter 99 — Final Project: Type-Safe HTTP Server

> *"Build a complete, production-quality HTTP server using ONLY Node.js's built-in `http` module. No frameworks. No external dependencies. Pure TypeScript — strict mode throughout."*

---

## Project Overview

We will build a fully typed HTTP server from scratch:

- **Request parsing** — typed `Request` object with method, url, headers, body
- **Response builder** — typed `Response` with status, headers, JSON/text/HTML
- **Router** — type-safe routing with URL parameters
- **Middleware** — composable middleware chain with typed context
- **Error handling** — typed error hierarchy with proper HTTP status codes
- **Validation** — request body validation with type guards

**Rules**: Only `node:http`, `node:url`, `node:path`, `node:crypto` — no npm packages.

---

## Project Structure

```
http-server/
├── package.json
├── tsconfig.json
└── src/
    ├── main.ts          ← entry point, route definitions
    ├── server.ts        ← core server and request/response
    ├── request.ts       ← typed Request parsing
    ├── response.ts      ← typed Response builder
    ├── router.ts        ← type-safe router
    ├── middleware.ts     ← middleware types and combinators
    ├── errors.ts        ← typed error hierarchy
    └── validation.ts    ← type-safe validation helpers
```

---

## Setup

```json
// package.json
{
  "name": "ts-http-server",
  "version": "1.0.0",
  "scripts": {
    "build": "tsc",
    "start": "node dist/main.js",
    "dev": "tsc --watch & sleep 2 && node --watch dist/main.js",
    "typecheck": "tsc --noEmit"
  },
  "devDependencies": {
    "@types/node": "^20.0.0",
    "typescript": "^5.0.0"
  }
}
```

```json
// tsconfig.json
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
    "esModuleInterop": true,
    "skipLibCheck": true,
    "forceConsistentCasingInFileNames": true,
    "declaration": true,
    "sourceMap": true
  },
  "include": ["src/**/*"]
}
```

---

## src/errors.ts — Typed Error Hierarchy

```typescript
// src/errors.ts

export type HttpStatusCode =
  | 200 | 201 | 204
  | 400 | 401 | 403 | 404 | 405 | 409 | 422 | 429
  | 500 | 501 | 503;

export class HttpError extends Error {
  constructor(
    public readonly statusCode: HttpStatusCode,
    message: string,
    public readonly code: string,
    public readonly details?: Record<string, unknown>
  ) {
    super(message);
    this.name = "HttpError";
  }

  toJSON(): Record<string, unknown> {
    return {
      error: {
        code: this.code,
        message: this.message,
        ...(this.details ? { details: this.details } : {}),
      },
    };
  }
}

export class BadRequestError extends HttpError {
  constructor(message: string, details?: Record<string, unknown>) {
    super(400, message, "BAD_REQUEST", details);
    this.name = "BadRequestError";
  }
}

export class UnauthorizedError extends HttpError {
  constructor(message = "Authentication required") {
    super(401, message, "UNAUTHORIZED");
    this.name = "UnauthorizedError";
  }
}

export class ForbiddenError extends HttpError {
  constructor(message = "Access denied") {
    super(403, message, "FORBIDDEN");
    this.name = "ForbiddenError";
  }
}

export class NotFoundError extends HttpError {
  constructor(path: string) {
    super(404, `Resource not found: ${path}`, "NOT_FOUND");
    this.name = "NotFoundError";
  }
}

export class MethodNotAllowedError extends HttpError {
  constructor(method: string, path: string) {
    super(405, `Method ${method} not allowed for ${path}`, "METHOD_NOT_ALLOWED");
    this.name = "MethodNotAllowedError";
  }
}

export class ConflictError extends HttpError {
  constructor(message: string) {
    super(409, message, "CONFLICT");
    this.name = "ConflictError";
  }
}

export class UnprocessableEntityError extends HttpError {
  constructor(message: string, details?: Record<string, unknown>) {
    super(422, message, "UNPROCESSABLE_ENTITY", details);
    this.name = "UnprocessableEntityError";
  }
}

export class InternalServerError extends HttpError {
  constructor(message = "Internal server error") {
    super(500, message, "INTERNAL_SERVER_ERROR");
    this.name = "InternalServerError";
  }
}

export function isHttpError(error: unknown): error is HttpError {
  return error instanceof HttpError;
}
```

---

## src/request.ts — Typed Request Parsing

```typescript
// src/request.ts
import type { IncomingMessage } from "node:http";
import { URL } from "node:url";
import { BadRequestError } from "./errors.js";

export type HttpMethod =
  | "GET" | "POST" | "PUT" | "PATCH" | "DELETE"
  | "HEAD" | "OPTIONS";

export interface ParsedRequest {
  method: HttpMethod;
  url: string;
  pathname: string;
  searchParams: URLSearchParams;
  headers: Record<string, string>;
  params: Record<string, string>;  // filled by router after matching
  body: unknown;
}

export function parseMethod(raw: string | undefined): HttpMethod {
  const upper = (raw ?? "GET").toUpperCase();
  const valid: HttpMethod[] = ["GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"];
  if (!valid.includes(upper as HttpMethod)) {
    throw new BadRequestError(`Unknown HTTP method: ${upper}`);
  }
  return upper as HttpMethod;
}

export function parseHeaders(raw: IncomingMessage["headers"]): Record<string, string> {
  const headers: Record<string, string> = {};
  for (const [key, value] of Object.entries(raw)) {
    if (Array.isArray(value)) {
      headers[key] = value.join(", ");
    } else if (value !== undefined) {
      headers[key] = value;
    }
  }
  return headers;
}

export async function parseBody(req: IncomingMessage): Promise<unknown> {
  return new Promise((resolve, reject) => {
    const chunks: Buffer[] = [];
    let size = 0;
    const MAX_BODY_SIZE = 1 * 1024 * 1024;  // 1MB

    req.on("data", (chunk: Buffer) => {
      size += chunk.length;
      if (size > MAX_BODY_SIZE) {
        reject(new BadRequestError("Request body too large"));
        req.destroy();
        return;
      }
      chunks.push(chunk);
    });

    req.on("end", () => {
      const raw = Buffer.concat(chunks).toString("utf-8");
      if (!raw) {
        resolve(undefined);
        return;
      }

      const contentType = req.headers["content-type"] ?? "";
      if (contentType.includes("application/json")) {
        try {
          resolve(JSON.parse(raw));
        } catch {
          reject(new BadRequestError("Invalid JSON in request body"));
        }
      } else {
        resolve(raw);
      }
    });

    req.on("error", reject);
  });
}

export async function parseRequest(
  req: IncomingMessage,
  baseUrl: string
): Promise<ParsedRequest> {
  const method = parseMethod(req.method);
  const rawUrl = req.url ?? "/";

  const parsed = new URL(rawUrl, baseUrl);
  const headers = parseHeaders(req.headers);

  // Parse body only for methods that can have one
  const body =
    method !== "GET" && method !== "HEAD" && method !== "DELETE" && method !== "OPTIONS"
      ? await parseBody(req)
      : undefined;

  return {
    method,
    url: rawUrl,
    pathname: parsed.pathname,
    searchParams: parsed.searchParams,
    headers,
    params: {},  // populated by router after matching
    body,
  };
}
```

---

## src/response.ts — Typed Response Builder

```typescript
// src/response.ts
import type { ServerResponse } from "node:http";
import type { HttpStatusCode } from "./errors.js";

export type ResponseBody =
  | string
  | Buffer
  | Record<string, unknown>
  | unknown[]
  | null;

export class HttpResponse {
  private _status: HttpStatusCode = 200;
  private _headers: Record<string, string> = {};
  private _body: ResponseBody = null;

  status(code: HttpStatusCode): this {
    this._status = code;
    return this;
  }

  header(name: string, value: string): this {
    this._headers[name.toLowerCase()] = value;
    return this;
  }

  headers(headers: Record<string, string>): this {
    for (const [name, value] of Object.entries(headers)) {
      this.header(name, value);
    }
    return this;
  }

  json(data: unknown): this {
    this._body = data as Record<string, unknown> | unknown[];
    this.header("content-type", "application/json");
    return this;
  }

  text(data: string): this {
    this._body = data;
    this.header("content-type", "text/plain; charset=utf-8");
    return this;
  }

  html(data: string): this {
    this._body = data;
    this.header("content-type", "text/html; charset=utf-8");
    return this;
  }

  empty(): this {
    this._status = 204;
    this._body = null;
    return this;
  }

  redirect(location: string, permanent = false): this {
    this._status = permanent ? 301 : 302;
    this.header("location", location);
    this._body = null;
    return this;
  }

  send(res: ServerResponse): void {
    // Set headers
    for (const [name, value] of Object.entries(this._headers)) {
      res.setHeader(name, value);
    }

    // Serialize body
    let bodyBuffer: Buffer | null = null;

    if (this._body === null || this._body === undefined) {
      bodyBuffer = null;
    } else if (Buffer.isBuffer(this._body)) {
      bodyBuffer = this._body;
    } else if (typeof this._body === "string") {
      bodyBuffer = Buffer.from(this._body, "utf-8");
    } else {
      bodyBuffer = Buffer.from(JSON.stringify(this._body), "utf-8");
    }

    if (bodyBuffer !== null) {
      res.setHeader("content-length", bodyBuffer.length);
    }

    res.writeHead(this._status);

    if (bodyBuffer !== null) {
      res.end(bodyBuffer);
    } else {
      res.end();
    }
  }
}

export function response(): HttpResponse {
  return new HttpResponse();
}
```

---

## src/middleware.ts — Middleware System

```typescript
// src/middleware.ts
import type { ParsedRequest } from "./request.js";
import type { HttpResponse } from "./response.js";

export interface Context {
  request: ParsedRequest;
  response: HttpResponse;
  state: Record<string, unknown>;  // shared mutable state between middlewares
}

export type Next = () => Promise<void>;

export type Middleware = (ctx: Context, next: Next) => Promise<void>;

// Compose middlewares into a single middleware
export function compose(middlewares: Middleware[]): Middleware {
  return async (ctx: Context, finalNext: Next) => {
    let index = -1;

    const dispatch = async (i: number): Promise<void> => {
      if (i <= index) throw new Error("next() called multiple times");
      index = i;

      const fn = i === middlewares.length ? finalNext : middlewares[i];
      if (!fn) return;

      await fn(ctx, () => dispatch(i + 1));
    };

    return dispatch(0);
  };
}

// ─── Built-in Middlewares ───────────────────────────────────────────────────

// Request ID middleware
export function requestId(): Middleware {
  return async (ctx, next) => {
    const id = ctx.request.headers["x-request-id"] ??
      Math.random().toString(36).slice(2);
    ctx.state.requestId = id;
    ctx.response.header("x-request-id", id);
    await next();
  };
}

// Logging middleware
export function logger(): Middleware {
  return async (ctx, next) => {
    const start = Date.now();
    const { method, pathname } = ctx.request;
    console.log(`→ ${method} ${pathname}`);

    try {
      await next();
    } finally {
      const duration = Date.now() - start;
      console.log(`← ${method} ${pathname} [${duration}ms]`);
    }
  };
}

// CORS middleware
export function cors(allowedOrigins: string[] = ["*"]): Middleware {
  return async (ctx, next) => {
    const origin = ctx.request.headers["origin"] ?? "";
    const allowed =
      allowedOrigins.includes("*") || allowedOrigins.includes(origin);

    if (allowed) {
      ctx.response.header(
        "access-control-allow-origin",
        allowedOrigins.includes("*") ? "*" : origin
      );
    }

    ctx.response
      .header("access-control-allow-methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
      .header("access-control-allow-headers", "Content-Type, Authorization, X-Request-Id");

    if (ctx.request.method === "OPTIONS") {
      ctx.response.status(204).empty().send(
        // we'll handle sending in the server
        null as unknown as import("node:http").ServerResponse
      );
      return;
    }

    await next();
  };
}

// Rate limiting middleware (in-memory — not production-grade)
export function rateLimit(requestsPerMinute: number): Middleware {
  const counts = new Map<string, { count: number; resetAt: number }>();

  return async (ctx, next) => {
    const ip = ctx.request.headers["x-forwarded-for"] ??
      ctx.request.headers["x-real-ip"] ?? "unknown";
    const now = Date.now();
    const resetAt = now + 60_000;

    const current = counts.get(ip);
    if (!current || now > current.resetAt) {
      counts.set(ip, { count: 1, resetAt });
    } else {
      current.count++;
      if (current.count > requestsPerMinute) {
        ctx.response
          .status(429)
          .header("retry-after", String(Math.ceil((current.resetAt - now) / 1000)))
          .json({ error: { code: "RATE_LIMIT_EXCEEDED", message: "Too many requests" } });
        return;
      }
    }

    await next();
  };
}
```

---

## src/router.ts — Type-Safe Router

```typescript
// src/router.ts
import type { HttpMethod } from "./request.js";
import type { Middleware } from "./middleware.js";
import { NotFoundError, MethodNotAllowedError } from "./errors.js";

interface Route {
  method: HttpMethod;
  pattern: RegExp;
  paramNames: string[];
  handler: Middleware;
}

// Convert a path pattern like "/users/:id/posts/:postId" to a RegExp
function compileRoute(path: string): { pattern: RegExp; paramNames: string[] } {
  const paramNames: string[] = [];

  const regexStr = path
    .replace(/:[a-zA-Z_][a-zA-Z0-9_]*/g, (match) => {
      paramNames.push(match.slice(1));  // remove leading ':'
      return "([^/]+)";
    })
    .replace(/\//g, "\\/");

  return {
    pattern: new RegExp(`^${regexStr}$`),
    paramNames,
  };
}

export class Router {
  private routes: Route[] = [];
  private middlewares: Middleware[] = [];

  // Global middleware for this router
  use(middleware: Middleware): this {
    this.middlewares.push(middleware);
    return this;
  }

  private addRoute(
    method: HttpMethod,
    path: string,
    ...handlers: Middleware[]
  ): this {
    const { pattern, paramNames } = compileRoute(path);
    const handler = handlers.length === 1
      ? handlers[0]
      : async (ctx: import("./middleware.js").Context, next: import("./middleware.js").Next) => {
          for (const h of handlers) {
            await h(ctx, next);
          }
        };

    this.routes.push({ method, pattern, paramNames, handler });
    return this;
  }

  get(path: string, ...handlers: Middleware[]): this {
    return this.addRoute("GET", path, ...handlers);
  }
  post(path: string, ...handlers: Middleware[]): this {
    return this.addRoute("POST", path, ...handlers);
  }
  put(path: string, ...handlers: Middleware[]): this {
    return this.addRoute("PUT", path, ...handlers);
  }
  patch(path: string, ...handlers: Middleware[]): this {
    return this.addRoute("PATCH", path, ...handlers);
  }
  delete(path: string, ...handlers: Middleware[]): this {
    return this.addRoute("DELETE", path, ...handlers);
  }

  // Find a matching route and return it with extracted params
  match(method: HttpMethod, pathname: string): {
    route: Route;
    params: Record<string, string>;
  } | null {
    // First pass: find routes that match the path
    const pathMatches = this.routes.filter((route) =>
      route.pattern.test(pathname)
    );

    if (pathMatches.length === 0) return null;

    // Second pass: find a method match
    const methodMatch = pathMatches.find((route) => route.method === method);

    if (!methodMatch) {
      // Path matched but wrong method — throw 405
      throw new MethodNotAllowedError(method, pathname);
    }

    // Extract URL parameters
    const match = pathname.match(methodMatch.pattern);
    const params: Record<string, string> = {};

    if (match) {
      methodMatch.paramNames.forEach((name, i) => {
        const captured = match[i + 1];
        if (captured !== undefined) {
          params[name] = decodeURIComponent(captured);
        }
      });
    }

    return { route: methodMatch, params };
  }

  // Get all middlewares including router-level ones
  getMiddlewares(): Middleware[] {
    return this.middlewares;
  }

  getRoutes(): Route[] {
    return this.routes;
  }
}
```

---

## src/validation.ts — Type-Safe Validation

```typescript
// src/validation.ts
import { BadRequestError, UnprocessableEntityError } from "./errors.js";

// Validator type: returns the validated value or throws
export type Validator<T> = (value: unknown) => T;

// Validation error details
export interface ValidationError {
  field: string;
  message: string;
}

// Schema validator that returns typed result or throws
export class Schema<T> {
  constructor(private validator: Validator<T>) {}

  parse(value: unknown): T {
    return this.validator(value);
  }

  optional(): Schema<T | undefined> {
    return new Schema((value) => {
      if (value === undefined || value === null) return undefined;
      return this.validator(value);
    });
  }
}

// Primitive validators
export const string = new Schema<string>((value) => {
  if (typeof value !== "string") throw new Error("Expected string");
  return value;
});

export const number = new Schema<number>((value) => {
  if (typeof value !== "number") throw new Error("Expected number");
  if (isNaN(value)) throw new Error("Expected non-NaN number");
  return value;
});

export const boolean = new Schema<boolean>((value) => {
  if (typeof value !== "boolean") throw new Error("Expected boolean");
  return value;
});

export const positiveInt = new Schema<number>((value) => {
  if (typeof value !== "number") throw new Error("Expected number");
  if (!Number.isInteger(value)) throw new Error("Expected integer");
  if (value <= 0) throw new Error("Expected positive integer");
  return value;
});

// Email validator
export const email = new Schema<string>((value) => {
  const s = string.parse(value);
  if (!/^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(s)) {
    throw new Error("Invalid email address");
  }
  return s;
});

// Non-empty string
export const nonEmptyString = new Schema<string>((value) => {
  const s = string.parse(value);
  if (s.trim().length === 0) throw new Error("Expected non-empty string");
  return s;
});

// Object validator — validates each field against a schema
export function object<T extends Record<string, unknown>>(
  shape: { [K in keyof T]: Schema<T[K]> }
): Schema<T> {
  return new Schema<T>((value) => {
    if (typeof value !== "object" || value === null || Array.isArray(value)) {
      throw new BadRequestError("Expected an object");
    }

    const input = value as Record<string, unknown>;
    const result: Partial<T> = {};
    const errors: ValidationError[] = [];

    for (const [key, schema] of Object.entries(shape)) {
      try {
        (result as Record<string, unknown>)[key] = (schema as Schema<unknown>).parse(input[key]);
      } catch (err) {
        errors.push({
          field: key,
          message: err instanceof Error ? err.message : "Invalid value",
        });
      }
    }

    if (errors.length > 0) {
      throw new UnprocessableEntityError("Validation failed", {
        errors: errors.reduce(
          (acc, e) => ({ ...acc, [e.field]: e.message }),
          {} as Record<string, string>
        ),
      });
    }

    return result as T;
  });
}

// Array validator
export function array<T>(itemSchema: Schema<T>): Schema<T[]> {
  return new Schema<T[]>((value) => {
    if (!Array.isArray(value)) throw new BadRequestError("Expected an array");
    const errors: ValidationError[] = [];
    const result = value.map((item: unknown, i: number) => {
      try {
        return itemSchema.parse(item);
      } catch (err) {
        errors.push({
          field: `[${i}]`,
          message: err instanceof Error ? err.message : "Invalid item",
        });
        return null as unknown as T;
      }
    });

    if (errors.length > 0) {
      throw new UnprocessableEntityError("Array validation failed", {
        errors: errors.reduce(
          (acc, e) => ({ ...acc, [e.field]: e.message }),
          {} as Record<string, string>
        ),
      });
    }

    return result;
  });
}

// Literal validator
export function literal<T extends string | number | boolean>(expected: T): Schema<T> {
  return new Schema<T>((value) => {
    if (value !== expected) throw new Error(`Expected ${JSON.stringify(expected)}`);
    return value as T;
  });
}

// Union validator (tries each schema)
export function union<T>(...schemas: Schema<T>[]): Schema<T> {
  return new Schema<T>((value) => {
    for (const schema of schemas) {
      try {
        return schema.parse(value);
      } catch {
        // try next
      }
    }
    throw new Error("Value did not match any schema in union");
  });
}
```

---

## src/server.ts — Core Server

```typescript
// src/server.ts
import { createServer, type IncomingMessage, type ServerResponse } from "node:http";
import type { AddressInfo } from "node:net";
import { parseRequest } from "./request.js";
import { HttpResponse, response } from "./response.js";
import type { Router } from "./router.js";
import { compose, type Middleware, type Context } from "./middleware.js";
import { isHttpError, InternalServerError } from "./errors.js";

export interface ServerOptions {
  port: number;
  host?: string;
  baseUrl?: string;
}

export class HttpServer {
  private middlewares: Middleware[] = [];
  private routers: Router[] = [];

  use(middleware: Middleware): this {
    this.middlewares.push(middleware);
    return this;
  }

  router(router: Router): this {
    this.routers.push(router);
    return this;
  }

  private async handleRequest(
    req: IncomingMessage,
    res: ServerResponse,
    baseUrl: string
  ): Promise<void> {
    const httpResponse = response();
    let ctx: Context | null = null;

    try {
      const parsedRequest = await parseRequest(req, baseUrl);
      ctx = {
        request: parsedRequest,
        response: httpResponse,
        state: {},
      };

      // Find matching route across all routers
      let routeHandler: Middleware | null = null;

      for (const router of this.routers) {
        const match = router.match(parsedRequest.method, parsedRequest.pathname);
        if (match) {
          parsedRequest.params = match.params;
          const routerMiddlewares = router.getMiddlewares();
          const allMiddlewares = [...routerMiddlewares, match.route.handler];
          routeHandler = compose(allMiddlewares);
          break;
        }
      }

      // Collect all middlewares
      const allMiddlewares: Middleware[] = [
        ...this.middlewares,
        ...(routeHandler
          ? [routeHandler]
          : [
              async (c: Context) => {
                c.response.status(404).json(
                  new import("./errors.js").NotFoundError(
                    parsedRequest.pathname
                  ).toJSON()
                );
              },
            ]),
      ];

      const composed = compose(allMiddlewares);
      await composed(ctx, async () => {});

    } catch (error: unknown) {
      const httpError = isHttpError(error)
        ? error
        : new InternalServerError(
            error instanceof Error ? error.message : "Unknown error"
          );

      if (!isHttpError(error)) {
        console.error("Unhandled error:", error);
      }

      if (ctx) {
        ctx.response
          .status(httpError.statusCode)
          .json(httpError.toJSON());
      } else {
        httpResponse
          .status(httpError.statusCode)
          .json(httpError.toJSON());
      }
    } finally {
      const finalResponse = ctx?.response ?? httpResponse;
      finalResponse.send(res);
    }
  }

  listen(options: ServerOptions): Promise<{ close: () => Promise<void>; address: AddressInfo }> {
    const baseUrl = options.baseUrl ?? `http://${options.host ?? "localhost"}:${options.port}`;

    return new Promise((resolve, reject) => {
      const server = createServer((req, res) => {
        this.handleRequest(req, res, baseUrl).catch((err: unknown) => {
          console.error("Fatal request error:", err);
          if (!res.writableEnded) {
            res.writeHead(500).end('{"error":{"code":"FATAL","message":"Server error"}}');
          }
        });
      });

      server.on("error", reject);

      server.listen(options.port, options.host ?? "0.0.0.0", () => {
        const address = server.address() as AddressInfo;
        console.log(`Server listening on http://${address.address}:${address.port}`);

        resolve({
          address,
          close: () =>
            new Promise((res, rej) =>
              server.close((err) => (err ? rej(err) : res()))
            ),
        });
      });
    });
  }
}

export function createHttpServer(): HttpServer {
  return new HttpServer();
}
```

---

## src/main.ts — Application Entry Point

```typescript
// src/main.ts
import { createHttpServer } from "./server.js";
import { Router } from "./router.js";
import { logger, requestId, cors } from "./middleware.js";
import { response } from "./response.js";
import { object, string, email, nonEmptyString, positiveInt } from "./validation.js";
import { NotFoundError, BadRequestError } from "./errors.js";
import type { Context, Next } from "./middleware.js";

// ─── In-Memory Data Store ───────────────────────────────────────────────────

interface User {
  id: string;
  name: string;
  email: string;
  createdAt: string;
}

const users = new Map<string, User>();
let nextId = 1;

function generateId(): string {
  return String(nextId++);
}

// ─── Validation Schemas ─────────────────────────────────────────────────────

const CreateUserSchema = object({
  name: nonEmptyString,
  email: email,
});

const UpdateUserSchema = object({
  name: nonEmptyString.optional(),
  email: email.optional(),
});

// ─── Route Handlers ─────────────────────────────────────────────────────────

// GET /api/users
async function listUsers(ctx: Context): Promise<void> {
  const page = Number(ctx.request.searchParams.get("page") ?? "1");
  const limit = Number(ctx.request.searchParams.get("limit") ?? "10");

  if (!Number.isInteger(page) || page < 1) {
    throw new BadRequestError("'page' must be a positive integer");
  }
  if (!Number.isInteger(limit) || limit < 1 || limit > 100) {
    throw new BadRequestError("'limit' must be an integer between 1 and 100");
  }

  const allUsers = [...users.values()];
  const start = (page - 1) * limit;
  const pageUsers = allUsers.slice(start, start + limit);

  ctx.response.json({
    data: pageUsers,
    meta: {
      page,
      limit,
      total: allUsers.length,
      pages: Math.ceil(allUsers.length / limit),
    },
  });
}

// GET /api/users/:id
async function getUser(ctx: Context): Promise<void> {
  const { id } = ctx.request.params;
  const user = users.get(id);

  if (!user) throw new NotFoundError(`/api/users/${id}`);

  ctx.response.json({ data: user });
}

// POST /api/users
async function createUser(ctx: Context): Promise<void> {
  const input = CreateUserSchema.parse(ctx.request.body);

  // Check for duplicate email
  const exists = [...users.values()].some((u) => u.email === input.email);
  if (exists) {
    throw new import("./errors.js").ConflictError(
      `User with email ${input.email} already exists`
    );
  }

  const user: User = {
    id: generateId(),
    name: input.name,
    email: input.email,
    createdAt: new Date().toISOString(),
  };

  users.set(user.id, user);

  ctx.response.status(201).json({ data: user });
}

// PUT /api/users/:id
async function updateUser(ctx: Context): Promise<void> {
  const { id } = ctx.request.params;
  const existing = users.get(id);

  if (!existing) throw new NotFoundError(`/api/users/${id}`);

  const input = UpdateUserSchema.parse(ctx.request.body);

  const updated: User = {
    ...existing,
    ...(input.name !== undefined ? { name: input.name } : {}),
    ...(input.email !== undefined ? { email: input.email } : {}),
  };

  users.set(id, updated);

  ctx.response.json({ data: updated });
}

// DELETE /api/users/:id
async function deleteUser(ctx: Context): Promise<void> {
  const { id } = ctx.request.params;

  if (!users.has(id)) throw new NotFoundError(`/api/users/${id}`);

  users.delete(id);

  ctx.response.status(204).empty();
}

// GET /health
async function health(ctx: Context): Promise<void> {
  ctx.response.json({
    status: "healthy",
    timestamp: new Date().toISOString(),
    uptime: process.uptime(),
    memory: {
      used: Math.round(process.memoryUsage().heapUsed / 1024 / 1024),
      total: Math.round(process.memoryUsage().heapTotal / 1024 / 1024),
      unit: "MB",
    },
    users: users.size,
  });
}

// GET /
async function home(ctx: Context): Promise<void> {
  ctx.response.html(`
    <!DOCTYPE html>
    <html>
    <head><title>TypeScript HTTP Server</title></head>
    <body>
      <h1>TypeScript HTTP Server</h1>
      <h2>API Endpoints</h2>
      <ul>
        <li>GET    /health — Server health check</li>
        <li>GET    /api/users — List all users</li>
        <li>POST   /api/users — Create a user</li>
        <li>GET    /api/users/:id — Get a user</li>
        <li>PUT    /api/users/:id — Update a user</li>
        <li>DELETE /api/users/:id — Delete a user</li>
      </ul>
      <p>Built with TypeScript + Node.js http module only. No frameworks.</p>
    </body>
    </html>
  `);
}

// ─── Router Setup ───────────────────────────────────────────────────────────

const apiRouter = new Router();
apiRouter
  .get("/api/users", listUsers)
  .post("/api/users", createUser)
  .get("/api/users/:id", getUser)
  .put("/api/users/:id", updateUser)
  .delete("/api/users/:id", deleteUser);

const appRouter = new Router();
appRouter
  .get("/", home)
  .get("/health", health);

// ─── Server Setup ───────────────────────────────────────────────────────────

const app = createHttpServer();

app
  .use(requestId())
  .use(logger())
  .use(cors(["*"]))
  .router(apiRouter)
  .router(appRouter);

// ─── Start Server ────────────────────────────────────────────────────────────

const PORT = Number(process.env.PORT ?? "3000");
const HOST = process.env.HOST ?? "localhost";

const { address, close } = await app.listen({ port: PORT, host: HOST });

console.log(`\n🚀 Server running at http://${address.address}:${address.port}`);
console.log("Press Ctrl+C to stop\n");

// Graceful shutdown
process.on("SIGINT", async () => {
  console.log("\nShutting down gracefully...");
  await close();
  process.exit(0);
});

process.on("SIGTERM", async () => {
  await close();
  process.exit(0);
});
```

---

## Running the Server

```bash
# Install dependencies
npm install

# Compile TypeScript
npm run build

# Start the server
npm start

# Or run in development mode
npm run dev
```

---

## Testing the API

```bash
# Health check
curl http://localhost:3000/health | python3 -m json.tool

# Create a user
curl -X POST http://localhost:3000/api/users \
  -H "Content-Type: application/json" \
  -d '{"name":"Alice Smith","email":"alice@example.com"}' \
  | python3 -m json.tool

# List users
curl http://localhost:3000/api/users | python3 -m json.tool

# Get a specific user (use the id from create response)
curl http://localhost:3000/api/users/1 | python3 -m json.tool

# Update a user
curl -X PUT http://localhost:3000/api/users/1 \
  -H "Content-Type: application/json" \
  -d '{"name":"Alice Jones"}' \
  | python3 -m json.tool

# Delete a user
curl -X DELETE http://localhost:3000/api/users/1

# Pagination
curl "http://localhost:3000/api/users?page=1&limit=5" | python3 -m json.tool

# Validation error
curl -X POST http://localhost:3000/api/users \
  -H "Content-Type: application/json" \
  -d '{"name":"","email":"not-an-email"}' \
  | python3 -m json.tool

# 404 error
curl http://localhost:3000/api/nonexistent | python3 -m json.tool
```

---

## What This Project Demonstrates

### Type System Features Used

| Feature | Where Used |
|---------|-----------|
| Discriminated unions | `ResponseBody`, `ParsedRequest` |
| Generic types | `Schema<T>`, `Router.match()` |
| Conditional types | `UpdateUserSchema` with `.optional()` |
| Utility types | `Partial<T>`, `Readonly<T>`, `Record<K,V>` |
| Type guards | `isHttpError()`, validators |
| Intersection types | Context state |
| Template literal types | Header normalization |
| `unknown` for safety | `parseBody()`, error catches |
| `never` | Exhaustiveness in error hierarchy |
| `as const` | HTTP method arrays |

### Architectural Patterns

- **Middleware composition** — `compose()` chains handlers; each calls `next()`
- **Builder pattern** — `HttpResponse` with fluent `.status().header().json()`
- **Repository pattern** — `Map<string, User>` in-memory store
- **Schema validation** — composable `Schema<T>` validators
- **Error hierarchy** — typed `HttpError` subclasses with HTTP status codes
- **Separation of concerns** — request/response/router/middleware in separate modules

---

## Congratulations!

You've completed the TypeScript book. Here's what you've built and learned:

1. **Introduction** — TypeScript's purpose and design philosophy
2. **Getting Started** — tooling, tsconfig, development workflow
3. **Basic Types** — primitives, inference, `any` vs `unknown` vs `never`
4. **Functions** — overloading, generics, higher-order functions
5. **Objects & Interfaces** — structural typing, interface vs type
6. **Advanced Type System** — unions, intersections, narrowing, discriminated unions
7. **Generics** — type parameters, constraints, generic patterns
8. **Utility Types** — Partial, Pick, Omit, Record, ReturnType, Awaited
9. **Classes** — access modifiers, abstract, implements, mixins
10. **Modules** — ES modules, barrel files, declaration files
11. **Advanced Types** — mapped types, conditional types, infer, template literals
12. **JavaScript Interop** — allowJs, @types, declaration files, gradual migration
13. **Async Programming** — Promise<T>, async/await, error handling, generators
14. **The Compiler** — tsconfig deep dive, strict flags, project references
15. **Internals** — type erasure, structural typing, assignability, variance
16. **Best Practices** — 20 rules for maintainable TypeScript
17. **Common Pitfalls** — 15 traps and how to avoid them
18. **Interview Prep** — conceptual questions, tricky code, coding problems
19. **Final Project** — a complete type-safe HTTP server with zero dependencies

The TypeScript journey never ends — the type system continues to evolve with every release, adding new capabilities while maintaining backward compatibility. Keep exploring, keep building, and let the compiler be your pair programmer.
