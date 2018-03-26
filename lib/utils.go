package lib

import "strings"

const invalidChars = string(uint(0)) + string(uint(1)) + " "

func Trim(s string) string {
	return strings.Trim(s, invalidChars)
}

func ByteToInt(b []byte) int {
	size:= 0
	for i:= range b {
		size =size<<8 + int(b[i])
	}

	return size
}
