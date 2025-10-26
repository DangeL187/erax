package erax

import (
	"errors"
	"fmt"
	"sort"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

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

func SetValueColor(color lipgloss.Color) {
	valueColor = color
	valueText = lipgloss.NewStyle().Foreground(valueColor)
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

func formatValue(text string, isLast bool) string {
	lines := strings.Split(text, "\n")
	var sb strings.Builder

	if len(lines) > 1 {
		sb.WriteString("\n")
	}

	for i, line := range lines {
		var prefix string
		if len(lines) > 1 {
			prefix = branch2 + " "
			if !isLast {
				prefix += branch2
			} else {
				prefix += "   "
			}
			prefix += "  "
		} else if i != 0 {
			prefix = branch2 + "      "
		} else {
			prefix = ""
		}

		sb.WriteString(prefix)
		sb.WriteString(valueText.Render(line))

		if i < len(lines)-1 {
			sb.WriteString("\n")
		}
	}

	return sb.String()
}

type unwrappableError interface {
	error

	Unwrap() error
}

func formatErrorChain(err unwrappableError, isFirst bool) string {
	var sb strings.Builder

	prefix := branch1
	if isFirst {
		prefix = message + "\n" + branch1
	}

	sb.WriteString(prefix + formatError(err.Error()) + "\n")

	allMeta := GetMetas(err)
	keys := make([]string, 0, len(allMeta))
	for key := range allMeta {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	for i, key := range keys {
		value, _ := GetMeta(err, key)
		connector := " " + branch1
		isLast := i == len(keys)-1
		if isLast {
			connector = " " + branch3
		}
		sb.WriteString(fmt.Sprintf("%s%s%s: %v\n", branch2, connector, keyText.Render(key), formatValue(value, isLast)))
	}

	if unwrapped := err.Unwrap(); unwrapped != nil {
		var next unwrappableError
		if errors.As(unwrapped, &next) {
			sb.WriteString(fmt.Sprintf("%+v", next))
		} else {
			sb.WriteString(branch3 + formatError(unwrapped.Error()))
		}
	}

	return sb.String()
}

var (
	errorColor  lipgloss.Color = "#f38ba8"
	keyColor    lipgloss.Color = "#cba6f7"
	valueColor  lipgloss.Color = "#a6e3a1"
	normalColor lipgloss.Color = "#585b70"

	errorText = lipgloss.NewStyle().Foreground(errorColor)
	keyText   = lipgloss.NewStyle().Foreground(keyColor)
	valueText = lipgloss.NewStyle().Foreground(valueColor)

	branch1 = lipgloss.NewStyle().Foreground(normalColor).Render(" ├─ ")
	branch2 = lipgloss.NewStyle().Foreground(normalColor).Render(" │ ")
	branch3 = lipgloss.NewStyle().Foreground(normalColor).Render(" ╰─ ")
	message = lipgloss.NewStyle().Foreground(normalColor).Render(" ▼ [ERROR TRACE]")
)
