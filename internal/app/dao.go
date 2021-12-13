package app

import (
	"huangc28/ios-inapp-trade/db"
	"huangc28/ios-inapp-trade/internal/models"
)

type ProdInfoDAO struct {
	conn db.Conn
}

type CreateProdInfoParams struct {
	BundleID string
	ProdID   string
	ProdName string
	ProdDesc string
	Price    float64
}

func NewProdInfoDAO(conn db.Conn) *ProdInfoDAO {
	return &ProdInfoDAO{
		conn: conn,
	}
}

func (dao *ProdInfoDAO) IsProdInfoExists(prodID, bundleID string) (bool, error) {
	query := `
SELECT EXISTS (
	SELECT 1
	FROM product_info
	WHERE prod_id = $1
	AND game_bundle_id = $2
);
	`
	var exists bool

	if err := dao.conn.QueryRowx(query, prodID, bundleID).Scan(&exists); err != nil {
		return exists, err
	}

	return exists, nil
}

func (dao *ProdInfoDAO) CreateProdInfo(params CreateProdInfoParams) (*models.ProductInfo, error) {
	query := `
INSERT INTO product_info (
	prod_id,
	prod_name,
	prod_desc,
	price,
	game_bundle_id
) VALUES ($1, $2, $3, $4, $5)
RETURNING
	prod_name,
	prod_desc,
	price,
	game_bundle_id
;
	`
	var m models.ProductInfo

	if err := dao.conn.QueryRowx(
		query,
		params.ProdID,
		params.ProdName,
		params.ProdDesc,
		params.Price,
		params.BundleID,
	).StructScan(&m); err != nil {
		return nil, err
	}

	return &m, nil
}

func (dao *ProdInfoDAO) fetchInventoryByBundleID(bundleID string) ([]*models.InventoryInfo, error) {
	query := `
SELECT
	product_info.prod_id,
	product_info.prod_name,
	product_info.price,
	COUNT(CASE WHEN inventory.available THEN 1 END) AS quantity
FROM
	product_info
LEFT JOIN inventory ON inventory.prod_id = product_info.id
WHERE
	product_info.game_bundle_id = $1
GROUP BY
	product_info.prod_id,
	product_info.prod_name,
	product_info.price
`

	ms := make([]*models.InventoryInfo, 0)

	rows, err := dao.conn.Queryx(query, bundleID)

	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var m models.InventoryInfo
		if err := rows.StructScan(&m); err != nil {
			return nil, err
		}

		ms = append(ms, &m)
	}

	return ms, nil
}
