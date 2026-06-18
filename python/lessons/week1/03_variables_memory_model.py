"""
WEEK 1 — DAY 3: Variables and the Memory Model
===============================================
Topic: How Python manages memory, what namespaces are, how assignment
       really works, and how CPython's allocator and GC operate.

Key ideas:
  - Assignment binds a name to an object in a namespace
  - Python has multiple layers of namespaces (LEGB rule)
  - CPython uses reference counting + cyclic garbage collection
  - del removes the name binding, not the object
"""

import sys
import gc
import ctypes


# ─── 1. WHAT ASSIGNMENT ACTUALLY DOES ────────────────────────────────────────
#
# x = 42  does NOT:
#   - Create a box called x and put 42 in it
#
# x = 42  DOES:
#   1. Create an int object with value 42 in heap memory
#   2. Bind the name "x" to that object in the current namespace
#   3. Increment the object's reference count by 1

x = 42
print(f"id(x) = {id(x)}")   # memory address of the int object

# Multiple assignment forms — all bind names to the same object mechanisms:
a = b = c = []          # all three names point to ONE list
print(f"a is b is c: {a is b is c}")   # True

a.append(1)
print(f"b after a.append: {b}")        # [1] — same object!


# ─── 2. AUGMENTED ASSIGNMENT AND MUTABILITY ───────────────────────────────────
#
# += behaves differently depending on mutability.
# For immutable types: creates a new object, rebinds the name
# For mutable types:   modifies in place (calls __iadd__)

print("\n=== += on immutable (int) ===")
x = 10
print(f"before: id={id(x)}")
x += 1
print(f"after:  id={id(x)}")   # different — new object

print("\n=== += on mutable (list) ===")
lst = [1, 2]
print(f"before: id={id(lst)}")
lst += [3]                     # calls lst.__iadd__([3]) → modifies in place
print(f"after:  id={id(lst)}")  # SAME id — extended in place

# This subtlety causes real bugs when lists are shared between names:
a = [1]
b = a
a += [2]            # modifies in place — b sees it
print(f"b = {b}")   # [1, 2]


# ─── 3. NAMESPACES ───────────────────────────────────────────────────────────
#
# A namespace is a mapping from names to objects (implemented as a dict).
# Python has multiple namespaces, searched in LEGB order:
#
#   L — Local:    names inside the current function
#   E — Enclosing: names in any enclosing function scopes (closures)
#   G — Global:   names at the module (file) level
#   B — Built-in: names in the builtins module (print, len, range, ...)
#
# When you use a name, Python searches L → E → G → B and uses the first match.

MODULE_LEVEL = "global"     # lives in the module's global namespace

def outer():
    enclosing = "enclosing"

    def inner():
        local = "local"
        # inner can see: local (L), enclosing (E), MODULE_LEVEL (G), print (B)
        print(local, enclosing, MODULE_LEVEL)

    inner()

outer()

# You can inspect namespaces directly:
print("\n=== Namespace inspection ===")
print(f"type(globals()): {type(globals())}")   # dict
print(f"MODULE_LEVEL in globals(): {'MODULE_LEVEL' in globals()}")


# ─── 4. del — REMOVES THE NAME, NOT THE OBJECT ───────────────────────────────
#
# del x does NOT destroy the object.
# It removes the name binding from the namespace.
# The object is freed only when its reference count drops to 0.

print("\n=== del semantics ===")

x = [1, 2, 3]
y = x               # y also points to the list

del x               # removes name x from namespace; list refcount: 1→... still 1

# x is gone:
try:
    print(x)
except NameError as e:
    print(f"NameError: {e}")

print(f"y still works: {y}")   # [1, 2, 3] — object still alive through y


# ─── 5. REFERENCE COUNTING IN DETAIL ─────────────────────────────────────────
#
# Every PyObject has a ob_refcnt field.
# It increments when a new name binds to the object.
# It decrements when a name is rebound, deleted, or goes out of scope.
# When ob_refcnt reaches 0, the object is immediately freed.
#
# sys.getrefcount(x) returns the count + 1 (the function call itself adds one)

print("\n=== Reference counting ===")

obj = [1, 2, 3]
print(f"refcount (just obj):    {sys.getrefcount(obj)}")   # 2 (obj + call)

alias = obj
print(f"refcount (obj + alias): {sys.getrefcount(obj)}")   # 3

stored = [obj]
print(f"refcount (+ in list):   {sys.getrefcount(obj)}")   # 4

del stored
del alias
print(f"refcount (back to obj): {sys.getrefcount(obj)}")   # 2


# ─── 6. CYCLIC GARBAGE COLLECTION ────────────────────────────────────────────
#
# Reference counting cannot free circular references:
#
#   a = []
#   a.append(a)   ← a refers to itself, refcount never drops to 0
#
# CPython has a cyclic GC (the gc module) that runs periodically
# and finds groups of objects that only reference each other.
# It uses a generational algorithm (0, 1, 2) — most objects die young.

print("\n=== Cyclic GC ===")

gc.disable()        # turn off automatic collection to demonstrate manually

a = []
a.append(a)         # circular reference
del a               # refcount doesn't hit 0 — object is NOT freed yet

collected = gc.collect()   # force a full collection cycle
print(f"objects collected: {collected}")   # at least 1

gc.enable()


# ─── 7. STACK VS HEAP ────────────────────────────────────────────────────────
#
# In CPython:
#
#   Stack: Python's call stack (frames) — managed automatically
#          Local variables live in frames, not the C stack
#
#   Heap:  All Python objects live here — managed by CPython's allocator
#          CPython has its own memory pool (pymalloc) on top of malloc
#          for objects < 512 bytes, to avoid fragmentation
#
# This is different from C/C++ where stack allocation is explicit.
# In Python you don't control where an object is allocated — always heap.


# ─── 8. INTERNING: WHEN PYTHON REUSES OBJECTS ────────────────────────────────
#
# Python interns (caches and reuses) certain objects for performance:
#
#   - Integers from -5 to 256
#   - Short strings that look like identifiers (alphanumeric, no spaces)
#   - None, True, False (singletons — only one ever exists)
#
# Interned objects: `is` returns True even for separately created values

print("\n=== Object interning ===")

# Integer interning
a, b = 256, 256
print(f"256 is 256: {a is b}")    # True — cached

a, b = 257, 257
print(f"257 is 257: {a is b}")    # False — not cached (in separate statements)

# String interning
s1 = "hello"
s2 = "hello"
print(f"'hello' is 'hello': {s1 is s2}")    # True — interned (identifier-like)

s1 = "hello world"
s2 = "hello world"
print(f"'hello world' is: {s1 is s2}")      # may be True or False — implementation detail

# None is always a singleton
print(f"None is None: {None is None}")      # always True


# ─── EXERCISES ────────────────────────────────────────────────────────────────
#
# 1. Create a list. Create two aliases to it. Delete both aliases.
#    Does the list disappear? How do you know?
#    Hint: track id() and use sys.getrefcount() before/after each step.
#
# 2. Create a circular reference:
#       a = {}
#       b = {}
#       a["b"] = b
#       b["a"] = a
#    After del a and del b, use gc.collect() to free them.
#    How many objects were collected?
#
# 3. Explore globals() and locals() inside a function.
#    Are they the same dict or different? What does that tell you?
#
# 4. Predict the output before running:
#       x = [1, 2]
#       y = x
#       x = x + [3]    # note: + creates a new list, += does not
#       print(y)
#
# THOUGHT QUESTION:
#   Python says "everything is an object." If a function is an object,
#   it must have an id(). What is stored at that memory address?
#   What attributes does a function object have? (Try dir(some_function))
