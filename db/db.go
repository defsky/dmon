package db

import (
	"database/sql"

	"github.com/gomodule/redigo/redis"
	"github.com/xormplus/xorm"
)

// Mysql ...
func Mysql(name string) *xorm.Engine {
	db, ok := mysqlDBs[name]
	if ok {
		return db
	}
	panic("mysql db not configured")
}

// Mssql ...
func Mssql(name string) *xorm.Engine {
	db, ok := mssqlDBs[name]
	if ok {
		return db
	}
	panic("mssql db not configured")
}

// Redis ...
func Redis() redis.Conn {
	if redisPool == nil {
		panic("redis db not configured")
	}
	return redisPool.Get()
}

// FetchRows ...
func FetchRows(rows *sql.Rows) ([]map[string]interface{}, error) {
	cols, _ := rows.Columns()

	records := make([]map[string]interface{}, 0)
	for rows.Next() {
		columns := make([]interface{}, len(cols))
		columnPtrs := make([]interface{}, len(cols))

		for i := range columns {
			columnPtrs[i] = &columns[i]
		}
		if err := rows.Scan(columnPtrs...); err != nil {
			return nil, err
		}

		m := make(map[string]interface{})
		for i, name := range cols {
			m[name] = columns[i]
		}
		// log.Println(m)
		records = append(records, m)
	}

	return records, nil
}

// Init ...
func Init() {
	initMssql()
	initMysql()
	initRedis()
}
