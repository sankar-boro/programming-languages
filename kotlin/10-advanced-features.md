# Chapter 10 — Extensions, Delegation, and Generics

> *"The most powerful abstractions in software engineering are those that let you extend systems without modifying them."*

---

## 10.1 Extension Functions

Extension functions allow you to **add new functions to existing classes** without inheriting from them or modifying their source code. This is one of Kotlin's most expressive features.

### Basic Extension Functions

```kotlin
// Adding a function to String
fun String.isPalindrome(): Boolean {
    val cleaned = this.lowercase().filter { it.isLetter() }
    return cleaned == cleaned.reversed()
}

println("racecar".isPalindrome())     // true
println("Hello".isPalindrome())       // false
println("A man a plan a canal Panama".isPalindrome())  // true

// Adding to Int
fun Int.isEven() = this % 2 == 0
fun Int.factorial(): Long {
    require(this >= 0) { "Factorial not defined for negative numbers" }
    return if (this == 0) 1L else this.toLong() * (this - 1).factorial()
}

println(4.isEven())     // true
println(5.factorial())  // 120

// Adding to a list
fun <T> List<T>.secondOrNull(): T? = if (size >= 2) this[1] else null
fun <T> List<T>.swap(i: Int, j: Int): List<T> {
    val result = toMutableList()
    val temp = result[i]
    result[i] = result[j]
    result[j] = temp
    return result
}

println(listOf(1, 2, 3).secondOrNull())  // 2
println(listOf<Int>().secondOrNull())    // null
println(listOf(1, 2, 3, 4).swap(0, 3)) // [4, 2, 3, 1]
```

### How Extension Functions Work

Extension functions don't actually modify the class. They're syntactic sugar for static functions:

```kotlin
// You write:
fun String.shout() = this.uppercase() + "!!!"

// Kotlin compiles to (approximately):
// static String shout(String $receiver) { return $receiver.toUpperCase() + "!!!"; }

// So this:
"hello".shout()

// Becomes this under the hood:
// shout("hello")
```

This means:
- Extension functions can only access the **public API** of the class
- They cannot override existing member functions
- If there's a member function with the same signature, the **member always wins**

```kotlin
class MyClass {
    fun greet() = "Member greet"
}

fun MyClass.greet() = "Extension greet"

val obj = MyClass()
println(obj.greet())  // Member greet — member wins!
```

### Extension Functions on Nullable Types

You can write extension functions on nullable types:

```kotlin
fun String?.orDefault(default: String = ""): String = this ?: default

val name: String? = null
println(name.orDefault("Anonymous"))  // Anonymous

// The standard library does this for you:
println(null.toString())  // "null" — extension on Any?
```

### Extension Functions in the Standard Library

Kotlin's standard library is built heavily on extension functions. This is why you can call methods on `String`, `List`, `Int`, etc. that don't exist in the Java API:

```kotlin
// These are all extension functions:
"Hello".reversed()        // extension on String
listOf(1,2,3).sum()       // extension on Iterable<Int>
42.coerceIn(0, 100)       // extension on Comparable<Int>
"hello".capitalize()      // extension on String
```

---

## 10.2 Extension Properties

Properties can be added to existing classes too:

```kotlin
val String.lastChar: Char
    get() = this[this.length - 1]

val String.wordCount: Int
    get() = this.trim().split(Regex("\\s+")).size

val <T> List<T>.penultimate: T?
    get() = if (size >= 2) this[size - 2] else null

println("Kotlin".lastChar)          // n
println("Hello World Kotlin".wordCount)  // 3
println(listOf(1, 2, 3).penultimate)     // 2
```

Extension properties cannot have backing fields — they're always computed from the receiver.

---

## 10.3 Scope Functions

Scope functions (`let`, `run`, `with`, `apply`, `also`) are higher-order extension functions that execute a block in the context of an object.

### let

Transforms an object or executes on a non-null value.
- Receiver reference: `it`
- Returns: lambda result

```kotlin
val name: String? = "Alice"

// Null check + transformation
val length = name?.let { it.length } ?: 0

// Transformation chain
val result = "Hello World"
    .let { it.split(" ") }
    .let { words -> words.map { it.length } }
    .let { lengths -> lengths.sum() }

println(result)  // 10

// Convert to a different type
val user = getUser()?.let { user ->
    UserDTO(user.name, user.email)
}
```

### run

Execute a block and return the result. Good for initialization.
- Receiver reference: `this`
- Returns: lambda result

```kotlin
val text = "Hello, World!"

val wordCount = text.run {
    val words = split(" ")
    words.size
}

// Or without a receiver — just a block for scoping
val result = run {
    val x = computeX()
    val y = computeY()
    x + y  // result of the block
}
```

### with

Execute multiple operations on an object.
- Takes the object as a parameter
- Receiver reference: `this`
- Returns: lambda result

```kotlin
val numbers = listOf(1, 2, 3, 4, 5)

val description = with(numbers) {
    val total = sum()
    val avg = average()
    "Count: $size, Sum: $total, Avg: $avg"
}
println(description)  // Count: 5, Sum: 15, Avg: 3.0
```

### apply

Configure an object; return the object itself.
- Receiver reference: `this`
- Returns: the original object

```kotlin
// Builder pattern
val person = Person().apply {
    name = "Alice"
    age = 30
    email = "alice@example.com"
}

// Creating and configuring a collection
val list = mutableListOf<Int>().apply {
    add(1); add(2); add(3)
    addAll(4..6)
}
println(list)  // [1, 2, 3, 4, 5, 6]
```

### also

Perform side effects; return the original object.
- Receiver reference: `it`
- Returns: the original object

```kotlin
val numbers = mutableListOf(1, 2, 3)
    .also { println("Initial: $it") }
    .apply { add(4) }
    .also { println("After add: $it") }
    .apply { sort() }
    .also { println("After sort: $it") }

println("Final: $numbers")
```

### Choosing the Right Scope Function

| Goal | Use |
|------|-----|
| Null-safe transformation / execute on non-null | `let` |
| Initialize + return result | `run` |
| Multiple operations on one object, return result | `with` |
| Configure object, return object | `apply` |
| Side effects, return object | `also` |

---

## 10.4 Class Delegation

Kotlin supports **delegation** as a language feature with the `by` keyword, implementing the delegation design pattern without boilerplate.

### The Problem: Implementing Delegation Manually

```kotlin
// You have a large interface
interface Repository {
    fun findById(id: Int): String
    fun findAll(): List<String>
    fun save(item: String): Boolean
    fun delete(id: Int): Boolean
}

// You want to add logging to an existing implementation
class LoggingRepository(private val delegate: Repository) : Repository {
    override fun findById(id: Int): String {
        println("findById($id)")
        return delegate.findById(id)  // manually delegate
    }
    override fun findAll(): List<String> {
        println("findAll()")
        return delegate.findAll()
    }
    override fun save(item: String): Boolean {
        println("save($item)")
        return delegate.save(item)
    }
    override fun delete(id: Int): Boolean {
        println("delete($id)")
        return delegate.delete(id)
    }
}
```

### Class Delegation with `by`

```kotlin
// Kotlin generates all the delegation boilerplate
class LoggingRepository(private val delegate: Repository) : Repository by delegate {
    override fun findById(id: Int): String {
        println("findById($id)")
        return delegate.findById(id)  // only override what you need
    }
    // All other methods are automatically delegated to 'delegate'
}
```

The `by delegate` clause tells Kotlin: "For all interface methods I don't override, call them on `delegate`."

### Practical Example: Combining Interfaces

```kotlin
interface Printable {
    fun print()
}

interface Saveable {
    fun save(path: String)
}

class Document(val content: String) : Printable, Saveable {
    override fun print() = println("Document: $content")
    override fun save(path: String) = println("Saved to $path")
}

// A wrapper that adds logging via delegation
class LoggingDocument(
    private val doc: Document
) : Printable by doc, Saveable by doc {
    // print() and save() are automatically delegated to doc
    // We can add before/after logic by overriding:
    override fun print() {
        println("[LOG] Printing started")
        doc.print()
        println("[LOG] Printing done")
    }
    // save() still delegates to doc without override
}

val doc = LoggingDocument(Document("Hello, Kotlin!"))
doc.print()
doc.save("output.txt")
```

---

## 10.5 Property Delegation

Properties can also delegate their getter/setter logic to a **delegate** object. The delegate handles the storage and retrieval of the value.

```kotlin
// Delegate pattern: the object with getValue() and setValue()
class Delegate {
    operator fun getValue(thisRef: Any?, property: KProperty<*>): String {
        return "Delegate getValue called"
    }
    
    operator fun setValue(thisRef: Any?, property: KProperty<*>, value: String) {
        println("Setting ${property.name} = $value")
    }
}

class Example {
    var p: String by Delegate()
}

val e = Example()
println(e.p)   // Delegate getValue called
e.p = "Hello"  // Setting p = Hello
```

---

## 10.6 Built-in Delegates

Kotlin's standard library includes several powerful property delegates:

### lazy

Computes and caches the value on first access:

```kotlin
class HeavyComputation {
    val expensiveResult: String by lazy {
        println("Computing...")
        Thread.sleep(1000)  // simulate slow computation
        "The answer is 42"
    }
}

val obj = HeavyComputation()
println("Object created")
println(obj.expensiveResult)  // Computing... The answer is 42
println(obj.expensiveResult)  // The answer is 42 (cached, no "Computing...")
```

### observable

Calls a callback whenever the property is changed:

```kotlin
import kotlin.properties.Delegates

class User {
    var name: String by Delegates.observable("<no name>") { property, old, new ->
        println("${property.name} changed: $old -> $new")
    }
}

val user = User()
user.name = "Alice"    // name changed: <no name> -> Alice
user.name = "Alicia"   // name changed: Alice -> Alicia
```

### vetoable

Like `observable`, but the callback can reject the change:

```kotlin
import kotlin.properties.Delegates

class Temperature {
    var celsius: Double by Delegates.vetoable(20.0) { _, _, new ->
        new >= -273.15  // returns true to accept, false to reject
    }
}

val temp = Temperature()
temp.celsius = 25.0     // accepted
println(temp.celsius)   // 25.0

temp.celsius = -300.0   // rejected (below absolute zero)
println(temp.celsius)   // 25.0 — unchanged
```

### Map Delegation

Properties can be backed by a Map — useful for dynamic or configuration objects:

```kotlin
class Config(map: Map<String, Any>) {
    val host: String by map
    val port: Int by map
    val debug: Boolean by map
}

val config = Config(mapOf(
    "host" to "localhost",
    "port" to 8080,
    "debug" to true
))

println(config.host)   // localhost
println(config.port)   // 8080
println(config.debug)  // true

// Mutable version with MutableMap
class Settings(private val map: MutableMap<String, Any> = mutableMapOf()) {
    var theme: String by map
    var fontSize: Int by map
}

val settings = Settings()
settings.theme = "dark"
settings.fontSize = 14
println(settings.theme)     // dark
println(settings.map)       // {theme=dark, fontSize=14}
```

---

## 10.7 Generics

Generics allow you to write code that works with multiple types while maintaining type safety.

### Basic Generics

```kotlin
// Generic function
fun <T> first(list: List<T>): T = list[0]

println(first(listOf(1, 2, 3)))       // 1
println(first(listOf("a", "b", "c"))) // a

// Generic class
class Box<T>(val value: T) {
    fun unwrap(): T = value
    override fun toString() = "Box($value)"
}

val intBox = Box(42)
val stringBox = Box("Hello")
println(intBox.unwrap())    // 42
println(stringBox.unwrap()) // Hello
```

### Type Constraints (Upper Bounds)

```kotlin
// T must be Comparable<T>
fun <T : Comparable<T>> max(a: T, b: T): T = if (a > b) a else b

println(max(3, 7))         // 7
println(max("apple", "banana"))  // banana

// Multiple constraints with where
fun <T> copyAndSort(list: List<T>): List<T>
    where T : Comparable<T>, T : Any {
    return list.sorted()
}

println(copyAndSort(listOf(3, 1, 4, 1, 5)))  // [1, 1, 3, 4, 5]
```

---

## 10.8 Variance: in and out

Variance defines how generic types relate when their type arguments are related.

### Covariance: out

If `T` is a subtype of `S`, then `Producer<T>` should be a subtype of `Producer<S>`.

Mark a type parameter as `out` to declare it as covariant:

```kotlin
interface Producer<out T> {
    fun produce(): T
}

class AnimalProducer : Producer<Animal> {
    override fun produce(): Animal = Animal()
}

class DogProducer : Producer<Dog> {
    override fun produce(): Dog = Dog()
}

// With out T: DogProducer is a Producer<Animal>
val producer: Producer<Animal> = DogProducer()  // OK — Dog is Animal
val animal: Animal = producer.produce()         // returns Dog, works as Animal
```

`out T` means:
- `T` can only appear in **out positions** (return types)
- `T` cannot appear in **in positions** (parameter types)

### Contravariance: in

If `T` is a supertype of `S`, then `Consumer<T>` should be a subtype of `Consumer<S>`.

Mark a type parameter as `in` to declare it as contravariant:

```kotlin
interface Consumer<in T> {
    fun consume(item: T)
}

class AnimalConsumer : Consumer<Animal> {
    override fun consume(item: Animal) = println("Consuming animal: ${item.name}")
}

// With in T: AnimalConsumer can be used as Consumer<Dog>
val dogConsumer: Consumer<Dog> = AnimalConsumer()  // OK
dogConsumer.consume(Dog("Rex"))  // works — Dog IS an Animal
```

`in T` means:
- `T` can only appear in **in positions** (parameter types)
- `T` cannot appear in **out positions** (return types)

### Use-Site Variance (Projection)

When you can't change the class declaration, you can use projections at the use site:

```kotlin
class MutableBox<T>(var value: T)

// out projection — read-only, allows covariance
fun copy(from: MutableBox<out Animal>, to: MutableBox<Animal>) {
    to.value = from.value  // OK — reading from out-projected box
    // from.value = Animal()  // ERROR — can't write to out-projected box
}

// in projection — write-only, allows contravariance
fun fill(box: MutableBox<in Dog>, dog: Dog) {
    box.value = dog  // OK — writing to in-projected box
}
```

### Star Projection

`*` means "I don't know the type argument and don't need to":

```kotlin
fun printContents(box: Box<*>) {
    // Can read as Any?
    val value: Any? = box.value  // OK
    // box.value = something  // ERROR — can't write to star-projected
}

// Useful for checking type without caring about type argument
fun isBoxOf(obj: Any): Boolean = obj is Box<*>
```

---

## 10.9 Reified Type Parameters

Normally, type parameters are **erased** at runtime (JVM type erasure). `reified` type parameters work around this by making the type available at runtime — but only inside `inline` functions.

```kotlin
// Without reified — type erased, can't use T at runtime
// fun <T> isType(obj: Any): Boolean = obj is T  // ERROR

// With reified — type is available at runtime
inline fun <reified T> isType(obj: Any): Boolean = obj is T

println(isType<String>("hello"))  // true
println(isType<Int>("hello"))     // false
println(isType<List<*>>(listOf(1,2,3)))  // true

// Powerful: filterIsInstance with reified
inline fun <reified T> Iterable<*>.filterType(): List<T> =
    filterIsInstance<T>()

val mixed = listOf(1, "two", 3, "four", 5.0)
println(mixed.filterType<String>())  // [two, four]
println(mixed.filterType<Int>())     // [1, 3]

// Getting the class
inline fun <reified T> typeNameOf(): String = T::class.simpleName ?: "Unknown"

println(typeNameOf<String>())   // String
println(typeNameOf<List<Int>>()) // List
```

---

## 10.10 Type Aliases

Type aliases give new names to existing types — useful for readability:

```kotlin
// Function type aliases
typealias Predicate<T> = (T) -> Boolean
typealias Transform<T, R> = (T) -> R
typealias Handler = (event: String, data: Any) -> Unit

val isEven: Predicate<Int> = { it % 2 == 0 }
val toString: Transform<Int, String> = { it.toString() }

// Collection aliases
typealias UserMap = MutableMap<String, User>
typealias IdList = List<Long>

// Making complex generics readable
typealias Result<T> = Either<Throwable, T>
typealias StringMatrix = List<List<String>>

val matrix: StringMatrix = listOf(
    listOf("a", "b", "c"),
    listOf("d", "e", "f")
)
```

Type aliases are compile-time only — the underlying type is unchanged at runtime.

---

## 10.11 Reflection Basics

Kotlin provides a reflection API to inspect and interact with types and members at runtime.

### KClass

```kotlin
// Getting a class reference
val stringClass: KClass<String> = String::class
val intClass: KClass<Int> = Int::class

println(stringClass.simpleName)      // String
println(stringClass.qualifiedName)   // kotlin.String
println(stringClass.isData)          // false
println(stringClass.isFinal)         // true

// From an instance
val obj = "Hello"
println(obj::class.simpleName)       // String
println(obj::class == String::class) // true

// Checking relationship
println(String::class.isSubclassOf(Any::class))  // true
```

### Data Class Reflection

```kotlin
data class Person(val name: String, val age: Int)

val person = Person("Alice", 30)

// Get all properties
val properties = Person::class.memberProperties
for (prop in properties) {
    println("${prop.name}: ${prop.get(person)}")
}
// age: 30
// name: Alice
```

### KProperty and KFunction

```kotlin
class Counter(var count: Int = 0)

val counter = Counter()
val prop = Counter::count

// Get value via property reference
println(prop.get(counter))  // 0

// Set value via property reference
prop.set(counter, 42)
println(counter.count)  // 42

// Callable references
fun add(a: Int, b: Int) = a + b
val addRef = ::add
println(addRef.call(3, 4))    // 7
println(addRef.invoke(3, 4))  // 7
```

---

## Complete Example: A Simple Dependency Injection Container

```kotlin
import kotlin.reflect.KClass

class SimpleContainer {
    private val bindings = mutableMapOf<KClass<*>, () -> Any>()
    private val singletons = mutableMapOf<KClass<*>, Any>()
    
    inline fun <reified T : Any> bind(noinline factory: () -> T) {
        bindings[T::class] = factory
    }
    
    inline fun <reified T : Any> singleton(noinline factory: () -> T) {
        val existing = singletons[T::class]
        if (existing == null) {
            val instance = factory()
            singletons[T::class] = instance
        }
        bindings[T::class] = { singletons[T::class]!! }
    }
    
    @Suppress("UNCHECKED_CAST")
    inline fun <reified T : Any> resolve(): T {
        val factory = bindings[T::class]
            ?: throw IllegalStateException("No binding for ${T::class.simpleName}")
        return factory() as T
    }
}

interface Logger {
    fun log(message: String)
}

class ConsoleLogger : Logger {
    override fun log(message: String) = println("[LOG] $message")
}

class UserService(private val logger: Logger) {
    fun createUser(name: String) {
        logger.log("Creating user: $name")
    }
}

fun main() {
    val container = SimpleContainer()
    
    container.singleton<Logger> { ConsoleLogger() }
    container.bind<UserService> { UserService(container.resolve()) }
    
    val service = container.resolve<UserService>()
    service.createUser("Alice")  // [LOG] Creating user: Alice
}
```

---

## Summary

Extension functions add behavior to existing types without modifying them. Scope functions (`let`, `run`, `with`, `apply`, `also`) execute blocks in the context of an object with different return behaviors. Class delegation (`by`) automatically delegates interface methods to another object. Property delegation (`by`) allows custom get/set logic, with built-in delegates like `lazy`, `observable`, and `vetoable`. Generics use type parameters for type-safe polymorphism; variance (`in`/`out`) describes how generic types relate to each other. `reified` type parameters preserve type information at runtime inside `inline` functions. Type aliases give readable names to complex types. Reflection allows runtime inspection of types and members.

---

## Key Takeaways

- Extension functions are **static** — they cannot override member functions
- `?.let { }` is the idiomatic null-safe transformation pattern
- `apply` configures and returns the object; `also` runs side effects and returns the object
- Class delegation with `by` eliminates delegation boilerplate
- `lazy` property delegate: `val prop by lazy { init }` — thread-safe cached initialization
- `out T` = covariant (Producer) = T only in return positions
- `in T` = contravariant (Consumer) = T only in parameter positions
- `reified` requires `inline` and enables runtime type checks in generic functions
- Type aliases are zero-cost — same type at runtime

---

## Practice Questions

### Conceptual
1. Why can an extension function only access the public API of a class?
2. What is the difference between `apply` and `also`?
3. When would you use `in` variance vs `out` variance?
4. Why does `reified` require `inline`?
5. What does `Delegates.observable` do?

### Code Exercises

**Exercise 1:** Write extension functions on `String`:
- `words(): List<String>` — split by whitespace
- `titleCase(): String` — capitalize first letter of each word
- `isPalindrome(): Boolean`

**Exercise 2:** Implement a `Cache<K, V>` class using `Delegates.observable` that prints a log message whenever any entry changes.

**Exercise 3:** Write a generic function `swap<T>(list: MutableList<T>, i: Int, j: Int)` using generics. Then write an overloaded version for `IntArray`.

**Exercise 4:** Create a type-safe builder using `apply`:
```kotlin
data class HttpRequest(
    val url: String = "",
    val method: String = "GET",
    val headers: Map<String, String> = emptyMap(),
    val body: String? = null
)
// Usage: request { url = "..."; method = "POST"; body = "..." }
```

**Exercise 5:** Using `reified`, write a function `parse<T>(json: String): T` that parses a simple JSON string into either `String`, `Int`, `Double`, or `Boolean` based on the reified type parameter.

---

*Next: [Chapter 11 — Coroutines and Concurrency](11-coroutines.md)*
