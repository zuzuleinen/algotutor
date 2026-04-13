# Algorithmic Training Project

## Project Structure
- `main.go` — the only Go file. Always replaced with the current problem template.
- `problems/` — one `.md` file per problem, saved as `001.md`, `002.md`, etc.
- `current.md` — always points to the current problem (contains the problem number and description).
- `progress.md` — tracks the user's level (0–N) for each concept.
- `progress.template.md` — blank progress table (all zeros), used for initialization.
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
Track the user's level for each concept in `progress.md`. Levels start at 0 (never seen) and increase as the user solves problems. Each level maps to harder problems within that concept.

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
11. **linked-lists** — traversal, reversal, fast/slow pointers, merge, cycle detection
12. **heaps** — priority queues, top-k, merge-k-sorted, median tracking

**Core Techniques**
13. **sorting** — sort.Slice, custom comparators, sorted order reasoning
14. **binary-search** — search in sorted data, search on answer, rotated arrays
15. **two-pointers** — left/right on sorted data, pair/triplet finding
16. **sliding-window** — fixed and variable size windows, substring problems
17. **prefix-sums** — range sums, subarray sums, running totals
18. **bit-manipulation** — bitwise ops, masks, XOR tricks, counting bits

**Recursion and Trees**
19. **recursion** — base cases, recursive thinking, call stack
20. **trees** — traversals (inorder/preorder/postorder), DFS, BFS, BST operations
21. **tries** — prefix trees, autocomplete, word search

**Graph Algorithms**
22. **graphs** — DFS, BFS, connected components, adjacency lists
23. **topological-sort** — dependency ordering, DAGs, cycle detection
24. **union-find** — disjoint sets, connected components, ranking
25. **shortest-path** — Dijkstra, BFS unweighted, Bellman-Ford

**Advanced Techniques**
26. **greedy** — local optimal choices, interval scheduling, activity selection
27. **intervals** — merge, insert, overlap detection, sweep line
28. **backtracking** — permutations, combinations, constraint satisfaction, N-queens
29. **divide-and-conquer** — merge sort pattern, split and combine
30. **dynamic-programming** — memoization, tabulation, subproblems, knapsack, LCS, LIS
31. **monotonic-stacks** — next greater element, histogram problems, stock span
32. **design** — LRU cache, iterator design, data structure composition

### Level Progression Within a Concept
- **Level 0**: Never seen. Start with the simplest possible problem for this concept.
- **Level 1**: Can do the basic pattern. Give a slightly harder variation.
- **Level 2**: Comfortable. Introduce edge cases or combine with a previously learned concept.
- **Level 3**: Strong. Give problems that require this concept as a tool within a larger problem.
- **Level 4+**: Mastery. Interview-level problems featuring this concept.

## Training Flow
When the user says **"train"**:
1. Read `progress.md` to see concept levels.
2. Pick the concept to train: the earliest concept in the list that is below level 3, prioritizing concepts at level 0 first, then level 1, then level 2.
3. Generate a problem appropriate for the user's current level in that concept.
4. Save it to `problems/NNN.md` with the concept noted.
5. Update `current.md` to point to it.
6. Write the problem template into `main.go` and present it.

## Solving a Specific Problem
When the user says **"I want to solve [problem name]"**:
1. Create the problem at the right difficulty level.
2. Save it to `problems/NNN.md` with status `pending`.
3. Update `current.md` to point to it.
4. Write the problem template into `main.go` and present it.

## Checking Flow
When the user says **"check"**:
1. Read `current.md` to find the current problem.
2. Read `problems/NNN.md` for the expected behavior.
3. Evaluate the user's solution in `main.go`.
4. If correct: mark the problem as `solved` in `problems/NNN.md`, update the concept level in `progress.md`, congratulate briefly, and stop.
5. If incorrect: give the minimal nudge needed. Do not give hints unless asked.

## Scaffolding Flow
When the user says **"I don't know"**:
1. Break the current problem into a simpler sub-problem.
2. Save the sub-problem as a new entry (e.g. `003a.md`, `003b.md`).
3. Update `current.md` to point to the sub-problem.
4. Keep going simpler until the user can solve it.
5. Once solved, step back up toward the original problem.

## Rules
- Never give hints unless the user asks.
- Never add helpful remarks or commentary unless asked.
- Always put `main` first in every Go file.
- Always add the problem description as a comment at the top of `main`.
- In `main`, every `fmt.Println` call must have an inline comment showing the expected output, e.g. `fmt.Println(reverseString("hello")) // "olleh"`.
- Always produce a valid, runnable `.go` file.
- Never use more than one file for Go code.
- Problems increase in complexity gradually, always building on what was just learned.
- The end goal is always to solve the originally requested problem.
- When training, assume the user starts from zero knowledge. Do not skip basics.
