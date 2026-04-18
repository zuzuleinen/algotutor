# Concept List

32 concepts in teaching order. Prerequisites in parentheses are the *conceptual* prerequisites — the vocabulary and
intuitions the user must already have in place to even understand the explanation of the new concept. Before training
any concept, every prerequisite must be at level ≥ 1 in `progress.md`. If a prerequisite is below level 1, do not teach
the new concept — reroute training to the missing prerequisite first, even if it comes later in this list.

The list is ordered so that prerequisites always appear before dependents. If you find yourself about to teach a
concept before one of its listed prerequisites, stop — something is wrong. Fix the ordering or the progress file; do
not teach the concept out of order.

## Fundamentals

1. **arrays** — indexing, in-place operations, rotation, subarrays
2. **strings** — manipulation, substrings, palindromes, character counting
3. **loops** — for loops, while loops, iterating over slices
4. **nested-loops** — pairs, triplets, brute-force enumeration
5. **math** — GCD, primes, modular arithmetic, integer properties

## Core Data Structures

6. **maps** — hash maps for O(1) lookup, frequency counting, grouping
7. **sets** — uniqueness, intersection, union, membership testing
8. **matrix** — 2D arrays, traversal, rotation, spiral order
9. **stacks** — LIFO, matching brackets, monotonic stacks, expression evaluation
10. **queues** — FIFO, BFS usage, deques
11. **linked-lists** — traversal, reversal, fast/slow pointers, merge, cycle detection *(requires: arrays)*

## Core Techniques

12. **sorting** — sort.Slice, custom comparators, sorted order reasoning *(requires: arrays, loops)*
13. **binary-search** — search in sorted data, search on answer, rotated arrays *(requires: arrays, sorting)*
14. **two-pointers** — left/right on sorted data, pair/triplet finding *(requires: arrays, sorting)*
15. **sliding-window** — fixed and variable size windows, substring problems *(requires: arrays, loops)*
16. **prefix-sums** — range sums, subarray sums, running totals *(requires: arrays, loops)*

## Recursion and Trees

17. **recursion** — base cases, recursive thinking, call stack *(requires: loops, math)*
18. **trees** — binary trees, traversals (inorder/preorder/postorder), DFS, BFS, BST operations *(requires: recursion, linked-lists)*
19. **heaps** — priority queues, top-k, merge-k-sorted, median tracking *(requires: trees, arrays)*
20. **tries** — prefix trees, autocomplete, word search *(requires: trees, maps)*

## Graph Algorithms

21. **graphs** — DFS, BFS, connected components, adjacency lists *(requires: trees, queues, maps)*
22. **topological-sort** — dependency ordering, DAGs, cycle detection *(requires: graphs)*
23. **union-find** — disjoint sets, connected components, ranking *(requires: arrays, recursion)*
24. **shortest-path** — Dijkstra, BFS unweighted, Bellman-Ford *(requires: graphs, heaps)*

## Advanced Techniques

25. **greedy** — local optimal choices, interval scheduling, activity selection *(requires: sorting)*
26. **intervals** — merge, insert, overlap detection, sweep line *(requires: sorting, arrays)*
27. **backtracking** — permutations, combinations, constraint satisfaction, N-queens *(requires: recursion)*
28. **divide-and-conquer** — merge sort pattern, split and combine *(requires: recursion, arrays)*
29. **dynamic-programming** — memoization, tabulation, subproblems, knapsack, LCS, LIS *(requires: recursion, maps, arrays)*
30. **monotonic-stacks** — next greater element, histogram problems, stock span *(requires: stacks)*
31. **design** — LRU cache, iterator design, data structure composition *(requires: linked-lists, maps, heaps)*
32. **bit-manipulation** — bitwise ops, masks, XOR tricks, counting bits *(requires: math)*
