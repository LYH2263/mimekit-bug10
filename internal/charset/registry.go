package charset

import "strings"

var aliases = map[string]string{
	"utf8":      "utf-8",
	"latin1":    "iso-8859-1",
	"iso8859-1": "iso-8859-1",
	"gb2312":    "gbk",
	"cp936":     "gbk",
}

func Normalize(label string) string {
	l := strings.ToLower(strings.TrimSpace(label))
	if l == "" {
		return "utf-8"
	}
	if canon, ok := aliases[l]; ok {
		return canon
	}
	return l
}

func Supported(label string) bool {
	switch Normalize(label) {
	case "utf-8", "iso-8859-1", "gbk", "us-ascii":
		return true
	default:
		return false
	}
}
