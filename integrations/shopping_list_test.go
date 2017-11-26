package integrations

import (
	"net/http"
	"testing"

	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
	"github.com/stretchr/testify/assert"
	"github.com/wgplaner/wg_planer_server/models"
)

func TestGetShoppinglistUnauthorized(t *testing.T) {
	prepareTestEnv(t)
	var (
		authInGroup = "1234567890fakefirebaseid0003"
		req         = NewRequest(t, "GET", authInGroup, "/shoppinglist/00112233-4455-6677-8899-aabbccddeeff")
	)
	MakeRequest(t, req, http.StatusUnauthorized)
}

func TestGetShoppinglist(t *testing.T) {
	prepareTestEnv(t)
	var (
		shopList    models.ShoppingList
		authInGroup = "1234567890fakefirebaseid0001"
		req         = NewRequest(t, "GET", authInGroup, "/shoppinglist/00112233-4455-6677-8899-aabbccddeeff")
		resp        = MakeRequest(t, req, http.StatusOK)
	)
	DecodeJSON(t, resp, &shopList)

	assert.Len(t, shopList.ListItems, 3)
	assert.Equal(t, shopList.Count, int64(3))
}

func TestCreateListItemUnauthorized(t *testing.T) {
	prepareTestEnv(t)
	var (
		authInGroup = "1234567890fakefirebaseid0003"
		item        = models.ListItem{
			Title:        swag.String("Eggs"),
			Category:     swag.String("Groceries"),
			Count:        swag.Int64(1),
			RequestedFor: []string{authInGroup},
		}
		req = NewRequestWithJSON(t, "POST", authInGroup,
			"/shoppinglist/00112233-4455-6677-8899-aabbccddeeff", item)
	)
	MakeRequest(t, req, http.StatusUnauthorized)
}

func TestCreateListItemInvalid(t *testing.T) {
	prepareTestEnv(t)
	var (
		item = models.ListItem{Title: swag.String("Eggs")}
		req  = NewRequestWithJSON(t, "POST", "1234567890fakefirebaseid0003",
			"/shoppinglist/00112233-4455-6677-8899-aabbccddeeff", item)
	)
	MakeRequest(t, req, http.StatusUnprocessableEntity)
}

func TestCreateListItem(t *testing.T) {
	prepareTestEnv(t)
	var (
		authInGroup = "1234567890fakefirebaseid0001"
		groupUID    = "00112233-4455-6677-8899-aabbccddeeff"
		item        = models.ListItem{
			Title:        swag.String("Eggs"),
			Category:     swag.String("Groceries"),
			Count:        swag.Int64(1),
			RequestedFor: []string{authInGroup},
		}
		req = NewRequestWithJSON(t, "POST", authInGroup,
			"/shoppinglist/"+groupUID, item)
		resp = MakeRequest(t, req, http.StatusOK)
	)
	// Check that the item was created.
	var shopList = models.ShoppingList{}
	req = NewRequest(t, "GET", authInGroup, "/shoppinglist/"+groupUID)
	resp = MakeRequest(t, req, http.StatusOK)
	DecodeJSON(t, resp, &shopList)
	assert.Len(t, shopList.ListItems, 4)
	assert.Equal(t, shopList.Count, int64(4))
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
			"/shoppinglist/"+groupUID, item)
		resp = MakeRequest(t, req, http.StatusOK)
	)
	DecodeJSON(t, resp, &uItem)
	assert.Equal(t, "New Milk", *uItem.Title)
	assert.Equal(t, "New Groceries", *uItem.Category)
	assert.Equal(t, int64(0), uItem.Price)
	assert.Equal(t, int64(2), *uItem.Count)
	assert.NotEqual(t, uItem.CreatedAt, uItem.UpdatedAt)
}

func TestUpdateListItemUnauthorized(t *testing.T) {
	prepareTestEnv(t)
	var (
		authInGroup = "1234567890fakefirebaseid0003"
		groupUID    = "00112233-4455-6677-8899-aabbccddeeff"
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
			"/shoppinglist/"+groupUID, item)
	)
	MakeRequest(t, req, http.StatusUnauthorized)
}

func TestUpdateListItemInvalid(t *testing.T) {
	prepareTestEnv(t)
	var (
		item = models.ListItem{Title: swag.String("Eggs")}
		req  = NewRequestWithJSON(t, "PUT", "1234567890fakefirebaseid0003",
			"/shoppinglist/00112233-4455-6677-8899-aabbccddeeff", item)
	)
	MakeRequest(t, req, http.StatusUnprocessableEntity)
}

func TestBuyListItems(t *testing.T) {
	prepareTestEnv(t)
	var (
		items = []string{"00112233-4455-6677-8899-000000000002", "00112233-4455-6677-8899-000000000003"}
		req   = NewRequestWithJSON(t, "POST", "1234567890fakefirebaseid0002",
			"/shoppinglist/00112233-4455-6677-8899-aabbccddeeff/buy-items", items)
	)
	MakeRequest(t, req, http.StatusOK)
}
