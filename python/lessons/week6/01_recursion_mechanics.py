"""
WEEK 6 — DAY 1: Recursion — Mechanics
========================================
Topic: How recursion works at the frame level, the base-case contract,
       tracing recursive calls, and the cost of recursion in Python.

Key ideas:
  - Every recursive call creates a new stack frame
  - Recursion has a natural base case — the frame that returns without recursing
  - Python does NOT optimize tail calls (unlike Scheme/Haskell)
  - For deep recursion, Python raises RecursionError at ~1000 frames
  - The "call tree" is the mental model — draw it before coding
"""

import sys
import functools
import inspect


# ─── 1. WHAT RECURSION IS — AT THE FRAME LEVEL ───────────────────────────────
#
# A recursive function calls ITSELF, creating a new frame each time.
# The frames stack up until the base case is hit — then they unwind.
#
# factorial(3):
#   frame: factorial(3)  → 3 * factorial(2)
#     frame: factorial(2)  → 2 * factorial(1)
#       frame: factorial(1)  → base case → returns 1
#     frame: factorial(2) resumes → 2 * 1 = 2 → returns 2
#   frame: factorial(3) resumes → 3 * 2 = 6 → returns 6

def factorial(n):
    # Base case: stop recursion
    if n <= 1:
        return 1
    # Recursive case: reduce the problem
    return n * factorial(n - 1)

print("=== Factorial ===")
for n in range(0, 8):
    print(f"  factorial({n}) = {factorial(n)}")


# ─── 2. TRACING RECURSIVE CALLS ──────────────────────────────────────────────
#
# Visualize each frame being created and destroyed.
# The indentation mirrors the call depth.

def factorial_traced(n, depth=0):
    indent = "  " * depth
    print(f"{indent}→ factorial({n})")
    if n <= 1:
        print(f"{indent}← base case: return 1")
        return 1
    result = n * factorial_traced(n - 1, depth + 1)
    print(f"{indent}← return {n} * ... = {result}")
    return result

print("\n=== Traced factorial(4) ===")
factorial_traced(4)


# ─── 3. THE TWO REQUIREMENTS: BASE CASE + REDUCTION ──────────────────────────
#
# Every correct recursive function MUST have:
#
# 1. BASE CASE:   a condition that stops recursion and returns directly
#                 without making another recursive call
#
# 2. REDUCTION:   each recursive call must move CLOSER to the base case
#                 (the argument must shrink or simplify toward the base case)
#
# Missing either one causes infinite recursion (RecursionError).

def bad_no_base(n):
    """Missing base case — infinite recursion."""
    return n * bad_no_base(n - 1)   # no stopping condition

def bad_no_reduction(n):
    """Doesn't reduce toward base case — also infinite."""
    if n == 0:
        return 1
    return n * bad_no_reduction(n)  # n never changes

print("\n=== Missing base case (RecursionError) ===")
try:
    bad_no_base(5)
except RecursionError:
    print("  RecursionError: no base case!")


# ─── 4. THINKING WITH CALL TREES ─────────────────────────────────────────────
#
# For recursion with branching (multiple recursive calls), draw the call tree.
# Fibonacci is the classic example — calls branch into a binary tree.
#
#         fib(4)
#        /       \
#     fib(3)    fib(2)
#    /     \    /    \
#  fib(2) fib(1) fib(1) fib(0)
#  /   \
# fib(1) fib(0)
#
# Nodes: 9 function calls just to compute fib(4).
# fib(n) makes ~2^n calls — exponential time complexity!

def fib_naive(n, call_count=[0]):
    call_count[0] += 1
    if n <= 1:
        return n
    return fib_naive(n - 1) + fib_naive(n - 2)

print("\n=== Naive fibonacci call count ===")
for n in [5, 10, 15, 20]:
    fib_naive.__defaults__[0].__setitem__(0, 0)   # reset count
    from types import MethodType
    count_box = [0]
    def fib_counted(n, box=count_box):
        box[0] += 1
        if n <= 1: return n
        return fib_counted(n-1) + fib_counted(n-2)
    fib_counted(n)
    print(f"  fib({n:<2}): {count_box[0]:>7,} calls")


# ─── 5. MUTUAL RECURSION ─────────────────────────────────────────────────────
#
# Two functions calling each other — still recursion, still needs a base case.

def is_even(n):
    if n == 0:
        return True
    return is_odd(n - 1)

def is_odd(n):
    if n == 0:
        return False
    return is_even(n - 1)

print("\n=== Mutual recursion ===")
for n in range(7):
    print(f"  is_even({n}): {is_even(n)}")


# ─── 6. COST OF RECURSION IN PYTHON ──────────────────────────────────────────
#
# Each recursive call in Python:
#   1. Creates a new PyFrameObject (~200 bytes)
#   2. Evaluates the function call overhead (argument binding, etc.)
#   3. Adds an entry to the C call stack
#
# Python does NOT perform tail-call optimization.
# Even a perfect tail call creates a new frame.
# Deep recursion is expensive both in time and memory.

import timeit

def sum_recursive(n):
    if n == 0: return 0
    return n + sum_recursive(n - 1)

def sum_iterative(n):
    total = 0
    while n > 0:
        total += n
        n -= 1
    return total

def sum_builtin(n):
    return sum(range(n + 1))

n = 900   # keep below recursion limit
t_rec = timeit.timeit(lambda: sum_recursive(n), number=1000)
t_itr = timeit.timeit(lambda: sum_iterative(n), number=1000)
t_blt = timeit.timeit(lambda: sum_builtin(n),   number=1000)

print(f"\n=== Sum 1..{n} performance ===")
print(f"  recursive: {t_rec:.4f}s")
print(f"  iterative: {t_itr:.4f}s")
print(f"  builtin:   {t_blt:.4f}s")
print(f"  recursion is {t_rec/t_itr:.1f}× slower than loop")
print(f"  builtin is  {t_itr/t_blt:.1f}× faster than loop")


# ─── 7. WHEN RECURSION IS THE RIGHT TOOL ─────────────────────────────────────
#
# Recursion shines when the PROBLEM STRUCTURE is recursive:
#
#   - Tree traversal (file systems, HTML DOM, binary trees)
#   - Divide and conquer algorithms (merge sort, quicksort, binary search)
#   - Backtracking (maze solving, N-queens, permutations)
#   - Parsing nested/hierarchical structures (JSON, XML, expressions)
#
# For flat sequences with a simple reduction, iteration is better.

def flatten(nested):
    """Recursively flatten a nested list of arbitrary depth."""
    result = []
    for item in nested:
        if isinstance(item, list):
            result.extend(flatten(item))   # recurse into sub-lists
        else:
            result.append(item)
    return result

print("\n=== Recursive flatten ===")
data = [1, [2, [3, [4, [5]]]], 6, [7, 8]]
print(f"  input:  {data}")
print(f"  output: {flatten(data)}")

# Binary tree traversal — recursion is the natural expression:
class TreeNode:
    def __init__(self, val, left=None, right=None):
        self.val = val
        self.left = left
        self.right = right

def inorder(node):
    """In-order traversal: left → root → right."""
    if node is None:
        return []
    return inorder(node.left) + [node.val] + inorder(node.right)

tree = TreeNode(4,
    TreeNode(2, TreeNode(1), TreeNode(3)),
    TreeNode(6, TreeNode(5), TreeNode(7))
)
print(f"\n=== Binary tree inorder traversal ===")
print(f"  result: {inorder(tree)}")   # [1, 2, 3, 4, 5, 6, 7]


# ─── 8. RECURSION LIMIT AND WORKAROUND ───────────────────────────────────────
#
# Python's default limit is sys.getrecursionlimit() (~1000).
# For genuinely deep recursion (e.g., deep trees), you need an iterative
# approach using an explicit stack — mimicking what the call stack would do.

def flatten_iterative(nested):
    """Flatten using an explicit stack instead of call stack."""
    result = []
    stack  = list(nested)   # start with top-level items

    while stack:
        item = stack.pop()
        if isinstance(item, list):
            stack.extend(item)    # push sub-items for later processing
        else:
            result.append(item)

    return result[::-1]   # reverse because stack pops in LIFO order

print("\n=== Iterative flatten (explicit stack) ===")
data = [1, [2, [3, [4, [5]]]], 6]
print(f"  {flatten_iterative(data)}")


# ─── EXERCISES ────────────────────────────────────────────────────────────────
#
# 1. Write power(base, exp) using recursion. Do NOT use ** or math.pow.
#    Then write an optimized version using "fast exponentiation":
#       if exp is even: power(base, exp) = power(base*base, exp//2)
#       if exp is odd:  power(base, exp) = base * power(base, exp-1)
#    Compare call counts for power(2, 64) — naive vs optimized.
#
# 2. Write a recursive function count_nodes(tree) where tree is a dict:
#       {"val": 1, "left": {"val": 2, "left": None, "right": None}, "right": None}
#    Return the total number of nodes.
#
# 3. Write permutations(lst) that returns all permutations of a list.
#    Do NOT use itertools. Think recursively:
#    "For each element e, prepend e to all permutations of the remaining elements."
#
# 4. Convert flatten() to use an explicit stack (as shown above).
#    Verify both versions produce identical results on deeply nested lists.
#    Time both on a list nested 500 levels deep: [[[...[[1]]...]]].
#
# THOUGHT QUESTION:
#   Languages like Scheme, Haskell, and Erlang perform tail-call optimization:
#   a function in tail position reuses the current frame instead of creating
#   a new one, making recursion as memory-efficient as a loop.
#   Python's Guido van Rossum explicitly rejected tail-call optimization.
#   His reason: it would destroy tracebacks (you'd lose the call history).
#   Do you agree with this trade-off? What does it reveal about Python's
#   philosophy — debugging clarity over raw performance?
