package db

import (
	"net/http"

	"github.com/jmoiron/sqlx"
)

type FormatResp struct {
	Err            error
	ErrCode        string
	HttpStatusCode int
	Response       interface{}
}

const (
	FailedToBeginTx  = "0000010"
	FailedToCommitTx = "0000011"
)

type TxFuncFormatResp func(tx *sqlx.Tx) FormatResp

func TransactWithFormatStruct(db *sqlx.DB, txFunc TxFuncFormatResp) FormatResp {
	tx, err := db.Beginx()

	if err != nil {
		return FormatResp{
			Err:     err,
			ErrCode: FailedToBeginTx,
		}
	}

	fnResp := txFunc(tx)

	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p)
		} else if fnResp.Err != nil {
			tx.Rollback()

			// If http status code is not set, default to be `500`
			if fnResp.HttpStatusCode == 0 {
				fnResp.HttpStatusCode = http.StatusInternalServerError
			}
		} else {
			fnResp.Err = tx.Commit()

			if fnResp.Err != nil {
				fnResp.ErrCode = FailedToCommitTx
				fnResp.HttpStatusCode = http.StatusInternalServerError
			}
		}
	}()

	return fnResp
}
