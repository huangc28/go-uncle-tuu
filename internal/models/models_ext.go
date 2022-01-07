package models

import (
	"database/sql"
	"time"
)

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

type InventoryProduct struct {
	ID              int          `json:"id"`
	UUID            string       `json:"uuid"`
	ProdID          string       `json:"prod_id"`
	TransactionID   string       `json:"transaction_id"`
	Receipt         string       `json:"receipt"`
	TransactionTime time.Time    `json:"transaction_time"`
	Available       bool         `json:"available"`
	CreatedAt       time.Time    `json:"created_at"`
	UpdatedAt       time.Time    `json:"updated_at"`
	DeletedAt       sql.NullTime `json:"deleted_at"`
	Delivered       bool         `json:"delivered"`
}
