/*
WEEK 2 — DAY 1: Types in Depth
================================
Topic: Go's type system — numeric types, strings, booleans, and type aliases.

Key ideas:
  - Go's type system is strong and static — no implicit coercions
  - Numeric types have specific sizes and overflow behavior
  - Strings are byte sequences (UTF-8), not character sequences
  - Type aliases vs type definitions are fundamentally different
  - The `any` type (interface{}) is Go's escape hatch for dynamic typing
*/

package main

import (
	"fmt"
	"math"
	"strings"
	"unicode/utf8"
)

// ─── 1. INTEGERS: SIZES, OVERFLOW, AND UNSIGNED ───────────────────────────────
//
// Go integers have explicit sizes. Overflow wraps around silently (no panic).
// This is a deliberate design decision — overflow detection is your responsibility.
//
// int and uint are platform-dependent: 32 bits on 32-bit systems, 64 on 64-bit.
// For portable code, use explicit sizes (int32, int64).
// For general use, int is idiomatic (the "default" integer type).
//
// Integer literals can be written as:
//   decimal: 1000000      → can use _ as separator: 1_000_000
//   hex:     0xFF
//   octal:   0o77 (or 077, deprecated)
//   binary:  0b1010

func integers() {
	fmt.Println("=== 1. Integers ===")

	// Sizes and ranges
	fmt.Printf("int8:   %d to %d\n", math.MinInt8, math.MaxInt8)
	fmt.Printf("int16:  %d to %d\n", math.MinInt16, math.MaxInt16)
	fmt.Printf("int32:  %d to %d\n", math.MinInt32, math.MaxInt32)
	fmt.Printf("int64:  %d to %d\n", math.MinInt64, math.MaxInt64)
	fmt.Printf("uint8:  %d to %d\n", 0, math.MaxUint8)
	fmt.Printf("uint16: %d to %d\n", 0, math.MaxUint16)

	// Overflow wraps around — no panic
	var x uint8 = 255
	x++
	fmt.Printf("uint8(255)++: %d (wraps to 0)\n", x)

	var y int8 = 127
	y++
	fmt.Printf("int8(127)++: %d (wraps to -128)\n", y)

	// Integer literals
	decimal := 1_000_000     // readability underscores
	hex := 0xFF              // 255
	octal := 0o77            // 63
	binary := 0b1010         // 10

	fmt.Printf("literals: %d %d %d %d\n", decimal, hex, octal, binary)

	// Integer division: truncates toward zero
	fmt.Printf("7/2 = %d (integer division)\n", 7/2)
	fmt.Printf("7%%2 = %d (remainder)\n", 7%2)
	fmt.Printf("-7/2 = %d (truncates toward zero, not floor)\n", -7/2)
}

// ─── 2. FLOATING POINT: IEEE 754 AND PRECISION ────────────────────────────────
//
// Go uses IEEE 754 standard for float32 and float64.
//
// float64 (double precision):
//   - 64 bits: 1 sign, 11 exponent, 52 mantissa
//   - Approximately 15-17 significant decimal digits
//
// float32 (single precision):
//   - 32 bits: 1 sign, 8 exponent, 23 mantissa
//   - Approximately 6-9 significant decimal digits
//   - Use when memory matters (e.g., large float arrays)
//
// Special values: +Inf, -Inf, NaN (Not a Number)
// NaN is not equal to anything, including itself: NaN != NaN

func floats() {
	fmt.Println("\n=== 2. Floating Point ===")

	// Precision
	var f32 float32 = 1.00000001
	var f64 float64 = 1.00000001
	fmt.Printf("float32: %.10f (precision lost)\n", f32)
	fmt.Printf("float64: %.10f (preserved)\n", f64)

	// Floating-point is NOT exact for decimal fractions
	fmt.Printf("0.1 + 0.2 = %.20f (not exactly 0.3)\n", 0.1+0.2)
	fmt.Printf("0.1 + 0.2 == 0.3: %v (false!)\n", 0.1+0.2 == 0.3)

	// How to compare floats correctly:
	const epsilon = 1e-9
	a, b := 0.1+0.2, 0.3
	fmt.Printf("Within epsilon: %v\n", math.Abs(a-b) < epsilon)

	// Special values
	posInf := math.Inf(1)
	negInf := math.Inf(-1)
	nan := math.NaN()

	fmt.Printf("+Inf: %v\n", posInf)
	fmt.Printf("-Inf: %v\n", negInf)
	fmt.Printf("NaN: %v\n", nan)
	fmt.Printf("NaN == NaN: %v (always false!)\n", nan == nan)
	fmt.Printf("math.IsNaN(nan): %v\n", math.IsNaN(nan))
	fmt.Printf("1.0 / 0.0 would panic; use math.Inf\n")
}

// ─── 3. STRINGS: BYTES, RUNES, AND UTF-8 ──────────────────────────────────────
//
// This is one of the most misunderstood areas in Go.
//
// A Go string is:
//   - An IMMUTABLE sequence of BYTES (not characters)
//   - Internally: a header with (pointer to byte array, length in bytes)
//   - Always valid UTF-8 by convention (not enforced)
//
// A rune (int32) is a Unicode code point — one "character" in the
// Unicode sense. A single rune can be 1–4 bytes in UTF-8.
//
// len(s)    → number of BYTES, not characters
// len([]rune(s)) → number of Unicode code points (characters)
//
// s[i]      → byte at index i (uint8)
// for i, ch := range s  → iterates over RUNES (code points), not bytes

func strings_() {
	fmt.Println("\n=== 3. Strings: Bytes, Runes, UTF-8 ===")

	s := "Hello, 世界"  // 7 ASCII bytes + 6 bytes for 2 Chinese chars (3 each)

	fmt.Printf("string: %q\n", s)
	fmt.Printf("len(s) = %d bytes\n", len(s))
	fmt.Printf("utf8.RuneCountInString = %d runes\n", utf8.RuneCountInString(s))

	// Indexing gives bytes, not characters
	fmt.Printf("s[0] = %d (%c) — first byte\n", s[0], s[0])
	fmt.Printf("s[7] = %d — first byte of '世' (not a valid char on its own)\n", s[7])

	// Range over string iterates by RUNE (Unicode code point)
	fmt.Println("\nRange over string (by rune):")
	for i, r := range s {
		fmt.Printf("  byte_index=%d rune=%c code_point=%d bytes=%d\n",
			i, r, r, utf8.RuneLen(r))
	}

	// Convert string to []byte for byte-level manipulation
	bytes := []byte(s)
	fmt.Printf("\n[]byte: %v\n", bytes[:7]) // first 7 bytes (ASCII part)

	// Convert string to []rune for character-level manipulation
	runes := []rune(s)
	fmt.Printf("[]rune: %v\n", runes)
	fmt.Printf("runes[7] = %c (second Chinese char)\n", runes[7])

	// String concatenation
	// + creates a new string (strings are immutable)
	greeting := "Hello" + ", " + "World"
	fmt.Println(greeting)

	// Efficient concatenation: strings.Builder
	var sb strings.Builder
	for i := 0; i < 5; i++ {
		sb.WriteString(fmt.Sprintf("item%d ", i))
	}
	fmt.Println(sb.String())

	// String comparison is lexicographic, byte-by-byte
	fmt.Printf("\"abc\" < \"abd\": %v\n", "abc" < "abd")
	fmt.Printf("\"abc\" == \"abc\": %v\n", "abc" == "abc")

	// Raw string literals — no escape sequences interpreted
	raw := `Line 1
Line 2
Tab:	here`
	fmt.Println(raw)
}

// ─── 4. BOOLEANS AND SHORT-CIRCUIT EVALUATION ─────────────────────────────────

func booleans() {
	fmt.Println("\n=== 4. Booleans ===")

	t := true
	f := false

	fmt.Printf("AND: %v, OR: %v, NOT: %v\n", t && f, t || f, !t)

	// Short-circuit evaluation — second operand may not be evaluated
	x := 0
	divideByX := func() bool {
		fmt.Println("  (divideByX evaluated)")
		return 10/x > 0  // would panic if x is 0
	}

	// Short-circuits — divideByX never called
	fmt.Println("false && divideByX():")
	result := false && divideByX()
	fmt.Println("result:", result)

	fmt.Println("true || divideByX():")
	result = true || divideByX()
	fmt.Println("result:", result)

	// Booleans are NOT integers — no if 1 or if 0
	// if 1 { }  // COMPILE ERROR: non-boolean condition
}

// ─── 5. TYPE ALIASES vs TYPE DEFINITIONS ──────────────────────────────────────
//
// Go 1.9+ has two different constructs that look similar but differ fundamentally:
//
//   type MyInt int         → TYPE DEFINITION: creates a NEW, distinct type
//                            MyInt and int are not the same type
//                            MyInt does NOT inherit int's methods (has none here)
//                            explicit conversion required: MyInt(x) or int(x)
//
//   type MyInt = int       → TYPE ALIAS: MyInt IS int, just another name
//                            MyInt and int are interchangeable
//                            all int methods are available on MyInt
//                            NO conversion needed
//
// Type aliases exist mainly for gradual code migration and stdlib augmentation.
// Type definitions exist to create domain-specific types with type safety.

type NewInt int         // new type
type AliasInt = int     // alias — just another name for int

// NewInt doesn't automatically have arithmetic operators that cross types
// but can define its own methods:
func (n NewInt) Double() NewInt { return n * 2 }

func typeAliasVsDefinition() {
	fmt.Println("\n=== 5. Type Alias vs Definition ===")

	var n NewInt = 10
	var a AliasInt = 10
	var i int = 10

	// NewInt ≠ int — explicit conversion required
	// i = n        // COMPILE ERROR
	i = int(n)     // OK with explicit conversion
	n = NewInt(i)  // OK with explicit conversion

	// AliasInt = int — interchangeable
	i = a           // OK — they ARE the same type
	a = i           // OK

	fmt.Printf("NewInt: %v, Double: %v\n", n, n.Double())
	fmt.Printf("AliasInt: %v (interchangeable with int)\n", a)
	_ = i
}

// ─── 6. THE `any` TYPE (interface{}) ──────────────────────────────────────────
//
// `any` is an alias for `interface{}` introduced in Go 1.18.
// It represents a value of any type.
//
// When you assign a value to an `any`:
//   - A small "interface box" is created on the heap
//   - It stores: (pointer to type descriptor, pointer to value)
//   - This is called "boxing" — it has a performance cost
//
// You can't do anything useful with an `any` value without a type assertion.
// Avoid `any` when possible — it sacrifices type safety.

func anyType() {
	fmt.Println("\n=== 6. The any Type ===")

	var x any = 42
	fmt.Printf("any holding int: %v (%T)\n", x, x)

	x = "now it's a string"
	fmt.Printf("any holding string: %v (%T)\n", x, x)

	x = []int{1, 2, 3}
	fmt.Printf("any holding slice: %v (%T)\n", x, x)

	// To use the value, you must assert its type
	x = 100
	if n, ok := x.(int); ok {
		fmt.Printf("Asserted int: %d * 2 = %d\n", n, n*2)
	}

	// Type switch — handle multiple types cleanly
	values := []any{42, "hello", 3.14, true, []int{1, 2}}
	for _, v := range values {
		switch t := v.(type) {
		case int:
			fmt.Printf("int: %d\n", t)
		case string:
			fmt.Printf("string: %q\n", t)
		case float64:
			fmt.Printf("float64: %.2f\n", t)
		case bool:
			fmt.Printf("bool: %v\n", t)
		default:
			fmt.Printf("unknown: %T = %v\n", t, t)
		}
	}
}

func main() {
	integers()
	floats()
	strings_()
	booleans()
	typeAliasVsDefinition()
	anyType()
}

/*
THOUGHT QUESTIONS:

1. What happens when uint8(255) is incremented? Is this a bug or a feature?
   How does this differ from Python?

2. Why is `0.1 + 0.2 != 0.3` in Go (and every other language using IEEE 754)?
   How should you compare floating-point numbers?

3. `len("hello, 世界")` returns 13, not 9. Why? How do you get the
   number of characters (runes)?

4. What is the difference between `type Celsius float64` and `type Celsius = float64`?
   When would you use each?

5. What is "boxing" when you assign a value to `interface{}`/`any`?
   Why does this have a performance cost?

EXERCISES:

1. Write a function that reverses a string correctly, handling multi-byte
   Unicode characters (e.g., "Hello, 世界" → "界世 ,olleH").

2. Demonstrate integer overflow: write a loop that adds to an int8 starting
   from 120. Print values until you've gone past the maximum and wrapped around.

3. Create a function `typeDescribe(v any)` that accepts any value and prints
   a human-friendly description: "positive int: 42", "non-empty string: hello", etc.

4. Write a string builder that concatenates 10,000 strings using both `+`
   and `strings.Builder`. Compare their performance with a simple timer.
*/
