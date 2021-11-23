package lib

import (
	"fmt"
	"strings"
)

var (
	dSizes = []string{"Bytes", "KB", "MB", "GB"}
	bSizes = []string{"Bytes", "KiB", "MiB", "GiB"}
)

func HumanDecimalSize(n int) string {
	i := 0
	f := float32(n)

	for ; i < len(dSizes); i++ {
		if f < 1000 {
			break
		}

		f /= 1000
	}

	return fmt.Sprintf("%s %s", strings.TrimRight(strings.TrimRight(fmt.Sprintf("%.2f", f), "0"), "."), dSizes[i])
}

func HumanBinarySize(n int) string {
	i := 0
	f := float32(n)

	for ; i < len(bSizes); i++ {
		if f < 1024 {
			break
		}

		f /= 1024
	}

	return fmt.Sprintf("%s %s", strings.TrimRight(strings.TrimRight(fmt.Sprintf("%.2f", f), "0"), "."), bSizes[i])
}
