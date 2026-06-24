/*
WEEK 5 — DAY 3: Goroutines and Channels
=========================================
Topic: Go's concurrency model — goroutines, channels, select, and common patterns.

Key ideas:
  - Goroutines are lightweight cooperatively-scheduled coroutines
  - Channels are the primary communication mechanism (share memory by communicating)
  - Buffered vs unbuffered channels have different blocking semantics
  - select handles multiple channel operations
  - The GMP scheduler manages goroutines on OS threads
*/

package main

import (
	"fmt"
	"sync"
	"time"
)

// ─── 1. GOROUTINES ─────────────────────────────────────────────────────────────
//
// A goroutine is a lightweight thread managed by the Go runtime.
// `go fn()` starts fn as a new goroutine.
//
// Goroutines vs OS threads:
//   - OS thread: ~1–8 MB stack, kernel-managed, expensive context switch
//   - Goroutine: ~2–8 KB stack (grows), user-space scheduled, very cheap
//   - Go can run millions of goroutines on hundreds of OS threads
//
// The Go scheduler uses GMP:
//   G = Goroutine  (the goroutine itself)
//   M = Machine    (OS thread)
//   P = Processor  (logical CPU — holds the run queue)
//
// When a goroutine blocks (I/O, channel, mutex), the P can run other Gs.
// GOMAXPROCS controls how many Ps (and thus Ms running Go code) exist.

func goroutines() {
	fmt.Println("=== 1. Goroutines ===")

	// Launch goroutines
	var wg sync.WaitGroup

	for i := 1; i <= 5; i++ {
		wg.Add(1)
		go func(n int) {
			defer wg.Done()
			fmt.Printf("  goroutine %d running\n", n)
			time.Sleep(time.Duration(n) * 10 * time.Millisecond)
			fmt.Printf("  goroutine %d done\n", n)
		}(i)  // pass i to avoid loop variable capture
	}

	wg.Wait()
	fmt.Println("all goroutines finished")
}

// ─── 2. CHANNELS — COMMUNICATION PRIMITIVES ────────────────────────────────────
//
// A channel is a FIFO queue with optional buffering.
//
// Unbuffered channel (make(chan T)):
//   - Send BLOCKS until a receiver is ready
//   - Receive BLOCKS until a sender is ready
//   - Synchronization point: sender and receiver rendezvous
//
// Buffered channel (make(chan T, N)):
//   - Send BLOCKS only when the buffer is FULL
//   - Receive BLOCKS only when the buffer is EMPTY
//   - N items can be queued without blocking
//
// Channel operations:
//   ch <- val   send
//   val = <-ch  receive
//   val, ok := <-ch  receive with close check
//   close(ch)   signal no more sends; receivers drain remaining values

func channels() {
	fmt.Println("\n=== 2. Channels ===")

	// Unbuffered channel — rendezvous
	ch := make(chan string)

	go func() {
		time.Sleep(50 * time.Millisecond)
		ch <- "hello from goroutine"  // blocks until main receives
	}()

	msg := <-ch  // blocks until the goroutine sends
	fmt.Println("received:", msg)

	// Buffered channel
	buf := make(chan int, 3)
	buf <- 1  // doesn't block — buffer has space
	buf <- 2
	buf <- 3
	// buf <- 4  // would block — buffer is full

	fmt.Println("buffered:", <-buf, <-buf, <-buf)

	// Channel direction in function signatures
	// chan<- T : send-only channel
	// <-chan T : receive-only channel
	// This documents intent and prevents accidental misuse
	producer := func(out chan<- int) {
		for i := 0; i < 5; i++ {
			out <- i
		}
		close(out)
	}

	nums := make(chan int)
	go producer(nums)

	for n := range nums {  // range over channel — exits when channel closed
		fmt.Printf("  received: %d\n", n)
	}
}

// ─── 3. CHANNEL PATTERNS ──────────────────────────────────────────────────────

// Pipeline: stages connected by channels
func generate(nums ...int) <-chan int {
	out := make(chan int)
	go func() {
		for _, n := range nums {
			out <- n
		}
		close(out)
	}()
	return out
}

func square(in <-chan int) <-chan int {
	out := make(chan int)
	go func() {
		for n := range in {
			out <- n * n
		}
		close(out)
	}()
	return out
}

func filter(in <-chan int, predicate func(int) bool) <-chan int {
	out := make(chan int)
	go func() {
		for n := range in {
			if predicate(n) {
				out <- n
			}
		}
		close(out)
	}()
	return out
}

func pipeline() {
	fmt.Println("\n=== 3. Pipeline Pattern ===")

	// Pipeline: generate → square → filter (even) → print
	nums := generate(1, 2, 3, 4, 5, 6, 7, 8, 9, 10)
	squares := square(nums)
	evens := filter(squares, func(n int) bool { return n%2 == 0 })

	for n := range evens {
		fmt.Printf("  %d\n", n)
	}
}

// Fan-out: send work to multiple goroutines
// Fan-in: merge results from multiple channels

func fanOutFanIn() {
	fmt.Println("\n=== 3b. Fan-Out / Fan-In ===")

	work := make(chan int, 10)
	results := make(chan int, 10)

	// Fan-out: 3 workers
	var wg sync.WaitGroup
	for i := 0; i < 3; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for n := range work {
				// simulate work
				result := n * n
				fmt.Printf("  worker %d: %d² = %d\n", id, n, result)
				results <- result
			}
		}(i)
	}

	// Send work
	for i := 1; i <= 9; i++ {
		work <- i
	}
	close(work)

	// Wait for all workers, then close results
	go func() {
		wg.Wait()
		close(results)
	}()

	// Collect results (fan-in)
	total := 0
	for r := range results {
		total += r
	}
	fmt.Println("Total:", total)
}

// ─── 4. SELECT — MULTIPLEXING CHANNELS ────────────────────────────────────────
//
// select is like a switch but for channel operations.
// It blocks until ONE of its cases can proceed, then executes that case.
// If multiple cases are ready, one is chosen at RANDOM.
// A default case makes select non-blocking.

func selectStatement() {
	fmt.Println("\n=== 4. Select ===")

	ch1 := make(chan string, 1)
	ch2 := make(chan string, 1)

	go func() {
		time.Sleep(20 * time.Millisecond)
		ch1 <- "from ch1"
	}()

	go func() {
		time.Sleep(10 * time.Millisecond)
		ch2 <- "from ch2"
	}()

	// Wait for whichever comes first
	select {
	case msg1 := <-ch1:
		fmt.Println("received ch1:", msg1)
	case msg2 := <-ch2:
		fmt.Println("received ch2:", msg2)  // ch2 is faster
	}

	// Non-blocking check with default
	select {
	case msg := <-ch1:
		fmt.Println("got from ch1:", msg)
	default:
		fmt.Println("ch1 not ready (non-blocking)")
	}
}

// Timeout pattern
func withTimeout(ch <-chan string, timeout time.Duration) (string, bool) {
	select {
	case val := <-ch:
		return val, true
	case <-time.After(timeout):
		return "", false
	}
}

// Done/cancellation pattern
func doWorkWithCancel(done <-chan struct{}) {
	for {
		select {
		case <-done:
			fmt.Println("  cancelled!")
			return
		default:
			fmt.Println("  working...")
			time.Sleep(50 * time.Millisecond)
		}
	}
}

func selectPatterns() {
	fmt.Println("\n=== 4b. Select Patterns ===")

	// Timeout
	slowCh := make(chan string)
	go func() {
		time.Sleep(200 * time.Millisecond)
		slowCh <- "slow result"
	}()

	if val, ok := withTimeout(slowCh, 50*time.Millisecond); ok {
		fmt.Println("got:", val)
	} else {
		fmt.Println("timed out!")
	}

	// Cancellation
	done := make(chan struct{})
	go func() {
		time.Sleep(120 * time.Millisecond)
		close(done)  // signal cancellation
	}()
	doWorkWithCancel(done)
}

// ─── 5. DONE CHANNEL AND CLOSE SEMANTICS ──────────────────────────────────────
//
// Closing a channel:
//   - Sends a zero value to ALL receivers immediately
//   - Subsequent receives return (zero, false)
//   - It's a BROADCAST — all goroutines waiting on the channel are woken
//   - Only the SENDER should close a channel (closing from receiver is a bug)
//   - Closing a nil channel panics
//   - Closing an already-closed channel panics

func closeSemantics() {
	fmt.Println("\n=== 5. Close Semantics ===")

	ch := make(chan int, 5)
	ch <- 1; ch <- 2; ch <- 3
	close(ch)

	// Drain closed channel — range stops when channel is closed and empty
	for v := range ch {
		fmt.Printf("  drained: %d\n", v)
	}

	// After close, receive returns (zero, false)
	v, ok := <-ch
	fmt.Printf("after drain: v=%d ok=%v\n", v, ok)

	// Broadcast via close
	quit := make(chan struct{})
	var wg sync.WaitGroup
	for i := 0; i < 3; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			<-quit  // all goroutines wait here
			fmt.Printf("  goroutine %d received quit signal\n", id)
		}(i)
	}
	time.Sleep(10 * time.Millisecond)
	close(quit)  // broadcasts to ALL waiting goroutines simultaneously
	wg.Wait()
}

// ─── 6. COMMON CONCURRENCY MISTAKES ──────────────────────────────────────────

func commonMistakes() {
	fmt.Println("\n=== 6. Common Concurrency Mistakes ===")

	// MISTAKE 1: Goroutine leak — goroutine blocked on channel, nobody drains it
	// Goroutine would be stuck forever if nobody reads from leakedCh:
	// leakedCh := make(chan int)
	// go func() { leakedCh <- 1 }()
	// Don't run this without a receiver!

	// MISTAKE 2: Data race — two goroutines access shared data without sync
	// var counter int
	// go func() { counter++ }()
	// go func() { counter++ }()
	// This is a data race — run `go test -race` to detect it

	// MISTAKE 3: Closing a channel twice → panic
	// ch := make(chan int)
	// close(ch)
	// close(ch)  // PANIC: close of closed channel

	// CORRECT: Use sync.Mutex for shared counter
	var counter int
	var mu sync.Mutex
	var wg sync.WaitGroup

	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			mu.Lock()
			counter++
			mu.Unlock()
		}()
	}
	wg.Wait()
	fmt.Printf("safe counter: %d (expected 1000)\n", counter)
}

func main() {
	goroutines()
	channels()
	pipeline()
	fanOutFanIn()
	selectStatement()
	selectPatterns()
	closeSemantics()
	commonMistakes()
}

/*
THOUGHT QUESTIONS:

1. What is the GMP scheduler model? What does each letter stand for?

2. Unbuffered channels: "the sender blocks until a receiver is ready."
   What synchronization guarantee does this provide?

3. Why is `close(ch)` a broadcast to all receivers? How does this enable
   the "done channel" cancellation pattern?

4. What is a goroutine leak? How can you detect and prevent one?

5. "Don't communicate by sharing memory; share memory by communicating."
   What does this mean in practice?

EXERCISES:

1. Implement a worker pool: given a list of URLs to "fetch", distribute
   the work among N goroutines, collect all results, print them.

2. Write a `merge(cs ...<-chan int) <-chan int` that merges multiple
   input channels into a single output channel (fan-in).

3. Implement a rate-limited channel: wrap a channel so at most N messages
   per second are sent through, using `time.Ticker`.

4. Write a concurrent map with read and write safety using sync.RWMutex.
   Implement Get, Set, Delete, and Keys methods.
*/
