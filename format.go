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

func errorToMap(err *errorType) map[string]any {
	m := map[string]any{
		"message": err.msg,
	}

	if len(err.meta) > 0 {
		m["meta"] = err.meta
	}

	hasCause := err.cause != nil
	hasErrs := len(err.errs) > 0

	if !hasCause && !hasErrs {
		return m
	}

	if hasCause && !hasErrs {
		m["cause"] = FormatToJSONMap(err.cause)
	} else {
		totalLen := len(err.errs)
		if hasCause {
			totalLen++
		}

		causeSlice := make([]map[string]any, 0, totalLen)
		if hasCause {
			causeSlice = append(causeSlice, FormatToJSONMap(err.cause))
		}
		for _, ue := range err.errs {
			causeSlice = append(causeSlice, FormatToJSONMap(ue))
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

	if (meta == nil || !metaOk) && !causeOk {
		return errors.New(msg)
	}

	err := &errorType{
		msg: msg,
	}

	if meta != nil && metaOk {
		err.meta = meta
	}

	if causeOk {
		if value, ok := cause.(map[string]any); ok {
			if childErr := mapToError(value); childErr != nil {
				err.cause = childErr
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

	sb.WriteString(prefix)
	writeFormattedError(&sb, err.Error())
	sb.WriteByte('\n')

	sb.WriteString(formatMeta(GetMetas(err), false))

	hasCause := err.cause != nil
	hasErrs := len(err.errs) > 0

	if !hasCause && !hasErrs {
		return sb.String()
	}

	nestingLevel++

	if hasCause {
		next, isErax := asErax(err.cause)
		if isErax {
			sb.WriteString(formatDefault(next, nestingLevel))
		} else {
			if !hasErrs {
				prefix = branch3
			} else {
				prefix = branch1
			}
			sb.WriteString(prefix)
			writeFormattedError(&sb, err.cause.Error())
		}
	}

	for i, ue := range err.errs {
		if i > 0 || hasCause {
			sb.WriteString("\n")
		}

		next, isErax := asErax(ue)
		if isErax {
			sb.WriteString(formatDefault(next, nestingLevel))
		} else {
			if i == len(err.errs)-1 {
				prefix = branch3
			} else {
				prefix = branch1
			}
			sb.WriteString(prefix)
			writeFormattedError(&sb, ue.Error())
		}
	}

	return sb.String()
}

func formatValue(text string, isLastPair, isLast bool) string {
	if !strings.Contains(text, "\n") {
		return valueText.Render(text)
	}

	var sb strings.Builder
	sb.WriteByte('\n')
	lines := strings.Split(text, "\n")
	linesLen := len(lines)

	for i, line := range lines {
		var prefix string
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

		sb.WriteString(prefix)
		sb.WriteString(valueText.Render(line))

		if i < linesLen-1 {
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

		sb.WriteString(connector1)
		sb.WriteString(connector2)
		sb.WriteString(keyText.Render(field.Key))
		sb.WriteString(": ")
		sb.WriteString(formatValue(field.Value, isLastPair, isLast))
		sb.WriteString(ending)
	}

	return sb.String()
}

func formatAlienError(err *errorType, isLast bool) string {
	var sb strings.Builder

	var targets []error
	if err.cause != nil {
		targets = append(targets, err.cause)
	}
	targets = append(targets, err.errs...)
	targetsLen := len(targets)

	for i, cause := range targets {
		lines := strings.Split(fmt.Sprintf("%+v", cause), "\n")
		linesLen := len(lines)
		isCurrentLast := isLast && i == targetsLen-1

		for lineIdx, line := range lines {
			if lineIdx == 0 {
				if isCurrentLast {
					sb.WriteString(branch3)
				} else {
					sb.WriteString(branch1)
				}
			} else {
				if isCurrentLast {
					sb.WriteString("    ")
				} else {
					sb.WriteString(branch2 + " ")
				}
			}
			sb.WriteString(errorText.Render(line))
			if lineIdx < linesLen-1 {
				sb.WriteString("\n")
			}
		}
		if i < targetsLen-1 {
			sb.WriteString("\n")
		}
	}

	metaStr := formatMeta(GetMetas(err), isLast)
	if metaStr != "" {
		sb.WriteString("\n" + metaStr)
	}

	return sb.String()
}

func writeFormattedError(sb *strings.Builder, text string) {
	if !strings.Contains(text, "\n") {
		sb.WriteString(errorText.Render(text))
		return
	}

	lines := strings.Split(text, "\n")
	linesLen := len(lines)
	for lineIdx, line := range lines {
		if lineIdx != 0 {
			sb.WriteString(branch2)
			sb.WriteByte(' ')
		}
		sb.WriteString(errorText.Render(line))
		if lineIdx < linesLen-1 {
			sb.WriteByte('\n')
		}
	}
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
