package transfer

import (
	"io"
	"strings"
)

func Encode(w io.Writer, raw []byte, encoding string) error {
	enc := strings.ToLower(strings.TrimSpace(encoding))
	switch enc {
	case "base64", "b64":
		return EncodeBase64(w, raw)
	case "quoted-printable", "qp":
		return EncodeQuotedPrintable(w, raw)
	default:
		return EncodeSevenBit(w, raw)
	}
}
