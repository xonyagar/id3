package lib

import (
	"bytes"
	"unicode/utf16"
	"unicode/utf8"
)

type Encoding struct {
	Title string
	Size  int
}

var Encodings = []Encoding{
	{"ISO-8859-1", 1},
	{"UTF-16", 2},
	{"UTF-16BE", 2},
	{"UTF-8", 1},
}

func ToUTF8(data []byte, enc Encoding) string {
	switch enc.Title {
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
		i := 0

		if lb%2 != 0 && data[i] == 0 {
			i++
		} else {
			lb--
		}

		for ; i < lb; i += 2 {
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
