package controllers

import (
	"github.com/wgplaner/wg_planer_server/restapi/operations/info"
	"github.com/wgplaner/wg_planer_server/wgplaner"

	"github.com/go-openapi/runtime/middleware"
	"github.com/op/go-logging"
)

var infoLog = logging.MustGetLogger("Info")

func GetVersionInfo(_ info.GetLatestVersionParams) middleware.Responder {
	infoLog.Debug(`Get version info`)
	return info.NewGetLatestVersionOK().WithPayload(&wgplaner.VersionInfo)
}
