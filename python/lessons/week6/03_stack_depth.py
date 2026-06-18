"""
WEEK 6 — DAY 3: Stack Depth, RecursionError, and Deep Algorithms
=================================================================
Topic: Python's recursion limit in practice, how to work with deeply
       recursive problems safely, and algorithms that require explicit
       stack management.

Key ideas:
  - RecursionError is a safety guard, not a feature limit to work around
  - Deep recursion in Python almost always signals an iterative solution exists
  - Backtracking and search algorithms need careful depth management
  - sys.setrecursionlimit() raises the limit but doesn't fix the root cause
  - Generator-based recursion can yield results without deep frame nesting
"""

import sys
import functools
from collections import deque


# ─── 1. THE RECURSION LIMIT IN DETAIL ────────────────────────────────────────
#
# sys.getrecursionlimit() returns the maximum call stack depth.
# Default: 1000. This is a SOFT limit enforced by CPython in ceval.c.
# Every Python call increments a depth counter. When it hits the limit,
# RecursionError is raised — even inside __repr__, __str__, etc.
#
# Why 1000? It's conservative to prevent the C stack from overflowing.
# The C stack is finite and operating-system dependent (~1MB–8MB typically).
# Each Python frame uses ~200–400 bytes of C stack space.

print(f"=== Recursion limit ===")
print(f"sys.getrecursionlimit(): {sys.getrecursionlimit()}")

# How deep can we actually go from here?
def measure_depth():
    def recurse(n):
        try:
            return recurse(n + 1)
        except RecursionError:
            return n
    return recurse(0)

available = measure_depth()
print(f"frames available from current depth: ~{available}")


# ─── 2. WHEN RECURSIONERROR IS A DESIGN SIGNAL ───────────────────────────────
#
# If you need more than ~900 levels of recursion, ask:
#   "Is the algorithm correct? Can it be iterative?"
#
# Cases where deep recursion legitimately happens:
#   - Very deep file system traversal (>1000 levels — extremely rare)
#   - Pathological graph structures passed as test cases
#   - Parsing deeply nested data (e.g., malicious JSON)
#
# In all these cases, an iterative solution with an explicit stack is safer.

def safe_depth_check(func, *args, depth_estimate=None):
    """Warn if we're dangerously close to the recursion limit."""
    current_depth = len(__import__("inspect").stack())
    remaining = sys.getrecursionlimit() - current_depth
    if depth_estimate and depth_estimate > remaining * 0.8:
        print(f"  WARNING: estimated depth {depth_estimate} may hit limit ({remaining} remaining)")
    return func(*args)


# ─── 3. BACKTRACKING — RECURSION'S KILLER APP ────────────────────────────────
#
# Backtracking = try a choice, recurse, if it fails → undo and try another.
# The call stack naturally stores "what choices have been made so far."
# This is recursion's strongest use case — the stack IS the search state.

def solve_n_queens(n):
    """
    Place n queens on an n×n chessboard so no two attack each other.
    Returns a list of all solutions (each solution is a list of column positions).
    """
    solutions = []

    def is_safe(queens, row, col):
        """Check if placing a queen at (row, col) is safe."""
        for r, c in enumerate(queens):
            if c == col:               return False   # same column
            if abs(r - row) == abs(c - col): return False   # diagonal
        return True

    def place(row, queens):
        if row == n:
            solutions.append(queens[:])   # found a complete solution
            return
        for col in range(n):
            if is_safe(queens, row, col):
                queens.append(col)        # make the choice
                place(row + 1, queens)    # recurse
                queens.pop()              # undo the choice (backtrack)

    place(0, [])
    return solutions

print("\n=== N-Queens backtracking ===")
for n in [4, 5, 6]:
    solutions = solve_n_queens(n)
    print(f"  {n}-queens: {len(solutions)} solutions")
    if n == 4:
        for sol in solutions:
            print(f"    {sol}")


# ─── 4. DEPTH-FIRST SEARCH — GRAPH TRAVERSAL ─────────────────────────────────
#
# DFS naturally uses a stack (the call stack for recursive DFS,
# or an explicit stack for iterative DFS).
# Key: mark visited nodes to avoid cycles in cyclic graphs.

def dfs_recursive(graph, start, visited=None):
    """Recursive DFS — uses call stack as the traversal stack."""
    if visited is None:
        visited = set()
    visited.add(start)
    result = [start]
    for neighbor in sorted(graph.get(start, [])):
        if neighbor not in visited:
            result.extend(dfs_recursive(graph, neighbor, visited))
    return result

def dfs_iterative(graph, start):
    """Iterative DFS — uses explicit stack."""
    visited = set()
    stack   = [start]
    result  = []

    while stack:
        node = stack.pop()
        if node in visited:
            continue
        visited.add(node)
        result.append(node)
        for neighbor in sorted(graph.get(node, []), reverse=True):
            if neighbor not in visited:
                stack.append(neighbor)

    return result

graph = {
    "A": ["B", "C"],
    "B": ["D", "E"],
    "C": ["F"],
    "D": [],
    "E": ["F"],
    "F": []
}

print("\n=== DFS: recursive vs iterative ===")
print(f"  recursive: {dfs_recursive(graph, 'A')}")
print(f"  iterative: {dfs_iterative(graph, 'A')}")


# ─── 5. BFS — WHEN ITERATION IS NATURAL ──────────────────────────────────────
#
# Breadth-first search visits nodes level by level.
# BFS is naturally iterative (uses a queue, not a stack).
# A recursive BFS is awkward — it's a sign to use iteration.

def bfs(graph, start):
    """BFS using a deque (efficient queue)."""
    visited = set([start])
    queue   = deque([start])
    result  = []

    while queue:
        node = queue.popleft()
        result.append(node)
        for neighbor in sorted(graph.get(node, [])):
            if neighbor not in visited:
                visited.add(neighbor)
                queue.append(neighbor)

    return result

print(f"\n=== BFS (natural iteration) ===")
print(f"  bfs from A: {bfs(graph, 'A')}")


# ─── 6. GENERATOR-BASED DEEP TRAVERSAL ───────────────────────────────────────
#
# Generators can yield from recursive structures without deep frame nesting.
# `yield from` delegates to a sub-generator — still uses frames but can be
# combined with itertools.islice to limit results without full traversal.

def walk_nested(obj):
    """Generator: yield all non-dict leaf values in a nested dict."""
    if isinstance(obj, dict):
        for value in obj.values():
            yield from walk_nested(value)   # delegate to sub-generator
    else:
        yield obj

data = {
    "a": 1,
    "b": {"c": 2, "d": {"e": 3, "f": {"g": 4}}},
    "h": 5
}

print(f"\n=== Generator-based traversal ===")
print(f"  leaves: {list(walk_nested(data))}")

# Can consume lazily — no full tree traversal until needed:
import itertools
first_two = list(itertools.islice(walk_nested(data), 2))
print(f"  first two leaves: {first_two}")


# ─── 7. PATHFINDING — COMBINING DEPTH AND BACKTRACKING ───────────────────────
#
# Find all paths from source to destination in a directed graph.
# Classic backtracking: explore a path, if it reaches the goal record it,
# backtrack and try other options.

def find_all_paths(graph, start, end, path=None):
    """Find all simple paths (no repeated nodes) from start to end."""
    if path is None:
        path = []
    path = path + [start]   # + creates a new list (safe for backtracking)

    if start == end:
        return [path]

    paths = []
    for neighbor in graph.get(start, []):
        if neighbor not in path:   # avoid cycles
            new_paths = find_all_paths(graph, neighbor, end, path)
            paths.extend(new_paths)
    return paths

dag = {
    "A": ["B", "C"],
    "B": ["D", "E"],
    "C": ["E"],
    "D": ["F"],
    "E": ["F"],
    "F": []
}

print(f"\n=== All paths A → F ===")
paths = find_all_paths(dag, "A", "F")
for p in paths:
    print(f"  {' → '.join(p)}")


# ─── 8. STACK DEPTH PROFILING ─────────────────────────────────────────────────
#
# A decorator to track maximum recursion depth reached by a function.
# Useful for understanding real-world depth requirements.

def track_depth(func):
    """Decorator: track and report maximum recursion depth."""
    max_depth = [0]
    current_depth = [0]

    @functools.wraps(func)
    def wrapper(*args, **kwargs):
        current_depth[0] += 1
        max_depth[0] = max(max_depth[0], current_depth[0])
        try:
            return func(*args, **kwargs)
        finally:
            current_depth[0] -= 1

    wrapper.max_depth  = max_depth
    wrapper.reset      = lambda: max_depth.__setitem__(0, 0)
    return wrapper

@track_depth
def fib(n):
    if n <= 1: return n
    return fib(n - 1) + fib(n - 2)

print(f"\n=== Depth tracking ===")
fib(10)
print(f"  fib(10) max depth: {fib.max_depth[0]}")
fib.reset()
fib(20)
print(f"  fib(20) max depth: {fib.max_depth[0]}")


# ─── EXERCISES ────────────────────────────────────────────────────────────────
#
# 1. Implement solve_maze(maze, start, end) using backtracking.
#    maze is a 2D list: 0 = open, 1 = wall.
#    Return the path as a list of (row, col) tuples, or None if no path exists.
#    Convert it to an iterative DFS after.
#
# 2. Write generate_subsets(lst) that returns all 2^n subsets of a list.
#    Use recursion. Trace the call tree for generate_subsets([1,2,3]).
#    Then write an iterative version.
#
# 3. Profile the N-queens solver:
#    Use the track_depth decorator to find the max recursion depth for n=8.
#    Is it safe (well within the 1000 limit)?
#    What is the maximum n where you'd worry about the depth limit?
#
# 4. Implement an iterative version of find_all_paths using an explicit stack.
#    Each stack entry should be (current_node, current_path).
#    Verify it finds the same paths as the recursive version.
#
# THOUGHT QUESTION:
#   The N-queens backtracking algorithm has time complexity O(n!) in the
#   worst case. But pruning (is_safe check) means most branches are cut early.
#   What determines how much pruning happens?
#   Why is backtracking often more practical than its worst-case complexity suggests?
