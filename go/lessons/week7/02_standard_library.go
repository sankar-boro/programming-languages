/*
WEEK 7 — DAY 2: Standard Library Deep Dive
===========================================
Topic: Go's rich standard library — the packages you use every day.

Key ideas:
  - Go ships a comprehensive stdlib — "batteries included"
  - fmt, os, io, bufio, strings, strconv, time, encoding/json
  - net/http for HTTP servers and clients
  - testing for unit tests
  - context for cancellation and deadlines
*/

package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"
	"unicode"
)

// ─── 1. fmt — FORMATTED I/O ───────────────────────────────────────────────────
//
// fmt is one of the most-used packages. Key functions:
//
//   Printing:
//     fmt.Print, fmt.Println, fmt.Printf     — stdout
//     fmt.Fprint, fmt.Fprintln, fmt.Fprintf  — to an io.Writer
//     fmt.Sprint, fmt.Sprintln, fmt.Sprintf  — to a string
//
//   Scanning:
//     fmt.Scan, fmt.Scanln, fmt.Scanf   — read from stdin
//     fmt.Fscan, fmt.Fscanln            — read from io.Reader
//     fmt.Sscan, fmt.Sscanf             — parse from string
//
//   Errors:
//     fmt.Errorf("msg: %w", err)  — create formatted error (with wrapping)
//
// Format verbs:
//   %v   default format        %+v  struct with field names
//   %T   type                  %#v  Go syntax representation
//   %d   integer (decimal)     %x   hex   %b binary   %o octal
//   %f   float (default)       %e   scientific  %.2f  2 decimal places
//   %s   string                %q   quoted string     %c rune as char
//   %p   pointer address       %t   bool

func fmtPackage() {
	fmt.Println("=== 1. fmt Package ===")

	// Basic formatting
	x := 42
	f := 3.14159
	s := "hello"

	fmt.Printf("decimal: %d, hex: %x, binary: %b, octal: %o\n", x, x, x, x)
	fmt.Printf("float default: %f\n", f)
	fmt.Printf("float precision: %.2f\n", f)
	fmt.Printf("float scientific: %e\n", f)
	fmt.Printf("string: %s, quoted: %q\n", s, s)

	type Point struct{ X, Y int }
	p := Point{3, 4}
	fmt.Printf("default: %v, with fields: %+v, Go syntax: %#v\n", p, p, p)
	fmt.Printf("type: %T\n", p)

	// Sprintf — format to string
	msg := fmt.Sprintf("User %d has %d credits", 42, 100)
	fmt.Println(msg)

	// Fprintf — format to writer
	var buf bytes.Buffer
	fmt.Fprintf(&buf, "buffered: %d", x)
	fmt.Println(buf.String())

	// Sscanf — parse from string
	var name string
	var age int
	fmt.Sscanf("Alice 30", "%s %d", &name, &age)
	fmt.Printf("parsed: name=%s, age=%d\n", name, age)
}

// ─── 2. strings — STRING MANIPULATION ────────────────────────────────────────
//
// The strings package provides functions for string operations.
// Remember: strings are IMMUTABLE — every operation returns a new string.

func stringsPackage() {
	fmt.Println("\n=== 2. strings Package ===")

	s := "  Hello, World!  "

	fmt.Println("TrimSpace:", strings.TrimSpace(s))
	fmt.Println("ToUpper:", strings.ToUpper(s))
	fmt.Println("ToLower:", strings.ToLower(s))
	fmt.Println("Contains:", strings.Contains(s, "World"))
	fmt.Println("HasPrefix:", strings.HasPrefix(strings.TrimSpace(s), "Hello"))
	fmt.Println("HasSuffix:", strings.HasSuffix(strings.TrimSpace(s), "!"))
	fmt.Println("Count:", strings.Count(s, "l"))  // count occurrences
	fmt.Println("Index:", strings.Index(s, "World"))  // -1 if not found
	fmt.Println("Replace:", strings.Replace(s, "l", "L", 2))  // replace first 2
	fmt.Println("ReplaceAll:", strings.ReplaceAll(s, "l", "L"))

	// Split and Join
	csv := "a,b,c,d,e"
	parts := strings.Split(csv, ",")
	fmt.Println("Split:", parts)
	fmt.Println("Join:", strings.Join(parts, " | "))

	// Fields — split by whitespace
	words := strings.Fields("  hello   world   go  ")
	fmt.Println("Fields:", words)

	// Trim — remove characters from both ends
	fmt.Println("Trim(./path/./,.):", strings.Trim("./path/./", "./"))
	fmt.Println("TrimLeft:", strings.TrimLeft("xxxhello", "x"))
	fmt.Println("TrimRight:", strings.TrimRight("helloyyy", "y"))

	// strings.Builder — efficient string construction
	var sb strings.Builder
	words2 := []string{"Go", "is", "fast", "and", "safe"}
	for i, w := range words2 {
		if i > 0 { sb.WriteString(" ") }
		sb.WriteString(w)
	}
	fmt.Println("Builder:", sb.String())

	// strings.NewReader — use a string as io.Reader
	r := strings.NewReader("hello world")
	buf := make([]byte, 5)
	n, _ := r.Read(buf)
	fmt.Printf("Read 5 bytes: %q\n", string(buf[:n]))

	// Map — transform each rune
	encode := strings.Map(func(r rune) rune {
		if unicode.IsLetter(r) {
			return r + 1  // shift each letter by 1
		}
		return r
	}, "Hello World!")
	fmt.Println("Map:", encode)
}

// ─── 3. strconv — STRING CONVERSIONS ──────────────────────────────────────────
//
// Converting between strings and other types.
// NEVER use fmt.Sprintf for simple number→string conversion — use strconv.

func strconvPackage() {
	fmt.Println("\n=== 3. strconv Package ===")

	// int → string
	n := 42
	s := strconv.Itoa(n)  // int to ASCII
	fmt.Printf("Itoa(%d) = %q\n", n, s)

	// string → int
	i, err := strconv.Atoi("123")
	fmt.Printf("Atoi(\"123\") = %d, err=%v\n", i, err)

	_, err = strconv.Atoi("not a number")
	fmt.Printf("Atoi(\"not a number\"): err=%v\n", err)

	// ParseInt with base
	hex, _ := strconv.ParseInt("FF", 16, 64)  // base 16, 64-bit
	fmt.Printf("ParseInt(\"FF\", 16) = %d\n", hex)

	bin, _ := strconv.ParseInt("1010", 2, 64)  // base 2
	fmt.Printf("ParseInt(\"1010\", 2) = %d\n", bin)

	// float64 conversions
	f, _ := strconv.ParseFloat("3.14159", 64)
	fmt.Printf("ParseFloat = %f\n", f)
	fs := strconv.FormatFloat(f, 'f', 2, 64)  // format, precision, bit size
	fmt.Printf("FormatFloat (2 decimal) = %q\n", fs)

	// bool
	b, _ := strconv.ParseBool("true")
	fmt.Printf("ParseBool(\"true\") = %v\n", b)
	fmt.Printf("FormatBool(false) = %q\n", strconv.FormatBool(false))

	// AppendInt — append number to existing byte slice (no allocation)
	buf := []byte("count: ")
	buf = strconv.AppendInt(buf, 42, 10)
	fmt.Printf("AppendInt: %q\n", string(buf))
}

// ─── 4. os — OPERATING SYSTEM INTERFACE ──────────────────────────────────────

func osPackage() {
	fmt.Println("\n=== 4. os Package ===")

	// Environment variables
	home := os.Getenv("HOME")
	fmt.Printf("HOME: %s\n", home)

	os.Setenv("MY_VAR", "hello")
	fmt.Printf("MY_VAR: %s\n", os.Getenv("MY_VAR"))

	// Command-line arguments
	fmt.Printf("os.Args[0] (binary): %s\n", os.Args[0])

	// Working directory
	wd, _ := os.Getwd()
	fmt.Printf("Working dir: %s\n", wd)

	// File operations
	tmpFile := "/tmp/go_stdlib_demo.txt"
	content := []byte("Hello from Go stdlib!\nLine 2\nLine 3\n")

	// Write file
	err := os.WriteFile(tmpFile, content, 0644)
	if err != nil {
		fmt.Println("WriteFile error:", err)
		return
	}
	fmt.Println("Wrote file:", tmpFile)

	// Read file
	data, err := os.ReadFile(tmpFile)
	if err == nil {
		fmt.Printf("ReadFile: %q\n", string(data))
	}

	// File info
	info, _ := os.Stat(tmpFile)
	if info != nil {
		fmt.Printf("Size: %d bytes, Mode: %s\n", info.Size(), info.Mode())
	}

	// Cleanup
	os.Remove(tmpFile)

	// os.Exit — terminates immediately, defers do NOT run
	// os.Exit(0)  — don't call this in demos
}

// ─── 5. bufio — BUFFERED I/O ─────────────────────────────────────────────────
//
// bufio wraps an io.Reader or io.Writer with a buffer.
// This reduces the number of system calls for reading/writing.
// Essential for reading line-by-line from files or stdin.

func bufioPackage() {
	fmt.Println("\n=== 5. bufio Package ===")

	// Scanner — read line-by-line
	input := "line one\nline two\nline three\n"
	scanner := bufio.NewScanner(strings.NewReader(input))
	scanner.Split(bufio.ScanLines)  // default — split on newlines

	lineNum := 1
	for scanner.Scan() {
		fmt.Printf("  %d: %s\n", lineNum, scanner.Text())
		lineNum++
	}
	if err := scanner.Err(); err != nil {
		fmt.Println("scan error:", err)
	}

	// Scanner can split on words or custom tokens
	wordScanner := bufio.NewScanner(strings.NewReader("hello world go"))
	wordScanner.Split(bufio.ScanWords)
	for wordScanner.Scan() {
		fmt.Printf("  word: %q\n", wordScanner.Text())
	}

	// Writer — buffered writes
	var buf bytes.Buffer
	w := bufio.NewWriter(&buf)
	fmt.Fprint(w, "first ")
	fmt.Fprint(w, "second ")
	// Data may not be in buf yet (buffered)
	w.Flush()  // flush buffer to underlying writer
	fmt.Printf("buffered write: %q\n", buf.String())

	// Reader — buffered reads (efficient for many small reads)
	r := bufio.NewReader(strings.NewReader("hello\nworld"))
	line, _ := r.ReadString('\n')  // read until delimiter
	fmt.Printf("ReadString: %q\n", line)
}

// ─── 6. time — DATE AND TIME ─────────────────────────────────────────────────

func timePackage() {
	fmt.Println("\n=== 6. time Package ===")

	now := time.Now()
	fmt.Println("Now:", now)
	fmt.Println("UTC:", now.UTC())
	fmt.Println("Unix timestamp:", now.Unix())
	fmt.Println("UnixNano:", now.UnixNano())

	// Format — Go uses a reference time: Mon Jan 2 15:04:05 MST 2006
	fmt.Println("Formatted:", now.Format("2006-01-02 15:04:05"))
	fmt.Println("RFC3339:", now.Format(time.RFC3339))
	fmt.Println("Kitchen:", now.Format(time.Kitchen))

	// Parsing
	t, _ := time.Parse("2006-01-02", "2024-01-15")
	fmt.Println("Parsed:", t)

	// Arithmetic
	future := now.Add(24 * time.Hour)
	past := now.Add(-7 * 24 * time.Hour)
	diff := future.Sub(past)
	fmt.Printf("Future: %s, Past: %s, Diff: %v\n",
		future.Format("2006-01-02"), past.Format("2006-01-02"), diff)

	// Duration
	d := 2*time.Hour + 30*time.Minute + 15*time.Second
	fmt.Println("Duration:", d)
	fmt.Println("Hours:", d.Hours())

	// Timer and Ticker (non-blocking demo)
	timer := time.NewTimer(1 * time.Millisecond)
	<-timer.C  // wait for timer
	fmt.Println("Timer fired!")

	// time.Sleep
	start := time.Now()
	time.Sleep(5 * time.Millisecond)
	fmt.Printf("Sleep took: %v\n", time.Since(start))
}

// ─── 7. encoding/json ────────────────────────────────────────────────────────
//
// encoding/json is essential for working with JSON.
// It uses struct tags to control field names and behavior.

type Address struct {
	Street string `json:"street"`
	City   string `json:"city"`
	Zip    string `json:"zip"`
}

type Person struct {
	Name      string    `json:"name"`
	Age       int       `json:"age"`
	Email     string    `json:"email,omitempty"` // omit if empty
	Address   *Address  `json:"address,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	Password  string    `json:"-"`                // never marshal
}

func jsonPackage() {
	fmt.Println("\n=== 7. encoding/json ===")

	p := Person{
		Name:      "Alice",
		Age:       30,
		Email:     "alice@example.com",
		Address:   &Address{"123 Main St", "Springfield", "12345"},
		CreatedAt: time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC),
		Password:  "secret",
	}

	// Marshal: Go struct → JSON bytes
	data, err := json.Marshal(p)
	if err != nil {
		fmt.Println("marshal error:", err)
		return
	}
	fmt.Println("JSON:", string(data))

	// MarshalIndent: pretty-printed
	pretty, _ := json.MarshalIndent(p, "", "  ")
	fmt.Println("Pretty JSON:")
	fmt.Println(string(pretty))

	// Unmarshal: JSON bytes → Go struct
	jsonStr := `{"name":"Bob","age":25,"address":{"city":"NYC","street":"456 Oak Ave","zip":"10001"}}`
	var p2 Person
	err = json.Unmarshal([]byte(jsonStr), &p2)
	if err == nil {
		fmt.Printf("Unmarshaled: %+v\n", p2)
	}

	// Generic JSON: decode into map[string]interface{}
	var data2 map[string]interface{}
	json.Unmarshal([]byte(jsonStr), &data2)
	fmt.Printf("Generic JSON: name=%v, age=%v\n", data2["name"], data2["age"])

	// Streaming JSON with Encoder/Decoder
	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)
	enc.Encode(p)   // writes JSON + newline
	fmt.Println("Encoded:", strings.TrimSpace(buf.String()))
}

// ─── 8. sort PACKAGE ─────────────────────────────────────────────────────────

func sortPackage() {
	fmt.Println("\n=== 8. sort Package ===")

	// Sort primitives
	ints := []int{5, 2, 8, 1, 9, 3}
	sort.Ints(ints)
	fmt.Println("Sorted ints:", ints)

	strs := []string{"banana", "apple", "cherry"}
	sort.Strings(strs)
	fmt.Println("Sorted strings:", strs)

	// Sort structs
	type Employee struct{ Name string; Salary int }
	emps := []Employee{
		{"Charlie", 70000}, {"Alice", 90000}, {"Bob", 80000},
	}
	sort.Slice(emps, func(i, j int) bool {
		return emps[i].Salary > emps[j].Salary  // descending
	})
	for _, e := range emps {
		fmt.Printf("  %s: $%d\n", e.Name, e.Salary)
	}

	// Binary search
	data := []int{1, 3, 5, 7, 9, 11, 13}
	idx := sort.SearchInts(data, 7)
	fmt.Printf("Binary search for 7: index=%d, value=%d\n", idx, data[idx])

	// sort.Search — generic binary search
	target := 9
	pos := sort.Search(len(data), func(i int) bool { return data[i] >= target })
	fmt.Printf("sort.Search for %d: found at index %d\n", target, pos)
}

func main() {
	fmtPackage()
	stringsPackage()
	strconvPackage()
	osPackage()
	bufioPackage()
	timePackage()
	jsonPackage()
	sortPackage()

	// Verify all readers work
	_ = io.EOF  // io.EOF is the standard "end of file" error sentinel
}

/*
THOUGHT QUESTIONS:

1. Why is strings.Builder preferred over repeated string concatenation (+)?
   What allocation difference does it make?

2. strconv.Itoa(42) vs fmt.Sprintf("%d", 42) — they produce the same result.
   When should you use each?

3. bufio.Scanner reads line-by-line by default. What happens if a line is
   longer than the default buffer size (64 KB)?

4. Go's time.Format uses a reference time (Jan 2 15:04:05 MST 2006) instead
   of % codes like C's strftime. What are the advantages and disadvantages?

5. json.Marshal ignores fields tagged with `json:"-"`. What is this useful for?
   What other fields does it skip by default?

EXERCISES:

1. Write a function `parseCSV(r io.Reader) ([][]string, error)` that reads
   a CSV file line by line, splitting each line on commas.

2. Write a function `formatDuration(d time.Duration) string` that returns
   a human-friendly string like "2 hours, 30 minutes, 15 seconds".

3. Write a `JSONMerge(a, b []byte) ([]byte, error)` function that merges
   two JSON objects (b's keys override a's).

4. Implement a `SortBy[T any, K cmp.Ordered](items []T, key func(T) K) []T`
   using the sort package — should produce a sorted copy without modifying original.
*/
