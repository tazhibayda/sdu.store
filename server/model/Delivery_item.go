package model

type DeliveryItem struct {
	DeliveryID int64 `json:"delivery_id"`
	ItemID     int64 `json:"item_id"`
	Quantity   int64 `json:"quantity"`
}
