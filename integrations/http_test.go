package integrations

import (
	"net/http"
	"testing"

	"github.com/go-openapi/swag"
	"github.com/wgplaner/wg_planer_server/models"
)

func TestNotFoundResponse(t *testing.T) {
	prepareTestEnv(t)
	var (
		req      = NewRequest(t, "GET", AuthEmpty, "/unknown/path")
		resp     = MakeRequest(t, req, http.StatusNotFound)
		apiError = &models.ErrorResponse{}
	)
	if DecodeJSON(t, resp, apiError) {
		apiError.Status = swag.Int64(http.StatusNotFound)
	}
}

func TestGetVersionInfo(t *testing.T) {
	prepareTestEnv(t)
	req := NewRequest(t, "GET", AuthEmpty, "/latest-version")
	resp := MakeRequest(t, req, http.StatusOK)
	DecodeJSON(t, resp, &models.VersionInfo{})
}
