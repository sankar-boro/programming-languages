"""
WEEK 8 — DAY 2: *args and **kwargs — Advanced Patterns
========================================================
Topic: Beyond the basics — how *args/**kwargs enable generic wrappers,
       decorator factories, argument forwarding, and the full grammar
       of Python's call and definition syntax.

Key ideas:
  - *args and **kwargs enable functions that work with any signature
  - They are the foundation of decorators, middleware, and generic wrappers
  - Argument forwarding (*args, **kwargs) is the key pattern for wrapping
  - The call-site * and ** unpack sequences/dicts into positional/keyword args
  - Combining / * and *args/**kwargs gives full control over call style
"""

import functools
import inspect
import timeit


# ─── 1. THE FULL ARGUMENT GRAMMAR — RECAP ─────────────────────────────────────
#
# def f(pos_only, /, standard, *args, kw_only, **kwargs):
#
# In a CALL:
#   f(1, 2, 3, 4, kw_only=5, extra=6)
#     └─────┘  └─┘           └──────┘
#   positional  *args      **kwargs
#
# In a DEFINITION (what gets bound):
#   pos_only=1, standard=2, args=(3,4), kw_only=5, kwargs={"extra":6}

def full_signature(pos_only, /, standard, *args, kw_only, **kwargs):
    print(f"  pos_only:  {pos_only}")
    print(f"  standard:  {standard}")
    print(f"  args:      {args}")
    print(f"  kw_only:   {kw_only}")
    print(f"  kwargs:    {kwargs}")

print("=== Full argument grammar ===")
full_signature(1, 2, 3, 4, kw_only=5, extra=6, another=7)


# ─── 2. ARGUMENT FORWARDING — THE WRAPPER PATTERN ─────────────────────────────
#
# The most important use of *args/**kwargs: forwarding arguments unchanged.
# This is how decorators, middleware, and proxy functions work.
# The wrapper doesn't need to know what arguments the wrapped function takes.

def log_calls(func):
    """Decorator: log every call with its arguments."""
    @functools.wraps(func)
    def wrapper(*args, **kwargs):
        # Capture args/kwargs without knowing the function's signature:
        arg_str = ", ".join([repr(a) for a in args] +
                            [f"{k}={v!r}" for k, v in kwargs.items()])
        print(f"  CALL: {func.__name__}({arg_str})")
        result = func(*args, **kwargs)   # forward everything unchanged
        print(f"  RETURN: {result!r}")
        return result
    return wrapper

@log_calls
def add(a, b):
    return a + b

@log_calls
def greet(name, *, greeting="Hello"):
    return f"{greeting}, {name}!"

print("\n=== Argument forwarding ===")
add(3, 4)
greet("Alice", greeting="Hi")


# ─── 3. DECORATOR FACTORIES — *args/**kwargs IN LAYERS ────────────────────────
#
# A decorator factory is a function that RETURNS a decorator.
# It adds a layer: factory(*config) → decorator(func) → wrapper(*args, **kwargs)
# Each layer uses *args/**kwargs differently.

def retry(max_attempts=3, exceptions=(Exception,)):
    """Decorator factory: retry on failure up to max_attempts times."""
    def decorator(func):
        @functools.wraps(func)
        def wrapper(*args, **kwargs):
            last_error = None
            for attempt in range(1, max_attempts + 1):
                try:
                    return func(*args, **kwargs)
                except exceptions as e:
                    last_error = e
                    print(f"  attempt {attempt}/{max_attempts} failed: {e}")
            raise last_error
        return wrapper
    return decorator

import random
random.seed(0)

@retry(max_attempts=4, exceptions=(ValueError,))
def flaky_operation(x):
    if random.random() < 0.5:
        raise ValueError(f"failed on {x}")
    return x * 2

print("\n=== Decorator factory ===")
try:
    result = flaky_operation(10)
    print(f"  result: {result}")
except ValueError:
    print("  all attempts failed")


# ─── 4. BUILDING GENERIC WRAPPERS ─────────────────────────────────────────────
#
# Generic wrappers work regardless of the wrapped function's signature.
# Patterns: timing, caching, rate limiting, circuit breaking, tracing.

def timed(func):
    """Measure execution time of any function."""
    @functools.wraps(func)
    def wrapper(*args, **kwargs):
        import time
        start = time.perf_counter()
        result = func(*args, **kwargs)
        elapsed = time.perf_counter() - start
        print(f"  {func.__name__} took {elapsed:.6f}s")
        return result
    return wrapper

def cached(func):
    """Simple cache for any function with hashable arguments."""
    cache = {}
    @functools.wraps(func)
    def wrapper(*args, **kwargs):
        # kwargs must be sorted for consistent cache key:
        key = (args, tuple(sorted(kwargs.items())))
        if key not in cache:
            cache[key] = func(*args, **kwargs)
        return cache[key]
    wrapper.cache = cache
    return wrapper

@timed
@cached
def slow_compute(n, *, multiplier=1):
    import time; time.sleep(0.01)
    return n * multiplier

print("\n=== Generic wrappers ===")
slow_compute(5, multiplier=2)   # cache miss → sleeps
slow_compute(5, multiplier=2)   # cache hit → instant
slow_compute(5, multiplier=3)   # cache miss → sleeps


# ─── 5. COLLECTING AND INSPECTING ARGUMENTS ───────────────────────────────────
#
# Advanced: use inspect.signature + bind to validate and introspect
# what arguments would be passed before actually calling the function.

def validate_args(func):
    """Validate types based on annotations before calling."""
    sig = inspect.signature(func)
    hints = func.__annotations__

    @functools.wraps(func)
    def wrapper(*args, **kwargs):
        bound = sig.bind(*args, **kwargs)
        bound.apply_defaults()
        for param_name, value in bound.arguments.items():
            if param_name in hints and param_name != "return":
                expected = hints[param_name]
                if not isinstance(value, expected):
                    raise TypeError(
                        f"{func.__name__}: {param_name} expected {expected.__name__}, "
                        f"got {type(value).__name__}"
                    )
        return func(*args, **kwargs)
    return wrapper

@validate_args
def calculate(a: int, b: int, *, op: str = "add") -> int:
    if op == "add": return a + b
    if op == "mul": return a * b
    raise ValueError(f"unknown op: {op}")

print("\n=== Type validation via signature inspection ===")
print(f"  calculate(3, 4):          {calculate(3, 4)}")
print(f"  calculate(3, 4, op='mul'): {calculate(3, 4, op='mul')}")
try:
    calculate(3, "four")
except TypeError as e:
    print(f"  TypeError: {e}")


# ─── 6. UNPACKING INTO CALLS — * AND ** AT CALL SITE ─────────────────────────
#
# At the CALL SITE:
#   *sequence → unpacks into positional arguments
#   **mapping → unpacks into keyword arguments
#
# This mirrors the definition syntax but works in reverse:
#   definition: *args   collects extra positionals INTO a tuple
#   call site:  *seq    SPREADS a sequence OUT into positional arguments

def add3(a, b, c):
    return a + b + c

args   = [1, 2, 3]
kwargs = {"a": 1, "b": 2, "c": 3}

print(f"\n=== Unpacking at call site ===")
print(f"  add3(*[1,2,3]):            {add3(*args)}")
print(f"  add3(**{{a:1,b:2,c:3}}):   {add3(**kwargs)}")
print(f"  add3(1, *[2, 3]):          {add3(1, *[2, 3])}")
print(f"  add3(*[1, 2], c=3):        {add3(*[1, 2], c=3)}")

# Multiple * unpacks in one call (Python 3.5+):
first  = [1, 2]
second = [3, 4, 5]
combined = [*first, *second]
print(f"\n  [*first, *second]: {combined}")

d1 = {"a": 1, "b": 2}
d2 = {"c": 3}
merged = {**d1, **d2}
print(f"  {{**d1, **d2}}: {merged}")


# ─── 7. ARGUMENT PASSING PATTERNS ─────────────────────────────────────────────

print("\n=== Practical patterns ===")

# Pattern 1: Function that behaves differently based on argument presence
def connect(host, port=None, *, timeout=30):
    if port is None:
        return f"connect({host}, timeout={timeout})"
    return f"connect({host}:{port}, timeout={timeout})"

print(f"  {connect('db.internal')}")
print(f"  {connect('db.internal', 5432)}")

# Pattern 2: Accept dict OR keyword args interchangeably
def create_record(data=None, **fields):
    """Accept {'key': val} dict OR key=val kwargs."""
    record = {}
    if data is not None:
        record.update(data)
    record.update(fields)
    return record

r1 = create_record({"name": "Alice", "age": 30})
r2 = create_record(name="Bob", age=25)
r3 = create_record({"name": "Charlie"}, age=40, active=True)   # mixed
print(f"  r1: {r1}")
print(f"  r2: {r2}")
print(f"  r3: {r3}")

# Pattern 3: Thread-safe argument forwarding
import threading

def run_in_thread(func, *args, **kwargs):
    """Run a function in a thread, forwarding all arguments."""
    t = threading.Thread(target=func, args=args, kwargs=kwargs)
    t.start()
    return t

results = []
def worker(n, multiplier=1):
    results.append(n * multiplier)

threads = [run_in_thread(worker, i, multiplier=2) for i in range(5)]
for t in threads: t.join()
print(f"  thread results: {sorted(results)}")


# ─── 8. EDGE CASES AND PITFALLS ───────────────────────────────────────────────

print("\n=== Edge cases ===")

# 1. *args is ALWAYS a tuple, even for one argument:
def check(*args):
    print(f"  args type: {type(args)}, value: {args}")

check(42)         # (42,) — a tuple, not just 42
check(1, 2, 3)

# 2. **kwargs preserves insertion order (Python 3.7+):
def show_order(**kwargs):
    print(f"  kwargs: {list(kwargs.keys())}")

show_order(z=3, a=1, m=2)   # z, a, m — insertion order preserved

# 3. You cannot pass positional arg after keyword arg:
try:
    eval("f(a=1, 2)")   # SyntaxError
except SyntaxError as e:
    print(f"  SyntaxError: {e}")

# 4. **kwargs dict is a NEW dict — modifying it doesn't affect caller:
def modify_kwargs(**kwargs):
    kwargs["new_key"] = "added"
    return kwargs

original = {"x": 1}
result = modify_kwargs(**original)
print(f"  original: {original}")   # unchanged
print(f"  result:   {result}")     # has new_key


# ─── EXERCISES ────────────────────────────────────────────────────────────────
#
# 1. Write a decorator @memoize that works correctly with both positional
#    and keyword arguments. The challenge: f(1, b=2) and f(a=1, b=2) and f(b=2, a=1)
#    should all map to the same cache key if the function signature allows it.
#    Hint: use inspect.signature + bind + apply_defaults to normalize args.
#
# 2. Build a middleware chain (like Flask/Django middleware):
#       def make_chain(*middlewares):
#           returns a function that applies each middleware in order
#    Each middleware is: def mw(next_handler): return lambda *a, **kw: ...
#    Test with logging_mw + auth_mw + the actual handler.
#
# 3. Write a function partial_right(func, *args, **kwargs) that pre-fills
#    the LAST positional arguments (unlike functools.partial which fills first).
#    Example: partial_right(print, end="\n---\n") wraps print with a custom end.
#
# 4. Use inspect.signature + bind_partial to write a function that checks
#    which arguments are still MISSING from a function call (not yet provided).
#    Example: missing_args(open, "file.txt") → ["mode"] still needed? etc.
#
# THOUGHT QUESTION:
#   Python's *args is always a tuple (immutable), but **kwargs is always a dict
#   (mutable). Since you receive a copy of the dict (not the caller's dict),
#   why is it still "safe" to mutate kwargs inside a function?
#   What would change if **kwargs gave you a view of the caller's data instead?
