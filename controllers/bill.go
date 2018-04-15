package controllers

import (
	"github.com/wgplaner/wg_planer_server/models"
	"github.com/wgplaner/wg_planer_server/restapi/operations/bill"

	"github.com/go-openapi/runtime/middleware"
	"github.com/op/go-logging"
)

var billLog = logging.MustGetLogger("Bill")

func GetBillList(params bill.GetBillListParams, principal *models.User) middleware.Responder {
	groupLog.Debugf(`User %q gets bills for group "%s"`, *principal.UID, params.GroupUID)

	var g *models.Group
	var err error

	// Database - Check group
	if g, err = models.GetGroupByUID(params.GroupUID); models.IsErrGroupNotExist(err) {
		groupLog.Debugf(`Can't find database group with id "%s"!`, params.GroupUID)
		return NewNotFoundResponse("Group not found on server")
	}
	if models.IsErrGroupInvalidUUID(err) {
		groupLog.Debugf(err.Error())
		return NewNotFoundResponse("invalid group uid")
	}
	if err != nil {
		groupLog.Critical(`Database Error!`, err)
		return NewInternalServerError("Internal Database Error")
	}
	// Check if group has member
	if !g.HasMember(*principal.UID) {
		return NewUnauthorizedResponse("User is not a member of the specified group")
	}

	bills, err := models.GetBillsByGroupUIDWithBillItems(params.GroupUID)
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
