/*
WEEK 3 — DAY 1: Control Flow — if, for, switch, defer
=======================================================
Topic: Go's control flow constructs and their internal behavior.

Key ideas:
  - Go has ONE loop construct: `for` (no while, do-while)
  - `if` can have an initialization statement
  - `switch` is powerful: no fallthrough by default, can switch on any comparable
  - `defer` is LIFO, captures arguments eagerly but runs body lazily
  - `break`, `continue`, and `goto` support labels
*/

package main

import (
	"fmt"
	"os"
)

// ─── 1. IF STATEMENTS ─────────────────────────────────────────────────────────
//
// Go's if is straightforward with one unique feature: an initialization
// statement before the condition, scoped to the if/else block.

func ifStatements() {
	fmt.Println("=== 1. If Statements ===")

	x := 42

	// Simple if
	if x > 40 {
		fmt.Println("x > 40")
	}

	// if-else
	if x%2 == 0 {
		fmt.Println("x is even")
	} else {
		fmt.Println("x is odd")
	}

	// if-else if chain
	score := 75
	if score >= 90 {
		fmt.Println("A")
	} else if score >= 80 {
		fmt.Println("B")
	} else if score >= 70 {
		fmt.Println("C")
	} else {
		fmt.Println("F")
	}

	// Init statement — scopes val to the if/else block
	// This is idiomatic Go — used with function calls that return errors
	if val, err := divide(10, 3); err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Printf("10/3 = %.4f\n", val)
	}
	// val is not accessible here — scoped to the if/else block

	// Common pattern: error handling
	if err := writeFile("test.txt"); err != nil {
		fmt.Println("write failed:", err)
	}
}

func divide(a, b float64) (float64, error) {
	if b == 0 {
		return 0, fmt.Errorf("division by zero")
	}
	return a / b, nil
}

func writeFile(name string) error {
	return os.WriteFile(name, []byte("hello"), 0644)
}

// ─── 2. FOR LOOPS — GO'S ONLY LOOP ───────────────────────────────────────────
//
// Go has ONE loop keyword: `for`. It covers:
//   - C-style 3-clause loop: for init; condition; post { }
//   - While-style: for condition { }
//   - Infinite loop: for { }
//   - Range loop: for i, v := range collection { }

func forLoops() {
	fmt.Println("\n=== 2. For Loops ===")

	// Classic C-style for
	for i := 0; i < 5; i++ {
		fmt.Printf("%d ", i)
	}
	fmt.Println()

	// While-style
	n := 1
	for n < 100 {
		n *= 2
	}
	fmt.Printf("first power of 2 >= 100: %d\n", n)

	// Infinite loop with break
	count := 0
	for {
		count++
		if count >= 5 {
			break
		}
	}
	fmt.Printf("counted to: %d\n", count)

	// Range over slice
	fruits := []string{"apple", "banana", "cherry"}
	for i, fruit := range fruits {
		fmt.Printf("  [%d] %s\n", i, fruit)
	}

	// Range with only index
	for i := range fruits {
		fmt.Printf("  index: %d\n", i)
	}

	// Range with blank identifier (value only)
	for _, fruit := range fruits {
		fmt.Printf("  fruit: %s\n", fruit)
	}

	// Range over string — iterates by RUNE (not byte)
	for i, r := range "Hello, 世界" {
		fmt.Printf("  byte_index=%d rune=%c\n", i, r)
	}

	// Range over map (random order)
	m := map[string]int{"a": 1, "b": 2}
	for k, v := range m {
		fmt.Printf("  %s=%d\n", k, v)
	}

	// Range over channel
	ch := make(chan int, 3)
	ch <- 10; ch <- 20; ch <- 30
	close(ch)
	for v := range ch {
		fmt.Printf("  from channel: %d\n", v)
	}
}

// ─── 3. BREAK AND CONTINUE WITH LABELS ────────────────────────────────────────
//
// Labels allow break/continue to target outer loops — useful for nested loops.

func labelsAndBreak() {
	fmt.Println("\n=== 3. Labels, break, continue ===")

	// Without labels — break exits the INNERMOST loop
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			if j == 1 {
				break  // only breaks inner loop
			}
			fmt.Printf("(%d,%d) ", i, j)
		}
	}
	fmt.Println("(break inner only)")

	// With label — break exits the OUTER loop
outer:
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			if i == 1 && j == 1 {
				break outer  // exits the outer loop
			}
			fmt.Printf("(%d,%d) ", i, j)
		}
	}
	fmt.Println("(break outer)")

	// continue with label — skip to next iteration of outer loop
	fmt.Println("continue with label:")
loop:
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			if j == 1 {
				continue loop  // skip to next i
			}
			fmt.Printf("(%d,%d) ", i, j)
		}
	}
	fmt.Println()
}

// ─── 4. SWITCH STATEMENTS ─────────────────────────────────────────────────────
//
// Go's switch is more powerful than C's:
//   - No fallthrough by default (cases don't fall into the next case)
//   - Each case can have multiple values
//   - Cases can be expressions, not just constants
//   - Explicit fallthrough is possible with the `fallthrough` keyword
//   - Switch on no expression = switch true (used as if-else chain)

func switchStatements() {
	fmt.Println("\n=== 4. Switch Statements ===")

	// Basic switch
	day := "Tuesday"
	switch day {
	case "Monday", "Tuesday", "Wednesday", "Thursday", "Friday":
		fmt.Println("Weekday")
	case "Saturday", "Sunday":
		fmt.Println("Weekend")
	default:
		fmt.Println("Unknown")
	}

	// Switch with expressions in cases
	x := 15
	switch {
	case x < 0:
		fmt.Println("negative")
	case x == 0:
		fmt.Println("zero")
	case x < 10:
		fmt.Println("single digit")
	case x < 100:
		fmt.Println("double digit")
	default:
		fmt.Println("large")
	}

	// Switch with init statement
	switch n := len("hello"); {
	case n < 3:
		fmt.Println("short")
	case n < 7:
		fmt.Println("medium")  // matches
	default:
		fmt.Println("long")
	}

	// Explicit fallthrough (uncommon in Go)
	switch 1 {
	case 1:
		fmt.Print("one ")
		fallthrough  // explicitly falls to next case
	case 2:
		fmt.Print("two ")
		fallthrough
	case 3:
		fmt.Println("three")
	}

	// Type switch — matches on concrete type of an interface value
	typeSwitch(42)
	typeSwitch("hello")
	typeSwitch(3.14)
	typeSwitch([]int{1, 2})
}

func typeSwitch(v any) {
	switch t := v.(type) {
	case int:
		fmt.Printf("int: %d\n", t)
	case string:
		fmt.Printf("string: %q\n", t)
	case float64:
		fmt.Printf("float64: %.2f\n", t)
	default:
		fmt.Printf("other: %T = %v\n", t, t)
	}
}

// ─── 5. DEFER — LIFO, EAGER ARG EVALUATION ───────────────────────────────────
//
// defer schedules a function call to run WHEN THE SURROUNDING FUNCTION RETURNS.
// This makes it perfect for cleanup: close files, release locks, recover panics.
//
// CRITICAL rules:
//
//   1. Deferred calls are executed in LIFO (last-in, first-out) order.
//      Think of them as a stack of cleanup functions.
//
//   2. Deferred function ARGUMENTS are evaluated IMMEDIATELY when defer is called,
//      but the function BODY runs when the enclosing function returns.
//
//   3. Deferred functions can READ AND MODIFY named return values.
//      This is used for "cleanup" that modifies the return value.

func deferBasics() {
	fmt.Println("\n=== 5. Defer ===")

	// LIFO order
	fmt.Println("Defer LIFO order:")
	defer fmt.Println("  first defer (runs last)")
	defer fmt.Println("  second defer (runs second)")
	defer fmt.Println("  third defer (runs first)")
	fmt.Println("  (function body done)")
}

func deferArgEvaluation() {
	fmt.Println("\n--- Defer arg evaluation ---")
	x := 1
	defer fmt.Printf("  deferred x: %d (captured when defer ran, x was 1)\n", x)
	x = 99
	fmt.Printf("  x at end of function: %d\n", x)
	// The deferred call sees x=1 (captured at defer time), not 99
}

func deferLoop() {
	fmt.Println("\n--- Defer in loops (WRONG!) ---")
	// If you defer in a loop, ALL defers run when the function returns — not per iteration
	// This is a common mistake. Use a closure or immediately-invoked function instead.
	for i := 0; i < 3; i++ {
		// defer fmt.Println("cleanup", i)  // would run 3 times when function returns
		func(i int) {
			// use a closure to scope the defer per iteration
			defer fmt.Printf("  cleanup for iteration %d\n", i)
		}(i)
	}
}

func deferWithNamedReturn() (result string) {
	defer func() {
		// This closure can read AND modify the named return value
		if result == "" {
			result = "default"  // set a default if function returned ""
		}
		fmt.Println("  deferred: result =", result)
	}()

	// Simulate not setting result
	return  // returns "" → defer modifies it to "default"
}

func deferCleanup() {
	fmt.Println("\n--- Defer for resource cleanup ---")
	// The idiomatic Go pattern: acquire resource, immediately defer release
	file, err := os.Open("test.txt")
	if err != nil {
		fmt.Println("  (test.txt not found, simulating cleanup pattern)")
		return
	}
	defer file.Close()  // guaranteed to run even if we panic below
	fmt.Println("  file opened, defer will close it")
}

// ─── 6. GOTO (RARELY USED) ────────────────────────────────────────────────────
//
// Go has goto. It's rarely used. Restrictions:
//   - Cannot jump over variable declarations
//   - Cannot jump into a block
// Mostly seen in machine-generated code or low-level error handling.

func gotoExample() {
	fmt.Println("\n=== 6. goto (rare) ===")
	i := 0
loop:
	if i < 3 {
		fmt.Printf("  i = %d\n", i)
		i++
		goto loop
	}
	fmt.Println("  done")
}

func main() {
	ifStatements()
	forLoops()
	labelsAndBreak()
	switchStatements()
	deferBasics()
	deferArgEvaluation()
	deferLoop()
	fmt.Println("\n--- Defer with named return ---")
	r := deferWithNamedReturn()
	fmt.Println("  final result:", r)
	deferCleanup()
	gotoExample()
}

/*
THOUGHT QUESTIONS:

1. Why does Go have only ONE loop keyword (`for`) instead of `while`, `do-while`,
   and `for`? What are the trade-offs?

2. Go's switch doesn't fall through by default. What are the benefits of this
   choice over C's behavior?

3. Explain: "Deferred arguments are evaluated eagerly, but the body runs lazily."
   Write a code example that demonstrates this distinction.

4. Why is `defer` preferred over try/finally for resource cleanup in Go?
   What advantage does it have?

5. Defers run in LIFO order. Why is this order sensible for resource cleanup?
   Give an example where FIFO order would cause a bug.

EXERCISES:

1. Write a function that opens two files, reads the first, and writes to the
   second. Use defer to ensure both files are closed in the correct order,
   even if an error occurs.

2. Write a function `measureTime(name string)` that uses defer to print
   how long a function took to execute. Usage: `defer measureTime("myFunc")()`.

3. Write a function that uses a labeled break to search for a target value
   in a 2D matrix and immediately exits all loops when found.

4. Using a type switch, write a `prettyPrint(v any)` function that handles
   int, float64, string, bool, []any, and map[string]any (JSON-like structure)
   with proper indentation for nested values.
*/
