"""
WEEK 2 — DAY 2: Numbers — Deep Behavior
========================================
Topic: Python's numeric types from the inside — how int, float, complex,
       and Decimal work, what IEEE 754 means in practice, operator
       behavior, and numeric gotchas that cause real bugs.

Key ideas:
  - Python int is arbitrary precision (backed by an array of C digits)
  - Python float is a C double: IEEE 754, 64-bit — limited precision
  - Decimal is exact decimal arithmetic (fixed-point, configurable precision)
  - Type coercion in arithmetic follows a strict hierarchy
"""

import sys
import math
import decimal
import fractions
import struct


# ─── 1. INT: ARBITRARY PRECISION ─────────────────────────────────────────────
#
# Python int is NOT a fixed 32 or 64-bit integer.
# CPython stores it as an array of 30-bit digits in base 2³⁰.
# It grows as needed — no overflow, no undefined behavior.
#
# Consequence: arithmetic on large ints is slower than on small ints.
# For 64-bit range values, CPython has a fast path.

print("=== int: arbitrary precision ===")
x = 10 ** 100          # a googol — no problem
print(f"10^100 = {x}")
print(f"digits: {len(str(x))}")

# Memory grows with magnitude
for exp in [10, 50, 100, 500, 1000]:
    n = 10 ** exp
    print(f"  10^{exp:<5}: {sys.getsizeof(n)} bytes")

# int operations: all return new int objects
a, b = 17, 5
print(f"\n17 + 5  = {a + b}")     # 22
print(f"17 - 5  = {a - b}")      # 12
print(f"17 * 5  = {a * b}")      # 85
print(f"17 / 5  = {a / b}")      # 3.4   ← returns FLOAT (true division)
print(f"17 // 5 = {a // b}")     # 3     ← floor division (int result)
print(f"17 % 5  = {a % b}")      # 2     ← modulo
print(f"17 ** 5 = {a ** b}")     # 1419857

# Floor division and modulo obey: a == (a // b) * b + (a % b)
assert a == (a // b) * b + (a % b), "floor division identity broken"

# Negative modulo: Python always returns non-negative for positive divisor
print(f"\n-7 % 3  = {-7 % 3}")   # 2   (not -1 as in C)
print(f"-7 // 3 = {-7 // 3}")   # -3  (floor toward -∞, not toward 0)


# ─── 2. BIT OPERATIONS ON INT ────────────────────────────────────────────────
#
# Since int is arbitrary precision, bitwise operations work on all bits.
# Useful for flags, masks, low-level protocols, hashing.

print("\n=== Bitwise operations ===")
a, b = 0b1010, 0b1100   # 10, 12

print(f"a      = {a:08b} ({a})")
print(f"b      = {b:08b} ({b})")
print(f"a & b  = {(a & b):08b} ({a & b})")    # AND: bits set in BOTH
print(f"a | b  = {(a | b):08b} ({a | b})")    # OR:  bits set in EITHER
print(f"a ^ b  = {(a ^ b):08b} ({a ^ b})")    # XOR: bits set in ONE only
print(f"~a     = {(~a)}")                       # NOT: flips all bits (returns -(a+1))
print(f"a << 2 = {a << 2}")                    # left shift: multiply by 2^n
print(f"a >> 1 = {a >> 1}")                    # right shift: floor divide by 2^n

# Common pattern: check if a number is a power of 2
def is_power_of_two(n):
    # n & (n-1) clears the lowest set bit; if n is a power of 2, result is 0
    return n > 0 and (n & (n - 1)) == 0

for n in [1, 2, 3, 4, 8, 15, 16, 100]:
    print(f"  is_power_of_two({n:<4}): {is_power_of_two(n)}")


# ─── 3. FLOAT: IEEE 754 DOUBLE PRECISION ─────────────────────────────────────
#
# Python float is exactly a C double: 64 bits total
#
#   1 bit  — sign
#   11 bits — exponent (biased by 1023)
#   52 bits — mantissa (fraction, implicit leading 1)
#
# Representable range: ±5×10⁻³²⁴ to ±1.8×10³⁰⁸
# Significant digits:  ~15–17 decimal digits
#
# KEY FACT: Most decimal fractions cannot be represented exactly in binary.
# 0.1 in binary is 0.0001100110011... (repeating) — it's truncated.

print("\n=== Float IEEE 754 ===")

# The classic float trap:
print(f"0.1 + 0.2        = {0.1 + 0.2}")           # 0.30000000000000004
print(f"0.1 + 0.2 == 0.3 = {0.1 + 0.2 == 0.3}")   # False

# See the actual stored value
print(f"repr(0.1)        = {repr(0.1)}")            # shows full precision
print(f"0.1 stored as:   {0.1:.20f}")               # shows the approximation

# Float special values
inf     = float("inf")
neg_inf = float("-inf")
nan     = float("nan")

print(f"\ninf:       {inf}")
print(f"-inf:      {neg_inf}")
print(f"nan:       {nan}")
print(f"nan == nan: {nan == nan}")      # False! NaN is not equal to itself
print(f"math.isnan(nan): {math.isnan(nan)}")   # correct way to check

# Inspect IEEE 754 bits using struct
bits = struct.pack("d", 0.1)
print(f"\n0.1 as 8 bytes: {bits.hex()}")

# Float limits
print(f"\nsys.float_info: {sys.float_info}")


# ─── 4. COMPARING FLOATS CORRECTLY ────────────────────────────────────────────
#
# Never use == with floats.
# Use math.isclose() which checks relative and absolute tolerance.

print("\n=== Correct float comparison ===")

a = 0.1 + 0.2
b = 0.3

print(f"a == b:                   {a == b}")                # False
print(f"math.isclose(a, b):       {math.isclose(a, b)}")   # True
print(f"abs(a - b) < 1e-9:        {abs(a - b) < 1e-9}")    # True (manual)

# math.isclose defaults: rel_tol=1e-9, abs_tol=0.0
# For values near zero, use abs_tol:
x = 1e-15
print(f"\nmath.isclose(x, 0):           {math.isclose(x, 0)}")              # False
print(f"math.isclose(x, 0, abs_tol=1e-9): {math.isclose(x, 0, abs_tol=1e-9)}")  # True


# ─── 5. DECIMAL: EXACT DECIMAL ARITHMETIC ─────────────────────────────────────
#
# Use decimal.Decimal when you need exact decimal results.
# Financial calculations, tax computations — anywhere 0.1 must be exactly 0.1.
#
# Decimal is NOT float — it stores numbers as base-10 internally.

print("\n=== decimal.Decimal ===")

from decimal import Decimal, getcontext

# Default precision: 28 significant digits
a = Decimal("0.1")
b = Decimal("0.2")
print(f"Decimal: 0.1 + 0.2 = {a + b}")    # exactly 0.3
print(f"Float:   0.1 + 0.2 = {0.1 + 0.2}")  # 0.30000000000000004

# Control precision
getcontext().prec = 50
result = Decimal(1) / Decimal(3)
print(f"\n1/3 to 50 digits: {result}")

# Decimal is slower than float — use only when exactness matters
import timeit
t_float   = timeit.timeit("0.1 + 0.2", number=1_000_000)
t_decimal = timeit.timeit("Decimal('0.1') + Decimal('0.2')",
                           setup="from decimal import Decimal", number=1_000_000)
print(f"\nfloat speed:   {t_float:.3f}s")
print(f"Decimal speed: {t_decimal:.3f}s")
print(f"Decimal is {t_decimal / t_float:.1f}× slower")


# ─── 6. FRACTIONS: EXACT RATIONAL ARITHMETIC ──────────────────────────────────
#
# fractions.Fraction stores numerator and denominator as ints.
# Exact for all rational numbers. Slower and more memory-heavy.

print("\n=== fractions.Fraction ===")

from fractions import Fraction

a = Fraction(1, 3)   # exactly 1/3
b = Fraction(1, 6)   # exactly 1/6
print(f"1/3 + 1/6 = {a + b}")   # 1/2 — exact

# Can construct from float (reveals the hidden approximation):
print(f"Fraction(0.1) = {Fraction(0.1)}")  # not 1/10!
print(f"Fraction('0.1') = {Fraction('0.1')}")  # exactly 1/10 (string form)


# ─── 7. TYPE COERCION IN ARITHMETIC ──────────────────────────────────────────
#
# Python numeric type hierarchy (from narrow to wide):
#   bool → int → float → complex
#
# When you mix types, Python promotes to the widest type.

print("\n=== Type promotion ===")
print(f"True + 1     = {True + 1},  type: {type(True + 1).__name__}")   # int
print(f"1 + 1.0      = {1 + 1.0},  type: {type(1 + 1.0).__name__}")     # float
print(f"1.0 + (1+0j) = {1.0 + (1+0j)}, type: {type(1.0 + (1+0j)).__name__}")  # complex

# Division always produces float:
print(f"4 / 2        = {4 / 2},   type: {type(4 / 2).__name__}")         # 2.0 float
print(f"4 // 2       = {4 // 2},   type: {type(4 // 2).__name__}")       # 2   int


# ─── 8. MATH MODULE ESSENTIALS ────────────────────────────────────────────────

print("\n=== math module ===")
print(f"math.pi       = {math.pi}")
print(f"math.e        = {math.e}")
print(f"math.tau      = {math.tau}")         # 2π
print(f"math.sqrt(2)  = {math.sqrt(2)}")
print(f"math.log(100, 10) = {math.log(100, 10)}")   # log base 10
print(f"math.log2(8)  = {math.log2(8)}")
print(f"math.ceil(3.2)  = {math.ceil(3.2)}")
print(f"math.floor(3.9) = {math.floor(3.9)}")
print(f"math.trunc(3.9) = {math.trunc(3.9)}")  # toward zero


# ─── EXERCISES ────────────────────────────────────────────────────────────────
#
# 1. Compute 2^1000. How many decimal digits does it have?
#    Now compute it modulo 97. Which is faster — computing then taking mod,
#    or using Python's built-in pow(2, 1000, 97)?
#    Time both with timeit.
#
# 2. Show that 0.1, 0.2, and 0.3 are all approximations.
#    Print each to 30 decimal places. What do you observe?
#
# 3. Use Fraction to prove that:
#       1/2 + 1/3 + 1/6 == 1  (exactly)
#    Then do the same with floats. Do you get exactly 1.0?
#
# 4. Write a function safe_divide(a, b) that:
#    - Returns a / b as a Fraction when both inputs are int
#    - Returns float('inf') when b == 0 and a > 0
#    - Returns float('nan') when both a and b are 0
#
# THOUGHT QUESTION:
#   Python's // (floor division) always rounds toward -∞, not toward 0.
#   Why might this be more mathematically consistent than C-style truncation?
#   Consider what -7 // 2 gives in Python vs C, and how modulo relates.
