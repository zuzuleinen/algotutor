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

	defaultAgent, err := promptAgent()
	if err != nil {
		return err
	}

	pickedCourse, err := promptCourse()
	if err != nil {
		return err
	}

	var enrollSlugs []string
	if pickedCourse != courseLater {
		enrollSlugs = []string{pickedCourse}
	}

	state := &courses.State{}
	for _, slug := range enrollSlugs {
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

	for _, slug := range enrollSlugs {
		if err := courses.EnsureCourseDir(slug); err != nil {
			return fmt.Errorf("init course %s: %w", slug, err)
		}
	}

	if pickedCourse == courseLater {
		if defaultAgent != "" {
			fmt.Printf("\nDefault agent: %s\n", defaultAgent)
		}
		fmt.Println("\nReady. Run `make enroll` when you want to enroll in a course.")
		return nil
	}

	enrollLabels := make([]string, len(enrollSlugs))
	for i, slug := range enrollSlugs {
		enrollLabels[i] = courses.DisplayLabel(slug)
	}
	fmt.Printf("\nEnrolled in: %s\n", strings.Join(enrollLabels, ", "))
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

	if err := state.SetActive(pickedCourse); err != nil {
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
	fmt.Printf("Enrolled in %s. Run `make train %s` to start.\n", c.Label, picked)
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

// courseLater is the sentinel value returned by promptCourse when the user
// wants to defer the choice. No course is enrolled — the user runs
// `make enroll` to pick one when they're ready.
const courseLater = "__later__"

func promptCourse() (string, error) {
	options := make([]huh.Option[string], 0, len(courses.Known)+1)
	for _, c := range courses.Known {
		options = append(options, huh.NewOption(c.Name, c.Slug))
	}
	options = append(options, huh.NewOption("I'll decide later", courseLater))

	var picked string
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Pick a course to start with").
				Description("You can enroll in more anytime with `make enroll`.").
				Options(options...).
				Value(&picked),
		),
	)
	if err := form.Run(); err != nil {
		return "", err
	}
	return picked, nil
}

func promptAgent() (string, error) {
	available := launch.DetectAvailable()
	if len(available) == 0 {
		fmt.Println("\nNo AI agents detected on PATH. You'll launch manually.")
		return "", nil
	}

	var options []huh.Option[string]
	for _, a := range available {
		options = append(options, huh.NewOption(a.DisplayName, a.Slug))
	}

	var picked string
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Default AI agent for your tutor?").
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
