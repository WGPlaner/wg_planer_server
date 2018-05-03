package models

import (
	"github.com/go-openapi/errors"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
	"github.com/go-openapi/validate"
)

// BillList bill list
// swagger:model BillList
type BillList struct {
	// bills
	// Required: true
	Bills []*Bill `json:"bills"`

	// count
	// Required: true
	// Read Only: true
	Count int64 `json:"count"`
}

// Validate validates this bill list
func (m *BillList) Validate(formats strfmt.Registry) error {
	var res []error
	if err := m.validateBills(formats); err != nil {
		// prop
		res = append(res, err)
	}
	if err := m.validateCount(formats); err != nil {
		// prop
		res = append(res, err)
	}
	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *BillList) validateBills(formats strfmt.Registry) error {
	return validate.Required("bills", "body", m.Bills)
}

func (m *BillList) validateCount(formats strfmt.Registry) error {
	return validate.Required("count", "body", int64(m.Count))
}

// MarshalBinary interface implementation
func (m *BillList) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *BillList) UnmarshalBinary(b []byte) error {
	var res BillList
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
