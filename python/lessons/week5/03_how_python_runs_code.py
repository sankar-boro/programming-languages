"""
WEEK 5 — DAY 3: How Python Runs Your Code — End-to-End
========================================================
Topic: The full pipeline from source text to execution — tokenization,
       parsing, AST, compilation, bytecode, and the PVM. Each stage
       examined with live Python tools.

Key ideas:
  - Python transforms source through 5 stages before executing a single line
  - The AST is a tree of nodes representing your program's structure
  - You can inspect and even modify the AST before compilation
  - .pyc files cache the bytecode to skip re-compilation on unchanged files
  - Understanding this pipeline explains import behavior, eval, exec, and macros
"""

import ast
import dis
import sys
import tokenize
import io
import marshal
import importlib.util


# ─── 1. THE 5-STAGE PIPELINE ─────────────────────────────────────────────────
#
# Stage 1: Source text (your .py file)
# Stage 2: Tokens (keywords, names, numbers, operators, whitespace markers)
# Stage 3: AST (Abstract Syntax Tree — grammar structure)
# Stage 4: Code object (compiled bytecode + metadata)
# Stage 5: PVM execution (the interpreter loop)
#
# You can intercept at any stage with Python's standard library tools.

SOURCE = """
x = 10
y = x * 2 + 1
if y > 15:
    print(y)
"""


# ─── 2. STAGE 1 → 2: TOKENIZATION ───────────────────────────────────────────
#
# The tokenizer breaks source text into a stream of tokens.
# Tokens are the atomic units: NAME, NUMBER, OP, STRING, NEWLINE, INDENT, etc.
# Indentation is converted to INDENT/DEDENT tokens — this is how Python
# enforces structure without braces.

print("=== Stage 2: Tokens ===")
tokens = list(tokenize.generate_tokens(io.StringIO(SOURCE).readline))
for tok in tokens[:20]:
    if tok.type not in (tokenize.NEWLINE, tokenize.NL, tokenize.ENCODING):
        print(f"  {tokenize.tok_name[tok.type]:<12} {tok.string!r}")


# ─── 3. STAGE 2 → 3: ABSTRACT SYNTAX TREE (AST) ─────────────────────────────
#
# The parser takes the token stream and builds an AST.
# Each node is an instance of an ast.* class.
# The AST captures STRUCTURE — operator precedence, block nesting, etc.
# The AST is still language-level (variable names, not memory addresses).

print("\n=== Stage 3: AST ===")
tree = ast.parse(SOURCE)

print(ast.dump(tree, indent=2))

# Walk all nodes:
print("\n--- AST node types ---")
for node in ast.walk(tree):
    print(f"  {type(node).__name__}", end=" ")
print()


# ─── 4. AST MANIPULATION ─────────────────────────────────────────────────────
#
# You can modify the AST before compilation.
# This is how libraries like pytest rewrite assert statements to give
# detailed failure messages.

class DoubleNumbers(ast.NodeTransformer):
    """AST transform: replace all integer literals with double their value."""
    def visit_Constant(self, node):
        if isinstance(node.value, int):
            return ast.Constant(value=node.value * 2)
        return node

print("\n=== AST manipulation (double all int literals) ===")
original  = ast.parse("result = 3 + 7")
modified  = DoubleNumbers().visit(original)
ast.fix_missing_locations(modified)

code_obj  = compile(modified, "<ast>", "exec")
namespace = {}
exec(code_obj, namespace)
print(f"  3 + 7 with doubled literals = {namespace['result']}")   # (3*2) + (7*2) = 20


# ─── 5. STAGE 3 → 4: COMPILATION TO BYTECODE ─────────────────────────────────
#
# The compiler walks the AST and emits bytecode instructions into a code object.
# This is where:
#   - Variable classifications happen (local/global/free)
#   - Constant folding occurs (2 * 3 → 6 at compile time)
#   - Jump targets are resolved
#
# Constant folding example:

print("\n=== Constant folding ===")

def uses_constant():
    return 2 * 3 * 7    # computed at compile time, not runtime

for instr in dis.get_instructions(uses_constant):
    print(f"  {instr.opname:<25} {instr.argval!r}")
# You'll see LOAD_CONST 42 — not three separate multiply operations!


# ─── 6. STAGE 4 → 5: THE PVM AND .pyc FILES ──────────────────────────────────
#
# Before the PVM runs the code, CPython checks for a cached .pyc file.
# .pyc lives in __pycache__/ and contains:
#   - Magic number (Python version identifier)
#   - Source file modification timestamp + size
#   - Marshalled (serialized) code object
#
# If source hasn't changed, CPython loads the cached bytecode → faster startup.
# marshal is the serialization format for code objects.

print("\n=== .pyc structure ===")

# Compile and marshal a simple code object to see what .pyc contains:
source = "x = 1 + 2"
code_obj = compile(source, "<demo>", "exec")

# marshal.dumps serializes the code object to bytes:
marshalled = marshal.dumps(code_obj)
print(f"  code object marshalled size: {len(marshalled)} bytes")
print(f"  first 16 bytes (hex):        {marshalled[:16].hex()}")

# Recover it:
recovered = marshal.loads(marshalled)
print(f"  recovered type: {type(recovered)}")
print(f"  same bytecode: {recovered.co_code == code_obj.co_code}")


# ─── 7. IMPORT MECHANICS ─────────────────────────────────────────────────────
#
# When you write `import math`, Python:
#   1. Checks sys.modules (cache) — if already imported, return cached module
#   2. Finds the module file (searches sys.path)
#   3. Compiles it (or loads .pyc cache)
#   4. Executes the module-level code in a fresh namespace
#   5. Stores the module object in sys.modules
#   6. Binds the name "math" in the importing namespace
#
# ALL of a module's top-level code runs on FIRST import.
# Subsequent imports return the cached module — code does NOT re-run.

print("\n=== Import mechanics ===")
print(f"  'sys' already in sys.modules: {'sys' in sys.modules}")
print(f"  'math' in sys.modules before import: {'math' in sys.modules}")

import math
print(f"  'math' in sys.modules after import:  {'math' in sys.modules}")
print(f"  id(math): {id(sys.modules['math'])}")

import math as math2
print(f"  same object? {math is math2}")   # True — cached, not re-executed


# ─── 8. EVAL, EXEC, COMPILE — THE PIPELINE ON DEMAND ─────────────────────────
#
# eval(expr):     evaluates a single expression, returns its value
# exec(code):     executes statements, returns None
# compile(src, ...): turns source/AST into a code object manually
#
# These expose the pipeline at the Python level — useful for:
#   - Dynamic code generation
#   - Sandboxed evaluation (with restricted globals)
#   - Building DSLs and template engines

print("\n=== eval / exec / compile ===")

# eval — expression mode:
result = eval("2 ** 8 + len('hello')")
print(f"  eval('2**8 + len(hello)'): {result}")

# exec — statement mode with custom namespace:
ns = {"__builtins__": {}}   # restricted: no built-ins
ns["x"] = 5
exec("y = x * x", ns)
print(f"  exec in restricted ns: y = {ns['y']}")

# Stages exposed manually:
src  = "a + b"
tree = ast.parse(src, mode="eval")
code = compile(tree, "<manual>", "eval")
val  = eval(code, {"a": 3, "b": 4})
print(f"  manual pipeline: {src} = {val}")


# ─── 9. sys.path AND MODULE SEARCH ORDER ─────────────────────────────────────
#
# When importing, Python searches sys.path in order:
#   1. '' (current directory)
#   2. PYTHONPATH environment variable entries
#   3. Standard library directories
#   4. Site-packages (third-party packages)
#
# You can add to sys.path at runtime to import from non-standard locations.

print("\n=== sys.path (module search) ===")
for i, p in enumerate(sys.path[:5]):
    print(f"  [{i}] {p!r}")
print(f"  ... ({len(sys.path)} total entries)")


# ─── EXERCISES ────────────────────────────────────────────────────────────────
#
# 1. Write a function tokenize_source(code_str) that returns a list of
#    (token_type_name, token_string) for non-whitespace tokens.
#    Test it on a few lines of Python code.
#
# 2. Use ast.parse() and ast.walk() to count how many function definitions
#    (ast.FunctionDef nodes) are in a given source string.
#    Write: count_functions(source: str) -> int
#
# 3. Write an AST transformer that renames all occurrences of variable "x"
#    to "renamed_x". Test by compiling and executing the transformed AST.
#    Hint: visit ast.Name nodes and check id field.
#
# 4. Demonstrate the import cache:
#    - Import a module
#    - Manually delete it from sys.modules
#    - Import it again
#    - Show that the module-level code runs again (use a print statement
#      in a small custom module, or use importlib.reload)
#
# THOUGHT QUESTION:
#   Python's `import` executes the entire module body on first import.
#   This means a slow or buggy import can block your entire program startup.
#   How could you defer or lazily execute module-level code?
#   What mechanism does Python provide for lazy imports in 3.12+?
#   (Hint: look up importlib.util.LazyLoader)
