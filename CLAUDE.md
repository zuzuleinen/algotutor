# Algorithmic Training Project

## Project Structure

- `main.go` — the only Go file for problems. Always replaced with the current problem template.
- `problems/` — one `.md` file per problem, saved as `001.md`, `002.md`, etc.
- `current.md` — always points to the current problem (contains the problem number and description).
- `progress.md` — tracks the user's level (0–N) for each concept.
- `progress.template.md` — blank progress table (all zeros), used for initialization.
- `cards.json` — spaced repetition review cards (created automatically during practice).
- `mistakes.json` — log of recurring error categories (created automatically during practice).
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

`current.md` contains only the number of the current problem, e.g. `003`.

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
17. **bit-manipulation** — bitwise ops, masks, XOR tricks, counting bits *(requires: math)*

**Recursion and Trees**

18. **recursion** — base cases, recursive thinking, call stack *(requires: loops, math)*
19. **trees** — binary trees, traversals (inorder/preorder/postorder), DFS, BFS, BST operations *(requires: recursion, linked-lists)*
20. **heaps** — priority queues, top-k, merge-k-sorted, median tracking *(requires: trees, arrays)*
21. **tries** — prefix trees, autocomplete, word search *(requires: trees, maps)*

**Graph Algorithms**

22. **graphs** — DFS, BFS, connected components, adjacency lists *(requires: trees, queues, maps)*
23. **topological-sort** — dependency ordering, DAGs, cycle detection *(requires: graphs)*
24. **union-find** — disjoint sets, connected components, ranking *(requires: arrays, recursion)*
25. **shortest-path** — Dijkstra, BFS unweighted, Bellman-Ford *(requires: graphs, heaps)*

**Advanced Techniques**

26. **greedy** — local optimal choices, interval scheduling, activity selection *(requires: sorting)*
27. **intervals** — merge, insert, overlap detection, sweep line *(requires: sorting, arrays)*
28. **backtracking** — permutations, combinations, constraint satisfaction, N-queens *(requires: recursion)*
29. **divide-and-conquer** — merge sort pattern, split and combine *(requires: recursion, arrays)*
30. **dynamic-programming** — memoization, tabulation, subproblems, knapsack, LCS, LIS *(requires: recursion, maps, arrays)*
31. **monotonic-stacks** — next greater element, histogram problems, stock span *(requires: stacks)*
32. **design** — LRU cache, iterator design, data structure composition *(requires: linked-lists, maps, heaps)*

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

When the user says **"train"**, first run the **mistake-tracking preamble**, then fall through to the picker:

**Preamble (may preempt the picker):**

- **Drill check.** Read `mistakes.json` (treat as `{"digest_at": null, "mistakes": []}` if missing). Take the last 20
  entries with `resolved_at: null`. Group by `category`. If any category has ≥ 3 unresolved entries in that window,
  generate a **single-category drill** for it (see "Mistake Tracking"), save it as the current problem, write the
  template into `main.go`, and present it — **skip the picker below for this turn.** Tell the user in one sentence:
  "You've hit `<category>` <N>× recently — drill first."
- **Digest.** Else, if `digest_at` is null or older than 7 days and there is at least one unresolved mistake, print a
  short digest (top 3 unresolved categories over the last 30 days with counts + one-sentence recommendation), update
  `digest_at` to the current RFC3339 timestamp, then continue to the picker.

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

1. Read `current.md` to find the current problem.
2. Read `problems/NNN.md` for the expected behavior.
3. Evaluate the user's solution in `main.go`.
4. **Algorithm fidelity.** If the problem names a specific algorithm or data structure (e.g. "bubble sort",
   "implement using a stack", "recursive solution"), the user's solution MUST implement *that* algorithm or
   technique. Correct output via a different algorithm does NOT count as solved — that's teaching the wrong
   concept. Treat it as incorrect and nudge them toward the named approach.
5. **Never offer "accept as-is" as an option.** Do not present the user with a multiple-choice menu like
   "1. accept 2. redo". The contract is: either they solve the stated problem with the stated technique,
   or they say "I don't know" to get a scaffolded easier version. There is no third door.
6. If correct (and algorithm matches): mark the problem as `solved` in `problems/NNN.md`, update the concept
   level in `progress.md`, congratulate briefly, then create spaced repetition cards (see "Spaced Repetition
   Cards" section below) and, for each mistake category this problem *could have* exercised, mark the most
   recent unresolved entry of that category in `mistakes.json` as resolved (see "Mistake Tracking"). Drills
   are an exception — see the drill rules for how they resolve mistakes and skip level promotion.
7. If incorrect: **name what is wrong** (e.g. "your loop condition is off by one", "you're not
   updating the sum") but **never supply the fix directly** — no corrected expressions, no formulas, no
   rewritten lines. If the user says "I don't know" in response, route them to a scaffolded sub-problem
   that, once solved, will make the fix obvious. The user must derive every expression themselves.
   **Log the mistake** in `mistakes.json` — one entry per distinct error category, max three per failed
   check, using the fixed taxonomy (see "Mistake Tracking").
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
- `trigger` — `"check"` (surfaced by a failed `check`) or `"scaffold"` (surfaced while scaffolding).
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
