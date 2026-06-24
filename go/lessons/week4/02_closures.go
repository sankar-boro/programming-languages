/*
WEEK 4 — DAY 2: Closures
==========================
Topic: How closures work in Go — captured variables, heap escape, and patterns.

Key ideas:
  - A closure is a function + the variables it captures from its enclosing scope
  - Captured variables are shared — modifications are visible across all closures
  - Captured variables escape to the heap (they must outlive the stack frame)
  - Closures are used for: factory functions, callbacks, iterators, middleware
  - The "loop variable capture" gotcha is a classic Go bug
*/

package main

import (
	"fmt"
	"sync"
)

// ─── 1. WHAT IS A CLOSURE? ────────────────────────────────────────────────────
//
// A closure = a function + a reference to its captured environment.
//
// When a function references a variable from an outer scope, it "closes over"
// that variable. The variable is not copied — it is SHARED by reference.
//
// Under the hood, the Go compiler creates a "closure struct" on the heap
// that holds pointers to all captured variables. The function itself is
// also stored as part of this struct (or pointed to by it).
//
// This is why captured variables escape to the heap — they must outlive
// the function that originally declared them.

func basicClosure() {
	fmt.Println("=== 1. Basic Closure ===")

	x := 10  // x lives on the stack of basicClosure

	// inner captures x — x must escape to the heap
	inner := func() {
		fmt.Printf("  x = %d (captured from outer scope)\n", x)
		x++  // modifies the SAME x
	}

	inner()     // x becomes 11
	inner()     // x becomes 12
	fmt.Printf("x after two calls: %d\n", x)
	// x is shared between basicClosure and inner
}

// ─── 2. CLOSURES AS FACTORY FUNCTIONS ────────────────────────────────────────
//
// The most common closure use: generate specialized functions that share state.

func makeCounter(start int) func() int {
	count := start  // each call to makeCounter creates a SEPARATE count
	return func() int {
		count++
		return count
	}
}

func makeAdder(x int) func(int) int {
	return func(y int) int {
		return x + y  // captures x
	}
}

func makeMultiplier(factor int) func(int) int {
	return func(n int) int {
		return n * factor  // captures factor
	}
}

func factoryFunctions() {
	fmt.Println("\n=== 2. Factory Functions ===")

	// Two independent counters — each has its OWN count variable
	counter1 := makeCounter(0)
	counter2 := makeCounter(100)

	fmt.Println("counter1:", counter1(), counter1(), counter1())  // 1 2 3
	fmt.Println("counter2:", counter2(), counter2())              // 101 102
	fmt.Println("counter1:", counter1())                          // 4 (independent)

	// Adders
	add5 := makeAdder(5)
	add10 := makeAdder(10)
	fmt.Printf("add5(3) = %d, add10(3) = %d\n", add5(3), add10(3))

	// Pipeline of multipliers
	double := makeMultiplier(2)
	triple := makeMultiplier(3)
	fmt.Printf("double(7) = %d, triple(7) = %d\n", double(7), triple(7))
}

// ─── 3. MUTABLE SHARED STATE ─────────────────────────────────────────────────
//
// Multiple closures can close over the SAME variable.
// This allows them to share state — like private mutable state in OOP.

func makeStack() (push func(int), pop func() (int, bool), peek func() (int, bool)) {
	stack := []int{}  // shared by all three closures

	push = func(v int) {
		stack = append(stack, v)
	}

	pop = func() (int, bool) {
		if len(stack) == 0 {
			return 0, false
		}
		v := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		return v, true
	}

	peek = func() (int, bool) {
		if len(stack) == 0 {
			return 0, false
		}
		return stack[len(stack)-1], true
	}

	return
}

func sharedState() {
	fmt.Println("\n=== 3. Shared State via Closures ===")

	push, pop, peek := makeStack()

	push(10)
	push(20)
	push(30)

	if v, ok := peek(); ok {
		fmt.Println("peek:", v)  // 30
	}

	for {
		v, ok := pop()
		if !ok { break }
		fmt.Printf("  popped: %d\n", v)
	}
}

// ─── 4. THE LOOP VARIABLE CAPTURE GOTCHA ──────────────────────────────────────
//
// This is one of the most common Go bugs. In Go 1.21 and earlier:
//
// In a for loop, the loop variable is a SINGLE variable that changes each iteration.
// If you capture it in a closure (goroutine, defer), ALL closures see the FINAL value.
//
// In Go 1.22+, each iteration of a range loop creates a NEW variable per iteration.
// The gotcha only applies to Go 1.21 and earlier (or classic for loops).

func loopVariableGotcha() {
	fmt.Println("\n=== 4. Loop Variable Capture Gotcha ===")

	// In Go 1.22+, range loop variables are per-iteration (fixed)
	// But classic for loops still have this issue even in 1.22+

	// WRONG WAY (classic for loop — ALL closures capture the SAME i)
	funcsWrong := make([]func(), 3)
	for i := 0; i < 3; i++ {
		funcsWrong[i] = func() {
			fmt.Printf("  wrong: %d\n", i)  // captures the variable i, not its value
		}
	}
	fmt.Println("Wrong (all print 3):")
	for _, f := range funcsWrong {
		f()
	}

	// FIX 1: Create a new variable inside the loop (shadows the loop variable)
	funcsFixed1 := make([]func(), 3)
	for i := 0; i < 3; i++ {
		i := i  // NEW variable i per iteration — shadows the outer i
		funcsFixed1[i] = func() {
			fmt.Printf("  fixed1: %d\n", i)
		}
	}
	fmt.Println("Fixed 1 (correct):")
	for _, f := range funcsFixed1 {
		f()
	}

	// FIX 2: Pass as function argument (arguments are copied, not shared)
	funcsFixed2 := make([]func(), 3)
	for i := 0; i < 3; i++ {
		func(j int) {  // j is a fresh copy
			funcsFixed2[j] = func() {
				fmt.Printf("  fixed2: %d\n", j)
			}
		}(i)
	}
	fmt.Println("Fixed 2 (correct):")
	for _, f := range funcsFixed2 {
		f()
	}

	// Range loop in Go 1.22+ is fixed — each iteration gets its own variable
	funcsRange := make([]func(), 3)
	vals := []int{10, 20, 30}
	for i, v := range vals {
		_ = i
		funcsRange[i] = func() {
			fmt.Printf("  range (1.22+): %d\n", v)  // each iteration has its own v
		}
	}
	fmt.Println("Range loop (Go 1.22+, each v is independent):")
	for _, f := range funcsRange {
		f()
	}
}

// ─── 5. CLOSURES AND GOROUTINES ───────────────────────────────────────────────
//
// When goroutines capture variables, you need to be careful about data races.
// The loop gotcha is amplified with goroutines.

func closuresAndGoroutines() {
	fmt.Println("\n=== 5. Closures and Goroutines ===")

	var wg sync.WaitGroup

	// CORRECT: pass i as argument to avoid sharing
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func(n int) {
			defer wg.Done()
			fmt.Printf("  goroutine %d\n", n)
		}(i)
	}
	wg.Wait()

	// Closures sharing a channel
	results := make(chan int, 5)
	for i := 1; i <= 5; i++ {
		go func(n int) {
			results <- n * n  // each goroutine captures its own n
		}(i)
	}
	for i := 0; i < 5; i++ {
		fmt.Printf("  result: %d\n", <-results)
	}
}

// ─── 6. PRACTICAL CLOSURE PATTERNS ────────────────────────────────────────────

// Middleware pattern: wrap a handler with pre/post logic
func withLogging(name string, fn func() error) func() error {
	return func() error {
		fmt.Printf("[start] %s\n", name)
		err := fn()
		if err != nil {
			fmt.Printf("[error] %s: %v\n", name, err)
		} else {
			fmt.Printf("[done]  %s\n", name)
		}
		return err
	}
}

// Once: execute a function at most once
func once(fn func()) func() {
	var done bool
	return func() {
		if !done {
			done = true
			fn()
		}
	}
}

// Memoize: cache results of an expensive function
func memoize(fn func(int) int) func(int) int {
	cache := map[int]int{}
	return func(n int) int {
		if v, ok := cache[n]; ok {
			fmt.Printf("  [cache hit] %d\n", n)
			return v
		}
		result := fn(n)
		cache[n] = result
		return result
	}
}

func closurePatterns() {
	fmt.Println("\n=== 6. Closure Patterns ===")

	// Middleware
	doWork := func() error {
		fmt.Println("  doing work...")
		return nil
	}
	logged := withLogging("doWork", doWork)
	logged()

	// Once
	initOnce := once(func() {
		fmt.Println("  initialized!")
	})
	initOnce()
	initOnce()  // doesn't run again
	initOnce()  // doesn't run again

	// Memoize
	expensiveFib := memoize(func(n int) int {
		// naive fib — very slow without memoization
		if n <= 1 { return n }
		return n  // simplified for demo
	})
	fmt.Println("fib(5):", expensiveFib(5))
	fmt.Println("fib(5):", expensiveFib(5))  // cache hit
	fmt.Println("fib(3):", expensiveFib(3))
}

func main() {
	basicClosure()
	factoryFunctions()
	sharedState()
	loopVariableGotcha()
	closuresAndGoroutines()
	closurePatterns()
}

/*
THOUGHT QUESTIONS:

1. When a closure captures a variable, where does that variable live?
   Why must it be on the heap instead of the stack?

2. Two closures close over the same variable x. If closure A modifies x,
   does closure B see the change? Explain.

3. What is the loop variable capture gotcha? Does Go 1.22+ fix it completely?
   What about classic 3-clause for loops?

4. In the `makeStack` example, three closures share a slice. What would
   happen if two goroutines called push and pop simultaneously?

5. Closures can implement stateful objects (like the stack example).
   What are the trade-offs vs using a struct with methods?

EXERCISES:

1. Implement `makeRateLimiter(n int, period time.Duration) func() bool`
   that returns true if the caller can proceed (within the rate limit) or
   false if rate-limited. Uses a closure over shared state.

2. Implement `debounce(fn func(), delay time.Duration) func()` using
   closures and time. The returned function only calls fn if it hasn't
   been called in the last `delay` duration.

3. Write `compose(fns ...func(int) int) func(int) int` that creates a
   function applying each fn in sequence. Use a closure over fns.

4. Write a concurrent-safe counter using closures (not sync.Mutex directly,
   but using a channel to serialize access). The returned functions should
   be safe to call from multiple goroutines.
*/
