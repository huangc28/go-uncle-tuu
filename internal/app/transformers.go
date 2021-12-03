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
