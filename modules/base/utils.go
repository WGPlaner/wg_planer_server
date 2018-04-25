package base

import (
	"crypto/rand"
	"math/big"
	"os"
	"os/exec"
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

func RemoveStringFromSlice(s []string, str string) []string {
	index := -1
	for i, part := range s {
		if part == str {
			index = i
		}
	}
	if index > -1 {
		if index >= len(s)-1 {
			s = s[:index]
		} else {
			s = append(s[:index], s[index+1:]...)
		}
	}
	return s
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

// IsTestCrasher is used for unit tests and checks if the current process is a crasher test.
// See https://stackoverflow.com/a/33404435
func IsTestCrasher() bool {
	return os.Getenv("UNIT_CRASHER") == "1"
}

// DoesFuncCrash is used for unit tests and checks if the function crashes.
// See https://stackoverflow.com/a/33404435
func DoesFuncCrash(f string) bool {
	cmd := exec.Command(os.Args[0], "-test.run="+f)
	cmd.Env = append(os.Environ(), "UNIT_CRASHER=1")
	err := cmd.Run()
	if e, ok := err.(*exec.ExitError); ok && !e.Success() {
		return true
	}
	return false
}
