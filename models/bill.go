package models

import (
	"github.com/go-openapi/errors"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
	"github.com/go-openapi/validate"
	"github.com/satori/go.uuid"
)

// Bill bill
// swagger:model Bill

type Bill struct {
	// uid
	UID strfmt.UUID `xorm:"pk" json:"uid,omitempty"`

	// bill items
	// Required: true
	BoughtItems []string `xorm:"-" json:"boughtItems"`

	// created by
	// Required: true
	CreatedBy *string `xorm:"VARCHAR(28)" json:"createdBy,omitempty"`

	// sent to
	// Required: true
	SentTo []string `xorm:"VARCHAR(28)" json:"sentTo,omitempty"`

	// payed by
	PayedBy []string `xorm:"VARCHAR(28)" json:"payedBy,omitempty"`

	// due date
	DueDate string `json:"dueDate,omitempty"`

	// group uid
	// Read Only: true
	GroupUID strfmt.UUID `json:"groupUID,omitempty"`

	// state (can be e.g. paid, todo)
	// Required: true
	State *string `xorm:"VARCHAR(5)" json:"state"`

	// sum
	Sum int64 `xorm:"-" json:"sum,omitempty"`

	// created at
	// Read Only: true
	CreatedAt strfmt.DateTime `xorm:"created" json:"createdAt,omitempty"`

	// updated at
	// Read Only: true
	UpdatedAt strfmt.DateTime `xorm:"updated" json:"updatedAt,omitempty"`
}

// Validate validates this bill
func (m *Bill) Validate(formats strfmt.Registry) error {
	var res []error
	if err := m.validateBoughtItems(formats); err != nil {
		// prop
		res = append(res, err)
	}
	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *Bill) validateBoughtItems(formats strfmt.Registry) error {
	if err := validate.Required("boughtItems", "body", m.BoughtItems); err != nil {
		return err
	}
	return nil
}

// MarshalBinary interface implementation
func (m *Bill) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *Bill) UnmarshalBinary(b []byte) error {
	var res Bill
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}

func GetBillsByGroupUID(guid strfmt.UUID) ([]*Bill, error) {
	bills := make([]*Bill, 0, 5)

	err := x.AllCols().Where(`group_uid=?`, guid).Find(&bills)
	if err != nil {
		return nil, err
	}

	return bills, nil
}

// GetBillsByGroupUIDWithBoughtItems get bills with bought items for given group
func GetBillsByGroupUIDWithBoughtItems(guid strfmt.UUID) ([]*Bill, error) {
	bills, err := GetBillsByGroupUID(guid)
	if err != nil {
		return nil, err
	}

	// Get items for each bill
	for _, b := range bills {
		var boughtItems []ListItem
		err = x.Cols("id", "price", "count").Where(`bill_uid=?`, b.UID).Find(&boughtItems)
		if err != nil {
			return nil, err
		}
		b.Sum = 0
		for _, i := range boughtItems {
			b.BoughtItems = append(b.BoughtItems, string(i.ID))
			if i.Count != nil {
				b.Sum += i.Price * *i.Count
			}
		}
	}

	return bills, nil
}

// CreateBillForUser create a bill for a user
func CreateBillForUser(u *User, billWithItems *Bill) (*Bill, error) {
	billUID, err := uuid.NewV4()
	if err != nil {
		return nil, err
	}

	b := &Bill{
		UID:       strfmt.UUID(billUID.String()),
		CreatedBy: u.UID,
		GroupUID:  u.GroupUID,
		SentTo:    []string{},
		PayedBy:   []string{},
		DueDate:   billWithItems.DueDate,
		State:     swag.String("todo"),
		// TODO: Other fields
	}

	if _, err := x.InsertOne(b); err != nil {
		return nil, err
	}

	_, err = x.
		Cols(`bill_uid`).
		Where(`bill_uid IS NULL`).
		And(`bought_by = ?`, *u.UID).
		In(`id`, billWithItems.BoughtItems).
		Update(ListItem{BillUID: b.UID})

	if err != nil {
		return nil, err
	}

	// Get Items
	var items []ListItem
	err = x.Cols(`id`).Where(`bill_uid = ?`, b.UID).Find(&items)
	for _, item := range items {
		b.BoughtItems = append(b.BoughtItems, string(item.ID))
	}

	if err != nil {
		return nil, err
	}

	return b, err
}
