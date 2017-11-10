package wgplaner

import (
	"log"

	"github.com/BurntSushi/toml"
)

type serverConfig struct {
	Port int
}

type dataConfig struct {
	UserImageDir     string `toml:"user_image_dir"`
	UserImageDefault string `toml:"user_image_default"`
}

type databaseConfig struct {
	Driver            string
	LogSQL            bool   `toml:"log_sql"`
	SqliteFile        string `toml:"sqlite_file"`
	MysqlServer       string `toml:"mysql_server"`
	MysqlPort         int    `toml:"mysql_port"`
	MysqlUser         string `toml:"mysql_user"`
	MysqlPassword     string `toml:"mysql_password"`
	MysqlDatabaseName string `toml:"mysql_db_name"`
}

type mailConfig struct {
	SendTestMail bool   `toml:"send_testmail"`
	SMTPPort     int    `toml:"smtp_port"`
	SMTPHost     string `toml:"smtp_host"`
	SMTPIdentity string `toml:"smtp_identity"`
	SMTPUser     string `toml:"smtp_user"`
	SMTPPassword string `toml:"smtp_password"`
}

type appConfigType struct {
	Server   serverConfig
	Data     dataConfig
	Database databaseConfig
	Mail     mailConfig
}

func validateConfiguration(config *appConfigType) ErrorList {
	err := ErrorList{}

	if config.Server.Port < 80 {
		err.Add("[Config] Portnumber is not valid (must be > 80)")
	}

	if configErr := ValidateDataConfig(config.Data); configErr.HasErrors() {
		err.Add("[Config] Invalid data config:")
		err.AddList(&configErr)
	}

	if configErr := ValidateDriverConfig(config.Database); configErr.HasErrors() {
		err.Add("[Config] Invalid driver:")
		err.AddList(&configErr)
	}

	if configErr := ValidateMailConfig(config.Mail); configErr.HasErrors() {
		err.Add("[Config] Invalid Mail Config:")
		err.AddList(&configErr)
	}

	return err
}

func LoadAppConfigurationOrFail() *appConfigType {
	var appConfig = &appConfigType{}

	// Path is relative to executable.
	if _, err := toml.DecodeFile("config/config.toml", appConfig); err != nil {
		log.Fatal("[Configuration] Error loading configuration! ", err)
		return nil
	}

	if err := validateConfiguration(appConfig); err.HasErrors() {
		log.Fatal("[Configuration] Error validating configuration: \n" + err.String())
		return nil
	}

	log.Println("[Configuration] Configuration successfully loaded!")

	return appConfig
}

var AppConfig *appConfigType
