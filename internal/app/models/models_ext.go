package models

type PurchaseRecord struct {
	Inventory
	ProdName string `json:"prod_name"`
}

type ReservedStockInfo struct {
	ReservedNum int     `json:"reserved_num"`
	ProdName    string  `json:"prod_name"`
	Price       float64 `json:"price"`
}
