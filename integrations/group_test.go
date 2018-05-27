package integrations

import (
	"net/http"
	"testing"

	"github.com/wgplaner/wg_planer_server/models"

	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
	"github.com/stretchr/testify/assert"
)

func TestGetGroup(t *testing.T) {
	prepareTestEnv(t)
	authInGroup := "1234567890fakefirebaseid0001"
	req := NewRequest(t, "GET", authInGroup, "/group")
	MakeRequest(t, req, http.StatusOK)
}

func TestGetGroupImage(t *testing.T) {
	prepareTestEnv(t)
	authInGroup := "1234567890fakefirebaseid0001"
	req := NewRequest(t, "GET", authInGroup, "/group/image")
	MakeRequest(t, req, http.StatusOK)
}

func TestUpdateGroupImage(t *testing.T) {
	prepareTestEnv(t)

	groupUID := "00112233-4455-6677-8899-aabbccddeeff"
	request := NewRequestWithImage(t,
		"PUT",
		AuthValid,
		"/group/image",
		"profileImage",
		models.GetGroupImageDefaultPath())
	resp := MakeRequest(t, request, http.StatusOK)

	// Check JSON response
	successResp := models.SuccessResponse{}
	DecodeJSON(t, resp, &successResp)
	assert.NotEmpty(t, *successResp.Message)
	assert.Equal(t, int64(200), *successResp.Status)

	// Check that image is uploaded
	_, fileErr := models.GetGroupImage(strfmt.UUID(groupUID))
	assert.NoError(t, fileErr, "Uploaded image not found")
}

func TestCreateGroup(t *testing.T) {
	prepareTestEnv(t)
	var (
		authInGroup  = "1234567890fakefirebaseid0003"
		createdGroup = models.Group{}
		newGroup     = models.Group{
			DisplayName: swag.String("New Group"),
		}
		req  = NewRequestWithJSON(t, "POST", authInGroup, "/group", newGroup)
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
	req := NewRequestWithJSON(t, "POST", "1234567890fakefirebaseid0003", "/group", models.Group{})
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
		req  = NewRequestWithJSON(t, "PUT", "1234567890fakefirebaseid0001", "/group", g)
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
		req = NewRequestWithJSON(t, "PUT", "1234567890fakefirebaseid0001", "/group", g)
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
		req = NewRequestWithJSON(t, "PUT", "1234567890fakefirebaseid0002", "/group", g)
	)
	MakeRequest(t, req, http.StatusUnauthorized)
}

func TestCreateGroupCode(t *testing.T) {
	prepareTestEnv(t)
	var (
		code = models.GroupCode{}
		uid  = strfmt.UUID("00112233-4455-6677-8899-aabbccddeeff")
		url  = "/group/create-code"
		req  = NewRequest(t, "GET", "1234567890fakefirebaseid0001", url)
		resp = MakeRequest(t, req, http.StatusOK)
	)
	DecodeJSON(t, resp, &code)
	assert.Equal(t, *code.GroupUID, uid)
	assert.Len(t, *code.Code, 12)
}

func TestGetJoinGroupHelp(t *testing.T) {
	prepareTestEnv(t)
	req := NewRequest(t, "GET", "1234567890fakefirebaseid0003", "/group/join/123456789ABC")
	MakeRequest(t, req, http.StatusOK)
}

func TestGetJoinGroupHelpInvalid(t *testing.T) {
	prepareTestEnv(t)
	// Invalid code => invalid format
	req := NewRequest(t, "GET", "1234567890fakefirebaseid0003", "/group/join/1234")
	MakeRequest(t, req, http.StatusBadRequest)
}

func TestJoinGroupThroughCodeInvalid(t *testing.T) {
	prepareTestEnv(t)
	req := NewRequest(t, "POST", "1234567890fakefirebaseid0001", "/group/join/123123123123")
	MakeRequest(t, req, http.StatusBadRequest)
}

func TestJoinGroupThroughCode(t *testing.T) {
	prepareTestEnv(t)
	var (
		userUID = "1234567890fakefirebaseid0001"
		code    = models.GroupCode{}
		reqCode = NewRequest(t, "GET", userUID, "/group/create-code")
		resp    = MakeRequest(t, reqCode, http.StatusOK)
	)
	DecodeJSON(t, resp, &code)
	// Join with generated code
	req := NewRequest(t, "POST", userUID, "/group/join/"+*code.Code)
	MakeRequest(t, req, http.StatusOK)
}

func TestLeaveGroup(t *testing.T) {
	prepareTestEnv(t)
	authInGroup := "1234567890fakefirebaseid0001"
	req := NewRequest(t, "POST", authInGroup, "/group/leave")
	MakeRequest(t, req, http.StatusOK)

	u := &models.User{}
	req = NewRequest(t, "GET", authInGroup, "/users/"+authInGroup)
	resp := MakeRequest(t, req, http.StatusOK)
	DecodeJSON(t, resp, u)
	assert.Empty(t, u.GroupUID)
}
