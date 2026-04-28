# Problem Bank — Go Concurrency

Curated problems organized by concept and difficulty.

When training, prefer problems from this bank at the appropriate level. Each entry is a
short title; the agent expands it into a full problem file with statement, signature,
example I/O, and expected behavior.

**Level mapping:**
- Level 0 → mechanics warmup (5–15 lines of solution code, "use the primitive directly")
- Level 1 → standard pattern, no surprises
- Level 2 → edge cases, multiple primitives, realistic shape
- Level 3 → larger problem where this concept is one tool among several
- Level 4+ → interview / production-level

**Validation:** `go test -race ./...`. The student writes the solution in `main.go` and the
oracle test in `main_test.go` (the agent provides the test file).

---

## concurrency-fundamentals

### Level 0
- Print CPU and GOMAXPROCS — call `runtime.NumCPU()` and `runtime.GOMAXPROCS(0)` and print both
- Force single-threaded — set GOMAXPROCS(1) and print before/after, observe what changes (nothing in output, but conceptually parallelism is constrained)

### Level 1
- 10k goroutines counter — launch 10k goroutines that sleep briefly, observe `runtime.NumGoroutine()` before / during / after
- Concurrency without parallelism — design a program that is concurrent (multiple goroutines) but explicitly serializes work using GOMAXPROCS(1)

### Level 2
- Predict scheduler behavior — given a small program, predict whether output is deterministic and explain why

---

## goroutines

### Level 0
- Hello goroutine — launch a goroutine that prints "hello", make sure main waits long enough to see it
- Two goroutines — launch two goroutines, each prints its own message; observe non-determinism

### Level 1
- Anonymous goroutine with arg — launch `go func(i int) { ... }(i)` to demonstrate parameter binding
- Goroutine that returns a value via shared variable — show why this is a race (no synchronization)

### Level 2
- Pre-1.22 loop variable bug — write the buggy version (`go func() { fmt.Println(i) }()` in a loop) and the fixed version

### Level 3
- Self-spawning workers — goroutine that spawns more goroutines based on input

---

## waitgroup

### Level 0
- Wait for one goroutine — `wg.Add(1)`, goroutine calls `wg.Done()`, main calls `wg.Wait()`
- Wait for N goroutines — fan out N goroutines, each does work and calls Done

### Level 1
- WaitGroup with closure — closure captures wg by pointer, demonstrates the pointer-not-value rule
- Sum of squares — N goroutines compute squares of inputs into a results slice (with mutex or channel)

### Level 2
- Add inside goroutine bug — show the wait-without-add race, then fix it
- Forgotten Done — show the deadlock from forgotten Done, fix with `defer wg.Done()`

### Level 3
- Bounded launches — launch up to N goroutines, wait for them all, then launch the next batch

---

## channels-unbuffered

### Level 0
- Send and receive — main creates `chan int`, goroutine sends 42, main receives and prints
- Two-way handshake — two goroutines, one sends, one receives, verify rendezvous semantics
- Compute via channel — compute `n*n` in a goroutine, send through channel, receive in main

### Level 1
- Producer-consumer — producer sends 1..10, consumer prints them; both run as goroutines
- Channel-based mutex — show that an unbuffered channel of capacity 1 can simulate a mutex

### Level 2
- Hand-off chain — three goroutines passing a token in a ring
- Detect channel direction — given a `chan int`, decide based on use whether it's a sender or receiver role

### Level 3
- Echo server pattern — receive request, send response, both via channels

---

## channels-buffered

### Level 0
- Buffered send without goroutine — `make(chan int, 1); ch <- 42; <-ch` (no goroutine needed)
- Detect buffer capacity — implement `IsBuffered(ch chan int) bool` using `cap`

### Level 1
- Producer faster than consumer — buffered channel as a queue between fast producer and slow consumer
- Bounded queue — implement a fixed-capacity FIFO using a buffered channel

### Level 2
- Capacity edge cases — full buffer blocks send; empty buffer blocks receive; demonstrate both

### Level 3
- Backpressure — design a producer that pauses when the consumer falls behind

---

## channel-direction

### Level 0
- Send-only function — `func produce(out chan<- int)` that produces 5 ints
- Receive-only function — `func consume(in <-chan int)` that prints what it receives

### Level 1
- Pipeline stage signature — write a stage that takes `<-chan int` and returns `<-chan int`
- Compile-time check — show a function that accepts `chan<- int` rejects a receive operation at compile time

---

## channel-close-range

### Level 0
- Range over channel — producer sends 1..5 then closes; consumer ranges and prints; observe loop exit
- Comma-ok — receive from a closed channel and check `ok`

### Level 1
- Broadcast via close — many goroutines receive on the same channel; closing it wakes them all
- Detect closure — distinguish "channel closed" from "real zero value"

### Level 2
- Close idiom in pipeline — sender closes when done so downstream `range` exits naturally

---

## channel-ownership

### Level 1
- Owner-closes pattern — refactor a "shared close" bug into the canonical owner-closes pattern
- Two-sender coordinator — when two goroutines send to one channel, design a coordinator goroutine that owns the close

### Level 2
- Avoid send-on-closed — given a buggy program where a receiver closes the channel, show why it panics and fix it

---

## select

### Level 0
- Wait on two channels — `select { case <-a: case <-b: }`
- First to fire — print which of two channels delivered first

### Level 1
- Merge two streams — combine two `<-chan int` into one output channel using select
- Pick on N channels — select over a dynamic number of channels (loop)

### Level 2
- Random selection — when multiple cases are ready, observe that select picks randomly

### Level 3
- State machine via select — a goroutine whose behavior depends on which channel fires

---

## select-default

### Level 0
- Try-send — send if a receiver is ready, otherwise drop and continue
- Try-receive — receive a value if available, otherwise return a sentinel

### Level 1
- Non-blocking poll — periodically check a channel without blocking the main loop
- Drop-on-overflow — producer drops values when the channel buffer is full

---

## select-timeout

### Level 0
- Timeout via time.After — wait at most 1s for a value
- Deadline pattern — wait for a value, give up after a fixed duration

### Level 1
- Per-iteration timer leak — show the time.After-in-loop leak; fix with a Timer.Reset
- Timeout with context — replace `time.After` with `<-ctx.Done()` and explain when to use which

### Level 2
- Multiple timeouts — different cases have different deadlines

---

## confinement

### Level 0
- Single-writer counter — one goroutine owns the counter; others ask it via a channel; no mutex

### Level 1
- Confined slice — pass a slice to one goroutine, no other goroutine touches it; verify no race
- Read-only confinement — one goroutine writes once, then never again; others may read concurrently

---

## mutex

### Level 0
- Protect a counter — N goroutines increment a counter under a mutex
- Critical section — read-modify-write of a map under a mutex

### Level 1
- defer mu.Unlock — refactor explicit Unlock to defer; understand why
- Mutex on struct — embed `sync.Mutex` in a struct; method-based locking

### Level 2
- Mutex copy bug — write a value-receiver method that copies the mutex, observe `go vet` flagging it
- Reentry deadlock — call a locking method from within a locking method on the same mutex

### Level 3
- Granularity — refactor one big lock into multiple smaller locks; reason about correctness

---

## rwmutex

### Level 0
- Many readers, one writer — concurrent readers under RLock, writer takes Lock

### Level 1
- When RWMutex helps — benchmark RWMutex vs Mutex on a read-heavy workload
- When RWMutex hurts — same primitives, write-heavy workload; show overhead

### Level 2
- RWMutex non-reentrancy — RLock-while-RLock-while-write-waiting deadlock pattern

---

## once

### Level 0
- Lazy init — `sync.Once.Do` to initialize a singleton

### Level 1
- Concurrent first call — many goroutines call the same Once.Do; only one runs the function
- Once with error — capture an error from the init function via closure

### Level 2
- Once with panic — what happens if the function panics; subsequent Do calls

---

## memory-model

### Level 0
- Visible write — show that a write in goroutine A is not guaranteed visible to goroutine B without sync
- Happens-before via mutex — mutex unlock happens-before next lock; explain visibility

### Level 1
- Happens-before via channel — a send happens-before the corresponding receive; explain visibility

### Level 2
- Reordering — show how the compiler/CPU may reorder unsynchronized writes

---

## atomic

### Level 0
- Atomic counter — replace mutex-protected counter with `atomic.AddInt64`
- Atomic flag — `atomic.LoadInt32` / `StoreInt32` for a boolean flag

### Level 1
- Compare-and-swap — implement a lock-free toggle with CAS
- Counter with CAS — implement increment via CAS loop

### Level 2
- Mixed atomic / plain bug — show the undefined behavior of mixing atomic and plain access

---

## atomic-value

### Level 1
- Hot-swap config — `atomic.Value` to publish new config snapshots without locking readers

### Level 2
- Type consistency rule — show what happens when you Store two different concrete types in the same Value

---

## race-detection

### Level 0
- Run a known-racy program with `go test -race` and read the report
- Fix the race by adding synchronization, re-run, observe clean output

### Level 1
- Iteration-count sensitivity — show that low iteration counts may not trip the detector even with a real race
- False sense of security — two tests, one with the race triggered, one without; only one report

---

## deadlocks

### Level 0
- All goroutines asleep — produce the canonical "fatal error: all goroutines are asleep" panic via send-with-no-receiver

### Level 1
- Lock-order inversion — two goroutines, two mutexes, opposite acquisition orders; observe deadlock
- Fix with consistent ordering — refactor to a canonical lock order

### Level 2
- Channel-induced deadlock — circular wait among goroutines via channels
- Bounded buffer deadlock — producer + consumer where the buffer fills but the consumer is also waiting on the producer

---

## goroutine-leaks

### Level 0
- Forgotten close — producer never closes; consumer's range loop hangs forever
- Forgotten cancel — context.WithCancel without `defer cancel()`; observe goroutine count over time

### Level 1
- Audit via test cleanup — capture `runtime.NumGoroutine()` before and after a test
- Fix-the-leak — given a program where goroutines accumulate, find and fix the leak

### Level 2
- Server with leaks — handler launches a goroutine per request, never bounds; demonstrate growth under load

---

## context-basics

### Level 0
- Pass a context — call a function with `context.Background()`, have it call another that takes a context
- Background vs TODO — when to use which

### Level 1
- Honor Done — goroutine that loops; teach it to exit when `<-ctx.Done()` fires

---

## context-cancel

### Level 0
- WithCancel — create a cancellable context, cancel it, observe `<-ctx.Done()` firing

### Level 1
- Cancel propagates — parent cancel cancels child contexts
- Always defer cancel — show the leak from skipping `cancel()`

### Level 2
- Cancel with multiple goroutines — fan out N goroutines, cancel the parent, all exit

---

## context-timeout

### Level 0
- WithTimeout — context that auto-cancels after 1s; goroutine respects it
- WithDeadline — same but with absolute time

### Level 1
- HTTP-style call — simulate a slow operation; show how to bound it with context-timeout
- Nested timeouts — outer 5s, inner 2s; the inner wins

---

## context-value

### Level 1
- Trace ID propagation — attach a trace ID via WithValue; read it in a deeply nested function
- Type-key idiom — define a private type for the key to avoid collisions

### Level 2
- Anti-pattern — refactor a function that uses context.Value for a required argument back to a real parameter

---

## errgroup

### Level 1
- Parallel fetch — fetch 3 URLs concurrently using errgroup; first error cancels the rest
- Error propagation — show that the first error returned wins; subsequent errors are dropped

### Level 2
- errgroup with limit — `g.SetLimit(N)` for bounded parallelism

---

## semaphore

### Level 1
- Token bucket — buffered channel of tokens, acquire by receive, release by send
- Limit concurrency — at most K goroutines processing items from a stream

### Level 2
- semaphore.Weighted — using `golang.org/x/sync/semaphore` for weighted acquisition

---

## pattern-pipeline

### Level 1
- Two-stage pipeline — produce numbers → square them; channels connect the stages
- Three-stage pipeline — produce → square → sum

### Level 2
- Pipeline with cancellation — every stage respects context cancellation
- Pipeline backpressure — slow consumer slows down the producer naturally

---

## pattern-fan-out-fan-in

### Level 1
- Fan-out workers — N workers reading from one input channel
- Fan-in merge — merge N channels into one with WaitGroup-then-close idiom

### Level 2
- Fan-out + fan-in — combined; show end-to-end ordering guarantees (or lack thereof)

---

## pattern-worker-pool

### Level 1
- Fixed pool — N workers, M jobs; collect results
- Pool with shutdown — pool that exits cleanly when input closes

### Level 2
- Pool with context — workers exit on cancel; in-flight jobs drained
- Adaptive pool — grow / shrink based on queue depth (advanced)

---

## pattern-graceful-shutdown

### Level 1
- SIGTERM handler — main goroutine listens for signal, cancels root context, waits for workers
- Drain in-flight — finish currently-processing jobs, refuse new ones

### Level 2
- Shutdown deadline — give workers up to T seconds to drain; force-exit after
- Shutdown ordering — shutdown stages in pipeline order, not reverse
