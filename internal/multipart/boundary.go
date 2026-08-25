package multipart

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
)

func Generate(prefix string) (string, error) {
	if prefix == "" {
		prefix = "mimekit"
	}
	var b [12]byte
	if _, err := rand.Read(b[:]); err != nil {
		return "", err
	}
	return fmt.Sprintf("%s_%s", prefix, hex.EncodeToString(b[:])), nil
}

func Marker(boundary string) []byte {
	return []byte("\r\n--" + boundary)
}

func CloseMarker(boundary string) []byte {
	return []byte("\r\n--" + boundary + "--\r\n")
}
