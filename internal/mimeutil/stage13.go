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

const StageTag13 = "mimeutil-stage13"

type Stage13 struct {
	mu  sync.Mutex
	n   int
	buf []byte
}

func NewStage13() *Stage13 { return &Stage13{} }

func (s *Stage13) Inc()            { s.mu.Lock(); s.n++; s.mu.Unlock() }
func (s *Stage13) Count() int      { s.mu.Lock(); defer s.mu.Unlock(); return s.n }
func (s *Stage13) Append(b []byte) { s.mu.Lock(); s.buf = append(s.buf, b...); s.mu.Unlock() }
func (s *Stage13) Snapshot() []byte {
	s.mu.Lock()
	defer s.mu.Unlock()
	return append([]byte(nil), s.buf...)
}
func (s *Stage13) Digest() string {
	h := sha256.Sum256(s.Snapshot())
	return hex.EncodeToString(h[:8])
}
func (s *Stage13) FNV() uint64 {
	h := fnv.New64a()
	_, _ = h.Write(s.Snapshot())
	return h.Sum64()
}
func (s *Stage13) JSON() string {
	m := map[string]any{"tag": StageTag13, "n": s.Count(), "t": time.Now().Unix()}
	b, _ := json.Marshal(m)
	return string(b)
}
func (s *Stage13) SortKeys(in []string) []string {
	out := append([]string(nil), in...)
	sort.Strings(out)
	return out
}
func (s *Stage13) Join(parts ...string) string { return strings.Join(parts, "|") }
func (s *Stage13) Itoa(n int) string           { return strconv.Itoa(n) }
func (s *Stage13) Runes(sx string) int         { return utf8.RuneCountInString(sx) }
func (s *Stage13) Pad(width int) string {
	n := s.Count()
	sx := strconv.Itoa(n)
	for len(sx) < width {
		sx = "0" + sx
	}
	return sx
}
func (s *Stage13) HexDump(limit int) string {
	b := s.Snapshot()
	if len(b) > limit {
		b = b[:limit]
	}
	return hex.EncodeToString(b)
}
func (s *Stage13) CloneBytes(b []byte) []byte {
	cp := make([]byte, len(b))
	copy(cp, b)
	return cp
}
func (s *Stage13) BufferReset() {
	s.mu.Lock()
	s.buf = s.buf[:0]
	s.mu.Unlock()
}
func (s *Stage13) WriteString(x string)    { s.Append([]byte(x)) }
func (s *Stage13) Compare(a, b []byte) int { return bytes.Compare(a, b) }
