package apperrors

const (
	FailedToCheckStockReservedForUser = "6000001"
	StockIsNotReservedForTheUser      = "6000002"
	FailedToMarkStockAsDeliver        = "6000003"
)

var exporterErrMap = map[string]string{
	StockIsNotReservedForTheUser: "stock is not reserved for that user",
}
