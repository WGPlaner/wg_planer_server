package controllers

import (
	"time"

	"github.com/go-openapi/strfmt"
	"github.com/satori/go.uuid"
	"github.com/wgplaner/wg_planer_server/models"
	"github.com/wgplaner/wg_planer_server/restapi/operations/shoppinglist"

	"github.com/go-openapi/runtime/middleware"
	"github.com/op/go-logging"
)

var shoppingLog = logging.MustGetLogger("Shop")

func GetListItems(params shoppinglist.GetListItemsParams, principal *models.User) middleware.Responder {
	var (
		err   error
		g     *models.Group
		items []*models.ListItem
	)

	if g, err = models.GetGroupByUID(params.GroupUID); models.IsErrGroupNotExist(err) {
		return NewNotFoundResponse("Group not found")
	} else if err != nil {
		return NewInternalServerError("Internal Server Error")
	}

	if !g.HasMember(*principal.UID) {
		return NewUnauthorizedResponse("Not member of group")
	}

	if items, err = g.GetActiveShoppingListItems(); err != nil {
		shoppingLog.Criticalf(`Database error finding list items for group "%s"`, g.UID)
		return NewInternalServerError("Database Error")
	}

	// TODO: Add filters (limit), etc.

	return shoppinglist.NewGetListItemsOK().WithPayload(&models.ShoppingList{
		Count:     int64(len(items)),
		ListItems: items,
	})
}

func UpdateListItem(params shoppinglist.UpdateListItemParams, principal *models.User) middleware.Responder {
	shoppingLog.Debugf(`Updating shopping list item. User "%s" for group "%s"`,
		*principal.UID, params.GroupUID)

	var (
		err error
		g   *models.Group
	)

	if g, err = models.GetGroupByUID(params.GroupUID); models.IsErrGroupNotExist(err) {
		return NewNotFoundResponse("Group not found")

	} else if err != nil {
		shoppingLog.Debugf(`Error validating group "%s": "%s"`, params.GroupUID, err.Error())
		return NewBadRequest(err.Error())
	}

	if !strfmt.IsUUID(string(params.Body.ID)) {
		return NewBadRequest("Invalid item ID")
	}

	if len(params.Body.RequestedFor) == 0 {
		return NewBadRequest("RequestedFor must contain at least one user")
	}

	// TODO: Check if user is unique
	if exists, err := models.AreUsersExist(params.Body.RequestedFor); err != nil {
		return NewInternalServerError(err.Error())

	} else if !exists {
		return NewBadRequest("A requestedFor user does not exist")
	}

	// TODO: This is ugly.
	listItem := &models.ListItem{
		ID:           params.Body.ID,
		GroupUID:     g.UID,
		Title:        params.Body.Title,
		Category:     params.Body.Category,
		Count:        params.Body.Count,
		Price:        params.Body.Price,
		RequestedFor: params.Body.RequestedFor,
	}

	// TODO: Check that the item exists

	// Insert new code into database
	if err := models.UpdateListItemCols(listItem, `title`, `category`, `count`, `price`, `requested_for`); err != nil {
		shoppingLog.Critical("Database error updating list item!", err)
		return NewInternalServerError("Internal Database Error")
	}

	// Get list item with its data
	if listItem, err = models.GetListItemByUIDs(listItem.GroupUID, listItem.ID); err != nil {
		shoppingLog.Critical("Database error querying list item!", err)
		return NewInternalServerError("Internal Database Error")
	}

	return shoppinglist.NewUpdateListItemOK().WithPayload(listItem)
}

func CreateListItem(params shoppinglist.CreateListItemParams, principal *models.User) middleware.Responder {
	shoppingLog.Debugf(`Creating shopping list item. User "%s" for group "%s"`,
		*principal.UID, params.GroupUID)

	var (
		err error
		g   *models.Group
	)

	if g, err = models.GetGroupByUID(params.GroupUID); models.IsErrGroupNotExist(err) {
		return NewNotFoundResponse("Group not found")

	} else if err != nil {
		shoppingLog.Debugf(`Error validating group "%s": "%s"`, params.GroupUID, err.Error())
		return NewBadRequest(err.Error())
	}

	if len(params.Body.RequestedFor) == 0 {
		return NewBadRequest("RequestedFor must contain at least one user")
	}

	// TODO: Check if user is unique
	if exists, err := models.AreUsersExist(params.Body.RequestedFor); err != nil {
		return NewInternalServerError(err.Error())

	} else if !exists {
		return NewBadRequest("A requestedFor user does not exist")
	}

	listItem := models.ListItem{
		ID:           strfmt.UUID(uuid.NewV4().String()),
		Title:        params.Body.Title,
		Category:     params.Body.Category,
		Count:        params.Body.Count,
		Price:        params.Body.Price,
		RequestedFor: params.Body.RequestedFor,
		RequestedBy:  *principal.UID,
		GroupUID:     g.UID,
		BoughtAt:     strfmt.DateTime(time.Time{}), // Not bought, yet
	}

	// Insert new code into database
	if err := models.CreateListItem(&listItem); err != nil {
		shoppingLog.Critical("Database error inserting list item!", err)
		return NewInternalServerError("Internal Database Error")
	}

	return shoppinglist.NewCreateListItemOK().WithPayload(&listItem)
}