package app

// DataItem ...
type DataItem struct {
	Name     string `json:"name"`
	Value    int    `json:"value"`
	DrillKey string `json:"drillkey"`
}

// ColHeadSet ...
type ColHeadSet struct {
	Name  string `json:"name"`
	Width int    `json:"width"`
}

// BadDocAgg ...
type BadDocAgg struct {
	ColNames []*ColHeadSet `json:"colNames"`
	Data     [][]string    `json:"data"`
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
