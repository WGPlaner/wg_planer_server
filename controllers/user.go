package controllers

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/go-openapi/swag"
	"github.com/wgplaner/wg_planer_server/models"
	"github.com/wgplaner/wg_planer_server/modules/base"
	"github.com/wgplaner/wg_planer_server/modules/mailer"
	"github.com/wgplaner/wg_planer_server/modules/setting"
	"github.com/wgplaner/wg_planer_server/restapi/operations/user"

	"github.com/acoshift/go-firebase-admin"
	"github.com/go-openapi/runtime/middleware"
	"github.com/op/go-logging"
)

var userLog = logging.MustGetLogger("User")

func createUser(params user.CreateUserParams, principal *models.User) middleware.Responder {
	userLog.Debugf(`Start creating user %q`, *params.Body.UID)

	if *principal.UID != *params.Body.UID {
		userLog.Infof(`Authorized user "%s" tried to create account for "%s"`,
			*principal.UID, *params.Body.UID)
		return NewUnauthorizedResponse(`Can't create user for others.`)
	}

	// Check if the user is already registered
	if _, err := models.GetUserByUID(*params.Body.UID); err == nil {
		userLog.Debugf(`User "%s" already exists!`, *params.Body.UID)
		return NewBadRequest("User already exists.")

	} else if !models.IsErrUserNotExist(err) {
		userLog.Critical("Database Error!", err)
		return newInternalServerError("Internal Database Error")
	}

	// Create new user
	userBuilder := models.UserBuilder{}
	userBuilder.SetDisplayName(params.Body.DisplayName)
	userBuilder.SetUID(params.Body.UID)
	userBuilder.SetEmail(params.Body.Email)
	userBuilder.SetFirebaseInstanceID(params.Body.FirebaseInstanceID)
	u, err := userBuilder.Construct()

	if err != nil {
		userLog.Critical("Created invalid user!", err)
		return NewBadRequest(err.Error())
	}

	// Insert new user into database
	if err := models.CreateUser(&u); err != nil {
		userLog.Critical("Database error!", err)
		return newInternalServerError("Internal Database Error")
	}

	userLog.Infof(`Created user "%s"`, *u.UID)

	return user.NewCreateUserOK().WithPayload(&u)
}

func updateUser(params user.UpdateUserParams, principal *models.User) middleware.Responder {
	userLog.Debugf(`Start updating user %q`, *params.Body.UID)

	var theUser *models.User

	if *principal.UID != *params.Body.UID {
		userLog.Infof(`Authorized user "%s" tried to update account for "%s"`,
			*principal.UID, *params.Body.UID)
		return NewUnauthorizedResponse(`Can't update user for others.`)
	}

	var err error
	// Check if the user is already registered
	if theUser, err = models.GetUserByUID(*params.Body.UID); models.IsErrUserNotExist(err) {
		userLog.Infof(`User "%s" does not exist!`, *params.Body.UID)
		return NewBadRequest("User does not exist!")

	} else if err != nil {
		userLog.Critical("Database Error!", err)
		return newInternalServerError("Internal Database Error")
	}

	// Create new user
	theUser = &models.User{
		UID:                params.Body.UID,
		DisplayName:        params.Body.DisplayName,
		Email:              params.Body.Email,
		FirebaseInstanceID: params.Body.FirebaseInstanceID,
	}

	// Insert new user into database
	err = models.UpdateUserCols(theUser, "display_name", "email", "firebase_instance_id")
	if err != nil {
		userLog.Critical("Database error!", err)
		return newInternalServerError("Internal Database Error")
	}

	// Get the updated user
	if theUser, err = models.GetUserByUID(*theUser.UID); err != nil {
		return newInternalServerError("Internal Database Error")
	}

	// Send a notification to all members of the user's group.
	if principal.GroupUID != "" {
		userLog.Debugf(`Updated user. Send message to members of group %q`, principal.GroupUID)

		UIDs, err := models.GetGroupMemberUIDs(principal.GroupUID)
		if err != nil {
			userLog.Criticalf("Error getting group members %q", principal.GroupUID)
			return newInternalServerError("Internal Server Error")
		}

		mailer.SendPushUpdateToUserIDs(UIDs, mailer.PushUserUpdate, []string{
			string(*principal.UID),
		})
	}

	return user.NewUpdateUserOK().WithPayload(theUser)
}

func getUser(params user.GetUserParams, principal *models.User) middleware.Responder {
	userLog.Debugf(`User %q gets user %q`, *principal.UID, params.UserID)

	var (
		err error
		u   *models.User
	)

	if !models.IsValidUserIDFormat(params.UserID) {
		return NewBadRequest(fmt.Sprintf("Invalid user id format"))
	}

	// Firebase Auth
	if !setting.AppConfig.Auth.IgnoreFirebase {
		_, err := setting.FireBaseApp.Auth().
			GetUser(params.HTTPRequest.Context(), params.UserID)

		if err == firebase.ErrUserNotFound {
			userLog.Debugf(`Can't find firebase user with id "%s"!`, params.UserID)
			return NewUnauthorizedResponse("User not authorized!")

		} else if err != nil {
			userLog.Critical("Firebase SDK Error!", err)
			return newInternalServerError("Internal Firebase Error")
		}
	}

	// Database
	if u, err = models.GetUserByUID(params.UserID); models.IsErrUserNotExist(err) {
		userLog.Debugf(`Can't find database user with id "%s"!`, params.UserID)
		return newNotFoundResponse("User not found on server")

	} else if err != nil {
		userLog.Critical("Database Error!", err)
		return newInternalServerError("Internal Database Error")
	}

	return user.NewGetUserOK().WithPayload(u)
}

func getUserBoughtItems(params user.GetUserBoughtItemsParams, principal *models.User) middleware.Responder {
	userLog.Debugf("Get bought items for user %q", *principal.UID)
	// TODO: Check if user is authorized to get bought items
	var err error
	var itemUser *models.User

	if itemUser, err = models.GetUserByUID(params.UserID); models.IsErrUserNotExist(err) {
		userLog.Debugf(`Can't find database user with id "%s"!`, params.UserID)
		return newNotFoundResponse("User not found on server")

	} else if err != nil {
		userLog.Critical("Database Error!", err)
		return newInternalServerError("Internal Database Error")
	}

	items, err := itemUser.GetBoughtItems()
	if err != nil {
		userLog.Critical("Error getting bought items!", err)
		return newInternalServerError("Internal Database Error")
	}

	return user.NewGetUserBoughtItemsOK().WithPayload(&items)
}

func getUserImage(params user.GetUserImageParams, principal *models.User) middleware.Responder {
	userLog.Debugf("Get user image for user %q", *principal.UID)

	// TODO: Maybe "IsUserExist"
	if _, err := models.GetUserByUID(params.UserID); models.IsErrUserNotExist(err) {
		userLog.Debugf(`Can't find database user with id "%s"!`, params.UserID)
		return newNotFoundResponse("User not found on server")

	} else if err != nil {
		userLog.Critical("Database Error!", err)
		return newInternalServerError("Internal Database Error")
	}

	var imgFile *os.File
	var fileErr error

	// Get default image if normal one does no exist
	if imgFile, fileErr = models.GetUserImage(params.UserID); os.IsNotExist(fileErr) {
		imgFile, fileErr = models.GetUserImageDefault()
	}

	if fileErr != nil {
		userLog.Error("Error getting profile image ", fileErr.Error())
		return newInternalServerError("Internal Server Error with profile image")
	}

	return user.NewGetUserImageOK().WithPayload(imgFile)
}

func updateUserImage(params user.UpdateUserImageParams, principal *models.User) middleware.Responder {
	userLog.Debugf("Start put user image for user %q", *principal.UID)

	// Check if auth and userID are the same.
	// We don't have to get the user again since principal contains the loaded user
	if params.UserID != swag.StringValue(principal.UID) {
		return NewUnauthorizedResponse("Can't change profile image of other users")
	}

	data, err := ioutil.ReadAll(params.ProfileImage)
	if err != nil {
		return newInternalServerError("Internal Server Error")
	}

	if !base.IsFileJPG(data) {
		msg := fmt.Sprintf(
			`Invalid file type. Only "image/jpeg" allowed. Mime was "%s"`,
			base.GetMimeType(data),
		)
		userLog.Debugf(`Invalid mime type "%s"`, msg)
		return NewBadRequest(msg)
	}

	if err = principal.UploadUserImage(data); err != nil {
		userLog.Critical(`Error uploading user avatar.`)
		return newInternalServerError("Internal Server Error")
	}

	// Send a notification to all members of the user's group.
	if principal.GroupUID != "" {
		UIDs, err := models.GetGroupMemberUIDs(principal.GroupUID)
		if err != nil {
			userLog.Criticalf("Error getting group members %q", principal.GroupUID)
			return newInternalServerError("Internal Server Error")
		}
		mailer.SendPushUpdateToUserIDs(UIDs, mailer.PushUserUpdateImage, []string{
			string(*principal.UID),
		})
	}

	return user.NewUpdateUserImageOK().WithPayload(&models.SuccessResponse{
		Message: swag.String("Successfully uploaded image file"),
		Status:  swag.Int64(http.StatusOK),
	})
}
