package collector

import (
	"huangc28/go-ios-iap-vendor/db"
	"huangc28/go-ios-iap-vendor/internal/apperrors"
	"huangc28/go-ios-iap-vendor/internal/pkg/requestbinder"
	"net/http"

	"github.com/gin-gonic/gin"
)

type CollectProductInfoBody struct {
	ProdID   string  `form:"prod_id"`
	ProdName string  `form:"prod_name"`
	ProdDesc string  `form:"prod_desc"`
	Quantity int     `form:"quantity"`
	Price    float64 `form:"price"`
	BundleID string  `form:"bundle_id"`
}

func collectProductInfoHandler(c *gin.Context) {
	body := CollectProductInfoBody{}

	// bind request handler
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

	// Store those information to DB.
	prodInfoDao := NewProdInfoDAO(db.GetDB())

	exists, err := prodInfoDao.IsProdInfoExists(body.ProdID, body.BundleID)

	if exists {
		c.AbortWithError(
			http.StatusInternalServerError,
			apperrors.NewErr(apperrors.DuplicatedProductInfo),
		)

		return
	}

	if err != nil {
		c.AbortWithError(
			http.StatusInternalServerError,
			apperrors.NewErr(apperrors.FailedToCheckProductExists, err.Error()),
		)

		return
	}

	prodInfo, err := prodInfoDao.CreateProdInfo(
		CreateProdInfoParams{
			BundleID: body.BundleID,
			ProdID:   body.ProdID,
			ProdName: body.ProdName,
			ProdDesc: body.ProdDesc,
			Price:    body.Price,
		},
	)

	if err != nil {
		c.AbortWithError(
			http.StatusInternalServerError,
			apperrors.NewErr(
				apperrors.FailedToCreateProdInfo,
				err.Error(),
			),
		)

		return
	}

	c.JSON(http.StatusOK, prodInfo)
}
