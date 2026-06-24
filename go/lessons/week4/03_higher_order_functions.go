/*
WEEK 4 — DAY 3: Higher-Order Functions
========================================
Topic: Functions that take or return functions — map, filter, reduce, and Go idioms.

Key ideas:
  - Higher-order functions (HOF) treat functions as data
  - Go doesn't have built-in map/filter/reduce (pre-generics they required any)
  - With generics (Go 1.18+), we can write truly type-safe HOFs
  - HOFs enable composition, reuse, and declarative code
  - sort.Slice, http.HandlerFunc, sync.Once are standard library HOFs
*/

package main

import (
	"fmt"
	"sort"
	"strings"
)

// ─── 1. MAP, FILTER, REDUCE — GENERIC VERSIONS ────────────────────────────────
//
// Before generics, you'd use []interface{} or code generation.
// With generics (Go 1.18+), these are clean and type-safe.

// Map applies a transformation function to each element
func Map[T, U any](s []T, fn func(T) U) []U {
	result := make([]U, len(s))
	for i, v := range s {
		result[i] = fn(v)
	}
	return result
}

// Filter returns elements matching a predicate
func Filter[T any](s []T, predicate func(T) bool) []T {
	var result []T
	for _, v := range s {
		if predicate(v) {
			result = append(result, v)
		}
	}
	return result
}

// Reduce accumulates a result by applying fn to each element
func Reduce[T, U any](s []T, initial U, fn func(U, T) U) U {
	acc := initial
	for _, v := range s {
		acc = fn(acc, v)
	}
	return acc
}

// ForEach applies a side-effect function to each element (no return)
func ForEach[T any](s []T, fn func(T)) {
	for _, v := range s {
		fn(v)
	}
}

// Any returns true if ANY element matches the predicate
func Any[T any](s []T, predicate func(T) bool) bool {
	for _, v := range s {
		if predicate(v) { return true }
	}
	return false
}

// All returns true if ALL elements match the predicate
func All[T any](s []T, predicate func(T) bool) bool {
	for _, v := range s {
		if !predicate(v) { return false }
	}
	return true
}

// GroupBy groups elements by a key function
func GroupBy[T any, K comparable](s []T, key func(T) K) map[K][]T {
	result := make(map[K][]T)
	for _, v := range s {
		k := key(v)
		result[k] = append(result[k], v)
	}
	return result
}

func mapFilterReduce() {
	fmt.Println("=== 1. Map, Filter, Reduce (Generic) ===")

	nums := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}

	evens := Filter(nums, func(n int) bool { return n%2 == 0 })
	fmt.Println("evens:", evens)

	squares := Map(evens, func(n int) int { return n * n })
	fmt.Println("squares of evens:", squares)

	total := Reduce(squares, 0, func(acc, n int) int { return acc + n })
	fmt.Println("sum of squares:", total)

	// Chain
	result := Reduce(
		Map(
			Filter(nums, func(n int) bool { return n%2 == 0 }),
			func(n int) int { return n * n },
		),
		0,
		func(acc, n int) int { return acc + n },
	)
	fmt.Println("chained result:", result)

	// String operations
	words := []string{"hello", "world", "go", "is", "great"}
	upper := Map(words, strings.ToUpper)
	fmt.Println("upper:", upper)

	long := Filter(words, func(s string) bool { return len(s) > 3 })
	fmt.Println("long words:", long)

	concat := Reduce(words, "", func(acc, s string) string {
		if acc == "" { return s }
		return acc + " " + s
	})
	fmt.Println("concatenated:", concat)

	// GroupBy
	type Person struct{ Name, Dept string }
	people := []Person{
		{"Alice", "Eng"}, {"Bob", "Sales"}, {"Carol", "Eng"},
		{"Dave", "Sales"}, {"Eve", "Eng"},
	}
	byDept := GroupBy(people, func(p Person) string { return p.Dept })
	for dept, members := range byDept {
		names := Map(members, func(p Person) string { return p.Name })
		fmt.Printf("  %s: %v\n", dept, names)
	}
}

// ─── 2. FUNCTION COMPOSITION ──────────────────────────────────────────────────
//
// Compose creates a new function that applies f after g (right-to-left).
// Pipe applies functions left-to-right (more intuitive for data pipelines).

func Compose[T any](fns ...func(T) T) func(T) T {
	return func(v T) T {
		for i := len(fns) - 1; i >= 0; i-- {
			v = fns[i](v)
		}
		return v
	}
}

func Pipe[T any](fns ...func(T) T) func(T) T {
	return func(v T) T {
		for _, fn := range fns {
			v = fn(v)
		}
		return v
	}
}

func composition() {
	fmt.Println("\n=== 2. Function Composition ===")

	trim := strings.TrimSpace
	lower := strings.ToLower
	normalize := func(s string) string {
		return strings.ReplaceAll(s, " ", "_")
	}

	// Pipe: trim → lower → normalize (left to right)
	processName := Pipe(trim, lower, normalize)

	inputs := []string{"  Hello World  ", "ALICE", "  Go Programming  "}
	for _, s := range inputs {
		fmt.Printf("  %q → %q\n", s, processName(s))
	}

	// Numeric pipeline
	addOne := func(n int) int { return n + 1 }
	double := func(n int) int { return n * 2 }
	square := func(n int) int { return n * n }

	transform := Pipe(addOne, double, square)
	// (5+1)*2 = 12, 12^2 = 144
	fmt.Printf("  Pipe(+1, *2, ^2)(5) = %d\n", transform(5))
}

// ─── 3. PARTIAL APPLICATION AND CURRYING ──────────────────────────────────────
//
// Partial application: fix some arguments of a function, return a new function.
// Currying: transform a function of N args into N nested single-arg functions.
// (Pure currying is uncommon in Go — partial application is more idiomatic.)

func partial2[A, B, C any](fn func(A, B) C, a A) func(B) C {
	return func(b B) C {
		return fn(a, b)
	}
}

func partialApplication() {
	fmt.Println("\n=== 3. Partial Application ===")

	multiply := func(a, b int) int { return a * b }
	double := partial2(multiply, 2)
	triple := partial2(multiply, 3)

	nums := []int{1, 2, 3, 4, 5}
	fmt.Println("doubled:", Map(nums, double))
	fmt.Println("tripled:", Map(nums, triple))

	// Practical: partial application with string prefix
	addPrefix := func(prefix, s string) string { return prefix + s }
	addHTTPS := partial2(addPrefix, "https://")
	addWWW := partial2(addPrefix, "www.")

	urls := []string{"example.com", "go.dev", "golang.org"}
	fmt.Println("with https:", Map(urls, addHTTPS))
	fmt.Println("with www:", Map(urls, addWWW))
}

// ─── 4. STANDARD LIBRARY HOFs ─────────────────────────────────────────────────
//
// Go's standard library uses higher-order functions extensively.
// sort.Slice, sort.SliceStable: pass a less function
// strings.Map: transform each rune
// strings.FieldsFunc: split by a predicate
// filepath.Walk: callback for each file
// http.HandleFunc: register handler functions

func stdlibHOFs() {
	fmt.Println("\n=== 4. Standard Library HOFs ===")

	// sort.Slice with custom comparator
	type Item struct {
		Name  string
		Price float64
		Stock int
	}

	items := []Item{
		{"apple", 0.5, 100},
		{"banana", 0.3, 50},
		{"cherry", 1.2, 30},
		{"date", 2.5, 10},
	}

	// Sort by price (ascending)
	sort.Slice(items, func(i, j int) bool {
		return items[i].Price < items[j].Price
	})
	fmt.Println("By price:")
	for _, item := range items {
		fmt.Printf("  %s: $%.2f\n", item.Name, item.Price)
	}

	// Sort by stock (descending)
	sort.Slice(items, func(i, j int) bool {
		return items[i].Stock > items[j].Stock
	})
	fmt.Println("By stock (desc):")
	for _, item := range items {
		fmt.Printf("  %s: %d\n", item.Name, item.Stock)
	}

	// strings.Map — transform each rune
	rot13 := strings.Map(func(r rune) rune {
		switch {
		case r >= 'A' && r <= 'Z':
			return 'A' + (r-'A'+13)%26
		case r >= 'a' && r <= 'z':
			return 'a' + (r-'a'+13)%26
		}
		return r
	}, "Hello, World!")
	fmt.Println("ROT13:", rot13)

	// strings.FieldsFunc — split by predicate
	sentence := "  hello,world;go is great!  "
	words := strings.FieldsFunc(sentence, func(r rune) bool {
		return r == ' ' || r == ',' || r == ';' || r == '!'
	})
	fmt.Println("Split:", words)
}

// ─── 5. TABLE-DRIVEN FUNCTION DISPATCH ───────────────────────────────────────
//
// Storing functions in maps or slices is idiomatic Go for dispatching
// without large switch statements.

type Handler func(args []string) error

func tableDispatch() {
	fmt.Println("\n=== 5. Table-Driven Dispatch ===")

	commands := map[string]Handler{
		"hello": func(args []string) error {
			name := "world"
			if len(args) > 0 { name = args[0] }
			fmt.Printf("  Hello, %s!\n", name)
			return nil
		},
		"upper": func(args []string) error {
			for _, a := range args {
				fmt.Println(" ", strings.ToUpper(a))
			}
			return nil
		},
		"count": func(args []string) error {
			fmt.Printf("  %d arguments\n", len(args))
			return nil
		},
	}

	inputs := []struct {
		cmd  string
		args []string
	}{
		{"hello", []string{"Gopher"}},
		{"upper", []string{"go", "is", "great"}},
		{"count", []string{"a", "b", "c", "d"}},
		{"unknown", nil},
	}

	for _, input := range inputs {
		handler, ok := commands[input.cmd]
		if !ok {
			fmt.Printf("  unknown command: %q\n", input.cmd)
			continue
		}
		handler(input.args)
	}
}

// ─── 6. LAZY EVALUATION WITH CLOSURES ─────────────────────────────────────────
//
// Go evaluates eagerly by default. But you can defer computation using closures
// (also called "thunks").

type Lazy[T any] func() T

func LazyValue[T any](fn func() T) Lazy[T] {
	var computed bool
	var value T
	return func() T {
		if !computed {
			value = fn()
			computed = true
		}
		return value
	}
}

func lazyEvaluation() {
	fmt.Println("\n=== 6. Lazy Evaluation ===")

	// Expensive computation only runs once, and only if accessed
	expensiveValue := LazyValue(func() int {
		fmt.Println("  [computing expensive value]")
		// simulate expensive work
		return 42
	})

	fmt.Println("before access")
	fmt.Println("value:", expensiveValue())  // computes here
	fmt.Println("value:", expensiveValue())  // cached — no recomputation
	fmt.Println("value:", expensiveValue())  // cached

	// Lazy sequence generator
	naturals := func() func() int {
		n := 0
		return func() int {
			n++
			return n
		}
	}()

	// Take first 5
	for i := 0; i < 5; i++ {
		fmt.Printf("  natural: %d\n", naturals())
	}
}

func main() {
	mapFilterReduce()
	composition()
	partialApplication()
	stdlibHOFs()
	tableDispatch()
	lazyEvaluation()
}

/*
THOUGHT QUESTIONS:

1. Before Go generics, how would you implement a type-safe map() function?
   What are the trade-offs of the pre-generics approaches?

2. What is the difference between function composition and a pipeline?
   When would you use each?

3. sort.Slice takes a `less func(i, j int) bool` — indices, not values.
   Why does it work this way? What does it enable?

4. What is a "thunk"? How does Go's closure mechanism enable lazy evaluation?

5. The `GroupBy` function returns `map[K][]T`. What constraint does K need?
   Why can't K be `any`?

EXERCISES:

1. Implement `FlatMap[T, U any](s []T, fn func(T) []U) []U` that maps
   each element to a slice and flattens the result.

2. Write `Partition[T any](s []T, predicate func(T) bool) ([]T, []T)` that
   splits a slice into two: elements matching and not matching the predicate.

3. Implement a simple query builder using higher-order functions:
   `query := NewQuery(data).Where(predicate).Select(transform).Limit(10).Exec()`

4. Write a function `Retry[T any](n int, fn func() (T, error)) (T, error)` using
   generics that retries fn up to n times, returning the first success or the last error.
*/
