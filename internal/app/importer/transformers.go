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
