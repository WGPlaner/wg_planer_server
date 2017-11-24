package setting

import (
	"encoding/json"
	"log"
	"math/rand"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/wgplaner/wg_planer_server/modules/base"

	"github.com/BurntSushi/toml"
	"github.com/go-openapi/loads"
	"github.com/op/go-logging"
)

var settingLog = logging.MustGetLogger("Config")

const (
	DRIVER_SQLITE = "sqlite"
	DRIVER_MYSQL  = "mysql"
)

type serverConfig struct {
	Port int `toml:"port"`
}

type authConfig struct {
	IgnoreFirebase    bool   `toml:"ignore_firebase"`
	FirebaseProjectId string `toml:"firebase_project_id"`
	FirebaseServerKey string `toml:"firebase_server_key"`
}

type dataConfig struct {
	UserImageDir      string `toml:"user_image_dir"`
	UserImageDefault  string `toml:"user_image_default"`
	GroupImageDir     string `toml:"group_image_dir"`
	GroupImageDefault string `toml:"group_image_default"`
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

	// Seed the random number generator (needed for group codes)
	rand.Seed(time.Now().UTC().UnixNano())
}

func NewConfigContext() {
	AppConfig = &appConfigType{}

	configPath := path.Join(AppWorkPath, "config/config.toml")

	// Path is relative to executable.
	if _, err := toml.DecodeFile(configPath, AppConfig); err != nil {
		settingLog.Fatal("Error loading configuration! ", err)
		return
	}

	validateConfiguration()

	settingLog.Info("Configuration successfully loaded!")

	FireBaseApp = CreateFirebaseConnection()

	if AppConfig.Mail.SendTestMail {
		SendTestMail()
	}

}

func LoadSwaggerSpec(msg json.RawMessage) *loads.Document {
	if swaggerSpec, errSpec := loads.Analyzed(msg, ""); errSpec != nil {
		settingLog.Fatal(errSpec)
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

func validateConfiguration() {
	validateServerConfig()
	validateAuthConfig()
	validateDataConfig()
	validateDriverConfig()
	validateMailConfig()
}

func validateServerConfig() {
	var e []string

	if AppConfig.Server.Port < 80 {
		e = append(e, "[Config] Port number is not valid (must be > 80)")
	}

	if len(e) > 0 {
		settingLog.Fatal("[Config] Error with server config:\n" + strings.Join(e, "\n"))
	}
}

func validateDataConfig() {
	var e []string

	if AppConfig.Data.UserImageDir == "" {
		e = append(e, "[Config][Data] 'user_image_dir' must not be empty!")

	} else if stat, err := os.Stat(path.Join(AppWorkPath, AppConfig.Data.UserImageDir)); err != nil {
		if os.IsNotExist(err) {
			e = append(e, "[Config][Data] 'user_image_dir' does not exist!")
		} else if os.IsPermission(err) {
			e = append(e, "[Config][Data] Permission denied for 'user_image_dir'!")
		} else {
			e = append(e, "[Config][Data] Unknown error with 'user_image_dir'! "+err.Error())
		}

	} else if !stat.IsDir() {
		e = append(e, "[Config][Data] 'user_image_dir' is not a directory!")
	}

	if len(e) > 0 {
		settingLog.Fatal("[Config] Error with data config:\n" + strings.Join(e, "\n"))
	}
}

func validateAuthConfig() {
	var e []string

	if !AppConfig.Auth.IgnoreFirebase {
		if AppConfig.Auth.FirebaseProjectId == "" {
			e = append(e, "[Config] Firebase Project ID is required if firebase is not deactivated")
		}
		if AppConfig.Auth.FirebaseServerKey == "" {
			e = append(e, "[Config] Firebase Server Key is required if firebase is not deactivated")
		}
	}

	if len(e) > 0 {
		settingLog.Fatal("[Config] Error with auth config:\n" + strings.Join(e, "\n"))
	}
}

func validateDriverConfig() {
	var e []string

	switch AppConfig.Database.Driver {
	case DRIVER_MYSQL:
		if AppConfig.Database.MysqlServer == "" {
			e = append(e, "[Config][MySQL] Server is empty!")
		}
		if AppConfig.Database.MysqlPort == 0 {
			e = append(e, "[Config][MySQL] Port is empty!")
		}
		if AppConfig.Database.MysqlUser == "" {
			e = append(e, "[Config][MySQL] User is empty!")
		}
		if AppConfig.Database.MysqlDatabaseName == "" {
			e = append(e, "[Config][MySQL] Databasename is empty!")
		}

	case DRIVER_SQLITE:
		if AppConfig.Database.SqliteFile == "" {
			e = append(e, "[Config][SQLite] File is empty! Must specify a filename!")
		}

	default:
		e = append(e, "[Driver] Drivername is not valid!")
	}

	if len(e) > 0 {
		settingLog.Fatal("[Config] Error with database config:\n" + strings.Join(e, "\n"))
	}
}

func validateMailConfig() {
	var e []string

	if !base.IntInSlice(AppConfig.Mail.SMTPPort, []int{25, 465, 587}) {
		mailLog.Warning("SMTP Port is not a default port!")
	}

	if len(e) > 0 {
		settingLog.Fatal("[Config] Error with mail config:\n" + strings.Join(e, "\n"))
	}
}
