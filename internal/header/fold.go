package header

import (
	"strings"
	"unicode/utf8"
)

const defaultFoldAt = 76

func UnfoldLines(block string) string {
	lines := strings.Split(block, "\n")
	var out []string
	for _, line := range lines {
		line = strings.TrimRight(line, "\r")
		if len(out) > 0 && (line != "" && (line[0] == ' ' || line[0] == '\t')) {
			out[len(out)-1] += strings.TrimLeft(line, " \t")
			continue
		}
		out = append(out, line)
	}
	return strings.Join(out, "\n")
}

func FoldLine(line string, at int) string {
	if at <= 0 {
		at = defaultFoldAt
	}
	if len(line) <= at {
		return line
	}
	var b strings.Builder
	b.WriteString(line[:at])
	rest := line[at:]
	for len(rest) > 0 {
		b.WriteString("\r\n ")
		chunk := rest
		if len(chunk) > at-1 {
			chunk = rest[:at-1]
			rest = rest[at-1:]
		} else {
			rest = ""
		}
		b.WriteString(chunk)
	}
	return b.String()
}

func ValidHeaderLine(line string) bool {
	if line == "" {
		return false
	}
	for i := 0; i < len(line); i++ {
		if line[i] == ':' {
			return i > 0
		}
	}
	return false
}

func RuneLen(s string) int {
	return utf8.RuneCountInString(s)
}
