// Command start switches the active course and auto-launches the configured
// AI agent with a starting prompt (`train` or `review`).
//
//	go run ./cmd/start train          # use default course, send `train`
//	go run ./cmd/start train conc     # set active=conc, send `train`
//	go run ./cmd/start review         # send `review` for active course
package main

import (
	"errors"
	"fmt"
	"os"
	"os/exec"

	"algotutor/internal/courses"
	"algotutor/internal/launch"
	"algotutor/internal/migrate"
)

func main() {
	if err := run(os.Args[1:]); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run(args []string) error {
	if len(args) == 0 {
		return errors.New("usage: start <train|review> [course]")
	}
	prompt := args[0]
	if prompt != "train" && prompt != "review" {
		return fmt.Errorf("unknown flow: %s (want train or review)", prompt)
	}

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

	if len(args) > 1 {
		slug := args[1]
		if !state.IsEnrolled(slug) {
			return fmt.Errorf("course %q is not enrolled (enrolled: %v)", slug, state.Enrolled)
		}
		if err := state.SetActive(slug); err != nil {
			return err
		}
		if err := state.Save(); err != nil {
			return err
		}
	}

	if state.Active == "" {
		if state.Default == "" {
			return errors.New("no active course in state.json — run `make init`")
		}
		state.Active = state.Default
		if err := state.Save(); err != nil {
			return err
		}
	}

	if state.DefaultAgent == nil || *state.DefaultAgent == "" {
		fmt.Printf("Active course: %s\nOpen your AI agent in this directory and type `%s`.\n", state.Active, prompt)
		return nil
	}

	cmd, err := launch.Command(*state.DefaultAgent, prompt)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not build launch command: %v\n", err)
		fmt.Printf("Active course: %s\nOpen your AI agent in this directory and type `%s`.\n", state.Active, prompt)
		return nil
	}
	c := exec.Command(cmd[0], cmd[1:]...)
	c.Stdin = os.Stdin
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	return c.Run()
}
