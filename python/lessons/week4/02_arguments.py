"""
WEEK 4 — DAY 2: Functions — Arguments In Depth
================================================
Topic: Every kind of Python function argument — positional, keyword,
       *args, **kwargs, keyword-only, positional-only — how they work
       internally and when to use each.

Key ideas:
  - Python passes objects by reference (assignment semantics)
  - *args collects extra positionals into a tuple
  - **kwargs collects extra keyword args into a dict
  - Argument order: positional → *args → keyword-only → **kwargs
  - / and * in signatures enforce call style
"""

import inspect


# ─── 1. HOW PYTHON PASSES ARGUMENTS ──────────────────────────────────────────
#
# Python is often called "pass by object reference" (also called "call by sharing").
# It is NOT "pass by value" (C-style copy) nor "pass by reference" (C++ &).
#
# What actually happens:
#   - The REFERENCE (name binding) is passed, not a copy of the value
#   - Inside the function, the parameter is a new LOCAL name pointing to the SAME object
#   - If the object is mutable, changes inside the function affect the caller's object
#   - If the object is immutable, you can't change it — you can only rebind the name

print("=== Pass by object reference ===")

def try_modify(x):
    x = x + 1    # rebinds local x to a new int — caller's x is unchanged
    print(f"  inside: {x}")

n = 10
try_modify(n)
print(f"  outside: {n}")    # still 10

def modify_list(lst):
    lst.append(99)    # mutates the same list object — caller sees the change

items = [1, 2, 3]
modify_list(items)
print(f"\n  after modify_list: {items}")   # [1, 2, 3, 99]

def replace_list(lst):
    lst = [99]        # rebinds local lst — caller's list unchanged
    print(f"  inside: {lst}")

items = [1, 2, 3]
replace_list(items)
print(f"  after replace_list: {items}")   # [1, 2, 3] unchanged


# ─── 2. POSITIONAL AND KEYWORD ARGUMENTS ─────────────────────────────────────
#
# When calling a function, you can pass arguments:
#   - Positionally: matched by position (left to right)
#   - By keyword:   matched by name (order doesn't matter)
# You can mix both, but positional args must come BEFORE keyword args in the call.

def describe(name, age, city):
    print(f"  {name}, age {age}, from {city}")

print("\n=== Positional vs keyword ===")
describe("Alice", 30, "Paris")                   # all positional
describe(age=30, name="Alice", city="Paris")     # all keyword (any order)
describe("Alice", city="Paris", age=30)          # mixed: positional first


# ─── 3. DEFAULT PARAMETER VALUES ─────────────────────────────────────────────
#
# Parameters with defaults are optional in the call.
# Parameters without defaults are required.
# Required parameters must come BEFORE optional ones in the signature.

def connect(host, port=5432, timeout=30, ssl=True):
    print(f"  connect({host!r}, port={port}, timeout={timeout}, ssl={ssl})")

print("\n=== Defaults ===")
connect("localhost")                            # only required arg
connect("localhost", 3306)                      # override port
connect("localhost", ssl=False)                 # skip to a specific kwarg
connect("prod.db", 5432, 10, False)            # all positional


# ─── 4. *ARGS — VARIADIC POSITIONAL ARGUMENTS ────────────────────────────────
#
# *args collects any extra positional arguments into a TUPLE.
# The name "args" is conventional — any name after * works.
# The * can appear after required positional parameters.

def total(*args):
    """Sum any number of numbers."""
    return sum(args)

print("\n=== *args ===")
print(f"total(): {total()}")
print(f"total(1, 2, 3): {total(1, 2, 3)}")
print(f"total(1, 2, 3, 4, 5): {total(1, 2, 3, 4, 5)}")

# args is always a tuple (immutable):
def show_args(*args):
    print(f"  type(args): {type(args)}, value: {args}")

show_args(1, "hello", [3])

# *args after required params:
def log(level, *messages):
    for msg in messages:
        print(f"  [{level}] {msg}")

log("INFO", "server started", "listening on port 8080")
log("ERROR", "connection refused")


# ─── 5. **KWARGS — VARIADIC KEYWORD ARGUMENTS ────────────────────────────────
#
# **kwargs collects any extra keyword arguments into a DICT.
# The name "kwargs" is conventional.
# Always comes last in the parameter list.

def create_user(**kwargs):
    """Accept any user attributes."""
    print(f"  creating user: {kwargs}")
    return kwargs

print("\n=== **kwargs ===")
create_user(name="Alice", age=30, admin=True)
create_user(name="Bob")

# kwargs is always a dict:
def show_kwargs(**kwargs):
    print(f"  type: {type(kwargs)}, keys: {list(kwargs.keys())}")

show_kwargs(a=1, b=2, c=3)

# Common pattern: pass kwargs through to another function
def styled_print(text, **style_options):
    color  = style_options.get("color", "default")
    weight = style_options.get("weight", "normal")
    print(f"  [{color}/{weight}] {text}")

styled_print("hello", color="red", weight="bold")
styled_print("world")


# ─── 6. COMBINING ALL ARGUMENT TYPES ─────────────────────────────────────────
#
# Full parameter order:
#   def f(pos1, pos2, *args, kw_only1, kw_only2=default, **kwargs)
#
# Rules:
#   - Positional parameters: before *args
#   - *args: catches remaining positionals
#   - Keyword-only: after *args, before **kwargs (must be passed by name)
#   - **kwargs: catches remaining keyword args

def full_example(a, b, *args, option1, option2=False, **kwargs):
    print(f"  a={a}, b={b}")
    print(f"  args={args}")
    print(f"  option1={option1}, option2={option2}")
    print(f"  kwargs={kwargs}")

print("\n=== Full combination ===")
full_example(1, 2, 3, 4, 5, option1="required", extra="surprise")
#             ^  ^  ^^^^^^   ^^^^^^^^^^^^^^^^    ^^^^^^^^^^^^^^^
#             a  b   args    keyword-only        **kwargs


# ─── 7. KEYWORD-ONLY ARGUMENTS (*) ───────────────────────────────────────────
#
# Any parameter after a bare * must be passed by keyword.
# This enforces clarity at call sites — caller must name the argument.
# Very useful for boolean flags that would be confusing positionally.

def resize_image(path, *, width, height, keep_aspect=True):
    print(f"  resize {path!r}: {width}×{height}, aspect={keep_aspect}")

print("\n=== Keyword-only (after *) ===")
resize_image("photo.jpg", width=800, height=600)
resize_image("photo.jpg", width=1920, height=1080, keep_aspect=False)

# This is now a SyntaxError:
try:
    resize_image("photo.jpg", 800, 600)   # positional not allowed after *
except TypeError as e:
    print(f"  TypeError: {e}")


# ─── 8. POSITIONAL-ONLY ARGUMENTS (/) ────────────────────────────────────────
#
# Any parameter before / must be passed by position (Python 3.8+).
# Callers cannot use the name. Useful for:
#   - When the parameter name is an implementation detail
#   - When you want to rename parameters without breaking call sites
#   - For performance (CPython can optimize positional-only arg passing)

def circle_area(r, /):    # r must be positional
    import math
    return math.pi * r ** 2

print("\n=== Positional-only (before /) ===")
print(f"  circle_area(5): {circle_area(5):.4f}")

# This is now a SyntaxError:
try:
    circle_area(r=5)      # keyword not allowed before /
except TypeError as e:
    print(f"  TypeError: {e}")

# / and * can both appear in the same signature:
def full_control(pos_only, /, normal, *, kw_only):
    print(f"  {pos_only=}, {normal=}, {kw_only=}")

full_control(1, 2, kw_only=3)
full_control(1, normal=2, kw_only=3)


# ─── 9. UNPACKING INTO FUNCTION CALLS ────────────────────────────────────────
#
# * unpacks a list/tuple into positional arguments
# ** unpacks a dict into keyword arguments
#
# This is the "call-site" counterpart to *args/**kwargs in definitions.

def add(a, b, c):
    return a + b + c

print("\n=== Argument unpacking ===")
args = [1, 2, 3]
print(f"  add(*args):  {add(*args)}")         # same as add(1, 2, 3)

kwargs = {"a": 1, "b": 2, "c": 3}
print(f"  add(**kwargs): {add(**kwargs)}")    # same as add(a=1, b=2, c=3)

# Mixing:
print(f"  add(1, *[2, 3]): {add(1, *[2, 3])}")

# Useful for forwarding args to another function:
def wrapper(*args, **kwargs):
    print(f"  wrapping call with args={args}, kwargs={kwargs}")
    return add(*args, **kwargs)

wrapper(1, 2, 3)
wrapper(a=1, b=2, c=3)


# ─── EXERCISES ────────────────────────────────────────────────────────────────
#
# 1. Write a function safe_call(func, *args, **kwargs) that:
#    - Calls func(*args, **kwargs)
#    - Returns (result, None) on success
#    - Returns (None, error_message) on any exception
#    Test it with ZeroDivisionError, TypeError, etc.
#
# 2. Write a function that accepts an arbitrary number of numbers and
#    returns (min, max, mean, median) as a tuple.
#    Use *args. Handle the edge case of zero arguments.
#
# 3. Demonstrate that Python is NOT pass-by-value using three examples:
#    a) A function that appears to modify an int (but doesn't)
#    b) A function that genuinely modifies a list
#    c) A function that rebinds a list name (caller list unchanged)
#    Print before and after for each.
#
# 4. Write a function create_endpoint(path, /, *, method="GET", auth=False)
#    that demonstrates both positional-only and keyword-only in one signature.
#    Write 3 valid calls and 2 invalid calls (show the TypeErrors).
#
# THOUGHT QUESTION:
#   *args gives you a TUPLE, not a list. **kwargs gives you a DICT, not an
#   OrderedDict (but dicts are ordered in Python 3.7+). Why might these
#   specific types have been chosen? What would change if *args gave a list?
