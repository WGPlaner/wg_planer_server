package integrations

import (
	"net/http"
	"testing"

	"github.com/go-openapi/strfmt"
	"github.com/stretchr/testify/assert"
	"github.com/wgplaner/wg_planer_server/models"
)

func TestCreateBill(t *testing.T) {
	prepareTestEnv(t)
	var (
		authValid = "1234567890fakefirebaseid0002"
		bill      = models.Bill{}
		items     = []string{"00112233-4455-6677-8899-000000000004"}
		newBill   = models.Bill{BoughtItems: items, DueDate: "2019-06-07"}
		req       = NewRequestWithJSON(t, "POST", authValid, "/group/bills/create", newBill)
		resp      = MakeRequest(t, req, http.StatusOK)
	)
	DecodeJSON(t, resp, &bill)
	assert.Len(t, bill.BoughtItems, 1)
	assert.Equal(t, "00112233-4455-6677-8899-000000000004", bill.BoughtItems[0])

	// Assert that items exist in database
	item := models.AssertExistsAndLoadBean(t, &models.ListItem{ID: "00112233-4455-6677-8899-000000000004"}).(*models.ListItem)
	assert.Equal(t, bill.UID, item.BillUID)
	models.AssertCount(t, &models.Bill{}, 2)
}

func TestGetBills(t *testing.T) {
	prepareTestEnv(t)
	var (
		billList  = models.BillList{}
		authValid = "1234567890fakefirebaseid0001"
		req       = NewRequest(t, "GET", authValid, "/group/bills")
		resp      = MakeRequest(t, req, http.StatusOK)
	)
	DecodeJSON(t, resp, &billList)

	assert.Len(t, billList.Bills, 1)
	assert.Equal(t, int64(1), billList.Count)
	assert.Len(t, billList.Bills[0].BoughtItems, 2)
	assert.Len(t, billList.Bills[0].SentTo, 2)
	assert.Equal(t, "todo", *billList.Bills[0].State)
	assert.Equal(t, int64(270), billList.Bills[0].Sum)
	assert.Equal(t, strfmt.UUID("00112233-4455-6677-8899-aabbccddeeff"), billList.Bills[0].GroupUID)
	assert.Equal(t, strfmt.UUID("00112233-4455-6677-8899-000000000001"), billList.Bills[0].BoughtListItems[0].ID)
	assert.Equal(t, strfmt.UUID("00112233-4455-6677-8899-123000000001"), billList.Bills[0].BoughtListItems[0].BillUID)
}
