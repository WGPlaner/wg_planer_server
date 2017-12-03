package models

import (
	"testing"

	"github.com/go-openapi/swag"
	"github.com/stretchr/testify/assert"
)

func TestSuccessResponse_MarshalBinary(t *testing.T) {
	assert.NoError(t, PrepareTestDatabase())

	resp := SuccessResponse{Message: swag.String("Msg"), Status: swag.Int64(999)}
	b1, err1 := resp.MarshalBinary()
	assert.NoError(t, err1)
	assert.NotEmpty(t, b1)
}

func TestSuccessResponse_UnmarshalBinary(t *testing.T) {
	assert.NoError(t, PrepareTestDatabase())

	resp1a := SuccessResponse{Message: swag.String("Msg"), Status: swag.Int64(999)}
	resp1b := SuccessResponse{}
	b1, _ := resp1a.MarshalBinary()
	err1b := resp1b.UnmarshalBinary(b1)
	assert.NoError(t, err1b)
	assert.Equal(t, resp1a, resp1b)
}
