package models

import (
	"bytes"
	"fmt"
	"image"
	"image/jpeg"
	"os"
	"path"
	"strings"

	"github.com/wgplaner/wg_planer_server/modules/avatar"
	"github.com/wgplaner/wg_planer_server/modules/base"
	"github.com/wgplaner/wg_planer_server/modules/setting"

	apiErrors "github.com/go-openapi/errors"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
	"github.com/go-openapi/validate"
	"github.com/nfnt/resize"
	"github.com/op/go-logging"
)

var groupLog = logging.MustGetLogger("Group")

const (
	GROUP_PROFILE_IMAGE_FILE_NAME = "group_image.jpg"
)

// Group group
// swagger:model Group
type Group struct {
	// admins
	// Read Only: true
	Admins []string `json:"admins"`

	// currency
	// Max Length: 4
	Currency string `xorm:"varchar(5) default '€'" json:"currency,omitempty"`

	// display name
	// Required: true
	DisplayName *string `json:"displayName"`

	// members
	// Read Only: true
	Members []string `xorm:"-" json:"members"`

	// photo Url
	PhotoURL strfmt.URI `xorm:"-" json:"photoUrl,omitempty"`

	// uid
	// Read Only: true
	UID strfmt.UUID `xorm:"varchar(36) pk" json:"uid,omitempty"`

	// created at
	// Read Only: true
	CreatedAt strfmt.DateTime `xorm:"created" json:"createdAt,omitempty"`

	// updated at
	// Read Only: true
	UpdatedAt strfmt.DateTime `xorm:"updated" json:"updatedAt,omitempty"`
}

// AfterLoad is invoked from XORM after setting the values of all fields of this object.
func (g *Group) AfterLoad() {
	members, err := GetGroupMemberUIDs(g.UID)
	if err == nil {
		g.Members = members
	} else {
		groupLog.Critical(`Error loading members`, err.Error())
	}
}

// Validate validates this group
func (g *Group) Validate(formats strfmt.Registry) error {
	var res []error
	if err := g.validateAdmins(formats); err != nil {
		// prop
		res = append(res, err)
	}
	if err := g.validateCurrency(formats); err != nil {
		// prop
		res = append(res, err)
	}
	if err := g.validateDisplayName(formats); err != nil {
		// prop
		res = append(res, err)
	}
	if err := g.validateMembers(formats); err != nil {
		// prop
		res = append(res, err)
	}
	if len(res) > 0 {
		return apiErrors.CompositeValidationError(res...)
	}
	return nil
}

func (g *Group) validateAdmins(formats strfmt.Registry) error {
	if swag.IsZero(g.Admins) { // not required
		return nil
	}
	return nil
}

func (g *Group) validateCurrency(formats strfmt.Registry) error {
	if swag.IsZero(g.Currency) { // not required
		return nil
	}
	if err := validate.MaxLength("currency", "body", string(g.Currency), 4); err != nil {
		return err
	}
	return nil
}

func (g *Group) validateDisplayName(formats strfmt.Registry) error {
	if err := validate.Required("displayName", "body", g.DisplayName); err != nil {
		return err
	}
	return nil
}

func (g *Group) validateMembers(formats strfmt.Registry) error {
	if swag.IsZero(g.Members) { // not required
		return nil
	}
	return nil
}

// MarshalBinary interface implementation
func (g *Group) MarshalBinary() ([]byte, error) {
	if g == nil {
		return nil, nil
	}
	return swag.WriteJSON(g)
}

// UnmarshalBinary interface implementation
func (g *Group) UnmarshalBinary(b []byte) error {
	var res Group
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*g = res
	return nil
}

func IsGroupExist(uid strfmt.UUID) (bool, error) {
	if !strfmt.IsUUID(string(uid)) {
		return false, ErrGroupInvalidUUID{}
	}
	return x.Exist(&Group{UID: uid})
}

func (g *Group) HasMember(uid string) bool {
	has, _ := x.
		Where("group_uid=?", g.UID).
		And("uid=?", uid).
		Exist(new(User))
	return has
}

func (g *Group) HasAdmin(uid string) bool {
	return base.StringInSlice(uid, g.Admins)
}

func (g *Group) GetActiveShoppingListItems() ([]*ListItem, error) {
	items := make([]*ListItem, 0, 10)
	err := x.
		AllCols().
		Where(`group_uid=?`, g.UID).
		And(`bought_at IS NULL`).
		Find(&items)
	return items, err
}

func GetGroupMemberUIDs(groupUID strfmt.UUID) ([]string, error) {
	users := make([]*User, 0, 5)
	userUids := make([]string, 0, 5)
	err := x.
		Cols(`uid`).
		Where("group_uid=?", groupUID).
		Find(&users)

	for _, u := range users {
		userUids = append(userUids, *u.UID)
	}
	return userUids, err
}

func CreateGroup(g *Group) error {
	g.DisplayName = swag.String(strings.TrimSpace(swag.StringValue(g.DisplayName)))
	g.Currency = strings.TrimSpace(g.Currency)

	_, err := x.InsertOne(g)
	return err
}

func UpdateGroup(g *Group) error {
	g.DisplayName = swag.String(strings.TrimSpace(swag.StringValue(g.DisplayName)))
	g.Currency = strings.TrimSpace(g.Currency)

	_, err := x.ID(g.UID).AllCols().Update(g)
	return err
}

func UpdateGroupCols(g *Group, cols ...string) error {
	g.DisplayName = swag.String(strings.TrimSpace(swag.StringValue(g.DisplayName)))
	g.Currency = strings.TrimSpace(g.Currency)

	_, err := x.ID(g.UID).Cols(cols...).Update(g)
	return err
}

func GetGroupByUID(uid strfmt.UUID) (*Group, error) {
	if !strfmt.IsUUID(string(uid)) {
		return nil, ErrGroupInvalidUUID{UID: string(uid)}
	}

	g := new(Group)

	if has, err := x.ID(uid).Get(g); err != nil {
		return nil, err

	} else if !has {
		return nil, ErrGroupNotExist{UID: uid}
	}

	return g, nil
}

func GetGroupImagePath(uid strfmt.UUID) string {
	return path.Join(
		setting.AppWorkPath,
		setting.AppConfig.Data.GroupImageDir,
		string(uid),
		GROUP_PROFILE_IMAGE_FILE_NAME,
	)
}

func GetGroupImage(uid strfmt.UUID) (*os.File, error) {
	return os.Open(GetGroupImagePath(uid))
}

func GetGroupImageDefault() (*os.File, error) {
	return os.Open(path.Join(
		setting.AppWorkPath,
		setting.AppConfig.Data.UserImageDefault,
	))
}

func GetGroupImageURL(uid string) string {
	return strings.Replace(`/users/{userID}/image`, `{userID}`, uid, -1)
}

func (g *Group) UploadGroupImage(data []byte) error {
	filePath := GetGroupImagePath(g.UID)
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
