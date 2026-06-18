# Best Practices — Writing Idiomatic Kotlin

> *"Idiomatic code is code that uses the language's features in the way the designers intended."*

---

## Overview

Writing Kotlin that "compiles and works" is easy. Writing Kotlin that is **idiomatic, readable, safe, and efficient** takes understanding of the language's conventions and philosophy. This chapter collects the most important best practices for writing professional Kotlin code.

---

## 1. Variable Declarations

### Use val by Default

```kotlin
// BAD — using var when val suffices
var name = "Alice"    // never reassigned below
var count = items.size

// GOOD — prefer val
val name = "Alice"
val count = items.size
```

### Let the Compiler Infer Types

```kotlin
// BAD — redundant type declarations
val name: String = "Alice"
val list: List<String> = listOf("a", "b")

// GOOD — inferred types
val name = "Alice"
val list = listOf("a", "b")

// EXCEPTION — declare type when it adds clarity or widens the type
val timeout: Long = 30_000    // makes it clear this is milliseconds (Long)
val base: Number = 42          // widening is intentional
```

### Use Descriptive Names

```kotlin
// BAD
val n = users.size
val l = items.filter { it.active }
fun f(x: String) = x.length

// GOOD
val userCount = users.size
val activeItems = items.filter { it.isActive }
fun stringLength(text: String) = text.length
```

---

## 2. Null Safety

### Prefer Non-Null Types

```kotlin
// BAD — unnecessary nullable
fun greet(name: String?): String = "Hello, ${name ?: "Guest"}!"

// GOOD — use a default in the signature
fun greet(name: String = "Guest"): String = "Hello, $name!"
```

### Handle Null at Boundaries

```kotlin
// BAD — nullable flowing through business logic
fun processOrder(userId: String?, orderId: String?) {
    if (userId == null || orderId == null) return
    val user = findUser(userId) ?: return
    // ...
}

// GOOD — non-null in business logic, null handled at entry point
fun handleOrderRequest(userId: String?, orderId: String?) {
    val resolvedUserId = userId ?: run { logError("Missing user ID"); return }
    val resolvedOrderId = orderId ?: run { logError("Missing order ID"); return }
    processOrder(resolvedUserId, resolvedOrderId)
}

fun processOrder(userId: String, orderId: String) {
    // all non-null from here
}
```

### Avoid !! in Production Code

```kotlin
// BAD — !! everywhere
val name = user!!.profile!!.name!!

// GOOD — safe chain with defaults
val name = user?.profile?.name ?: "Unknown"

// ACCEPTABLE !! only when you KNOW it's non-null and want a crash otherwise
val config = loadConfig()
    ?: throw IllegalStateException("Configuration file is required but missing")
// OR
val config = checkNotNull(loadConfig()) { "Config required" }
```

### Return Empty Collections, Not Null

```kotlin
// BAD
fun findUsers(query: String): List<User>? = 
    if (results.isEmpty()) null else results

// GOOD — empty list is never null
fun findUsers(query: String): List<User> = 
    users.filter { query.lowercase() in it.name.lowercase() }
```

---

## 3. Functions

### Use Expression Bodies for Simple Functions

```kotlin
// BAD — unnecessary block body
fun double(x: Int): Int {
    return x * 2
}

// GOOD — expression body
fun double(x: Int) = x * 2

// Also good for slightly complex expressions
fun max(a: Int, b: Int) = if (a > b) a else b
fun describe(n: Int) = when {
    n < 0 -> "negative"
    n == 0 -> "zero"
    else -> "positive"
}
```

### Name Arguments for Clarity

```kotlin
// BAD — ambiguous Boolean arguments
createUser("Alice", "alice@example.com", true, false)

// GOOD — named arguments clarify intent
createUser(
    name = "Alice",
    email = "alice@example.com",
    isAdmin = true,
    isActive = false
)
```

### Use Default Parameters Over Overloads

```kotlin
// BAD — Java-style overloads
fun connect(host: String) = connect(host, 8080)
fun connect(host: String, port: Int) = connect(host, port, 30)
fun connect(host: String, port: Int, timeout: Int) { /* actual logic */ }

// GOOD — default parameters
fun connect(host: String, port: Int = 8080, timeout: Int = 30) {
    /* actual logic */
}
```

---

## 4. Classes and OOP

### Use Data Classes for Data Holders

```kotlin
// BAD — manual class with boilerplate
class Point(val x: Double, val y: Double) {
    override fun equals(other: Any?): Boolean { /* ... */ }
    override fun hashCode(): Int { /* ... */ }
    override fun toString(): String = "Point($x, $y)"
}

// GOOD — data class does it all
data class Point(val x: Double, val y: Double)
```

### Design Classes to Be Final (Don't use open Unnecessarily)

```kotlin
// BAD — open by default thinking
open class DatabaseService {
    fun query(sql: String): List<Any> { /* ... */ }
}

// GOOD — closed unless inheritance is explicitly designed for
class DatabaseService {
    fun query(sql: String): List<Any> { /* ... */ }
}

// When inheritance IS the intent:
abstract class AbstractRepository<T> {
    abstract fun findById(id: Long): T?
    fun findAll(): List<T> = /* default implementation */
}
```

### Use Sealed Classes for Finite Hierarchies

```kotlin
// BAD — open class with unclear set of subtypes
open class ApiResult
class Success(val data: String) : ApiResult()
class Failure(val error: String) : ApiResult()
// Someone could add more subclasses anywhere!

// GOOD — sealed class = closed, exhaustive
sealed class ApiResult {
    data class Success(val data: String) : ApiResult()
    data class Failure(val error: String, val code: Int) : ApiResult()
    object Loading : ApiResult()
}

// when is exhaustive — no else needed
fun handle(result: ApiResult) = when (result) {
    is ApiResult.Success -> show(result.data)
    is ApiResult.Failure -> showError(result.error)
    ApiResult.Loading    -> showSpinner()
}
```

### Prefer Composition Over Inheritance

```kotlin
// BAD — fragile inheritance
open class Logger {
    open fun log(msg: String) = println(msg)
}

class PrefixLogger(val prefix: String) : Logger() {
    override fun log(msg: String) = super.log("$prefix: $msg")
}

// GOOD — composition is clearer and more flexible
class PrefixLogger(
    private val delegate: Logger,
    private val prefix: String
) {
    fun log(msg: String) = delegate.log("$prefix: $msg")
}

// Or with class delegation:
interface Loggable {
    fun log(msg: String)
}

class PrefixLoggable(
    private val delegate: Loggable,
    private val prefix: String
) : Loggable by delegate {
    override fun log(msg: String) = delegate.log("$prefix: $msg")
}
```

---

## 5. Functional Style

### Use Collection Operations Instead of Loops

```kotlin
// BAD — imperative style
val result = mutableListOf<String>()
for (user in users) {
    if (user.isActive && user.age >= 18) {
        result.add(user.name.uppercase())
    }
}

// GOOD — declarative, functional style
val result = users
    .filter { it.isActive && it.age >= 18 }
    .map { it.name.uppercase() }
```

### Prefer Functional Transformations for Readability

```kotlin
// BAD — manual accumulation
var total = 0.0
for (order in orders) {
    if (order.status == "completed") {
        total += order.amount
    }
}

// GOOD — clear, one-liner
val total = orders
    .filter { it.status == "completed" }
    .sumOf { it.amount }
```

### Use Sequences for Large Data

```kotlin
// BAD — eager evaluation creates intermediate collections
val result = (1..10_000_000)
    .toList()
    .map { it * 2 }
    .filter { it % 3 == 0 }
    .take(10)

// GOOD — lazy evaluation with sequences
val result = (1..10_000_000)
    .asSequence()
    .map { it * 2 }
    .filter { it % 3 == 0 }
    .take(10)
    .toList()
```

---

## 6. Scope Functions — Use the Right Tool

```kotlin
// apply — configure an object, return object
val request = HttpRequest().apply {
    url = "https://api.example.com"
    method = "POST"
    body = """{"key": "value"}"""
}

// let — transform, or run when non-null
val length = username?.let { it.trim().length } ?: 0

// also — side effects (logging, debugging)
val user = createUser()
    .also { logger.log("Created user: ${it.id}") }
    .also { analytics.track("user_created") }

// run — execute block and return result (useful for scoping)
val config = run {
    val rawValue = System.getProperty("config.path")
    val path = rawValue ?: "/etc/app/config"
    Config.load(path)
}

// with — multiple operations on same object, return result
val stats = with(numbers) {
    Stats(
        count = size,
        sum = sum(),
        average = average(),
        min = min(),
        max = max()
    )
}
```

---

## 7. Extension Functions

### Add Extensions to Clarify Intent, Not to be Clever

```kotlin
// GOOD — adds genuine clarity
fun String.isValidEmail() = contains("@") && contains(".")

fun LocalDate.isWeekend() = dayOfWeek == DayOfWeek.SATURDAY || 
                            dayOfWeek == DayOfWeek.SUNDAY

fun <T> List<T>.randomElement(): T = this[Random.nextInt(size)]

// BAD — too clever, hurts readability
fun Int.timesRepeat(block: () -> Unit) = repeat(this) { block() }
5.timesRepeat { println("hello") }  // reads strangely
// Just use: repeat(5) { println("hello") }
```

### Prefer Member Functions When Modifying State

```kotlin
// BAD — extension that mutates internal state
fun Counter.reset() { this.count = 0 }  // should be a member function

// GOOD — use member function for state mutation
class Counter {
    var count = 0
    fun reset() { count = 0 }  // member — has direct access and clear ownership
}

// Extensions are good for adding functionality without access to internals
fun Counter.toPercent(total: Int) = count.toDouble() / total * 100
```

---

## 8. Coroutines

### Prefer coroutineScope over GlobalScope

```kotlin
// BAD — GlobalScope coroutine is unscoped, might leak
fun processData() {
    GlobalScope.launch {
        // ...
    }
}

// GOOD — use coroutineScope or a properly scoped CoroutineScope
suspend fun processData() = coroutineScope {
    launch { /* ... */ }
}
```

### Use async/await for Concurrent Independent Operations

```kotlin
// BAD — sequential when operations are independent
suspend fun getProfile(): Profile {
    val user = fetchUser()        // waits
    val orders = fetchOrders()    // waits after user
    val prefs = fetchPrefs()      // waits after orders
    return Profile(user, orders, prefs)
}

// GOOD — concurrent independent operations
suspend fun getProfile(): Profile = coroutineScope {
    val user = async { fetchUser() }
    val orders = async { fetchOrders() }
    val prefs = async { fetchPrefs() }
    Profile(user.await(), orders.await(), prefs.await())
}
```

### Handle Exceptions Explicitly in Coroutines

```kotlin
// BAD — exception is silently ignored
launch {
    riskyOperation()  // if this throws, it propagates to Job handler
}

// GOOD — handle exceptions explicitly
launch {
    try {
        riskyOperation()
    } catch (e: NetworkException) {
        logger.error("Network failed", e)
        showRetryUI()
    }
}

// Or use supervisorScope for independent child failures:
supervisorScope {
    launch { riskyOperation1() }  // failure doesn't cancel sibling
    launch { riskyOperation2() }  // each handles its own failure
}
```

---

## 9. Coding Style

### Naming Conventions

```kotlin
// Classes: PascalCase
class UserRepository
class HttpClient
data class ApiResponse

// Functions and properties: camelCase
fun calculateTax(income: Double): Double
val userCount: Int
var isConnected: Boolean

// Constants: SCREAMING_SNAKE_CASE (in companion objects or top-level)
const val MAX_RETRY_COUNT = 3
const val DEFAULT_TIMEOUT_MS = 5000L

// Packages: lowercase.dot.separated
package com.example.utils

// Type parameters: single uppercase letter or descriptive name
fun <T> process(item: T): T
fun <Key, Value> buildMap(): Map<Key, Value>
```

### Formatting

```kotlin
// Line length: keep under 100 characters
// Function with many parameters — one per line
fun createReport(
    title: String,
    author: String,
    date: LocalDate,
    sections: List<Section>,
    footer: String = ""
): Report { /* ... */ }

// Chained calls — each on its own line
val result = items
    .filter { it.isActive }
    .map { it.name }
    .sorted()
    .joinToString(", ")

// Trailing commas are idiomatic (and help diffs)
val colors = listOf(
    "red",
    "green",
    "blue",   // trailing comma OK
)
```

---

## 10. Common Idioms

### Idiomatic Patterns You Should Know

```kotlin
// Swap without temp variable
var a = 1
var b = 2
a = b.also { b = a }

// Conditional assignment
val status = if (isLoggedIn) "active" else "anonymous"

// Execute once
val initialized by lazy { initialize() }

// Convert nullable to empty list
val items = nullableItems ?: emptyList()

// Get or compute default
val value = map.getOrPut(key) { computeDefault() }

// Safe navigation chain
val city = user?.address?.city?.uppercase()

// Scope for multi-line initialization
val config = Config().apply {
    timeout = 30
    maxRetries = 3
}

// Multiple assignments from function
val (user, token) = loginUser(credentials)  // requires component functions

// Filtering nulls
val nonNulls = listOfNullables.filterNotNull()

// Check all conditions
val valid = listOf(isNotEmpty, isUnique, isAuthorized).all { it }
```

---

## Summary Checklist

Use this as a code review guide:

- [ ] Variables: prefer `val` over `var`
- [ ] Nullability: non-null by default; handle null at boundaries
- [ ] Functions: use expression bodies where appropriate
- [ ] Functions: use named arguments when calling with multiple booleans or same-type params
- [ ] Classes: data classes for data holders
- [ ] Classes: sealed classes for closed hierarchies
- [ ] Classes: `open` only when inheritance is explicitly intended
- [ ] Collections: use `filter`/`map`/`fold` over loops for transformations
- [ ] Sequences: use `asSequence()` for large data or chains with early termination
- [ ] Coroutines: use `coroutineScope`, not `GlobalScope`
- [ ] Coroutines: use `async`/`await` for concurrent independent work
- [ ] Style: follow Kotlin naming conventions
- [ ] Avoid: `!!` in business logic
- [ ] Avoid: mutable state when `val` and functional style work
- [ ] Avoid: null collections — return empty collections instead

---

*Next: [Common Pitfalls](15-common-pitfalls.md)*
