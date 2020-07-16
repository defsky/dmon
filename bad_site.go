package main

import (
	"encoding/json"
	"log"
	"time"
)

func getBadSiteDoc() *DataItem {
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
	) a LEFT JOIN CTE_Customer b on b.ID = a.Customer`)

	records, err := sess.QueryString()
	if err != nil {
		log.Println(err)
		return nil
	}

	drillkey := "dashboard:baddoc:badsite"
	badSite := &DataItem{Name: "位置类型", DrillKey: drillkey}
	badSite.Value = len(records)

	data := make([][]string, 0)
	for _, row := range records {
		datestr := row["DisableDate"]
		d, err := time.Parse("2006-01-02T15:04:05Z", datestr)
		if err == nil {
			datestr = d.Format("2006-01-02")
		}

		data = append(data, []string{row["Code"], row["Name"], row["IsEffective"], datestr, row["Erro"]})
	}

	detail := &BadDocAgg{
		ColNames: []*ColHeadSet{
			{},
		},
		Data: data,
	}
	detailJSON, err := json.Marshal(detail)
	if err != nil {
		log.Println(err)
		return nil
	}

	log.Println("bad site doc", string(detailJSON))
	rds.Do("SET", drillkey, string(detailJSON))

	return badSite
}
