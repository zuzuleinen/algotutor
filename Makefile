.PHONY: init enroll train review run sync-agents check-agents

# Lets `make train conc` and `make review conc` accept positional args. The
# `%:` rule swallows the extra arg as a no-op target so make doesn't complain.
%:
	@:

init:
	go run ./cmd/init

enroll:
	go run ./cmd/init -enroll

train:
	@go run ./cmd/start train $(filter-out train,$(MAKECMDGOALS))

review:
	@go run ./cmd/review $(filter-out review,$(MAKECMDGOALS))

run:
	@go run ./cmd/run

sync-agents:
	cp AGENTS.md CLAUDE.md
	cp AGENTS.md GEMINI.md

check-agents:
	@diff -q AGENTS.md CLAUDE.md > /dev/null || { echo "CLAUDE.md out of sync — run 'make sync-agents'"; exit 1; }
	@diff -q AGENTS.md GEMINI.md > /dev/null || { echo "GEMINI.md out of sync — run 'make sync-agents'"; exit 1; }
	@echo "Agent mirrors in sync."
