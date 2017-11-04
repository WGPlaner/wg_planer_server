package wgplaner

import "testing"

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
