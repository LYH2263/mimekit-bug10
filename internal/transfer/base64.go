package transfer

import (
	"context"
	"encoding/base64"
	"io"
)

func DecodeBase64(ctx context.Context, r io.Reader) ([]byte, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	dec := base64.NewDecoder(base64.StdEncoding, r)
	return readTransferAll(ctx, dec)
}

func EncodeBase64(w io.Writer, raw []byte) error {
	enc := base64.NewEncoder(base64.StdEncoding, w)
	if _, err := enc.Write(raw); err != nil {
		return err
	}
	return enc.Close()
}

func readTransferAll(ctx context.Context, r io.Reader) ([]byte, error) {
	var buf []byte
	tmp := make([]byte, 4096)
	for {
		if err := ctx.Err(); err != nil {
			return nil, err
		}
		n, err := r.Read(tmp)
		if n > 0 {
			buf = append(buf, tmp[:n]...)
			// 中途再查一次，避免大块读完才发现取消。
			if err := ctx.Err(); err != nil {
				return nil, err
			}
		}
		if err == io.EOF {
			return buf, nil
		}
		if err != nil {
			return nil, err
		}
	}
}
