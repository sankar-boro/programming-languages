"""
WEEK 4 — DAY 1: Functions — Definition and First-Class Nature
==============================================================
Topic: How functions are defined, what makes them "first-class objects,"
       how Python stores and calls them internally, and what the function
       object contains.

Key ideas:
  - def is a statement that creates a function object and binds it to a name
  - Functions are first-class: assignable, passable, storable
  - A function object carries code, defaults, closure cells, and metadata
  - Lambda creates an anonymous function expression (limited to one expression)
"""

import dis
import inspect
import types


# ─── 1. WHAT DEF ACTUALLY DOES ────────────────────────────────────────────────
#
# def greet(name): ...
#
# Is NOT a declaration. It is a STATEMENT that:
#   1. Compiles the function body to a code object
#   2. Creates a function object wrapping that code
#   3. Binds the name "greet" in the current namespace to that function object
#
# This happens at RUNTIME, not at import time (well, import runs the module top-level).
# You can define functions conditionally, inside loops, inside other functions.

def show_function_as_object():
    def add(a, b):
        return a + b

    print("=== Function is an object ===")
    print(f"type:       {type(add)}")           # <class 'function'>
    print(f"id:         {id(add)}")             # memory address
    print(f"name:       {add.__name__}")         # "add"
    print(f"module:     {add.__module__}")       # current module
    print(f"code:       {add.__code__}")         # the compiled code object
    print(f"defaults:   {add.__defaults__}")     # None (no defaults)
    print(f"globals:    {type(add.__globals__)}") # the module's global namespace

show_function_as_object()


# ─── 2. FUNCTIONS ARE FIRST-CLASS OBJECTS ────────────────────────────────────
#
# "First-class" means: can be assigned to variables, passed as arguments,
# returned from functions, stored in data structures. Just like ints or strings.

print("\n=== First-class functions ===")

def square(x):     return x * x
def cube(x):       return x * x * x
def negate(x):     return -x

# Assign to a variable:
transform = square
print(f"transform(5): {transform(5)}")   # 25

# Store in a data structure:
operations = [square, cube, negate]
for op in operations:
    print(f"  {op.__name__}(3) = {op(3)}")

# Pass as an argument:
def apply(func, value):
    return func(value)

print(f"\napply(square, 7): {apply(square, 7)}")
print(f"apply(cube, 3):   {apply(cube, 3)}")

# Return from a function:
def make_multiplier(factor):
    def multiplier(x):
        return x * factor   # factor is captured from the enclosing scope
    return multiplier

double = make_multiplier(2)
triple = make_multiplier(3)
print(f"\ndouble(5): {double(5)}")
print(f"triple(5): {triple(5)}")


# ─── 3. THE CODE OBJECT ────────────────────────────────────────────────────────
#
# Every function has a __code__ attribute — a code object.
# The code object contains everything about the compiled function:
#   - bytecode instructions
#   - variable names (co_varnames)
#   - argument count (co_argcount)
#   - constant values (co_consts)
#   - local variable count (co_nlocals)
#   - source file and line numbers

def add(a, b, c=0):
    x = a + b
    return x + c

code = add.__code__
print("\n=== Code object ===")
print(f"co_argcount:   {code.co_argcount}")    # 3 (includes default params)
print(f"co_varnames:   {code.co_varnames}")    # ('a', 'b', 'c', 'x')
print(f"co_consts:     {code.co_consts}")      # (None, 0) or similar
print(f"co_nlocals:    {code.co_nlocals}")     # 4
print(f"co_filename:   {code.co_filename}")
print(f"co_firstlineno:{code.co_firstlineno}")

print("\n=== Bytecode ===")
dis.dis(add)


# ─── 4. DEFAULT ARGUMENTS — EVALUATED ONCE ───────────────────────────────────
#
# Default values are evaluated ONCE when the def statement executes.
# NOT each time the function is called.
#
# This is Python's most infamous gotcha for beginners.

print("\n=== Default argument trap ===")

def append_to(element, target=[]):    # [] is created ONCE
    target.append(element)
    return target

print(append_to(1))    # [1]
print(append_to(2))    # [1, 2]  — not [2]! same list reused
print(append_to(3))    # [1, 2, 3]

# The default list lives in:
print(f"default stored in: {append_to.__defaults__}")

# CORRECT pattern: use None as default, create inside function
def append_to_fixed(element, target=None):
    if target is None:
        target = []    # new list on each call
    target.append(element)
    return target

print(f"\nfixed: {append_to_fixed(1)}")    # [1]
print(f"fixed: {append_to_fixed(2)}")    # [2] — independent


# ─── 5. LAMBDA — ANONYMOUS FUNCTION EXPRESSIONS ───────────────────────────────
#
# lambda args: expression
#
# Creates a function object without a def statement.
# Restricted to a SINGLE expression (no statements, no multi-line body).
# Used for short, throwaway functions — especially as arguments.
#
# Lambda is NOT faster than def. It creates the same function object.

print("\n=== Lambda ===")

# Equivalent:
def square_def(x): return x * x
square_lambda = lambda x: x * x

print(f"def:    {square_def(5)}")
print(f"lambda: {square_lambda(5)}")
print(f"same type: {type(square_def) == type(square_lambda)}")   # True

# Primary use: short functions as arguments
numbers = [3, -1, 4, -1, 5, -9, 2, 6]
print(f"\nsorted by abs: {sorted(numbers, key=lambda x: abs(x))}")
print(f"max by abs:    {max(numbers, key=lambda x: abs(x))}")

# Lambda captures the enclosing scope (closures apply):
multipliers = [lambda x, n=n: x * n for n in range(1, 4)]
print(f"\nlambda(5): {[m(5) for m in multipliers]}")  # [5, 10, 15]
# Note: n=n captures the current value — without it, all lambdas use last n


# ─── 6. FUNCTION INTROSPECTION ────────────────────────────────────────────────
#
# inspect module provides high-level introspection tools.
# These are used by IDEs, debuggers, testing frameworks, and decorators.

print("\n=== Function introspection ===")

def greet(name: str, greeting: str = "Hello") -> str:
    """Return a greeting string."""
    return f"{greeting}, {name}!"

sig = inspect.signature(greet)
print(f"signature: {sig}")

for param_name, param in sig.parameters.items():
    print(f"  param: {param_name}")
    print(f"    default:   {param.default}")
    print(f"    annotation:{param.annotation}")
    print(f"    kind:      {param.kind.name}")

# Type hints are stored in __annotations__:
print(f"\nannotations: {greet.__annotations__}")

# Docstring:
print(f"docstring:   {greet.__doc__}")


# ─── 7. NESTED FUNCTIONS ──────────────────────────────────────────────────────
#
# Functions can be defined inside other functions.
# Inner functions are created fresh on EACH CALL to the outer function.
# They can see the outer function's local variables (lexical scoping).

print("\n=== Nested functions ===")

def counter(start=0):
    count = start      # local to counter()

    def increment(step=1):
        nonlocal count     # declare we want to modify the enclosing variable
        count += step
        return count

    def reset():
        nonlocal count
        count = start
        return count

    return increment, reset

inc, rst = counter(10)
print(f"inc():   {inc()}")     # 11
print(f"inc(5):  {inc(5)}")    # 16
print(f"rst():   {rst()}")     # 10
print(f"inc():   {inc()}")     # 11 again


# ─── 8. FUNCTIONS DEFINED CONDITIONALLY ──────────────────────────────────────
#
# Since def is a runtime statement, you can define functions conditionally.
# This is used in feature detection, platform-specific code, monkey-patching.

import sys

if sys.platform == "win32":
    def get_separator():
        return "\\"
else:
    def get_separator():
        return "/"

print(f"\nseparator: {get_separator()}")

# Factory pattern — different implementations based on config:
def make_formatter(style="text"):
    if style == "text":
        def fmt(value):
            return str(value)
    elif style == "json":
        import json
        def fmt(value):
            return json.dumps(value)
    else:
        raise ValueError(f"Unknown style: {style}")
    return fmt

text_fmt = make_formatter("text")
json_fmt  = make_formatter("json")
print(f"text: {text_fmt({'a': 1})}")
print(f"json: {json_fmt({'a': 1})}")


# ─── EXERCISES ────────────────────────────────────────────────────────────────
#
# 1. Inspect a function's __code__ object:
#    Write a function that takes TWO functions as arguments and returns True
#    if they have the same argument count and same number of local variables.
#    Hint: use co_argcount and co_nlocals.
#
# 2. The default argument trap:
#    Write a function make_default_dict() that is supposed to return
#    an empty dict each time but accidentally shares one dict.
#    Then fix it using the None pattern.
#
# 3. Write a function compose(*funcs) that returns a new function applying
#    each function in sequence (right to left):
#       double_then_negate = compose(negate, double)
#       double_then_negate(5) → -10
#    Use lambda or a nested def.
#
# 4. Use inspect.signature() to write a function describe(func) that prints:
#    - The function name
#    - All parameter names with their types and defaults
#    - The return annotation if present
#    Test it on built-in functions like len, print, sorted.
#
# THOUGHT QUESTION:
#   When you write `double = make_multiplier(2)`, the inner function
#   `multiplier` "remembers" that factor=2. Where is this value stored?
#   It can't be on the call stack (make_multiplier already returned).
#   What mechanism keeps `factor` alive? (Preview of Week 7: closures)
