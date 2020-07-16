package app

import (
	"encoding/json"
	"fmt"
	"log"
)

func getBadMO() *DataItem {
	sess := u9db.SQL(`
	select 
	a.DocNo,b.DocLineNO,
	1 as 'ErrorType',
	'备料未计算成本' as 'ErrorDesc'
from MO_MO a 
	inner join MO_MOPickList b on b.MO=a.ID
	left join Base_Organization c on c.ID = a.Org
	left join CBO_ItemMaster d on d.ID = b.ItemMaster
where c.Code = '2017' and a.IsStartMO=0
	and a.DocState<>3 -- 完工
	and b.IssueStyle=0 -- 推式
	and b.IssuedQty < b.ActualReqQty
	and b.IsCalcCost=0
	and b.ModifiedOn > ?
union all
select 
	a.DocNo,b.DocLineNO,
	2 as 'ErrorType',
	'备料为自制材料' as 'ErrorDesc'
from MO_MO a 
	inner join MO_MOPickList b on b.MO=a.ID
	left join Base_Organization c on c.ID = a.Org
	left join CBO_ItemMaster d on d.ID = b.ItemMaster
	left join CBO_Category e on e.ID = d.AssetCategory
where c.Code = '2017' and a.IsStartMO=0
	and a.DocState<>3 -- 完工
	and b.IssueStyle=2 -- 完工倒冲
	and e.Code = '0106' -- 财务分类-自制材料
	and a.DocNo like 'DZ%'`, "2020-07-01")

	records, err := sess.QueryString()
	if err != nil {
		log.Println(err)
		return nil
	}

	drillkey := "dashboard:baddoc:badmoagg"
	badMO := &DataItem{Name: "生产备料", DrillKey: drillkey}

	aggCount := make(map[string]int)
	aggLineno := make(map[string][]string)

	oldErrType := ""
	for _, row := range records {
		badMO.Value++
		docno := row["DocNo"]
		lineno := row["DocLineNO"]
		errTypeCurrent := row["ErrorType"]
		if _, ok := aggCount[docno]; ok {
			aggCount[docno]++
			if errTypeCurrent != oldErrType {
				oldErrType = errTypeCurrent
				aggLineno[docno] = append(aggLineno[docno], row["ErrorDesc"])
			}
			aggLineno[docno] = append(aggLineno[docno], lineno)
			continue
		}
		aggCount[docno] = 1
		oldErrType = row["ErrorType"]
		aggLineno[docno] = make([]string, 0)
		aggLineno[docno] = append(aggLineno[docno], row["ErrorDesc"])
		aggLineno[docno] = append(aggLineno[docno], lineno)
	}

	data := make([][]string, 0)

	for k, v := range aggCount {
		data = append(data, []string{k, fmt.Sprintf("合计 %d 行有问题, %s", v, concatSlice(aggLineno[k], ","))})
	}

	detail := &BadDocAgg{
		ColNames: []*ColHeadSet{
			{Name: "单号", Width: 200},
			{Name: "问题描述", Width: 500}},
		Data: data,
	}

	badmoaggjson, err := json.Marshal(detail)
	if err != nil {
		log.Println(err)
		return nil
	}

	log.Println("bad mo agg", string(badmoaggjson))
	rds.Do("SET", drillkey, string(badmoaggjson))

	return badMO
}
