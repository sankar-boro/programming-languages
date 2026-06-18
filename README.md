# Programming Languages

A collection of in-depth books on programming languages, written as structured Markdown files with runnable code examples.

## Books

### Rust
**20 chapters** — `rust/`

Covers ownership, borrowing, lifetimes, traits, generics, concurrency, and internals. Final project: a complete HTTP/1.1 server built using only the `std` library.

| File | Topic |
|------|-------|
| `00` | Title, Preface, Table of Contents |
| `01` | Introduction — why Rust, who it's for |
| `02` | Getting Started — toolchain, Cargo, Hello World |
| `03` | Basic Syntax — variables, types, stack vs heap |
| `04` | Control Flow — if, match, loops, pattern matching |
| `05` | Functions — closures, expressions vs statements |
| `06` | Ownership — the three rules, move semantics |
| `07` | Borrowing — shared vs exclusive references, slices |
| `08` | Lifetimes — annotations, elision, `'static` |
| `09` | Structs & Enums — ADTs, Option, Result |
| `10` | Error Handling — panic, Result, `?` operator |
| `11` | Generics & Traits — monomorphization, trait objects |
| `12` | Collections — Vec, HashMap, iterators |
| `13` | Advanced Features — closures, Box, Rc, RefCell, Arc |
| `14` | Concurrency — threads, channels, Mutex, Send/Sync |
| `15` | Modules & Crates — visibility, Cargo workspaces |
| `16` | Internals — MIR, vtables, zero-cost abstractions |
| `90` | Best Practices |
| `91` | Common Pitfalls |
| `92` | Interview Preparation |
| `99` | Final Project: HTTP server (std only) |

---

### TypeScript
**20 chapters** — `typescript/`

Covers the type system deeply, from basics through mapped types, conditional types, and template literal types. Final project: a type-safe HTTP server using only Node.js's built-in `http` module.

| File | Topic |
|------|-------|
| `00` | Title, Preface, Table of Contents |
| `01` | Introduction — TypeScript vs JavaScript, structural typing |
| `02` | Getting Started — tsconfig, ts-node, project setup |
| `03` | Basic Types — primitives, `any`, `unknown`, `never`, literals |
| `04` | Functions — overloading, generics, higher-order functions |
| `05` | Objects & Interfaces — `interface` vs `type`, readonly, structural typing |
| `06` | Advanced Type System — unions, narrowing, discriminated unions |
| `07` | Generics — type parameters, constraints, `keyof` |
| `08` | Utility Types — Partial, Pick, Omit, Record, ReturnType, Awaited |
| `09` | Classes — access modifiers, abstract, implements, mixins |
| `10` | Modules — ES modules, barrel files, declaration files |
| `11` | Advanced Types — mapped types, conditional types, `infer`, template literals |
| `12` | JS Interop — `allowJs`, `@types`, gradual migration |
| `13` | Async Programming — Promise, async/await, generators |
| `14` | The Compiler — tsconfig deep dive, strict flags, project references |
| `15` | Internals — type erasure, structural typing, assignability, variance |
| `90` | Best Practices |
| `91` | Common Pitfalls |
| `92` | Interview Preparation |
| `99` | Final Project: HTTP server (Node.js `http` module only) |

---

## Other Languages

Folders exist for: `android`, `c`, `clojure`, `cobol`, `cpp`, `csharp`, `dart`, `elixir`, `fortran`, `go`, `haskell`, `java`, `javascript`, `julia`, `kotlin`, `lua`, `matlab`, `perl`, `php`, `python`, `r`, `ruby`, `scala`, `swift`.

## Format

Each chapter follows this structure:
- Conceptual explanation with diagrams where relevant
- Runnable code examples with inline comments
- **Summary** — key ideas in plain English
- **Key Takeaways** — bullet points for quick review
- **Practice Questions** — test your understanding
- **Exercises** — hands-on coding problems
