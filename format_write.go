package erax

import "strings"

// writeFormattedError formats and writes an error message, handling multi-line messages with proper indentation
func writeFormattedError(sb *strings.Builder, text string, isParentNested, hasCause, isAlien bool, levels []bool) {
	if indexByte(text, '\n') == -1 {
		if isAlien {
			sb.WriteString(alienText.Render(text))
		} else {
			sb.WriteString(errorText.Render(text))
		}
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

		if isAlien {
			sb.WriteString(alienText.Render(line))
		} else {
			sb.WriteString(errorText.Render(line))
		}
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
			} else if hasCause != isParentNested {
				sb.WriteString(branchMid)
				sb.WriteString("  ")
			}
			sb.WriteByte(' ')
		}
		lineIdx++
	}
}

func writeIndent(sb *strings.Builder, levels []bool) {
	for _, isLast := range levels {
		if isLast {
			sb.WriteString("     ")
		} else {
			sb.WriteString(branchMid)
			sb.WriteString("   ")
		}
	}
}

// indexByte returns the index of the first occurrence of a byte in a string.
func indexByte(s string, c byte) int {
	for i := 0; i < len(s); i++ {
		if s[i] == c {
			return i
		}
	}
	return -1
}
