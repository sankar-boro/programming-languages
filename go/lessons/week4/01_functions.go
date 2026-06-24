/*
WEEK 4 — DAY 1: Functions — Signatures, Multiple Returns, Named Returns
========================================================================
Topic: Go functions in depth — signatures, the call stack, variadic functions,
       multiple return values, and how Go passes arguments.

Key ideas:
  - Go passes ALL arguments by value (copying)
  - Functions can return multiple values — a core Go feature
  - Named return values can be used for documentation and bare returns
  - Variadic functions accept any number of trailing arguments
  - Functions are first-class values in Go
*/

package main

import (
	"fmt"
	"math"
	"sort"
)

// ─── 1. FUNCTION SIGNATURES ────────────────────────────────────────────────────
//
// func funcName(param1 Type1, param2 Type2) ReturnType { }
//
// Parameters of the same type can be grouped:
//   func add(x, y int) int   (instead of func add(x int, y int) int)
//
// Return types can also be grouped:
//   func swap(x, y int) (int, int)

func add(x, y int) int {
	return x + y
}

func divide(a, b float64) (float64, error) {
	if b == 0 {
		return 0, fmt.Errorf("cannot divide by zero")
	}
	return a / b, nil
}

// Multiple parameters of the same type, grouped
func minMax(a, b, c int) (min, max int) {
	if a < b { min = a } else { min = b }
	if c < min { min = c }
	if a > b { max = a } else { max = b }
	if c > max { max = c }
	return  // bare return — returns named values
}

func functionSignatures() {
	fmt.Println("=== 1. Function Signatures ===")

	fmt.Println("add(3, 4):", add(3, 4))

	q, err := divide(10, 3)
	if err == nil {
		fmt.Printf("10/3 = %.4f\n", q)
	}

	min, max := minMax(5, 2, 8)
	fmt.Printf("minMax(5,2,8): min=%d max=%d\n", min, max)
}

// ─── 2. MULTIPLE RETURN VALUES IN DEPTH ───────────────────────────────────────
//
// Multiple return values are a key Go feature — used pervasively for errors.
//
// Under the hood: multiple return values are returned via the stack.
// The caller allocates space on its stack frame for all return values.
// The callee writes to that space directly.
//
// This is more efficient than allocating a tuple/struct on the heap,
// and it makes error handling a first-class concern without exceptions.

func splitAndValidate(s string) (before, after string, found bool) {
	// Named returns: auto-initialized to zero values
	for i, c := range s {
		if c == ':' {
			before = s[:i]
			after = s[i+1:]
			found = true
			return  // bare return — returns named values
		}
	}
	before = s
	return  // found is still false (zero value)
}

func parseCoord(s string) (x, y float64, err error) {
	_, err = fmt.Sscanf(s, "%f,%f", &x, &y)
	if err != nil {
		err = fmt.Errorf("invalid coordinate %q: %w", s, err)
	}
	return
}

func multipleReturns() {
	fmt.Println("\n=== 2. Multiple Return Values ===")

	b, a, ok := splitAndValidate("key:value")
	fmt.Printf("split 'key:value': before=%q after=%q found=%v\n", b, a, ok)

	b, a, ok = splitAndValidate("nocodon")
	fmt.Printf("split 'nocodon': before=%q after=%q found=%v\n", b, a, ok)

	x, y, err := parseCoord("3.5,7.2")
	if err == nil {
		fmt.Printf("coord: x=%.1f y=%.1f\n", x, y)
	}

	_, _, err = parseCoord("invalid")
	fmt.Printf("invalid coord error: %v\n", err)
}

// ─── 3. THE CALL STACK ─────────────────────────────────────────────────────────
//
// When a function is called:
//   1. A NEW stack frame is pushed onto the goroutine's stack
//   2. The frame contains: parameters (copied from caller), local variables,
//      return value space, and the return address
//   3. When the function returns: frame is popped, return values are copied
//
// Goroutine stack growth:
//   - Goroutine stacks start small (2KB in Go 1.4+)
//   - Go uses "stack copying" — when a goroutine's stack runs out,
//     a new, larger stack is allocated, and all stack data is copied
//   - You generally don't need to worry about stack depth (unlike C)
//   - Deep recursion will eventually exhaust the system limit (not the goroutine limit)
//
// Key insight: ALL function arguments are copied onto the new stack frame.
// This is why mutating function parameters doesn't affect callers.

func showCopying(s []int) {
	// s is a COPY of the slice header — but the backing array is shared
	fmt.Printf("  inside showCopying: s=%v, &s=%p\n", s, &s)
	s[0] = 999  // modifies the SHARED backing array
	s = append(s, 100)  // may create a NEW backing array — but header copy is lost
}

func callStack() {
	fmt.Println("\n=== 3. The Call Stack ===")

	original := []int{1, 2, 3}
	fmt.Printf("before call: original=%v &original=%p\n", original, &original)
	showCopying(original)
	fmt.Printf("after call:  original=%v\n", original)
	// original[0] is 999 (shared backing array)
	// but the appended 100 is NOT visible (header was a copy)
}

// ─── 4. VARIADIC FUNCTIONS ────────────────────────────────────────────────────
//
// A variadic function accepts any number of trailing arguments of type T.
// Inside the function, they're received as a []T slice.
// Call with: func(a, b, c) or func(slice...) to unpack a slice.

func sum(nums ...int) int {
	total := 0
	for _, n := range nums {
		total += n
	}
	return total
}

func printf(format string, args ...any) {
	// Internally, fmt.Fprintf uses this pattern
	fmt.Printf("[custom] "+format, args...)
}

func max(first int, rest ...int) int {
	m := first
	for _, v := range rest {
		if v > m { m = v }
	}
	return m
}

func variadics() {
	fmt.Println("\n=== 4. Variadic Functions ===")

	fmt.Println("sum():          ", sum())
	fmt.Println("sum(1):         ", sum(1))
	fmt.Println("sum(1,2,3,4,5): ", sum(1, 2, 3, 4, 5))

	// Unpack a slice
	nums := []int{10, 20, 30, 40}
	fmt.Println("sum(nums...):   ", sum(nums...))

	printf("hello %s, you are %d years old\n", "world", 42)

	fmt.Println("max(1):          ", max(1))
	fmt.Println("max(3, 1, 4, 1, 5):", max(3, 1, 4, 1, 5))
}

// ─── 5. FUNCTIONS AS FIRST-CLASS VALUES ──────────────────────────────────────
//
// Functions are values in Go. You can:
//   - Assign them to variables
//   - Pass them as arguments
//   - Return them from functions
//   - Store them in data structures
//   - Call them via the variable

type Predicate func(int) bool
type Transform func(int) int

func filter(nums []int, predicate Predicate) []int {
	var result []int
	for _, n := range nums {
		if predicate(n) {
			result = append(result, n)
		}
	}
	return result
}

func mapSlice(nums []int, transform Transform) []int {
	result := make([]int, len(nums))
	for i, n := range nums {
		result[i] = transform(n)
	}
	return result
}

func reduce(nums []int, initial int, fn func(acc, val int) int) int {
	acc := initial
	for _, n := range nums {
		acc = fn(acc, n)
	}
	return acc
}

func firstClassFunctions() {
	fmt.Println("\n=== 5. Functions as First-Class Values ===")

	nums := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}

	// Pass functions as arguments
	evens := filter(nums, func(n int) bool { return n%2 == 0 })
	fmt.Println("evens:", evens)

	squared := mapSlice(nums, func(n int) int { return n * n })
	fmt.Println("squared:", squared)

	total := reduce(nums, 0, func(acc, val int) int { return acc + val })
	fmt.Println("sum:", total)

	// Store functions in a map (dispatch table)
	ops := map[string]func(float64, float64) float64{
		"add": func(a, b float64) float64 { return a + b },
		"sub": func(a, b float64) float64 { return a - b },
		"mul": func(a, b float64) float64 { return a * b },
		"div": func(a, b float64) float64 { return a / b },
	}

	for name, op := range ops {
		fmt.Printf("  10 %s 3 = %.2f\n", name, op(10, 3))
	}

	// sort.Slice uses a function argument
	words := []string{"banana", "apple", "cherry", "date"}
	sort.Slice(words, func(i, j int) bool {
		return len(words[i]) < len(words[j])  // sort by length
	})
	fmt.Println("sorted by length:", words)
}

// ─── 6. FUNCTION VALUES AND THE CALL OVERHEAD ─────────────────────────────────
//
// Calling a function through a variable (indirect call) vs calling a known
// function (direct call):
//
//   Direct call:   CALL instruction to a known address — very fast
//   Indirect call: CALL through a pointer — requires one extra memory load
//                  (CPU may not be able to inline or predict the branch)
//
// This is why hot paths avoid function values. But for most code,
// the overhead is negligible — profile before optimizing.
//
// A function value in Go is:
//   - A pointer to the function's machine code
//   - If it's a closure, also a pointer to a "closure context" (the captured vars)

type Middleware func(http_handler func(string) string) func(string) string

func logging(next func(string) string) func(string) string {
	return func(req string) string {
		fmt.Printf("  [LOG] handling %q\n", req)
		resp := next(req)
		fmt.Printf("  [LOG] response: %q\n", resp)
		return resp
	}
}

func auth(next func(string) string) func(string) string {
	return func(req string) string {
		if req == "/admin" {
			return "403 Forbidden"
		}
		return next(req)
	}
}

func functionValues() {
	fmt.Println("\n=== 6. Function Values — Middleware Pattern ===")

	// Base handler
	handler := func(req string) string {
		return "200 OK: " + req
	}

	// Wrap with middleware
	wrapped := logging(auth(handler))

	wrapped("/home")
	wrapped("/admin")
}

// ─── 7. RECURSIVE FUNCTIONS AND STACK DEPTH ───────────────────────────────────

func fibonacci(n int) int {
	if n <= 1 { return n }
	return fibonacci(n-1) + fibonacci(n-2)
}

// Tail-recursive version (Go does NOT optimize tail calls — same stack usage)
func fibTail(n, a, b int) int {
	if n == 0 { return a }
	return fibTail(n-1, b, a+b)
}

// Iterative — most efficient
func fibIter(n int) int {
	if n <= 1 { return n }
	a, b := 0, 1
	for i := 2; i <= n; i++ {
		a, b = b, a+b
	}
	return b
}

func recursion() {
	fmt.Println("\n=== 7. Recursion (with caveats) ===")
	fmt.Println("fib(10) naive:    ", fibonacci(10))
	fmt.Println("fib(10) tail:     ", fibTail(10, 0, 1))
	fmt.Println("fib(10) iterative:", fibIter(10))

	// Note: Go does NOT perform tail call optimization.
	// Deep recursion WILL overflow the goroutine's stack.
	// The goroutine stack will GROW to accommodate, but eventually
	// the system will run out of address space (around 1 billion frames on 64-bit).
	fmt.Println("fib(40):", fibonacci(40))  // slow, 2^40 calls
	fmt.Println("fibIter(40):", fibIter(40)) // fast
}

// Demo of a math utility to make imports used
func circleArea(r float64) float64 { return math.Pi * r * r }

func main() {
	functionSignatures()
	multipleReturns()
	callStack()
	variadics()
	firstClassFunctions()
	functionValues()
	recursion()
	_ = circleArea
}

/*
THOUGHT QUESTIONS:

1. "Go passes all arguments by value." What does this mean for a slice argument?
   What is shared and what is copied?

2. What are named return values? When are they helpful, and when can they
   make code harder to understand?

3. A variadic parameter `nums ...int` is a `[]int` inside the function.
   What happens if you do `nums = append(nums, 99)` inside the function?
   Does it affect the caller?

4. When you store a function in a variable and call it, what is the
   performance cost vs a direct function call?

5. Go does not optimize tail calls. What does this mean for deeply recursive
   Go code? What's the alternative?

EXERCISES:

1. Write a `pipeline(fns ...func(int) int) func(int) int` function that
   composes a series of int transformations into one function.

2. Implement a `memoize(fn func(int) int) func(int) int` that caches results
   of an expensive function. Use it to speed up the naive fibonacci.

3. Write a `curry2(fn func(int, int) int) func(int) func(int) int` that
   curries a 2-argument function into two 1-argument functions.

4. Implement a retry mechanism: `retry(n int, fn func() error) error` that
   calls fn up to n times with exponential backoff (using time.Sleep).
*/
