package app

import (
	"net/http"

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
