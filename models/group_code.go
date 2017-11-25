package models

import (
	"time"

	apiErrors "github.com/go-openapi/errors"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
	"github.com/go-openapi/validate"
	"github.com/wgplaner/wg_planer_server/modules/base"
)

const (
	GROUP_CODE_LENGTH     = 12
	GROUP_CODE_VALID_DAYS = 3
)

// GroupCode group code
// swagger:model GroupCode
type GroupCode struct {
	// code
	// Required: true
	// Pattern: ^[A-Z0-9]{12}$
	Code *string `xorm:"pk" json:"code"`

	// group Uid
	// Required: true
	GroupUID *strfmt.UUID `json:"groupUID"`

	// valid until
	// Required: true
	// Read Only: true
	ValidUntil strfmt.DateTime `json:"validUntil"`
}

// Validate validates this group code
func (m *GroupCode) Validate(formats strfmt.Registry) error {
	var res []error
	if err := m.validateCode(formats); err != nil {
		// prop
		res = append(res, err)
	}
	if err := m.validateGroupUID(formats); err != nil {
		// prop
		res = append(res, err)
	}
	if err := m.validateValidUntil(formats); err != nil {
		// prop
		res = append(res, err)
	}
	if len(res) > 0 {
		return apiErrors.CompositeValidationError(res...)
	}
	return nil
}

func (m *GroupCode) validateCode(formats strfmt.Registry) error {
	if err := validate.Required("code", "body", m.Code); err != nil {
		return err
	}
	if err := validate.Pattern("code", "body", string(*m.Code), `^[A-Z0-9]{12}$`); err != nil {
		return err
	}
	return nil
}

func (m *GroupCode) validateGroupUID(formats strfmt.Registry) error {
	if err := validate.Required("groupUID", "body", m.GroupUID); err != nil {
		return err
	}
	if err := validate.FormatOf("groupUID", "body", "uuid", m.GroupUID.String(), formats); err != nil {
		return err
	}
	return nil
}

func (m *GroupCode) validateValidUntil(formats strfmt.Registry) error {
	if err := validate.Required("validUntil", "body", strfmt.DateTime(m.ValidUntil)); err != nil {
		return err
	}
	if err := validate.FormatOf("validUntil", "body", "date-time", m.ValidUntil.String(), formats); err != nil {
		return err
	}
	return nil
}

// MarshalBinary interface implementation
func (m *GroupCode) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *GroupCode) UnmarshalBinary(b []byte) error {
	var res GroupCode
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}

func CreateGroupCode(uid strfmt.UUID) (*GroupCode, error) {
	// Invalidate old codes
	oldCodes := GroupCode{ValidUntil: strfmt.DateTime(time.Now().UTC())}
	if _, err := x.Where("group_uid = ?", uid).Update(&oldCodes); err != nil {
		groupLog.Errorf(`Can't update validUntil date of other (old) codes for group "%s"'; Err: "%s"`, uid, err)
		return nil, err
	}

	var (
		code       = base.GetRandomAlphaNumCode(GROUP_CODE_LENGTH, true)
		validUntil = time.Now().UTC().AddDate(0, 0, GROUP_CODE_VALID_DAYS)
	)

	groupCode := &GroupCode{
		GroupUID:   &uid,
		Code:       &code,
		ValidUntil: strfmt.DateTime(validUntil),
	}

	// Insert new code into database
	if _, err := x.InsertOne(groupCode); err != nil {
		return nil, err
	}

	return groupCode, nil
}

func GetGroupCode(c string) (*GroupCode, error) {
	code := new(GroupCode)
	if has, err := x.ID(c).Get(code); err != nil {
		return nil, err
	} else if !has {
		return nil, ErrGroupCodeNotExist{Code: c}
	}
	return code, nil
}

func IsGroupCodeValid(c string) (bool, *GroupCode) {
	if len(c) != 12 {
		return false, nil
	}

	var err error
	code := new(GroupCode)

	if code, err = GetGroupCode(c); err != nil || code == nil {
		return false, nil
	} else if time.Now().After(time.Time(code.ValidUntil)) {
		return false, nil
	}
	return true, code
}
