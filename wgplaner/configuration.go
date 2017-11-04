package wgplaner

import (
	"errors"
	"log"

	"github.com/BurntSushi/toml"
)

const (
	DRIVER_SQLITE = "sqlite"
	DRIVER_MYSQL  = "mysql"
)

type serverConfig struct {
	Port int
}

type databaseConfig struct {
	Driver            string
	SqliteFile        string `toml:"sqlite_file"`
	MysqlServer       string `toml:"mysql_server"`
	MysqlPort         int    `toml:"mysql_port"`
	MysqlUser         string `toml:"mysql_user"`
	MysqlPassword     string `toml:"mysql_password"`
	MysqlDatabaseName string `toml:"mysql_db_name"`
}

type appConfigType struct {
	Server   serverConfig
	Database databaseConfig
}

func validateConfiguration(config *appConfigType) error {
	if !isValidDriverName(config.Database.Driver) {
		return errors.New("error in configuration: Invalid driver name")
	}
	//if config.Database.Driver == DRIVER_SQLITE {
	//	if _, err := os.Stat(config.Database.SqliteFile); os.IsNotExist(err) {
	//		log.Fatal("File is missing: config/serviceAccountKey.json") // exit program
	//	}
	//}
	if config.Server.Port < 80 {
		return errors.New("error in configuration: Portnumber is not valid (must be > 80)")
	}
	return nil
}

func LoadAppConfiguration() *appConfigType {
	var appConfig = &appConfigType{}

	// Path is relative to executable.
	if _, err := toml.DecodeFile("config/config.toml", appConfig); err != nil {
		log.Fatal("[Configuration] Error loading configuration! ", err)
		return nil
	}

	if err := validateConfiguration(appConfig); err != nil {
		log.Fatal("[Configuration] Error validating configuration! ", err.Error())
		return nil
	}

	log.Println("[Configuration] Configuration successfully loaded!")

	return appConfig
}

var AppConfig *appConfigType
