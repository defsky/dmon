package app

// DataItem ...
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

// JobFunc is defination of job function
type JobFunc func() *DataItem

func registerJob(j ...JobFunc) {
	if jobs == nil {
		jobs = make([]JobFunc, 0, 5)
	}
	jobs = append(jobs, j...)
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
