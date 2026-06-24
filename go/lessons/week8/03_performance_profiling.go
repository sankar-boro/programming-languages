/*
WEEK 8 — DAY 3: Performance and Profiling
==========================================
Topic: Making Go code fast — benchmarking, profiling, and common optimizations.

Key ideas:
  - Measure before optimizing: "premature optimization is the root of all evil"
  - go test -bench and pprof are the primary tools
  - The most impactful optimizations are usually algorithmic, then allocation-reducing
  - Compiler optimizations: inlining, escape analysis, bounds check elimination
  - GOGC, GOMEMLIMIT, GOMAXPROCS are runtime tuning knobs

Go performance toolchain:
  go test -bench=.               — run benchmarks
  go test -bench=. -benchmem    — show allocations per op
  go test -bench=. -cpuprofile=cpu.out && go tool pprof cpu.out
  go test -bench=. -memprofile=mem.out && go tool pprof mem.out
  go tool trace trace.out        — goroutine execution trace
*/

package main

import (
	"fmt"
	"math"
	"math/bits"
	"runtime"
	"sort"
	"strings"
	"sync"
	"sync/atomic"
	"time"
	"unsafe"
)

// ─── 1. BENCHMARKING PATTERNS ─────────────────────────────────────────────────
//
// Benchmarks reveal actual performance — intuition is unreliable.
// Three things a benchmark measures: time/op, B/op (bytes), allocs/op.
//
// Rules for good benchmarks:
//   1. Put work inside the b.N loop
//   2. Use b.ResetTimer() if setup is expensive
//   3. Prevent the compiler from optimizing away results (use a sink variable)
//   4. Use -count=5 to stabilize results
//   5. Use benchstat tool to compare benchmark runs

// Sink prevents the compiler from eliminating "dead" computations
var Sink interface{}

func benchmarkingPatterns() {
	fmt.Println("=== 1. Benchmarking Patterns ===")
	fmt.Println(`
func BenchmarkStringConcat(b *testing.B) {
    words := []string{"hello", "world", "go", "is", "fast"}
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        var result string
        for _, w := range words {
            result += w + " "    // BAD: O(n²) allocations
        }
        Sink = result  // prevent dead-code elimination
    }
}

func BenchmarkStringsBuilder(b *testing.B) {
    words := []string{"hello", "world", "go", "is", "fast"}
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        var sb strings.Builder
        sb.Grow(32)              // GOOD: pre-allocate
        for _, w := range words {
            sb.WriteString(w)
            sb.WriteByte(' ')
        }
        Sink = sb.String()
    }
}

// Example output with -benchmem:
// BenchmarkStringConcat-8      1000000     1234 ns/op    432 B/op    6 allocs/op
// BenchmarkStringsBuilder-8   10000000      123 ns/op     32 B/op    1 allocs/op
`)

	// Live demo: measure string concatenation vs Builder
	n := 1000
	words := make([]string, n)
	for i := range words { words[i] = "word" }

	start := time.Now()
	s := ""
	for _, w := range words { s += w }
	concatTime := time.Since(start)

	start = time.Now()
	var sb strings.Builder
	sb.Grow(n * 4)
	for _, w := range words { sb.WriteString(w) }
	_ = sb.String()
	builderTime := time.Since(start)

	fmt.Printf("concat (%d words): %v\n", n, concatTime)
	fmt.Printf("Builder (%d words): %v\n", n, builderTime)
}

// ─── 2. COMMON ALLOCATION HOT SPOTS ─────────────────────────────────────────

func allocationHotSpots() {
	fmt.Println("\n=== 2. Common Allocation Hot Spots ===")

	// HOT SPOT 1: Growing slices without capacity hint
	fmt.Println("1. Slice growth:")
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	mallocs1 := m.Mallocs

	// BAD
	var s []int
	for i := 0; i < 10000; i++ { s = append(s, i) }
	runtime.ReadMemStats(&m)
	fmt.Printf("   without cap: %d allocations\n", m.Mallocs-mallocs1)

	runtime.ReadMemStats(&m)
	mallocs2 := m.Mallocs
	// GOOD
	s2 := make([]int, 0, 10000)
	for i := 0; i < 10000; i++ { s2 = append(s2, i) }
	runtime.ReadMemStats(&m)
	fmt.Printf("   with cap:    %d allocations\n", m.Mallocs-mallocs2)
	_ = s; _ = s2

	// HOT SPOT 2: fmt.Sprintf for simple conversions
	fmt.Println("\n2. String conversion:")
	fmt.Println("   BAD:  fmt.Sprintf(\"%d\", n)     → 1 allocation")
	fmt.Println("   GOOD: strconv.Itoa(n)            → 0 extra allocations")

	// HOT SPOT 3: Interface boxing
	fmt.Println("\n3. Interface boxing:")
	type Adder interface{ Add(int) int }
	// Each call to an interface method is cheap, but BOXING is expensive
	// (assigning a concrete value to an interface allocates when value doesn't fit in pointer)
	fmt.Println("   Storing small values in interface{} boxes them to heap")

	// HOT SPOT 4: Closures that capture variables
	fmt.Println("\n4. Closure captures:")
	fmt.Println("   Closures that capture variables force those vars to heap")
	fmt.Println("   Use: go build -gcflags=\"-m\" to see 'moved to heap' messages")

	// HOT SPOT 5: Map with frequent insertions
	fmt.Println("\n5. Map growth:")
	fmt.Println("   Like slices, pre-size maps: make(map[K]V, expectedSize)")
}

// ─── 3. CPU PROFILING WITH pprof ──────────────────────────────────────────────
//
// pprof is the standard Go profiler. It profiles:
//   - CPU time (where is the program spending time?)
//   - Heap allocations (what is allocating memory?)
//   - Goroutine stacks (where are goroutines blocked?)
//   - Mutex contention (where are goroutines waiting for locks?)
//
// Two ways to profile:
//   1. Test-based: go test -cpuprofile=cpu.out; go tool pprof cpu.out
//   2. HTTP: import _ "net/http/pprof"; go to :6060/debug/pprof

func pprofGuide() {
	fmt.Println("\n=== 3. CPU Profiling with pprof ===")
	fmt.Println(`
Method 1 — benchmark profiling:
  go test -bench=BenchmarkFoo -cpuprofile=cpu.out
  go tool pprof cpu.out

  (pprof) top10      — top 10 functions by CPU time
  (pprof) list Foo   — annotated source code for function Foo
  (pprof) web        — open flamegraph in browser (requires graphviz)

Method 2 — HTTP profiling (production servers):
  import _ "net/http/pprof"
  // Then in main: go http.ListenAndServe(":6060", nil)

  Endpoints:
    /debug/pprof/          — index
    /debug/pprof/heap      — heap profile
    /debug/pprof/goroutine — goroutine stacks
    /debug/pprof/mutex     — mutex contention
    /debug/pprof/trace     — execution trace

  go tool pprof http://localhost:6060/debug/pprof/heap
  go tool trace trace.out   — visual goroutine timeline

Reading pprof output:
  flat:   time/allocs IN this function
  cum:    time/allocs in this function AND everything it calls
  Look for high 'flat' values — those are the bottlenecks.
`)
}

// ─── 4. INLINING AND ESCAPE ANALYSIS OPTIMIZATION ────────────────────────────
//
// The compiler automatically inlines small functions (avoids function call overhead).
// You can see inlining decisions with: go build -gcflags="-m"
//
// "can inline" = inlining will happen
// "inlining call to X" = X was inlined at this call site
// "does not inline" = function is too complex to inline
//
// Functions that prevent inlining:
//   - closures that capture variables
//   - goroutines started within the function
//   - panic/recover
//   - very long functions (budget-based)
//
// You can force no-inlining: //go:noinline
// You can check inlining: //go:nosplit (different thing but shows size)

func add(a, b int) int { return a + b }  // will be inlined

//go:noinline
func addNoInline(a, b int) int { return a + b }  // won't be inlined

func inliningDemo() {
	fmt.Println("\n=== 4. Inlining ===")
	fmt.Println("go build -gcflags=\"-m\" shows: 'can inline add'")
	fmt.Println("Inlining eliminates function call overhead and enables further optimization")

	// Demonstrate that inlining works correctly
	result1 := add(3, 4)         // likely inlined to just `7`
	result2 := addNoInline(3, 4) // always a function call
	fmt.Printf("add(3,4)=%d, addNoInline(3,4)=%d\n", result1, result2)

	// Bounds check elimination
	// The compiler eliminates bounds checks when it can prove the index is valid
	s := []int{1, 2, 3, 4, 5}
	if len(s) >= 3 {
		// Compiler knows s[0], s[1], s[2] are valid — no bounds checks
		fmt.Printf("No bounds check: %d %d %d\n", s[0], s[1], s[2])
	}
}

// ─── 5. CACHE-FRIENDLY DATA STRUCTURES ───────────────────────────────────────
//
// Modern CPUs are much faster than RAM. Cache misses are expensive.
// L1 cache hit: ~4 cycles. RAM access: ~300 cycles.
//
// Principles:
//   1. Access memory sequentially (arrays >> linked lists for iteration)
//   2. Keep hot data together (struct of arrays vs array of structs)
//   3. Avoid pointer chasing (pointers = cache misses)
//   4. False sharing: two goroutines writing to nearby cache lines (64 bytes apart)

type Point struct{ X, Y float64 }  // 16 bytes

// Array of Structs (AoS) — common but cache-unfriendly for X-only operations
type AoS struct{ Points []Point }

// Struct of Arrays (SoA) — cache-friendly when operating on one field at a time
type SoA struct {
	X []float64
	Y []float64
}

func cacheFriendlyDemo() {
	fmt.Println("\n=== 5. Cache-Friendly Data Structures ===")

	n := 1_000_000
	aos := AoS{Points: make([]Point, n)}
	soa := SoA{X: make([]float64, n), Y: make([]float64, n)}

	for i := range aos.Points {
		aos.Points[i] = Point{float64(i), float64(i) * 2}
		soa.X[i] = float64(i)
		soa.Y[i] = float64(i) * 2
	}

	// Compute sum of X values only
	start := time.Now()
	sumAoS := 0.0
	for _, p := range aos.Points { sumAoS += p.X }
	aosTime := time.Since(start)

	start = time.Now()
	sumSoA := 0.0
	for _, x := range soa.X { sumSoA += x }
	soaTime := time.Since(start)

	fmt.Printf("AoS (sum X): %v  result=%.0f\n", aosTime, sumAoS)
	fmt.Printf("SoA (sum X): %v  result=%.0f\n", soaTime, sumSoA)
	fmt.Printf("SoA speedup: %.2fx\n", float64(aosTime)/float64(soaTime))
	fmt.Println("SoA is faster: iterating X skips Y data (better cache utilization)")
}

// ─── 6. ATOMIC OPERATIONS VS MUTEX ───────────────────────────────────────────
//
// For simple counter updates, atomic ops are faster than mutexes.
// Mutex: ~20-30ns per lock/unlock pair
// Atomic: ~2-5ns per operation
//
// When to use atomic:
//   - Simple counter (Add, Load, Store, Swap, CAS)
//   - Single numeric value
//   - No complex invariants
//
// When to use mutex:
//   - Multiple fields that must be updated together
//   - Complex data structures
//   - Read-write patterns with many readers

type AtomicCounter struct{ n int64 }
func (c *AtomicCounter) Inc() { atomic.AddInt64(&c.n, 1) }
func (c *AtomicCounter) Get() int64 { return atomic.LoadInt64(&c.n) }

type MutexCounter struct{
	mu sync.Mutex
	n  int64
}
func (c *MutexCounter) Inc() { c.mu.Lock(); c.n++; c.mu.Unlock() }
func (c *MutexCounter) Get() int64 { c.mu.Lock(); defer c.mu.Unlock(); return c.n }

func atomicVsMutex() {
	fmt.Println("\n=== 6. Atomic vs Mutex Performance ===")

	const ops = 1_000_000
	const goroutines = 4

	// Atomic counter
	var ac AtomicCounter
	var wg sync.WaitGroup
	start := time.Now()
	for g := 0; g < goroutines; g++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for i := 0; i < ops/goroutines; i++ { ac.Inc() }
		}()
	}
	wg.Wait()
	atomicTime := time.Since(start)

	// Mutex counter
	var mc MutexCounter
	start = time.Now()
	for g := 0; g < goroutines; g++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for i := 0; i < ops/goroutines; i++ { mc.Inc() }
		}()
	}
	wg.Wait()
	mutexTime := time.Since(start)

	fmt.Printf("Atomic counter (%d ops, %d goroutines): %v (result=%d)\n",
		ops, goroutines, atomicTime, ac.Get())
	fmt.Printf("Mutex counter  (%d ops, %d goroutines): %v (result=%d)\n",
		ops, goroutines, mutexTime, mc.Get())
	fmt.Printf("Atomic speedup: %.2fx\n", float64(mutexTime)/float64(atomicTime))
}

// ─── 7. RUNTIME TUNING KNOBS ──────────────────────────────────────────────────
//
// Go provides several runtime environment variables for tuning:
//
//   GOGC=N          — heap growth % before GC (default 100)
//                     GOGC=200 = less frequent GC, more memory usage
//                     GOGC=50  = more frequent GC, less memory
//                     GOGC=off = disable GC (for benchmarking only!)
//
//   GOMEMLIMIT=N    — maximum heap size (Go 1.19+)
//                     GOMEMLIMIT=512MiB = GC more aggressively to stay under 512MB
//                     Useful for containerized environments
//
//   GOMAXPROCS=N    — number of OS threads that can run Go code in parallel
//                     Default: number of logical CPUs
//                     GOMAXPROCS=1 = single-threaded (useful for debugging races)
//
//   GOTRACEBACK=all — show all goroutine stacks on panic
//   GODEBUG=...     — enable debug output for gc, scheduler, etc.

func runtimeTuning() {
	fmt.Println("\n=== 7. Runtime Tuning ===")

	fmt.Printf("GOMAXPROCS:      %d (logical CPUs: %d)\n",
		runtime.GOMAXPROCS(0), runtime.NumCPU())
	fmt.Printf("NumGoroutine:    %d\n", runtime.NumGoroutine())

	fmt.Println("\nKey tuning variables:")
	fmt.Println("  GOGC=200      — less GC, more memory (latency-sensitive services)")
	fmt.Println("  GOGC=50       — more GC, less memory (memory-constrained environments)")
	fmt.Println("  GOMEMLIMIT=1GiB — absolute memory cap (Go 1.19+)")
	fmt.Println("  GOMAXPROCS=N  — parallelism limit")

	// Changing GOMAXPROCS programmatically
	old := runtime.GOMAXPROCS(0)
	fmt.Printf("\nCurrent GOMAXPROCS: %d\n", old)

	// Set to 1 for sequential execution (useful for debugging)
	// runtime.GOMAXPROCS(1)
	// ...
	// runtime.GOMAXPROCS(old)  // restore

	fmt.Println("\nExample GOGC/GOMEMLIMIT interaction:")
	fmt.Println("  Container with 512MiB limit:")
	fmt.Println("  GOMEMLIMIT=450MiB GOGC=100 go run server.go")
	fmt.Println("  → GC kicks in when heap approaches 450MiB")
}

// ─── 8. OPTIMIZATION CHECKLIST ────────────────────────────────────────────────

func optimizationChecklist() {
	fmt.Println("\n=== 8. Optimization Checklist ===")

	steps := []struct {
		step string
		desc string
	}{
		{"1. MEASURE", "go test -bench=. -benchmem — get a baseline first"},
		{"2. PROFILE", "go tool pprof cpu.out — find the actual bottleneck"},
		{"3. ALGORITHM", "O(n log n) vs O(n²) is 1000x faster at n=1M — fix this first"},
		{"4. ALLOCATIONS", "reduce allocs/op — fewer allocations = less GC pressure"},
		{"5. CACHE", "sequential memory access — arrays beat maps/linked lists"},
		{"6. CONCURRENCY", "use goroutines for I/O-bound, GOMAXPROCS for CPU-bound"},
		{"7. PRIMITIVES", "atomic ops for counters, RWMutex for read-heavy paths"},
		{"8. COMPILER", "inlining and escape analysis — small functions, avoid interface on hot paths"},
		{"9. TUNE", "GOGC, GOMEMLIMIT — only after all else is optimized"},
	}

	for _, s := range steps {
		fmt.Printf("  %-20s %s\n", s.step, s.desc)
	}

	fmt.Println("\nCommon micro-optimizations:")
	// Integer math tricks
	x := 42
	fmt.Printf("  x*2 = x<<1: %d = %d\n", x*2, x<<1)
	fmt.Printf("  x/2 = x>>1: %d = %d\n", x/2, x>>1)
	fmt.Printf("  x%%2 == 0 is even: %v\n", x&1 == 0)

	// Math/bits for popcount
	n := uint64(255)
	fmt.Printf("  bits.OnesCount64(%b) = %d\n", n, bits.OnesCount64(n))
	fmt.Printf("  bits.Len64(%d) = %d (position of highest bit)\n", n, bits.Len64(n))

	// Fast sqrt approximation check
	fmt.Printf("  math.Sqrt(9) = %g\n", math.Sqrt(9))

	// unsafe.Sizeof without reflection
	fmt.Printf("  unsafe.Sizeof(int64) = %d bytes\n", unsafe.Sizeof(int64(0)))
	fmt.Printf("  unsafe.Sizeof([100]byte) = %d bytes\n", unsafe.Sizeof([100]byte{}))

	// Sorting: sort.Search for binary search
	data := []int{1, 3, 5, 7, 9, 11}
	target := 7
	i := sort.SearchInts(data, target)
	fmt.Printf("  Binary search for %d: found at index %d\n", target, i)
}

func main() {
	benchmarkingPatterns()
	allocationHotSpots()
	pprofGuide()
	inliningDemo()
	cacheFriendlyDemo()
	atomicVsMutex()
	runtimeTuning()
	optimizationChecklist()
}

/*
THOUGHT QUESTIONS:

1. Why do you need a sink variable (`var Sink interface{}`) in benchmarks?
   What happens if the compiler proves a computation's result is never used?

2. Struct of Arrays (SoA) outperforms Array of Structs (AoS) when operating
   on one field. But when would AoS be better?

3. Atomic operations are faster than mutexes for counters. But they only work
   on a single memory location at a time. What does this mean for a struct
   with two fields that must always be consistent with each other?

4. GOGC=off disables GC entirely. When is this appropriate, and what are
   the risks?

5. pprof's flat vs cumulative time: if a function has high 'cum' but low 'flat',
   what does that tell you? What should you look at next?

EXERCISES:

1. Write a benchmark comparing three map implementations for frequency counting:
   map[string]int, sync.Map, and a shard-locked map (array of N maps).
   Run with -race -benchmem.

2. Implement the Sieve of Eratosthenes two ways:
   (a) Using a map[int]bool
   (b) Using a []bool (bitset)
   Benchmark both for n=1,000,000 and explain the difference.

3. Write a CPU-profiled benchmark of a JSON marshal/unmarshal loop.
   Identify the top 3 hot functions in the pprof output.
   Suggest optimizations.

4. Implement a false sharing demonstration: have two goroutines increment
   two adjacent fields in a struct. Then pad the struct to separate them
   to different cache lines (64 bytes apart) and measure the speedup.
*/

/*
─────────────────────────────────────────────────────────────────
COMPLETE GO BOOK — SUMMARY OF ALL 24 FILES
─────────────────────────────────────────────────────────────────

WEEK 1 — Core Language Foundations
  01_how_go_works.go       — compilation pipeline, runtime, init order
  02_syntax_values_types.go — 25 keywords, type system, zero values
  03_variables_memory_model.go — value semantics, escape analysis, iota

WEEK 2 — Types and Data Structures
  01_types_in_depth.go     — int overflow, IEEE 754, strings as bytes, any
  02_composite_types.go    — slice internals, map hash tables, append growth
  03_structs.go            — memory layout, padding, tags, embedding, options

WEEK 3 — Control Flow and Pointers
  01_control_flow.go       — for loop, defer LIFO, switch, type switch, goto
  02_pointers_in_depth.go  — mutation vs copy, nil safety, unsafe.Pointer
  03_error_handling.go     — error interface, wrapping, Is/As, panic/recover

WEEK 4 — Functions and Closures
  01_functions.go          — named returns, variadic, first-class, no TCO
  02_closures.go           — heap escape, factory, loop gotcha, goroutine closures
  03_higher_order_functions.go — generic Map/Filter/Reduce, Pipe, lazy thunks

WEEK 5 — Interfaces and Goroutines
  01_call_stack_goroutine.go — stack copying, GMP scheduler, 10k goroutines
  02_interfaces_in_depth.go  — iface=(type,value), nil gotcha, io composability
  03_goroutines_channels.go  — GMP, pipelines, fan-out, select, cancellation

WEEK 6 — Advanced Features
  01_recursion.go          — memoization, BST, mutual recursion, no TCO
  02_generics.go           — type params, constraints, Stack[T], Result[T]
  03_sync_package.go       — Mutex, RWMutex, WaitGroup, Once, atomic, sync.Map

WEEK 7 — Production Go
  01_modules_packages.go   — visibility, imports, init(), go.mod, internal, build tags
  02_standard_library.go   — fmt, strings, strconv, os, bufio, time, json, sort
  03_testing.go            — table-driven tests, benchmarks, examples, -race, coverage

WEEK 8 — Internals and Performance
  01_memory_gc.go          — tri-color GC, escape analysis, sync.Pool, MemStats
  02_concurrency_patterns.go — context, errgroup, worker pool, rate limiter
  03_performance_profiling.go — pprof, cache, atomic vs mutex, GOGC, tuning

Total: 24 runnable Go files, ~5000 lines of code and documentation
─────────────────────────────────────────────────────────────────
*/
