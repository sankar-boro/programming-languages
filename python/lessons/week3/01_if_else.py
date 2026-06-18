"""
WEEK 3 — DAY 1: Control Flow — if / elif / else
=================================================
Topic: How Python evaluates conditionals at the bytecode level, the
       execution model behind branching, and patterns from basic to advanced.

Key ideas:
  - if is a statement, not an expression (unlike ternary ?: in C/JS)
  - Python has a ternary form: value_if_true if condition else value_if_false
  - Bytecode uses JUMP instructions — conditional is a jump table
  - Pattern matching (match/case, Python 3.10+) replaces complex elif chains
"""

import dis
import sys


# ─── 1. HOW IF WORKS INTERNALLY ──────────────────────────────────────────────
#
# When Python compiles an if statement, it generates bytecode jump instructions.
#
#   if x > 0:
#       do_something()
#
# Compiles to roughly:
#   LOAD x
#   LOAD 0
#   COMPARE_OP >          ← produce True or False
#   POP_JUMP_IF_FALSE L1  ← if False, jump past the body
#   <body code>
#   L1: <rest of code>
#
# There is no "magic" — just conditional jumps in the bytecode.

def show_if_bytecode():
    def check(x):
        if x > 0:
            return "positive"
        else:
            return "non-positive"

    print("=== if/else bytecode ===")
    dis.dis(check)

show_if_bytecode()


# ─── 2. BASIC IF / ELIF / ELSE ───────────────────────────────────────────────
#
# Python evaluates conditions top to bottom and executes the FIRST matching
# branch. Once a branch runs, the rest are skipped entirely.

def classify_number(n):
    """Demonstrate basic branching."""
    if n > 0:
        return "positive"
    elif n < 0:
        return "negative"
    else:
        return "zero"

for n in [5, -3, 0]:
    print(f"classify_number({n:>2}) = {classify_number(n)}")


# ─── 3. TRUTHINESS IN IF CONDITIONS ──────────────────────────────────────────
#
# Python evaluates the truthiness of any expression in an if condition.
# You don't need to write `if x == True` or `if len(lst) > 0`.
# The idiomatic forms are simpler and more readable.

print("\n=== Idiomatic truthiness checks ===")

# List check
items = []
if not items:          # not `if len(items) == 0`
    print("  list is empty")

items = [1, 2]
if items:              # not `if len(items) > 0`
    print("  list has items")

# None check
value = None
if value is None:      # always use `is`, not `== None`
    print("  value is None")

# String check
name = ""
if not name:
    print("  name is empty")


# ─── 4. THE TERNARY EXPRESSION ────────────────────────────────────────────────
#
# Python's "ternary" is an expression (produces a value), unlike an if statement.
# Syntax:  value_if_true  if  condition  else  value_if_false
#
# Use when you need a value based on a condition in one line.
# Avoid nesting ternaries — it destroys readability.

x = 7
label = "odd" if x % 2 != 0 else "even"
print(f"\n{x} is {label}")

# Equivalent if statement:
if x % 2 != 0:
    label = "odd"
else:
    label = "even"

# Ternary in a function argument:
numbers = [3, -1, 4, -1, 5]
absolutes = [n if n >= 0 else -n for n in numbers]   # in a list comprehension
print(f"absolutes: {absolutes}")


# ─── 5. COMMON PATTERNS AND ANTI-PATTERNS ────────────────────────────────────

print("\n=== Patterns ===")

# ANTI-PATTERN: comparing to True/False explicitly
flag = True
if flag == True:      # works but wrong style
    pass
if flag is True:      # also wrong — use truthiness
    pass
if flag:              # correct — idiomatic
    print("  flag is truthy")

# ANTI-PATTERN: nested ternary (unreadable)
n = 5
# result = "big" if n > 10 else "medium" if n > 3 else "small"  # confusing
# BETTER: use if/elif/else or a dict:
def classify(n):
    if n > 10:   return "big"
    if n > 3:    return "medium"
    return "small"

print(f"  classify(5): {classify(5)}")

# PATTERN: early return instead of deep nesting ("guard clause")
def process(user):
    # Without guard clauses:
    # if user:
    #     if user.get("active"):
    #         if user.get("admin"):
    #             do_admin_stuff()
    #
    # With guard clauses (flat is better than nested):
    if not user:
        return "no user"
    if not user.get("active"):
        return "inactive"
    if not user.get("admin"):
        return "not admin"
    return "admin access granted"

print(f"  process(None): {process(None)}")
print(f"  process(active,admin): {process({'active': True, 'admin': True})}")


# ─── 6. DICT AS A DISPATCH TABLE (REPLACING LONG ELIF CHAINS) ─────────────────
#
# Long chains of `elif x == "a": ... elif x == "b":` are slow and verbose.
# A dict lookup is O(1) and separates data from logic.

def handle_add(x, y):      return x + y
def handle_sub(x, y):      return x - y
def handle_mul(x, y):      return x * y
def handle_div(x, y):      return x / y if y != 0 else float("inf")

operations = {
    "+": handle_add,
    "-": handle_sub,
    "*": handle_mul,
    "/": handle_div,
}

def calculate(op, x, y):
    handler = operations.get(op)
    if handler is None:
        raise ValueError(f"Unknown operation: {op}")
    return handler(x, y)

print("\n=== Dict dispatch ===")
for op in ["+", "-", "*", "/"]:
    print(f"  10 {op} 3 = {calculate(op, 10, 3):.4g}")


# ─── 7. MATCH / CASE (PYTHON 3.10+) ──────────────────────────────────────────
#
# Structural pattern matching — far more powerful than switch/case in C.
# It matches structure, not just equality. Works on tuples, lists, dicts, classes.

print("\n=== match / case ===")

def describe_point(point):
    match point:
        case (0, 0):
            return "origin"
        case (x, 0):
            return f"on x-axis at {x}"
        case (0, y):
            return f"on y-axis at {y}"
        case (x, y):
            return f"at ({x}, {y})"
        case _:
            return "not a point"

for p in [(0, 0), (3, 0), (0, -4), (2, 5)]:
    print(f"  {p} → {describe_point(p)}")

# Match with guards (additional conditions):
def classify_http(status):
    match status:
        case 200:
            return "OK"
        case code if 200 <= code < 300:
            return f"Success ({code})"
        case code if 300 <= code < 400:
            return f"Redirect ({code})"
        case code if 400 <= code < 500:
            return f"Client error ({code})"
        case code if 500 <= code < 600:
            return f"Server error ({code})"
        case _:
            return "Unknown"

print()
for code in [200, 201, 301, 404, 500]:
    print(f"  HTTP {code}: {classify_http(code)}")


# ─── 8. PERFORMANCE: IF VS DICT VS MATCH ─────────────────────────────────────
#
# For few branches (≤5): if/elif is fine — straightforward, readable.
# For many branches on a single value: dict is O(1), elif is O(n).
# match/case is compiled to optimized jumps — comparable to dict for simple cases.

import timeit

choices = list("abcdefghij")

def via_elif(c):
    if   c == "a": return 1
    elif c == "b": return 2
    elif c == "c": return 3
    elif c == "d": return 4
    elif c == "e": return 5
    elif c == "f": return 6
    elif c == "g": return 7
    elif c == "h": return 8
    elif c == "i": return 9
    elif c == "j": return 10
    return 0

dispatch = {"a":1,"b":2,"c":3,"d":4,"e":5,"f":6,"g":7,"h":8,"i":9,"j":10}

def via_dict(c):
    return dispatch.get(c, 0)

import random
test_val = "j"   # worst case for elif (last branch)

t_elif = timeit.timeit(lambda: via_elif(test_val), number=500_000)
t_dict = timeit.timeit(lambda: via_dict(test_val), number=500_000)

print(f"\n=== Performance (worst-case branch 'j', n=500k) ===")
print(f"  elif chain: {t_elif:.4f}s")
print(f"  dict lookup: {t_dict:.4f}s")
print(f"  dict is {t_elif / t_dict:.1f}× faster at 10 branches")


# ─── EXERCISES ────────────────────────────────────────────────────────────────
#
# 1. Use dis.dis() on a function with an if/elif/else chain.
#    Identify JUMP_IF_FALSE and JUMP instructions. Map them back to your code.
#
# 2. Write classify_triangle(a, b, c) that returns:
#    "invalid" if the sides don't form a triangle
#    "equilateral", "isosceles", or "scalene"
#    Use only if/elif/else. Then rewrite using match/case.
#
# 3. Rewrite this nested mess using guard clauses:
#       def find_discount(user, cart):
#           if user:
#               if cart:
#                   if len(cart) > 5:
#                       if user.get("premium"):
#                           return 0.20
#                           return 0.10
#                       return 0.05
#                   return 0
#               return 0
#           return 0
#
# 4. Build a dict-based calculator that handles: +, -, *, /, **, %
#    Use lambda for inline handlers. What are the trade-offs vs named functions?
#
# THOUGHT QUESTION:
#   Python has no switch statement (until match in 3.10). Before 3.10,
#   the dict-dispatch pattern was the idiomatic replacement.
#   What does this reveal about Python's philosophy: "data over control flow"?
