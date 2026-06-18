"""
WEEK 1 — DAY 2: Syntax, Values, and Types
==========================================
Topic: Python's type system, how types work at runtime, and what
       "dynamic typing" really means under the hood.

Key ideas:
  - Types live on objects, not on variable names
  - Python's built-in types are implemented in C (in CPython)
  - Type checking happens at runtime, not compile time
"""


# ─── 1. EXPRESSIONS VS STATEMENTS ────────────────────────────────────────────
#
# Expression: any code that PRODUCES a value
# Statement:  any code that PERFORMS an action (may not produce a value)
#
# The distinction matters because you can nest expressions but not statements.

# Expressions — all of these produce a value
result = 2 + 3          # arithmetic expression → 5
name   = "py" + "thon"  # string expression     → "python"
check  = 5 > 3          # comparison expression  → True

# Statements — these perform actions
x = 10                  # assignment statement
print(x)                # call statement (print returns None)

# A statement cannot be used where a value is expected:
# y = (x = 10)   ← SyntaxError in Python (unlike C/JS)
# Python 3.8+ introduced the walrus operator := for assignment in expressions:
data = [1, 2, 3, 4, 5]
if (n := len(data)) > 3:
    print(f"List has {n} elements")   # n was assigned AND tested


# ─── 2. PYTHON'S CORE TYPES ──────────────────────────────────────────────────
#
# Python has a hierarchy of built-in types. Every type is itself an object
# of type `type`. (Yes — type is its own type. More on this in Week 4.)

print("\n=== Core built-in types ===")

# Numeric types
i = 42            # int   — arbitrary precision (no overflow in Python)
f = 3.14          # float — IEEE 754 double precision (64-bit)
c = 2 + 3j        # complex — two floats (real + imaginary)

# Text type
s = "hello"       # str   — immutable sequence of Unicode code points

# Boolean type
t = True          # bool  — subclass of int (True == 1, False == 0)
ff = False

# None type
n = None          # NoneType — singleton, only one None exists

# Sequence types
lst  = [1, 2, 3]      # list  — mutable, ordered
tup  = (1, 2, 3)      # tuple — immutable, ordered
rng  = range(5)       # range — lazy sequence, generates on demand

# Mapping type
dct  = {"a": 1}       # dict  — mutable, key-value pairs

# Set types
st   = {1, 2, 3}      # set         — mutable, unordered, unique
fst  = frozenset({1}) # frozenset   — immutable set

for obj in [i, f, c, s, t, n, lst, tup, dct, st]:
    print(f"  {repr(obj):<20} → {type(obj).__name__}")


# ─── 3. DYNAMIC TYPING IN DEPTH ──────────────────────────────────────────────
#
# Static typing  (Java, C): the VARIABLE has a type, checked at compile time
# Dynamic typing (Python):  the OBJECT has a type, checked at runtime
#
# The variable name is just a label — it can point at any type at any time.

print("\n=== Dynamic typing ===")

value = 10
print(f"value is: {value}, type: {type(value).__name__}")

value = "ten"
print(f"value is: {value}, type: {type(value).__name__}")

value = [1, 2, 3]
print(f"value is: {value}, type: {type(value).__name__}")

# Python does NOT prevent you from reusing names for different types.
# This is flexible but can hide bugs. Type hints (Week 5+) help with this.


# ─── 4. TYPE CHECKING AT RUNTIME ─────────────────────────────────────────────
#
# isinstance() checks the type hierarchy (preferred)
# type() checks exact type (stricter — misses subclasses)

print("\n=== Type checking ===")

x = True
print(f"type(x) is bool:          {type(x) is bool}")       # True
print(f"type(x) is int:           {type(x) is int}")        # False — exact check
print(f"isinstance(x, bool):      {isinstance(x, bool)}")   # True
print(f"isinstance(x, int):       {isinstance(x, int)}")    # True — bool IS an int

# bool is a subclass of int — True and False are literally 1 and 0
print(f"\nTrue + True = {True + True}")    # 2
print(f"True * 5    = {True * 5}")         # 5
print(f"False + 1   = {False + 1}")        # 1


# ─── 5. MUTABILITY VS IMMUTABILITY ───────────────────────────────────────────
#
# Immutable: the object's value cannot change after creation
#            int, float, complex, str, tuple, bool, frozenset, bytes
#
# Mutable:   the object can be changed in place
#            list, dict, set, bytearray, and most user-defined objects
#
# Why does this matter?
#   - Immutable objects are safe to share (Python caches them freely)
#   - Mutable objects shared between names can cause unexpected changes

print("\n=== Mutability ===")

# Immutable: "changing" a string creates a new object
s = "hello"
original_id = id(s)
s = s + " world"
print(f"same object? {id(s) == original_id}")   # False — new object created

# Mutable: appending to a list modifies the SAME object
lst = [1, 2, 3]
original_id = id(lst)
lst.append(4)
print(f"same object? {id(lst) == original_id}")  # True — modified in place


# ─── 6. INT INTERNALS: ARBITRARY PRECISION ───────────────────────────────────
#
# Python's int has no fixed size — it grows as needed.
# There is no integer overflow in Python (unlike C/Java).
# Large integers use multiple C longs internally.

print("\n=== Arbitrary precision integers ===")
big = 2 ** 100
print(f"2^100 = {big}") # 2 to the power of 100 is a very large number
print(f"type:  {type(big)}")   # still just int

# sys.getsizeof shows how memory grows with magnitude
import sys
for exp in [1, 10, 100, 1000]:
    n = 2 ** exp
    print(f"  2^{exp:<4}: {sys.getsizeof(n)} bytes")


# ─── 7. FLOAT INTERNALS: IEEE 754 ────────────────────────────────────────────
#
# Python floats are C doubles: 64-bit IEEE 754
# 1 sign bit + 11 exponent bits + 52 mantissa bits
# Range: ~±1.8×10^308, but precision is limited (~15–17 significant digits)

print("\n=== Float precision ===")
print(0.1 + 0.2)           # not 0.3 — floating point representation error
print(0.1 + 0.2 == 0.3)    # False!

# Correct way to compare floats
import math
print(math.isclose(0.1 + 0.2, 0.3))   # True


# ─── EXERCISES ────────────────────────────────────────────────────────────────
#
# 1. Prove that bool is a subclass of int without using isinstance().
#    Hint: look at bool.__mro__
#
# 2. Create a large integer and a large float both representing 10^50.
#    Are they equal? What does this tell you about precision?
#
# 3. What is the result of: (1, 2, 3) + (4, 5)?
#    Does this violate immutability? Explain what actually happened.
#
# 4. Use sys.getsizeof() on an empty list, then a list with 1 element,
#    then 2, 4, 8 elements. What pattern do you notice?
#
# THOUGHT QUESTION:
#   Python strings are immutable. So how does this work efficiently?
#       result = ""
#       for i in range(10000):
#           result += str(i)
#   What is Python actually doing? Is there a better way?
