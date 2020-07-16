package main

import (
	"encoding/json"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/defsky/dmon/config"
	"github.com/defsky/dmon/db"
	"github.com/gomodule/redigo/redis"
	"github.com/xormplus/xorm"
)

// DataItem ...
type DataItem struct {
	Name     string `json:"name"`
	Value    int    `json:"value"`
	DrillKey string `json:"drillkey"`
}

var concatSlice = func(s []string, sp string) string {
	r := ""
	for _, v := range s {
		if len(r) > 0 {
			r += sp
		}
		r += v
	}

	return r
}

func init() {
	config.Init()
	db.Init()

	u9db = db.Mssql("u928")
	rds = db.Redis()
}

var u9db *xorm.Engine
var rds redis.Conn

func main() {
	go waitSignal()

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

func waitSignal() {
	c := make(chan os.Signal)
	signal.Notify(c)

	for {
		s := <-c
		switch s {
		case os.Interrupt:
			log.Println("User Interrupt")
			os.Exit(0)
		}
	}
}
