/*
WEEK 8 — DAY 1: Memory Management and the Garbage Collector
=============================================================
Topic: How Go manages memory — the GC algorithm, escape analysis,
       heap vs stack, and how to profile memory usage.

Key ideas:
  - Go uses a concurrent tri-color mark-sweep garbage collector
  - Escape analysis determines stack vs heap allocation (compiler decides)
  - sync.Pool can reduce GC pressure for short-lived allocations
  - runtime/pprof and runtime.MemStats let you inspect memory usage
  - GOGC controls GC frequency (default: 100 = double the heap before GC)
*/

package main

import (
	"fmt"
	"runtime"
	"sync"
	"time"
)

// ─── 1. STACK vs HEAP: THE COMPILER DECIDES ───────────────────────────────────
//
// The Go compiler uses ESCAPE ANALYSIS to decide where to allocate:
//
//   Stack allocation:
//     - No GC involvement — freed when the function returns
//     - Very fast (just adjusts the stack pointer)
//     - Size must be known at compile time (usually)
//     - Variable must not outlive the function
//
//   Heap allocation:
//     - GC managed — freed when no references remain
//     - Slower (involves the allocator and eventually the GC)
//     - Required when: variable outlives the function, size unknown, interface boxing
//
// See escape analysis decisions:
//   go build -gcflags="-m" main.go

func stackExample() int {
	x := 42        // on the stack — never escapes
	y := x + 1     // on the stack
	return y        // copied to caller; x and y freed when function returns
}

func heapExample() *int {
	x := 42        // x ESCAPES to heap — we return its address
	return &x       // x must outlive this function → heap allocated
}

func interfaceEscape(v interface{}) {
	// v is stored in an interface — typically heap allocated
	fmt.Println(v)
}

func escapeAnalysisDemo() {
	fmt.Println("=== 1. Stack vs Heap (Escape Analysis) ===")
	fmt.Println("Run: go build -gcflags=\"-m\" to see escape decisions\n")

	v := stackExample()
	p := heapExample()
	fmt.Printf("Stack value: %d\n", v)
	fmt.Printf("Heap pointer: %d at %p\n", *p, p)

	interfaceEscape(42)  // 42 is boxed and escapes to heap

	// Slices escape to heap when their size is unknown or they grow
	n := 5
	s := make([]int, n)  // might escape — depends on how it's used
	for i := range s { s[i] = i }
	fmt.Println("Slice:", s)
}

// ─── 2. THE GC ALGORITHM: TRI-COLOR MARK-SWEEP ────────────────────────────────
//
// Go uses a concurrent, tri-color mark-sweep GC:
//
// Three colors for heap objects:
//   White: not yet visited — garbage candidates
//   Grey:  discovered but not fully processed — in the work queue
//   Black: fully processed — retained, children also processed
//
// Algorithm:
//   1. STW (stop-the-world): mark all root objects grey (globals, goroutine stacks)
//   2. CONCURRENT mark: goroutines run while GC processes grey objects
//      - Pop grey object → mark black → mark its children grey
//      - Write barrier: track any pointers written during marking (Dijkstra)
//   3. STW: final mark pass (handle any mutations during marking)
//   4. SWEEP: reclaim all white (unmarked) objects
//   5. Repeat when heap doubles (GOGC=100 means: GC when heap = 2x live set)
//
// Go's GC pauses have been reduced to <1ms in modern versions.
// Most pause time is "STW" for scanning goroutine stacks.

func gcAlgorithm() {
	fmt.Println("\n=== 2. GC Algorithm (Tri-Color Mark-Sweep) ===")

	// Force GC and print stats
	runtime.GC()

	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	fmt.Printf("Alloc (live heap):     %d KB\n", m.Alloc/1024)
	fmt.Printf("TotalAlloc (lifetime): %d KB\n", m.TotalAlloc/1024)
	fmt.Printf("Sys (from OS):         %d KB\n", m.Sys/1024)
	fmt.Printf("NumGC (GC runs):       %d\n", m.NumGC)
	fmt.Printf("PauseTotal:            %v\n", time.Duration(m.PauseTotalNs))
	fmt.Printf("GOGC setting:          %d\n", readGOGC())

	// Show pause times for last few GC cycles
	fmt.Println("\nRecent GC pause times:")
	numPauses := int(m.NumGC)
	if numPauses > 10 { numPauses = 10 }
	for i := 0; i < numPauses; i++ {
		idx := (int(m.NumGC)+255-i) % 256  // circular buffer
		if m.PauseNs[idx] > 0 {
			fmt.Printf("  GC %d pause: %v\n", int(m.NumGC)-i, time.Duration(m.PauseNs[idx]))
		}
	}
}

func readGOGC() int {
	// GOGC defaults to 100 — GC when heap grows by 100% of live set
	// runtime.GOGC() is available in Go 1.19+
	return 100 // default
}

// ─── 3. ALLOCATION PRESSURE — HOW TO REDUCE IT ────────────────────────────────
//
// More allocations = more GC work = more pause time.
// Techniques to reduce allocations:
//
//   1. Reuse buffers / objects with sync.Pool
//   2. Pass by value (avoid pointer when value is small)
//   3. Pre-allocate slices with known capacity
//   4. Avoid interface boxing on hot paths
//   5. Use []byte instead of strings for mutation

func allocationPressure() {
	fmt.Println("\n=== 3. Reducing Allocation Pressure ===")

	// BAD: allocates a new buffer every call
	makeBufBad := func() []byte {
		return make([]byte, 1024)
	}

	// GOOD: reuse buffers from pool
	pool := &sync.Pool{
		New: func() interface{} {
			buf := make([]byte, 1024)
			return &buf
		},
	}

	getBufGood := func() *[]byte {
		return pool.Get().(*[]byte)
	}
	returnBuf := func(buf *[]byte) {
		// Clear the buffer before returning (security)
		for i := range *buf { (*buf)[i] = 0 }
		pool.Put(buf)
	}

	// Measure allocations
	before := allocCount()
	for i := 0; i < 100; i++ {
		buf := makeBufBad()
		_ = buf
	}
	afterBad := allocCount()

	before2 := allocCount()
	for i := 0; i < 100; i++ {
		buf := getBufGood()
		returnBuf(buf)
	}
	afterGood := allocCount()

	fmt.Printf("Without pool: %d allocations\n", afterBad-before)
	fmt.Printf("With pool:    %d allocations\n", afterGood-before2)

	// Pre-allocate slices
	var allocsThen, allocsNow uint64
	runtime.ReadMemStats(&struct{ Alloc uint64 }{allocsThen})

	// BAD: grows repeatedly
	s1 := []int{}
	for i := 0; i < 1000; i++ {
		s1 = append(s1, i)  // multiple reallocations
	}

	// GOOD: pre-allocate
	s2 := make([]int, 0, 1000)  // known capacity
	for i := 0; i < 1000; i++ {
		s2 = append(s2, i)  // NO reallocations
	}
	_ = allocsNow; _ = allocsThen; _ = s1; _ = s2
	fmt.Println("Pre-allocated slice avoids repeated reallocations")
}

func allocCount() uint64 {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	return m.Mallocs
}

// ─── 4. sync.Pool — OBJECT REUSE ─────────────────────────────────────────────
//
// sync.Pool is a cache of objects that can be reused.
// Objects are NOT guaranteed to persist between GC cycles.
// The GC may clear the pool at any time.
//
// Use for: temporary buffers, encoders, parsers — objects that are:
//   - Expensive to create
//   - Short-lived (created, used, discarded)
//   - Needed from many goroutines
//
// Don't use for: connection pools, persistent objects, or items that must
// survive GC cycles (use a channel-based pool instead).

type BigBuffer struct {
	data [64 * 1024]byte  // 64 KB buffer
}

var bufPool = &sync.Pool{
	New: func() interface{} {
		fmt.Println("  [pool: allocating new buffer]")
		return &BigBuffer{}
	},
}

func processwithPool(data []byte) {
	buf := bufPool.Get().(*BigBuffer)
	defer bufPool.Put(buf)

	copy(buf.data[:], data)
	// ... process buf.data ...
}

func syncPoolDemo() {
	fmt.Println("\n=== 4. sync.Pool ===")

	// First three calls might allocate
	processwithPool([]byte("request 1"))
	processwithPool([]byte("request 2"))
	processwithPool([]byte("request 3"))

	fmt.Println("Buffers returned to pool — subsequent calls reuse them")
	processwithPool([]byte("request 4"))  // may reuse from pool
	processwithPool([]byte("request 5"))  // may reuse from pool

	// After GC, pool may be cleared
	runtime.GC()
	fmt.Println("After GC — pool may be empty:")
	processwithPool([]byte("request 6"))  // may allocate again
}

// ─── 5. MEMORY PROFILING ──────────────────────────────────────────────────────
//
// Two ways to profile memory:
//
//   1. runtime.MemStats: a snapshot of current GC statistics
//      Call runtime.ReadMemStats(&m) to fill the struct.
//
//   2. pprof: production-grade profiling
//      import _ "net/http/pprof"  in your main, then:
//      go tool pprof http://localhost:6060/debug/pprof/heap
//
// Key MemStats fields:
//   Alloc:       bytes currently allocated on the heap (live objects)
//   TotalAlloc:  cumulative bytes allocated (includes freed)
//   HeapAlloc:   bytes allocated by the GC (same as Alloc)
//   HeapSys:     bytes obtained from OS for the heap
//   HeapInuse:   bytes in use by the GC (includes objects + their alignment)
//   HeapIdle:    bytes in idle spans (usable by GC but not currently)
//   Mallocs:     cumulative number of heap allocations
//   Frees:       cumulative number of heap frees
//   NumGC:       number of GC cycles completed
//   PauseNs:     circular buffer of recent GC pause durations

func memoryProfiling() {
	fmt.Println("\n=== 5. Memory Profiling ===")

	printMemStats := func(label string) {
		var m runtime.MemStats
		runtime.ReadMemStats(&m)
		fmt.Printf("[%s] Alloc=%dKB TotalAlloc=%dKB Mallocs=%d NumGC=%d\n",
			label, m.Alloc/1024, m.TotalAlloc/1024, m.Mallocs, m.NumGC)
	}

	printMemStats("start")

	// Allocate a bunch of objects
	var objects []*[]byte
	for i := 0; i < 10000; i++ {
		b := make([]byte, 1024)
		objects = append(objects, &b)
	}
	printMemStats("after 10k allocs")

	// Release references — objects become eligible for GC
	objects = nil
	runtime.GC()
	printMemStats("after GC")

	fmt.Println("\nFor production profiling:")
	fmt.Println("  import _ \"net/http/pprof\"")
	fmt.Println("  go http.ListenAndServe(\":6060\", nil)")
	fmt.Println("  go tool pprof http://localhost:6060/debug/pprof/heap")
}

// ─── 6. FINALIZERS ────────────────────────────────────────────────────────────
//
// runtime.SetFinalizer attaches a function that runs when the GC reclaims
// an object. Use only as a safety net — don't rely on it for correctness.
//
// Problems with finalizers:
//   - Run time is non-deterministic
//   - May run in a different goroutine
//   - Two objects with finalizers cannot reference each other (cycle prevention)
//   - Delay GC by one extra cycle
//
// Prefer defer/close for resource cleanup. Use finalizers only as debugging aid.

type Resource struct {
	name string
}

func NewResource(name string) *Resource {
	r := &Resource{name: name}
	runtime.SetFinalizer(r, func(r *Resource) {
		fmt.Printf("  [finalizer] Resource %q collected by GC\n", r.name)
		// WARNING: don't do important cleanup here — not guaranteed to run
	})
	return r
}

func finalizersDemo() {
	fmt.Println("\n=== 6. Finalizers (use sparingly) ===")

	{
		r := NewResource("important")
		_ = r
		// r goes out of scope here — eligible for GC
	}

	// Force GC — finalizer MIGHT run
	runtime.GC()
	time.Sleep(10 * time.Millisecond)  // finalizers run in a background goroutine
	runtime.GC()

	fmt.Println("Prefer defer r.Close() over finalizers for resource cleanup")
}

func main() {
	escapeAnalysisDemo()
	gcAlgorithm()
	allocationPressure()
	syncPoolDemo()
	memoryProfiling()
	finalizersDemo()
}

/*
THOUGHT QUESTIONS:

1. What is escape analysis? Who performs it (programmer, compiler, or runtime)?
   Can you always predict which allocations go to the stack vs heap?

2. Go's GC is "concurrent" — what does that mean? Does it mean there are
   NO stop-the-world pauses? What requires STW?

3. sync.Pool objects may be discarded during GC. What does this mean for
   a connection pool? Can you use sync.Pool for database connections?

4. Why do finalizers NOT provide a reliable cleanup mechanism in Go?
   What is the correct pattern for releasing resources?

5. GOGC=100 means "trigger GC when heap grows to 2x the live set."
   What are the trade-offs of setting GOGC lower vs higher?

EXERCISES:

1. Write a benchmark that measures the allocation improvement from using
   sync.Pool for a JSON encoder/decoder. Run with -benchmem.

2. Write a program that allocates large objects repeatedly. Use runtime.MemStats
   to print allocation and GC statistics at each step. Observe when GC triggers.

3. Use go build -gcflags="-m" to analyze a function with several variables.
   For each variable, explain why it escaped to the heap or stayed on the stack.

4. Write a `Buffer` type that implements sync.Pool patterns internally —
   callers call Get() and must call Put() when done. Make it safe for
   concurrent use without external synchronization.
*/
