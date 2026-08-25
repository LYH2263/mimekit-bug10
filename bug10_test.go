package mimekit

import (
	"context"
	"errors"
	"testing"
)

func TestBug10_CharsetFailNoHalfMessage(t *testing.T) {
	raw := []byte("Content-Type: text/plain; charset=x-unknown\r\n\r\n" + string([]byte{0xff, 0xfe, 0xfd}))
	msg, err := ParseBytes(context.Background(), raw)
	if err == nil {
		t.Fatal("expected charset failure")
	}
	if !errors.Is(err, ErrCharset) {
		t.Fatalf("want ErrCharset chain, got %v", err)
	}
	if msg != nil {
		t.Fatalf("must not return half Message on charset fail, body=%q", msg.Root.Body)
	}
}
