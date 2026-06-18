"""
WEEK 7 — DAY 1: Closures
==========================
Topic: What closures are at the memory level, how free variables are stored
       in cell objects, when closures are created, and practical patterns.

Key ideas:
  - A closure = function + its enclosing scope's free variables (cell objects)
  - Free variables live in "cells" — heap-allocated boxes that outlive the frame
  - Closures enable stateful functions without classes
  - Python closures are READ by default; use `nonlocal` to WRITE
  - The __closure__ attribute exposes the captured cells
"""

import dis
import inspect


# ─── 1. WHAT IS A CLOSURE? ───────────────────────────────────────────────────
#
# A closure is a function that "closes over" (captures) variables from
# its enclosing scope — variables that aren't local to it and aren't global.
#
# The enclosing function has RETURNED. Its frame is gone.
# But the closed-over variables are still alive — stored in cell objects.
# The inner function holds a reference to those cells.
#
# Without closures, inner functions could only see global state.
# With closures, each call to the outer function creates INDEPENDENT state.

def make_greeting(salutation):
    """Outer function: salutation becomes a free variable in greet."""
    def greet(name):
        return f"{salutation}, {name}!"   # salutation is a FREE variable
    return greet

hello = make_greeting("Hello")
hi    = make_greeting("Hi")
hey   = make_greeting("Hey")

print("=== Basic closure ===")
print(hello("Alice"))   # "Hello, Alice!"
print(hi("Bob"))        # "Hi, Bob!"
print(hey("Charlie"))   # "Hey, Charlie!"

# Each closure has its OWN captured salutation:
print(f"\nhello.__closure__: {hello.__closure__}")
print(f"captured value:    {hello.__closure__[0].cell_contents}")
print(f"hi captured value: {hi.__closure__[0].cell_contents}")


# ─── 2. HOW CELLS WORK INTERNALLY ────────────────────────────────────────────
#
# When Python detects that a variable in outer() is used by inner(),
# it allocates a "cell" object on the HEAP (not the C stack).
# The outer frame's local variable slot POINTS to the cell.
# The inner function's __closure__ tuple also POINTS to the same cell.
#
# When outer() returns and its frame is destroyed, the cell SURVIVES
# because inner() still holds a reference to it.
#
#   make_greeting frame (destroyed after return):
#     salutation_slot → cell {"Hello"}
#
#   hello function object (persists):
#     __closure__[0]  → cell {"Hello"}   ← same cell, still alive

def inspect_closure():
    x = 42

    def inner():
        return x    # x is a free variable — will be stored in a cell

    return inner

f = inspect_closure()

print("\n=== Closure internals ===")
print(f"f.__closure__:              {f.__closure__}")
print(f"type of cell:               {type(f.__closure__[0])}")
print(f"cell_contents:              {f.__closure__[0].cell_contents}")
print(f"f.__code__.co_freevars:     {f.__code__.co_freevars}")   # ('x',)
print(f"f.__code__.co_cellvars:     {'(none in inner — x is free, not cell)'}")

# The outer function's code object lists x as a cellvar:
print(f"\nouter co_cellvars:          {inspect_closure.__code__.co_cellvars}")   # ('x',)


# ─── 3. CLOSURES VS CLASSES — STATEFUL FUNCTIONS ─────────────────────────────
#
# Closures create stateful functions without defining a class.
# Each call to the factory creates a fresh, independent closure.
# This is the functional equivalent of a single-method class.

def make_counter(start=0, step=1):
    """Factory: returns a stateful counter function."""
    count = start

    def counter():
        nonlocal count    # write to the captured cell
        current = count
        count += step
        return current

    return counter

c1 = make_counter()
c2 = make_counter(start=100, step=10)

print("\n=== Closure as stateful function ===")
print(f"c1: {c1()}, {c1()}, {c1()}")         # 0, 1, 2
print(f"c2: {c2()}, {c2()}, {c2()}")         # 100, 110, 120
print(f"c1 again: {c1()}")                    # 3 — independent from c2


# ─── 4. CLOSURES READ VARIABLES BY DEFAULT ────────────────────────────────────
#
# Without `nonlocal`, an inner function can READ but not WRITE to a captured var.
# Attempting to assign without nonlocal creates a NEW local variable instead —
# and the read of the old value before assignment causes UnboundLocalError.

def broken_counter():
    count = 0
    def increment():
        count += 1       # attempts to read AND write count
        return count     # Python marks count as local → UnboundLocalError on read
    return increment

inc = broken_counter()
try:
    inc()
except UnboundLocalError as e:
    print(f"\n=== Without nonlocal: ===")
    print(f"  UnboundLocalError: {e}")

# With nonlocal — works:
def working_counter():
    count = 0
    def increment():
        nonlocal count
        count += 1
        return count
    return increment

inc = working_counter()
print(f"  with nonlocal: {inc()}, {inc()}, {inc()}")


# ─── 5. MULTIPLE FUNCTIONS SHARING ONE CLOSURE ────────────────────────────────
#
# Multiple inner functions can share the SAME cell.
# This lets you build objects with multiple methods sharing private state.

def make_stack():
    """A stack implementation using a closure — no class needed."""
    items = []   # shared by push, pop, and peek

    def push(x):
        items.append(x)

    def pop():
        if not items:
            raise IndexError("pop from empty stack")
        return items.pop()

    def peek():
        if not items:
            raise IndexError("peek at empty stack")
        return items[-1]

    def size():
        return len(items)

    return push, pop, peek, size

push, pop, peek, size = make_stack()

print("\n=== Shared closure (stack) ===")
push(1); push(2); push(3)
print(f"  size:  {size()}")
print(f"  peek:  {peek()}")
print(f"  pop:   {pop()}")
print(f"  size:  {size()}")


# ─── 6. CLOSURES AND THE LATE-BINDING PITFALL ────────────────────────────────
#
# (Revisited with full understanding now that we know what a cell is.)
#
# The cell stores a REFERENCE to the variable's slot, not the VALUE.
# When the loop variable changes, ALL closures see the new value —
# because they all reference the SAME cell.

print("\n=== Late-binding pitfall (explained) ===")

# All closures share ONE cell for `i`:
funcs = [lambda: i for i in range(5)]
print(f"wrong (all share cell): {[f() for f in funcs]}")   # [4,4,4,4,4]

# Fix 1: default argument captures the VALUE (not the cell):
funcs = [lambda i=i: i for i in range(5)]
print(f"fixed (default arg):    {[f() for f in funcs]}")   # [0,1,2,3,4]

# Fix 2: factory function creates a new cell per call:
def make_lambda(i):
    return lambda: i   # each call creates a new cell with its own i

funcs = [make_lambda(i) for i in range(5)]
print(f"fixed (factory):        {[f() for f in funcs]}")   # [0,1,2,3,4]


# ─── 7. BYTECODE: HOW CLOSURES ARE COMPILED ───────────────────────────────────
#
# Free variables use LOAD_DEREF / STORE_DEREF opcodes (not LOAD_FAST).
# LOAD_DEREF: look up the value in the cell object referenced by __closure__
# STORE_DEREF: store a value into the cell (what nonlocal enables)

def outer_for_dis():
    captured = 10
    def inner():
        return captured + 1
    return inner

inner_func = outer_for_dis()

print("\n=== Closure bytecode (LOAD_DEREF) ===")
dis.dis(inner_func)
# Look for LOAD_DEREF — it accesses the cell, not a local variable slot


# ─── 8. PRACTICAL PATTERNS ───────────────────────────────────────────────────

print("\n=== Practical closure patterns ===")

# Pattern 1: Partial application (before functools.partial)
def multiply(a, b):
    return a * b

def make_multiplier(factor):
    return lambda x: multiply(x, factor)

double = make_multiplier(2)
triple = make_multiplier(3)
print(f"  double(7): {double(7)}")
print(f"  triple(7): {triple(7)}")

# Pattern 2: Memoization via closure (manual lru_cache):
def memoize(func):
    cache = {}    # captured in the closure
    def wrapper(*args):
        if args not in cache:
            cache[args] = func(*args)
        return cache[args]
    wrapper.cache = cache   # expose cache for inspection
    return wrapper

@memoize
def expensive_square(n):
    return n * n

print(f"\n  expensive_square(5): {expensive_square(5)}")
print(f"  expensive_square(5): {expensive_square(5)}")   # from cache
print(f"  cache: {expensive_square.cache}")

# Pattern 3: Event handler with context
def make_button_handler(button_name, log):
    """Closure captures button_name and log without classes."""
    def handle_click():
        log.append(f"{button_name} clicked")
    return handle_click

event_log = []
submit_handler = make_button_handler("Submit", event_log)
cancel_handler = make_button_handler("Cancel", event_log)

submit_handler()
cancel_handler()
submit_handler()
print(f"\n  event log: {event_log}")


# ─── EXERCISES ────────────────────────────────────────────────────────────────
#
# 1. Inspect __closure__ on a function that closes over TWO variables.
#    How many cells does it have? What order are they in?
#    (Hint: compare with __code__.co_freevars for the order)
#
# 2. Build make_adder(n) that returns a function adding n to its argument.
#    Then build make_pipeline(*funcs) that returns a function applying
#    each function in sequence. Use closures (no classes).
#
# 3. Write a closure-based rate limiter:
#       make_rate_limited(func, max_calls_per_second)
#    Returns a wrapped function that raises RuntimeError if called more
#    than max_calls_per_second times in a 1-second window.
#    Use time.time() and a captured deque of call timestamps.
#
# 4. Use dis.dis() to compare a function using LOAD_FAST (local) vs
#    LOAD_DEREF (closure). Count the instructions.
#    Which is more expensive? Why?
#
# THOUGHT QUESTION:
#   A closure captures the variable's CELL, not the variable's VALUE.
#   In Python, you can't have two closures share a mutable counter WITHOUT
#   some shared mutable container (a cell via nonlocal, or a list/dict).
#   In functional languages like Haskell, closures capture VALUES (immutable).
#   What are the advantages and disadvantages of each approach?
