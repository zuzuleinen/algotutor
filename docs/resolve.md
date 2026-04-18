# Re-solve Mode

Re-solve mode tests whether the user can *reproduce* a solution, not just recognize it. Every successfully solved
problem goes on a Leitner ladder — when it comes due, `train` hands it back with a fresh template and the prior
solution hidden, and the user must solve it again from scratch. This is separate from card review: cards test atomic
recall, re-solve tests end-to-end execution.

## `resolve.json` format

If the file doesn't exist, treat it as `{"schedule": {}, "concept_failures": {}}` and create it on first write.

Schema:

```json
{
  "schedule": {
    "014": {
      "step": 0,
      "due": "2026-04-23T00:00:00Z",
      "first_solved_at": "2026-04-16T00:00:00Z",
      "concept": "binary-search",
      "level": 1,
      "history": []
    }
  },
  "concept_failures": {
    "binary-search": 0
  }
}
```

Per-problem fields (`schedule[NNN]`):

- `step` — integer 0..4, index into the interval ladder. Starts at 0.
- `due` — RFC3339 timestamp when this problem is next eligible for re-solve.
- `first_solved_at` — RFC3339 timestamp of the original successful `check`.
- `concept` — the concept the problem trained. Cached here to avoid re-reading the problem file during the picker.
- `level` — the problem's level at time of first solve. Cached.
- `history` — array of re-solve attempts: `{"at": "<RFC3339>", "outcome": "clean|recovered|failed"}`.

Per-concept field (`concept_failures[<concept>]`):

- Integer. Counts *consecutive* failed re-solves on that concept. Any non-failed re-solve outcome on the same concept
  resets it to 0. Used by the level-drop rule.

## The ladder

Intervals in days, indexed by `step`:

| step | interval (days) |
|------|-----------------|
| 0    | 7               |
| 1    | 21              |
| 2    | 60              |
| 3    | 180             |
| 4    | 365             |

A problem at `step = 4` stays there on clean solves and keeps rescheduling every 365 days.

## Registration (on first solve)

When the Checking Flow marks a problem as `solved` for the first time (not a re-solve, not a drill, not a scaffold
sub-problem), add an entry to `resolve.json.schedule` keyed by the problem number:

```json
{
  "step": 0,
  "due": "<solved_at + 7 days, RFC3339>",
  "first_solved_at": "<solved_at, RFC3339>",
  "concept": "<concept from problem file>",
  "level": <level from problem file>,
  "history": []
}
```

## Exclusions

The following are NOT registered in `resolve.json` on solve:

- Drills (`kind: drill` in the problem file).
- Scaffolded sub-problems (filename has a letter suffix: `007a.md`, `007b.md`, etc.).
- Mix problems (see `docs/mix.md`).

## Picker (inside `train`)

Run as the **Re-solve check** step of the Training Flow preamble (after drill check, before digest):

1. Read `resolve.json` (treat as empty if missing). Find the single entry with the earliest `due`. If `due > now`, do
   nothing — fall through to digest/picker.
2. Otherwise, let `NNN` be that entry's problem number. Read `problems/NNN.md` for the statement.
3. Write a **fresh `main.go` template**: problem description as a comment at the top, `main` first with
   `fmt.Println(...)` calls that have inline expected-output comments, and an empty target function body. **Do not
   include the user's prior solution.** If you don't have the prior solution in context, good — you shouldn't use it
   even if you did.
4. Do NOT re-read or re-present scaffold sub-problems (`NNNa.md`, `NNNb.md`). Those were one-shot aids for the first
   solve.
5. Write `NNN:resolve` to `current.md`.
6. Announce in one sentence: "Problem NNN (<title>) is due for re-solve."

## Outcomes (on re-solve `check`)

Classify the outcome based on what happened during this re-solve session:

- **Clean** — `check` passes on the first call; no "I don't know" was issued; no prior failed `check` this session.
- **Recovered** — `check` eventually passes, but at least one of the following occurred: scaffolding was used, a
  `check` call failed before the passing one, or the user asked for any form of hint.
- **Failed** — the user explicitly gives up before a successful `check` (see "Giving up" below).

State updates per outcome:

| outcome   | `step`                       | `due`                            | `concept_failures[concept]`          | concept level                                |
|-----------|------------------------------|----------------------------------|--------------------------------------|----------------------------------------------|
| clean     | `min(step + 1, 4)`           | `now + interval[new step]`       | reset to 0                           | unchanged                                    |
| recovered | unchanged                    | `now + interval[same step]`      | reset to 0                           | unchanged                                    |
| failed    | reset to 0                   | `now + 7 days`                   | `concept_failures[concept] + 1`      | drop by 1 ONLY if the new counter reached 2, then reset counter to 0; floor at 0 |

Always append to `history`:

```json
{"at": "<now, RFC3339>", "outcome": "clean" | "recovered" | "failed"}
```

After writing `resolve.json`, reset `current.md` to empty. Leave `main.go` alone — the user may want to look at their
own work before running the next `train`.

## Mistake tracking on re-solve

- On a **failed** re-solve, log a mistake in `mistakes.json` with `trigger: "resolve"` and a category picked from the
  taxonomy. If the user never attempted any code, use `other` with a `note` describing the stall.
- On a **recovered** re-solve, mistakes logged by the scaffolding flow during this session use `trigger: "resolve"`
  (not `"scaffold"`) to distinguish re-solve gaps from first-solve gaps.
- On a **clean** re-solve, resolve open mistakes in any category this problem *could have* exercised, same rules as a
  normal clean solve.

## Scaffolding during re-solve

Scaffolding is allowed during re-solve and follows the normal Scaffolding Flow, with two additions:

1. Any use of scaffolding this session downgrades the outcome to **recovered** at best (even if the parent check
   eventually passes cleanly afterward).
2. Mistake entries logged during re-solve scaffolding use `trigger: "resolve"` instead of `"scaffold"`.

Sub-problem files generated during re-solve scaffolding use the normal `NNNa.md`, `NNNb.md` convention and are
excluded from the re-solve ladder (same as first-solve scaffolds).

## Giving up

The user signals a failed re-solve by saying "give up", "fail this", "skip re-solve", or any clear equivalent.
Do NOT auto-fail a re-solve based on scaffold depth or elapsed time — only the user decides when to stop.

When the user gives up:

1. Log the failure to `resolve.json` per the outcome table above.
2. Log a mistake entry with `trigger: "resolve"`.
3. Reset `current.md` to empty.
4. Tell the user: "Marked re-solve failed. `<concept>` consecutive failures: <N>. Level drop: <yes/no>."

The next `train` re-runs the preamble from scratch.
