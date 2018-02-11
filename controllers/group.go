package controllers

import (
	"bytes"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"regexp"

	"github.com/wgplaner/wg_planer_server/models"
	"github.com/wgplaner/wg_planer_server/modules/base"
	"github.com/wgplaner/wg_planer_server/modules/mailer"
	"github.com/wgplaner/wg_planer_server/modules/setting"
	"github.com/wgplaner/wg_planer_server/restapi/operations/group"

	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
	"github.com/op/go-logging"
	"github.com/satori/go.uuid"
)

var groupLog = logging.MustGetLogger("Group")

func GetGroup(params group.GetGroupParams, principal *models.User) middleware.Responder {
	groupLog.Debugf(`User %q gets group "%s"`, *principal.UID, params.GroupUID)

	var g *models.Group
	var err error

	// Database
	if g, err = models.GetGroupByUID(params.GroupUID); models.IsErrGroupNotExist(err) {
		groupLog.Debugf(`Can't find database group with id "%s"!`, params.GroupUID)
		return NewNotFoundResponse("Group not found on server")

	} else if models.IsErrGroupInvalidUUID(err) {
		groupLog.Debugf(err.Error())
		return NewNotFoundResponse("invalid group uid")

	} else if err != nil {
		groupLog.Critical(`Database Error!`, err)
		return NewInternalServerError("Internal Database Error")
	}

	if !g.HasMember(*principal.UID) {
		return NewUnauthorizedResponse("User is a member of the specified group")
	}

	return group.NewGetGroupOK().WithPayload(g)
}

func GetGroupImage(params group.GetGroupImageParams, principal *models.User) middleware.Responder {
	groupLog.Debugf(`User %q gets image for group "%s"`, *principal.UID, params.GroupUID)
	var g *models.Group
	var err error

	if g, err = models.GetGroupByUID(params.GroupUID); models.IsErrGroupNotExist(err) {
		groupLog.Debugf(`Can't find database group with id "%s"!`, params.GroupUID)
		return NewNotFoundResponse("Group not found on server")

	} else if models.IsErrGroupInvalidUUID(err) {
		groupLog.Debugf(err.Error())
		return NewNotFoundResponse("invalid group uid")

	} else if err != nil {
		groupLog.Critical(`Database Error getting group!`, err)
		return NewInternalServerError("Internal Database Error")
	}

	//if !base.StringInSlice(*principal.UID, g.Members) {
	//	return NewUnauthorizedResponse("User not a member of the group.")
	//}

	var imgFile *os.File
	var fileErr error

	// Get default image if normal one does no exist
	if imgFile, fileErr = models.GetGroupImage(g.UID); os.IsNotExist(fileErr) {
		imgFile, fileErr = models.GetGroupImageDefault()
	}

	if fileErr != nil {
		groupLog.Error("Error getting group's profile image ", fileErr.Error())
		return NewInternalServerError("Internal Server Error with profile image")
	}

	return group.NewGetGroupImageOK().WithPayload(imgFile)

}

func CreateGroupCode(params group.CreateGroupCodeParams, principal *models.User) middleware.Responder {
	groupLog.Debugf(`User %q generates code for group %q!`, *principal.UID, params.GroupUID)

	var (
		c   *models.GroupCode
		err error

		groupUID = strfmt.UUID(params.GroupUID)
	)

	if principal.GroupUID != groupUID {
		return NewUnauthorizedResponse("Can't create group code for other groups")
	}

	// Group MUST exist or we have inconsistencies
	if _, err = models.GetGroupByUID(groupUID); err != nil {
		groupLog.Debugf(`Error validating group "%s": "%s"`, params.GroupUID, err.Error())
		return NewInternalServerError("Internal Server Error")
	}

	// TODO: Check authorization for user in the group

	if c, err = models.CreateGroupCode(groupUID); err != nil {
		groupLog.Critical("Database error!", err)
		return NewInternalServerError("Internal Database Error")
	}

	return group.NewCreateGroupCodeOK().WithPayload(c)
}

func CreateGroup(params group.CreateGroupParams, principal *models.User) middleware.Responder {
	groupLog.Debugf(`User %q starts creating a group`, *principal.UID)

	var err error

	groupUid, err := uuid.NewV4()
	if err != nil {
		groupLog.Critical("Error generating NewV4 UID!", err)
		return NewInternalServerError("Internal Error")
	}

	// Create new group
	newGroupUid := strfmt.UUID(groupUid.String())

	theUser := &models.User{
		UID:      principal.UID,
		GroupUID: newGroupUid,
	}

	theGroup := &models.Group{
		UID:         newGroupUid,
		Admins:      []string{*principal.UID},
		DisplayName: params.Body.DisplayName,
		Currency:    params.Body.Currency,
	}

	// TODO: Check if user has already a group

	if err = models.UpdateUserCols(theUser, "group_uid"); err != nil {
		groupLog.Critical("Database error!", err)
		return NewInternalServerError("Internal Database Error")
	}

	// Insert new user into database
	if err = models.CreateGroup(theGroup); err != nil {
		groupLog.Critical("Database error!", err)
		return NewInternalServerError("Internal Database Error")
	}

	if theGroup, err = models.GetGroupByUID(theGroup.UID); err != nil {
		groupLog.Critical("Database error!", err)
		return NewInternalServerError("Internal Database Error")
	}

	groupLog.Infof(`Created group "%s"`, theGroup.UID)

	return group.NewCreateGroupOK().WithPayload(theGroup)
}

func UpdateGroup(params group.UpdateGroupParams, principal *models.User) middleware.Responder {
	groupLog.Debugf(`User %q starts updating group %q`, *principal.UID, params.Body.UID)

	var g *models.Group
	var err error

	if g, err = models.GetGroupByUID(params.Body.UID); models.IsErrGroupNotExist(err) {
		groupLog.Debugf(`Update group: "%s" does not exist: %s"`, params.Body.UID, err)
		return NewNotFoundResponse(`Group does not exist.`)

	} else if models.IsErrGroupInvalidUUID(err) {
		groupLog.Debugf(`Update group: %q"`, err.Error())
		return NewBadRequest(`invalid group uid format`)

	} else if err != nil {
		groupLog.Critical("Database error!", err)
		return NewInternalServerError("Internal Database Error")
	}

	if !g.HasAdmin(*principal.UID) {
		return NewUnauthorizedResponse("Not an admin")
	}

	g.DisplayName = params.Body.DisplayName
	g.Currency = params.Body.Currency

	// Update user into database
	if err := models.UpdateGroupCols(g, `display_name`, `currency`); err != nil {
		groupLog.Critical("Database error!", err)
		return NewInternalServerError("Internal Database Error")
	}

	mailer.SendPushUpdateToUserIDs(g.Members, mailer.PushUpdateGroupData, []string{
		string(g.UID),
	})

	groupLog.Infof(`Updated group "%s"`, g.UID)

	return group.NewCreateGroupOK().WithPayload(g)
}

func JoinGroup(params group.JoinGroupParams, principal *models.User) middleware.Responder {
	groupLog.Debugf(`User %q watns to join a group with code %q`, *principal.UID, params.GroupCode)

	g, err := principal.JoinGroupWithCode(params.GroupCode)

	if models.IsErrGroupCodeNotExist(err) {
		return NewBadRequest("Invalid group code")

	} else if err != nil {
		// TODO: Handle different errors
		groupLog.Error(`Unknown Internal Server Error: `, err)
		return NewInternalServerError("Unknown Server Error")
	}

	mailer.SendPushUpdateToUserIDs(g.Members, mailer.PushUpdateGroupNewMember, []string{
		string(*principal.UID),
	})

	return group.NewJoinGroupOK().WithPayload(g)
}

func JoinGroupHelp(params group.JoinGroupHelpParams) middleware.Responder {
	groupLog.Debug(`Get help site for joining group`)

	var (
		filepath = path.Join(setting.AppWorkPath, "views/group_code.html")
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

func LeaveGroup(params group.LeaveGroupParams, principal *models.User) middleware.Responder {
	groupLog.Debugf(`user %q leaves his group`, *principal.UID)

	// TODO: Own function "u.HasGroupAndLoad()"?
	g, err := models.GetGroupByUID(principal.GroupUID)
	if models.IsErrGroupNotExist(err) {
		return NewBadRequest("user does not have a group")

	} else if err != nil {
		groupLog.Critical("Unknown Group Error")
		return NewInternalServerError("Internal Error")
	}

	if len(g.Members) <= 1 {
		// TODO: Delete group

	} else if principal.IsAdmin() {
		// Remove user as admin and add another one
		// TODO: Put in another function
		if len(g.Admins) == 1 {
			g.Admins = []string{}
		} else {
			base.RemoveStringFromSlice(g.Admins, *principal.UID)
		}

		for _, m := range g.Members {
			if m != *principal.UID {
				g.Admins = []string{m}
				break
			}
		}
		err := models.UpdateGroupCols(g, `admins`)
		if err != nil {
			groupLog.Critical("Error updating group")
		}
	}

	if err := principal.LeaveGroup(); err != nil {
		groupLog.Critical("Database error updating group!", err)
		return NewInternalServerError("Internal Database Error")
	}

	mailer.SendPushUpdateToUserIDs(g.Members, mailer.PushUpdateGroupMemberLeft, []string{
		string(*principal.UID),
	})

	return group.NewLeaveGroupOK().WithPayload(&models.SuccessResponse{
		Message: swag.String("Successfully left group"),
		Status:  swag.Int64(http.StatusOK),
	})
}

func UpdateGroupImage(params group.UpdateGroupImageParams, principal *models.User) middleware.Responder {
	groupLog.Debugf(`User %q starts updating image of group %q`, *principal.UID, params.GroupUID)

	var (
		err error
		g   *models.Group
	)

	// Database
	if g, err = models.GetGroupByUID(params.GroupUID); models.IsErrUserNotExist(err) {
		groupLog.Debugf(`Can't find database group with id "%s"!`, params.GroupUID)
		return NewNotFoundResponse("Unknown group")

	} else if models.IsErrGroupInvalidUUID(err) {
		groupLog.Debugf(err.Error())
		return NewNotFoundResponse("invalid group uid")

	} else if err != nil {
		groupLog.Critical("Database Error getting group!", err)
		return NewInternalServerError("Internal Server Error")
	}

	// TODO: use isAdmin
	if !g.HasMember(*principal.UID) {
		return NewUnauthorizedResponse("User not a member of the group.")
	}

	data, err := ioutil.ReadAll(params.ProfileImage)
	if err != nil {
		return NewInternalServerError("Internal Server Error")
	}

	if !base.IsFileJPG(data) {
		msg := fmt.Sprintf(
			`Invalid file type. Only "image/jpeg" allowed. Mime was "%s"`,
			base.GetMimeType(data),
		)
		userLog.Debugf(`Invalid mime type "%s"`, msg)
		return NewBadRequest(msg)
	}

	if err = g.UploadGroupImage(data); err != nil {
		userLog.Critical(`Error uploading group avatar.`)
		return NewInternalServerError("Internal Server Error")
	}

	mailer.SendPushUpdateToUserIDs(g.Members, mailer.PushUpdateGroupImage, []string{
		string(g.UID),
	})

	return group.NewUpdateGroupImageOK().WithPayload(&models.SuccessResponse{
		Message: swag.String("Successfully uploaded image file"),
		Status:  swag.Int64(http.StatusOK),
	})
}
