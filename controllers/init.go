package controllers

import (
	"errors"
	"io"

	"github.com/wgplaner/wg_planer_server/models"
	"github.com/wgplaner/wg_planer_server/modules/setting"
	"github.com/wgplaner/wg_planer_server/restapi/operations"
	"github.com/wgplaner/wg_planer_server/restapi/operations/bill"
	"github.com/wgplaner/wg_planer_server/restapi/operations/group"
	"github.com/wgplaner/wg_planer_server/restapi/operations/info"
	"github.com/wgplaner/wg_planer_server/restapi/operations/shoppinglist"
	"github.com/wgplaner/wg_planer_server/restapi/operations/user"

	"github.com/go-openapi/runtime"
	"github.com/op/go-logging"
)

var initLog = logging.MustGetLogger("Auth")

// GlobalInit is for global configuration reload-able.
func GlobalInit() {
	setting.NewConfigContext()
	initLog.Infof("AppPath: %s", setting.AppPath)
	initLog.Infof("AppWorkPath: %s", setting.AppWorkPath)

	if err := models.NewEngine(); err != nil {
		initLog.Fatalf("Failed to initialize ORM engine: %v", err)
	}
}

func InitializeControllers(api *operations.WgplanerAPI) {
	// Producers
	api.HTMLProducer = runtime.ProducerFunc(func(writer io.Writer, data interface{}) error {
		if html, ok := data.(string); ok {
			_, err := writer.Write([]byte(html))
			return err
		}
		return errors.New("error in HTML producer")
	})

	// Authentication
	api.UserIDAuthAuth = userIDAuth
	api.FirebaseIDAuthAuth = firebaseIDAuth

	// Create API handlers
	api.InfoGetVersionHandler = info.GetVersionHandlerFunc(getVersionInfo)

	api.BillCreateBillHandler = bill.CreateBillHandlerFunc(createBill)
	api.BillGetBillListHandler = bill.GetBillListHandlerFunc(getBillList)

	api.GroupCreateGroupHandler = group.CreateGroupHandlerFunc(createGroup)
	api.GroupCreateGroupCodeHandler = group.CreateGroupCodeHandlerFunc(createGroupCode)
	api.GroupGetGroupHandler = group.GetGroupHandlerFunc(getGroup)
	api.GroupGetGroupImageHandler = group.GetGroupImageHandlerFunc(getGroupImage)
	api.GroupUpdateGroupHandler = group.UpdateGroupHandlerFunc(updateGroup)
	api.GroupUpdateGroupImageHandler = group.UpdateGroupImageHandlerFunc(updateGroupImage)
	api.GroupJoinGroupHandler = group.JoinGroupHandlerFunc(joinGroup)
	api.GroupJoinGroupHelpHandler = group.JoinGroupHelpHandlerFunc(joinGroupHelp)
	api.GroupLeaveGroupHandler = group.LeaveGroupHandlerFunc(leaveGroup)

	api.UserCreateUserHandler = user.CreateUserHandlerFunc(createUser)
	api.UserGetUserHandler = user.GetUserHandlerFunc(getUser)
	api.UserGetUserImageHandler = user.GetUserImageHandlerFunc(getUserImage)
	api.UserUpdateUserHandler = user.UpdateUserHandlerFunc(updateUser)
	api.UserUpdateUserImageHandler = user.UpdateUserImageHandlerFunc(updateUserImage)

	api.ShoppinglistCreateListItemHandler = shoppinglist.CreateListItemHandlerFunc(createListItem)
	api.ShoppinglistGetListItemsHandler = shoppinglist.GetListItemsHandlerFunc(getListItems)
	api.ShoppinglistUpdateListItemHandler = shoppinglist.UpdateListItemHandlerFunc(updateListItem)
	api.ShoppinglistBuyListItemsHandler = shoppinglist.BuyListItemsHandlerFunc(buyListItems)
}
