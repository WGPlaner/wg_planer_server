package wgplaner

import (
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/core"
	"github.com/go-xorm/xorm"
	_ "github.com/mattn/go-sqlite3"
	"github.com/wgplaner/wg_planer_server/gen/models"
)

var implementedDrivers = []string{
	DRIVER_MYSQL,
	DRIVER_SQLITE,
}
var ormEngine *xorm.Engine

//func sync(engine *xorm.Engine) error {
//	return engine.Sync(&SyncLoginInfo2{}, &SyncUser2{}, &models.User{}, &models.Group{})
//}

var OrmEngine *xorm.Engine

func isValidDriverName(driverName string) bool {
	return StringInSlice(driverName, implementedDrivers)
}

func CreateOrmEngine(dbConfig *databaseConfig) *xorm.Engine {
	var err error
	var engine *xorm.Engine

	switch AppConfig.Database.Driver {
	case DRIVER_MYSQL:
		engine, err = getMysqlEngine(dbConfig)
	case DRIVER_SQLITE:
		engine, err = xorm.NewEngine("sqlite3", dbConfig.SqliteFile)
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
	engine.Logger().SetLevel(core.LOG_DEBUG)

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
	)
}
