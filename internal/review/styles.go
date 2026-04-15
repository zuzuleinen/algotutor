package review

import "github.com/charmbracelet/lipgloss"

var (
	accentColor  = lipgloss.Color("215")
	borderColor  = lipgloss.Color("61")
	successColor = lipgloss.Color("78")
	warningColor = lipgloss.Color("221")
	dangerColor  = lipgloss.Color("203")
	mutedColor   = lipgloss.Color("241")
	brightColor  = lipgloss.Color("255")
	easyColor    = lipgloss.Color("117")
	codeColor    = lipgloss.Color("180")

	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(accentColor)

	subtitleStyle = lipgloss.NewStyle().
			Foreground(brightColor)

	mutedStyle = lipgloss.NewStyle().
			Foreground(mutedColor)

	conceptStyle = lipgloss.NewStyle().
			Foreground(accentColor).
			Background(lipgloss.Color("236")).
			Padding(0, 1)

	cardStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(borderColor).
			Padding(1, 3)

	questionStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(brightColor)

	answerStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("252"))

	codeStyle = lipgloss.NewStyle().
			Foreground(codeColor).
			Bold(true)

	dividerStyle = lipgloss.NewStyle().
			Foreground(mutedColor)

	againStyle = lipgloss.NewStyle().Foreground(dangerColor).Bold(true)
	hardStyle  = lipgloss.NewStyle().Foreground(warningColor).Bold(true)
	goodStyle  = lipgloss.NewStyle().Foreground(successColor).Bold(true)
	easyStyle  = lipgloss.NewStyle().Foreground(easyColor).Bold(true)

	progressFullStyle  = lipgloss.NewStyle().Foreground(successColor)
	progressEmptyStyle = lipgloss.NewStyle().Foreground(mutedColor)

	statsLabelStyle = lipgloss.NewStyle().
			Foreground(mutedColor).
			Width(12).
			Align(lipgloss.Right)

	statsBarAgain = lipgloss.NewStyle().Foreground(dangerColor)
	statsBarHard  = lipgloss.NewStyle().Foreground(warningColor)
	statsBarGood  = lipgloss.NewStyle().Foreground(successColor)
	statsBarEasy  = lipgloss.NewStyle().Foreground(easyColor)
)
