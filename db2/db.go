package db2

import (
	"github.com/xormplus/xorm"
)

// Mssql ...
func Mssql(name string) *xorm.Engine {
	db, ok := mssqlDBs[name]
	if !ok {
		panic("Mssql not configured")
	}
	return db
}

// Init ...
func Init() {
	initMssql()
}
