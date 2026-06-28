package erax

import (
	"fmt"
	"strings"
)

// Format pretty-prints error-trace.
//
// You're a good person, don't put non-erax errors here, ok?
func Format(err error) string {
	return fmt.Sprintf("%f", err)
}

func FormatV(err error) string {
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

	writeFormattedError(sb, err.msg, isParentNested, hasCause, levels)
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

			formatDefault(sb, next, isNested, append(levels, isLast))
		} else {
			writeIndent(sb, levels)

			if isLast {
				sb.WriteString(branchEndBig)
			} else {
				sb.WriteString(branchNextBig)
			}

			writeFormattedError(sb, ue.Error(), isNested, false, append(levels, isLast))
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

			formatDefault(sb, next, isParentNested, childLevels)
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
			writeFormattedError(sb, err.cause.Error(), isParentNested, false, childLevels)
		}
	}
}

func formatAlienError(sb *strings.Builder, err *errorType, isLast bool, levels []bool) {
	totalTargets := len(err.errs)
	if err.cause != nil {
		totalTargets++
	}
	targets := make([]error, 0, totalTargets)
	if err.cause != nil {
		targets = append(targets, err.cause)
	}
	targets = append(targets, err.errs...)

	for i, cause := range targets {
		if i > 0 {
			sb.WriteByte('\n')
		}
		isCurrentLast := isLast && i == totalTargets-1
		writeAlienErrorLines(sb, cause, isCurrentLast)
	}

	writeMeta(sb, err.meta, isLast, false, levels)
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
