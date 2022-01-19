package inventory

import (
	"database/sql"
	"huangc28/go-ios-iap-vendor/db"
	"huangc28/go-ios-iap-vendor/internal/apperrors"
	"huangc28/go-ios-iap-vendor/internal/pkg/requestbinder"
	"net/http"

	"github.com/gin-gonic/gin"
)

type GetAvailableStockBody struct {
	// ProdID unique product id user intend to export from inventory.
	ProdID string `json:"prod_id" form:"prod_id" binding:"required"`
}

func GetAvailableStock(c *gin.Context) {
	body := GetAvailableStockBody{}

	if err := requestbinder.Bind(c, &body); err != nil {
		c.JSON(
			http.StatusBadRequest,
			apperrors.NewErr(
				apperrors.FailedToBindAPIBody,
				err.Error(),
			),
		)

		return
	}

	dao := NewInventoryDAO(db.GetDB())

	stock, err := dao.GetAvailableStock(body.ProdID)

	if err == sql.ErrNoRows {
		c.JSON(
			http.StatusInternalServerError,
			apperrors.NewErr(
				apperrors.NoAvailableProductFound,
			),
		)

		return
	}

	if err != nil {
		c.JSON(
			http.StatusInternalServerError,
			apperrors.NewErr(
				apperrors.FailedToGetAvailableStock,
				err.Error(),
			),
		)

		return

	}

	c.JSON(http.StatusOK, TrfAvailableStock(*stock))
}
