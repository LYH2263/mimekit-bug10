package transfer

import "io"

func DecodeSevenBit(r io.Reader) ([]byte, error) {
	return io.ReadAll(r)
}

func EncodeSevenBit(w io.Writer, raw []byte) error {
	_, err := w.Write(raw)
	return err
}

func DecodeEightBit(r io.Reader) ([]byte, error) {
	return io.ReadAll(r)
}
