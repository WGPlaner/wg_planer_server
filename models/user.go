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
	ProfileImageFileName = "profile_image.jpg"
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

// Load returns whether the user exists on the database.
func (u *User) Load() (bool, error) {
	isRegistered, err := x.Get(u)
	if err == nil {
		return isRegistered, nil
	}
	return false, err
}

func (u *User) IsAdmin() bool {
	g, err := GetGroupByUID(u.GroupUID)
	if err != nil {
		return false
	}
	return g.HasAdmin(*u.UID)
}

func (u *User) LeaveGroup() error {
	// Delete "bought" with no bill.
	x.Where(``).
		Where(`bought_by=?`, *u.UID).
		And(`bill_uid IS NULL OR bill_uid=""`).
		Delete(&ListItem{})

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

// RevertListItemPurchaseByUID reverts the buying action for given list items.
func (u *User) RevertListItemPurchaseByUID(itemUID strfmt.UUID) error {
	var item ListItem
	// Check if items exist
	found, errCount := x.Where(`group_uid=?`, u.GroupUID).
		And(`id=?`, itemUID).
		Get(&item)

	if errCount != nil {
		return errCount

	} else if !found {
		return ErrListItemNotExist{ID: itemUID, GroupUID: u.GroupUID}
	}

	if item.BillUID != "" {
		return ErrListItemHasBill{ID: itemUID, GroupUID: u.GroupUID}
	}

	_, err := x.Cols(`bought_by`, `bought_at`).
		Where(`group_uid=?`, u.GroupUID).
		And(`id=?`, itemUID).
		Update(&ListItem{
			BoughtAt: nil,
			BoughtBy: "",
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
	// Validate using builder
	userBuilder := UserBuilder{*u}
	user, builderErr := userBuilder.Construct()
	if builderErr != nil {
		return builderErr
	}
	_, err := x.InsertOne(&user)
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
	_, err := x.ID(*u.UID).AllCols().Update(u)
	u.PhotoURL = strfmt.URI(GetUserImageURL(*u.UID))
	return err
}

func UpdateUserCols(u *User, cols ...string) error {
	u.DisplayName = swag.String(strings.TrimSpace(swag.StringValue(u.DisplayName)))
	_, err := x.ID(*u.UID).Cols(cols...).Update(u)
	u.PhotoURL = strfmt.URI(GetUserImageURL(*u.UID))
	return err
}

func GetUserImagePath(uid string) string {
	return path.Join(
		setting.AppWorkPath,
		setting.AppConfig.Data.UserImageDir,
		uid,
		ProfileImageFileName,
	)
}

// GetUserImage returns the user image or an error if user does not exist
func GetUserImage(uid string) (*os.File, error) {
	return os.Open(GetUserImagePath(uid))
}

// GetUserImageDefault returns the default user image or an error if it does not exist
func GetUserImageDefault() (*os.File, error) {
	return os.Open(GetUserImageDefaultPath())
}

// GetUserImageDefaultPath returns the path of the default user image as a string.
func GetUserImageDefaultPath() string {
	return path.Join(
		setting.AppWorkPath,
		setting.AppConfig.Data.UserImageDefault,
	)
}

// GetUserImageURL returns the API URL of the user image as a string.
func GetUserImageURL(uid string) string {
	return strings.Replace(`/users/{userID}/image`, `{userID}`, uid, -1)
}

// UploadUserImage takes "data" and saves it on disk.
// Returns an error which is nil if everything worked.
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

// GetBoughtItems
func (u *User) GetBoughtItems() (ShoppingList, error) {
	items := make([]*ListItem, 0, 10)

	// Check if items exist
	err := x.AllCols().
		Where(`bought_by=?`, *u.UID).
		And(`bill_uid IS NULL OR bill_uid=""`).
		Find(&items)

	if err != nil {
		return ShoppingList{}, err
	}

	shoppingList := ShoppingList{
		Count:     int64(len(items)),
		ListItems: items,
	}

	return shoppingList, nil
}

// UserBuilder is used to build a new user (and validate it)
type UserBuilder struct {
	user User
}

// Construct checks if user is valid an returns it.
func (u *UserBuilder) Construct() (User, error) {
	if u.user.UID == nil || *u.user.UID == "" {
		return u.user, ErrUserMissingProperty{"user needs an UID"}
	}
	if u.user.DisplayName == nil || *u.user.DisplayName == "" {
		return u.user, ErrUserMissingProperty{"user needs a DisplayName"}
	}
	// if u.user.FirebaseInstanceID == "" {
	//     return u.user, ErrUserMissingProperty{"user needs a firebase instance ID"}
	// }
	*u.user.DisplayName = strings.TrimSpace(*u.user.DisplayName)
	u.user.PhotoURL = strfmt.URI(GetUserImageURL(*u.user.UID))
	return u.user, nil
}

// SetUID sets the user's UID
func (u *UserBuilder) SetUID(uid *string) {
	u.user.UID = uid
}

// SetDisplayName sets the user's DisplayName
func (u *UserBuilder) SetDisplayName(name *string) {
	u.user.DisplayName = name
}

// SetEmail sets the user's Email
func (u *UserBuilder) SetEmail(email strfmt.Email) {
	u.user.Email = email
}

// SetFirebaseInstanceID sets the user's FirebaseInstanceID
func (u *UserBuilder) SetFirebaseInstanceID(id string) {
	u.user.FirebaseInstanceID = id
}

// SetGroupUID sets the user's GroupUID
func (u *UserBuilder) SetGroupUID(id strfmt.UUID) {
	u.user.GroupUID = id
}

// SetLocale sets the user's Locale
func (u *UserBuilder) SetLocale(loc string) {
	u.user.Locale = loc
}

// SetPhotoURL sets the user's PhotoURL
func (u *UserBuilder) SetPhotoURL(uri strfmt.URI) {
	u.user.PhotoURL = uri
}
