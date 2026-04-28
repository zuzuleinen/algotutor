# Algotutor — Multi-Course Training Project

algotutor hosts independent training tracks ("courses"). Each course has its own concepts,
progress, problems, cards, mistakes, re-solve schedule, and retention. The user is enrolled
in one or more courses and trains in **one course at a time**.

## Active Course Resolution

**Always run this before any other flow.**

1. Read `state.json` at the repo root. It looks like:
   ```json
   {
     "enrolled": ["algos", "conc"],
     "active": "algos",
     "default": "algos",
     "default_agent": null
   }
   ```
2. The course slug at `state.active` is the **active course**. Throughout this document,
   `<active>` is a placeholder that means *that slug*. So `courses/<active>/progress.md`
   resolves to `courses/algos/progress.md` when active is `algos`.
3. If `state.json` is missing, ask the user to run `make init`. Do not invent state.
4. If the user types `train <course>`, `review <course>`, `mistakes <course>`, or
   `reset <course>` with a course slug, **first set `state.active` to that slug** (the slug
   must be in `state.enrolled`), then continue. Persist the change to `state.json`.
5. If the user types `train` with no argument and `state.enrolled` has more than one course,
   resolve to `state.active` (last-used), but tell them in one line which course you're
   training: "Training `<active>` (say `train conc` to switch)." If they have only one
   enrolled course, train that one silently.

## Project Structure

```
state.json                       — active course, enrolled list, default agent
main.go                          — current problem template (active course)
main_test.go                     — present only when active course is `conc`
docs/                            — project-wide + shared mechanics docs
  agents.md                      — per-agent setup
  cards.md                       — spaced-repetition card format (shared)
  mix.md                         — mix-mode mechanics (shared)
  resolve.md                     — re-solve mechanics (shared)
courses/<slug>/                  — per-course universe
  progress.md                    — concept levels for this course (user state)
  progress.template.md           — blank progress (zeros)
  current.md                     — current problem pointer
  problems/                      — NNN.md per problem
  problem-bank.md                — curated problems by concept and level
  cards.json                     — review cards (auto-created)
  mistakes.json                  — mistake log (auto-created)
  resolve.json                   — re-solve schedule (auto-created)
  retention.json                 — per-concept retention (auto-created)
  mix.json                       — mix-session state (auto-created)
  docs/
    concepts.md                  — concept list with prerequisites (course-specific)
    go-gotchas.md                — language traps relevant to this course
    mistakes.md                  — course-specific mistake taxonomy + drill rules
cmd/
  init/                          — huh-driven onboarding (`make init`)
  start/                         — flips active course + launches agent (`make train`)
  review/                        — review TUI (`make review [course]`)
internal/
  courses/                       — state.json read/write, path helpers
  migrate/                       — one-shot legacy → multi-course migration
  cards/                         — FSRS card storage
  review/                        — review TUI Bubble Tea model
```

When this document references a path like `courses/<active>/progress.md`, substitute
`<active>` with the active course slug from `state.json`.

## Working Files Per Course

The user works on `main.go` (and `main_test.go` for concurrency) at the **repo root**, not
inside the course directory. The agent always edits the root files, but reads the problem
statement from `courses/<active>/problems/NNN.md`.

| Active course | Files at root              | Validation (under the hood)   |
|---------------|----------------------------|-------------------------------|
| `algos`       | `main.go`                  | `go run .` + `fmt.Println`    |
| `conc`        | `main.go` + `main_test.go` | `go test -race .`             |

The user validates with **`make run`** — a single dispatcher that picks the right command
based on the active course. Tell the user "run `make run` to sanity-check before
`check`," not the underlying Go command. The agent's `check` flow is the *evaluation*
(grade, log mistakes, advance levels); `make run` is the *local smoke test* (does it
compile / pass tests).

When the active course changes, swap the contents of these files to match the new course's
current problem template. **Never carry over `main_test.go` from a `conc` session into an
`algos` session** — delete it when active becomes `algos`.

## Problem Presentation

When the agent presents a problem to the user — in any flow (Picker, Solving a Specific
Problem, Active-mix resume, Drill check, Re-solve check, Scaffolding step-down or
step-back-up) — the **text response** that announces the problem must follow this shape:

1. **No preamble.** Do NOT begin with "I've written…", "I've created…", "Here's your
   problem…", or any commentary about the tool calls that just ran. The user can already
   see those.
2. **Lead with a strong divider** so the eye lands on a clear "new section starts here"
   marker:

   ```
   ═══════════════════════════════════════════════════════════════════════════
   ```

   No blank-line padding before it — terminal markdown renderers collapse repeated blank
   lines unreliably. The divider itself is the visual anchor.
3. **Problem heading on the next line** as a bold title: `**Problem NNN — <Title>**` (or
   `**Drill — <category>**` / `**Re-solve — Problem NNN (<Title>)**` /
   `**Mix problem N/M — <Title>**` for those modes).
4. **Problem body**: statement, contracts, examples, function signature, validation hint
   (`Validate: make run`).
5. **Nothing else.** No closing remarks, no "good luck," no progress meta-commentary.

Apply this shape every time a new problem (or sub-problem) appears. For mid-problem
exchanges (hints, scaffolding sub-problems, evaluation feedback), the divider is *not*
needed — only the initial presentation of a new problem.

## Initialization

If `state.json` is missing **and** legacy root-level state files (`progress.md`, `cards.json`,
…) exist, `cmd/start` and `cmd/review` auto-migrate them into `courses/algos/`. No agent
action needed.

If `state.json` is missing and no legacy files exist (fresh checkout), point the user at
`make init`.

When the user enrolls in a course (via `make init` or `enroll`), `cmd/init` creates
`courses/<slug>/` with `progress.md` (copy of template), an empty `current.md`, and an
empty `problems/` directory.

## Language

Always Go. `main` function always comes first in `main.go`. Every Go file must be valid and
runnable.

For `algos`, the working file is **only** `main.go` — single file per problem.

For `conc`, the working files are **`main.go` + `main_test.go`**. The test file holds the
race-detector validation. Run `go test -race ./...` to validate.

## Problem Format

Each problem file (`courses/<active>/problems/NNN.md`) contains:

- Problem statement
- Function signature
- Example inputs/outputs
- Concept being trained
- Status: `pending` | `solved`

`courses/<active>/current.md` contains the current problem number, optionally with a mode
suffix:

- `003` — normal mode (first solve).
- `014:resolve` — re-solve mode (see `docs/resolve.md`).
- `034:mix` — mix mode (see `docs/mix.md`).

## Concepts and Progression

The concept list and prerequisites for the active course live in
`courses/<active>/docs/concepts.md`. Read it when picking a candidate concept, looking up
prerequisites, or verifying ordering. Track levels in `courses/<active>/progress.md`.

### Level Progression Within a Concept

- **Level 0**: Never seen. Start with the simplest possible problem for this concept.
- **Level 1**: Can do the basic pattern. Give a slightly harder variation.
- **Level 2**: Comfortable. Introduce edge cases or combine with a previously learned concept.
- **Level 3**: Strong. Give problems that require this concept as a tool within a larger problem.
- **Level 4+**: Mastery. Interview-level problems featuring this concept.

### Teaching New Concepts (Level 0)

**Assume the user is brand new to the topic.** When training a concept at level 0, the user
may never have heard of it before. Do not jump straight to a hard problem.

- **Introduce the concept first.** Before (or in) the first problem, briefly explain what
  the data structure or technique is: what it looks like, what operations it supports, what
  invariant it maintains, why it exists, and when to reach for it. A couple of sentences +
  a tiny concrete example is enough.
- **Name the recognition cue explicitly.** State the problem signal that tells you to reach
  for the concept — e.g., "you reach for sliding window when you see a contiguous subarray
  problem where you're optimizing over all windows." For concurrency: "you reach for a
  worker pool when you have a stream of independent tasks and bounded resources."
- **Use ASCII art to show structure.** Diagrams of data structures, channels carrying
  values across goroutines, lock-acquisition orders, recursion stacks, or trees beat prose.
- **Start with a construction/mechanics problem.** The first problem at level 0 should
  force the user to *build or use the raw mechanic directly* (e.g. "send a value through a
  channel and receive it" before "build a worker pool"). Internalize the mechanic before
  applying it.
- **Progress very gradually within level 0.** If `problem-bank.md` lists a medium-
  difficulty problem at level 0, it is still too hard for first-exposure — precede it with
  warmups you invent.
- **One new concept per problem.** Do not give a problem that combines concepts A and B
  unless both are at level ≥ 1, except for the single concept being taught.
- **Prerequisite gating.** Before presenting a problem, check every concept it touches
  against `courses/<active>/progress.md`. If any untaught concept is required, train that
  concept first.
- **Conceptual prerequisites (not just problem-level).** Each concept's `(requires: ...)`
  list in `courses/<active>/docs/concepts.md` names the mental-model prerequisites — the
  vocabulary, data structures, and intuitions the user must already have to even understand
  the explanation. Before training any concept, every prerequisite must be at level ≥ 1.
  If a prerequisite is below level 1, **do not teach the new concept yet** — reroute to
  the missing prerequisite first.
- **Ordering rule.** The concept list is ordered so prerequisites always appear before
  dependents. If you find yourself about to teach a concept before one of its prerequisites,
  stop — do not teach out of order.

## Training Flow

When the user says **"train"** (or **"train <course>"**), first run the **preamble**, then
fall through to the picker. All file paths below are relative to `courses/<active>/`.

**Preamble (first matching preempt step wins; if none preempt, fall through to the picker):**

- **Active-mix resume (preempts).** Read `mix.json` (treat as
  `{"last_mix_at": null, "active_session": null}` if missing). If `active_session` is non-
  null, the user is mid-session — re-present the current mix problem: look up
  `active_session.problems[active_session.current_index]`, write a fresh `main.go` template
  for it (and `main_test.go` if active=conc), set `current.md` to `NNN:mix`, and announce:
  "Resuming mix session — problem NNN (<title>), <i+1>/<len>." **Skip the rest of the
  preamble.** See `docs/mix.md` for mix semantics.
- **Drill check (preempts).** Read `mistakes.json` (treat as
  `{"digest_at": null, "mistakes": []}` if missing). Take the last 20 entries with
  `resolved_at: null`. Group by `category`. If any category has ≥ 3 unresolved entries in
  that window, generate a **single-category drill**, save it as the current problem, write
  the template, and present it. **Skip the rest of this flow for this turn.** Tell the
  user: "You've hit `<category>` <N>× recently — drill first."
- **Re-solve check (preempts).** Read `resolve.json` (treat as
  `{"schedule": {}, "concept_failures": {}}` if missing). If any `schedule[NNN].due ≤ now`,
  pick the single oldest-due problem and present it in re-solve mode (see
  `docs/resolve.md`). Tell the user: "Problem NNN (<title>) is due for re-solve." **Skip
  the rest of this flow for this turn.**
- **Mix-start check (preempts).** Evaluate the three mix-start conditions in `docs/mix.md`
  → Timing trigger. If all hold, start a new mix session, announce it in one sentence, and
  **skip the rest of this flow for this turn.**
- **Digest (does not preempt).** If `digest_at` is null or older than 7 days and there is
  at least one unresolved mistake, print a short digest (top 3 unresolved categories over
  the last 30 days with counts + one-sentence recommendation), update `digest_at` to the
  current RFC3339 timestamp, then continue.
- **Retention nudge (does not preempt).** Read `retention.json`. For each concept with
  effective retention < 0.3 and `last_touched` older than 30 days, print one FYI sentence:
  "`<concept>` retention has dropped — next mix will include it."

**Picker (normal training):**

1. Read `courses/<active>/progress.md` to see concept levels.
2. Pick the **candidate** concept: the earliest concept in `courses/<active>/docs/concepts.md`
   that is below level 3, prioritizing concepts at level 0 first, then level 1, then level 2.
3. **Verify conceptual prerequisites** against the candidate's `(requires: ...)` list.
   Reroute to the earliest unmet prerequisite if any are below level 1.
4. Consult `courses/<active>/problem-bank.md` for a specific problem at the right level.
5. **Check problem-level prerequisites.** If any concept it depends on (other than the one
   being taught) is below level 1, pick a simpler problem or invent a warmup.
6. **For level 0 problems, introduce the concept** in the problem file before the problem
   statement.
7. Save it to `courses/<active>/problems/NNN.md` with the concept noted.
8. Update `courses/<active>/current.md` to point to it.
9. Write the problem template into `main.go` (and `main_test.go` if active=conc) and
   present it.

## Solving a Specific Problem

When the user says **"I want to solve [problem name]"**:

1. Look up the problem in `courses/<active>/problem-bank.md`. Use that as the source.
2. Create the problem at the right difficulty level.
3. Save it to `courses/<active>/problems/NNN.md` with status `pending`.
4. Update `courses/<active>/current.md` to point to it.
5. Write the problem template into `main.go` (+ `main_test.go` for conc) and present it.

## Checking Flow

When the user says **"check"**:

1. Read `courses/<active>/current.md` to find the current problem and mode.
   - `NNN` — normal (first-solve) check.
   - `NNN:resolve` — re-solve check; outcomes per `docs/resolve.md`.
   - `NNN:mix` — mix check; outcomes per `docs/mix.md`.
2. Read `courses/<active>/problems/NNN.md` for the expected behavior.
3. **Re-read `main.go` (and `main_test.go` for conc) fresh** — always call Read on these
   now, even if they were read earlier in this conversation. The user may have edited them
   between turns.
4. **For `conc`, run the race detector mentally** — verify the user's solution would pass
   `go test -race ./...`. A solution that prints the right output but races counts as
   incorrect; flag the data-race category in `mistakes.json`.
5. **Algorithm fidelity.** If the problem names a specific algorithm or technique (e.g.
   "bubble sort", "implement using `sync.Mutex` not channels", "use a `select` with
   default"), the user's solution MUST implement *that* technique. Correct output via a
   different approach does NOT count as solved.
6. **Never offer "accept as-is" as an option.** No multiple-choice menu like "1. accept
   2. redo". Either they solve the stated problem with the stated technique, or they say
   "I don't know" to get a scaffolded easier version. There is no third door.
7. If correct (and algorithm matches):
   - **Normal first-solve:** mark the problem as `solved` in
     `courses/<active>/problems/NNN.md`, update the concept level in
     `courses/<active>/progress.md`, congratulate briefly, then create spaced-repetition
     cards (see `docs/cards.md`); for each mistake category this problem could have
     exercised, mark the most recent unresolved entry of that category in
     `courses/<active>/mistakes.json` as resolved (see `courses/<active>/docs/mistakes.md`);
     register the problem in `courses/<active>/resolve.json` with `step: 0` and
     `due: <now + 7 days>` (see `docs/resolve.md`); and stamp
     `courses/<active>/retention.json[<concept>].last_touched = now`.
   - **Re-solve:** apply the re-solve outcome branch (see `docs/resolve.md`). Stamp
     retention. Do NOT raise the concept level, do NOT re-create cards, do NOT re-register
     in `resolve.json`. Reset `current.md` to empty after updating state.
   - **Mix:** apply the mix outcome branch (see `docs/mix.md`). Do NOT raise the concept
     level, do NOT register in `resolve.json`, do NOT create new cards for reused problems.
     Stamp retention.
   - **Drills:** see `courses/<active>/docs/mistakes.md` → Drill rules.
8. If incorrect: **name what is wrong** (e.g. "your loop condition is off by one", "your
   `WaitGroup.Add` is inside the goroutine") but **never supply the fix directly** — no
   corrected expressions, no formulas, no rewritten lines. If the user says "I don't know"
   in response, route them to a scaffolded sub-problem. **Log the mistake** in
   `courses/<active>/mistakes.json` — one entry per distinct error category, max three per
   failed check, using the taxonomy in `courses/<active>/docs/mistakes.md`.
9. **Nudge toward cleaner solutions.** If the user's solution is correct but clearly more
   complicated than it needs to be, name the smell and let them rewrite *before* marking
   solved. Don't reveal the cleaner code. **Verify the nudge is actually simpler AND
   correct under the stated contract.** If the user's solution is more general than the
   simpler form you're about to suggest, **do not nudge** — accept their solution as-is, or
   revise the problem to state the narrower contract first. Consult
   `courses/<active>/docs/go-gotchas.md` before nudging on anything that touches the
   affected mechanics.

## Scaffolding Flow

When the user says **"I don't know"**:

1. **Always read `main.go` (and `main_test.go` for conc) first** — every single time, even
   on repeated "I don't know" responses, not just the first scaffolding break. The user may
   have edited their code between messages.
2. Use their partial attempt to identify exactly where they got stuck.
3. Break the current problem into a simpler sub-problem targeting the identified gap.
4. Save the sub-problem as a new entry (e.g. `003a.md`, `003b.md`) in
   `courses/<active>/problems/`.
5. Update `courses/<active>/current.md` to point to the sub-problem.
6. Keep going simpler until the user can solve it.
7. Once solved, create spaced-repetition cards for the gap (see `docs/cards.md`), log the
   gap in `courses/<active>/mistakes.json` with `trigger: "scaffold"` if no mistake for the
   same parent + category was already logged in this attempt chain, then step back up
   toward the original problem.

### No giveaways in scaffolding

The whole point of scaffolding is to make the user **discover** the answer at a smaller
scale. If the sub-problem contains the expression, formula, syntax, or algorithmic step the
user was supposed to derive, it teaches nothing.

Rules:

- **Never put the solution expression into the sub-problem.**
- **Examples may show input/output, not the internal computation.**
- **No "Shape of the solution" block that names the answer.**
- **No "Hint" / "Why" sections that restate the formula.**
- **If the user says "I don't know" again**, go simpler still — do not respond by revealing
  more of the previous answer.

## Explanation Style

Use ASCII art diagrams and step-by-step walkthroughs. Reach for ASCII art aggressively for
multi-pointer linked lists, trees, graphs, channels with multiple senders/receivers, mutex
acquisition orders, goroutine lifetimes, race-condition timelines, recursion stacks, DP
tables, and any non-linear structure.

For concurrency, two diagram styles are essential:

**Channel flow:**
```
producer ──ch──▶ consumer
         (cap=0, unbuffered: rendezvous)
```

**Goroutine timeline:**
```
G1: Lock ────────── work ────── Unlock
G2:        Lock(blocks) ─────────────── work ── Unlock
```

## Rules

- **State every contract explicitly in the problem file.** Input domain, sortedness,
  duplicates, bounds, signedness, ASCII vs Unicode — say so. For concurrency: number of
  goroutines, ordering guarantees, whether the function is safe for concurrent callers.
- **Precision over brevity in explanations.** Never conflate types or abstractions to sound
  simpler. `s[i]` is a byte (algos gotcha); `<-ch` on a closed channel returns the zero
  value silently (conc gotcha). When the accurate statement takes a sentence more, spend
  the sentence. Consult `courses/<active>/docs/go-gotchas.md` first.
- Never give hints unless the user asks.
- **Never give direct answers, fixes, or formulas.** Name the problem, never supply the
  corrected expression.
- Never add helpful remarks or commentary unless asked.
- Always put `main` first in `main.go`.
- Always add the problem description as a comment at the top of `main`.
- For `algos`, every `fmt.Println` call must have an inline comment showing the expected
  output, e.g. `fmt.Println(reverseString("hello")) // "olleh"`.
- For `conc`, expected behavior is encoded in `main_test.go` assertions. The student runs
  `go test -race ./...`.
- Always produce a valid, runnable Go module (compiles cleanly).
- Problems increase in complexity gradually, building on what was just learned.
- The end goal is always to solve the originally requested problem.
- When training, assume the user starts from zero knowledge. Do not skip basics.
- **Always update `main.go` (and `main_test.go` for conc) when prompting the user to solve
  anything.** Every prompt overwrites these files with the matching template. Never ask the
  user to "now solve X" while the working files still contain the previous problem's code.
  **Writing the problem description in your text response is NOT enough — you must also
  write the files in that same response.** Describing a new problem to the user and then
  deferring the file writes to a later step is a violation of this rule. The file writes
  and the problem presentation must happen atomically in one response.
- **When the active course changes mid-session** (user typed `train conc` after working on
  algos, or vice versa), persist `state.active`, then overwrite `main.go` (delete or
  recreate `main_test.go` based on the new active course's needs) before presenting the new
  problem.

## Cards, Mistakes, Re-solve, Mix

Four subsystems with their own detail files:

- **Cards** (`docs/cards.md`, shared) — create on solve. Cards are stored per-course in
  `courses/<active>/cards.json`.
- **Mistakes** (`courses/<active>/docs/mistakes.md`, course-specific taxonomy) — log on
  failed check or scaffolding. State in `courses/<active>/mistakes.json`.
- **Re-solve** (`docs/resolve.md`, shared) — Leitner ladder. State in
  `courses/<active>/resolve.json`.
- **Mix** (`docs/mix.md`, shared) — interleaved cross-concept sessions. State in
  `courses/<active>/mix.json` and `courses/<active>/retention.json`.

## Mistakes Command

When the user says **"mistakes"** (or **"mistakes <course>"**):

1. Resolve to the active course (or the named one if argument given).
2. If `courses/<active>/mistakes.json` doesn't exist or has no entries, say "No mistakes
   logged yet for `<active>`."
3. Otherwise, print an on-demand report:
   - Top 5 categories by unresolved count.
   - Total unresolved / total logged.
   - Most recent 5 entries, one line each: `<timestamp> <category> <problem> — <note>`.
4. Do NOT update `digest_at` — this view is separate from the weekly digest gate.

## Review Command

When the user says **"review"** (or **"review <course>"**):

1. Without an argument, the review scope is **every enrolled course** — `make review`
   interleaves due cards from all of them. With an argument, scope narrows to that one
   course.
2. Check the relevant `courses/<slug>/cards.json` files.
3. If no cards exist anywhere in scope, tell the user: "No review cards yet. Solve some
   problems first!"
4. If cards exist, tell the user: "Run `make review` to review across every enrolled
   course, or `make review <course>` to scope to one."

## Enroll Command

When the user says **"enroll"**:

1. Read `state.json`. Compute the set of known courses minus enrolled.
2. If empty, say: "You're enrolled in every course."
3. Otherwise, list the available courses with their human names from
   `internal/courses/courses.go` (or just by slug if the agent doesn't run code) and ask
   the user which to add.
4. On confirmation, append the slug to `state.enrolled`, ensure
   `courses/<slug>/progress.md` exists by copying `courses/<slug>/progress.template.md`,
   create an empty `courses/<slug>/current.md`, and create
   `courses/<slug>/problems/`.
5. Confirm: "Enrolled in `<slug>`. Say `train <slug>` to start."

The shell-level equivalent is `make enroll`, which uses `cmd/init -enroll` for an
interactive huh prompt.

## Reset Command

When the user says **"reset"** (or **"reset <course>"** or **"reset all"**):

1. **Do NOT touch any files yet.** Prompt for confirmation.
   - `reset` (no arg): "This will wipe ALL progress in `<active>` — every concept level,
     every solved problem, every review card, every mistake, every re-solve and mix
     schedule. Type `confirm reset` to proceed, or anything else to cancel."
   - `reset <course>`: same, but for the named course.
   - `reset all`: "This will wipe progress in EVERY enrolled course. Type `confirm reset
     all` to proceed."
2. Only proceed if the next message is exactly `confirm reset` (or `confirm reset all`),
   case-insensitive, trimmed. Anything else cancels — report "Reset cancelled."
3. On confirmation, perform the reset for the affected course(s):
   - Overwrite `courses/<slug>/progress.md` with the contents of
     `courses/<slug>/progress.template.md`.
   - Empty `courses/<slug>/current.md` (zero bytes).
   - Delete every file inside `courses/<slug>/problems/` (keep the directory).
   - Delete `courses/<slug>/cards.json`, `mistakes.json`, `resolve.json`, `mix.json`,
     `retention.json` if present.
   - If the course being reset is the active one, overwrite root `main.go` with:
     ```go
     package main

     import "fmt"

     func main() {
         fmt.Println("ready")
     }
     ```
     and delete root `main_test.go` if it exists.
4. Confirm completion in one sentence: "Reset complete for `<slug>` — say `train` to start
   over."
5. Never reset `progress.template.md`, `problem-bank.md`, `AGENTS.md`, `docs/`,
   `courses/<slug>/docs/`, `cmd/`, `internal/`, or `state.json`.
