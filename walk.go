package mimekit

import "github.com/LYH2263/go-mimekit/internal/buffer"

type WalkFunc func(*Part) error

func Walk(m *Message, fn WalkFunc) error {
	if m == nil || m.Root == nil {
		return nil
	}
	return walkPart(m.Root, fn)
}

func walkPart(p *Part, fn WalkFunc) error {
	view := buffer.Clone(p.Body)
	part := &Part{
		Headers:   p.Headers.Clone(),
		Body:      view,
		RawBody:   append([]byte(nil), p.RawBody...),
		MediaType: p.MediaType,
		Boundary:  p.Boundary,
	}
	if err := fn(part); err != nil {
		return err
	}
	for _, c := range p.Children {
		if err := walkPart(c, fn); err != nil {
			return err
		}
	}
	return nil
}

func CollectParts(m *Message) []*Part {
	var out []*Part
	_ = Walk(m, func(p *Part) error {
		out = append(out, p)
		return nil
	})
	return out
}
