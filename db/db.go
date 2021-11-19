package db

import (
	"database/sql"
	"fmt"
	"strings"

	log "github.com/sirupsen/logrus"

	"github.com/jmoiron/sqlx"
	"github.com/jmoiron/sqlx/reflectx"

	_ "github.com/lib/pq"
)

type Conn interface {
	Exec(query string, args ...interface{}) (sql.Result, error)
	Query(query string, args ...interface{}) (*sql.Rows, error)
	Queryx(query string, args ...interface{}) (*sqlx.Rows, error)
	QueryRow(query string, args ...interface{}) *sql.Row
	QueryRowx(query string, args ...interface{}) *sqlx.Row
	Prepare(query string) (*sql.Stmt, error)
	Preparex(query string) (*sqlx.Stmt, error)
}

var dbInstance *sqlx.DB

func InitDB() {
	dsn := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		"127.0.0.1",
		5435,
		"postgres",
		"1234",
		"inapp_trade",
	)

	driver, err := sqlx.Open("postgres", dsn)

	if err != nil {
		log.Fatalf("failed to open connection %s", err.Error)
	}

	dbInstance = driver
	dbInstance.Mapper = reflectx.NewMapperFunc("json", strings.ToLower)

	log.Info("Database connected!")
}

func GetDB() *sqlx.DB {
	return dbInstance
}
