package models

import (
	"github.com/go-openapi/errors"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
	"github.com/go-openapi/validate"
)

type ListItem struct {
	// bought at
	// Read Only: true
	BoughtAt strfmt.DateTime `json:"boughtAt,omitempty"`

	// category
	// Required: true
	Category *string `json:"category"`

	// count
	// Required: true
	Count *int64 `xorm:"DEFAULT 0" json:"count"`

	// group Uid
	// Read Only: true
	GroupUID strfmt.UUID `xorm:"index(uid) unique(uid)" json:"groupUid,omitempty"`

	// id
	// Read Only: true
	ID strfmt.UUID `xorm:"index(uid) unique(uid)" json:"id,omitempty"`

	// price
	Price int64 `xorm:"DEFAULT 0" json:"price,omitempty"`

	// requested by
	// Read Only: true
	RequestedBy string `xorm:"NOT NULL" json:"requestedBy,omitempty"`

	// requested for
	RequestedFor []string `xorm:"NOT NULL" json:"requestedFor"`

	// title
	// Required: true
	// Max Length: 150
	Title *string `xorm:"NOT NULL" json:"title"`

	// created at
	// Read Only: true
	CreatedAt strfmt.DateTime `xorm:"created" json:"createdAt,omitempty"`

	// updated at
	// Read Only: true
	UpdatedAt strfmt.DateTime `xorm:"updated" json:"updatedAt,omitempty"`
}

// Validate validates this list item
func (m *ListItem) Validate(formats strfmt.Registry) error {
	var res []error
	if err := m.validateCategory(formats); err != nil {
		// prop
		res = append(res, err)
	}
	if err := m.validateCount(formats); err != nil {
		// prop
		res = append(res, err)
	}
	if err := m.validateRequestedFor(formats); err != nil {
		// prop
		res = append(res, err)
	}
	if err := m.validateTitle(formats); err != nil {
		// prop
		res = append(res, err)
	}
	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *ListItem) validateCategory(formats strfmt.Registry) error {
	if err := validate.Required("category", "body", m.Category); err != nil {
		return err
	}
	return nil
}

func (m *ListItem) validateCount(formats strfmt.Registry) error {
	if err := validate.Required("count", "body", m.Count); err != nil {
		return err
	}
	return nil
}

func (m *ListItem) validateRequestedFor(formats strfmt.Registry) error {
	if swag.IsZero(m.RequestedFor) { // not required
		return nil
	}
	return nil
}

func (m *ListItem) validateTitle(formats strfmt.Registry) error {
	if err := validate.Required("title", "body", m.Title); err != nil {
		return err
	}
	if err := validate.MaxLength("title", "body", string(*m.Title), 150); err != nil {
		return err
	}
	return nil
}

// MarshalBinary interface implementation
func (m *ListItem) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *ListItem) UnmarshalBinary(b []byte) error {
	var res ListItem
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}

func GetListItemByUIDs(guid, luid strfmt.UUID) (*ListItem, error) {
	l := &ListItem{
		GroupUID: guid,
		ID:       luid,
	}

	if has, err := x.Get(l); err != nil {
		return nil, err

	} else if !has {
		return nil, ErrListItemNotExist{GroupUID: guid, ID: luid}
	}

	return l, nil
}

func CreateListItem(item *ListItem) error {
	_, err := x.InsertOne(item)
	return err
}

func UpdateListItem(l *ListItem) error {
	_, err := x.AllCols().Update(l)
	return err
}

func UpdateListItemCols(l *ListItem, cols ...string) error {
	_, err := x.Cols(cols...).Update(l)
	return err
}
