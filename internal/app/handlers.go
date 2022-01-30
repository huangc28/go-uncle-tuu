package app

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"huangc28/go-ios-iap-vendor/db"
	"huangc28/go-ios-iap-vendor/internal/apperrors"
	"huangc28/go-ios-iap-vendor/internal/pkg/requestbinder"
)

type FetchInventoryBody struct {
	BundleID string `form:"bundle_id" binding:"required"`
}

func fetchInventoryHandler(c *gin.Context) {
	body := FetchInventoryBody{}

	if err := requestbinder.Bind(c, &body); err != nil {
		c.JSON(
			http.StatusInternalServerError,
			apperrors.NewErr(
				apperrors.FailedToBindAPIBody,
				err.Error(),
			),
		)

		return
	}

	// Retrieve inventory status according given bundle_id
	dao := NewProdInfoDAO(db.GetDB())
	ms, err := dao.fetchInventoryByBundleID(body.BundleID)

	if err != nil {
		c.JSON(
			http.StatusInternalServerError,
			apperrors.NewErr(
				apperrors.FailedToFetchInventoryInfo,
				err.Error(),
			),
		)

		return
	}

	c.JSON(http.StatusOK, TrfInventory(ms))
}

type AddItemToInventoryBody struct {
	// ProdName example: "arktw_diamond_1".
	ProdID string `json:"prod_id" form:"prod_id" binding:"required"`

	// Receipt receipt string after successful transaction.
	Receipt         string    `json:"receipt" form:"receipt" binding:"required"`
	TransactionID   string    `json:"transaction_id" form:"transaction_id" binding:"required"`
	TransactionDate time.Time `json:"transaction_date" form:"transaction_date" binding:"required"`
}

func addItemToInventory(c *gin.Context) {
	body := AddItemToInventoryBody{}

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

	// Add game item to inventory.
	dao := NewInventoryDAO(db.GetDB())
	if err := dao.AddItemToInventory(GameItem{
		ProdID:          body.ProdID,
		Receipt:         body.Receipt,
		TransactionID:   body.TransactionID,
		TransactionDate: body.TransactionDate,
	}); err != nil {
		c.JSON(
			http.StatusInternalServerError,
			apperrors.NewErr(
				apperrors.FailedToAddItemToInventory,
				err.Error(),
			),
		)

		return
	}

	c.JSON(http.StatusOK, struct{}{})
}
