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
	BillItems []string `xorm:"-" json:"billItems"`

	// created by
	// Required: true
	CreatedBy *string `xorm:"VARCHAR(28)" json:"createdBy"`

	// due date
	DueDate string `json:"dueDate,omitempty"`

	// group uid
	// Read Only: true
	GroupUID strfmt.UUID `json:"groupUID,omitempty"`

	// state
	// Required: true
	State *string `json:"state"`

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
	if err := m.validateBillItems(formats); err != nil {
		// prop
		res = append(res, err)
	}
	if err := m.validateCreatedBy(formats); err != nil {
		// prop
		res = append(res, err)
	}
	if err := m.validateState(formats); err != nil {
		// prop
		res = append(res, err)
	}
	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *Bill) validateBillItems(formats strfmt.Registry) error {
	if err := validate.Required("billItems", "body", m.BillItems); err != nil {
		return err
	}
	return nil
}

func (m *Bill) validateCreatedBy(formats strfmt.Registry) error {
	if err := validate.Required("createdBy", "body", m.CreatedBy); err != nil {
		return err
	}
	return nil
}

func (m *Bill) validateState(formats strfmt.Registry) error {
	if err := validate.Required("state", "body", m.State); err != nil {
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

func CreateBillForGroup(g *Group, u *User) (*Bill, error) {
	billUid, err := uuid.NewV4()
	if err != nil {
		return nil, err
	}

	b := &Bill{
		UID:       strfmt.UUID(billUid.String()),
		State:     swag.String("TODO"),
		CreatedBy: u.UID,
		GroupUID:  g.UID,
		// TODO: Other fields
	}

	if _, err := x.InsertOne(b); err != nil {
		return nil, err
	}

	_, err = x.
		Cols(`bill_uid`).
		Where(`group_uid=?`, g.UID).
		And(`bill_uid IS NULL`).
		Update(ListItem{BillUID: b.UID})

	if err != nil {
		return nil, err
	}

	return b, err
}
