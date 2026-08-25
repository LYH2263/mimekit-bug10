package multipart

import (
	"bytes"
	"context"
	"errors"
	"io"
	"strings"
)

var (
	ErrPartDone = errors.New("multipart: done")
	ErrBadPart  = errors.New("multipart: bad part")
)

type Reader struct {
	ctx      context.Context
	data     []byte
	boundary string
	off      int
	done     bool
}

func NewReader(ctx context.Context, raw []byte, boundary string) *Reader {
	return &Reader{ctx: ctx, data: raw, boundary: boundary}
}

func (r *Reader) Next() ([]byte, map[string][]string, error) {
	if err := r.ctx.Err(); err != nil {
		return nil, nil, err
	}
	if r.done {
		return nil, nil, ErrPartDone
	}
	delim := []byte("--" + r.boundary)
	start := bytes.Index(r.data[r.off:], delim)
	if start < 0 {
		return nil, nil, ErrPartDone
	}
	r.off += start + len(delim)
	if r.off+2 <= len(r.data) && r.data[r.off] == '-' && r.data[r.off+1] == '-' {
		r.done = true
		return nil, nil, ErrPartDone
	}
	if r.off < len(r.data) && r.data[r.off] == '\r' {
		r.off += 2
	} else if r.off < len(r.data) && r.data[r.off] == '\n' {
		r.off++
	}
	end := bytes.Index(r.data[r.off:], delim)
	chunk := r.data[r.off:]
	if end >= 0 {
		chunk = r.data[r.off : r.off+end]
	}
	split := bytes.Index(chunk, []byte("\r\n\r\n"))
	if split < 0 {
		split = bytes.Index(chunk, []byte("\n\n"))
	}
	if split < 0 {
		return nil, nil, ErrBadPart
	}
	hdrBlock := chunk[:split]
	body := chunk[split+2:]
	if bytes.HasPrefix(body, []byte("\r\n")) {
		body = body[2:]
	}
	headers, err := headerParse(chunk[:split+2])
	if err != nil {
		return body, headers, err
	}
	_ = hdrBlock
	if end >= 0 {
		r.off += end
	}
	return body, headers, nil
}

func headerParse(block []byte) (map[string][]string, error) {
	return parseHeaderBytes(block)
}

func parseHeaderBytes(block []byte) (map[string][]string, error) {
	text := string(block)
	text = strings.ReplaceAll(text, "\r\n", "\n")
	lines := strings.Split(text, "\n")
	h := make(map[string][]string)
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		colon := strings.IndexByte(line, ':')
		if colon <= 0 {
			continue
		}
		key := strings.TrimSpace(line[:colon])
		val := strings.TrimSpace(line[colon+1:])
		h[key] = append(h[key], val)
	}
	return h, nil
}

func ReadAllParts(ctx context.Context, raw []byte, boundary string) ([]Section, error) {
	r := NewReader(ctx, raw, boundary)
	var parts []Section
	for {
		body, hdr, err := r.Next()
		if err == ErrPartDone {
			break
		}
		if err != nil {
			return parts, err
		}
		p := Section{
			Headers:   hdr,
			RawBody:   append([]byte(nil), body...),
			Body:      append([]byte(nil), body...),
			MediaType: mediaFromHeaders(hdr),
		}
		parts = append(parts, p)
	}
	return parts, nil
}

type Section struct {
	Headers   map[string][]string
	RawBody   []byte
	Body      []byte
	MediaType string
}

func mediaFromHeaders(h map[string][]string) string {
	for k, v := range h {
		if strings.EqualFold(k, "Content-Type") && len(v) > 0 {
			media, _, _ := parseMediaType(v[0])
			return media
		}
	}
	return "text/plain"
}

func parseMediaType(ct string) (string, map[string]string, error) {
	parts := strings.SplitN(ct, ";", 2)
	media := strings.TrimSpace(parts[0])
	params := map[string]string{}
	if len(parts) == 2 {
		for _, seg := range strings.Split(parts[1], ";") {
			kv := strings.SplitN(strings.TrimSpace(seg), "=", 2)
			if len(kv) == 2 {
				params[strings.ToLower(strings.TrimSpace(kv[0]))] = strings.Trim(strings.TrimSpace(kv[1]), "\"")
			}
		}
	}
	return media, params, nil
}

func Drain(r io.Reader) ([]byte, error) {
	return io.ReadAll(r)
}
