package controllers

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/mail"
	"os"
	"path"
	"strings"
	"time"

	"github.com/wgplaner/wg_planer_server/gen/models"
	"github.com/wgplaner/wg_planer_server/gen/restapi/operations/user"
	"github.com/wgplaner/wg_planer_server/wgplaner"

	"github.com/acoshift/go-firebase-admin"
	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
	"github.com/op/go-logging"
)

var userLog = logging.MustGetLogger("User")

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
		userLog.Debugf(`Can't find firebase user with id "%s"!`, *theUser.UID)
		return false, nil

	} else if err != nil {
		userLog.Critical("Firebase SDK Error!", err)
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
	userLog.Debugf(`Start creating User "%s"`, *params.Body.UID)

	theUser := models.User{
		UID: params.Body.UID,
	}

	if authUser, ok := principal.(models.User); !ok || *authUser.UID != *theUser.UID {
		userLog.Infof(`Authorized user "%s" tried to create account for "%s"`,
			*authUser.UID, *theUser.UID)
		return NewUnauthorizedResponse(`Can't create user for others.`)
	}

	// Check if the user is already registered
	if isRegistered, err := wgplaner.OrmEngine.Get(&theUser); err != nil {
		userLog.Critical("Database Error!", err)
		return userInternalServerError

	} else if isRegistered {
		userLog.Debug("User already exists!")
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
		CreatedAt:   creationTime,
		UpdatedAt:   creationTime,
	}

	// Validate user
	if isValid, err := validateUser(&theUser); !isValid {
		userLog.Debug("Error validating user!", err)
		return NewBadRequest(fmt.Sprintf(`invalid user: "%s"`, err.Error()))
	}

	// Insert new user into database
	if _, err := wgplaner.OrmEngine.InsertOne(&theUser); err != nil {
		userLog.Critical("Database error!", err)
		return userInternalServerError
	}

	theUser.PhotoURL = strfmt.URI(photoURL.String())

	userLog.Infof(`Created user "%s"`, *theUser.UID)

	return user.NewCreateUserOK().WithPayload(&theUser)
}

func UpdateUser(params user.UpdateUserParams, principal interface{}) middleware.Responder {
	userLog.Debug("Start updating user")

	theUser := models.User{UID: params.Body.UID}

	// Check if the user is already registered
	if isRegistered, err := wgplaner.OrmEngine.Get(&theUser); err != nil {
		userLog.Critical("Database Error!", err)
		return userInternalServerError

	} else if !isRegistered {
		userLog.Infof(`User "%s" does not exist!`, *theUser.UID)
		return NewBadRequest("User does not exist!")
	}

	// Create new user
	var (
		displayName   = strings.TrimSpace(swag.StringValue(params.Body.DisplayName))
		creationTime  = strfmt.DateTime(time.Now().UTC())
		imageURL      = user.GetUserImageURL{UserID: swag.StringValue(theUser.UID)}
		photoURL, err = imageURL.Build()
	)

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
		userLog.Debug("Error validating user!", err)
		return NewBadRequest(fmt.Sprintf(`Invalid user: "%s"`, err.Error()))
	}

	// Insert new user into database
	if _, err := wgplaner.OrmEngine.Update(&theUser); err != nil {
		userLog.Critical("Database error!", err)
		return userInternalServerError
	}

	return user.NewUpdateUserOK().WithPayload(&theUser)
}

func GetUser(params user.GetUserParams, principal interface{}) middleware.Responder {
	userLog.Debugf(`Get user "%s"`, params.UserID)

	theUser := models.User{UID: &params.UserID}

	if isValid := isValidUserID(theUser.UID); !isValid {
		return NewBadRequest(fmt.Sprintf("Invalid user id format"))
	}

	// Firebase Auth
	if !wgplaner.AppConfig.Auth.IgnoreFirebase {
		_, err := wgplaner.FireBaseApp.Auth().
			GetUser(params.HTTPRequest.Context(), *theUser.UID)

		if err == firebase.ErrUserNotFound {
			userLog.Debugf(`Can't find firebase user with id "%s"!`, *theUser.UID)
			return NewUnauthorizedResponse("User not authorized!")

		} else if err != nil {
			userLog.Critical("Firebase SDK Error!", err)
			return userInternalServerError
		}
	}

	// Database
	if isRegistered, err := wgplaner.OrmEngine.Get(&theUser); err != nil {
		userLog.Critical("Database Error!", err)
		return userInternalServerError

	} else if !isRegistered {
		userLog.Debugf(`Can't find database user with id "%s"!`, params.UserID)
		return NewNotFoundResponse("User not found on server")
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
	userLog.Debug("Get user image")

	// TODO: Check authorization; Check if user exists
	theUser := models.User{UID: &params.UserID}

	var imgFile *os.File
	var fileErr error

	// Get default image if normal one does no exist
	if imgFile, fileErr = wgplaner.GetUserProfileImage(&theUser); os.IsNotExist(fileErr) {
		imgFile, fileErr = wgplaner.GetUserProfileImageDefault()
	}

	if fileErr != nil {
		userLog.Error("Error getting profile image ", fileErr.Error())
		return NewInternalServerError("Internal Server Error with profile image")
	}

	return user.NewGetUserImageOK().WithPayload(imgFile)
}

func UpdateUserImage(params user.UpdateUserImageParams, principal interface{}) middleware.Responder {
	userLog.Debug("Start put user image")

	var (
		theUser       = models.User{UID: &params.UserID}
		internalError = NewInternalServerError("Internal Server Error")
	)

	// Database
	if isRegistered, err := wgplaner.OrmEngine.Get(&theUser); err != nil {
		userLog.Critical("Database Error!", err)
		return userInternalServerError

	} else if !isRegistered {
		userLog.Debugf(`Can't find database user with id "%s"!`, params.UserID)
		// TODO: Maybe 404?
		return user.NewUpdateUserImageBadRequest().WithPayload(&models.ErrorResponse{
			Message: swag.String("Unknown user"),
			Status:  swag.Int64(http.StatusBadRequest),
		})
	}

	// Check if auth and userId are the same
	if params.UserID != swag.StringValue(principal.(models.User).UID) {
		return user.NewUpdateUserImageDefault(http.StatusUnauthorized).
			WithPayload(&models.ErrorResponse{
				Message: swag.String("Can't change profile image of other users"),
				Status:  swag.Int64(http.StatusUnauthorized),
			})
	}

	// We need the first 512 Bytes for "IsValidJpeg". Because "params.ProfileImage.Data"
	// is only a reader, there is no way around extracting them.
	first512Bytes := make([]byte, 512)
	if _, err := params.ProfileImage.Data.Read(first512Bytes); err != nil {
		return internalError
	}

	if isValid, mime := wgplaner.IsValidJpeg(first512Bytes); !isValid {
		userLog.Debugf(`Invalid mime type "%s"`, mime)
		return user.NewUpdateUserImageBadRequest().WithPayload(&models.ErrorResponse{
			Message: swag.String(fmt.Sprintf(
				`Invalid file type. Only "image/jpeg" allowed. Mime was "%s"`,
				mime,
			)),
			Status: swag.Int64(http.StatusBadRequest),
		})
	}

	// Write profile image

	filePath := wgplaner.GetUserProfileImageFilePath(&theUser)

	// Create directory
	if dirErr := os.MkdirAll(path.Dir(filePath), 0700); dirErr != nil {
		userLog.Error("Can't create directory ", dirErr.Error())
		return internalError
	}

	// Create or overwrite file
	imgFile, err := os.Create(filePath)
	defer imgFile.Close()

	if err != nil {
		userLog.Debug("Can't create new file ", err.Error())
		return internalError

	} else {
		if _, err := imgFile.Write(first512Bytes); err != nil {
			userLog.Error("Couldn't write first 512Byte", err.Error())
			return internalError
		}
		if _, writeErr := io.Copy(imgFile, params.ProfileImage.Data); writeErr != nil {
			userLog.Error("Can't copy file content ",
				writeErr.Error())
			return internalError
		}
	}

	return user.NewUpdateUserImageOK().WithPayload(&models.SuccessResponse{
		Message: swag.String("Successfully uploaded image file"),
		Status:  swag.Int64(http.StatusOK),
	})
}
