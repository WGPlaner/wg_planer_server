package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateGroupCode(t *testing.T) {
	assert.NoError(t, PrepareTestDatabase())

	code1, err1 := CreateGroupCode("00112233-4455-6677-8899-000000000000")
	assert.Error(t, err1)
	assert.Nil(t, code1)
	AssertCount(t, new(GroupCode), 1)

	code2, err2 := CreateGroupCode("00112233-4455-6677-8899-aabbccddeeff")
	assert.NoError(t, err2)
	assert.NotNil(t, code2)
	AssertExistsAndLoadBean(t, &GroupCode{Code: code2.Code})
	AssertCount(t, new(GroupCode), 2)
}

func TestGetGroupCode(t *testing.T) {
	assert.NoError(t, PrepareTestDatabase())

	code1, err1 := GetGroupCode("invalid")
	assert.True(t, IsErrGroupCodeNotExist(err1))
	assert.Nil(t, code1)

	code2, err2 := GetGroupCode("EZ14BAG6T3RG")
	assert.True(t, IsErrGroupCodeNotExist(err2))
	assert.Nil(t, code2)

	code3, err3 := GetGroupCode("ABCDEFGHI123")
	assert.NoError(t, err3)
	assert.NotNil(t, code3)
}

func TestIsGroupCodeValid(t *testing.T) {
	assert.NoError(t, PrepareTestDatabase())

	exists1, code1 := IsGroupCodeValid("invalid")
	assert.False(t, exists1)
	assert.Nil(t, code1)

	exists2, code2 := IsGroupCodeValid("EZ14BAG6T3RG")
	assert.False(t, exists2)
	assert.Nil(t, code2)

	exists3, code3 := IsGroupCodeValid("ABCDEFGHI123")
	assert.True(t, exists3)
	assert.NotNil(t, code3)
}
