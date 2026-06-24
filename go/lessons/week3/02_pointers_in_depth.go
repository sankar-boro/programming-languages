/*
WEEK 3 — DAY 2: Pointers In Depth
====================================
Topic: Pointers, pointer arithmetic (lack of), nil safety, and when to use them.

Key ideas:
  - Pointers allow mutation without copying and express ownership
  - Go pointers cannot be arithmetic-operated (safer than C)
  - nil pointer dereference causes a panic (runtime, not compile-time)
  - unsafe.Pointer and uintptr exist for low-level operations
  - Pointer receivers allow methods to mutate their receivers
*/

package main

import (
	"fmt"
	"unsafe"
)

// ─── 1. POINTER FUNDAMENTALS REVISITED ────────────────────────────────────────

func fundamentals() {
	fmt.Println("=== 1. Pointer Fundamentals ===")

	// & = address-of; * = dereference
	x := 42
	p := &x             // p is *int — "pointer to int"
	fmt.Printf("x = %d, &x = %p\n", x, &x)
	fmt.Printf("p = %p, *p = %d\n", p, *p)

	// Modifying through pointer
	*p = 100
	fmt.Printf("After *p=100: x = %d\n", x)

	// Pointer to pointer
	pp := &p           // pp is **int
	**pp = 200
	fmt.Printf("After **pp=200: x = %d\n", x)

	// Pointers to different types
	s := "hello"
	sp := &s
	fmt.Printf("*sp = %q\n", *sp)
	*sp = "world"
	fmt.Printf("After *sp=world: s = %q\n", s)
}

// ─── 2. WHY PASS POINTERS? (mutation vs copy) ────────────────────────────────
//
// Go passes arguments BY VALUE. If you want a function to modify its argument,
// pass a pointer. This makes mutation EXPLICIT — you can always see at the
// call site that a function might modify the argument (it's a pointer).

type Player struct {
	Name  string
	Score int
	Level int
}

// Takes a copy — cannot modify original
func awardBonusValue(p Player, bonus int) {
	p.Score += bonus  // modifies the copy, not the original
}

// Takes a pointer — CAN modify original
func awardBonusPointer(p *Player, bonus int) {
	p.Score += bonus  // modifies the original
}

func mutationVsCopy() {
	fmt.Println("\n=== 2. Mutation vs Copy ===")

	p := Player{Name: "Alice", Score: 100, Level: 5}

	awardBonusValue(p, 50)  // passes a copy
	fmt.Printf("After value call: score=%d (unchanged)\n", p.Score)

	awardBonusPointer(&p, 50)  // passes pointer to p
	fmt.Printf("After pointer call: score=%d (modified)\n", p.Score)

	// Also useful: avoiding large copy overhead
	// A Player struct with 10 fields would copy all 10 fields on every call
	// A *Player is just 8 bytes (pointer size) regardless of struct size
}

// ─── 3. NIL POINTERS ──────────────────────────────────────────────────────────
//
// The zero value of a pointer is nil — it points to nothing.
// Dereferencing a nil pointer causes a PANIC at runtime.
// Always check for nil before dereferencing an unknown pointer.

func findPlayer(players []Player, name string) *Player {
	for i := range players {
		if players[i].Name == name {
			return &players[i]
		}
	}
	return nil  // not found
}

func nilPointers() {
	fmt.Println("\n=== 3. Nil Pointers ===")

	players := []Player{
		{Name: "Alice", Score: 100},
		{Name: "Bob", Score: 200},
	}

	// Safe nil check
	if p := findPlayer(players, "Alice"); p != nil {
		fmt.Printf("Found: %s (score=%d)\n", p.Name, p.Score)
	}

	if p := findPlayer(players, "Charlie"); p == nil {
		fmt.Println("Charlie not found")
	}

	// Nil pointer dereference would panic:
	// var nilPlayer *Player
	// fmt.Println(nilPlayer.Score)  // PANIC

	// Demonstrate safe nil method call (with nil check in the method)
	var nilPlayer *Player
	fmt.Println("name of nil player:", nilPlayer.SafeName())
}

func (p *Player) SafeName() string {
	if p == nil {
		return "(nil)"
	}
	return p.Name
}

// ─── 4. POINTERS AND INTERFACES ───────────────────────────────────────────────
//
// Interfaces are already reference types (they contain a pointer internally).
// But pointer vs value receiver determines which values SATISFY an interface.
//
// Method set rules (critical):
//   Value type T:    can call methods with VALUE receiver on T
//   Pointer type *T: can call methods with BOTH value AND pointer receiver
//
// If an interface requires a pointer receiver method, only *T satisfies it.

type Stringer interface {
	String() string
}

type Box struct {
	Width, Height int
}

// Value receiver — both Box and *Box satisfy Stringer
func (b Box) String() string {
	return fmt.Sprintf("Box(%dx%d)", b.Width, b.Height)
}

type Counter struct {
	n int
}

// Pointer receiver — only *Counter satisfies Incrementer
type Incrementer interface {
	Increment()
}

func (c *Counter) Increment() { c.n++ }

func pointersAndInterfaces() {
	fmt.Println("\n=== 4. Pointers and Interfaces ===")

	// Value receiver — both T and *T satisfy the interface
	var s1 Stringer = Box{10, 20}    // OK — value receiver
	var s2 Stringer = &Box{30, 40}   // also OK — *T can call value receiver methods
	fmt.Println(s1.String())
	fmt.Println(s2.String())

	// Pointer receiver — only *T satisfies the interface
	// var inc Incrementer = Counter{} // COMPILE ERROR: Counter does not implement Incrementer
	var inc Incrementer = &Counter{}  // OK — pointer
	inc.Increment()
	fmt.Printf("Counter after Increment: %d\n", inc.(*Counter).n)
}

// ─── 5. unsafe.Pointer and uintptr ────────────────────────────────────────────
//
// The `unsafe` package breaks Go's type safety for low-level operations.
// Use only when absolutely necessary (e.g., C interop, performance-critical code).
//
// unsafe.Pointer — a special pointer type that can hold any pointer
//   - Can convert between pointer types
//   - The GC knows about it (won't move the pointed-to data)
//
// uintptr — an integer that holds an address
//   - Can be arithmetic-operated (pointer arithmetic!)
//   - The GC does NOT know about it — cannot hold as a "reference"
//   - Converting unsafe.Pointer → uintptr → unsafe.Pointer is dangerous
//     because the GC might move the object between the two conversions

type Pair struct {
	First  int32
	Second int32
}

func unsafePointers() {
	fmt.Println("\n=== 5. unsafe.Pointer ===")

	// Reading struct fields by offset
	pair := Pair{First: 100, Second: 200}
	pairPtr := unsafe.Pointer(&pair)

	// Access Second field directly via pointer arithmetic
	secondPtr := (*int32)(unsafe.Pointer(uintptr(pairPtr) + unsafe.Offsetof(pair.Second)))
	fmt.Printf("pair.Second via unsafe: %d\n", *secondPtr)

	// Reinterpret bits: read a float64 as uint64
	f := 3.14
	u := *(*uint64)(unsafe.Pointer(&f))
	fmt.Printf("3.14 as uint64 bits: %016x\n", u)
	// This is how you inspect IEEE 754 representation

	// Sizeof and Alignof
	fmt.Printf("\nunsafe.Sizeof(int32):  %d\n", unsafe.Sizeof(int32(0)))
	fmt.Printf("unsafe.Sizeof(Pair):   %d\n", unsafe.Sizeof(Pair{}))
	fmt.Printf("unsafe.Alignof(int32): %d\n", unsafe.Alignof(int32(0)))
	fmt.Printf("unsafe.Alignof(Pair):  %d\n", unsafe.Alignof(Pair{}))
}

// ─── 6. COMMON POINTER PATTERNS ───────────────────────────────────────────────

// Optional values — pointer as nullable
type Config struct {
	Timeout *int  // nil means "not set, use default"
	Debug   *bool // nil means "not set, use default"
}

func withTimeout(t int) *int     { return &t }
func withDebug(d bool) *bool     { return &d }

// Pointer for optional fields in function signatures
type SearchOptions struct {
	Limit  *int
	Offset *int
}

func search(query string, opts *SearchOptions) []string {
	limit := 10  // default
	offset := 0  // default
	if opts != nil {
		if opts.Limit != nil { limit = *opts.Limit }
		if opts.Offset != nil { offset = *opts.Offset }
	}
	_ = limit; _ = offset
	return []string{fmt.Sprintf("results for '%s'", query)}
}

func pointerPatterns() {
	fmt.Println("\n=== 6. Common Pointer Patterns ===")

	// Nullable config
	t := 30
	d := true
	cfg := Config{Timeout: &t, Debug: &d}
	fmt.Printf("timeout: %d, debug: %v\n", *cfg.Timeout, *cfg.Debug)

	cfgDefault := Config{}  // all nil — use defaults
	timeout := 10           // default
	if cfgDefault.Timeout != nil {
		timeout = *cfgDefault.Timeout
	}
	fmt.Printf("default timeout: %d\n", timeout)

	// Optional search options
	results := search("golang", nil)  // no options
	fmt.Println(results)

	limit := 5
	results = search("golang", &SearchOptions{Limit: &limit})
	fmt.Println(results)
}

func main() {
	fundamentals()
	mutationVsCopy()
	nilPointers()
	pointersAndInterfaces()
	unsafePointers()
	pointerPatterns()
}

/*
THOUGHT QUESTIONS:

1. Go doesn't allow pointer arithmetic (p++ doesn't work). What safety
   benefit does this provide? What low-level capabilities does it prevent?

2. Explain the method set rule: "A value T cannot call pointer receiver methods,
   but *T can call value receiver methods." Why is this rule designed this way?

3. When would you use a pointer to a primitive (like *int or *bool) in a struct?
   Give a real-world example.

4. What is the difference between unsafe.Pointer and uintptr? Why is converting
   between them dangerous?

5. Why is nil pointer dereference a runtime panic in Go rather than a compile-time
   error? What would it take to make it a compile-time error?

EXERCISES:

1. Write a generic swap function: `func swap(a, b *int)` that swaps the values
   at the two pointers. Then write a generic version using any pointer type.

2. Implement a linked list using pointers:
   type Node struct { Value int; Next *Node }
   Write Insert, Delete, and Print methods.

3. Write a function `deepCopy(p *Player) *Player` that creates a completely
   independent copy of a Player (so modifying the copy doesn't affect the original).

4. Using unsafe.Pointer, write a function that converts a []byte to a string
   without allocating new memory. Explain why this is safe (or not).
*/
