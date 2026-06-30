package erax

import "strings"

// writeMeta formats and writes metadata fields to a string builder with proper indentation.
func writeMeta(sb *strings.Builder, meta []MetaField, isNested bool, levels []bool) {
	metaLen := len(meta)
	if metaLen == 0 {
		return
	}

	var childLevels []bool
	if len(levels) > 1 {
		childLevels = levels[:len(levels)-1]
	}

	isLastLevel := len(levels) > 0 && levels[len(levels)-1]

	for i := 0; i < metaLen; i++ {
		field := &meta[i]
		isLastPair := i == metaLen-1
		writeIndent(sb, childLevels)

		if isLastLevel {
			sb.WriteByte(' ')
			sb.WriteString(branchMid)
			sb.WriteByte(' ')
		} else if isNested {
			sb.WriteString(branchTwix)
		} else {
			sb.WriteString(branchMid)
			sb.WriteString("  ")
		}
		sb.WriteString("  ")
		if isLastPair {
			sb.WriteString(branchEnd)
		} else {
			sb.WriteString(branchNext)
		}

		sb.WriteString(keyText.Render(field.Key))
		sb.WriteString(": ")
		writeValue(sb, field.Value, isLastPair, isNested, levels)
		sb.WriteByte('\n')
	}
}

// writeValue formats and writes a metadata value, handling multi-line values with proper indentation.
func writeValue(sb *strings.Builder, text string, isLastPair, isNested bool, levels []bool) {
	if indexByte(text, '\n') == -1 {
		sb.WriteString(valueText.Render(text))
		return
	}

	sb.WriteByte('\n')
	start := 0
	lineIdx := 0
	textLen := len(text)

	var childLevels []bool
	if len(levels) > 1 {
		childLevels = levels[:len(levels)-1]
	}

	isLastLevel := len(levels) > 0 && levels[len(levels)-1]

	for start < textLen {
		idx := indexByte(text[start:], '\n')
		var line string
		if idx == -1 {
			line = text[start:]
			start = textLen
		} else {
			line = text[start : start+idx]
			start += idx + 1
		}

		if lineIdx > 0 {
			sb.WriteByte('\n')
		}

		writeIndent(sb, childLevels)

		if isLastLevel {
			sb.WriteByte(' ')
			sb.WriteString(branchMid)
			sb.WriteByte(' ')
		} else if isNested {
			sb.WriteString(branchTwix)
		} else {
			sb.WriteString(branchMid)
			sb.WriteString("  ")
		}
		sb.WriteByte(' ')

		if isLastPair {
			sb.WriteString("   ")
		} else {
			sb.WriteString(branchMid)
			sb.WriteByte(' ')
		}
		sb.WriteString("  ")

		sb.WriteString(valueText.Render(line))
		lineIdx++
	}
}
