package app

import "log"

// DataItem is top aggregate data for one job
type DataItem struct {
	Name     string `json:"name"`
	Value    int    `json:"value"`
	DrillKey string `json:"drillkey"`
}

// ColHeadSet is head property set of data table columns
type ColHeadSet struct {
	Name  string `json:"name"`
	Width int    `json:"width"`
}

// BadDocAgg is top class aggregate data of bad docs
type BadDocAgg struct {
	ColNames []*ColHeadSet `json:"colNames"`
	Data     [][]string    `json:"data"`
}

// JobHandler is defination of job handler function
type JobHandler func() *DataItem

// Job is a task definations
type Job struct {
	name    string
	handler JobHandler
}

var concatSlice = func(s []string, sp string) string {
	r := ""
	for _, v := range s {
		if len(r) > 0 {
			r += sp
		}
		r += v
	}

	return r
}

// AggFunc ...
type AggFunc func(RealData) string

// RealData ...
type RealData []*map[string]string

// AggData ...
type AggData []string

// GroupTree ...
type GroupTree map[string]interface{}

// DataTable ...
type DataTable struct {
	Columns []string
	Data    [][]string
}

// Agg ...
func (kt GroupTree) Agg(fn map[string]AggFunc) *DataTable {
	keys := make([]string, 0)
	fnList := make([]AggFunc, 0)
	for k, f := range fn {
		keys = append(keys, k)
		fnList = append(fnList, f)
	}
	dt := kt.aggregate(fnList)

	dt.Columns = keys
	return dt
}

func (kt GroupTree) aggregate(f []AggFunc) *DataTable {
	data := make([][]string, 0)
	for k, v := range kt {

		switch v.(type) {
		case RealData:
			row := []string{k}
			for _, f := range f {
				row = append(row, f(v.(RealData)))
			}
			data = append(data, row)
		case GroupTree:
			dt := v.(GroupTree).aggregate(f)

			for _, r := range dt.Data {
				row := []string{k}
				row = append(row, r...)
				data = append(data, row)
			}
		}
	}

	return &DataTable{
		Data: data,
	}
}

// GroupBy ...
func (kt GroupTree) GroupBy(key string) {
	for k, v := range kt {
		switch v.(type) {
		case RealData:
			gt := make(GroupTree)
			data := v.(RealData)
			for i, row := range data {
				v, ok := (*row)[key]
				if !ok {
					log.Printf("group key not exists: %s", key)
					break
				}

				if gt[v] == nil {
					gt[v] = make(RealData, 0)
				}
				gt[v] = append(gt[v].(RealData), data[i])
			}
			kt[k] = gt
		case GroupTree:
			v.(GroupTree).GroupBy(key)
		}
	}
}

// DataFrame ...
type DataFrame struct {
	groupBy []string
	aggFunc map[string]AggFunc
	data    []map[string]string
	gtree   GroupTree
}

// NewDataFrame ...
func NewDataFrame(d []map[string]string) *DataFrame {
	gt := make(GroupTree)
	rd := make(RealData, 0)
	for i := range d {
		rd = append(rd, &d[i])
	}
	gt["rootkey"] = rd

	return &DataFrame{
		data:  d,
		gtree: gt,
	}
}

// GroupBy ...
func (df *DataFrame) GroupBy(keys ...string) *DataFrame {
	if len(keys) == 0 {
		return df
	}

	for _, key := range keys {
		df.gtree.GroupBy(key)
	}
	df.groupBy = append(df.groupBy, keys...)

	return df
}

// Agg ...
func (df *DataFrame) Agg(ag map[string]AggFunc) *DataTable {
	if len(df.groupBy) <= 0 {
		log.Println("DataFrame has not been grouped")
		return nil
	}

	tree := df.gtree["rootkey"].(GroupTree)

	dt := tree.Agg(ag)

	dt.Columns = append(df.groupBy, dt.Columns...)

	return dt
}

// Select ...
func (df *DataFrame) Select(keys ...string) *DataTable {
	dt := &DataTable{
		Columns: make([]string, 0),
		Data:    make([][]string, 0),
	}
	for _, key := range keys {
		dt.Columns = append(dt.Columns, key)
	}
	return dt
}

func sum(n ...int) int {
	count := len(n)

	if count > 1 {
		c := n[count-1]
		return c + sum(n[0:count-1]...)
	}
	return n[0]
}
