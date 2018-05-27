package integrations

import (
	"net/http"
	"testing"

	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
	"github.com/stretchr/testify/assert"
	"github.com/wgplaner/wg_planer_server/models"
)

func TestGetShoppinglist(t *testing.T) {
	prepareTestEnv(t)
	var (
		shopList    models.ShoppingList
		authInGroup = "1234567890fakefirebaseid0001"
		req         = NewRequest(t, "GET", authInGroup, "/shoppinglist")
		resp        = MakeRequest(t, req, http.StatusOK)
	)
	DecodeJSON(t, resp, &shopList)

	assert.Len(t, shopList.ListItems, 2)
	assert.Equal(t, shopList.Count, int64(2))
}

func TestCreateListItemInvalid(t *testing.T) {
	prepareTestEnv(t)
	var (
		item = models.ListItem{Title: swag.String("Eggs")}
		req  = NewRequestWithJSON(t, "POST", "1234567890fakefirebaseid0003",
			"/shoppinglist", item)
	)
	MakeRequest(t, req, http.StatusUnprocessableEntity)
}

func TestCreateListItem(t *testing.T) {
	prepareTestEnv(t)
	var (
		authInGroup = "1234567890fakefirebaseid0001"
		item        = models.ListItem{
			Title:        swag.String("Eggs"),
			Category:     swag.String("Groceries"),
			Count:        swag.Int64(1),
			RequestedFor: []string{authInGroup},
		}
		req = NewRequestWithJSON(t, "POST", authInGroup,
			"/shoppinglist", item)
		resp = MakeRequest(t, req, http.StatusOK)
	)
	// Check that the item was created.
	var shopList = models.ShoppingList{}
	req = NewRequest(t, "GET", authInGroup, "/shoppinglist")
	resp = MakeRequest(t, req, http.StatusOK)
	DecodeJSON(t, resp, &shopList)
	assert.Len(t, shopList.ListItems, 3)
	assert.Equal(t, int64(3), shopList.Count)
}

func TestUpdateListItem(t *testing.T) {
	prepareTestEnv(t)
	var (
		authInGroup = "1234567890fakefirebaseid0001"
		groupUID    = "00112233-4455-6677-8899-aabbccddeeff"
		uItem       = models.ListItem{}
		item        = models.ListItem{
			ID:           "00112233-4455-6677-8899-000000000001",
			GroupUID:     strfmt.UUID(groupUID),
			Title:        swag.String("New Milk"),
			Category:     swag.String("New Groceries"),
			Count:        swag.Int64(2),
			Price:        0,
			RequestedFor: []string{authInGroup},
		}
		req = NewRequestWithJSON(t, "PUT", authInGroup,
			"/shoppinglist", item)
		resp = MakeRequest(t, req, http.StatusOK)
	)
	DecodeJSON(t, resp, &uItem)
	assert.Equal(t, "New Milk", *uItem.Title)
	assert.Equal(t, "New Groceries", *uItem.Category)
	assert.Equal(t, int64(0), uItem.Price)
	assert.Equal(t, int64(2), *uItem.Count)
	assert.NotEqual(t, uItem.CreatedAt, uItem.UpdatedAt)
}

func TestUpdateListItemInvalid(t *testing.T) {
	prepareTestEnv(t)
	var (
		item = models.ListItem{Title: swag.String("Eggs")}
		req  = NewRequestWithJSON(t, "PUT", "1234567890fakefirebaseid0003",
			"/shoppinglist", item)
	)
	MakeRequest(t, req, http.StatusUnprocessableEntity)
}

func TestBuyListItems(t *testing.T) {
	prepareTestEnv(t)
	var (
		boughByID = "1234567890fakefirebaseid0002"
		items     = []string{"00112233-4455-6677-8899-000000000002", "00112233-4455-6677-8899-000000000003"}
		req       = NewRequestWithJSON(t, "POST", boughByID,
			"/shoppinglist/buy-items", items)
	)
	MakeRequest(t, req, http.StatusOK)
	// Check database
	listItem := models.AssertExistsAndLoadBean(t,
		&models.ListItem{ID: strfmt.UUID(items[0])}).(*models.ListItem)
	assert.Equal(t, boughByID, listItem.BoughtBy)
}

func TestBuyListItemsThatDoNotExist(t *testing.T) {
	prepareTestEnv(t)
	var (
		items = []string{"00112233-4455-6677-8899-ccbbaa000000"}
		req   = NewRequestWithJSON(t, "POST", "1234567890fakefirebaseid0002",
			"/shoppinglist/buy-items", items)
	)
	MakeRequest(t, req, http.StatusBadRequest)
}

func TestUnBuyListItems(t *testing.T) {
	prepareTestEnv(t)
	var (
		boughtByID = "1234567890fakefirebaseid0002"
		item       = "00112233-4455-6677-8899-000000000004"
		req        = NewRequestWithJSON(t, "POST", boughtByID,
			"/shoppinglist/revert-purchase", item)
	)
	MakeRequest(t, req, http.StatusOK)
	// Check database
	listItem := models.AssertExistsAndLoadBean(t,
		&models.ListItem{ID: strfmt.UUID(item)}).(*models.ListItem)
	assert.Equal(t, "", listItem.BoughtBy)
}

func TestUnBuyListItemsWithBill(t *testing.T) {
	prepareTestEnv(t)
	var (
		boughtByID = "1234567890fakefirebaseid0002"
		item       = "00112233-4455-6677-8899-000000000003"
		req        = NewRequestWithJSON(t, "POST", boughtByID,
			"/shoppinglist/revert-purchase", item)
	)
	MakeRequest(t, req, http.StatusBadRequest)
	// Check database
	listItem := models.AssertExistsAndLoadBean(t,
		&models.ListItem{ID: strfmt.UUID(item)}).(*models.ListItem)
	assert.Equal(t, boughtByID, listItem.BoughtBy)
}
