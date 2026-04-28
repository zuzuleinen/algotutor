# Concept List

35 concepts in teaching order. Prerequisites in parentheses are the *conceptual* prerequisites — the vocabulary and
intuitions the user must already have in place to even understand the explanation of the new concept. Before training
any concept, every prerequisite must be at level ≥ 1 in `progress.md`. If a prerequisite is below level 1, do not teach
the new concept — reroute training to the missing prerequisite first, even if it comes later in this list.

The list is ordered so that prerequisites always appear before dependents. If you find yourself about to teach a
concept before one of its listed prerequisites, stop — something is wrong. Fix the ordering or the progress file; do
not teach the concept out of order.

## Fundamentals

1. **complexity-analysis** — Big-O notation, time and space complexity, comparing algorithm efficiency, reading constraints to pick an approach
2. **arrays** — indexing, in-place operations, rotation, subarrays
3. **strings** — manipulation, substrings, palindromes, character counting
4. **loops** — for loops, while loops, iterating over slices
5. **nested-loops** — pairs, triplets, brute-force enumeration
6. **math** — GCD, primes, modular arithmetic, integer properties

## Core Data Structures

7. **maps** — hash maps for O(1) lookup, frequency counting, grouping
8. **sets** — uniqueness, intersection, union, membership testing
9. **matrix** — 2D arrays, traversal, rotation, spiral order
10. **stacks** — LIFO, matching brackets, monotonic stacks, expression evaluation
11. **queues** — FIFO, BFS usage, deques
12. **linked-lists** — traversal, reversal, fast/slow pointers, merge, cycle detection *(requires: arrays)*

## Core Techniques

13. **sorting** — sort.Slice, custom comparators, sorted order reasoning *(requires: arrays, loops)*
14. **binary-search** — search in sorted data, search on answer, rotated arrays *(requires: arrays, sorting)*
15. **two-pointers** — left/right on sorted data, pair/triplet finding *(requires: arrays, sorting)*
16. **sliding-window** — fixed and variable size windows, substring problems *(requires: arrays, loops)*
17. **prefix-sums** — range sums, subarray sums, running totals *(requires: arrays, loops)*

## Recursion and Trees

18. **recursion** — base cases, recursive thinking, call stack *(requires: loops, math)*
19. **trees** — binary trees, traversals (inorder/preorder/postorder), DFS, BFS, BST operations *(requires: recursion, linked-lists)*
20. **heaps** — priority queues, top-k, merge-k-sorted, median tracking *(requires: trees, arrays)*
21. **tries** — prefix trees, autocomplete, word search *(requires: trees, maps)*

## Graph Algorithms

22. **graph-modeling** — identifying nodes and edges in word problems, adjacency list representation, directed vs undirected, weighted vs unweighted *(requires: trees, queues, maps)*
23. **graphs** — DFS, BFS, connected components *(requires: graph-modeling)*
24. **topological-sort** — dependency ordering, DAGs, cycle detection *(requires: graphs)*
25. **union-find** — disjoint sets, connected components, ranking *(requires: arrays, recursion)*
26. **shortest-path** — Dijkstra, BFS unweighted, Bellman-Ford *(requires: graphs, heaps)*

## Advanced Techniques

27. **greedy** — local optimal choices, interval scheduling, activity selection *(requires: sorting)*
28. **intervals** — merge, insert, overlap detection, sweep line *(requires: sorting, arrays)*
29. **backtracking** — permutations, combinations, constraint satisfaction, N-queens *(requires: recursion)*
30. **divide-and-conquer** — merge sort pattern, split and combine *(requires: recursion, arrays)*
31. **dynamic-programming-1d** — memoization and tabulation on linear state: Fibonacci, house robber, coin change, LIS *(requires: recursion, arrays)*
32. **dynamic-programming-2d** — 2D and multi-dimensional state: grid paths, LCS, edit distance, knapsack *(requires: dynamic-programming-1d)*
33. **monotonic-stacks** — next greater element, histogram problems, stock span *(requires: stacks)*
34. **design** — LRU cache, iterator design, data structure composition *(requires: linked-lists, maps, heaps)*
35. **bit-manipulation** — bitwise ops, masks, XOR tricks, counting bits *(requires: math)*
