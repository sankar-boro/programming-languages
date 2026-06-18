"""
WEEK 7 — DAY 2: Lexical Scoping — Deep Dive
=============================================
Topic: Why Python uses lexical (static) scoping, how the LEGB rule is
       implemented at the bytecode level, scope edge cases, class scope
       surprises, and the difference between dynamic and lexical scoping.

Key ideas:
  - Lexical scoping: scope is determined by WHERE code is WRITTEN, not called
  - Dynamic scoping: scope is determined by WHERE code is CALLED from
  - Python is lexically scoped — the compiler resolves names statically
  - Class bodies are NOT a scope for nested functions (the biggest surprise)
  - global/nonlocal are declarations, not assignments — they affect compilation
"""

import dis


# ─── 1. LEXICAL VS DYNAMIC SCOPING ───────────────────────────────────────────
#
# Lexical (static) scoping — Python's model:
#   A function's free variables are resolved in the scope where the
#   function was DEFINED, regardless of where it is called from.
#
# Dynamic scoping — NOT Python's model (used in early Lisps, some shells):
#   A function's free variables are resolved in the scope of the CALLER.
#
# The difference matters when a function is passed between scopes:

x = "global"

def get_x():
    return x   # x is free — lexical: finds "global" where get_x was defined

def caller():
    x = "caller_local"   # this does NOT affect get_x under lexical scoping
    return get_x()

print("=== Lexical scoping ===")
print(f"caller(): {caller()}")   # "global" — get_x sees the definition scope
# Under dynamic scoping (hypothetical), caller() would return "caller_local"


# ─── 2. HOW THE COMPILER CLASSIFIES VARIABLES ────────────────────────────────
#
# At COMPILE TIME (not runtime), Python's compiler analyzes each function and
# classifies every variable name as one of:
#
#   local:    assigned anywhere in the function → FAST slot
#   free:     referenced but not assigned → DEREF (from closure cell)
#   global:   declared with `global` → GLOBAL dict lookup
#   cell:     referenced by an inner function → DEREF (provides cell to inner)
#   builtin:  everything else → BUILTIN dict lookup
#
# This classification happens BEFORE execution. The opcodes encode it.

def analyze(x):
    # x = FAST (local — it's a parameter, which counts as local)
    y = x + 1           # y = FAST (local — assigned here)
    # z = FAST (local — assigned below)
    z = y * 2
    return z

print("\n=== Variable classification ===")
print(f"co_varnames (locals): {analyze.__code__.co_varnames}")
print(f"co_freevars (free):   {analyze.__code__.co_freevars}")

# A function with free variables:
def outer():
    captured = 100

    def inner():
        # captured is FREE — not assigned in inner, assigned in outer
        return captured * 2

    print(f"\n  inner co_freevars:  {inner.__code__.co_freevars}")
    print(f"  outer co_cellvars:  {outer.__code__.co_cellvars}")
    return inner

inner = outer()


# ─── 3. THE ASSIGNMENT RULE — MOST COMMON SURPRISE ───────────────────────────
#
# IF a name is ASSIGNED ANYWHERE in a function body (without global/nonlocal),
# Python marks it as LOCAL for the ENTIRE function body — even before the assignment.
#
# This catches people off guard when they try to use a global variable
# and then reassign it in the same function.

x = 10

def gotcha():
    print(x)    # you expect to print the global x = 10
    x = 20      # but this assignment makes x LOCAL for the whole function
                 # so the print above reads a local x that doesn't exist yet
                 # → UnboundLocalError

try:
    gotcha()
except UnboundLocalError as e:
    print(f"\n=== The assignment rule ===")
    print(f"  UnboundLocalError: {e}")

# Fix options:
def fixed_with_global():
    global x    # declare x refers to the global
    print(x)
    x = 20

def fixed_no_assignment():
    print(x)    # just read — no assignment, no problem
    y = 20      # assign a DIFFERENT name


# ─── 4. BYTECODE PROOF OF STATIC CLASSIFICATION ──────────────────────────────

def show_opcodes():
    g = 42   # local variable

    def reads_local():
        return g    # g is free var here

    def reads_global():
        return len   # len is built-in

    print("=== reads_local (LOAD_DEREF for free variable) ===")
    dis.dis(reads_local)

    print("\n=== reads_global (LOAD_GLOBAL for built-in) ===")
    dis.dis(reads_global)

show_opcodes()


# ─── 5. CLASS SCOPE — THE BIG SURPRISE ───────────────────────────────────────
#
# This is the most confusing scope rule in Python:
#
# A class body creates a NAMESPACE, but NOT a scope for nested functions.
# Methods inside a class CANNOT see class-level names as free variables.
# They must access them through `self` or `ClassName.attribute`.
#
# Why? Because class bodies are executed once as a special dict-building block.
# The class namespace is NOT in the enclosing scope of methods —
# it's a temporary dict that becomes `Class.__dict__` after the class is built.

class Counter:
    count = 0       # class attribute

    def increment(self):
        # This WON'T work:
        # count += 1   ← count is not in local, enclosing, global, or builtin
        # Must use self.count or Counter.count:
        Counter.count += 1
        return Counter.count

    # Demonstration of the class scope not being visible to lambdas:
    multiplier = 3
    # This does NOT work:
    # values = [x * multiplier for x in range(5)]   # NameError: multiplier
    # Must use explicit reference:
    values = [x * 3 for x in range(5)]   # literal, not variable

c = Counter()
print(f"\n=== Class scope ===")
print(f"  increment: {c.increment()}, {c.increment()}, {c.increment()}")
print(f"  values:    {Counter.values}")

# The subtlety: class-level names are visible INSIDE the class body
# but only for direct expressions — NOT inside nested function/lambda/comprehension
class Demo:
    x = 10
    y = x + 5     # works: x is directly visible in the class body statement
    # z = [x for _ in range(3)]  # NameError in Python 3 — comprehension has own scope

print(f"\n  Demo.y: {Demo.y}")


# ─── 6. COMPREHENSION SCOPE (PYTHON 3) ───────────────────────────────────────
#
# In Python 3, list/dict/set comprehensions and generator expressions
# have their OWN implicit scope (they are compiled to their own code objects).
# The iteration variable is LOCAL to the comprehension.
# This differs from Python 2 where comprehension vars leaked into the enclosing scope.

x = "outer"
result = [x for x in range(5)]   # this x is LOCAL to the comprehension
print(f"\n=== Comprehension scope ===")
print(f"  x after comprehension: {x!r}")   # "outer" — not overwritten in Python 3
print(f"  result: {result}")

# Proof: the comprehension's x is a separate local:
import dis
code_str = "[x for x in range(5)]"
compiled = compile(code_str, "<demo>", "eval")
print(f"\n  comprehension creates its own code object:")
for const in compiled.co_consts:
    if hasattr(const, "co_name"):
        print(f"    {const.co_name}: varnames={const.co_varnames}")


# ─── 7. NESTED SCOPES AND SHADOWING IN PRACTICE ──────────────────────────────

# LEGB chain — innermost always wins:
x = "global"

def level1():
    x = "level1"

    def level2():
        x = "level2"

        def level3():
            print(f"    level3 sees: {x}")   # "level2" — closest enclosing

        def level3_no_local():
            print(f"    level3 (no local) sees: {x}")   # "level2"

        level3()
        level3_no_local()

    level2()

print(f"\n=== Nested scope chain ===")
print(f"  global x: {x!r}")
level1()

# Shadowing built-ins — a common mistake:
def shadow_len():
    len = lambda x: "I broke len!"   # shadows the built-in
    print(f"  len([1,2,3]) inside: {len([1,2,3])}")
    # The real len is gone for this scope — access via builtins module:
    import builtins
    print(f"  builtins.len:        {builtins.len([1,2,3])}")

print(f"\n=== Shadowing built-in ===")
shadow_len()
print(f"  len([1,2,3]) outside: {len([1,2,3])}")   # restored outside


# ─── 8. SCOPE AND EXEC/EVAL ───────────────────────────────────────────────────
#
# exec() and eval() accept optional globals and locals dicts.
# This lets you create custom namespaces — useful for sandboxing and DSLs.
# The scope lookup inside exec'd code follows the passed dicts.

print(f"\n=== Custom scope with exec ===")

custom_globals = {
    "x": 100,
    "__builtins__": {"print": print, "range": range}   # restricted builtins
}

exec("y = x * 2\nprint(y)", custom_globals)
print(f"  custom_globals after exec: y = {custom_globals.get('y')}")

# eval with restricted scope:
result = eval("x + 1", {"x": 5, "__builtins__": {}})
print(f"  eval('x+1', x=5): {result}")


# ─── EXERCISES ────────────────────────────────────────────────────────────────
#
# 1. Verify the assignment rule with dis.dis():
#    Write two nearly identical functions:
#       def f1(): print(x)               # reads global
#       def f2(): print(x); x = 0       # UnboundLocalError
#    Use dis.dis() on both. Find the opcode difference (LOAD_GLOBAL vs LOAD_FAST).
#    This shows that the compiler decides at compile time, not runtime.
#
# 2. Create a function where the same name is: global in one function,
#    local in a sibling function, and free in a nested function of the second.
#    Use dis.dis() to confirm the opcodes for each.
#
# 3. Demonstrate the class scope comprehension issue:
#    Write a class where you try to use a class variable inside a comprehension.
#    Show the NameError. Then show THREE ways to fix it:
#    a) Use a default argument in the comprehension
#    b) Use a regular for loop instead
#    c) Define the list outside the class and assign it
#
# 4. Write a function find_scope(name, frame) that, given a variable name
#    and a frame object, determines whether the name is local, free,
#    global, or builtin in that frame. Use frame.f_code attributes.
#
# THOUGHT QUESTION:
#   Python's lexical scoping means a function "locks in" its enclosing scope
#   at DEFINITION time. This makes code easier to reason about — you can
#   understand a function just by reading its definition, without knowing
#   where it will be called.
#   But dynamic scoping has uses: a called function could pick up configuration
#   from its caller without explicit parameter passing.
#   What Python feature partially achieves dynamic-scope-like behavior
#   while remaining safe? (Hint: think about context variables — contextvars)
