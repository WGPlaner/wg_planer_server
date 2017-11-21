package controllers

import (
	"fmt"
	"time"

	"github.com/wgplaner/wg_planer_server/gen/models"
	"github.com/wgplaner/wg_planer_server/gen/restapi/operations/shoppinglist"
	"github.com/wgplaner/wg_planer_server/wgplaner"

	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/strfmt"
	"github.com/op/go-logging"
	"github.com/satori/go.uuid"
)

var shoppingLog = logging.MustGetLogger("Shop")

func GetListItems(params shoppinglist.GetListItemsParams, principal interface{}) middleware.Responder {
	var items []*models.ListItem
	var groupUid = strfmt.UUID(params.GroupUID)

	if err := wgplaner.OrmEngine.Find(&items, &models.ListItem{GroupUID: groupUid}); err != nil {
		shoppingLog.Criticalf(`Database error finding list items for group "%s"`, groupUid)
		return NewInternalServerError("Database Error")
	}

	// TODO: Check authorization, add filters (limit), etc.

	shoppingList := models.ShoppingList{
		Count:     int64(len(items)),
		ListItems: items,
	}

	return shoppinglist.NewGetListItemsOK().WithPayload(&shoppingList)
}

func UpdateListItem(params shoppinglist.UpdateListItemParams, principal interface{}) middleware.Responder {
	var (
		theUser   = principal.(models.User)
		groupUid  = strfmt.UUID(params.GroupUID)
		createdAt = time.Now().UTC()
	)
	shoppingLog.Debugf(`Updating shopping list item. User "%s" for group "%s"`,
		*theUser.UID, groupUid)

	if err := validateGroupUuid(groupUid); err != nil {
		shoppingLog.Debugf(`Error validating group "%s": "%s"`, groupUid, err.Error())
		return NewBadRequest(err.Error())
	}

	if len(params.Body.RequestedFor) == 0 {
		return NewBadRequest("RequestedFor must contain at least one user")
	}

	// TODO: Check if user is unique

	for _, userId := range params.Body.RequestedFor {
		if exists, err := wgplaner.OrmEngine.Exist(&models.User{UID: &userId}); err != nil {
			shoppingLog.Debugf(`Database error checking existence for user`, err.Error())
			return NewInternalServerError(err.Error())

		} else if !exists {
			shoppingLog.Debugf(`User in RequestedFor does not exist: "%s"`, userId)
			return NewBadRequest(
				fmt.Sprintf(`User "%s" in RequestedFor does not exist`, userId),
			)
		}
	}

	// TODO: This is ugly.
	listItem := models.ListItem{
		ID:           params.Body.ID,
		Title:        params.Body.Title,
		Category:     params.Body.Category,
		Count:        params.Body.Count,
		Price:        params.Body.Price,
		RequestedFor: params.Body.RequestedFor,
		RequestedBy:  *theUser.UID,
		GroupUID:     groupUid,
		BoughtAt:     strfmt.DateTime(time.Time{}), // Not bought, yet
		CreatedAt:    strfmt.DateTime(createdAt),
		UpdatedAt:    strfmt.DateTime(createdAt),
	}

	// Insert new code into database
	if _, err := wgplaner.OrmEngine.Update(&listItem); err != nil {
		shoppingLog.Critical("Database error updating list item!", err)
		return userInternalServerError
	}

	return shoppinglist.NewUpdateListItemOK().WithPayload(&listItem)
}

func CreateListItem(params shoppinglist.CreateListItemParams, principal interface{}) middleware.Responder {
	var (
		theUser   = principal.(models.User)
		groupUid  = strfmt.UUID(params.GroupUID)
		createdAt = time.Now().UTC()
	)
	shoppingLog.Debugf(`Creating shopping list item. User "%s" for group "%s"`,
		*theUser.UID, groupUid)

	if err := validateGroupUuid(groupUid); err != nil {
		shoppingLog.Debugf(`Error validating group "%s": "%s"`, groupUid, err.Error())
		return NewBadRequest(err.Error())
	}

	if len(params.Body.RequestedFor) == 0 {
		return NewBadRequest("RequestedFor must contain at least one user")
	}

	// TODO: Check if user is unique

	for _, userId := range params.Body.RequestedFor {
		if exists, err := wgplaner.OrmEngine.Exist(&models.User{UID: &userId}); err != nil {
			shoppingLog.Debugf(`Database error checking existence for user`, err.Error())
			return NewInternalServerError(err.Error())

		} else if !exists {
			shoppingLog.Debugf(`User in RequestedFor does not exist: "%s"`, userId)
			return NewBadRequest(
				fmt.Sprintf(`User "%s" in RequestedFor does not exist`, userId),
			)
		}
	}

	listItem := models.ListItem{
		ID:           strfmt.UUID(uuid.NewV4().String()),
		Title:        params.Body.Title,
		Category:     params.Body.Category,
		Count:        params.Body.Count,
		Price:        params.Body.Price,
		RequestedFor: params.Body.RequestedFor,
		RequestedBy:  *theUser.UID,
		GroupUID:     groupUid,
		BoughtAt:     strfmt.DateTime(time.Time{}), // Not bought, yet
		CreatedAt:    strfmt.DateTime(createdAt),
		UpdatedAt:    strfmt.DateTime(createdAt),
	}

	// Insert new code into database
	if _, err := wgplaner.OrmEngine.InsertOne(&listItem); err != nil {
		shoppingLog.Critical("Database error inserting list item!", err)
		return userInternalServerError
	}

	return shoppinglist.NewCreateListItemOK().WithPayload(&listItem)
}
