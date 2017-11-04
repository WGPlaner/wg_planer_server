package controllers

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/acoshift/go-firebase-admin"
	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
	"github.com/wgplaner/wg_planer_server/gen/models"
	"github.com/wgplaner/wg_planer_server/gen/restapi/operations/user"
	"github.com/wgplaner/wg_planer_server/wgplaner"
)

var userInternalServerError = user.NewGetUserDefault(500).WithPayload(&models.ErrorResponse{
	Message: swag.String("Internal Server error!"),
	Status:  swag.Int64(500),
})

// TODO: Validate user id (length, etc)
func validateUser(theUser *models.User) (bool, error) {
	//return theUser.Validate()
	return true, nil
}

func CreateUser(params user.CreateUserParams) middleware.Responder {
	log.Println("[Controller][User] Creating User")

	displayName := swag.StringValue(params.Body.DisplayName)
	creationTime := strfmt.DateTime(time.Now().UTC())

	// TODO: Trim strings
	if displayName == "" {
		displayName = "John Doe"
	}

	// TODO: Remove Example Data
	theUser := models.User{
		UID:         params.Body.UID,
		DisplayName: &displayName,
		Email:       "test@example.com",
		GroupUID:    "",
		PhotoURL:    "",
		CreatedAt:   creationTime,
		UpdatedAt:   creationTime,
	}

	return user.NewCreateUserOK().WithPayload(&theUser)
}

func GetUser(params user.GetUserParams) middleware.Responder {
	theUser := models.User{UID: params.UserID}

	if isValid, err := validateUser(&theUser); !isValid {
		return user.NewGetUserBadRequest().WithPayload(&models.ErrorResponse{
			Message: swag.String(fmt.Sprintf("Invalid userId: \"%s\"", err.Error())),
			Status:  swag.Int64(400),
		})
	}

	// Firebase Auth
	_, err := wgplaner.FireBaseApp.Auth().
		GetUser(params.HTTPRequest.Context(), theUser.UID)

	if err == firebase.ErrUserNotFound {
		log.Printf("[User][GET] Can't find firebase user with id \"%s\"!", params.UserID)
		return user.NewGetUserUnauthorized().WithPayload(&models.ErrorResponse{
			Message: swag.String("User not authorized!"),
			Status:  swag.Int64(401),
		})
	} else if err != nil {
		log.Println("[User][GET] Firebase SDK Error!", err)
		return userInternalServerError
	}

	// Database
	if isRegistered, err := wgplaner.OrmEngine.Get(&theUser); err != nil {
		log.Println("[User][GET] Database Error!", err)
		return userInternalServerError
	} else if !isRegistered {
		log.Printf("[User][GET] Can't find databse user with id \"%s\"!", params.UserID)
		return user.NewGetUserNotFound().WithPayload(&models.ErrorResponse{
			Message: swag.String("User not found on server"),
			Status:  swag.Int64(404),
		})
	}

	return user.NewGetUserOK().WithPayload(&theUser)
}
