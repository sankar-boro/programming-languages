/*
WEEK 5 — DAY 2: Interfaces In Depth
=====================================
Topic: How interfaces work internally — the iface/eface pair, duck typing,
       interface satisfaction, and the power of small interfaces.

Key ideas:
  - An interface value is a (type, value) pair internally
  - Interface satisfaction is implicit — no `implements` keyword
  - The empty interface (any) holds any value (boxing)
  - Small interfaces are idiomatic Go (io.Reader, io.Writer = 1 method each)
  - Interface composition builds complex abstractions from simple ones
*/

package main

import (
	"bytes"
	"fmt"
	"io"
	"math"
	"os"
	"strings"
)

// ─── 1. INTERFACE INTERNALS: THE (TYPE, VALUE) PAIR ──────────────────────────
//
// An interface value in Go is TWO words (16 bytes on 64-bit):
//
//   Non-empty interface (iface):
//     word 1: *itab  — pointer to interface table (type descriptor + method pointers)
//     word 2: *data  — pointer to the actual value (or the value itself if ≤ ptr size)
//
//   Empty interface (eface, aka any):
//     word 1: *_type — pointer to runtime type information
//     word 2: *data  — pointer to the actual value
//
// The itab (interface table) contains:
//   - The concrete type
//   - The interface type
//   - Pointers to the concrete type's methods that satisfy the interface
//
// Calling an interface method = one extra indirection (load method pointer from itab)
// This is slightly slower than a direct call but enables polymorphism.
//
// nil interface check:
//   var r io.Reader = nil   // both words are nil → r == nil → true
//   var b *bytes.Buffer = nil
//   var r io.Reader = b    // word1 = *itab (non-nil!), word2 = nil → r == nil → FALSE!
//   This is the famous "nil interface" gotcha.

func interfaceInternals() {
	fmt.Println("=== 1. Interface Internals: (type, value) pair ===")

	var r io.Reader

	// Both parts are nil
	fmt.Printf("var r io.Reader: r==nil is %v\n", r == nil)

	// Assign a concrete type
	r = strings.NewReader("hello")
	fmt.Printf("r = strings.NewReader: r==nil is %v, type=%T\n", r == nil, r)

	// The nil interface gotcha
	var buf *bytes.Buffer = nil  // nil pointer to a concrete type
	r = buf                      // type part is non-nil (*bytes.Buffer), value part is nil
	fmt.Printf("r = (*bytes.Buffer)(nil): r==nil is %v (FALSE — gotcha!)\n", r == nil)

	// Correct fix: check both
	if r != nil {
		// Even though buf is nil, this branch is taken
		// Calling r.Read would panic (nil pointer dereference)
		fmt.Println("r is non-nil (but value is nil — dangerous!)")
	}
}

// ─── 2. IMPLICIT SATISFACTION — DUCK TYPING ──────────────────────────────────
//
// In Go, a type satisfies an interface if it has ALL the interface's methods.
// There is no `implements` keyword.
//
// This is structural typing (for interfaces): "if it walks like a duck..."
// You can satisfy an interface you didn't know existed at the time of writing.

type Shape interface {
	Area() float64
	Perimeter() float64
}

type Circle struct {
	Radius float64
}

func (c Circle) Area() float64      { return math.Pi * c.Radius * c.Radius }
func (c Circle) Perimeter() float64 { return 2 * math.Pi * c.Radius }

type Rectangle struct {
	Width, Height float64
}

func (r Rectangle) Area() float64      { return r.Width * r.Height }
func (r Rectangle) Perimeter() float64 { return 2 * (r.Width + r.Height) }

type Triangle struct {
	A, B, C float64  // side lengths
}

func (t Triangle) Area() float64 {
	s := (t.A + t.B + t.C) / 2  // Heron's formula
	return math.Sqrt(s * (s - t.A) * (s - t.B) * (s - t.C))
}

func (t Triangle) Perimeter() float64 { return t.A + t.B + t.C }

func printShapeInfo(s Shape) {
	fmt.Printf("  %T: area=%.2f perimeter=%.2f\n", s, s.Area(), s.Perimeter())
}

func implicitSatisfaction() {
	fmt.Println("\n=== 2. Implicit Satisfaction ===")

	shapes := []Shape{
		Circle{Radius: 5},
		Rectangle{Width: 4, Height: 6},
		Triangle{A: 3, B: 4, C: 5},
	}

	totalArea := 0.0
	for _, s := range shapes {
		printShapeInfo(s)
		totalArea += s.Area()
	}
	fmt.Printf("Total area: %.2f\n", totalArea)
}

// ─── 3. SMALL INTERFACES — THE GO PHILOSOPHY ──────────────────────────────────
//
// "The bigger the interface, the weaker the abstraction." — Rob Pike
//
// Go's most powerful interfaces have 1–2 methods:
//
//   io.Reader:  Read(p []byte) (n int, err error)
//   io.Writer:  Write(p []byte) (n int, err error)
//   io.Closer:  Close() error
//   error:      Error() string
//   fmt.Stringer: String() string
//   sort.Interface: Len() int; Less(i, j int) bool; Swap(i, j int)
//
// This lets you write functions that accept ANY type that can be read from,
// written to, closed, etc. — regardless of what that type actually is.

func smallInterfaces() {
	fmt.Println("\n=== 3. Small Interfaces ===")

	// io.Reader works with: files, strings, bytes, network connections, HTTP bodies...
	readers := []io.Reader{
		strings.NewReader("from a string"),
		bytes.NewReader([]byte("from bytes")),
	}

	for _, r := range readers {
		buf, _ := io.ReadAll(r)
		fmt.Printf("  Read from %T: %q\n", r, string(buf))
	}

	// io.Writer works with: files, bytes.Buffer, stdout, network connections...
	writers := []io.Writer{
		os.Stdout,
		&bytes.Buffer{},
	}

	for _, w := range writers {
		fmt.Fprintf(w, "writing to %T\n", w)
	}

	// fmt.Stringer — your types can be printed naturally
	c := Circle{Radius: 3}
	fmt.Println("Circle:", c.String())  // wait, we need to add String()
}

func (c Circle) String() string {
	return fmt.Sprintf("Circle(r=%.1f)", c.Radius)
}

func (r Rectangle) String() string {
	return fmt.Sprintf("Rect(%.1fx%.1f)", r.Width, r.Height)
}

// ─── 4. INTERFACE COMPOSITION ─────────────────────────────────────────────────
//
// Interfaces can embed other interfaces to build larger contracts.
// This is how io.ReadWriter, io.ReadWriteCloser are defined.

type ReadWriter interface {
	io.Reader
	io.Writer
}

type ReadWriteCloser interface {
	io.Reader
	io.Writer
	io.Closer
}

// Define our own composed interfaces
type Sizer interface {
	Size() int64
}

type SizedReader interface {
	io.Reader
	Sizer
}

// A buffer that knows its own size
type SizedBuffer struct {
	buf  bytes.Buffer
	size int64
}

func (s *SizedBuffer) Write(p []byte) (n int, err error) {
	n, err = s.buf.Write(p)
	s.size += int64(n)
	return
}

func (s *SizedBuffer) Read(p []byte) (int, error) { return s.buf.Read(p) }
func (s *SizedBuffer) Size() int64                 { return s.size }

func interfaceComposition() {
	fmt.Println("\n=== 4. Interface Composition ===")

	sb := &SizedBuffer{}

	// Write some data
	fmt.Fprintf(sb, "Hello, World! ")
	fmt.Fprintf(sb, "This is a test.")
	fmt.Printf("Written %d bytes\n", sb.Size())

	// Read it back
	data, _ := io.ReadAll(sb)
	fmt.Printf("Read back: %q\n", string(data))
}

// ─── 5. TYPE ASSERTIONS AND TYPE SWITCHES ─────────────────────────────────────

type Logger interface {
	Log(msg string)
}

type DebugLogger interface {
	Logger
	Debug(msg string)
}

type AppLogger struct{ prefix string }

func (l *AppLogger) Log(msg string)   { fmt.Printf("[%s] %s\n", l.prefix, msg) }
func (l *AppLogger) Debug(msg string) { fmt.Printf("[%s][DEBUG] %s\n", l.prefix, msg) }

func logSomething(l Logger, msg string) {
	l.Log(msg)

	// Type assertion: does l also implement DebugLogger?
	if dl, ok := l.(DebugLogger); ok {
		dl.Debug("(additional debug info)")
	}
}

func typeAssertions() {
	fmt.Println("\n=== 5. Type Assertions and Type Switches ===")

	var l Logger = &AppLogger{prefix: "APP"}
	logSomething(l, "starting application")

	// Type switch on interface
	values := []any{42, "hello", Circle{5}, Rectangle{3, 4}, nil}
	for _, v := range values {
		switch t := v.(type) {
		case int:
			fmt.Printf("  int: %d\n", t)
		case string:
			fmt.Printf("  string: %q\n", t)
		case Shape:
			fmt.Printf("  Shape: area=%.2f\n", t.Area())
		case nil:
			fmt.Println("  nil")
		default:
			fmt.Printf("  unknown: %T\n", t)
		}
	}
}

// ─── 6. INTERFACE SATISFACTION CHECK AT COMPILE TIME ─────────────────────────
//
// To verify that a type satisfies an interface at compile time
// (without creating an instance), use a blank assignment:

// Compile-time assertion
var _ Shape = Circle{}       // Circle must implement Shape
var _ Shape = Rectangle{}    // Rectangle must implement Shape
var _ io.Writer = (*bytes.Buffer)(nil)  // *bytes.Buffer must implement io.Writer

// This pattern is used in libraries to ensure types implement interfaces correctly.
// If Circle were missing the Perimeter() method, the compiler would catch it here.

func compileTimeChecks() {
	fmt.Println("\n=== 6. Compile-Time Interface Checks ===")
	fmt.Println("var _ Shape = Circle{} // enforced at compile time")
	fmt.Println("If Circle were missing a method, this file wouldn't compile")
}

// ─── 7. THE POWER OF io.Reader / io.Writer ────────────────────────────────────
//
// Demonstrate the extreme composability of io.Reader and io.Writer.
// io.TeeReader, io.MultiWriter, io.LimitReader etc all work on the interface.

func ioComposability() {
	fmt.Println("\n=== 7. io.Reader / io.Writer Composability ===")

	// Write to multiple destinations at once (stdout + buffer)
	var buf bytes.Buffer
	w := io.MultiWriter(os.Stdout, &buf)
	fmt.Fprintln(w, "This goes to stdout AND the buffer simultaneously")

	// Read and capture (TeeReader reads r, writes to w, returns a reader)
	src := strings.NewReader("source data")
	var capture bytes.Buffer
	tee := io.TeeReader(src, &capture)

	// Read from tee — data also flows into capture
	data, _ := io.ReadAll(tee)
	fmt.Printf("Read: %q\n", string(data))
	fmt.Printf("Captured: %q\n", capture.String())

	// Limit reader
	limited := io.LimitReader(strings.NewReader("hello world"), 5)
	small, _ := io.ReadAll(limited)
	fmt.Printf("LimitReader(5): %q\n", string(small))
}

func main() {
	interfaceInternals()
	implicitSatisfaction()
	smallInterfaces()
	interfaceComposition()
	typeAssertions()
	compileTimeChecks()
	ioComposability()
}

/*
THOUGHT QUESTIONS:

1. An interface value is a (type, value) pair. What is the difference between:
   var r io.Reader = nil  (r == nil is TRUE)
   var b *bytes.Buffer = nil; var r io.Reader = b  (r == nil is FALSE)

2. Why does Go use implicit interface satisfaction instead of explicit `implements`?
   What is the main advantage? What is lost?

3. The "bigger the interface, the weaker the abstraction." What does this mean?
   Give an example of a well-designed single-method interface.

4. io.Reader can be satisfied by files, strings, bytes, network connections, HTTP bodies.
   What makes this possible? What single constraint ties them all together?

5. What is a compile-time interface check (`var _ Shape = Circle{}`)?
   Why is this useful in a library?

EXERCISES:

1. Write a `Cache` interface with Get(key string) (string, bool) and
   Set(key, value string). Implement it with both a map-based (in-memory)
   and a file-based version. Write a function that works with either.

2. Implement a `ProgressReader` that wraps an io.Reader and calls a callback
   with the number of bytes read so far (for tracking download progress).

3. Write a `MultiError` type that implements the error interface and holds
   multiple errors. Add an `Unwrap() []error` method (Go 1.20+).

4. Create a `Validator` interface with `Validate() error`. Implement it for
   User, Product, and Order structs. Write `ValidateAll(items []Validator) []error`.
*/
