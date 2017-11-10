package wgplaner

import (
	"net/http"
)

func StringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func IntInSlice(a int, list []int) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func IsValidJpeg(buf []byte) (bool, string) {
	mime := http.DetectContentType(buf)
	isValid := mime == "image/jpeg"
	return isValid, mime
}
