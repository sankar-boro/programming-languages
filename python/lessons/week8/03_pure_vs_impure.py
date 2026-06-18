"""
WEEK 8 — DAY 3: Pure vs Impure Functions
==========================================
Topic: What purity means, why it matters, how Python makes impurity
       easy (and sometimes necessary), and the practical trade-offs
       between pure functional style and Python's pragmatic approach.

Key ideas:
  - A pure function: same inputs → same output, no side effects
  - Impure function: depends on or modifies external state
  - Purity enables: testability, caching, parallelism, reasoning
  - Python is NOT a pure functional language — impurity is sometimes right
  - Idempotency, referential transparency, and functional core patterns
"""

import functools
import random
import time
import copy


# ─── 1. DEFINING PURITY ──────────────────────────────────────────────────────
#
# A function is PURE if:
#   1. DETERMINISTIC: given the same inputs, it always returns the same output
#   2. NO SIDE EFFECTS: it doesn't modify any state outside itself
#                       (no I/O, no mutation of arguments, no global changes)
#
# A function is IMPURE if it violates either condition.

# PURE: depends only on its argument, modifies nothing
def add(a, b):
    return a + b

# PURE: builds a new list, doesn't modify the input
def doubled(lst):
    return [x * 2 for x in lst]

# IMPURE: reads from global state (result depends on external variable)
multiplier = 3
def scale(x):
    return x * multiplier   # depends on global `multiplier` — NOT deterministic!

# IMPURE: modifies its argument
def append_in_place(lst, value):
    lst.append(value)   # side effect: mutates the caller's list

# IMPURE: I/O side effect
def log_and_return(x):
    print(f"value: {x}")   # side effect: writes to stdout
    return x

print("=== Pure vs Impure ===")
print(f"add(3, 4):     {add(3, 4)}")       # always 3+4=7
print(f"doubled([1,2]): {doubled([1,2])}")  # always [2,4]

multiplier = 3
print(f"scale(5):      {scale(5)}")         # 15 — but if multiplier changes...
multiplier = 10
print(f"scale(5):      {scale(5)}")         # now 50! same input, different output


# ─── 2. SIDE EFFECTS — A TAXONOMY ────────────────────────────────────────────
#
# Side effects (things that make a function impure):
#
#   Reads:                             Writes:
#   - global/module variables          - global/module variables
#   - mutable arguments                - mutable arguments (mutation)
#   - current time / random            - files / databases / network
#   - I/O (stdin, files, network)      - stdout / stderr / logging
#   - system state (env vars, clock)   - modifying shared data structures
#
# None of these are forbidden — they're just not "pure."
# Python code does all of these. The question is: do you know which
# functions are pure and which aren't?

print("\n=== Side effect taxonomy ===")

# Reading current time → impure (different result each call)
def current_timestamp():
    return time.time()

t1 = current_timestamp()
t2 = current_timestamp()
print(f"  same function, different results: {t1 == t2}")   # False

# Random → impure
def roll_die():
    return random.randint(1, 6)

# Mutating argument → impure (caller's data changes)
def zero_negatives(lst):
    """Impure: modifies the list in place."""
    for i, v in enumerate(lst):
        if v < 0:
            lst[i] = 0

data = [1, -2, 3, -4]
zero_negatives(data)
print(f"  data after zero_negatives: {data}")   # original is changed


# ─── 3. WHY PURITY MATTERS ────────────────────────────────────────────────────
#
# 1. TESTABILITY: pure functions are trivially testable — no setup/teardown,
#    no mocking, no state to reset. assert f(x) == expected is the whole test.
#
# 2. CACHEABILITY: pure functions CAN be memoized safely.
#    Impure functions (reading time, random, globals) cannot.
#
# 3. PARALLELISM: pure functions can run concurrently — no race conditions,
#    no shared mutable state.
#
# 4. REASONING: you can understand a pure function in isolation.
#    No need to trace global state or argument mutation history.

# Demonstrating safe memoization only for pure functions:
@functools.lru_cache(maxsize=None)
def fib(n):   # pure — correct to cache
    if n <= 1: return n
    return fib(n-1) + fib(n-2)

# If roll_die were cached — WRONG behavior:
# @functools.lru_cache(maxsize=None)
# def roll_die(): return random.randint(1,6)
# → always returns same result after first call!

print(f"\n=== Purity and caching ===")
print(f"  fib(10) = {fib(10)}")
print(f"  cache:   {fib.cache_info()}")


# ─── 4. REFERENTIAL TRANSPARENCY ─────────────────────────────────────────────
#
# A pure function call is "referentially transparent" — you can replace
# the call with its return value without changing program behavior.
#
# Pure:    add(3, 4)  →  can always be replaced with  7
# Impure:  print(x)  →  cannot be replaced with None (the side effect matters)
#
# Referential transparency enables:
#   - Compiler optimizations (constant folding across calls)
#   - Safe refactoring (extract/inline without behavior change)
#   - Equational reasoning (prove properties algebraically)

def discount_price(price, rate):
    """Pure: referentially transparent."""
    return round(price * (1 - rate), 2)

# These are equivalent (referential transparency):
# total = discount_price(100, 0.1) + discount_price(200, 0.1)
# total = 90.0 + 180.0
# total = 270.0

print(f"\n=== Referential transparency ===")
print(f"  discount_price(100, 0.1) = {discount_price(100, 0.1)}")
print(f"  discount_price(200, 0.1) = {discount_price(200, 0.1)}")
print(f"  always the same values — can reason algebraically")


# ─── 5. PURE CORE, IMPURE SHELL ───────────────────────────────────────────────
#
# Best practice: push impure code to the EDGES (I/O, random, time, config),
# keep the CORE logic pure. This maximizes testability and reasoning.
#
# Pattern: "Functional Core, Imperative Shell"
#   - Pure functions do all computation
#   - Impure code at the boundary reads input and writes output

# IMPURE shell: reads from "external source"
def fetch_prices():
    """Impure: would read from DB/API in real code."""
    return {"apple": 1.20, "banana": 0.50, "cherry": 3.00}

# PURE core: all computation, no I/O
def apply_discount(prices: dict, category_discounts: dict) -> dict:
    """Pure: deterministic transformation, no side effects."""
    return {
        item: round(price * (1 - category_discounts.get(item, 0)), 2)
        for item, price in prices.items()
    }

def format_receipt(prices: dict) -> str:
    """Pure: deterministic string building."""
    lines = [f"  {item:<12} ${price:.2f}" for item, price in sorted(prices.items())]
    total = sum(prices.values())
    lines.append(f"  {'TOTAL':<12} ${total:.2f}")
    return "\n".join(lines)

# IMPURE shell: orchestrates I/O
def run_checkout():
    prices    = fetch_prices()                          # impure: I/O
    discounts = {"banana": 0.20, "cherry": 0.10}
    final     = apply_discount(prices, discounts)       # pure
    receipt   = format_receipt(final)                   # pure
    print(receipt)                                      # impure: I/O

print("\n=== Functional core, imperative shell ===")
run_checkout()


# ─── 6. IDEMPOTENCY ───────────────────────────────────────────────────────────
#
# A function is IDEMPOTENT if calling it multiple times with the same input
# produces the same result as calling it once.
#
# f(f(x)) == f(x)  for all x
#
# This is related to but weaker than purity.
# Idempotent operations are safe to retry — crucial for distributed systems.

print("\n=== Idempotency ===")

# Idempotent (calling multiple times = calling once):
def normalize(s: str) -> str:
    """Strip and lowercase."""
    return s.strip().lower()

s = "  HELLO  "
print(f"  normalize once:   {normalize(s)!r}")
print(f"  normalize twice:  {normalize(normalize(s))!r}")   # same result
print(f"  idempotent:       {normalize(s) == normalize(normalize(s))}")

# NOT idempotent (calling multiple times changes the result):
counter = [0]
def increment():
    counter[0] += 1
    return counter[0]

# increment() → 1
# increment(increment()) → 2  — different result on second call


# ─── 7. PRACTICAL IMPURITY — WHEN SIDE EFFECTS ARE CORRECT ───────────────────
#
# Pure code is not always the goal. Some problems REQUIRE impurity:
#   - Logging and monitoring (must write to files/stdout)
#   - Caching in mutable structures (desired side effect)
#   - User interaction (must read input)
#   - Database transactions (must persist state)
#
# Python is a pragmatic language. The goal is CONTROLLED impurity:
# - Know which functions are pure vs impure
# - Keep impure functions thin
# - Don't mix I/O logic with computational logic

# Good: thin impure wrapper around a pure core
def save_result(computation, *args, filepath=None, **kwargs):
    """Impure: I/O wrapper. Core logic (computation) stays pure."""
    result = computation(*args, **kwargs)   # pure computation
    if filepath:
        with open(filepath, "w") as f:      # impure I/O
            f.write(str(result))
    return result

# The computation function itself stays pure and testable.


# ─── 8. TESTING PURE VS IMPURE ────────────────────────────────────────────────
#
# Pure functions: test directly, no mocks needed
# Impure functions: need fixtures, mocks, or dependency injection

print("\n=== Testing characteristics ===")

# Pure function test — trivial:
assert add(2, 3) == 5
assert doubled([1, 2, 3]) == [2, 4, 6]
assert discount_price(100, 0.1) == 90.0
print("  pure function tests: pass (no setup needed)")

# For impure functions, best practice is to inject dependencies:
def process_data(data, writer=print):
    """writer is injected — makes the function testable."""
    result = [x * 2 for x in data]
    writer(result)    # default: print (impure), but can inject a mock
    return result

captured = []
process_data([1, 2, 3], writer=captured.append)   # inject mock writer
print(f"  captured via injection: {captured}")


# ─── EXERCISES ────────────────────────────────────────────────────────────────
#
# 1. Classify these functions as pure or impure, and state WHY:
#    a) def f(x): return x ** 2
#    b) def f(lst): lst.sort(); return lst
#    c) def f(lst): return sorted(lst)
#    d) def f(x): return random.choice(x)
#    e) def f(x, y=[]): y.append(x); return y
#    f) def f(x): return x if x > 0 else abs(x)
#
# 2. Refactor this impure function into a pure core + impure shell:
#       def process_orders(orders):
#           result = []
#           for order in orders:
#               if order["status"] == "pending":
#                   order["status"] = "processed"    # mutates input!
#                   print(f"Processing order {order['id']}")  # I/O
#                   result.append(order)
#           return result
#
# 3. Write a function is_pure_ish(func) that heuristically detects impurity by:
#    - Checking if co_names (global lookups) contains any mutable globals
#    - Checking if the function uses LOAD_GLOBAL for non-builtin names
#    - Checking co_flags for whether it's a generator (not necessarily impure)
#    Use dis module. This won't be perfect — explain the limitations.
#
# 4. Implement a memoize decorator that ONLY caches if the function is pure.
#    Add a `pure=True` parameter: if False, skip caching.
#    Bonus: detect impurity automatically using a heuristic from exercise 3.
#
# THOUGHT QUESTION:
#   Python has `functools.lru_cache` which caches function results.
#   But it can only be safely applied to PURE functions.
#   Yet Python doesn't enforce this — you can apply it to impure functions.
#   What are the failure modes? Give three concrete examples of bugs
#   introduced by incorrectly caching an impure function.
