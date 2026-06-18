# Common Pitfalls — Mistakes to Avoid

> *"Every mistake is a learning opportunity. Learning from others' mistakes is even better."*

---

## Overview

Even experienced developers run into Kotlin-specific traps. This chapter documents the most common mistakes — from beginner gotchas to subtle performance issues — so you can recognize and avoid them.

---

## Pitfall 1: Misusing val for Mutable Objects

**The Mistake:**
```kotlin
val list = mutableListOf(1, 2, 3)
list.add(4)   // No error! val doesn't mean the object is immutable
list.add(5)

println(list)  // [1, 2, 3, 4, 5]
```

**Why It's a Pitfall:**
Developers coming from Java think `val` means "constant" or "immutable." In Kotlin, `val` means the **reference** cannot be reassigned — the object it points to can still be mutable.

**The Fix:**
If you want a truly immutable list, use `listOf()` (not `mutableListOf()`):
```kotlin
val list = listOf(1, 2, 3)
// list.add(4)  // compile error — no add() on List

// For true immutability with a copy:
val list2 = listOf(1, 2, 3)
val list3 = list2 + listOf(4, 5)  // creates a new list
```

---

## Pitfall 2: The !! Operator Overuse

**The Mistake:**
```kotlin
// Suppressing all null safety with !!
val user = getUser(id)!!
val address = user.getAddress()!!
val city = address.getCity()!!
println(city.length)
```

**Why It's a Pitfall:**
`!!` throws a NullPointerException if the value is null — exactly what Kotlin's null safety is designed to prevent. Overusing `!!` removes all null safety guarantees.

**The Fix:**
```kotlin
// Use safe navigation and defaults
val cityLength = getUser(id)?.address?.city?.length ?: 0

// Or handle null explicitly with meaningful errors
val user = getUser(id) 
    ?: throw IllegalArgumentException("User $id not found")
val address = user.address 
    ?: throw IllegalStateException("User has no address")
val city = address.city 
    ?: throw IllegalStateException("Address has no city")
```

---

## Pitfall 3: Mutable Shared State in Coroutines

**The Mistake:**
```kotlin
var counter = 0

fun main() = runBlocking {
    val jobs = (1..1000).map {
        launch(Dispatchers.Default) {
            counter++  // RACE CONDITION — not thread-safe!
        }
    }
    jobs.forEach { it.join() }
    println(counter)  // Not 1000! — some increments are lost
}
```

**Why It's a Pitfall:**
`counter++` is NOT atomic. Multiple coroutines on different threads can read the same value, both increment it, and write the same new value — losing an increment.

**The Fix:**
```kotlin
import java.util.concurrent.atomic.AtomicInteger

val counter = AtomicInteger(0)

fun main() = runBlocking {
    val jobs = (1..1000).map {
        launch(Dispatchers.Default) {
            counter.incrementAndGet()  // atomic — thread-safe
        }
    }
    jobs.forEach { it.join() }
    println(counter.get())  // Always 1000
}

// Or use Mutex for complex operations:
val mutex = Mutex()
var counter = 0

launch {
    mutex.withLock {
        counter++  // protected by lock
    }
}
```

---

## Pitfall 4: Smart Cast Fails on var

**The Mistake:**
```kotlin
class Example {
    var name: String? = "Alice"
    
    fun printLength() {
        if (name != null) {
            println(name.length)  // COMPILE ERROR: Smart cast to 'String' is impossible
                                   // because 'name' is a mutable property
        }
    }
}
```

**Why It's a Pitfall:**
Kotlin's compiler doesn't trust that `name` (a `var` property) won't change between the check and the use — another thread or another call could set it to null.

**The Fix:**
```kotlin
fun printLength() {
    val localName = name  // copy to a local val
    if (localName != null) {
        println(localName.length)  // smart cast works on val
    }
}

// Or use let:
fun printLength() {
    name?.let { println(it.length) }
}

// Or use a local val with Elvis:
fun printLength() {
    val localName = name ?: return
    println(localName.length)
}
```

---

## Pitfall 5: Integer Division Surprise

**The Mistake:**
```kotlin
val percentage = 1 / 3 * 100
println(percentage)  // 0, not 33!
```

**Why It's a Pitfall:**
In Kotlin (and Java), division of two integers is **integer division** — the fractional part is truncated. `1 / 3` is `0`, and `0 * 100 = 0`.

**The Fix:**
```kotlin
// Use Double for decimal division
val percentage = 1.0 / 3.0 * 100  // 33.33...
val percentage2 = 1.toDouble() / 3 * 100
val percentage3 = 1.0 / 3 * 100

println("%.1f%%".format(percentage))  // 33.3%
```

---

## Pitfall 6: Modifying a Collection While Iterating

**The Mistake:**
```kotlin
val list = mutableListOf(1, 2, 3, 4, 5)

// ConcurrentModificationException!
for (item in list) {
    if (item % 2 == 0) {
        list.remove(item)
    }
}
```

**Why It's a Pitfall:**
You can't modify a collection while a for-loop is iterating over it (using the underlying iterator). This throws `ConcurrentModificationException`.

**The Fix:**
```kotlin
// Option 1: filter into a new list (preferred)
val filtered = list.filter { it % 2 != 0 }
println(filtered)  // [1, 3, 5]

// Option 2: use iterator's remove
val iterator = list.iterator()
while (iterator.hasNext()) {
    if (iterator.next() % 2 == 0) {
        iterator.remove()  // safe removal through iterator
    }
}

// Option 3: removeAll
list.removeAll { it % 2 == 0 }
```

---

## Pitfall 7: Nullable in Generic Types — Unexpected Behavior

**The Mistake:**
```kotlin
fun <T> printFirstOrNull(list: List<T>) {
    val first = list.firstOrNull()
    println(first ?: "empty")  // works for empty list
    
    // But what if T is already nullable?
    val nullableList = listOf<String?>(null, "hello")
    println(nullableList.firstOrNull() ?: "empty")  // prints "empty" — but null IS a valid element!
}
```

**Why It's a Pitfall:**
`firstOrNull()` returns `null` in two cases: empty list, OR the first element is `null`. You can't distinguish between them with `?: "default"`.

**The Fix:**
```kotlin
// Use indices or isEmpty checks
val list = listOf<String?>(null, "hello")
if (list.isEmpty()) {
    println("List is empty")
} else {
    println("First element: ${list[0]}")  // explicit index access
}

// Or use indexOfFirst and check
val hasElements = list.isNotEmpty()
```

---

## Pitfall 8: Forgetting That String Templates Call toString()

**The Mistake:**
```kotlin
data class User(val name: String, val role: String)

val user: User? = null
println("User: $user")  // User: null — not an NPE, but "null" string
println("User: ${user?.name}")  // User: null — also "null" string
```

**Why It's a Pitfall:**
String templates call `.toString()` on the value. For null, this returns the string `"null"` — not an exception. This can mask nullability issues in output.

**The Fix:**
```kotlin
// Be explicit about null handling in output
println("User: ${user?.name ?: "no user"}")
println("User: ${user ?: "not logged in"}")
```

---

## Pitfall 9: equals() and hashCode() with Mutable Data Classes

**The Mistake:**
```kotlin
data class MutableUser(var name: String, var age: Int)

val user = MutableUser("Alice", 30)
val set = mutableSetOf(user)

println(user in set)  // true

user.name = "Bob"  // mutate after adding to set!
println(user in set)  // false! — hash changed, can't find it
println(set.size)  // 1 — it's still in there, just unfindable
```

**Why It's a Pitfall:**
Data classes generate `equals()` and `hashCode()` based on properties. If you mutate a data class that's in a `Set` or `Map`, its hash code changes, and the collection can't find it anymore.

**The Fix:**
```kotlin
// Option 1: Don't use mutable properties in data classes
data class User(val name: String, val age: Int)  // all val

// Option 2: If you need mutation, don't put the object in sets/maps
// Use a stable ID instead:
data class MutableUser(val id: UUID, var name: String, var age: Int)
val userMap = mutableMapOf<UUID, MutableUser>()
// Key is the stable ID, not the mutable object itself
```

---

## Pitfall 10: Lazy Initialization in Multi-Threaded Contexts

**The Mistake:**
```kotlin
class Config {
    val settings: Map<String, String> by lazy(LazyThreadSafetyMode.NONE) {
        loadSettingsFromFile()  // expensive but called once
    }
}

// Used from multiple threads simultaneously:
val config = Config()
thread1 { println(config.settings["key1"]) }  // might initialize
thread2 { println(config.settings["key2"]) }  // might also initialize!
// Race condition: loadSettingsFromFile() might be called twice
```

**Why It's a Pitfall:**
`LazyThreadSafetyMode.NONE` is faster but NOT thread-safe. If accessed concurrently, the initializer can run multiple times.

**The Fix:**
```kotlin
// Default lazy is thread-safe
val settings: Map<String, String> by lazy {  // SYNCHRONIZED by default
    loadSettingsFromFile()
}

// Or explicitly:
val settings: Map<String, String> by lazy(LazyThreadSafetyMode.SYNCHRONIZED) {
    loadSettingsFromFile()
}
```

---

## Pitfall 11: Capturing Variables in Lambdas (Closures)

**The Mistake:**
```kotlin
val actions = mutableListOf<() -> Unit>()
var x = 0

for (i in 0..4) {
    x = i
    actions.add { println(x) }  // captures x, not i!
}

actions.forEach { it() }  // prints 4, 4, 4, 4, 4 — not 0, 1, 2, 3, 4
```

**Why It's a Pitfall:**
The lambda captures the **variable** `x`, not its value at the time the lambda was created. All lambdas share the same `x`, which ends up as 4.

**The Fix:**
```kotlin
// Option 1: Use the loop variable directly (it's a new val each iteration)
val actions = mutableListOf<() -> Unit>()
for (i in 0..4) {
    actions.add { println(i) }  // captures 'i' — a new val each iteration
}
actions.forEach { it() }  // 0, 1, 2, 3, 4 ✓

// Option 2: Copy the mutable variable to a local val
val actions2 = mutableListOf<() -> Unit>()
var x = 0
for (i in 0..4) {
    x = i
    val captured = x  // local val — each iteration gets its own
    actions2.add { println(captured) }
}
```

---

## Pitfall 12: == vs === with Nullable Wrappers

**The Mistake:**
```kotlin
val a: Int = 1000
val b: Int = 1000
println(a == b)   // true — expected
println(a === b)  // true — OK (primitives compared by value)

val x: Int? = 1000
val y: Int? = 1000
println(x == y)   // true — structural equality
println(x === y)  // FALSE in some JVMs! (Integer cache only covers -128 to 127)
```

**Why It's a Pitfall:**
`Int?` compiles to Java's `Integer` (boxed). The JVM caches `Integer` objects for values -128 to 127 — outside this range, `===` (reference equality) may return false even for equal values.

**The Fix:**
Always use `==` (structural equality) for comparing values:
```kotlin
val x: Int? = 1000
val y: Int? = 1000
println(x == y)  // true — always correct for value comparison
// Never use === for value comparison — only for reference identity checks
```

---

## Pitfall 13: Not Handling Coroutine Cancellation

**The Mistake:**
```kotlin
suspend fun computeHeavy(): Int {
    var result = 0
    for (i in 1..1_000_000) {
        result += i  // pure computation, no suspension points
        // coroutine cancellation is IGNORED here
    }
    return result
}

// Cancelling this coroutine has no effect until it finishes
val job = launch { computeHeavy() }
delay(100)
job.cancel()  // request cancellation
job.join()    // still waits the full time...
```

**Why It's a Pitfall:**
Coroutines cooperatively check for cancellation at suspension points. If your code never suspends, it never checks for cancellation.

**The Fix:**
```kotlin
suspend fun computeHeavy(): Int {
    var result = 0
    for (i in 1..1_000_000) {
        if (!isActive) return result  // check cancellation
        // or: yield()  — checks and suspends briefly
        result += i
    }
    return result
}

// Or use ensureActive() for cleaner checking:
suspend fun computeHeavy(): Int {
    var result = 0
    for (i in 1..1_000_000) {
        ensureActive()  // throws CancellationException if cancelled
        result += i
    }
    return result
}
```

---

## Pitfall 14: Wrong Collection Mutability Type

**The Mistake:**
```kotlin
// This creates a MUTABLE list but returns it as immutable
fun getItems(): List<String> {
    val items = mutableListOf("a", "b", "c")
    return items  // returned as List, but caller might cast it back!
}

val items = getItems()
(items as MutableList<String>).add("d")  // works! unsafe downcast
```

**Why It's a Pitfall:**
Returning a `MutableList` as a `List` doesn't protect it — callers can downcast. This breaks encapsulation.

**The Fix:**
```kotlin
// Return a truly non-mutable copy
fun getItems(): List<String> {
    val mutableItems = mutableListOf("a", "b", "c")
    return mutableItems.toList()  // creates an immutable copy
}

// Or use toUnmodifiableList() for Java interop:
import java.util.Collections
fun getItems(): List<String> = Collections.unmodifiableList(mutableListOf("a", "b", "c"))
```

---

## Pitfall 15: Operator Overloading Abuse

**The Mistake:**
```kotlin
data class Color(val r: Int, val g: Int, val b: Int) {
    operator fun plus(other: Color) = Color(r + other.r, g + other.g, b + other.b)
    operator fun times(factor: Double) = Color(
        (r * factor).toInt(),
        (g * factor).toInt(),
        (b * factor).toInt()
    )
    operator fun minus(other: Color) = Color(r - other.r, g - other.g, b - other.b)
    operator fun div(divisor: Int) = Color(r / divisor, g / divisor, b / divisor)
    operator fun rem(other: Color) = Color(r % other.r, g % other.g, b % other.b)
    operator fun rangeTo(other: Color) = TODO("what does a color range even mean?")
}

val mixed = (Color(255, 0, 0) + Color(0, 255, 0)) / 2  // confusing!
```

**Why It's a Pitfall:**
Operator overloading is powerful but can make code confusing when the operators don't intuitively map to the operation (e.g., what does `%` mean for colors?).

**The Fix:**
Only overload operators when the meaning is **obvious and natural**:
```kotlin
data class Color(val r: Int, val g: Int, val b: Int) {
    operator fun plus(other: Color) = Color(r + other.r, g + other.g, b + other.b)
    // ^ "mixing" colors by addition is arguably intuitive
    
    // Better: use named functions for non-obvious operations
    fun blend(other: Color, ratio: Double): Color = Color(
        (r * (1 - ratio) + other.r * ratio).toInt(),
        (g * (1 - ratio) + other.g * ratio).toInt(),
        (b * (1 - ratio) + other.b * ratio).toInt()
    )
}

val blended = Color(255, 0, 0).blend(Color(0, 255, 0), 0.5)  // clear intent
```

---

## Quick Reference: Pitfall Checklist

| # | Pitfall | Quick Fix |
|---|---------|-----------|
| 1 | `val` ≠ immutable object | Use `listOf()` not `mutableListOf()` |
| 2 | `!!` overuse | Use `?.`, `?:`, `checkNotNull` |
| 3 | Shared mutable state in coroutines | Use `AtomicXxx` or `Mutex` |
| 4 | Smart cast fails on `var` | Copy to local `val` first |
| 5 | Integer division truncation | Use `1.0 / 3` not `1 / 3` |
| 6 | ConcurrentModificationException | Use `filter` or iterator |
| 7 | Nullable elements in `firstOrNull` | Check `isEmpty()` separately |
| 8 | `null` prints as "null" | Use `?: "default"` in templates |
| 9 | Mutating data class in set/map | Use `val` properties in data classes |
| 10 | Unsafe lazy in multithreading | Use default `lazy` (synchronized) |
| 11 | Capturing var in lambda | Use `val` copy or loop variable |
| 12 | `===` with boxed integers | Always use `==` for value comparison |
| 13 | Ignoring coroutine cancellation | Check `isActive` or call `ensureActive()` |
| 14 | Leaking mutable collection | Return `toList()` copy |
| 15 | Operator overloading abuse | Only overload when intuitively obvious |

---

*Next: [Interview Preparation](16-interview-prep.md)*
