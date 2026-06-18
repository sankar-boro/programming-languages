# Interview Preparation — Kotlin Language Questions

> *"An interview is a conversation about what you know, how you think, and what you can build."*

---

## Overview

This chapter covers common Kotlin interview questions at three levels: **conceptual understanding**, **code analysis**, and **coding problems**. These focus exclusively on the Kotlin language itself, not on Android or backend frameworks.

---

## Section 1: Conceptual Questions and Answers

### Basics

**Q1: What is the difference between `val` and `var`?**

`val` declares a read-only reference — once assigned, it cannot be reassigned. `var` declares a mutable reference that can be reassigned at any time.

Important: `val` means the reference is immutable, not the object it points to. A `val` can refer to a `MutableList`, and you can still add items to that list.

```kotlin
val list = mutableListOf(1, 2, 3)
list.add(4)      // OK — the list is mutable
// list = mutableListOf()  // ERROR — the reference is immutable
```

---

**Q2: What is the difference between `==` and `===` in Kotlin?**

- `==` calls `equals()` — structural equality (content comparison)
- `===` checks reference identity — whether two references point to the same object

```kotlin
val s1 = "Hello"
val s2 = "Hello"
val s3 = s1

println(s1 == s2)   // true — same content
println(s1 === s2)  // true — JVM string interning (implementation detail)
println(s1 === s3)  // true — same reference

data class Point(val x: Int, val y: Int)
val p1 = Point(1, 2)
val p2 = Point(1, 2)

println(p1 == p2)   // true — equals() by content
println(p1 === p2)  // false — different objects
```

---

**Q3: What is a data class and what does it automatically generate?**

A `data class` is designed to hold data. The compiler automatically generates:
- `equals()` — based on all primary constructor properties
- `hashCode()` — consistent with equals()
- `toString()` — readable representation
- `copy()` — creates a copy with some properties changed
- `componentN()` functions — for destructuring

```kotlin
data class Person(val name: String, val age: Int)
// One line replaces 50+ lines of Java boilerplate
```

---

**Q4: What is a sealed class?**

A sealed class restricts its hierarchy — all subclasses must be in the same package (Kotlin 1.5+). This gives the compiler complete knowledge of all subtypes, enabling exhaustive `when` matching without an `else` branch.

```kotlin
sealed class Result<out T> {
    data class Success<T>(val value: T) : Result<T>()
    data class Error(val message: String) : Result<Nothing>()
}

fun handle(result: Result<String>) = when (result) {
    is Result.Success -> println(result.value)
    is Result.Error   -> println("Error: ${result.message}")
    // No else needed — compiler knows all subclasses
}
```

---

**Q5: Explain Kotlin's null safety system.**

Kotlin separates types into nullable (`String?`) and non-null (`String`). The compiler prevents operations on nullable types that could cause NPE:

- `?.` (safe call): returns null if receiver is null
- `?:` (Elvis): provides a default if expression is null
- `!!` (not-null assertion): throws NPE if null
- Smart casts: compiler tracks null checks and automatically casts

```kotlin
val name: String? = getName()
println(name?.length ?: 0)  // safe — no NPE possible
```

---

**Q6: What is the difference between `object`, `companion object`, and `data object`?**

- `object` declaration: creates a singleton class
- `companion object`: object inside a class, replaces Java's `static` members
- `data object` (Kotlin 1.9+): singleton with auto-generated `toString()`

```kotlin
object Database { /* singleton */ }

class MyClass {
    companion object {
        fun create() = MyClass()  // called as MyClass.create()
    }
}

data object Disconnected  // toString() returns "Disconnected"
```

---

**Q7: What are extension functions and how do they work?**

Extension functions allow adding methods to existing classes without modifying them or using inheritance. They're syntactic sugar — at the JVM level, they're static methods that take the receiver as the first argument.

```kotlin
fun String.isPalindrome() = this == this.reversed()

"racecar".isPalindrome()  // true — looks like a String method
// Compiles to: isPalindrome("racecar")  — static call
```

Key limitation: extension functions cannot override member functions. If a member function and extension function have the same signature, the member always wins.

---

**Q8: What is the difference between `List` and `MutableList`?**

`List` is a read-only interface — it has no `add()`, `remove()`, or other mutation methods. `MutableList` extends `List` and adds mutation operations.

Important: read-only is not the same as immutable. A `List` reference to a `MutableList` object is read-only through that reference, but the underlying object can still be mutated through another reference.

```kotlin
val mutable = mutableListOf(1, 2, 3)
val readOnly: List<Int> = mutable  // read-only VIEW of the same list
mutable.add(4)
println(readOnly)  // [1, 2, 3, 4] — changed!
```

---

### Intermediate

**Q9: Explain `inline` functions. Why and when should you use them?**

`inline` functions copy their body (and passed lambdas) to every call site. Benefits:

1. **No lambda object allocation**: removes overhead of creating function objects
2. **Non-local returns**: lambdas inside inline functions can return from the enclosing function
3. **Reified type parameters**: type information is preserved at call sites

Use inline when:
- The function is a small wrapper that primarily passes lambdas
- Performance matters in hot code paths
- You need non-local returns or reified types

Don't use inline for large function bodies (causes code bloat).

```kotlin
inline fun <reified T> isA(value: Any) = value is T

println(isA<String>("hello"))  // true — only works with reified
```

---

**Q10: What is the difference between `lazy` and `lateinit`?**

| Aspect | `lazy` | `lateinit` |
|--------|--------|------------|
| Modifier type | `val` | `var` |
| Type | Any type | Non-primitive, non-null only |
| Initialized when | First access | Manually before use |
| Thread safety | Synchronized by default | No |
| Access before init | Returns initialized value | `UninitializedPropertyAccessException` |

```kotlin
val heavy: DatabaseConnection by lazy {
    createExpensiveConnection()  // called once, on first access
}

lateinit var dependency: Service  // must be set before use
```

---

**Q11: Explain co-variance (`out`) and contra-variance (`in`) in Kotlin generics.**

**Covariance (`out`)**: if `Dog` is a subtype of `Animal`, `Producer<Dog>` is a subtype of `Producer<Animal>`. Type can only appear in "out" (return) positions.

**Contravariance (`in`)**: if `Animal` is a supertype of `Dog`, `Consumer<Animal>` is a subtype of `Consumer<Dog>`. Type can only appear in "in" (parameter) positions.

```kotlin
interface Producer<out T> { fun produce(): T }     // covariant
interface Consumer<in T> { fun consume(item: T) }  // contravariant

val dogProducer: Producer<Dog> = DogFactory()
val animalProducer: Producer<Animal> = dogProducer  // OK — covariance

val animalConsumer: Consumer<Animal> = AnimalShelter()
val dogConsumer: Consumer<Dog> = animalConsumer  // OK — contravariance
```

---

**Q12: What is structured concurrency in Kotlin coroutines?**

Structured concurrency means every coroutine belongs to a scope. When the scope ends:
- All child coroutines are cancelled (if still running)
- All child coroutines must complete before the scope ends normally
- If a child fails, the failure propagates to the parent

This prevents coroutine leaks and ensures predictable lifecycle management.

```kotlin
suspend fun doWork() = coroutineScope {
    val job1 = launch { /* work 1 */ }
    val job2 = launch { /* work 2 */ }
    // Both must complete before doWork() returns
    // If job1 fails, job2 is cancelled and the exception propagates
}
```

---

**Q13: What is a `Flow` in Kotlin?**

`Flow<T>` is a cold, asynchronous stream of values. Unlike a `List`, values are computed and emitted over time. Unlike a `Sequence`, operations can suspend without blocking a thread.

Cold means: the producer code doesn't run until a consumer calls `collect()`.

```kotlin
fun evenNumbers(): Flow<Int> = flow {
    for (i in 1..100) {
        if (i % 2 == 0) emit(i)  // emit only even numbers
        delay(10)                  // non-blocking pause
    }
}

evenNumbers().collect { value -> println(value) }
```

---

**Q14: What are coroutine dispatchers and when would you use each?**

Dispatchers determine which thread(s) a coroutine runs on:

- `Dispatchers.Main`: UI thread (Android/Swing) — for UI updates
- `Dispatchers.Default`: CPU-intensive work — shared pool sized to CPU cores
- `Dispatchers.IO`: I/O operations — larger pool, designed for blocking operations
- `Dispatchers.Unconfined`: runs on the current thread until first suspension

```kotlin
suspend fun fetchAndProcess() {
    val data = withContext(Dispatchers.IO) { readFromDatabase() }
    val processed = withContext(Dispatchers.Default) { processData(data) }
}
```

---

## Section 2: Code Analysis Questions

**Q15: What is the output of this code?**

```kotlin
var x = 5
val result = when {
    x > 3  -> "greater"
    x == 5 -> "equal"
    else   -> "other"
}
println(result)
```

**Answer:** `"greater"` — `when` evaluates branches in order and takes the first match. `x > 3` is true (5 > 3), so it returns "greater" even though `x == 5` is also true.

---

**Q16: What is the output?**

```kotlin
fun main() {
    val list = listOf(1, 2, 3, null, 5)
    val result = list.filterNotNull().sum()
    println(result)
}
```

**Answer:** `11` — `filterNotNull()` removes the null element, leaving `[1, 2, 3, 5]`, which sums to 11.

---

**Q17: Will this code compile? If not, why?**

```kotlin
class MyClass {
    var name: String? = null
    
    fun printName() {
        if (name != null) {
            println(name.length)
        }
    }
}
```

**Answer:** No, it won't compile. `name` is a `var` — the compiler can't guarantee it won't change between the null check and the `name.length` call (e.g., another thread could set it to null). The compiler requires a smart cast condition, which fails for mutable properties. Fix: `val localName = name; if (localName != null) println(localName.length)` or `name?.let { println(it.length) }`.

---

**Q18: What is printed?**

```kotlin
val numbers = listOf(1, 2, 3, 4, 5)
val result = numbers
    .filter { it > 2 }
    .map { it * it }
    .fold(0) { acc, n -> acc + n }
println(result)
```

**Answer:** `50` — filter gives [3, 4, 5], map gives [9, 16, 25], fold sums them: 9 + 16 + 25 = 50.

---

**Q19: Is this a tail-recursive function? Can `tailrec` be applied?**

```kotlin
fun sum(n: Int): Int {
    if (n <= 0) return 0
    return n + sum(n - 1)
}
```

**Answer:** No, this is NOT tail-recursive. The recursive call `sum(n - 1)` is followed by an addition (`n + ...`). In tail recursion, the recursive call must be the **last** operation. To make it tail-recursive, use an accumulator:

```kotlin
tailrec fun sum(n: Int, acc: Int = 0): Int {
    if (n <= 0) return acc
    return sum(n - 1, acc + n)  // recursive call IS the last operation
}
```

---

**Q20: What is the difference between these two?**

```kotlin
// Version A
val a = async { longOperation() }
val b = async { anotherOperation() }
val result = a.await() + b.await()

// Version B
val a = longOperation()
val b = anotherOperation()
val result = a + b
```

**Answer:** Version A runs both operations **concurrently** — `anotherOperation()` starts before `longOperation()` finishes. Total time ≈ max(timeA, timeB).

Version B is **sequential** — `anotherOperation()` doesn't start until `longOperation()` returns. Total time ≈ timeA + timeB.

---

## Section 3: Coding Problems

### Problem 1: Null Safety

Write a function `safeGet(map: Map<String, Any?>, key: String): String` that:
- Returns the value as a String if it exists and is a String
- Returns "missing" if the key doesn't exist
- Returns "null" if the value is null
- Returns "wrong type" if the value is not a String

**Solution:**
```kotlin
fun safeGet(map: Map<String, Any?>, key: String): String {
    return when (val value = map[key]) {
        null -> if (key !in map) "missing" else "null"
        is String -> value
        else -> "wrong type"
    }
}
```

---

### Problem 2: Extension Functions

Write extension functions on `List<Int>`:
- `median(): Double` — returns the median value
- `mode(): List<Int>` — returns the most frequent value(s)

**Solution:**
```kotlin
fun List<Int>.median(): Double {
    if (isEmpty()) return 0.0
    val sorted = sorted()
    return if (size % 2 == 0) {
        (sorted[size / 2 - 1] + sorted[size / 2]) / 2.0
    } else {
        sorted[size / 2].toDouble()
    }
}

fun List<Int>.mode(): List<Int> {
    if (isEmpty()) return emptyList()
    val frequency = groupBy { it }.mapValues { it.value.size }
    val maxFreq = frequency.values.max()
    return frequency.filter { it.value == maxFreq }.keys.sorted()
}

// Test
val numbers = listOf(1, 2, 2, 3, 3, 4, 5)
println(numbers.median())  // 3.0
println(numbers.mode())    // [2, 3]
```

---

### Problem 3: Sealed Classes

Model a calculator that handles:
- `Add(a: Double, b: Double)`
- `Subtract(a: Double, b: Double)`
- `Multiply(a: Double, b: Double)`
- `Divide(a: Double, b: Double)` — can fail

**Solution:**
```kotlin
sealed class Operation {
    data class Add(val a: Double, val b: Double) : Operation()
    data class Subtract(val a: Double, val b: Double) : Operation()
    data class Multiply(val a: Double, val b: Double) : Operation()
    data class Divide(val a: Double, val b: Double) : Operation()
}

sealed class CalculatorResult {
    data class Value(val result: Double) : CalculatorResult()
    data class Error(val message: String) : CalculatorResult()
}

fun calculate(op: Operation): CalculatorResult = when (op) {
    is Operation.Add      -> CalculatorResult.Value(op.a + op.b)
    is Operation.Subtract -> CalculatorResult.Value(op.a - op.b)
    is Operation.Multiply -> CalculatorResult.Value(op.a * op.b)
    is Operation.Divide   ->
        if (op.b == 0.0) CalculatorResult.Error("Division by zero")
        else CalculatorResult.Value(op.a / op.b)
}

fun main() {
    listOf(
        Operation.Add(5.0, 3.0),
        Operation.Divide(10.0, 0.0),
        Operation.Multiply(2.5, 4.0)
    ).forEach { op ->
        val result = calculate(op)
        when (result) {
            is CalculatorResult.Value -> println("Result: ${result.result}")
            is CalculatorResult.Error -> println("Error: ${result.message}")
        }
    }
}
```

---

### Problem 4: Higher-Order Functions

Implement a `pipeline` function that applies a list of transformations to a value:

```kotlin
fun <T> pipeline(initial: T, vararg transforms: (T) -> T): T
```

**Solution:**
```kotlin
fun <T> pipeline(initial: T, vararg transforms: (T) -> T): T =
    transforms.fold(initial) { acc, transform -> transform(acc) }

fun main() {
    val result = pipeline(
        "  hello, world!  ",
        String::trim,
        { it.split(", ").joinToString(" and ") },
        String::uppercase,
        { "[$it]" }
    )
    println(result)  // [HELLO AND WORLD!]
}
```

---

### Problem 5: Coroutines

Write a `parallelMap` extension function on `List<T>` that maps each element using a suspend function, running all mappings concurrently:

```kotlin
suspend fun <T, R> List<T>.parallelMap(transform: suspend (T) -> R): List<R>
```

**Solution:**
```kotlin
import kotlinx.coroutines.*

suspend fun <T, R> List<T>.parallelMap(transform: suspend (T) -> R): List<R> =
    coroutineScope {
        map { async { transform(it) } }.awaitAll()
    }

// Usage:
fun main() = runBlocking {
    val ids = listOf("user1", "user2", "user3")
    
    val users = ids.parallelMap { id ->
        delay(100)  // simulate network fetch
        "User($id)"
    }
    
    println(users)  // [User(user1), User(user2), User(user3)]
    // Takes ~100ms, not 300ms
}
```

---

### Problem 6: Generics with Reified

Write a function `filterAndCast<T>(list: List<Any>): List<T>` that filters a mixed list to only items of type T:

**Solution:**
```kotlin
inline fun <reified T> filterAndCast(list: List<Any>): List<T> =
    list.filterIsInstance<T>()

// Or manually:
inline fun <reified T> List<Any>.ofType(): List<T> = filterIsInstance<T>()

fun main() {
    val mixed = listOf(1, "hello", 2.0, "world", 3, true)
    
    println(filterAndCast<String>(mixed))  // [hello, world]
    println(filterAndCast<Int>(mixed))     // [1, 3]
    println(mixed.ofType<String>())        // [hello, world]
}
```

---

## Section 4: Quick-Fire Q&A

| Question | Answer |
|----------|--------|
| What does `it` refer to in a lambda? | Implicit parameter name when the lambda has exactly one parameter |
| What is the difference between `run` and `let`? | `run` uses `this`, `let` uses `it`; both return the lambda result |
| What is `Nothing` type? | A type with no values — represents a function that never returns normally |
| Can you inherit from a data class? | Yes, but a data class cannot be inherited FROM (it's final by default) |
| What is the spread operator (`*`)? | Spreads an array into vararg position |
| What is `by` keyword used for? | Class delegation and property delegation |
| What does `crossinline` do? | Prevents non-local returns in lambdas passed to inline functions |
| What is a platform type? | A type from Java with unknown nullability (shown as `T!`) |
| Can sealed class have non-sealed subclasses? | Yes, in the same package; they can be open or abstract |
| What is `@JvmInline value class`? | Zero-cost wrapper — the wrapper is eliminated in bytecode |

---

## Key Interview Tips

1. **Know the why**: Don't just say what `sealed class` is — explain WHY it's useful (exhaustive when, closed hierarchy, compiler guarantees)

2. **Use concrete examples**: Demonstrate everything with code; interviewers want to see you can write Kotlin, not just talk about it

3. **Know the pitfalls**: Interviewers love to ask about val vs immutability, smart cast limitations, and null safety edge cases

4. **Understand the internals**: Extension functions are static methods. Data classes generate 5+ methods. Coroutines are state machines. This shows depth.

5. **Connect to Java**: Know how Kotlin maps to Java — this is usually relevant for teams migrating from Java

6. **Performance awareness**: Know when to use sequences vs collections, when inline matters, and how boxing affects generics

---

*Next: [Appendix — Cheat Sheets and References](17-appendix.md)*
