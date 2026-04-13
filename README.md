# algotutor

An AI-powered algorithmic training system. Open a Claude Code session in this directory, type `train`, and start solving problems.

## How it works

Claude acts as a tutor that generates progressively harder algorithm problems in Go. It tracks your skill level across 32 concepts — from arrays and strings up through dynamic programming and system design — and picks the next problem based on where you are.

### Commands

| Command | What it does |
|---|---|
| `train` | Get the next problem based on your progress |
| `check` | Submit your solution for evaluation |
| `I don't know` | Break the problem into simpler sub-problems |
| `I want to solve [problem name]` | Request a specific problem |

### Concepts covered

**Fundamentals** — arrays, strings, loops, nested loops, math

**Core Data Structures** — maps, sets, matrix, stacks, queues, linked lists, heaps

**Core Techniques** — sorting, binary search, two pointers, sliding window, prefix sums, bit manipulation

**Recursion and Trees** — recursion, trees, tries

**Graph Algorithms** — graphs, topological sort, union-find, shortest path

**Advanced Techniques** — greedy, intervals, backtracking, divide and conquer, dynamic programming, monotonic stacks, design

## Requirements

- [Claude Code](https://docs.anthropic.com/en/docs/claude-code)
- [Go](https://go.dev/)

## Getting started

1. Clone the repo
2. Open a Claude Code session in the directory
3. Type `train`

On first run, Claude will initialize your progress file and problem directory. Your progress is local and gitignored.
