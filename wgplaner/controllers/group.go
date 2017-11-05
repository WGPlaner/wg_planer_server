package controllers

import (
	"fmt"
	"log"
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

func validateGroup(_ *models.Group) (bool, error) {
	// TODO
	return true, nil
}

func CreateGroup(params group.CreateGroupParams, principal interface{}) middleware.Responder {
	log.Println("[Group][POST] Creating group")

	theGroup := models.Group{}

	// Create new group
	displayName := strings.TrimSpace(swag.StringValue(params.Body.DisplayName))
	creationTime := strfmt.DateTime(time.Now().UTC())

	theGroup = models.Group{
		UID:         uuid.NewV4().String(),
		Admins:      []string{principal.(models.User).UID},
		DisplayName: &displayName,
		Currency:    "€",
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

	return group.NewCreateGroupOK().WithPayload(&theGroup)
}
