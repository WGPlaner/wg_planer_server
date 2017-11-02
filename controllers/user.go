package controllers

import (
	"github.com/acoshift/go-firebase-admin"
	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/swag"
	"github.com/wgplaner/wg_planer_server/gen/models"
	"github.com/wgplaner/wg_planer_server/gen/restapi/operations/user"
	"log"
	"time"
)

func CreateUser(params user.CreateUserParams) middleware.Responder {
	log.Println("[Controller][User] Creating User")

	displayName := swag.StringValue(params.Body.DisplayName)
	creationTimeStr := time.Now().String()

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
		CreatedAt:   creationTimeStr,
		UpdatedAt:   creationTimeStr,
	}

	return user.NewCreateUserOK().WithPayload(&theUser)
}

func GetUser(params user.GetUserParams) middleware.Responder {
	// Firebase Auth
	myUser, err := fireBaseAuth.GetUser(params.HTTPRequest.Context(), params.UserID)

	// TODO: Validate user id (length, etc)

	if err == firebase.ErrUserNotFound || myUser == nil {
		log.Printf("[Controller][User][GET] Can't find user with id \"%s\"!", params.UserID)
		return user.NewGetUserNotFound().WithPayload(&models.ErrorResponse{
			Message: swag.String("User not found!"),
			Status:  swag.Int64(404),
		})
	} else if err != nil {
		log.Printf("[Controller][User][GET] Firebase SDK Error!", err)
		return user.NewGetUserDefault(500).WithPayload(&models.ErrorResponse{
			Message: swag.String("Internal Server error!"),
			Status:  swag.Int64(500),
		})
	}

	theUser := models.User{
		UID:         myUser.UserID,
		DisplayName: &myUser.DisplayName,
		Email:       myUser.Email,
		GroupUID:    "<TODO>",
		PhotoURL:    myUser.PhotoURL, // "https://api.wggplaner.ameyering.de/users/$userId
		CreatedAt:   myUser.Metadata.CreatedAt.String(),
		UpdatedAt:   myUser.Metadata.LastSignedInAt.String(),
	}

	return user.NewGetUserOK().WithPayload(&theUser)
}
