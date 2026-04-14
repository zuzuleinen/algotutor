# Algorithmic Training Project

## Project Structure

- `main.go` — the only Go file for problems. Always replaced with the current problem template.
- `problems/` — one `.md` file per problem, saved as `001.md`, `002.md`, etc.
- `current.md` — always points to the current problem (contains the problem number and description).
- `progress.md` — tracks the user's level (0–N) for each concept.
- `progress.template.md` — blank progress table (all zeros), used for initialization.
- `cards.json` — spaced repetition review cards (created automatically during practice).
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

When the user says **"train"**:

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
4. If correct: mark the problem as `solved` in `problems/NNN.md`, update the concept level in `progress.md`,
   congratulate briefly, then create spaced repetition cards (see "Spaced Repetition Cards" section below).
5. If incorrect: give the minimal nudge needed. Do not give hints unless asked.

## Scaffolding Flow

When the user says **"I don't know"**:

1. Read `main.go` first to see what the user has written so far. Use their partial attempt to identify exactly where
   they got stuck, and tailor the sub-problem to that specific gap.
2. Break the current problem into a simpler sub-problem targeting the identified gap.
3. Save the sub-problem as a new entry (e.g. `003a.md`, `003b.md`).
4. Update `current.md` to point to the sub-problem.
5. Keep going simpler until the user can solve it.
6. Once solved, create spaced repetition cards for the gap that was identified (see "Spaced Repetition Cards" section
   below), then step back up toward the original problem.

## Rules

- Never give hints unless the user asks.
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

## Review Command

When the user says **"review"**:

1. Check if `cards.json` exists and has cards.
2. If no cards exist, tell the user: "No review cards yet. Solve some problems first!"
3. If cards exist, tell the user to start their review session: "Run `go run ./cmd/review` to start your review
   session."
