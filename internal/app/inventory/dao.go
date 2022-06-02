package inventory

import (
	"errors"
	"fmt"
	"huangc28/go-ios-iap-vendor/db"
	"huangc28/go-ios-iap-vendor/internal/app/contracts"
	"huangc28/go-ios-iap-vendor/internal/app/models"
	"log"
	"strings"
	"time"

	cintrnal "github.com/golobby/container/pkg/container"
	"github.com/jmoiron/sqlx"
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
func (dao *InventoryDAO) GetAvailableStock(prodID string, userID int64) (*models.Inventory, error) {
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
AND
	reserved_for_user = $2
ORDER BY created_at ASC LIMIT 1 FOR UPDATE OF inventory;

`
	m := models.Inventory{}

	if err := dao.conn.QueryRowx(rowLockQuery, prodID, userID).StructScan(&m); err != nil {
		return &m, err
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
// AddItemToInventory add stock to inventory IF the transaction id of the stock isn't already exists in the inventory.
func (dao *InventoryDAO) AddItemToInventory(item GameItem) error {
	query := `
INSERT INTO inventory (
	prod_id,
	transaction_id,
	receipt,
	temp_receipt,
	available,
	transaction_time
) SELECT id, $1, $2, $3, TRUE, $4 FROM product_info WHERE prod_id = $5 ON CONFLICT (transaction_id) DO NOTHING;
	`

	if _, err := dao.conn.Exec(query, item.TransactionID, item.Receipt, item.TempReceipt, item.TransactionDate, item.ProdID); err != nil {
		return err
	}

	return nil
}

type ProdInfoIDProdIDKeyValuePair struct {
	ID     int    `json:"id"`
	ProdID string `json:"prod_id"`
}

func (dao *InventoryDAO) GetProdInfoIDProdIDKeyValuePair() (map[string]int, error) {
	idProdIDMap := make(map[string]int)

	query := `
SELECT
	id,
	prod_id
FROM
	product_info
	`

	rows, err := dao.conn.Queryx(query)

	if err != nil {
		return idProdIDMap, err
	}

	for rows.Next() {
		o := ProdInfoIDProdIDKeyValuePair{}

		if err := rows.StructScan(&o); err != nil {
			return idProdIDMap, err
		}

		idProdIDMap[o.ProdID] = o.ID
	}

	return idProdIDMap, nil
}

func (dao *InventoryDAO) BatchAddItemsToInventory(gameItems []*GameItem, prodIDIDMap map[string]int) error {
	sqlStr := "INSERT INTO inventory(prod_id, receipt, temp_receipt, transaction_id, transaction_time) VALUES "
	vals := []interface{}{}

	for _, gameItem := range gameItems {
		id, prodIDExists := prodIDIDMap[gameItem.ProdID]

		if !prodIDExists {
			log.Printf("DEBUG product id: %s does not exist, skipping", gameItem.ProdID)

			continue
		}

		sqlStr += "(?, ?, ?, ?, ?),"
		vals = append(
			vals,
			id,
			gameItem.Receipt,
			gameItem.TempReceipt,
			gameItem.TransactionID,
			gameItem.TransactionDate,
		)
	}

	if len(gameItems) <= 0 {
		log.Println("after prod id filtering, there are no game items to import")

		return nil
	}

	sqlStr = strings.TrimSuffix(sqlStr, ",")
	pgStr := db.ReplaceSQLPlaceHolderWithPG(sqlStr, "?")
	pgStr = fmt.Sprintf("%s ON CONFLICT(transaction_id) DO NOTHING", pgStr)

	stmt, err := dao.conn.Prepare(pgStr)

	if err != nil {
		return err
	}

	_, err = stmt.Exec(vals...)

	if err != nil {
		return err
	}

	return nil
}

type AvailableStockInfo struct {
	ProdID          string    `json:"prod_id"`
	UUID            string    `json:"uuid"`
	TransactionTime time.Time `json:"transaction_time"`
}

func (dao *InventoryDAO) GetAvailableStocksForProdIDs(gameIDs []string) ([]*AvailableStockInfo, error) {
	query, args, err := sqlx.In(
		`
SELECT
	product_info.prod_id,
	inventory.uuid,
	inventory.transaction_time
FROM product_info
INNER JOIN inventory ON product_info.id = inventory.prod_id
WHERE product_info.prod_id IN (?)
AND inventory.available = TRUE
AND reserved_for_user IS NULL
ORDER BY transaction_time ASC;
	`, gameIDs)

	if err != nil {
		return nil, err
	}

	query = db.GetDB().Rebind(query)

	rows, err := db.GetDB().Queryx(query, args...)

	if err != nil {
		return nil, err
	}

	prods := make([]*AvailableStockInfo, 0)

	for rows.Next() {
		var prod AvailableStockInfo
		if err := rows.StructScan(&prod); err != nil {
			return nil, err
		}
		prods = append(prods, &prod)
	}

	return prods, nil
}

func (dao *InventoryDAO) AssignStockToUser(assigneeID int, prodUUIDs []string) error {
	query, args, err := sqlx.In(`
UPDATE inventory
SET reserved_for_user=?
WHERE uuid IN(?)
`, assigneeID, prodUUIDs)

	if err != nil {
		return err
	}

	query = db.GetDB().Rebind(query)

	if _, err := db.GetDB().Exec(query, args...); err != nil {
		return err
	}

	return nil
}
