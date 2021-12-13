package app

import "huangc28/ios-inapp-trade/internal/models"

type TrfedInventoryInfo struct {
	Inventory []*models.InventoryInfo `json:"inventory"`
}

func TrfInventory(ms []*models.InventoryInfo) *TrfedInventoryInfo {
	return &TrfedInventoryInfo{
		Inventory: ms,
	}
}

func TrfAvailableStock(m models.InventoryProduct) interface{} {
	return struct {
		UUID          string `json:"uuid"`
		ProdID        string `json:"prod_id"`
		TransactionID string `json:"transaction_id"`
		Receipt       string `json:"receipt"`
	}{
		m.UUID,
		m.ProdID,
		m.TransactionID,
		m.Receipt,
	}
}
