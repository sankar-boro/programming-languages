"""
WEEK 3 — DAY 3: Control Flow — while Loops
============================================
Topic: When and how to use while loops, the internal execution model,
       common patterns, infinite loops, and when while beats for.

Key ideas:
  - while is for condition-driven repetition (not collection traversal)
  - for iterates known sequences; while handles unknown/dynamic termination
  - Python has no do-while — but the while True: ... break pattern is idiomatic
  - Sentinel patterns and state machines are natural fits for while
"""

import dis
import time
import itertools


# ─── 1. HOW WHILE WORKS INTERNALLY ───────────────────────────────────────────
#
# while condition:
#     body
#
# Bytecode:
#   L_start:
#     <evaluate condition>
#     POP_JUMP_IF_FALSE L_end
#     <body>
#     JUMP_ABSOLUTE L_start
#   L_end:
#
# The loop re-evaluates the condition at the TOP of every iteration.
# If the condition is False on the first check, the body never runs.

def show_while_bytecode():
    def countdown(n):
        while n > 0:
            n -= 1

    print("=== while loop bytecode ===")
    dis.dis(countdown)

show_while_bytecode()


# ─── 2. WHILE VS FOR ─────────────────────────────────────────────────────────
#
# Use FOR when:   you know the collection or the count upfront
# Use WHILE when: the termination condition depends on state that changes
#                 inside the loop — user input, convergence, events, reading

# FOR is cleaner for known sequences:
for i in range(5):
    pass   # O(1) loop variable update, clear intent

# WHILE is right for unknown termination:
import random
random.seed(42)

attempts = 0
while True:
    roll = random.randint(1, 6)
    attempts += 1
    if roll == 6:
        break

print(f"Rolled 6 after {attempts} attempts")


# ─── 3. WHILE WITH A CONDITION ────────────────────────────────────────────────

def binary_search(sorted_list, target):
    """
    Classic algorithm that needs while: we don't know how many iterations
    before finding the element or narrowing to empty.
    """
    low, high = 0, len(sorted_list) - 1

    while low <= high:
        mid = (low + high) // 2
        if sorted_list[mid] == target:
            return mid
        elif sorted_list[mid] < target:
            low = mid + 1
        else:
            high = mid - 1

    return -1   # not found

data = list(range(0, 100, 2))   # [0, 2, 4, ..., 98]
print(f"\nbinary_search(data, 64): index {binary_search(data, 64)}")
print(f"binary_search(data, 63): index {binary_search(data, 63)}")


# ─── 4. WHILE TRUE / BREAK — THE DO-WHILE PATTERN ───────────────────────────
#
# Python has no do-while loop. The idiomatic replacement:
#
#   while True:
#       <body always runs at least once>
#       if exit_condition:
#           break
#
# This guarantees the body runs at least once — unlike while condition: which
# might never run if the condition starts False.

print("\n=== while True / break pattern ===")

def get_positive_number(prompt_values):
    """Simulate a loop that must run at least once (like a do-while)."""
    values = iter(prompt_values)
    while True:
        value = next(values, None)
        print(f"  got: {value}")
        if value is not None and value > 0:
            return value
        print("  invalid, retry")

result = get_positive_number([-1, 0, -5, 7])
print(f"  accepted: {result}")


# ─── 5. WHILE WITH ELSE ───────────────────────────────────────────────────────
#
# Like for/else: the else block runs only if the while condition
# became False naturally — NOT if the loop exited via break.

def find_factor(n):
    """Find the smallest factor of n greater than 1."""
    divisor = 2
    while divisor * divisor <= n:
        if n % divisor == 0:
            print(f"  {n} has factor {divisor}")
            break
        divisor += 1
    else:
        print(f"  {n} is prime")

print("\n=== while / else ===")
for n in [10, 13, 17, 25]:
    find_factor(n)


# ─── 6. SENTINEL PATTERN ─────────────────────────────────────────────────────
#
# A sentinel is a special value that signals "stop."
# Classic use: reading lines until an empty one, processing until None, etc.

print("\n=== Sentinel pattern ===")

def process_until_sentinel(stream, sentinel=None):
    """Process items from stream until sentinel value is seen."""
    results = []
    for item in stream:
        if item == sentinel:
            break
        results.append(item * 2)
    return results

data_stream = [1, 2, 3, None, 4, 5]   # None is the sentinel
result = process_until_sentinel(data_stream)
print(f"processed: {result}")

# walrus operator (:=) makes sentinel patterns clean in while:
# while (line := file.readline()) != "":
#     process(line)


# ─── 7. CONVERGENCE LOOPS ─────────────────────────────────────────────────────
#
# Numerical algorithms run until a value is "close enough."
# For loops are wrong here — we don't know the iteration count ahead of time.

def newton_sqrt(n, tolerance=1e-10):
    """Newton's method: iteratively improve a guess until convergence."""
    guess = n / 2.0
    iterations = 0

    while True:
        improved = (guess + n / guess) / 2
        if abs(improved - guess) < tolerance:
            return improved, iterations
        guess = improved
        iterations += 1

import math
for n in [2, 9, 144, 0.5]:
    result, iters = newton_sqrt(n)
    print(f"  sqrt({n}) ≈ {result:.10f} (actual: {math.sqrt(n):.10f}) in {iters} iters")


# ─── 8. STATE MACHINES ────────────────────────────────────────────────────────
#
# State machines model systems that move between discrete states.
# while + if/elif on a state variable is the classic implementation.

print("\n=== Simple state machine: traffic light ===")

def traffic_light_simulation(cycles):
    """Simulate a traffic light through N full cycles."""
    states = ["GREEN", "YELLOW", "RED"]
    durations = {"GREEN": 3, "YELLOW": 1, "RED": 2}   # seconds (simulated)
    transitions = {"GREEN": "YELLOW", "YELLOW": "RED", "RED": "GREEN"}

    state = "GREEN"
    completed_cycles = 0
    steps = []

    while completed_cycles < cycles:
        steps.append(state)
        state = transitions[state]
        if state == "GREEN":
            completed_cycles += 1

    return steps

log = traffic_light_simulation(2)
print(" → ".join(log))


# ─── 9. COMMON WHILE BUGS ─────────────────────────────────────────────────────

print("\n=== Common bugs ===")

# BUG 1: Forgetting to update the loop variable (infinite loop)
# n = 10
# while n > 0:
#     print(n)
#     # forgot: n -= 1   ← this would loop forever

# BUG 2: Off-by-one in condition
# while n >= 0:  vs  while n > 0:
# Test with n=0 — which one should handle zero?

n = 5
count = 0
while n > 0:          # n=0 exits — "positive integers only"
    count += 1
    n -= 1
print(f"  iterations (n > 0):  {count}")   # 5

n = 5
count = 0
while n >= 0:         # n=0 also runs — includes zero
    count += 1
    n -= 1
print(f"  iterations (n >= 0): {count}")   # 6

# BUG 3: Mutating the collection while iterating with while
# items = [1, 2, 3, 4, 5]
# i = 0
# while i < len(items):      # len() changes as you remove items!
#     if items[i] % 2 == 0:
#         items.remove(items[i])   # shifts indices — may skip elements
#     else:
#         i += 1
# BETTER: filter with comprehension


# ─── EXERCISES ────────────────────────────────────────────────────────────────
#
# 1. Implement the Collatz conjecture using while:
#    - Start with any positive integer n
#    - If n is even: n = n // 2
#    - If n is odd:  n = 3*n + 1
#    - Stop when n == 1
#    Count and return the number of steps.
#    Which starting value under 100 takes the most steps?
#
# 2. Implement a simple lexer using a state machine:
#    States: START, IN_NUMBER, IN_WORD
#    Input: a string like "abc 123 def 456"
#    Output: list of tokens like [("WORD","abc"), ("NUM","123"), ...]
#
# 3. Write a while loop that implements the Euclidean algorithm:
#       gcd(a, b): while b != 0: a, b = b, a % b; return a
#    Verify it matches math.gcd(). Trace the steps for gcd(48, 18).
#
# 4. Simulate reading lines from a "file" (a list of strings ending with "")
#    using the walrus operator:
#       lines = ["hello", "world", "", "never reached"]
#       i = 0
#       while (line := lines[i]) != "":
#           process it
#           i += 1
#
# THOUGHT QUESTION:
#   A for loop in Python cannot be "paused" mid-iteration and resumed later.
#   But a generator can. How does the iterator protocol make generators
#   possible? What does `yield` actually save and restore?
#   (Hint: think about what state a for loop has at each step.)
