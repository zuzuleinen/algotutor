package cards

import (
	"encoding/json"
	"os"
	"time"

	fsrs "github.com/open-spaced-repetition/go-fsrs/v4"
)

type ReviewEntry struct {
	Rating int       `json:"rating"`
	Review time.Time `json:"review"`
	State  int       `json:"state"`
}

type FSRSData struct {
	Due           time.Time `json:"due"`
	Stability     float64   `json:"stability"`
	Difficulty    float64   `json:"difficulty"`
	ElapsedDays   uint64    `json:"elapsed_days"`
	ScheduledDays uint64    `json:"scheduled_days"`
	Reps          uint64    `json:"reps"`
	Lapses        uint64    `json:"lapses"`
	State         int       `json:"state"`
	LastReview    time.Time `json:"last_review"`
}

func (d *FSRSData) ToFSRSCard() fsrs.Card {
	return fsrs.Card{
		Due:           d.Due,
		Stability:     d.Stability,
		Difficulty:    d.Difficulty,
		ElapsedDays:   d.ElapsedDays,
		ScheduledDays: d.ScheduledDays,
		Reps:          d.Reps,
		Lapses:        d.Lapses,
		State:         fsrs.State(d.State),
		LastReview:    d.LastReview,
	}
}

func (d *FSRSData) FromFSRSCard(c fsrs.Card) {
	d.Due = c.Due
	d.Stability = c.Stability
	d.Difficulty = c.Difficulty
	d.ElapsedDays = c.ElapsedDays
	d.ScheduledDays = c.ScheduledDays
	d.Reps = c.Reps
	d.Lapses = c.Lapses
	d.State = int(c.State)
	d.LastReview = c.LastReview
}

type ReviewCard struct {
	ID            string        `json:"id"`
	Front         string        `json:"front"`
	Back          string        `json:"back"`
	Concept       string        `json:"concept"`
	SourceProblem string        `json:"source_problem"`
	CreatedAt     time.Time     `json:"created_at"`
	FSRS          FSRSData      `json:"fsrs"`
	ReviewLog     []ReviewEntry `json:"review_log"`
}

type CardStore struct {
	Cards []ReviewCard `json:"cards"`
	path  string
}

func Load(path string) (*CardStore, error) {
	store := &CardStore{path: path}

	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return store, nil
		}
		return nil, err
	}

	if err := json.Unmarshal(data, store); err != nil {
		return nil, err
	}
	return store, nil
}

func (s *CardStore) Save() error {
	data, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(s.path, data, 0644)
}

func (s *CardStore) DueCards(now time.Time) []*ReviewCard {
	var due []*ReviewCard
	for i := range s.Cards {
		if !s.Cards[i].FSRS.Due.After(now) {
			due = append(due, &s.Cards[i])
		}
	}
	return due
}

func (s *CardStore) Stats(now time.Time) (total, due, newCount int) {
	total = len(s.Cards)
	for i := range s.Cards {
		if !s.Cards[i].FSRS.Due.After(now) {
			due++
		}
		if s.Cards[i].FSRS.State == int(fsrs.New) {
			newCount++
		}
	}
	return
}

func (s *CardStore) DueByConcept(now time.Time) map[string]int {
	counts := make(map[string]int)
	for i := range s.Cards {
		if !s.Cards[i].FSRS.Due.After(now) {
			counts[s.Cards[i].Concept]++
		}
	}
	return counts
}
