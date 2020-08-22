package app

import (
	"encoding/json"
	"log"
)

func getBadRMA() *DataItem {
	sess := u9db.SQL(`
	WITH CTE_BOM as (
		select 
			a.ItemMaster as 'ParentItemMaster',
			b.ItemMaster as 'ComponentItemMaster'
		from CBO_BOMMaster a
		inner join CBO_BOMComponent b on b.BOMMaster = a.ID
		where a.Org = 1001703126479896
	)
	select * from (
		select
			a.DocNo,b.DocLineNo,c1.Code,
			c1.Name,'' as 'LotCode',
			case b.CostSource
				when 2 then c1.RefrenceCost -- 参考成本
				when 3 then b.CostPrice -- 手工录入
			end as 'CostPrice'
		from SM_RMA a
			inner join SM_RMALine b on b.RMA = a.ID
			inner join CTE_BOM c on c.ParentItemMaster = b.ItemInfo_ItemID
			left join CBO_ItemMaster b1 on b1.ID = b.ItemInfo_ItemID
			left join CBO_ItemMaster c1 on c1.ID = c.ComponentItemMaster
		where a.[Status]<>3 
			and b1.ItemFormAttribute=12 -- 套件
			--and c1.RefrenceCost <=0
		union
		select
			a.DocNo,b.DocLineNo,b1.Code ,b1.Name,b.LotInfo_LotCode as 'LotCode', 
			case b.CostSource
				when 2 then b1.RefrenceCost -- 参考成本
				when 3 then b.CostPrice -- 手工录入
			end as 'CostPrice'
		from SM_RMA a
			inner join SM_RMALine b on b.RMA = a.ID
			left join CBO_ItemMaster b1 on b1.ID = b.ItemInfo_ItemID
		where a.[Status]<>3 
			and b1.ItemFormAttribute<>12 
			--and b1.RefrenceCost <=0
	) a
	where a.CostPrice <= 0
	order by DocNo, DocLineNo`)

	records, err := sess.QueryString()
	if err != nil {
		log.Println(err)
		return nil
	}

	drillkey := "dashboard:baddoc:badrmaagg"
	badRMA := &DataItem{Name: "退回处理", DrillKey: drillkey}

	data := make([][]string, 0)
	for _, row := range records {
		badRMA.Value++
		data = append(data, []string{row["DocNo"], row["DocLineNo"], row["Code"], row["Name"], row["LotCode"], row["CostPrice"]})
	}

	detail := &BadDocAgg{
		ColNames: []*ColHeadSet{
			{Name: "单号", Width: 150},
			{Name: "行号", Width: 50},
			{Name: "料号", Width: 150},
			{Name: "品名", Width: 150},
			{Name: "批号", Width: 150},
			{Name: "成本价", Width: 100}},
		Data: data,
	}

	badrmajson, err := json.Marshal(detail)
	if err != nil {
		log.Println(err)
		return nil
	}

	log.Println("bad RMA agg", string(badrmajson))
	rds.Do("SET", drillkey, string(badrmajson))

	return badRMA
}
