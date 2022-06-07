package apperrors

const (
	FailedToGetReservedStock             = "4000001"
	NoReservedStockAvailable             = "4000002"
	FailedToGetAvailableStocksForProdIDs = "4000003"
	NotEnoughStocks                      = "4000004"
	FailedToAssignStocksToUser           = "4000005"
	FailedToCreateStockAssignment        = "4000006"
	FailedToGetAssignmentExportStatus    = "4000007	"
)

var inventoryErrMap = map[string]string{
	NoReservedStockAvailable: "沒有預留的庫存",
}
