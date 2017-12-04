package models

import (
	"bytes"
	"context"
	"fmt"
	"image"
	"image/jpeg"
	"os"
	"path"
	"strings"
	"time"

	"github.com/wgplaner/wg_planer_server/modules/avatar"
	"github.com/wgplaner/wg_planer_server/modules/setting"

	"github.com/acoshift/go-firebase-admin"
	"github.com/go-openapi/errors"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
	"github.com/go-openapi/validate"
	"github.com/nfnt/resize"
	"github.com/op/go-logging"
)

var userLog = logging.MustGetLogger("Group")

const (
	PROFILE_IMAGE_FILE_NAME = "profile_image.jpg"
)

// User user
// swagger:model User
type User struct {
	// display name
	// Required: true
	// Max Length: 20
	// Min Length: 3
	DisplayName *string `json:"displayName"`

	// email
	Email strfmt.Email `json:"email,omitempty"`

	// firebase instance Id
	// Pattern: ^[-_:a-zA-Z0-9]{152}$
	FirebaseInstanceID string `xorm:"VARCHAR(152)" json:"firebaseInstanceID,omitempty"`

	// group UID
	GroupUID strfmt.UUID `xorm:"VARCHAR(36) INDEX" json:"groupUID,omitempty"`

	// locale
	Locale string `xorm:"VARCHAR(5)" json:"locale,omitempty"`

	// photo Url
	PhotoURL strfmt.URI `xorm:"-" json:"photoUrl,omitempty"`

	// uid
	// Required: true
	UID *string `xorm:"varchar(28) pk" json:"uid"`

	// created at
	// Read Only: true
	CreatedAt strfmt.DateTime `xorm:"created" json:"createdAt,omitempty"`

	// updated at
	// Read Only: true
	UpdatedAt strfmt.DateTime `xorm:"updated" json:"updatedAt,omitempty"`
}

// Validate validates this user
func (u *User) Validate(formats strfmt.Registry) error {
	var res []error
	if err := u.validateDisplayName(formats); err != nil {
		res = append(res, err)
	}
	if err := u.validateFirebaseInstanceID(formats); err != nil {
		res = append(res, err)
	}
	if err := u.validateUID(formats); err != nil {
		res = append(res, err)
	}
	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (u *User) validateDisplayName(formats strfmt.Registry) error {
	if err := validate.Required("displayName", "body", u.DisplayName); err != nil {
		return err
	}
	if err := validate.MinLength("displayName", "body", string(*u.DisplayName), 3); err != nil {
		return err
	}
	if err := validate.MaxLength("displayName", "body", string(*u.DisplayName), 20); err != nil {
		return err
	}
	return nil
}

func (u *User) validateFirebaseInstanceID(formats strfmt.Registry) error {
	if swag.IsZero(u.FirebaseInstanceID) { // not required
		return nil
	}
	if err := validate.Pattern("firebaseInstanceID", "body", string(u.FirebaseInstanceID), `^[-_:a-zA-Z0-9]{152}$`); err != nil {
		return err
	}
	return nil
}

func (u *User) validateUID(formats strfmt.Registry) error {
	if err := validate.Required("uid", "body", u.UID); err != nil {
		return err
	}
	return nil
}

// MarshalBinary interface implementation
func (u *User) MarshalBinary() ([]byte, error) {
	if u == nil {
		return nil, nil
	}
	return swag.WriteJSON(u)
}

// UnmarshalBinary interface implementation
func (u *User) UnmarshalBinary(b []byte) error {
	var res User
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*u = res
	return nil
}

func (u *User) Load() (bool, error) {
	if isRegistered, err := x.Get(u); err != nil {
		return false, err
	} else {
		return isRegistered, nil
	}
}

func (u *User) IsAdmin() bool {
	g, err := GetGroupByUID(u.GroupUID)
	if err != nil {
		return false
	}
	return g.HasAdmin(*u.UID)
}

func (u *User) LeaveGroup() error {
	u.GroupUID = ""
	return UpdateUserCols(u, "group_uid")
}

func (u *User) JoinGroupWithCode(groupCode string) (*Group, error) {
	groupLog.Debugf(`User "%s" joins a group with code "%s"`, *u.UID, groupCode)

	var (
		exists  bool
		err     error
		g       *Group
		theCode *GroupCode
	)

	// Check the code and get the group uid
	if exists, theCode = IsGroupCodeValid(groupCode); !exists {
		return nil, ErrGroupCodeNotExist{Code: groupCode}
	}

	// Check group
	if g, err = GetGroupByUID(*theCode.GroupUID); err != nil {
		return nil, err
	} else if g == nil {
		return nil, ErrGroupNotExist{UID: *theCode.GroupUID}
	}

	// user joins the group.
	u.GroupUID = *theCode.GroupUID
	if _, err := x.ID(*u.UID).Update(u); err != nil {
		return nil, err
	}

	// get updated group
	// TODO: Is this necessary?
	if g, err = GetGroupByUID(g.UID); err != nil {
		return nil, err
	}

	return g, nil
}

func (u *User) BuyListItemsByUIDs(itemUIDs []strfmt.UUID) error {
	// TODO: Test if already bought

	// Check if items exist
	count, errCount := x.Where(`group_uid=?`, u.GroupUID).
		In(`id`, itemUIDs).
		Count(new(ListItem))

	if errCount != nil {
		return errCount
	} else if int(count) != len(itemUIDs) {
		return ErrListItemNotExist{}
	}

	_, err := x.Cols(`bought_by`, `bought_at`).
		Where(`group_uid=?`, u.GroupUID).
		In(`id`, itemUIDs).
		Update(&ListItem{
			BoughtAt: swag.Time(time.Now().UTC()),
			BoughtBy: *u.UID,
		})
	return err
}

func IsValidUserIDFormat(uid string) bool {
	return len(uid) == 28 // Firebase IDs are 28 characters long
}

func IsUserExist(uid string) (bool, error) {
	if !IsValidUserIDFormat(uid) {
		return false, ErrUserInvalidUID{}
	}
	return x.Exist(&User{UID: &uid})
}

func IsUserOnFirebase(uid string) (bool, error) {
	_, err := setting.FireBaseApp.Auth().GetUser(context.Background(), uid)

	if err == firebase.ErrUserNotFound {
		userLog.Debugf(`Can't find firebase user with id "%s"!`, uid)
		return false, nil

	} else if err != nil {
		userLog.Critical("Firebase SDK Error!", err)
		return false, err
	}

	return true, nil
}

func AreUsersExist(uids []string) (bool, error) {
	// TODO: Better use where in and count
	for _, uid := range uids {
		if exists, err := IsUserExist(uid); err != nil {
			return false, err
		} else if !exists {
			return false, nil
		}
	}
	return true, nil
}

func CreateUser(u *User) error {
	if u.DisplayName == nil {
		return ErrUserMissingProperty{}
	}
	*u.DisplayName = strings.TrimSpace(*u.DisplayName)
	_, err := x.InsertOne(u)
	u.PhotoURL = strfmt.URI(GetUserImageURL(*u.UID))
	return err
}

func GetUserByUID(uid string) (*User, error) {
	u := new(User)

	if has, err := x.ID(uid).Get(u); err != nil {
		return nil, err

	} else if !has {
		return nil, ErrUserNotExist{UID: uid}
	}

	u.PhotoURL = strfmt.URI(GetUserImageURL(*u.UID))

	return u, nil
}

func UpdateUser(u *User) error {
	u.DisplayName = swag.String(strings.TrimSpace(*u.DisplayName))
	_, err := x.ID(u.UID).AllCols().Update(u)
	u.PhotoURL = strfmt.URI(GetUserImageURL(*u.UID))
	return err
}

func UpdateUserCols(u *User, cols ...string) error {
	u.DisplayName = swag.String(strings.TrimSpace(swag.StringValue(u.DisplayName)))
	_, err := x.ID(u.UID).Cols(cols...).Update(u)
	u.PhotoURL = strfmt.URI(GetUserImageURL(*u.UID))
	return err
}

func GetUserImagePath(uid string) string {
	return path.Join(
		setting.AppWorkPath,
		setting.AppConfig.Data.UserImageDir,
		uid,
		PROFILE_IMAGE_FILE_NAME,
	)
}

func GetUserImage(uid string) (*os.File, error) {
	return os.Open(GetUserImagePath(uid))
}

func GetUserImageDefault() (*os.File, error) {
	return os.Open(path.Join(
		setting.AppWorkPath,
		setting.AppConfig.Data.UserImageDefault,
	))
}

func GetUserImageURL(uid string) string {
	return strings.Replace(`/users/{userID}/image`, `{userID}`, uid, -1)
}

func (u *User) UploadUserImage(data []byte) error {
	filePath := GetUserImagePath(*u.UID)
	img, _, err := image.Decode(bytes.NewReader(data))

	if err != nil {
		return fmt.Errorf("decode: %v", err)
	}

	m := resize.Resize(avatar.AvatarSize, avatar.AvatarSize, img, resize.NearestNeighbor)

	// Create directory
	if err = os.MkdirAll(path.Dir(filePath), 0700); err != nil {
		return fmt.Errorf("failed to create dir %s: %v", setting.AppConfig.Data.UserImageDir, err)
	}

	// Create or overwrite file
	fw, err := os.Create(filePath)

	if err != nil {
		return fmt.Errorf("create: %v", err)

	}
	defer fw.Close()

	if err = jpeg.Encode(fw, m, &jpeg.Options{Quality: 95}); err != nil {
		return fmt.Errorf("encode: %v", err)
	}

	return nil
}
