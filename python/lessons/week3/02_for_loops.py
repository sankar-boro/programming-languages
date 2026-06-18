"""
WEEK 3 — DAY 2: Control Flow — for Loops
==========================================
Topic: How for loops work in Python — the iterator protocol, how Python
       actually implements iteration, range internals, and every loop
       pattern from basic to advanced.

Key ideas:
  - for loops use the iterator protocol: __iter__ and __next__
  - range() is a lazy object — it does NOT build a list
  - enumerate(), zip(), and itertools are fundamental tools
  - List comprehensions compile to faster bytecode than equivalent for loops
"""

import dis
import itertools
import timeit


# ─── 1. HOW FOR LOOPS WORK INTERNALLY ────────────────────────────────────────
#
# for x in collection:  is NOT a simple array index loop.
#
# What Python actually does:
#   1. Call iter(collection)  → returns an iterator object
#   2. Repeatedly call next(iterator) → gets each value
#   3. When next() raises StopIteration → the loop ends
#
# Any object implementing __iter__ and __next__ works with for.
# This is the Iterator Protocol.

print("=== For loop desugared ===")

items = [10, 20, 30]

# What Python does under the hood:
iterator = iter(items)          # calls items.__iter__()
while True:
    try:
        value = next(iterator)  # calls iterator.__next__()
        print(f"  got: {value}")
    except StopIteration:
        break                   # loop ends naturally

# The for loop above is exactly equivalent to:
for value in items:
    print(f"  for: {value}")


# ─── 2. BUILDING YOUR OWN ITERATOR ───────────────────────────────────────────
#
# Any class with __iter__ and __next__ works with for loops.
# This lets you iterate over anything — files, sockets, custom data structures.

class CountUp:
    """Iterator that counts from start to stop."""
    def __init__(self, start, stop):
        self.current = start
        self.stop = stop

    def __iter__(self):
        return self    # the object IS its own iterator

    def __next__(self):
        if self.current >= self.stop:
            raise StopIteration
        value = self.current
        self.current += 1
        return value

print("\n=== Custom iterator ===")
for n in CountUp(1, 6):
    print(f"  {n}", end=" ")
print()


# ─── 3. RANGE — LAZY, NOT A LIST ─────────────────────────────────────────────
#
# range(n) does NOT build a list of n integers.
# It stores only: start, stop, step.
# It computes each value on demand.
# This makes range(1_000_000_000) use ~48 bytes, not 8 GB.

import sys

r = range(1_000_000_000)
print(f"\nrange(1_000_000_000) size: {sys.getsizeof(r)} bytes")
print(f"list(range(5)):            {list(range(5))}")
print(f"range membership O(1):     {999_999_999 in r}")   # no iteration needed!

# range() supports:
print(f"range(0, 10, 2):  {list(range(0, 10, 2))}")    # step by 2
print(f"range(10, 0, -1): {list(range(10, 0, -1))}")   # countdown
print(f"range(5, 5):      {list(range(5, 5))}")         # empty range


# ─── 4. LOOP CONTROL: BREAK, CONTINUE, ELSE ──────────────────────────────────
#
# break:    exits the loop immediately
# continue: skips to the next iteration
# else:     runs ONLY if the loop completed WITHOUT a break
#           (this surprises most people — think of it as "loop succeeded")

print("\n=== break / continue / else ===")

# continue: skip even numbers
print("Odd numbers 1-10:")
for n in range(1, 11):
    if n % 2 == 0:
        continue
    print(f"  {n}", end=" ")
print()

# break: stop at first negative
numbers = [3, 7, 1, -2, 5, 9]
print("\nFirst negative:")
for n in numbers:
    if n < 0:
        print(f"  found: {n}")
        break

# else clause — runs only if no break occurred:
def find_prime(candidates, target):
    for c in candidates:
        if target % c == 0:
            print(f"  {target} is divisible by {c} — not prime")
            break
    else:
        print(f"  {target} is prime")

print()
find_prime(range(2, 7), 6)   # divisible by 2
find_prime(range(2, 7), 7)   # no divisor found → else runs


# ─── 5. ENUMERATE — INDEX + VALUE ─────────────────────────────────────────────
#
# Never use `range(len(lst))` to get index + value.
# enumerate() is the idiomatic, readable, faster approach.

fruits = ["apple", "banana", "cherry"]

print("\n=== enumerate ===")
# WRONG (index manually):
for i in range(len(fruits)):
    print(f"  {i}: {fruits[i]}")

# RIGHT (enumerate):
for i, fruit in enumerate(fruits):
    print(f"  {i}: {fruit}")

# Start from a different index:
for i, fruit in enumerate(fruits, start=1):
    print(f"  {i}: {fruit}")


# ─── 6. ZIP — ITERATE MULTIPLE SEQUENCES TOGETHER ────────────────────────────
#
# zip() pairs elements from multiple iterables, stopping at the shortest.
# zip_longest() (from itertools) pads with a fill value instead.

names  = ["Alice", "Bob", "Charlie"]
scores = [85, 92, 78]
grades = ["B", "A", "C"]

print("\n=== zip ===")
for name, score, grade in zip(names, scores, grades):
    print(f"  {name}: {score} ({grade})")

# Unzipping: * unpacking on zip gives back the originals
pairs = [(1, "a"), (2, "b"), (3, "c")]
numbers, letters = zip(*pairs)
print(f"\nnumbers: {numbers}")
print(f"letters: {letters}")

# zip_longest fills missing values
from itertools import zip_longest
short = [1, 2]
long  = ["a", "b", "c", "d"]
for pair in zip_longest(short, long, fillvalue=0):
    print(f"  {pair}")


# ─── 7. LIST COMPREHENSIONS ───────────────────────────────────────────────────
#
# Comprehensions are a compact, readable form of for-loop + append.
# They are compiled to optimized bytecode (LIST_APPEND vs regular APPEND).
# They are typically 10–30% faster than the equivalent for loop.

print("\n=== List comprehensions ===")

# Equivalent forms:
squares_loop = []
for n in range(1, 6):
    squares_loop.append(n ** 2)

squares_comp = [n ** 2 for n in range(1, 6)]   # cleaner, faster

print(f"squares: {squares_comp}")

# With a condition (filter):
evens = [n for n in range(20) if n % 2 == 0]
print(f"evens:   {evens}")

# Nested (cartesian product — outer loop first):
pairs = [(x, y) for x in range(1, 4) for y in range(1, 4) if x != y]
print(f"pairs:   {pairs}")

# Performance comparison:
t_loop = timeit.timeit("[n**2 for n in range(1000)]", number=10_000)
t_comp = timeit.timeit(
    "result = []\nfor n in range(1000): result.append(n**2)",
    number=10_000
)
print(f"\nComprehension: {t_loop:.4f}s")
print(f"For loop:      {t_comp:.4f}s")
print(f"Comprehension is {t_comp / t_loop:.2f}× faster")


# ─── 8. GENERATOR EXPRESSIONS ────────────────────────────────────────────────
#
# Same syntax as list comprehension but with () instead of [].
# Does NOT build the list in memory — generates values one at a time.
# Use when you only iterate once and don't need random access.

import sys

big_list = [n ** 2 for n in range(100_000)]             # builds full list
big_gen  = (n ** 2 for n in range(100_000))             # lazy generator

print(f"\nList size:      {sys.getsizeof(big_list)} bytes")
print(f"Generator size: {sys.getsizeof(big_gen)} bytes")

# sum() works great with generators — no intermediate list built:
total = sum(n ** 2 for n in range(1_000_000))
print(f"sum of squares (1M): {total}")


# ─── 9. ITERTOOLS ESSENTIALS ─────────────────────────────────────────────────
#
# itertools provides efficient, composable iterator building blocks.
# All return lazy iterators — they generate values on demand.

print("\n=== itertools ===")

# chain: flatten multiple iterables into one
from itertools import chain
result = list(chain([1, 2], [3, 4], [5, 6]))
print(f"chain:       {result}")

# islice: slice an iterator (without building a list)
from itertools import islice
first_five = list(islice((n**2 for n in range(100)), 5))
print(f"islice:      {first_five}")

# product: cartesian product (like nested for loops)
from itertools import product
combos = list(product("AB", [1, 2]))
print(f"product:     {combos}")

# groupby: group consecutive equal values (data must be sorted first)
from itertools import groupby
data = [("fruit", "apple"), ("fruit", "banana"), ("veg", "carrot")]
for key, group in groupby(data, key=lambda x: x[0]):
    items = [x[1] for x in group]
    print(f"  {key}: {items}")

# accumulate: running totals
from itertools import accumulate
import operator
data = [1, 2, 3, 4, 5]
print(f"accumulate(sum): {list(accumulate(data))}")            # running sum
print(f"accumulate(mul): {list(accumulate(data, operator.mul))}")  # running product


# ─── EXERCISES ────────────────────────────────────────────────────────────────
#
# 1. Implement FibonacciIterator — an iterator class that yields Fibonacci
#    numbers indefinitely. Use it with islice to get the first 15 numbers.
#
# 2. Write a function flatten(nested_list) that flattens one level of nesting:
#    [[1,2],[3,[4,5]],6] → [1,2,3,[4,5],6]
#    Use a list comprehension. Then write a version that fully flattens
#    arbitrary depth using a for loop and recursion.
#
# 3. Given two lists of equal length, use zip and a comprehension to produce
#    a dict: keys from list1, values from list2. One line.
#
# 4. Use itertools.product to generate all possible 3-letter strings using
#    only "a", "b", "c". How many are there? Verify with len().
#
# 5. Use dis.dis() on a list comprehension vs equivalent for+append loop.
#    Find the LIST_APPEND vs APPEND instruction. Which has fewer bytecode ops?
#
# THOUGHT QUESTION:
#   Python's for loop uses StopIteration to signal the end of iteration.
#   Why is using an exception for normal control flow (not an error)
#   an interesting design choice? What are the trade-offs vs returning
#   a sentinel value like None or -1?
