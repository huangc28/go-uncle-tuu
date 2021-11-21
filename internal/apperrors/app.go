package apperrors

const (
	DuplicatedProductInfo      = "1000001"
	FailedToCheckProductExists = "1000002"
	FailedToBindAPIBody        = "1000003"
	FailedToCreateProdInfo     = "1000004"
)

var appErrMap = map[string]string{
	DuplicatedProductInfo:  "product info has been collected",
	FailedToCreateProdInfo: "failed to create product info",
}
