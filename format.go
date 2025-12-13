package erax

import (
	"encoding/json"
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

func Format(err error) string {
	return fmt.Sprintf("%f", err)
}

func FormatV(err error) string {
	return fmt.Sprintf("%+v", err)
}

func FormatToJSONString(err error) (string, error) {
	if err == nil {
		return "", nil
	}

	data := FormatToJSONMap(err)

	jsonBytes, e := json.Marshal(data)
	if e != nil {
		return "", e
	}

	return string(jsonBytes), nil
}

type unwrappableError interface {
	error

	Unwrap() error
}

func FormatToJSONMap(err error) map[string]any {
	if err == nil {
		return nil
	}

	var next unwrappableError
	if errors.As(err, &next) {
		return errorToMap(next)
	}

	return map[string]any{
		"message": err.Error(),
	}
}

func FromJSONMap(m map[string]any) error {
	if m == nil || len(m) == 0 {
		return nil
	}

	return mapToError(m)
}

func errorToMap(err unwrappableError) map[string]any {
	m := map[string]any{}

	m["message"] = err.Error()

	meta := GetMetas(err)
	if len(meta) != 0 {
		m["meta"] = meta
	}

	unwrapped := err.Unwrap()
	if unwrapped == nil {
		return m
	}

	var list []error
	if uw, ok := unwrapped.(interface{ Unwrap() []error }); ok {
		list = uw.Unwrap()
	} else {
		list = []error{unwrapped}
	}

	var cause []map[string]any
	for _, ue := range list {
		cause = append(cause, FormatToJSONMap(ue))
	}

	if len(cause) == 1 {
		m["cause"] = cause[0]
	} else if len(cause) > 1 {
		m["cause"] = cause
	}

	return m
}

func mapToError(m map[string]any) error {
	msg, msgOk := m["message"].(string)
	if !msgOk {
		return nil
	}

	meta, metaOk := m["meta"].(map[string]string)
	cause, causeOk := m["cause"]

	if !metaOk && !causeOk {
		return errors.New(msg)
	}

	err := &errorType{
		msg: msg,
	}

	if metaOk {
		err.meta = meta
	}

	if causeOk {
		if value, ok := cause.(map[string]any); ok {
			err.err = mapToError(value)
		} else if value, ok := cause.([]map[string]any); ok {
			var errs []error
			for _, c := range value {
				errs = append(errs, mapToError(c))
			}
			err.err = errors.Join(errs...)
		}
	}

	return err
}

func formatErrorChain(err unwrappableError, isFirst bool) string {
	var sb strings.Builder

	prefix := branch1
	if isFirst {
		prefix = message + "\n" + branch1
	}

	sb.WriteString(prefix + formatError(err.Error()) + "\n")

	sb.WriteString(formatMeta(GetMetas(err), false))

	unwrapped := err.Unwrap()
	if unwrapped == nil {
		return sb.String()
	}

	var list []error
	if uw, ok := unwrapped.(interface{ Unwrap() []error }); ok {
		list = uw.Unwrap()
	} else {
		list = []error{unwrapped}
	}

	prefix = branch1

	for i, ue := range list {
		if i > 0 {
			sb.WriteString("\n")
		}

		var next *errorType
		if errors.As(ue, &next) {
			sb.WriteString(fmt.Sprintf("%+v", next))
		} else {
			if i == len(list)-1 {
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

func formatMeta(meta map[string]string, isLast bool) string {
	sb := strings.Builder{}

	keys := make([]string, 0, len(meta))
	for key := range meta {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	for i, key := range keys {
		value := meta[key]
		isLastPair := i == len(keys)-1

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

		sb.WriteString(fmt.Sprintf("%s%s%s: %v%s", connector1, connector2, keyText.Render(key), formatValue(value, isLastPair, isLast), ending))
	}

	return sb.String()
}

func formatAlienError(err unwrappableError, isLast bool) string {
	sb := strings.Builder{}

	lines := strings.Split(fmt.Sprintf("%+v", err.Unwrap()), "\n")
	for lineIdx, line := range lines {
		if lineIdx == 0 {
			if isLast {
				sb.WriteString(branch3)
			} else {
				sb.WriteString(branch1)
			}
		} else {
			if isLast {
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

	sb.WriteString("\n" + formatMeta(GetMetas(err), isLast))

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
