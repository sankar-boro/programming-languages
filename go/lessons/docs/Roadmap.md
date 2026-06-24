# Go Mastery — Full 6-Month Roadmap

---

## PHASE 1 — Core Language Foundations (Weeks 1–4)

| Week | Topics |
|------|--------|
| 1 | How Go works, the compilation model, Go runtime, `go build` pipeline |
| 2 | Types in depth — basic types, zero values, type system, `reflect` |
| 3 | Variables, constants, short declarations, pointers and memory model |
| 4 | Control flow — `if`, `for`, `switch`, `defer` — how they execute internally |

---

## PHASE 2 — Deep Dive into Functions (Weeks 5–8)

| Week | Topics |
|------|--------|
| 5 | Functions — signatures, multiple return values, named returns, call stack |
| 6 | Closures — how they capture variables, heap escape, practical patterns |
| 7 | Recursion — mechanics, stack frames, tail call optimization (Go has none) |
| 8 | Higher-order functions, function types, `func` as first-class values |

---

## PHASE 3 — Composite Types & Data Structures (Weeks 9–12)

| Week | Topics |
|------|--------|
| 9  | Arrays and slices — internal representation (`ptr`, `len`, `cap`), `append` mechanics |
| 10 | Maps — hash table internals, key requirements, iteration order |
| 11 | Structs — memory layout, embedding, anonymous fields, tags |
| 12 | Pointers — stack vs heap, escape analysis, `new` vs `&`, unsafe |

---

## PHASE 4 — Interfaces & OOP in Go (Weeks 13–17)

| Week | Topics |
|------|--------|
| 13 | Methods — value receivers vs pointer receivers, method sets |
| 14 | Interfaces — implicit implementation, duck typing, interface internals (type + value pair) |
| 15 | Composition over inheritance — embedding interfaces, struct embedding |
| 16 | Empty interface `any`/`interface{}`, type assertions, type switches |
| 17 | Error handling — `error` interface, sentinel errors, `errors.Is`/`As`, custom errors |

---

## PHASE 5 — Concurrency (Weeks 18–21)

| Week | Topics |
|------|--------|
| 18 | Goroutines — OS threads vs goroutines, the Go scheduler (GMP model) |
| 19 | Channels — unbuffered vs buffered, `select`, `close`, direction constraints |
| 20 | `sync` package — `Mutex`, `RWMutex`, `WaitGroup`, `Once`, `Map` |
| 21 | Race conditions, `go test -race`, data race patterns and prevention |

---

## PHASE 6 — Advanced Internals (Weeks 22–26)

| Week | Topics |
|------|--------|
| 22 | Modules and packages — `go.mod`, versioning, internal packages, build tags |
| 23 | The Go runtime — stack growth, goroutine preemption, GOMAXPROCS |
| 24 | Memory and GC — tri-color mark-sweep, write barriers, `runtime/pprof` |
| 25 | Generics — type parameters, constraints, `comparable`, when to use |
| 26 | `reflect`, `unsafe`, `cgo` — how Go interops with the outside world |

---
