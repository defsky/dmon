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

var jobs []*Job

func registerJob(j ...*Job) {
	if jobs == nil {
		jobs = make([]*Job, 0, 5)
	}
	jobs = append(jobs, j...)
}

// Start will startup app
//  intervel is timer sleep time
func Start(interval int) {
	for {
		dataset := make([]*DataItem, 0)

		log.Println("开始遍历执行任务清单 ...")

		s := time.Now()
		for _, job := range jobs {
			log.Printf("执行任务 [%s] ...", job.name)

			st := time.Now()
			data := job.handler()
			elapsed := time.Since(st).Seconds()

			log.Printf("任务 [%s] 执行完毕，耗时 %fs", job.name, elapsed)

			if data != nil {
				dataset = append(dataset, data)
			}
		}
		el := time.Since(s).Seconds()

		log.Printf("任务清单执行完毕，总耗时 %fs", el)

		datasetjson, err := json.Marshal(dataset)
		if err != nil {
			log.Println(err)
		} else {
			log.Println("baddoc:", string(datasetjson))
			rds.Do("SET", "dashboard:baddoc", string(datasetjson))
		}

		log.Printf("等待下次触发，休眠 %ds ...", interval)
		time.Sleep(time.Duration(interval) * time.Second)
	}
}

// Init ...
func Init() {
	config.Init()
	db.Init()

	u9db = db.Mssql("u928")
	rds = db.Redis()

	registerJob(
		&Job{name: "检查问题生产订单", handler: getBadMO},
		&Job{name: "检查重复单据", handler: getRepeatedDoc},
		&Job{name: "检查客户位置问题", handler: getBadSiteDoc},
		&Job{name: "查询未审核的单据", handler: getNotApprovedDoc},
		&Job{name: "检查问题BOM", handler: getBadBom},
		&Job{name: "检查销售退货成本价", handler: getBadRMA})
}
