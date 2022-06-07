package models

import "time"

type PurchaseRecord struct {
	Inventory
	ProdName string `json:"prod_name"`
}

type ReservedStockInfo struct {
	ReservedNum int     `json:"reserved_num"`
	ProdName    string  `json:"prod_name"`
	Price       float64 `json:"price"`
}

type ProductListOption struct {
	ProdName     string `json:"prod_name"`
	ProdBundleID string `json:"prod_bundle_id"`
	NumInStock   int    `json:"num_in_stock"`
}

type StockAssigmentStatus struct {
	AssigneeName   string `json:"assignee_name"`
	AssignmentUUID string `json:"assignment_uuid"`
	StockUUID      string `json:"stock_uuid"`

	CreatedAt time.Time       `json:"created_at"`
	ProdName  string          `json:"prod_name"`
	GameName  string          `json:"game_name"`
	Delivered DeliveredStatus `json:"delivered"`
	Available bool            `json:"available"`
}
