package workbench

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/LYH2263/go-mimekit"
)

type API struct {
	Pipeline *mimekit.Pipeline
}

func (a *API) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case "/api/parse":
		a.handleParse(w, r)
	case "/api/encode":
		a.handleEncode(w, r)
	case "/api/stats":
		a.handleStats(w, r)
	default:
		http.NotFound(w, r)
	}
}

func (a *API) handleParse(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method", http.StatusMethodNotAllowed)
		return
	}
	raw, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	msg, err := mimekit.Parse(r.Context(), bytesNewReader(raw))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	writeJSON(w, map[string]any{
		"parts": msg.PartCount(),
		"text":  msg.PlainText(),
	})
}

func (a *API) handleEncode(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method", http.StatusMethodNotAllowed)
		return
	}
	msg := mimekit.NewTextPlain(string(mustRead(r)))
	out, err := mimekit.EncodeString(msg)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	writeJSON(w, map[string]any{"raw": out})
}

func (a *API) handleStats(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, map[string]any{"pipeline": "mimekit", "ok": true})
}

func writeJSON(w http.ResponseWriter, v any) {
	w.Header().Set("Content-Type", "application/json")
	enc := json.NewEncoder(w)
	_ = enc.Encode(v)
	if f, ok := w.(http.Flusher); ok {
		f.Flush()
	}
}

func bytesNewReader(b []byte) *sliceReader { return &sliceReader{b: b} }

type sliceReader struct {
	b []byte
	i int
}

func (s *sliceReader) Read(p []byte) (int, error) {
	if s.i >= len(s.b) {
		return 0, io.EOF
	}
	n := copy(p, s.b[s.i:])
	s.i += n
	return n, nil
}

func mustRead(r *http.Request) []byte {
	b, _ := io.ReadAll(r.Body)
	return b
}
