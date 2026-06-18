"""
WEEK 2 — DAY 1: Strings — Deep Behavior
========================================
Topic: Python strings from the ground up — Unicode internals, immutability
       consequences, every operation and what it costs, and common pitfalls.

Key ideas:
  - str is an immutable sequence of Unicode code points (not bytes)
  - CPython uses 3 internal encodings depending on content (Latin-1, UCS-2, UCS-4)
  - String operations almost always create new objects
  - String interning is an optimization, not a guarantee
"""

import sys
import timeit


# ─── 1. STRINGS ARE UNICODE, NOT BYTES ───────────────────────────────────────
#
# In Python 3, a str is a sequence of Unicode code points.
# A code point is an integer assigned to every character in the Unicode standard.
# Unicode has 1,114,112 possible code points (U+0000 to U+10FFFF).
#
# This is different from bytes (b"..."), which are raw binary data.

s = "hello"
print(f"String:     {s}")
print(f"Type:       {type(s)}")
print(f"Length:     {len(s)}")          # number of code points, not bytes

# Access individual code points by index
print(f"First char: {s[0]}")            # "h"
print(f"Code point: {ord(s[0])}")       # 104 — the integer for 'h'
print(f"Back again: {chr(104)}")        # "h" — integer → character

# Non-ASCII characters are just code points too
emoji = "🐍"
print(f"\nEmoji:      {emoji}")
print(f"Code point: {ord(emoji)}")      # 128013
print(f"Length:     {len(emoji)}")      # 1 — ONE code point
print(f"Byte size:  {sys.getsizeof(emoji)} bytes in memory")


# ─── 2. CPYTHON'S INTERNAL STRING ENCODING ───────────────────────────────────
#
# CPython uses one of three internal representations (PEP 393, "compact strings"):
#
#   Latin-1 (1 byte/char):   when all code points fit in U+0000–U+00FF
#   UCS-2   (2 bytes/char):  when any code point is U+0100–U+FFFF
#   UCS-4   (4 bytes/char):  when any code point is U+10000–U+10FFFF (emojis etc.)
#
# Python picks the smallest representation that fits all characters.
# This saves memory for ASCII-heavy text.

ascii_str   = "hello"
ucs2_str    = "héllo"      # é is U+00E9 — but é < U+0100, still Latin-1
ucs4_str    = "hello🐍"    # emoji forces UCS-4

print("\n=== Internal memory per character ===")
# Rough bytes per character (excluding object overhead)
for s, label in [(ascii_str, "ASCII"), (ucs2_str, "Latin-1+"), (ucs4_str, "UCS-4")]:
    size = sys.getsizeof(s)
    per_char = (size - 49) / len(s) if len(s) else 0   # 49 = overhead estimate
    print(f"  {label:<10}: {size} bytes total, ~{per_char:.1f} bytes/char")


# ─── 3. IMMUTABILITY — WHAT IT REALLY MEANS ──────────────────────────────────
#
# You cannot change a character in a string:
#   s[0] = "H"   ← TypeError
#
# Every string "operation" creates a NEW string object.
# The original is unchanged.

print("\n=== String immutability ===")

s = "hello"
original_id = id(s)

s_upper = s.upper()
print(f"s:           {s}        id: {id(s)}")
print(f"s.upper():   {s_upper}  id: {id(s_upper)}")
print(f"Same object? {id(s) == id(s_upper)}")   # False

# += concatenation on strings: rebinds the name, creates a new object
s2 = "hello"
id_before = id(s2)
s2 += " world"
print(f"\ns2 after +=: id changed? {id(s2) != id_before}")   # True — new object


# ─── 4. STRING CONCATENATION COST ────────────────────────────────────────────
#
# Repeated += in a loop is O(n²) because each step copies the full string.
# "hello" + "world" → allocate n+m bytes, copy both → discard old
#
# The right approach for building strings: join()

def slow_concat(n):
    """O(n²): each iteration copies the growing string."""
    result = ""
    for i in range(n):
        result += str(i)
    return result

def fast_join(n):
    """O(n): collect parts, join once."""
    parts = []
    for i in range(n):
        parts.append(str(i))
    return "".join(parts)

# Also valid with a generator expression (most Pythonic):
def fast_genexpr(n):
    return "".join(str(i) for i in range(n))

print("\n=== Concatenation performance (n=10000) ===")
n = 10_000
t_slow = timeit.timeit(lambda: slow_concat(n), number=10)
t_fast = timeit.timeit(lambda: fast_join(n), number=10)
print(f"  += loop: {t_slow:.4f}s")
print(f"  join():  {t_fast:.4f}s")
print(f"  speedup: {t_slow / t_fast:.1f}×")


# ─── 5. STRING METHODS — THE IMPORTANT ONES ──────────────────────────────────
#
# All string methods return NEW strings. None modify in place.

s = "  Hello, World!  "

print(f"\n=== String methods ===")
print(f"strip():       '{s.strip()}'")          # remove whitespace both ends
print(f"lstrip():      '{s.lstrip()}'")         # left only
print(f"rstrip():      '{s.rstrip()}'")         # right only
print(f"lower():       '{s.strip().lower()}'")
print(f"upper():       '{s.strip().upper()}'")
print(f"title():       '{s.strip().title()}'")

words = "apple,banana,cherry"
print(f"\nsplit(','):    {words.split(',')}")    # returns a list
print(f"split limit:   {words.split(',', 1)}")  # stop after 1 split

parts = ["one", "two", "three"]
print(f"join:          {'|'.join(parts)}")       # "one|two|three"

sentence = "the cat sat on the mat"
print(f"\nreplace:       '{sentence.replace('at', 'og')}'")
print(f"count('at'):   {sentence.count('at')}")
print(f"find('cat'):   {sentence.find('cat')}")    # index of first match
print(f"find('dog'):   {sentence.find('dog')}")    # -1 if not found


# ─── 6. STRING SLICING ────────────────────────────────────────────────────────
#
# s[start:stop:step]
# All indices are optional. Negative indices count from the end.
# Slicing always returns a new string.

s = "Python"
print(f"\n=== Slicing ===")
print(f"s[0:3]:    {s[0:3]}")      # "Pyt"
print(f"s[2:]:     {s[2:]}")       # "thon"
print(f"s[:3]:     {s[:3]}")       # "Pyt"
print(f"s[-3:]:    {s[-3:]}")      # "hon"  (last 3)
print(f"s[::2]:    {s[::2]}")      # "Pto"  (every 2nd)
print(f"s[::-1]:   {s[::-1]}")     # "nohtyP" (reversed)


# ─── 7. F-STRINGS VS OTHER FORMATTING ────────────────────────────────────────
#
# f-strings (Python 3.6+) are compiled to efficient bytecode.
# They are faster than % formatting and .format() for most cases.

name = "Python"
version = 3.12

# All equivalent:
s1 = "I love %s %.2f" % (name, version)         # old style — avoid
s2 = "I love {} {}".format(name, version)        # str.format — verbose
s3 = f"I love {name} {version:.2f}"              # f-string — preferred

print(f"\n=== Formatting ===")
print(s1)
print(s2)
print(s3)

# f-strings can contain expressions
print(f"  2 + 2 = {2 + 2}")
print(f"  upper = {name.upper()}")
print(f"  repr  = {name!r}")       # !r applies repr()


# ─── 8. BYTES VS STR ─────────────────────────────────────────────────────────
#
# str  = sequence of Unicode code points (text)
# bytes = sequence of raw byte values 0–255 (binary data)
#
# You must encode to go from str → bytes, decode to go from bytes → str.
# The encoding specifies HOW to represent code points as bytes.

print("\n=== str vs bytes ===")

s = "café"
b = s.encode("utf-8")       # str → bytes using UTF-8 encoding
print(f"str:         {s!r}")
print(f"utf-8 bytes: {b!r}")   # b'caf\xc3\xa9' — é takes 2 bytes in UTF-8
print(f"back to str: {b.decode('utf-8')!r}")

# Wrong encoding → error or garbage
try:
    b.decode("ascii")
except UnicodeDecodeError as e:
    print(f"UnicodeDecodeError: {e}")


# ─── EXERCISES ────────────────────────────────────────────────────────────────
#
# 1. Find the Unicode code point for: €, 中, 🦊
#    Use ord(). What does this tell you about "length" of a string?
#
# 2. Time these three ways to build a string of 1..1000 joined by "-":
#      a) s = ""; for i in ...: s += str(i) + "-"
#      b) "-".join(str(i) for i in range(1000))
#      c) "-".join(map(str, range(1000)))
#    Use timeit. Which is fastest and why?
#
# 3. Decode b'\xe4\xb8\xad\xe6\x96\x87' as UTF-8. What string do you get?
#    Now try decoding it as Latin-1. What happens?
#
# 4. Write a function is_palindrome(s) that returns True if s reads the
#    same forwards and backwards (ignore case, ignore spaces).
#    Use only string methods and slicing — no loops.
#
# THOUGHT QUESTION:
#   If Python strings are immutable, how does this line work without copying
#   the entire string on every iteration?
#       for char in "hello world":
#           print(char)
#   What is Python actually doing when it iterates over a string?
