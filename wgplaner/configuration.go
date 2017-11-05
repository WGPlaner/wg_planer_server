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

type mailConfig struct {
	SmtpPort     int    `toml:"smtp_port"`
	SmtpHost     string `toml:"smtp_host"`
	SmtpIdentity string `toml:"smtp_identity"`
	SmtpUser     string `toml:"smtp_user"`
	SmtpPassword string `toml:"smtp_password"`
}

type appConfigType struct {
	Server   serverConfig
	Database databaseConfig
	Mail     mailConfig
}

func validateConfiguration(config *appConfigType) error {
	if !isValidDriverName(config.Database.Driver) {
		return errors.New("error in configuration: Invalid driver name")
	}
	if config.Server.Port < 80 {
		return errors.New("error in configuration: Portnumber is not valid (must be > 80)")
	}
	if !IntInSlice(config.Mail.SmtpPort, []int{25, 465, 587}) {
		log.Println("[WARNING][Configuration] SMTP Port is not a default port!")
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
