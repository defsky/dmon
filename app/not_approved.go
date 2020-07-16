package app

import (
	"encoding/json"
	"log"
	"strconv"
)

func getNotApprovedDoc() *DataItem {
	sess := u9db.SQL(`
select '应收单' as 'doctype', COUNT(1) as 'count' 
from AR_ARBillHead a
where a.CreatedBy like 'U9%'
	and a.DocStatus <> 2
union
select '完工报告' as 'doctype', COUNT(1) as 'count' 
from MO_CompleteRpt a
	left join CBO_Wh b on b.ID = a.RcvWh
where a.DocState <> 3 and b.Code = '1001'
union
select '杂收单' as 'doctype',COUNT(1) as 'count'
from InvDoc_MiscRcvTrans a
where (a.CreatedBy like 'U9%' or a.CreatedBy like '曹良成%')
	and a.[Status]<>2
union
select '杂发单' as 'doctype',COUNT(1) as 'count'
from InvDoc_MiscShip a
where (a.CreatedBy like 'U9%' or a.CreatedBy like '曹良成%')
	and a.[Status] <> 2
union
select '采购收货' as 'doctype',COUNT(1) as 'count'
from PM_Receivement a
where a.ReceivementType = 0 and a.CreatedBy like 'U9%' and a.[Status]<>5
union
select '采购退货' as 'doctype',COUNT(1) as 'count'
from PM_Receivement a
where a.ReceivementType = 1 and a.CreatedBy like 'U9%' and a.[Status]=0
union
select '退回处理' as 'doctype',COUNT(1) as 'count'
from SM_RMA a
where a.[Status]<>3
union
select '调入单' as 'doctype',COUNT(1) as 'count'
from InvDoc_TransferIn a
where a.[Status]<>2 and a.CreatedBy='111'
union
select '形态转换' as 'doctype',COUNT(1) as 'count'
from InvDoc_TransferForm a
where a.CreatedBy in ('门店系统抛单账号','曹良成1')
	and a.[Status]<>2
union
select '出货单' as 'doctype',COUNT(1) as 'count'
from SM_Ship a
where a.[Status] <> 3`)

	records, err := sess.QueryString()
	if err != nil {
		log.Println(err)
		return nil
	}

	drillkey := "dashboard:baddoc:notapproved"
	doc := &DataItem{Name: "未审核单", DrillKey: drillkey}
	data := make([][]string, 0)

	for _, row := range records {
		if c, err := strconv.Atoi(row["count"]); err == nil {
			if c > 0 {
				doc.Value += c
				data = append(data, []string{row["doctype"], row["count"]})
			}

		} else {
			log.Println(err)
			return nil
		}
	}

	detail := &BadDocAgg{
		ColNames: []*ColHeadSet{{Name: "单据类型", Width: 100}, {Name: "数量", Width: 100}},
		Data:     data,
	}

	detailJSON, err := json.Marshal(detail)
	if err != nil {
		log.Println(err)
		return nil
	}

	log.Println("not approved doc", string(detailJSON))
	rds.Do("SET", drillkey, string(detailJSON))

	return doc
}
