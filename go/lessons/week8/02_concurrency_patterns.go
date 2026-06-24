/*
WEEK 8 — DAY 2: Advanced Concurrency Patterns
==============================================
Topic: Production-grade concurrency — context, errgroup, worker pools,
       rate limiting, circuit breakers, and channel patterns.

Key ideas:
  - context.Context is the standard way to propagate cancellation and deadlines
  - errgroup simplifies managing groups of goroutines that return errors
  - Worker pools bound concurrency and control resource usage
  - Rate limiting and circuit breakers protect downstream systems
  - The "share memory by communicating" philosophy: channels for ownership transfer
*/

package main

import (
	"context"
	"fmt"
	"math/rand"
	"sync"
	"time"
)

// ─── 1. context.Context — CANCELLATION AND DEADLINES ─────────────────────────
//
// context.Context is Go's standard mechanism for:
//   - Cancellation: "stop what you're doing"
//   - Deadlines: "stop if not done by time T"
//   - Timeouts: "stop if not done in duration D"
//   - Values: passing request-scoped data (trace IDs, auth tokens)
//
// Rules:
//   - ALWAYS pass context as the FIRST parameter: func Do(ctx context.Context, ...)
//   - ALWAYS check ctx.Done() in long-running operations
//   - NEVER store a context in a struct (pass it explicitly)
//   - NEVER pass a nil context — use context.Background() or context.TODO()
//
// The context tree:
//   Background → WithCancel → WithTimeout → WithValue
//   Cancelling a parent cancels all its children.

func doWork(ctx context.Context, name string, duration time.Duration) error {
	fmt.Printf("[%s] starting, will take %v\n", name, duration)

	select {
	case <-time.After(duration):
		fmt.Printf("[%s] done!\n", name)
		return nil
	case <-ctx.Done():
		fmt.Printf("[%s] cancelled: %v\n", name, ctx.Err())
		return ctx.Err()
	}
}

func contextDemo() {
	fmt.Println("=== 1. context.Context ===")

	// WithCancel: explicit cancellation
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		time.Sleep(100 * time.Millisecond)
		fmt.Println("  cancelling context...")
		cancel()
	}()

	err := doWork(ctx, "task-1", 500*time.Millisecond)
	fmt.Println("  result:", err)

	// WithTimeout: automatically cancelled after duration
	ctx2, cancel2 := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel2()  // always defer cancel to avoid leak

	err = doWork(ctx2, "task-2", 200*time.Millisecond)
	fmt.Println("  timeout result:", err)

	// WithDeadline: cancel at a specific time
	deadline := time.Now().Add(30 * time.Millisecond)
	ctx3, cancel3 := context.WithDeadline(context.Background(), deadline)
	defer cancel3()

	err = doWork(ctx3, "task-3", 100*time.Millisecond)
	fmt.Println("  deadline result:", err)

	// Context values — for request-scoped metadata
	type key string
	ctx4 := context.WithValue(context.Background(), key("requestID"), "req-abc-123")
	requestID := ctx4.Value(key("requestID")).(string)
	fmt.Printf("  request ID from context: %s\n", requestID)
}

// ─── 2. errgroup — PARALLEL GOROUTINES WITH ERROR HANDLING ───────────────────
//
// errgroup is from golang.org/x/sync — it's not in stdlib but widely used.
// We'll implement the core pattern manually since we can't import external packages.
//
// errgroup.Group:
//   g.Go(func() error { ... })  — launch goroutine
//   g.Wait()                    — wait for all, return first non-nil error

// Simple errgroup implementation
type Group struct {
	wg   sync.WaitGroup
	mu   sync.Mutex
	err  error
}

func (g *Group) Go(fn func() error) {
	g.wg.Add(1)
	go func() {
		defer g.wg.Done()
		if err := fn(); err != nil {
			g.mu.Lock()
			if g.err == nil {
				g.err = err
			}
			g.mu.Unlock()
		}
	}()
}

func (g *Group) Wait() error {
	g.wg.Wait()
	return g.err
}

func errgroupDemo() {
	fmt.Println("\n=== 2. errgroup Pattern ===")

	g := &Group{}

	// Launch 5 tasks in parallel
	for i := 0; i < 5; i++ {
		i := i  // capture
		g.Go(func() error {
			time.Sleep(time.Duration(rand.Intn(50)) * time.Millisecond)
			if i == 3 {
				return fmt.Errorf("task %d failed", i)
			}
			fmt.Printf("  task %d succeeded\n", i)
			return nil
		})
	}

	if err := g.Wait(); err != nil {
		fmt.Println("  at least one task failed:", err)
	} else {
		fmt.Println("  all tasks succeeded")
	}
}

// ─── 3. WORKER POOL PATTERN ───────────────────────────────────────────────────
//
// Bound concurrency to N workers instead of spawning unlimited goroutines.
// Use when: downstream service has rate limits, resources are expensive.

type Job struct {
	ID   int
	Data string
}

type Result struct {
	JobID int
	Value string
	Err   error
}

func worker(id int, jobs <-chan Job, results chan<- Result) {
	for job := range jobs {
		// Simulate work
		time.Sleep(time.Duration(rand.Intn(20)) * time.Millisecond)

		var result Result
		if rand.Intn(10) == 0 {  // 10% failure rate
			result = Result{JobID: job.ID, Err: fmt.Errorf("job %d failed", job.ID)}
		} else {
			result = Result{JobID: job.ID, Value: fmt.Sprintf("processed(%s)", job.Data)}
		}
		fmt.Printf("  worker %d: job %d → %q err=%v\n", id, job.ID, result.Value, result.Err)
		results <- result
	}
}

func workerPool() {
	fmt.Println("\n=== 3. Worker Pool ===")

	const numWorkers = 3
	const numJobs = 10

	jobs := make(chan Job, numJobs)
	results := make(chan Result, numJobs)

	// Start N workers
	var wg sync.WaitGroup
	for i := 1; i <= numWorkers; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			worker(id, jobs, results)
		}(i)
	}

	// Send jobs
	for i := 1; i <= numJobs; i++ {
		jobs <- Job{ID: i, Data: fmt.Sprintf("data-%d", i)}
	}
	close(jobs)  // signals workers to stop when queue is empty

	// Wait for workers, then close results
	go func() {
		wg.Wait()
		close(results)
	}()

	// Collect results
	var succeeded, failed int
	for r := range results {
		if r.Err != nil { failed++ } else { succeeded++ }
	}
	fmt.Printf("Results: %d succeeded, %d failed\n", succeeded, failed)
}

// ─── 4. RATE LIMITER ──────────────────────────────────────────────────────────
//
// Limit the rate of operations using a token bucket via time.Ticker.

type RateLimiter struct {
	ticker  *time.Ticker
	tokens  chan struct{}
	done    chan struct{}
}

func NewRateLimiter(rps int) *RateLimiter {
	rl := &RateLimiter{
		ticker: time.NewTicker(time.Second / time.Duration(rps)),
		tokens: make(chan struct{}, rps),
		done:   make(chan struct{}),
	}
	go func() {
		for {
			select {
			case <-rl.ticker.C:
				select {
				case rl.tokens <- struct{}{}:
				default:  // bucket full — drop token
				}
			case <-rl.done:
				return
			}
		}
	}()
	return rl
}

func (rl *RateLimiter) Wait(ctx context.Context) error {
	select {
	case <-rl.tokens:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

func (rl *RateLimiter) Stop() {
	rl.ticker.Stop()
	close(rl.done)
}

func rateLimiterDemo() {
	fmt.Println("\n=== 4. Rate Limiter ===")

	// 5 requests per second
	rl := NewRateLimiter(5)
	defer rl.Stop()

	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Millisecond)
	defer cancel()

	var count int
	for {
		if err := rl.Wait(ctx); err != nil {
			fmt.Printf("  stopped after %d requests: %v\n", count, err)
			break
		}
		count++
		fmt.Printf("  request %d at %v\n", count, time.Now().Format("15:04:05.000"))
	}
}

// ─── 5. PIPELINE WITH CANCELLATION ───────────────────────────────────────────
//
// A production-grade pipeline where all stages respect context cancellation.

func generateWithCtx(ctx context.Context, nums ...int) <-chan int {
	out := make(chan int)
	go func() {
		defer close(out)
		for _, n := range nums {
			select {
			case out <- n:
			case <-ctx.Done():
				return
			}
		}
	}()
	return out
}

func squareWithCtx(ctx context.Context, in <-chan int) <-chan int {
	out := make(chan int)
	go func() {
		defer close(out)
		for n := range in {
			select {
			case out <- n * n:
			case <-ctx.Done():
				return
			}
		}
	}()
	return out
}

func pipelineWithCancellation() {
	fmt.Println("\n=== 5. Pipeline with Cancellation ===")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Millisecond)
	defer cancel()

	nums := generateWithCtx(ctx, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10)
	squares := squareWithCtx(ctx, nums)

	for sq := range squares {
		fmt.Printf("  square: %d\n", sq)
	}
	fmt.Println("  pipeline done (may be early due to timeout)")
}

// ─── 6. DONE CHANNEL PATTERN FOR GRACEFUL SHUTDOWN ───────────────────────────
//
// Real servers need to handle OS signals (Ctrl+C, SIGTERM) and shut down
// gracefully — finish in-flight requests, close connections, flush data.

type Server struct {
	done chan struct{}
	wg   sync.WaitGroup
}

func NewServer() *Server {
	return &Server{done: make(chan struct{})}
}

func (s *Server) handleRequest(id int) {
	s.wg.Add(1)
	go func() {
		defer s.wg.Done()
		fmt.Printf("  [req %d] started\n", id)
		select {
		case <-time.After(time.Duration(rand.Intn(100)) * time.Millisecond):
			fmt.Printf("  [req %d] completed\n", id)
		case <-s.done:
			fmt.Printf("  [req %d] abandoned (shutting down)\n", id)
		}
	}()
}

func (s *Server) Shutdown() {
	close(s.done)    // signal all goroutines to stop
	s.wg.Wait()      // wait for in-flight requests to finish (or abort)
	fmt.Println("  server shutdown complete")
}

func gracefulShutdown() {
	fmt.Println("\n=== 6. Graceful Shutdown ===")

	server := NewServer()

	// Simulate incoming requests
	for i := 1; i <= 5; i++ {
		server.handleRequest(i)
		time.Sleep(10 * time.Millisecond)
	}

	// Trigger shutdown
	time.Sleep(30 * time.Millisecond)
	fmt.Println("  initiating shutdown...")
	server.Shutdown()
}

// ─── 7. CONTEXT PROPAGATION THROUGH A STACK ──────────────────────────────────
//
// Best practice: every function in a call chain that does I/O or blocking
// work should accept a context as its first argument.

type UserService struct{}
type OrderService struct{}

func (s *UserService) GetUser(ctx context.Context, id int) (string, error) {
	select {
	case <-ctx.Done():
		return "", ctx.Err()
	case <-time.After(10 * time.Millisecond):
		return fmt.Sprintf("user-%d", id), nil
	}
}

func (s *OrderService) GetOrders(ctx context.Context, userID int) ([]string, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case <-time.After(20 * time.Millisecond):
		return []string{fmt.Sprintf("order-1-for-%d", userID), "order-2"}, nil
	}
}

type Handler struct {
	users  *UserService
	orders *OrderService
}

func (h *Handler) GetUserWithOrders(ctx context.Context, userID int) error {
	// Context flows through the entire call chain
	user, err := h.users.GetUser(ctx, userID)
	if err != nil {
		return fmt.Errorf("get user: %w", err)
	}

	orders, err := h.orders.GetOrders(ctx, userID)
	if err != nil {
		return fmt.Errorf("get orders for %s: %w", user, err)
	}

	fmt.Printf("  User: %s, Orders: %v\n", user, orders)
	return nil
}

func contextPropagation() {
	fmt.Println("\n=== 7. Context Propagation ===")

	h := &Handler{users: &UserService{}, orders: &OrderService{}}

	// Enough timeout
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()
	err := h.GetUserWithOrders(ctx, 42)
	fmt.Println("  with enough timeout:", err)

	// Too short timeout
	ctx2, cancel2 := context.WithTimeout(context.Background(), 5*time.Millisecond)
	defer cancel2()
	err = h.GetUserWithOrders(ctx2, 42)
	fmt.Println("  with short timeout:", err)
}

func main() {
	rand.New(rand.NewSource(time.Now().UnixNano()))

	contextDemo()
	errgroupDemo()
	workerPool()
	rateLimiterDemo()
	pipelineWithCancellation()
	gracefulShutdown()
	contextPropagation()
}

/*
THOUGHT QUESTIONS:

1. "Always pass context as the first parameter." Why is this a convention
   and not enforced by the language? What would Go need to enforce it?

2. What is a goroutine leak? How does context cancellation prevent goroutine leaks
   in the pipeline pattern?

3. A worker pool has N=3 workers but receives 100 jobs. How many goroutines
   are created? How does this compare to creating one goroutine per job?

4. context.WithValue stores values using interface{} keys. Why is it recommended
   to use unexported type keys rather than plain strings?

5. Why is `defer cancel()` important after `context.WithTimeout` even though
   the timeout will eventually cancel the context anyway?

EXERCISES:

1. Implement a true token bucket rate limiter with burst capacity:
   Allow B burst requests immediately, then 1 request per second thereafter.

2. Write a circuit breaker that opens (stops sending) after N consecutive
   failures and attempts to half-open after a timeout.

3. Implement a parallel HTTP fetcher: given N URLs, fetch them concurrently
   with max M parallel requests, respecting a context timeout.

4. Write a "singleflight" implementation: if 5 goroutines all call the same
   expensive function simultaneously, only 1 actually runs — others wait and
   share the result. (This is golang.org/x/sync/singleflight's behavior.)
*/
