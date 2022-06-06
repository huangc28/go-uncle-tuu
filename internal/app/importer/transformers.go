package importer

import (
	"huangc28/go-ios-iap-vendor/internal/app/models"
	"time"
)

type TrfPurchaseRecord struct {
	UUID            string    `json:"uuid"`
	ProdName        string    `json:"prod_name"`
	TransactionID   string    `json:"transaction_id"`
	TransactionTime time.Time `json:"transaction_time"`
}

func TtfPurchaseRecords(ms []models.PurchaseRecord) []TrfPurchaseRecord {
	trfms := make([]TrfPurchaseRecord, 0)

	for _, ms := range ms {
		trfm := TrfPurchaseRecord{
			UUID:            ms.Uuid,
			ProdName:        ms.ProdName,
			TransactionID:   ms.TransactionID.String,
			TransactionTime: ms.TransactionTime,
		}

		trfms = append(trfms, trfm)
	}

	return trfms
}

type TrfmedProcurement struct {
	Filename     string    `json:"filename"`
	UUID         string    `json:"uuid"`
	Status       string    `json:"status"`
	FailedReason *string   `json:"failed_reason"`
	CreatedAt    time.Time `json:"created_at"`
}

func TrfProcurement(p *models.Procurement) TrfmedProcurement {
	trfmp := TrfmedProcurement{
		Filename:  p.Filename,
		UUID:      p.Uuid,
		Status:    string(p.Status),
		CreatedAt: p.CreatedAt,
	}

	if p.FailedReason.Valid {
		trfmp.FailedReason = &p.FailedReason.String
	}

	return trfmp
}

func TrfProcurements(ps []*models.Procurement) []TrfmedProcurement {
	trfedps := make([]TrfmedProcurement, 0)

	for _, p := range ps {
		trfmp := TrfProcurement(p)

		trfedps = append(trfedps, trfmp)
	}

	return trfedps
}
