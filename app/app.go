package app

import (
	"encoding/json"
	"log"
	"time"

	"github.com/defsky/dmon/config"
	"github.com/defsky/dmon/db"
	"github.com/gomodule/redigo/redis"
	"github.com/xormplus/xorm"
)

var u9db *xorm.Engine
var rds redis.Conn

// Start ...
func Start() {
	for {
		dataset := make([]*DataItem, 0)

		if mo := getBadMO(); mo != nil {
			dataset = append(dataset, mo)
		}
		if repeated := getRepeatedDoc(); repeated != nil {
			dataset = append(dataset, repeated)
		}
		if badsite := getBadSiteDoc(); badsite != nil {
			dataset = append(dataset, badsite)
		}
		if item := getNotApprovedDoc(); item != nil {
			dataset = append(dataset, item)
		}
		if item := getBadBom(); item != nil {
			dataset = append(dataset, item)
		}

		datasetjson, err := json.Marshal(dataset)
		if err != nil {
			log.Println(err)
		} else {
			log.Println("baddoc:", string(datasetjson))
			rds.Do("SET", "dashboard:baddoc", string(datasetjson))
		}

		time.Sleep(30 * time.Second)
	}
}

func init() {
	config.Init()
	db.Init()

	u9db = db.Mssql("u928")
	rds = db.Redis()
}
