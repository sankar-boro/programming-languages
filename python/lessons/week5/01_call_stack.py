"""
WEEK 5 — DAY 1: The Call Stack
================================
Topic: How Python builds and tears down stack frames, what lives inside
       a frame, and how to inspect the call stack at runtime.

Key ideas:
  - Every function call creates a new frame object on the call stack
  - A frame holds: local variables, the code object, a pointer to the
    enclosing frame, and the current instruction pointer
  - The stack is a LIFO structure — last in, first out
  - CPython's stack depth is limited (default ~1000 frames)
  - inspect and sys give you live access to the call stack
"""

import sys
import dis
import inspect
import traceback


# ─── 1. WHAT IS THE CALL STACK? ──────────────────────────────────────────────
#
# The call stack is a sequence of active function calls.
# Each call adds a frame on top. When the function returns, its frame is removed.
#
# Stack at any point during execution looks like:
#
#   [bottom]  module-level (global frame)
#             main()
#             process()
#             helper()          ← currently executing
#   [top]
#
# "top of stack" = the currently running function
# Python only executes ONE frame at a time (the top one).

def a():
    print("  inside a() — about to call b()")
    b()
    print("  back in a()")

def b():
    print("  inside b() — about to call c()")
    c()
    print("  back in b()")

def c():
    print("  inside c() — printing call stack:")
    # Walk the current call stack:
    frame = inspect.currentframe()
    depth = 0
    while frame is not None:
        info = inspect.getframeinfo(frame)
        print(f"    [{depth}] {info.function}() at line {info.lineno}")
        frame = frame.f_back   # f_back points to the caller's frame
        depth += 1

print("=== Call stack walkthrough ===")
a()


# ─── 2. WHAT IS A FRAME OBJECT? ──────────────────────────────────────────────
#
# A frame (PyFrameObject in CPython) contains:
#
#   f_code       — the code object being executed
#   f_locals     — dict of local variables
#   f_globals    — reference to the module's global namespace
#   f_back       — the calling frame (None for the top-level frame)
#   f_lasti      — index of last attempted bytecode instruction
#   f_lineno     — current line number in source
#   f_builtins   — reference to builtins namespace
#
# Frames are Python objects — you can inspect them at runtime.

print("\n=== Frame object internals ===")

def inspect_my_frame(x, y):
    frame = inspect.currentframe()
    print(f"  function:   {frame.f_code.co_name}")
    print(f"  filename:   {frame.f_code.co_filename.split('/')[-1]}")
    print(f"  lineno:     {frame.f_lineno}")
    print(f"  locals:     {frame.f_locals}")
    print(f"  local vars: {frame.f_code.co_varnames}")
    print(f"  arg count:  {frame.f_code.co_argcount}")

inspect_my_frame(10, 20)


# ─── 3. FRAME CREATION AND DESTRUCTION ───────────────────────────────────────
#
# CALL:   Python creates a new frame, pushes it onto the stack, transfers
#         control to the new frame's code object at instruction 0.
#
# RETURN: Python pops the frame off the stack, restores the caller's frame,
#         places the return value where the caller expects it.
#
# Each frame has its OWN evaluation stack (a small operand stack) where
# bytecode instructions push/pop intermediate values.

def show_frame_lifecycle():
    def inner():
        # This frame exists only while inner() is running
        frame = inspect.currentframe()
        print(f"  inner frame id: {id(frame)}")
        return "done"

    print(f"  stack depth before call: {len(inspect.stack())}")
    result = inner()
    print(f"  stack depth after call:  {len(inspect.stack())}")
    # inner's frame is now GONE — dereferenced and eligible for GC

print("\n=== Frame lifecycle ===")
show_frame_lifecycle()


# ─── 4. THE PYTHON EVALUATION STACK ──────────────────────────────────────────
#
# Inside each frame, there is a small VALUE STACK (operand stack).
# Bytecode instructions pop operands from and push results onto this stack.
# This is how expression evaluation works — it's a stack machine.
#
# Example: a + b * c
#
#   LOAD_FAST   a        stack: [a]
#   LOAD_FAST   b        stack: [a, b]
#   LOAD_FAST   c        stack: [a, b, c]
#   BINARY_MULTIPLY      stack: [a, b*c]
#   BINARY_ADD           stack: [a + b*c]
#   RETURN_VALUE         stack: []  → value returned to caller

def expression_demo(a, b, c):
    return a + b * c

print("\n=== Bytecode operand stack ===")
dis.dis(expression_demo)


# ─── 5. STACK DEPTH LIMIT ─────────────────────────────────────────────────────
#
# CPython limits the call stack depth to prevent stack overflow.
# Default limit: 1000 frames (sys.getrecursionlimit()).
# Exceeding it raises RecursionError.

print(f"\n=== Stack limits ===")
print(f"sys.getrecursionlimit(): {sys.getrecursionlimit()}")

def measure_actual_depth():
    """Find actual available depth from current position."""
    def recurse(n):
        try:
            return recurse(n + 1)
        except RecursionError:
            return n
    return recurse(0)

actual = measure_actual_depth()
print(f"actual depth available from here: ~{actual}")

# You CAN increase the limit (use carefully — risks C stack overflow):
# sys.setrecursionlimit(5000)


# ─── 6. TRACEBACK — THE CALL STACK ON EXCEPTION ───────────────────────────────
#
# When an exception is raised, Python captures the current call stack
# and attaches it to the exception as a traceback object.
# The traceback IS a linked list of frame snapshots — most recent call last.

print("\n=== Traceback as call stack ===")

def top():
    middle()

def middle():
    bottom()

def bottom():
    raise ValueError("something went wrong")

try:
    top()
except ValueError:
    # Walk the traceback manually:
    tb = sys.exc_info()[2]   # (type, value, traceback)
    print("Traceback frames (innermost last):")
    while tb is not None:
        frame = tb.tb_frame
        lineno = tb.tb_lineno
        name = frame.f_code.co_name
        print(f"  {name}() at line {lineno}")
        tb = tb.tb_next

# Or use traceback module for formatted output:
try:
    top()
except ValueError:
    import io
    buf = io.StringIO()
    traceback.print_exc(file=buf)
    print("\nFormatted traceback:")
    print(buf.getvalue())


# ─── 7. INSPECT.STACK() — LIVE CALL STACK ─────────────────────────────────────
#
# inspect.stack() returns a list of FrameInfo objects for the current call stack.
# [0] is the current frame, [-1] is the module-level frame.
# Useful for debugging, logging, and building developer tools.

def log_with_caller():
    """A logging function that knows who called it."""
    stack = inspect.stack()
    caller = stack[1]   # [0] is log_with_caller itself, [1] is its caller
    print(f"  called from: {caller.function}() at line {caller.lineno}")
    print(f"  context: {caller.code_context[0].strip() if caller.code_context else 'N/A'}")

def some_business_logic():
    log_with_caller()

print("=== inspect.stack() caller detection ===")
some_business_logic()


# ─── 8. FRAMES AND PERFORMANCE ────────────────────────────────────────────────
#
# Frame creation is NOT free. Each call:
#   - Allocates a PyFrameObject (or uses a free-list cache in CPython)
#   - Copies argument values into the new frame's locals array
#   - Updates the stack pointer
#
# This is why Python function calls have overhead.
# CPython 3.11+ introduced "frame interning" and "zero-cost frames"
# to reduce this overhead significantly.
#
# In hot loops, inlining logic or using built-in C functions avoids
# Python frame overhead entirely.

import timeit

def add_python(a, b):
    return a + b

# Calling a pure Python function vs using the operator directly:
t_call  = timeit.timeit("add_python(1, 2)", globals=globals(), number=5_000_000)
t_inline = timeit.timeit("1 + 2",            number=5_000_000)

print(f"\n=== Frame creation overhead ===")
print(f"  Python call add(1, 2):  {t_call:.4f}s")
print(f"  Inline 1 + 2:           {t_inline:.4f}s")
print(f"  Call overhead: {(t_call - t_inline) / t_inline * 100:.0f}% slower")


# ─── EXERCISES ────────────────────────────────────────────────────────────────
#
# 1. Write a function call_depth() that returns the current call stack depth
#    without using inspect.stack(). Hint: walk f_back from currentframe().
#
# 2. Write a decorator @trace that prints the function name and arguments
#    every time a decorated function is called, and prints the return value
#    when it returns. Use inspect.currentframe() to get the call site line.
#
# 3. Write a function that intentionally causes RecursionError, catches it,
#    and reports how deep the stack actually got before crashing.
#
# 4. Use dis.dis() on this code and trace the operand stack manually:
#       def calc(x):
#           return (x + 1) * (x - 1)
#    Draw the stack state after each bytecode instruction.
#
# THOUGHT QUESTION:
#   CPython stores frames on the C stack AND the Python call stack.
#   When Python is running recursion, it uses BOTH.
#   Why might a RecursionError occur BEFORE sys.getrecursionlimit() is reached
#   when each Python function call also uses some C stack space?
