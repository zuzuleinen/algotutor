// Command init runs the huh-driven onboarding flow for algotutor.
//
//	go run ./cmd/init           # first-time setup
//	go run ./cmd/init -enroll   # add another course to existing enrollment
package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"

	"algotutor/internal/courses"
	"algotutor/internal/launch"
	"algotutor/internal/migrate"

	"github.com/charmbracelet/huh"
)

const (
	requiredGoMajor = 1
	requiredGoMinor = 26
)

func main() {
	enroll := flag.Bool("enroll", false, "Add a course to existing enrollment")
	flag.Parse()

	var err error
	if *enroll {
		err = runEnroll()
	} else {
		err = runInit()
	}
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func runInit() error {
	if err := preflight(); err != nil {
		return err
	}

	if err := migrate.MaybeRun(); err != nil && !migrate.IsNoLegacyState(err) {
		return fmt.Errorf("migration failed: %w", err)
	}

	if _, err := os.Stat(courses.StateFile); err == nil {
		// Already initialized — but if no default agent is set, offer to pick one now.
		state, loadErr := courses.Load()
		if loadErr == nil && (state.DefaultAgent == nil || *state.DefaultAgent == "") {
			agent, pickErr := promptAgent()
			if pickErr != nil {
				return pickErr
			}
			if agent != "" {
				state.DefaultAgent = &agent
				if saveErr := state.Save(); saveErr != nil {
					return fmt.Errorf("save state.json: %w", saveErr)
				}
				fmt.Printf("Default agent set to %s.\n", agent)
			}
		} else {
			fmt.Println("Already initialized. Use `make enroll` to add a course.")
		}
		return nil
	}

	selectedCourses, err := promptCourses(nil)
	if err != nil {
		return err
	}
	if len(selectedCourses) == 0 {
		fmt.Println("No courses selected. Aborting.")
		return nil
	}

	defaultAgent, err := promptAgent()
	if err != nil {
		return err
	}

	state := &courses.State{}
	for _, slug := range selectedCourses {
		if err := state.Enroll(slug); err != nil {
			return err
		}
	}
	state.Active = "" // no course selected until the user explicitly trains
	if defaultAgent != "" {
		a := defaultAgent
		state.DefaultAgent = &a
	}
	if err := state.Save(); err != nil {
		return fmt.Errorf("save state.json: %w", err)
	}

	for _, slug := range selectedCourses {
		if err := courses.EnsureCourseDir(slug); err != nil {
			return fmt.Errorf("init course %s: %w", slug, err)
		}
	}

	fmt.Printf("\nEnrolled in: %s\n", strings.Join(selectedCourses, ", "))
	if defaultAgent != "" {
		fmt.Printf("Default agent: %s\n", defaultAgent)
	}

	trainNow, err := promptTrainNow()
	if err != nil {
		return err
	}
	if !trainNow {
		fmt.Println("\nReady. Run `make train` when you want to start.")
		return nil
	}

	pick := state.Default
	if len(state.Enrolled) > 1 {
		pick, err = promptPickCourse(state.Enrolled, "Which course do you want to train?")
		if err != nil {
			return err
		}
	}
	if err := state.SetActive(pick); err != nil {
		return err
	}
	if err := state.Save(); err != nil {
		return err
	}
	return launchOrInstruct(state, "train")
}

func runEnroll() error {
	state, err := courses.Load()
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return errors.New("no state.json — run `make init` first")
		}
		return err
	}

	available := make([]string, 0, len(courses.Known))
	for _, c := range courses.Known {
		if !state.IsEnrolled(c.Slug) {
			available = append(available, c.Slug)
		}
	}
	if len(available) == 0 {
		fmt.Println("You're enrolled in every available course.")
		return nil
	}

	picked, err := promptPickCourse(available, "Which course do you want to enroll in?")
	if err != nil {
		return err
	}
	if err := state.Enroll(picked); err != nil {
		return err
	}
	if err := courses.EnsureCourseDir(picked); err != nil {
		return err
	}
	if err := state.Save(); err != nil {
		return err
	}
	c, _ := courses.LookupKnown(picked)
	fmt.Printf("Enrolled in %s. Run `make train %s` to start.\n", c.Name, picked)
	return nil
}

func preflight() error {
	major, minor, err := parseGoVersion(runtime.Version())
	if err != nil {
		return fmt.Errorf("could not parse Go version %q: %w", runtime.Version(), err)
	}
	if major < requiredGoMajor || (major == requiredGoMajor && minor < requiredGoMinor) {
		return fmt.Errorf("algotutor requires Go %d.%d+, found Go %d.%d", requiredGoMajor, requiredGoMinor, major, minor)
	}
	return nil
}

func parseGoVersion(v string) (int, int, error) {
	v = strings.TrimPrefix(v, "go")
	parts := strings.Split(v, ".")
	if len(parts) < 2 {
		return 0, 0, fmt.Errorf("malformed version")
	}
	major, err := strconv.Atoi(parts[0])
	if err != nil {
		return 0, 0, err
	}
	minor, err := strconv.Atoi(parts[1])
	if err != nil {
		return 0, 0, err
	}
	return major, minor, nil
}

func promptCourses(preselect []string) ([]string, error) {
	preselectSet := map[string]bool{}
	if preselect == nil {
		// Pre-select all by default on first init.
		for _, c := range courses.Known {
			preselectSet[c.Slug] = true
		}
	} else {
		for _, s := range preselect {
			preselectSet[s] = true
		}
	}

	options := make([]huh.Option[string], 0, len(courses.Known))
	for _, c := range courses.Known {
		opt := huh.NewOption(c.Name, c.Slug)
		if preselectSet[c.Slug] {
			opt = opt.Selected(true)
		}
		options = append(options, opt)
	}

	var selected []string
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewMultiSelect[string]().
				Title("Pick the courses you want to enroll in").
				Description("Space to toggle, enter to confirm").
				Options(options...).
				Validate(func(s []string) error {
					if len(s) == 0 {
						return errors.New("pick at least one course")
					}
					return nil
				}).
				Value(&selected),
		),
	)
	if err := form.Run(); err != nil {
		return nil, err
	}
	return selected, nil
}

func promptAgent() (string, error) {
	available := launch.DetectAvailable()
	if len(available) == 0 {
		fmt.Println("\nNo AI agents detected on PATH. You'll launch manually.")
		return "", nil
	}

	options := []huh.Option[string]{huh.NewOption("None — I'll launch manually", "")}
	for _, a := range available {
		options = append(options, huh.NewOption(a.DisplayName, a.Slug))
	}

	var picked string
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Default AI agent for `make train` / `make review`?").
				Description("Used to auto-launch your agent in this directory.").
				Options(options...).
				Value(&picked),
		),
	)
	if err := form.Run(); err != nil {
		return "", err
	}
	return picked, nil
}

func promptTrainNow() (bool, error) {
	var v bool
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewConfirm().
				Title("Train now?").
				Affirmative("Yes, start training").
				Negative("Not now").
				Value(&v),
		),
	)
	if err := form.Run(); err != nil {
		return false, err
	}
	return v, nil
}

func promptPickCourse(slugs []string, title string) (string, error) {
	options := make([]huh.Option[string], len(slugs))
	for i, slug := range slugs {
		c, ok := courses.LookupKnown(slug)
		label := slug
		if ok {
			label = c.Name
		}
		options[i] = huh.NewOption(label, slug)
	}
	var picked string
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title(title).
				Options(options...).
				Value(&picked),
		),
	)
	if err := form.Run(); err != nil {
		return "", err
	}
	return picked, nil
}

func launchOrInstruct(state *courses.State, prompt string) error {
	if state.DefaultAgent == nil || *state.DefaultAgent == "" {
		fmt.Printf("\nOpen your AI agent in this directory and type `%s`.\n", prompt)
		return nil
	}
	cmd, err := launch.Command(*state.DefaultAgent, prompt)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not build launch command for %s: %v\n", *state.DefaultAgent, err)
		fmt.Printf("Open your AI agent in this directory and type `%s`.\n", prompt)
		return nil
	}
	c := exec.Command(cmd[0], cmd[1:]...)
	c.Stdin = os.Stdin
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	return c.Run()
}
