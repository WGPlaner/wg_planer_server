package controllers

import (
	"net/http"

	"github.com/wgplaner/wg_planer_server/models"
	"github.com/wgplaner/wg_planer_server/modules/setting"

	"github.com/go-openapi/errors"
	"github.com/op/go-logging"
)

var authLog = logging.MustGetLogger("Auth")

func UserIDAuth(token string) (*models.User, error) {
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

func FirebaseIDAuth(token string) (*models.User, error) {
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
