package inventory

import (
	"database/sql"
	"huangc28/go-ios-iap-vendor/db"
	"huangc28/go-ios-iap-vendor/internal/app/contracts"
	"huangc28/go-ios-iap-vendor/internal/apperrors"
	"huangc28/go-ios-iap-vendor/internal/pkg/requestbinder"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golobby/container/pkg/container"
)

type GetReservedStockBody struct {
	ProdID string `form:"prod_id"`
}

// Retrieve all user reserved stock.
func GetReservedStock(c *gin.Context, depCon container.Container) {
	userUUID := c.GetString("user_uuid")

	body := GetReservedStockBody{}

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

	var userDAO contracts.UserDAOer
	depCon.Make(&userDAO)

	user, err := userDAO.GetUserByUUID(userUUID, "id")

	if err != nil {
		c.AbortWithError(
			http.StatusInternalServerError,
			apperrors.NewErr(
				apperrors.FailedToGetUserByUUID,
				err.Error(),
			),
		)

		return
	}

	invDAO := NewInventoryDAO(db.GetDB())

	// Only white listed user can access to all products even reserved products.
	reservedStockInfo, err := invDAO.GetUserReservedStockByUUID(body.ProdID, int(user.ID))

	if err == sql.ErrNoRows {
		c.AbortWithError(
			http.StatusBadRequest,
			apperrors.NewErr(apperrors.NoReservedStockAvailable),
		)

		return
	}

	if err != nil {
		c.AbortWithError(
			http.StatusInternalServerError,
			apperrors.NewErr(
				apperrors.FailedToGetReservedStock,
				err.Error(),
			),
		)

		return

	}

	c.JSON(http.StatusOK, reservedStockInfo)
}

type GetAvailableStockBody struct {
	// ProdID unique product id user intend to export from inventory.
	ProdID string `json:"prod_id" form:"prod_id" binding:"required"`
}

func GetAvailableStock(c *gin.Context) {
	body := GetAvailableStockBody{}

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

	dao := NewInventoryDAO(db.GetDB())

	stock, err := dao.GetAvailableStock(body.ProdID)

	if err == sql.ErrNoRows {
		c.AbortWithError(
			http.StatusInternalServerError,
			apperrors.NewErr(
				apperrors.NoAvailableProductFound,
			),
		)

		return
	}

	if err != nil {
		c.AbortWithError(
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
