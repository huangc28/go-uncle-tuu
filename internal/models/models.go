package models

type ProductInfo struct {
	ProdName     string  `json:"prod_name"`
	ProdDesc     string  `json:"prod_desc"`
	Price        float64 `json:"price"`
	GameBundleID string  `json:"game_bundle_id"`
}

type InventoryInfo struct {
	ProdID   string  `json:"prod_id"`
	ProdName string  `json:"prod_name"`
	Price    float64 `json:"price"`
	Quantity int     `json:"quantity"`
}
