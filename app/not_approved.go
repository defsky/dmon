package app

import (
	"log"
	"strconv"
)

func getNotApprovedDoc() (*DataItem, *BadDocAgg) {
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

	lotAdd := u9db.SQL(`
	with CTE_Supplier as (
		select a.ID,a.Code,b.Name from CBO_Supplier a inner join CBO_Supplier_Trl b on b.ID = a.ID
	),
	CTE_Lot as (
		select a.id,a.lotcode,a.docno from Lot_LotMaster a
	)
	update b
	set b.invlot = b1.ID
	from PM_Receivement a 
		inner join PM_RcvLine b on b.receivement = a.id
		inner join CTE_Lot b1 on b1.LotCode = b.InvLotCode
	where a.ReceivementType = 1
		and a.DescFlexField_PubDescSeg26 <> ''
		and b.invlot is null and b.InvLotCode <> ''
		and b1.LotCode is not null`)

	records, err := sess.QueryString()
	if err != nil {
		log.Println(err)
		return nil, nil
	}

	drillkey := "dashboard:baddoc:notapproved"
	doc := &DataItem{Name: "未审核单", DrillKey: drillkey}
	data := make([][]string, 0)

	for _, row := range records {
		if c, err := strconv.Atoi(row["count"]); err == nil {
			if c > 0 {
				doc.Value += c
				data = append(data, []string{row["doctype"], row["count"]})

				if row["doctype"] == "采购退货" {
					if r, err := lotAdd.Execute(); err == nil {
						if affected, err := r.RowsAffected(); err == nil {
							log.Printf("采购退货批号已更新，受影响行数: %d", affected)
						} else {
							log.Println(err)
						}
					} else {
						log.Println(err)
					}
				}
			}

		} else {
			log.Println(err)
			return nil, nil
		}
	}

	detail := &BadDocAgg{
		ColNames: []*ColHeadSet{{Name: "单据类型", Width: 100}, {Name: "数量", Width: 100}},
		Data:     data,
	}

	return doc, detail
}
