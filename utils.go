package id3

import "strings"

const invalidChars = string(uint(0)) + string(uint(1)) + " "

func trim(s string) string {
	return strings.Trim(s, invalidChars)
}
