package multipart

import (
	"errors"
	"io"
	"strings"
)

var ErrWriterClosed = errors.New("multipart: writer closed")

type Writer struct {
	w        io.Writer
	boundary string
	opened   bool
	closed   bool
	broken   bool
}

func NewWriter(w io.Writer, prefix string) *Writer {
	b, _ := Generate(prefix)
	return &Writer{w: w, boundary: b}
}

func (mw *Writer) Boundary() string { return mw.boundary }

func (mw *Writer) OpenPart(headers map[string][]string) error {
	if mw.closed || mw.broken {
		return ErrWriterClosed
	}
	if !mw.opened {
		if _, err := io.WriteString(mw.w, "--"+mw.boundary+"\r\n"); err != nil {
			return err
		}
		mw.opened = true
	} else {
		if _, err := io.WriteString(mw.w, "\r\n--"+mw.boundary+"\r\n"); err != nil {
			return err
		}
	}
	for k, vals := range headers {
		for _, v := range vals {
			line := k + ": " + v + "\r\n"
			if _, err := io.WriteString(mw.w, line); err != nil {
				return err
			}
		}
	}
	if _, err := io.WriteString(mw.w, "\r\n"); err != nil {
		return err
	}
	return nil
}

func (mw *Writer) ClosePart() error { return nil }

func (mw *Writer) Flush() error {
	if mw.broken {
		return ErrWriterClosed
	}
	_, err := io.WriteString(mw.w, "")
	return err
}

// Abort marks the writer broken after a mid-part failure so Close must not emit a closing boundary.
func (mw *Writer) Abort() {
	mw.broken = true
	mw.closed = true
}

func (mw *Writer) Close() error {
	if mw.broken {
		return ErrWriterClosed
	}
	if mw.closed {
		return nil
	}
	if !mw.opened {
		return ErrBadPart
	}
	mw.closed = true
	_, err := io.WriteString(mw.w, "\r\n--"+mw.boundary+"--\r\n")
	return err
}

func (mw *Writer) Broken() bool { return mw.broken }

func FormatContentType(media, boundary string) string {
	return strings.TrimSpace(media) + "; boundary=\"" + boundary + "\""
}
