package apperrors

const (
	FailedToGetReservedStock = "4000001"
	NoReservedStockAvailable = "4000002"
)

var inventoryErrMap = map[string]string{
	NoReservedStockAvailable: "沒有預留的庫存",
}
