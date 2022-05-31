package games

import (
	"huangc28/go-ios-iap-vendor/db"
	"huangc28/go-ios-iap-vendor/internal/app/models"
)

type GameDAO struct {
	conn db.Conn
}

func NewGameDAO(conn db.Conn) *GameDAO {
	return &GameDAO{
		conn: conn,
	}
}

// GetGames get all supported games.
func (dao *GameDAO) GetGames() ([]*models.Game, error) {
	query := `
SELECT
	game_bundle_id,
	readable_name
FROM
	games
WHERE
	supported = true;
	`

	rows, err := dao.conn.Queryx(query)

	if err != nil {
		return nil, err
	}

	games := make([]*models.Game, 0)
	for rows.Next() {
		var game models.Game
		if err := rows.StructScan(&game); err != nil {
			return nil, err
		}

		games = append(games, &game)
	}

	return games, nil
}

func (dao *GameDAO) GetProductInfoByGameBundleID(gameBundleID string) ([]*models.ProductListOption, error) {
	query := `
SELECT
	product_info.prod_name,
	product_info.prod_id AS prod_bundle_id,
	COUNT(inv.prod_id) AS num_in_stock
FROM
	product_info
LEFT JOIN inventory AS inv ON product_info.id = inv.prod_id AND inv.available = true
WHERE
	product_info.game_bundle_id=$1
GROUP BY
	product_info.prod_id, product_info.prod_name;
`
	rows, err := dao.conn.Queryx(query, gameBundleID)

	if err != nil {
		return nil, err
	}

	prodOptions := make([]*models.ProductListOption, 0)

	for rows.Next() {
		var prodOption models.ProductListOption
		if err := rows.StructScan(&prodOption); err != nil {
			return nil, err
		}
		prodOptions = append(prodOptions, &prodOption)
	}

	return prodOptions, nil
}
