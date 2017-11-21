package wgplaner

import (
	"log"
	"math/rand"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/wgplaner/wg_planer_server/gen/restapi"

	"github.com/BurntSushi/toml"
	"github.com/go-openapi/loads"
	"github.com/op/go-logging"
)

var configLog = logging.MustGetLogger("Config")

type serverConfig struct {
	Port int `toml:"port"`
}

type authConfig struct {
	IgnoreFirebase    bool   `toml:"ignore_firebase"`
	FirebaseProjectId string `toml:"firebase_project_id"`
	FirebaseServerKey string `toml:"firebase_server_key"`
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
	Auth     authConfig
	Data     dataConfig
	Database databaseConfig
	Mail     mailConfig
}

var (
	// Global settings
	AppConfig   *appConfigType
	AppPath     string // Path to executable
	AppVersion  string
	AppWorkPath string // Working directory

	IsWindows bool
)

func init() {
	IsWindows = runtime.GOOS == "windows"

	var err error
	if AppPath, err = getAppPath(); err != nil {
		log.Fatal(4, "Failed to get app path: %v", err)
	}
	AppWorkPath = getWorkPath(AppPath)
}

func NewConfigContext() {
	AppConfig = &appConfigType{}

	configPath := path.Join(AppWorkPath, "config/config.toml")

	// Path is relative to executable.
	if _, err := toml.DecodeFile(configPath, AppConfig); err != nil {
		configLog.Fatal("Error loading configuration! ", err)
		return
	}

	if err := validateConfiguration(AppConfig); err.HasErrors() {
		configLog.Fatal("Error validating configuration: \n" + err.String())
		return
	}

	configLog.Info("Configuration successfully loaded!")

	OrmEngine = CreateOrmEngine(&AppConfig.Database)
	FireBaseApp = CreateFirebaseConnection()

	if AppConfig.Mail.SendTestMail {
		SendTestMail()
	}

	// Seed the random number generator (needed for group codes)
	rand.Seed(time.Now().UTC().UnixNano())
}

func LoadSwaggerSpec() *loads.Document {
	if swaggerSpec, errSpec := loads.Analyzed(restapi.SwaggerJSON, ""); errSpec != nil {
		configLog.Fatal(errSpec)
		return nil

	} else {
		return swaggerSpec
	}
}

func getAppPath() (string, error) {
	var appPath string
	var err error

	if IsWindows && filepath.IsAbs(os.Args[0]) {
		appPath = filepath.Clean(os.Args[0])

	} else if appPath, err = exec.LookPath(os.Args[0]); err != nil {
		return "", err
	}

	appPath, err = filepath.Abs(appPath)
	if err != nil {
		return "", err
	}

	// Note: we don't use path.Dir here because it does not handle case
	// which path starts with two "/" in Windows: "//psf/Home/..."
	return strings.Replace(appPath, "\\", "/", -1), err
}

// Get the working directory path
func getWorkPath(appPath string) string {
	var workPath string

	if i := strings.LastIndex(appPath, "/"); i == -1 {
		workPath = appPath
	} else {
		workPath = appPath[:i]
	}

	return strings.Replace(workPath, "\\", "/", -1)
}

func validateConfiguration(config *appConfigType) ErrorList {
	err := ErrorList{}

	if configErr := validateServerConfig(config.Server); configErr.HasErrors() {
		err.Add("[Config] Invalid server config:")
		err.AddList(&configErr)
	}

	if configErr := validateAuthConfig(config.Auth); configErr.HasErrors() {
		err.Add("[Config] Invalid auth config:")
		err.AddList(&configErr)
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

func validateServerConfig(config serverConfig) ErrorList {
	errList := ErrorList{}

	if config.Port < 80 {
		errList.Add("[Config] Port number is not valid (must be > 80)")
	}

	return errList
}

func validateAuthConfig(config authConfig) ErrorList {
	errList := ErrorList{}

	// Nothing to do at the moment

	return errList
}
