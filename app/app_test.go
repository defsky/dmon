package app

import (
	"fmt"
	"testing"
)

func TestAgg(t *testing.T) {
	data := []map[string]string{
		{"docno": "Test001", "lineno": "1", "code": "1001", "lot": "aaa"},
		{"docno": "Test001", "lineno": "2", "code": "1001", "lot": "bbb"},
		{"docno": "Test001", "lineno": "3", "code": "1002", "lot": "ccc"},
		{"docno": "Test001", "lineno": "4", "code": "1002", "lot": "ddd"},

		{"docno": "Test002", "lineno": "1", "code": "1001", "lot": "eee"},
		{"docno": "Test002", "lineno": "2", "code": "1001", "lot": "fff"},
		{"docno": "Test002", "lineno": "3", "code": "1002", "lot": "ggg"},
		{"docno": "Test002", "lineno": "4", "code": "1002", "lot": "hhh"},
		{"docno": "Test002", "lineno": "5", "code": "1002", "lot": "jjj"},
	}

	df := NewDataFrame(data)

	dt := df.GroupBy("docno", "code").Agg(map[string]AggFunc{
		"info": func(d RealData) string {
			info := ""
			for _, rowptr := range d {
				if len(info) > 0 {
					info += ","
				}
				row := (*rowptr)
				info += fmt.Sprintf("%s:%s", row["lineno"], row["lot"])
			}
			return info
		},
		"linenos": func(d RealData) string {
			info := ""
			for _, rowptr := range d {
				if len(info) > 0 {
					info += ","
				}
				row := (*rowptr)
				info += row["lineno"]
			}
			return info
		},
	})

	t.Log(dt)
}
