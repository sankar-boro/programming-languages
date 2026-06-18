# Chapter 9 — Collections, Transformations, and Sequences

> *"The most powerful programming tool in your toolkit is the ability to transform data simply and expressively."*

---

## 9.1 The Collections Hierarchy

Kotlin's collections are organized in a clear hierarchy:

```
Iterable<T>
├── Collection<T>
│   ├── List<T>        (ordered, indexed, allows duplicates)
│   ├── Set<T>         (unordered, no duplicates)
│   └── MutableCollection<T>
│       ├── MutableList<T>
│       └── MutableSet<T>
└── Map<K, V>          (key-value pairs, keys are unique)
    └── MutableMap<K, V>
```

**Key insight:** Kotlin separates **read interfaces** (`Collection`, `List`, `Set`, `Map`) from **write interfaces** (`MutableCollection`, `MutableList`, `MutableSet`, `MutableMap`). This enforces the principle of least privilege at the type level.

---

## 9.2 Lists, Sets, and Maps

### Lists

A `List` is an ordered collection with indexed access. Elements can be duplicated.

```kotlin
// Immutable List
val fruits = listOf("Apple", "Banana", "Cherry", "Apple")
println(fruits)           // [Apple, Banana, Cherry, Apple]
println(fruits[0])        // Apple
println(fruits.size)      // 4
println(fruits.contains("Banana"))  // true
println(fruits.indexOf("Apple"))    // 0
println(fruits.lastIndexOf("Apple")) // 3

// Mutable List
val mutableFruits = mutableListOf("Apple", "Banana")
mutableFruits.add("Cherry")
mutableFruits.add(0, "Avocado")  // insert at index 0
mutableFruits.removeAt(1)        // remove by index
mutableFruits.remove("Cherry")   // remove by value
println(mutableFruits)  // [Avocado, Banana]
```

### Sets

A `Set` is a collection with no duplicate elements. Order is not guaranteed for `HashSet`.

```kotlin
// Immutable Set
val nums = setOf(1, 2, 3, 2, 1)  // duplicates removed
println(nums)          // [1, 2, 3]
println(nums.size)     // 3
println(2 in nums)     // true
println(5 in nums)     // false

// Mutable Set
val primes = mutableSetOf(2, 3, 5, 7)
primes.add(11)
primes.add(7)   // duplicate — not added
primes.remove(2)
println(primes) // [3, 5, 7, 11]

// Set operations
val a = setOf(1, 2, 3, 4)
val b = setOf(3, 4, 5, 6)

println(a union b)        // {1, 2, 3, 4, 5, 6}
println(a intersect b)    // {3, 4}
println(a subtract b)     // {1, 2}

// LinkedHashSet — maintains insertion order
val ordered = linkedSetOf("Charlie", "Alice", "Bob")
println(ordered)  // [Charlie, Alice, Bob]

// TreeSet-equivalent — sorted order
val sorted = sortedSetOf(3, 1, 4, 1, 5, 9, 2, 6)
println(sorted)  // [1, 2, 3, 4, 5, 6, 9]
```

### Maps

A `Map` stores key-value pairs. Keys are unique; values can repeat.

```kotlin
// Immutable Map
val capitals = mapOf(
    "France" to "Paris",
    "Germany" to "Berlin",
    "Japan" to "Tokyo"
)

println(capitals["France"])           // Paris
println(capitals.getOrDefault("UK", "Unknown"))  // Unknown
println(capitals.containsKey("Japan"))  // true
println(capitals.size)               // 3

// Iterating
for ((country, capital) in capitals) {
    println("$country → $capital")
}

// Mutable Map
val scores = mutableMapOf("Alice" to 90, "Bob" to 85)
scores["Charlie"] = 92
scores["Alice"] = 95  // update existing
scores.remove("Bob")
println(scores)  // {Alice=95, Charlie=92}

// getOrPut — get existing or compute and store
val wordCount = mutableMapOf<String, Int>()
val words = "the quick brown fox jumps over the lazy dog".split(" ")
for (word in words) {
    wordCount[word] = wordCount.getOrDefault(word, 0) + 1
}
println(wordCount)  // {the=2, quick=1, brown=1, ...}
```

---

## 9.3 Mutable vs Immutable Collections

This is one of Kotlin's most important design decisions. The **read-only view** of a collection does NOT mean the collection is **truly immutable** — it just means you can't modify it through that reference.

```kotlin
// mutableList is MutableList<Int>
val mutableList = mutableListOf(1, 2, 3)

// readOnlyView is List<Int> — read-only VIEW of the same list
val readOnlyView: List<Int> = mutableList

println(readOnlyView)  // [1, 2, 3]
mutableList.add(4)
println(readOnlyView)  // [1, 2, 3, 4] — reflects the change!

// To get a truly immutable copy:
val immutableCopy = readOnlyView.toList()
mutableList.add(5)
println(readOnlyView)   // [1, 2, 3, 4, 5] — changed again
println(immutableCopy)  // [1, 2, 3, 4] — unchanged copy
```

### Creating Collections

```kotlin
// Creating Lists
val empty = emptyList<Int>()          // Empty, truly immutable
val single = listOf(42)               // Single element
val from1to10 = (1..10).toList()      // From range
val filled = List(5) { it * it }      // [0, 1, 4, 9, 16] — generated

// Creating Sets
val emptySet = emptySet<String>()
val fromList = listOf(1, 1, 2, 2, 3).toSet()  // removes duplicates: {1, 2, 3}

// Creating Maps
val emptyMap = emptyMap<String, Int>()
val fromPairs = mapOf(Pair("a", 1), Pair("b", 2))
val fromTo = mapOf("a" to 1, "b" to 2)     // 'to' creates Pair
val fromKeys = listOf("a", "b", "c").associateWith { it.length }  // {a=1, b=1, c=1}
```

---

## 9.4 Collection Transformations

The Kotlin standard library provides a rich set of transformation functions. These always return **new collections** — they don't modify the original.

### map — Transform Each Element

```kotlin
val numbers = listOf(1, 2, 3, 4, 5)

val squares = numbers.map { it * it }
println(squares)  // [1, 4, 9, 16, 25]

val strings = numbers.map { "item_$it" }
println(strings)  // [item_1, item_2, item_3, item_4, item_5]

// mapIndexed — also provides the index
val indexed = numbers.mapIndexed { index, value -> "[$index]=$value" }
println(indexed)  // [[0]=1, [1]=2, [2]=3, [3]=4, [4]=5]

// mapNotNull — maps and filters out nulls
val maybeNumbers = listOf("1", "two", "3", "four", "5")
val parsedNumbers = maybeNumbers.mapNotNull { it.toIntOrNull() }
println(parsedNumbers)  // [1, 3, 5]
```

### flatMap — Map and Flatten

```kotlin
val sentences = listOf("Hello World", "Kotlin is great")

// map gives a List of Lists
val words1 = sentences.map { it.split(" ") }
println(words1)  // [[Hello, World], [Kotlin, is, great]]

// flatMap flattens the result
val words2 = sentences.flatMap { it.split(" ") }
println(words2)  // [Hello, World, Kotlin, is, great]

// Another example: expand numbers to ranges
val ranges = listOf(1..3, 5..7, 9..10)
val allNums = ranges.flatMap { it.toList() }
println(allNums)  // [1, 2, 3, 5, 6, 7, 9, 10]
```

---

## 9.5 Filtering

### filter — Keep Elements Matching a Condition

```kotlin
val numbers = listOf(1, 2, 3, 4, 5, 6, 7, 8, 9, 10)

val evens = numbers.filter { it % 2 == 0 }
println(evens)  // [2, 4, 6, 8, 10]

val large = numbers.filter { it > 5 }
println(large)  // [6, 7, 8, 9, 10]

// filterNot — keep elements NOT matching
val odds = numbers.filterNot { it % 2 == 0 }
println(odds)  // [1, 3, 5, 7, 9]

// filterIsInstance — filter and cast by type
val mixed: List<Any> = listOf(1, "two", 3, "four", 5.0, true)
val strings = mixed.filterIsInstance<String>()
println(strings)  // [two, four]

val ints = mixed.filterIsInstance<Int>()
println(ints)  // [1, 3]
```

### take and drop

```kotlin
val numbers = listOf(1, 2, 3, 4, 5, 6, 7, 8, 9, 10)

println(numbers.take(3))           // [1, 2, 3]
println(numbers.drop(7))           // [8, 9, 10]
println(numbers.takeLast(3))       // [8, 9, 10]
println(numbers.dropLast(7))       // [1, 2, 3]
println(numbers.takeWhile { it < 5 })  // [1, 2, 3, 4]
println(numbers.dropWhile { it < 5 })  // [5, 6, 7, 8, 9, 10]
```

### Partition

Split a collection into two based on a predicate:

```kotlin
val numbers = listOf(1, 2, 3, 4, 5, 6, 7, 8, 9, 10)
val (evens, odds) = numbers.partition { it % 2 == 0 }
println(evens)  // [2, 4, 6, 8, 10]
println(odds)   // [1, 3, 5, 7, 9]
```

---

## 9.6 Aggregation

### reduce — Accumulate Using First Element as Seed

```kotlin
val numbers = listOf(1, 2, 3, 4, 5)

val sum = numbers.reduce { acc, n -> acc + n }
println(sum)  // 15

val product = numbers.reduce { acc, n -> acc * n }
println(product)  // 120

val max = numbers.reduce { acc, n -> if (n > acc) n else acc }
println(max)  // 5

// reduce throws if collection is empty!
// Use reduceOrNull for safety:
val empty = emptyList<Int>()
val safeSum = empty.reduceOrNull { acc, n -> acc + n }
println(safeSum)  // null
```

### fold — Accumulate with an Initial Value

```kotlin
val numbers = listOf(1, 2, 3, 4, 5)

// fold is safer than reduce — takes initial value
val sum = numbers.fold(0) { acc, n -> acc + n }
println(sum)  // 15

// Works on empty collections
val emptySum = emptyList<Int>().fold(0) { acc, n -> acc + n }
println(emptySum)  // 0

// Build a complex result
val wordLengths = listOf("Hello", "World", "Kotlin")
val result = wordLengths.fold(mutableMapOf<Int, MutableList<String>>()) { acc, word ->
    acc.getOrPut(word.length) { mutableListOf() }.add(word)
    acc
}
println(result)  // {5=[Hello, World], 6=[Kotlin]}

// foldRight — from right to left
val sentence = listOf("Hello", "World", "Kotlin")
val joined = sentence.foldRight("") { word, acc ->
    if (acc.isEmpty()) word else "$word $acc"
}
println(joined)  // Hello World Kotlin
```

### Built-in Aggregation Functions

```kotlin
val numbers = listOf(3, 1, 4, 1, 5, 9, 2, 6, 5, 3, 5)

println(numbers.sum())            // 44
println(numbers.count())          // 11
println(numbers.average())        // 4.0
println(numbers.min())            // 1
println(numbers.max())            // 9
println(numbers.minOrNull())      // 1 (null-safe version)
println(numbers.maxOrNull())      // 9

// sumOf, countOf, minOf, maxOf with selectors
data class Product(val name: String, val price: Double, val quantity: Int)

val cart = listOf(
    Product("Apple", 1.50, 3),
    Product("Milk", 2.99, 2),
    Product("Bread", 3.49, 1)
)

val totalCost = cart.sumOf { it.price * it.quantity }
println("Total: \$${"%.2f".format(totalCost)}")  // Total: $14.97

val mostExpensive = cart.maxByOrNull { it.price }
println("Most expensive: ${mostExpensive?.name}")  // Most expensive: Bread

val countAffordable = cart.count { it.price < 3.0 }
println("Affordable items: $countAffordable")  // Affordable items: 2
```

---

## 9.7 Grouping and Partitioning

### groupBy — Group Elements by a Key

```kotlin
data class Person(val name: String, val age: Int, val city: String)

val people = listOf(
    Person("Alice", 30, "New York"),
    Person("Bob", 25, "London"),
    Person("Charlie", 30, "New York"),
    Person("Diana", 25, "Tokyo"),
    Person("Eve", 35, "London")
)

// Group by city
val byCity = people.groupBy { it.city }
for ((city, persons) in byCity) {
    println("$city: ${persons.map { it.name }}")
}
// New York: [Alice, Charlie]
// London: [Bob, Eve]
// Tokyo: [Diana]

// Group by age, extract only names
val namesByAge = people.groupBy(
    keySelector = { it.age },
    valueTransform = { it.name }
)
println(namesByAge)  // {30=[Alice, Charlie], 25=[Bob, Diana], 35=[Eve]}
```

### associate — Create Maps from Collections

```kotlin
val words = listOf("apple", "banana", "cherry")

// Map each word to its length
val wordLengths = words.associateWith { it.length }
println(wordLengths)  // {apple=5, banana=6, cherry=6}

// Map length to first word with that length
val lengthToWord = words.associateBy { it.length }
println(lengthToWord)  // {5=apple, 6=cherry} — note: banana overwritten

// Full control with associate
val indexedWords = words.associate { word ->
    word.first() to word.uppercase()
}
println(indexedWords)  // {a=APPLE, b=BANANA, c=CHERRY}
```

### Sorted Collections

```kotlin
val numbers = listOf(3, 1, 4, 1, 5, 9, 2, 6)

println(numbers.sorted())           // [1, 1, 2, 3, 4, 5, 6, 9]
println(numbers.sortedDescending()) // [9, 6, 5, 4, 3, 2, 1, 1]

data class Person(val name: String, val age: Int)
val people = listOf(
    Person("Charlie", 30),
    Person("Alice", 25),
    Person("Bob", 30)
)

// Sort by age, then by name
val sorted = people
    .sortedBy { it.age }
    .thenSortedBy { it.name }  // chaining doesn't work — use sortedWith

val properSorted = people.sortedWith(compareBy({ it.age }, { it.name }))
println(properSorted)
// [Person(Alice, 25), Person(Bob, 30), Person(Charlie, 30)]
```

---

## 9.8 Sequences and Lazy Evaluation

Collections are **eager** — every operation creates a new collection immediately. For large collections or long chains, this is wasteful.

**Sequences** are **lazy** — operations are not executed until the result is actually needed.

### The Problem with Eager Collections

```kotlin
val numbers = (1..1_000_000)

// Eager evaluation: creates three complete intermediate lists
val result = numbers.toList()
    .map { it * 2 }        // creates list of 1,000,000 elements
    .filter { it > 100 }   // creates another big list
    .take(5)               // finally takes just 5

// That's 2+ million unnecessary objects created!
```

### Sequences to the Rescue

```kotlin
val numbers = (1..1_000_000)

// Lazy evaluation: processes elements one at a time
val result = numbers.asSequence()
    .map { it * 2 }        // no list created yet
    .filter { it > 100 }   // no list created yet
    .take(5)               // no list created yet — terminal operation triggers
    .toList()              // NOW it processes — and STOPS after finding 5

println(result)  // [102, 104, 106, 108, 110]
// Processed only ~56 elements, not 1,000,000!
```

### Creating Sequences

```kotlin
// From collections
val seq1 = listOf(1, 2, 3).asSequence()

// Generate a sequence
val seq2 = sequence {
    yield(1)          // yield one value
    yield(2)
    yieldAll(listOf(3, 4, 5))  // yield multiple values
}

println(seq2.toList())  // [1, 2, 3, 4, 5]

// Infinite sequence (terminated by take)
val naturalNumbers = generateSequence(1) { it + 1 }
println(naturalNumbers.take(10).toList())  // [1, 2, 3, 4, 5, 6, 7, 8, 9, 10]

// Fibonacci sequence
val fibonacci = generateSequence(Pair(0L, 1L)) { (a, b) -> Pair(b, a + b) }
    .map { it.first }

println(fibonacci.take(15).toList())
// [0, 1, 1, 2, 3, 5, 8, 13, 21, 34, 55, 89, 144, 233, 377]

// File lines as a sequence (memory efficient)
// File("large.txt").bufferedReader().lineSequence().filter { "error" in it }.take(10)
```

### Sequence vs Collection Performance

```kotlin
import kotlin.system.measureTimeMillis

val data = (1..10_000_000).toList()

val eagerTime = measureTimeMillis {
    data.map { it * 2 }.filter { it % 3 == 0 }.take(10).toList()
}

val lazyTime = measureTimeMillis {
    data.asSequence().map { it * 2 }.filter { it % 3 == 0 }.take(10).toList()
}

println("Eager: ${eagerTime}ms")   // Eager: ~500ms (varies)
println("Lazy:  ${lazyTime}ms")    // Lazy:  ~1ms
```

### When to Use Sequences

**Use sequences when:**
- Processing large collections
- Using multiple chained transformations
- Potentially short-circuiting (take, first, any, find)
- Working with infinite data

**Use collections when:**
- Collection is small (< ~1000 elements typically)
- Single transformation
- Random access is needed
- Multiple traversals of the result

### Terminal Operations on Sequences

A sequence is lazy until you call a **terminal operation** that needs the result:

```kotlin
val seq = generateSequence(1) { it + 1 }

// These are terminal operations — they trigger evaluation:
seq.take(5).toList()           // [1, 2, 3, 4, 5]
seq.take(5).sum()              // 15
seq.first { it > 100 }         // 101
seq.take(1000).count()         // 1000
seq.take(10).average()         // 5.5
```

---

## 9.9 Useful Collection Utilities

```kotlin
val numbers = listOf(3, 1, 4, 1, 5, 9, 2, 6, 5, 3, 5)

// Finding elements
println(numbers.first())               // 3
println(numbers.last())                // 5
println(numbers.firstOrNull { it > 10 })  // null
println(numbers.find { it > 4 })       // 5
println(numbers.findLast { it > 4 })   // 5

// Checking conditions
println(numbers.any { it > 8 })        // true
println(numbers.all { it > 0 })        // true
println(numbers.none { it < 0 })       // true

// Transforming to different structures
println(numbers.distinct())            // [3, 1, 4, 5, 9, 2, 6]
println(numbers.toSet())               // {3, 1, 4, 5, 9, 2, 6}

// Combining collections
val a = listOf(1, 2, 3)
val b = listOf(4, 5, 6)
println(a + b)           // [1, 2, 3, 4, 5, 6]
println(a.zip(b))        // [(1, 4), (2, 5), (3, 6)]

val zipped = a.zip(b) { x, y -> x + y }
println(zipped)          // [5, 7, 9]

// Unzip
val pairs = listOf(1 to "one", 2 to "two", 3 to "three")
val (ints, strings) = pairs.unzip()
println(ints)    // [1, 2, 3]
println(strings) // [one, two, three]

// Chunking
val nums = (1..10).toList()
println(nums.chunked(3))
// [[1, 2, 3], [4, 5, 6], [7, 8, 9], [10]]

println(nums.windowed(3))
// [[1, 2, 3], [2, 3, 4], [3, 4, 5], ..., [8, 9, 10]]

println(nums.windowed(3, step = 2))
// [[1, 2, 3], [3, 4, 5], [5, 6, 7], [7, 8, 9]]
```

---

## Complete Example: Data Processing Pipeline

```kotlin
data class Sale(
    val date: String,
    val product: String,
    val category: String,
    val amount: Double,
    val quantity: Int
)

fun analyzeSales(sales: List<Sale>) {
    println("=== Sales Analysis ===\n")
    
    // Total revenue
    val totalRevenue = sales.sumOf { it.amount * it.quantity }
    println("Total Revenue: \$${"%.2f".format(totalRevenue)}")
    
    // Revenue by category
    val byCategory = sales
        .groupBy { it.category }
        .mapValues { (_, categorySales) ->
            categorySales.sumOf { it.amount * it.quantity }
        }
        .toSortedMap()
    
    println("\nRevenue by Category:")
    byCategory.forEach { (cat, rev) ->
        println("  $cat: \$${"%.2f".format(rev)}")
    }
    
    // Top 3 products
    val topProducts = sales
        .groupBy { it.product }
        .mapValues { (_, prodSales) -> prodSales.sumOf { it.amount * it.quantity } }
        .entries
        .sortedByDescending { it.value }
        .take(3)
    
    println("\nTop 3 Products:")
    topProducts.forEachIndexed { i, (product, revenue) ->
        println("  ${i + 1}. $product: \$${"%.2f".format(revenue)}")
    }
    
    // Average order value
    val avgOrderValue = sales
        .map { it.amount * it.quantity }
        .average()
    println("\nAverage Order Value: \$${"%.2f".format(avgOrderValue)}")
    
    // Products sold more than 5 times
    val popularProducts = sales
        .groupBy { it.product }
        .filter { (_, sales) -> sales.sumOf { it.quantity } > 5 }
        .keys
        .sorted()
    println("\nPopular Products (>5 units): $popularProducts")
}

fun main() {
    val sales = listOf(
        Sale("2024-01", "Apple", "Fruit", 1.50, 10),
        Sale("2024-01", "Bread", "Bakery", 3.49, 3),
        Sale("2024-01", "Milk", "Dairy", 2.99, 5),
        Sale("2024-02", "Apple", "Fruit", 1.50, 8),
        Sale("2024-02", "Cheese", "Dairy", 5.99, 2),
        Sale("2024-02", "Bread", "Bakery", 3.49, 4),
        Sale("2024-03", "Apple", "Fruit", 1.50, 12),
        Sale("2024-03", "Milk", "Dairy", 2.99, 6)
    )
    
    analyzeSales(sales)
}
```

---

## Summary

Kotlin's collection framework separates read-only (`List`, `Set`, `Map`) from mutable (`MutableList`, `MutableSet`, `MutableMap`) interfaces. Collection operations (`map`, `filter`, `fold`, `groupBy`, etc.) are eager by default — creating new collections for each step. `asSequence()` switches to lazy evaluation — perfect for large collections and long transformation chains. `generateSequence` and the `sequence {}` builder create potentially infinite lazy sequences.

---

## Key Takeaways

- `listOf`, `setOf`, `mapOf` create read-only collections; `mutableListOf` etc. create mutable ones
- Read-only ≠ truly immutable — another reference may still modify the underlying data
- `map` transforms every element; `filter` keeps matching elements; `flatMap` transforms and flattens
- `fold` accumulates with an initial value (safe for empty collections); `reduce` uses the first element (throws on empty)
- `groupBy` returns a `Map<K, List<V>>`; `associateWith` creates a value map; `associateBy` creates a key map
- Sequences are lazy — no evaluation until a terminal operation is called
- Use sequences for large data or chains with early termination; use collections for small data

---

## Practice Questions

### Conceptual
1. What is the difference between `List<T>` and `MutableList<T>` in Kotlin?
2. Why is a read-only `List` not necessarily immutable?
3. What is the difference between `map` and `flatMap`?
4. What is the difference between `reduce` and `fold`?
5. When should you prefer `asSequence()` over regular collection operations?

### Code Exercises

**Exercise 1:** Given a list of strings, write a pipeline that:
- Removes empty strings
- Trims whitespace
- Converts to lowercase
- Removes duplicates
- Sorts alphabetically

**Exercise 2:** Given a list of `Transaction(amount: Double, type: String)`, compute:
- Total credits (type = "credit")
- Total debits (type = "debit")
- Net balance

**Exercise 3:** Implement `mostFrequent(words: List<String>): String` using `groupBy`, returning the word that appears most often.

**Exercise 4:** Using `generateSequence`, create:
- A sequence of powers of 2: 1, 2, 4, 8, 16, ...
- The first 10 powers greater than 100

**Exercise 5:** Rewrite this eager chain as a sequence and predict the performance difference:
```kotlin
(1..10_000_000).toList()
    .map { it.toString() }
    .filter { it.startsWith("9") }
    .first()
```

---

*Next: [Chapter 10 — Advanced Kotlin Features](10-advanced-features.md)*
