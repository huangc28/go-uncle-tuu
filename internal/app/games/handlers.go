package games

import (
	"huangc28/go-ios-iap-vendor/db"
	"huangc28/go-ios-iap-vendor/internal/app/models"
	"huangc28/go-ios-iap-vendor/internal/apperrors"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetGames(c *gin.Context) {
	gameDAO := NewGameDAO(db.GetDB())
	games, err := gameDAO.GetGames()

	if err != nil {
		c.AbortWithError(
			http.StatusInternalServerError,
			apperrors.NewErr(
				apperrors.FailedToGetGames,
				err.Error(),
			),
		)

		return
	}

	c.JSON(http.StatusOK, transformProducts(games))
}

type GetProductsParams struct {
	GameBundleID string `uri:"game_bundle_id" binding:"required"`
}

func GetProducts(c *gin.Context) {
	var params GetProductsParams
	if err := c.ShouldBindUri(&params); err != nil {
		c.AbortWithError(
			http.StatusBadRequest,
			apperrors.NewErr(
				apperrors.FailedToBindAPIBody,
				err.Error(),
			),
		)

		return
	}

	gameDAO := NewGameDAO(db.GetDB())
	prodOptions, err := gameDAO.GetProductInfoByGameBundleID(params.GameBundleID)

	if err != nil {
		c.AbortWithError(
			http.StatusBadRequest,
			apperrors.NewErr(
				apperrors.FailedToGetProductOptions,
				err.Error(),
			),
		)

		return
	}

	c.JSON(http.StatusOK, struct {
		Products []*models.ProductListOption `json:"products"`
	}{
		prodOptions,
	})
}
