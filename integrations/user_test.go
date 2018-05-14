package integrations

import (
	"net/http"
	"testing"

	"github.com/wgplaner/wg_planer_server/models"
	"github.com/wgplaner/wg_planer_server/restapi/operations/user"

	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
	"github.com/stretchr/testify/assert"
)

func TestGetUserUnauthorized(t *testing.T) {
	prepareTestEnv(t)
	req := NewRequest(t, "GET", AuthInvalid, "/users/1234567890fakefirebaseid0001")
	MakeRequest(t, req, http.StatusUnauthorized)
}

func TestGetOwnUser(t *testing.T) {
	prepareTestEnv(t)
	var (
		userID  = "1234567890fakefirebaseid0001"
		req     = NewRequest(t, "GET", userID, "/users/"+userID)
		resp    = MakeRequest(t, req, http.StatusOK)
		apiUser *models.User
	)

	if DecodeJSON(t, resp, &apiUser) {
		assert.Equal(t, userID, *apiUser.UID)
		assert.Equal(t, "John Doe", *apiUser.DisplayName)
	}
}

func TestGetOtherUser(t *testing.T) {
	prepareTestEnv(t)
	var (
		ownUserID = "1234567890fakefirebaseid0001"
		userID    = "1234567890fakefirebaseid0002"
		req       = NewRequest(t, "GET", ownUserID, "/users/"+userID)
		resp      = MakeRequest(t, req, http.StatusOK)
		apiUser   *models.User
	)

	if DecodeJSON(t, resp, &apiUser) {
		assert.Equal(t, userID, *apiUser.UID)
		assert.Equal(t, "Max Meier", *apiUser.DisplayName)
	}
}

func TestGetUserImage(t *testing.T) {
	prepareTestEnv(t)
	var (
		authID      = "1234567890fakefirebaseid0001"
		userID      = "1234567890fakefirebaseid0002"
		urlBuilder  = user.GetUserImageURL{UserID: userID}
		imgURL, err = urlBuilder.Build()
	)

	if !assert.NoError(t, err) {
		return
	}

	var (
		req  = NewRequest(t, "GET", authID, imgURL.String())
		resp = MakeRequest(t, req, http.StatusOK)
	)

	assert.Equal(t, resp.Headers.Get("Content-Type"), "application/octet-stream")
}

func TestGetUserBoughtItems(t *testing.T) {
	prepareTestEnv(t)
	var (
		authID = "1234567890fakefirebaseid0001"
		userID = "1234567890fakefirebaseid0002"
		req    = NewRequest(t, "GET", authID, "/users/"+userID+"/bought")
		resp   = MakeRequest(t, req, http.StatusOK)
	)

	boughtItems := models.ShoppingList{}
	DecodeJSON(t, resp, &boughtItems)
	assert.Equal(t, int64(1), boughtItems.Count)
	assert.Equal(t, int64(1510337621), boughtItems.ListItems[0].BoughtAt.Unix())
}

func TestUpdateUserImage(t *testing.T) {
	prepareTestEnv(t)

	request := NewRequestWithImage(t,
		"PUT",
		AuthValid,
		"/users/"+AuthValid+"/image",
		"profileImage",
		models.GetUserImageDefaultPath())
	resp := MakeRequest(t, request, http.StatusOK)

	// Check JSON response
	successResp := models.SuccessResponse{}
	DecodeJSON(t, resp, &successResp)
	assert.NotEmpty(t, *successResp.Message)
	assert.Equal(t, int64(200), *successResp.Status)

	// Check that image is uploaded
	_, fileErr := models.GetUserImage(AuthValid)
	assert.NoError(t, fileErr, "Uploaded image not found")
}

func TestCreateUserUnauthorized(t *testing.T) {
	prepareTestEnv(t)
	var (
		errResp = models.ErrorResponse{}
		newUser = models.User{
			UID:         swag.String("1234567890fakefirebaseid0010"),
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

func TestCreateExistingUser(t *testing.T) {
	prepareTestEnv(t)
	var (
		uid     = "1234567890fakefirebaseid0001"
		newUser = models.User{
			UID:         &uid,
			DisplayName: swag.String("Andre"),
		}
		req  = NewRequestWithJSON(t, "POST", uid, "/users", newUser)
		resp = MakeRequest(t, req, http.StatusBadRequest)
	)

	DecodeJSON(t, resp, &models.ErrorResponse{})
}

func TestCreateUser(t *testing.T) {
	prepareTestEnv(t)
	var (
		uid     = "1234567890fakefirebaseid0010"
		newUser = models.User{
			UID:         &uid,
			DisplayName: swag.String("Andre"),
			GroupUID:    strfmt.UUID("0ec972c9-6c7a-40c8-82c3-000000000000"), // Random UID
		}
		createdUser = models.User{}
		req         = NewRequestWithJSON(t, "POST", uid, "/users", newUser)
		resp        = MakeRequest(t, req, http.StatusOK)
	)

	if DecodeJSON(t, resp, &createdUser) {
		assert.Equal(t, *createdUser.UID, *newUser.UID)
		assert.Equal(t, *createdUser.DisplayName, *newUser.DisplayName)
		// Group should not be saved. Only through Group Code
		assert.Equal(t, createdUser.GroupUID, strfmt.UUID(""))
	}

	// TODO: Load Beans
}

func TestUpdateUserUnauthorized(t *testing.T) {
	prepareTestEnv(t)
	var (
		errResp = models.ErrorResponse{}
		newUser = models.User{
			UID:         swag.String("1234567890fakefirebaseid0002"),
			DisplayName: swag.String("Maxi Meier"),
		}
		req  = NewRequestWithJSON(t, "PUT", "1234567890fakefirebaseid0001", "/users", newUser)
		resp = MakeRequest(t, req, http.StatusUnauthorized)
	)

	DecodeJSON(t, resp, &errResp)
	assert.NotEmpty(t, *errResp.Message)
}

func TestUpdateUser(t *testing.T) {
	prepareTestEnv(t)
	var (
		uid         = "1234567890fakefirebaseid0002"
		oldGroupUID = strfmt.UUID("00112233-4455-6677-8899-aabbccddeeff")
		uUser       = models.User{
			UID:         &uid,
			DisplayName: swag.String("Maxi Meier"),
			GroupUID:    strfmt.UUID("0ec972c9-6c7a-40c8-82c3-000000000000"), // New uid
		}
		updatedUser = models.User{}
		req         = NewRequestWithJSON(t, "PUT", uid, "/users", uUser)
		resp        = MakeRequest(t, req, http.StatusOK)
	)

	if DecodeJSON(t, resp, &updatedUser) {
		assert.Equal(t, *uUser.UID, *updatedUser.UID)
		assert.Equal(t, *uUser.DisplayName, *updatedUser.DisplayName)
		// Group should not be updated. Only through Group Code
		assert.Equal(t,
			oldGroupUID,
			updatedUser.GroupUID,
		)
	}
}
