/*
WEEK 2 — DAY 3: Structs — Memory Layout, Methods, and Embedding
================================================================
Topic: How structs work in Go — value semantics, memory layout, embedding.

Key ideas:
  - Structs group related data; they are value types (copied on assignment)
  - Memory layout: fields are laid out sequentially with alignment padding
  - Struct tags are metadata read at runtime (JSON, SQL, etc.)
  - Embedding provides composition — Go's substitute for inheritance
  - Anonymous fields promote methods and fields to the outer struct
*/

package main

import (
	"encoding/json"
	"fmt"
	"unsafe"
)

// ─── 1. STRUCT BASICS AND MEMORY LAYOUT ───────────────────────────────────────
//
// A struct's memory layout:
//   - Fields are laid out in declaration order
//   - Each field must be ALIGNED to a multiple of its size
//   - Padding bytes are inserted between fields to satisfy alignment
//   - The struct's total size is padded to a multiple of its largest field
//
// This matters for:
//   - Memory efficiency (reorder fields to minimize padding)
//   - C interop (must match C struct layout exactly)
//   - Cache behavior (smaller structs = more fit in cache)

type BadLayout struct {  // wastes space
	A bool    // 1 byte, then 7 bytes padding
	B float64 // 8 bytes (needs 8-byte alignment)
	C bool    // 1 byte, then 7 bytes padding
	D float64 // 8 bytes
} // total: 1+7+8+1+7+8 = 32 bytes

type GoodLayout struct {  // efficient
	B float64 // 8 bytes
	D float64 // 8 bytes
	A bool    // 1 byte
	C bool    // 1 byte, then 6 bytes padding (struct must be multiple of 8)
} // total: 8+8+1+1+6 = 24 bytes

func memoryLayout() {
	fmt.Println("=== 1. Struct Memory Layout ===")

	fmt.Printf("BadLayout:  %d bytes\n", unsafe.Sizeof(BadLayout{}))
	fmt.Printf("GoodLayout: %d bytes (same data, better ordering)\n", unsafe.Sizeof(GoodLayout{}))

	// Visualize field offsets
	b := BadLayout{}
	fmt.Printf("\nBadLayout field offsets:\n")
	fmt.Printf("  A: offset=%d, size=%d\n", unsafe.Offsetof(b.A), unsafe.Sizeof(b.A))
	fmt.Printf("  B: offset=%d, size=%d\n", unsafe.Offsetof(b.B), unsafe.Sizeof(b.B))
	fmt.Printf("  C: offset=%d, size=%d\n", unsafe.Offsetof(b.C), unsafe.Sizeof(b.C))
	fmt.Printf("  D: offset=%d, size=%d\n", unsafe.Offsetof(b.D), unsafe.Sizeof(b.D))
}

// ─── 2. STRUCT DECLARATION AND INITIALIZATION ─────────────────────────────────

type Point struct {
	X, Y float64  // two fields, same type, same line
}

type Rectangle struct {
	TopLeft     Point
	BottomRight Point
	Label       string
}

func structInit() {
	fmt.Println("\n=== 2. Struct Initialization ===")

	// Positional initialization (order-dependent, fragile)
	p1 := Point{1.0, 2.0}

	// Named field initialization (preferred — order-independent)
	p2 := Point{X: 3.0, Y: 4.0}

	// Partial initialization — unspecified fields get zero values
	p3 := Point{X: 5.0}  // Y = 0.0

	fmt.Printf("p1: %+v\n", p1)  // %+v includes field names
	fmt.Printf("p2: %+v\n", p2)
	fmt.Printf("p3: %+v\n", p3)

	// Nested struct
	rect := Rectangle{
		TopLeft:     Point{0, 0},
		BottomRight: Point{100, 50},
		Label:       "window",
	}
	fmt.Printf("rect: %+v\n", rect)

	// Anonymous struct (one-off, no name)
	config := struct {
		Host string
		Port int
	}{
		Host: "localhost",
		Port: 8080,
	}
	fmt.Printf("config: %+v\n", config)

	// Struct comparison: structs are comparable if all fields are comparable
	q1 := Point{1, 2}
	q2 := Point{1, 2}
	fmt.Printf("p1 == p2: %v\n", p1 == p2)
	fmt.Printf("q1 == q2: %v\n", q1 == q2)
}

// ─── 3. STRUCT TAGS ───────────────────────────────────────────────────────────
//
// Tags are string metadata attached to struct fields.
// They are read at RUNTIME using the reflect package.
// Common uses: JSON encoding, database mapping, validation.

type User struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	Email     string `json:"email,omitempty"` // omitted if empty
	Password  string `json:"-"`               // always omitted from JSON
	CreatedAt string `json:"created_at"`
}

func structTags() {
	fmt.Println("\n=== 3. Struct Tags ===")

	u := User{
		ID:        1,
		Name:      "Alice",
		Email:     "alice@example.com",
		Password:  "secret123",
		CreatedAt: "2024-01-01",
	}

	// Marshal to JSON — tags control the output
	data, _ := json.Marshal(u)
	fmt.Printf("JSON: %s\n", data)
	// Notice: password is omitted, fields use snake_case names from tags

	// Unmarshal from JSON
	jsonData := `{"id":2,"name":"Bob","created_at":"2024-06-01"}`
	var u2 User
	json.Unmarshal([]byte(jsonData), &u2)
	fmt.Printf("Unmarshaled: %+v\n", u2)
	// email is empty (omitempty doesn't affect unmarshaling)
}

// ─── 4. METHODS ───────────────────────────────────────────────────────────────
//
// Methods are functions with a receiver — they're associated with a type.
// Two kinds of receivers:
//
//   Value receiver: func (p Point) Method()
//     - p is a COPY of the original
//     - changes to p don't affect the original
//     - can call on values AND pointers
//
//   Pointer receiver: func (p *Point) Method()
//     - p is a pointer to the original
//     - changes to *p modify the original
//     - can call on values (auto-addressed) AND pointers
//
// When to use pointer receiver:
//   1. The method needs to modify the receiver
//   2. The receiver is large (avoid copying)
//   3. Consistency: if some methods use pointer receivers, all should
//
// When to use value receiver:
//   1. The type is small and copyable (like Point)
//   2. The type should be immutable from methods

func (p Point) Distance(other Point) float64 {
	dx := p.X - other.X
	dy := p.Y - other.Y
	return dx*dx + dy*dy  // squared distance (no math.Sqrt import needed here)
}

func (p Point) String() string {
	return fmt.Sprintf("(%.1f, %.1f)", p.X, p.Y)
}

type Counter struct {
	value int
}

func (c *Counter) Increment() { c.value++ }      // pointer receiver — modifies c
func (c *Counter) Add(n int) { c.value += n }
func (c Counter) Value() int { return c.value }   // value receiver — read-only

func methods() {
	fmt.Println("\n=== 4. Methods ===")

	p1 := Point{0, 0}
	p2 := Point{3, 4}
	fmt.Printf("Distance² from %v to %v: %.1f\n", p1, p2, p1.Distance(p2))
	// p1.String() is auto-called by fmt when using %v or %s

	c := Counter{}
	c.Increment()
	c.Increment()
	c.Add(10)
	fmt.Printf("Counter value: %d\n", c.Value())

	// Go auto-addresses for pointer receiver calls on value variables
	cp := &Counter{}   // explicit pointer
	cp.Increment()     // pointer receiver — natural
	c.Increment()      // auto-address: Go reads this as (&c).Increment()

	// Calling value receiver on pointer also works:
	var cp2 *Counter = &Counter{value: 5}
	fmt.Printf("value via pointer: %d\n", cp2.Value())  // Go dereferences: (*cp2).Value()
}

// ─── 5. STRUCT EMBEDDING ──────────────────────────────────────────────────────
//
// Embedding is Go's composition mechanism. It is NOT inheritance.
//
// When you embed a type T in a struct S:
//   - T's fields and methods are PROMOTED to S
//   - You can call them directly on S values
//   - But S and T are different types — no polymorphism
//   - S does NOT "extend" T in the OOP sense
//
// Think of embedding as "delegation" — S forwards calls to the embedded T.
// You still have direct access to the embedded field by its type name.

type Animal struct {
	Name string
}

func (a Animal) Speak() string {
	return a.Name + " makes a sound"
}

type Dog struct {
	Animal        // embedded — no field name, type name IS the field name
	Breed string
}

// Dog can override Speak:
func (d Dog) Speak() string {
	return d.Name + " barks!"  // d.Name promoted from Animal
}

type GuardDog struct {
	Dog           // embed Dog (which embeds Animal)
	GuardLevel int
}

func embedding() {
	fmt.Println("\n=== 5. Struct Embedding ===")

	d := Dog{
		Animal: Animal{Name: "Rex"},
		Breed:  "Labrador",
	}

	// Promoted field access
	fmt.Println("Name via promotion:", d.Name)       // d.Animal.Name
	fmt.Println("Name via full path:", d.Animal.Name) // same thing

	// Promoted method — Dog has its own Speak, overrides Animal's
	fmt.Println(d.Speak())         // Dog.Speak()
	fmt.Println(d.Animal.Speak())  // Animal.Speak() — still accessible

	// Multi-level embedding
	gd := GuardDog{
		Dog:        d,
		GuardLevel: 5,
	}
	fmt.Println("Guard dog name:", gd.Name)       // promoted through two levels
	fmt.Println("Guard dog speaks:", gd.Speak())  // Dog.Speak()
	fmt.Printf("Guard level: %d\n", gd.GuardLevel)

	// Embedding an interface is also valid — see interfaces chapter
}

// ─── 6. CONSTRUCTOR FUNCTIONS (Go IDIOM) ─────────────────────────────────────
//
// Go has no constructors. Instead, use `NewXxx` factory functions.
// Return *T (pointer) when:
//   - The struct is large
//   - The struct has methods with pointer receivers
//   - You want to be able to return nil (optional value)
// Return T (value) when:
//   - The struct is small and all methods use value receivers

type Config struct {
	Host     string
	Port     int
	MaxConns int
	Debug    bool
}

// Option pattern — functional options (common Go idiom)
type Option func(*Config)

func WithHost(h string) Option     { return func(c *Config) { c.Host = h } }
func WithPort(p int) Option        { return func(c *Config) { c.Port = p } }
func WithMaxConns(n int) Option    { return func(c *Config) { c.MaxConns = n } }
func WithDebug(d bool) Option      { return func(c *Config) { c.Debug = d } }

func NewConfig(opts ...Option) *Config {
	// defaults
	c := &Config{
		Host:     "localhost",
		Port:     8080,
		MaxConns: 100,
		Debug:    false,
	}
	for _, opt := range opts {
		opt(c)
	}
	return c
}

func constructors() {
	fmt.Println("\n=== 6. Constructor Functions and Option Pattern ===")

	// Default config
	c1 := NewConfig()
	fmt.Printf("default: %+v\n", c1)

	// Custom config — only override what you need
	c2 := NewConfig(
		WithHost("0.0.0.0"),
		WithPort(9000),
		WithDebug(true),
	)
	fmt.Printf("custom: %+v\n", c2)
}

func main() {
	memoryLayout()
	structInit()
	structTags()
	methods()
	embedding()
	constructors()
}

/*
THOUGHT QUESTIONS:

1. Why does the order of fields in a struct affect its memory size?
   How would you reorder fields to minimize padding?

2. When should you use a pointer receiver vs a value receiver for methods?
   What is the "consistency" rule?

3. Is struct embedding the same as inheritance? What can you do with
   inheritance in Java/Python that you cannot do with Go embedding?

4. What is the functional options pattern? What problem does it solve
   compared to a struct literal or variadic arguments?

5. Why does Go not have a `constructor` keyword? What pattern does it
   use instead?

EXERCISES:

1. Create a `Stack[int]` struct with Push, Pop, Peek, and Len methods.
   Use a pointer receiver for all mutating methods. Verify behavior.

2. Design a `Logger` struct with embedded `io.Writer`. Write a method
   `Log(level, message string)` that formats and writes to the writer.

3. Create a `Shape` interface with `Area() float64` and `Perimeter() float64`.
   Implement it for `Circle`, `Rectangle`, and `Triangle`. Then create a
   `Polygon` struct that embeds multiple shapes and aggregates their areas.

4. Write a function `sizeof(v interface{}) string` that uses reflect to
   print a struct's field names, types, sizes, and offsets.
*/
