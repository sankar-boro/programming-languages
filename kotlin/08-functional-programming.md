# Chapter 8 — Functional Programming in Kotlin

> *"Functional programming is a way of programming where you define what to compute, not how to compute it."*

---

## 8.1 Functions as First-Class Citizens

In Kotlin, functions are **first-class citizens**. This means:
- Functions can be stored in variables
- Functions can be passed as arguments to other functions
- Functions can be returned from other functions
- Functions can be created without names (lambdas)

This is a fundamental shift from Java (pre-Java 8), where everything had to be wrapped in objects.

```kotlin
// Storing a function in a variable
val greet: (String) -> String = { name -> "Hello, $name!" }
println(greet("Alice"))  // Hello, Alice!

// Passing a function as an argument
fun execute(operation: () -> Unit) {
    println("Before operation")
    operation()
    println("After operation")
}

execute { println("Doing the work") }
// Before operation
// Doing the work
// After operation

// Returning a function from a function
fun multiplierOf(factor: Int): (Int) -> Int {
    return { number -> number * factor }
}

val double = multiplierOf(2)
val triple = multiplierOf(3)

println(double(5))   // 10
println(triple(5))   // 15
```

---

## 8.2 Lambdas

A **lambda** is an anonymous function — a function without a name, written inline.

### Lambda Syntax

```kotlin
// Full lambda syntax
val add: (Int, Int) -> Int = { x: Int, y: Int -> x + y }

// With type inferred from context
val add: (Int, Int) -> Int = { x, y -> x + y }

// With type on the lambda (no need to declare variable type)
val add = { x: Int, y: Int -> x + y }

println(add(3, 4))  // 7
```

### The `it` Parameter

When a lambda has **exactly one parameter**, you can use the implicit name `it`:

```kotlin
val square: (Int) -> Int = { it * it }
val isEven: (Int) -> Boolean = { it % 2 == 0 }
val shout: (String) -> String = { it.uppercase() + "!" }

println(square(5))      // 25
println(isEven(4))      // true
println(shout("hello")) // HELLO!
```

### Lambda Body

The **last expression** in a lambda is its return value:

```kotlin
val classify: (Int) -> String = { n ->
    val abs = if (n < 0) -n else n
    when {
        abs == 0    -> "zero"
        abs < 10    -> "small"
        abs < 100   -> "medium"
        else        -> "large"
    }
    // last expression is the result — no 'return' keyword
}

println(classify(-3))   // small
println(classify(50))   // medium
println(classify(-500)) // large
```

### Trailing Lambda Syntax

When a lambda is the **last** argument to a function, it can be placed **outside** the parentheses:

```kotlin
// Normal call
listOf(1, 2, 3).forEach({ println(it) })

// Trailing lambda — lambda outside parentheses
listOf(1, 2, 3).forEach { println(it) }

// When lambda is the ONLY argument, parentheses can be omitted entirely
listOf(1, 2, 3).forEach { println(it) }

// Multi-line trailing lambda
listOf("Alice", "Bob", "Charlie").forEach { name ->
    val greeting = "Hello, $name!"
    println(greeting)
}
```

This is the most common way to write lambdas in Kotlin.

### Returning from a Lambda

`return` in a lambda performs a **local return** (returns from the lambda, not the enclosing function). To return from the enclosing function, use a **labeled return**:

```kotlin
fun findFirstEven(numbers: List<Int>): Int? {
    numbers.forEach { n ->
        if (n % 2 == 0) return n  // returns from findFirstEven (non-local return)
        // Note: this only works with inline functions (forEach is inline)
    }
    return null
}

// Labeled return — returns from just the lambda
fun printNonEvens(numbers: List<Int>) {
    numbers.forEach { n ->
        if (n % 2 == 0) return@forEach  // return from lambda only
        println(n)
    }
    println("Done")
}

printNonEvens(listOf(1, 2, 3, 4, 5))
// 1
// 3
// 5
// Done
```

---

## 8.3 Function Types

Function types describe the signature of functions. They're written as `(ParamTypes) -> ReturnType`.

```kotlin
// No parameters, returns Unit
val noParams: () -> Unit = { println("Hello!") }

// One Int parameter, returns String
val intToString: (Int) -> String = { it.toString() }

// Two parameters, returns Boolean
val compareInts: (Int, Int) -> Boolean = { a, b -> a > b }

// Nullable function type
val maybeFunction: ((Int) -> String)? = null

// Calling a nullable function type
maybeFunction?.invoke(42)  // or maybeFunction?(42) — safe call
```

### Function Types as Parameters

```kotlin
fun operate(x: Int, y: Int, operation: (Int, Int) -> Int): Int {
    return operation(x, y)
}

println(operate(5, 3) { a, b -> a + b })  // 8
println(operate(5, 3) { a, b -> a * b })  // 15
println(operate(10, 4) { a, b -> a - b }) // 6
```

### Function Types as Return Values

```kotlin
fun makeCounter(start: Int = 0, step: Int = 1): () -> Int {
    var current = start
    return {
        val result = current
        current += step
        result
    }
}

val counter = makeCounter()
println(counter())  // 0
println(counter())  // 1
println(counter())  // 2

val byFives = makeCounter(0, 5)
println(byFives())  // 0
println(byFives())  // 5
println(byFives())  // 10
```

### Function References

Use `::` to get a reference to an existing function:

```kotlin
fun isEven(n: Int) = n % 2 == 0
fun square(n: Int) = n * n
fun double(n: Int) = n * 2

val numbers = listOf(1, 2, 3, 4, 5, 6, 7, 8, 9, 10)

// Using function references (equivalent to lambdas)
val evens = numbers.filter(::isEven)         // same as filter { isEven(it) }
val squares = numbers.map(::square)          // same as map { square(it) }
val doubled = numbers.map(::double)          // same as map { double(it) }

println(evens)    // [2, 4, 6, 8, 10]
println(squares)  // [1, 4, 9, 16, 25, 36, 49, 64, 81, 100]
println(doubled)  // [2, 4, 6, 8, 10, 12, 14, 16, 18, 20]

// Member function reference
val words = listOf("hello", "world", "kotlin")
val lengths = words.map(String::length)  // [5, 5, 6]
val uppercased = words.map(String::uppercase)
println(uppercased)  // [HELLO, WORLD, KOTLIN]
```

---

## 8.4 Higher-Order Functions

A **higher-order function** is a function that takes another function as a parameter, or returns a function.

### Simple Higher-Order Functions

```kotlin
fun transform(numbers: List<Int>, transformation: (Int) -> Int): List<Int> {
    return numbers.map(transformation)
}

println(transform(listOf(1, 2, 3, 4, 5)) { it * it })
// [1, 4, 9, 16, 25]

println(transform(listOf(1, 2, 3, 4, 5), ::double))
// [2, 4, 6, 8, 10]
```

### Combining Functions (Function Composition)

```kotlin
fun <T, R, V> compose(f: (T) -> R, g: (R) -> V): (T) -> V = { x -> g(f(x)) }

val addOne: (Int) -> Int = { it + 1 }
val double: (Int) -> Int = { it * 2 }

val addOneThenDouble = compose(addOne, double)
val doubleThenAddOne = compose(double, addOne)

println(addOneThenDouble(3))  // (3+1)*2 = 8
println(doubleThenAddOne(3))  // (3*2)+1 = 7
```

### Practical Higher-Order Functions

```kotlin
// Retry logic
fun <T> retry(times: Int, operation: () -> T): T {
    var lastException: Exception? = null
    repeat(times) { attempt ->
        try {
            return operation()  // non-local return from inline function
        } catch (e: Exception) {
            lastException = e
            println("Attempt ${attempt + 1} failed: ${e.message}")
        }
    }
    throw lastException ?: RuntimeException("All attempts failed")
}

// Timing
fun <T> timed(label: String, operation: () -> T): T {
    val start = System.currentTimeMillis()
    val result = operation()
    val elapsed = System.currentTimeMillis() - start
    println("$label took ${elapsed}ms")
    return result
}

val result = timed("Heavy calculation") {
    (1..1_000_000).sum()  // 500000500000
}
println(result)

// Conditional execution
fun <T> runIf(condition: Boolean, block: () -> T): T? =
    if (condition) block() else null

runIf(System.getProperty("debug") != null) {
    println("Debug mode enabled")
}
```

---

## 8.5 Closures

A **closure** is a function that captures variables from its surrounding scope. The captured variables exist as long as the closure exists, even after the enclosing function returns.

```kotlin
fun makeAdder(amount: Int): (Int) -> Int {
    // 'amount' is captured by the returned lambda
    return { n -> n + amount }
}

val add5 = makeAdder(5)
val add10 = makeAdder(10)

println(add5(3))   // 8
println(add10(3))  // 13
println(add5(7))   // 12

// Each closure captures its own 'amount'
println(add5(0))   // 5
println(add10(0))  // 10
```

### Capturing Mutable State

Lambdas in Kotlin can capture and modify `var` variables from their enclosing scope:

```kotlin
fun makeCounter(): () -> Int {
    var count = 0  // captured by the lambda below
    return {
        count++    // modifies the captured variable
        count
    }
}

val counter1 = makeCounter()
val counter2 = makeCounter()  // separate closure, separate 'count'

println(counter1())  // 1
println(counter1())  // 2
println(counter1())  // 3
println(counter2())  // 1 — independent counter
println(counter1())  // 4
```

### Closures in Loops

```kotlin
// Each lambda captures its own copy of 'i' in the loop parameter
val actions = mutableListOf<() -> Unit>()

for (i in 1..5) {
    actions.add { println("Action $i") }  // captures the current value of i
}

actions.forEach { it() }
// Action 1
// Action 2
// Action 3
// Action 4
// Action 5
```

---

## 8.6 Inline Functions

When you pass a lambda to a function, Kotlin creates an **object** for the lambda. This has a performance cost (object creation, virtual dispatch). For performance-critical code, `inline` functions solve this.

### The Problem

```kotlin
fun repeat(times: Int, action: (Int) -> Unit) {
    for (i in 0 until times) action(i)
}

// Every call to repeat() creates a lambda object
repeat(1000) { println(it) }
// Creates an object for the lambda — allocation in a hot loop!
```

### The Solution: inline

```kotlin
inline fun repeat(times: Int, action: (Int) -> Unit) {
    for (i in 0 until times) action(i)
}
```

With `inline`, the compiler copies the function body and the lambda body directly into the call site — no object created, no function call overhead.

```kotlin
// What you write:
repeat(3) { println(it) }

// What the compiler generates (approximately):
for (i in 0 until 3) println(i)
```

### Non-Local Returns Enable

`inline` also enables **non-local returns** — returning from the enclosing function inside a lambda:

```kotlin
inline fun forEach(list: List<Int>, action: (Int) -> Unit) {
    for (item in list) action(item)
}

fun findFirst(numbers: List<Int>, predicate: (Int) -> Boolean): Int? {
    forEach(numbers) { n ->
        if (predicate(n)) return n  // returns from findFirst — only possible because forEach is inline
    }
    return null
}

println(findFirst(listOf(1, 3, 5, 4, 6)) { it % 2 == 0 })  // 4
```

### When to Use inline

Use `inline` when:
- The function is called frequently in hot code paths
- The function primarily passes lambdas to other places
- You need non-local returns

Don't use `inline` for:
- Very large function bodies (code bloat from copying)
- Recursive functions (can't inline recursive functions)

---

## 8.7 crossinline and noinline

### noinline

When a function parameter is `noinline`, the corresponding lambda is NOT inlined. Use this when you need to store the lambda or pass it to another function:

```kotlin
inline fun processWith(
    noinline storedAction: (String) -> Unit,   // NOT inlined — can be stored
    immediateAction: (String) -> Unit           // inlined
) {
    val stored = storedAction  // OK — noinline lambda can be treated as object
    immediateAction("processed")
    stored("stored")
}

processWith(
    storedAction = { println("Stored: $it") },
    immediateAction = { println("Immediate: $it") }
)
// Immediate: processed
// Stored: stored
```

### crossinline

When a lambda marked `crossinline` cannot have non-local returns, because it's called in a different execution context (e.g., inside another lambda):

```kotlin
inline fun runOnMainThread(crossinline block: () -> Unit) {
    // This lambda runs in a different context (e.g., event queue)
    executeOnMain { block() }  // crossinline: block can't use non-local return
}

fun example() {
    runOnMainThread {
        println("On main thread")
        // return  // ERROR: can't non-locally return from crossinline lambda
    }
}
```

---

## 8.8 SAM Conversions

**SAM** stands for **Single Abstract Method**. A SAM interface is an interface with exactly one abstract method (like Java's functional interfaces).

Kotlin allows you to pass a lambda wherever a SAM interface is expected:

```kotlin
// Java-defined functional interface
// interface Runnable { void run(); }

// Java way of using Runnable
val runnable = Runnable { println("Running!") }
Thread(runnable).start()

// Kotlin SAM conversion — pass lambda directly
Thread { println("Running in thread!") }.start()
```

### Kotlin fun Interface

Kotlin allows defining SAM interfaces using the `fun interface` keyword:

```kotlin
fun interface Validator<T> {
    fun validate(value: T): Boolean
}

// Create with lambda
val ageValidator = Validator<Int> { age -> age in 0..150 }
val emailValidator = Validator<String> { email -> "@" in email }

println(ageValidator.validate(25))      // true
println(ageValidator.validate(-5))      // false
println(emailValidator.validate("alice@example.com"))  // true
println(emailValidator.validate("notanemail"))          // false
```

### Composing SAM Interfaces

```kotlin
fun interface Predicate<T> {
    fun test(value: T): Boolean
    
    fun and(other: Predicate<T>): Predicate<T> = Predicate { test(it) && other.test(it) }
    fun or(other: Predicate<T>): Predicate<T> = Predicate { test(it) || other.test(it) }
    fun negate(): Predicate<T> = Predicate { !test(it) }
}

val isPositive = Predicate<Int> { it > 0 }
val isEven = Predicate<Int> { it % 2 == 0 }
val isPositiveEven = isPositive.and(isEven)

println(isPositiveEven.test(4))   // true
println(isPositiveEven.test(-4))  // false
println(isPositiveEven.test(3))   // false
```

---

## Scope Functions: apply, let, run, with, also

Kotlin's standard library includes scope functions that are higher-order functions used to execute a block on an object. They differ in how they refer to the object and what they return.

| Function | Object reference | Return value | Extension function? |
|----------|-----------------|--------------|---------------------|
| `let`    | `it`            | Lambda result | Yes |
| `run`    | `this`          | Lambda result | Yes |
| `with`   | `this`          | Lambda result | No (regular function) |
| `apply`  | `this`          | The object   | Yes |
| `also`   | `it`            | The object   | Yes |

```kotlin
// let — transform or run with nullable
val name: String? = "Alice"
val length = name?.let {
    println("Processing: $it")
    it.length  // return value
}
println(length)  // 5

// run — execute block and return result
val result = "Hello, World!".run {
    val upper = uppercase()
    val words = split(", ")
    words.size  // return value
}
println(result)  // 2

// with — useful when you have many operations on one object
val message = with(StringBuilder()) {
    append("Hello")
    append(", ")
    append("World!")
    toString()  // return value
}
println(message)  // Hello, World!

// apply — configure an object
val person = Person("", 0).apply {
    name = "Alice"
    age = 30
}
// person is the return value

// also — side effects, returns original object
val numbers = mutableListOf(1, 2, 3)
    .also { println("Original: $it") }    // prints list
    .also { it.add(4) }                   // adds to list
    .also { println("Modified: $it") }    // prints modified list

println(numbers)  // [1, 2, 3, 4]
```

---

## Complete Example: Functional Pipeline

```kotlin
data class Employee(
    val name: String,
    val department: String,
    val salary: Double,
    val yearsOfExperience: Int
)

fun main() {
    val employees = listOf(
        Employee("Alice", "Engineering", 95000.0, 8),
        Employee("Bob", "Marketing", 72000.0, 3),
        Employee("Charlie", "Engineering", 115000.0, 12),
        Employee("Diana", "HR", 68000.0, 5),
        Employee("Eve", "Engineering", 88000.0, 6),
        Employee("Frank", "Marketing", 79000.0, 7)
    )
    
    // Find senior engineers earning above average
    val avgSalary = employees.map { it.salary }.average()
    
    val seniorEngineers = employees
        .filter { it.department == "Engineering" }
        .filter { it.yearsOfExperience >= 7 }
        .filter { it.salary > avgSalary }
        .sortedByDescending { it.salary }
    
    println("Senior engineers above average salary (${"%.0f".format(avgSalary)}):")
    seniorEngineers.forEach { emp ->
        println("  ${emp.name}: \$${emp.salary} (${emp.yearsOfExperience} years)")
    }
    
    // Department statistics
    val deptStats = employees
        .groupBy { it.department }
        .mapValues { (dept, emps) ->
            val avgSal = emps.map { it.salary }.average()
            val avgExp = emps.map { it.yearsOfExperience }.average()
            "avg salary: ${"%.0f".format(avgSal)}, avg experience: ${"%.1f".format(avgExp)} years"
        }
    
    println("\nDepartment stats:")
    deptStats.forEach { (dept, stats) -> println("  $dept: $stats") }
    
    // Compose transformations
    val raiseCalculator: (Double) -> (Double) -> Double = { percent ->
        { salary -> salary * (1 + percent / 100) }
    }
    
    val giveRaise = raiseCalculator(10.0)  // 10% raise
    val raisedSalaries = employees.map { emp ->
        emp.copy(salary = giveRaise(emp.salary))
    }
    
    println("\nAfter 10% raise:")
    raisedSalaries.forEach { println("  ${it.name}: \$${it.salary}") }
}
```

---

## Summary

Kotlin treats functions as first-class values — they can be stored, passed, and returned. Lambdas are anonymous functions written inline with `{ params -> body }` syntax. Function types (`(Int) -> String`) describe function signatures. Higher-order functions take or return functions, enabling powerful composition patterns. Closures capture their surrounding scope. `inline` functions eliminate lambda object creation overhead and enable non-local returns. `crossinline` prevents non-local returns in specific contexts; `noinline` opts specific lambdas out of inlining. SAM conversions allow passing lambdas where Java functional interfaces are expected.

---

## Key Takeaways

- Functions are first-class values in Kotlin — store, pass, and return them
- Lambda syntax: `{ param -> body }` or `{ it }` for single parameters
- Trailing lambda: `f { body }` when lambda is the last argument
- Function references: `::functionName` or `Type::methodName`
- `inline` eliminates lambda allocation overhead — use in hot code paths
- Closures capture their enclosing scope, including mutable variables
- `fun interface` creates SAM interfaces usable with lambda syntax
- Scope functions (`let`, `run`, `with`, `apply`, `also`) differ in how they pass the object and what they return

---

## Practice Questions

### Conceptual
1. What is the difference between a lambda and a function reference?
2. When would you use `inline` and what are the trade-offs?
3. What is the `it` convention in Kotlin lambdas?
4. What is the difference between `apply` and `also`?
5. What is a non-local return and why does it require `inline`?

### Code Exercises

**Exercise 1:** Implement `compose(f, g)` that returns `g(f(x))`. Test with a chain of three transformations.

**Exercise 2:** Write a function `memoize` that takes a function `(Int) -> Int` and returns a memoized version that caches results. Use a closure and a map.

**Exercise 3:** Implement your own `myFilter` and `myMap` as higher-order functions (without using the standard library versions).

**Exercise 4:** Using scope functions, rewrite this code idiomatically:
```kotlin
val sb = StringBuilder()
sb.append("Name: ")
sb.append(user.name)
sb.append("\n")
sb.append("Age: ")
sb.append(user.age)
val result = sb.toString()
```

**Exercise 5:** Create a `Pipeline` class that chains transformations:
```kotlin
val pipeline = Pipeline<Int>()
    .then { it * 2 }
    .then { it + 1 }
    .then { it.toString() }

println(pipeline.execute(5))  // "11"
```

---

*Next: [Chapter 9 — Collections, Transformations, and Sequences](09-collections.md)*
