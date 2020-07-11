package db

import (
	"fmt"
	"log"
	"pmon/config"

	"github.com/jinzhu/gorm"

	// mssql driver
	_ "github.com/jinzhu/gorm/dialects/mssql"
)

var (
	mssqlDBs map[string]*gorm.DB
)

func initMssql() {
	mssqlDBs = make(map[string]*gorm.DB)

	dbcfgs := config.GetConfig().DB.Mssql

	if len(dbcfgs) <= 0 {
		return
	}

	log.Println("Init Mssql databases ...")

	for name, cfg := range dbcfgs {
		url := fmt.Sprintf("sqlserver://%s:%s@%s:%s?database=%s&charset=%s&encrypt=disable",
			cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.Database, cfg.Charset)

		log.Printf("Connecting %s ...", url)

		db, err := gorm.Open("mssql", url)
		if err != nil {
			panic(err)
		}

		mssqlDBs[name] = db
	}
}
