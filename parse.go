package mimekit

import (
	"bytes"
	"context"
	"io"
)

func Parse(ctx context.Context, r io.Reader) (*Message, error) {
	return ParseWithOptions(ctx, r, Options{})
}

func ParseWithOptions(ctx context.Context, r io.Reader, opts Options) (*Message, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	opts = opts.withDefaults()
	raw, err := readAllLimited(r, opts.MaxHeaderBytes+opts.MaxBodyBytes)
	if err != nil {
		return nil, err
	}
	p := NewPipeline(opts)
	return p.runStages(ctx, raw)
}

func ParseBytes(ctx context.Context, raw []byte) (*Message, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	p := NewPipeline(Options{})
	return p.runStages(ctx, raw)
}

func splitHeaderBody(raw []byte) ([]byte, []byte, error) {
	idx := bytes.Index(raw, []byte("\r\n\r\n"))
	if idx < 0 {
		idx = bytes.Index(raw, []byte("\n\n"))
		if idx < 0 {
			return nil, nil, ErrInvalidMIME
		}
		return raw[:idx], raw[idx+2:], nil
	}
	return raw[:idx], raw[idx+4:], nil
}

func parseHeaderBlock(block []byte) (Header, error) {
	return headerParseBlock(block)
}

func readAllLimited(r io.Reader, limit int) ([]byte, error) {
	var buf bytes.Buffer
	_, err := io.CopyN(&buf, r, int64(limit)+1)
	if err != nil && err != io.EOF {
		return nil, err
	}
	if buf.Len() > limit {
		return nil, ErrInvalidMIME
	}
	return buf.Bytes(), nil
}
