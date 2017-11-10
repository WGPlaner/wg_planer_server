package controllers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/mail"
	"os"
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

// True, if the user is registered in the database. The users' data will be written
// to `theUser`
func isUserRegistered(theUser *models.User) (bool, error) {
	if isRegistered, err := wgplaner.OrmEngine.Get(theUser); err != nil {
		return false, err
	} else {
		return isRegistered, nil
	}
}

func isUserOnFirebase(theUser *models.User) (bool, error) {
	_, err := wgplaner.FireBaseApp.Auth().GetUser(context.Background(), *theUser.UID)

	if err == firebase.ErrUserNotFound {
		log.Printf("Can't find firebase user with id \"%s\"!", *theUser.UID)
		return false, nil
	} else if err != nil {
		log.Println("Firebase SDK Error!", err)
		return false, err
	}

	return true, nil
}

func isValidUserID(id *string) bool {
	return len(*id) == 28 // Firebase IDs are 28 characters long
}

// TODO: Validate user id (length, etc)
func validateUser(theUser *models.User) (bool, error) {
	//return theUser.Validate()
	if theUser.Email != "" {
		if _, err := mail.ParseAddress(string(theUser.Email)); err != nil {
			return false, err
		}
	}
	return true, nil
}

func CreateUser(params user.CreateUserParams, principal interface{}) middleware.Responder {
	log.Println("[User][POST] Creating User")

	theUser := models.User{
		UID: params.Body.UID,
	}

	// Check if the user is already registered
	if isRegistered, err := wgplaner.OrmEngine.Get(&theUser); err != nil {
		log.Println("[User][POST] Database Error!", err)
		return userInternalServerError

	} else if isRegistered {
		log.Println("[User][POST] User already exists!")
		return user.NewCreateUserOK().WithPayload(&theUser)
	}

	// Create new user
	displayName := strings.TrimSpace(swag.StringValue(params.Body.DisplayName))
	creationTime := strfmt.DateTime(time.Now().UTC())
	imageURL := user.GetUserImageURL{UserID: swag.StringValue(theUser.UID)}
	photoURL, err := imageURL.Build()
	if err != nil {
		return userInternalServerError
	}

	theUser = models.User{
		UID:         theUser.UID,
		DisplayName: &displayName,
		Email:       params.Body.Email,
		GroupUID:    params.Body.GroupUID,
		PhotoURL:    strfmt.URI(photoURL.String()),
		CreatedAt:   creationTime,
		UpdatedAt:   creationTime,
	}

	// Validate user
	if isValid, err := validateUser(&theUser); !isValid {
		log.Println("[User][POST] Error validating user!", err)
		return user.NewCreateUserBadRequest().WithPayload(&models.ErrorResponse{
			Message: swag.String(fmt.Sprintf("invalid user: \"%s\"", err.Error())),
			Status:  swag.Int64(400),
		})
	}

	// Insert new user into database
	if _, err := wgplaner.OrmEngine.InsertOne(&theUser); err != nil {
		log.Println("[User][POST] Database error!", err)
		return userInternalServerError
	}

	return user.NewCreateUserOK().WithPayload(&theUser)
}

func UpdateUser(params user.UpdateUserParams, principal interface{}) middleware.Responder {
	log.Println("[User][PUT] Creating User")

	theUser := models.User{UID: params.Body.UID}

	// Check if the user is already registered
	if isRegistered, err := wgplaner.OrmEngine.Get(&theUser); err != nil {
		log.Println("[User][PUT] Database Error!", err)
		return userInternalServerError

	} else if !isRegistered {
		log.Println("[User][PUT] User does not exist!")
		return user.NewUpdateUserDefault(400).WithPayload(&models.ErrorResponse{
			Message: swag.String("User does not exist!"),
			Status:  swag.Int64(400),
		})
	}

	// Create new user
	displayName := strings.TrimSpace(swag.StringValue(params.Body.DisplayName))
	creationTime := strfmt.DateTime(time.Now().UTC())
	imageURL := user.GetUserImageURL{UserID: swag.StringValue(theUser.UID)}
	photoURL, err := imageURL.Build()
	if err != nil {
		return userInternalServerError
	}

	theUser = models.User{
		UID:         params.Body.UID,
		DisplayName: &displayName,
		Email:       params.Body.Email,
		GroupUID:    params.Body.GroupUID,
		PhotoURL:    strfmt.URI(photoURL.String()),
		CreatedAt:   creationTime,
		UpdatedAt:   creationTime,
	}

	// Validate user
	if isValid, err := validateUser(&theUser); !isValid {
		log.Println("[User][PUT] Error validating user!", err)
		return user.NewUpdateUserBadRequest().WithPayload(&models.ErrorResponse{
			Message: swag.String(fmt.Sprintf("Invalid user: \"%s\"", err.Error())),
			Status:  swag.Int64(400),
		})
	}

	// Insert new user into database
	if _, err := wgplaner.OrmEngine.Update(&theUser); err != nil {
		log.Println("[User][PUT] Database error!", err)
		return userInternalServerError
	}

	return user.NewUpdateUserOK().WithPayload(&theUser)
}

func GetUser(params user.GetUserParams, principal interface{}) middleware.Responder {
	theUser := models.User{UID: &params.UserID}

	if isValid := isValidUserID(theUser.UID); !isValid {
		return user.NewGetUserBadRequest().WithPayload(&models.ErrorResponse{
			Message: swag.String(fmt.Sprintf("Invalid user id format")),
			Status:  swag.Int64(400),
		})
	}

	// Firebase Auth
	_, err := wgplaner.FireBaseApp.Auth().
		GetUser(params.HTTPRequest.Context(), *theUser.UID)

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
			Status:  swag.Int64(http.StatusNotFound),
		})
	}

	imageURL := user.GetUserImageURL{UserID: swag.StringValue(theUser.UID)}
	photoURL, err := imageURL.Build()
	if err != nil {
		return userInternalServerError
	}

	theUser.PhotoURL = strfmt.URI(photoURL.String())

	return user.NewGetUserOK().WithPayload(&theUser)
}

func GetUserImage(params user.GetUserImageParams, principal interface{}) middleware.Responder {
	log.Println("[Controller][User] Get User Image")

	theUser := models.User{UID: &params.UserID}

	var imgFile *os.File
	var fileErr error

	// Get default image if normal one does no exist
	if imgFile, fileErr = wgplaner.GetUserProfileImage(&theUser); os.IsNotExist(fileErr) {
		imgFile, fileErr = wgplaner.GetUserProfileImageDefault()
	}

	if fileErr != nil {
		log.Println("[Controller][User] Error getting profile image ", fileErr.Error())
		return user.NewGetUserImageDefault(http.StatusInternalServerError).
			WithPayload(&models.ErrorResponse{
				Message: swag.String("Internal Server Error"),
				Status:  swag.Int64(http.StatusInternalServerError),
			})
	}

	return user.NewGetUserImageOK().WithPayload(imgFile)
}
