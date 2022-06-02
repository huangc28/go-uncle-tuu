package inventory

import (
	"fmt"
)

type StockNotEnoughError struct {
	ProdID string
}

func (e *StockNotEnoughError) Error() string {
	return fmt.Sprintf("stock not enough for %s", e.ProdID)
}

// Organize stockInInv to following structure.
//
// {
//   arktw_diamond_2: [] // length is the number of available inventory.
//   arktw_diamond_1: []
// }
func IsQuantityInStockEnoughForAssigning(reqStocks []StockParam, stocks []*AvailableStockInfo) error {
	availableStockMap := make(map[string][]*AvailableStockInfo, 0)

	for _, stock := range stocks {
		if _, exists := availableStockMap[stock.ProdID]; !exists {
			availableStockMap[stock.ProdID] = []*AvailableStockInfo{stock}
		} else {
			availableStockMap[stock.ProdID] = append(availableStockMap[stock.ProdID], stock)
		}
	}

	// Check if stocks has enough for request.
	for _, reqStock := range reqStocks {
		// if req stock does not exists in `availableStockMap`, that means we do not have available stock
		if _, stockExists := availableStockMap[reqStock.ProdID]; !stockExists {
			return &StockNotEnoughError{
				ProdID: reqStock.ProdID,
			}
		}

		availableStocks := availableStockMap[reqStock.ProdID]
		if len(availableStocks) < reqStock.Quantity {
			return &StockNotEnoughError{
				ProdID: reqStock.ProdID,
			}
		}
	}

	return nil
}
