package mimekit

import "errors"

var (
	ErrClosed      = errors.New("mimekit: pipeline closed")
	ErrInvalidMIME = errors.New("mimekit: invalid mime")
	ErrNoBoundary  = errors.New("mimekit: missing boundary")
	ErrPartDone    = errors.New("mimekit: part stream done")
	ErrCanceled    = errors.New("mimekit: canceled")
	ErrTransfer    = errors.New("mimekit: transfer decode failed")
	ErrCharset     = errors.New("mimekit: charset conversion failed")
)
