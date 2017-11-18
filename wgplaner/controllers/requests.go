package controllers

import (
	"net/http"

	"github.com/wgplaner/wg_planer_server/gen/models"

	"github.com/go-openapi/runtime"
	"github.com/go-openapi/swag"
)

// Error HTTP Responder

//  _  _      ___     ___
// | || |    / _ \   / _ \
// | || |_  | | | | | | | |
// |__   _| | | | | | | | |
//    | |   | |_| | | |_| |
//    |_|    \___/   \___/
//

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

//  _  _      ___    __
// | || |    / _ \  /_ |
// | || |_  | | | |  | |
// |__   _| | | | |  | |
//    | |   | |_| |  | |
//    |_|    \___/   |_|
//

type UnauthorizedReponse struct {
	Payload *models.ErrorResponse `json:"body,omitempty"`
}

func NewUnauthorizedResponse(msg string) *UnauthorizedReponse {
	return &UnauthorizedReponse{
		Payload: &models.ErrorResponse{
			Message: swag.String(msg),
			Status:  swag.Int64(http.StatusNotFound),
		},
	}
}

func (o *UnauthorizedReponse) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {
	rw.WriteHeader(http.StatusUnauthorized)
	payload := o.Payload

	if payload == nil {
		payload = &models.ErrorResponse{
			Message: swag.String("Unauthorized"),
			Status:  swag.Int64(http.StatusUnauthorized),
		}
	}

	if err := producer.Produce(rw, payload); err != nil {
		panic(err) // let the recovery middleware deal with this
	}
}

//  _  _      ___    _  _
// | || |    / _ \  | || |
// | || |_  | | | | | || |_
// |__   _| | | | | |__   _|
//    | |   | |_| |    | |
//    |_|    \___/     |_|
//

type NotFoundResponse struct {
	Payload *models.ErrorResponse `json:"body,omitempty"`
}

func NewNotFoundResponse(msg string) *NotFoundResponse {
	return &NotFoundResponse{
		Payload: &models.ErrorResponse{
			Message: swag.String(msg),
			Status:  swag.Int64(http.StatusNotFound),
		},
	}
}

func (o *NotFoundResponse) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {
	rw.WriteHeader(http.StatusNotFound)
	payload := o.Payload

	if payload == nil {
		payload = &models.ErrorResponse{
			Message: swag.String("Not Found"),
			Status:  swag.Int64(http.StatusNotFound),
		}
	}

	if err := producer.Produce(rw, payload); err != nil {
		panic(err) // let the recovery middleware deal with this
	}
}

//
//  _____    ___     ___
// | ____|  / _ \   / _ \
// | |__   | | | | | | | |
// |___ \  | | | | | | | |
//  ___) | | |_| | | |_| |
// |____/   \___/   \___/
//

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
