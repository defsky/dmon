package main

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