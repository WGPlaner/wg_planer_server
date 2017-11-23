package main

import (
	"log"

	"github.com/wgplaner/wg_planer_server/restapi"
	"github.com/wgplaner/wg_planer_server/restapi/operations"
	"github.com/wgplaner/wg_planer_server/wgplaner"
	"github.com/wgplaner/wg_planer_server/wgplaner/controllers"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
)

// Version holds the current WGPlaner version
var Version = "0.0.1"

func init() {
	wgplaner.AppVersion = Version
}

func main() {
	var (
		api    = operations.NewWgplanerAPI(wgplaner.LoadSwaggerSpec())
		server = restapi.NewServer(api)
	)

	defer server.Shutdown()

	wgplaner.NewConfigContext()
	controllers.InitializeControllers(api)

	server.Port = wgplaner.AppConfig.Server.Port

	// serve API
	if err := server.Serve(); err != nil {
		log.Fatalln(err)
	}
}
