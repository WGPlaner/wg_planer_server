package models

import (
	"errors"
	"fmt"
	"log"
	"path"

	"github.com/wgplaner/wg_planer_server/modules/setting"

	// Load MySQL driver
	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/core"
	"github.com/go-xorm/xorm"
	// Load PostgreSQL driver
	_ "github.com/lib/pq"
	// Load SQLite3 driver
	_ "github.com/mattn/go-sqlite3"
)

var (
	x      *xorm.Engine
	tables []interface{}
)

func init() {
	tables = []interface{}{
		new(Bill),
		new(User),
		new(Group),
		new(GroupCode),
		new(ListItem),
	}

	gonicNames := []string{"ID", "UID"}
	for _, name := range gonicNames {
		core.LintGonicMapper[name] = true
	}
}

// SetEngine sets the xorm.Engine
func SetEngine() error {
	x = getEngine()
	x.ShowSQL(true)
	return nil
}

// NewEngine initializes a new xorm.Engine
func NewEngine() (err error) {
	if err = SetEngine(); err != nil {
		return err
	}
	if err = x.Ping(); err != nil {
		return err
	}
	//if err = x.StoreEngine("InnoDB").Sync2(tables...); err != nil {
	//	return fmt.Errorf("sync database struct error: %v", err)
	//}
	return nil
}

func getEngine() *xorm.Engine {
	var err error
	var engine *xorm.Engine

	switch setting.AppConfig.Database.Driver {
	case setting.DriverMySQL:
		engine, err = getMysqlEngine()

	case setting.DriverSQLite:
		filePath := path.Join(setting.AppWorkPath, setting.AppConfig.Database.SqliteFile)
		engine, err = xorm.NewEngine("sqlite3", filePath)

	default:
		err = errors.New("unknown SQL driver")
	}

	if err != nil {
		log.Fatal("[Configuration] SQL Connection failed! ", err)
		return nil
	}

	engine.SetMapper(core.GonicMapper{})
	engine.ShowSQL(true)

	if err = engine.Sync(tables...); err != nil {
		log.Fatal("[SQL] Synchronization failed! ", err)
		return nil
	}

	if setting.AppConfig.Database.LogSQL {
		engine.Logger().SetLevel(core.LOG_DEBUG)
	} else {
		engine.Logger().SetLevel(core.LOG_WARNING)
	}

	return engine
}

func getMysqlEngine() (*xorm.Engine, error) {
	dataSource := fmt.Sprintf("%s:%s@/%s?charset=utf-8", setting.AppConfig.Database.MysqlServer,
		setting.AppConfig.Database.MysqlPassword, setting.AppConfig.Database.MysqlDatabaseName)

	return xorm.NewEngine("mysql", dataSource)
}
