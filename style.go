package erax

import "github.com/charmbracelet/lipgloss"

func SetBranchColor(color lipgloss.Color) {
	branchColor = color
	branchS = lipgloss.NewStyle().Foreground(branchColor).Render(" ╰╮")
	branchH = lipgloss.NewStyle().Foreground(branchColor).Render(" ├╮")
	branchTwix = lipgloss.NewStyle().Foreground(branchColor).Render(" ││ ")
	branchNextBig = lipgloss.NewStyle().Foreground(branchColor).Render(" ├── ")
	branchMid = lipgloss.NewStyle().Foreground(branchColor).Render(" │")
	branchEndBig = lipgloss.NewStyle().Foreground(branchColor).Render(" ╰── ")
	branchNext = lipgloss.NewStyle().Foreground(branchColor).Render("├─ ")
	branchEnd = lipgloss.NewStyle().Foreground(branchColor).Render("╰─ ")
	message = lipgloss.NewStyle().Foreground(branchColor).Render(" ▼ [ERROR TRACE]")
}

func SetErrorColor(color lipgloss.Color) {
	errorColor = color
	errorText = lipgloss.NewStyle().Foreground(errorColor)
}

func SetKeyColor(color lipgloss.Color) {
	keyColor = color
	keyText = lipgloss.NewStyle().Foreground(keyColor)
}

func SetValueColor(color lipgloss.Color) {
	valueColor = color
	valueText = lipgloss.NewStyle().Foreground(valueColor)
}

var (
	branchColor lipgloss.Color = "#585b70"
	errorColor  lipgloss.Color = "#f38ba8"
	keyColor    lipgloss.Color = "#cba6f7"
	valueColor  lipgloss.Color = "#a6e3a1"

	branchS       = lipgloss.NewStyle().Foreground(branchColor).Render(" ╰╮")
	branchH       = lipgloss.NewStyle().Foreground(branchColor).Render(" ├╮")
	branchTwix    = lipgloss.NewStyle().Foreground(branchColor).Render(" ││ ")
	branchNextBig = lipgloss.NewStyle().Foreground(branchColor).Render(" ├── ")
	branchMid     = lipgloss.NewStyle().Foreground(branchColor).Render(" │")
	branchEndBig  = lipgloss.NewStyle().Foreground(branchColor).Render(" ╰── ")
	branchNext    = lipgloss.NewStyle().Foreground(branchColor).Render("├─ ")
	branchEnd     = lipgloss.NewStyle().Foreground(branchColor).Render("╰─ ")
	message       = lipgloss.NewStyle().Foreground(branchColor).Render(" ▼ [ERROR TRACE]")

	errorText = lipgloss.NewStyle().Foreground(errorColor)
	keyText   = lipgloss.NewStyle().Foreground(keyColor)
	valueText = lipgloss.NewStyle().Foreground(valueColor)
)
