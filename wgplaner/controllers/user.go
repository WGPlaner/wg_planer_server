package controllers

import (
	"fmt"
	"log"
	"net/http"
	"net/mail"
	"strings"
	"time"

	"github.com/acoshift/go-firebase-admin"
	"github.com/go-openapi/errors"
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

func UserIDAuth(token string) (interface{}, error) {
	theUser := models.User{UID: token}
	if isRegistered, err := isUserRegistered(&theUser); err != nil {
		log.Println("[Controller][User][Auth] Error with isUserRegistered ", err.Error())
		return nil, errors.New(http.StatusInternalServerError, "Internal Server Error")

	} else if !isRegistered {
		return nil, errors.Unauthenticated("invalid credentials")
	}
	return theUser, nil
}

// True, if the user is registered in the database. The users' data will be written
// to `theUser`
func isUserRegistered(theUser *models.User) (bool, error) {
	if isRegistered, err := wgplaner.OrmEngine.Get(theUser); err != nil {
		return false, err
	} else {
		return isRegistered, nil
	}
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

func CreateUser(params user.CreateUserParams) middleware.Responder {
	log.Println("[User][POST] Creating User")

	uid := strings.TrimSpace(params.Body.UID)
	theUser := models.User{UID: uid}

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

	theUser = models.User{
		UID:         uid,
		DisplayName: &displayName,
		Email:       params.Body.Email,
		GroupUID:    params.Body.GroupUID,
		PhotoURL:    "TODO",
		CreatedAt:   creationTime,
		UpdatedAt:   creationTime,
	}

	// Validate user
	if isValid, err := validateUser(&theUser); !isValid {
		log.Println("[User][POST] Error validating user!", err)
		return user.NewCreateUserBadRequest().WithPayload(&models.ErrorResponse{
			Message: swag.String(fmt.Sprintf("Invalid userId: \"%s\"", err.Error())),
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
