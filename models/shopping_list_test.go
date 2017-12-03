package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestShoppingList_MarshalBinary(t *testing.T) {
	assert.NoError(t, PrepareTestDatabase())

	s := ShoppingList{}
	b1, err1 := s.MarshalBinary()
	assert.NoError(t, err1)
	assert.NotEmpty(t, b1)
}

func TestShoppingList_UnmarshalBinary(t *testing.T) {
	assert.NoError(t, PrepareTestDatabase())

	s1a := ShoppingList{Count: 2, ListItems: []*ListItem{{Price: 10}, {Price: 33}}}
	s1b := ShoppingList{}
	b1, _ := s1a.MarshalBinary()
	err1b := s1b.UnmarshalBinary(b1)
	assert.NoError(t, err1b)
	assert.Equal(t, s1a, s1b)
}
