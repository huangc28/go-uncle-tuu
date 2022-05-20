package apperrors

const (
	FailedToGetPurchasedRecords       = "3000001"
	FailedToCreateImportFailedFileDir = "3000002"
	FailedToInitGoogleStorageClient   = "3000003"
	FailedToOpenUploadedFile          = "3000004"
	FailedToUploadFileToGCS           = "3000005"
)

var importerErrMap = map[string]string{}
