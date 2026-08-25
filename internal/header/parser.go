package header

import (
	"bytes"
	"errors"
	"strings"
)

var ErrBadHeader = errors.New("header: invalid line")

func ParseBlock(block []byte) (map[string][]string, error) {
	text := UnfoldLines(string(block))
	lines := strings.Split(text, "\n")
	h := make(map[string][]string)
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		if !ValidHeaderLine(line) {
			return nil, ErrBadHeader
		}
		colon := strings.IndexByte(line, ':')
		key := strings.TrimSpace(line[:colon])
		val := strings.TrimSpace(line[colon+1:])
		valBytes := []byte(val)
		h[key] = append(h[key], string(valBytes))
	}
	return h, nil
}

func Serialize(h map[string][]string) []byte {
	var buf bytes.Buffer
	for k, vals := range h {
		for _, v := range vals {
			buf.WriteString(k)
			buf.WriteString(": ")
			buf.WriteString(v)
			buf.WriteString("\r\n")
		}
	}
	return buf.Bytes()
}
