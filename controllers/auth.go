package controllers

import (
	"net/http"

	"github.com/go-openapi/errors"
	"github.com/op/go-logging"
	"github.com/wgplaner/wg_planer_server/models"
	"github.com/wgplaner/wg_planer_server/modules/setting"
)

var authLog = logging.MustGetLogger("Auth")

// userIDAuth takes an auth token and validates that token against the database.
// It returns the user if the auth token is valid and an error otherwise.
func userIDAuth(token string) (*models.User, error) {
	authLog.Debugf(`Check userID authorization for user id "%s"`, token)

	var u *models.User
	var err error

	if u, err = models.GetUserByUID(token); models.IsErrUserNotExist(err) {
		authLog.Debugf(`Unauthorized database user "%s"`, token)
		return nil, errors.Unauthenticated("invalid credentials (wgplaner account)")

	} else if err != nil {
		authLog.Error(`DB error with isUserRegistered`, err.Error())
		return nil, errors.New(http.StatusInternalServerError, "Internal Server Error")
	}

	return u, nil
}

// firebaseIDAuth takes an auth token and validates that token against firebase.
// It returns a user with only its ID set if the auth token is valid and an error otherwise.
func firebaseIDAuth(token string) (*models.User, error) {
	authLog.Debugf(`Check firebaseId authorization for user id "%s"`, token)

	if !models.IsValidUserIDFormat(token) {
		return nil, errors.Unauthenticated("invalid credentials (format)")
	}

	u := &models.User{UID: &token}

	if setting.AppConfig.Auth.IgnoreFirebase {
		authLog.Debugf(`Ignore firebase auth`)
		return u, nil
	}

	if isRegistered, err := models.IsUserOnFirebase(token); err != nil {
		authLog.Error(`DB error with IsUserOnFirebase`, err.Error())
		return nil, errors.New(http.StatusInternalServerError, "Internal Server Error")

	} else if !isRegistered {
		authLog.Debugf(`Unauthorized firebase user "%s"`, token)
		return nil, errors.Unauthenticated("invalid credentials (firebase account)")
	}

	return u, nil
}
