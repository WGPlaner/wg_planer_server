package models

import (
	"github.com/go-openapi/strfmt"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/swag"
	"github.com/go-openapi/validate"
)

// User user
// swagger:model User
type User struct {

	// created at
	// Read Only: true
	CreatedAt strfmt.DateTime `json:"createdAt,omitempty"`

	// display name
	// Required: true
	// Max Length: 20
	// Min Length: 3
	DisplayName *string `json:"displayName"`

	// email
	Email strfmt.Email `json:"email,omitempty"`

	// firebase instance Id
	// Pattern: ^[-_:a-zA-Z0-9]{152}$
	FirebaseInstanceID string `json:"firebaseInstanceId,omitempty"`

	// group Uid
	GroupUID strfmt.UUID `json:"groupUid,omitempty"`

	// locale
	Locale string `json:"locale,omitempty"`

	// photo Url
	PhotoURL strfmt.URI `json:"photoUrl,omitempty"`

	// uid
	// Required: true
	UID *string `json:"uid"`

	// updated at
	// Read Only: true
	UpdatedAt strfmt.DateTime `json:"updatedAt,omitempty"`
}

// Validate validates this user
func (m *User) Validate(formats strfmt.Registry) error {
	var res []error
	if err := m.validateDisplayName(formats); err != nil {
		res = append(res, err)
	}
	if err := m.validateFirebaseInstanceID(formats); err != nil {
		res = append(res, err)
	}
	if err := m.validateUID(formats); err != nil {
		res = append(res, err)
	}
	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *User) validateDisplayName(formats strfmt.Registry) error {
	if err := validate.Required("displayName", "body", m.DisplayName); err != nil {
		return err
	}
	if err := validate.MinLength("displayName", "body", string(*m.DisplayName), 3); err != nil {
		return err
	}
	if err := validate.MaxLength("displayName", "body", string(*m.DisplayName), 20); err != nil {
		return err
	}
	return nil
}

func (m *User) validateFirebaseInstanceID(formats strfmt.Registry) error {
	if swag.IsZero(m.FirebaseInstanceID) { // not required
		return nil
	}
	if err := validate.Pattern("firebaseInstanceId", "body", string(m.FirebaseInstanceID), `^[-_:a-zA-Z0-9]{152}$`); err != nil {
		return err
	}
	return nil
}

func (m *User) validateUID(formats strfmt.Registry) error {
	if err := validate.Required("uid", "body", m.UID); err != nil {
		return err
	}
	return nil
}

// MarshalBinary interface implementation
func (m *User) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *User) UnmarshalBinary(b []byte) error {
	var res User
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
