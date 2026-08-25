package mimekit

import (
	"mime"
	"net/textproto"
	"strings"
)

type Header map[string][]string

func (h Header) Get(key string) string {
	if h == nil {
		return ""
	}
	k := textproto.CanonicalMIMEHeaderKey(key)
	for hk, vals := range h {
		if textproto.CanonicalMIMEHeaderKey(hk) == k && len(vals) > 0 {
			return vals[0]
		}
	}
	return ""
}

func (h Header) Set(key, val string) {
	if h == nil {
		return
	}
	k := textproto.CanonicalMIMEHeaderKey(key)
	h[k] = []string{val}
}

func (h Header) Add(key, val string) {
	k := textproto.CanonicalMIMEHeaderKey(key)
	h[k] = append(h[k], val)
}

func (h Header) Clone() Header {
	out := make(Header, len(h))
	for k, v := range h {
		cp := make([]string, len(v))
		copy(cp, v)
		out[k] = cp
	}
	return out
}

type Part struct {
	Headers   Header
	Body      []byte
	RawBody   []byte
	Children  []*Part
	MediaType string
	Boundary  string
}

type Message struct {
	Headers Header
	Root    *Part
}

func MediaTypeOf(h Header) string {
	ct := h.Get("Content-Type")
	if ct == "" {
		return "text/plain"
	}
	media, _, _ := mime.ParseMediaType(ct)
	if media == "" {
		return "text/plain"
	}
	return media
}

func TransferEncoding(h Header) string {
	te := strings.ToLower(strings.TrimSpace(h.Get("Content-Transfer-Encoding")))
	if te == "" {
		return "7bit"
	}
	return te
}

func CharsetLabel(h Header) string {
	ct := h.Get("Content-Type")
	if ct == "" {
		return "utf-8"
	}
	_, params, err := mime.ParseMediaType(ct)
	if err != nil {
		return "utf-8"
	}
	cs, ok := params["charset"]
	if !ok || cs == "" {
		return "utf-8"
	}
	return strings.ToLower(cs)
}
