package erax

import "strings"

func writeMeta(sb *strings.Builder, meta []MetaField, isLast, isNested bool, levels []bool) {
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
		writeValue(sb, field.Value, isLastPair, isLast, isNested, levels)

		if !isLast || !isLastPair {
			sb.WriteByte('\n')
		}
	}
}

func writeValue(sb *strings.Builder, text string, isLastPair, isLast, isNested bool, levels []bool) {
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

		if isLast {
			sb.WriteString("(LOL)") // TODO: do we need it??
		} else {
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

func writeAlienErrorLines(sb *strings.Builder, err error, isLast bool) {
	text := FormatV(err)
	start := 0
	textLen := len(text)
	lineIdx := 0

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

		if lineIdx == 0 {
			sb.WriteByte(' ')
			if isLast {
				sb.WriteString(branchEnd)
			} else {
				sb.WriteString(branchNext)
			}
		} else {
			if isLast {
				sb.WriteString("    ")
			} else {
				sb.WriteString(branchMid)
				sb.WriteString("  ")
			}
		}
		sb.WriteString(errorText.Render(line))
		lineIdx++
	}
}

func writeFormattedError(sb *strings.Builder, text string, isParentNested, hasCause bool, levels []bool) {
	if indexByte(text, '\n') == -1 {
		sb.WriteString(errorText.Render(text))
		return
	}

	start := 0
	textLen := len(text)
	lineIdx := 0

	var childLevels []bool
	if len(levels) > 1 {
		childLevels = levels[:len(levels)-1]
	}

	isLast := len(levels) > 0 && levels[len(levels)-1]

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

		sb.WriteString(errorText.Render(line))
		if start < textLen {
			sb.WriteByte('\n')

			writeIndent(sb, childLevels)
			if hasCause && isParentNested {
				if isLast {
					sb.WriteByte(' ')
					sb.WriteString(branchMid)
					sb.WriteByte(' ')
				} else {
					sb.WriteString(branchTwix)
				}
			} else if !hasCause && isLast {
				sb.WriteString("    ")
			} else if hasCause != isParentNested { // hasCause != isParentNested
				sb.WriteString(branchMid)
				sb.WriteString("  ")
			}
			sb.WriteByte(' ')
		}
		lineIdx++
	}
}

func indexByte(s string, c byte) int {
	for i := 0; i < len(s); i++ {
		if s[i] == c {
			return i
		}
	}
	return -1
}
