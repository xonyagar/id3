package id3

import (
	"bytes"
	"unicode/utf16"
	"unicode/utf8"
)

type Encoding string

var Encodings = []Encoding{"ISO-8859-1", "UTF-16", "UTF-16BE", "UTF-8"}

func toUTF8(data []byte, enc Encoding) string {
	switch enc {
	case "ISO-8859-1":
		buf := make([]rune, len(data))
		for i, b := range data {
			buf[i] = rune(b)
		}

		return string(buf)
	case "UTF-16":
		u16s := make([]uint16, 1)

		ret := &bytes.Buffer{}

		b8buf := make([]byte, 4)

		lb := len(data)
		for i := 0; i < lb; i += 2 {
			u16s[0] = uint16(data[i]) + (uint16(data[i+1]) << 8)
			r := utf16.Decode(u16s)
			n := utf8.EncodeRune(b8buf, r[0])
			ret.Write(b8buf[:n])
		}

		return ret.String()
		// TODO: check other encodings
	default:
		return string(data)
	}
}
