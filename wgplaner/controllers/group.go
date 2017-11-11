package controllers

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
	"github.com/satori/go.uuid"
	"github.com/wgplaner/wg_planer_server/gen/models"
	"github.com/wgplaner/wg_planer_server/gen/restapi/operations/group"
	"github.com/wgplaner/wg_planer_server/wgplaner"
)

const (
	GROUP_CODE_LENGTH     = 9
	GROUP_CODE_VALID_DAYS = 3
)

func validateGroup(_ *models.Group) (bool, error) {
	// TODO
	return true, nil
}

func GetGroup(params group.GetGroupParams, principal interface{}) middleware.Responder {
	theGroup := models.Group{UID: strfmt.UUID(params.GroupID)}

	// TODO: Validate
	// validateGroup(&theGroup)

	// Database
	if isRegistered, err := wgplaner.OrmEngine.Get(&theGroup); err != nil {
		log.Println("[Group][GET] Database Error!", err)
		return userInternalServerError
	} else if !isRegistered {
		log.Printf("[Group][GET] Can't find databse group with id \"%s\"!", theGroup.UID)
		return group.NewGetGroupNotFound().WithPayload(&models.ErrorResponse{
			Message: swag.String("Group not found on server"),
			Status:  swag.Int64(http.StatusNotFound),
		})
	}

	return group.NewGetGroupOK().WithPayload(&theGroup)
}

func CreateGroupCode(params group.CreateGroupCodeParams, principal interface{}) middleware.Responder {
	log.Println("[Group Code][GET] Generate group code!")

	// TODO: Check authorization for user in the group

	groupUid := strfmt.UUID(params.GroupID)
	code := wgplaner.RandomAlphaNumCode(GROUP_CODE_LENGTH, false)
	validDateTime := strfmt.DateTime(
		time.Now().UTC().AddDate(0, 0, GROUP_CODE_VALID_DAYS),
	)

	groupCode := models.GroupCode{
		GroupUID:   &groupUid,
		Code:       &code,
		ValidUntil: validDateTime,
	}

	// TODO: Delete old codes

	// Insert new code into database
	if _, err := wgplaner.OrmEngine.InsertOne(&groupCode); err != nil {
		log.Println("[Group Code][GET] Database error!", err)
		return userInternalServerError
	}

	return group.NewCreateGroupCodeOK().WithPayload(&groupCode)
}

func CreateGroup(params group.CreateGroupParams, principal interface{}) middleware.Responder {
	log.Println("[Group][POST] Creating group")

	theGroup := models.Group{}

	// Create new group
	displayName := strings.TrimSpace(swag.StringValue(params.Body.DisplayName))
	creationTime := strfmt.DateTime(time.Now().UTC())
	currency := strings.TrimSpace(params.Body.Currency)

	if currency == "" {
		currency = "€"
	}

	theGroup = models.Group{
		UID:         strfmt.UUID(uuid.NewV4().String()),
		Admins:      []string{*principal.(models.User).UID},
		Members:     []string{*principal.(models.User).UID},
		DisplayName: &displayName,
		Currency:    currency,
		CreatedAt:   creationTime,
		UpdatedAt:   creationTime,
	}

	// Validate group
	if isValid, err := validateGroup(&theGroup); !isValid {
		log.Println("[Group][POST] Error validating user!", err)
		return group.NewCreateGroupBadRequest().WithPayload(&models.ErrorResponse{
			Message: swag.String(fmt.Sprintf("Invalid group data: \"%s\"", err.Error())),
			Status:  swag.Int64(400),
		})
	}

	// TODO: Check if user has already a group

	// Insert new user into database
	if _, err := wgplaner.OrmEngine.InsertOne(&theGroup); err != nil {
		log.Println("[Group][POST] Database error!", err)
		return userInternalServerError
	}

	log.Println("[Group][POST] Created group")

	return group.NewCreateGroupOK().WithPayload(&theGroup)
}
