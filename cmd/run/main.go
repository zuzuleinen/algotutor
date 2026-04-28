// Command run is the per-course validator dispatcher.
//
// Reads state.json, looks up the active course's Validator, and execs it.
// Local sanity check before saying `check` to the agent.
//
//	go run ./cmd/run            # validate active course
//	make run                    # same, via Makefile
package main

import (
	"errors"
	"fmt"
	"os"
	"os/exec"

	"algotutor/internal/courses"
	"algotutor/internal/migrate"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run() error {
	if err := migrate.MaybeRun(); err != nil && !migrate.IsNoLegacyState(err) {
		return fmt.Errorf("migration failed: %w", err)
	}

	state, err := courses.Load()
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return errors.New("no state.json — run `make init` first")
		}
		return err
	}
	if state.Active == "" {
		return errors.New("no active course set; run `make train`")
	}

	course, ok := courses.LookupKnown(state.Active)
	if !ok {
		return fmt.Errorf("active course %q is not a known course", state.Active)
	}
	if len(course.Validator) == 0 {
		return fmt.Errorf("course %q has no validator configured", state.Active)
	}

	c := exec.Command(course.Validator[0], course.Validator[1:]...)
	c.Stdin = os.Stdin
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	// Surface the underlying command's exit code unchanged.
	if err := c.Run(); err != nil {
		var exitErr *exec.ExitError
		if errors.As(err, &exitErr) {
			os.Exit(exitErr.ExitCode())
		}
		return err
	}
	return nil
}
