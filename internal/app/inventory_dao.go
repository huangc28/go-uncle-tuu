package app

import (
	"huangc28/go-ios-iap-vendor/db"
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
