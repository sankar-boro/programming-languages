# Appendix — Cheat Sheets and References

---

## Appendix A: Kotlin Syntax Cheat Sheet

### Variables and Types

```kotlin
val name: String = "Alice"         // immutable reference
var count: Int = 0                 // mutable reference
val inferred = "type inferred"     // type inferred as String
val multiLine = """                // raw string
    Line 1
    Line 2
""".trimIndent()

// Type conversions (all explicit)
val i = 42
val l = i.toLong()
val d = i.toDouble()
val s = i.toString()
val back = "42".toInt()
val safe = "abc".toIntOrNull()  // null if parse fails
```

### Nullable Types

```kotlin
val name: String? = null           // nullable String
val len = name?.length             // safe call → Int?
val len2 = name?.length ?: 0       // Elvis → Int
val len3 = name!!.length           // assertion → Int (throws NPE if null)
val cast = name as? String         // safe cast → String?

// Smart cast
if (name != null) {
    println(name.length)           // name is String here
}
name?.let { nonNull ->             // nonNull is String inside
    println(nonNull.length)
}
```

### Control Flow

```kotlin
// if expression
val max = if (a > b) a else b

// when expression
val desc = when (x) {
    1       -> "one"
    in 2..9 -> "small"
    is String -> "a string"
    else    -> "other"
}

// when without argument
when {
    x > 0 && y > 0 -> println("both positive")
    x < 0           -> println("x is negative")
    else            -> println("other")
}

// Ranges
1..10           // 1, 2, ..., 10 (inclusive)
1 until 10      // 1, 2, ..., 9 (exclusive upper)
10 downTo 1     // 10, 9, ..., 1
1..10 step 2    // 1, 3, 5, 7, 9

// Loops
for (i in 1..10) { }
for ((index, value) in list.withIndex()) { }
for ((key, value) in map) { }
while (condition) { }
do { } while (condition)
```

### Functions

```kotlin
// Basic function
fun add(a: Int, b: Int): Int = a + b

// Default parameters
fun greet(name: String = "World") = "Hello, $name!"

// Named arguments
greet(name = "Alice")

// Varargs
fun sum(vararg nums: Int) = nums.sum()
sum(1, 2, 3)

// Spread operator
val arr = intArrayOf(1, 2, 3)
sum(*arr)

// Extension function
fun String.isPalindrome() = this == this.reversed()

// Local function
fun outer() {
    fun inner() { /* sees outer's variables */ }
    inner()
}

// Tail recursion
tailrec fun fib(n: Int, a: Long = 0, b: Long = 1): Long =
    if (n == 0) a else fib(n - 1, b, a + b)

// Lambda
val double: (Int) -> Int = { x -> x * 2 }
val square = { x: Int -> x * x }
val negate: (Int) -> Int = { -it }  // implicit 'it'

// Higher-order function
fun apply(x: Int, f: (Int) -> Int) = f(x)
apply(5) { it * 3 }  // trailing lambda

// Function reference
::greet                  // top-level function
String::uppercase        // member function
```

### Classes

```kotlin
// Basic class
class Person(val name: String, var age: Int) {
    init {
        require(age >= 0) { "Age must be non-negative" }
    }
    
    fun greet() = "Hi, I'm $name"
}

// Data class — equals, hashCode, toString, copy, componentN generated
data class Point(val x: Double, val y: Double)

// Enum class
enum class Direction { NORTH, SOUTH, EAST, WEST }

// Sealed class
sealed class Result<out T> {
    data class Success<T>(val value: T) : Result<T>()
    data class Error(val msg: String) : Result<Nothing>()
}

// Object (singleton)
object Config { val timeout = 30 }

// Companion object
class Factory {
    companion object {
        fun create() = Factory()
        const val VERSION = "1.0"
    }
}

// Inheritance
open class Animal(val name: String) {
    open fun speak(): String = "..."
}

class Dog(name: String) : Animal(name) {
    override fun speak() = "Woof!"
}

// Abstract class
abstract class Shape {
    abstract fun area(): Double
}

// Interface
interface Printable {
    fun print()
    fun log() = println(this)  // default implementation
}

// Class delegation
class CachingList<T>(private val delegate: MutableList<T>) : MutableList<T> by delegate {
    // Only override methods you want to change
}

// Value class
@JvmInline
value class Email(val value: String)
```

### Collections

```kotlin
// Immutable
listOf(1, 2, 3)
setOf("a", "b")
mapOf("k" to "v")
emptyList<Int>()

// Mutable
mutableListOf(1, 2, 3)
mutableSetOf("a", "b")
mutableMapOf("k" to "v")

// Common operations
list.filter { it > 0 }
list.map { it * 2 }
list.flatMap { it.toString().toList() }
list.fold(0) { acc, n -> acc + n }
list.reduce { acc, n -> acc + n }
list.groupBy { it.category }
list.associateBy { it.id }
list.associateWith { it.length }
list.partition { it > 0 }     // returns Pair<List, List>
list.sortedBy { it.name }
list.sortedWith(compareBy({ it.age }, { it.name }))
list.take(5)
list.drop(5)
list.takeWhile { it < 10 }
list.dropWhile { it < 10 }
list.first()
list.firstOrNull()
list.find { it.active }
list.any { it > 0 }
list.all { it > 0 }
list.none { it < 0 }
list.count { it > 0 }
list.sumOf { it.amount }
list.maxByOrNull { it.score }
list.distinct()
list.zip(other)
list.chunked(3)
list.windowed(3)
list + other          // new combined list
list.toSet()
list.toMutableList()

// Sequences (lazy)
list.asSequence().filter { }.map { }.take(10).toList()
generateSequence(0) { it + 1 }  // infinite sequence
sequence { yield(1); yieldAll(listOf(2, 3)) }
```

### Scope Functions

```kotlin
// let — it, returns lambda result
val len = name?.let { it.length }

// run — this, returns lambda result
val result = "hello".run { uppercase() }

// with — this, returns lambda result  
val joined = with(sb) { append("a"); append("b"); toString() }

// apply — this, returns receiver
val obj = MyClass().apply { x = 1; y = 2 }

// also — it, returns receiver
list.also { println("Before: $it") }.add(4)
```

### Generics

```kotlin
// Basic generic function
fun <T> first(list: List<T>): T = list[0]

// Upper bound constraint
fun <T : Comparable<T>> max(a: T, b: T) = if (a > b) a else b

// Multiple constraints
fun <T> f(x: T) where T : Comparable<T>, T : Cloneable { }

// Out (covariant)
interface Producer<out T> { fun produce(): T }

// In (contravariant)
interface Consumer<in T> { fun consume(item: T) }

// Star projection
fun printList(list: List<*>) = list.forEach { println(it) }

// Reified (requires inline)
inline fun <reified T> isInstance(value: Any) = value is T
```

### Coroutines (Conceptual)

```kotlin
// Coroutine builders
runBlocking { }             // blocks current thread
launch { }                  // fire-and-forget, returns Job
async { }                   // returns Deferred<T>

// Suspend functions
suspend fun fetch(): String { delay(1000); return "data" }

// Awaiting results
val deferred = async { fetch() }
val result = deferred.await()

// Switching contexts
withContext(Dispatchers.IO) { /* I/O work */ }
withContext(Dispatchers.Default) { /* CPU work */ }

// Timeout
withTimeout(1000L) { fetch() }          // throws on timeout
withTimeoutOrNull(1000L) { fetch() }    // returns null on timeout

// Flow
flow { emit(1); emit(2) }
.filter { it > 0 }
.map { it * 2 }
.collect { println(it) }

(1..10).asFlow()
generateSequence(0) { it + 1 }.asFlow().take(100)
```

---

## Appendix B: Operator Reference

### Arithmetic Operators

| Operator | Function | Example |
|----------|----------|---------|
| `+` | `plus` | `a + b`, `a.plus(b)` |
| `-` | `minus` | `a - b` |
| `*` | `times` | `a * b` |
| `/` | `div` | `a / b` |
| `%` | `rem` | `a % b` |
| `-` (unary) | `unaryMinus` | `-a` |
| `+` (unary) | `unaryPlus` | `+a` |

### Augmented Assignment

| Operator | Equivalent |
|----------|-----------|
| `a += b` | `a = a + b` |
| `a -= b` | `a = a - b` |
| `a *= b` | `a = a * b` |
| `a /= b` | `a = a / b` |
| `a %= b` | `a = a % b` |

### Comparison Operators

| Operator | Function | Notes |
|----------|----------|-------|
| `==` | `equals` | Structural equality |
| `!=` | `!equals` | Structural inequality |
| `===` | N/A | Reference equality |
| `!==` | N/A | Reference inequality |
| `<` | `compareTo` | Returns Int; `< 0` means less |
| `>` | `compareTo` | |
| `<=` | `compareTo` | |
| `>=` | `compareTo` | |

### Other Operators

| Operator | Function / Meaning |
|----------|-------------------|
| `a[i]` | `get(i)` |
| `a[i] = v` | `set(i, v)` |
| `a in b` | `b.contains(a)` |
| `a !in b` | `!b.contains(a)` |
| `a..b` | `a.rangeTo(b)` |
| `a()` | `invoke()` |
| `++a`, `a++` | `inc()` |
| `--a`, `a--` | `dec()` |
| `!a` | `not()` |
| `a && b` | Short-circuit AND |
| `a \|\| b` | Short-circuit OR |

---

## Appendix C: Standard Library Quick Reference

### String Extensions

```kotlin
str.length
str.isEmpty() / isNotEmpty()
str.isBlank() / isNotBlank()
str.trim() / trimStart() / trimEnd()
str.uppercase() / lowercase()
str.capitalize() / decapitalize()
str.contains(substr)
str.startsWith(prefix) / endsWith(suffix)
str.indexOf(substr) / lastIndexOf(substr)
str.substring(start, end)
str.replace(old, new) / replaceFirst(old, new)
str.split(delimiter)
str.toInt() / toIntOrNull()
str.toDouble() / toDoubleOrNull()
str.toCharArray()
str.reversed()
str.repeat(n)
str.padStart(len, char) / padEnd(len, char)
str.format(vararg args)
str.matches(regex)
str.filter { predicate }
str.map { transform }
str[index]
str.first() / last()
```

### Number Extensions

```kotlin
n.absoluteValue  // (import kotlin.math.*)
n.sign
n.coerceIn(min, max)
n.coerceAtLeast(min)
n.coerceAtMost(max)
n.toInt() / toLong() / toDouble() / toFloat() / toByte() / toShort()
n.toString()
n.toString(radix)

// Math functions (kotlin.math)
abs(n) / sqrt(n) / pow(base, exp)
floor(n) / ceil(n) / round(n)
min(a, b) / max(a, b)
log(n) / log2(n) / log10(n)
sin(n) / cos(n) / tan(n)
PI / E
```

### Collection Extensions (selected)

```kotlin
// Creation
(1..10).toList()
List(n) { index -> value }
buildList { add(1); addAll(otherList) }
buildMap { put("k", "v") }
buildSet { add(1); add(2) }

// Transformation
map, flatMap, mapIndexed, mapNotNull, mapKeys, mapValues
filter, filterNot, filterIsInstance, filterNotNull, filterIndexed
associate, associateBy, associateWith
groupBy, groupingBy
sortedBy, sortedByDescending, sortedWith
distinctBy
zip, unzip, zipWithNext
chunked, windowed, flatten

// Aggregation
fold, reduce, foldRight, reduceRight
sum, sumOf
count, countBy
average
min, max, minOrNull, maxOrNull
minBy, maxBy, minByOrNull, maxByOrNull
minOf, maxOf

// Search
find, findLast
first, firstOrNull, last, lastOrNull
single, singleOrNull
any, all, none
contains, containsAll
indexOf, lastIndexOf

// Manipulation
plus, minus
take, drop, takeLast, dropLast
takeWhile, dropWhile
partition
reversed
shuffled
random
```

---

## Appendix D: Keywords Reference

| Keyword | Meaning |
|---------|---------|
| `val` | Read-only variable/property |
| `var` | Mutable variable/property |
| `fun` | Function declaration |
| `class` | Class declaration |
| `object` | Singleton object declaration |
| `interface` | Interface declaration |
| `enum` | Enumeration class |
| `sealed` | Sealed class/interface |
| `data` | Data class modifier |
| `abstract` | Abstract class/function |
| `open` | Allows inheritance/override |
| `final` | Prevents further override |
| `override` | Overrides a parent member |
| `inner` | Inner class (references outer) |
| `companion` | Companion object |
| `lateinit` | Deferred non-null initialization |
| `by` | Delegation |
| `in` | Contravariance / collection membership |
| `out` | Covariance |
| `reified` | Preserves type at runtime (in inline) |
| `inline` | Copies function body to call site |
| `noinline` | Excludes a lambda from inlining |
| `crossinline` | Prevents non-local return in lambda |
| `suspend` | Suspendable function marker |
| `tailrec` | Tail recursion optimization |
| `operator` | Operator overloading |
| `infix` | Infix function notation |
| `external` | Implemented externally (native/JS) |
| `expect/actual` | Multiplatform declarations |
| `vararg` | Variable number of arguments |
| `typealias` | Type alias declaration |
| `get/set` | Custom property accessor |
| `field` | Backing field in accessor |
| `it` | Implicit lambda parameter (one param) |
| `this` | Current receiver |
| `super` | Parent class/interface reference |
| `return` | Return from function |
| `throw` | Throw an exception |
| `try/catch/finally` | Exception handling |
| `if/else` | Conditional expression |
| `when` | Multi-branch expression |
| `for/while/do` | Loop constructs |
| `break/continue` | Loop control |
| `is/!is` | Type check |
| `as/as?` | Type cast (safe/unsafe) |
| `in/!in` | Range/collection membership |
| `..` | Range creation |
| `?:` | Elvis operator |
| `?.` | Safe call |
| `!!` | Not-null assertion |
| `*` | Spread operator |

---

## Appendix E: Quick Comparisons

### Kotlin vs Java Equivalents

| Java | Kotlin |
|------|--------|
| `String` | `String` |
| `int` | `Int` |
| `long` | `Long` |
| `boolean` | `Boolean` |
| `void` | `Unit` |
| `Object` | `Any` |
| `null` can be anywhere | `String?` explicitly |
| `instanceof` | `is` |
| `(String)obj` | `obj as String` |
| `obj == null ? default : obj` | `obj ?: default` |
| `if (obj != null) obj.x` | `obj?.x` |
| `static` | `companion object` |
| Anonymous class for SAM | Lambda directly |
| `new MyClass()` | `MyClass()` |
| `final` (local variable) | `val` |
| Getter/setter methods | Properties |
| `switch` (limited) | `when` (powerful) |
| Checked exceptions | No (use `@Throws` for interop) |
| No default params | Default parameters |
| No named args | Named arguments |

---

## Appendix F: Frequently Confused Concepts

### 1. val vs const val

```kotlin
val x = 42        // computed at runtime, can be any type
const val Y = 42  // computed at compile time, must be primitive or String
                  // Only allowed in companion object or top-level
```

### 2. object vs class

```kotlin
object Singleton { }  // one instance, created lazily on first access
class NotASingleton { }  // new instance each time with MyClass()
```

### 3. data class copy vs clone

```kotlin
data class Point(var x: Int, var y: Int)
val p1 = Point(1, 2)
val p2 = p1.copy()   // shallow copy — new object, same values
// For truly independent deep copies, implement your own
```

### 4. filter vs find

```kotlin
list.filter { it > 0 }    // returns ALL matching elements as a List
list.find { it > 0 }      // returns FIRST matching element or null
list.first { it > 0 }     // returns FIRST matching or throws NoSuchElementException
```

### 5. map vs forEach

```kotlin
list.map { transform(it) }      // returns a NEW list of transformed elements
list.forEach { process(it) }    // runs side effect, returns Unit
```

### 6. Sequence vs List

```kotlin
list.filter { }.map { }           // eager — creates intermediate List
list.asSequence().filter { }.map { }.toList()  // lazy — processes element by element
```

### 7. launch vs async

```kotlin
launch { /* fire and forget */ }         // returns Job (no result)
val result = async { compute() }.await() // returns Deferred<T> (has result)
```

---

## Appendix G: Useful Idioms at a Glance

```kotlin
// Swap two variables
a = b.also { b = a }

// Execute if not null
value?.let { doSomething(it) }

// Execute if null
value ?: doSomethingIfNull()

// Get non-null or throw with message
val v = value ?: error("Value must not be null")

// Get non-null or throw with context
val v = checkNotNull(value) { "Expected value to be set" }

// Validate argument
require(age >= 0) { "Age must be non-negative: $age" }

// Validate state
check(isConnected) { "Must be connected first" }

// Build a string
buildString {
    append("Hello")
    append(", ")
    append(name)
}

// Build a list
buildList {
    add(1); add(2)
    addAll(moreItems)
}

// Conditional add to list
val list = buildList<Int> {
    add(1)
    if (includeTwo) add(2)
    addAll(rest)
}

// Convert nullable to default
val display = name ?: "Anonymous"

// Repeat N times
repeat(5) { index -> println("Item $index") }

// Measure time
val elapsed = measureTimeMillis { doWork() }

// Apply multiple operations to one object
val configured = MyClass().apply {
    property1 = "value1"
    property2 = 42
    configure()
}

// Create a map from list
val map = list.associateBy { it.id }
val map2 = list.associateWith { it.computeValue() }

// Flatten nested lists
val flat = listOfLists.flatten()
val flatMapped = list.flatMap { it.subItems }

// Get or compute default in map
val value = cache.getOrPut(key) { computeExpensiveValue(key) }

// Safe index access
list.getOrNull(index) ?: default

// Find first match or default
list.firstOrNull { predicate } ?: default
```

---

*End of Appendix — Happy Kotlin coding!*

---

## Book Index — By Concept

| Concept | Chapter |
|---------|---------|
| val / var | Ch. 3 |
| Null safety (`?`, `?.`, `?:`, `!!`) | Ch. 7 |
| Smart casts | Ch. 7 |
| lateinit, lazy | Ch. 7 |
| Data classes | Ch. 6 |
| Sealed classes | Ch. 6 |
| Enum classes | Ch. 6 |
| Extension functions | Ch. 10 |
| Lambdas | Ch. 8 |
| Higher-order functions | Ch. 8 |
| inline, noinline, crossinline | Ch. 8 |
| Coroutines, suspend | Ch. 11 |
| Flow | Ch. 11 |
| Generics, variance (in/out) | Ch. 10 |
| reified types | Ch. 10 |
| Property delegation (by lazy, observable) | Ch. 10 |
| Class delegation (by) | Ch. 10 |
| Collections | Ch. 9 |
| Sequences | Ch. 9 |
| Java interop | Ch. 12 |
| Platform types | Ch. 12 |
| Kotlin internals / bytecode | Ch. 13 |
| tailrec | Ch. 5 |
| Scope functions | Ch. 8, 10 |
| when expressions | Ch. 4 |
| Ranges and progressions | Ch. 4 |
