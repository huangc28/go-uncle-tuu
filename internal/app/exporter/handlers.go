package exporter

import (
	"huangc28/go-ios-iap-vendor/internal/app/contracts"
	"huangc28/go-ios-iap-vendor/internal/apperrors"
	"huangc28/go-ios-iap-vendor/internal/pkg/requestbinder"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golobby/container/pkg/container"
)

type ReportReceivedBody struct {
	UUID string `form:"uuid" json:"uuid" binding:"required"`
}

func ReportReceived(c *gin.Context, depCon container.Container) {
	body := &ReportReceivedBody{}
	if err := requestbinder.Bind(c, &body); err != nil {
		c.AbortWithError(
			http.StatusBadRequest,
			apperrors.NewErr(
				apperrors.FailedToBindAPIBody,
				err.Error(),
			),
		)

		return
	}

	var (
		inventoryDao contracts.InventoryDAOer
		userDao      contracts.UserDAOer
	)

	depCon.Make(&inventoryDao)
	depCon.Make(&userDao)

	user, err := userDao.GetUserByUUID(c.GetString("user_uuid"), "id")

	if err != nil {
		c.AbortWithError(
			http.StatusBadRequest,
			apperrors.NewErr(apperrors.FailedToGetUserByUUID),
		)

		return
	}

	// Makesure the stock is reserved for that user.
	isReserved, err := inventoryDao.IsStockReservedForUser(body.UUID, user.ID)

	if err != nil {
		c.AbortWithError(
			http.StatusBadRequest,
			apperrors.NewErr(apperrors.FailedToCheckStockReservedForUser),
		)

		return
	}

	if !isReserved {
		c.AbortWithError(
			http.StatusBadRequest,
			apperrors.NewErr(apperrors.StockIsNotReservedForTheUser),
		)

		return
	}

	if err := inventoryDao.MarkStockAsDelivered(body.UUID); err != nil {
		c.AbortWithError(
			http.StatusInternalServerError,
			apperrors.NewErr(apperrors.FailedToMarkStockAsDeliver),
		)

		return
	}

	c.JSON(http.StatusOK, struct{}{})
}
