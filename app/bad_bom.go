package app

import (
	"log"
)

func getBadBom() (*DataItem, *BadDocAgg) {
	sess := u9db.SQL(`
	with CTE_ItemMaster as (
		select a.ID,a.Code,a.Name,a.AssetCategory,a.ItemFormAttribute
		from CBO_ItemMaster a 
		where a.Org=1001703126479896
	)
	select 
		a1.Code as 'ParentCode',a1.Name as 'ParentName',
		b.Sequence as 'LineNo',
		b1.Code as 'CompentCode',b1.Name as 'ComponentName',
		0 as 'ErrorType','发料方式错误' as 'Error'
	from CBO_BOMMaster a
		inner join CBO_BOMComponent b on b.BOMMaster = a.ID
		left join CTE_ItemMaster b1 on b1.ID = b.ItemMaster
		left join CTE_ItemMaster a1 on a1.ID = a.ItemMaster
	where a.Org=1001703126479896
		and b1.Code like '1101%' --工件
		and b.IssueStyle <> 2 --完工倒冲
	union
	select
		a1.Code as 'ParentCode',a1.Name as 'ParentName',
		b.Sequence as 'LineNo',
		b1.Code as 'CompentCode',b1.Name as 'ComponentName',
		1 as 'ErrorType','倒冲自制材料' as 'Error'
	from CBO_BOMMaster a
		inner join CBO_BOMComponent b on b.BOMMaster = a.ID
		left join CTE_ItemMaster b1 on b1.ID = b.ItemMaster
		left join CTE_ItemMaster a1 on a1.ID = a.ItemMaster
		left join CBO_Category c on c.ID = b1.AssetCategory
	where a.Org=1001703126479896
		and b.IssueStyle=2 --完工倒冲
		and c.Code='0106' --自制材料
	union
	select
		a1.Code as 'ParentCode',a1.Name as 'ParentName',
		b.Sequence as 'LineNo',
		b1.Code as 'CompentCode',b1.Name as 'ComponentName',
		2 as 'ErrorType','制造件子项总仓发料' as 'Error'
	from CBO_BOMMaster a
		inner join CBO_BOMComponent b on b.BOMMaster = a.ID
		left join CTE_ItemMaster b1 on b1.ID = b.ItemMaster
		left join CTE_ItemMaster a1 on a1.ID = a.ItemMaster
		left join CBO_Wh b2 on b2.ID = b.SupplyWareHouse
	where a.Org=1001703126479896
		and b2.Code = '1001' --总仓
		and a1.ItemFormAttribute = 10 --制造件`)

	records, err := sess.QueryString()
	if err != nil {
		log.Println(err)
		return nil, nil
	}

	drillkey := "dashboard:baddoc:badbomagg"
	badBOM := &DataItem{Name: "BOM", DrillKey: drillkey}
	badBOM.Value = len(records)

	df := NewDataFrame(records)
	dt := df.GroupBy("ParentCode", "Error").Agg(map[string]AggFunc{
		"lineno": func(d RealData) string {
			ret := ""
			for _, ptr := range d {
				if len(ret) > 0 {
					ret += ","
				}
				row := (*ptr)
				ret += row["LineNo"]
			}
			return ret
		},
	})

	detail := &BadDocAgg{
		ColNames: []*ColHeadSet{
			{Name: "母件料号", Width: 150},
			{Name: "问题类型", Width: 200},
			{Name: "关联行号", Width: 200}},
		Data: dt.Data,
	}

	return badBOM, detail
}
