package wgplaner

import (
	"errors"
	"fmt"
	"log"

	"github.com/wgplaner/wg_planer_server/gen/models"

	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/core"
	"github.com/go-xorm/xorm"
	_ "github.com/mattn/go-sqlite3"
)

const (
	DRIVER_SQLITE = "sqlite"
	DRIVER_MYSQL  = "mysql"
)

var OrmEngine *xorm.Engine

func ValidateDriverConfig(config databaseConfig) ErrorList {
	err := ErrorList{}

	switch config.Driver {
	case DRIVER_MYSQL:
		if config.MysqlServer == "" {
			err.Add("[Config][MySQL] Server is empty!")
		}
		if config.MysqlPort == 0 {
			err.Add("[Config][MySQL] Port is empty!")
		}
		if config.MysqlUser == "" {
			err.Add("[Config][MySQL] User is empty!")
		}
		if config.MysqlDatabaseName == "" {
			err.Add("[Config][MySQL] Databasename is empty!")
		}

	case DRIVER_SQLITE:
		if config.SqliteFile == "" {
			err.Add("[Config][SQLite] File is empty! Must specify a filename!")
		}

	default:
		err.Add("[Driver] Drivername is not valid!")
	}

	return err
}

func CreateOrmEngine(dbConfig *databaseConfig) *xorm.Engine {
	var err error
	var engine *xorm.Engine

	switch AppConfig.Database.Driver {
	case DRIVER_MYSQL:
		engine, err = getMysqlEngine(dbConfig)
	case DRIVER_SQLITE:
		engine, err = xorm.NewEngine("sqlite3", dbConfig.SqliteFile)
	default:
		err = errors.New("unknown SQL driver")
	}

	if err != nil {
		log.Fatal("[Configuration] SQL Connection failed! ", err)
		return nil
	}

	if err = syncDatabaseTables(engine); err != nil {
		log.Fatal("[SQL] Synchronization failed! ", err)
		return nil
	}

	engine.ShowSQL(true)

	if AppConfig.Database.LogSQL {
		engine.Logger().SetLevel(core.LOG_DEBUG)
	} else {
		engine.Logger().SetLevel(core.LOG_WARNING)
	}

	return engine
}

func getMysqlEngine(config *databaseConfig) (*xorm.Engine, error) {
	dataSource := fmt.Sprintf("%s:%s@/%s?charset=utf-8", config.MysqlServer,
		config.MysqlPassword, config.MysqlDatabaseName)

	return xorm.NewEngine("mysql", dataSource)
}

func syncDatabaseTables(engine *xorm.Engine) error {
	return engine.Sync(
		&models.User{},
		&models.Group{},
		&models.GroupCode{},
		&models.ListItem{},
	)
}
