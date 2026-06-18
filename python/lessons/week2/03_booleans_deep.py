"""
WEEK 2 — DAY 3: Booleans — Deep Behavior
==========================================
Topic: Python's bool type, truthiness, short-circuit evaluation,
       comparison operators, and identity vs equality.

Key ideas:
  - bool is a subclass of int: True == 1, False == 0
  - Every Python object has a truth value (truthiness)
  - and/or do NOT return booleans — they return one of their operands
  - is checks identity (same object), == checks equality (same value)
"""


# ─── 1. BOOL IS A SUBCLASS OF INT ────────────────────────────────────────────
#
# bool inherits from int. There are exactly two bool instances: True and False.
# They are singletons — Python never creates a second True or False object.
#
# True is the integer 1. False is the integer 0. Not approximately — exactly.

print("=== bool is int ===")
print(f"True == 1:       {True == 1}")       # True
print(f"False == 0:      {False == 0}")      # True
print(f"True + True:     {True + True}")     # 2
print(f"True * 42:       {True * 42}")       # 42
print(f"False * 42:      {False * 42}")      # 0
print(f"sum([True, True, False, True]): {sum([True, True, False, True])}")  # 3

# isinstance confirms inheritance
print(f"\nisinstance(True, int):  {isinstance(True, int)}")   # True
print(f"isinstance(True, bool): {isinstance(True, bool)}")   # True
print(f"type(True) is int:      {type(True) is int}")        # False (exact check)
print(f"type(True) is bool:     {type(True) is bool}")       # True

# MRO (method resolution order) shows the hierarchy:
print(f"bool.__mro__: {bool.__mro__}")   # bool → int → object


# ─── 2. TRUTHINESS — EVERY OBJECT HAS A TRUTH VALUE ──────────────────────────
#
# Python evaluates any object as True or False in a boolean context.
# This is called "truthiness" or "bool coercion."
#
# An object is FALSY if any of these hold:
#   - It IS False or 0 or 0.0 or 0j
#   - It IS None
#   - It is an EMPTY container: "", [], (), {}, set(), frozenset()
#   - Its __bool__() returns False, or __len__() returns 0
#
# Everything else is TRUTHY.

print("\n=== Truthiness ===")

falsy_values = [False, 0, 0.0, 0j, None, "", [], (), {}, set(), frozenset()]
truthy_values = [True, 1, -1, 0.1, "a", [0], (None,), {"k": 0}, {0}]

print("Falsy:")
for v in falsy_values:
    print(f"  bool({v!r:<15}) = {bool(v)}")

print("Truthy (sample):")
for v in truthy_values[:5]:
    print(f"  bool({v!r:<15}) = {bool(v)}")


# ─── 3. HOW __bool__ AND __len__ WORK ────────────────────────────────────────
#
# Python calls bool(obj) by:
#   1. Trying obj.__bool__()  → if it returns NotImplemented, go to step 2
#   2. Trying obj.__len__()   → True if nonzero, False if zero
#   3. If neither defined:    → True (all objects are truthy by default)

class AlwaysFalse:
    def __bool__(self):
        return False

class EmptyLike:
    def __len__(self):
        return 0   # no __bool__, but __len__ returns 0

class Weird:
    pass           # no __bool__ or __len__

print("\n=== Custom truthiness ===")
print(f"AlwaysFalse(): {bool(AlwaysFalse())}")   # False
print(f"EmptyLike():   {bool(EmptyLike())}")      # False
print(f"Weird():       {bool(Weird())}")           # True (default)


# ─── 4. COMPARISON OPERATORS ─────────────────────────────────────────────────
#
# Comparisons return bool. Python supports chained comparisons.
# Chaining is AND-logic, evaluated left to right (each value checked once).

print("\n=== Comparison operators ===")

a, b, c = 1, 2, 3
print(f"a < b:        {a < b}")
print(f"a <= b:       {a <= b}")
print(f"a == b:       {a == b}")
print(f"a != b:       {a != b}")

# Chained comparisons — Python-specific, very readable
print(f"\n1 < 2 < 3:    {1 < 2 < 3}")         # True — both must hold
print(f"1 < 2 > 3:    {1 < 2 > 3}")           # False — 2 > 3 fails
print(f"1 < 3 > 2:    {1 < 3 > 2}")           # True

# Chaining is NOT the same as writing two separate comparisons with `and`:
# 1 < 2 < 3   is   (1 < 2) and (2 < 3)   — 2 is evaluated ONCE
# This matters when the middle value has side effects.


# ─── 5. IS VS == ─────────────────────────────────────────────────────────────
#
# ==  → equality: calls obj.__eq__(other)  → checks VALUE
# is  → identity: checks if both names point to the SAME object in memory
#
# Use == for value comparison (almost always what you want).
# Use is only for: None checks, singleton checks (True/False).

print("\n=== is vs == ===")

a = [1, 2, 3]
b = [1, 2, 3]   # same value, different object

print(f"a == b:  {a == b}")    # True  — same value
print(f"a is b:  {a is b}")    # False — different objects in memory

# Correct None check:
value = None
print(f"\nvalue is None:  {value is None}")    # correct
print(f"value == None:  {value == None}")      # works but not idiomatic

# Dangerous: some objects define __eq__ to do unexpected things.
# is guarantees you're checking the actual object, not any overriding logic.


# ─── 6. AND, OR, NOT — THEY DON'T RETURN BOOL ───────────────────────────────
#
# This is one of Python's most misunderstood behaviors.
#
# `and` returns the FIRST falsy operand, or the LAST operand if all are truthy.
# `or`  returns the FIRST truthy operand, or the LAST operand if all are falsy.
# `not` ALWAYS returns a bool.
#
# This is called "short-circuit evaluation."

print("\n=== and / or return operands, not booleans ===")

# and: returns first falsy, or last if all truthy
print(f"1 and 2:        {1 and 2}")          # 2     (both truthy → last)
print(f"0 and 2:        {0 and 2}")          # 0     (first falsy)
print(f"'' and 'hello': {'' and 'hello'}")   # ''    (first falsy)
print(f"1 and 2 and 3:  {1 and 2 and 3}")   # 3     (all truthy → last)

# or: returns first truthy, or last if all falsy
print(f"\n1 or 2:         {1 or 2}")          # 1     (first truthy)
print(f"0 or 2:         {0 or 2}")            # 2     (skip 0, return 2)
print(f"0 or '' or []:  {0 or '' or []}")     # []    (all falsy → last)
print(f"0 or 'hi':      {0 or 'hi'}")         # 'hi'  (first truthy)

# not always returns bool
print(f"\nnot 0:          {not 0}")            # True
print(f"not 'hello':    {not 'hello'}")       # False
print(f"type(not 0):    {type(not 0)}")       # bool


# ─── 7. PRACTICAL PATTERNS USING AND/OR ──────────────────────────────────────

print("\n=== Practical and/or patterns ===")

# Default value pattern (Python 2 era — prefer walrus or if in modern code):
name = ""
display = name or "Anonymous"
print(f"display name: {display}")    # "Anonymous"

# Conditional expression (ternary) — cleaner than or-trick:
name = ""
display = name if name else "Anonymous"
print(f"display name: {display}")    # "Anonymous"

# Guard pattern: only call if truthy
def process(data):
    return data.upper()

data = None
result = data and process(data)      # short-circuits, never calls process
print(f"result: {result}")           # None — safe, no AttributeError


# ─── 8. SHORT-CIRCUIT EVALUATION — PERFORMANCE & CORRECTNESS ─────────────────
#
# and/or evaluate LEFT to RIGHT and STOP as soon as the result is determined.
# This means right-hand side expressions may NEVER execute.

def expensive():
    print("  (expensive called)")
    return True

print("\n=== Short-circuit: right side may not run ===")

# False and <anything> → False immediately, right side skipped
result = False and expensive()
print(f"False and expensive(): {result}")     # expensive() not called

# True or <anything> → True immediately, right side skipped
result = True or expensive()
print(f"True or expensive():   {result}")     # expensive() not called

# True and <anything> → must evaluate right side
result = True and expensive()
print(f"True and expensive():  {result}")     # expensive() IS called


# ─── EXERCISES ────────────────────────────────────────────────────────────────
#
# 1. Predict the output of each before running:
#       a) print([] or {} or 0 or "fallback")
#       b) print([] and {} and 0 and "fallback")
#       c) print(not not [])
#       d) print(bool(0.0), bool(-0.0), bool(0j))
#
# 2. Write a function all_truthy(*args) that returns True only if ALL
#    arguments are truthy. Implement it without using the built-in all().
#    Then write any_truthy(*args) without using any().
#
# 3. Create a class Temperature that is falsy when the temperature is
#    below 0 (freezing). Implement __bool__. Test it in an if statement.
#
# 4. Why is `is` dangerous to use for integers?
#    Try: a = 1000; b = 1000; print(a is b)
#    Then: print(1000 is 1000)
#    Explain the difference.
#
# THOUGHT QUESTION:
#   Python's `and` and `or` return operands rather than True/False.
#   What are the advantages and disadvantages of this design choice?
#   Can you think of a case where this behavior causes a hidden bug?
