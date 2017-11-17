package controllers

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
	"github.com/op/go-logging"
	"github.com/satori/go.uuid"
	"github.com/wgplaner/wg_planer_server/gen/models"
	"github.com/wgplaner/wg_planer_server/gen/restapi/operations/group"
	"github.com/wgplaner/wg_planer_server/wgplaner"
)

var groupLog = logging.MustGetLogger("Group")

const (
	GROUP_CODE_LENGTH     = 12
	GROUP_CODE_VALID_DAYS = 3
)

type groupErrorCode int

type groupError struct {
	code groupErrorCode
	msg  string
}

func (err *groupError) Error() string {
	return err.msg
}

func (err *groupError) Code() groupErrorCode {
	return err.code
}

func (err *groupError) is(code groupErrorCode) bool {
	return err.code == code
}

const (
	ERR_GROUP_NOT_FOUND groupErrorCode = iota
	ERR_GROUP_DATABASE
	ERR_GROUP_USER_NOT_AUTHORIZED
	ERR_GROUP_CODE_EXPIRED
	ERR_GROUP_CODE_INVALID
)

var (
	errGroupDatabase          = groupError{ERR_GROUP_DATABASE, "Internal Database Error"}
	errGroupNotFound          = groupError{ERR_GROUP_NOT_FOUND, "Group not found"}
	errGroupCodeInvalid       = groupError{ERR_GROUP_CODE_INVALID, "Invalid group code"}
	errGroupCodeExpired       = groupError{ERR_GROUP_CODE_EXPIRED, "Group code expired"}
	errGroupUserNotAuthorized = groupError{ERR_GROUP_USER_NOT_AUTHORIZED, "User not authorized"}
)

func validateGroup(_ *models.Group) (bool, error) {
	// TODO
	return true, nil
}

func joinGroupWithCode(theUser *models.User, groupCode string) (*models.Group, *groupError) {
	theCode := models.GroupCode{Code: swag.String(groupCode)}

	if keyExists, err := wgplaner.OrmEngine.Get(&theCode); err != nil {
		groupLog.Critical(`Database Error!`, err)
		return nil, &errGroupDatabase

	} else if !keyExists {
		groupLog.Debugf(`Can't find database group code with id "%s"!`, groupCode)
		return nil, &errGroupCodeInvalid
	}

	if time.Now().After(time.Time(theCode.ValidUntil)) {
		return nil, &errGroupCodeExpired
	}

	// TODO: Check group

	return &models.Group{
		UID: *theCode.GroupUID,
	}, &errGroupNotFound
}

func GetGroup(params group.GetGroupParams, principal interface{}) middleware.Responder {
	theGroup := models.Group{UID: strfmt.UUID(params.GroupID)}
	groupLog.Debugf(`Get group "%s"`, theGroup.UID)

	// TODO: Validate
	// validateGroup(&theGroup)

	// Database
	if isRegistered, err := wgplaner.OrmEngine.Get(&theGroup); err != nil {
		groupLog.Critical(`Database Error!`, err)
		return userInternalServerError

	} else if !isRegistered {
		groupLog.Debugf(`Can't find database group with id "%s"!`, theGroup.UID)
		return group.NewGetGroupNotFound().WithPayload(&models.ErrorResponse{
			Message: swag.String("Group not found on server"),
			Status:  swag.Int64(http.StatusNotFound),
		})
	}

	return group.NewGetGroupOK().WithPayload(&theGroup)
}

func CreateGroupCode(params group.CreateGroupCodeParams, principal interface{}) middleware.Responder {
	groupLog.Debugf(`Generate group code for group "%s"!`, params.GroupID)

	// TODO: Check authorization for user in the group

	var (
		groupUid   = strfmt.UUID(params.GroupID)
		code       = wgplaner.RandomAlphaNumCode(GROUP_CODE_LENGTH, true)
		validUntil = time.Now().UTC().AddDate(0, 0, GROUP_CODE_VALID_DAYS)
	)

	groupCode := models.GroupCode{
		GroupUID:   &groupUid,
		Code:       &code,
		ValidUntil: strfmt.DateTime(validUntil),
	}

	// Set valid-until date of old codes to just now. This makes only the latest code valid.
	oldCodes := models.GroupCode{ValidUntil: strfmt.DateTime(time.Now().UTC())}
	_, err := wgplaner.OrmEngine.Where("group_u_i_d = ?", groupUid).Update(&oldCodes)
	if err != nil {
		groupLog.Errorf(`Can't update validUntil date of other (old) codes for group "%s"'; Err: "%s"`,
			params.GroupID, err)
	}

	// Insert new code into database
	if _, err := wgplaner.OrmEngine.InsertOne(&groupCode); err != nil {
		groupLog.Critical("Database error!", err)
		return userInternalServerError
	}

	return group.NewCreateGroupCodeOK().WithPayload(&groupCode)
}

func CreateGroup(params group.CreateGroupParams, principal interface{}) middleware.Responder {
	groupLog.Debug(`Start creating group`)

	theGroup := models.Group{}

	// Create new group
	displayName := strings.TrimSpace(swag.StringValue(params.Body.DisplayName))
	creationTime := strfmt.DateTime(time.Now().UTC())
	currency := strings.TrimSpace(params.Body.Currency)

	if currency == "" {
		currency = "€"
	}

	theGroup = models.Group{
		UID:         strfmt.UUID(uuid.NewV4().String()),
		Admins:      []string{*principal.(models.User).UID},
		Members:     []string{*principal.(models.User).UID},
		DisplayName: &displayName,
		Currency:    currency,
		CreatedAt:   creationTime,
		UpdatedAt:   creationTime,
	}

	// Validate group
	if isValid, err := validateGroup(&theGroup); !isValid {
		groupLog.Notice("Error validating user!", err)
		return group.NewCreateGroupBadRequest().WithPayload(&models.ErrorResponse{
			Message: swag.String(fmt.Sprintf(`Invalid group data: "%s"`, err.Error())),
			Status:  swag.Int64(400),
		})
	}

	// TODO: Check if user has already a group

	// Insert new user into database
	if _, err := wgplaner.OrmEngine.InsertOne(&theGroup); err != nil {
		groupLog.Critical("Database error!", err)
		return userInternalServerError
	}

	groupLog.Infof(`Created group "%s"`, theGroup.UID)

	return group.NewCreateGroupOK().WithPayload(&theGroup)
}

func JoinGroup(params group.JoinGroupParams, principal interface{}) middleware.Responder {
	// TODO: Acutally Implement This

	theUser := principal.(models.User)
	theGroup, err := joinGroupWithCode(&theUser, params.GroupCode)

	if err != nil {
		return group.NewJoinGroupOK().WithPayload(theGroup)
	}

	switch err.Code() {
	case ERR_GROUP_CODE_EXPIRED:
		groupLog.Debugf(`Group code "%s" expired`, params.GroupCode)
		return group.NewJoinGroupDefault(http.StatusBadRequest).
			WithPayload(&models.ErrorResponse{
				Message: swag.String(err.Error()),
				Status:  swag.Int64(http.StatusBadRequest),
			})

	case ERR_GROUP_CODE_INVALID:
		groupLog.Debugf(`Invalid group code "%s"`, params.GroupCode)
		return group.NewJoinGroupDefault(http.StatusInternalServerError).
			WithPayload(&models.ErrorResponse{
				Message: swag.String(err.Error()),
				Status:  swag.Int64(http.StatusInternalServerError),
			})

	case ERR_GROUP_NOT_FOUND:
		groupLog.Debugf(`Group was deleted but the code "%s" is still valid: %s`,
			params.GroupCode, err.Error())
		// TODO: This should not happen
		return group.NewJoinGroupDefault(http.StatusNotFound).
			WithPayload(&models.ErrorResponse{
				Message: swag.String(err.Error()),
				Status:  swag.Int64(http.StatusNotFound),
			})

	default:
		groupLog.Error(`Unknown Internal Server Error`, err)
		return group.NewJoinGroupDefault(http.StatusInternalServerError).
			WithPayload(&models.ErrorResponse{
				Message: swag.String("Unknown Server Error"),
				Status:  swag.Int64(http.StatusInternalServerError),
			})
	}

}

func JoinGroupHelp(params group.JoinGroupHelpParams) middleware.Responder {
	groupLog.Debug(`Get help site for joining group`)
	return group.NewJoinGroupHelpOK().WithPayload(`<!DOCTYPE html>
<html>
	<head>
		<title>Join Group</title>
	</head>
	<body>
		not implemented yet
	</body>
</html>`)
}

func LeaveGroup(params group.LeaveGroupParams, principal interface{}) middleware.Responder {
	return group.NewLeaveGroupDefault(501)
}
