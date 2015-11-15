package utils

import (
	"fmt"
	"os/exec"
	"runtime"
)

func CommaFormat(num int64) string {
	src := []byte(fmt.Sprintf("%d", num))

	sign := src[0]
	if sign == '-' || sign == '+' {
		src = src[1:]
	} else {
		sign = 0
	}

	l := len(src)

	buff := make([]byte, 0, l+(l-1)/3+1)

	if sign > 0 {
		buff = append(buff, sign)
	}

	var size int
	for i := 0; i < l; i += size {
		if i == 0 {
			size = l % 3
			if size == 0 {
				size = 3
			}
		} else {
			size = 3
			buff = append(buff, ',')
		}

		buff = append(buff, src[i:i+size]...)
	}

	return string(buff)
}

func OpenBrowser(url string) error {
	switch runtime.GOOS {
	case "linux":
		exec.Command("xdg-open", url).Start()
	case "windows":
		exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		exec.Command("open", url).Start()
	default:
		return fmt.Errorf("%s is not supported.", runtime.GOOS)
	}

	return nil
}
