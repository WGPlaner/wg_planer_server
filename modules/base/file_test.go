package base

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetMimeType(t *testing.T) {
	// PNG
	d1 := []byte{0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A}
	assert.Equal(t, "image/png", GetMimeType(d1))

	d2 := []byte{0x10, 0x10, 0x10, 0x10, 0x10, 0x10, 0x10, 0x10}
	assert.Equal(t, "application/octet-stream", GetMimeType(d2))
}

func TestIsFileJPG(t *testing.T) {
	// PNG
	d1 := []byte{0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A}
	assert.False(t, IsFileJPG(d1))

	// JPEG
	d2 := []byte{0xFF, 0xD8, 0xFF}
	assert.True(t, IsFileJPG(d2))
}

func TestIsFileImage(t *testing.T) {
	// PNG
	d1 := []byte{0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A}
	assert.True(t, IsFileImage(d1))

	// JPEG
	d2 := []byte{0xFF, 0xD8, 0xFF}
	assert.True(t, IsFileImage(d2))

	// Random data
	d3 := []byte{0x10, 0x10, 0x10, 0x10, 0x10, 0x10, 0x10, 0x10}
	assert.False(t, IsFileImage(d3))
}

func TestFileMustExist(t *testing.T) {
	if IsTestCrasher() {
		FileMustExist("/random/absolute/path")
		return
	}
	assert.True(t, DoesFuncCrash("TestFileMustExist"))
}
