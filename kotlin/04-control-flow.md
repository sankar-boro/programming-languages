# Chapter 4 — Control Flow

> *"A language that doesn't affect the way you think about programming is not worth knowing."*
> — Alan Perlis

---

## 4.1 if as an Expression

In Kotlin, `if` is not a statement — it is an **expression**. This means it returns a value, and you can use it wherever a value is expected.

### Basic if-else

```kotlin
val age = 20

if (age >= 18) {
    println("Adult")
} else {
    println("Minor")
}
```

### if as an Expression

```kotlin
val age = 20
val status = if (age >= 18) "Adult" else "Minor"
println(status)  // Adult

// Direct use in expressions
val max = if (a > b) a else b

// Multi-line if expression
val description = if (age < 13) {
    "child"
} else if (age < 18) {
    "teenager"
} else if (age < 65) {
    "adult"
} else {
    "senior"
}
println(description)
```

When `if` is used as an expression with multiple lines in each branch, the **last expression** in each block is the value of that branch:

```kotlin
val a = 5
val b = 10

val max = if (a > b) {
    println("a is larger")
    a   // this is the value of this branch
} else {
    println("b is larger or equal")
    b   // this is the value of this branch
}

println("Max is: $max")
// b is larger or equal
// Max is: 10
```

### The Ternary Operator Replacement

Java has a ternary operator (`condition ? a : b`). Kotlin replaces it with inline `if`:

```java
// Java
String label = score >= 60 ? "Pass" : "Fail";
```

```kotlin
// Kotlin
val label = if (score >= 60) "Pass" else "Fail"
```

The Kotlin version is slightly more verbose but also more readable.

### if with No else — Unit Returns

When `if` is used as a statement (result discarded), the `else` branch is optional:

```kotlin
val n = -5

if (n < 0) {
    println("Negative number")
}
// OK — no else needed when used as statement

// When used as expression with both branches needed:
val abs = if (n < 0) -n else n  // else IS required here
```

---

## 4.2 when Expressions (Deep Dive)

`when` is Kotlin's powerful replacement for Java's `switch` statement. But where Java's `switch` is limited, Kotlin's `when` is rich and expressive.

### Basic when

```kotlin
val day = 3

when (day) {
    1 -> println("Monday")
    2 -> println("Tuesday")
    3 -> println("Wednesday")
    4 -> println("Thursday")
    5 -> println("Friday")
    6 -> println("Saturday")
    7 -> println("Sunday")
    else -> println("Invalid day")
}
// Wednesday
```

### when as an Expression

Like `if`, `when` can return a value:

```kotlin
val day = 6
val dayName = when (day) {
    1 -> "Monday"
    2 -> "Tuesday"
    3 -> "Wednesday"
    4 -> "Thursday"
    5 -> "Friday"
    6 -> "Saturday"
    7 -> "Sunday"
    else -> "Invalid day"
}
println(dayName)  // Saturday

// When used as expression, else is required if not all cases are covered
```

### Combining Conditions

Multiple values can be combined with a comma:

```kotlin
val day = 6
val dayType = when (day) {
    1, 2, 3, 4, 5 -> "Weekday"
    6, 7 -> "Weekend"
    else -> "Invalid"
}
println(dayType)  // Weekend
```

### Checking Ranges

```kotlin
val score = 85

val grade = when (score) {
    in 90..100 -> "A"
    in 80..89  -> "B"
    in 70..79  -> "C"
    in 60..69  -> "D"
    in 0..59   -> "F"
    else -> "Invalid score"
}
println(grade)  // B
```

### Checking Types (Type Smart Cast)

```kotlin
fun describe(obj: Any): String = when (obj) {
    is Int    -> "Integer: $obj"
    is String -> "String of length ${obj.length}"  // obj is smart-cast to String here
    is Boolean -> "Boolean: $obj"
    is List<*> -> "List with ${obj.size} elements"  // obj is smart-cast to List here
    else       -> "Unknown type"
}

println(describe(42))            // Integer: 42
println(describe("Hello"))       // String of length 5
println(describe(true))          // Boolean: true
println(describe(listOf(1,2,3))) // List with 3 elements
println(describe(3.14))          // Unknown type
```

### when Without an Argument

When `when` has no argument, it acts like an `if-else` chain with arbitrary conditions:

```kotlin
val temperature = 25

when {
    temperature < 0   -> println("Freezing")
    temperature < 10  -> println("Very cold")
    temperature < 20  -> println("Cold")
    temperature < 25  -> println("Comfortable")
    temperature < 35  -> println("Warm")
    else              -> println("Hot")
}
// Warm
```

This is very powerful — each branch can have any Boolean expression:

```kotlin
val name = "Alice"
val age = 25

when {
    name.startsWith("A") && age < 30 -> println("Young Alice!")
    name.length > 5                  -> println("Long name: $name")
    age > 50                         -> println("Experienced person")
    else                             -> println("Hello, $name")
}
// Young Alice!
```

### when with Multi-Line Blocks

```kotlin
val x = 5

val result = when (x) {
    1 -> {
        val computed = x * x
        "x is one, computed: $computed"
    }
    in 2..10 -> {
        val doubled = x * 2
        "x is between 2 and 10, doubled: $doubled"
    }
    else -> "x is out of range"
}
println(result)  // x is between 2 and 10, doubled: 10
```

### when with Sealed Classes (Preview)

`when` is especially powerful with sealed classes — the compiler can verify exhaustiveness:

```kotlin
sealed class Shape {
    data class Circle(val radius: Double) : Shape()
    data class Rectangle(val width: Double, val height: Double) : Shape()
    data class Triangle(val base: Double, val height: Double) : Shape()
}

fun area(shape: Shape): Double = when (shape) {
    is Shape.Circle    -> Math.PI * shape.radius * shape.radius
    is Shape.Rectangle -> shape.width * shape.height
    is Shape.Triangle  -> 0.5 * shape.base * shape.height
    // No 'else' needed! The compiler knows all subclasses of sealed Shape
}

println(area(Shape.Circle(5.0)))                 // 78.53...
println(area(Shape.Rectangle(4.0, 6.0)))         // 24.0
println(area(Shape.Triangle(3.0, 8.0)))          // 12.0
```

We'll cover sealed classes in depth in Chapter 6.

### Exhaustiveness in when

When `when` is used as an **expression** (its value is used), the `else` branch is required UNLESS:
1. You're matching on a `Boolean` (covering both true and false)
2. You're matching on an `enum` (covering all values)
3. You're matching on a `sealed class` (covering all subtypes)

```kotlin
enum class Direction { NORTH, SOUTH, EAST, WEST }

fun describe(d: Direction): String = when (d) {
    Direction.NORTH -> "Going north"
    Direction.SOUTH -> "Going south"
    Direction.EAST  -> "Going east"
    Direction.WEST  -> "Going west"
    // No else needed — all enum values are covered
}
```

---

## 4.3 for Loops and Ranges

### Basic for Loop

Kotlin's `for` loop iterates over anything that provides an iterator:

```kotlin
val fruits = listOf("Apple", "Banana", "Cherry")

for (fruit in fruits) {
    println(fruit)
}
// Apple
// Banana
// Cherry
```

### Iterating with Index

```kotlin
val fruits = listOf("Apple", "Banana", "Cherry")

// Using withIndex()
for ((index, fruit) in fruits.withIndex()) {
    println("$index: $fruit")
}
// 0: Apple
// 1: Banana
// 2: Cherry

// Using indices property
for (i in fruits.indices) {
    println("$i: ${fruits[i]}")
}
```

### Iterating a Range

```kotlin
// Inclusive range: 1 to 10
for (i in 1..10) {
    print("$i ")
}
// 1 2 3 4 5 6 7 8 9 10

// Exclusive upper bound: 1 until 10
for (i in 1 until 10) {
    print("$i ")
}
// 1 2 3 4 5 6 7 8 9

// Step: every other number
for (i in 0..20 step 5) {
    print("$i ")
}
// 0 5 10 15 20

// Counting down
for (i in 10 downTo 1) {
    print("$i ")
}
// 10 9 8 7 6 5 4 3 2 1

// Counting down with step
for (i in 20 downTo 0 step 4) {
    print("$i ")
}
// 20 16 12 8 4 0
```

### Iterating Strings and Characters

```kotlin
val word = "Kotlin"

for (ch in word) {
    print("$ch-")
}
// K-o-t-l-i-n-

for ((index, ch) in word.withIndex()) {
    println("$index: $ch")
}
```

### Iterating Maps

```kotlin
val capitals = mapOf(
    "France" to "Paris",
    "Germany" to "Berlin",
    "Japan" to "Tokyo"
)

for ((country, capital) in capitals) {
    println("$country: $capital")
}
// France: Paris
// Germany: Berlin
// Japan: Tokyo
```

---

## 4.4 while and do-while

### while Loop

```kotlin
var count = 0

while (count < 5) {
    print("$count ")
    count++
}
// 0 1 2 3 4
```

### do-while Loop

The body executes at least once before the condition is checked:

```kotlin
var count = 10

do {
    println("Count: $count")
    count--
} while (count > 0)
// Prints Count: 10, Count: 9, ... Count: 1

// Even if condition is initially false:
var x = 100
do {
    println("This runs once: $x")  // Prints once
} while (x < 0)  // Condition false, but body already ran
```

### Practical while Example

```kotlin
fun main() {
    var attempts = 0
    val maxAttempts = 3
    var success = false
    
    while (attempts < maxAttempts && !success) {
        print("Enter password: ")
        val input = readLine()
        
        if (input == "secret123") {
            success = true
            println("Access granted!")
        } else {
            attempts++
            println("Wrong password. ${maxAttempts - attempts} attempts remaining.")
        }
    }
    
    if (!success) {
        println("Account locked.")
    }
}
```

---

## 4.5 Ranges and Progressions

Ranges are a powerful feature in Kotlin. They represent a sequence of values with a defined start, end, and optional step.

### Range Types

```kotlin
// IntRange
val intRange = 1..10
println(intRange.first)  // 1
println(intRange.last)   // 10
println(intRange.step)   // 1

// CharRange
val charRange = 'a'..'z'

// LongRange
val longRange = 1L..1_000_000L

// Closed range (includes last value)
val closed = 1..10  // 1, 2, 3, ..., 10

// Half-open range (excludes last value)
val halfOpen = 1 until 10  // 1, 2, 3, ..., 9
```

### Range Operations

```kotlin
val range = 1..100

println(50 in range)        // true
println(101 in range)       // false
println(0 !in range)        // true

println(range.count())      // 100
println(range.sum())        // 5050
println(range.average())    // 50.5
println(range.min())        // 1
println(range.max())        // 100

// Convert to list
val list = (1..5).toList()
println(list)  // [1, 2, 3, 4, 5]
```

### Progression

A progression is a range with a step:

```kotlin
// downTo creates a descending progression
val descending = 10 downTo 1  // 10, 9, 8, ..., 1

// step modifies the interval
val evens = 2..20 step 2    // 2, 4, 6, ..., 20
val odds = 1..20 step 2     // 1, 3, 5, ..., 19
val bigStep = 0..100 step 25  // 0, 25, 50, 75, 100

for (n in evens) print("$n ")  // 2 4 6 8 10 12 14 16 18 20
println()
for (n in bigStep) print("$n ")  // 0 25 50 75 100
```

### Custom Comparable Ranges

Any `Comparable` type supports ranges:

```kotlin
val nameRange = "Alice".."George"  // Alphabetical range
println("Bob" in nameRange)        // true
println("Henry" in nameRange)      // false

// Date ranges (with a comparable type)
val letters = 'a'..'f'
for (c in letters) print("$c ")   // a b c d e f
```

---

## 4.6 break and continue

### break

Exits the nearest enclosing loop:

```kotlin
for (i in 1..10) {
    if (i == 5) break
    print("$i ")
}
// 1 2 3 4

var x = 0
while (true) {
    if (x >= 5) break
    println(x++)
}
// 0 1 2 3 4
```

### continue

Skips the current iteration and proceeds to the next:

```kotlin
for (i in 1..10) {
    if (i % 2 == 0) continue  // skip even numbers
    print("$i ")
}
// 1 3 5 7 9
```

### Labels — Breaking Outer Loops

Kotlin supports **labeled break** and **labeled continue** to target specific loops in nested structures:

```kotlin
outer@ for (i in 1..5) {
    for (j in 1..5) {
        if (j == 3) continue@outer  // continue the outer loop
        print("($i,$j) ")
    }
}
println()
// (1,1) (1,2) (2,1) (2,2) (3,1) (3,2) (4,1) (4,2) (5,1) (5,2)

outer@ for (i in 1..5) {
    for (j in 1..5) {
        if (i == 3 && j == 2) break@outer  // break the outer loop entirely
        print("($i,$j) ")
    }
}
println()
// (1,1) (1,2) (1,3) (1,4) (1,5) (2,1) (2,2) (2,3) (2,4) (2,5) (3,1)
```

---

## Comprehensive Example: Putting It All Together

```kotlin
fun analyzeNumbers(numbers: List<Int>) {
    println("=== Number Analysis ===")
    println("Count: ${numbers.size}")
    
    // Basic stats using ranges and loops
    var sum = 0
    var min = numbers.first()
    var max = numbers.first()
    
    for (n in numbers) {
        sum += n
        if (n < min) min = n
        if (n > max) max = n
    }
    
    val avg = sum.toDouble() / numbers.size
    println("Sum: $sum, Min: $min, Max: $max, Avg: $avg")
    
    // Categorize each number
    for (n in numbers) {
        val category = when {
            n < 0          -> "negative"
            n == 0         -> "zero"
            n in 1..9      -> "single digit"
            n in 10..99    -> "double digit"
            else           -> "large"
        }
        println("  $n -> $category")
    }
    
    // Find first even number using labeled break
    var firstEven: Int? = null
    search@ for (n in numbers) {
        if (n % 2 == 0) {
            firstEven = n
            break@search
        }
    }
    
    println("First even: ${firstEven ?: "none found"}")
}

fun main() {
    analyzeNumbers(listOf(3, -2, 15, 0, 8, 42, 7, -11, 100))
}
```

Output:
```
=== Number Analysis ===
Count: 9
Sum: 162, Min: -11, Max: 100, Avg: 18.0
  3 -> single digit
  -2 -> negative
  15 -> double digit
  0 -> zero
  8 -> single digit
  42 -> double digit
  7 -> single digit
  -11 -> negative
  100 -> large
First even: -2
```

---

## Summary

Kotlin's control flow is more expressive than Java's:

- `if` is an **expression** that returns a value — replacing the ternary operator
- `when` is a powerful **expression** that replaces `switch` with: arbitrary conditions, type checking, range checking, and pattern matching
- `for` loops iterate over any `Iterable`, with built-in support for ranges, progressions, and indexed access
- `while` and `do-while` work as in other languages
- Ranges (`1..10`, `1 until 10`, `10 downTo 1 step 2`) are a first-class feature
- Labeled `break@label` and `continue@label` allow targeting specific nested loops

---

## Key Takeaways

- `if` returns a value — prefer inline `if` over the Java ternary `?:`
- `when` is exhaustive when matching sealed classes and enums (no `else` needed)
- `when` without an argument acts as a chain of `if-else` conditions
- Ranges use `..` for inclusive, `until` for exclusive upper bound, `downTo` for descending
- Smart casts apply inside `is` checks in `when` branches
- Labels (`outer@`) allow `break` and `continue` to target specific loops

---

## Practice Questions

### Conceptual
1. How is Kotlin's `if` different from Java's `if`?
2. When is the `else` branch required in a `when` expression?
3. What is the difference between `1..10` and `1 until 10`?
4. How does `when` with no argument differ from `when` with an argument?
5. When would you use `break@label` instead of a simple `break`?

### Code Exercises

**Exercise 1:** Write a function `fizzBuzz(n: Int)` that prints numbers from 1 to n, replacing:
- Multiples of 3 with "Fizz"
- Multiples of 5 with "Buzz"  
- Multiples of both with "FizzBuzz"
Use `when` without an argument.

**Exercise 2:** Given an `Int` value, use a `when` expression to return its category as a `String`:
- Negative: "negative"
- 0: "zero"
- 1-9: "single digit"
- 10-99: "double digit"
- 100+: "triple digit or more"

**Exercise 3:** Write a program that:
- Generates numbers from 1 to 100
- Skips multiples of 7 (use `continue`)
- Stops if the sum exceeds 500 (use `break`)
- Prints the final sum and how many numbers were added

**Exercise 4:** Write a nested loop that prints a multiplication table (1-5 × 1-5) using labeled `continue` to skip cells where both factors are odd.

**Exercise 5:** Rewrite this Java-style code in idiomatic Kotlin:
```java
int x = getValue();
String result;
if (x > 0) {
    result = "positive";
} else if (x < 0) {
    result = "negative";
} else {
    result = "zero";
}
```

---

*Next: [Chapter 5 — Functions](05-functions.md)*
