package integrations

import (
	"net/http"
	"testing"
)

func TestUserIDAuth(t *testing.T) {
	prepareTestEnv(t)
	// This required a registered user
	req1 := NewRequest(t, "GET", "1234567890fakefirebaseid0001",
		"/users/1234567890fakefirebaseid0001/image")
	MakeRequest(t, req1, http.StatusOK)

	req2 := NewRequest(t, "GET", "0000000000fakefirebaseid0000",
		"/users/1234567890fakefirebaseid0001/image")
	MakeRequest(t, req2, http.StatusUnauthorized)
}
