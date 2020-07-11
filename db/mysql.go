package db

import (
	"fmt"
	"log"
	"pmon/config"

	"github.com/jinzhu/gorm"

	// mysql driver
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

var (
	mysqlDBs map[string]*gorm.DB
)

func initMysql() {
	mysqlDBs = make(map[string]*gorm.DB)

	dbcfgs := config.GetConfig().DB.Mysql

	if len(dbcfgs) <= 0 {
		return
	}

	log.Println("Init Mysql databases ...")

	for name, cfg := range dbcfgs {
		url := fmt.Sprintf("%s:%s@(%s:%s)/%s?charset=%s&parseTime=True&loc=Local",
			cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.Database, cfg.Charset)

		log.Printf("Connecting %s ...", url)
		db, err := gorm.Open("mysql", url)
		if err != nil {
			panic(err)
		}
		mysqlDBs[name] = db
	}
}
