// Package courses manages multi-course state for algotutor.
//
// State lives in state.json at the repo root and tracks which courses the user
// has enrolled in, which course is currently active, the fallback default, and
// an optional default AI agent for auto-launch flows. All per-course files
// (progress.md, current.md, problems/, cards.json, etc.) live under
// courses/<slug>/.
package courses

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
)

// StateFile is the path to the root-level state file.
const StateFile = "state.json"

// CoursesDir is the path to the per-course directory tree.
const CoursesDir = "courses"

// Known is the canonical list of course slugs algotutor ships with. Order
// matters: the first entry is the install-time default.
var Known = []Course{
	{
		Slug:      "algos",
		Label:     "Algos",
		Name:      "Algorithms & Data Structures (Go)",
		Validator: []string{"go", "run", "."},
	},
	{
		Slug:      "conc",
		Label:     "Concurrency",
		Name:      "Go Concurrency",
		Validator: []string{"go", "test", "-race", "."},
	},
}

// Course identifies a course by its slug, human-readable name, and the shell
// command used to validate the user's working files. Validator is the argv
// `make run` invokes for problems in this course.
type Course struct {
	Slug      string
	Label     string   // short display name, e.g. "Algos", "Concurrency"
	Name      string
	Validator []string
}

// State is the on-disk schema for state.json.
type State struct {
	Enrolled     []string `json:"enrolled"`
	Active       string   `json:"active"`
	Default      string   `json:"default"`
	DefaultAgent *string  `json:"default_agent"`
}

// Load reads state.json. If the file is missing, returns os.ErrNotExist so
// callers can decide to migrate or run init.
func Load() (*State, error) {
	data, err := os.ReadFile(StateFile)
	if err != nil {
		return nil, err
	}
	var s State
	if err := json.Unmarshal(data, &s); err != nil {
		return nil, fmt.Errorf("parse %s: %w", StateFile, err)
	}
	return &s, nil
}

// Save writes the state atomically to state.json.
func (s *State) Save() error {
	data, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return err
	}
	data = append(data, '\n')
	tmp := StateFile + ".tmp"
	if err := os.WriteFile(tmp, data, 0o644); err != nil {
		return err
	}
	return os.Rename(tmp, StateFile)
}

// IsEnrolled reports whether the user has enrolled in slug.
func (s *State) IsEnrolled(slug string) bool {
	for _, e := range s.Enrolled {
		if e == slug {
			return true
		}
	}
	return false
}

// Enroll adds slug to Enrolled if not already present and keeps the slice
// sorted by Known order (so listing order is stable).
func (s *State) Enroll(slug string) error {
	if !IsKnown(slug) {
		return fmt.Errorf("unknown course slug: %s", slug)
	}
	if s.IsEnrolled(slug) {
		return nil
	}
	s.Enrolled = append(s.Enrolled, slug)
	sort.SliceStable(s.Enrolled, func(i, j int) bool {
		return knownIndex(s.Enrolled[i]) < knownIndex(s.Enrolled[j])
	})
	if s.Default == "" {
		s.Default = slug
	}
	if s.Active == "" {
		s.Active = slug
	}
	return nil
}

// SetActive sets the active course. The course must be enrolled.
func (s *State) SetActive(slug string) error {
	if !s.IsEnrolled(slug) {
		return fmt.Errorf("course %q is not enrolled", slug)
	}
	s.Active = slug
	return nil
}

// IsKnown reports whether slug is a known course.
func IsKnown(slug string) bool {
	return knownIndex(slug) >= 0
}

// LookupKnown returns the Course definition for slug, or false if unknown.
func LookupKnown(slug string) (Course, bool) {
	idx := knownIndex(slug)
	if idx < 0 {
		return Course{}, false
	}
	return Known[idx], true
}

// DisplayLabel returns the human-readable short label for slug (e.g. "Algos",
// "Concurrency"), falling back to the slug itself when unknown. Use this in
// user-facing output instead of printing slugs directly.
func DisplayLabel(slug string) string {
	if c, ok := LookupKnown(slug); ok {
		return c.Label
	}
	return slug
}

func knownIndex(slug string) int {
	for i, c := range Known {
		if c.Slug == slug {
			return i
		}
	}
	return -1
}

// Dir returns the directory path for course slug, e.g. courses/algos.
func Dir(slug string) string {
	return filepath.Join(CoursesDir, slug)
}

// Path returns a path inside the course's directory tree.
func Path(slug string, parts ...string) string {
	return filepath.Join(append([]string{CoursesDir, slug}, parts...)...)
}

// ActiveDir returns the directory for the currently-active course.
func (s *State) ActiveDir() string {
	return Dir(s.Active)
}

// ActivePath returns a path inside the active course's directory.
func (s *State) ActivePath(parts ...string) string {
	return Path(s.Active, parts...)
}

// EnsureCourseDir creates the directory tree for a course if it doesn't exist
// and seeds progress.md / current.md / problems/ from the template.
func EnsureCourseDir(slug string) error {
	if !IsKnown(slug) {
		return fmt.Errorf("unknown course slug: %s", slug)
	}
	dir := Dir(slug)
	if err := os.MkdirAll(filepath.Join(dir, "problems"), 0o755); err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Join(dir, "docs"), 0o755); err != nil {
		return err
	}
	template := filepath.Join(dir, "progress.template.md")
	progress := filepath.Join(dir, "progress.md")
	if _, err := os.Stat(progress); errors.Is(err, os.ErrNotExist) {
		if data, terr := os.ReadFile(template); terr == nil {
			if werr := os.WriteFile(progress, data, 0o644); werr != nil {
				return werr
			}
		}
	}
	current := filepath.Join(dir, "current.md")
	if _, err := os.Stat(current); errors.Is(err, os.ErrNotExist) {
		if werr := os.WriteFile(current, nil, 0o644); werr != nil {
			return werr
		}
	}
	return nil
}
