# Concept List — Go Concurrency

31 concepts in teaching order. Prerequisites in parentheses are the *conceptual*
prerequisites — the vocabulary and intuitions the user must already have in place to even
understand the explanation of the new concept. Before training any concept, every
prerequisite must be at level ≥ 1 in `progress.md`. If a prerequisite is below level 1, do
not teach the new concept — reroute training to the missing prerequisite first, even if it
comes later in this list.

The list is ordered so that prerequisites always appear before dependents.

## Fundamentals

1. **concurrency-fundamentals** — concurrency vs parallelism, the Go runtime, goroutine
   cost, GOMAXPROCS. Recognition cue: "this program needs to do many independent things
   that overlap in time."
2. **goroutines** — `go` keyword, lifetime, scheduler, why goroutines are not OS threads.
   Recognition cue: "I want this work to run in the background while my program continues."
   *(requires: concurrency-fundamentals)*

## Synchronization Without Channels

3. **waitgroup** — `sync.WaitGroup` Add/Done/Wait, completion signaling, Add-before-go
   discipline, common patterns and pitfalls. Recognition cue: "I need to wait for N
   goroutines to finish, but they don't return values I need." *(requires: goroutines)*

## Channels

4. **channels-unbuffered** — `make(chan T)`, send/receive, the rendezvous semantics
   (sender blocks until receiver is ready and vice versa), the slogan "share memory by
   communicating." Recognition cue: "I want one goroutine to hand a value to another
   goroutine safely." *(requires: goroutines)*
5. **channels-buffered** — `make(chan T, n)`, capacity, when send blocks (full) vs when
   receive blocks (empty), decoupling producer/consumer rates. Recognition cue: "I have a
   bursty producer and a steady consumer (or vice versa) and don't want them lockstepped."
   *(requires: channels-unbuffered)*
6. **channel-direction** — `chan<- T` (send-only) and `<-chan T` (receive-only) in
   function signatures, expressing ownership and intent at the API boundary. Recognition
   cue: "I'm writing a function that takes a channel and want the type system to enforce
   what role it plays." *(requires: channels-unbuffered)*
7. **channel-close-range** — `close(ch)`, `for v := range ch`, the comma-ok idiom
   (`v, ok := <-ch`), what happens when you close (subsequent receives drain then yield
   zero values; subsequent sends panic), broadcast-via-close. Recognition cue: "I want
   receivers to know the stream is finished." *(requires: channels-unbuffered)*
8. **channel-ownership** — the discipline of "the goroutine that owns a channel is the
   one that closes it" (typically the sender), why receivers should not close, why
   shared-write channels need an extra coordinator. *(requires: channel-close-range)*

## Select

9. **select** — wait on multiple channel operations, the random-pick semantics when
   several cases are ready, why `select` is *the* concurrency control flow primitive in
   Go. Recognition cue: "I'm waiting on more than one channel and want to react to
   whichever fires first." *(requires: channels-unbuffered)*
10. **select-default** — `default:` makes the select non-blocking (try-send / try-receive).
    Recognition cue: "I want to send if anyone's listening, otherwise drop and move on."
    *(requires: select)*
11. **select-timeout** — `case <-time.After(d):`, `context.Done()` with select, building
    deadlines. Recognition cue: "I'm willing to wait at most D for this channel op to
    succeed." *(requires: select)*

## Coordination Primitives

12. **confinement** — the single-writer principle: if only one goroutine ever writes to a
    given piece of state, you need no synchronization for that state. The cleanest way to
    avoid races: don't share. *(requires: goroutines)*
13. **mutex** — `sync.Mutex`, Lock/Unlock, the `defer mu.Unlock()` idiom, what a critical
    section is, why double-locking by the same goroutine deadlocks (Mutex is not
    reentrant in Go). Recognition cue: "Multiple goroutines update the same in-memory
    structure and confinement is impractical." *(requires: goroutines)*
14. **rwmutex** — `sync.RWMutex`, RLock/RUnlock, when reads dramatically outnumber writes,
    when a plain Mutex is actually faster. *(requires: mutex)*
15. **once** — `sync.Once.Do(fn)`, lazy initialization, why this is *not* the same as a
    nil check + mutex (Once handles re-entry from `fn` itself differently). *(requires:
    mutex)*

## Memory Model

16. **memory-model** — happens-before, why a write in goroutine A may not be visible to
    goroutine B without synchronization, what synchronization primitives establish
    happens-before, why naive shared-variable communication breaks. *(requires: goroutines,
    mutex)*
17. **atomic** — `sync/atomic` Load/Store/CAS for primitive types, when atomics are
    appropriate (counters, flags) and when they are not (compound state). Recognition
    cue: "I have one integer being updated from many goroutines." *(requires: memory-model)*
18. **atomic-value** — `atomic.Value` for swapping immutable snapshots of arbitrary
    types, the type-consistency rule. *(requires: atomic)*
19. **race-detection** — `go test -race` and `go run -race`, reading race-detector
    reports, the false-negative limitation (only races that actually fire are reported).
    *(requires: memory-model, mutex, channels-unbuffered)*

## Failure Modes

20. **deadlocks** — recognition (`fatal error: all goroutines are asleep - deadlock!`),
    common shapes (send with no receiver, receive with no sender, lock-order inversion,
    circular channel waits), prevention discipline (consistent lock ordering, timeouts,
    bounded buffers). *(requires: mutex, channels-unbuffered, select)*
21. **goroutine-leaks** — symptoms (`runtime.NumGoroutine` climbs over time), causes
    (forgotten close, blocked sends, unbounded receivers), audits via test cleanups,
    prevention via context-cancellation discipline. *(requires: channels-unbuffered,
    goroutines)*

## Context

22. **context-basics** — what `context.Context` is, why it propagates, the contract of
    "the goroutine that received a context must respect Done", the difference between
    Background and TODO. *(requires: channels-unbuffered, select)*
23. **context-cancel** — `context.WithCancel`, the cancellation tree (canceling a parent
    cancels all children), explicit `cancel()` to free resources. *(requires: context-basics)*
24. **context-timeout** — `context.WithTimeout` and `context.WithDeadline`, integrating
    with `select` and network calls, why you almost always want a timeout on outbound
    calls. *(requires: context-cancel)*
25. **context-value** — `context.WithValue`, what it's for (request-scoped data), what
    it's *not* for (function arguments), the type-key idiom. *(requires: context-basics)*

## Higher-Level Patterns

26. **errgroup** — `golang.org/x/sync/errgroup`, parallel work that may fail, first-error
    semantics, integration with context cancellation. Recognition cue: "I want to run N
    things concurrently and stop everything if any of them fails." *(requires: waitgroup,
    context-cancel)*
27. **semaphore** — bounded concurrency via a buffered channel of tokens, or via
    `golang.org/x/sync/semaphore`. Recognition cue: "I have N tasks but I want at most K
    of them running at a time." *(requires: channels-buffered)*
28. **pattern-pipeline** — connecting stages with channels, each stage owns its output,
    stages close their output when done so downstream `range` exits. *(requires: channel-
    close-range, channel-ownership, channel-direction)*
29. **pattern-fan-out-fan-in** — multiple goroutines reading from one input channel
    (fan-out), multiple goroutines merging into one output channel (fan-in), the merge
    function with WaitGroup-then-close. *(requires: pattern-pipeline, waitgroup)*
30. **pattern-worker-pool** — fixed N workers consuming jobs from a channel, dispatching
    pattern, results collection. Recognition cue: "Stream of independent tasks, bounded
    parallelism." *(requires: channels-unbuffered, waitgroup)*
31. **pattern-graceful-shutdown** — context cancellation + WaitGroup + drain, the rule
    "stop accepting work, drain in-flight, then exit." Recognition cue: "Server needs to
    handle SIGTERM cleanly." *(requires: context-cancel, waitgroup, channel-close-range)*

## Optional Advanced (not gated; level 3+)

These can be picked when the user is interview-prepping or asks specifically:

- **sync-cond** — `sync.Cond.Wait/Signal/Broadcast`, condition predicates, the
  while-not-loop idiom. Rarely needed; usually channels are clearer.
- **sync-pool** — `sync.Pool` for object reuse to reduce GC pressure, the "may evict at
  any time" semantics.
- **sync-map** — `sync.Map` and the narrow read-heavy / write-rare scenarios it actually
  beats `map+RWMutex`.
- **singleflight** — `golang.org/x/sync/singleflight` for request coalescing.
- **pattern-rate-limit** — token-bucket rate limiting, `golang.org/x/time/rate`.
- **or-channel** — `or(channels...) <-chan T` returning a channel that closes when any of
  the inputs close.
- **tee-channel** — splitting one channel's values to two consumers.
- **bridge-channel** — flattening a channel-of-channels into a single channel.
