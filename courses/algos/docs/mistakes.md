# Mistake Tracking

During practice, log recurring error categories to `mistakes.json` so that drills can target them. Without this, every
failed problem is a one-shot lesson whose signal evaporates after the session. With it, the user's weakest patterns
surface and can be isolated before they compound into bad habits.

## `mistakes.json` format

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
  (surfaced by a failed or scaffolded re-solve — see `docs/resolve.md`).
- `resolved_at` — `null` until a later clean solve exercises the same category without reproducing it, then set to the
  RFC3339 timestamp of that resolution.

## Taxonomy

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

## When to log

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

## When to resolve

After a successful `check` on a problem that *could have* exercised a category (e.g. `off-by-one` only resolves on
problems with a loop; `duplicates-missed` only resolves on problems whose inputs actually had duplicates), mark the
most recent unresolved entry of that category as `resolved_at: <now>`. Resolve at most one entry per category per
solve. If unsure whether the current problem exercised the category, leave it open — do not auto-resolve.

## Drill rules

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
