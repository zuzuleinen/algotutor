<h2 align="center">
    <picture>
        <source media="(prefers-color-scheme: dark)" srcset="logo-dark.svg">
        <source media="(prefers-color-scheme: light)" srcset="logo-light.svg">
        <img height="80" alt="algotutor" src="logo.svg">
    </picture>
    <br>
    AI-powered training for Go developers
</h2>

<div align="center">

[![Agent-agnostic](https://img.shields.io/badge/Agent-agnostic-FA9BFA?style=flat)](docs/agents.md)
[![Go](https://img.shields.io/badge/Go-1.26-4B78E6?style=flat&logo=go&logoColor=white)](https://go.dev/)
[![Spaced repetition](https://img.shields.io/badge/Review-FSRS-73DC8C?style=flat)](https://github.com/open-spaced-repetition/go-fsrs)

**[algotutor.ai](https://algotutor.ai)**

</div>

Your agent acts as your personal tutor, generating progressively harder Go problems. It tracks your skill level on each concept and picks the next problem
based on where you are. 

Current courses: **Algorithms & Data Structures** and **Go Concurrency**.

<div align="center">

<img src="img_2.png" width="700" alt="algotutor in action"/>

</div>

## Get started

Clone the project and `cd` into it:

```bash
git clone git@github.com:zuzuleinen/algotutor.git && cd algotutor
```

Run initial set-up to choose your agent and optionally enroll in a course:

```bash
make init
```

## Train > Check > Review

Start a training session with `make train`. 

Submit your solution by telling your tutor `check`. It's OK to say, `"I don't know"`. They will pick an easier problem for you.

Start a daily review session with `make review`.

## Commands

### Local commands — terminal

| Command                | What it does                                                              |
|------------------------|---------------------------------------------------------------------------|
| `make init`            | First-time setup — enroll in courses, pick a default agent                |
| `make enroll`          | Add another course to your enrollment                                     |
| `make train`           | Launch your agent in training mode for the active course                  |
| `make train <slug>`    | Switch active course to `<slug>` and launch the agent there               |
| `make review`          | Open the review TUI across every enrolled course                          |
| `make review <slug>`   | Open the review TUI scoped to one course                                  |
| `make list`            | List all available courses and your enrollment status                     |
| `make run`             | Sanity-check your solution before `check` (active-course aware)           |

### Agent commands — chat

Type these to your agent once it's running:

| Command                          | What it does                                                                  |
|----------------------------------|-------------------------------------------------------------------------------|
| `check`                          | Submit your solution for evaluation (grading, mistake logging, level updates) |
| `I don't know`                   | Break the problem into simpler sub-problems                                   |
| `I want to solve [problem desc]` | Request a specific problem with description or screenshot                     |
| `reset`                          | Wipe progress in the active course (with `confirm reset` gate)                |
| `reset all`                      | Wipe progress in every enrolled course (with `confirm reset all` gate)        |


## Requirements

- An AI coding agent — see [docs/agents.md](docs/agents.md) for the supported list and setup
- [Go](https://go.dev/) ≥ 1.26
