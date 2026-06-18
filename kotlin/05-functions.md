# Chapter 5 — Functions

> *"Functions are the building blocks of programs. In Kotlin, they are first-class citizens."*

---

## 5.1 Defining Functions

Functions in Kotlin are declared with the `fun` keyword. They can be top-level (outside any class), inside a class, or even inside another function.

### Basic Function Syntax

```kotlin
fun functionName(param1: Type1, param2: Type2): ReturnType {
    // body
    return value
}
```

### Simple Examples

```kotlin
fun add(a: Int, b: Int): Int {
    return a + b
}

fun greet(name: String): String {
    return "Hello, $name!"
}

fun printDivider() {
    println("-------------------")
}

fun main() {
    println(add(3, 4))        // 7
    println(greet("Alice"))   // Hello, Alice!
    printDivider()            // -------------------
}
```

### Return Types

The return type is declared after the parameter list with `:`. If a function doesn't return a meaningful value, the return type is `Unit` (which can be omitted):

```kotlin
fun doSomething(): Unit {
    println("doing something")
}

// Equivalent — Unit return type is optional
fun doSomethingElse() {
    println("doing something else")
}
```

### Explicit Returns

Inside a function body, `return` exits the function and optionally provides the return value:

```kotlin
fun absolute(n: Int): Int {
    if (n < 0) return -n
    return n
}

fun firstPositive(numbers: List<Int>): Int? {
    for (n in numbers) {
        if (n > 0) return n   // early return
    }
    return null  // no positive found
}
```

---

## 5.2 Named Arguments

When calling a function, you can specify argument names explicitly. This improves readability and allows you to pass arguments in any order.

```kotlin
fun createUser(name: String, age: Int, email: String, isAdmin: Boolean) {
    println("User: $name, Age: $age, Email: $email, Admin: $isAdmin")
}

// Without named arguments — order matters, hard to read
createUser("Alice", 30, "alice@example.com", false)

// With named arguments — order doesn't matter, meaning is clear
createUser(
    name = "Alice",
    age = 30,
    email = "alice@example.com",
    isAdmin = false
)

// Named arguments in any order
createUser(
    isAdmin = true,
    email = "bob@example.com",
    name = "Bob",
    age = 25
)
```

### Mixing Positional and Named Arguments

```kotlin
fun power(base: Double, exponent: Int): Double {
    // ...
}

// First argument positional, second named
power(2.0, exponent = 10)

// Both named
power(base = 2.0, exponent = 10)
```

**Rule:** Named arguments must follow positional arguments. You cannot use a positional argument after a named one:

```kotlin
// OK: positional first, named after
createUser("Alice", age = 30, email = "a@b.com", isAdmin = false)

// ERROR: named then positional — illegal
// createUser(name = "Alice", 30, "a@b.com", false)
```

### Why Named Arguments Matter

Named arguments are especially valuable for:
1. Functions with many parameters of the same type (avoiding confusion)
2. Boolean flag parameters (making intent clear)
3. Readability at the call site

```kotlin
// BAD — which Boolean is which?
sendEmail("alice@example.com", true, false, true)

// GOOD — named arguments make it clear
sendEmail(
    to = "alice@example.com",
    html = true,
    encrypted = false,
    urgent = true
)
```

---

## 5.3 Default Parameters

Kotlin allows function parameters to have **default values**. This reduces the need for overloaded functions.

### Basic Default Parameters

```kotlin
fun greet(name: String, greeting: String = "Hello", punctuation: String = "!") {
    println("$greeting, $name$punctuation")
}

greet("Alice")                           // Hello, Alice!
greet("Bob", "Hi")                       // Hi, Bob!
greet("Charlie", "Hey", ".")             // Hey, Charlie.
greet("Dave", punctuation = "?")         // Hello, Dave? (skip middle default)
```

### Defaults from Earlier Parameters

Default values can reference earlier parameters:

```kotlin
fun createRange(start: Int, end: Int = start + 10) = start..end

println(createRange(5))      // 5..15
println(createRange(5, 20))  // 5..20
```

### Defaults with Complex Expressions

```kotlin
fun log(
    message: String,
    level: String = "INFO",
    timestamp: Long = System.currentTimeMillis(),
    tag: String = "[${level}]"
) {
    println("$tag $timestamp - $message")
}

log("Application started")
log("Error occurred", level = "ERROR")
log("Custom", level = "DEBUG", timestamp = 0L)
```

### Default vs Overloading

Java's way: overloaded methods

```java
// Java — verbose overloads
void connect(String host) { connect(host, 8080); }
void connect(String host, int port) { connect(host, port, 30); }
void connect(String host, int port, int timeout) { ... }
```

Kotlin's way: default parameters

```kotlin
// Kotlin — one function handles all cases
fun connect(host: String, port: Int = 8080, timeout: Int = 30) {
    println("Connecting to $host:$port with timeout $timeout")
}

connect("example.com")           // host:8080, timeout:30
connect("example.com", 443)      // host:443, timeout:30
connect("example.com", 443, 60)  // host:443, timeout:60
connect("example.com", timeout = 60)  // host:8080, timeout:60
```

### @JvmOverloads for Java Interop

When calling Kotlin from Java, default parameters don't work natively (Java doesn't have them). Use `@JvmOverloads` to generate all overloaded versions:

```kotlin
@JvmOverloads
fun connect(host: String, port: Int = 8080, timeout: Int = 30) {
    println("Connecting to $host:$port with timeout $timeout")
}
// Generates: connect(String), connect(String, int), connect(String, int, int)
```

---

## 5.4 Single-Expression Functions

When a function consists of a single expression, you can use the `=` shorthand and omit the return type (if it can be inferred):

```kotlin
// Full form
fun double(x: Int): Int {
    return x * 2
}

// Single-expression form
fun double(x: Int) = x * 2

// More examples
fun square(x: Int) = x * x
fun isEven(n: Int) = n % 2 == 0
fun max(a: Int, b: Int) = if (a > b) a else b
fun greet(name: String) = "Hello, $name!"

// With explicit return type (when you want to document it)
fun cube(x: Int): Int = x * x * x
```

Single-expression functions are idiomatic Kotlin. They're concise, and the `=` signals "this function IS this expression" rather than "this function DOES this block of code."

### When to Use Single-Expression Style

Use it when:
- The function body is a single, readable expression
- The return type is obvious from the expression

Avoid it when:
- The expression is complex or spans many lines
- You're doing multiple operations

```kotlin
// Good single-expression
fun circleArea(radius: Double) = Math.PI * radius * radius

// Better as a block — too complex for single-expression
fun complexCalculation(a: Int, b: Int, c: Int): Int {
    val intermediate = a * b + c
    val adjusted = if (intermediate < 0) -intermediate else intermediate
    return adjusted * 2 + 1
}
```

---

## 5.5 Varargs

Varargs (variable-length arguments) allow you to pass a variable number of arguments to a function. Declared with `vararg`:

```kotlin
fun sum(vararg numbers: Int): Int {
    var total = 0
    for (n in numbers) total += n
    return total
}

println(sum())               // 0
println(sum(1))              // 1
println(sum(1, 2, 3))        // 6
println(sum(1, 2, 3, 4, 5))  // 15
```

### Vararg Type

Inside the function, the vararg parameter is an `Array`:

```kotlin
fun printAll(vararg items: String) {
    println("Count: ${items.size}")
    for (item in items) println("  - $item")
}

printAll("Apple", "Banana", "Cherry")
// Count: 3
//   - Apple
//   - Banana
//   - Cherry
```

### Spread Operator

To pass an existing array to a vararg function, use the **spread operator** (`*`):

```kotlin
fun sum(vararg numbers: Int) = numbers.sum()

val nums = intArrayOf(1, 2, 3, 4, 5)
println(sum(*nums))  // 15

// Combining with individual elements
println(sum(0, *nums, 6))  // 0 + 1+2+3+4+5 + 6 = 21
```

### Vararg Position

Vararg can be placed anywhere in the parameter list, but only one vararg is allowed per function:

```kotlin
fun format(prefix: String, vararg items: String, suffix: String): String {
    val joined = items.joinToString(", ")
    return "$prefix $joined $suffix"
}

// When calling, items after vararg must be named
println(format("List:", "a", "b", "c", suffix = "."))
// List: a, b, c .
```

---

## 5.6 Local Functions

Kotlin allows defining functions inside other functions. These are called **local functions**:

```kotlin
fun processData(data: List<Int>): List<Int> {
    // Local function — only visible inside processData
    fun validate(n: Int): Boolean {
        return n > 0 && n <= 1000
    }
    
    // Local function can use variables from the enclosing scope
    fun transform(n: Int): Int {
        return n * 2
    }
    
    return data
        .filter { validate(it) }
        .map { transform(it) }
}

println(processData(listOf(-1, 5, 1001, 42, 0, 100)))
// [10, 84, 200]
```

### Local Functions Capturing Outer Scope

Local functions capture variables from their enclosing function (they are closures):

```kotlin
fun buildReport(title: String, items: List<String>): String {
    val sb = StringBuilder()
    
    fun addLine(line: String) {
        sb.appendLine(line)  // captures sb from outer scope
    }
    
    fun addSeparator() {
        addLine("-".repeat(40))  // local functions can call each other
    }
    
    addLine(title)
    addSeparator()
    for (item in items) addLine("• $item")
    addSeparator()
    
    return sb.toString()
}

println(buildReport("Shopping List", listOf("Apples", "Milk", "Bread")))
```

### When to Use Local Functions

- When helper logic is only relevant to one function
- To avoid polluting the class/module with single-use helpers
- When the helper needs direct access to the outer function's variables
- To give meaningful names to complex sub-operations

---

## 5.7 Tail Recursion (tailrec)

Recursive functions can cause stack overflow for deep recursion. Kotlin's `tailrec` modifier converts **tail-recursive** functions into iterative loops, avoiding the stack overflow.

### What is Tail Recursion?

A function is **tail-recursive** if the recursive call is the **last operation** in the function — nothing else happens after the recursive call returns.

```kotlin
// NOT tail-recursive — multiplication happens AFTER the recursive call
fun factorial(n: Long): Long {
    if (n <= 1) return 1
    return n * factorial(n - 1)  // must multiply after returning
}

// Tail-recursive — recursive call IS the last operation
tailrec fun factorialTail(n: Long, accumulator: Long = 1): Long {
    if (n <= 1) return accumulator
    return factorialTail(n - 1, n * accumulator)  // last operation
}
```

### Why tailrec Matters

```kotlin
// Without tailrec: stack overflow for large n
fun factorial(n: Long): Long {
    if (n <= 1) return 1
    return n * factorial(n - 1)
}

// Will throw StackOverflowError for n > ~10000
// println(factorial(100_000))  // CRASH

// With tailrec: converted to a loop, no stack overflow
tailrec fun safeFactorial(n: Long, acc: Long = 1): Long {
    if (n <= 1) return acc
    return safeFactorial(n - 1, n * acc)
}

println(safeFactorial(100_000))  // Works fine (astronomically large number)
```

### More tailrec Examples

```kotlin
// Fibonacci with tailrec
tailrec fun fib(n: Int, a: Long = 0, b: Long = 1): Long {
    if (n == 0) return a
    return fib(n - 1, b, a + b)
}

println(fib(0))   // 0
println(fib(1))   // 1
println(fib(10))  // 55
println(fib(50))  // 12586269025

// Sum with tailrec
tailrec fun sumRange(from: Int, to: Int, acc: Int = 0): Int {
    if (from > to) return acc
    return sumRange(from + 1, to, acc + from)
}

println(sumRange(1, 100))  // 5050

// Finding an element with tailrec
tailrec fun findIndex(list: List<Int>, target: Int, index: Int = 0): Int {
    if (index >= list.size) return -1
    if (list[index] == target) return index
    return findIndex(list, target, index + 1)
}

println(findIndex(listOf(3, 1, 4, 1, 5, 9, 2, 6), 9))  // 5
```

### tailrec Requirements

For `tailrec` to work:
1. The recursive call must be the **last operation** in the function
2. The function must call **itself** (not mutual recursion)
3. The compiler will warn you if it cannot optimize (when the recursion is not actually in tail position)

```kotlin
// DOES NOT WORK — not in tail position
tailrec fun badFactorial(n: Long): Long {
    if (n <= 1) return 1
    return n * badFactorial(n - 1)  // WARNING: function is not tail-recursive
}
```

---

## Functions as Values (Preview)

Kotlin functions are **first-class citizens** — they can be stored in variables and passed as arguments. This is covered fully in Chapter 8 (Functional Programming), but here's a taste:

```kotlin
// Storing a function in a variable
val doubleIt: (Int) -> Int = { x -> x * 2 }
println(doubleIt(5))  // 10

// Function reference using ::
fun greet(name: String) = "Hello, $name!"
val greetFn = ::greet
println(greetFn("Alice"))  // Hello, Alice!

// Passing a function as an argument
fun applyTwice(value: Int, transform: (Int) -> Int): Int {
    return transform(transform(value))
}

println(applyTwice(3, ::square))  // 81 (3 → 9 → 81)
println(applyTwice(3) { it * 2 })  // 12 (3 → 6 → 12)
```

---

## Complete Example: A Mini Math Library

```kotlin
// Top-level math functions
fun factorial(n: Long): Long {
    tailrec fun go(n: Long, acc: Long): Long =
        if (n <= 1) acc else go(n - 1, n * acc)
    return go(n, 1)
}

fun fibonacci(n: Int): Long {
    tailrec fun go(n: Int, a: Long, b: Long): Long =
        if (n == 0) a else go(n - 1, b, a + b)
    return go(n, 0, 1)
}

fun power(base: Double, exponent: Int = 2): Double {
    tailrec fun go(exp: Int, acc: Double): Double =
        if (exp == 0) acc else go(exp - 1, acc * base)
    return go(exponent, 1.0)
}

fun clamp(value: Double, min: Double = 0.0, max: Double = 1.0) =
    when {
        value < min -> min
        value > max -> max
        else -> value
    }

fun main() {
    println(factorial(10))           // 3628800
    println(fibonacci(10))           // 55
    println(power(2.0))              // 4.0
    println(power(2.0, exponent = 10))  // 1024.0
    println(clamp(1.5))              // 1.0
    println(clamp(-0.5))             // 0.0
    println(clamp(0.7, min = 0.5, max = 0.8))  // 0.7
}
```

---

## Summary

Kotlin functions are declared with `fun` and can appear at the top level, in classes, or nested inside other functions. Named arguments make call sites readable when functions have many parameters of the same type. Default parameters replace Java's pattern of overloaded functions. Single-expression syntax (`fun f() = expr`) is idiomatic for simple functions. Varargs accept variable numbers of arguments, with the spread operator (`*`) for passing arrays. Local functions can capture outer scope variables and help structure complex functions. The `tailrec` modifier converts tail-recursive functions to loops, preventing stack overflow.

---

## Key Takeaways

- Functions are top-level citizens — no class wrapper needed
- Named arguments (`f(x = 5)`) improve readability at call sites
- Default parameters (`fun f(x: Int = 0)`) replace most overloaded function patterns
- Single-expression functions (`fun f() = expr`) are idiomatic for simple operations
- Varargs use the `vararg` modifier; use `*` to spread an array
- Local functions see their enclosing scope (they are closures)
- `tailrec` enables safe, efficient deep recursion by converting to a loop
- A function is eligible for `tailrec` only if the recursive call is the last operation

---

## Practice Questions

### Conceptual
1. What is the difference between named arguments and default parameters?
2. When would you prefer `tailrec` over a plain recursive function?
3. What is a local function and when is it appropriate?
4. What does the spread operator (`*`) do?
5. Why can Kotlin use default parameters to replace many overloaded Java methods?

### Code Exercises

**Exercise 1:** Write a function `stats(vararg numbers: Double)` that returns a data class with:
- min, max, sum, average, count
Call it with 5, 10, or no arguments.

**Exercise 2:** Convert this recursive function to a `tailrec` version:
```kotlin
fun sumDigits(n: Int): Int {
    if (n < 10) return n
    return n % 10 + sumDigits(n / 10)
}
```

**Exercise 3:** Write a function `buildString` with these parameters and defaults:
- `prefix: String = ""`
- `value: String`
- `suffix: String = ""`
- `repeat: Int = 1`

**Exercise 4:** Write a local function solution to: given a list of strings, find all strings that are palindromes. The palindrome check should be a local function.

**Exercise 5:** Create a function `format(template: String, vararg args: Any)` that replaces `{}` placeholders in the template with the provided arguments in order:
```kotlin
println(format("Hello, {}! You are {} years old.", "Alice", 30))
// Hello, Alice! You are 30 years old.
```

---

*Next: [Chapter 6 — Object-Oriented Programming](06-oop.md)*
