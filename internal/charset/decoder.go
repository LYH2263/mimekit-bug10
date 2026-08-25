package charset

import (
	"context"
	"sync"
	"unicode/utf8"
)

var (
	decMu    sync.RWMutex
	decCache = map[string]func([]byte) ([]byte, error){}
)

func Decode(ctx context.Context, raw []byte, label string) ([]byte, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	key := Normalize(label)
	decMu.RLock()
	fn, ok := decCache[key]
	decMu.RUnlock()
	if !ok {
		fn = buildDecoder(key)
	}
	out, err := fn(raw)
	if err != nil {
		// 失败不得把「吞错 identity」写进缓存，避免污染后续同标签解码。
		return nil, err
	}
	if !ok {
		decMu.Lock()
		decCache[key] = fn
		decMu.Unlock()
	}
	return out, nil
}

func buildDecoder(label string) func([]byte) ([]byte, error) {
	switch label {
	case "utf-8", "us-ascii":
		return func(b []byte) ([]byte, error) {
			return append([]byte(nil), b...), nil
		}
	case "iso-8859-1":
		return func(b []byte) ([]byte, error) {
			var out []byte
			for _, c := range b {
				out = append(out, byte(c))
			}
			return out, nil
		}
	case "gbk":
		return func(b []byte) ([]byte, error) {
			return decodeGBK(b)
		}
	default:
		return func(b []byte) ([]byte, error) {
			if !utf8.Valid(b) {
				return nil, ErrCharset
			}
			return append([]byte(nil), b...), nil
		}
	}
}

var ErrCharset = errCharset{}

type errCharset struct{}

func (errCharset) Error() string { return "charset: conversion failed" }

func decodeGBK(b []byte) ([]byte, error) {
	var out []byte
	for i := 0; i < len(b); {
		if b[i] < 0x80 {
			out = append(out, b[i])
			i++
			continue
		}
		if i+1 >= len(b) {
			return nil, ErrCharset
		}
		r := rune(int(b[i])<<8 | int(b[i+1]))
		buf := make([]byte, 4)
		n := utf8.EncodeRune(buf, r)
		out = append(out, buf[:n]...)
		i += 2
	}
	return out, nil
}

func ClearCache() {
	decMu.Lock()
	decCache = map[string]func([]byte) ([]byte, error){}
	decMu.Unlock()
}
