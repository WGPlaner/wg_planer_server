package wgplaner

import (
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
)

var ormEngine *xorm.Engine

func CreateSqlConnection() {
	var err error
	dbConfig := GetAppConfig().Database
	dataSource := fmt.Sprintf("%s:%s@/%s?charset=utf-8", dbConfig.Server, dbConfig.Password, dbConfig.DatabaseName)

	if ormEngine, err = xorm.NewEngine("mysql", dataSource); err != nil {
		log.Fatal("[Configuration] SQL Connection failed! ", err)
	}
}

func GetOrmEngine() *xorm.Engine {
	if ormEngine == nil {
		CreateSqlConnection()
	}
	return ormEngine
}
