package models

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
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

func TestErrorList_String(t *testing.T) {
	errList := ErrorList{}
	errList.Add("first")
	errList.Add("second")
	assert.Equal(t, "first\nsecond\n", errList.String())
}

func TestErrorList_HasErrors(t *testing.T) {
	errList1 := ErrorList{}
	errList1.Add("first")
	assert.True(t, errList1.HasErrors())

	errList2 := ErrorList{}
	assert.False(t, errList2.HasErrors())
}
