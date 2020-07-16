package main

import (
	"encoding/json"
	"log"
	"time"
)

func getRepeatedDoc() *DataItem {
	sess := u9db.SQL(`
declare @start_date nvarchar(100)
set @start_date = '2020-01-01'
select 
	max(CAST( convert(varchar(10), a.CreatedOn,112) AS DATE)) as 'date',
	'完工报告' as 'doctype',
	a.DescFlexField_PubDescSeg26 as 'srcdocno',
	count(a.docno) as 'count'
from MO_CompleteRpt a
where a.DescFlexField_PubDescSeg26 <> ''
	and a.CreatedOn >= @start_date
group by a.DescFlexField_PubDescSeg26
having count(a.docno) <> 1
union all
select
	max(CAST( convert(varchar(10), a.CreatedOn,112) AS DATE)) as 'date',
	'杂收单' as 'doctype',
	a.DescFlexField_PubDescSeg26 as 'srcdocno',
	count(a.docno) as 'count'
from InvDoc_MiscRcvTrans a
where a.DescFlexField_PubDescSeg26 <> ''
	and a.CreatedOn >= @start_date
group by a.DescFlexField_PubDescSeg26
having count(a.docno) <> 1
union all
select
	max(CAST( convert(varchar(10), a.CreatedOn,112) AS DATE)) as 'date',
	'调入单' as 'doctype',
	a.DescFlexField_PubDescSeg26 as 'srcdocno',
	count(a.docno) as 'count'
from InvDoc_TransferIn a
where a.DescFlexField_PubDescSeg26 <> ''
	and a.CreatedOn >= @start_date
group by a.DescFlexField_PubDescSeg26
having count(a.docno) <> 1
union all
select
	max(CAST( convert(varchar(10), a.CreatedOn,112) AS DATE)) as 'date',
	'形态转换单' as 'doctype',
	a.DescFlexField_PubDescSeg26 as 'srcdocno',
	count(a.docno) as 'count'
from InvDoc_TransferForm a
where a.DescFlexField_PubDescSeg26 <> ''
	and a.CreatedOn >= @start_date
group by a.DescFlexField_PubDescSeg26
having count(a.docno) <> 1
union all
select
	max(CAST( convert(varchar(10), a.CreatedOn,112) AS DATE)) as 'date',
	'采购退货' as 'doctype',
	a.DescFlexField_PubDescSeg26 as 'srcdocno',
	count(a.docno) as 'count'
from PM_Receivement a
where a.ReceivementType = 1 -- 0 采购收货 1 采购退货 2 销售退回收货
	and a.DescFlexField_PubDescSeg26 <> ''
	and a.CreatedOn >= @start_date
group by a.DescFlexField_PubDescSeg26
having count(a.docno) <> 1
union all
select
	max(CAST( convert(varchar(10), a.CreatedOn,112) AS DATE)) as 'date',
	'杂发单' as 'doctype',
	a.DescFlexField_PubDescSeg26 as 'srcdocno',
	count(a.docno) as 'count'
from InvDoc_MiscShip a
where a.DescFlexField_PubDescSeg26 <> ''
	and a.CreatedOn >= @start_date
group by a.DescFlexField_PubDescSeg26
having count(a.docno) <> 1`)

	records, err := sess.QueryString()
	if err != nil {
		log.Println(err)
		return nil
	}

	drillkey := "dashboard:baddoc:repeated"
	repeatedDoc := &DataItem{Name: "重复单据", DrillKey: drillkey}
	repeatedDoc.Value = len(records)

	data := make([][]string, 0)
	for _, row := range records {
		datestr := row["date"]
		d, err := time.Parse("2006-01-02T15:04:05Z", datestr)
		if err == nil {
			datestr = d.Format("2006-01-02")
		}
		data = append(data, []string{datestr, row["doctype"], row["srcdocno"], row["count"]})
	}

	detail := &BadDocAgg{
		ColNames: []*ColHeadSet{
			{Name: "业务日期", Width: 120},
			{Name: "单据类型", Width: 100},
			{Name: "来源单号", Width: 180},
			{Name: "数量", Width: 100},
		},
		Data: data,
	}
	detailJSON, err := json.Marshal(detail)
	if err != nil {
		log.Println(err)
		return nil
	}

	log.Println("repeated doc", string(detailJSON))
	rds.Do("SET", drillkey, string(detailJSON))

	return repeatedDoc
}
