"""
WEEK 6 — DAY 2: Recursion vs Iteration
========================================
Topic: When to choose recursion vs iteration, how to convert between them,
       memoization as a bridge, and tail recursion simulation in Python.

Key ideas:
  - Every recursion can be converted to iteration using an explicit stack
  - Every loop can be expressed as recursion (but shouldn't always be)
  - Memoization converts exponential recursion to linear — key technique
  - The call stack IS a stack — recursion uses it implicitly
  - Use recursion when the problem tree is the solution structure
"""

import sys
import functools
import timeit
from collections import deque


# ─── 1. THE EQUIVALENCE — RECURSION IS A HIDDEN STACK ────────────────────────
#
# When you recurse, the call stack stores:
#   - Return address (where to resume after return)
#   - Local variables at each level
#   - Arguments passed to each level
#
# An iterative approach with an explicit stack stores the SAME information
# — you just manage it yourself instead of letting Python manage it.
#
# This means recursion ↔ iteration is always possible (Church-Turing thesis).

# RECURSIVE version:
def sum_recursive(n):
    if n == 0: return 0
    return n + sum_recursive(n - 1)

# ITERATIVE version (same logic, explicit loop):
def sum_iterative(n):
    total = 0
    while n > 0:
        total += n
        n -= 1
    return total

# ITERATIVE with EXPLICIT STACK (mechanically equivalent to recursion):
def sum_explicit_stack(n):
    stack  = []
    total  = 0
    current = n

    while current > 0 or stack:
        if current > 0:
            stack.append(current)   # "push frame" with current value
            current -= 1           # "recurse"
        else:
            total += stack.pop()   # "return and accumulate"

    return total

print("=== Three approaches, same result ===")
for f in [sum_recursive, sum_iterative, sum_explicit_stack]:
    print(f"  {f.__name__}(10) = {f(10)}")


# ─── 2. MEMOIZATION — FIXING EXPONENTIAL RECURSION ───────────────────────────
#
# Naive Fibonacci is O(2^n) because it recomputes the same subproblems.
# Memoization caches results — if we've computed fib(k) before, reuse it.
# This collapses the exponential call tree to a linear chain.
#
# Result: O(n) time, O(n) space.

# Manual memoization with a dict:
def fib_memo(n, cache={}):
    if n in cache:
        return cache[n]
    if n <= 1:
        return n
    cache[n] = fib_memo(n - 1) + fib_memo(n - 2)
    return cache[n]

# Using functools.lru_cache (decorator-based memoization):
@functools.lru_cache(maxsize=None)
def fib_lru(n):
    if n <= 1:
        return n
    return fib_lru(n - 1) + fib_lru(n - 2)

# Iterative (best for pure computation):
def fib_iterative(n):
    if n <= 1: return n
    a, b = 0, 1
    for _ in range(n - 1):
        a, b = b, a + b
    return b

print("\n=== Fibonacci comparison ===")
n = 35
t_naive  = timeit.timeit(lambda: fib_memo.__wrapped__(n) if hasattr(fib_memo, '__wrapped__') else None, number=1)

# Measure memoized vs iterative:
fib_lru.cache_clear()
t_lru = timeit.timeit(lambda: fib_lru(n), number=1000)
t_itr = timeit.timeit(lambda: fib_iterative(n), number=1000)

print(f"  fib({n}) = {fib_iterative(n)}")
print(f"  lru_cache: {t_lru:.6f}s (1000 calls, cached after first)")
print(f"  iterative: {t_itr:.6f}s (1000 calls)")

# lru_cache info:
fib_lru.cache_clear()
fib_lru(20)
print(f"\n  lru_cache info: {fib_lru.cache_info()}")


# ─── 3. WHEN RECURSION WINS — TREE PROBLEMS ──────────────────────────────────
#
# Recursion expresses tree-shaped problems naturally.
# The structure of the code mirrors the structure of the problem.
# Iterative solutions for trees require explicit stacks — more code, less clear.

# Problem: sum all values in a nested dict (arbitrary depth)
data = {
    "a": 1,
    "b": {"c": 2, "d": {"e": 3}},
    "f": 4,
    "g": {"h": {"i": {"j": 5}}}
}

def sum_nested_dict(d):
    """Recursive: mirrors the nested structure directly."""
    total = 0
    for value in d.values():
        if isinstance(value, dict):
            total += sum_nested_dict(value)   # recurse into nested dict
        else:
            total += value
    return total

print(f"\n=== Nested dict sum ===")
print(f"  recursive: {sum_nested_dict(data)}")   # 15


# ─── 4. CONVERTING TREE RECURSION TO ITERATION ───────────────────────────────
#
# Use a stack (or deque) to simulate the call stack explicitly.
# The pattern is always: push initial work, then pop + process + push more work.

def sum_nested_dict_iterative(d):
    """Iterative version using an explicit stack."""
    total = 0
    stack = list(d.values())   # start with top-level values

    while stack:
        item = stack.pop()
        if isinstance(item, dict):
            stack.extend(item.values())   # push nested values for later
        else:
            total += item

    return total

print(f"  iterative: {sum_nested_dict_iterative(data)}")


# ─── 5. TAIL RECURSION — AND WHY PYTHON DOESN'T OPTIMIZE IT ──────────────────
#
# A tail call is when a function's LAST action is a recursive call with
# NO further computation after it returns.
#
# Tail recursive:
#   def factorial_tail(n, acc=1):
#       if n <= 1: return acc
#       return factorial_tail(n - 1, n * acc)   ← last action is the call
#
# NOT tail recursive (must multiply after return):
#   def factorial(n):
#       if n <= 1: return 1
#       return n * factorial(n - 1)    ← must multiply after call returns
#
# Python does NOT optimize tail calls — every call still creates a frame.
# But understanding tail recursion helps you reason about state accumulation.

def factorial_tail(n, acc=1):
    """Tail-recursive form — accumulator carries the result."""
    if n <= 1:
        return acc
    return factorial_tail(n - 1, n * acc)   # accumulator grows, no post-call work

print("\n=== Tail recursive factorial ===")
for n in range(1, 8):
    print(f"  factorial_tail({n}) = {factorial_tail(n)}")

# To make it truly O(1) space in Python, convert to a loop:
def factorial_loop(n):
    """The loop Python SHOULD compile tail recursion to (but doesn't)."""
    acc = 1
    while n > 1:
        acc *= n
        n -= 1
    return acc

print(f"  factorial_loop(7) = {factorial_loop(7)}")


# ─── 6. TRAMPOLINING — SIMULATING TAIL CALL OPTIMIZATION ─────────────────────
#
# A trampoline avoids stack overflow for tail-recursive functions.
# Instead of calling directly, return a callable (a "thunk").
# The trampoline loop calls it, gets another thunk or a final value.
# Only ONE frame is ever active at a time — the trampoline's frame.

def trampoline(func):
    """Decorator: convert a tail-recursive function into an iterative loop."""
    @functools.wraps(func)
    def wrapper(*args, **kwargs):
        result = func(*args, **kwargs)
        while callable(result):
            result = result()
        return result
    return wrapper

@trampoline
def factorial_trampoline(n, acc=1):
    if n <= 1:
        return acc
    # Return a THUNK (lambda with no args) instead of calling directly:
    return lambda: factorial_trampoline(n - 1, n * acc)

print("\n=== Trampoline (avoids stack overflow for tail calls) ===")
sys.setrecursionlimit(100)   # artificially low to demonstrate
try:
    factorial_tail(150)
except RecursionError:
    print("  factorial_tail(150): RecursionError (as expected)")

# Trampoline uses only 1 Python frame regardless of n:
result = factorial_trampoline(150)
print(f"  factorial_trampoline(150): {result}")
sys.setrecursionlimit(1000)   # restore


# ─── 7. DYNAMIC PROGRAMMING — RECURSION + MEMOIZATION SYSTEMATIZED ───────────
#
# Dynamic programming (DP) = break into overlapping subproblems + cache results.
# Two approaches:
#   Top-down: recursive + memoization (intuitive, matches problem definition)
#   Bottom-up: iterative, build solution from small subproblems up (no stack)

# Classic: coin change problem
# Given coin denominations and a target, find minimum coins needed.

# Top-down (memoized recursion):
@functools.lru_cache(maxsize=None)
def min_coins_top_down(coins, amount):
    if amount == 0: return 0
    if amount < 0:  return float("inf")
    return 1 + min(min_coins_top_down(coins, amount - c) for c in coins)

# Bottom-up (iterative DP table):
def min_coins_bottom_up(coins, amount):
    dp = [float("inf")] * (amount + 1)
    dp[0] = 0
    for amt in range(1, amount + 1):
        for coin in coins:
            if coin <= amt:
                dp[amt] = min(dp[amt], 1 + dp[amt - coin])
    return dp[amount]

coins  = (1, 5, 10, 25)
amount = 87

print(f"\n=== Coin change: {amount} cents with {coins} ===")
td = min_coins_top_down(coins, amount)
bu = min_coins_bottom_up(coins, amount)
print(f"  top-down (memoized): {td} coins")
print(f"  bottom-up (table):   {bu} coins")


# ─── EXERCISES ────────────────────────────────────────────────────────────────
#
# 1. Implement merge_sort(lst) recursively.
#    Then implement it iteratively using explicit merging passes (bottom-up).
#    Time both on a list of 10,000 random integers.
#
# 2. Memoize this function manually (without lru_cache):
#       def climbing_stairs(n):
#           # How many ways to climb n stairs taking 1 or 2 steps at a time?
#           if n <= 1: return 1
#           return climbing_stairs(n-1) + climbing_stairs(n-2)
#    Then convert to bottom-up iterative. Verify all three give the same answer.
#
# 3. Implement trampoline yourself (without the decorator form):
#    Write a standalone trampoline(f, *args) function that drives any
#    thunk-returning tail-recursive function.
#
# 4. Tree depth problem:
#    Given a nested dict (any depth), return the maximum nesting depth.
#    Implement recursively. Then implement iteratively using a stack that
#    tracks (item, current_depth). Verify they match.
#
# THOUGHT QUESTION:
#   Memoization works by caching results keyed on function arguments.
#   lru_cache requires arguments to be HASHABLE.
#   Why? What happens if you try to memoize a function that takes a list?
#   How would you work around this limitation?
