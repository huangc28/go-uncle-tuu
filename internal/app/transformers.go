package app

import (
	"fmt"
	"huangc28/ios-inapp-trade/internal/models"
)

type TrfedInventoryInfo struct {
	Inventory []*models.InventoryInfo `json:"inventory"`
}

func TrfInventory(ms []*models.InventoryInfo) *TrfedInventoryInfo {
	return &TrfedInventoryInfo{
		Inventory: ms,
	}
}

func TrfAvailableStock(m models.InventoryProduct) interface{} {
	tt := m.TransactionTime

	ct := fmt.Sprintf(
		"%d-%02d-%02dT%02d:%02d:%02d",
		tt.Year(),
		tt.Month(),
		tt.Day(),
		tt.Hour(),
		tt.Minute(),
		tt.Second(),
	)

	return struct {
		UUID            string `json:"uuid"`
		ProdID          string `json:"prod_id"`
		TransactionID   string `json:"transaction_id"`
		TransactionTime string `json:"transaction_time"`
		Receipt         string `json:"receipt"`
	}{
		m.UUID,
		m.ProdID,
		m.TransactionID,
		ct,
		m.Receipt,
	}
}
