package charset

func Encode(raw []byte, label string) ([]byte, error) {
	switch Normalize(label) {
	case "utf-8", "us-ascii":
		return append([]byte(nil), raw...), nil
	case "iso-8859-1":
		out := make([]byte, len(raw))
		copy(out, raw)
		return out, nil
	case "gbk":
		return encodeGBK(raw)
	default:
		return append([]byte(nil), raw...), nil
	}
}

func encodeGBK(raw []byte) ([]byte, error) {
	var out []byte
	for _, r := range string(raw) {
		if r < 0x80 {
			out = append(out, byte(r))
			continue
		}
		out = append(out, byte(r>>8), byte(r&0xff))
	}
	return out, nil
}
