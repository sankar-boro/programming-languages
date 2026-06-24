/*
WEEK 1 — DAY 3: Variables, Constants, and the Memory Model
============================================================
Topic: How Go stores values in memory, the stack vs heap distinction,
       and how Go's variable declaration forms work.

Key ideas:
  - Variables in Go are memory locations with a type
  - Unlike Python, Go variables are NOT labels to objects
  - The compiler decides stack vs heap via escape analysis
  - Constants are compile-time values, not memory locations
  - Pointers give you direct access to memory addresses
*/

package main

import (
	"fmt"
	"unsafe"
)

// ─── 1. VARIABLE DECLARATION FORMS ────────────────────────────────────────────
//
// Go has several ways to declare variables. They are equivalent but used
// in different contexts:
//
//   var x int             → zero value, explicit type, used at package level
//   var x int = 10        → explicit type + initializer
//   var x = 10            → inferred type (int), used at package level
//   x := 10               → short declaration, ONLY inside functions
//   var x, y int          → multiple variables of same type
//   var x, y = 1, "hello" → multiple variables with inferred types
//
// `:=` is the most common inside functions. `var` is required at package level.

var packageLevel int = 42  // must use `var` at package level
// packageLevel := 42      // COMPILE ERROR at package level

func variableDeclarations() {
	fmt.Println("=== 1. Variable Declaration Forms ===")

	// Long form
	var a int = 10
	var b = 20         // type inferred as int
	c := 30            // short declaration

	// Multiple
	var x, y int = 1, 2
	p, q := "hello", true

	// Grouped declaration block
	var (
		width  int     = 800
		height int     = 600
		title  string  = "Go Window"
	)

	fmt.Println(a, b, c)
	fmt.Println(x, y)
	fmt.Println(p, q)
	fmt.Printf("window: %dx%d '%s'\n", width, height, title)

	// Blank identifier: _ discards a value
	val, _ := returnTwo()
	fmt.Println("first return only:", val)
}

func returnTwo() (int, string) {
	return 42, "discarded"
}

// ─── 2. THE MEMORY MODEL: VALUES vs REFERENCES ────────────────────────────────
//
// This is CRITICAL to understand in Go. Unlike Python where variables
// are always references to heap objects, Go variables are:
//
//   DIRECT VALUES for basic types (int, float64, bool, struct, array):
//     - The variable IS the data
//     - Assignment COPIES the data
//     - No sharing unless you use a pointer
//
//   HEADER/REFERENCE for slice, map, channel, function:
//     - The variable holds a small "header" struct
//     - The header points to heap-allocated backing data
//     - Assignment copies the header (both now point to the same backing data)
//
// Python model:
//   x = [1,2,3]     # x is a reference to a list object on the heap
//   y = x           # y is ANOTHER reference to the SAME object
//
// Go value model:
//   x := [3]int{1,2,3}  // x IS the array (on the stack)
//   y := x              // y is a COPY of the array
//   y[0] = 99           // only y changes; x is unchanged
//
// Go slice model (reference-like):
//   x := []int{1,2,3}   // x is a slice HEADER pointing to heap data
//   y := x              // y copies the header — both point to SAME backing array
//   y[0] = 99           // BOTH x and y see the change

func memoryModel() {
	fmt.Println("\n=== 2. The Memory Model ===")

	// Arrays are values — copy semantics
	arrX := [3]int{1, 2, 3}
	arrY := arrX         // full copy
	arrY[0] = 99
	fmt.Printf("arrX: %v (unchanged)\n", arrX)
	fmt.Printf("arrY: %v (modified)\n", arrY)

	// Slices share backing array — reference semantics
	sliceX := []int{1, 2, 3}
	sliceY := sliceX     // copies the header — shared backing array
	sliceY[0] = 99
	fmt.Printf("sliceX: %v (also changed!)\n", sliceX)
	fmt.Printf("sliceY: %v\n", sliceY)

	// Structs are values — copy semantics
	type Point struct{ X, Y int }
	p1 := Point{1, 2}
	p2 := p1         // full copy
	p2.X = 100
	fmt.Printf("p1: %v (unchanged)\n", p1)
	fmt.Printf("p2: %v (modified)\n", p2)

	// Strings are immutable values (header = ptr + len)
	s1 := "hello"
	s2 := s1       // copies the string header, NOT the bytes
	_ = s2         // but since strings are immutable, it doesn't matter
}

// ─── 3. POINTERS — EXPLICIT REFERENCE SEMANTICS ───────────────────────────────
//
// Go has pointers. Unlike C, Go pointers:
//   - Cannot be arithmetic-operated (no ptr++)
//   - Are nil-safe (the language prevents dangling pointers)
//   - Are garbage-collected (you never manually free pointed-to memory)
//
// & = "address of" — gives you a pointer to a variable
// * = "dereference" — follows the pointer to get/set the value
//
// When to use pointers:
//   1. To mutate a value inside a function (otherwise Go copies the argument)
//   2. To avoid copying large structs
//   3. To express "optional" or "nullable" values

func pointers() {
	fmt.Println("\n=== 3. Pointers ===")

	x := 42
	p := &x           // p is *int — pointer to x's memory location
	fmt.Printf("x = %d\n", x)
	fmt.Printf("&x = %p (memory address)\n", &x)
	fmt.Printf("p = %p (same address)\n", p)
	fmt.Printf("*p = %d (value at that address)\n", *p)

	// Mutate x through the pointer
	*p = 100
	fmt.Printf("After *p = 100: x = %d\n", x) // x is now 100

	// Pointer to struct — Go auto-dereferences
	type User struct {
		Name string
		Age  int
	}
	u := User{"Alice", 30}
	up := &u
	up.Name = "Bob"  // same as (*up).Name = "Bob" — Go handles this
	fmt.Printf("After pointer mutation: %+v\n", u)

	// new() allocates a zeroed value on the heap and returns a pointer
	np := new(int)  // *int, points to zeroed int on heap
	*np = 55
	fmt.Printf("new(int): %d at %p\n", *np, np)

	// Nil pointer
	var nilPtr *int  // zero value for pointer is nil
	fmt.Printf("nil pointer: %v\n", nilPtr)
	// *nilPtr = 1  // PANIC: nil pointer dereference
}

// ─── 4. ESCAPE ANALYSIS: STACK vs HEAP ────────────────────────────────────────
//
// Go's compiler decides where to allocate memory using ESCAPE ANALYSIS.
// You don't control this directly — the compiler does.
//
// Stack allocation:
//   - Fast: just move the stack pointer
//   - Automatic: freed when function returns
//   - Limited: goroutine stacks start at 2KB (grow dynamically)
//
// Heap allocation:
//   - Slower: requires GC to eventually free
//   - Used when a value must outlive the function that created it
//
// A variable "escapes to the heap" when:
//   1. You take its address and return the pointer (lifetime > stack frame)
//   2. It's stored in an interface (interface values are heap-allocated)
//   3. It's too large for the stack
//   4. Its size is unknown at compile time
//
// To see escape analysis:
//   go build -gcflags="-m" main.go

func stackAlloc() int {
	x := 42    // x lives on the stack — freed when stackAlloc returns
	return x   // value is COPIED to the caller
}

func heapAlloc() *int {
	x := 42   // x escapes to heap — the pointer must outlive this function
	return &x // returning a pointer to a local → x is heap-allocated
}

func escapeAnalysis() {
	fmt.Println("\n=== 4. Escape Analysis: Stack vs Heap ===")

	v1 := stackAlloc()  // v1 gets a copy of 42
	v2 := heapAlloc()   // v2 is a pointer to a heap-allocated int

	fmt.Printf("stackAlloc: %d (copied from stack)\n", v1)
	fmt.Printf("heapAlloc: %d (pointer to heap, addr=%p)\n", *v2, v2)

	// Interface boxing always causes heap allocation
	var iface interface{} = 42  // the 42 escapes to heap to satisfy interface{}
	fmt.Printf("interface boxed value: %v\n", iface)
}

// ─── 5. CONSTANTS — COMPILE-TIME VALUES ───────────────────────────────────────
//
// Constants in Go are:
//   - Evaluated at compile time
//   - NOT stored in memory (no address, can't take &const)
//   - Untyped unless you declare a type explicitly
//   - Much more flexible than C/Java constants
//
// Untyped constants have "ideal" types — they take on the type they're
// used in, as long as the value fits.

const Pi = 3.14159265358979  // untyped floating-point constant
const MaxRetries = 5          // untyped integer constant

const TypedPi float32 = 3.14159  // typed constant — fixed to float32

// iota — auto-incrementing integer constant generator
type Weekday int

const (
	Sunday Weekday = iota  // 0
	Monday                 // 1
	Tuesday                // 2
	Wednesday              // 3
	Thursday               // 4
	Friday                 // 5
	Saturday               // 6
)

type ByteSize float64

const (
	_           = iota // discard first value (0)
	KB ByteSize = 1 << (10 * iota) // 1 << 10 = 1024
	MB                              // 1 << 20
	GB                              // 1 << 30
	TB                              // 1 << 40
)

func constants() {
	fmt.Println("\n=== 5. Constants ===")

	// Untyped Pi works as float32 OR float64 — adapts to context
	var f32 float32 = Pi     // OK — untyped constant adapts
	var f64 float64 = Pi     // also OK
	fmt.Printf("Pi as float32: %.7f\n", f32)
	fmt.Printf("Pi as float64: %.15f\n", f64)

	fmt.Printf("MaxRetries: %d\n", MaxRetries)
	fmt.Printf("Days: Sun=%d Mon=%d Sat=%d\n", Sunday, Monday, Saturday)
	fmt.Printf("KB=%.0f MB=%.0f GB=%.0f TB=%.0f\n", float64(KB), float64(MB), float64(GB), float64(TB))

	// Can't take address of a constant
	// &Pi  // COMPILE ERROR: cannot take the address of Pi
}

// ─── 6. SIZEOF — HOW MUCH MEMORY DOES A TYPE USE? ────────────────────────────
//
// unsafe.Sizeof returns the number of bytes a type occupies.
// This helps understand memory layout.

func memorySize() {
	fmt.Println("\n=== 6. Memory Sizes ===")

	fmt.Printf("bool:       %d bytes\n", unsafe.Sizeof(bool(false)))
	fmt.Printf("int8:       %d bytes\n", unsafe.Sizeof(int8(0)))
	fmt.Printf("int16:      %d bytes\n", unsafe.Sizeof(int16(0)))
	fmt.Printf("int32:      %d bytes\n", unsafe.Sizeof(int32(0)))
	fmt.Printf("int64:      %d bytes\n", unsafe.Sizeof(int64(0)))
	fmt.Printf("int:        %d bytes (platform-dependent)\n", unsafe.Sizeof(int(0)))
	fmt.Printf("float32:    %d bytes\n", unsafe.Sizeof(float32(0)))
	fmt.Printf("float64:    %d bytes\n", unsafe.Sizeof(float64(0)))
	fmt.Printf("string:     %d bytes (header: ptr+len)\n", unsafe.Sizeof(string("")))
	fmt.Printf("[]int:      %d bytes (header: ptr+len+cap)\n", unsafe.Sizeof([]int{}))
	fmt.Printf("map[s]int:  %d bytes (pointer to hash table)\n", unsafe.Sizeof(map[string]int{}))

	type Point struct{ X, Y int }
	type Point3D struct{ X, Y, Z float32 }
	fmt.Printf("Point{x,y int}: %d bytes\n", unsafe.Sizeof(Point{}))
	fmt.Printf("Point3D{x,y,z float32}: %d bytes\n", unsafe.Sizeof(Point3D{}))
}

func main() {
	variableDeclarations()
	memoryModel()
	pointers()
	escapeAnalysis()
	constants()
	memorySize()
}

/*
THOUGHT QUESTIONS:

1. In Python, `x = 42; y = x; y = 100` — does x change? What about in Go?

2. What is escape analysis? Who performs it — you, or the compiler?
   What triggers a variable to escape to the heap?

3. Why is the zero value for a pointer `nil` — not some random address?

4. What are "untyped constants" in Go? How can the same constant be used
   as both float32 and float64 without conversion?

5. A slice header is 24 bytes (ptr + len + cap). What does this mean for
   `y := x` where x is a slice? What is shared? What is copied?

EXERCISES:

1. Write a function `increment(p *int)` that increments the value at
   the pointer by 1. Call it and verify the original variable changed.

2. Demonstrate that arrays are copied but slices share backing data.
   Create an array, copy it, modify the copy, print both arrays.
   Then repeat with a slice.

3. Create a struct with several fields of different types.
   Use unsafe.Sizeof + unsafe.Offsetof to print the memory layout.
   Then add a single byte field between two 8-byte fields and
   observe how padding affects the total size.

4. Use `go build -gcflags="-m" .` on this file. Which variables does
   the compiler report as escaping to the heap? Why?
*/
