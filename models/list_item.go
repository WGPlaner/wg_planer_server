package models

import (
	"time"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
	"github.com/go-openapi/validate"
)

type ListItem struct {
	// bill UID
	// Read Only: true
	BillUID strfmt.UUID `json:"billUID,omitempty"`

	// bought at
	// Read Only: true
	BoughtAt *time.Time `xorm:"NULL" json:"boughtAt,omitempty"`

	// category
	// Required: true
	Category *string `json:"category"`

	// count
	// Required: true
	Count *int64 `xorm:"DEFAULT 0" json:"count"`

	// group UID
	// Read Only: true
	GroupUID strfmt.UUID `xorm:"index(uid) unique(uid)" json:"groupUID,omitempty"`

	// id
	// Read Only: true
	ID strfmt.UUID `xorm:"index(uid) unique(uid)" json:"id,omitempty"`

	// price
	Price int64 `xorm:"DEFAULT 0" json:"price,omitempty"`

	// requested by
	// Read Only: true
	RequestedBy string `xorm:"NOT NULL" json:"requestedBy,omitempty"`

	// bought by
	// Read Only: true
	BoughtBy string `xorm:"NULL" json:"boughtBy,omitempty"`

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

// AfterLoad is invoked from XORM after setting the values of all fields of this object.
func (l *ListItem) AfterLoad() {
	// Clear date if it is before 2000
	if l.BoughtAt != nil && time.Date(2000, 1, 0, 0, 0, 0, 0, time.Local).After(*l.BoughtAt) {
		l.BoughtAt = nil
	}
}

// Validate validates this list item
func (l *ListItem) Validate(formats strfmt.Registry) error {
	var res []error
	if err := l.validateCategory(formats); err != nil {
		// prop
		res = append(res, err)
	}
	if err := l.validateCount(formats); err != nil {
		// prop
		res = append(res, err)
	}
	if err := l.validateRequestedFor(formats); err != nil {
		// prop
		res = append(res, err)
	}
	if err := l.validateTitle(formats); err != nil {
		// prop
		res = append(res, err)
	}
	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (l *ListItem) validateCategory(formats strfmt.Registry) error {
	if err := validate.Required("category", "body", l.Category); err != nil {
		return err
	}
	return nil
}

func (l *ListItem) validateCount(formats strfmt.Registry) error {
	if err := validate.Required("count", "body", l.Count); err != nil {
		return err
	}
	return nil
}

func (l *ListItem) validateRequestedFor(formats strfmt.Registry) error {
	if swag.IsZero(l.RequestedFor) { // not required
		return nil
	}
	return nil
}

func (l *ListItem) validateTitle(formats strfmt.Registry) error {
	if err := validate.Required("title", "body", l.Title); err != nil {
		return err
	}
	if err := validate.MaxLength("title", "body", string(*l.Title), 150); err != nil {
		return err
	}
	return nil
}

// MarshalBinary interface implementation
func (l *ListItem) MarshalBinary() ([]byte, error) {
	if l == nil {
		return nil, nil
	}
	return swag.WriteJSON(l)
}

// UnmarshalBinary interface implementation
func (l *ListItem) UnmarshalBinary(b []byte) error {
	var res ListItem
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*l = res
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
	_, err := x.Cols(cols...).
		Where(`group_uid=?`, l.GroupUID).
		And(`id=?`, l.ID).
		Update(l)
	return err
}
