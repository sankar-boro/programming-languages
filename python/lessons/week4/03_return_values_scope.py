"""
WEEK 4 — DAY 3: Functions — Return Values and Scope
=====================================================
Topic: How return works internally, multiple return values, None as implicit
       return, and a deep treatment of Python's scoping rules (LEGB) including
       global, nonlocal, and common scope pitfalls.

Key ideas:
  - return always returns exactly one object (tuples pack multiple values)
  - A function without return returns None
  - Scope is LEXICAL — determined by where code is written, not where it runs
  - LEGB: Local → Enclosing → Global → Built-in search order
  - global and nonlocal are explicit declarations for mutation across scopes
"""

import dis


# ─── 1. HOW RETURN WORKS INTERNALLY ──────────────────────────────────────────
#
# return expr:
#   1. Evaluates expr
#   2. Stores the result in the frame's return slot
#   3. Unwinds the call stack frame
#   4. The caller receives the value from the return slot
#
# There is exactly ONE return value (an object reference).
# No copies. The same object is handed back to the caller.

def show_return_bytecode():
    def add(a, b):
        return a + b

    print("=== return bytecode ===")
    dis.dis(add)
    # Look for RETURN_VALUE — it pops the top of the stack and returns it

show_return_bytecode()


# ─── 2. IMPLICIT RETURN — NONE ────────────────────────────────────────────────
#
# If a function reaches the end without a return statement,
# Python automatically returns None.
# A bare `return` (no expression) also returns None.
#
# This is why print() appears to return something — it doesn't.

print("\n=== Implicit None return ===")

def no_return():
    x = 1 + 1
    # no return — falls off the end

def bare_return():
    return      # explicit None return

result1 = no_return()
result2 = bare_return()
print(f"no_return():   {result1!r}")   # None
print(f"bare_return(): {result2!r}")   # None
print(f"print():       {print('hi')!r}")   # prints "hi", returns None


# ─── 3. MULTIPLE RETURN VALUES ────────────────────────────────────────────────
#
# Python has no syntax for "multiple return values."
# `return a, b` returns ONE tuple (a, b).
# Tuple unpacking at the call site makes it LOOK like multiple returns.

print("\n=== Multiple return values (are tuples) ===")

def min_max(numbers):
    return min(numbers), max(numbers)    # returns ONE tuple

result = min_max([3, 1, 4, 1, 5, 9])
print(f"raw result: {result!r}")         # (1, 9) — a tuple
print(f"type:       {type(result)}")

# Unpack at call site:
low, high = min_max([3, 1, 4, 1, 5, 9])
print(f"low={low}, high={high}")

# Return a named tuple for clarity (preserves both sequence and named access):
from collections import namedtuple

Stats = namedtuple("Stats", ["minimum", "maximum", "mean"])

def describe(numbers):
    n = len(numbers)
    return Stats(
        minimum=min(numbers),
        maximum=max(numbers),
        mean=sum(numbers) / n
    )

s = describe([1, 2, 3, 4, 5])
print(f"\nStats: {s}")
print(f"min={s.minimum}, max={s.maximum}, mean={s.mean:.2f}")
print(f"still unpackable: {s[0]}, {s[1]}, {s[2]}")


# ─── 4. EARLY RETURN — GUARD CLAUSES ──────────────────────────────────────────
#
# Returning early when preconditions fail keeps functions flat and readable.
# Each return is a "this case is handled, stop here."

def safe_divide(a, b):
    if b == 0:
        return None    # guard: invalid input
    return a / b

def process_user(user):
    if not user:
        return "error: no user"
    if not user.get("active"):
        return "error: inactive"
    if not user.get("email"):
        return "error: no email"
    # main logic — only reached if all guards pass
    return f"processing {user['email']}"

print("\n=== Guard clause returns ===")
print(process_user(None))
print(process_user({"active": True}))
print(process_user({"active": True, "email": "a@b.com"}))


# ─── 5. THE LEGB SCOPE RULE ──────────────────────────────────────────────────
#
# When Python looks up a name, it searches these namespaces in order:
#
#   L — Local:     the current function's local scope
#   E — Enclosing: any enclosing function scopes (outer functions)
#   G — Global:    the module-level namespace
#   B — Built-in:  Python's built-in namespace (print, len, range, ...)
#
# The FIRST match wins. Shadowing happens when a local name hides an outer one.

GLOBAL_VAR = "global"

def outer_func():
    enclosing_var = "enclosing"

    def inner_func():
        local_var = "local"

        # LEGB search order for each name:
        print(f"  local_var:     {local_var}")      # L — found in local
        print(f"  enclosing_var: {enclosing_var}")  # E — found in enclosing
        print(f"  GLOBAL_VAR:    {GLOBAL_VAR}")     # G — found in global
        print(f"  len:           {len}")             # B — found in built-in

    inner_func()

print("\n=== LEGB lookup ===")
outer_func()


# ─── 6. SHADOWING ────────────────────────────────────────────────────────────
#
# A local name with the same name as a global one SHADOWS the global.
# Inside the function, the local takes precedence.
# The global is unchanged — you just can't see it through the local name.

x = "global x"

def shadow_example():
    x = "local x"        # this is a NEW local variable, not the global
    print(f"  inside:  {x}")   # "local x"

shadow_example()
print(f"outside: {x}")         # "global x" — unchanged

# Accidental shadowing of built-ins (common mistake):
def bad_practice():
    list = [1, 2, 3]    # shadows the built-in `list`
    # Now you can't use list() as a constructor inside this function!
    # list([1, 2, 3])   ← would fail: TypeError ('list' is not callable)
    return list

print(f"\nbad_practice result: {bad_practice()}")
# After the function returns, the built-in `list` is fine again


# ─── 7. GLOBAL — MODIFYING MODULE-LEVEL VARIABLES ─────────────────────────────
#
# Without global, any assignment inside a function creates a LOCAL variable.
# global explicitly declares that a name refers to the module-level variable.
#
# Use global sparingly — it creates hidden coupling between functions.

count = 0

def increment():
    global count       # declare: "count" refers to the module-level name
    count += 1         # now modifies the global, not a local

print("\n=== global ===")
print(f"count before: {count}")
increment()
increment()
increment()
print(f"count after:  {count}")

# Without global — this FAILS:
def broken_increment():
    try:
        count += 1     # UnboundLocalError: count referenced before assignment
    except UnboundLocalError as e:
        print(f"  UnboundLocalError: {e}")

broken_increment()
# Why? Python sees the assignment `count += 1` and marks `count` as local.
# But it was never assigned locally before the read — hence UnboundLocalError.


# ─── 8. NONLOCAL — MODIFYING ENCLOSING SCOPE VARIABLES ───────────────────────
#
# nonlocal lets an inner function modify a variable from an enclosing function.
# Without nonlocal, assignment in the inner function creates a new local.
#
# This is the mechanism that makes closures with mutable state possible.

print("\n=== nonlocal ===")

def make_counter():
    count = 0

    def increment():
        nonlocal count    # refers to make_counter's local `count`
        count += 1
        return count

    def get():
        return count      # reading is fine without nonlocal

    return increment, get

inc, get = make_counter()
inc(); inc(); inc()
print(f"counter: {get()}")   # 3

# Without nonlocal — fails:
def broken_counter():
    count = 0
    def increment():
        try:
            count += 1     # UnboundLocalError
        except UnboundLocalError as e:
            return f"error: {e}"
    return increment

broken_inc = broken_counter()
print(f"\nbroken_counter: {broken_inc()}")


# ─── 9. SCOPE PITFALL: LATE BINDING IN CLOSURES ───────────────────────────────
#
# This is the second most common Python gotcha after the mutable default argument.
#
# Functions defined in a loop capture the VARIABLE (name), not its VALUE.
# By the time the functions run, the loop variable has its final value.

print("\n=== Late binding closure pitfall ===")

# WRONG:
funcs = []
for i in range(5):
    funcs.append(lambda: i)   # all lambdas refer to the same `i`

print("Wrong (all same):", [f() for f in funcs])   # [4, 4, 4, 4, 4]

# FIX 1: default argument captures the value at definition time
funcs = []
for i in range(5):
    funcs.append(lambda i=i: i)   # i=i evaluates NOW

print("Fixed (default):", [f() for f in funcs])    # [0, 1, 2, 3, 4]

# FIX 2: factory function explicitly closes over the value
def make_func(value):
    return lambda: value

funcs = [make_func(i) for i in range(5)]
print("Fixed (factory):", [f() for f in funcs])    # [0, 1, 2, 3, 4]


# ─── 10. LOCALS() AND GLOBALS() ──────────────────────────────────────────────
#
# locals() — returns a COPY of the current local scope's dict
# globals() — returns the ACTUAL module global namespace dict (live reference)
#
# Modifying globals() modifies the real namespace. locals() is a snapshot.

print("\n=== locals() and globals() ===")

def scope_inspector():
    a = 1
    b = 2
    print(f"  locals(): {locals()}")

scope_inspector()

# globals() is the real thing — you can add to it (but don't):
print(f"  'GLOBAL_VAR' in globals(): {'GLOBAL_VAR' in globals()}")


# ─── EXERCISES ────────────────────────────────────────────────────────────────
#
# 1. Write a function that demonstrates all three of:
#    a) A function modifying a mutable argument (list)
#    b) A function failing to modify an immutable argument (int)
#    c) A function using return to "give back" a new value
#
# 2. Fix the late-binding closure bug for this case:
#       actions = [lambda: print(f"action {n}") for n in range(3)]
#       for act in actions: act()
#    Make each lambda print its own number.
#
# 3. Build a make_accumulator() function using nonlocal:
#    - Each call to the returned function adds a number to a running total
#    - Returns the running total
#       acc = make_accumulator()
#       acc(10)  → 10
#       acc(5)   → 15
#       acc(20)  → 35
#
# 4. Inspect the LEGB resolution of a name using a deliberate shadow:
#    - Shadow the built-in `max` with a local variable inside a function
#    - Show that the built-in is inaccessible inside the function
#    - Show that the built-in is still accessible outside
#    - How would you access the built-in inside if you needed to?
#      (Hint: import builtins)
#
# THOUGHT QUESTION:
#   Python uses LEXICAL scoping (where the code is written determines scope),
#   not DYNAMIC scoping (where the code is called from).
#   In dynamic scoping, calling a function from inside another function would
#   give it access to the caller's variables. Python deliberately chose not to.
#   What are the advantages of lexical scoping for reasoning about code?
