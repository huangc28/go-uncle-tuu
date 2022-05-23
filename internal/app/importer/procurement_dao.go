package importer

import (
	"huangc28/go-ios-iap-vendor/db"
	"huangc28/go-ios-iap-vendor/internal/app/models"
)

type ProcurementDAO struct {
	conn db.Conn
}

func NewProcurementDAO(conn db.Conn) *ProcurementDAO {
	return &ProcurementDAO{
		conn: conn,
	}
}

func (dao *ProcurementDAO) CreateProcurement(filename string) (models.Procurement, error) {
	query := `
INSERT INTO procurements (
	filename
) VALUES($1)
RETURNING *;
	`
	var procurement models.Procurement

	if err := dao.conn.QueryRowx(query, filename).StructScan(&procurement); err != nil {
		return procurement, err
	}

	return procurement, nil
}
