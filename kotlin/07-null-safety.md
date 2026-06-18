# Chapter 7 — Null Safety: Kotlin's Superpower

> *"I call it my billion-dollar mistake. It was the invention of the null reference in 1965."*
> — Tony Hoare

---

## 7.1 The Billion-Dollar Mistake

Tony Hoare invented null references in 1965 for ALGOL W, and later estimated that the resulting bugs and vulnerabilities have cost the software industry over a billion dollars in total damage. Every Java, C, C++, Python, and Ruby developer has encountered the `NullPointerException` (NPE) — the crash caused by attempting to use a null reference.

### The Problem in Java

```java
// Java — every object reference could be null at any time
String name = getUsername();  // Could return null
int length = name.length();   // NullPointerException if name is null!

// Java developers write defensive checks everywhere
if (name != null) {
    int length = name.length();
    // But they still forget sometimes
}

// Even with checks, null can travel through your code unexpectedly
User user = database.findUser(id);  // null if not found?
String email = user.getEmail();     // NPE if user is null
String domain = email.split("@")[1]; // NPE if email is null
```

### The Kotlin Solution

Kotlin addresses null at the **type system level**. A type is either nullable (can hold null) or non-null (guaranteed to never be null).

```kotlin
// Non-null: compiler guarantees this is never null
var name: String = "Alice"
// name = null  // COMPILE ERROR — String cannot hold null

// Nullable: explicitly allows null
var maybeNull: String? = null
// maybeNull = "Bob"  // OK

// Compiler catches the mistake before runtime:
val length = maybeNull.length  // COMPILE ERROR — need null check first
```

The NPE can still happen in Kotlin, but only in a few specific situations:
- Calling `!!` (not-null assertion) on a null value
- Java interop (platform types)
- Uninitialized `lateinit var`

---

## 7.2 Nullable vs Non-Null Types

Every type in Kotlin has two variants: nullable (`T?`) and non-null (`T`).

```kotlin
val nonNull: String = "Hello"     // String — never null
val nullable: String? = null      // String? — can be null or String

val num: Int = 42                 // Int — never null
val maybeNum: Int? = null         // Int? — can be null or Int

val list: List<String> = listOf("a", "b")  // never null
val maybeList: List<String>? = null        // can be null

// Nullable Int? vs Int — are different types
fun double(n: Int): Int = n * 2         // takes non-null
fun maybeSomething(n: Int?): Int? = n?.let { it * 2 }  // takes nullable, returns nullable
```

### The Type Hierarchy

```
                  Any?
                  / \
                Any  null
               / | \
          String Int Boolean (and all other types)
```

`Any?` is the supertype of all types (including nullable). `Any` is the supertype of all non-null types.

```kotlin
val a: Any = "Hello"  // OK — String is a subtype of Any
val b: Any = 42       // OK — Int is a subtype of Any
// val c: Any = null  // COMPILE ERROR — null is not a subtype of Any

val d: Any? = null    // OK — Any? accepts null
val e: Any? = "Hi"   // OK — Any? accepts non-null too
```

### Null Safety Is Compile-Time

```kotlin
fun getLength(s: String?): Int {
    // This won't compile — must handle null
    // return s.length  // ERROR: Only safe (?.) or non-null asserted (!!.) calls
    
    // These compile:
    return s?.length ?: 0      // Option 1: safe call with default
    // return s!!.length       // Option 2: assert non-null (risky)
    // return if (s != null) s.length else 0  // Option 3: explicit check
}
```

---

## 7.3 Safe Calls (?.)

The safe call operator `?.` calls a method or accesses a property only if the receiver is not null. If the receiver IS null, the entire expression returns null.

```kotlin
val name: String? = "Alice"
val nullName: String? = null

println(name?.length)      // 5
println(nullName?.length)  // null

// The result type of a safe call is always nullable
val length: Int? = name?.length  // Int? not Int
```

### Chaining Safe Calls

Safe calls can be chained. The chain short-circuits at the first null:

```kotlin
data class Address(val street: String?, val city: String?)
data class Person(val name: String, val address: Address?)

val person = Person("Alice", Address("123 Main St", "Springfield"))
val nullAddressPerson = Person("Bob", null)
val noStreetPerson = Person("Charlie", Address(null, "Portland"))

// Chain of safe calls
println(person.address?.city)              // Springfield
println(nullAddressPerson.address?.city)   // null (address is null)
println(noStreetPerson.address?.street)    // null (street is null)

// Deep chain
val city: String? = person.address?.city
println(city?.uppercase())  // SPRINGFIELD
```

### Safe Calls with Methods

```kotlin
val text: String? = "Hello, World!"
val nullText: String? = null

println(text?.length)                 // 13
println(text?.uppercase())            // HELLO, WORLD!
println(text?.replace(",", ""))       // Hello World!
println(nullText?.uppercase())        // null
println(text?.filter { it.isLetter() })  // HelloWorld
```

### Safe Calls with let

Combining `?.` with `let` allows running a block only when the value is non-null:

```kotlin
val name: String? = "Alice"

// Run block only if name is non-null
name?.let { nonNullName ->
    println("Processing: $nonNullName")
    println("Length: ${nonNullName.length}")
}

// With implicit 'it'
name?.let {
    println("Name is not null: $it")
}

// Chaining
val nullableName: String? = null
nullableName?.let { println("This won't print") }
```

---

## 7.4 The Elvis Operator (?:)

The Elvis operator `?:` provides a **default value** when the expression on the left is null:

```kotlin
val name: String? = null
val displayName = name ?: "Anonymous"  // "Anonymous" if name is null
println(displayName)  // Anonymous

val nonNull: String? = "Alice"
val display = nonNull ?: "Default"  // "Alice" — not null, so Elvis not used
println(display)  // Alice
```

### Elvis in Expressions

```kotlin
fun getLength(s: String?): Int = s?.length ?: 0

println(getLength(null))    // 0
println(getLength("Hello")) // 5

// Chaining
val len: Int = text?.trim()?.length ?: 0

// Elvis with complex expressions
val user: User? = getUser()
val name: String = user?.name ?: user?.username ?: "Guest"
```

### Elvis for Early Return or Throw

The right side of `?:` can be a `return` or `throw`:

```kotlin
fun processUser(userId: String): String {
    val user = findUser(userId) ?: return "User not found"
    val email = user.email ?: throw IllegalStateException("User has no email")
    return "Processing $email"
}

// Pattern: validate at the start of a function
fun processOrder(orderId: String, userId: String) {
    val order = findOrder(orderId) ?: return  // early return if not found
    val user = findUser(userId) ?: throw IllegalArgumentException("Invalid user")
    
    // From here on, order and user are non-null
    println("Processing order ${order.id} for ${user.name}")
}
```

---

## 7.5 The Not-Null Assertion (!!)

The `!!` operator **asserts** that a value is not null. If it IS null, it throws a `NullPointerException`.

```kotlin
val name: String? = "Alice"
val length = name!!.length  // safe — name is not null

val nullName: String? = null
// val crashLength = nullName!!.length  // throws NullPointerException at runtime!
```

### When to Use !!

The `!!` is **deliberately verbose and ugly**. This is intentional — Kotlin is signaling "here be dragons."

Use `!!` only when:
1. You are **absolutely certain** the value cannot be null at this point
2. You want to **signal a bug** if it ever is null (a programming error, not an expected condition)
3. Working with Java interop where nullability is unknown

```kotlin
// ACCEPTABLE: you know the database always returns this ID if the user exists
val user = database.findById(id)  // returns User?
val verifiedUser = user!!  // acceptable if you KNOW the ID is valid

// BETTER: use an assertion with a message
val verifiedUser = checkNotNull(user) { "User with ID $id not found" }
// throws IllegalStateException with message if null, instead of bare NPE
```

### !! Anti-Patterns

```kotlin
// BAD: using !! when you could use ?. and ?:
val name: String? = getName()
println(name!!.length)      // bad — NPE if null
println(name?.length ?: 0)  // good — handled gracefully

// BAD: double !! in a chain
user!!.address!!.street!!.length  // cascade of potential NPEs

// GOOD: safe chain
user?.address?.street?.length ?: 0
```

---

## 7.6 Smart Casts

Kotlin's **smart cast** system tracks when a nullable value has been checked for null. After a check, the compiler treats the value as non-null.

### Smart Cast with if

```kotlin
fun printLength(name: String?) {
    if (name != null) {
        // Inside this block, name is smart-cast to String (non-null)
        println(name.length)  // no ?. needed
        println(name.uppercase())
    }
}

// Smart cast on the else branch
fun printOrDefault(text: String?) {
    if (text == null) {
        println("No text provided")
        return
    }
    // After returning early for null, text is smart-cast to non-null here
    println(text.trim().uppercase())
}
```

### Smart Cast with &&

```kotlin
fun process(s: String?) {
    if (s != null && s.isNotEmpty()) {
        // Both conditions checked — s is String here
        println(s.uppercase())
    }
}
```

### Smart Cast with when

```kotlin
fun describe(value: Any?): String = when {
    value == null      -> "null"
    value is Int       -> "Int: $value"       // smart-cast to Int
    value is String    -> "String: $value"     // smart-cast to String
    value is List<*>   -> "List of ${value.size} items"  // smart-cast to List
    else               -> "Unknown: $value"
}

println(describe(null))        // null
println(describe(42))          // Int: 42
println(describe("Hello"))     // String: Hello
println(describe(listOf(1,2))) // List of 2 items
```

### Smart Cast Requirements

Smart casts only work when the compiler can guarantee the value hasn't changed between the check and the use:

```kotlin
var name: String? = "Alice"

// DOES NOT WORK for var — compiler can't guarantee 'name' wasn't changed by another thread
// if (name != null) {
//     println(name.length)  // Still treats as String? — var could change
// }

// Works if you capture in a val
val safeName = name
if (safeName != null) {
    println(safeName.length)  // smart-cast works — val can't change
}

// Works directly with val properties
class User(val name: String?)
val user = User("Alice")
if (user.name != null) {
    println(user.name.length)  // smart-cast — val property is stable
}
```

---

## 7.7 lateinit var

Sometimes you can't initialize a property in the constructor — it gets initialized later (e.g., by a dependency injection framework). `lateinit` defers the initialization obligation:

```kotlin
class MyClass {
    lateinit var name: String
    
    fun initialize() {
        name = "Initialized!"
    }
    
    fun printName() {
        println(name)
    }
}

val obj = MyClass()
// println(obj.name)  // throws UninitializedPropertyAccessException
obj.initialize()
obj.printName()  // Initialized!
```

### isInitialized Check

```kotlin
class Service {
    lateinit var connection: DatabaseConnection
    
    fun isReady() = ::connection.isInitialized  // property reference check
    
    fun connect() {
        connection = DatabaseConnection()
    }
    
    fun query(sql: String): List<Any> {
        if (!::connection.isInitialized) {
            throw IllegalStateException("Call connect() first")
        }
        return connection.execute(sql)
    }
}
```

### lateinit Rules and Restrictions

- Only works with `var` (not `val`)
- Only works with non-null types
- Only works with non-primitive types (can't use `lateinit var x: Int`)
- If accessed before initialized: `UninitializedPropertyAccessException`

```kotlin
class Example {
    lateinit var text: String  // OK
    // lateinit var count: Int  // ERROR — Int is a primitive type
    // lateinit val name: String  // ERROR — must be var
    // lateinit var maybe: String? = null  // ERROR — must be non-null type
}
```

---

## 7.8 lazy Initialization

`lazy` creates a property that is initialized on first access and then cached. Unlike `lateinit`, `lazy` works with `val`.

```kotlin
class ExpensiveResource {
    val database: DatabaseConnection by lazy {
        println("Opening database connection...")
        DatabaseConnection()  // created only once, on first access
    }
}

val resource = ExpensiveResource()
println("Resource created, but database not connected yet")
// Accessing .database for the first time:
println(resource.database)  // prints "Opening database connection..."
println(resource.database)  // uses cached value — no second initialization
```

### lazy is Thread-Safe by Default

By default, `lazy` uses `LazyThreadSafetyMode.SYNCHRONIZED`, ensuring the initialization block runs at most once even under concurrent access:

```kotlin
val heavyObject by lazy {
    // Thread-safe by default — initialized exactly once
    computeExpensiveValue()
}

// For single-threaded code, you can opt out of the overhead:
val fastLazy by lazy(LazyThreadSafetyMode.NONE) {
    computeValue()  // faster, but not thread-safe
}
```

### lazy vs lateinit

| Aspect | `lazy` | `lateinit` |
|--------|--------|------------|
| Declaration | `val` | `var` |
| Initialization | On first access | Manually, any time |
| Type restriction | Any type | Non-primitive, non-null only |
| Thread safety | Synchronized by default | No |
| Access before init | Returns initialized value | Throws exception |

---

## 7.9 Null Safety Best Practices

### Practice 1: Prefer Non-Null Types

Design your code to use non-null types by default. Only use nullable types when null has meaningful semantic value.

```kotlin
// BAD: Using null when a default makes more sense
fun getGreeting(name: String?): String? = name?.let { "Hello, $it!" }

// GOOD: Use a default value
fun getGreeting(name: String = "Friend"): String = "Hello, $name!"
```

### Practice 2: Push Null Handling to the Edges

Deal with nullable types at the point where they enter your system (from APIs, databases, user input). Make your core logic work with non-null types:

```kotlin
// BAD: nullable types flowing through business logic
class OrderService {
    fun processOrder(user: User?, order: Order?) {
        if (user == null || order == null) return
        // process
    }
}

// GOOD: handle null at the boundary
class OrderService {
    fun processOrder(user: User, order: Order) {
        // Always non-null here — caller handles null checks
    }
}

// At the boundary (e.g., API handler):
fun handleRequest(userId: String, orderId: String) {
    val user = userRepo.findById(userId) ?: run {
        sendError("User not found")
        return
    }
    val order = orderRepo.findById(orderId) ?: run {
        sendError("Order not found")
        return
    }
    orderService.processOrder(user, order)  // non-null guaranteed
}
```

### Practice 3: Use requireNotNull and checkNotNull

These functions provide better error messages than `!!`:

```kotlin
// Less helpful: bare NPE
val config = loadConfig()!!

// More helpful: clear error message
val config = requireNotNull(loadConfig()) { "Config file missing or invalid" }

// require — throws IllegalArgumentException (for input validation)
fun divide(a: Int, b: Int): Int {
    require(b != 0) { "Divisor cannot be zero" }
    return a / b
}

// check — throws IllegalStateException (for state validation)
fun processPayment(amount: Double) {
    check(isConnected) { "Payment processor not connected" }
    // process
}
```

### Practice 4: Avoid Double Bangs (!!) in Production Code

Prefer safe calls and Elvis over `!!`. If you see `!!` in code review, it deserves scrutiny.

```kotlin
// RED FLAG
user!!.address!!.city!!.uppercase()

// GREEN
user?.address?.city?.uppercase() ?: "UNKNOWN"
```

### Practice 5: Understand When Null Is Appropriate

Not everything should be non-null. Null is appropriate when:
- An optional configuration value is not provided
- A database record is not found
- A field is genuinely optional in a domain model

```kotlin
data class User(
    val name: String,       // required — never null
    val email: String,      // required — never null
    val phone: String?,     // optional — user may not have provided it
    val middleName: String? // optional — not everyone has one
)
```

### Practice 6: Use Collections Over Nullable Collections

Prefer empty collections over null collections:

```kotlin
// BAD: returning null when no results
fun findUsers(query: String): List<User>? = 
    if (results.isEmpty()) null else results

// GOOD: return empty list
fun findUsers(query: String): List<User> = 
    results.filter { it.matches(query) }

// The caller doesn't need null checks:
val users = findUsers("alice")
users.forEach { println(it.name) }  // Empty forEach if no results — no NPE
```

---

## Comprehensive Example: Null Safety in Practice

```kotlin
data class Address(
    val street: String,
    val city: String,
    val country: String,
    val zipCode: String?  // optional
)

data class User(
    val id: String,
    val name: String,
    val email: String,
    val address: Address?,   // optional
    val phone: String?       // optional
)

class UserRepository {
    private val users = mutableMapOf<String, User>()
    
    fun findById(id: String): User? = users[id]
    
    fun save(user: User) { users[user.id] = user }
}

fun formatUserLocation(user: User): String {
    // Chain of safe calls
    val city = user.address?.city ?: "Unknown city"
    val country = user.address?.country ?: "Unknown country"
    return "$city, $country"
}

fun getUserContact(userId: String, repo: UserRepository): String {
    val user = repo.findById(userId) ?: return "User not found"
    
    // Smart cast — user is non-null after the ?: return
    val contactInfo = when {
        user.email.isNotEmpty() -> "Email: ${user.email}"
        user.phone != null -> "Phone: ${user.phone}"
        else -> "No contact info"
    }
    
    return "${user.name} — $contactInfo (${formatUserLocation(user)})"
}

fun main() {
    val repo = UserRepository()
    
    repo.save(User(
        id = "1",
        name = "Alice Smith",
        email = "alice@example.com",
        address = Address("123 Main St", "Springfield", "US", "12345"),
        phone = "+1-555-0100"
    ))
    
    repo.save(User(
        id = "2",
        name = "Bob Jones",
        email = "bob@example.com",
        address = null,  // no address provided
        phone = null
    ))
    
    println(getUserContact("1", repo))
    // Alice Smith — Email: alice@example.com (Springfield, US)
    
    println(getUserContact("2", repo))
    // Bob Jones — Email: bob@example.com (Unknown city, Unknown country)
    
    println(getUserContact("999", repo))
    // User not found
}
```

---

## Summary

Kotlin's null safety system is enforced by the compiler at the type level. Non-null types (`String`) can never hold null; nullable types (`String?`) can. Safe calls (`?.`) return null instead of throwing; the Elvis operator (`?:`) provides defaults; the not-null assertion (`!!`) is a last resort that throws NPE. Smart casts eliminate redundant null checks when the compiler can verify safety. `lateinit var` defers initialization of non-null types; `lazy` initializes on first access and caches the result.

---

## Key Takeaways

- Kotlin's type system distinguishes `String` (non-null) from `String?` (nullable)
- Safe call `?.` returns null instead of throwing when receiver is null
- Elvis `?:` provides a default value (or `return`/`throw`) when left side is null
- `!!` is intentionally ugly — it signals "trust me, this is not null"
- Smart casts eliminate `?.` after compiler-verified null checks
- `lateinit var` for late-initialized non-null types (non-primitive only)
- `by lazy` for properties initialized on first access, cached forever
- Prefer non-null types; handle null at system boundaries; avoid `!!` in business logic

---

## Practice Questions

### Conceptual
1. What is the difference between `String` and `String?` in Kotlin?
2. When does a safe call (`?.`) return null vs the actual value?
3. What does the Elvis operator do when the left side is non-null?
4. Why does smart cast fail for `var` properties?
5. What is the difference between `lateinit var` and `by lazy`?

### Code Exercises

**Exercise 1:** Write a function `safeLength(s: String?): Int` using:
- Version A: Safe call + Elvis
- Version B: Explicit null check
- Version C: Smart cast pattern

**Exercise 2:** Given:
```kotlin
data class Config(val host: String?, val port: Int?, val timeout: Int?)
val config: Config? = loadConfig()  // may return null
```
Write code to safely extract host, port, and timeout with defaults: "localhost", 8080, 30.

**Exercise 3:** Convert this "nullable-hell" chain to idiomatic Kotlin:
```kotlin
val city = if (user != null && user.address != null && user.address.city != null)
    user.address.city.uppercase()
else
    "UNKNOWN"
```

**Exercise 4:** Implement a simple cache class using `lateinit var` for a backing store that must be initialized before use. Add an `isInitialized` check function.

**Exercise 5:** Identify the issues with this code and fix them:
```kotlin
fun processUser(user: User?) {
    val name = user!!.name!!  // issue 1
    val address = user!!.address!!.city!!  // issue 2
    println("Processing $name from $address")
}
```

---

*Next: [Chapter 8 — Functional Programming in Kotlin](08-functional-programming.md)*
