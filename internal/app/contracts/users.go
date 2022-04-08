package contracts

import "huangc28/go-ios-iap-vendor/internal/app/models"

type UserDAOer interface {
	GetUserByUUID(uuid string, fields ...string) (*models.User, error)
	DisableExport(userID int64) error
}
