package app

import (
	"encoding/json"
	"log"
	"sync"
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

		var wg sync.WaitGroup
		wg.Add(len(jobs))

		log.Println("开始遍历执行任务清单 ...")
		s := time.Now()

		for _, job := range jobs {

			go func(j *Job) {
				log.Printf("任务启动 [%s] ...", j.name)

				st := time.Now()
				item, detail := j.handler()
				elapsed := time.Since(st).Seconds()

				if item != nil {
					dataset = append(dataset, item)
					if item.Value > 0 {
						// uploadToRedis(item.DrillKey, detail)
						d, err := json.Marshal(detail)
						if err != nil {
							log.Printf("任务详情 [%s] : %s", j.name, err)
							return
						}
						log.Printf("任务详情 [%s] : %s", j.name, string(d))
						// var fmtJSON bytes.Buffer
						// if err = json.Indent(&fmtJSON, d, "", "  "); err != nil {
						// 	log.Printf("任务详情 [%s] : %s", j.name, err)
						// } else {
						// 	log.Printf("任务详情 [%s] : %s", j.name, fmtJSON.String())
						// }

						if _, err := rds.Do("SET", item.DrillKey, string(d)); err != nil {
							log.Printf("任务详情 [%s] : 上传失败: %s", j.name, err)
						} else {
							log.Printf("任务详情 [%s] : 上传成功", j.name)
						}
					} else {
						log.Printf("任务详情 [%s] : no data", j.name)
					}
				}
				log.Printf("任务结束 [%s]，耗时 %fs", j.name, elapsed)
				wg.Done()
			}(job)
		}
		wg.Wait()
		uploadToRedis("dashboard:baddoc", dataset)
		el := time.Since(s).Seconds()

		// datasetjson, err := json.Marshal(dataset)
		// if err != nil {
		// 	log.Println(err)
		// } else {
		// 	log.Println("baddoc:", string(datasetjson))
		// 	rds.Do("SET", "dashboard:baddoc", string(datasetjson))
		// }

		log.Printf("任务清单执行完毕，总耗时 %fs", el)
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
