"""
WEEK 8 — DAY 1: Higher-Order Functions
========================================
Topic: Functions that take or return functions — what makes them powerful,
       the standard library HOFs, building your own, and the patterns
       they enable (map, filter, reduce, compose, partial application).

Key ideas:
  - A higher-order function (HOF) takes a function as argument OR returns one
  - HOFs enable abstraction over BEHAVIOR, not just data
  - map/filter/reduce are the three fundamental HOFs
  - functools.partial, functools.reduce, functools.wraps are essential tools
  - HOFs compose — stacking them is how functional pipelines are built
"""

import functools
import operator
import timeit
from typing import Callable, TypeVar

T = TypeVar("T")
U = TypeVar("U")


# ─── 1. WHAT MAKES A FUNCTION HIGHER-ORDER ────────────────────────────────────
#
# A function is higher-order if it:
#   a) Takes one or more functions as arguments, OR
#   b) Returns a function as its result (or both)
#
# This is only possible because functions are first-class objects in Python.

def apply_twice(func, value):
    """HOF: takes a function, applies it twice."""
    return func(func(value))

def square(x): return x * x
def negate(x): return -x

print("=== Higher-order functions ===")
print(f"apply_twice(square, 3):  {apply_twice(square, 3)}")   # square(square(3)) = 81
print(f"apply_twice(negate, 5):  {apply_twice(negate, 5)}")   # negate(negate(5)) = 5

# Functions returning functions:
def make_power(exponent):
    """HOF: returns a function that raises to a given power."""
    return lambda x: x ** exponent

square_fn = make_power(2)
cube_fn   = make_power(3)
print(f"square_fn(4): {square_fn(4)}")   # 16
print(f"cube_fn(3):   {cube_fn(3)}")     # 27


# ─── 2. MAP — TRANSFORM EVERY ELEMENT ────────────────────────────────────────
#
# map(func, iterable) → applies func to each element, returns a lazy iterator.
# Does NOT build a list in memory — generates values on demand.
# Equivalent to: (func(x) for x in iterable)

print("\n=== map ===")
numbers = [1, 2, 3, 4, 5]

# map returns a map object (lazy iterator):
doubled = map(lambda x: x * 2, numbers)
print(f"type: {type(doubled)}")
print(f"list(map(double, numbers)): {list(doubled)}")

# With a named function (cleaner for complex transforms):
def celsius_to_fahrenheit(c):
    return c * 9 / 5 + 32

temps_c = [0, 20, 37, 100]
temps_f = list(map(celsius_to_fahrenheit, temps_c))
print(f"celsius → fahrenheit: {temps_f}")

# map with multiple iterables (zips them):
a = [1, 2, 3]
b = [10, 20, 30]
sums = list(map(operator.add, a, b))
print(f"map(add, a, b): {sums}")   # [11, 22, 33]

# Performance: map(func, iterable) vs list comprehension:
t_map  = timeit.timeit(lambda: list(map(lambda x: x*2, range(1000))), number=10000)
t_comp = timeit.timeit(lambda: [x*2 for x in range(1000)],            number=10000)
print(f"\nmap:           {t_map:.4f}s")
print(f"comprehension: {t_comp:.4f}s")
# Comprehension is often slightly faster for lambdas; map wins with named C functions


# ─── 3. FILTER — SELECT ELEMENTS ─────────────────────────────────────────────
#
# filter(predicate, iterable) → yields elements where predicate(element) is truthy.
# predicate=None filters by truthiness.
# Also lazy — returns an iterator.

print("\n=== filter ===")

numbers = range(-5, 6)
positives = list(filter(lambda x: x > 0, numbers))
print(f"positives: {positives}")

# filter with None — keeps truthy values:
mixed = [0, 1, "", "hello", None, [], [1, 2], False, True]
truthy = list(filter(None, mixed))
print(f"truthy:    {truthy}")

# Equivalent comprehension:
positives_comp = [x for x in range(-5, 6) if x > 0]
print(f"comp:      {positives_comp}")


# ─── 4. REDUCE — FOLD / ACCUMULATE ───────────────────────────────────────────
#
# functools.reduce(func, iterable[, initial])
# Applies func cumulatively: func(func(func(a, b), c), d) ...
# The result "accumulates" into a single value.
#
# This is the most powerful and general of the three HOFs.
# map and filter can be expressed as reduce (but shouldn't be in practice).

from functools import reduce

print("\n=== reduce ===")

numbers = [1, 2, 3, 4, 5]

# Sum: reduce((acc, x) → acc + x, [1,2,3,4,5]) = ((((1+2)+3)+4)+5)
total = reduce(operator.add, numbers)
print(f"reduce(add, [1..5]):     {total}")   # 15

# Product:
product = reduce(operator.mul, numbers)
print(f"reduce(mul, [1..5]):     {product}")  # 120

# Max (don't do this in practice — use max()):
maximum = reduce(lambda acc, x: acc if acc > x else x, numbers)
print(f"reduce(max, [1..5]):     {maximum}")  # 5

# With initial value (important when iterable may be empty):
total_with_init = reduce(operator.add, [], 0)   # empty iterable, initial=0
print(f"reduce(add, [], init=0): {total_with_init}")   # 0

# Reduce to build a data structure:
words = ["the", "quick", "brown", "fox"]
frequency = reduce(
    lambda acc, word: {**acc, word: acc.get(word, 0) + 1},
    words,
    {}
)
print(f"word freq: {frequency}")


# ─── 5. FUNCTION COMPOSITION ─────────────────────────────────────────────────
#
# Composition: apply multiple functions in sequence.
# compose(f, g)(x) = f(g(x))   — right to left (mathematical convention)
# pipe(f, g)(x)    = g(f(x))   — left to right (data flow convention)

def compose(*funcs):
    """Apply functions right to left: compose(f, g, h)(x) = f(g(h(x)))"""
    def composed(value):
        result = value
        for func in reversed(funcs):
            result = func(result)
        return result
    return composed

def pipe(*funcs):
    """Apply functions left to right: pipe(f, g, h)(x) = h(g(f(x)))"""
    def piped(value):
        result = value
        for func in funcs:
            result = func(result)
        return result
    return piped

# Can also implement compose using reduce:
def compose_reduce(*funcs):
    return reduce(lambda f, g: lambda x: f(g(x)), funcs)

add1    = lambda x: x + 1
double  = lambda x: x * 2
square2 = lambda x: x ** 2

pipeline = pipe(add1, double, square2)    # (x+1)*2 then squared
print(f"\n=== Function composition ===")
print(f"pipe(add1, double, square)(3):    {pipeline(3)}")     # ((3+1)*2)^2 = 64

transform = compose(square2, double, add1)   # same order, right to left
print(f"compose(square, double, add1)(3): {transform(3)}")   # same: 64


# ─── 6. FUNCTOOLS.PARTIAL — PARTIAL APPLICATION ───────────────────────────────
#
# functools.partial(func, *args, **kwargs)
# Returns a new function with some arguments pre-filled.
# This is partial application — NOT currying (currying is one arg at a time).

from functools import partial

def power(base, exponent):
    return base ** exponent

square_p = partial(power, exponent=2)    # fix exponent=2
cube_p   = partial(power, exponent=3)
double_p = partial(operator.mul, 2)     # fix first arg of mul

print(f"\n=== functools.partial ===")
print(f"square_p(5): {square_p(5)}")    # 25
print(f"cube_p(3):   {cube_p(3)}")      # 27
print(f"double_p(7): {double_p(7)}")    # 14

# partial preserves the original function reference:
print(f"square_p.func:     {square_p.func}")
print(f"square_p.keywords: {square_p.keywords}")

# Use case: customizing generic functions:
import json
compact_json = partial(json.dumps, separators=(",", ":"), sort_keys=True)
pretty_json  = partial(json.dumps, indent=2, sort_keys=True)

data = {"b": 2, "a": 1, "c": [1, 2, 3]}
print(f"\ncompact: {compact_json(data)}")
print(f"pretty:\n{pretty_json(data)}")


# ─── 7. FUNCTOOLS.WRAPS — PRESERVING FUNCTION METADATA ───────────────────────
#
# When you write a decorator (a HOF that wraps a function), the wrapper
# replaces the original function. This loses __name__, __doc__, __module__.
# functools.wraps(original) copies this metadata to the wrapper.
# Always use @wraps in decorators.

print("\n=== functools.wraps ===")

def bad_decorator(func):
    def wrapper(*args, **kwargs):
        print(f"  calling {func.__name__}")
        return func(*args, **kwargs)
    return wrapper   # wrapper loses original metadata

def good_decorator(func):
    @functools.wraps(func)   # copies __name__, __doc__, __module__, __annotations__
    def wrapper(*args, **kwargs):
        print(f"  calling {func.__name__}")
        return func(*args, **kwargs)
    return wrapper

@bad_decorator
def my_function_bad():
    """I do something important."""
    pass

@good_decorator
def my_function_good():
    """I do something important."""
    pass

print(f"bad_decorator:  __name__={my_function_bad.__name__!r}, __doc__={my_function_bad.__doc__!r}")
print(f"good_decorator: __name__={my_function_good.__name__!r}, __doc__={my_function_good.__doc__!r}")


# ─── 8. BUILDING A FUNCTIONAL PIPELINE ───────────────────────────────────────
#
# HOFs compose naturally into data transformation pipelines.
# This is the core idea of functional programming: data flows through
# a sequence of pure transformations.

print("\n=== Functional pipeline ===")

# Data: raw transaction records
transactions = [
    {"amount": 150.0, "category": "food",     "valid": True},
    {"amount":  -20.0, "category": "refund",  "valid": True},
    {"amount": 300.0, "category": "food",     "valid": False},  # invalid
    {"amount":  80.0, "category": "transport","valid": True},
    {"amount": 120.0, "category": "food",     "valid": True},
]

# Pipeline using HOFs:
valid_txns      = filter(lambda t: t["valid"], transactions)
food_txns       = filter(lambda t: t["category"] == "food", valid_txns)
food_amounts    = map(lambda t: t["amount"], food_txns)
total_food_spend = reduce(operator.add, food_amounts, 0.0)

print(f"  total valid food spend: ${total_food_spend:.2f}")

# Same pipeline using pipe() + lambdas (declarative style):
get_valid = partial(filter, lambda t: t["valid"])
get_food  = partial(filter, lambda t: t["category"] == "food")
get_amount = partial(map,   lambda t: t["amount"])
sum_amounts = partial(reduce, operator.add, initial=0.0)

pipeline = pipe(get_valid, get_food, get_amount, sum_amounts)
print(f"  pipeline result:        ${pipeline(transactions):.2f}")


# ─── EXERCISES ────────────────────────────────────────────────────────────────
#
# 1. Implement your own map() and filter() as pure Python HOFs that return
#    generators. They should work with any iterable, not just lists.
#    Verify against the built-in versions with timeit.
#
# 2. Use reduce to:
#    a) Flatten [[1,2],[3,4],[5,6]] into [1,2,3,4,5,6]
#    b) Group a list of dicts by a key (like itertools.groupby but returns a dict)
#    c) Build a nested dict from a flat list of (key, value) pairs
#
# 3. Write a robust compose() that:
#    - Accepts any number of functions
#    - Raises TypeError with a clear message if any argument is not callable
#    - Returns the identity function (lambda x: x) if called with no arguments
#    - Has a __name__ describing the composed functions
#
# 4. Build a data validation pipeline using HOFs:
#    Given: records = [{"name": "Alice", "age": 30}, {"name": "", "age": -1}, ...]
#    Chain: filter(has_name) → filter(valid_age) → map(normalize) → list
#    where normalize capitalizes name and ensures age is int.
#
# THOUGHT QUESTION:
#   map and filter return LAZY iterators. reduce() consumes its iterator eagerly.
#   Why does reduce HAVE to be eager? Can you think of a case where lazy reduce
#   would be semantically problematic?
#   (Hint: think about what reduce(f, [a, b, c]) means — can f be lazy?)
