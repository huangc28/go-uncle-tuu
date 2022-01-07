package auth

import (
	"huangc28/go-ios-iap-vendor/db"
	"huangc28/go-ios-iap-vendor/internal/app/models"
)

type AuthDao struct {
	conn db.Conn
}

func NewAuthDao(conn db.Conn) *AuthDao {
	return &AuthDao{
		conn: conn,
	}
}

func (dao *AuthDao) GetUserByUsername(username string) (models.User, error) {
	query := `
SELECT password, uuid FROM users
WHERE username = $1;
	`

	m := models.User{}

	if err := dao.conn.QueryRowx(query, username).StructScan(&m); err != nil {
		return m, err
	}

	return m, nil
}
