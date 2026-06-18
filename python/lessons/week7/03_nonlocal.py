"""
WEEK 7 — DAY 3: nonlocal — Mutating Enclosing Scope
=====================================================
Topic: What `nonlocal` does at the bytecode level, when you need it,
       patterns it enables, and its relationship to closures and cells.

Key ideas:
  - `nonlocal` is a COMPILE-TIME declaration — it changes how a name is classified
  - Without nonlocal: assignment in inner function creates a new local
  - With nonlocal: assignment writes to the enclosing function's cell
  - nonlocal walks the scope chain (not just one level up)
  - nonlocal and global have the same relationship: both override the default local classification
"""

import dis
import sys


# ─── 1. WHY NONLOCAL EXISTS ───────────────────────────────────────────────────
#
# The assignment rule: any name assigned in a function is LOCAL to that function.
# Problem: what if you want to assign to a variable in an ENCLOSING function?
#
#   `global` handles module-level variables.
#   `nonlocal` handles enclosing-function variables.
#
# Before nonlocal (Python 2): developers used mutable containers as a workaround:
#   count = [0]   # list used as a mutable box
#   def inc(): count[0] += 1   # mutate the list, not rebind the name

def demo_without_nonlocal():
    """What happens WITHOUT nonlocal — a new local is created silently."""
    count = 0

    def increment():
        count = count + 1   # assignment → count is LOCAL to increment
                            # but count is read before assignment → UnboundLocalError
        return count

    try:
        increment()
    except UnboundLocalError as e:
        print(f"  Without nonlocal: {e}")

def demo_with_nonlocal():
    """nonlocal makes `count` refer to the enclosing scope's cell."""
    count = 0

    def increment():
        nonlocal count   # count now refers to the enclosing cell
        count += 1       # STORE_DEREF — writes to the cell
        return count

    return increment()

print("=== Why nonlocal ===")
demo_without_nonlocal()
print(f"  With nonlocal: {demo_with_nonlocal()}")


# ─── 2. BYTECODE: STORE_DEREF ─────────────────────────────────────────────────
#
# nonlocal changes the bytecode emitted for the name:
#   Without nonlocal: assignment → STORE_FAST (local slot)
#   With nonlocal:    assignment → STORE_DEREF (cell)

def outer_no_nonlocal():
    x = 0
    def inner():
        x = 1   # STORE_FAST — creates new local
    return inner

def outer_with_nonlocal():
    x = 0
    def inner():
        nonlocal x
        x = 1   # STORE_DEREF — writes to cell
    return inner

print("\n=== Bytecode comparison ===")
print("--- WITHOUT nonlocal (STORE_FAST) ---")
dis.dis(outer_no_nonlocal())

print("\n--- WITH nonlocal (STORE_DEREF) ---")
dis.dis(outer_with_nonlocal())


# ─── 3. NONLOCAL WALKS THE SCOPE CHAIN ───────────────────────────────────────
#
# nonlocal searches the enclosing scopes from inner to outer,
# stopping at the FIRST scope that has a cell/local with that name.
# It skips the module global scope — that's `global`'s job.

def level_1():
    x = "level1"

    def level_2():
        # x is NOT assigned here — so it's a free variable in level_2

        def level_3():
            nonlocal x      # walks up: level_2 has no x, level_1 has x → found
            x = "modified by level_3"

        level_3()
        print(f"  level_2 sees x: {x!r}")

    level_2()
    print(f"  level_1 sees x: {x!r}")

print("\n=== nonlocal walks chain ===")
level_1()


# ─── 4. NONLOCAL VS GLOBAL ────────────────────────────────────────────────────
#
# global:   skips ALL enclosing function scopes, targets the MODULE global
# nonlocal: targets the NEAREST enclosing scope that defines the variable
#
# They cannot both apply to the same name in the same scope.

MODULE_VAR = "module"

def outer():
    outer_var = "outer"

    def inner():
        global MODULE_VAR   # targets module level
        nonlocal outer_var  # targets outer()

        MODULE_VAR = "modified by inner (global)"
        outer_var  = "modified by inner (nonlocal)"

    inner()
    print(f"  outer_var after inner(): {outer_var!r}")

print("\n=== global vs nonlocal ===")
print(f"  MODULE_VAR before: {MODULE_VAR!r}")
outer()
print(f"  MODULE_VAR after:  {MODULE_VAR!r}")


# ─── 5. PATTERNS ENABLED BY NONLOCAL ─────────────────────────────────────────

print("\n=== Practical nonlocal patterns ===")

# Pattern 1: Accumulator / running total
def make_accumulator():
    total = 0
    def add(value):
        nonlocal total
        total += value
        return total
    return add

acc = make_accumulator()
print(f"  acc: {acc(10)}, {acc(5)}, {acc(20)}")   # 10, 15, 35

# Pattern 2: Toggle
def make_toggle(initial=False):
    state = initial
    def toggle():
        nonlocal state
        state = not state
        return state
    return toggle

toggle = make_toggle()
print(f"  toggle: {toggle()}, {toggle()}, {toggle()}")   # True, False, True

# Pattern 3: Once — a function that runs only on first call
def once(func):
    """Decorator: function runs only the first time it's called."""
    has_run = False
    result  = None

    def wrapper(*args, **kwargs):
        nonlocal has_run, result
        if not has_run:
            result  = func(*args, **kwargs)
            has_run = True
        return result

    return wrapper

@once
def expensive_init():
    print("  (initializing — expensive!)")
    return 42

print(f"\n  first call:  {expensive_init()}")
print(f"  second call: {expensive_init()}")   # no re-execution
print(f"  third call:  {expensive_init()}")

# Pattern 4: Retry logic with attempt tracking
def make_retry(max_attempts=3):
    attempts = 0
    def retry(func, *args):
        nonlocal attempts
        while attempts < max_attempts:
            attempts += 1
            try:
                return func(*args)
            except Exception as e:
                print(f"  attempt {attempts} failed: {e}")
        raise RuntimeError(f"Failed after {max_attempts} attempts")
    return retry

import random
random.seed(42)

retry = make_retry(max_attempts=4)

def flaky():
    if random.random() < 0.6:
        raise ValueError("random failure")
    return "success"

try:
    print(f"\n  retry result: {retry(flaky)}")
except RuntimeError as e:
    print(f"\n  {e}")


# ─── 6. NONLOCAL WITH MULTIPLE VARIABLES ──────────────────────────────────────

def make_bidirectional_counter(start=0):
    """Counter that can increment AND decrement — needs nonlocal for count."""
    count = start

    def up():
        nonlocal count
        count += 1
        return count

    def down():
        nonlocal count
        count -= 1
        return count

    def value():
        return count   # read-only — no nonlocal needed

    return up, down, value

up, down, val = make_bidirectional_counter(10)
print(f"\n=== Bidirectional counter ===")
print(f"  up:   {up()}, {up()}, {up()}")    # 11, 12, 13
print(f"  down: {down()}, {down()}")         # 12, 11
print(f"  val:  {val()}")                    # 11


# ─── 7. WHAT NONLOCAL CANNOT DO ───────────────────────────────────────────────
#
# 1. nonlocal cannot create a new variable in the enclosing scope.
#    The variable must ALREADY EXIST in an enclosing function scope.
#
# 2. nonlocal cannot target the module global scope — use `global` for that.
#
# 3. nonlocal cannot skip scopes — it binds to the NEAREST match.

y = "global y"

def bad_nonlocal():
    def inner():
        try:
            nonlocal y   # SyntaxError: no binding for y in enclosing scope
        except SyntaxError:
            pass
    inner()

# This is a SyntaxError caught at compile time, not runtime.
# Let's verify with a string + compile():
code = """
def outer():
    def inner():
        nonlocal missing_var   # not defined in outer
        missing_var = 1
    inner()
"""
try:
    compile(code, "<demo>", "exec")
except SyntaxError as e:
    print(f"\n=== nonlocal requires existing variable ===")
    print(f"  SyntaxError: {e}")


# ─── 8. THE MUTABLE-CONTAINER WORKAROUND (HISTORICAL) ─────────────────────────
#
# Before Python 3.0 introduced `nonlocal`, developers used mutable containers
# (lists or dicts) as a workaround. Understanding this helps when reading
# older code or Python 2 codebases.

print("\n=== Python 2 era workaround (mutable container) ===")

def old_style_counter():
    count = [0]   # list as a mutable box — not a rebinding, a mutation

    def increment():
        count[0] += 1   # mutates the list — no nonlocal needed
        return count[0]

    return increment

old_inc = old_style_counter()
print(f"  {old_inc()}, {old_inc()}, {old_inc()}")   # 1, 2, 3

# Why this works: count is never REBOUND (the name always refers to the same list).
# count[0] = ... is a mutation (STORE_SUBSCR), not a rebinding (STORE_FAST).
# So Python doesn't classify count as a local — it's a free variable (LOAD_DEREF).


# ─── EXERCISES ────────────────────────────────────────────────────────────────
#
# 1. Use dis.dis() to confirm that nonlocal changes STORE_FAST to STORE_DEREF.
#    Write two identical functions — one with nonlocal, one without.
#    Disassemble the inner function in each case.
#
# 2. Write make_history_counter() — a counter that also remembers every
#    value it has held. Returns (increment, decrement, get_history).
#    Use nonlocal for count AND history list (or explain why you don't
#    need nonlocal for the list mutation).
#
# 3. Implement make_lru_cache(maxsize) using closures and nonlocal.
#    Use an OrderedDict (from collections) as the cache.
#    On cache hit: move to most-recently-used end.
#    On cache miss + full: evict least-recently-used.
#
# 4. Explain, with code examples, the difference between these two patterns:
#       # Pattern A: rebinding (needs nonlocal)
#       def f():
#           x = 0
#           def g(): nonlocal x; x = x + 1
#
#       # Pattern B: mutation (no nonlocal needed)
#       def f():
#           x = [0]
#           def g(): x[0] = x[0] + 1
#    When would you choose A over B? What does each communicate to a reader?
#
# THOUGHT QUESTION:
#   `nonlocal` only works with ENCLOSING FUNCTION scopes — not class scopes,
#   not module scopes (use `global` for module scope).
#   Why do you think class bodies were excluded from nonlocal's scope?
#   (Hint: class bodies execute once and become __dict__ — they're not persistent frames)
