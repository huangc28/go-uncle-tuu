package inventory

import (
	"huangc28/go-ios-iap-vendor/db"
	"huangc28/go-ios-iap-vendor/internal/app/models"
)

type StockAssignmentDAO struct {
	conn db.Conn
}

func NewStockAssignmentDAO(conn db.Conn) *StockAssignmentDAO {
	return &StockAssignmentDAO{
		conn: conn,
	}
}

func (dao *StockAssignmentDAO) SetConn(conn db.Conn) {
	dao.conn = conn
}

func (dao *StockAssignmentDAO) CreateAssignment() (*models.StockAssignment, error) {
	query := `
INSERT INTO stock_assignments (updated_at)
VALUES (null);
RETURNING *;
	`
	var sa *models.StockAssignment

	if err := dao.conn.QueryRowx(query).StructScan(sa); err != nil {
		return nil, err
	}

	return sa, nil
}
