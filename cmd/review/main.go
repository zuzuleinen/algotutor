package main

import (
	"errors"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"

	"algotutor/internal/cards"
	"algotutor/internal/courses"
	"algotutor/internal/migrate"
	"algotutor/internal/review"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	if err := migrate.MaybeRun(); err != nil && !migrate.IsNoLegacyState(err) {
		fmt.Fprintf(os.Stderr, "migration failed: %v\n", err)
		os.Exit(1)
	}

	state, err := courses.Load()
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			fmt.Fprintln(os.Stderr, "no state.json — run `make init` to set up courses")
			os.Exit(1)
		}
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	scope, err := resolveScope(state, os.Args[1:])
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	now := time.Now()
	storeOf := map[*cards.ReviewCard]*cards.CardStore{}
	var dueCards []*cards.ReviewCard
	welcome := review.WelcomeStats{DueByConcept: map[string]int{}}

	for _, slug := range scope {
		path := courses.Path(slug, "cards.json")
		store, err := cards.Load(path)
		if err != nil {
			fmt.Fprintf(os.Stderr, "load %s: %v\n", path, err)
			os.Exit(1)
		}
		total, due, newCount := store.Stats(now)
		welcome.Total += total
		welcome.Due += due
		welcome.New += newCount
		for concept, count := range store.DueByConcept(now) {
			welcome.DueByConcept[concept] += count
		}
		for _, c := range store.DueCards(now) {
			storeOf[c] = store
			dueCards = append(dueCards, c)
		}
	}

	if len(dueCards) == 0 {
		labels := make([]string, len(scope))
		for i, slug := range scope {
			labels[i] = courses.DisplayLabel(slug)
		}
		label := strings.Join(labels, " + ")
		if welcome.Total == 0 {
			fmt.Printf("No review cards in %s yet. Solve some problems first!\n", label)
		} else {
			fmt.Printf("No cards due in %s. Come back later!\n", label)
		}
		return
	}

	rand.Shuffle(len(dueCards), func(i, j int) {
		dueCards[i], dueCards[j] = dueCards[j], dueCards[i]
	})

	m := review.NewModel(storeOf, dueCards, welcome)
	p := tea.NewProgram(m, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

// resolveScope returns the list of course slugs to review.
//   - With no argument: all enrolled courses (interleaved review).
//   - With a course argument: just that course.
func resolveScope(state *courses.State, args []string) ([]string, error) {
	if len(args) > 0 {
		slug := args[0]
		if !state.IsEnrolled(slug) {
			return nil, fmt.Errorf("course %q is not enrolled (enrolled: %v)", slug, state.Enrolled)
		}
		return []string{slug}, nil
	}
	if len(state.Enrolled) == 0 {
		return nil, errors.New("no courses enrolled — run `make init`")
	}
	scope := make([]string, len(state.Enrolled))
	copy(scope, state.Enrolled)
	return scope, nil
}
