"""
WEEK 5 — DAY 2: Execution Frames — How Python Runs Your Code
=============================================================
Topic: The CPython interpreter loop, code objects, bytecode in depth,
       how frames execute instructions, and the role of the GIL.

Key ideas:
  - CPython's ceval.c contains the main interpreter loop (a giant switch)
  - Code objects are immutable, shareable; frames are mutable, per-call
  - The interpreter dispatches on opcode — one instruction per loop tick
  - The GIL ensures only one thread runs Python bytecode at a time
  - Python 3.11+ uses "specializing adaptive interpreter" for speedups
"""

import dis
import sys
import code
import types
import opcode


# ─── 1. CODE OBJECTS VS FRAME OBJECTS ────────────────────────────────────────
#
# Code object (PyCodeObject):
#   - Created ONCE per function definition
#   - Immutable — shared across all calls to the function
#   - Contains: bytecode, constants, variable names, metadata
#   - Lives as long as the function object lives
#
# Frame object (PyFrameObject):
#   - Created on EVERY function call
#   - Mutable — holds the actual runtime state for one call
#   - Contains: locals dict, reference to code object, instruction pointer
#   - Destroyed when the function returns

def demonstrate():
    x = 10
    y = 20
    return x + y

# One code object — shared
code_obj = demonstrate.__code__
print("=== Code object (shared, immutable) ===")
print(f"  type:          {type(code_obj)}")
print(f"  co_name:       {code_obj.co_name}")
print(f"  co_varnames:   {code_obj.co_varnames}")
print(f"  co_consts:     {code_obj.co_consts}")
print(f"  co_stacksize:  {code_obj.co_stacksize}")   # max operand stack depth needed
print(f"  co_flags:      {code_obj.co_flags:#010b}")  # bit flags (generator, coroutine, etc.)
print(f"  co_nlocals:    {code_obj.co_nlocals}")


# ─── 2. READING BYTECODE ─────────────────────────────────────────────────────
#
# co_code is the raw bytecode as bytes.
# Each instruction is 2 bytes: 1 byte opcode + 1 byte argument.
# dis.get_instructions() decodes them into human-readable form.

print("\n=== Bytecode instructions ===")
for instr in dis.get_instructions(demonstrate):
    print(f"  {instr.offset:>3}  {instr.opname:<25} {instr.argval!r}")

# Raw bytes:
print(f"\n  raw co_code: {demonstrate.__code__.co_code.hex()}")


# ─── 3. OPCODE CATEGORIES ────────────────────────────────────────────────────
#
# Python bytecode has ~150 opcodes (varies by version). Main categories:
#
#   LOAD_*        — push a value onto the operand stack
#   STORE_*       — pop a value from the stack, store it somewhere
#   BINARY_*      — pop two values, compute, push result
#   UNARY_*       — pop one value, compute, push result
#   CALL_*        — call a function (multiple variants since 3.11)
#   JUMP_*        — conditional/unconditional jumps (control flow)
#   BUILD_*       — build a list, tuple, dict, set, slice from stack items
#   RETURN_VALUE  — pop top of stack, return to caller

print("\n=== Opcode categories (sample) ===")

def example(a, b):
    lst = [a, b]
    return lst[0] + lst[1]

dis.dis(example)


# ─── 4. THE INTERPRETER LOOP (SIMPLIFIED) ────────────────────────────────────
#
# CPython's main loop in ceval.c (simplified pseudocode):
#
#   while True:
#       opcode, arg = fetch_next_instruction(frame)
#       switch(opcode):
#           case LOAD_FAST:
#               push(frame.locals[arg])
#           case STORE_FAST:
#               frame.locals[arg] = pop()
#           case BINARY_ADD:
#               right = pop()
#               left  = pop()
#               push(left + right)   # calls left.__add__(right)
#           case RETURN_VALUE:
#               retval = pop()
#               restore_previous_frame()
#               return retval
#           ...
#
# In Python 3.11+, the loop uses "computed gotos" (on supported compilers)
# for a direct jump table instead of a C switch — faster dispatch.

# You can see this with a simple addition:
def add_two(a, b):
    return a + b

print("\n=== add_two bytecode (shows interpreter loop steps) ===")
dis.dis(add_two)
# LOAD_FAST 'a'   → push a onto operand stack
# LOAD_FAST 'b'   → push b onto operand stack
# BINARY_OP +     → pop both, call a.__add__(b), push result
# RETURN_VALUE    → pop result, return to caller


# ─── 5. HOW VARIABLES ARE STORED IN FRAMES ───────────────────────────────────
#
# CPython uses THREE different mechanisms for variable storage, each with a
# different LOAD/STORE opcode pair:
#
#   LOAD_FAST / STORE_FAST   — LOCAL variables (array indexed by position)
#   LOAD_GLOBAL / STORE_GLOBAL — MODULE-level globals (dict lookup)
#   LOAD_DEREF / STORE_DEREF  — FREE variables (closure cells)
#
# FAST is fastest — direct array index into the frame's fastlocals array.
# GLOBAL is slower — dict lookup in the module namespace.
# DEREF is for closed-over variables — indirection through a cell object.

def show_variable_opcodes():
    global_var = "I'm in module scope"   # stored with STORE_GLOBAL

    def outer():
        enclosing = "I'm in enclosing"   # will be DEREF'd by inner

        def inner():
            local = "I'm local"          # STORE_FAST / LOAD_FAST
            print(local)                 # LOAD_FAST
            print(enclosing)             # LOAD_DEREF

        print("\n--- inner() bytecode ---")
        dis.dis(inner)
        inner()

    outer()

show_variable_opcodes()


# ─── 6. THE GIL — GLOBAL INTERPRETER LOCK ────────────────────────────────────
#
# CPython has a GIL: a mutex that allows only ONE thread to execute
# Python bytecode at a time.
#
# Why does the GIL exist?
#   - CPython's memory management (reference counting) is NOT thread-safe
#   - Without the GIL, two threads could simultaneously decrement the same
#     refcount, causing use-after-free or double-free bugs
#   - The GIL makes the CPython runtime simpler and safer
#
# What the GIL means in practice:
#   - CPU-bound Python code does NOT run truly in parallel (threads)
#   - I/O-bound code CAN overlap (GIL is released during I/O waits)
#   - True CPU parallelism in Python → use multiprocessing (separate processes)
#   - Python 3.13 introduced "free-threaded" mode (PEP 703) — GIL optional

import threading
import time

counter = 0

def increment_counter(n):
    global counter
    for _ in range(n):
        counter += 1   # NOT atomic — race condition possible without GIL

# With GIL: threads work correctly for simple int operations in practice
# (though counter += 1 is not technically atomic in all cases)
threads = [threading.Thread(target=increment_counter, args=(100_000,)) for _ in range(4)]
counter = 0
for t in threads: t.start()
for t in threads: t.join()
print(f"\n=== GIL demo ===")
print(f"  expected: 400000, got: {counter}")
# Often correct because GIL switches happen between bytecode instructions


# ─── 7. CO_FLAGS — FRAME TYPE DETECTION ──────────────────────────────────────
#
# The co_flags bitfield tells CPython what kind of function this is.
# Key flags:
#   0x04  CO_VARARGS       — has *args
#   0x08  CO_VARKEYWORDS   — has **kwargs
#   0x20  CO_GENERATOR     — is a generator (uses yield)
#   0x100 CO_COROUTINE     — is an async def function

CO_VARARGS     = 0x04
CO_VARKEYWORDS = 0x08
CO_GENERATOR   = 0x20
CO_COROUTINE   = 0x100

def check_flags(func):
    flags = func.__code__.co_flags
    print(f"  {func.__name__}:")
    print(f"    has *args:    {bool(flags & CO_VARARGS)}")
    print(f"    has **kwargs: {bool(flags & CO_VARKEYWORDS)}")
    print(f"    is generator: {bool(flags & CO_GENERATOR)}")

def normal(a, b): pass
def variadic(*args, **kwargs): pass
def gen_func(): yield 1

print("\n=== co_flags ===")
check_flags(normal)
check_flags(variadic)
check_flags(gen_func)


# ─── 8. CREATING AND EXECUTING CODE OBJECTS MANUALLY ─────────────────────────
#
# Code objects can be created from strings with compile().
# This is how exec/eval and REPLs work internally.
# Understanding this reveals what Python is doing when it "runs" your code.

print("\n=== compile() and exec() ===")

source = """
x = 10
y = 20
result = x + y
print(f"x + y = {result}")
"""

code_obj = compile(source, filename="<demo>", mode="exec")
print(f"compiled code type: {type(code_obj)}")
print(f"instructions: {len(list(dis.get_instructions(code_obj)))}")

namespace = {}
exec(code_obj, namespace)
print(f"namespace after exec: { {k: v for k, v in namespace.items() if not k.startswith('__')} }")

# eval() works on expressions, returns the value:
expr = compile("2 ** 10 + len('hello')", "<expr>", "eval")
print(f"eval result: {eval(expr)}")


# ─── EXERCISES ────────────────────────────────────────────────────────────────
#
# 1. Write a function bytecode_summary(func) that prints:
#    - Number of bytecode instructions
#    - Number of unique opcodes used
#    - Whether the function uses any LOAD_GLOBAL (accesses globals)
#    Test it on: a pure function, a function using global vars, a generator.
#
# 2. Use compile() to dynamically create a function from a string template:
#       code = compile("lambda x: x ** 2 + " + str(offset), ...)
#    Create 5 different "shifted square" functions with offset 0–4.
#
# 3. Use co_flags to write a function is_generator(func) that returns True
#    if the function is a generator without calling it.
#    Verify with: def g(): yield 1  vs  def f(): return 1
#
# 4. Profile the frame creation overhead more precisely:
#    Compare calling a 1-line Python function vs accessing a dict value
#    vs calling a built-in (len). Use timeit with n=1_000_000.
#    Rank them by speed and explain WHY each is faster/slower.
#
# THOUGHT QUESTION:
#   Python 3.11 introduced "specializing adaptive interpreter" (PEP 659).
#   When the interpreter notices the same opcode is called with the same
#   types repeatedly, it replaces the generic opcode with a specialized one
#   (e.g., BINARY_ADD → BINARY_ADD_INT_INT).
#   How is this similar to JIT compilation? How is it different?
