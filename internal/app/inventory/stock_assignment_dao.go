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

func (dao *StockAssignmentDAO) CreateAssignment(assigneeID int) (*models.StockAssignment, error) {
	query := `
INSERT INTO stock_assignments (assignee_id)
VALUES ($1)
RETURNING *;
	`
	var sa models.StockAssignment

	if err := dao.conn.QueryRowx(query, assigneeID).StructScan(&sa); err != nil {
		return nil, err
	}

	return &sa, nil
}

func (dao *StockAssignmentDAO) GetAssignmentStatus() ([]*models.StockAssigmentStatus, error) {
	query := `
SELECT
	u.username AS assignee_name,
	sa.uuid AS assignment_uuid,
	sa.created_at,
	pi2.prod_name,
	g.readable_name AS game_name,
	i.uuid as stock_uuid,
	i.delivered,
	i.available
FROM
	stock_assignments sa
INNER JOIN users u ON sa.assignee_id = u.id
INNER JOIN inventory i ON i.assignment_id = sa.id
INNER JOIN product_info pi2 ON pi2.id = i.prod_id
INNER JOIN games g ON g.game_bundle_id = pi2.game_bundle_id;
	`

	rows, err := dao.conn.Queryx(query)

	if err != nil {
		return nil, err
	}

	as := make([]*models.StockAssigmentStatus, 0)

	for rows.Next() {
		var a models.StockAssigmentStatus
		if err := rows.StructScan(&a); err != nil {
			return nil, err
		}
		as = append(as, &a)
	}

	return as, nil
}
