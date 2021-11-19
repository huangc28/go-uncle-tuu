package models

type ProductInfo struct {
	ProdName     string  `json:"prod_name"`
	ProdDesc     string  `json:"prod_desc"`
	Price        float64 `json:"price"`
	GameBundleID string  `json:"game_bundle_id"`
}
