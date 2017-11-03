package main

import (
	"flag"
	"github.com/go-openapi/loads"
	"github.com/wgplaner/wg_planer_server/controllers"
	"github.com/wgplaner/wg_planer_server/gen/restapi"
	"github.com/wgplaner/wg_planer_server/gen/restapi/operations"
	"github.com/wgplaner/wg_planer_server/gen/restapi/operations/user"
	"github.com/wgplaner/wg_planer_server/wgplaner"

	"log"
)

var portFlag = flag.Int("port", 3000, "Port to run this service on")

func initializeControllers(api *operations.WgplanerAPI) {
	// Create API handlers
	api.UserCreateUserHandler = user.CreateUserHandlerFunc(controllers.CreateUser)
	api.UserGetUserHandler = user.GetUserHandlerFunc(controllers.GetUser)
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

	// parse flags --------------------------------------------------------------
	flag.Parse()
	// set the port this service will be run on
	server.Port = *portFlag

	// load configuration and initialize ----------------------------------------
	wgplaner.LoadAppConfiguration()
	initializeControllers(api)
	controllers.InitialiseFirebaseConnection()

	// serve API
	if err := server.Serve(); err != nil {
		log.Fatalln(err)
	}
}
