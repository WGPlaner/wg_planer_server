package controllers

import (
	"bytes"
	"fmt"
	"html/template"
	"net/http"
	"path"
	"regexp"
	"strings"
	"time"

	"github.com/wgplaner/wg_planer_server/gen/models"
	"github.com/wgplaner/wg_planer_server/gen/restapi/operations/group"
	"github.com/wgplaner/wg_planer_server/wgplaner"

	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
	"github.com/op/go-logging"
	"github.com/satori/go.uuid"
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
	errGroupInternalError     = &groupError{5000, "Internal Server Error"}
)

func validateGroup(_ *models.Group) (bool, error) {
	// TODO
	return true, nil
}

// Validate the given group UUID. This means checking if the UUID is valid and the group exists.
func validateGroupUuid(groupUid strfmt.UUID) *groupError {
	if !strfmt.IsUUID(string(groupUid)) {
		groupLog.Debugf(`Invalid group ID pattern "%s"!`, groupUid)
		return errGroupInvalidUUID
	}

	groupExists, err := wgplaner.OrmEngine.Exist(&models.Group{UID: groupUid})

	if err != nil {
		groupLog.Critical(`Database Error querying group!`, err)
		return errGroupDatabase

	} else if !groupExists {
		return errGroupNotFound
	}

	return nil
}

// Validate the given group code model. Checks if the code is valid.
// Queries the database. "theCode" will be updated.
func validateGroupCode(theCode *models.GroupCode) *groupError {
	if keyExists, err := wgplaner.OrmEngine.Get(theCode); err != nil {
		groupLog.Critical(`Database Error querying group code!`, err)
		return errGroupDatabase

	} else if !keyExists {
		groupLog.Debugf(`Can't find database group code with id "%s"!`, theCode.Code)
		return errGroupCodeInvalid
	}

	if time.Now().After(time.Time(theCode.ValidUntil)) {
		return errGroupCodeExpired
	}

	return nil
}

// Add "member" to the member-field of group. Loads the group from the database
// and updates it.
func groupAddMember(theGroup *models.Group, member models.User) *groupError {
	// Get the group
	if exists, err := wgplaner.OrmEngine.Get(theGroup); err != nil {
		groupLog.Critical("Database error querying group!", err)
		return errGroupDatabase

	} else if !exists {
		groupLog.Critical(`Group to join not found!`)
		return errGroupNotFound
	}

	theGroup.Members = wgplaner.AppendUniqueString(theGroup.Members, *member.UID)

	// Update the group
	if _, err := wgplaner.OrmEngine.Update(theGroup); err != nil {
		groupLog.Critical("Database error updating the group!", err)
		return errGroupDatabase
	}

	return nil
}

func joinGroupWithCode(theUser *models.User, groupCode string) (*models.Group, *groupError) {
	groupLog.Debugf(`User "%s" joins a group with code "%s"`, *theUser.UID, groupCode)
	theCode := models.GroupCode{Code: swag.String(groupCode)}

	// Check the code and get the group uid
	if err := validateGroupCode(&theCode); err != nil {
		return nil, err
	}

	// TODO: Check group (if it exists)

	// user joins the group.
	theUser.GroupUID = *theCode.GroupUID
	if _, err := wgplaner.OrmEngine.Update(theUser); err != nil {
		groupLog.Critical("Database error updating user!", err, theUser)
		return nil, errGroupDatabase
	}

	// user is added to group
	// TODO: This should not be necessary; members should be read dynamically
	theGroup := models.Group{UID: *theCode.GroupUID}
	if err := groupAddMember(&theGroup, *theUser); err != nil {
		return nil, errGroupInternalError
	}

	return &theGroup, nil
}

func GetGroup(params group.GetGroupParams, principal interface{}) middleware.Responder {
	theUser := principal.(models.User)
	theGroup := models.Group{UID: strfmt.UUID(params.GroupID)}
	groupLog.Debugf(`Get group "%s"`, theGroup.UID)

	// TODO: Validate
	// validateGroup(&theGroup)

	// Database
	if exists, err := wgplaner.OrmEngine.Get(&theGroup); err != nil {
		groupLog.Critical(`Database Error!`, err)
		return userInternalServerError

	} else if !exists {
		groupLog.Debugf(`Can't find database group with id "%s"!`, theGroup.UID)
		return NewNotFoundResponse("Group not found on server")
	}

	if !wgplaner.StringInSlice(*theUser.UID, theGroup.Members) {
		return NewUnauthorizedResponse("User is a member of the specified group")
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
		groupLog.Debugf(`Error validating group "%s": "%s"`, params.GroupID, err.Error())
		return NewBadRequest(err.Error())
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

	// Create new group
	var (
		displayName  = strings.TrimSpace(swag.StringValue(params.Body.DisplayName))
		creationTime = strfmt.DateTime(time.Now().UTC())
		currency     = strings.TrimSpace(params.Body.Currency)
	)

	if currency == "" {
		currency = "€"
	}

	newGroupUid := strfmt.UUID(uuid.NewV4().String())

	theUser := models.User{
		UID:      principal.(models.User).UID,
		GroupUID: newGroupUid,
	}

	theGroup := models.Group{
		UID:         newGroupUid,
		Admins:      []string{*theUser.UID},
		Members:     []string{*theUser.UID},
		DisplayName: &displayName,
		Currency:    currency,
		CreatedAt:   creationTime,
		UpdatedAt:   creationTime,
	}

	// Validate group
	if isValid, err := validateGroup(&theGroup); !isValid {
		groupLog.Notice("Error validating user!", err)
		return NewBadRequest(fmt.Sprintf(`Invalid group data: "%s"`, err.Error()))
	}

	// TODO: Check if user has already a group
	if _, err := wgplaner.OrmEngine.Update(&theUser); err != nil {
		groupLog.Critical("Database error!", err)
		return userInternalServerError
	}

	// Insert new user into database
	if _, err := wgplaner.OrmEngine.InsertOne(&theGroup); err != nil {
		groupLog.Critical("Database error!", err)
		return userInternalServerError
	}

	groupLog.Infof(`Created group "%s"`, theGroup.UID)

	return group.NewCreateGroupOK().WithPayload(&theGroup)
}

func JoinGroup(params group.JoinGroupParams, principal interface{}) middleware.Responder {
	theUser := principal.(models.User)
	theGroup, err := joinGroupWithCode(&theUser, params.GroupCode)

	if err == nil {
		return group.NewJoinGroupOK().WithPayload(theGroup)
	}

	switch err {
	case errGroupCodeExpired:
		groupLog.Debugf(`Group code "%s" expired`, params.GroupCode)
		return NewBadRequest(err.Error())

	case errGroupCodeInvalid:
		groupLog.Debugf(`Invalid group code "%s"`, params.GroupCode)
		return NewBadRequest(err.Error())

	case errGroupNotFound:
		groupLog.Debugf(`Group was deleted but the code "%s" is still valid: %s`,
			params.GroupCode, err.Error())
		// TODO: This should not happen
		return NewNotFoundResponse(err.Error())

	default:
		groupLog.Error(`Unknown Internal Server Error: `, err)
		return NewBadRequest("Unknown Server Error")
	}

}

func JoinGroupHelp(params group.JoinGroupHelpParams) middleware.Responder {
	groupLog.Debug(`Get help site for joining group`)

	var (
		filepath = path.Join(wgplaner.AppWorkPath, "views/group_code.html")
		templ    = template.Must(template.ParseFiles(filepath))
		buf      = bytes.NewBuffer([]byte{})
		content  = map[string]string{"GroupCode": params.GroupCode}
	)

	r := regexp.MustCompile(`^[A-Z0-9]{12}$`)
	if !r.MatchString(params.GroupCode) {
		return group.NewJoinGroupHelpDefault(http.StatusBadRequest).WithPayload("Error. Your Code is invalid!")
	}

	if err := templ.Execute(buf, content); err != nil {
		groupLog.Error(`Can't execute template'`, err)
		return group.NewJoinGroupHelpOK().WithPayload("Internal Server Error")
	}

	return group.NewJoinGroupHelpOK().WithPayload(buf.String())
}

func LeaveGroup(params group.LeaveGroupParams, principal interface{}) middleware.Responder {
	theUser := principal.(models.User)
	theUser.GroupUID = ""

	// Insert new user into database
	if _, err := wgplaner.OrmEngine.Cols("group_u_i_d").Update(&theUser); err != nil {
		userLog.Critical("Database error updating group!", err)
		return userInternalServerError
	}

	return group.NewLeaveGroupOK().WithPayload(&models.SuccessResponse{
		Message: swag.String("Successfully left group"),
		Status:  swag.Int64(http.StatusOK),
	})
}
