package erax

import (
	"errors"
	"fmt"
	"sort"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

func Trace(err Error) string {
	if err == nil {
		return ""
	}

	return formatErrorChain(err, true)
}

func SetErrorColor(color lipgloss.Color) {
	errorColor = color
	errorText = lipgloss.NewStyle().Foreground(errorColor)
}

func SetKeyColor(color lipgloss.Color) {
	keyColor = color
	keyText = lipgloss.NewStyle().Foreground(keyColor)
}

func SetNormalColor(color lipgloss.Color) {
	normalColor = color
	branch1 = lipgloss.NewStyle().Foreground(normalColor).Render(" ├─ ")
	branch2 = lipgloss.NewStyle().Foreground(normalColor).Render(" │ ")
	branch3 = lipgloss.NewStyle().Foreground(normalColor).Render(" ╰─ ")
	message = lipgloss.NewStyle().Foreground(normalColor).Render(" ▼ [ERROR TRACE]\n")
}

func formatError(text string) string {
	lines := strings.Split(text, "\n")
	output := ""
	for lineIdx, line := range lines {
		if lineIdx != 0 {
			output += branch2 + " "
		}
		output += errorText.Render(line)
		if lineIdx < len(lines)-1 {
			output += "\n"
		}
	}
	return output
}

func formatText(text string) string {
	lines := strings.Split(text, "\n")
	output := ""
	for lineIdx, line := range lines {
		if lineIdx != 0 {
			output += branch2 + " "
		}
		output += line
		if lineIdx < len(lines)-1 {
			output += "\n"
		}
	}
	return output
}

func formatErrorChain(err Error, isFirst bool) string {
	var sb strings.Builder

	prefix := branch1
	if isFirst {
		prefix = message + "\n" + branch1
	}

	sb.WriteString(prefix + formatError(err.Msg()) + "\n")

	keys := make([]string, 0, len(err.Metas()))
	for key := range err.Metas() {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	for i, key := range keys {
		value := err.Metas()[key]
		connector := " " + branch1
		if i == len(keys)-1 {
			connector = " " + branch3
		}
		sb.WriteString(fmt.Sprintf("%s%s%s: %v\n", branch2, connector, keyText.Render(key), formatText(value.(string))))
	}

	if unwrapped := err.Unwrap(); unwrapped != nil {
		var next Error
		if errors.As(unwrapped, &next) {
			sb.WriteString(formatErrorChain(next, false))
		} else {
			sb.WriteString(branch3 + formatError(unwrapped.Error()))
		}
	}

	return sb.String()
}

var (
	errorColor  lipgloss.Color = "#f38ba8"
	keyColor    lipgloss.Color = "#cba6f7"
	normalColor lipgloss.Color = "#585b70"

	errorText = lipgloss.NewStyle().Foreground(errorColor)
	keyText   = lipgloss.NewStyle().Foreground(keyColor)

	branch1 = lipgloss.NewStyle().Foreground(normalColor).Render(" ├─ ")
	branch2 = lipgloss.NewStyle().Foreground(normalColor).Render(" │ ")
	branch3 = lipgloss.NewStyle().Foreground(normalColor).Render(" ╰─ ")
	message = lipgloss.NewStyle().Foreground(normalColor).Render(" ▼ [ERROR TRACE]")
)
