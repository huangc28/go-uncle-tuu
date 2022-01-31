package inventory

import (
	"errors"
	"huangc28/go-ios-iap-vendor/db"
	"huangc28/go-ios-iap-vendor/internal/app/contracts"
	"huangc28/go-ios-iap-vendor/internal/app/models"
	"log"
	"time"

	cintrnal "github.com/golobby/container/pkg/container"
)

type InventoryDAO struct {
	conn db.Conn
}

func NewInventoryDAO(conn db.Conn) *InventoryDAO {
	return &InventoryDAO{
		conn: conn,
	}
}

func InventoryDaoServiceProvider(c cintrnal.Container) func() error {
	return func() error {
		c.Transient(func() contracts.InventoryDAOer {
			return NewInventoryDAO(db.GetDB())
		})

		return nil
	}
}

func (dao *InventoryDAO) GetUserReservedStockByUUID(prodID string, userID int) (*models.ReservedStockInfo, error) {
	query := `
SELECT
	COUNT(inventory.id) AS reserved_num,
	product_info.prod_name,
	product_info.price
FROM
	inventory
INNER JOIN product_info ON inventory.prod_id = product_info.id
WHERE
	inventory.available=true
AND
	inventory.reserved_for_user = $2
AND
	product_info.prod_id = $1
GROUP BY
	product_info.prod_name,
	product_info.price
	`
	var m models.ReservedStockInfo

	if err := dao.conn.QueryRowx(query, prodID, userID).StructScan(&m); err != nil {

		return nil, err
	}

	return &m, nil
}

// Find the first available stock in inventory. If no stock available, return no stock error.
// If we find the first available stock, row lock it first, update available to false, then return the stock info.
func (dao *InventoryDAO) GetAvailableStock(prodID string) (*models.Inventory, error) {
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

func (dao *InventoryDAO) IsStockReservedForUser(stockUUID string, userID int64) (bool, error) {
	query := `
SELECT
	available,
	delivered,
	reserved_for_user
FROM
	inventory
WHERE
	uuid = $1;
	`

	var m models.Inventory

	if err := dao.conn.QueryRowx(query, stockUUID).StructScan(&m); err != nil {
		return false, err
	}

	if m.Available.Bool {
		return false, errors.New("stock is still available, it's not exported yet")
	}

	if m.Delivered == models.DeliveredStatusDelivered {
		return false, errors.New("stock has already delivered")
	}

	return int64(m.ReservedForUser.Int32) == userID, nil
}

func (dao *InventoryDAO) markStockDeliverStatus(stockUUID string, status models.DeliveredStatus) error {
	query := `
UPDATE inventory
SET delivered = $1
WHERE UUID = $2;
`
	if _, err := dao.conn.Exec(query, status, stockUUID); err != nil {
		return err
	}

	return nil
}

func (dao *InventoryDAO) MarkStockAsDelivered(stockUUID string) error {
	return dao.markStockDeliverStatus(stockUUID, models.DeliveredStatusDelivered)
}

func (dao *InventoryDAO) MarkStockAsNotDelivered(stockUUID string) error {
	return dao.markStockDeliverStatus(stockUUID, models.DeliveredStatusNotDelivered)
}

type GameItem struct {
	ProdID          string
	Receipt         string
	TempReceipt     string
	TransactionID   string
	TransactionDate time.Time
}

// https://stackoverflow.com/questions/51002790/locking-a-specific-row-in-postgres/52557413

func (dao *InventoryDAO) AddItemToInventory(item GameItem) error {
	log.Printf("temp receipt %v", item)
	query := `
INSERT INTO inventory (
	prod_id,
	transaction_id,
	receipt,
	temp_receipt,
	available,
	transaction_time
) SELECT id, $1, $2, $3, TRUE, $4 FROM product_info WHERE prod_id = $5;
	`
	if _, err := dao.conn.Exec(query, item.TransactionID, item.Receipt, item.TempReceipt, item.TransactionDate, item.ProdID); err != nil {
		return err
	}

	return nil
}
