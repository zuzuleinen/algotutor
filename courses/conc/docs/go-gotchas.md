# Go Concurrency Gotchas

Language-level traps relevant to the concurrency course. **Consult before writing any
problem statement, example, or nudge that touches the affected mechanic.** Nudging the
user toward a "simpler" form that is actually broken bakes a latent bug into their mental
model.

## Channels

### Receive on a closed channel returns the zero value silently

`<-ch` after `close(ch)` does not block and does not error — it returns `T{}` (the zero
value of the channel's element type). Use the comma-ok form to distinguish:

```go
v, ok := <-ch
if !ok {
    // channel closed and drained
}
```

When teaching `range ch`, mention this — the loop exits when the comma-ok form would
return `ok == false`.

### Send on a closed channel panics

`ch <- v` after `close(ch)` panics with `send on closed channel`. This is *runtime*
behavior — no compile-time guard. The conventional defense is the **channel ownership**
discipline: the sender owns the channel and the sender is the only one allowed to close
it.

### A nil channel blocks forever — sometimes intentionally

`var ch chan int; ch <- 1` blocks the goroutine permanently. So does `<-ch`. This sounds
like a bug, but it's a **feature in `select`**: setting a case channel to `nil` disables
that case, which is the idiomatic way to dynamically turn off a select arm.

```go
var dataCh chan int           // start nil — select arm disabled
// ... later, after we open the source:
dataCh = make(chan int)
// ... when the source is exhausted:
dataCh = nil                   // disable the arm again
```

If a problem statement says "the function should not block," explicitly say so — don't
assume nil-channel behavior is wrong.

### Sending on an unbuffered channel from the same goroutine deadlocks

```go
ch := make(chan int)
ch <- 42        // deadlock — no receiver
fmt.Println(<-ch)
```

This is a common first-time bug. The fix is either a goroutine for the send or a buffered
channel of capacity ≥ 1.

### `len(ch)` and `cap(ch)` are racy

`len(ch)` returns the number of elements currently buffered; `cap(ch)` returns the
buffer capacity. The capacity is fixed, so `cap` is fine; `len` reports a value that may
already be stale by the time the caller acts on it. Don't use `len(ch)` to decide whether
a send will block — use `select` with `default` instead.

### Receive can drain a closed channel before signaling closed

After `close(ch)`, any values still in the buffer are received normally — receivers see
those values first, *then* see `ok == false`. This is the right behavior for "producer
pushed N items then closed" but matters for "stop signal" patterns.

## Goroutines and the loop variable

### Pre-Go 1.22: capture by reference

```go
for i := 0; i < 3; i++ {
    go func() { fmt.Println(i) }()    // probably prints 3, 3, 3
}
```

Pre-Go 1.22, the closure captures `i` by reference. Fix by passing as argument or shadowing:

```go
for i := 0; i < 3; i++ {
    go func(i int) { fmt.Println(i) }(i)
}
```

Go 1.22+ scopes loop variables per iteration, so the original code prints 0, 1, 2 in some
order. Algotutor targets Go 1.26, so by default the new behavior applies. Still mention
this gotcha — the user will read older codebases where the bug is live.

### `range` over a slice gives the same `v` variable each iteration

Same family as the loop-var bug. Pre-Go 1.22 `for _, v := range xs { go use(&v) }` gave
all goroutines the same `&v`. Same fix; Go 1.22+ fixes the default.

## Mutex and locking

### `sync.Mutex` is not reentrant

Locking a mutex from a goroutine that already holds it deadlocks. Go has no recursive
mutex on purpose — if you're tempted to lock the same mutex twice, restructure.

### `defer mu.Unlock()` is the right idiom — almost always

Use `defer` so unlock runs even on panic or early return. The exception is a short
critical section in a long function where you want the unlock earlier; in that case,
unlock explicitly and don't defer.

```go
mu.Lock()
defer mu.Unlock()
// critical section
```

### Lock-order inversion → deadlock

If goroutine A locks `m1` then `m2`, and goroutine B locks `m2` then `m1`, you have a
classic deadlock. Establish a **canonical lock order** (e.g., always lock by increasing
pointer address or by name) and follow it everywhere.

### Copying a `sync.Mutex` is wrong

Once a mutex has been used (or even just exists in a struct), copying it is undefined
behavior — `go vet` will flag it. Pass mutex-containing structs by pointer.

```go
type S struct { mu sync.Mutex; n int }
func (s S) Inc()  { s.mu.Lock(); ... }   // BUG — value receiver copies the mutex
func (s *S) Inc() { s.mu.Lock(); ... }   // OK — pointer receiver
```

### `sync.RWMutex.RLock` is reentrant only if no writer is waiting

A reader holding the RLock can deadlock if it tries to RLock again *while a writer is
waiting on Lock*. Treat RLock as non-reentrant in practice.

## WaitGroup

### `Add` after `Wait` is racy

```go
var wg sync.WaitGroup
wg.Add(1)
go func() { defer wg.Done(); wg.Add(1); ... }()    // BUG
wg.Wait()
```

The inner `Add(1)` may race against the `Wait()` returning. **Discipline: call `Add`
before launching the goroutine, in the goroutine that does the launching, never inside
the goroutine itself.**

### Copying a `WaitGroup` is wrong

Same family as Mutex copying. Pass by pointer.

### Forgetting `Done` leaks the WaitGroup

If a goroutine panics or returns early without calling `Done`, `Wait` blocks forever. Use
`defer wg.Done()` as the first line of the goroutine body.

## sync.Once

### `Once.Do(f)` runs `f` exactly once across all callers

If `f` panics, the Once is still considered "done" — subsequent `Do` calls return
immediately without re-running. If `f` may fail, return the error from a closure-captured
variable rather than expecting Once to retry.

## Atomics

### Atomic ops compose, but compound state does not

`atomic.AddInt64(&n, 1)` is safe. But `if atomic.LoadInt64(&n) == 0 { atomic.StoreInt64(&n, 1) }`
is not — that's a TOCTOU race. Use `CompareAndSwap` for compound state, or just use a
mutex.

### Mixing atomic and non-atomic access on the same variable is undefined

If any access to a variable uses atomics, **all** accesses must.

## Context

### `context.Background()` vs `context.TODO()`

`Background()` is the root for long-lived servers. `TODO()` is a placeholder when you
haven't decided yet — its presence signals "this should probably plumb through a real
context but I haven't done the refactor." Use `TODO` deliberately, not as default.

### Forgetting to call `cancel`

```go
ctx, cancel := context.WithCancel(parent)
// using ctx but never calling cancel — leaks the timer/goroutine until parent cancels
```

Always `defer cancel()` immediately after creation, even if you think the context will be
canceled by a timeout — the docs explicitly require it.

### `context.Value` is for request-scoped data, not function arguments

If something is required to do the function's job, it should be a parameter. `WithValue`
is for things like trace IDs that flow through *every* function on the call path without
appearing in every signature.

### Goroutines that don't check `ctx.Done()` are leaks waiting to happen

The contract of receiving a context is that you *honor* its cancellation. A goroutine
that takes a context but never selects on `ctx.Done()` will keep running after the user
cancels — that's a leak.

## select

### Empty select blocks forever

`select {}` is the idiomatic "block this goroutine indefinitely" — no surprise, but worth
naming.

### Cases with side effects are evaluated when select runs

```go
select {
case ch <- expensive():    // expensive() runs *before* select picks
}
```

Each case's send-value or receive-target is evaluated **before** select picks a winner.
If you don't want that, compute the value first.

## Race detector

### `-race` only catches races that actually fire

A racy program that happens to interleave benignly during the test run won't trip the
detector. Tests must be designed to *exercise* the concurrency — large iteration counts,
goroutines that block on each other, varied scheduling.

### `-race` adds significant overhead

Roughly 5-10× slowdown and 5-10× memory overhead. Use in CI and dev, not production
binaries.

## time.After in tight select loops leaks until expiration

```go
for {
    select {
    case <-time.After(d):     // each iteration creates a new timer
    case v := <-ch:
        process(v)
    }
}
```

Each iteration allocates a new timer that lives until `d` elapses. For a hot loop, use
`time.NewTimer` and `Reset`, or use `time.NewTicker`.
