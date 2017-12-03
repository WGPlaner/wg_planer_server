package models

import (
	"testing"

	"github.com/go-openapi/strfmt"
	"github.com/stretchr/testify/assert"
)

func TestGetListItemByUIDs(t *testing.T) {
	assert.NoError(t, PrepareTestDatabase())
	groupUID := strfmt.UUID("00112233-4455-6677-8899-aabbccddeeff")

	itemUID1 := strfmt.UUID("00112233-4455-6677-8899-000000000001")
	item1, err1 := GetListItemByUIDs(groupUID, itemUID1)
	assert.NoError(t, err1)
	assert.Equal(t, "Milk", *item1.Title)

	itemUID2 := strfmt.UUID("00112233-4455-6677-0000-000000000000")
	item2, err2 := GetListItemByUIDs(groupUID, itemUID2)
	assert.Error(t, err2)
	assert.Nil(t, item2)
}
