package main

import (
	"log"
	"math/rand"
	"time"

	"github.com/wgplaner/wg_planer_server/gen/restapi"
	"github.com/wgplaner/wg_planer_server/gen/restapi/operations"
	"github.com/wgplaner/wg_planer_server/wgplaner"
	"github.com/wgplaner/wg_planer_server/wgplaner/controllers"

	"github.com/go-openapi/loads"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
)

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
	controllers.InitializeControllers(api)

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
