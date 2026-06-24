/*
WEEK 2 — DAY 2: Composite Types — Arrays, Slices, Maps
========================================================
Topic: How Go's composite types work internally.

Key ideas:
  - Arrays are values with fixed size known at compile time
  - Slices are dynamic views — a header (ptr, len, cap) over an array
  - append may or may not allocate — understanding capacity is critical
  - Maps are hash tables — iteration order is intentionally random
*/

package main

import (
	"fmt"
	"sort"
)

// ─── 1. ARRAYS: FIXED-SIZE, VALUE SEMANTICS ───────────────────────────────────
//
// Arrays in Go are RARELY used directly (slices are preferred).
// They are the foundation that slices are built on.
//
// [N]T and [M]T are DIFFERENT types — you can't mix them.
// Arrays are passed by VALUE — they are fully copied.
// Use arrays when: you need a fixed-size buffer, or as map keys.

func arrays() {
	fmt.Println("=== 1. Arrays ===")

	// Declaration forms
	var a1 [3]int              // zero-initialized
	a2 := [3]int{10, 20, 30}  // explicit values
	a3 := [...]int{1, 2, 3, 4, 5}  // compiler counts elements

	fmt.Printf("a1 (zero): %v\n", a1)
	fmt.Printf("a2: %v, len=%d\n", a2, len(a2))
	fmt.Printf("a3: %v, len=%d\n", a3, len(a3))

	// [3]int and [5]int are different types
	// a1 = a3  // COMPILE ERROR: cannot use [5]int as [3]int

	// Arrays are values — copying is a full copy
	b := a2       // full copy
	b[0] = 999
	fmt.Printf("a2 after b[0]=999: %v (unchanged)\n", a2)
	fmt.Printf("b: %v\n", b)

	// Multidimensional arrays
	matrix := [2][3]int{{1, 2, 3}, {4, 5, 6}}
	fmt.Printf("matrix[1][2] = %d\n", matrix[1][2])

	// Arrays as map keys (slices cannot be map keys)
	type Point [2]int
	visited := map[Point]bool{}
	visited[Point{3, 4}] = true
	fmt.Printf("visited[{3,4}]: %v\n", visited[Point{3, 4}])
}

// ─── 2. SLICES: INTERNAL STRUCTURE ────────────────────────────────────────────
//
// A slice is a struct with three fields:
//   type slice struct {
//       ptr unsafe.Pointer  // pointer to first element in underlying array
//       len int             // number of elements accessible via the slice
//       cap int             // total elements in underlying array from ptr
//   }
//
// This header is 24 bytes on 64-bit systems.
//
// When you write:
//   s := []int{1, 2, 3}
//
// Go allocates an array [3]int on the heap (or stack, compiler decides),
// then creates a slice header pointing to it.
//
// When you write:
//   t := s
//
// Only the header is copied. Both s and t point to the SAME underlying array.

func sliceInternals() {
	fmt.Println("\n=== 2. Slice Internals ===")

	// Create a slice
	s := []int{10, 20, 30, 40, 50}
	fmt.Printf("s:   len=%d cap=%d %v\n", len(s), cap(s), s)

	// Sub-slice: shares the same backing array
	t := s[1:4]  // elements at index 1, 2, 3
	fmt.Printf("t:   len=%d cap=%d %v\n", len(t), cap(t), t)
	// cap(t) = cap(s) - 1 = 4 because t starts at index 1 of s's array

	// Modifying t modifies s too — same backing array
	t[0] = 999
	fmt.Printf("After t[0]=999:\n")
	fmt.Printf("  s: %v\n", s) // s[1] is now 999
	fmt.Printf("  t: %v\n", t)

	// make([]T, len, cap) — create a slice with specific length and capacity
	u := make([]int, 3, 10)  // len=3, cap=10
	fmt.Printf("\nu:   len=%d cap=%d %v\n", len(u), cap(u), u)
	// u has 3 elements (all 0) but underlying array has space for 10

	// Nil slice vs empty slice
	var nilSlice []int      // nil — no backing array
	emptySlice := []int{}   // non-nil but empty — has a backing array
	made := make([]int, 0)  // also non-nil but empty

	fmt.Printf("\nnil slice:   %v, nil=%v, len=%d\n", nilSlice, nilSlice == nil, len(nilSlice))
	fmt.Printf("empty slice: %v, nil=%v, len=%d\n", emptySlice, emptySlice == nil, len(emptySlice))
	fmt.Printf("made slice:  %v, nil=%v, len=%d\n", made, made == nil, len(made))
	// Both nil and empty slices have len=0 and can be ranged over safely
}

// ─── 3. APPEND: THE GROWTH ALGORITHM ─────────────────────────────────────────
//
// append(s, elem) adds elements to a slice:
//
//   If len(s) < cap(s): adds element in-place, increments len, no allocation
//   If len(s) == cap(s): allocates a NEW larger array, copies data, returns new slice
//
// The growth factor:
//   - If cap < 256: double the capacity
//   - If cap >= 256: grow by 25% + constant
//   (This changed in Go 1.18 — previously always doubled up to some threshold)
//
// CRITICAL: append MAY return a different slice than you passed in.
//   Always use: s = append(s, elem)
//   Never: append(s, elem) ignoring the return value

func appendGrowth() {
	fmt.Println("\n=== 3. append and Growth ===")

	var s []int
	fmt.Printf("initial:  len=%d cap=%d ptr=%p\n", len(s), cap(s), s)

	for i := 0; i < 20; i++ {
		oldCap := cap(s)
		s = append(s, i)
		if cap(s) != oldCap {
			fmt.Printf("grew at len=%2d: oldCap=%2d → newCap=%2d\n",
				len(s), oldCap, cap(s))
		}
	}
	// Observe: capacity doubles (1, 2, 4, 8, 16, 32...)

	// Appending slices together
	a := []int{1, 2, 3}
	b := []int{4, 5, 6}
	c := append(a, b...)  // ... unpacks b into individual elements
	fmt.Printf("append(a, b...): %v\n", c)

	// The sharing hazard:
	// If you sub-slice and then append within capacity, you corrupt the original
	original := make([]int, 3, 6)  // len=3, cap=6
	original[0], original[1], original[2] = 1, 2, 3

	sub := original[1:3]  // shares backing array, cap = 4
	fmt.Printf("\noriginal: len=%d cap=%d %v\n", len(original), cap(original), original)
	fmt.Printf("sub:      len=%d cap=%d %v\n", len(sub), cap(sub), sub)

	// append into sub — fits in cap, no new allocation — OVERWRITES original[3]
	sub = append(sub, 99)
	fmt.Printf("After append(sub, 99):\n")
	fmt.Printf("  original: %v (original[3] is now 99!)\n", original[:4])
	fmt.Printf("  sub: %v\n", sub)

	// Fix: use full slice expression to limit capacity
	sub2 := original[1:3:3]  // cap = 3-1 = 2 (third index limits capacity)
	sub2 = append(sub2, 88)  // forces new allocation — safe
	fmt.Printf("sub2 with limited cap: %v (original unchanged)\n", sub2)
}

// ─── 4. MAPS: HASH TABLES ─────────────────────────────────────────────────────
//
// Maps in Go are hash tables. A map variable is a pointer to a runtime.hmap struct.
// This means:
//   - Maps are always reference types — assignment shares the map
//   - map[K]V requires K to be comparable (==)
//   - Iteration order is randomized intentionally (prevents relying on order)
//   - Concurrently writing to a map is a data race — use sync.Map or a mutex
//   - Reading from a nil map is safe (returns zero value)
//   - Writing to a nil map panics

func maps() {
	fmt.Println("\n=== 4. Maps ===")

	// Declaration forms
	var m1 map[string]int          // nil map — can read, can't write
	m2 := map[string]int{}         // empty map — can read and write
	m3 := make(map[string]int)     // also empty, equivalent to m2
	m4 := map[string]int{          // literal with initial values
		"alice": 95,
		"bob":   87,
		"carol": 92,
	}

	fmt.Printf("m1 (nil): %v, nil=%v\n", m1, m1 == nil)
	fmt.Printf("m2 (empty): %v, nil=%v\n", m2, m2 == nil)
	_ = m3

	// Reading from a nil map is safe
	val := m1["missing"]
	fmt.Printf("read from nil map: %d (zero value)\n", val)

	// Two-value read: comma ok idiom
	score, exists := m4["alice"]
	fmt.Printf("alice: score=%d, exists=%v\n", score, exists)

	score, exists = m4["nobody"]
	fmt.Printf("nobody: score=%d, exists=%v\n", score, exists)

	// Maps are reference types — copying the variable copies the pointer
	m5 := m4
	m5["alice"] = 0  // modifies the SAME map as m4
	fmt.Printf("m4[alice] after m5 mutation: %d\n", m4["alice"])

	// Delete
	delete(m4, "carol")
	fmt.Printf("after delete: %v\n", m4)

	// Iteration order is random — MUST sort for deterministic output
	scores := map[string]int{"a": 3, "b": 1, "c": 2}
	keys := make([]string, 0, len(scores))
	for k := range scores {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		fmt.Printf("  %s: %d\n", k, scores[k])
	}
}

// ─── 5. SLICE PATTERNS ────────────────────────────────────────────────────────

func slicePatterns() {
	fmt.Println("\n=== 5. Slice Patterns ===")

	nums := []int{3, 1, 4, 1, 5, 9, 2, 6, 5, 3}

	// Filter in-place (no extra allocation)
	n := 0
	for _, v := range nums {
		if v > 3 {
			nums[n] = v
			n++
		}
	}
	nums = nums[:n]
	fmt.Printf("filtered (>3): %v\n", nums)

	// Stack using a slice
	stack := []int{}
	push := func(v int) { stack = append(stack, v) }
	pop := func() int {
		v := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		return v
	}
	push(1); push(2); push(3)
	fmt.Printf("pop: %d %d %d\n", pop(), pop(), pop())

	// Queue using a slice (inefficient — use container/list for large queues)
	queue := []string{}
	enqueue := func(s string) { queue = append(queue, s) }
	dequeue := func() string {
		s := queue[0]
		queue = queue[1:]
		return s
	}
	enqueue("first"); enqueue("second"); enqueue("third")
	fmt.Printf("dequeue: %s %s %s\n", dequeue(), dequeue(), dequeue())

	// Copy to avoid sharing
	original := []int{1, 2, 3, 4, 5}
	copied := make([]int, len(original))
	copy(copied, original)
	copied[0] = 999
	fmt.Printf("original after copy mutation: %v (unchanged)\n", original)
}

// ─── 6. MAP PATTERNS ──────────────────────────────────────────────────────────

func mapPatterns() {
	fmt.Println("\n=== 6. Map Patterns ===")

	// Counting frequencies
	words := []string{"go", "is", "great", "go", "is", "fast", "go"}
	freq := make(map[string]int)
	for _, w := range words {
		freq[w]++  // first access returns 0 (zero value), then we increment
	}
	fmt.Printf("word frequencies: %v\n", freq)

	// Grouping
	type Person struct{ Name, City string }
	people := []Person{
		{"Alice", "NYC"}, {"Bob", "LA"}, {"Carol", "NYC"}, {"Dave", "LA"},
	}
	byCity := make(map[string][]Person)
	for _, p := range people {
		byCity[p.City] = append(byCity[p.City], p)
	}
	for city, ps := range byCity {
		fmt.Printf("%s: ", city)
		for _, p := range ps { fmt.Printf("%s ", p.Name) }
		fmt.Println()
	}

	// Set using map[T]struct{}
	// struct{} is a zero-size type — uses no memory for values
	set := map[string]struct{}{}
	set["a"] = struct{}{}
	set["b"] = struct{}{}
	_, inSet := set["a"]
	fmt.Printf("'a' in set: %v\n", inSet)
	fmt.Printf("'c' in set: %v\n", func() bool { _, ok := set["c"]; return ok }())
}

func main() {
	arrays()
	sliceInternals()
	appendGrowth()
	maps()
	slicePatterns()
	mapPatterns()
}

/*
THOUGHT QUESTIONS:

1. A slice header has ptr, len, and cap. What is the difference between len
   and cap? What happens when you append past cap?

2. After `t := s[1:4]`, what is cap(t) if cap(s) = 8 and t starts at index 1?

3. Why does Go randomize map iteration order? What problem does this prevent?

4. Why is `map[T]struct{}` preferred over `map[T]bool` for sets?

5. What is the "sharing hazard" with sub-slices and append? How does the
   full slice expression `s[low:high:max]` prevent it?

EXERCISES:

1. Implement a function `unique(s []int) []int` that returns a new slice
   with duplicate elements removed, preserving order.

2. Write a function `mergeSlices(a, b []int) []int` that merges two sorted
   slices into one sorted slice without using sort.Sort.

3. Implement a basic LRU cache using a map and a doubly-linked list (or
   a simpler version using a map + slice).

4. Write a function that takes a `map[string]interface{}` (like a parsed JSON
   object) and prints it in a nested, indented format.
*/
