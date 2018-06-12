package base

import (
	"strings"
	"testing"
)

func TestStringInSlice(t *testing.T) {
	s := []string{"Hello", "World", "!", "This", "is", "a", "test"}

	if StringInSlice("hello", s) != false {
		t.Error("Matching should be case-sensitive")
	}
	if StringInSlice("!", s) != true {
		t.Error("Single character string should match")
	}
	if StringInSlice("what", s) != false {
		t.Error("Expected false")
	}
}

func TestIntInSlice(t *testing.T) {
	ints := []int{-10, 80, 443, 25, 22, 20, 21}

	if IntInSlice(-10, ints) != true {
		t.Error("Matching should work for negative numbers")
	}
	if IntInSlice(0, ints) != false || IntInSlice(100, ints) != false {
		t.Error("Expected false")
	}
}

func TestRandomAlphaNumCode(t *testing.T) {
	if str := GetRandomAlphaNumCode(10, false); len(str) != 10 {
		t.Error("Length should be as specified")
	}
	if str := GetRandomAlphaNumCode(10, true); len(str) != 10 {
		t.Error("Length should be as specified")
	}
	if str := GetRandomAlphaNumCode(10, true); str != strings.ToUpper(str) {
		t.Error("Only uppercase must only contain uppercase letters")
	}
}

func TestAppendUniqueString(t *testing.T) {
	initialList := []string{"One", "Two", "Three"}

	if len(AppendUniqueString(initialList, "Three")) != 3 {
		t.Error("Should only add unique values")
	}

	if len(AppendUniqueString(initialList, "Four")) != 4 {
		t.Error("Should add non-unique values")
	}
}

func TestUnique(t *testing.T) {
	initialList := []string{"one", "two", "one", "two", "one", "three", "four", "three", "three"}
	uniqueList := Unique(initialList)
	if len(uniqueList) != 4 {
		t.Error()
	}
}
