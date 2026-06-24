/*
WEEK 3 — DAY 3: Error Handling
================================
Topic: Go's error philosophy — explicit errors, not exceptions.

Key ideas:
  - Errors are VALUES — functions return them like any other value
  - The error interface has one method: Error() string
  - Go has NO exceptions (panic/recover is different — not for normal control flow)
  - errors.Is and errors.As enable structured error inspection
  - Wrap errors with %w to preserve the error chain
*/

package main

import (
	"errors"
	"fmt"
	"os"
	"strconv"
)

// ─── 1. THE error INTERFACE ────────────────────────────────────────────────────
//
// error is a built-in interface:
//   type error interface {
//       Error() string
//   }
//
// Any type that implements Error() string satisfies error.
// This simple design allows errors to carry any information.

func errorInterface() {
	fmt.Println("=== 1. The error Interface ===")

	// errors.New creates a simple string error
	err := errors.New("something went wrong")
	fmt.Printf("error: %v\n", err)
	fmt.Printf("type: %T\n", err)  // *errors.errorString

	// fmt.Errorf creates a formatted error
	name := "alice"
	err2 := fmt.Errorf("user %q not found", name)
	fmt.Printf("formatted error: %v\n", err2)

	// nil error means no error
	var noErr error = nil
	fmt.Printf("no error: %v\n", noErr)
	fmt.Printf("no error == nil: %v\n", noErr == nil)
}

// ─── 2. RETURNING AND CHECKING ERRORS ────────────────────────────────────────
//
// Idiomatic Go: return the error as the LAST return value.
// Check errors immediately after the call.

func readUserAge(input string) (int, error) {
	age, err := strconv.Atoi(input)
	if err != nil {
		return 0, fmt.Errorf("invalid age %q: %w", input, err)
	}
	if age < 0 || age > 150 {
		return 0, fmt.Errorf("age %d is out of valid range [0, 150]", age)
	}
	return age, nil
}

func returningErrors() {
	fmt.Println("\n=== 2. Returning and Checking Errors ===")

	inputs := []string{"25", "abc", "-5", "200", "42"}
	for _, input := range inputs {
		age, err := readUserAge(input)
		if err != nil {
			fmt.Printf("  error for %q: %v\n", input, err)
			continue
		}
		fmt.Printf("  valid age: %d\n", age)
	}
}

// ─── 3. CUSTOM ERROR TYPES ────────────────────────────────────────────────────
//
// Implement the error interface to create errors that carry structured data.
// This allows callers to inspect the error and take specific actions.

// ValidationError carries field-specific error information
type ValidationError struct {
	Field   string
	Message string
	Value   any
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("validation failed: field=%q value=%v: %s", e.Field, e.Value, e.Message)
}

// NotFoundError for missing resources
type NotFoundError struct {
	Resource string
	ID       any
}

func (e *NotFoundError) Error() string {
	return fmt.Sprintf("%s with id=%v not found", e.Resource, e.ID)
}

type User struct {
	ID    int
	Name  string
	Email string
	Age   int
}

func validateUser(u User) error {
	if u.Name == "" {
		return &ValidationError{Field: "Name", Message: "cannot be empty", Value: u.Name}
	}
	if u.Age < 0 || u.Age > 150 {
		return &ValidationError{Field: "Age", Message: "must be 0-150", Value: u.Age}
	}
	if u.Email == "" {
		return &ValidationError{Field: "Email", Message: "cannot be empty", Value: u.Email}
	}
	return nil
}

var userDB = map[int]User{
	1: {1, "Alice", "alice@example.com", 30},
}

func getUser(id int) (*User, error) {
	u, ok := userDB[id]
	if !ok {
		return nil, &NotFoundError{Resource: "User", ID: id}
	}
	return &u, nil
}

func customErrorTypes() {
	fmt.Println("\n=== 3. Custom Error Types ===")

	// Validate
	err := validateUser(User{Name: "", Age: -5})
	if err != nil {
		fmt.Println("validation:", err)
	}

	err = validateUser(User{Name: "Bob", Email: "bob@example.com", Age: 25})
	if err == nil {
		fmt.Println("user valid")
	}

	// Get user
	_, err = getUser(99)
	if err != nil {
		fmt.Println("get user:", err)
	}

	u, err := getUser(1)
	if err == nil {
		fmt.Printf("found user: %+v\n", *u)
	}
}

// ─── 4. errors.Is AND errors.As ──────────────────────────────────────────────
//
// When errors are wrapped with %w, the chain is preserved.
// errors.Is(err, target) checks if ANY error in the chain matches target
// errors.As(err, &target) extracts the first error in the chain of the given type
//
// These are the correct way to check for specific errors — NOT == comparison
// on wrapped errors.

var ErrPermission = errors.New("permission denied")
var ErrNotFound = errors.New("not found")

func accessFile(path string, hasPermission bool) error {
	if !hasPermission {
		return fmt.Errorf("accessing %q: %w", path, ErrPermission)
	}
	_, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("opening %q: %w", path, ErrNotFound)
	}
	return nil
}

func errorsIsAs() {
	fmt.Println("\n=== 4. errors.Is and errors.As ===")

	// Sentinel error with errors.Is
	err := accessFile("/secret", false)
	if err != nil {
		fmt.Println("error:", err)

		// Check if any error in the chain IS ErrPermission
		if errors.Is(err, ErrPermission) {
			fmt.Println("  → permission denied!")
		}
		if errors.Is(err, ErrNotFound) {
			fmt.Println("  → not found (won't print)")
		}
	}

	// Custom error type with errors.As
	err2 := validateUser(User{Age: -1, Name: "test"})
	if err2 != nil {
		var valErr *ValidationError
		if errors.As(err2, &valErr) {
			fmt.Printf("  → validation error on field %q: %s\n", valErr.Field, valErr.Message)
		}
	}

	// Wrapped custom error
	outerErr := fmt.Errorf("process failed: %w", &NotFoundError{Resource: "Config", ID: "main"})
	var notFound *NotFoundError
	if errors.As(outerErr, &notFound) {
		fmt.Printf("  → extracted: %s id=%v\n", notFound.Resource, notFound.ID)
	}
}

// ─── 5. WRAPPING AND UNWRAPPING ERRORS ────────────────────────────────────────
//
// fmt.Errorf with %w wraps an error (preserves the chain).
// errors.Unwrap peels one layer.
// errors.Is/As traverse the entire chain.
//
// Build error chains to add context at each layer:
//   "http request failed: connection timeout: dial failed: no route to host"

func readConfig(path string) error {
	_, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("readConfig: %w", err)  // wrap with context
	}
	return nil
}

func initApp(configPath string) error {
	if err := readConfig(configPath); err != nil {
		return fmt.Errorf("initApp: %w", err)  // add another layer
	}
	return nil
}

func errorWrapping() {
	fmt.Println("\n=== 5. Error Wrapping ===")

	err := initApp("/nonexistent/config.json")
	if err != nil {
		fmt.Println("error:", err)

		// Unwrap one layer at a time
		for err != nil {
			fmt.Printf("  chain: %T: %v\n", err, err)
			err = errors.Unwrap(err)
		}
	}
}

// ─── 6. PANIC AND RECOVER ─────────────────────────────────────────────────────
//
// panic: stops the current goroutine's normal execution, unwinds the stack,
//        and runs deferred functions. If not recovered, crashes the program.
//
// recover: called INSIDE a deferred function to stop the panic and get
//          the panic value. Outside of defer, recover always returns nil.
//
// Use panic/recover for:
//   - Truly unexpected conditions (programming errors, impossible states)
//   - Library code that needs to convert panics to errors for callers
// DO NOT use for normal control flow (that's what errors are for).

func safeDivide(a, b int) (result int, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("recovered from panic: %v", r)
		}
	}()
	if b == 0 {
		panic("division by zero")  // unreachable in practice with check
	}
	return a / b, nil
}

func panicAndRecover() {
	fmt.Println("\n=== 6. Panic and Recover ===")

	// Normal case
	result, err := safeDivide(10, 2)
	fmt.Printf("10/2 = %d, err=%v\n", result, err)

	// Would panic, but recover catches it
	result, err = safeDivide(10, 0)
	fmt.Printf("10/0 = %d, err=%v\n", result, err)

	// Demonstrate panic propagation
	fmt.Println("Calling mustPositive(5):", mustPositive(5))
	fmt.Println("Calling with recovery:")
	func() {
		defer func() {
			if r := recover(); r != nil {
				fmt.Println("  recovered:", r)
			}
		}()
		mustPositive(-1)
	}()
	fmt.Println("Program continues after recovery")
}

func mustPositive(n int) int {
	if n < 0 {
		panic(fmt.Sprintf("expected positive, got %d", n))
	}
	return n
}

func main() {
	errorInterface()
	returningErrors()
	customErrorTypes()
	errorsIsAs()
	errorWrapping()
	panicAndRecover()
}

/*
THOUGHT QUESTIONS:

1. Why does Go use return values for errors instead of exceptions?
   What are the trade-offs of each approach?

2. What is the difference between `err == ErrPermission` and
   `errors.Is(err, ErrPermission)` when using wrapped errors?

3. When should you use panic vs returning an error?
   What is the "library code should not panic" convention?

4. What is the errors.As function for? How is it different from errors.Is?

5. The error chain: "initApp: readConfig: open /file: no such file"
   What does each layer add? Why is context at each level valuable?

EXERCISES:

1. Create a custom error type `HTTPError` with StatusCode and Body fields.
   Write functions that return these errors at different layers, wrapping them.
   Then use errors.As to extract the status code at the top level.

2. Write a `retry(n int, fn func() error) error` function that calls fn up to
   n times, stopping on success. Return the last error if all attempts fail.

3. Implement a simple file parser that returns structured validation errors
   for each line that doesn't match the expected format. Collect multiple
   errors and return them all at once (hint: use a slice of errors).

4. Write a function `mustRun(fn func() error)` that calls fn and panics
   if it returns an error. Then wrap it in a Safe version using recover.
*/
