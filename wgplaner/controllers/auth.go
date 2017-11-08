package controllers

import (
	"log"
	"net/http"

	"github.com/go-openapi/errors"
	"github.com/wgplaner/wg_planer_server/gen/models"
)

func UserIDAuth(token string) (interface{}, error) {
	theUser := models.User{UID: &token}

	if isRegistered, err := isUserRegistered(&theUser); err != nil {
		log.Println("[Controller][Auth] Error with isUserRegistered ", err.Error())
		return nil, errors.New(http.StatusInternalServerError, "Internal Server Error")

	} else if !isRegistered {
		return nil, errors.Unauthenticated("invalid credentials")
	}

	return theUser, nil
}

func FirebaseIDAuth(token string) (interface{}, error) {
	theUser := models.User{UID: &token}

	if isRegistered, err := isUserOnFirebase(&theUser); err != nil {
		log.Println("[Controller][Auth] Error with isUserOnFirebase ", err.Error())
		return nil, errors.New(http.StatusInternalServerError, "Internal Server Error")

	} else if !isRegistered {
		return nil, errors.Unauthenticated("invalid credentials")
	}

	return theUser, nil
}
