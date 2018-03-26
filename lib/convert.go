package lib

func ByteToInt(b []byte) int {
	size := 0
	for i := range b {
		size = size<<8 + int(b[i])
	}

	return size
}
