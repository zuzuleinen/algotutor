# algotutor

An AI-powered algorithmic training system for Go developers. Open a Claude Code session in this directory, type `train`, and start solving
problems.

## How it works

Claude acts as a tutor that generates progressively harder algorithm problems in Go. It tracks your skill level across
32 concepts — from arrays and strings up through dynamic programming and system design — and picks the next problem
based on where you are.

### Commands

| Command                          | What it does                                |
|----------------------------------|---------------------------------------------|
| `train`                          | Get the next problem — drill, re-solve, mix, or new, based on state |
| `check`                          | Submit your solution for evaluation         |
| `I don't know`                   | Break the problem into simpler sub-problems |
| `I want to solve [problem name]` | Request a specific problem                  |
| `review`                         | Check if you have cards due for review      |
| `mistakes`                       | Show your recurring-error report            |

### Concepts covered

**Fundamentals** — arrays, strings, loops, nested loops, math

**Core Data Structures** — maps, sets, matrix, stacks, queues, linked lists, heaps

**Core Techniques** — sorting, binary search, two pointers, sliding window, prefix sums, bit manipulation

**Recursion and Trees** — recursion, trees, tries

**Graph Algorithms** — graphs, topological sort, union-find, shortest path

**Advanced Techniques** — greedy, intervals, backtracking, divide and conquer, dynamic programming, monotonic stacks,
design

### Spaced repetition review

As you solve problems, Claude automatically creates review cards capturing what you learned — algorithmic patterns, Go
syntax, data structure properties. Cards follow
the [SuperMemo 20 Rules for effective memorization](https://www.supermemo.com/en/blog/twenty-rules-of-formulating-knowledge).

Run `make review` (or `go run ./cmd/review`) to start an Anki-style review session. The review TUI uses
the [FSRS](https://github.com/open-spaced-repetition/go-fsrs) algorithm to schedule cards. Rate each card 1–4 (
Again/Hard/Good/Easy) and it will reappear at the optimal interval.

<img src="img.png" width="600" alt="img.png">

<img src="img_1.png" width="600" alt="img_1.png">

### Mistake tracking

Every failed `check` is tagged with a fixed error taxonomy (off-by-one, forgotten-update, missed base case, empty-input
missed, wrong-algorithm, and ~25 more) and logged to `mistakes.json`. Gaps that would otherwise evaporate at the end of
a session stick around as data.

When any category accumulates ≥ 3 unresolved entries in your recent history, `train` stops picking a new concept and
instead hands you a tiny single-category drill — five-line problems stripped of surrounding concept, aimed at exactly
that failure mode. Solve it and the oldest open mistakes in that category close out.

Every 7 days, `train` prints a short digest of your top recurring categories. Run `mistakes` any time to see the full
report on demand. Drills do not raise concept levels — their only effect is to patch the pattern.

### Re-solve

Solving a problem once isn't mastery. Every successfully solved problem enters a Leitner schedule (7 / 21 / 60 / 180 /
365 days) in `resolve.json`. When a problem comes due, `train` hands it back with a fresh `main.go` template — your
previous solution is hidden — and you re-solve it from scratch.

A clean re-solve pushes the next due date further out. Needing scaffolding holds the step. Giving up (`give up`,
`fail this`, `skip re-solve`) resets the ladder, and **two consecutive failed re-solves on the same concept** drop
its level by one. Re-solves preempt new training the moment anything is due.

### Mix

Research shows interleaved practice beats grinding one concept at a time. Once you have 5+ concepts at level 2+ and
at least 3 have gone cold (untouched for 14+ days), `train` starts a mix session — 3 problems from 3 different
concepts, one after the other, each drawn one level below your working level so the context switching itself is the
challenge.

Mix doesn't raise concept levels. It updates a per-concept retention score in `retention.json` — clean mix solves push
retention up (+0.2), failures push it down (−0.3), and retention decays 0.1 per 14 days a concept sits untouched. Low
retention shows up as a nudge on `train`.

Mix sessions trigger at most once a week and run one at a time — there's no command to force one. Inside a mix,
scaffolding and `I don't know` work normally; say `skip` to drop a problem and move on, or `end mix` to abandon the
session.

## Requirements

- [Claude Code](https://docs.anthropic.com/en/docs/claude-code)
- [Go](https://go.dev/)

## Getting started

1. Clone the repo
2. Open a Claude Code session in the directory
3. Type `train`

On first run, Claude will initialize your progress file and problem directory. Your progress is local and gitignored.

## Recommendations

You can use `claude --dangerously-skip-permissions` to not be prompted all the time.

I had good experience with `Claude Sonnet 4.6` — it's set as the default in `.claude/settings.json`.

The working problem is always inside `main.go`. You can validate with `go run .` before asking `claude check`.

Try to make as much progress as you can before saying `I don't know`. This way Claude can better assess your gaps and
missing prerequisites.

It should feel effortful.
