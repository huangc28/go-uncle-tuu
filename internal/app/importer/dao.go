package importer

import (
	"huangc28/go-ios-iap-vendor/db"
	"huangc28/go-ios-iap-vendor/internal/app/models"
)

type ImporterDAO struct {
	conn db.Conn
}

func NewImporterDAO(conn db.Conn) *ImporterDAO {
	return &ImporterDAO{
		conn: conn,
	}
}

func (dao *ImporterDAO) GetPurchasedRecords(bundleID string, perPage, offset int) ([]models.PurchaseRecord, error) {
	query := `
SELECT
	inventory.transaction_id,
	inventory.transaction_time,
	inventory.uuid,
	product_info.prod_name
FROM
	inventory
INNER JOIN product_info ON inventory.prod_id = product_info.id
WHERE
	product_info.game_bundle_id=$1
ORDER BY transaction_time DESC
LIMIT $2
OFFSET $3;
`
	rows, err := dao.conn.Queryx(query, bundleID, perPage, offset)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	ms := make([]models.PurchaseRecord, 0)

	for rows.Next() {
		m := models.PurchaseRecord{}

		if err := rows.StructScan(&m); err != nil {
			return nil, err
		}

		ms = append(ms, m)
	}

	return ms, nil
}
