package mimekit

import (
	"context"
	"fmt"
	"strings"

	"github.com/LYH2263/go-mimekit/internal/charset"
	"github.com/LYH2263/go-mimekit/internal/transfer"
)

var sharedPipeline = NewPipeline(Options{})

func DecodeBody(ctx context.Context, p *Part) ([]byte, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	if err := sharedPipeline.checkOpen(); err != nil {
		return nil, err
	}
	if p == nil {
		return nil, ErrInvalidMIME
	}
	raw := p.RawBody
	if len(raw) == 0 {
		raw = p.Body
	}
	te := TransferEncoding(p.Headers)
	decoded, err := transfer.Decode(ctx, raw, te)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrTransfer, err)
	}
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	if err := sharedPipeline.checkOpen(); err != nil {
		return nil, err
	}
	label := CharsetLabel(p.Headers)
	out, err := charset.Decode(ctx, decoded, label)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrCharset, err)
	}
	return out, nil
}

// MaterializeBody 仅在 transfer+charset 全链成功后写回 Part.Body；失败路径不得污染部件。
func MaterializeBody(ctx context.Context, p *Part) error {
	if p == nil {
		return ErrInvalidMIME
	}
	out, err := DecodeBody(ctx, p)
	if err != nil {
		return err
	}
	p.Body = append([]byte(nil), out...)
	return nil
}

func DecodeBodyUTF8(ctx context.Context, p *Part) (string, error) {
	b, err := DecodeBody(ctx, p)
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(b)), nil
}

func ResetSharedPipeline(opts Options) {
	sharedPipeline = NewPipeline(opts)
}
