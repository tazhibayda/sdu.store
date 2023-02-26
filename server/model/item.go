package model

type Item struct {
	ID         int64  `json:"id"`
	CategoryId int64  `json:"category_id"`
	Color      string `json:"color"`
	Size       string `json:"size"`
	Quantity   int64  `json:"quantity"`
}
