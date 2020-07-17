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
var jobs []JobFunc

// Start will startup app
func Start() {
	for {
		dataset := make([]*DataItem, 0)

		for _, job := range jobs {
			if data := job(); data != nil {
				dataset = append(dataset, data)
			}
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

	registerJob(getBadMO, getRepeatedDoc, getBadSiteDoc,
		getNotApprovedDoc, getBadBom)
}
