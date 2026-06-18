# Chapter 11 — Coroutines and Concurrency

> *"Coroutines are computations that can be suspended and resumed. They make asynchronous code look and feel like sequential code."*

---

## 11.1 The Problem with Threads

### Threads: The Traditional Approach

Traditional concurrency in Java and Kotlin uses threads. Threads are powerful but expensive:

```kotlin
// Creating threads directly — costly and unmanaged
fun fetchDataWithThreads() {
    val thread1 = Thread {
        // Simulate network request
        Thread.sleep(1000)
        println("Data from server 1")
    }
    
    val thread2 = Thread {
        Thread.sleep(500)
        println("Data from server 2")
    }
    
    thread1.start()
    thread2.start()
    thread1.join()
    thread2.join()
}
```

### The Problem

Threads have significant costs:
- **Memory**: Each thread needs a stack (~1MB by default)
- **Context switching**: OS-level switching is expensive
- **Scalability**: A server handling 10,000 requests needs 10,000 threads — infeasible
- **Complexity**: Shared state, synchronization, deadlocks, race conditions

Blocking a thread while waiting for I/O wastes the thread. If the thread is blocked, it can't do other work.

### The Callback Approach

One alternative is callbacks, but they lead to "callback hell":

```kotlin
// Callback-based code — hard to read, hard to maintain
fun fetchUser(id: String, callback: (User?, Error?) -> Unit) {
    // async call...
    callback(user, null)
}

fun fetchOrders(user: User, callback: (List<Order>?, Error?) -> Unit) {
    // async call...
    callback(orders, null)
}

// Deeply nested — "callback hell"
fetchUser("alice") { user, error ->
    if (error != null) { handle(error); return }
    fetchOrders(user!!) { orders, error ->
        if (error != null) { handle(error); return }
        orders?.forEach { order ->
            processOrder(order) { result, error ->
                // even deeper...
            }
        }
    }
}
```

---

## 11.2 What Is a Coroutine?

A **coroutine** is a lightweight concurrency primitive that can be **suspended** and **resumed** without blocking a thread.

Key concepts:
- Coroutines are **much cheaper than threads** — you can run millions of them
- Suspended coroutines don't block any thread — they're just paused
- Coroutines are **sequential in appearance** but concurrent in execution

```kotlin
// Conceptual model of coroutine suspension:

// Coroutine 1:
// --[running]--[suspended (waiting for network)]--[resumed]--[running]--
//                      |                              |
// Thread:  ------------|-----[does other work]--------|------------------

// The thread is freed during suspension and reused for other work
```

### Coroutines vs Threads

| Aspect | Threads | Coroutines |
|--------|---------|------------|
| Cost | ~1MB per thread | ~few bytes per coroutine |
| Creation speed | Slow (OS level) | Fast (library level) |
| Switching | Preemptive (OS) | Cooperative (coroutine) |
| Blocking | Blocks thread | Suspends coroutine |
| Scale | ~thousands | ~millions |
| Complexity | High | Low (looks sequential) |

---

## 11.3 suspend Functions

The `suspend` keyword marks a function as suspendable — it can be paused and resumed without blocking a thread.

```kotlin
suspend fun fetchUser(id: String): User {
    // This can suspend the coroutine while waiting
    delay(100)  // suspend for 100ms, but the thread is free
    return User(id = id, name = "Alice")
}

suspend fun fetchOrders(user: User): List<Order> {
    delay(200)  // another suspension point
    return listOf(Order("ORDER-1"), Order("ORDER-2"))
}
```

### Rules of suspend

1. `suspend` functions can only be called from **other suspend functions** or from a **coroutine**
2. Calling a suspend function is not automatically asynchronous — it runs sequentially unless you use coroutine builders

```kotlin
// This is WRONG — calling suspend function from non-suspend context:
fun main() {
    val user = fetchUser("alice")  // ERROR: suspend function can only be called from coroutine or another suspend function
}

// This is CORRECT — using a coroutine builder:
fun main() = runBlocking {
    val user = fetchUser("alice")  // OK — we're inside a coroutine
    println(user)
}
```

### Suspension Points

A function that calls `delay()` or other suspension functions has **suspension points** — places where the coroutine can pause:

```kotlin
suspend fun processAll(ids: List<String>) {
    for (id in ids) {
        val user = fetchUser(id)  // potential suspension point
        println("Processed: ${user.name}")
    }
}

// Even though each fetchUser() is called sequentially in the code,
// the thread is free during each suspension
```

---

## 11.4 Coroutine Builders

Coroutine builders are functions that create and start coroutines. The most common ones (from `kotlinx.coroutines`):

### runBlocking

Starts a coroutine and **blocks the current thread** until it completes. Used mainly in `main()` functions and tests:

```kotlin
import kotlinx.coroutines.*

fun main() = runBlocking {
    println("Start")
    delay(1000)  // coroutine suspends, but the thread is blocked (because of runBlocking)
    println("End")
}
```

### launch

Starts a new coroutine without blocking. Returns a `Job` (handle to the coroutine):

```kotlin
fun main() = runBlocking {
    println("Main starts")
    
    val job = launch {
        delay(1000)
        println("Coroutine finished")
    }
    
    println("Main continues")  // runs before the coroutine finishes
    job.join()  // wait for the coroutine to complete
    println("Main ends")
}
// Output:
// Main starts
// Main continues
// Coroutine finished
// Main ends
```

### async

Starts a coroutine and returns a `Deferred<T>` — a future/promise for a result:

```kotlin
fun main() = runBlocking {
    // Launch both computations concurrently
    val deferred1 = async {
        delay(1000)
        "Result 1"
    }
    
    val deferred2 = async {
        delay(500)
        "Result 2"
    }
    
    // Both run concurrently — total time ~1000ms, not 1500ms
    val result1 = deferred1.await()  // wait for result
    val result2 = deferred2.await()
    
    println("$result1 and $result2")  // Result 1 and Result 2
}
```

### Comparing launch and async

```kotlin
fun main() = runBlocking {
    // Sequential (naive approach)
    val time1 = measureTimeMillis {
        val user = fetchUser("alice")      // wait 1000ms
        val orders = fetchOrders(user)     // wait 500ms more
        println("Sequential: ${user.name}, ${orders.size} orders")
    }
    println("Time: ${time1}ms")  // ~1500ms
    
    // Concurrent with async
    val time2 = measureTimeMillis {
        val userDeferred = async { fetchUser("alice") }
        val statsDeferred = async { fetchUserStats("alice") }  // independent of user
        
        val user = userDeferred.await()
        val stats = statsDeferred.await()
        println("Concurrent: ${user.name}, ${stats.loginCount} logins")
    }
    println("Time: ${time2}ms")  // ~1000ms (both run in parallel)
}
```

---

## 11.5 Structured Concurrency

**Structured concurrency** is a key principle in Kotlin coroutines: coroutines are not orphaned — they have a well-defined lifecycle tied to a scope.

### CoroutineScope

Every coroutine runs inside a `CoroutineScope`. When the scope ends, all coroutines in it are cancelled:

```kotlin
fun main() = runBlocking {
    // This is a scope
    // All launched coroutines are children of this scope
    
    launch {
        delay(1000)
        println("Coroutine 1 done")
    }
    
    launch {
        delay(500)
        println("Coroutine 2 done")
    }
    
    // runBlocking waits for ALL children to complete before returning
}
// Coroutine 2 done
// Coroutine 1 done
```

### Why Structured Concurrency Matters

Without structured concurrency, a coroutine can become a "ghost" — running forever if its parent forgets about it:

```kotlin
// BAD — GlobalScope coroutine is unstructured
fun processSomething() {
    GlobalScope.launch {  // runs independently — might never be cancelled
        delay(1_000_000)
        println("Done?")  // might never run
    }
    // function returns, but coroutine keeps running!
}

// GOOD — structured: coroutine is tied to the caller's scope
suspend fun processSomething() = coroutineScope {
    launch {  // child of this scope
        delay(1000)
        println("Done")
    }
    // scope waits for all children
}
```

### Cancellation

Coroutines can be cancelled. Cancellation propagates to children:

```kotlin
fun main() = runBlocking {
    val job = launch {
        repeat(1000) { i ->
            println("Working... $i")
            delay(500)
        }
    }
    
    delay(2000)  // let it run for 2 seconds
    println("Cancelling...")
    job.cancel()  // request cancellation
    job.join()    // wait for cancellation to complete
    println("Cancelled")
}
// Working... 0
// Working... 1
// Working... 2
// Working... 3
// Cancelling...
// Cancelled
```

### Cancellation Is Cooperative

Suspension points (like `delay()`) check for cancellation. If your code never suspends, it won't respond to cancellation:

```kotlin
// This coroutine IGNORES cancellation — it never suspends
val job = launch {
    var i = 0
    while (true) {
        i++  // pure computation, no suspension points
    }
}
job.cancel()  // doesn't actually stop the loop!

// Fix: check isActive or use yield()
val job = launch {
    var i = 0
    while (isActive) {  // check cancellation status
        i++
    }
    println("Stopped at $i")
}
```

### Exception Handling in Coroutines

```kotlin
// launch: exceptions are stored in Job, propagate to parent
val job = launch {
    throw RuntimeException("Coroutine crashed!")
}
job.join()  // exception propagates here

// async: exceptions are wrapped in Deferred
val deferred = async {
    throw RuntimeException("Deferred crashed!")
}
try {
    deferred.await()  // exception thrown here
} catch (e: RuntimeException) {
    println("Caught: ${e.message}")
}

// supervisorScope: child failures don't cancel siblings
supervisorScope {
    val job1 = launch {
        throw RuntimeException("Job1 failed")
    }
    val job2 = launch {
        delay(500)
        println("Job2 completed")  // still runs even though Job1 failed
    }
}
```

---

## 11.6 Coroutine Context and Dispatchers

Every coroutine has a **context** that determines how and where it runs. The most important context element is the **dispatcher**.

### Dispatchers

A `Dispatcher` determines what thread(s) a coroutine uses:

```kotlin
import kotlinx.coroutines.*

fun main() = runBlocking {
    // Dispatchers.Main — UI thread (Android/Swing) — not available in plain JVM
    
    // Dispatchers.Default — CPU-intensive work (shared thread pool)
    launch(Dispatchers.Default) {
        println("Default: ${Thread.currentThread().name}")
        // Good for: sorting, parsing, heavy computation
    }
    
    // Dispatchers.IO — I/O operations (larger thread pool, can block)
    launch(Dispatchers.IO) {
        println("IO: ${Thread.currentThread().name}")
        // Good for: database, file, network operations
    }
    
    // Dispatchers.Unconfined — starts on current thread, resumes on suspension thread
    launch(Dispatchers.Unconfined) {
        println("Unconfined: ${Thread.currentThread().name}")
    }
}
```

### withContext — Switching Contexts

`withContext` switches the dispatcher for a block without creating a new coroutine:

```kotlin
suspend fun fetchAndProcess(): String {
    val rawData = withContext(Dispatchers.IO) {
        // Run I/O on IO dispatcher
        readFromFile("data.txt")  // blocking I/O, but thread from IO pool absorbs it
    }
    
    val processed = withContext(Dispatchers.Default) {
        // CPU work on Default dispatcher
        processData(rawData)
    }
    
    return processed
}
```

---

## 11.7 Flow: Cold Asynchronous Streams

A `Flow<T>` is a cold, asynchronous stream of values. Where `suspend` returns a single value, `Flow` returns multiple values over time.

### Basic Flow

```kotlin
import kotlinx.coroutines.*
import kotlinx.coroutines.flow.*

// A flow that emits 1, 2, 3 with delays
fun numbers(): Flow<Int> = flow {
    for (i in 1..3) {
        delay(100)   // simulate async operation
        emit(i)      // emit a value
    }
}

fun main() = runBlocking {
    numbers().collect { value ->
        println("Received: $value")
    }
}
// Received: 1
// Received: 2
// Received: 3
```

### Flow is Cold

Unlike Kotlin channels, a flow is **cold** — it doesn't start producing values until `collect` is called:

```kotlin
val flow = numbers()  // flow defined — nothing runs yet
println("Flow defined")

flow.collect { println("Got: $it") }  // NOW it starts
// Got: 1
// Got: 2
// Got: 3
```

### Flow Transformations

Flow has the same operators as collections, but they work asynchronously:

```kotlin
fun evenSquares(): Flow<Int> = flow {
    for (i in 1..10) {
        emit(i)
    }
}
.filter { it % 2 == 0 }
.map { it * it }

fun main() = runBlocking {
    evenSquares().collect { println(it) }
}
// 4
// 16
// 36
// 64
// 100
```

### Terminal Flow Operators

```kotlin
fun main() = runBlocking {
    val nums = (1..10).asFlow()
    
    println(nums.first())                // 1
    println(nums.last())                 // 10
    println(nums.count())                // 10
    println(nums.sum())                  // 55
    println(nums.toList())               // [1, 2, 3, ..., 10]
    println(nums.reduce { a, b -> a + b })  // 55
    
    // collect with index
    nums.filter { it % 2 == 0 }
        .withIndex()
        .collect { (index, value) ->
            println("$index: $value")
        }
    // 0: 2
    // 1: 4
    // 2: 6
    // ...
}
```

### Flow vs Sequence

| Aspect | Sequence | Flow |
|--------|----------|------|
| Blocking | Yes (thread is blocked) | No (coroutine suspends) |
| Async support | No | Yes |
| Context switching | No | Yes (flowOn) |
| Cancellation | No | Yes |
| Best for | CPU-bound lazy eval | Async, I/O, events |

### StateFlow and SharedFlow

These are special flows for sharing state and events between coroutines:

```kotlin
import kotlinx.coroutines.flow.*

// StateFlow — always has a current value, replays last value to new collectors
class ViewModel {
    private val _state = MutableStateFlow(0)
    val state: StateFlow<Int> = _state.asStateFlow()
    
    fun increment() {
        _state.value++
    }
}

// SharedFlow — can configure replay and buffer
val sharedFlow = MutableSharedFlow<String>(replay = 1)

fun main() = runBlocking {
    val vm = ViewModel()
    
    launch {
        vm.state.collect { value ->
            println("State: $value")
        }
    }
    
    launch {
        repeat(3) {
            delay(100)
            vm.increment()
        }
    }
}
```

---

## Practical Patterns with Coroutines

### Sequential vs Concurrent

```kotlin
suspend fun fetchUserData(userId: String): UserData {
    // Sequential — one after another
    val user = fetchUser(userId)        // wait
    val orders = fetchOrders(userId)    // then wait
    val stats = fetchStats(userId)      // then wait
    return UserData(user, orders, stats)
}

// Total time: ~3 seconds (sequential)

suspend fun fetchUserDataConcurrently(userId: String): UserData = coroutineScope {
    // Concurrent — all at once
    val user = async { fetchUser(userId) }
    val orders = async { fetchOrders(userId) }
    val stats = async { fetchStats(userId) }
    
    // Wait for all results
    UserData(user.await(), orders.await(), stats.await())
}

// Total time: ~1 second (limited by slowest request)
```

### Timeout

```kotlin
suspend fun safeFetch(): String? {
    return withTimeoutOrNull(1000L) {  // 1 second timeout
        fetchFromServer()  // might take too long
    }
    // returns null if timeout occurs
}

suspend fun mustFetch(): String {
    return withTimeout(1000L) {  // throws TimeoutCancellationException
        fetchFromServer()
    }
}
```

### Retry with Coroutines

```kotlin
suspend fun <T> retry(
    times: Int,
    delay: Long = 1000L,
    block: suspend () -> T
): T {
    var lastException: Exception? = null
    repeat(times) { attempt ->
        try {
            return block()
        } catch (e: Exception) {
            lastException = e
            println("Attempt ${attempt + 1} failed: ${e.message}")
            delay(delay)  // wait before retry (non-blocking)
        }
    }
    throw lastException!!
}

fun main() = runBlocking {
    val result = retry(times = 3, delay = 500L) {
        unstableNetworkCall()
    }
    println("Result: $result")
}
```

---

## Summary

Coroutines are lightweight concurrency primitives that suspend rather than block. The `suspend` modifier marks functions that can be suspended. Coroutine builders (`launch`, `async`, `runBlocking`) start coroutines in different modes. Structured concurrency ensures coroutines have well-defined lifetimes tied to their scope. Dispatchers control which thread(s) a coroutine uses. `Flow<T>` is a cold asynchronous stream for emitting multiple values over time.

---

## Key Takeaways

- `suspend` marks a function as a potential suspension point — it must be called from a coroutine
- `launch` fires-and-forgets; `async` returns a `Deferred<T>` for a result
- Structured concurrency: every coroutine belongs to a scope; scope cancels all children on failure
- Cancellation is cooperative — check `isActive` or use suspension points
- `withContext(Dispatcher.X)` switches the execution context without creating a new coroutine
- `Flow<T>` is cold — nothing executes until `collect` is called
- `StateFlow` holds state with replay; `SharedFlow` is a general event broadcast mechanism
- Use `async { }.await()` for concurrent operations; sequential calls are just regular suspend calls

---

## Practice Questions

### Conceptual
1. What is the difference between a thread and a coroutine?
2. Why can a `suspend` function only be called from another `suspend` function or a coroutine?
3. What is the difference between `launch` and `async`?
4. What is structured concurrency and why does it matter?
5. What makes `Flow` different from `List`? What makes it different from `Sequence`?

### Code Exercises

**Exercise 1:** Write a `suspend fun` that fetches three pieces of data (simulate with `delay()`). Implement it:
- Version A: sequentially (all three sequential)
- Version B: concurrently (all three in parallel with `async`)
Use `measureTimeMillis` to compare timings.

**Exercise 2:** Create a flow that emits the Fibonacci sequence indefinitely. Collect only the first 10 values that are even.

**Exercise 3:** Implement a `withRetry(times, delay, block)` suspend function that retries `block` on exception.

**Exercise 4:** Write a producer-consumer pattern using `Flow`:
- Producer emits numbers 1..100 with a 10ms delay each
- Consumer processes only prime numbers

**Exercise 5:** Explain the difference in behavior:
```kotlin
// Version A
launch {
    val a = async { task1() }
    val b = async { task2() }
    process(a.await(), b.await())
}

// Version B
launch {
    val a = task1()
    val b = task2()
    process(a, b)
}
```

---

*Next: [Chapter 12 — Interoperability](12-interoperability.md)*
