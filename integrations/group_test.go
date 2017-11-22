package integrations

import (
	"net/http"
	"testing"
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
