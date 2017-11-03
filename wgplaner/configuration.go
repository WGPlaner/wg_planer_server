package wgplaner

import (
	"log"

	"github.com/BurntSushi/toml"
)

type databaseConfig struct {
	Server       string
	Port         int
	User         string
	Password     string
	DatabaseName string
}

type appConfigType struct {
	Database databaseConfig
}

var appConfig = &appConfigType{}

func LoadAppConfiguration() {
	// Path is relative to executable.
	if _, err := toml.DecodeFile("config/config.toml", appConfig); err != nil {
		log.Fatal("[Configuration] Error loading configuration! ", err)
		return
	}
	log.Println("[Configuration] Configuration successfully loaded!")
}

func GetAppConfig() *appConfigType {
	if appConfig == nil {
		LoadAppConfiguration()
	}
	return appConfig
}
