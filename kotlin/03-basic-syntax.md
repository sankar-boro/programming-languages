# Chapter 3 — Variables, Types, and Operators

> *"A type system is a tractable syntactic method for proving the absence of certain program behaviors by classifying phrases according to the kinds of values they compute."*
> — Benjamin C. Pierce, Types and Programming Languages

---

## 3.1 val vs var

Kotlin has two keywords for declaring variables: `val` and `var`. This distinction is fundamental and matters a great deal in Kotlin.

### val — Immutable Reference

`val` declares a **read-only** (immutable) variable. Once assigned, it cannot be reassigned.

```kotlin
val name = "Alice"
val age = 30
val pi = 3.14159

// This would be a compile error:
// name = "Bob"  // Error: Val cannot be reassigned
```

`val` is similar to `final` in Java or `const` in JavaScript.

### var — Mutable Reference

`var` declares a **mutable** variable that can be reassigned:

```kotlin
var score = 0
score = 10
score += 5
println(score)  // 15
```

### Which to Use?

**Prefer `val` over `var`** — this is a core Kotlin idiom.

Using `val` by default:
- Makes code easier to reason about (no unexpected mutations)
- Communicates intent clearly
- Works better with functional programming patterns
- Is safer in concurrent code

Only use `var` when you genuinely need to reassign a variable.

```kotlin
// Good — prefer val
val users = listOf("Alice", "Bob", "Charlie")
val count = users.size

// Acceptable — mutation is genuinely needed
var sum = 0
for (n in 1..10) {
    sum += n
}
println(sum)  // 55

// Bad pattern — using var when val would do
var greeting = "Hello"
// greeting never changes after this... should be val
```

### val vs Immutability

Important: `val` means the **reference is immutable**, not the object it points to.

```kotlin
val list = mutableListOf(1, 2, 3)
list.add(4)       // OK — the list object is mutable, we just can't reassign `list`
list.add(5)       // OK
// list = mutableListOf()  // Error — can't reassign the val reference

println(list)  // [1, 2, 3, 4, 5]
```

This is an important distinction. For true immutability, use immutable collection types (covered in Chapter 9).

### Declaring Type Explicitly

While type inference usually handles this, you can declare the type explicitly:

```kotlin
val name: String = "Alice"
val age: Int = 30
val score: Double = 9.5
val isActive: Boolean = true
```

---

## 3.2 Basic Data Types

Kotlin is a statically typed language — every value has a type known at compile time. Unlike Java, Kotlin does **not distinguish between primitive types and wrapper types** at the language level (though the compiler optimizes to primitives on the JVM).

### Integer Types

| Type | Size | Range |
|------|------|-------|
| `Byte` | 8 bits | -128 to 127 |
| `Short` | 16 bits | -32,768 to 32,767 |
| `Int` | 32 bits | -2,147,483,648 to 2,147,483,647 |
| `Long` | 64 bits | -9,223,372,036,854,775,808 to 9,223,372,036,854,775,807 |

```kotlin
val byteVal: Byte = 100
val shortVal: Short = 30000
val intVal: Int = 2_000_000    // underscores for readability
val longVal: Long = 9_000_000_000L  // L suffix for Long literals

// Default integer literal type is Int
val x = 42        // Int
val y = 42L       // Long (L suffix)

// Underscores in numeric literals (no effect on value)
val million = 1_000_000
val creditCard = 1234_5678_9012_3456L
val hex = 0xFF_EC_DE_5E
val binary = 0b11010010_01101001
```

### Floating-Point Types

| Type | Size | Precision |
|------|------|-----------|
| `Float` | 32 bits | ~6-7 decimal digits |
| `Double` | 64 bits | ~15-16 decimal digits |

```kotlin
val doubleVal = 3.14        // Double (default)
val floatVal = 3.14f        // Float (f suffix)
val scientific = 1.23e10    // Double, scientific notation

// Double is preferred — more precise
val pi = 3.14159265358979  // Double
val approxPi = 3.14f       // Float
```

### Characters

`Char` represents a single Unicode character:

```kotlin
val letter: Char = 'A'
val digit: Char = '5'
val special: Char = '\n'   // newline escape
val unicode: Char = 'A'  // 'A' in Unicode

// Characters are NOT numbers in Kotlin (unlike Java)
// You can't do: val x: Int = 'A'  — compile error
val code: Int = 'A'.code   // 65 — explicit conversion
val char: Char = 65.toChar()  // 'A'
```

### Booleans

```kotlin
val isKotlinAwesome: Boolean = true
val isBoringLanguage: Boolean = false

// Boolean operators
val and = true && false   // false
val or = true || false    // true
val not = !true           // false
```

### Strings

Strings are sequences of characters. In Kotlin, `String` is immutable:

```kotlin
val hello = "Hello, World!"
val multiLine = """
    Line 1
    Line 2
    Line 3
""".trimIndent()

println(multiLine)
// Line 1
// Line 2
// Line 3
```

#### String Templates

One of Kotlin's most convenient features:

```kotlin
val name = "Alice"
val age = 30

// Simple variable reference
println("Name: $name")

// Expressions in templates
println("In 5 years: ${age + 5}")
println("Uppercase: ${name.uppercase()}")

// Nested templates (rare but possible)
val items = listOf("a", "b", "c")
println("Items: ${items.joinToString(", ")}")

// Escaping the $ sign
println("Price: \$99.99")  // Price: $99.99
```

#### Raw Strings (Triple-Quoted)

```kotlin
val json = """
    {
        "name": "Alice",
        "age": 30
    }
""".trimIndent()

println(json)
// {
//     "name": "Alice",
//     "age": 30
// }

// String templates work inside raw strings too
val name = "Bob"
val greeting = """Hello, $name!
Welcome to Kotlin."""
println(greeting)
```

#### String Operations

```kotlin
val str = "Hello, Kotlin!"

println(str.length)           // 14
println(str.uppercase())      // HELLO, KOTLIN!
println(str.lowercase())      // hello, kotlin!
println(str.contains("Kotlin"))  // true
println(str.startsWith("Hello")) // true
println(str.endsWith("!"))       // true
println(str.replace("Kotlin", "World"))  // Hello, World!
println(str.substring(7, 13))   // Kotlin
println(str.split(", "))        // [Hello, Kotlin!]
println(str.trim())             // removes whitespace from both ends
println(str[0])                 // H — index access
println(str.isEmpty())          // false
println("".isBlank())           // true
println("  ".isBlank())         // true — blank checks for whitespace-only
```

### The Any, Unit, and Nothing Types

```kotlin
// Any — the root of the Kotlin type hierarchy (like Object in Java)
val anything: Any = "I can be anything"
val alsoAnything: Any = 42
val booleanToo: Any = true

// Unit — equivalent to void in Java
// Functions that don't return a meaningful value return Unit
fun printMessage(msg: String): Unit {  // Unit is usually omitted
    println(msg)
}
// Same as:
fun printMessage(msg: String) {
    println(msg)
}

// Nothing — a type with no values; represents "never returns normally"
fun fail(message: String): Nothing {
    throw IllegalStateException(message)
}

fun infiniteLoop(): Nothing {
    while (true) { }
}
```

`Nothing` is useful because the compiler knows that after a `Nothing`-returning function, no code is reachable.

---

## 3.3 Type Inference

Kotlin's type inference is powerful — the compiler deduces types from context so you don't have to repeat yourself.

```kotlin
val name = "Alice"      // String
val age = 30            // Int
val height = 5.9        // Double
val active = true       // Boolean
val score = 100L        // Long

// The compiler infers the return type of functions too:
fun double(n: Int) = n * 2  // return type inferred as Int

// Complex inference
val numbers = listOf(1, 2, 3)  // List<Int>
val mapped = numbers.map { it * 2 }  // List<Int>
val mixed = listOf(1, "two", 3.0)    // List<Any>
```

### When to Declare Types Explicitly

Type inference handles most cases, but explicit types are valuable when:

1. The type is not obvious from the expression:
```kotlin
val result: Double = computeValue()  // clearer to read
```

2. You want a wider type than what would be inferred:
```kotlin
val number: Number = 42  // inferred as Int, but you want Number
```

3. At API boundaries (public function signatures):
```kotlin
// Good practice: always declare return types for public functions
fun calculateTax(income: Double): Double {
    return income * 0.2
}
```

4. When initializing without a value:
```kotlin
// Must declare type when not initializing
val count: Int  // declaration without initialization
// ... later in a block:
count = 5  // must be initialized before use
```

### Type Checking and Casting

```kotlin
val obj: Any = "Hello"

// is — type check (like instanceof in Java)
if (obj is String) {
    // Smart cast: obj is automatically treated as String here
    println(obj.length)  // No explicit cast needed!
}

// !is — negation
if (obj !is Int) {
    println("Not an integer")
}

// as — explicit (unsafe) cast
val str = obj as String   // throws ClassCastException if obj is not String

// as? — safe cast (returns null if cast fails, not exception)
val str2 = obj as? String  // String? (nullable)
val num = obj as? Int      // null, since obj is String
```

---

## 3.4 Type Conversions

Unlike Java, Kotlin does **NOT** do implicit numeric conversions. Every conversion must be explicit:

```kotlin
val intValue: Int = 100
// val longValue: Long = intValue  // COMPILE ERROR in Kotlin!

// Explicit conversions:
val longValue: Long = intValue.toLong()
val doubleValue: Double = intValue.toDouble()
val byteValue: Byte = intValue.toByte()
val floatValue: Float = intValue.toFloat()
val stringValue: String = intValue.toString()

// Converting String to number
val parsed: Int = "42".toInt()
val parsedOrNull: Int? = "not a number".toIntOrNull()  // returns null instead of exception
println(parsedOrNull)  // null

// Parsing with radix
val hex = "FF".toInt(16)  // 255
val binary = "1010".toInt(2)  // 10
```

---

## 3.5 Operators

### Arithmetic Operators

```kotlin
val a = 10
val b = 3

println(a + b)   // 13  — addition
println(a - b)   // 7   — subtraction
println(a * b)   // 30  — multiplication
println(a / b)   // 3   — integer division (truncates)
println(a % b)   // 1   — remainder (modulo)

// Float division
println(10.0 / 3.0)  // 3.3333333333333335
println(10 / 3.0)    // 3.3333333333333335 (Int auto-widened to Double)

// Augmented assignment
var x = 10
x += 5   // x = 15
x -= 3   // x = 12
x *= 2   // x = 24
x /= 4   // x = 6
x %= 4   // x = 2
```

### Comparison Operators

```kotlin
val x = 5
val y = 10

println(x == y)   // false — structural equality
println(x != y)   // true
println(x < y)    // true
println(x > y)    // false
println(x <= y)   // true
println(x >= y)   // false

// Referential equality (same object in memory)
val s1 = "Hello"
val s2 = "Hello"
val s3 = s1

println(s1 == s2)    // true  — structural equality (same content)
println(s1 === s2)   // true  — JVM string interning makes this true for literals
println(s1 === s3)   // true  — same reference
```

**Important:** In Kotlin, `==` calls the `equals()` method (structural equality). `===` checks reference equality. This is different from Java, where `==` on objects checks reference equality.

### Logical Operators

```kotlin
val a = true
val b = false

println(a && b)   // false — AND (short-circuits)
println(a || b)   // true  — OR (short-circuits)
println(!a)       // false — NOT
```

Short-circuit evaluation: in `a && b`, if `a` is false, `b` is never evaluated. In `a || b`, if `a` is true, `b` is never evaluated.

```kotlin
fun riskyOperation(): Boolean {
    println("riskyOperation called")
    return true
}

val result = false && riskyOperation()
// "riskyOperation called" is NOT printed — short-circuited
println(result)  // false
```

### Bitwise Operators

Kotlin uses named functions instead of symbols for bitwise operations:

```kotlin
val a = 0b1010  // 10 in decimal
val b = 0b0110  // 6 in decimal

println(a and b)    // 0b0010 = 2  — bitwise AND
println(a or b)     // 0b1110 = 14 — bitwise OR
println(a xor b)    // 0b1100 = 12 — bitwise XOR
println(a.inv())    // bitwise inversion
println(a shl 1)    // 0b10100 = 20 — shift left
println(a shr 1)    // 0b0101 = 5  — shift right (signed)
println(a ushr 1)   // unsigned shift right
```

### Range Operator

```kotlin
val range = 1..10       // IntRange from 1 to 10 (inclusive)
val charRange = 'a'..'z'  // CharRange

println(5 in range)   // true
println(11 in range)  // false
println('c' in charRange)  // true

// Used in for loops
for (i in 1..5) print("$i ")  // 1 2 3 4 5
println()
for (c in 'a'..'e') print("$c ")  // a b c d e
```

### String Plus Operator

```kotlin
val greeting = "Hello" + ", " + "World!"  // concatenation
println(greeting)  // Hello, World!

// Preferred: use string templates instead
val name = "Alice"
val better = "Hello, $name!"  // more readable and efficient
```

### Operator Overloading

Kotlin allows you to define operators for your own classes using the `operator` keyword:

```kotlin
data class Vector(val x: Double, val y: Double) {
    operator fun plus(other: Vector) = Vector(x + other.x, y + other.y)
    operator fun times(scalar: Double) = Vector(x * scalar, y * scalar)
    override fun toString() = "Vector($x, $y)"
}

val v1 = Vector(1.0, 2.0)
val v2 = Vector(3.0, 4.0)
val sum = v1 + v2   // calls v1.plus(v2)
val scaled = v1 * 2.0  // calls v1.times(2.0)

println(sum)    // Vector(4.0, 6.0)
println(scaled) // Vector(2.0, 4.0)
```

### Increment and Decrement

```kotlin
var x = 5
println(x++)   // 5 — post-increment: returns then increments
println(x)     // 6
println(++x)   // 7 — pre-increment: increments then returns
println(x--)   // 7 — post-decrement
println(x)     // 6
println(--x)   // 5 — pre-decrement
```

---

## 3.6 Input and Output Basics

### Output

```kotlin
println("Hello, World!")      // print with newline
print("Hello, ")              // print without newline
print("World!")
println()                     // just a newline

// Formatted output
val pi = 3.14159
println("Pi is approximately %.2f".format(pi))  // Pi is approximately 3.14

// Using String.format (Java-style)
val formatted = String.format("Name: %-10s Age: %d", "Alice", 30)
println(formatted)  // Name: Alice      Age: 30

// System.out directly (Java interop)
System.out.println("Using Java's System.out")
```

### Input

Reading from the command line:

```kotlin
// Read a single line as String
val line = readLine()           // returns String? (null at end of input)

// Read with a prompt (no automatic newline in prompt)
print("Enter your name: ")
val name = readLine() ?: "Anonymous"  // Elvis: default if null
println("Hello, $name!")

// Read an integer
print("Enter a number: ")
val number = readLine()?.toIntOrNull()
if (number != null) {
    println("Double: ${number * 2}")
} else {
    println("Invalid number")
}
```

### Reading Multiple Values

```kotlin
fun main() {
    print("Enter two numbers separated by space: ")
    val input = readLine() ?: ""
    val parts = input.split(" ")
    
    if (parts.size == 2) {
        val a = parts[0].toIntOrNull()
        val b = parts[1].toIntOrNull()
        
        if (a != null && b != null) {
            println("Sum: ${a + b}")
        } else {
            println("Invalid input")
        }
    }
}
```

---

## Putting It All Together

Here's a complete program demonstrating the concepts from this chapter:

```kotlin
fun main() {
    // Variables
    val firstName = "Alice"    // String — immutable
    val lastName = "Smith"     // String — immutable
    var age = 28               // Int — mutable
    val height = 165.5         // Double
    val isStudent = true       // Boolean
    
    // String templates
    println("Name: $firstName $lastName")
    println("Age: $age")
    println("Height: ${height}cm")
    println("Student: $isStudent")
    
    // Age changes (hence var)
    age++
    println("Next year: $age")
    
    // Type conversions
    val ageAsDouble = age.toDouble()
    val ageAsString = age.toString()
    println("Age as Double: $ageAsDouble")
    println("Age as String: '$ageAsString'")
    
    // Operators
    val bmi = 70.0 / (height / 100.0).let { it * it }
    println("BMI: ${"%.1f".format(bmi)}")
    
    // Ranges
    val passingGrades = 60..100
    val examScore = 85
    println("Passed: ${examScore in passingGrades}")
    
    // Type checking
    val value: Any = "I am a String"
    if (value is String) {
        println("String length: ${value.length}")  // smart cast
    }
}
```

---

## Summary

Kotlin has two variable declaration keywords: `val` (immutable reference) and `var` (mutable). The type system includes numeric types (`Int`, `Long`, `Double`, `Float`, `Byte`, `Short`), `Char`, `Boolean`, and `String`. Type inference allows the compiler to deduce types in most situations, reducing boilerplate. Kotlin requires **explicit type conversions** — there are no implicit narrowing or widening conversions. String templates (`$var` and `${expression}`) make string formatting clean and readable. Operators work as expected, with the important distinction that `==` does structural equality (calls `equals()`) and `===` does reference equality.

---

## Key Takeaways

- `val` is preferred over `var` — use immutability by default
- `val` means the reference is immutable, not the object itself
- All types are objects in Kotlin — no primitives at the language level
- Type inference reduces boilerplate without sacrificing type safety
- Kotlin does NOT do implicit numeric conversions — use `.toInt()`, `.toLong()`, etc.
- `==` calls `equals()` in Kotlin; `===` checks reference identity
- String templates are more readable and efficient than concatenation
- Raw strings (triple-quoted) avoid escaping complexity

---

## Practice Questions

### Conceptual
1. What is the difference between `val` and `var`? When should you use each?
2. If `val list = mutableListOf(1, 2, 3)`, can you call `list.add(4)`? Why or why not?
3. Why does Kotlin require explicit type conversions instead of implicit ones?
4. What is the difference between `==` and `===` in Kotlin?
5. What is the type of `42` in Kotlin? What about `42L`? What about `42.0`?

### Code Exercises

**Exercise 1:** Without running the code, predict the output:
```kotlin
var x = 10
val y = 3
println(x / y)
println(x % y)
println(x.toDouble() / y)
x = x * y
println(x)
```

**Exercise 2:** Write a program that reads a temperature in Celsius from the user and prints it in Fahrenheit. Formula: F = (C × 9/5) + 32. Handle invalid input gracefully.

**Exercise 3:** Given:
```kotlin
val obj: Any = 42
```
Write code that:
- Checks if `obj` is an `Int`
- If so, prints its square
- Otherwise, prints "Not an integer"

**Exercise 4:** Create a `String` variable with your full name. Using string operations (not templates), print:
- The name in uppercase
- The length of the name
- The first character
- Whether it contains a space

**Exercise 5:** Explain the difference between these two declarations and when you'd use each:
```kotlin
val number: Number = 42
val number = 42
```

---

*Next: [Chapter 4 — Control Flow](04-control-flow.md)*
