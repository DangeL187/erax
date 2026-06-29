package erax

import (
	"fmt"
	"strings"
)

// Format pretty-prints the error trace.
//
// You're a good person, so please don't feed it non-Erax errors, okay?
//
// It won't break anything (it's just %+v), but... why would you even do that?
func Format(err error) string {
	return fmt.Sprintf("%+v", err)
}

func formatErrorChain(sb *strings.Builder, err *errorType, isParentNested bool, levels []bool) {
	hasCause := err.cause != nil
	hasErrs := len(err.errs) > 0
	isNested := !hasCause && hasErrs

	if levels == nil {
		if isNested {
			sb.WriteString(branchEndBig)
		} else {
			sb.WriteString(branchNextBig)
		}
		levels = append(levels, isNested)
	}

	writeFormattedError(sb, err.msg, isParentNested, hasCause, false, levels)
	sb.WriteByte('\n')

	writeMeta(sb, err.meta, false, isParentNested, levels)

	if !hasCause && !hasErrs {
		return
	}

	for i, ue := range err.errs {
		if i > 0 {
			sb.WriteByte('\n')
		}

		isLast := i == len(err.errs)-1

		next, isErax := asErax(ue)
		if isErax {
			writeIndent(sb, levels)
			if isNested {
				if isLast {
					sb.WriteString(branchS)
					sb.WriteByte('\n')
					writeIndent(sb, levels)
					sb.WriteString("  ")
				} else {
					sb.WriteString(branchH)
					sb.WriteByte('\n')
					writeIndent(sb, levels)
					sb.WriteString(branchMid)
				}
				if next.cause == nil {
					sb.WriteString(branchEnd)
				} else {
					sb.WriteString(branchNext)
				}
			}

			formatErrorChain(sb, next, isNested, append(levels, isLast))
		} else {
			writeIndent(sb, levels)

			if isLast {
				sb.WriteString(branchEndBig)
			} else {
				sb.WriteString(branchNextBig)
			}

			writeFormattedError(sb, fmt.Sprintf("%+v", ue), isNested, false, false, append(levels, isLast))
		}
	}

	if hasCause {
		var childLevels []bool
		if len(levels) > 1 {
			childLevels = levels[:len(levels)-1]
		}

		writeIndent(sb, childLevels)

		next, isErax := asErax(err.cause)
		if isErax {
			if len(levels) > 0 && levels[len(levels)-1] {
				sb.WriteString("  ")
				sb.WriteString(branchNext)
				childLevels = append(childLevels, true)
			} else if isParentNested {
				sb.WriteString(branchMid)
				if next.cause == nil {
					sb.WriteString(branchEnd)
				} else {
					sb.WriteString(branchNext)
				}
				childLevels = append(childLevels, false)
			}

			formatErrorChain(sb, next, isParentNested, childLevels)
		} else {
			if len(levels) > 0 && levels[len(levels)-1] {
				sb.WriteString("  ")
				sb.WriteString(branchEnd)
				childLevels = levels
			} else if isParentNested {
				sb.WriteString(branchMid)
				sb.WriteString(branchEnd)
				childLevels = levels
			} else {
				sb.WriteString(branchEndBig)
				childLevels = append(childLevels, true)
			}
			writeFormattedError(sb, fmt.Sprintf("%+v", err.cause), isParentNested, false, false, childLevels)
		}
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
