package repository

type WarehouseItem struct {
	Id       string `json:"id"`
	Min      int32  `json:"min"`
	Max      int32  `json:"max"`
	Quantity int32  `json:"quantity"`
}
