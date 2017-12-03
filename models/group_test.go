package models

import (
	"testing"

	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
	"github.com/stretchr/testify/assert"
)

func TestIsGroupExist(t *testing.T) {
	assert.NoError(t, PrepareTestDatabase())
	uid1 := strfmt.UUID("invalid")
	uid2 := strfmt.UUID("00112233-4455-6677-8899-000000000000")
	uid3 := strfmt.UUID("00112233-4455-6677-8899-aabbccddeeff")

	exist1, err1 := IsGroupExist(uid1)
	exist2, err2 := IsGroupExist(uid2)
	exist3, err3 := IsGroupExist(uid3)

	assert.Error(t, err1)
	assert.False(t, exist1)

	assert.NoError(t, err2)
	assert.False(t, exist2)

	assert.NoError(t, err3)
	assert.True(t, exist3)
}

func TestGroup_MarshalBinary(t *testing.T) {
	assert.NoError(t, PrepareTestDatabase())

	g := Group{DisplayName: swag.String("G1"), Currency: "€"}
	b1, err1 := g.MarshalBinary()
	assert.NoError(t, err1)
	assert.NotEmpty(t, b1)
}

func TestGroup_UnmarshalBinary(t *testing.T) {
	assert.NoError(t, PrepareTestDatabase())

	g1a := Group{DisplayName: swag.String("G1"), Currency: "€"}
	g1b := Group{}
	b1, _ := g1a.MarshalBinary()
	err1b := g1b.UnmarshalBinary(b1)
	assert.NoError(t, err1b)
	assert.Equal(t, g1a, g1b)
}
