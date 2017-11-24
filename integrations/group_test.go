package integrations

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/wgplaner/wg_planer_server/models"

	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
	"github.com/stretchr/testify/assert"
)

func TestGetGroupUnauthorized(t *testing.T) {
	prepareTestEnv(t)
	authNotInGroup := "1234567890fakefirebaseid0003"
	req := NewRequest(t, "GET", authNotInGroup, "/groups/00112233-4455-6677-8899-aabbccddeeff")
	MakeRequest(t, req, http.StatusUnauthorized)
}

func TestGetGroup(t *testing.T) {
	prepareTestEnv(t)
	authInGroup := "1234567890fakefirebaseid0001"
	req := NewRequest(t, "GET", authInGroup, "/groups/00112233-4455-6677-8899-aabbccddeeff")
	MakeRequest(t, req, http.StatusOK)
}

func TestCreateGroup(t *testing.T) {
	prepareTestEnv(t)
	var (
		authInGroup  = "1234567890fakefirebaseid0003"
		createdGroup = models.Group{}
		newGroup     = models.Group{
			DisplayName: swag.String("New Group"),
		}
		req  = NewRequestWithJSON(t, "POST", authInGroup, "/groups", newGroup)
		resp = MakeRequest(t, req, http.StatusOK)
	)
	DecodeJSON(t, resp, &createdGroup)
	assert.Equal(t, *newGroup.DisplayName, *createdGroup.DisplayName)
	assert.NotEmpty(t, createdGroup.UID)
	assert.Equal(t, createdGroup.CreatedAt, createdGroup.UpdatedAt)

	// TODO check database beans
}

func TestCreateGroupInvalid(t *testing.T) {
	prepareTestEnv(t)
	req := NewRequestWithJSON(t, "POST", "1234567890fakefirebaseid0003", "/groups", models.Group{})
	MakeRequest(t, req, http.StatusUnprocessableEntity)
}

func TestCreateGroupCode(t *testing.T) {
	prepareTestEnv(t)
	var (
		code = models.GroupCode{}
		uid  = strfmt.UUID("00112233-4455-6677-8899-aabbccddeeff")
		url  = fmt.Sprintf("/groups/%s/create-code", uid)
		req  = NewRequest(t, "GET", "1234567890fakefirebaseid0001", url)
		resp = MakeRequest(t, req, http.StatusOK)
	)
	DecodeJSON(t, resp, &code)
	assert.Equal(t, *code.GroupUID, uid)
	assert.Len(t, *code.Code, 12)
}

func TestCreateGroupCodeUnauthorized(t *testing.T) {
	prepareTestEnv(t)
	var (
		err  = models.ErrorResponse{}
		uid  = strfmt.UUID("00112233-4455-6677-8899-aabbccddeeff")
		url  = fmt.Sprintf("/groups/%s/create-code", uid)
		req  = NewRequest(t, "GET", "1234567890fakefirebaseid0003", url)
		resp = MakeRequest(t, req, http.StatusUnauthorized)
	)
	DecodeJSON(t, resp, &err)
	assert.NotEmpty(t, *err.Message)
}

func TestGetJoinGroupHelp(t *testing.T) {
	prepareTestEnv(t)
	req := NewRequest(t, "GET", "1234567890fakefirebaseid0003", "/groups/join/123456789ABC")
	MakeRequest(t, req, http.StatusOK)
}

func TestGetJoinGroupHelpInvalid(t *testing.T) {
	prepareTestEnv(t)
	// Invalid code => invalid format
	req := NewRequest(t, "GET", "1234567890fakefirebaseid0003", "/groups/join/1234")
	MakeRequest(t, req, http.StatusBadRequest)
}
