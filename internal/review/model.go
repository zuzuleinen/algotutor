package review

import (
	"time"

	"algotutor/internal/cards"

	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
	fsrs "github.com/open-spaced-repetition/go-fsrs/v4"
)

type screen int

const (
	screenWelcome screen = iota
	screenFront
	screenBack
	screenSummary
)

type sessionStats struct {
	reviewed int
	again    int
	hard     int
	good     int
	easy     int
}

type Model struct {
	screen     screen
	dueCards   []*cards.ReviewCard
	current    int
	stats      sessionStats
	width      int
	height     int
	store      *cards.CardStore
	fsrs       *fsrs.FSRS
	textarea   textarea.Model
	userAnswer string
}

func newTextarea() textarea.Model {
	ta := textarea.New()
	ta.Placeholder = "Type your answer..."
	ta.ShowLineNumbers = false
	ta.SetWidth(50)
	ta.SetHeight(3)
	ta.FocusedStyle.CursorLine = ta.FocusedStyle.CursorLine.Copy().UnsetBackground()
	ta.Focus()
	return ta
}

func NewModel(store *cards.CardStore, dueCards []*cards.ReviewCard) Model {
	return Model{
		screen:   screenWelcome,
		dueCards: dueCards,
		store:    store,
		fsrs:     fsrs.NewFSRS(fsrs.DefaultParam()),
		textarea: newTextarea(),
	}
}

func (m Model) Init() tea.Cmd {
	return textarea.Blink
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		}

		if m.screen == screenWelcome || m.screen == screenSummary {
			if msg.String() == "q" || msg.String() == "esc" {
				return m, tea.Quit
			}
		}

		switch m.screen {
		case screenWelcome:
			return m.updateWelcome(msg)
		case screenFront:
			return m.updateFront(msg)
		case screenBack:
			return m.updateBack(msg)
		case screenSummary:
			return m.updateSummary(msg)
		}
	}

	if m.screen == screenFront {
		var cmd tea.Cmd
		m.textarea, cmd = m.textarea.Update(msg)
		return m, cmd
	}

	return m, nil
}

func (m Model) updateWelcome(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	if msg.String() == "enter" {
		m.screen = screenFront
		m.textarea.Reset()
		m.textarea.Focus()
		return m, textarea.Blink
	}
	return m, nil
}

func (m Model) updateFront(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "tab":
		m.userAnswer = m.textarea.Value()
		m.screen = screenBack
		return m, nil
	}

	var cmd tea.Cmd
	m.textarea, cmd = m.textarea.Update(msg)
	return m, cmd
}

func (m Model) updateBack(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	var rating fsrs.Rating
	switch msg.String() {
	case "1":
		rating = fsrs.Again
		m.stats.again++
	case "2":
		rating = fsrs.Hard
		m.stats.hard++
	case "3":
		rating = fsrs.Good
		m.stats.good++
	case "4":
		rating = fsrs.Easy
		m.stats.easy++
	default:
		return m, nil
	}

	card := m.dueCards[m.current]
	now := time.Now()

	fsrsCard := card.FSRS.ToFSRSCard()
	recordLog := m.fsrs.Repeat(fsrsCard, now)
	info := recordLog[rating]

	card.FSRS.FromFSRSCard(info.Card)
	card.ReviewLog = append(card.ReviewLog, cards.ReviewEntry{
		Rating: int(rating),
		Review: now,
		State:  int(info.ReviewLog.State),
	})

	m.stats.reviewed++
	_ = m.store.Save()

	m.current++
	if m.current >= len(m.dueCards) {
		m.screen = screenSummary
	} else {
		m.screen = screenFront
		m.textarea.Reset()
		m.textarea.Focus()
		return m, textarea.Blink
	}
	return m, nil
}

func (m Model) updateSummary(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	if msg.String() == "enter" || msg.String() == "q" {
		return m, tea.Quit
	}
	return m, nil
}

func (m Model) View() string {
	switch m.screen {
	case screenWelcome:
		return m.viewWelcome()
	case screenFront:
		return m.viewFront()
	case screenBack:
		return m.viewBack()
	case screenSummary:
		return m.viewSummary()
	}
	return ""
}
