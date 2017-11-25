package controllers

import (
	"github.com/wgplaner/wg_planer_server/models"
	"github.com/wgplaner/wg_planer_server/restapi/operations/info"

	"github.com/go-openapi/runtime/middleware"
	"github.com/op/go-logging"
)

var infoLog = logging.MustGetLogger("Info")

var VersionInfo = models.VersionInfo{
	AndroidVersionCode:   1,
	AndroidVersionString: "v0.0.5",
	APIVersionCode:       1,
	APIVersionString:     "v0.0.5",
}

func GetVersionInfo(_ info.GetVersionParams) middleware.Responder {
	infoLog.Debug(`Get version info`)
	return info.NewGetVersionOK().WithPayload(&VersionInfo)
}
