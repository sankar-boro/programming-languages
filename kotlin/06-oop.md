# Chapter 6 — Object-Oriented Programming in Kotlin

> *"Kotlin is not 'Java without the pain.' It's a language that embraces both OOP and FP, and uses the right tool for the right job."*

---

## 6.1 Classes and Objects

A class is a blueprint for creating objects. Kotlin's class syntax is concise but fully featured.

### Basic Class

```kotlin
class Person {
    var name: String = ""
    var age: Int = 0
    
    fun introduce() {
        println("Hi, I'm $name, and I'm $age years old.")
    }
}

fun main() {
    val person = Person()    // No 'new' keyword in Kotlin
    person.name = "Alice"
    person.age = 30
    person.introduce()  // Hi, I'm Alice, and I'm 30 years old.
}
```

Note: Kotlin has **no `new` keyword**. You create objects by calling the class name as if it were a function.

---

## 6.2 Constructors

### Primary Constructor

The **primary constructor** is part of the class header:

```kotlin
class Person(val name: String, val age: Int) {
    fun introduce() = println("Hi, I'm $name, and I'm $age years old.")
}

val alice = Person("Alice", 30)
println(alice.name)   // Alice
println(alice.age)    // 30
alice.introduce()     // Hi, I'm Alice, and I'm 30 years old.
```

When you declare a parameter with `val` or `var` in the primary constructor, it automatically becomes a **property** of the class.

```kotlin
// Without val/var: constructor parameter, not a property
class Example(x: Int) {
    val doubled = x * 2  // x is only accessible during initialization
}

// With val: a read-only property
class Box(val width: Int, val height: Int) {
    val area = width * height  // can use properties in other properties
}
```

### init Block

The `init` block runs as part of the primary constructor:

```kotlin
class Temperature(val celsius: Double) {
    val fahrenheit: Double
    val description: String
    
    init {
        fahrenheit = celsius * 9.0 / 5.0 + 32.0
        description = when {
            celsius < 0    -> "Freezing"
            celsius < 15   -> "Cold"
            celsius < 25   -> "Comfortable"
            celsius < 35   -> "Warm"
            else           -> "Hot"
        }
        
        require(celsius >= -273.15) {
            "Temperature cannot be below absolute zero: $celsius°C"
        }
    }
}

val temp = Temperature(25.0)
println("${temp.celsius}°C = ${temp.fahrenheit}°F (${temp.description})")
// 25.0°C = 77.0°F (Comfortable)

// Temperature(-300.0)  // throws IllegalArgumentException
```

Multiple `init` blocks are allowed and run in order:

```kotlin
class MultiInit(val x: Int) {
    init {
        println("First init block: x = $x")
    }
    
    val doubled = x * 2
    
    init {
        println("Second init block: doubled = $doubled")
    }
}

val m = MultiInit(5)
// First init block: x = 5
// Second init block: doubled = 10
```

### Secondary Constructors

Secondary constructors are defined with `constructor` keyword inside the class body:

```kotlin
class User(val username: String, val email: String, val age: Int) {
    
    // Secondary constructor — must delegate to primary
    constructor(username: String, email: String) : this(username, email, 0)
    
    // Another secondary constructor
    constructor(username: String) : this(username, "$username@default.com")
    
    override fun toString() = "User($username, $email, $age)"
}

println(User("alice", "alice@example.com", 30))  // User(alice, alice@example.com, 30)
println(User("bob", "bob@example.com"))           // User(bob, bob@example.com, 0)
println(User("charlie"))                          // User(charlie, charlie@default.com, 0)
```

**Prefer primary constructors with default parameters** over secondary constructors. Secondary constructors should only be used when you need genuinely different initialization logic.

---

## 6.3 Properties and Backing Fields

Properties in Kotlin are more powerful than Java fields — they come with built-in getter and setter support.

### Simple Properties

```kotlin
class Circle(val radius: Double) {
    val area: Double       // computed property (read-only)
        get() = Math.PI * radius * radius
    
    val circumference: Double
        get() = 2 * Math.PI * radius
}

val circle = Circle(5.0)
println(circle.radius)        // 5.0
println(circle.area)          // 78.53...
println(circle.circumference) // 31.41...
```

### Custom Getter and Setter

```kotlin
class Rectangle(var width: Double, var height: Double) {
    var area: Double
        get() = width * height
        set(value) {
            // When setting area, keep the aspect ratio and scale both dimensions
            val ratio = value / area
            width *= Math.sqrt(ratio)
            height *= Math.sqrt(ratio)
        }
}

val rect = Rectangle(4.0, 3.0)
println(rect.area)   // 12.0
rect.area = 48.0     // setter called
println(rect.width)  // ~8.0
println(rect.height) // ~6.0
```

### The Backing Field

When you need to store a value AND have a custom getter/setter, use the **backing field** (`field` keyword):

```kotlin
class Counter {
    var count: Int = 0
        private set    // setter is private — only class itself can set
    
    fun increment() {
        count++
    }
    
    fun reset() {
        count = 0
    }
}

val counter = Counter()
counter.increment()
counter.increment()
counter.increment()
println(counter.count)  // 3
// counter.count = 0    // ERROR: setter is private
counter.reset()
println(counter.count)  // 0
```

```kotlin
class Person(name: String) {
    var name: String = name
        get() = field.trim().replaceFirstChar { it.uppercase() }  // normalize name
        set(value) {
            field = value.trim()  // 'field' refers to the backing field
        }
}

val p = Person("  alice  ")
println(p.name)  // Alice (trimmed and capitalized)
p.name = "  bob  "
println(p.name)  // Bob
```

### Property Visibility

```kotlin
class BankAccount(initialBalance: Double) {
    var balance: Double = initialBalance
        private set    // only this class can set balance
    
    fun deposit(amount: Double) {
        require(amount > 0) { "Deposit must be positive" }
        balance += amount
    }
    
    fun withdraw(amount: Double) {
        require(amount > 0) { "Withdrawal must be positive" }
        require(amount <= balance) { "Insufficient funds" }
        balance -= amount
    }
}

val account = BankAccount(100.0)
account.deposit(50.0)
println(account.balance)   // 150.0
// account.balance = 0.0   // ERROR: private setter
```

---

## 6.4 Visibility Modifiers

Kotlin has four visibility modifiers, with slightly different meanings than Java:

| Modifier | Class Member | Top-Level |
|----------|-------------|-----------|
| `public` | Visible everywhere (default) | Visible everywhere (default) |
| `private` | Visible within the class | Visible within the file |
| `protected` | Visible in class + subclasses | Not applicable |
| `internal` | Visible within the module | Visible within the module |

```kotlin
class Example {
    public val publicProp = "Everyone can see this"
    private val privateProp = "Only this class"
    protected val protectedProp = "Class and subclasses"
    internal val internalProp = "Anywhere in this module"
    
    private fun helper() { /* ... */ }
    
    fun publicMethod() {
        helper()  // can access private members
    }
}
```

### The `internal` Modifier

`internal` is unique to Kotlin. A module is a set of files compiled together (e.g., a Gradle module or Maven artifact). This is useful for library authors who want to expose APIs within the library but not to consumers:

```kotlin
// In a library module:
internal class InternalHelper {  // Library consumers can't use this
    fun doHelperStuff() { }
}

public class PublicAPI {  // Library consumers CAN use this
    private val helper = InternalHelper()  // Internal class used internally
    
    fun publicFunction() = helper.doHelperStuff()
}
```

---

## 6.5 Inheritance

By default, Kotlin classes are **closed** (cannot be extended). You must explicitly mark a class as `open` to allow inheritance.

```kotlin
open class Animal(val name: String) {
    open fun speak(): String = "..."
    
    fun eat() = println("$name is eating")  // not open — cannot override
}

class Dog(name: String) : Animal(name) {
    override fun speak() = "Woof!"
}

class Cat(name: String) : Animal(name) {
    override fun speak() = "Meow!"
}

val animals: List<Animal> = listOf(Dog("Rex"), Cat("Whiskers"), Dog("Buddy"))
for (animal in animals) {
    println("${animal.name} says: ${animal.speak()}")
}
// Rex says: Woof!
// Whiskers says: Meow!
// Buddy says: Woof!
```

### Why `open` by Default?

In Java, all classes are open by default. Kotlin made the opposite choice — classes are **final by default** — because:
- It prevents unintended extension that breaks the parent class's invariants
- It forces the author to explicitly decide when inheritance is appropriate
- It aligns with Joshua Bloch's "Effective Java" advice: "Design and document for inheritance or else prohibit it"

### Calling the Parent

```kotlin
open class Shape(val color: String) {
    open fun describe(): String = "A $color shape"
}

class Circle(color: String, val radius: Double) : Shape(color) {
    override fun describe(): String {
        val parentDesc = super.describe()  // call parent's implementation
        return "$parentDesc (circle, radius=$radius)"
    }
}

println(Circle("red", 5.0).describe())
// A red shape (circle, radius=5.0)
```

### Overriding Properties

Properties can also be overridden:

```kotlin
open class Vehicle {
    open val maxSpeed: Int = 100
    open val name: String = "Vehicle"
}

class SportsCar : Vehicle() {
    override val maxSpeed: Int = 300
    override val name: String = "Sports Car"
}

class Bicycle : Vehicle() {
    override val maxSpeed: Int = 40
    override val name: String = "Bicycle"
}

val vehicles = listOf(SportsCar(), Bicycle(), Vehicle())
for (v in vehicles) {
    println("${v.name}: max ${v.maxSpeed} km/h")
}
// Sports Car: max 300 km/h
// Bicycle: max 40 km/h
// Vehicle: max 100 km/h
```

### Preventing Further Override

Mark an override with `final` to prevent further overriding:

```kotlin
open class A {
    open fun foo() = "A.foo"
}

open class B : A() {
    final override fun foo() = "B.foo"  // cannot be overridden further
}

class C : B() {
    // override fun foo() = "C.foo"  // ERROR: foo is final in B
}
```

---

## 6.6 Abstract Classes

An `abstract` class cannot be instantiated directly. It can have abstract members (no implementation) and concrete members (with implementation):

```kotlin
abstract class Shape {
    abstract val name: String
    abstract fun area(): Double
    abstract fun perimeter(): Double
    
    // Concrete method — shared by all shapes
    fun describe() = "$name — Area: ${"%.2f".format(area())}, Perimeter: ${"%.2f".format(perimeter())}"
}

class Circle(val radius: Double) : Shape() {
    override val name = "Circle"
    override fun area() = Math.PI * radius * radius
    override fun perimeter() = 2 * Math.PI * radius
}

class Rectangle(val width: Double, val height: Double) : Shape() {
    override val name = "Rectangle"
    override fun area() = width * height
    override fun perimeter() = 2 * (width + height)
}

val shapes: List<Shape> = listOf(Circle(5.0), Rectangle(4.0, 6.0))
shapes.forEach { println(it.describe()) }
// Circle — Area: 78.54, Perimeter: 31.42
// Rectangle — Area: 24.00, Perimeter: 20.00
```

---

## 6.7 Interfaces

Interfaces define a contract. Unlike abstract classes:
- A class can implement multiple interfaces
- Interfaces cannot hold state (no backing fields)
- But interfaces CAN have default implementations

```kotlin
interface Printable {
    fun print()                  // abstract
    fun prettyPrint() {          // default implementation
        println("=== Pretty Print ===")
        print()
        println("===================")
    }
}

interface Saveable {
    fun save(filename: String)
    val defaultExtension: String  // abstract property
}

class Document(val content: String) : Printable, Saveable {
    override fun print() = println(content)
    override fun save(filename: String) = println("Saving to $filename.$defaultExtension")
    override val defaultExtension = "txt"
}

val doc = Document("Hello, Kotlin!")
doc.print()             // Hello, Kotlin!
doc.prettyPrint()       // === Pretty Print === / Hello, Kotlin! / ===================
doc.save("output")      // Saving to output.txt
```

### Interface Default Methods

Interfaces can provide default implementations that can be optionally overridden:

```kotlin
interface Logger {
    fun log(message: String)  // must be implemented
    
    fun logError(message: String) {  // default — can be overridden
        log("ERROR: $message")
    }
    
    fun logInfo(message: String) {   // default
        log("INFO: $message")
    }
}

class ConsoleLogger : Logger {
    override fun log(message: String) = println("[${'$'}{System.currentTimeMillis()}] $message")
    // logError and logInfo use default implementations
}

class FileLogger(val filename: String) : Logger {
    override fun log(message: String) = appendToFile(filename, message)
    override fun logError(message: String) = log("CRITICAL ERROR: $message")  // custom
    
    private fun appendToFile(name: String, content: String) {
        println("Writing to $name: $content")
    }
}
```

### Interface vs Abstract Class

| Aspect | Interface | Abstract Class |
|--------|-----------|----------------|
| Multiple inheritance | Yes (implements multiple) | No (extends one) |
| State (fields) | No backing fields | Yes |
| Constructor | No | Yes |
| Default implementations | Yes | Yes |
| `open` by default | Yes | Yes |

**Use interfaces when** you're defining a contract or capability that many unrelated types can fulfill.  
**Use abstract classes when** you're defining a common base with shared state and behavior.

---

## 6.8 Data Classes

Data classes are one of Kotlin's most-loved features. They're designed to hold data, and the compiler automatically generates common methods.

```kotlin
data class Person(val name: String, val age: Int)
```

The compiler generates:
- `equals()` — structural equality based on all properties
- `hashCode()` — consistent with equals()
- `toString()` — readable representation
- `copy()` — create a copy with some properties changed
- `componentN()` functions — for destructuring

```kotlin
data class Person(val name: String, val age: Int)

val alice = Person("Alice", 30)
val bob = Person("Bob", 25)
val aliceCopy = Person("Alice", 30)

// toString
println(alice)               // Person(name=Alice, age=30)

// equals
println(alice == aliceCopy)  // true
println(alice == bob)        // false

// copy — create modified copy
val olderAlice = alice.copy(age = 31)
println(olderAlice)         // Person(name=Alice, age=31)

val aliceWithNewName = alice.copy(name = "Alicia")
println(aliceWithNewName)   // Person(name=Alicia, age=30)

// hashCode consistency
println(alice.hashCode() == aliceCopy.hashCode())  // true

// Destructuring
val (name, age) = alice
println("Name: $name, Age: $age")  // Name: Alice, Age: 30
```

### Data Class Rules

1. Primary constructor must have at least one parameter
2. All primary constructor parameters must be `val` or `var`
3. Data classes cannot be `abstract`, `open`, `sealed`, or `inner`
4. (Since Kotlin 1.1) Data classes can extend other classes

```kotlin
// Valid data class
data class Point(val x: Double, val y: Double)

// Valid — extending another class
abstract class Named(open val name: String)
data class NamedPoint(override val name: String, val x: Double, val y: Double) : Named(name)
```

### Properties Outside the Constructor

Properties declared in the body are NOT included in `equals`, `hashCode`, `toString`, or `copy`:

```kotlin
data class User(val id: Int, val name: String) {
    var lastLogin: Long = 0  // NOT part of data class equality/copy
}

val u1 = User(1, "Alice")
u1.lastLogin = 1000L

val u2 = User(1, "Alice")
u2.lastLogin = 9999L

println(u1 == u2)  // true — lastLogin is not considered
println(u1)        // User(id=1, name=Alice) — lastLogin not shown
```

---

## 6.9 Enum Classes

Enum classes represent a fixed set of constants. Each constant is a singleton instance of the enum class.

### Basic Enum

```kotlin
enum class Direction {
    NORTH, SOUTH, EAST, WEST
}

val dir = Direction.NORTH
println(dir)          // NORTH
println(dir.name)     // NORTH
println(dir.ordinal)  // 0

// All values
println(Direction.values().toList())
// [NORTH, SOUTH, EAST, WEST]

// From string
val parsed = Direction.valueOf("EAST")
println(parsed)  // EAST
```

### Enum with Properties

```kotlin
enum class Planet(val mass: Double, val radius: Double) {
    MERCURY(3.303e+23, 2.4397e6),
    VENUS(4.869e+24, 6.0518e6),
    EARTH(5.976e+24, 6.37814e6),
    MARS(6.421e+23, 3.3972e6);
    
    val surfaceGravity: Double
        get() = G * mass / (radius * radius)
    
    val surfaceWeight: (Double) -> Double
        get() = { mass -> mass * surfaceGravity }
    
    companion object {
        const val G = 6.67300E-11  // gravitational constant
    }
}

val earthWeight = 75.0
val mass = earthWeight / Planet.EARTH.surfaceGravity

for (p in Planet.values()) {
    println("Weight on ${p.name}: ${"%.2f".format(p.surfaceWeight(mass))}")
}
```

### Enum with Abstract Methods

```kotlin
enum class Operation(val symbol: String) {
    PLUS("+") {
        override fun apply(x: Double, y: Double) = x + y
    },
    MINUS("-") {
        override fun apply(x: Double, y: Double) = x - y
    },
    TIMES("*") {
        override fun apply(x: Double, y: Double) = x * y
    },
    DIVIDE("/") {
        override fun apply(x: Double, y: Double) = x / y
    };
    
    abstract fun apply(x: Double, y: Double): Double
    
    override fun toString() = symbol
}

val a = 10.0
val b = 3.0

for (op in Operation.values()) {
    println("$a $op $b = ${op.apply(a, b)}")
}
// 10.0 + 3.0 = 13.0
// 10.0 - 3.0 = 7.0
// 10.0 * 3.0 = 30.0
// 10.0 / 3.0 = 3.3333...
```

### Enum in when

```kotlin
enum class Color { RED, GREEN, BLUE }

fun mix(c1: Color, c2: Color): String = when {
    c1 == Color.RED && c2 == Color.YELLOW   -> "Orange"
    c1 == Color.YELLOW && c2 == Color.BLUE  -> "Green"
    c1 == Color.BLUE && c2 == Color.YELLOW  -> "Green"
    else -> "Unknown"
}

// Exhaustive when with enum — no else needed
fun describe(color: Color): String = when (color) {
    Color.RED   -> "Red — danger or stop"
    Color.GREEN -> "Green — safe or go"
    Color.BLUE  -> "Blue — calm or information"
}
```

---

## 6.10 Sealed Classes and Sealed Interfaces

**Sealed classes** restrict which classes can extend them. All subclasses must be in the **same package** (Kotlin 1.5+) or the **same file** (pre-1.5). This gives the compiler complete knowledge of all possible subclasses.

### Basic Sealed Class

```kotlin
sealed class Result<out T> {
    data class Success<T>(val value: T) : Result<T>()
    data class Failure(val error: Throwable) : Result<Nothing>()
    object Loading : Result<Nothing>()
}

fun processResult(result: Result<String>) = when (result) {
    is Result.Success -> println("Success: ${result.value}")
    is Result.Failure -> println("Error: ${result.error.message}")
    Result.Loading    -> println("Loading...")
    // No else needed — compiler knows all possible subclasses
}

processResult(Result.Success("Data loaded!"))
processResult(Result.Failure(RuntimeException("Network error")))
processResult(Result.Loading)
```

### Sealed Classes for Domain Modeling

```kotlin
sealed class PaymentStatus {
    object Pending : PaymentStatus()
    data class Processing(val transactionId: String) : PaymentStatus()
    data class Completed(val transactionId: String, val amount: Double) : PaymentStatus()
    data class Failed(val reason: String, val retryable: Boolean) : PaymentStatus()
    object Cancelled : PaymentStatus()
}

fun describeStatus(status: PaymentStatus): String = when (status) {
    PaymentStatus.Pending -> "Payment pending"
    is PaymentStatus.Processing -> "Processing (TX: ${status.transactionId})"
    is PaymentStatus.Completed -> "Paid \$${status.amount} (TX: ${status.transactionId})"
    is PaymentStatus.Failed -> if (status.retryable) "Failed, retrying..." else "Failed: ${status.reason}"
    PaymentStatus.Cancelled -> "Payment cancelled"
}

val status: PaymentStatus = PaymentStatus.Completed("TX-123", 99.99)
println(describeStatus(status))  // Paid $99.99 (TX: TX-123)
```

### Sealed Interface (Kotlin 1.5+)

```kotlin
sealed interface Notification {
    data class Email(val to: String, val subject: String) : Notification
    data class SMS(val to: String, val body: String) : Notification
    data class Push(val deviceId: String, val title: String, val body: String) : Notification
}

fun send(notification: Notification) = when (notification) {
    is Notification.Email -> println("Email to ${notification.to}: ${notification.subject}")
    is Notification.SMS -> println("SMS to ${notification.to}: ${notification.body}")
    is Notification.Push -> println("Push to ${notification.deviceId}: ${notification.title}")
}
```

### Sealed vs Enum

| Feature | Enum | Sealed Class |
|---------|------|-------------|
| Each constant | Same type | Different types possible |
| Each constant has state | Shared (all same class) | Each can have own properties |
| Use case | Fixed constants | Fixed hierarchy of types |
| Pattern matching | Yes | Yes |

---

## 6.11 Object Declarations and Companion Objects

### Object Declaration (Singleton)

The `object` keyword creates a singleton — a class with exactly one instance:

```kotlin
object DatabaseConfig {
    val host = "localhost"
    val port = 5432
    val databaseName = "mydb"
    
    fun connectionString() = "jdbc:postgresql://$host:$port/$databaseName"
}

println(DatabaseConfig.host)               // localhost
println(DatabaseConfig.connectionString()) // jdbc:postgresql://localhost:5432/mydb
```

No `getInstance()` — you access the object directly by its name.

### Companion Objects

A companion object is an `object` declaration inside a class. It replaces Java's `static` members:

```kotlin
class User(val name: String, val email: String) {
    
    companion object {
        // "Static" factory methods
        fun fromJson(json: String): User {
            // Parse JSON and create User
            val name = "Alice"  // simplified
            val email = "alice@example.com"
            return User(name, email)
        }
        
        // "Static" constant
        const val MAX_NAME_LENGTH = 50
    }
}

val user = User.fromJson("{...}")
println(User.MAX_NAME_LENGTH)  // 50
// Access companion object members directly on the class
```

### Named Companion Objects

```kotlin
class MyClass {
    companion object Factory {
        fun create(): MyClass = MyClass()
    }
}

val instance = MyClass.Factory.create()
val instance2 = MyClass.create()  // companion name is optional
```

### Object Expressions (Anonymous Objects)

Object expressions create anonymous objects, often used for one-off interface implementations:

```kotlin
interface Clickable {
    fun onClick()
    fun onLongClick() {
        println("Long click default")
    }
}

// Anonymous object implementing an interface
val button = object : Clickable {
    override fun onClick() = println("Clicked!")
    override fun onLongClick() = println("Long clicked!")
}

button.onClick()      // Clicked!
button.onLongClick()  // Long clicked!

// Anonymous object with no base type
val point = object {
    val x = 10
    val y = 20
}
println("(${point.x}, ${point.y})")  // (10, 20)
```

---

## Complete OOP Example: A Shape Hierarchy

```kotlin
// Abstract base
abstract class Shape(val color: String) {
    abstract fun area(): Double
    abstract fun perimeter(): Double
    
    override fun toString() = "${this::class.simpleName}($color) area=${"%.2f".format(area())}"
}

// Concrete classes
class Circle(color: String, val radius: Double) : Shape(color) {
    override fun area() = Math.PI * radius * radius
    override fun perimeter() = 2 * Math.PI * radius
}

class Rectangle(color: String, val width: Double, val height: Double) : Shape(color) {
    override fun area() = width * height
    override fun perimeter() = 2 * (width + height)
}

data class Triangle(val color: String, val a: Double, val b: Double, val c: Double) {
    fun area(): Double {
        val s = (a + b + c) / 2
        return Math.sqrt(s * (s - a) * (s - b) * (s - c))
    }
}

// Companion for factory
class ShapeFactory {
    companion object {
        fun unitCircle(color: String = "white") = Circle(color, 1.0)
        fun square(color: String = "white", side: Double) = Rectangle(color, side, side)
    }
}

fun main() {
    val shapes: List<Shape> = listOf(
        Circle("red", 5.0),
        Rectangle("blue", 4.0, 6.0),
        ShapeFactory.unitCircle("green"),
        ShapeFactory.square("yellow", side = 3.0)
    )
    
    shapes.forEach { println(it) }
    
    println("\nLargest area: ${shapes.maxByOrNull { it.area() }}")
    println("Total area: ${"%.2f".format(shapes.sumOf { it.area() })}")
}
```

---

## Summary

Kotlin's OOP system builds on Java's but improves it in key ways: classes are closed (final) by default, preventing unintended extension; the primary constructor is part of the class header; `val`/`var` in the constructor creates properties automatically; `data class` generates `equals`, `hashCode`, `toString`, `copy`, and destructuring; `sealed class` restricts inheritance and enables exhaustive `when` matching; `object` declarations create singletons; `companion object` replaces Java's static members.

---

## Key Takeaways

- Classes are final by default — use `open` to allow inheritance
- Primary constructor is in the class header; use `init` for initialization logic
- `val`/`var` in primary constructor automatically creates a property
- `data class` is for pure data holders — generates all boilerplate
- `sealed class` = closed hierarchy + exhaustive `when` = type-safe alternatives
- `enum class` is for fixed sets of constants, optionally with behavior
- `object` creates singletons; `companion object` is Kotlin's `static`
- Interfaces can have default implementations; classes can implement multiple interfaces

---

## Practice Questions

### Conceptual
1. Why are Kotlin classes closed (final) by default?
2. What is the difference between `object` declaration and `companion object`?
3. When would you use `sealed class` over `enum class`?
4. What methods does a `data class` automatically generate?
5. What is the difference between `abstract class` and `interface` in Kotlin?

### Code Exercises

**Exercise 1:** Model a `BankAccount` class with:
- `accountNumber: String` and `owner: String` (immutable)
- `balance: Double` with private setter
- `deposit(amount: Double)` and `withdraw(amount: Double)` methods
- Validation in `init` (balance >= 0)

**Exercise 2:** Create a sealed class `NetworkResponse<T>` with:
- `Success(data: T, statusCode: Int = 200)`
- `Error(message: String, statusCode: Int)`
- `Loading`
Write a function that handles each case with a `when` expression.

**Exercise 3:** Design an enum class `Season` with a property `monthRange: IntRange` (the months that fall in that season). Add a companion object method `fromMonth(month: Int): Season`.

**Exercise 4:** Create a `shape hierarchy` with an abstract class `Shape`, and concrete implementations `Circle`, `Rectangle`, and `Triangle`. Add a companion factory method on each concrete class. Use `when(shape)` for polymorphic area calculation.

**Exercise 5:** Write a `Logger` singleton using `object` that:
- Maintains a `MutableList<String>` of log entries
- Has `log(message: String)`, `clear()`, and `export(): List<String>` methods
- Ensures the log list is only accessible as read-only from outside

---

*Next: [Chapter 7 — Null Safety](07-null-safety.md)*
