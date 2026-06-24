/*
WEEK 6 — DAY 2: Generics (Go 1.18+)
======================================
Topic: Type parameters, constraints, and when to use generics in Go.

Key ideas:
  - Generics add type parameters to functions and types
  - Constraints limit what types can be used (using interfaces)
  - `comparable` is a built-in constraint for == / != operators
  - `any` = no constraint (like interface{} but for type params)
  - Generics enable type-safe reusable containers and algorithms
  - Prefer concrete types or interfaces when generics aren't needed
*/

package main

import (
	"cmp"
	"fmt"
	"sort"
)

// ─── 1. TYPE PARAMETERS — THE BASICS ─────────────────────────────────────────
//
// Syntax: func FuncName[T constraint](param T) T
//
// T is a type parameter — a placeholder for a concrete type.
// The caller provides the concrete type (or Go infers it).
// The constraint limits which types T can be.

// Without generics (pre-1.18): must repeat for each type or use interface{}
func maxInt(a, b int) int {
	if a > b { return a }
	return b
}

func maxFloat(a, b float64) float64 {
	if a > b { return a }
	return b
}

// With generics: one function works for all ordered types
func max[T cmp.Ordered](a, b T) T {
	if a > b { return a }
	return b
}

// cmp.Ordered is a constraint from the standard library (Go 1.21+):
// type Ordered interface {
//     ~int | ~int8 | ~int16 | ~int32 | ~int64 |
//     ~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr |
//     ~float32 | ~float64 | ~string
// }
// The ~ means "any type with this underlying type"

func typeParameterBasics() {
	fmt.Println("=== 1. Type Parameters ===")

	// Type inference — Go infers T from the arguments
	fmt.Println("max(3, 5):", max(3, 5))            // T = int
	fmt.Println("max(3.1, 2.9):", max(3.1, 2.9))   // T = float64
	fmt.Println("max(\"a\", \"b\"):", max("a", "b")) // T = string

	// Explicit type parameter
	fmt.Println("max[int](3, 5):", max[int](3, 5))

	// Without generics
	fmt.Println("maxInt(3, 5):", maxInt(3, 5))
	fmt.Println("maxFloat(3.1, 2.9):", maxFloat(3.1, 2.9))
}

// ─── 2. CONSTRAINTS — DEFINING WHAT T CAN BE ──────────────────────────────────
//
// A constraint is an interface used as a type bound.
// Three main categories:
//
//   any          = no constraint (the type can be anything)
//   comparable   = types that support == and != (built-in)
//   custom       = your own interface as a constraint

// Custom constraint using interface syntax
type Number interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 |
		~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 |
		~float32 | ~float64
}

func sum[T Number](nums []T) T {
	var total T
	for _, n := range nums {
		total += n
	}
	return total
}

// comparable constraint — for maps and equality checks
func contains[T comparable](slice []T, target T) bool {
	for _, v := range slice {
		if v == target { return true }
	}
	return false
}

func indexOf[T comparable](slice []T, target T) int {
	for i, v := range slice {
		if v == target { return i }
	}
	return -1
}

func constraints() {
	fmt.Println("\n=== 2. Constraints ===")

	ints := []int{1, 2, 3, 4, 5}
	floats := []float64{1.1, 2.2, 3.3}

	fmt.Println("sum(ints):", sum(ints))
	fmt.Println("sum(floats):", sum(floats))

	words := []string{"apple", "banana", "cherry"}
	fmt.Println("contains 'banana':", contains(words, "banana"))
	fmt.Println("contains 'grape':", contains(words, "grape"))
	fmt.Println("indexOf 'cherry':", indexOf(words, "cherry"))

	nums := []int{10, 20, 30, 40}
	fmt.Println("indexOf 30:", indexOf(nums, 30))
}

// ─── 3. GENERIC TYPES — CONTAINERS ────────────────────────────────────────────
//
// Types can also have type parameters.
// This enables type-safe containers without code generation.

// Stack — a type-safe generic stack
type Stack[T any] struct {
	items []T
}

func (s *Stack[T]) Push(v T) {
	s.items = append(s.items, v)
}

func (s *Stack[T]) Pop() (T, bool) {
	var zero T
	if len(s.items) == 0 {
		return zero, false
	}
	v := s.items[len(s.items)-1]
	s.items = s.items[:len(s.items)-1]
	return v, true
}

func (s *Stack[T]) Peek() (T, bool) {
	var zero T
	if len(s.items) == 0 {
		return zero, false
	}
	return s.items[len(s.items)-1], true
}

func (s *Stack[T]) Len() int { return len(s.items) }

// Optional — a generic wrapper for optional values
type Optional[T any] struct {
	value T
	valid bool
}

func Some[T any](v T) Optional[T]   { return Optional[T]{value: v, valid: true} }
func None[T any]() Optional[T]      { return Optional[T]{} }
func (o Optional[T]) IsPresent() bool { return o.valid }
func (o Optional[T]) Get() (T, bool)  { return o.value, o.valid }
func (o Optional[T]) OrElse(def T) T {
	if o.valid { return o.value }
	return def
}

// Result — either a value or an error
type Result[T any] struct {
	value T
	err   error
}

func Ok[T any](v T) Result[T]          { return Result[T]{value: v} }
func Err[T any](e error) Result[T]     { return Result[T]{err: e} }
func (r Result[T]) IsOk() bool         { return r.err == nil }
func (r Result[T]) Unwrap() T {
	if r.err != nil { panic(r.err) }
	return r.value
}
func (r Result[T]) UnwrapOr(def T) T {
	if r.err != nil { return def }
	return r.value
}

func genericTypes() {
	fmt.Println("\n=== 3. Generic Types ===")

	// Stack
	s := &Stack[int]{}
	s.Push(10); s.Push(20); s.Push(30)
	fmt.Println("Stack len:", s.Len())
	for s.Len() > 0 {
		v, _ := s.Pop()
		fmt.Printf("  popped: %d\n", v)
	}

	// String stack
	ss := &Stack[string]{}
	ss.Push("a"); ss.Push("b"); ss.Push("c")
	if v, ok := ss.Peek(); ok {
		fmt.Println("String stack peek:", v)
	}

	// Optional
	name := Some("Alice")
	empty := None[string]()
	fmt.Println("Optional Some:", name.OrElse("default"))
	fmt.Println("Optional None:", empty.OrElse("default"))

	// Result
	r1 := Ok(42)
	r2 := Err[int](fmt.Errorf("something failed"))
	fmt.Println("Ok result:", r1.Unwrap())
	fmt.Println("Err result:", r2.UnwrapOr(-1))
}

// ─── 4. GENERIC FUNCTIONS ON SLICES ──────────────────────────────────────────
//
// The slices package (Go 1.21+) provides generic slice utilities.
// Let's implement similar functions to understand the pattern.

func Reverse[T any](s []T) []T {
	result := make([]T, len(s))
	for i, v := range s {
		result[len(s)-1-i] = v
	}
	return result
}

func Unique[T comparable](s []T) []T {
	seen := make(map[T]bool)
	var result []T
	for _, v := range s {
		if !seen[v] {
			seen[v] = true
			result = append(result, v)
		}
	}
	return result
}

func Chunk[T any](s []T, size int) [][]T {
	var chunks [][]T
	for size < len(s) {
		s, chunks = s[size:], append(chunks, s[:size])
	}
	if len(s) > 0 {
		chunks = append(chunks, s)
	}
	return chunks
}

func Zip[T, U any](a []T, b []U) []struct{ First T; Second U } {
	n := len(a)
	if len(b) < n { n = len(b) }
	result := make([]struct{ First T; Second U }, n)
	for i := 0; i < n; i++ {
		result[i] = struct{ First T; Second U }{a[i], b[i]}
	}
	return result
}

func SortBy[T any, K cmp.Ordered](s []T, key func(T) K) []T {
	result := make([]T, len(s))
	copy(result, s)
	sort.Slice(result, func(i, j int) bool {
		return key(result[i]) < key(result[j])
	})
	return result
}

func genericSliceFunctions() {
	fmt.Println("\n=== 4. Generic Slice Functions ===")

	nums := []int{1, 2, 3, 4, 5}
	fmt.Println("Reverse:", Reverse(nums))
	fmt.Println("Unique:", Unique([]int{1, 2, 2, 3, 3, 3, 4}))
	fmt.Println("Chunk(3):", Chunk(nums, 3))

	names := []string{"Alice", "Bob", "Carol"}
	ages := []int{30, 25, 35}
	pairs := Zip(names, ages)
	for _, p := range pairs {
		fmt.Printf("  %s is %d\n", p.First, p.Second)
	}

	type Person struct{ Name string; Age int }
	people := []Person{{"Bob", 25}, {"Alice", 30}, {"Carol", 35}}
	sorted := SortBy(people, func(p Person) int { return p.Age })
	for _, p := range sorted {
		fmt.Printf("  %s: %d\n", p.Name, p.Age)
	}
}

// ─── 5. WHEN NOT TO USE GENERICS ─────────────────────────────────────────────
//
// Generics add complexity. Use them when:
//   ✓ You're writing a container (Stack, Queue, Set, Map, Optional)
//   ✓ You're writing algorithms over slices/maps that don't depend on the type
//   ✓ You have two+ implementations that are identical except for the type
//
// Don't use generics when:
//   ✗ A single concrete type works fine
//   ✗ You can use an interface (polymorphism at runtime)
//   ✗ The function needs to know the structure of the type
//   ✗ You'd need runtime type switching (just use interface{} + type switch)

// GOOD: generic — same logic, different types
func Keys[K comparable, V any](m map[K]V) []K {
	keys := make([]K, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}

// BAD idea as generic — needs to know type structure
// func serialize[T any](v T) string { ... }
// This would need reflection anyway — just use interface{}

// GOOD: interface — algorithm depends on behavior, not type
type Sortable interface {
	Len() int
	Less(i, j int) bool
	Swap(i, j int)
}

func whenNotToUseGenerics() {
	fmt.Println("\n=== 5. When NOT to Use Generics ===")

	m := map[string]int{"a": 1, "b": 2, "c": 3}
	fmt.Println("Keys:", Keys(m))  // good generic use

	fm := map[int]string{1: "one", 2: "two"}
	fmt.Println("Keys:", Keys(fm))  // same function, different types
}

func main() {
	typeParameterBasics()
	constraints()
	genericTypes()
	genericSliceFunctions()
	whenNotToUseGenerics()
}

/*
THOUGHT QUESTIONS:

1. What is a "constraint" in Go generics? How is it expressed?

2. What is the `~` operator in a constraint like `~int`? What does it
   enable that `int` alone would not?

3. `comparable` is a built-in constraint. What types satisfy it?
   Why can't `[]int` satisfy `comparable`?

4. When should you use generics vs interfaces for polymorphism?
   What is the key trade-off?

5. Go generics use monomorphization (like Rust) or type erasure (like Java)?
   What are the implications of each approach?

EXERCISES:

1. Implement a generic `Set[T comparable]` with Add, Remove, Contains,
   Union, Intersection, and Difference methods.

2. Write a generic `Cache[K comparable, V any]` with Get and Set methods,
   plus a maximum size and LRU eviction policy.

3. Implement generic versions of Map, Filter, and Reduce that work on
   any slice type. Use them to process a []Employee with salary filtering.

4. Write a generic `tree.BinarySearchTree[T cmp.Ordered]` with Insert,
   Contains, and InOrder methods.
*/
