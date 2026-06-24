/*
WEEK 6 — DAY 1: Recursion — Mechanics, Patterns, and Trade-offs
=================================================================
Topic: How recursion works in Go — stack frames, the base case, mutual recursion,
       memoization, and when iteration is preferred.

Key ideas:
  - Every recursive call creates a new stack frame
  - Go does NOT optimize tail calls (no TCO) — deep recursion grows the stack
  - Goroutine stacks grow dynamically, but there is a runtime limit
  - Memoization turns exponential recursion into linear
  - Iteration is usually preferred for performance-critical code in Go
*/

package main

import (
	"fmt"
	"sync"
)

// ─── 1. RECURSION MECHANICS ────────────────────────────────────────────────────
//
// Recursion = a function that calls itself.
// Every call allocates a new STACK FRAME.
// The base case stops the recursion.
// Without a base case → infinite recursion → stack grows until limit hit.
//
// Frame layout for factorial(5):
//   factorial(5) → frame: n=5
//     factorial(4) → frame: n=4
//       factorial(3) → frame: n=3
//         factorial(2) → frame: n=2
//           factorial(1) → frame: n=1 → returns 1
//         returns 2
//       returns 6
//     returns 24
//   returns 120

func factorial(n int) int {
	if n <= 1 {
		return 1  // base case — stops recursion
	}
	return n * factorial(n-1)  // recursive case — calls itself with smaller n
}

// Tracing version to visualize the call stack
func factorialTrace(n, depth int) int {
	indent := ""
	for i := 0; i < depth; i++ {
		indent += "  "
	}
	fmt.Printf("%s→ factorial(%d)\n", indent, n)
	if n <= 1 {
		fmt.Printf("%s← factorial(%d) = 1\n", indent, n)
		return 1
	}
	result := n * factorialTrace(n-1, depth+1)
	fmt.Printf("%s← factorial(%d) = %d\n", indent, n, result)
	return result
}

func recursionMechanics() {
	fmt.Println("=== 1. Recursion Mechanics ===")
	fmt.Printf("factorial(10) = %d\n\n", factorial(10))

	fmt.Println("Tracing factorial(4):")
	factorialTrace(4, 0)
}

// ─── 2. FIBONACCI — NAIVE VS MEMOIZED ────────────────────────────────────────
//
// The naive recursive Fibonacci is O(2^n) — exponential.
// Many subproblems are recomputed: fib(5) computes fib(3) twice.
//
// Memoization stores previously computed results:
//   Time: O(n)  — each value computed once
//   Space: O(n) — the cache

var fibCache sync.Map  // concurrent-safe cache

func fibMemo(n int) int {
	if n <= 1 {
		return n
	}
	if v, ok := fibCache.Load(n); ok {
		return v.(int)
	}
	result := fibMemo(n-1) + fibMemo(n-2)
	fibCache.Store(n, result)
	return result
}

// Bottom-up dynamic programming (no recursion, no cache overhead)
func fibDP(n int) int {
	if n <= 1 { return n }
	a, b := 0, 1
	for i := 2; i <= n; i++ {
		a, b = b, a+b
	}
	return b
}

func fibComparison() {
	fmt.Println("\n=== 2. Fibonacci: Naive vs Memoized vs DP ===")

	// Memoized
	for _, n := range []int{0, 1, 5, 10, 20, 40} {
		fmt.Printf("  fib(%2d) = %d\n", n, fibMemo(n))
	}

	// DP
	fmt.Printf("fibDP(50) = %d\n", fibDP(50))
}

// ─── 3. TREE TRAVERSAL ────────────────────────────────────────────────────────
//
// Recursion shines for tree-shaped data.
// A tree is naturally defined recursively: node + left subtree + right subtree.

type TreeNode struct {
	Value       int
	Left, Right *TreeNode
}

func insert(root *TreeNode, val int) *TreeNode {
	if root == nil {
		return &TreeNode{Value: val}
	}
	if val < root.Value {
		root.Left = insert(root.Left, val)
	} else {
		root.Right = insert(root.Right, val)
	}
	return root
}

func inorder(node *TreeNode, result *[]int) {
	if node == nil { return }
	inorder(node.Left, result)
	*result = append(*result, node.Value)
	inorder(node.Right, result)
}

func height(node *TreeNode) int {
	if node == nil { return 0 }
	lh := height(node.Left)
	rh := height(node.Right)
	if lh > rh { return lh + 1 }
	return rh + 1
}

func countNodes(node *TreeNode) int {
	if node == nil { return 0 }
	return 1 + countNodes(node.Left) + countNodes(node.Right)
}

func printTree(node *TreeNode, prefix string, isLeft bool) {
	if node == nil { return }
	connector := "└── "
	if isLeft { connector = "├── " }
	fmt.Printf("%s%s%d\n", prefix, connector, node.Value)
	childPrefix := prefix + "│   "
	if !isLeft { childPrefix = prefix + "    " }
	printTree(node.Left, childPrefix, true)
	printTree(node.Right, childPrefix, false)
}

func treeTraversal() {
	fmt.Println("\n=== 3. Tree Traversal ===")

	var root *TreeNode
	for _, v := range []int{5, 3, 7, 1, 4, 6, 8, 2} {
		root = insert(root, v)
	}

	fmt.Println("Tree structure:")
	printTree(root, "", false)

	var sorted []int
	inorder(root, &sorted)
	fmt.Println("In-order (sorted):", sorted)
	fmt.Printf("Height: %d, Nodes: %d\n", height(root), countNodes(root))
}

// ─── 4. MUTUAL RECURSION ──────────────────────────────────────────────────────
//
// Two functions that call each other recursively.
// Go supports forward references within a package, so this just works.

func isEven(n int) bool {
	if n == 0 { return true }
	return isOdd(n - 1)
}

func isOdd(n int) bool {
	if n == 0 { return false }
	return isEven(n - 1)
}

// Mutual recursion for parsing expressions (simplified)
// This pattern is used in recursive descent parsers.

type Token struct {
	Kind  string
	Value string
}

func parseNumber(tokens []Token, pos int) (int, int) {
	// Parse a number from tokens starting at pos
	// Returns (value, new_pos)
	if pos >= len(tokens) || tokens[pos].Kind != "NUM" {
		return 0, pos
	}
	var n int
	fmt.Sscanf(tokens[pos].Value, "%d", &n)
	return n, pos + 1
}

func mutualRecursion() {
	fmt.Println("\n=== 4. Mutual Recursion ===")

	for i := 0; i <= 10; i++ {
		fmt.Printf("  isEven(%d) = %v\n", i, isEven(i))
	}
}

// ─── 5. TAIL RECURSION AND WHY GO DOESN'T OPTIMIZE IT ────────────────────────
//
// Tail recursion: the recursive call is the LAST operation in the function.
// (No computation happens after the recursive call returns.)
//
// In languages that optimize tail calls (Scheme, Haskell, some Scala):
//   - The current frame is REUSED instead of a new one being pushed
//   - Deep tail recursion uses O(1) stack space
//
// Go does NOT optimize tail calls. The Go team decided:
//   - It complicates stack traces (debugging would be harder)
//   - Go has goroutines for concurrency, iteration for loops
//   - Most "tail recursive" code can be trivially rewritten as iteration
//
// This means tail-recursive Go code still grows the stack.
// PREFER iteration over recursion in Go for performance.

// Tail-recursive factorial (Go does NOT optimize this)
func factTail(n, acc int) int {
	if n <= 1 { return acc }
	return factTail(n-1, n*acc)  // tail call — but NOT optimized by Go
}

// Iterative factorial — always prefer this in Go
func factIter(n int) int {
	result := 1
	for i := 2; i <= n; i++ {
		result *= i
	}
	return result
}

func tailRecursion() {
	fmt.Println("\n=== 5. Tail Recursion (Go doesn't optimize) ===")

	fmt.Printf("factTail(10, 1) = %d\n", factTail(10, 1))
	fmt.Printf("factIter(10) = %d\n", factIter(10))
	fmt.Println("Both correct, but factIter uses O(1) stack space")
	fmt.Println("factTail creates 10 stack frames in Go (no TCO)")
}

// ─── 6. POWER SETS AND COMBINATORICS ─────────────────────────────────────────
//
// Recursion excels at generating combinatorial structures.

// Generate all subsets of a slice
func powerSet(items []int) [][]int {
	if len(items) == 0 {
		return [][]int{{}}  // base case: empty set has one subset (the empty set)
	}

	first := items[0]
	rest := items[1:]

	// Recursively get all subsets of the rest
	subsetsWithoutFirst := powerSet(rest)

	// Add first to each of those subsets
	var subsetsWithFirst [][]int
	for _, subset := range subsetsWithoutFirst {
		newSubset := make([]int, len(subset)+1)
		newSubset[0] = first
		copy(newSubset[1:], subset)
		subsetsWithFirst = append(subsetsWithFirst, newSubset)
	}

	return append(subsetsWithoutFirst, subsetsWithFirst...)
}

// Generate all permutations
func permutations(items []int) [][]int {
	if len(items) <= 1 {
		return [][]int{append([]int{}, items...)}
	}

	var result [][]int
	for i, item := range items {
		rest := append(append([]int{}, items[:i]...), items[i+1:]...)
		for _, perm := range permutations(rest) {
			result = append(result, append([]int{item}, perm...))
		}
	}
	return result
}

func combinatorics() {
	fmt.Println("\n=== 6. Combinatorics via Recursion ===")

	ps := powerSet([]int{1, 2, 3})
	fmt.Printf("Power set of {1,2,3} (%d subsets):\n", len(ps))
	for _, subset := range ps {
		fmt.Printf("  %v\n", subset)
	}

	perms := permutations([]int{1, 2, 3})
	fmt.Printf("\nPermutations of {1,2,3} (%d):\n", len(perms))
	for _, p := range perms {
		fmt.Printf("  %v\n", p)
	}
}

func main() {
	recursionMechanics()
	fibComparison()
	treeTraversal()
	mutualRecursion()
	tailRecursion()
	combinatorics()
}

/*
THOUGHT QUESTIONS:

1. Every recursive call pushes a new stack frame. What determines how many
   frames fit? What happens if you exceed the limit?

2. Go does not optimize tail calls. What is a tail call? What optimization
   would be applied if Go did support TCO?

3. Why is memoized Fibonacci O(n) time and O(n) space?
   Can you reduce the space complexity to O(1)?

4. The power set of n elements has 2^n subsets. What is the time complexity
   of the powerSet function? Why?

5. When is recursion PREFERRED over iteration? When is it HARMFUL?

EXERCISES:

1. Implement a recursive binary search:
   `func binarySearch(sorted []int, target int) int`
   Returns the index or -1. Then rewrite it iteratively.

2. Write `flatten(tree interface{}) []interface{}` that recursively flattens
   a nested slice structure: [[1,[2,3]],4] → [1,2,3,4].

3. Implement merge sort recursively:
   `func mergeSort(nums []int) []int`
   Then benchmark it against sort.Ints for various input sizes.

4. Implement a recursive descent parser for simple arithmetic expressions:
   expr  = term (('+' | '-') term)*
   term  = factor (('*' | '/') factor)*
   factor = number | '(' expr ')'
*/
