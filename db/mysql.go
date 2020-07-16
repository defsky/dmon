package db

import (
	"fmt"
	"log"

	"github.com/defsky/dmon/config"
	"github.com/xormplus/xorm"

	// mysql driver
	_ "github.com/go-sql-driver/mysql"
)

var (
	mysqlDBs map[string]*xorm.Engine
)

func initMysql() {
	mysqlDBs = make(map[string]*xorm.Engine)

	dbcfgs := config.GetConfig().DB.Mysql

	if len(dbcfgs) <= 0 {
		return
	}

	log.Println("Init Mysql databases ...")

	for name, cfg := range dbcfgs {
		dsn := fmt.Sprintf("%s:%s@(%s:%s)/%s?charset=%s&parseTime=True&loc=Local",
			cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.Database, cfg.Charset)

		log.Printf("Initiatine %s ...", name)
		db, err := xorm.NewEngine(xorm.MYSQL_DRIVER, dsn)
		if err != nil {
			log.Println(err)
		}
		mysqlDBs[name] = db
	}
}
