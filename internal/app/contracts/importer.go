package contracts

import "huangc28/go-ios-iap-vendor/internal/app/models"

type ProcurementDAOer interface {
	GetPendingProcurements() ([]*models.Procurement, error)
}
