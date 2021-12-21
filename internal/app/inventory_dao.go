package app

import (
	"huangc28/go-ios-iap-vendor/db"
	"huangc28/go-ios-iap-vendor/internal/models"
	"log"
	"time"
)

type InventoryDAO struct {
	conn db.Conn
}

func NewInventoryDAO(conn db.Conn) *InventoryDAO {
	return &InventoryDAO{
		conn: conn,
	}
}

type GameItem struct {
	ProdID          string
	Receipt         string
	TransactionID   string
	TransactionDate time.Time
}

// https://stackoverflow.com/questions/51002790/locking-a-specific-row-in-postgres/52557413

func (dao *InventoryDAO) AddItemToInventory(item GameItem) error {
	query := `
INSERT INTO inventory (
	prod_id,
	transaction_id,
	receipt,
	available,
	transaction_time
) SELECT id, $1, $2, TRUE, $3 FROM product_info WHERE prod_id = $4;
	`
	if _, err := dao.conn.Exec(query, item.TransactionID, item.Receipt, item.TransactionDate, item.ProdID); err != nil {
		return err
	}

	return nil
}

// Find the first available stock in inventory. If no stock available, return no stock error.
// If we find the first available stock, row lock it first, update available to false, then retu the stock info.
func (dao *InventoryDAO) GetAvailableStock(prodID string) (*models.InventoryProduct, error) {
	rowLockQuery := `
SELECT
	inventory.*
FROM
	inventory
INNER JOIN product_info ON inventory.prod_id = product_info.id
WHERE
	available=true
AND
	product_info.prod_id = $1
ORDER BY created_at ASC LIMIT 1 FOR UPDATE OF inventory;

`
	m := models.InventoryProduct{}

	if err := dao.conn.QueryRowx(rowLockQuery, prodID).StructScan(&m); err != nil {
		return &m, err
	}

	// If selected stock is not available, let's select another available stock.
	if !m.Available {
		return dao.GetAvailableStock(prodID)
	}

	// If selected stock is available, update available to 'false'.
	updateAvailabilityQuery := `
UPDATE inventory
SET available = false
WHERE inventory.id = $1
RETURNING *;
`
	if err := dao.conn.QueryRowx(updateAvailabilityQuery, m.ID).StructScan(&m); err != nil {
		return nil, err
	}

	log.Printf("avai stock %v", m)

	return &m, nil
}
