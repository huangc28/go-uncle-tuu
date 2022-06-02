package users

import (
	"fmt"
	"huangc28/go-ios-iap-vendor/db"
	"huangc28/go-ios-iap-vendor/internal/app/contracts"
	"huangc28/go-ios-iap-vendor/internal/app/models"

	cintrnal "github.com/golobby/container/pkg/container"
)

type UserDAO struct {
	conn db.Conn
}

func NewUserDAO(conn db.Conn) *UserDAO {
	return &UserDAO{
		conn: conn,
	}
}

func UserDaoServiceProvider(c cintrnal.Container) func() error {
	return func() error {
		c.Transient(func() contracts.UserDAOer {
			return NewUserDAO(db.GetDB())
		})

		return nil
	}
}

func (dao *UserDAO) GetUserByUUID(UUID string, fields ...string) (*models.User, error) {
	baseQuery := `
SELECT %s
FROM users
WHERE uuid = $1;
	`
	query := fmt.Sprintf(baseQuery, db.ComposeFieldsSQLString(fields...))

	var user models.User

	if err := dao.conn.QueryRowx(query, UUID).StructScan(&user); err != nil {
		return nil, err
	}

	return &user, nil
}

func (dao *UserDAO) DisableExport(userID int64) error {
	query := `
UPDATE
	users
SET
	can_be_exported = false
WHERE
	id = $1;
`
	_, err := dao.conn.Exec(query, userID)

	return err
}

func (dao *UserDAO) GetUserByUsername(username string, fields ...string) (*models.User, error) {

	baseQuery := `
SELECT %s
FROM users
WHERE username = $1;
	`
	query := fmt.Sprintf(baseQuery, db.ComposeFieldsSQLString(fields...))

	var user models.User

	if err := dao.conn.QueryRowx(query, username).StructScan(&user); err != nil {
		return nil, err
	}

	return &user, nil
}
