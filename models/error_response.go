package models

import (
	"github.com/go-openapi/errors"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
	"github.com/go-openapi/validate"
)

// ErrorResponse error response
// swagger:model ErrorResponse
type ErrorResponse struct {
	// message
	// Required: true
	Message *string `json:"message"`

	// status
	// Required: true
	Status *int64 `json:"status"`
}

// Validate validates this error response
func (r *ErrorResponse) Validate(formats strfmt.Registry) error {
	var res []error
	if err := r.validateMessage(formats); err != nil {
		// prop
		res = append(res, err)
	}
	if err := r.validateStatus(formats); err != nil {
		// prop
		res = append(res, err)
	}
	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (r *ErrorResponse) validateMessage(formats strfmt.Registry) error {
	if err := validate.Required("message", "body", r.Message); err != nil {
		return err
	}
	return nil
}

func (r *ErrorResponse) validateStatus(formats strfmt.Registry) error {
	if err := validate.Required("status", "body", r.Status); err != nil {
		return err
	}
	return nil
}

// MarshalBinary interface implementation
func (r *ErrorResponse) MarshalBinary() ([]byte, error) {
	if r == nil {
		return nil, nil
	}
	return swag.WriteJSON(r)
}

// UnmarshalBinary interface implementation
func (r *ErrorResponse) UnmarshalBinary(b []byte) error {
	var res ErrorResponse
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*r = res
	return nil
}
