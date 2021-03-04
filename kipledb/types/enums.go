package types

type PageLimit struct {
	PageNo   int `json:"page_no"`   // page_no begin from 1.
	PageSize int    `json:"page_size"` // set to -1, when you want to get the total records.
	Count    int    `json:"count"`     // set to -1 when you want to count the total records.
}
