.PHONY: review sync-agents check-agents

review:
	go run ./cmd/review

sync-agents:
	cp AGENTS.md CLAUDE.md
	cp AGENTS.md GEMINI.md

check-agents:
	@diff -q AGENTS.md CLAUDE.md > /dev/null || { echo "CLAUDE.md out of sync — run 'make sync-agents'"; exit 1; }
	@diff -q AGENTS.md GEMINI.md > /dev/null || { echo "GEMINI.md out of sync — run 'make sync-agents'"; exit 1; }
	@echo "Agent mirrors in sync."
