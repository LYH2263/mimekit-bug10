package mimekit

import (
	"context"
	"sync"
)

type Pipeline struct {
	mu     sync.Mutex
	closed bool
	opts   Options
}

func NewPipeline(opts Options) *Pipeline {
	return &Pipeline{opts: opts.withDefaults()}
}

func (p *Pipeline) Close() error {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.closed = true
	return nil
}

func (p *Pipeline) checkOpen() error {
	p.mu.Lock()
	defer p.mu.Unlock()
	if p.closed {
		return ErrClosed
	}
	return nil
}

type parseState struct {
	headers Header
	rawBody []byte
	root    *Part
}

func (p *Pipeline) runStages(ctx context.Context, raw []byte) (*Message, error) {
	if err := p.checkOpen(); err != nil {
		return nil, err
	}
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	st := &parseState{}
	hdrBlock, body, err := splitHeaderBody(raw)
	if err != nil {
		return nil, err
	}
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	st.headers, err = parseHeaderBlock(hdrBlock)
	if err != nil {
		return nil, err
	}
	st.rawBody = body
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	te := TransferEncoding(st.headers)
	decoded, err := decodeTransfer(ctx, body, te)
	if err != nil {
		return nil, err
	}
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	cs := CharsetLabel(st.headers)
	text, err := convertCharset(ctx, decoded, cs, "utf-8")
	if err != nil {
		text = decoded
	}
	root := &Part{
		Headers:   st.headers.Clone(),
		RawBody:   append([]byte(nil), body...),
		Body:      text,
		MediaType: MediaTypeOf(st.headers),
	}
	if err != nil {
		return &Message{Headers: st.headers.Clone(), Root: root}, err
	}
	if isMultipart(root.MediaType) {
		boundary, err := extractBoundary(st.headers.Get("Content-Type"))
		if err != nil {
			return nil, err
		}
		root.Boundary = boundary
		if err := ctx.Err(); err != nil {
			return nil, err
		}
		// multipart 段体必须是 transfer 解码后的字节流，不能回退到原始编码层。
		children, err := readMultipart(ctx, decoded, boundary, p.opts)
		if err != nil {
			return nil, err
		}
		root.Children = children
	}
	if err := p.checkOpen(); err != nil {
		return nil, err
	}
	st.root = root
	return &Message{Headers: st.headers.Clone(), Root: st.root}, nil
}

func isMultipart(media string) bool {
	return len(media) >= 10 && media[:10] == "multipart/"
}
