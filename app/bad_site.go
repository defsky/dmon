package app

import (
	"log"
	"time"
)

func getBadSiteDoc() (*DataItem, *BadDocAgg) {
	sess := u9db.SQL(`
	WITH CTE_Customer AS (
		SELECT a.ID,a.Code,a.Effective_IsEffective, a.Effective_DisableDate,a1.Name 
		FROM dbo.CBO_Customer a
			inner join dbo.CBO_Customer_Trl a1 on a1.ID = a.ID
		WHERE a.Org = 1001703126479896
	)
	SELECT 
		b.Code,b.Name,
		b.Effective_IsEffective as 'IsEffective',
		b.Effective_DisableDate as 'DisableDate',
		a.Error 
	FROM (
		SELECT 
			a.Customer,
			'位置类型重复' as 'Error'
		from dbo.CBO_CustomerSite a
			inner join dbo.CBO_Customer b on b.ID = a.Customer
		where b.Org = 1001703126479896
		group by a.DescFlexField_PrivateDescSeg1,a.Customer
		having COUNT(a.ID)>1
		UNION
		SELECT 
			a.Customer,
			'未知的位置类型' as 'Error'
		from dbo.CBO_CustomerSite a
			inner join dbo.CBO_Customer b on b.ID = a.Customer
		where b.Org = 1001703126479896
			and a.DescFlexField_PrivateDescSeg1 not in ('0','1')
		UNION
		select 
			a.ID,
			'没有默认位置' as 'Error'
		from (
			SELECT 
				c1.ID,c1.code,
				convert(int,(c2.IsDefaultBillTo & c2.IsDefaultClaim & c2.IsDefaultContrast & c2.IsDefaultPayment & c2.IsDefaultShipTo)) as 'siteWeight'
			from dbo.CBO_Customer c1
				inner join dbo.CBO_CustomerSite c2 on c2.Customer = c1.ID
			where c1.Org = 1001703126479896
		) a
		group by a.ID
		having SUM(a.siteWeight)<>1
	) a LEFT JOIN CTE_Customer b on b.ID = a.Customer`)

	records, err := sess.QueryString()
	if err != nil {
		log.Println(err)
		return nil, nil
	}

	drillkey := "dashboard:baddoc:badsite"
	badSite := &DataItem{Name: "客户档案", DrillKey: drillkey}
	badSite.Value = len(records)

	data := make([][]string, 0)
	for _, row := range records {
		datestr := row["DisableDate"]
		d, err := time.Parse("2006-01-02T15:04:05Z", datestr)
		if err == nil {
			datestr = d.Format("2006-01-02")
		}

		data = append(data, []string{row["Code"], row["Name"], row["IsEffective"], datestr, row["Error"]})
	}

	detail := &BadDocAgg{
		ColNames: []*ColHeadSet{
			{Name: "客户编码", Width: 100},
			{Name: "客户名称", Width: 250},
			{Name: "有效性", Width: 80},
			{Name: "失效日期", Width: 100},
			{Name: "问题描述", Width: 200},
		},
		Data: data,
	}

	return badSite, detail
}
