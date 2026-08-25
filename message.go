package mimekit

import (
	"io"
	"strings"
)

func NewTextPlain(body string) *Message {
	p := &Part{
		Headers: Header{
			"Content-Type": {"text/plain; charset=utf-8"},
		},
		Body:      []byte(body),
		MediaType: "text/plain",
	}
	return &Message{Headers: p.Headers.Clone(), Root: p}
}

func (m *Message) PlainText() string {
	if m == nil || m.Root == nil {
		return ""
	}
	return strings.TrimSpace(string(m.Root.Body))
}

func (m *Message) PartCount() int {
	n := 0
	_ = Walk(m, func(p *Part) error {
		n++
		return nil
	})
	return n
}

type EncodeOptions struct {
	BoundaryPrefix string
	FoldAt         int
}

func (m *Message) WriteTo(w io.Writer, opts EncodeOptions) error {
	return Encode(w, m, opts)
}
