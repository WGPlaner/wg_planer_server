package controllers

import (
	"net/http"

	"github.com/wgplaner/wg_planer_server/gen/models"

	"github.com/go-openapi/runtime"
	"github.com/go-openapi/swag"
)

type BadRequest struct {
	Payload *models.ErrorResponse `json:"body,omitempty"`
}

func NewBadRequest(msg string) *BadRequest {
	return &BadRequest{
		Payload: &models.ErrorResponse{
			Message: swag.String(msg),
			Status:  swag.Int64(http.StatusBadRequest),
		},
	}
}

func (o *BadRequest) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {
	rw.WriteHeader(http.StatusBadRequest)
	payload := o.Payload

	if payload == nil {
		payload = &models.ErrorResponse{
			Message: swag.String("Bad Request"),
			Status:  swag.Int64(http.StatusBadRequest),
		}
	}

	if err := producer.Produce(rw, payload); err != nil {
		panic(err) // let the recovery middleware deal with this
	}
}

type InternalServerError struct {
	Payload *models.ErrorResponse `json:"body,omitempty"`
}

func NewInternalServerError(msg string) *InternalServerError {
	return &InternalServerError{
		Payload: &models.ErrorResponse{
			Message: swag.String(msg),
			Status:  swag.Int64(http.StatusInternalServerError),
		},
	}
}

func (o *InternalServerError) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {
	rw.WriteHeader(http.StatusInternalServerError)
	payload := o.Payload

	if payload == nil {
		payload = &models.ErrorResponse{
			Message: swag.String("Internal Server Error"),
			Status:  swag.Int64(http.StatusInternalServerError),
		}
	}

	if err := producer.Produce(rw, payload); err != nil {
		panic(err) // let the recovery middleware deal with this
	}
}
