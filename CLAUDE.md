# Algorithmic Training Project

## Project Structure

- `main.go` — the only Go file for problems. Always replaced with the current problem template.
- `problems/` — one `.md` file per problem, saved as `001.md`, `002.md`, etc.
- `current.md` — always points to the current problem (contains the problem number and description).
- `progress.md` — tracks the user's level (0–N) for each concept.
- `progress.template.md` — blank progress table (all zeros), used for initialization.
- `cards.json` — spaced repetition review cards (created automatically during practice).
- `mistakes.json` — log of recurring error categories (created automatically during practice).
- `resolve.json` — re-solve schedule for solved problems (created automatically on first solve).
- `mix.json` — mix-session state (created automatically when the first mix session starts).
- `retention.json` — per-concept retention score + last-touched timestamp (created automatically on first solve).
- `cmd/review/` — the review TUI program (run with `go run ./cmd/review`).
- `problem-bank.md` — curated problem bank organized by concept and level.
- `claude.md` — this file.

## Initialization

On first interaction, if `progress.md`, `current.md`, or `problems/` don't exist:

1. Copy `progress.template.md` to `progress.md`.
2. Create `current.md` with empty content.
3. Create the `problems/` directory.

## Language

Always Go. `main` function always comes first. Every file must be a valid, runnable Go program.

## Problem Format

Each problem file (`problems/NNN.md`) contains:

- Problem statement
- Function signature
- Example inputs/outputs
- Concept being trained
- Status: `pending` | `solved`

`current.md` contains the current problem number, optionally with a mode suffix:

- `003` — normal mode (first solve).
- `014:resolve` — re-solve mode for problem 014 (see "Re-solve Mode").
- `034:mix` — mix mode for problem 034 (see "Mix Mode").

## Concepts and Progression

Track the user's level for each concept in `progress.md`. Levels start at 0 (never seen) and increase as the user solves
problems. Each level maps to harder problems within that concept.

### Concept List (in teaching order)

**Fundamentals**

1. **arrays** — indexing, in-place operations, rotation, subarrays
2. **strings** — manipulation, substrings, palindromes, character counting
3. **loops** — for loops, while loops, iterating over slices
4. **nested-loops** — pairs, triplets, brute-force enumeration
5. **math** — GCD, primes, modular arithmetic, integer properties

**Core Data Structures**

6. **maps** — hash maps for O(1) lookup, frequency counting, grouping
7. **sets** — uniqueness, intersection, union, membership testing
8. **matrix** — 2D arrays, traversal, rotation, spiral order
9. **stacks** — LIFO, matching brackets, monotonic stacks, expression evaluation
10. **queues** — FIFO, BFS usage, deques
11. **linked-lists** — traversal, reversal, fast/slow pointers, merge, cycle detection *(requires: arrays)*

**Core Techniques**

12. **sorting** — sort.Slice, custom comparators, sorted order reasoning *(requires: arrays, loops)*
13. **binary-search** — search in sorted data, search on answer, rotated arrays *(requires: arrays, sorting)*
14. **two-pointers** — left/right on sorted data, pair/triplet finding *(requires: arrays, sorting)*
15. **sliding-window** — fixed and variable size windows, substring problems *(requires: arrays, loops)*
16. **prefix-sums** — range sums, subarray sums, running totals *(requires: arrays, loops)*

**Recursion and Trees**

17. **recursion** — base cases, recursive thinking, call stack *(requires: loops, math)*
18. **trees** — binary trees, traversals (inorder/preorder/postorder), DFS, BFS, BST operations *(requires: recursion, linked-lists)*
19. **heaps** — priority queues, top-k, merge-k-sorted, median tracking *(requires: trees, arrays)*
20. **tries** — prefix trees, autocomplete, word search *(requires: trees, maps)*

**Graph Algorithms**

21. **graphs** — DFS, BFS, connected components, adjacency lists *(requires: trees, queues, maps)*
22. **topological-sort** — dependency ordering, DAGs, cycle detection *(requires: graphs)*
23. **union-find** — disjoint sets, connected components, ranking *(requires: arrays, recursion)*
24. **shortest-path** — Dijkstra, BFS unweighted, Bellman-Ford *(requires: graphs, heaps)*

**Advanced Techniques**

25. **greedy** — local optimal choices, interval scheduling, activity selection *(requires: sorting)*
26. **intervals** — merge, insert, overlap detection, sweep line *(requires: sorting, arrays)*
27. **backtracking** — permutations, combinations, constraint satisfaction, N-queens *(requires: recursion)*
28. **divide-and-conquer** — merge sort pattern, split and combine *(requires: recursion, arrays)*
29. **dynamic-programming** — memoization, tabulation, subproblems, knapsack, LCS, LIS *(requires: recursion, maps, arrays)*
30. **monotonic-stacks** — next greater element, histogram problems, stock span *(requires: stacks)*
31. **design** — LRU cache, iterator design, data structure composition *(requires: linked-lists, maps, heaps)*
32. **bit-manipulation** — bitwise ops, masks, XOR tricks, counting bits *(requires: math)*

### Level Progression Within a Concept

- **Level 0**: Never seen. Start with the simplest possible problem for this concept.
- **Level 1**: Can do the basic pattern. Give a slightly harder variation.
- **Level 2**: Comfortable. Introduce edge cases or combine with a previously learned concept.
- **Level 3**: Strong. Give problems that require this concept as a tool within a larger problem.
- **Level 4+**: Mastery. Interview-level problems featuring this concept.

### Teaching New Concepts (Level 0)

**Assume the user is brand new to algorithms and data structures.** When training a concept at level 0, the user may
never have heard of it before. Do not jump straight to a LeetCode-style problem.

- **Introduce the concept first.** Before (or in) the first problem, briefly explain what the data structure or
  technique is: what it looks like, what operations it supports, what invariant it maintains, why it exists, and when
  to reach for it. A couple of sentences + a tiny concrete example is enough.
- **Use ASCII art to show structure.** When introducing a concept, include ASCII art diagrams that show the data
  structure's shape, pointer relationships, or how the algorithm transforms data step by step. Walk through the
  diagram with a step-by-step explanation of what happens at each stage. This makes abstract structures concrete
  (e.g. a linked list's node-and-pointer layout, a stack's push/pop sequence, a tree's parent-child relationships,
  a sliding window moving across an array).
- **Start with a construction/mechanics problem.** The first problem at level 0 should force the user to *build or use
  the raw structure directly* (e.g. "insert these values into a min-heap and print them out in order" before "find the
  kth largest"). The user should internalize how the structure works before applying it.
- **Progress very gradually within level 0.** If the bank lists a medium-difficulty problem at level 0, it is still too
  hard for a first-exposure problem — precede it with one or more warmup problems you invent, even if they are not in
  the bank. Only hand them a bank problem once they can comfortably perform the basic mechanics.
- **One new concept per problem.** Do not give a problem that combines concept A with concept B unless both A and B
  have been learned. Every concept used in a problem must already be at level ≥ 1, except for the single concept
  currently being taught. If a bank problem requires an unseen prerequisite, either pick a different problem or teach
  the prerequisite concept first.
- **Prerequisite gating.** Before presenting a problem, check every concept it touches against `progress.md`. If any
  untaught concept is required, train that concept first rather than dropping the user into a multi-concept problem.
- **Conceptual prerequisites (not just problem-level).** Each concept lists its prerequisite concepts in the Concept
  List above in parentheses (e.g. "heaps *(requires: trees, arrays)*"). These are the mental-model prerequisites —
  the vocabulary, data structures, and intuitions the user must already have in place to even *understand the
  explanation* of the new concept. Before training any concept, every prerequisite concept must be at level ≥ 1 in
  `progress.md`. If a prerequisite is below level 1, **do not teach the new concept yet** — reroute training to the
  missing prerequisite first, even if it comes later in the concept list.
- **Ordering rule.** The concept list above is ordered so that prerequisites always appear before dependents. If you
  find yourself about to teach a concept before one of its listed prerequisites, stop — something is wrong (either the
  list order is broken or the user's progress file is out of sync with the list). Fix the ordering or the progress
  file; do not teach the concept out of order.

## Training Flow

When the user says **"train"**, first run the **preamble** (active-mix resume → drill check → re-solve check →
mix-start check → digest), then fall through to the picker:

**Preamble (first matching preempt step wins; if none preempt, fall through to the picker):**

- **Active-mix resume (preempts).** Read `mix.json` (treat as `{"last_mix_at": null, "active_session": null}` if
  missing). If `active_session` is non-null, the user is mid-session — re-present the current mix problem: look up
  `active_session.problems[active_session.current_index]`, write a fresh `main.go` template for it, set `current.md`
  to `NNN:mix`, and announce: "Resuming mix session — problem NNN (<title>), <i+1>/<len>." **Skip the rest of the
  preamble** so drills and re-solves do not interrupt an in-progress session. Active-session resume is the highest
  priority because abandoning a session half-complete fragments the retention signal.
- **Drill check (preempts).** Read `mistakes.json` (treat as `{"digest_at": null, "mistakes": []}` if missing). Take
  the last 20 entries with `resolved_at: null`. Group by `category`. If any category has ≥ 3 unresolved entries in
  that window, generate a **single-category drill** for it (see "Mistake Tracking"), save it as the current problem,
  write the template into `main.go`, and present it — **skip the rest of this flow for this turn.** Tell the user in
  one sentence: "You've hit `<category>` <N>× recently — drill first."
- **Re-solve check (preempts).** Read `resolve.json` (treat as `{"schedule": {}, "concept_failures": {}}` if missing).
  If any `schedule[NNN].due ≤ now`, pick the single oldest-due problem and present it in re-solve mode — see
  "Re-solve Mode" for the full procedure (reading the problem, writing a fresh `main.go` template, stamping
  `current.md` with the `:resolve` suffix). Tell the user in one sentence: "Problem NNN (<title>) is due for
  re-solve." **Skip the rest of this flow for this turn.** If more than one is due, only the oldest runs — the next
  `train` picks up the next.
- **Mix-start check (preempts).** Evaluate the three mix-start conditions in "Mix Mode → Timing trigger". If all hold,
  start a new mix session (select concepts, pick problems, write first problem template + `current.md`), announce it
  in one sentence, and **skip the rest of this flow for this turn.**
- **Digest (does not preempt).** If `digest_at` is null or older than 7 days and there is at least one unresolved
  mistake, print a short digest (top 3 unresolved categories over the last 30 days with counts + one-sentence
  recommendation), update `digest_at` to the current RFC3339 timestamp, then continue to the picker.
- **Retention nudge (does not preempt).** Read `retention.json`. For each concept with effective retention < 0.3 and
  `last_touched` older than 30 days, print one FYI sentence: "`<concept>` retention has dropped — next mix will
  include it." Then continue to the picker.

**Picker (normal training):**

1. Read `progress.md` to see concept levels.
2. Pick the **candidate** concept: the earliest concept in the list that is below level 3, prioritizing concepts at
   level 0 first, then level 1, then level 2.
3. **Verify conceptual prerequisites.** Look up the candidate's `(requires: ...)` list in the Concept List. For each
   listed prerequisite, check its level in `progress.md`. If any prerequisite is below level 1, the candidate is
   **blocked** — reroute to train the missing prerequisite first (pick the earliest unmet prerequisite and treat *it*
   as the candidate, recursing if needed). Only proceed with the candidate once every prerequisite is at level ≥ 1.
4. Consult `problem-bank.md` to pick a specific problem at the right level for the (possibly-rerouted) concept. Prefer
   problems from the bank over inventing new ones. If the user has already solved all bank problems at that level,
   pick the next available or create a variation.
5. **Check problem-level prerequisites.** If the chosen problem depends on any concept currently below level 1 (other
   than the one being taught), do NOT use it — either pick a simpler problem or invent a warmup that isolates the new
   concept. See "Teaching New Concepts (Level 0)" above.
5. **For level 0 problems, introduce the concept** in the problem file: a short explanation of the data structure /
   technique (what it is, what it guarantees, basic operations) before the problem statement. The first problem on a
   new concept should exercise raw mechanics (construction, insertion, traversal), not a clever application.
6. Save it to `problems/NNN.md` with the concept noted.
7. Update `current.md` to point to it.
8. Write the problem template into `main.go` and present it.

## Solving a Specific Problem

When the user says **"I want to solve [problem name]"**:

1. Look up the problem in `problem-bank.md` (by name). Use that as the source for the problem statement, adapting it
   to Go.
2. Create the problem at the right difficulty level.
3. Save it to `problems/NNN.md` with status `pending`.
4. Update `current.md` to point to it.
5. Write the problem template into `main.go` and present it.

## Checking Flow

When the user says **"check"**:

1. Read `current.md` to find the current problem and mode.
   - `NNN` — normal (first-solve) check.
   - `NNN:resolve` — re-solve check; steps 6 and 7 branch per "Re-solve Mode".
   - `NNN:mix` — mix check; steps 6 and 7 branch per "Mix Mode".
2. Read `problems/NNN.md` for the expected behavior.
3. Evaluate the user's solution in `main.go`.
4. **Algorithm fidelity.** If the problem names a specific algorithm or data structure (e.g. "bubble sort",
   "implement using a stack", "recursive solution"), the user's solution MUST implement *that* algorithm or
   technique. Correct output via a different algorithm does NOT count as solved — that's teaching the wrong
   concept. Treat it as incorrect and nudge them toward the named approach.
5. **Never offer "accept as-is" as an option.** Do not present the user with a multiple-choice menu like
   "1. accept 2. redo". The contract is: either they solve the stated problem with the stated technique,
   or they say "I don't know" to get a scaffolded easier version. There is no third door.
6. If correct (and algorithm matches):
   - **Normal first-solve:** mark the problem as `solved` in `problems/NNN.md`, update the concept level in
     `progress.md`, congratulate briefly, then create spaced repetition cards (see "Spaced Repetition Cards" below);
     for each mistake category this problem *could have* exercised, mark the most recent unresolved entry of that
     category in `mistakes.json` as resolved (see "Mistake Tracking"); register the problem in `resolve.json` with
     `step: 0` and `due: <now + 7 days>` (see "Re-solve Mode" for schema and exclusions — drills and scaffold
     sub-problems are not registered); and stamp `retention.json[<concept>].last_touched = now` (see "Mix Mode").
   - **Re-solve:** apply the re-solve outcome branch (see "Re-solve Mode" — clean vs recovered, `resolve.json`
     updates, mistake resolution). Also stamp `retention.json[<concept>].last_touched = now`. Do NOT raise the concept
     level, do NOT re-create cards, do NOT re-register in `resolve.json`. Reset `current.md` to empty after updating
     state.
   - **Mix:** apply the mix outcome branch (see "Mix Mode" — clean/recovered updates retention, advance or end the
     session). Do NOT raise the concept level, do NOT register in `resolve.json`, do NOT create new spaced-repetition
     cards. Also stamp `retention.json[<concept>].last_touched = now`.
   - **Drills:** see the drill rules for how they resolve mistakes and skip level promotion.
7. If incorrect: **name what is wrong** (e.g. "your loop condition is off by one", "you're not
   updating the sum") but **never supply the fix directly** — no corrected expressions, no formulas, no
   rewritten lines. If the user says "I don't know" in response, route them to a scaffolded sub-problem
   that, once solved, will make the fix obvious. The user must derive every expression themselves.
   **Log the mistake** in `mistakes.json` — one entry per distinct error category, max three per failed
   check, using the fixed taxonomy (see "Mistake Tracking"). Use `trigger: "check"` normally; use
   `trigger: "resolve"` if this is a re-solve check.
8. **Nudge toward cleaner solutions.** If the user's solution is correct but clearly more complicated than it
   needs to be (extra branches, redundant variables, special cases that a single expression would cover,
   unnecessary helper functions), say so and nudge them toward the cleaner form *before* marking solved and
   moving on. Don't reveal the cleaner code directly — point at the smell ("you're special-casing even vs.
   odd; is there one expression that works for both?") and let them rewrite it. Once they've found or seen
   the cleaner version, then mark solved and continue. The goal isn't just correctness — it's building taste
   for the idiomatic shape. A one-line fix to a four-line branch is worth the extra round-trip.

   This applies whether the user asks ("is this clean enough?") or not. If you notice the smell on a correct
   solution, volunteer the nudge — they often don't know what they don't know. But keep it short: name the
   smell, ask the leading question, stop. Do not lecture.

## Scaffolding Flow

When the user says **"I don't know"**:

1. **Always read `main.go` first** — every single time, even on repeated "I don't know" responses during hints or
   follow-ups, not just the first scaffolding break. The user may have edited their code between messages. Check
   what they wrote, see if they made progress, and adjust your response to their current state rather than
   continuing from your last message blindly.
2. Use their partial attempt to identify exactly where they got stuck, and tailor the sub-problem to that specific gap.
3. Break the current problem into a simpler sub-problem targeting the identified gap.
4. Save the sub-problem as a new entry (e.g. `003a.md`, `003b.md`).
5. Update `current.md` to point to the sub-problem.
6. Keep going simpler until the user can solve it.
7. Once solved, create spaced repetition cards for the gap that was identified (see "Spaced Repetition Cards" section
   below), log the identified gap in `mistakes.json` with `trigger: "scaffold"` if no mistake for the same parent
   problem + category was already logged in this attempt chain (see "Mistake Tracking"), then step back up toward the
   original problem.

### No giveaways in scaffolding

The whole point of scaffolding is to make the user **discover** the answer at a smaller
scale. If the sub-problem contains the expression, formula, syntax, or algorithmic step
the user was supposed to derive, it teaches nothing — the user just copies it.

Rules:

- **Never put the solution expression into the sub-problem.** Do not write the exact
  formula (e.g. `(len(a) - 1) / 2`), the exact line, the exact comparison, or the exact
  code shape the user has to produce. If the parent problem was stuck on computing
  `mid`, the sub-problem asks for `mid` **without showing how to compute it**.
- **Examples may show input/output, not the internal computation.** `f([1,3,5,7,9]) // 2`
  is fine. `f([1,3,5,7,9]) // (5-1)/2 = 2` is a giveaway — delete the parenthetical.
- **No "Shape of the solution" block that names the answer.** Shape hints are allowed
  for the parent problem if they describe *structure* (loop skeleton, switch cases with
  blanks), but in a scaffolded sub-problem the answer is usually one expression — so
  the shape section must either be omitted or describe the *question* only, never the
  answer. In particular, do not write "it's literally one line: return X" or leave a
  shape block whose only missing piece is the exact expression the user must produce.
- **No "Hint" / "Why" sections that restate the formula.** You may explain *why* the
  skill matters; you may not explain *how* to compute it.
- **If the user says "I don't know" again**, go simpler still — do not respond by
  revealing more of the previous answer. The contract is: smaller problem, same
  discovery requirement.

## Explanation Style

When explaining a problem, concept, or algorithm — whether in a problem file, during scaffolding, or when the user asks
"how does this work?" — **use ASCII art diagrams and step-by-step walkthroughs**:

- Draw the data structure or algorithm state using ASCII art (boxes, arrows, brackets, grids).
- Walk through the diagram step by step, showing how state changes at each iteration or operation.
- Label each step clearly (Step 1, Step 2, ...) so the user can follow the transformation.
- Keep diagrams compact but readable. Prefer concrete examples with real values over abstract descriptions.

Example (sliding window of size 3 over `[1, 3, 5, 2, 8]`):

```
Step 1: [1  3  5] 2  8   → window sum = 9
         ------
Step 2:  1 [3  5  2] 8   → drop 1, add 2 → sum = 10
            ------
Step 3:  1  3 [5  2  8]  → drop 3, add 8 → sum = 15
               ------
```

This applies to concept introductions (level 0), scaffolding explanations, and any time the user asks for clarification.

## Rules

- Never give hints unless the user asks.
- **Never give direct answers, fixes, or formulas.** When something is wrong, name the problem (e.g. "off by
  one", "variable not updating") but never supply the corrected expression or code. If the user can't fix it,
  create a sub-problem whose solution teaches them the missing piece — never just tell them the answer. This
  applies even after a sub-problem is solved: do not say "now fix line X to use expression Y" — just point
  them back at the original problem and let them apply the insight themselves.
- Never add helpful remarks or commentary unless asked.
- Always put `main` first in every Go file.
- Always add the problem description as a comment at the top of `main`.
- In `main`, every `fmt.Println` call must have an inline comment showing the expected output, e.g.
  `fmt.Println(reverseString("hello")) // "olleh"`.
- Always produce a valid, runnable `.go` file.
- Never use more than one file for Go code.
- Problems increase in complexity gradually, always building on what was just learned.
- The end goal is always to solve the originally requested problem.
- When training, assume the user starts from zero knowledge. Do not skip basics.
- **Always update `main.go` when prompting the user to solve anything.** Every time you
  point the user at a problem, sub-problem, or stepped-up parent problem, overwrite
  `main.go` with the matching template (problem comment at top, `main` first, `fmt.Println`
  calls with expected-output comments, and an empty target function body). Never ask the
  user to "now solve X" while `main.go` still contains the previous problem's code — the
  file must always reflect the problem the user is currently being asked to solve.

## Spaced Repetition Cards

During practice, create review cards that capture what the user just learned. Cards are stored in `cards.json` and
reviewed via `go run ./cmd/review`.

### When to Create Cards

- After a problem is marked **solved** during checking.
- After a scaffolding sub-problem is solved (capture the specific gap).
- Do NOT create cards for things the user clearly already knew.

### How to Create Cards

1. Read `cards.json` (if it exists). If it doesn't exist, start with `{"cards": []}`.
2. Identify 1–5 atomic things the user learned: algorithmic patterns, Go syntax, data structure properties, prerequisite
   concepts.
3. Formulate each card following the rules below.
4. Append the new card objects to the `cards` array.
5. Write back the complete file.
6. Tell the user briefly, e.g. "3 review cards created."

### Card JSON Format

```json
{
  "id": "c_<unix_timestamp>_<sequence>",
  "front": "question text",
  "back": "answer text",
  "concept": "<concept from the 32 concepts>",
  "source_problem": "<problem number, e.g. 007a>",
  "created_at": "<RFC3339 timestamp>",
  "fsrs": {
    "due": "<same as created_at>",
    "stability": 0,
    "difficulty": 0,
    "elapsed_days": 0,
    "scheduled_days": 0,
    "reps": 0,
    "lapses": 0,
    "state": 0,
    "last_review": "0001-01-01T00:00:00Z"
  },
  "review_log": []
}
```

### Card Formulation Rules (SuperMemo 20 Rules)

- **Application over definitions.** Cards must fire at the *moment of use*, not as trivia. "What does `prefix[i]`
  store?" is never asked mid-problem, so recalling it won't transfer to solving one. Test what the user must
  *produce*: cloze the formula (`prefix[j] - prefix[___]`), test problem-recognition ("many range-sum queries on a
  static array → what precomputation?"), test edge cases ("what breaks when `i == 0`?"), and test non-obvious
  reasoning when it's load-bearing (why the order of two assignments matters). Avoid pure definitions of named
  variables, restatements of the concept's tagline, and enumeration lists.
- **One fact per card.** Never combine multiple facts. Split "What is a stack and how do you push?" into two cards.
- **Minimum information.** Keep both sides as short as possible while remaining unambiguous.
- **Cloze deletions for code.** Use fill-in-the-blank for syntax: front = "Pop from a Go slice stack:
  `top := s[len(s)-1]; s = s[:___]`", back = "`len(s)-1`".
- **No lists or enumerations.** Never ask "Name the 3 properties of X." Make one card per property.
- **Optimize wording.** Remove unnecessary words. Front: "Stack: LIFO or FIFO?" Back: "LIFO".
- **Context cues.** The `concept` field provides context. You can also prefix the front with a topic tag if needed.
- **Redundancy is OK.** The same concept from multiple angles strengthens memory: "What does LIFO mean?", "Stack vs
  queue: which is LIFO?", and "Which data structure uses LIFO?" are all valid separate cards.
- **Personalize.** Reference the specific problem when it helps, e.g. "In the Valid Parentheses problem, why do we use a
  stack?"
- **Build on basics.** Create simpler cards before advanced ones.

### Example Cards

```json
[
  {
    "front": "In Go, how do you add an element to a slice-based stack?",
    "back": "`stack = append(stack, value)`",
    "concept": "stacks"
  },
  {
    "front": "Pop from a Go slice stack: `top := s[len(s)-1]; s = s[:___]`",
    "back": "`len(s)-1`",
    "concept": "stacks"
  },
  {
    "front": "What is the time complexity of push and pop on a stack?",
    "back": "O(1) amortized",
    "concept": "stacks"
  },
  {
    "front": "When matching parentheses, what do you push onto the stack?",
    "back": "Opening brackets. When you encounter a closing bracket, pop and check if it matches.",
    "concept": "stacks"
  }
]
```

### Avoiding Duplicates

Before creating cards, scan existing cards in `cards.json`. Do not create a card if an existing card already covers the
same fact (same question or equivalent knowledge). It is fine to create cards on the same concept from different angles.

## Mistake Tracking

During practice, log recurring error categories to `mistakes.json` so that drills can target them. Without this, every
failed problem is a one-shot lesson whose signal evaporates after the session. With it, the user's weakest patterns
surface and can be isolated before they compound into bad habits.

### `mistakes.json` format

If the file doesn't exist, treat it as `{"digest_at": null, "mistakes": []}` and create it on first write.

Each entry:

```json
{
  "id": "m_<unix_timestamp>_<seq>",
  "timestamp": "<RFC3339>",
  "problem": "007",
  "concept": "binary-search",
  "category": "off-by-one",
  "note": "loop condition stayed true at equality; no contraction",
  "trigger": "check",
  "resolved_at": null
}
```

Fields:

- `id` — unique, `m_<unix_timestamp>_<seq>` (bump `seq` to disambiguate same-second logs).
- `timestamp` — RFC3339 timestamp when the mistake was observed.
- `problem` — problem number; for `trigger: "scaffold"` this is the *parent* problem, not the sub-problem ID.
- `concept` — the concept being trained at the time of the mistake.
- `category` — exactly one value from the taxonomy below.
- `note` — one line, specific to this occurrence, describing *what* went wrong. No code. No formulas.
- `trigger` — `"check"` (surfaced by a failed `check`), `"scaffold"` (surfaced while scaffolding), or `"resolve"`
  (surfaced by a failed or scaffolded re-solve — see "Re-solve Mode").
- `resolved_at` — `null` until a later clean solve exercises the same category without reproducing it, then set to the
  RFC3339 timestamp of that resolution.

### Taxonomy

Pick exactly one category per mistake. If nothing fits, use `other` with a specific `note`. If `other` recurs with
similar notes three times, promote it into a new category and update this list.

**Loops & iteration**

- `off-by-one` — loop bound wrong by 1; slice bound off by 1; `<` vs `<=`.
- `wrong-loop-bound` — iterating to the wrong endpoint, independent of parity.
- `forgotten-update` — loop variable, accumulator, or pointer not advanced in the body.
- `mutation-during-iteration` — modifying the slice/map being iterated.
- `infinite-loop` — termination condition never reached.

**Initialization & control flow**

- `uninit-accumulator` — sum/count/min/max not initialized, or initialized to the wrong sentinel.
- `early-return` — returned before finishing the computation.
- `wrong-return-value` — returned the wrong variable, wrong scope, or wrong type.
- `wrong-comparison-operator` — `<` vs `<=`, `==` vs `!=`, `>` vs `>=`.

**Recursion**

- `missing-base-case` — recursion never terminates for some input.
- `wrong-base-case` — returns wrong value at the base.
- `wrong-recursive-combine` — recurses correctly but combines subresults wrong.

**Indexing**

- `out-of-bounds` — index ≥ length, or negative index.
- `index-vs-value` — used the index where the value was needed, or vice versa.

**Edge cases**

- `empty-input-missed` — didn't handle empty slice / nil / empty string.
- `single-element-missed` — didn't handle `len == 1` or single-node case.
- `duplicates-missed` — algorithm breaks on repeated values.
- `negative-input-missed` — algorithm assumed non-negative input.
- `overflow-missed` — didn't consider integer overflow for large inputs.

**Go mechanics**

- `pointer-vs-value-receiver` — mutation via value receiver, or pointer where unnecessary.
- `slice-aliasing` — two slices share backing array; modification leaks.
- `map-iteration-order` — assumed deterministic map iteration order.
- `nil-map-write` — wrote to an unallocated map.
- `shadowed-variable` — `:=` inside a block shadowed an outer variable.

**Algorithmic**

- `wrong-algorithm` — problem named algorithm X; user implemented Y.
- `wrong-complexity` — stated complexity is wrong (when asked).
- `brute-force-when-technique-expected` — problem was training a specific technique; user brute-forced.

**Other**

- `misread-problem` — solved a different problem than the one stated.
- `syntax` — persistent Go syntax errors (not typos); log only on repeat across sessions.
- `other` — last resort.

### When to log

**On `check` (incorrect):** log one entry per *distinct* error category observed, up to three per failed check. Pick
the load-bearing categories; do not stretch to fill. Use `trigger: "check"`.

**On scaffolding:** when you create a sub-problem and can name the gap as one taxonomy category, log one entry with
`trigger: "scaffold"`.

**Dedupe within an attempt chain:** if a mistake for the same problem + category was already logged in the current
session, do not re-log. Do log if a scaffold surfaces a *different* category than the earlier check did.

**Do NOT log:**

- Typos or momentary syntax slips (missing brace, missing import) — unless they keep recurring across sessions, then
  log as `syntax`.
- User asking a clarifying question or saying "I don't know" without having attempted code.
- A first-try clean solve — nothing went wrong.

### When to resolve

After a successful `check` on a problem that *could have* exercised a category (e.g. `off-by-one` only resolves on
problems with a loop; `duplicates-missed` only resolves on problems whose inputs actually had duplicates), mark the
most recent unresolved entry of that category as `resolved_at: <now>`. Resolve at most one entry per category per
solve. If unsure whether the current problem exercised the category, leave it open — do not auto-resolve.

### Drill rules

When the `train` preamble triggers a drill:

- **Single-category focus.** A drill tests exactly one category. Do not combine.
- **Stripped to mechanics.** No larger concept on top — raw mechanics only. Prefer problems that fit in five lines of
  solution code. The point is to force the user to confront the specific failure mode with nothing to hide behind.
- **Obvious oracle.** The examples should make the right answer verifiable at a glance.
- **Frontmatter.** The problem file starts with `kind: drill` and `drill_category: <category>` so drills are
  identifiable later.
- **No promotion.** Solving a drill does NOT raise any concept level in `progress.md`. Its only effect is to resolve
  up to 3 open mistakes of the drilled category (oldest first) and to teach the pattern.
- **No scaffolding on a drill.** If the user says "I don't know" on a drill, replace the current drill with a simpler
  drill in the same category — do not spawn `NNNa.md`. Drills are already the floor.
- **One drill per `train` turn.** After the user solves (or replaces) the drill, the next `train` re-runs the preamble
  from scratch — it may pick another drill or return to normal training.

## Re-solve Mode

Re-solve mode tests whether the user can *reproduce* a solution, not just recognize it. Every successfully solved
problem goes on a Leitner ladder — when it comes due, `train` hands it back with a fresh template and the prior
solution hidden, and the user must solve it again from scratch. This is separate from card review: cards test atomic
recall, re-solve tests end-to-end execution.

### `resolve.json` format

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

### The ladder

Intervals in days, indexed by `step`:

| step | interval (days) |
|------|-----------------|
| 0    | 7               |
| 1    | 21              |
| 2    | 60              |
| 3    | 180             |
| 4    | 365             |

A problem at `step = 4` stays there on clean solves and keeps rescheduling every 365 days.

### Registration (on first solve)

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

### Exclusions

The following are NOT registered in `resolve.json` on solve:

- Drills (`kind: drill` in the problem file).
- Scaffolded sub-problems (filename has a letter suffix: `007a.md`, `007b.md`, etc.).
- Mix problems (when mix mode ships — future).

### Picker (inside `train`)

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

### Outcomes (on re-solve `check`)

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

### Mistake tracking on re-solve

- On a **failed** re-solve, log a mistake in `mistakes.json` with `trigger: "resolve"` and a category picked from the
  taxonomy. If the user never attempted any code, use `other` with a `note` describing the stall.
- On a **recovered** re-solve, mistakes logged by the scaffolding flow during this session use `trigger: "resolve"`
  (not `"scaffold"`) to distinguish re-solve gaps from first-solve gaps.
- On a **clean** re-solve, resolve open mistakes in any category this problem *could have* exercised, same rules as a
  normal clean solve.

### Scaffolding during re-solve

Scaffolding is allowed during re-solve and follows the normal Scaffolding Flow, with two additions:

1. Any use of scaffolding this session downgrades the outcome to **recovered** at best (even if the parent check
   eventually passes cleanly afterward).
2. Mistake entries logged during re-solve scaffolding use `trigger: "resolve"` instead of `"scaffold"`.

Sub-problem files generated during re-solve scaffolding use the normal `NNNa.md`, `NNNb.md` convention and are
excluded from the re-solve ladder (same as first-solve scaffolds).

### Giving up

The user signals a failed re-solve by saying "give up", "fail this", "skip re-solve", or any clear equivalent.
Do NOT auto-fail a re-solve based on scaffold depth or elapsed time — only the user decides when to stop.

When the user gives up:

1. Log the failure to `resolve.json` per the outcome table above.
2. Log a mistake entry with `trigger: "resolve"`.
3. Reset `current.md` to empty.
4. Tell the user: "Marked re-solve failed. `<concept>` consecutive failures: <N>. Level drop: <yes/no>."

The next `train` re-runs the preamble from scratch.

## Mix Mode

Mix mode tests *transfer* — whether the user can switch contexts between concepts in a single session. Research shows
interleaved practice beats blocked practice for long-term retention, even though it feels harder in the moment. Mix
sessions do NOT raise concept levels; they update a separate **retention** score per concept in `retention.json`. Cold
concepts (not exercised recently) are weighted heavier when picking which concepts to mix.

### `mix.json` format

If the file doesn't exist, treat it as `{"last_mix_at": null, "active_session": null}` and create it on first write.

```json
{
  "last_mix_at": "2026-04-16T00:00:00Z",
  "active_session": {
    "id": "mix_1713268800",
    "started_at": "2026-04-16T10:00:00Z",
    "concepts": ["arrays", "sliding-window", "stacks"],
    "problems": ["034", "035", "036"],
    "current_index": 1,
    "outcomes": ["clean"]
  }
}
```

- `last_mix_at` — RFC3339 timestamp when the last mix session **completed or was abandoned**. Never updated on session
  start — so an in-progress session does not burn the cooldown.
- `active_session` — null when no mix is in progress. When non-null, the user is working through `problems` in order.
  `outcomes` records per-problem results (clean/recovered/failed) matching position in `problems`.

### `retention.json` format

If missing, treat as `{}`. Keyed by concept name. Missing concept behaves as `{retention: 0, last_touched: null}`.

```json
{
  "arrays": {"retention": 0.8, "last_touched": "2026-04-10T00:00:00Z"},
  "binary-search": {"retention": 0.4, "last_touched": "2026-03-22T00:00:00Z"}
}
```

- `retention` — float in [0, 1]. Updated **only** by mix outcomes.
- `last_touched` — RFC3339 timestamp of the last engagement with this concept (any solve: first-solve, re-solve, or
  mix). Updated broadly because any activity refreshes the cold-concept clock.

### Retention decay

Retention decays lazily on read: `0.1` per full 14-day period since `last_touched`.

```
days = (now - last_touched) / 86400
periods = floor(days / 14)            # 0 if days < 14
effective = max(0, stored_retention - 0.1 * periods)
```

On any write to retention (mix outcome), apply decay first, then the delta, then update `last_touched`:

```
new_retention = clamp(0, 1, effective_before + delta)
```

Missing concepts behave as retention 0; no decay applied.

### Unlock

Mix is available only when **≥ 5 concepts are at level ≥ 2** in `progress.md`. Below that threshold, the mix-start
check is a no-op.

### Timing trigger

The Mix-start check in the `train` preamble fires a new session when **all three** conditions are true:

1. Unlock threshold met (≥ 5 concepts at level ≥ 2).
2. `mix.json.last_mix_at` is null OR more than 7 days ago.
3. At least **3 eligible concepts** (level ≥ 2) are **cold** — `retention.json[<concept>].last_touched` is null or
   older than 14 days.

If any condition fails, skip to digest. Do NOT start a mix session.

### Starting a new session

1. **Pick 3 concepts.** Pool = concepts at level ≥ 2. Weight each cold concept (last_touched null or > 14 days ago) at
   **2×**, non-cold concepts at 1×. Draw 3 distinct concepts via weighted random.
2. **Pick a problem per concept.** For each selected concept:
   - Eligible level = `max(1, current_level - 1)` — one step below the user's working level. Interleaving is the
     challenge; don't stack novelty on top.
   - Pool = bank problems at that level the user has NOT solved in the last 30 days. This includes fresh bank
     problems the user has never seen (preferred) and previously solved ones that have cooled off.
   - Prefer fresh bank problems over reused ones; random pick from the eligible pool. If the pool is empty, generate
     a new bank-style variation for that concept at the same level.
3. **Generate session ID:** `mix_<unix_timestamp>`.
4. **Create problem files:**
   - Fresh bank problems → create `problems/NNN.md` as normal, with two extra header fields: `**Kind:** mix` and
     `**Mix Session:** <session_id>`.
   - Reused problems → do NOT duplicate the file; reuse the existing `NNN.md` in place.
5. **Write `mix.json.active_session`:**
   ```json
   {
     "id": "<session_id>",
     "started_at": "<now, RFC3339>",
     "concepts": ["<c1>", "<c2>", "<c3>"],
     "problems": ["<n1>", "<n2>", "<n3>"],
     "current_index": 0,
     "outcomes": []
   }
   ```
6. **Set `current.md`** to `<n1>:mix`.
7. **Write a fresh `main.go` template** for problem `<n1>`.
8. **Announce** in one sentence: "Mix of 3: `<c1>`, `<c2>`, `<c3>` — cold concepts surfacing. First up: problem `<n1>`
   (<title>)."

### Mix problem outcomes (on `check`)

When `current.md` ends in `:mix`, classify the outcome for this single problem:

- **Clean** — `check` passes first try; no "I don't know"; no prior failed `check` on this mix problem.
- **Recovered** — passes eventually, but scaffolding was used OR a prior `check` failed OR a hint was given.
- **Failed** — user says "skip", "next", "give up on this one", or equivalent before a successful `check`.

State updates per outcome, applied to `retention.json[<concept>]`:

| outcome   | `retention` delta   | `last_touched` |
|-----------|---------------------|----------------|
| clean     | +0.2 (cap 1.0)      | now            |
| recovered | 0                   | now            |
| failed    | −0.3 (floor 0)      | now            |

Apply decay first (see "Retention decay"), then the delta. Append the outcome to
`mix.json.active_session.outcomes`.

Do NOT raise the concept level. Do NOT register in `resolve.json`. Do NOT re-create spaced-repetition cards for
reused problems; fresh mix problems follow normal card-creation rules on clean solve.

Mistake tracking works normally for mix problems — failed `check` logs with `trigger: "check"`; scaffolding uses
`trigger: "scaffold"`.

### Advancing the session

After classifying a mix problem's outcome:

- If `current_index + 1 < len(problems)`:
  - Increment `current_index`.
  - Set `current.md` to `<next_problem>:mix`.
  - Write a fresh `main.go` template for the next problem.
  - Announce: "Next in mix: `<concept>` — problem `<num>` (<title>), <i+1>/<len>."
- Else (last problem in the session):
  - Set `mix.json.last_mix_at` = now.
  - Set `mix.json.active_session` = null.
  - Reset `current.md` to empty.
  - Summarize: "Mix session complete. Outcomes: clean×N, recovered×M, failed×K. Retention updated."

### Skipping one problem vs abandoning the session

- **Skip one** — "skip", "next", "move on this one" → classify as **failed**, update retention, advance to next.
- **Abandon session** — "end mix", "quit mix", "abandon session" → classify the current problem as **failed**, mark
  any remaining problems as **failed** too (applying retention deltas for each), set `last_mix_at = now`, clear
  `active_session`, reset `current.md`. Summarize as normal.

### Scaffolding during mix

Scaffolding is allowed during mix. It downgrades the current mix problem's outcome to **recovered** at best, even if
the user eventually passes cleanly. Mistake entries logged during mix scaffolding use `trigger: "scaffold"` (not a
new value — scaffolds during mix are ordinary scaffolds).

Sub-problem files generated during mix scaffolding follow the normal `NNNa.md`, `NNNb.md` convention. They are NOT
part of the mix session and are NOT registered in `resolve.json`.

### Resume semantics

If `mix.json.active_session` is non-null when `train` runs, the Active-mix-resume preamble step preempts everything
else — no drill check, no re-solve check, no digest. Rationale: abandoning a session half-complete fragments the
retention signal, and drills/re-solves will queue patiently until the user finishes (or explicitly abandons).

### Exclusions

- Mix problems are NOT added to `resolve.json` (they don't enter the reproduction ladder).
- Drills and scaffold sub-problems are never picked *for* a mix session.
- Reused problems that already have an upcoming re-solve within 3 days are skipped (avoid back-to-back re-exposure).

## Mistakes Command

When the user says **"mistakes"**:

1. If `mistakes.json` doesn't exist or has no entries, say "No mistakes logged yet."
2. Otherwise, print an on-demand report:
   - Top 5 categories by unresolved count (over all time).
   - Total unresolved / total logged.
   - Most recent 5 entries, one line each: `<timestamp> <category> <problem> — <note>`.
3. Do NOT update `digest_at` — this view is separate from the weekly digest gate.

## Review Command

When the user says **"review"**:

1. Check if `cards.json` exists and has cards.
2. If no cards exist, tell the user: "No review cards yet. Solve some problems first!"
3. If cards exist, tell the user to start their review session: "Run `go run ./cmd/review` to start your review
   session."
