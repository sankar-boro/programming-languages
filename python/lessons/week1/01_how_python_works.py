"""
WEEK 1 — DAY 1: How Python Works
=================================
Topic: The Python execution pipeline, CPython internals, and the object model.

Key ideas:
  - Python is NOT directly compiled to machine code
  - Every value is an object with type, identity, and value
  - Variables are names (labels), not memory boxes
"""

import dis
import sys


# ─── 1. THE EXECUTION PIPELINE ───────────────────────────────────────────────
#
# When you run: python script.py
#
#   1. CPython reads your source code as text
#   2. Tokenizes it (breaks into keywords, names, operators, literals)
#   3. Parses tokens into an Abstract Syntax Tree (AST)
#   4. Compiles the AST into bytecode (.pyc files in __pycache__)
#   5. The Python Virtual Machine (PVM) executes the bytecode
#
# You never see steps 1–4. They happen invisibly on every run.

def show_bytecode():
    """dis.dis() reveals the bytecode Python actually executes."""
    def add(a, b):
        return a + b

    print("=== Bytecode for add(a, b): ===")
    dis.dis(add)
    # You'll see instructions like LOAD_FAST, BINARY_ADD, RETURN_VALUE
    # These are the real operations the PVM runs — not your Python source

show_bytecode()


# ─── 2. EVERYTHING IS AN OBJECT ──────────────────────────────────────────────
#
# Every value in Python is an object. An object has three things:
#
#   identity  → where it lives in memory (id())
#   type      → what kind of thing it is (type())
#   value     → the data it holds
#
# This is true for integers, strings, functions, classes — everything.

print("\n=== Every value is an object ===")

x = 42
print(f"value:    {x}")
print(f"type:     {type(x)}")        # <class 'int'>
print(f"identity: {id(x)}")          # memory address (CPython)

# Even a function is an object
def greet():
    pass

print(f"\nfunction type:     {type(greet)}")    # <class 'function'>
print(f"function identity: {id(greet)}")        # it lives in memory too


# ─── 3. VARIABLES ARE LABELS, NOT BOXES ──────────────────────────────────────
#
# Wrong mental model:  [ x ] → holds 42 inside it
# Correct model:       x ──→ (object: int, value=42, id=140...)
#
# A variable is a name that references an object.
# Multiple names can reference the same object.

print("\n=== Variables as labels ===")

a = [1, 2, 3]
b = a               # b is another label for the SAME list object

print(f"id(a): {id(a)}")
print(f"id(b): {id(b)}")
print(f"a is b: {a is b}")   # True — same object in memory

b.append(4)
print(f"a after b.append(4): {a}")   # [1, 2, 3, 4] — a sees the change


# ─── 4. REBINDING A NAME ─────────────────────────────────────────────────────
#
# When you do x = 20 after x = 10:
#   - A new int object (20) is created
#   - x now points to the new object
#   - The old object (10) loses one reference
#   - If nothing else points to it, CPython's reference counter drops to 0
#   - The garbage collector frees that memory
#
# The old object isn't destroyed immediately in all cases — CPython
# caches small integers (-5 to 256) as a performance optimization.

print("\n=== Rebinding names ===")

x = 10
print(f"x = 10  → id: {id(x)}")

x = 20
print(f"x = 20  → id: {id(x)}")    # different id — different object

# Small integer caching: CPython pre-creates objects for -5 to 256
a = 5
b = 5
print(f"\nSmall int (5):  a is b = {a is b}")   # True — same cached object

a = 1000
b = 1000
print(f"Large int (1000): a is b = {a is b}")   # False — two separate objects


# ─── 5. REFERENCE COUNTING ───────────────────────────────────────────────────
#
# CPython tracks how many names point to each object.
# When the count hits 0, memory is freed immediately.
# sys.getrefcount() shows the count (adds 1 for the call itself).

print("\n=== Reference counting ===")

name = "python"
print(f"refcount of 'python': {sys.getrefcount(name)}")

alias = name
print(f"refcount after alias:  {sys.getrefcount(name)}")   # one higher

del alias
print(f"refcount after del:    {sys.getrefcount(name)}")   # back down


# ─── EXERCISES ────────────────────────────────────────────────────────────────
#
# 1. Run dis.dis() on a function that does: x = a * b + c
#    Read each bytecode instruction. What does LOAD_FAST do?
#
# 2. Create two variables pointing at the same list.
#    Delete one with `del`. Does the list disappear? Why or why not?
#
# 3. Check id() for the integer 256 vs 257 created in two separate variables.
#    What do you observe? What does it tell you about CPython's internals?
#
# THOUGHT QUESTION:
#   If Python uses reference counting, what happens with circular references?
#   (e.g., a = []; a.append(a))
#   Can reference counting alone free that memory?
