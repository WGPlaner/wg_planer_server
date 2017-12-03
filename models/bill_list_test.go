package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBillList_MarshalBinary(t *testing.T) {
	assert.NoError(t, PrepareTestDatabase())

	list := BillList{Count: 2, Bills: []*Bill{{Sum: 10}, {Sum: 33}}}
	b1, err1 := list.MarshalBinary()
	assert.NoError(t, err1)
	assert.NotEmpty(t, b1)
}

func TestBillList_UnmarshalBinary(t *testing.T) {
	assert.NoError(t, PrepareTestDatabase())

	list1a := BillList{Count: 2, Bills: []*Bill{{Sum: 10}, {Sum: 33}}}
	list1b := BillList{}
	b1, _ := list1a.MarshalBinary()
	err1b := list1b.UnmarshalBinary(b1)
	assert.NoError(t, err1b)
	assert.Equal(t, list1a, list1b)
}
