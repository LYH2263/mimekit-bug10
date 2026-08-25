package transfer

import (
	"bytes"
	"context"
	"strings"
)

func Decode(ctx context.Context, raw []byte, encoding string) ([]byte, error) {
	enc := strings.ToLower(strings.TrimSpace(encoding))
	r := bytes.NewReader(raw)
	switch enc {
	case "base64", "b64":
		return DecodeBase64(ctx, r)
	case "quoted-printable", "qp":
		return DecodeQuotedPrintable(ctx, r)
	case "8bit", "binary":
		return DecodeEightBit(r)
	default:
		return DecodeSevenBit(r)
	}
}
