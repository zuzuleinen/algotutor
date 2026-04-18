# Go Gotchas — Precision Traps When Teaching

A checklist of Go-specific semantics that, if glossed over, produce sloppy
teaching and bake latent bugs into the user's mental model. **Consult this file
before writing any problem statement, example, or nudge that depends on the
affected mechanic.** When a gotcha applies, either (a) state the narrowing
contract in the problem ("ASCII-only", "non-negative inputs", "length ≥ 1"), or
(b) use the fully general form. Never silently pick the ASCII-simple form as
"the answer" without saying so.

---

## Strings are UTF-8 bytes, not characters

- `string` is an immutable sequence of **bytes**, not characters.
- `s[i]` returns the `i`-th **byte** of the UTF-8 encoding. It equals the `i`-th
  character **only when `s` is pure ASCII**.
- One Unicode code point (rune) spans 1–4 bytes in UTF-8. Examples:
  - `'h'` → 1 byte, `'é'` → 2 bytes, `'中'` → 3 bytes, `'😀'` → 4 bytes.
- `len(s)` is the byte length, not the character count.
- `for i, r := range s` yields `(byte_offset, rune)` pairs — i jumps by 1–4
  depending on the rune. `utf8.RuneCountInString(s)` gives the rune count.
- To index by character, use `[]rune(s)[i]` (allocates, O(n)) or iterate with
  `range`.

**When teaching:**

- If a problem's solution uses `s[i]` as "the i-th character", state that
  input is ASCII-only in the problem file. Otherwise default to `[]rune(s)` or
  `range`.
- Never nudge the user from `[]rune(s)[i]` toward `string(s[i])` as "simpler"
  unless the ASCII contract is stated. The rune version is correct for all
  Unicode; the byte version is a narrower tool.

## Slice aliasing

- `b := a[lo:hi]` creates a new slice header but shares the **same backing
  array** as `a`. Writes through `b` are visible in `a` and vice versa.
- `append(b, x)` may or may not reallocate. If `cap(b) > len(b)`, the write
  lands in the shared array and can silently overwrite `a`'s data past the
  original slice.
- To decouple, copy: `c := append([]int(nil), a[lo:hi]...)` or `copy(dst, src)`.

**When teaching:** any problem that slices and then mutates needs to state
whether the input may be modified. Do not use `append` on a shared slice in
example code without explaining the aliasing risk.

## Map zero values and nil maps

- A `nil` map reads fine (`v, ok := m[k]` returns the zero value and `ok=false`)
  but **writing to a nil map panics**.
- Declaring `var m map[string]int` gives a nil map. Use `m := map[string]int{}`
  or `m := make(map[string]int)` before writing.
- Map iteration order is **intentionally randomized** per-range — never rely on
  insertion order or sorted order.

**When teaching:** if a problem requires deterministic map iteration, sort the
keys or use a slice alongside the map. Do not show example code that assumes
order.

## Integer division and modulo

- `a / b` for integers truncates toward **zero**, not toward negative infinity.
  `(-7) / 2 == -3`, not `-4`.
- `a % b` has the sign of `a`. `(-7) % 2 == -1`, not `1`.
- For modular arithmetic on possibly-negative values, use `((a % b) + b) % b`.
- Integer overflow is silent — `int` is 64-bit on most platforms but finite. For
  sums of large inputs, consider overflow (`overflow-missed` in the taxonomy).

**When teaching:** if a problem's inputs can be negative or the math involves
modular wrapping, state the domain explicitly. Don't show `i % 2` as "parity
check" without noting it behaves strangely for negatives.

## `append` returns a (possibly new) slice

- Always use the return value: `s = append(s, x)`. Discarding it loses the
  update when a realloc happened.
- Growth is amortized, not linear — the capacity typically doubles. Don't teach
  "append is O(1)" without the amortized caveat.

## Range over slices and maps

- `for i, v := range s` — `v` is a **copy** of the element. Mutating `v` does
  not change `s[i]`. To mutate in place, use `s[i]` directly.
- When `v` is a struct, the copy can be expensive. Index-based loops avoid it.

## Rune literals vs. string literals

- `'a'` is a `rune` (int32). `"a"` is a string.
- `s == "a"` compares strings; `s[0] == 'a'` compares byte-to-rune (valid for
  ASCII; Go will auto-convert the byte to rune for comparison).

---

## Teaching discipline

Two rules derived from the above:

1. **State the contract.** Every problem file must declare the input domain
   when the solution depends on it (ASCII vs Unicode, non-empty, non-negative,
   sorted, bounded, no-duplicates). If the contract is missing, the user's
   "overcomplicated" solution may actually be the correct general one.

2. **Verify the nudge.** Before nudging the user toward a simpler form, confirm
   that form is correct under the stated contract. If the simpler form is only
   correct in a narrower domain, either widen the problem (let the general
   solution stand) or narrow the contract (and then nudge). Never nudge a user
   away from a more-correct solution toward a narrower one.
