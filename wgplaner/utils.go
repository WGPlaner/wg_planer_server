package wgplaner

import (
	"math/rand"
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

func RandomAlphaNumCode(length int, onlyUpperCase bool) string {
	alphaUpper := "ABCDEFGHIJKLMNOPRSTUVWXYZ"
	alphaLower := "abcdefghiklmnoprstuvwxyz"
	num := "0123456789"
	values := alphaUpper + num

	if !onlyUpperCase {
		values += alphaLower
	}

	codeStr := ""

	for i := 0; i < length; i++ {
		j := rand.Intn(len(values))
		codeStr += string(values[j])
	}

	return codeStr
}

func AppendUniqueString(list []string, str string) []string {
	if !StringInSlice(str, list) {
		return append(list, str)
	}
	return list
}
