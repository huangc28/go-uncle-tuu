package contracts

type InventoryDAOer interface {
	IsStockReservedForUser(stockUUID string, userID int64) (bool, error)
	MarkStockAsDelivered(stockUUID string) error
	MarkStockAsNotDelivered(stockUUID string) error
}
