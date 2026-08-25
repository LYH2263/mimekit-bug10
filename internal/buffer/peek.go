package buffer

import "io"

type PeekReader struct {
	r   io.Reader
	buf []byte
}

func NewPeekReader(r io.Reader) *PeekReader {
	return &PeekReader{r: r}
}

func (p *PeekReader) Peek(n int) ([]byte, error) {
	for len(p.buf) < n {
		chunk := make([]byte, n-len(p.buf))
		k, err := p.r.Read(chunk)
		if k > 0 {
			p.buf = append(p.buf, chunk[:k]...)
		}
		if err != nil {
			return p.buf, err
		}
	}
	return p.buf[:n], nil
}

func (p *PeekReader) Read(b []byte) (int, error) {
	if len(p.buf) > 0 {
		n := copy(b, p.buf)
		p.buf = p.buf[n:]
		return n, nil
	}
	return p.r.Read(b)
}
