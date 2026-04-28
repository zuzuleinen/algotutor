package migrate

import (
	"os"
	"path/filepath"
	"testing"

	"algotutor/internal/courses"
)

func TestMaybeRunNoLegacyState(t *testing.T) {
	dir := t.TempDir()
	chdir(t, dir)

	err := MaybeRun()
	if !IsNoLegacyState(err) {
		t.Fatalf("want IsNoLegacyState, got %v", err)
	}
}

func TestMaybeRunSkipWhenStateExists(t *testing.T) {
	dir := t.TempDir()
	chdir(t, dir)

	if err := os.WriteFile(filepath.Join(dir, courses.StateFile), []byte("{}"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "progress.md"), []byte("legacy"), 0o644); err != nil {
		t.Fatal(err)
	}

	if err := MaybeRun(); err != nil {
		t.Fatalf("MaybeRun returned err with state.json present: %v", err)
	}
	// progress.md must NOT have moved.
	if _, err := os.Stat(filepath.Join(dir, "progress.md")); err != nil {
		t.Fatalf("progress.md was moved despite state.json present")
	}
}

func TestMaybeRunMovesLegacyLayout(t *testing.T) {
	dir := t.TempDir()
	chdir(t, dir)

	mustWrite(t, filepath.Join(dir, "progress.md"), "progress")
	mustWrite(t, filepath.Join(dir, "current.md"), "")
	mustWrite(t, filepath.Join(dir, "cards.json"), "[]")
	mustWrite(t, filepath.Join(dir, "problem-bank.md"), "bank")
	mustWrite(t, filepath.Join(dir, "progress.template.md"), "template")
	if err := os.MkdirAll(filepath.Join(dir, "docs"), 0o755); err != nil {
		t.Fatal(err)
	}
	mustWrite(t, filepath.Join(dir, "docs", "concepts.md"), "concepts")
	mustWrite(t, filepath.Join(dir, "docs", "go-gotchas.md"), "gotchas")
	mustWrite(t, filepath.Join(dir, "docs", "mistakes.md"), "mistakes")
	if err := os.MkdirAll(filepath.Join(dir, "problems"), 0o755); err != nil {
		t.Fatal(err)
	}
	mustWrite(t, filepath.Join(dir, "problems", "001.md"), "p1")

	if err := MaybeRun(); err != nil {
		t.Fatalf("MaybeRun: %v", err)
	}

	// state.json should now exist and be enrolled in algos.
	state, err := courses.Load()
	if err != nil {
		t.Fatalf("load state: %v", err)
	}
	if state.Active != "algos" || state.Default != "algos" {
		t.Errorf("state active/default = %s/%s; want algos/algos", state.Active, state.Default)
	}
	if !state.IsEnrolled("algos") {
		t.Errorf("algos not enrolled after migration")
	}

	// Files should have moved into courses/algos/.
	wantedAt := []string{
		"courses/algos/progress.md",
		"courses/algos/current.md",
		"courses/algos/cards.json",
		"courses/algos/problem-bank.md",
		"courses/algos/progress.template.md",
		"courses/algos/docs/concepts.md",
		"courses/algos/docs/go-gotchas.md",
		"courses/algos/docs/mistakes.md",
		"courses/algos/problems/001.md",
	}
	for _, p := range wantedAt {
		if _, err := os.Stat(filepath.Join(dir, p)); err != nil {
			t.Errorf("%s missing after migration: %v", p, err)
		}
	}

	// Originals should be gone.
	gone := []string{
		"progress.md", "current.md", "cards.json", "problem-bank.md",
		"progress.template.md", "docs/concepts.md", "docs/go-gotchas.md",
		"docs/mistakes.md", "problems/001.md",
	}
	for _, p := range gone {
		if _, err := os.Stat(filepath.Join(dir, p)); err == nil {
			t.Errorf("%s still at root after migration", p)
		}
	}
}

func chdir(t *testing.T, dir string) {
	t.Helper()
	prev, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	if err := os.Chdir(dir); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		_ = os.Chdir(prev)
	})
}

func mustWrite(t *testing.T, path, contents string) {
	t.Helper()
	if err := os.WriteFile(path, []byte(contents), 0o644); err != nil {
		t.Fatal(err)
	}
}
