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

func TestGetGroupNotFound(t *testing.T) {
	prepareTestEnv(t)
	authInGroup := "1234567890fakefirebaseid0001"
	req := NewRequest(t, "GET", authInGroup, "/groups/00112233-4455-6677-8899-000000000000")
	MakeRequest(t, req, http.StatusNotFound)
}

func TestGetGroupImageUnauthorized(t *testing.T) {
	prepareTestEnv(t)
	authInGroup := "1234567890fakefirebaseid0003"
	req := NewRequest(t, "GET", authInGroup, "/groups/00112233-4455-6677-8899-aabbccddeeff/image")
	MakeRequest(t, req, http.StatusOK)
}

func TestGetGroupImage(t *testing.T) {
	prepareTestEnv(t)
	authInGroup := "1234567890fakefirebaseid0001"
	req := NewRequest(t, "GET", authInGroup, "/groups/00112233-4455-6677-8899-aabbccddeeff/image")
	MakeRequest(t, req, http.StatusOK)
}

func TestGetGroupNotFoundImage(t *testing.T) {
	prepareTestEnv(t)
	authInGroup := "1234567890fakefirebaseid0001"
	req := NewRequest(t, "GET", authInGroup, "/groups/00112233-4455-6677-8899-000000000000/image")
	MakeRequest(t, req, http.StatusNotFound)
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
	assert.Contains(t, createdGroup.Members, authInGroup)

	// TODO check database beans
}

func TestCreateGroupInvalid(t *testing.T) {
	prepareTestEnv(t)
	req := NewRequestWithJSON(t, "POST", "1234567890fakefirebaseid0003", "/groups", models.Group{})
	MakeRequest(t, req, http.StatusUnprocessableEntity)
}

func TestUpdateGroup(t *testing.T) {
	prepareTestEnv(t)
	var (
		uG = models.Group{}
		g  = models.Group{
			UID:         "00112233-4455-6677-8899-aabbccddeeff",
			DisplayName: swag.String("Updated Group"),
		}
		req  = NewRequestWithJSON(t, "PUT", "1234567890fakefirebaseid0001", "/groups", g)
		resp = MakeRequest(t, req, http.StatusOK)
	)
	DecodeJSON(t, resp, &uG)
	assert.Equal(t, *g.DisplayName, *uG.DisplayName)
	assert.NotEqual(t, uG.UpdatedAt, uG.CreatedAt)
}

func TestUpdateGroupNotFound(t *testing.T) {
	prepareTestEnv(t)
	var (
		g = models.Group{
			UID:         "00112233-4455-6677-8899-000000000000",
			DisplayName: swag.String("Non existent group"),
		}
		req = NewRequestWithJSON(t, "PUT", "1234567890fakefirebaseid0001", "/groups", g)
	)
	MakeRequest(t, req, http.StatusNotFound)
}

func TestUpdateGroupNotAdmin(t *testing.T) {
	prepareTestEnv(t)
	var (
		g = models.Group{
			UID:         "00112233-4455-6677-8899-aabbccddeeff",
			DisplayName: swag.String("Updated Group"),
		}
		req = NewRequestWithJSON(t, "PUT", "1234567890fakefirebaseid0002", "/groups", g)
	)
	MakeRequest(t, req, http.StatusUnauthorized)
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

func TestCreateGroupCodeInvalid(t *testing.T) {
	prepareTestEnv(t)
	var (
		uid = strfmt.UUID("00112233-4455-6677-8899-0000000000")
		url = fmt.Sprintf("/groups/%s/create-code", uid)
		req = NewRequest(t, "GET", "1234567890fakefirebaseid0001", url)
	)
	MakeRequest(t, req, http.StatusBadRequest)
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

func TestJoinGroupThroughCodeInvalid(t *testing.T) {
	prepareTestEnv(t)
	req := NewRequest(t, "POST", "1234567890fakefirebaseid0001", "/groups/join/123123123123")
	MakeRequest(t, req, http.StatusBadRequest)
}

func TestJoinGroupThroughCode(t *testing.T) {
	prepareTestEnv(t)
	var (
		userUid = "1234567890fakefirebaseid0001"
		code    = models.GroupCode{}
		uid     = strfmt.UUID("00112233-4455-6677-8899-aabbccddeeff")
		reqCode = NewRequest(t, "GET", userUid, fmt.Sprintf("/groups/%s/create-code", uid))
		resp    = MakeRequest(t, reqCode, http.StatusOK)
	)
	DecodeJSON(t, resp, &code)
	// Join with generated code
	req := NewRequest(t, "POST", userUid, "/groups/join/"+*code.Code)
	MakeRequest(t, req, http.StatusOK)
}

func TestLeaveGroup(t *testing.T) {
	prepareTestEnv(t)
	authInGroup := "1234567890fakefirebaseid0001"
	req := NewRequest(t, "POST", authInGroup, "/groups/leave")
	MakeRequest(t, req, http.StatusOK)

	u := &models.User{}
	req = NewRequest(t, "GET", authInGroup, "/users/"+authInGroup)
	resp := MakeRequest(t, req, http.StatusOK)
	DecodeJSON(t, resp, u)
	assert.Empty(t, u.GroupUID)
}
