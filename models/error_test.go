package models

import (
	"errors"
	"testing"
)

func TestErrorList_Add(t *testing.T) {
	errList := ErrorList{}
	errList.Add("First error")
	if len(errList.errors) != 1 {
		t.Error("Should contain 1 error!")
	}
}

func TestErrorList_AddError(t *testing.T) {
	errList := ErrorList{}
	errList.AddError(errors.New("First error"))
	if len(errList.errors) != 1 {
		t.Error("Should contain 1 error!")
	}
}

func TestErrorList_AddList(t *testing.T) {
	destList := ErrorList{}
	srcList := ErrorList{
		errors: []error{
			errors.New("First"),
			errors.New("Second"),
		},
	}

	destList.AddList(&srcList)

	if len(destList.errors) != 2 {
		t.Error("Should contain 2 copied errors!")
	}

	if destList.errors[0] != srcList.errors[0] ||
		destList.errors[1] != srcList.errors[1] {
		t.Error("Elements should be equal!")
	}
}
