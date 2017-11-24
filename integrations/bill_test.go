package integrations

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/wgplaner/wg_planer_server/models"
)

func TestCreateBill(t *testing.T) {
	prepareTestEnv(t)
	authValid := "1234567890fakefirebaseid0001"
	req := NewRequest(t, "POST", authValid, "/group/00112233-4455-6677-8899-aabbccddeeff/bills/create")
	MakeRequest(t, req, http.StatusOK)
}

func TestGetBills(t *testing.T) {
	prepareTestEnv(t)
	var (
		billList  = &models.BillList{}
		authValid = "1234567890fakefirebaseid0001"
		req       = NewRequest(t, "GET", authValid, "/group/00112233-4455-6677-8899-aabbccddeeff/bills")
		resp      = MakeRequest(t, req, http.StatusOK)
	)
	DecodeJSON(t, resp, billList)

	assert.Len(t, billList.Bills, 1)
	assert.Equal(t, billList.Count, int64(1))
}
