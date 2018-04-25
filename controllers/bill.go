package controllers

import (
	"github.com/wgplaner/wg_planer_server/models"
	"github.com/wgplaner/wg_planer_server/restapi/operations/bill"

	"github.com/go-openapi/runtime/middleware"
	"github.com/op/go-logging"
)

var billLog = logging.MustGetLogger("Bill")

// getBillList returns a list of bills for the requested group.
func getBillList(params bill.GetBillListParams, principal *models.User) middleware.Responder {
	groupLog.Debugf(`User %q gets bills for group "%s"`, *principal.UID, params.GroupUID)

	var g *models.Group
	var errResp middleware.Responder

	if g, errResp = getGroupAuthorizedOrError(params.GroupUID, *principal.UID); errResp != nil {
		return errResp
	}

	bills, err := models.GetBillsByGroupUIDWithBillItems(g.UID)
	if err != nil {
		return newInternalServerError("Internal Server Error")
	}

	// TODO: Check authorization, etc

	billList := &models.BillList{
		Bills: bills,
		Count: int64(len(bills)),
	}

	return bill.NewGetBillListOK().WithPayload(billList)
}

// createBill creates a bill for the requested group.
func createBill(params bill.CreateBillParams, principal *models.User) middleware.Responder {
	billLog.Debugf(`Start creating bill for group "%s"`, params.GroupUID)

	// TODO: Check authorization, etc

	var g *models.Group
	var errResp middleware.Responder

	if g, errResp = getGroupAuthorizedOrError(params.GroupUID, *principal.UID); errResp != nil {
		return errResp
	}

	b, err := models.CreateBillForGroup(g, principal)
	if err != nil {
		return newInternalServerError("Internal Server Error")
	}

	return bill.NewCreateBillOK().WithPayload(b)
}
