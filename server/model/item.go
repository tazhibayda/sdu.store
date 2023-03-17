package model

type Item struct {
	ID        int64  `json:"id"`
	ProductID int64  `json:"product_id"`
	Color     string `json:"color"`
	Size      string `json:"size"`
	Quantity  int64  `json:"quantity"`
}
