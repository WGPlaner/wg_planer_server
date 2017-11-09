package main

import (
	"log"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"

	"github.com/go-openapi/loads"
	"github.com/wgplaner/wg_planer_server/gen/restapi"
	"github.com/wgplaner/wg_planer_server/gen/restapi/operations"
	"github.com/wgplaner/wg_planer_server/gen/restapi/operations/group"
	"github.com/wgplaner/wg_planer_server/gen/restapi/operations/info"
	"github.com/wgplaner/wg_planer_server/gen/restapi/operations/user"
	"github.com/wgplaner/wg_planer_server/wgplaner"
	"github.com/wgplaner/wg_planer_server/wgplaner/controllers"
)

func initializeControllers(api *operations.WgplanerAPI) {
	api.UserIDAuthAuth = controllers.UserIDAuth
	api.FirebaseIDAuthAuth = controllers.FirebaseIDAuth
	// Create API handlers
	api.InfoGetLatestVersionHandler = info.GetLatestVersionHandlerFunc(controllers.GetVersionInfo)
	api.GroupCreateGroupHandler = group.CreateGroupHandlerFunc(controllers.CreateGroup)
	api.GroupGetGroupHandler = group.GetGroupHandlerFunc(controllers.GetGroup)
	api.UserCreateUserHandler = user.CreateUserHandlerFunc(controllers.CreateUser)
	api.UserGetUserHandler = user.GetUserHandlerFunc(controllers.GetUser)
	api.UserUpdateUserHandler = user.UpdateUserHandlerFunc(controllers.UpdateUser)
}

func main() {
	// load embedded swagger file
	swaggerSpec, err := loads.Analyzed(restapi.SwaggerJSON, "")
	if err != nil {
		log.Fatalln(err)
	}

	// create new service API ---------------------------------------------------
	api := operations.NewWgplanerAPI(swaggerSpec)
	server := restapi.NewServer(api)
	defer server.Shutdown()

	// load configuration and initialize ----------------------------------------
	wgplaner.AppConfig = wgplaner.LoadAppConfigurationOrFail()
	wgplaner.OrmEngine = wgplaner.CreateOrmEngine(&wgplaner.AppConfig.Database)
	wgplaner.FireBaseApp = wgplaner.CreateFirebaseConnection()
	initializeControllers(api)

	if wgplaner.AppConfig.Mail.SendTestMail {
		wgplaner.SendTestMail()
	}

	// set the port this service will be run on ---------------------------------
	server.Port = wgplaner.AppConfig.Server.Port

	// serve API
	if err := server.Serve(); err != nil {
		log.Fatalln(err)
	}
}
