package controllers

import (
	"bytes"
	"fmt"
	"html/template"
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

var (
	errGroupDatabase          = &groupError{1001, "Internal Database Error"}
	errGroupInvalidUUID       = &groupError{2001, "Group UUID invalid"}
	errGroupNotFound          = &groupError{2004, "Group not found"}
	errGroupUserNotAuthorized = &groupError{3001, "User not authorized"}
	errGroupCodeInvalid       = &groupError{4001, "Invalid group code"}
	errGroupCodeExpired       = &groupError{4002, "Group code expired"}
)

func validateGroup(_ *models.Group) (bool, error) {
	// TODO
	return true, nil
}

func validateGroupUuid(groupUid strfmt.UUID) *groupError {
	groupExists, err := wgplaner.OrmEngine.Get(&models.Group{UID: groupUid})

	if err != nil {
		groupLog.Critical(`Database Error!`, err)
		return errGroupDatabase

	} else if !groupExists {
		return errGroupNotFound

	} else {
		return nil
	}
}

func joinGroupWithCode(theUser *models.User, groupCode string) (*models.Group, *groupError) {
	theCode := models.GroupCode{Code: swag.String(groupCode)}

	if keyExists, err := wgplaner.OrmEngine.Get(&theCode); err != nil {
		groupLog.Critical(`Database Error!`, err)
		return nil, errGroupDatabase

	} else if !keyExists {
		groupLog.Debugf(`Can't find database group code with id "%s"!`, groupCode)
		return nil, errGroupCodeInvalid
	}

	if time.Now().After(time.Time(theCode.ValidUntil)) {
		return nil, errGroupCodeExpired
	}

	// TODO: Check group

	return &models.Group{
		UID: *theCode.GroupUID,
	}, errGroupNotFound
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

// Set valid-until date of old codes to just now to invalidate them.
func invalidateCodesOfGroup(groupUuid strfmt.UUID) {
	var (
		oldCodes = models.GroupCode{ValidUntil: strfmt.DateTime(time.Now().UTC())}
		_, err   = wgplaner.OrmEngine.Where("group_u_i_d = ?", groupUuid).Update(&oldCodes)
	)

	if err != nil {
		groupLog.Errorf(`Can't update validUntil date of other (old) `+
			`codes for group "%s"'; Err: "%s"`, groupUuid, err)
	}
}

func CreateGroupCode(params group.CreateGroupCodeParams, principal interface{}) middleware.Responder {
	groupLog.Debugf(`Generate group code for group "%s"!`, params.GroupID)

	groupUid := strfmt.UUID(params.GroupID)

	if err := validateGroupUuid(groupUid); err != nil {
		groupLog.Debugf(`Error validating group "%s"`, params.GroupID)
		return group.NewCreateGroupBadRequest().WithPayload(&models.ErrorResponse{
			Message: swag.String(err.Error()),
			Status:  swag.Int64(http.StatusBadRequest),
		})
	}

	// TODO: Check authorization for user in the group

	invalidateCodesOfGroup(groupUid)

	var (
		code       = wgplaner.RandomAlphaNumCode(GROUP_CODE_LENGTH, true)
		validUntil = time.Now().UTC().AddDate(0, 0, GROUP_CODE_VALID_DAYS)
	)

	groupCode := models.GroupCode{
		GroupUID:   &groupUid,
		Code:       &code,
		ValidUntil: strfmt.DateTime(validUntil),
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

	switch err {
	case errGroupCodeExpired:
		groupLog.Debugf(`Group code "%s" expired`, params.GroupCode)
		return group.NewJoinGroupDefault(http.StatusBadRequest).
			WithPayload(&models.ErrorResponse{
				Message: swag.String(err.Error()),
				Status:  swag.Int64(http.StatusBadRequest),
			})

	case errGroupCodeInvalid:
		groupLog.Debugf(`Invalid group code "%s"`, params.GroupCode)
		return group.NewJoinGroupDefault(http.StatusInternalServerError).
			WithPayload(&models.ErrorResponse{
				Message: swag.String(err.Error()),
				Status:  swag.Int64(http.StatusInternalServerError),
			})

	case errGroupNotFound:
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

	var (
		templ   = template.Must(template.ParseFiles("./views/group_code.html"))
		buf     = bytes.NewBuffer([]byte{})
		content = map[string]string{"GroupCode": params.GroupCode}
	)

	if err := templ.Execute(buf, content); err != nil {
		groupLog.Error(`Can't execute template'`, err)
		return group.NewJoinGroupHelpOK().WithPayload("Error")
	}

	return group.NewJoinGroupHelpOK().WithPayload(buf.String())
}

func LeaveGroup(params group.LeaveGroupParams, principal interface{}) middleware.Responder {
	return group.NewLeaveGroupDefault(501)
}
