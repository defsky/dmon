package db

import (
	"pmon/config"
	"testing"
)

func TestDB(t *testing.T) {

	resp, err := Redis().Do("GET", config.GetConfig().DB.Redis.Keys["status"])
	if err != nil {
		t.Error(err)
	}
	t.Logf("redis resp: %s\n", resp)

	qlen := -1
	errs := Mssql("store").Table("U9.ProcingVoucher").Count(&qlen).GetErrors()
	if len(errs) > 0 {
		t.Logf("DB errors: %v\n", errs)
	}
	t.Logf("qlen: %d\n", qlen)
}
