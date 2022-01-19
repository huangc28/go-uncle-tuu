package inventory

import (
	"huangc28/go-ios-iap-vendor/db"
	"huangc28/go-ios-iap-vendor/internal/app/models"
	"log"
)

type InventoryDAO struct {
	conn db.Conn
}

func NewInventoryDAO(conn db.Conn) *InventoryDAO {
	return &InventoryDAO{
		conn: conn,
	}
}

// Find the first available stock in inventory. If no stock available, return no stock error.
// If we find the first available stock, row lock it first, update available to false, then return the stock info.
func (dao *InventoryDAO) GetAvailableStock(prodID string) (*models.Inventory, error) {
	log.Printf("prodID %s", prodID)

	rowLockQuery := `
SELECT
	inventory.*,
FROM
	inventory
INNER JOIN product_info ON inventory.prod_id = product_info.id
WHERE
	available=true
AND
	product_info.prod_id = $1
ORDER BY created_at ASC LIMIT 1 FOR UPDATE OF inventory;

`
	m := models.Inventory{}

	if err := dao.conn.QueryRowx(rowLockQuery, prodID).StructScan(&m); err != nil {
		return &m, err
	}

	// If selected stock is not available, let's select another available stock.
	if !m.Available.Bool {
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

	log.Printf("available stock %v", m)

	return &m, nil
}
