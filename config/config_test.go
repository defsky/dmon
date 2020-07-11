package config

import (
	"testing"
)

func TestConfig(t *testing.T) {
	t.Logf("appConfig.target: %s\n", appConfig.Target)
	t.Logf("appConfig.mssql.store: %v\n", appConfig.DB.Mssql["store"])
	t.Logf("appConfig.mssql.u928: %v\n", appConfig.DB.Mssql["u928"])
	t.Logf("appConfig.mysql.foo: %v\n", appConfig.DB.Mysql["foo"])
	t.Logf("appConfig.mysql.bar: %v\n", appConfig.DB.Mysql["bar"])
	t.Logf("appConfig.redis.keys.status: %s\n", appConfig.DB.Redis.Keys["status"])
}
