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
		if html, ok := data.(string); !ok {
			return errors.New("error in HTML producer")
		} else {
			_, err := writer.Write([]byte(html))
			return err
		}
	})

	// Authentication
	api.UserIDAuthAuth = UserIDAuth
	api.FirebaseIDAuthAuth = FirebaseIDAuth

	// Create API handlers
	api.InfoGetVersionHandler = info.GetVersionHandlerFunc(GetVersionInfo)

	api.BillCreateBillHandler = bill.CreateBillHandlerFunc(CreateBill)
	api.BillGetBillListHandler = bill.GetBillListHandlerFunc(GetBillList)

	api.GroupCreateGroupHandler = group.CreateGroupHandlerFunc(CreateGroup)
	api.GroupCreateGroupCodeHandler = group.CreateGroupCodeHandlerFunc(CreateGroupCode)
	api.GroupGetGroupHandler = group.GetGroupHandlerFunc(GetGroup)
	api.GroupGetGroupImageHandler = group.GetGroupImageHandlerFunc(GetGroupImage)
	api.GroupUpdateGroupHandler = group.UpdateGroupHandlerFunc(UpdateGroup)
	api.GroupUpdateGroupImageHandler = group.UpdateGroupImageHandlerFunc(UpdateGroupImage)
	api.GroupJoinGroupHandler = group.JoinGroupHandlerFunc(JoinGroup)
	api.GroupJoinGroupHelpHandler = group.JoinGroupHelpHandlerFunc(JoinGroupHelp)
	api.GroupLeaveGroupHandler = group.LeaveGroupHandlerFunc(LeaveGroup)

	api.UserCreateUserHandler = user.CreateUserHandlerFunc(CreateUser)
	api.UserGetUserHandler = user.GetUserHandlerFunc(GetUser)
	api.UserGetUserImageHandler = user.GetUserImageHandlerFunc(GetUserImage)
	api.UserUpdateUserHandler = user.UpdateUserHandlerFunc(UpdateUser)
	api.UserUpdateUserImageHandler = user.UpdateUserImageHandlerFunc(UpdateUserImage)

	api.ShoppinglistCreateListItemHandler = shoppinglist.CreateListItemHandlerFunc(CreateListItem)
	api.ShoppinglistGetListItemsHandler = shoppinglist.GetListItemsHandlerFunc(GetListItems)
	api.ShoppinglistUpdateListItemHandler = shoppinglist.UpdateListItemHandlerFunc(UpdateListItem)
	api.ShoppinglistBuyListItemsHandler = shoppinglist.BuyListItemsHandlerFunc(BuyListItems)
}
