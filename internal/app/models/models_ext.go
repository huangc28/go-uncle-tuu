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

type ProductListOption struct {
	ProdName     string `json:"prod_name"`
	ProdBundleID string `json:"prod_bundle_id"`
	NumInStock   int    `json:"num_in_stock"`
}
