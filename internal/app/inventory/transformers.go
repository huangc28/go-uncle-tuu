package inventory

import (
	"fmt"
	"huangc28/go-ios-iap-vendor/internal/app/models"
)

func TrfAvailableStock(m models.Inventory) interface{} {
	tt := m.TransactionTime

	ct := fmt.Sprintf(
		"%d-%02d-%02dT%02d:%02d:%02d",
		tt.Year(),
		tt.Month(),
		tt.Day(),
		tt.Hour(),
		tt.Minute(),
		tt.Second(),
	)

	return struct {
		UUID            string `json:"uuid"`
		ProdID          int    `json:"prod_id"`
		TransactionID   string `json:"transaction_id"`
		TransactionTime string `json:"transaction_time"`
		Receipt         string `json:"receipt"`
	}{
		m.Uuid,
		int(m.ProdID.Int32),
		m.TransactionID.String,
		ct,
		m.Receipt.String,
	}
}

type TrfedReservedStocks struct {
	Count int `json:"count"`
}

// @TODO: We are not responding stock information here. Respond only the number
// of reserved stocks.
func TrfReservedStocks(ms []*models.Inventory) TrfedReservedStocks {
	return TrfedReservedStocks{
		Count: len(ms),
	}
}
