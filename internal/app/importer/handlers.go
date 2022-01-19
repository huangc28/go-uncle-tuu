package importer

import (
	"huangc28/go-ios-iap-vendor/db"
	"huangc28/go-ios-iap-vendor/internal/apperrors"
	"huangc28/go-ios-iap-vendor/internal/pkg/requestbinder"
	"net/http"

	"github.com/gin-gonic/gin"
)

type GetPurchasedRecordsBody struct {
	BundleID string `form:"bundle_id" json:"bundle_id" binding:"required,gt=0"`
	Page     int    `form:"page,default=0" json:"page,default=0"`
	PerPage  int    `form:"per_page,default=5" json:"per_page,default=5"`
}

func GetPurchasedRecordsHandler(c *gin.Context) {
	body := GetPurchasedRecordsBody{}

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

	dao := NewImporterDAO(db.GetDB())

	recs, err := dao.GetPurchasedRecords(body.BundleID, body.PerPage, body.Page*body.PerPage)

	if err != nil {
		c.JSON(
			http.StatusInternalServerError,
			apperrors.NewErr(
				apperrors.FailedToGetPurchasedRecords,
				err.Error(),
			),
		)

		return
	}

	trfRecs := TtfPurchaseRecords(recs)

	c.JSON(http.StatusOK, struct {
		Data []TrfPurchaseRecord `json:"data"`
	}{trfRecs})
}
