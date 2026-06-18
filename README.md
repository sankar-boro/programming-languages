# Programming Languages

A collection of in-depth books and courses on programming languages, written as structured Markdown and code files with runnable examples.

## Books

### Rust — `rust/`

Covers ownership, borrowing, lifetimes, traits, generics, concurrency, and internals. Final project: a complete HTTP/1.1 server built using only the `std` library. 20 chapters (`00`–`16`, `90`–`92`, `99`).

### TypeScript — `typescript/`

Covers the type system deeply, from basics through mapped types, conditional types, and template literal types. Final project: a type-safe HTTP server using only Node.js's built-in `http` module. 20 chapters (`00`–`15`, `90`–`92`, `99`).

### Android (AndroidX) — `android/`

A structured curriculum for modern Android development with AndroidX. Organized into four progressive levels with mini-projects at each level.

| Level | Focus |
|-------|-------|
| `level-1-beginner` | AndroidX setup, ConstraintLayout, RecyclerView, Material — mini project: Task List app |
| `level-2-intermediate` | ViewModel, LiveData, StateFlow, Navigation, Room, Jetpack Compose intro — mini project: Notes app |
| `level-3-advanced` | Hilt, WorkManager, Coroutines, DataStore, Testing, Paging 3 — mini project: News app |
| `level-4-expert` | Clean Architecture, modularization, offline-first, performance, CI/CD — mini project: E-commerce module |

Also includes `roadmap.md` (30/60/90/120-day learning plan) and a `capstone/` project.

### Python — `python/`

An 8-week course structured as runnable `.py` files, going deep on how Python works internally.

| Week | Topic |
|------|-------|
| 1 | How Python works, syntax, values, types, variables and memory model |
| 2 | Strings, numbers, booleans — deep internals |
| 3 | Control flow — if/else, for loops, while loops |
| 4 | Functions — definition, arguments, return values, scope |
| 5 | Call stack, execution frames, how Python runs code |
| 6 | Recursion — mechanics, vs iteration, stack depth |
| 7 | Closures, lexical scoping, nonlocal |
| 8 | Higher-order functions, `*args`/`**kwargs`, pure vs impure functions |

---

## Other Languages

Folders exist for: `c`, `clojure`, `cobol`, `cpp`, `csharp`, `dart`, `elixir`, `fortran`, `go`, `haskell`, `java`, `javascript`, `julia`, `kotlin`, `lua`, `matlab`, `perl`, `php`, `r`, `ruby`, `scala`, `swift`.

---

## Chapter Format (Rust & TypeScript books)

Each chapter follows this structure:
- Conceptual explanation with diagrams where relevant
- Runnable code examples with inline comments
- **Summary** — key ideas in plain English
- **Key Takeaways** — bullet points for quick review
- **Practice Questions** — test your understanding
- **Exercises** — hands-on coding problems
