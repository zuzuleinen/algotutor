package review

import (
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"
)

func (m Model) viewWelcome() string {
	now := time.Now()
	total, due, newCount := m.store.Stats(now)
	byConcept := m.store.DueByConcept(now)

	title := titleStyle.Render("AlgoTutor Review")

	statsLine := subtitleStyle.Render(fmt.Sprintf("%d cards due today", due))
	totalLine := mutedStyle.Render(fmt.Sprintf("%d total cards  ·  %d new", total, newCount))

	var conceptLines []string
	concepts := sortedKeys(byConcept)
	for _, concept := range concepts {
		count := byConcept[concept]
		dots := strings.Repeat("·", max(1, 14-len(concept)))
		line := fmt.Sprintf("  %s %s %d due",
			conceptStyle.Render(concept),
			mutedStyle.Render(dots),
			count,
		)
		conceptLines = append(conceptLines, line)
	}

	var content strings.Builder
	content.WriteString("\n")
	content.WriteString(lipgloss.PlaceHorizontal(56, lipgloss.Center, title))
	content.WriteString("\n\n")
	content.WriteString(lipgloss.PlaceHorizontal(56, lipgloss.Center, statsLine))
	content.WriteString("\n")
	content.WriteString(lipgloss.PlaceHorizontal(56, lipgloss.Center, totalLine))
	content.WriteString("\n\n")
	if len(conceptLines) > 0 {
		content.WriteString(lipgloss.PlaceHorizontal(56, lipgloss.Center,
			mutedStyle.Render("By concept:")))
		content.WriteString("\n")
		for _, line := range conceptLines {
			content.WriteString(line)
			content.WriteString("\n")
		}
		content.WriteString("\n")
	}
	content.WriteString(lipgloss.PlaceHorizontal(56, lipgloss.Center,
		mutedStyle.Render("Press Enter to start")))
	content.WriteString("\n")

	box := cardStyle.Width(58).Render(content.String())
	return m.center(box)
}

func (m Model) viewFront() string {
	card := m.dueCards[m.current]
	total := len(m.dueCards)

	header := lipgloss.JoinHorizontal(lipgloss.Top,
		conceptStyle.Render(card.Concept),
		strings.Repeat(" ", max(1, 46-len(card.Concept))),
		mutedStyle.Render(fmt.Sprintf("%d / %d", m.current+1, total)),
	)

	question := questionStyle.Render(card.Front)
	answerLabel := mutedStyle.Render("Your answer:")
	box := cardStyle.Width(58).Render(
		"\n" + question + "\n\n" +
			answerLabel + "\n" +
			m.textarea.View() + "\n",
	)

	progress := m.renderProgress(m.current, total, 50)

	hint := mutedStyle.Render("Tab to reveal")

	return m.center(
		header + "\n\n" +
			box + "\n\n" +
			"  " + progress + "\n\n" +
			lipgloss.PlaceHorizontal(58, lipgloss.Center, hint),
	)
}

func (m Model) viewBack() string {
	card := m.dueCards[m.current]
	total := len(m.dueCards)

	header := lipgloss.JoinHorizontal(lipgloss.Top,
		conceptStyle.Render(card.Concept),
		strings.Repeat(" ", max(1, 46-len(card.Concept))),
		mutedStyle.Render(fmt.Sprintf("%d / %d", m.current+1, total)),
	)

	question := questionStyle.Render(card.Front)
	divider := dividerStyle.Render(strings.Repeat("─", 50))

	var answerSection string
	if m.userAnswer != "" {
		userLabel := mutedStyle.Render("Your answer:")
		userAns := answerStyle.Render(m.userAnswer)
		correctLabel := mutedStyle.Render("Correct answer:")
		correctAns := answerStyle.Render(card.Back)
		answerSection = userLabel + "\n" + userAns + "\n\n" + correctLabel + "\n" + correctAns
	} else {
		correctAns := answerStyle.Render(card.Back)
		answerSection = correctAns
	}

	box := cardStyle.Width(58).Render(
		"\n" + question + "\n" +
			divider + "\n" +
			answerSection + "\n",
	)

	progress := m.renderProgress(m.current, total, 50)

	ratings := lipgloss.JoinHorizontal(lipgloss.Top,
		againStyle.Render("1 Again"),
		"    ",
		hardStyle.Render("2 Hard"),
		"    ",
		goodStyle.Render("3 Good"),
		"    ",
		easyStyle.Render("4 Easy"),
	)

	return m.center(
		header + "\n\n" +
			box + "\n\n" +
			"  " + progress + "\n\n" +
			lipgloss.PlaceHorizontal(58, lipgloss.Center, ratings),
	)
}

func (m Model) viewSummary() string {
	title := titleStyle.Render("Session Complete!")

	reviewed := subtitleStyle.Render(fmt.Sprintf("%d cards reviewed", m.stats.reviewed))

	maxRating := max(m.stats.again, max(m.stats.hard, max(m.stats.good, m.stats.easy)))
	if maxRating == 0 {
		maxRating = 1
	}

	barWidth := 20
	bars := []string{
		statsLabelStyle.Render("Again ") + statsBarAgain.Render(bar(m.stats.again, maxRating, barWidth)) + fmt.Sprintf("  %d", m.stats.again),
		statsLabelStyle.Render("Hard ") + statsBarHard.Render(bar(m.stats.hard, maxRating, barWidth)) + fmt.Sprintf("  %d", m.stats.hard),
		statsLabelStyle.Render("Good ") + statsBarGood.Render(bar(m.stats.good, maxRating, barWidth)) + fmt.Sprintf("  %d", m.stats.good),
		statsLabelStyle.Render("Easy ") + statsBarEasy.Render(bar(m.stats.easy, maxRating, barWidth)) + fmt.Sprintf("  %d", m.stats.easy),
	}

	hint := mutedStyle.Render("Press Enter to exit")

	var content strings.Builder
	content.WriteString("\n")
	content.WriteString(lipgloss.PlaceHorizontal(56, lipgloss.Center, title))
	content.WriteString("\n\n")
	content.WriteString(lipgloss.PlaceHorizontal(56, lipgloss.Center, reviewed))
	content.WriteString("\n\n")
	for _, b := range bars {
		content.WriteString("  " + b + "\n")
	}
	content.WriteString("\n")
	content.WriteString(lipgloss.PlaceHorizontal(56, lipgloss.Center, hint))
	content.WriteString("\n")

	box := cardStyle.Width(58).Render(content.String())
	return m.center(box)
}

func (m Model) renderProgress(current, total, width int) string {
	if total == 0 {
		return ""
	}
	filled := (current * width) / total
	empty := width - filled
	return progressFullStyle.Render(strings.Repeat("█", filled)) +
		progressEmptyStyle.Render(strings.Repeat("░", empty))
}

func (m Model) center(content string) string {
	return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, content)
}

func bar(value, maxVal, width int) string {
	if maxVal == 0 {
		return strings.Repeat("░", width)
	}
	filled := (value * width) / maxVal
	empty := width - filled
	return strings.Repeat("█", filled) + strings.Repeat("░", empty)
}

func sortedKeys(m map[string]int) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}
