# Agent setup

algotutor's behaviour lives in `AGENTS.md` and is mirrored to `CLAUDE.md` and
`GEMINI.md` by `make sync-agents`. Any AI coding agent that can read files, edit
files, and run shell commands can drive the workflow once it loads those
instructions.

This page lists per-agent setup tips. Each agent is independent — pick one and go.

## Quick reference

| Agent             | Loads instructions from                                | Model config              | Permissions                                                              |
| ----------------- | ------------------------------------------------------ | ------------------------- | ------------------------------------------------------------------------ |
| Claude Code       | `CLAUDE.md` (auto)                                     | `.claude/settings.json`   | `.claude/settings.local.json`, or run with `--dangerously-skip-permissions` |
| OpenAI Codex CLI  | `AGENTS.md` (auto)                                     | `~/.codex/config.toml`    | `--sandbox <mode>` / `--ask-for-approval <policy>`                       |
| Cursor            | `AGENTS.md` (auto, recent builds) or `.cursor/rules/`  | IDE settings              | Auto via Composer / Agent mode                                           |
| Cline / Roo Code  | `AGENTS.md` (auto) or `.clinerules/`                   | Extension settings        | Extension permissions UI                                                 |
| OpenCode          | `AGENTS.md` (auto)                                     | `~/.config/opencode/`     | Per-launch flag                                                          |
| Aider             | `AGENTS.md` via `--read` flag                          | `--model <name>`          | `--yes` for auto-approve                                                 |
| Gemini CLI        | `GEMINI.md` (auto)                                     | Built-in                  | Per-launch flag                                                          |

## Per-agent setup

### Claude Code

```sh
claude
# or, to skip per-tool prompts:
claude --dangerously-skip-permissions
```

`CLAUDE.md` is auto-loaded on session start. Default model is set in
`.claude/settings.json` (`claude-sonnet-4-6`). Permission allow-list is in
`.claude/settings.local.json`.

### OpenAI Codex CLI

```sh
codex
# or, to choose permissions for the session:
codex --sandbox workspace-write --ask-for-approval on-request
```

`AGENTS.md` is auto-loaded from the project root. Model preference goes in
`~/.codex/config.toml`. Recommend a current reasoning-class model.

### Cursor

Open the project folder in Cursor, switch to Agent / Composer mode, and type
`train` in the chat. Recent Cursor builds auto-load `AGENTS.md`. If yours doesn't,
add a one-line `.cursor/rules/main.mdc` containing `@AGENTS.md`.

### Cline / Roo Code

Open the project in VS Code with the Cline (or Roo Code) extension installed,
activate the chat panel, and type `train`. Both extensions auto-load `AGENTS.md`.

### OpenCode

```sh
opencode
```

`AGENTS.md` is auto-loaded. Model is configured per-launch or in
`~/.config/opencode/`.

### Aider

```sh
aider --read AGENTS.md --model <your-model>
```

Aider does not auto-discover instruction files — pass `--read AGENTS.md`
explicitly. Add `--yes` to auto-approve edits.

### Gemini CLI

```sh
gemini
```

`GEMINI.md` is auto-loaded (it's a byte-identical mirror of `AGENTS.md`).

## Auto-launch from `make`

If you set a default agent during `make init`, `make train` and `make review` will
auto-launch it for you with the right prompt. Otherwise they print "Open your agent
and type `train`" and you do the launching.

## If your agent doesn't auto-load

Type this once at session start:

> Read AGENTS.md and follow it.

After that, all algotutor commands (`train`, `check`, `I don't know`, `mistakes`,
`review`, `reset`, `I want to solve [X]`) work normally.

## Switching agents mid-session

All state lives on disk in JSON / Markdown files: `progress.md`, `current.md`,
`cards.json`, `mistakes.json`, `resolve.json`, `mix.json`, `retention.json`. Stop
one agent, start another in the same directory, and the next `train` picks up
where you left off — same current problem, same concept levels, same mix /
re-solve schedule, same mistake log. No re-bootstrap is needed.

The only externally visible difference is response style — Sonnet, GPT-5, and
Gemini phrase things differently. The workflow, rules, and state are identical.

## Keeping mirrors in sync

`AGENTS.md` is the canonical instruction file. `CLAUDE.md` and `GEMINI.md` are
byte-identical copies. After editing `AGENTS.md`:

```sh
make sync-agents    # regenerate the mirrors
make check-agents   # verify mirrors match (use in CI / pre-commit)
```

If `check-agents` fails, your mirrors have drifted — run `sync-agents`.
