package db2

import (
	"fmt"
	"log"

	"github.com/defsky/dmon/config"

	"github.com/xormplus/xorm"

	// mssql driver
	_ "github.com/denisenkom/go-mssqldb"
)

var (
	mssqlDBs map[string]*xorm.Engine
)

func initMssql() {
	mssqlDBs = make(map[string]*xorm.Engine)

	dbcfgs := config.GetConfig().DB.Mssql

	if len(dbcfgs) <= 0 {
		return
	}

	log.Println("Init Mssql databases ...")

	for name, cfg := range dbcfgs {
		dsn := fmt.Sprintf("server=%s;port=%s;user id=%s;password=%s;database=%s;charset=%s",
			cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.Database, cfg.Charset)

		log.Printf("Connecting %s ...", dsn)

		db, err := xorm.NewMSSQL(xorm.MSSQL_DRIVER, dsn)
		if err != nil {
			log.Println(err)
		}

		mssqlDBs[name] = db
	}
}
