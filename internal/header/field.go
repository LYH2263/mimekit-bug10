package header

import (
	"net/textproto"
	"strings"
)

func CanonicalKey(k string) string {
	return textproto.CanonicalMIMEHeaderKey(k)
}

func Merge(dst, src map[string][]string) {
	for k, v := range src {
		key := CanonicalKey(k)
		cp := make([]string, len(v))
		copy(cp, v)
		dst[key] = append(dst[key], cp...)
	}
}

func First(h map[string][]string, key string) string {
	k := CanonicalKey(key)
	for hk, vals := range h {
		if CanonicalKey(hk) == k && len(vals) > 0 {
			return vals[0]
		}
	}
	return ""
}

func NormalizeValue(v string) string {
	return strings.TrimSpace(v)
}
