# Mistake Tracking — Go Concurrency

The `mistakes.json` schema and drill rules in `docs/mistakes.md` (algos) apply unchanged.
Only the **taxonomy** changes for this course — concurrency failures look completely
different from algorithmic failures.

## Schema reminder

```json
{
  "id": "m_<unix_timestamp>_<seq>",
  "timestamp": "<RFC3339>",
  "problem": "007",
  "concept": "channels-unbuffered",
  "category": "send-without-receiver",
  "note": "sent on unbuffered channel from main, no goroutine reading",
  "trigger": "check",
  "resolved_at": null
}
```

Same field semantics as the algos taxonomy: `category` is exactly one of the values below;
`note` is one line, no code, no formulas; `trigger` is `check`, `scaffold`, or `resolve`.

## Taxonomy

Pick exactly one category per mistake. If nothing fits, use `other` with a specific `note`.
If `other` recurs three times with similar notes, promote it.

### Goroutine lifecycle

- `goroutine-leaked` — a goroutine has no path to exit (forgotten `close`, blocked send,
  unbounded receive without context).
- `wait-without-add` — `WaitGroup.Wait` called before `Add`, or `Add` inside the goroutine
  body racing with `Wait`.
- `forgotten-done` — goroutine returns (or panics) without calling `wg.Done()`,
  `wg.Wait()` blocks forever.
- `goroutine-not-launched` — wrote `f()` instead of `go f()` (or vice versa where parallelism
  was required).

### Channel mechanics

- `send-without-receiver` — unbuffered send with no goroutine ready to receive (deadlock
  panic).
- `receive-without-sender` — unbuffered receive with no sender, no timeout, no context.
- `send-on-closed-channel` — sender closed the channel then sent again, or two senders share
  a channel with no coordinator and one closes.
- `receive-on-closed-channel-misread` — treated the zero value as a real value (forgot the
  comma-ok check or the `range` exit).
- `wrong-buffer-size` — picked unbuffered when buffered was needed (or vice versa) — the
  failure shape was solvable by adjusting capacity.
- `nil-channel-block` — accidentally sent or received on a nil channel (vs intentionally
  using nil to disable a select arm).
- `wrong-close-owner` — receiver closed the channel; or shared-write channel closed
  without coordinator.
- `forgotten-close` — sender finished but didn't `close(ch)`, downstream `range` never
  exits.
- `closed-channel-double-close` — closed an already-closed channel (panic).

### Select

- `select-no-default-when-needed` — should have been non-blocking but blocked.
- `select-default-when-not-needed` — burned CPU spinning because `default` was wrong.
- `select-missing-done` — long-lived `select` with no `case <-ctx.Done()` arm — leak.
- `select-arm-evaluated-too-eagerly` — case expression had a side effect that ran before
  select picked.

### Locking

- `lock-not-held` — read or wrote shared state without holding the mutex.
- `lock-not-released` — Lock without matching Unlock (or returned before Unlock without
  defer).
- `lock-order-inversion` — two goroutines acquired the same pair of locks in different
  orders.
- `mutex-copied` — copied a struct containing a mutex (also caught by `go vet`).
- `wrong-mutex-vs-rwmutex` — used Mutex where RWMutex would scale, or used RWMutex with
  rare reads (overhead).
- `mutex-reentry` — locked a non-reentrant mutex from a goroutine that already held it.
- `defer-unlock-wrong-place` — `defer mu.Unlock()` placed before `mu.Lock()` or in a scope
  that doesn't own the lock.

### Memory model & atomicity

- `data-race` — concurrent unsynchronized access to a shared variable; `go test -race`
  would flag it.
- `compound-atomic` — used atomic Load/Store separately when CAS was needed.
- `atomic-mixed-with-plain` — same variable accessed atomically in some paths and via
  plain read/write in others.

### Context

- `context-not-checked` — goroutine ignored `ctx.Done()`, kept running after cancel.
- `context-not-propagated` — call chain dropped the context (stuffed `context.Background()`
  partway down).
- `cancel-not-called` — created context with cancel but never deferred `cancel()` — leak.
- `value-as-arg` — used `context.WithValue` for what should have been a function argument.

### Patterns

- `unbounded-fanout` — launched a goroutine per item with no semaphore; resource explosion.
- `pipeline-stage-leaks` — pipeline stage didn't close its output, downstream stage hung.
- `worker-pool-no-drain` — closed jobs but didn't wait for workers to finish results.
- `graceful-shutdown-skipped-drain` — exited on signal without draining in-flight work.

### Algorithmic / fidelity

- `wrong-primitive` — problem said "use a mutex"; user used a channel (or vice versa).
- `wrong-pattern` — problem named a pattern (worker pool, pipeline, fan-in); user did
  something else with the same I/O.

### Other

- `misread-problem` — solved a different problem than stated.
- `syntax` — persistent Go syntax errors (not typos); log only on repeat across sessions.
- `other` — last resort.

## When to log

Same rules as algos: max three categories per failed check, deduplicated within an attempt
chain, on scaffolding when the gap is namable.

**Concurrency-specific note:** when `go test -race ./...` reports a race, log
`data-race` even if the test would otherwise pass. The race detector finding is itself
the failure.

## When to resolve

A category resolves when a later clean solve exercises the *same mechanic* without
reproducing the failure. Examples:

- `send-without-receiver` resolves when the user solves a later channel problem without
  triggering the deadlock panic.
- `lock-order-inversion` only resolves on a problem that actually involves two or more
  locks. Don't auto-resolve on a single-lock problem.
- `data-race` only resolves on a problem that actually has shared state under
  `go test -race ./...` and the test passed.

If unsure whether the current problem exercised the category, leave it open.

## Drill rules

Same as algos: single-category focus, stripped to mechanics, obvious oracle, no level
promotion, no scaffolding within a drill.

For concurrency drills, the "obvious oracle" usually means the test in `main_test.go`
either passes under `-race` or it doesn't. Five-line problem, deterministic test.
