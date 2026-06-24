/*
WEEK 5 — DAY 1: The Call Stack and Goroutine Stack
====================================================
Topic: How Go manages the stack — goroutine stacks vs OS threads,
       stack growth, frame layout, and stack traces.

Key ideas:
  - Each goroutine has its OWN small stack (starts at 2–8 KB)
  - Goroutine stacks GROW dynamically (stack copying or segmented stacks)
  - Stack frames store: locals, parameters, return values, return address
  - runtime.Stack and debug.PrintStack show the goroutine stack trace
  - runtime.Callers / runtime.CallersFrames give you the call stack programmatically
*/

package main

import (
	"fmt"
	"runtime"
	"runtime/debug"
)

// ─── 1. GOROUTINE STACK BASICS ────────────────────────────────────────────────
//
// When a goroutine is created, Go allocates a small stack for it.
// As the goroutine calls functions, frames are pushed onto the stack.
// When functions return, their frames are popped.
//
// Go uses "stack copying" (introduced in Go 1.4):
//   1. When the stack is too small for a new frame, a NEW, larger stack is allocated
//   2. All existing data is COPIED to the new stack
//   3. Pointers ON the stack are updated to point to new locations
//
// This is why you cannot store the address of a Go stack variable in C code
// that might use it later — the stack could have been copied (moved).
//
// Stack growth is transparent — you don't need to manage it.

func goroutineStack() {
	fmt.Println("=== 1. Goroutine Stack ===")

	// Every goroutine starts here: a small stack
	// This function itself is a stack frame

	// Show current goroutine count
	fmt.Printf("Active goroutines: %d\n", runtime.NumGoroutine())

	// Force GC and print memory stats
	runtime.GC()
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	fmt.Printf("Stack in use: %d bytes (all goroutines)\n", m.StackInuse)
	fmt.Printf("Stack from OS: %d bytes\n", m.StackSys)
}

// ─── 2. STACK FRAMES AND PARAMETER PASSING ────────────────────────────────────
//
// When function A calls function B:
//
//   Stack grows downward (conceptually):
//   ┌─────────────────┐  ← top of stack
//   │ A's locals       │
//   │ A's return addr  │
//   ├─────────────────┤
//   │ B's parameters   │  ← A writes these
//   │ B's locals       │  ← B initializes these
//   │ B's return val   │  ← B writes return values here
//   ├─────────────────┤
//   │ ...              │
//   └─────────────────┘  ← bottom of goroutine stack
//
// ALL parameters are copied onto the new frame (pass by value).
// The return value space is also on the stack.

func level3(x int) int {
	// This is frame 3 — deepest call
	return x * x
}

func level2(x int) int {
	// This is frame 2
	result := level3(x + 1)
	return result * 2
}

func level1(x int) int {
	// This is frame 1
	return level2(x + 1)
}

func stackFrames() {
	fmt.Println("\n=== 2. Stack Frames ===")

	// Calling level1 creates 3 nested frames: level1 → level2 → level3
	result := level1(1)
	fmt.Printf("level1(1) = %d\n", result)
	// Walk: level1(1) → level2(2) → level3(3) → 9, back: *2=18
}

// ─── 3. STACK TRACES ──────────────────────────────────────────────────────────
//
// Go's runtime can produce a stack trace at any time.
// This is how panic output works — it shows the call stack.

func printStackTrace() {
	fmt.Println("\n=== 3. Stack Traces ===")

	// debug.PrintStack prints the current goroutine's stack trace
	// (same as what you see in a panic)
	debug.PrintStack()
}

// runtime.Stack fills a buffer with the stack trace
func getStackTrace() string {
	buf := make([]byte, 4096)
	n := runtime.Stack(buf, false)  // false = only current goroutine
	return string(buf[:n])
}

func nestedForTrace() {
	trace := getStackTrace()
	fmt.Println("\n--- Stack trace from nested call ---")
	fmt.Println(trace)
}

// ─── 4. PROGRAMMATIC CALLER INSPECTION ────────────────────────────────────────
//
// runtime.Callers returns program counters for the call stack.
// runtime.CallersFrames converts them to human-readable frames.
// This is how logging libraries add "called from" information.

func callerInfo(skip int) string {
	_, file, line, ok := runtime.Caller(skip + 1)
	if !ok {
		return "unknown"
	}
	// Trim to just filename for brevity
	short := file
	for i := len(file) - 1; i > 0; i-- {
		if file[i] == '/' {
			short = file[i+1:]
			break
		}
	}
	return fmt.Sprintf("%s:%d", short, line)
}

func fullCallStack() []string {
	pcs := make([]uintptr, 32)
	n := runtime.Callers(2, pcs)  // skip runtime.Callers and fullCallStack
	frames := runtime.CallersFrames(pcs[:n])

	var stack []string
	for {
		frame, more := frames.Next()
		stack = append(stack, fmt.Sprintf("%s (%s:%d)", frame.Function, frame.File, frame.Line))
		if !more { break }
	}
	return stack
}

func callerC() []string { return fullCallStack() }
func callerB() []string { return callerC() }
func callerA() []string { return callerB() }

func programmaticCallers() {
	fmt.Println("\n=== 4. Programmatic Caller Inspection ===")
	fmt.Printf("Called from: %s\n", callerInfo(0))

	stack := callerA()
	fmt.Println("Full call stack:")
	for i, frame := range stack {
		fmt.Printf("  [%d] %s\n", i, frame)
	}
}

// ─── 5. GOROUTINE STACK GROWTH IN PRACTICE ────────────────────────────────────
//
// Watch goroutine stacks grow to handle deep recursion.
// In Go, deep recursion does NOT immediately exhaust the stack —
// the stack grows automatically. But there IS a limit (default 1 GB).

func deepRecurse(n int) int {
	if n == 0 { return 0 }
	return 1 + deepRecurse(n-1)
}

func stackGrowth() {
	fmt.Println("\n=== 5. Stack Growth ===")

	// This would require ~100,000 stack frames — Go handles it gracefully
	depth := 100_000
	result := deepRecurse(depth)
	fmt.Printf("Recursed %d levels deep: result=%d\n", depth, result)

	// Print stack memory stats before and after
	var m1, m2 runtime.MemStats
	runtime.ReadMemStats(&m1)
	deepRecurse(10_000)
	runtime.ReadMemStats(&m2)
	fmt.Printf("Stack growth for 10k frames: was=%d now=%d bytes\n", m1.StackInuse, m2.StackInuse)
}

// ─── 6. GOROUTINE vs OS THREAD STACKS ─────────────────────────────────────────
//
// OS thread stacks:
//   - Fixed size (typically 1–8 MB on Linux)
//   - Set at thread creation, cannot grow
//   - Overflow → stack overflow crash
//   - Heavy: OS context switch required to switch between threads
//
// Go goroutine stacks:
//   - Start small (2–8 KB)
//   - Grow dynamically via stack copying
//   - No overflow (within the 1 GB default max)
//   - Lightweight: Go's scheduler switches goroutines in user space
//
// This is WHY Go can run millions of goroutines but only thousands of OS threads.
//
// GOMAXPROCS = number of OS threads running Go code simultaneously
// (defaults to number of CPUs)

func goroutineVsThread() {
	fmt.Println("\n=== 6. Goroutine vs OS Thread Stacks ===")

	fmt.Printf("GOMAXPROCS: %d (OS threads running Go code)\n", runtime.GOMAXPROCS(0))
	fmt.Printf("NumCPU: %d\n", runtime.NumCPU())
	fmt.Printf("NumGoroutine: %d\n", runtime.NumGoroutine())

	// Spawn many goroutines cheaply
	done := make(chan struct{})
	const N = 10_000
	for i := 0; i < N; i++ {
		go func() {
			<-done  // park here
		}()
	}
	fmt.Printf("After spawning %d goroutines: %d active\n", N, runtime.NumGoroutine())

	// Clean up
	close(done)
	runtime.Gosched()  // yield to let goroutines finish
}

func main() {
	goroutineStack()
	stackFrames()
	nestedForTrace()
	programmaticCallers()
	stackGrowth()
	goroutineVsThread()
}

/*
THOUGHT QUESTIONS:

1. Why do Go goroutine stacks start small (2 KB) and grow dynamically,
   rather than starting large like OS threads?

2. What is "stack copying"? What problem does it solve vs segmented stacks?
   (Research: Go 1.3 moved from segmented stacks to stack copying — why?)

3. Why can't you store a pointer to a Go stack variable in a C struct
   that outlives the goroutine's function call?

4. How does `runtime.Caller(0)` work? What does the skip parameter do?

5. What is GOMAXPROCS? What happens if you set it to 1 vs N (number of CPUs)?

EXERCISES:

1. Write a function `withCallerInfo(fn func())` that wraps a function call
   and prints "called from <file>:<line>" before executing fn.

2. Write a recursive function that deliberately causes a stack overflow
   (no base case). Run it in a goroutine with recover, catch the panic,
   and print the recovered value.

3. Spawn 1,000,000 goroutines that each just sleep for 1 second, then exit.
   Measure how long it takes to spawn them all and the peak memory usage.

4. Write a logger that includes the calling file and line number in every
   log message, using runtime.Caller.
*/
