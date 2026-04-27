<h2 align="center">
    <picture>
        <source media="(prefers-color-scheme: dark)" srcset="logo-dark.svg">
        <source media="(prefers-color-scheme: light)" srcset="logo-light.svg">
        <img height="80" alt="algotutor" src="logo.svg">
    </picture>
    <br>
    AI-powered algorithmic training for Go developers
</h2>

<div align="center">

[![Agent-agnostic](https://img.shields.io/badge/Agent-agnostic-FA9BFA?style=flat)](docs/agents.md)
[![Go](https://img.shields.io/badge/Go-1.26-4B78E6?style=flat&logo=go&logoColor=white)](https://go.dev/)
[![Spaced repetition](https://img.shields.io/badge/Review-FSRS-73DC8C?style=flat)](https://github.com/open-spaced-repetition/go-fsrs)

</div>

**algotutor** turns an AI coding session into a personal algorithms tutor.
Open your agent in this directory, type `train`, and start solving Go problems.
It tracks your skill level across 32 concepts — from arrays and strings up through
dynamic programming and design — and picks the next problem based on where you are.

Spaced-repetition review, mistake tracking, re-solves, and interleaved mix sessions
are all built in.

<div align="center">

<img src="img_2.png" width="700" alt="algotutor in action"/>

</div>

## How it works

An AI agent acts as a tutor that generates progressively harder algorithm problems in Go. It tracks your skill level across
32 concepts — from arrays and strings up through dynamic programming and design — and picks the next problem
based on where you are.

The agent reads its instructions from `AGENTS.md` (mirrored to `CLAUDE.md` and `GEMINI.md` so each agent auto-loads the right file). Any agent that can read files, edit files, and run shell commands works — Claude Code, OpenAI Codex CLI, Cursor, Cline, Aider, OpenCode, Gemini CLI. See [Supported agents](#supported-agents) below, or [docs/agents.md](docs/agents.md) for the full setup matrix.

Read more
details: [algotutor: using AI to actually get better at algorithms](https://medium.com/@andreiboar/algotutor-using-ai-to-actually-get-better-at-algorithms-a2b7b96e054a)

### Commands

| Command                          | What it does                                                        |
|----------------------------------|---------------------------------------------------------------------|
| `train`                          | Get the next problem — drill, re-solve, mix, or new, based on state |
| `check`                          | Submit your solution for evaluation                                 |
| `I don't know`                   | Break the problem into simpler sub-problems                         |
| `I want to solve [problem name]` | Request a specific problem                                          |
| `review`                         | Check if you have cards due for review                              |
| `mistakes`                       | Show your recurring-error report                                    |

### Spaced repetition review

As you solve problems, the agent automatically creates review cards capturing what you learned — algorithmic patterns, Go
syntax, data structure properties. Cards follow
the [SuperMemo 20 Rules for effective memorization](https://www.supermemo.com/en/blog/twenty-rules-of-formulating-knowledge).

Run `make review` (or `go run ./cmd/review`) to start an Anki-style review session. The review TUI uses
the [FSRS](https://github.com/open-spaced-repetition/go-fsrs) algorithm to schedule cards. Rate each card 1–4 (
Again/Hard/Good/Easy) and it will reappear at the optimal interval.

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

- An AI coding agent — see [Supported agents](#supported-agents)
- [Go](https://go.dev/)

## Supported agents

algotutor works with any AI coding agent that can read files, edit files, and run shell commands. Most agents auto-load `AGENTS.md` (or `CLAUDE.md` / `GEMINI.md`, which are byte-identical mirrors). Pick one:

| Agent                                                         | How to use                                                       |
|---------------------------------------------------------------|------------------------------------------------------------------|
| [Claude Code](https://docs.anthropic.com/en/docs/claude-code) | `claude` — auto-loads `CLAUDE.md`                                |
| [OpenAI Codex CLI](https://github.com/openai/codex)           | `codex --auto-edit` — auto-loads `AGENTS.md`                     |
| [Cursor](https://cursor.com)                                  | Open folder, switch to Agent mode — auto-loads `AGENTS.md`       |
| [Cline](https://github.com/cline/cline)                       | VS Code extension; type `train` in chat — auto-loads `AGENTS.md` |
| [OpenCode](https://github.com/sst/opencode)                   | `opencode` — auto-loads `AGENTS.md`                              |
| [Aider](https://aider.chat)                                   | `aider --read AGENTS.md`                                         |
| Gemini CLI                                                    | `gemini` — auto-loads `GEMINI.md`                                |

See [docs/agents.md](docs/agents.md) for per-agent model selection, permission flags, and bootstrap notes for agents that don't auto-load.

You can switch agents mid-session — all state lives in JSON / Markdown files on disk, so the next agent picks up exactly where the previous one left off.

## Getting started

1. Clone the repo
2. Open your AI agent in this directory
3. Type `train`

On first run, the agent will initialize your progress file and problem directory. Your progress is local and gitignored.

## Recommendations

The working problem is always inside `main.go`. You can validate with `go run .` before saying `check`.

Try to make as much progress as you can before saying `I don't know`. This way the agent can better assess your gaps and
missing prerequisites.

If you use an IDE with AI auto-completion, disable it.

It should feel effortful. Don't be afraid to say `I don't know` multiple times. Practice regularly in sessions of 30-60
minutes.

For agent-specific tips (model selection, permission flags, defaults), see [docs/agents.md](docs/agents.md).
