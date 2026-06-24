/*
WEEK 1 — DAY 2: Syntax, Values, and the Type System
======================================================
Topic: Go's syntax philosophy, value model, and how the type system works.

Key ideas:
  - Go has a deliberately small, regular syntax
  - Everything has a type; there is no implicit type coercion
  - Go is strongly and statically typed
  - Types are checked at compile time, not runtime
  - Zero values: every type has a sensible default
*/

package main

import (
	"fmt"
	"math"
	"reflect"
	"strconv"
)

// ─── 1. GO'S SYNTAX PHILOSOPHY ────────────────────────────────────────────────
//
// Go's syntax was designed to be:
//   - Small: ~25 keywords (vs C++'s ~90, Python's ~35)
//   - Regular: few special cases, few exceptions
//   - Readable: code should be readable top-to-bottom without surprises
//
// Notable decisions:
//   - Semicolons are auto-inserted by the lexer (you don't write them)
//   - Opening brace MUST be on the same line (no style debates)
//   - Unused imports and variables are COMPILE ERRORS (enforced hygiene)
//   - Type comes AFTER the variable name: var x int (not int x)

func syntaxPhilosophy() {
	fmt.Println("=== 1. Syntax Philosophy ===")

	// Type after name (like "x is an int")
	var count int = 10
	var name string = "Go"
	var ratio float64 = 3.14

	fmt.Printf("count: %v (%T)\n", count, count)
	fmt.Printf("name:  %v (%T)\n", name, name)
	fmt.Printf("ratio: %v (%T)\n", ratio, ratio)

	_ = count // suppress "unused variable" error for demo
	_ = name
	_ = ratio
}

// ─── 2. THE TYPE HIERARCHY ────────────────────────────────────────────────────
//
// Go's types fall into categories:
//
//   Basic types:
//     bool
//     string
//     int, int8, int16, int32 (rune), int64
//     uint, uint8 (byte), uint16, uint32, uint64, uintptr
//     float32, float64
//     complex64, complex128
//
//   Composite types (built-in):
//     array    [N]T
//     slice    []T
//     map      map[K]V
//     struct   struct{...}
//
//   Reference types:
//     pointer  *T
//     function func(...)
//     channel  chan T
//     slice    []T  (slice header is a reference)
//     map      map[K]V (map header is a reference)
//
//   Interface type:
//     interface{...}
//
// IMPORTANT: In Go, integers are NOT automatically promoted to floats.
//   You must explicitly convert: float64(myInt)

func typeHierarchy() {
	fmt.Println("\n=== 2. The Type Hierarchy ===")

	// Signed integers
	var i8 int8 = 127       // max int8
	var i32 int32 = 2147483647
	var i64 int64 = 9223372036854775807

	// Unsigned
	var u8 uint8 = 255      // same as byte
	var b byte = 'A'        // byte is an alias for uint8

	// Floating point
	var f32 float32 = 3.14
	var f64 float64 = math.Pi  // 64-bit is the default for float literals

	// Rune — a Unicode code point (alias for int32)
	var r rune = '猫'  // a Chinese character, code point U+732B
	fmt.Printf("rune '猫' = %d (int32)\n", r)

	// String — sequence of bytes (UTF-8 encoded), NOT characters
	var s string = "Hello, 世界"
	fmt.Printf("string: %q, len=%d (bytes, not chars)\n", s, len(s))

	// No implicit conversion between numeric types
	// var bad float64 = i32  // COMPILE ERROR
	var good float64 = float64(i32)  // explicit conversion required

	fmt.Printf("int8: %d, int32: %d, int64: %d\n", i8, i32, i64)
	fmt.Printf("uint8/byte: %d, 'A'=%d\n", u8, b)
	fmt.Printf("float32: %f, float64: %.15f\n", f32, f64)
	fmt.Printf("float64(i32): %f\n", good)
}

// ─── 3. ZERO VALUES — A CORE GO PRINCIPLE ─────────────────────────────────────
//
// Every type in Go has a ZERO VALUE — the value a variable holds when
// declared without an explicit initializer.
//
// There are NO uninitialized variables in Go. This eliminates a whole class
// of bugs (reading uninitialized memory) at the language level.
//
// Zero values:
//   bool:    false
//   int:     0
//   float64: 0.0
//   string:  ""  (empty string)
//   pointer: nil
//   slice:   nil (but a nil slice is valid — length 0, can append to it)
//   map:     nil (nil map is readable but writing panics)
//   channel: nil
//   func:    nil
//   struct:  all fields at their zero values

func zeroValues() {
	fmt.Println("\n=== 3. Zero Values ===")

	var b bool
	var i int
	var f float64
	var s string
	var p *int
	var sl []int
	var m map[string]int

	fmt.Printf("bool:    %v\n", b)      // false
	fmt.Printf("int:     %v\n", i)      // 0
	fmt.Printf("float64: %v\n", f)      // 0
	fmt.Printf("string:  %q\n", s)      // ""
	fmt.Printf("*int:    %v\n", p)      // <nil>
	fmt.Printf("[]int:   %v (nil=%v)\n", sl, sl == nil)   // [] true
	fmt.Printf("map:     %v (nil=%v)\n", m, m == nil)     // map[] true

	// A nil slice is valid — you can range over it, append to it
	for _, v := range sl {
		fmt.Println(v) // never executes
	}
	sl = append(sl, 1, 2, 3) // works fine on nil slice
	fmt.Printf("after append: %v\n", sl)

	// A nil map panics on WRITE but not read
	// m["key"] = 1  // PANIC: assignment to entry in nil map
	val := m["key"] // safe — returns zero value (0)
	fmt.Printf("nil map read: %d\n", val)
}

// ─── 4. TYPE IDENTITY — NAMED VS UNNAMED TYPES ────────────────────────────────
//
// Go distinguishes between named types and unnamed (composite) types.
//
// Named type: defined with `type` keyword — creates a distinct new type.
// Even if the underlying type is the same, named types are not interchangeable.
//
// This is called NOMINAL typing in this narrow sense, but Go uses STRUCTURAL
// typing for interfaces (covered later).

type Celsius float64
type Fahrenheit float64
type Meters float64
type Feet float64

func (c Celsius) ToFahrenheit() Fahrenheit {
	return Fahrenheit(c*9/5 + 32)
}

func typeIdentity() {
	fmt.Println("\n=== 4. Type Identity ===")

	var temp Celsius = 100
	var distance Meters = 5.0

	fmt.Printf("100°C = %.1f°F\n", temp.ToFahrenheit())

	// These would be compile errors — type safety prevents unit confusion:
	// var f Fahrenheit = temp        // ERROR: cannot use Celsius as Fahrenheit
	// var ft Feet = distance         // ERROR: cannot use Meters as Feet
	_ = distance

	// Explicit conversion is allowed (same underlying type):
	var raw float64 = float64(temp)
	fmt.Printf("raw float64: %v\n", raw)
}

// ─── 5. REFLECT — INSPECTING TYPES AT RUNTIME ─────────────────────────────────
//
// The `reflect` package lets you inspect types at runtime. This is how
// fmt.Printf("%v", x) works for any type — it uses reflect internally.
//
// In normal code you rarely need reflect. It exists for:
//   - Writing generic utilities (before Go 1.18 generics)
//   - Serialization (JSON, XML, DB drivers)
//   - Frameworks and ORMs
//
// reflect has two key types:
//   reflect.Type  — describes the type
//   reflect.Value — holds the value + type, allows reading/writing

func reflectInspection() {
	fmt.Println("\n=== 5. Reflect: Inspecting Types at Runtime ===")

	values := []interface{}{
		42,
		3.14,
		"hello",
		true,
		[]int{1, 2, 3},
		map[string]int{"a": 1},
	}

	for _, v := range values {
		t := reflect.TypeOf(v)
		val := reflect.ValueOf(v)
		fmt.Printf("Type: %-20s Kind: %-10s Value: %v\n",
			t.String(), t.Kind().String(), val)
	}

	// reflect.Kind vs reflect.Type:
	//   Type  = the full name ("main.Celsius", "float64", "[]int")
	//   Kind  = the underlying category (float64, slice, struct, ptr, ...)
	var c Celsius = 37.5
	fmt.Printf("\nCelsius Type: %v, Kind: %v\n", reflect.TypeOf(c), reflect.TypeOf(c).Kind())
	// Type = main.Celsius, Kind = float64
}

// ─── 6. TYPE CONVERSION VS TYPE ASSERTION ─────────────────────────────────────
//
// Type CONVERSION: converting between compatible concrete types.
//   Must be explicit. Checked at compile time.
//   float64(myInt), string(myRune), int(myFloat64)
//
// Type ASSERTION: extracting the concrete type from an interface.
//   Must be explicit. Checked at runtime. Can panic if wrong.
//   x.(string)    — panics if x is not a string
//   s, ok := x.(string)  — safe form: ok=false if not a string

func conversionVsAssertion() {
	fmt.Println("\n=== 6. Conversion vs Assertion ===")

	// Conversion — compile-time, between concrete types
	var i int = 65
	var f float64 = float64(i)  // widen
	var b byte = byte(i)        // narrow (same value, both 65)
	var r rune = rune(i)        // int32, Unicode point 65 = 'A'
	var s string = string(r)    // rune → single-char string
	var n int = 42
	var ns string = strconv.Itoa(n)  // int → string representation (NOT string(n))

	fmt.Printf("int→float64: %f\n", f)
	fmt.Printf("int→byte: %c (%d)\n", b, b)
	fmt.Printf("rune→string: %q\n", s)
	fmt.Printf("int→string via Itoa: %q\n", ns)

	// Type assertion — runtime, from interface to concrete
	var iface interface{} = "hello"

	// Safe assertion
	if str, ok := iface.(string); ok {
		fmt.Printf("Assertion succeeded: %q\n", str)
	}

	// Unsafe assertion — panics if wrong type
	str := iface.(string)
	fmt.Printf("Unsafe assertion: %q (safe here because we know the type)\n", str)

	// Wrong assertion would panic:
	// num := iface.(int)  // PANIC: interface conversion: want int, have string
}

func main() {
	syntaxPhilosophy()
	typeHierarchy()
	zeroValues()
	typeIdentity()
	reflectInspection()
	conversionVsAssertion()
}

/*
THOUGHT QUESTIONS:

1. Why does Go require explicit type conversion between numeric types?
   What class of bugs does this prevent?

2. What is a zero value? Why is Go's "every variable is initialized" rule
   better than C's "uninitialized variables have garbage values"?

3. What is the difference between a type's Kind and its Type in reflect?

4. Why would you create a named type like `type Celsius float64`
   instead of just using `float64` everywhere?

5. What is the difference between type conversion and type assertion?
   When can each fail?

EXERCISES:

1. Create named types for Kilograms and Pounds. Write a conversion method.
   Try assigning a Kilograms value to a Pounds variable — what happens?

2. Declare variables of every basic type without initialization.
   Print their zero values. Verify they match the table above.

3. Write a function that takes an interface{} and uses reflect to print
   the Kind, Type, and Value of whatever is passed in.

4. Why does `string(65)` give you "A" and not "65"? How do you get "65"?
   (Hint: strconv)
*/
