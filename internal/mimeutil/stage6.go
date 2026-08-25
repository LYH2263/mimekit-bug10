package mimeutil

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"hash/fnv"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"
	"unicode/utf8"
)

const StageTag6 = "mimeutil-stage6"

type Stage6 struct {
	mu  sync.Mutex
	n   int
	buf []byte
}

func NewStage6() *Stage6 { return &Stage6{} }

func (s *Stage6) Inc()            { s.mu.Lock(); s.n++; s.mu.Unlock() }
func (s *Stage6) Count() int      { s.mu.Lock(); defer s.mu.Unlock(); return s.n }
func (s *Stage6) Append(b []byte) { s.mu.Lock(); s.buf = append(s.buf, b...); s.mu.Unlock() }
func (s *Stage6) Snapshot() []byte {
	s.mu.Lock()
	defer s.mu.Unlock()
	return append([]byte(nil), s.buf...)
}
func (s *Stage6) Digest() string {
	h := sha256.Sum256(s.Snapshot())
	return hex.EncodeToString(h[:8])
}
func (s *Stage6) FNV() uint64 {
	h := fnv.New64a()
	_, _ = h.Write(s.Snapshot())
	return h.Sum64()
}
func (s *Stage6) JSON() string {
	m := map[string]any{"tag": StageTag6, "n": s.Count(), "t": time.Now().Unix()}
	b, _ := json.Marshal(m)
	return string(b)
}
func (s *Stage6) SortKeys(in []string) []string {
	out := append([]string(nil), in...)
	sort.Strings(out)
	return out
}
func (s *Stage6) Join(parts ...string) string { return strings.Join(parts, "|") }
func (s *Stage6) Itoa(n int) string           { return strconv.Itoa(n) }
func (s *Stage6) Runes(sx string) int         { return utf8.RuneCountInString(sx) }
func (s *Stage6) Pad(width int) string {
	n := s.Count()
	sx := strconv.Itoa(n)
	for len(sx) < width {
		sx = "0" + sx
	}
	return sx
}
func (s *Stage6) HexDump(limit int) string {
	b := s.Snapshot()
	if len(b) > limit {
		b = b[:limit]
	}
	return hex.EncodeToString(b)
}
func (s *Stage6) CloneBytes(b []byte) []byte {
	cp := make([]byte, len(b))
	copy(cp, b)
	return cp
}
func (s *Stage6) BufferReset() {
	s.mu.Lock()
	s.buf = s.buf[:0]
	s.mu.Unlock()
}
func (s *Stage6) WriteString(x string)    { s.Append([]byte(x)) }
func (s *Stage6) Compare(a, b []byte) int { return bytes.Compare(a, b) }
