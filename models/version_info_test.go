package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestVersionInfo_MarshalBinary(t *testing.T) {
	assert.NoError(t, PrepareTestDatabase())

	info := VersionInfo{AndroidVersionString: "Msg", AndroidVersionCode: 999}
	b1, err1 := info.MarshalBinary()
	assert.NoError(t, err1)
	assert.NotEmpty(t, b1)
}

func TestVersionInfo_UnmarshalBinary(t *testing.T) {
	assert.NoError(t, PrepareTestDatabase())

	info1a := VersionInfo{AndroidVersionString: "Msg", AndroidVersionCode: 999}
	info1b := VersionInfo{}
	b1, _ := info1a.MarshalBinary()
	err1b := info1b.UnmarshalBinary(b1)
	assert.NoError(t, err1b)
	assert.Equal(t, info1a, info1b)
}
