package mimekit

import (
	"context"
	"fmt"

	"github.com/LYH2263/go-mimekit/internal/charset"
	"github.com/LYH2263/go-mimekit/internal/header"
	"github.com/LYH2263/go-mimekit/internal/multipart"
	"github.com/LYH2263/go-mimekit/internal/transfer"
)

func headerParseBlock(block []byte) (Header, error) {
	m, err := header.ParseBlock(block)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrInvalidMIME, err)
	}
	out := make(Header, len(m))
	for k, v := range m {
		cp := make([]string, len(v))
		copy(cp, v)
		out[k] = cp
	}
	return out, nil
}

func decodeTransfer(ctx context.Context, raw []byte, te string) ([]byte, error) {
	out, err := transfer.Decode(ctx, raw, te)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrTransfer, err)
	}
	return out, nil
}

func convertCharset(ctx context.Context, raw []byte, from, to string) ([]byte, error) {
	_ = to
	return charset.Decode(ctx, raw, from)
}

func readMultipart(ctx context.Context, raw []byte, boundary string, opts Options) ([]*Part, error) {
	_ = opts
	sections, err := multipart.ReadAllParts(ctx, raw, boundary)
	if err != nil {
		return nil, err
	}
	parts := make([]*Part, 0, len(sections))
	for _, sec := range sections {
		if err := ctx.Err(); err != nil {
			return nil, err
		}
		p := &Part{
			Headers:   Header(sec.Headers),
			RawBody:   append([]byte(nil), sec.RawBody...),
			Body:      append([]byte(nil), sec.Body...),
			MediaType: sec.MediaType,
		}
		te := TransferEncoding(p.Headers)
		dec, err := transfer.Decode(ctx, p.RawBody, te)
		if err != nil {
			// 失败不得带着半成品 parts 返回，避免上层误挂 Children。
			return nil, fmt.Errorf("%w: part transfer: %v", ErrTransfer, err)
		}
		cs := CharsetLabel(p.Headers)
		text, err := charset.Decode(ctx, dec, cs)
		if err != nil {
			return nil, fmt.Errorf("%w: part charset: %v", ErrCharset, err)
		}
		p.Body = text
		parts = append(parts, p)
	}
	return parts, nil
}
