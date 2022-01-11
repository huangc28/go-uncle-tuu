package models

type PurchaseRecord struct {
	Inventory
	ProdName string `json:"prod_name"`
}
