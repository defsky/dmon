package app

import (
	"fmt"
	"log"
	"strconv"
	"time"
)

func getBadMO() (*DataItem, *BadDocAgg) {
	now := time.Now()
	firstDayOfMonth := fmt.Sprintf("%02d-%02d-%02d", now.Year(), now.Month(), 1)
	sess := u9db.SQL(`
select * from (
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
		and b.IssuedQty=0
		and b.IsCalcCost=0
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
		and a.DocNo like 'DZ%'
	union
	select 
		a.DocNo,b.DocLineNO,
		3 as 'ErrorType',
		'工件发料方式错误' as 'ErrorDesc'
		--b.CreatedOn,b.CreatedBy,
		--b.ModifiedOn,b.ModifiedBy,
		--d.Code,d.Name, b.IsCalcCost,
		--case b.ConsignProcessItemSrc 
		--	when 0 then '受托方供料'
		--	when 1 then '委托方带料(制造业)'
		--	when 2 then '受托方领料'
		--	when 3 then '委托方带料(加工贸易)'
		--end as 'ConsignProcessItemSrc'
	from MO_MO a 
		inner join MO_MOPickList b on b.MO=a.ID
		left join Base_Organization c on c.ID = a.Org
		left join CBO_ItemMaster d on d.ID = b.ItemMaster
		left join CBO_Category e on e.ID = d.AssetCategory
	where c.Code = '2017' and a.IsStartMO=0
		and a.DocState<>3 -- 完工
		and b.IssueStyle not in (2,4) -- 完工倒冲，不发料
		and d.Code like '1101%' --工件
	union 
	select '订单提前关闭' as 'DocNo',-1 as 'DocLineNo',
		4 as 'ErrorType',
		a.DocNo as 'ErrorDesc'
	from MO_MO a 
	where a.IsStartMO=0 and IsMRPorMPS=1
		and a.DocState=3 and a.CreatedBy='admin' and a.ParentMO is not null
		and a.TotalCompleteQty <> a.ProductQty
		-- and a.CreatedOn >= '2020-04-28'
		and a.ClosedOn >= '` + firstDayOfMonth + `'
) as a order by a.DocNo,a.DocLineNO`)

	records, err := sess.QueryString()
	if err != nil {
		log.Println(err)
		return nil, nil
	}

	drillkey := "dashboard:baddoc:badmoagg"
	badMO := &DataItem{Name: "生产订单", DrillKey: drillkey}

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
			if n, err := strconv.Atoi(lineno); err == nil && n < 0 {
				aggLineno[docno] = append(aggLineno[docno], row["ErrorDesc"])
			}
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
		listStr := ""
		if v > 20 {
			listStr = concatSlice(aggLineno[k][:20], ",") + "......"
		} else {
			listStr = concatSlice(aggLineno[k], ",")
		}
		data = append(data, []string{k, fmt.Sprintf("合计 %d 行有问题, %s", v, listStr)})
	}

	detail := &BadDocAgg{
		ColNames: []*ColHeadSet{
			{Name: "单号", Width: 200},
			{Name: "问题描述", Width: 500}},
		Data: data,
	}

	return badMO, detail
}
