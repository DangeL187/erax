package erax

import (
	"errors"
	"fmt"
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

func Format(err error) string {
	return fmt.Sprintf("%f", err)
}

func FormatV(err error) string {
	return fmt.Sprintf("%+v", err)
}

type unwrappableError interface {
	error

	Unwrap() error
}

func errorToMap(err *errorType) map[string]any {
	m := map[string]any{
		"message": err.msg,
	}

	if len(err.meta) > 0 {
		m["meta"] = err.meta
	}

	if len(err.errs) == 0 {
		return m
	}

	if len(err.errs) == 1 {
		m["cause"] = FormatToJSONMap(err.errs[0])
	} else {
		causeSlice := make([]map[string]any, len(err.errs))
		for i, ue := range err.errs {
			causeSlice[i] = FormatToJSONMap(ue)
		}
		m["cause"] = causeSlice
	}

	return m
}

func mapToError(m map[string]any) error {
	msg, msgOk := m["message"].(string)
	if !msgOk {
		return nil
	}

	meta, metaOk := m["meta"].([]MetaField)
	cause, causeOk := m["cause"]

	if !metaOk && !causeOk {
		return errors.New(msg)
	}

	err := &errorType{
		msg:  msg,
		meta: meta,
	}

	if causeOk {
		if value, ok := cause.(map[string]any); ok {
			if childErr := mapToError(value); childErr != nil {
				err.errs = []error{childErr}
			}
		} else if value, ok := cause.([]map[string]any); ok {
			err.errs = make([]error, len(value))
			for i, c := range value {
				err.errs[i] = mapToError(c)
			}
		}
	}

	return err
}

func formatErrorChain(err *errorType, isFirst bool, nestingLevel int) string {
	var sb strings.Builder

	prefix := branch1
	if isFirst {
		prefix = message + "\n" + branch1
	}

	sb.WriteString(prefix + formatError(err.Error()) + "\n")
	sb.WriteString(formatMeta(GetMetas(err), false))

	if len(err.errs) == 0 {
		return sb.String()
	}

	prefix = branch1

	nestingLevel++

	for i, ue := range err.errs {
		if i > 0 {
			sb.WriteString("\n")
		}

		next, isErax := asErax(ue)
		if isErax {
			sb.WriteString(formatDefault(next, nestingLevel))
		} else {
			if i == len(err.errs)-1 {
				prefix = branch3
			}
			sb.WriteString(prefix + formatError(ue.Error()))
		}
	}

	return sb.String()
}

func formatValue(text string, isLastPair, isLast bool) string {
	lines := strings.Split(text, "\n")
	var sb strings.Builder

	if len(lines) > 1 {
		sb.WriteString("\n")
	}

	for i, line := range lines {
		var prefix string
		if len(lines) > 1 {
			prefix = "   "
			if !isLast {
				prefix = branch2
			}
			prefix += " "

			if isLastPair {
				prefix += "   "
			} else {
				prefix += branch2
			}
			prefix += "  "
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

func formatMeta(meta []MetaField, isLast bool) string {
	metaLen := len(meta)
	if metaLen == 0 {
		return ""
	}

	// too heavy to use :(
	/*sort.Slice(meta, func(i, j int) bool {
		return meta[i].Key < meta[j].Key
	})*/

	sb := strings.Builder{}

	for i, field := range meta {
		isLastPair := i == metaLen-1

		connector1 := branch2
		if isLast {
			connector1 = "   "
		}

		connector2 := " " + branch1
		if isLastPair {
			connector2 = " " + branch3
		}

		ending := "\n"
		if isLast && isLastPair {
			ending = ""
		}

		sb.WriteString(fmt.Sprintf("%s%s%s: %v%s",
			connector1,
			connector2,
			keyText.Render(field.Key),
			formatValue(field.Value, isLastPair, isLast),
			ending,
		))
	}

	return sb.String()
}

func formatAlienError(err *errorType, isLast bool) string {
	sb := strings.Builder{}

	for i, cause := range err.errs {
		lines := strings.Split(fmt.Sprintf("%+v", cause), "\n")
		for lineIdx, line := range lines {
			if lineIdx == 0 {
				if isLast && i == len(err.errs)-1 {
					sb.WriteString(branch3)
				} else {
					sb.WriteString(branch1)
				}
			} else {
				if isLast && i == len(err.errs)-1 {
					sb.WriteString("    ")
				} else {
					sb.WriteString(branch2 + " ")
				}
			}
			sb.WriteString(errorText.Render(line))
			if lineIdx < len(lines)-1 {
				sb.WriteString("\n")
			}
		}
		if i < len(err.errs)-1 {
			sb.WriteString("\n")
		}
	}

	metaStr := formatMeta(GetMetas(err), isLast)
	if metaStr != "" {
		sb.WriteString("\n" + metaStr)
	}

	return sb.String()
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
