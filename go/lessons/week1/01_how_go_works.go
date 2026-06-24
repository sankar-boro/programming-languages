/*
WEEK 1 — DAY 1: How Go Works
==============================
Topic: The Go compilation model, the Go runtime, and the execution pipeline.

Key ideas:
  - Go is a compiled language: source → machine code directly
  - The Go runtime is linked into every binary (GC, scheduler, goroutines)
  - Go has no VM, no interpreter — the binary is self-contained
  - go build, go run, and the toolchain pipeline
*/

package main

import (
	"fmt"
	"os"
	"runtime"
	"runtime/debug"
)

// ─── 1. THE COMPILATION PIPELINE ──────────────────────────────────────────────
//
// When you run: go build main.go
//
//   1. Lexer/Scanner — source text → tokens (keywords, identifiers, literals)
//   2. Parser         — tokens → Abstract Syntax Tree (AST)
//   3. Type checker   — walks the AST, resolves types, checks correctness
//   4. SSA IR         — converts AST into Static Single Assignment form
//   5. Backend        — SSA → machine code (amd64, arm64, etc.)
//   6. Linker         — combines packages + Go runtime → single binary
//
// Unlike Python/JVM, there is NO intermediate bytecode file you can inspect.
// The output is native machine code for your OS + architecture.
//
// go run main.go = compile to a temp dir + immediately execute the binary.
// It's NOT an interpreter. It's a compile-and-run shortcut.

func compilationPipeline() {
	fmt.Println("=== 1. Compilation Pipeline ===")
	fmt.Println("Go source → AST → SSA IR → machine code → binary")
	fmt.Println("The binary includes the Go runtime (GC, scheduler, stack management).")
	fmt.Println("No external runtime needed. The binary runs standalone.")
}

// ─── 2. THE GO RUNTIME ────────────────────────────────────────────────────────
//
// Every Go binary contains the Go runtime. The runtime provides:
//
//   - Garbage collector (tri-color mark-sweep, concurrent)
//   - Goroutine scheduler (M:N — many goroutines, few OS threads)
//   - Stack management (goroutine stacks start at 2KB, grow dynamically)
//   - Channel implementation
//   - Memory allocator
//   - panic / recover mechanism
//
// This is why Go binaries are larger than C binaries (~1–2 MB overhead).
// The tradeoff: you get GC, goroutines, and safety with zero external deps.

func goRuntime() {
	fmt.Println("\n=== 2. The Go Runtime ===")

	// What CPU architecture are we compiled for?
	fmt.Println("GOARCH:", runtime.GOARCH)   // e.g., "amd64", "arm64"
	fmt.Println("GOOS:", runtime.GOOS)       // e.g., "linux", "darwin", "windows"
	fmt.Println("NumCPU:", runtime.NumCPU()) // logical CPUs available
	fmt.Println("GOMAXPROCS:", runtime.GOMAXPROCS(0)) // OS threads Go uses (default = NumCPU)
	fmt.Println("NumGoroutine:", runtime.NumGoroutine()) // active goroutines right now

	// Runtime memory stats
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	fmt.Printf("HeapAlloc: %d bytes\n", m.HeapAlloc)
	fmt.Printf("HeapSys:   %d bytes\n", m.HeapSys)
}

// ─── 3. PACKAGES AND THE MAIN ENTRY POINT ────────────────────────────────────
//
// Every Go program is made of packages. Rules:
//
//   - Every file must declare its package at the top: `package <name>`
//   - The entry point of a Go program is always: func main() in package main
//   - There can only be ONE package main per binary
//   - Other packages are libraries (no main function)
//
// Go does NOT have a global scope beyond the package. All top-level identifiers
// belong to their package. Capitalized identifiers are exported; lowercase are not.

func packagesAndMain() {
	fmt.Println("\n=== 3. Packages and main ===")
	fmt.Println("This file is in: package main")
	fmt.Println("func main() is the program entry point")
	fmt.Println("Exported names start with uppercase: fmt.Println")
	fmt.Println("Unexported names start with lowercase: local to the package")
}

// ─── 4. THE BUILD SYSTEM ──────────────────────────────────────────────────────
//
// Go ships with a complete build system. No Makefile needed.
//
//   go build         → compile the current package into a binary
//   go run main.go   → compile + run immediately
//   go test          → compile + run tests
//   go vet           → static analysis for common mistakes
//   go fmt           → format code (enforces one style, no debates)
//   go doc           → show package documentation
//   go mod init      → initialize a module (go.mod)
//   go get           → add a dependency
//
// GCFLAGS — useful for understanding compilation:
//   go build -gcflags="-m" main.go   → shows escape analysis decisions
//   go build -gcflags="-S" main.go   → shows assembly output
//   go tool objdump binary           → disassemble the binary

func buildSystem() {
	fmt.Println("\n=== 4. The Build System ===")
	fmt.Println("go build   → native binary, zero external deps")
	fmt.Println("go run     → compile + execute (NOT an interpreter)")
	fmt.Println("go fmt     → enforces one canonical code style")
	fmt.Println("go vet     → static analysis, catches common bugs")
}

// ─── 5. BUILD INFO — WHAT'S IN YOUR BINARY ───────────────────────────────────
//
// Since Go 1.18, every binary embeds build metadata you can read at runtime.
// This includes the Go version used, module path, and VCS info.

func buildInfo() {
	fmt.Println("\n=== 5. Build Info ===")
	if info, ok := debug.ReadBuildInfo(); ok {
		fmt.Println("Go version:", info.GoVersion)
		fmt.Println("Module path:", info.Path)
		for _, dep := range info.Deps {
			fmt.Printf("  dep: %s@%s\n", dep.Path, dep.Version)
		}
	} else {
		fmt.Println("(build info only available in module-aware builds)")
	}
}

// ─── 6. OS ARGS AND THE PROGRAM BOUNDARY ──────────────────────────────────────
//
// When the OS executes a Go binary, it passes:
//   - os.Args[0]  = the path to the binary itself
//   - os.Args[1:] = command-line arguments
//
// Go's runtime initializes BEFORE main() runs:
//   1. Sets up the goroutine scheduler
//   2. Starts the GC background goroutine
//   3. Runs init() functions in all packages (in dependency order)
//   4. Calls main.main()

func programBoundary() {
	fmt.Println("\n=== 6. Program Boundary ===")
	fmt.Println("Binary path:", os.Args[0])
	fmt.Println("CLI args:", os.Args[1:])
	fmt.Println()
	fmt.Println("Execution order:")
	fmt.Println("  1. Runtime init (scheduler, GC)")
	fmt.Println("  2. Package-level var declarations (all packages)")
	fmt.Println("  3. init() functions (all packages, dependency order)")
	fmt.Println("  4. main.main()")
}

// ─── PACKAGE-LEVEL INIT EXAMPLE ───────────────────────────────────────────────
//
// init() runs automatically before main(). It's used for setup.
// You can have multiple init() functions in one file — they all run.

var packageVar = initPackageVar() // runs during package init

func initPackageVar() string {
	return "initialized at package load time"
}

func init() {
	// This runs before main()
	// Use for: registering drivers, setting up global state, validation
	fmt.Println("[init] Package initialized. packageVar =", packageVar)
}

func main() {
	compilationPipeline()
	goRuntime()
	packagesAndMain()
	buildSystem()
	buildInfo()
	programBoundary()
}

/*
THOUGHT QUESTIONS:

1. Why is `go run` not an interpreter? What actually happens when you call it?

2. What is the Go runtime, and why does every Go binary include it?

3. What is GOMAXPROCS? What happens if you set it to 1?

4. Why does Go enforce a single coding style (go fmt) rather than letting
   developers choose their own?

5. What runs before main()? In what order?

EXERCISES:

1. Run `go build -gcflags="-m" main.go` and observe which variables
   the compiler decides to allocate on the heap vs the stack.

2. Run `go tool objdump <binary> | head -100` to see the assembly
   of the first few functions. What do you notice?

3. Add two init() functions to this file. Predict what order they run in,
   then verify by running the program.
*/
