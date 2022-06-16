package apperrors

const (
	FailedToGetUserByUUID = "3000001"
	FailedToDisableExport = "3000002"
	UserCanNotExportStock = "3000003"
)

var userErrMap = map[string]string{
	UserCanNotExportStock: "因上次出庫未到帳，出庫安全鎖已觸發。請聯絡客服人員",
}
