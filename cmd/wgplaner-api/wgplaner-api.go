package main

import (
	"errors"
	"io"
	"log"
	"math/rand"
	"time"

	"github.com/wgplaner/wg_planer_server/gen/restapi"
	"github.com/wgplaner/wg_planer_server/gen/restapi/operations"
	"github.com/wgplaner/wg_planer_server/gen/restapi/operations/group"
	"github.com/wgplaner/wg_planer_server/gen/restapi/operations/info"
	"github.com/wgplaner/wg_planer_server/gen/restapi/operations/shoppinglist"
	"github.com/wgplaner/wg_planer_server/gen/restapi/operations/user"
	"github.com/wgplaner/wg_planer_server/wgplaner"
	"github.com/wgplaner/wg_planer_server/wgplaner/controllers"

	"github.com/go-openapi/loads"
	"github.com/go-openapi/runtime"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
)

func initializeControllers(api *operations.WgplanerAPI) {
	// Producers
	api.HTMLProducer = runtime.ProducerFunc(func(writer io.Writer, data interface{}) error {
		if html, ok := data.(string); !ok {
			return errors.New("Error in HTML Producer")
		} else {
			_, err := writer.Write([]byte(html))
			return err
		}
	})

	// Authentication
	api.UserIDAuthAuth = controllers.UserIDAuth
	api.FirebaseIDAuthAuth = controllers.FirebaseIDAuth

	// Create API handlers
	api.InfoGetLatestVersionHandler = info.GetLatestVersionHandlerFunc(controllers.GetVersionInfo)
	api.GroupCreateGroupHandler = group.CreateGroupHandlerFunc(controllers.CreateGroup)
	api.GroupGetGroupHandler = group.GetGroupHandlerFunc(controllers.GetGroup)
	api.GroupCreateGroupCodeHandler = group.CreateGroupCodeHandlerFunc(controllers.CreateGroupCode)
	api.GroupJoinGroupHandler = group.JoinGroupHandlerFunc(controllers.JoinGroup)
	api.GroupJoinGroupHelpHandler = group.JoinGroupHelpHandlerFunc(controllers.JoinGroupHelp)
	api.GroupLeaveGroupHandler = group.LeaveGroupHandlerFunc(controllers.LeaveGroup)
	api.UserCreateUserHandler = user.CreateUserHandlerFunc(controllers.CreateUser)
	api.UserGetUserHandler = user.GetUserHandlerFunc(controllers.GetUser)
	api.UserGetUserImageHandler = user.GetUserImageHandlerFunc(controllers.GetUserImage)
	api.UserUpdateUserHandler = user.UpdateUserHandlerFunc(controllers.UpdateUser)
	api.UserUpdateUserImageHandler = user.UpdateUserImageHandlerFunc(controllers.UpdateUserImage)
	api.ShoppinglistCreateListItemHandler = shoppinglist.CreateListItemHandlerFunc(controllers.CreateListItem)
	api.ShoppinglistGetListItemsHandler = shoppinglist.GetListItemsHandlerFunc(controllers.GetListItems)
}

func main() {
	var errSpec error
	var swaggerSpec *loads.Document

	// load embedded swagger file -----------------------------------------------
	if swaggerSpec, errSpec = loads.Analyzed(restapi.SwaggerJSON, ""); errSpec != nil {
		log.Fatalln(errSpec)
	}

	// create new service API ---------------------------------------------------
	api := operations.NewWgplanerAPI(swaggerSpec)
	server := restapi.NewServer(api)
	defer server.Shutdown()

	// load configuration and initialize ----------------------------------------
	wgplaner.AppConfig = wgplaner.LoadAppConfigOrFail()
	wgplaner.OrmEngine = wgplaner.CreateOrmEngine(&wgplaner.AppConfig.Database)
	wgplaner.FireBaseApp = wgplaner.CreateFirebaseConnection()
	initializeControllers(api)

	if wgplaner.AppConfig.Mail.SendTestMail {
		wgplaner.SendTestMail()
	}

	// Seed the random number generator (needed for group codes)
	rand.Seed(time.Now().UTC().UnixNano())

	// set the port this service will be run on ---------------------------------
	server.Port = wgplaner.AppConfig.Server.Port

	// serve API
	if err := server.Serve(); err != nil {
		log.Fatalln(err)
	}
}
