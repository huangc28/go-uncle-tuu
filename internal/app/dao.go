package app

import (
	"huangc28/ios-inapp-trade/db"
	"huangc28/ios-inapp-trade/internal/models"
	"log"
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

func (dao *ProdInfoDAO) IsProdInfoExists(prodName, bundleID string) {
	query := `
EXISTS (
	SELECT COUNT(id)
	FROM product_info
	WHERE prod_name = $1
	AND game_bundle_id = $2
);
	`

	log.Printf("DEBUG %v", query)

}

func (dao *ProdInfoDAO) CreateProdInfoIfNotExists(params CreateProdInfoParams) (*models.ProductInfo, error) {
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
