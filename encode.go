package mimekit

import (
	"bytes"
	"io"
	"strings"

	"github.com/LYH2263/go-mimekit/internal/header"
	"github.com/LYH2263/go-mimekit/internal/multipart"
	"github.com/LYH2263/go-mimekit/internal/transfer"
)

func Encode(w io.Writer, m *Message, opts EncodeOptions) error {
	if m == nil || m.Root == nil {
		return ErrInvalidMIME
	}
	return encodePart(w, m.Root, opts)
}

func encodePart(w io.Writer, p *Part, opts EncodeOptions) error {
	for k, vals := range p.Headers {
		for _, v := range vals {
			line := header.FoldLine(k+": "+v, opts.FoldAt)
			if _, err := io.WriteString(w, line+"\r\n"); err != nil {
				return err
			}
		}
	}
	if len(p.Children) > 0 {
		mw := multipart.NewWriter(w, opts.BoundaryPrefix)
		for _, c := range p.Children {
			if err := mw.OpenPart(c.Headers); err != nil {
				mw.Abort()
				return err
			}
			te := TransferEncoding(c.Headers)
			if err := transfer.Encode(w, c.Body, te); err != nil {
				mw.Abort()
				return err
			}
			if err := mw.ClosePart(); err != nil {
				mw.Abort()
				return err
			}
		}
		if err := mw.Flush(); err != nil {
			mw.Abort()
			return err
		}
		return mw.Close()
	}
	te := TransferEncoding(p.Headers)
	if err := transfer.Encode(w, p.Body, te); err != nil {
		return err
	}
	return nil
}

func EncodeString(m *Message) (string, error) {
	var buf bytes.Buffer
	if err := Encode(&buf, m, EncodeOptions{FoldAt: 76}); err != nil {
		return "", err
	}
	return buf.String(), nil
}

func extractBoundary(ct string) (string, error) {
	_, params, err := parseMediaParams(ct)
	if err != nil {
		return "", err
	}
	b, ok := params["boundary"]
	if !ok || b == "" {
		return "", ErrNoBoundary
	}
	return b, nil
}

func parseMediaParams(ct string) (string, map[string]string, error) {
	parts := strings.SplitN(ct, ";", 2)
	media := strings.TrimSpace(parts[0])
	params := map[string]string{}
	if len(parts) == 2 {
		for _, seg := range strings.Split(parts[1], ";") {
			seg = strings.TrimSpace(seg)
			if seg == "" {
				continue
			}
			kv := strings.SplitN(seg, "=", 2)
			if len(kv) != 2 {
				continue
			}
			k := strings.TrimSpace(kv[0])
			v := strings.Trim(strings.TrimSpace(kv[1]), "\"")
			params[strings.ToLower(k)] = v
		}
	}
	return media, params, nil
}
