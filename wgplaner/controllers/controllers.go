package controllers

import (
	"errors"
	"io"

	"github.com/wgplaner/wg_planer_server/gen/restapi/operations"
	"github.com/wgplaner/wg_planer_server/gen/restapi/operations/group"
	"github.com/wgplaner/wg_planer_server/gen/restapi/operations/info"
	"github.com/wgplaner/wg_planer_server/gen/restapi/operations/shoppinglist"
	"github.com/wgplaner/wg_planer_server/gen/restapi/operations/user"

	"github.com/go-openapi/runtime"
)

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
	api.InfoGetLatestVersionHandler = info.GetLatestVersionHandlerFunc(GetVersionInfo)

	api.GroupCreateGroupHandler = group.CreateGroupHandlerFunc(CreateGroup)
	api.GroupGetGroupHandler = group.GetGroupHandlerFunc(GetGroup)
	api.GroupCreateGroupCodeHandler = group.CreateGroupCodeHandlerFunc(CreateGroupCode)
	api.GroupJoinGroupHandler = group.JoinGroupHandlerFunc(JoinGroup)
	api.GroupJoinGroupHelpHandler = group.JoinGroupHelpHandlerFunc(JoinGroupHelp)
	api.GroupLeaveGroupHandler = group.LeaveGroupHandlerFunc(LeaveGroup)
	api.GroupGetGroupImageHandler = group.GetGroupImageHandlerFunc(GetGroupImage)
	api.GroupUpdateGroupImageHandler = group.UpdateGroupImageHandlerFunc(UpdateGroupImage)

	api.UserCreateUserHandler = user.CreateUserHandlerFunc(CreateUser)
	api.UserGetUserHandler = user.GetUserHandlerFunc(GetUser)
	api.UserGetUserImageHandler = user.GetUserImageHandlerFunc(GetUserImage)
	api.UserUpdateUserHandler = user.UpdateUserHandlerFunc(UpdateUser)
	api.UserUpdateUserImageHandler = user.UpdateUserImageHandlerFunc(UpdateUserImage)
	api.ShoppinglistCreateListItemHandler = shoppinglist.CreateListItemHandlerFunc(CreateListItem)
	api.ShoppinglistUpdateListItemHandler = shoppinglist.UpdateListItemHandlerFunc(UpdateListItem)
	api.ShoppinglistGetListItemsHandler = shoppinglist.GetListItemsHandlerFunc(GetListItems)
}
