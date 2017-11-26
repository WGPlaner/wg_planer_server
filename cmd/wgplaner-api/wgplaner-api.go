package main

import (
	"log"

	"github.com/wgplaner/wg_planer_server/controllers"
	"github.com/wgplaner/wg_planer_server/modules/setting"
	"github.com/wgplaner/wg_planer_server/restapi"
	"github.com/wgplaner/wg_planer_server/restapi/operations"
)

// Version holds the current WGPlaner version
var Version = "0.0.1"

func init() {
	setting.AppVersion = Version
}

func main() {
	var (
		api    = operations.NewWgplanerAPI(setting.LoadSwaggerSpec(restapi.SwaggerJSON))
		server = restapi.NewServer(api)
	)

	defer server.Shutdown()

	controllers.GlobalInit()
	controllers.InitializeControllers(api)

	server.Port = setting.AppConfig.Server.Port

	// serve API
	if err := server.Serve(); err != nil {
		log.Fatalln(err)
	}
}
