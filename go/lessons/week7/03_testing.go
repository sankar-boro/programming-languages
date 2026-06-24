/*
WEEK 7 — DAY 3: Testing in Go
===============================
Topic: Go's built-in testing framework — unit tests, benchmarks, examples, and table-driven tests.

Key ideas:
  - go test is built into the toolchain — no third-party runner needed
  - Test files end in _test.go and are never compiled into the binary
  - Test functions: func TestXxx(t *testing.T)
  - Benchmark functions: func BenchmarkXxx(b *testing.B)
  - Example functions: func ExampleXxx() with // Output: comments
  - Table-driven tests are the idiomatic Go testing pattern
  - -race flag detects data races during tests
*/

// NOTE: This is a demo/educational file — actual test code goes in *_test.go files.
// The patterns shown here work exactly as-is when placed in a *_test.go file.

package main

import "fmt"

// ─── 1. BASIC TEST STRUCTURE ──────────────────────────────────────────────────
//
// A test file is any file ending in _test.go
// Test functions: func TestXxx(t *testing.T) — must start with Test
//
// t.Error(msg)   — mark test as failed, continue running
// t.Fatal(msg)   — mark test as failed, STOP test immediately
// t.Errorf(fmt)  — like t.Error but with formatting
// t.Fatalf(fmt)  — like t.Fatal but with formatting
// t.Log(msg)     — print log (shown if test fails, or with -v)
// t.Skip(msg)    — skip this test
// t.Helper()     — mark as helper function (affects line numbers in errors)
//
// Run tests:
//   go test ./...           — run all tests
//   go test -v ./...        — verbose: show all test names
//   go test -run TestFoo    — run only tests matching TestFoo
//   go test -race ./...     — enable race detector

func basicTestStructure() {
	fmt.Println("=== 1. Test Structure ===")

	// This would be in math_test.go:
	example := `
package math_test

import (
    "testing"
    "yourmodule/math"
)

func TestAdd(t *testing.T) {
    result := math.Add(2, 3)
    if result != 5 {
        t.Errorf("Add(2, 3) = %d; want 5", result)
    }
}

func TestDivide(t *testing.T) {
    result, err := math.Divide(10, 2)
    if err != nil {
        t.Fatalf("Divide(10, 2) returned error: %v", err)
    }
    if result != 5 {
        t.Errorf("Divide(10, 2) = %f; want 5", result)
    }
}

func TestDivideByZero(t *testing.T) {
    _, err := math.Divide(10, 0)
    if err == nil {
        t.Error("expected error for division by zero, got nil")
    }
}
`
	fmt.Println(example)
}

// ─── 2. TABLE-DRIVEN TESTS ────────────────────────────────────────────────────
//
// The most idiomatic Go testing pattern: a table of inputs and expected outputs.
// Benefits:
//   - All test cases in one place (easy to add new cases)
//   - Less boilerplate
//   - Subtests make failures easy to identify by name

func tableDrivenTests() {
	fmt.Println("=== 2. Table-Driven Tests ===")

	example := `
func TestAdd(t *testing.T) {
    tests := []struct {
        name     string
        a, b     int
        expected int
    }{
        {"positive", 2, 3, 5},
        {"negative", -1, -2, -3},
        {"zero", 0, 0, 0},
        {"mixed", -5, 10, 5},
        {"large", 1000000, 2000000, 3000000},
    }

    for _, tt := range tests {
        // t.Run creates a SUBTEST — has its own pass/fail status
        t.Run(tt.name, func(t *testing.T) {
            got := Add(tt.a, tt.b)
            if got != tt.expected {
                t.Errorf("Add(%d, %d) = %d; want %d", tt.a, tt.b, got, tt.expected)
            }
        })
    }
}

// Run specific subtest:  go test -run TestAdd/positive
// Run all subtests:      go test -run TestAdd/
`
	fmt.Println(example)
}

// ─── 3. TESTING ERRORS AND PANICS ─────────────────────────────────────────────

func errorAndPanicTests() {
	fmt.Println("=== 3. Testing Errors and Panics ===")

	example := `
func TestDivide(t *testing.T) {
    tests := []struct {
        name      string
        a, b      float64
        want      float64
        wantErr   bool
    }{
        {"normal", 10, 2, 5, false},
        {"decimal", 7, 3, 0, false},   // we don't check exact float value
        {"by zero", 5, 0, 0, true},    // expect an error
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got, err := Divide(tt.a, tt.b)

            // Check error expectation
            if (err != nil) != tt.wantErr {
                t.Errorf("Divide() error = %v, wantErr %v", err, tt.wantErr)
                return  // don't check value if error expectation failed
            }

            if !tt.wantErr && math.Abs(got-tt.want) > 1e-9 {
                t.Errorf("Divide() = %v, want %v", got, tt.want)
            }
        })
    }
}

// Testing panics: use a helper
func assertPanics(t *testing.T, fn func()) {
    t.Helper()
    defer func() {
        if r := recover(); r == nil {
            t.Error("expected panic but didn't get one")
        }
    }()
    fn()
}

func TestMustPositivePanics(t *testing.T) {
    assertPanics(t, func() {
        MustPositive(-1)  // should panic
    })
}
`
	fmt.Println(example)
}

// ─── 4. BENCHMARKS ────────────────────────────────────────────────────────────
//
// Benchmarks measure performance. They run in a loop controlled by the framework.
//
// func BenchmarkXxx(b *testing.B)
//   b.N — the number of iterations (determined automatically by the framework)
//   Always put the thing being benchmarked in the b.N loop.
//
// Run:
//   go test -bench=.             — run all benchmarks
//   go test -bench=BenchmarkFoo  — run specific benchmark
//   go test -bench=. -benchmem   — also show memory allocations
//   go test -bench=. -count=3    — run 3 times for stability

func benchmarks() {
	fmt.Println("\n=== 4. Benchmarks ===")

	example := `
func BenchmarkFibIter(b *testing.B) {
    for i := 0; i < b.N; i++ {
        FibIter(30)  // benchmarks the iterative fibonacci
    }
}

func BenchmarkFibMemo(b *testing.B) {
    for i := 0; i < b.N; i++ {
        FibMemo(30)  // benchmarks the memoized fibonacci
    }
}

// Example output:
// BenchmarkFibIter-8    100000000     10.2 ns/op
// BenchmarkFibMemo-8     50000000     23.5 ns/op   2 allocs/op
//
// -8 = GOMAXPROCS, ns/op = nanoseconds per iteration, allocs/op = heap allocations

// Benchmark with setup:
func BenchmarkSort(b *testing.B) {
    // Setup outside the loop (not benchmarked)
    data := make([]int, 1000)
    for i := range data {
        data[i] = rand.Intn(1000)
    }

    b.ResetTimer()  // reset timer after setup
    for i := 0; i < b.N; i++ {
        // Make a copy each iteration (sort mutates)
        d := make([]int, len(data))
        copy(d, data)
        sort.Ints(d)
    }
}
`
	fmt.Println(example)
}

// ─── 5. EXAMPLES ──────────────────────────────────────────────────────────────
//
// Example functions serve as both tests and documentation.
// The // Output: comment is compared against actual stdout output.
// Examples appear in go doc output.

func examples() {
	fmt.Println("\n=== 5. Example Functions ===")

	example := `
func ExampleAdd() {
    fmt.Println(Add(2, 3))
    fmt.Println(Add(-1, 1))
    // Output:
    // 5
    // 0
}

// If output doesn't match, the test FAILS.
// This keeps documentation examples in sync with reality.

func ExampleStack_Push() {
    s := &Stack{}
    s.Push(10)
    s.Push(20)
    fmt.Println(s.Peek())
    // Output: 20 true
}

// Run: go test -run Example
// Or: go doc yourpkg.Add  — shows the example in documentation
`
	fmt.Println(example)
}

// ─── 6. TEST HELPERS AND CLEANUP ──────────────────────────────────────────────

func testHelpers() {
	fmt.Println("\n=== 6. Test Helpers and Cleanup ===")

	example := `
// t.Helper() marks the function as a test helper.
// When it fails, the error points to the CALLER, not the helper.
func assertNoError(t *testing.T, err error) {
    t.Helper()
    if err != nil {
        t.Fatalf("unexpected error: %v", err)
    }
}

func assertEqual[T comparable](t *testing.T, got, want T) {
    t.Helper()
    if got != want {
        t.Errorf("got %v, want %v", got, want)
    }
}

// t.Cleanup registers a function to run after the test (or subtest) completes.
// Works like defer but for tests.
func TestWithDatabase(t *testing.T) {
    db := setupTestDB(t)
    t.Cleanup(func() {
        db.Close()             // runs after test, even if test panics
        db.DropTestTables()
    })

    // Use db in the test...
}

// t.TempDir() creates a temporary directory that's deleted after the test.
func TestFileOperations(t *testing.T) {
    dir := t.TempDir()  // automatically cleaned up
    file := filepath.Join(dir, "test.txt")
    os.WriteFile(file, []byte("test"), 0644)
    // ...
}

// Parallel tests — run concurrently
func TestConcurrent(t *testing.T) {
    tests := []struct{ name string }{ ... }
    for _, tt := range tests {
        tt := tt  // capture range variable
        t.Run(tt.name, func(t *testing.T) {
            t.Parallel()  // this subtest runs in parallel with others
            // ...
        })
    }
}
`
	fmt.Println(example)
}

// ─── 7. RACE DETECTOR AND TEST COVERAGE ───────────────────────────────────────

func raceAndCoverage() {
	fmt.Println("\n=== 7. Race Detector and Coverage ===")

	fmt.Println("Race detector:")
	fmt.Println("  go test -race ./...")
	fmt.Println("  → Instruments all memory accesses")
	fmt.Println("  → Reports data races at runtime")
	fmt.Println("  → ~5-10x slower than normal — only for testing")

	fmt.Println("\nCoverage:")
	fmt.Println("  go test -cover ./...")
	fmt.Println("  → ok  mypackage  0.015s  coverage: 87.5% of statements")
	fmt.Println()
	fmt.Println("  go test -coverprofile=coverage.out ./...")
	fmt.Println("  go tool cover -html=coverage.out")
	fmt.Println("  → Open visual HTML report showing uncovered lines in red")

	fmt.Println("\nFuzzing (Go 1.18+):")
	fmt.Println("  func FuzzAdd(f *testing.F) {")
	fmt.Println("      f.Add(2, 3)  // seed corpus")
	fmt.Println("      f.Fuzz(func(t *testing.T, a, b int) {")
	fmt.Println("          Add(a, b)  // should never panic")
	fmt.Println("      })")
	fmt.Println("  }")
	fmt.Println("  go test -fuzz=FuzzAdd")
}

func main() {
	basicTestStructure()
	tableDrivenTests()
	errorAndPanicTests()
	benchmarks()
	examples()
	testHelpers()
	raceAndCoverage()
}

/*
THOUGHT QUESTIONS:

1. Why are test files (*_test.go) never compiled into the production binary?
   What benefits does this separation provide?

2. What is the advantage of table-driven tests over separate test functions
   for each case?

3. The benchmark loop `for i := 0; i < b.N; i++` — why does the framework
   control N rather than letting you choose a fixed number of iterations?

4. Why does `t.Helper()` matter? What happens to error line numbers without it?

5. What is fuzz testing? What class of bugs can it find that unit tests miss?

EXERCISES:

1. Write a complete test file for a `Calculator` struct with Add, Sub, Mul, Div.
   Use table-driven tests. Test all edge cases including division by zero.

2. Write a benchmark comparing three string concatenation methods:
   (a) + operator in a loop, (b) strings.Builder, (c) bytes.Buffer.
   Run with -benchmem to see allocations.

3. Write an Example function for a custom type that implements fmt.Stringer.
   Verify the output comment matches by running go test.

4. Write a parallel test that spawns 100 goroutines all calling a shared
   function. Run with -race to detect and fix any races.
*/
