package controllers

import (
	"github.com/go-openapi/runtime/middleware"
	"github.com/wgplaner/wg_planer_server/gen/restapi/operations/info"
	"github.com/wgplaner/wg_planer_server/wgplaner"
)

func GetVersionInfo(_ info.GetLatestVersionParams) middleware.Responder {
	return info.NewGetLatestVersionOK().WithPayload(&wgplaner.VersionInfo)
}
