package controllers

import (
	"github.com/wgplaner/wg_planer_server/models"
	"github.com/wgplaner/wg_planer_server/restapi/operations/bill"

	"github.com/go-openapi/runtime/middleware"
	"github.com/op/go-logging"
)

var billLog = logging.MustGetLogger("Bill")

func GetBillList(params bill.GetBillListParams, principal *models.User) middleware.Responder {
	bills, err := models.GetBillsByGroupUID(params.GroupUID)
	if err != nil {
		return NewInternalServerError("Internal Server Error")
	}

	// TODO: Check authorization, etc

	billList := &models.BillList{
		Bills: bills,
		Count: int64(len(bills)),
	}

	return bill.NewGetBillListOK().WithPayload(billList)
}

func CreateBill(params bill.CreateBillParams, principal *models.User) middleware.Responder {
	billLog.Debugf(`Start creating bill for group "%s"`, params.GroupUID)

	// TODO: Check authorization, etc

	g := &models.Group{
		UID: params.GroupUID,
	}

	b, err := models.CreateBillForGroup(g, principal)
	if err != nil {
		return NewInternalServerError("Internal Server Error")
	}

	return bill.NewCreateBillOK().WithPayload(b)
}
