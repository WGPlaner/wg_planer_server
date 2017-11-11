package wgplaner

import (
	"strings"
	"testing"
)

func TestStringInSlice(t *testing.T) {
	strings := []string{"Hello", "World", "!", "This", "is", "a", "test"}

	if StringInSlice("hello", strings) != false {
		t.Error("Matching should be case-sensitive")
	}
	if StringInSlice("!", strings) != true {
		t.Error("Single character string should match")
	}
	if StringInSlice("what", strings) != false {
		t.Error("Expected false")
	}
}

func TestIntInSlice(t *testing.T) {
	ints := []int{-10, 80, 443, 25, 22, 20, 21}

	if IntInSlice(-10, ints) != true {
		t.Error("Matching should work for negative numbers")
	}
	if IntInSlice(0, ints) != false && IntInSlice(100, ints) != false {
		t.Error("Expected false")
	}
}

func TestRandomAlphaNumCode(t *testing.T) {
	if str := RandomAlphaNumCode(10, false); len(str) != 10 {
		t.Error("Length should be as specified")
	}
	if str := RandomAlphaNumCode(10, true); len(str) != 10 {
		t.Error("Length should be as specified")
	}
	if str := RandomAlphaNumCode(10, true); str != strings.ToUpper(str) {
		t.Error("Only uppercase must only contain uppercase letters")
	}
}
