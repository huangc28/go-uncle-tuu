package importer

import (
	"huangc28/go-ios-iap-vendor/db"
	"huangc28/go-ios-iap-vendor/internal/app/contracts"
	"huangc28/go-ios-iap-vendor/internal/app/models"

	cinternal "github.com/golobby/container/pkg/container"
)

type ProcurementDAO struct {
	conn db.Conn
}

func NewProcurementDAO(conn db.Conn) *ProcurementDAO {
	return &ProcurementDAO{
		conn: conn,
	}
}

func ProcurementDAOServiceProvider(c cinternal.Container) func() error {
	return func() error {
		c.Transient(func() contracts.ProcurementDAOer {
			return NewProcurementDAO(db.GetDB())
		})

		return nil
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

func (dao *ProcurementDAO) GetPendingProcurements() ([]*models.Procurement, error) {
	query := `
SELECT
	filename,
	status,
	failed_reason,
	created_at
FROM
	procurements
WHERE
	status=$1;
`

	rows, err := dao.conn.Queryx(query, models.ImportStatusPending)

	if err != nil {
		return nil, err
	}

	procs := make([]*models.Procurement, 0)
	for rows.Next() {
		var proc models.Procurement

		if err := rows.StructScan(&proc); err != nil {
			return nil, err
		}

		procs = append(procs, &proc)
	}

	return procs, nil
}

func (dao *ProcurementDAO) GetProcurements() ([]*models.Procurement, error) {
	query := `
SELECT
	filename,
	status,
	failed_reason,
	created_at
FROM
	procurements
`

	rows, err := dao.conn.Queryx(query)

	if err != nil {
		return nil, err
	}

	procs := make([]*models.Procurement, 0)
	for rows.Next() {
		var proc models.Procurement

		if err := rows.StructScan(&proc); err != nil {
			return nil, err
		}

		procs = append(procs, &proc)
	}

	return procs, nil
}
