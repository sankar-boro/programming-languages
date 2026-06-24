/*
WEEK 7 — DAY 1: Modules, Packages, and the Build System
=========================================================
Topic: How Go organizes code — modules, packages, imports, visibility, and go.mod.

Key ideas:
  - A package is the unit of compilation and encapsulation
  - A module is a collection of packages with a shared go.mod
  - Visibility is controlled by capitalization (exported vs unexported)
  - init() functions run package initialization in dependency order
  - go.mod defines module identity, Go version, and dependencies
  - The import path must match the directory structure
*/

package main

import (
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"time"
)

// ─── 1. PACKAGES — UNIT OF ORGANIZATION ──────────────────────────────────────
//
// Every Go source file belongs to exactly ONE package (first line: package X).
// A directory may contain only one package (except test files: package X_test).
//
// Naming conventions:
//   - Package names are lowercase, short, no underscores
//   - The package name is the LAST element of the import path
//   - import "net/http" → you use http.Handler, http.Get, etc.
//
// A Go program has ONE package named "main" with ONE func main().
// All other packages are library packages (no main function).

// Everything in THIS file is in package main.
// Unexported identifiers (lowercase) are visible within this file only (package).
// Exported identifiers (uppercase) are visible to ALL packages.

type person struct {   // unexported — only accessible within this package
	name string
	age  int
}

type User struct {     // exported — accessible from other packages
	Name  string
	Email string
	Age   int
}

func newPerson(name string, age int) *person {
	return &person{name: name, age: age}  // unexported constructor
}

func CreateUser(name, email string, age int) *User {
	return &User{Name: name, Email: email, Age: age}  // exported constructor
}

func packages() {
	fmt.Println("=== 1. Packages ===")

	p := newPerson("Alice", 30)   // accessible in same package
	u := CreateUser("Bob", "bob@example.com", 25)

	fmt.Printf("person (unexported): %s, %d\n", p.name, p.age)
	fmt.Printf("User (exported): %s, %s, %d\n", u.Name, u.Email, u.Age)
}

// ─── 2. IMPORT MECHANICS ──────────────────────────────────────────────────────
//
// Go imports form a DAG (directed acyclic graph) — no circular imports.
// The compiler compiles dependencies before dependents.
//
// Import forms:
//   import "fmt"                    → standard import, use as fmt.Println
//   import myfmt "fmt"             → alias import, use as myfmt.Println
//   import . "fmt"                 → dot import (avoid): Println without fmt.
//   import _ "image/png"           → blank import: run init() only (side effect)
//
// Imports MUST be used — unused imports are compile errors.
// This prevents dead code from accumulating.

func importMechanics() {
	fmt.Println("\n=== 2. Import Mechanics ===")

	// stdlib imports used in this file: fmt, math/rand, os, path/filepath, time
	fmt.Printf("Random: %d\n", rand.Intn(100))
	fmt.Printf("Time: %s\n", time.Now().Format("15:04:05"))

	// path/filepath for OS-independent path operations
	path := filepath.Join("usr", "local", "bin")
	fmt.Println("Joined path:", path)

	// Blank import example (would be: import _ "image/png")
	// This is used to register the PNG image decoder without importing png.Decode
	// The png package's init() does the registration
	fmt.Println("Blank imports run init() to register drivers, codecs, plugins")
}

// ─── 3. init() FUNCTIONS ──────────────────────────────────────────────────────
//
// init() is called automatically after all package-level variables are initialized,
// and before main() is called.
//
// Rules:
//   - A single file can have MULTIPLE init() functions
//   - They run in the order they appear in the file
//   - Across packages, they run in dependency order
//   - You cannot call init() explicitly
//   - init() cannot take arguments or return values
//
// Common uses:
//   - Validate configuration
//   - Register drivers/decoders
//   - Set up package-level state that requires error checking
//   - Database driver registration (database/sql pattern)

var (
	startTime = time.Now()  // initialized before init()
	config    map[string]string
)

func init() {
	// First init in this file
	config = map[string]string{
		"debug": "false",
		"env":   "development",
	}
	fmt.Println("[init 1] config initialized")
}

func init() {
	// Second init in the same file — runs after the first
	fmt.Printf("[init 2] start time: %s\n", startTime.Format("15:04:05"))
	fmt.Printf("[init 2] config: %v\n", config)
}

func initFunctions() {
	fmt.Println("\n=== 3. init() Functions ===")
	fmt.Printf("config after all inits: %v\n", config)
}

// ─── 4. MODULE SYSTEM: go.mod ─────────────────────────────────────────────────
//
// A MODULE is a collection of related Go packages.
// go.mod defines:
//   - module path (the module's import path prefix)
//   - Go version compatibility
//   - Direct and indirect dependencies with exact versions
//
// Example go.mod:
//   module github.com/myname/myapp
//
//   go 1.21
//
//   require (
//       github.com/gin-gonic/gin v1.9.1
//       github.com/lib/pq v1.10.9
//   )
//
//   require (
//       // indirect dependencies (required by direct deps)
//       github.com/bytedance/sonic v1.10.2 // indirect
//   )
//
// go.sum: cryptographic checksums of all dependencies — ensures reproducibility
//
// Commands:
//   go mod init github.com/me/myapp   — create go.mod
//   go get github.com/pkg@v1.2.3      — add/update dependency
//   go mod tidy                       — remove unused deps, add missing
//   go mod vendor                     — copy deps to vendor/ directory
//   go mod download                   — download deps to module cache (~/.cache/go)

func moduleSystem() {
	fmt.Println("\n=== 4. Module System ===")
	fmt.Println("go.mod structure:")
	fmt.Println(`
  module github.com/myname/myapp  ← module path (unique name)

  go 1.21                         ← minimum Go version

  require (
      github.com/gin-gonic/gin v1.9.1  ← direct dependency
  )`)

	fmt.Println("\nKey commands:")
	fmt.Println("  go mod init       — create go.mod")
	fmt.Println("  go get pkg@v1.2   — add/upgrade dependency")
	fmt.Println("  go mod tidy       — clean up go.mod and go.sum")
	fmt.Println("  go mod vendor     — local copy of all deps")
}

// ─── 5. VISIBILITY AND ENCAPSULATION ─────────────────────────────────────────
//
// Go has TWO visibility levels:
//   Exported (uppercase first letter):   visible to ALL packages
//   Unexported (lowercase first letter): visible ONLY within the same package
//
// This applies to:
//   - Functions: func Exported() vs func unexported()
//   - Types: type Exported struct vs type unexported struct
//   - Fields: ExportedField vs unexportedField
//   - Variables: ExportedVar vs unexportedVar
//   - Constants: ExportedConst vs unexportedConst
//
// Note: there is NO "private", "public", "protected" — just exported/unexported.
// Note: there is NO "friend" access — packages are the encapsulation boundary.
// Note: ALL code within the same PACKAGE can access unexported identifiers.

type BankAccount struct {
	OwnerName string    // exported — callers can read
	balance   float64   // unexported — callers cannot access directly
}

func NewBankAccount(owner string, initialBalance float64) *BankAccount {
	return &BankAccount{OwnerName: owner, balance: initialBalance}
}

func (a *BankAccount) Deposit(amount float64) error {
	if amount <= 0 {
		return fmt.Errorf("deposit amount must be positive, got %.2f", amount)
	}
	a.balance += amount
	return nil
}

func (a *BankAccount) Withdraw(amount float64) error {
	if amount <= 0 {
		return fmt.Errorf("withdrawal amount must be positive")
	}
	if amount > a.balance {
		return fmt.Errorf("insufficient funds: have %.2f, need %.2f", a.balance, amount)
	}
	a.balance -= amount
	return nil
}

func (a *BankAccount) Balance() float64 {
	return a.balance  // read-only access via method
}

func visibility() {
	fmt.Println("\n=== 5. Visibility and Encapsulation ===")

	account := NewBankAccount("Alice", 1000.0)
	fmt.Printf("Account owner: %s\n", account.OwnerName) // OK — exported
	// fmt.Println(account.balance)  // COMPILE ERROR — unexported

	account.Deposit(500)
	account.Withdraw(200)
	fmt.Printf("Balance: $%.2f\n", account.Balance())

	if err := account.Withdraw(10000); err != nil {
		fmt.Println("Error:", err)
	}
}

// ─── 6. INTERNAL PACKAGES ─────────────────────────────────────────────────────
//
// The `internal` package convention restricts which packages can import it.
//
// A package at path a/b/internal/c can ONLY be imported by packages
// rooted at a/b. External packages get a compile error.
//
// Use case: you want to share code between packages in your module,
// but not expose it as a public API to external users.
//
// Example:
//   github.com/myapp/internal/db     — can be imported by myapp and myapp/api
//   github.com/myapp/auth/internal   — can only be imported by myapp/auth/...

func internalPackages() {
	fmt.Println("\n=== 6. Internal Packages ===")
	fmt.Println("internal/ packages restrict imports to same module subtree")
	fmt.Println("Example:")
	fmt.Println("  module/internal/db  → importable only within module/...")
	fmt.Println("  external packages get: 'use of internal package' compile error")
}

// ─── 7. BUILD TAGS ────────────────────────────────────────────────────────────
//
// Build tags (//go:build) conditionally include files in the build.
//
//   //go:build linux
//   → This file only compiles on Linux
//
//   //go:build linux && amd64
//   → This file only compiles on Linux x86-64
//
//   //go:build !windows
//   → This file compiles on all platforms except Windows
//
//   //go:build integration
//   → Run with: go test -tags=integration
//
// Common uses:
//   - Platform-specific implementations
//   - Feature flags
//   - Integration tests (don't run in CI without special flag)

func buildTags() {
	fmt.Println("\n=== 7. Build Tags ===")

	// Show current OS and architecture
	fmt.Printf("GOOS:   %s\n", os.Getenv("GOOS"))   // may be empty if not set
	fmt.Printf("GOARCH: %s\n", os.Getenv("GOARCH"))  // may be empty if not set

	fmt.Println("\nBuild tag examples:")
	fmt.Println("  //go:build linux          ← Linux only")
	fmt.Println("  //go:build !cgo           ← When CGO disabled")
	fmt.Println("  //go:build integration    ← go test -tags=integration")
}

func main() {
	packages()
	importMechanics()
	initFunctions()
	moduleSystem()
	visibility()
	internalPackages()
	buildTags()
}

/*
THOUGHT QUESTIONS:

1. Go enforces that all imported packages are used (unused imports = compile error).
   What is the rationale? What is the blank import `_` for?

2. Multiple init() functions can exist in the same file. What order do they run in?
   What order do init() functions run across packages?

3. Why does Go use capitalization for visibility rather than keywords like
   `public`/`private`? What are the advantages?

4. What is an "internal" package? What problem does it solve that
   normal unexported identifiers don't?

5. Go modules use exact versioning in go.sum (SHA256 checksums). Why is this
   important for build reproducibility?

EXERCISES:

1. Create a small library with:
   - An exported type (User) with some exported and some unexported fields
   - An exported constructor (NewUser) that validates inputs
   - Unexported helper functions
   Then write a main package that uses it.

2. Create a package with two init() functions. Add a package-level variable
   that is modified by each init(). Print the final value in main.
   Verify the order is as documented.

3. Create a package structure with an internal/ subdirectory. Verify that
   code outside the expected path cannot import the internal package.

4. Write a build-tag example: create two files with the same function name,
   one with //go:build debug and one without. Compile with and without
   -tags=debug and observe which implementation is used.
*/
