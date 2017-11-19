package integrations

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/wgplaner/wg_planer_server/gen/models"
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
