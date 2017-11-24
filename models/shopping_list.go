package models

import (
	"strconv"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
	"github.com/go-openapi/validate"
)

// ShoppingList shopping list
// swagger:model ShoppingList
type ShoppingList struct {
	// count
	// Required: true
	// Read Only: true
	Count int64 `json:"count"`

	// list items
	// Required: true
	// Read Only: true
	ListItems []*ListItem `json:"listItems"`
}

// Validate validates this shopping list
func (m *ShoppingList) Validate(formats strfmt.Registry) error {
	var res []error
	if err := m.validateCount(formats); err != nil {
		// prop
		res = append(res, err)
	}
	if err := m.validateListItems(formats); err != nil {
		// prop
		res = append(res, err)
	}
	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *ShoppingList) validateCount(formats strfmt.Registry) error {
	if err := validate.Required("count", "body", int64(m.Count)); err != nil {
		return err
	}
	return nil
}

func (m *ShoppingList) validateListItems(formats strfmt.Registry) error {
	if err := validate.Required("listItems", "body", m.ListItems); err != nil {
		return err
	}
	for i := 0; i < len(m.ListItems); i++ {
		if swag.IsZero(m.ListItems[i]) { // not required
			continue
		}
		if m.ListItems[i] != nil {
			if err := m.ListItems[i].Validate(formats); err != nil {
				if ve, ok := err.(*errors.Validation); ok {
					return ve.ValidateName("listItems" + "." + strconv.Itoa(i))
				}
				return err
			}
		}
	}
	return nil
}

// MarshalBinary interface implementation
func (m *ShoppingList) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *ShoppingList) UnmarshalBinary(b []byte) error {
	var res ShoppingList
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
