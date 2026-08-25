package transfer

import (
	"bytes"
	"context"
	"io"
	"mime/quotedprintable"
)

func DecodeQuotedPrintable(ctx context.Context, r io.Reader) ([]byte, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	dec := quotedprintable.NewReader(r)
	out, err := readTransferAll(ctx, dec)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func EncodeQuotedPrintable(w io.Writer, raw []byte) error {
	enc := quotedprintable.NewWriter(w)
	if _, err := enc.Write(raw); err != nil {
		return err
	}
	return enc.Close()
}

func SoftBreak(w io.Writer, line []byte, limit int) error {
	for len(line) > limit {
		if _, err := w.Write(line[:limit]); err != nil {
			return err
		}
		if _, err := w.Write([]byte("=\r\n")); err != nil {
			return err
		}
		line = line[limit:]
	}
	_, err := w.Write(line)
	return err
}

func IsQPLine(s []byte) bool {
	return bytes.Contains(s, []byte("="))
}
