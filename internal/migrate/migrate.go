// Package migrate moves pre-multi-course state into courses/algos/.
//
// Triggers when state.json is missing AND any root-level user-state file
// exists (progress.md, cards.json, etc.). Idempotent: a second run is a no-op
// because state.json now exists.
package migrate

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"algotutor/internal/courses"
)

// rootFiles is the set of files we move from the repo root into
// courses/algos/. Each entry maps source -> destination relative path.
var rootFiles = map[string]string{
	"progress.md":          "progress.md",
	"current.md":           "current.md",
	"cards.json":           "cards.json",
	"mistakes.json":        "mistakes.json",
	"resolve.json":         "resolve.json",
	"retention.json":       "retention.json",
	"mix.json":             "mix.json",
	"problem-bank.md":      "problem-bank.md",
	"progress.template.md": "progress.template.md",
}

// rootDirs is the set of directories we move from the repo root.
var rootDirs = map[string]string{
	"problems": "problems",
}

// rootDocs is the set of docs we move from the repo root's docs/ to the
// course's docs/. Project-wide docs (agents.md) and shared mechanics
// (cards.md, mix.md, resolve.md) stay at the root.
var rootDocs = map[string]string{
	"docs/concepts.md":   "docs/concepts.md",
	"docs/go-gotchas.md": "docs/go-gotchas.md",
	"docs/mistakes.md":   "docs/mistakes.md",
}

// MaybeRun executes the migration if pre-multi-course layout is detected. It
// returns nil if no migration was needed. It does NOT print anything — silent
// migration per design.
func MaybeRun() error {
	// If state.json already exists, we're past migration.
	if _, err := os.Stat(courses.StateFile); err == nil {
		return nil
	}
	// Detect a legacy file. If none, also nothing to do — fresh checkout.
	legacy := false
	for src := range rootFiles {
		if _, err := os.Stat(src); err == nil {
			legacy = true
			break
		}
	}
	if !legacy {
		for src := range rootDirs {
			if info, err := os.Stat(src); err == nil && info.IsDir() {
				legacy = true
				break
			}
		}
	}
	if !legacy {
		// Nothing to migrate. Caller should run init.
		return errNoLegacyState
	}
	return run()
}

// errNoLegacyState is returned when no migration is needed and no state.json
// exists — the caller should run init.
var errNoLegacyState = errors.New("no legacy state to migrate; run init")

// IsNoLegacyState reports whether err indicates the caller should run init.
func IsNoLegacyState(err error) bool {
	return errors.Is(err, errNoLegacyState)
}

func run() error {
	dst := courses.Dir("algos")
	if err := os.MkdirAll(filepath.Join(dst, "docs"), 0o755); err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Join(dst, "problems"), 0o755); err != nil {
		return err
	}
	for src, rel := range rootFiles {
		if err := moveFile(src, filepath.Join(dst, rel)); err != nil {
			return fmt.Errorf("move %s: %w", src, err)
		}
	}
	for src, rel := range rootDirs {
		if err := moveDir(src, filepath.Join(dst, rel)); err != nil {
			return fmt.Errorf("move %s: %w", src, err)
		}
	}
	for src, rel := range rootDocs {
		if err := moveFile(src, filepath.Join(dst, rel)); err != nil {
			return fmt.Errorf("move %s: %w", src, err)
		}
	}
	// Drop the now-empty docs/ entries that were course-specific. Keep the
	// directory itself (it still hosts agents.md, cards.md, mix.md, resolve.md).
	state := &courses.State{
		Enrolled: []string{"algos"},
		Active:   "algos",
		Default:  "algos",
	}
	return state.Save()
}

func moveFile(src, dst string) error {
	if _, err := os.Stat(src); errors.Is(err, os.ErrNotExist) {
		return nil
	} else if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(dst), 0o755); err != nil {
		return err
	}
	if err := os.Rename(src, dst); err == nil {
		return nil
	}
	// Cross-device or other rename failure → copy + remove.
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()
	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	if _, err := io.Copy(out, in); err != nil {
		out.Close()
		return err
	}
	if err := out.Close(); err != nil {
		return err
	}
	return os.Remove(src)
}

func moveDir(src, dst string) error {
	info, err := os.Stat(src)
	if errors.Is(err, os.ErrNotExist) {
		return nil
	} else if err != nil {
		return err
	}
	if !info.IsDir() {
		return fmt.Errorf("%s is not a directory", src)
	}
	if err := os.MkdirAll(filepath.Dir(dst), 0o755); err != nil {
		return err
	}
	// Try a fast rename first.
	if err := os.Rename(src, dst); err == nil {
		return nil
	}
	// Fall back: copy each entry across, then remove the source.
	if err := os.MkdirAll(dst, 0o755); err != nil {
		return err
	}
	entries, err := os.ReadDir(src)
	if err != nil {
		return err
	}
	for _, e := range entries {
		s := filepath.Join(src, e.Name())
		d := filepath.Join(dst, e.Name())
		if e.IsDir() {
			if err := moveDir(s, d); err != nil {
				return err
			}
		} else {
			if err := moveFile(s, d); err != nil {
				return err
			}
		}
	}
	return os.RemoveAll(src)
}
