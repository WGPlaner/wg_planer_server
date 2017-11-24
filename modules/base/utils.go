package base

import (
	"crypto/rand"
	"math/big"
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

func GetRandomAlphaNumCode(n int, onlyUpperCase bool) string {
	const (
		alphaUpper = "ABCDEFGHIJKLMNOPRSTUVWXYZ"
		alphaLower = "abcdefghiklmnoprstuvwxyz"
		num        = "0123456789"
	)
	var values = alphaUpper + num

	if !onlyUpperCase {
		values += alphaLower
	}

	buffer := make([]byte, n)
	max := big.NewInt(int64(len(values)))

	for i := 0; i < n; i++ {
		buffer[i] = values[randomInt(max)]
	}

	return string(buffer)
}

func randomInt(max *big.Int) int {
	r, err := rand.Int(rand.Reader, max)
	if err != nil {
		panic("ran out of randomness")
	}

	return int(r.Int64())
}

func AppendUniqueString(list []string, str string) []string {
	if !StringInSlice(str, list) {
		return append(list, str)
	}
	return list
}
