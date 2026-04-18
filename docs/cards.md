# Spaced Repetition Cards

During practice, create review cards that capture what the user just learned. Cards are stored in `cards.json` and
reviewed via `go run ./cmd/review`.

## When to Create Cards

- After a problem is marked **solved** during checking.
- After a scaffolding sub-problem is solved (capture the specific gap).
- Do NOT create cards for things the user clearly already knew.

## How to Create Cards

1. Read `cards.json` (if it exists). If it doesn't exist, start with `{"cards": []}`.
2. Identify 1–5 atomic things the user learned: algorithmic patterns, Go syntax, data structure properties, prerequisite
   concepts.
3. Formulate each card following the rules below.
4. Append the new card objects to the `cards` array.
5. Write back the complete file.
6. Tell the user briefly, e.g. "3 review cards created."

## Card JSON Format

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

## Card Formulation Rules (SuperMemo 20 Rules)

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

## Example Cards

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

## Avoiding Duplicates

Before creating cards, scan existing cards in `cards.json`. Do not create a card if an existing card already covers the
same fact (same question or equivalent knowledge). It is fine to create cards on the same concept from different angles.
