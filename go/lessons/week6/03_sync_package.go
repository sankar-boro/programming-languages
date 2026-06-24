/*
WEEK 6 — DAY 3: The sync Package — Mutexes, WaitGroups, and More
==================================================================
Topic: Go's synchronization primitives for concurrent programming.

Key ideas:
  - sync.Mutex protects shared data from concurrent access
  - sync.RWMutex allows multiple concurrent readers or one writer
  - sync.WaitGroup waits for a group of goroutines to finish
  - sync.Once ensures a function runs exactly once
  - sync.Map is a concurrent-safe map for specific use cases
  - atomic operations for lock-free counters
*/

package main

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

// ─── 1. THE DATA RACE PROBLEM ─────────────────────────────────────────────────
//
// A data race occurs when two goroutines access the same memory location
// concurrently and at least one of them writes.
//
// Data races are UNDEFINED BEHAVIOR — they may:
//   - Give wrong results
//   - Cause crashes
//   - Produce different results on every run
//   - Pass tests and fail in production
//
// Detect with: go test -race ./... or go run -race main.go

func dataRaceProblem() {
	fmt.Println("=== 1. Data Race Problem ===")

	// UNSAFE: data race (don't do this!)
	// var counter int
	// for i := 0; i < 1000; i++ {
	//     go func() { counter++ }()  // RACE: concurrent read-modify-write
	// }

	// Detection: go run -race will report:
	// WARNING: DATA RACE
	// Write at 0x... by goroutine N:
	// Previous write at 0x... by goroutine M:
	fmt.Println("Run 'go run -race' to detect data races")
	fmt.Println("Use sync.Mutex or atomic ops to fix them")
}

// ─── 2. sync.Mutex ────────────────────────────────────────────────────────────
//
// A mutex (mutual exclusion lock) ensures that only ONE goroutine can
// execute the critical section at a time.
//
// Lock()   — acquire the lock (blocks if already held by another goroutine)
// Unlock() — release the lock
//
// ALWAYS use `defer mu.Unlock()` immediately after `mu.Lock()` to prevent
// deadlocks if the function returns early or panics.

type SafeCounter struct {
	mu    sync.Mutex
	count int
}

func (c *SafeCounter) Increment() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.count++
}

func (c *SafeCounter) Add(n int) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.count += n
}

func (c *SafeCounter) Value() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.count
}

func (c *SafeCounter) Reset() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	v := c.count
	c.count = 0
	return v
}

func mutex() {
	fmt.Println("\n=== 2. sync.Mutex ===")

	counter := &SafeCounter{}
	var wg sync.WaitGroup

	// 1000 goroutines increment concurrently — all safe
	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			counter.Increment()
		}()
	}

	wg.Wait()
	fmt.Printf("Final count: %d (expected 1000)\n", counter.Value())

	// Demonstrate deadlock avoidance with defer
	fmt.Println("Using defer mu.Unlock() prevents deadlocks on early returns")
}

// ─── 3. sync.RWMutex ──────────────────────────────────────────────────────────
//
// RWMutex distinguishes between reads and writes:
//
//   RLock()   — acquire read lock (multiple readers can hold simultaneously)
//   RUnlock() — release read lock
//   Lock()    — acquire write lock (EXCLUSIVE — blocks all readers and writers)
//   Unlock()  — release write lock
//
// Use RWMutex when:
//   - Reads are much more frequent than writes
//   - Reads don't modify shared data
//
// A writer must wait for ALL readers to finish (and vice versa).

type Cache struct {
	mu    sync.RWMutex
	items map[string]string
}

func NewCache() *Cache {
	return &Cache{items: make(map[string]string)}
}

func (c *Cache) Get(key string) (string, bool) {
	c.mu.RLock()          // multiple goroutines can read simultaneously
	defer c.mu.RUnlock()
	v, ok := c.items[key]
	return v, ok
}

func (c *Cache) Set(key, value string) {
	c.mu.Lock()          // exclusive — no other goroutines can read or write
	defer c.mu.Unlock()
	c.items[key] = value
}

func (c *Cache) Delete(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.items, key)
}

func (c *Cache) Keys() []string {
	c.mu.RLock()
	defer c.mu.RUnlock()
	keys := make([]string, 0, len(c.items))
	for k := range c.items {
		keys = append(keys, k)
	}
	return keys
}

func rwMutex() {
	fmt.Println("\n=== 3. sync.RWMutex ===")

	cache := NewCache()

	// Many concurrent readers, occasional writer
	var wg sync.WaitGroup

	// Writers
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			key := fmt.Sprintf("key%d", i)
			cache.Set(key, fmt.Sprintf("value%d", i))
			fmt.Printf("  wrote %s\n", key)
		}(i)
	}

	// Readers (can run concurrently with each other)
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			time.Sleep(time.Duration(i) * time.Millisecond)
			key := fmt.Sprintf("key%d", i%5)
			v, ok := cache.Get(key)
			if ok {
				fmt.Printf("  read %s=%s\n", key, v)
			}
		}(i)
	}

	wg.Wait()
	fmt.Println("Keys:", cache.Keys())
}

// ─── 4. sync.WaitGroup ────────────────────────────────────────────────────────
//
// WaitGroup waits for a collection of goroutines to finish.
//
// wg.Add(n)  — increment the counter by n (BEFORE starting goroutines)
// wg.Done()  — decrement the counter by 1 (call from goroutine)
// wg.Wait()  — block until counter reaches 0
//
// Rules:
//   - wg.Add() must be called BEFORE the goroutine is started
//   - wg.Done() is typically deferred at the start of the goroutine
//   - Never call wg.Add() from inside the goroutine (race with wg.Wait())

func waitGroup() {
	fmt.Println("\n=== 4. sync.WaitGroup ===")

	var wg sync.WaitGroup
	results := make([]int, 5)

	for i := 0; i < 5; i++ {
		wg.Add(1)           // increment BEFORE launching goroutine
		go func(i int) {
			defer wg.Done() // guaranteed to be called
			time.Sleep(time.Duration(5-i) * 20 * time.Millisecond)
			results[i] = i * i
			fmt.Printf("  goroutine %d computed %d\n", i, results[i])
		}(i)
	}

	wg.Wait()  // blocks until all Done() calls
	fmt.Println("All results:", results)

	// Pattern: semaphore (limit concurrent goroutines)
	const maxConcurrent = 3
	sem := make(chan struct{}, maxConcurrent)
	var wg2 sync.WaitGroup

	for i := 0; i < 10; i++ {
		wg2.Add(1)
		sem <- struct{}{} // acquire slot (blocks if maxConcurrent active)
		go func(i int) {
			defer wg2.Done()
			defer func() { <-sem }() // release slot

			fmt.Printf("  processing %d\n", i)
			time.Sleep(50 * time.Millisecond)
		}(i)
	}
	wg2.Wait()
}

// ─── 5. sync.Once ─────────────────────────────────────────────────────────────
//
// sync.Once ensures a function is called exactly once, regardless of
// how many goroutines call it concurrently. Perfect for lazy initialization.

type Database struct {
	once sync.Once
	conn string  // simulated connection
}

func (db *Database) connect() {
	db.once.Do(func() {
		fmt.Println("  [connecting to database...]")
		time.Sleep(100 * time.Millisecond) // simulate connection setup
		db.conn = "connected"
		fmt.Println("  [database connected]")
	})
}

func (db *Database) Query(q string) string {
	db.connect()  // safe to call from any goroutine — connect runs only once
	return fmt.Sprintf("result of %q on %s", q, db.conn)
}

func syncOnce() {
	fmt.Println("\n=== 5. sync.Once ===")

	db := &Database{}
	var wg sync.WaitGroup

	// 5 goroutines all try to query — connect() runs exactly once
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			result := db.Query(fmt.Sprintf("SELECT * FROM table%d", i))
			fmt.Printf("  goroutine %d: %s\n", i, result)
		}(i)
	}

	wg.Wait()
}

// ─── 6. atomic OPERATIONS ─────────────────────────────────────────────────────
//
// For simple counter operations, atomic is faster than mutex.
// atomic operations are guaranteed to be indivisible (atomic).
// No other goroutine can observe a partial update.
//
// sync/atomic provides: Add, Load, Store, Swap, CompareAndSwap
// for int32, int64, uint32, uint64, uintptr, and unsafe.Pointer.
//
// Use when: the operation is a single read/write/increment
// Don't use for: compound operations (check-then-act needs mutex)

type AtomicCounter struct {
	n int64  // must use int64 (or int32) for atomic ops
}

func (c *AtomicCounter) Increment() int64 {
	return atomic.AddInt64(&c.n, 1)  // returns new value
}

func (c *AtomicCounter) Value() int64 {
	return atomic.LoadInt64(&c.n)    // safe concurrent read
}

func (c *AtomicCounter) Reset() {
	atomic.StoreInt64(&c.n, 0)
}

func atomicOps() {
	fmt.Println("\n=== 6. Atomic Operations ===")

	counter := &AtomicCounter{}
	var wg sync.WaitGroup

	for i := 0; i < 10000; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			counter.Increment()
		}()
	}

	wg.Wait()
	fmt.Printf("Atomic counter: %d (expected 10000)\n", counter.Value())

	// Compare-and-swap: update only if current value matches expected
	var val int64 = 0
	swapped := atomic.CompareAndSwapInt64(&val, 0, 42)  // if val==0, set to 42
	fmt.Printf("CAS(0→42): swapped=%v val=%d\n", swapped, val)

	swapped = atomic.CompareAndSwapInt64(&val, 0, 100)  // if val==0... but it's 42
	fmt.Printf("CAS(0→100): swapped=%v val=%d\n", swapped, val)  // not swapped
}

// ─── 7. sync.Map ──────────────────────────────────────────────────────────────
//
// sync.Map is a concurrent-safe map optimized for specific patterns:
//   - Many goroutines reading the same keys
//   - Each goroutine writes to a different set of keys
// It does NOT require a mutex — uses a different internal synchronization.
//
// For general concurrent map use, a mutex-protected map is often simpler and
// sometimes faster. Use sync.Map only when profiling shows benefit.

func syncMap() {
	fmt.Println("\n=== 7. sync.Map ===")

	var m sync.Map
	var wg sync.WaitGroup

	// Many writers writing to different keys
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			m.Store(fmt.Sprintf("key%d", i), i*10)
		}(i)
	}

	// Many readers
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			key := fmt.Sprintf("key%d", i%10)
			if v, ok := m.Load(key); ok {
				fmt.Printf("  loaded %s=%v\n", key, v)
			}
		}(i)
	}

	wg.Wait()

	// Range over all entries
	fmt.Println("All entries:")
	m.Range(func(key, value any) bool {
		fmt.Printf("  %v=%v\n", key, value)
		return true  // return false to stop
	})

	// Delete
	m.Delete("key0")
	// LoadOrStore: load if exists, store if not
	actual, loaded := m.LoadOrStore("key99", "new")
	fmt.Printf("LoadOrStore key99: actual=%v loaded=%v\n", actual, loaded)
}

func main() {
	dataRaceProblem()
	mutex()
	rwMutex()
	waitGroup()
	syncOnce()
	atomicOps()
	syncMap()
}

/*
THOUGHT QUESTIONS:

1. What is a data race? Why is it undefined behavior rather than just "a bug
   with a deterministic wrong answer"?

2. When would you choose sync.RWMutex over sync.Mutex?
   What workload characteristic makes RWMutex beneficial?

3. Why must wg.Add(n) be called BEFORE launching goroutines?
   What race condition would occur if you called it inside the goroutine?

4. sync.Once guarantees a function runs exactly once even with concurrent callers.
   How might this be implemented internally?

5. When is atomic.AddInt64 preferable to sync.Mutex? When is it NOT sufficient?

EXERCISES:

1. Implement a concurrent rate limiter using a Mutex and time.Ticker:
   Allow at most N requests per second. Return false when rate exceeded.

2. Write a concurrent `Pool[T any]` that manages a fixed-size pool of
   reusable objects, similar to sync.Pool but type-safe.

3. Implement a concurrent broadcast system: one publisher sends messages
   to N subscribers. Use channels and sync primitives.

4. Write a lock-free stack using atomic.CompareAndSwapPointer.
   Benchmark it against a mutex-based stack.
*/
