package integrations

import (
	"net/http"
	"testing"

	"github.com/go-openapi/swag"
	"github.com/stretchr/testify/assert"
	"github.com/wgplaner/wg_planer_server/gen/models"
	"github.com/wgplaner/wg_planer_server/gen/restapi/operations/user"
)

func TestGetUserUnauthorized(t *testing.T) {
	prepareTestEnv(t)
	req := NewRequest(t, "GET", AuthInvalid, "/users/1234567890fakefirebaseid0001")
	MakeRequest(t, req, http.StatusUnauthorized)
}

func TestGetOwnUser(t *testing.T) {
	prepareTestEnv(t)
	var (
		userId  = "1234567890fakefirebaseid0001"
		req     = NewRequest(t, "GET", userId, "/users/"+userId)
		resp    = MakeRequest(t, req, http.StatusOK)
		apiUser *models.User
	)

	if DecodeJSON(t, resp, &apiUser) {
		assert.Equal(t, userId, *apiUser.UID)
		assert.Equal(t, "John Doe", *apiUser.DisplayName)
	}
}

func TestGetOtherUser(t *testing.T) {
	prepareTestEnv(t)
	var (
		ownUserId = "1234567890fakefirebaseid0001"
		userId    = "1234567890fakefirebaseid0002"
		req       = NewRequest(t, "GET", ownUserId, "/users/"+userId)
		resp      = MakeRequest(t, req, http.StatusOK)
		apiUser   *models.User
	)

	if DecodeJSON(t, resp, &apiUser) {
		assert.Equal(t, userId, *apiUser.UID)
		assert.Equal(t, "Max Meier", *apiUser.DisplayName)
	}
}

func TestGetUserImage(t *testing.T) {
	prepareTestEnv(t)
	var (
		authId      = "1234567890fakefirebaseid0001"
		userId      = "1234567890fakefirebaseid0002"
		urlBuilder  = user.GetUserImageURL{UserID: userId}
		imgURL, err = urlBuilder.Build()
	)

	if !assert.NoError(t, err) {
		return
	}

	var (
		req  = NewRequest(t, "GET", authId, imgURL.String())
		resp = MakeRequest(t, req, http.StatusOK)
	)

	assert.Equal(t, resp.Headers.Get("Content-Type"), "application/octet-stream")
}

func TestCreateUserUnauthorized(t *testing.T) {
	prepareTestEnv(t)
	var (
		errResp = models.ErrorResponse{}
		newUser = models.User{
			UID:         swag.String("1234567890fakefirebaseid0003"),
			DisplayName: swag.String("Andre"),
		}
		req  = NewRequestWithJSON(t, "POST", "1234567890fakefirebaseid0001", "/users", newUser)
		resp = MakeRequest(t, req, http.StatusUnauthorized)
	)

	if DecodeJSON(t, resp, &errResp) {
		assert.NotEmpty(t, *errResp.Message)
	}
}

func TestCreateUserMissingRequired(t *testing.T) {
	prepareTestEnv(t)
	var (
		errResp = models.ErrorResponse{}
		req     = NewRequestWithJSON(t, "POST", AuthValid, "/users", models.User{})
		resp    = MakeRequest(t, req, http.StatusUnprocessableEntity)
	)

	if DecodeJSON(t, resp, &errResp) {
		assert.NotEmpty(t, *errResp.Message)
	}
}

func TestCreateUser(t *testing.T) {
	prepareTestEnv(t)
	var (
		uid     = "1234567890fakefirebaseid0003"
		newUser = models.User{
			UID:         &uid,
			DisplayName: swag.String("Andre"),
		}
		createdUser = models.User{}
		req         = NewRequestWithJSON(t, "POST", uid, "/users", newUser)
		resp        = MakeRequest(t, req, http.StatusOK)
	)

	if DecodeJSON(t, resp, &createdUser) {
		assert.Equal(t, *createdUser.UID, *newUser.UID)
		assert.Equal(t, *createdUser.DisplayName, *newUser.DisplayName)
	}

	// TODO: Load Beans
}