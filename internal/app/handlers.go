package app

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"

	"huangc28/ios-inapp-trade/db"
	"huangc28/ios-inapp-trade/internal/pkg/requestbinder"
)

type CollectProductInfoBody struct {
	ProdID   string  `form:"prod_id"`
	ProdName string  `form:"prod_name"`
	ProdDesc string  `form:"prod_desc"`
	Quantity int     `form:"quantity"`
	Price    float64 `form:"price"`
	BundleID string  `form:"bundle_id"`
}

type ErrorMessage struct {
	Err string `json:"error"`
}

func CollectProductInfoHandler(c *gin.Context) {
	body := CollectProductInfoBody{}

	// bind request handler
	if err := requestbinder.Bind(c, &body); err != nil {
		c.JSON(http.StatusInternalServerError, ErrorMessage{
			Err: err.Error(),
		})

		return
	}

	log.Printf("DEBUG body %v", body)

	// Store those information to DB.

	prodInfoDao := NewProdInfoDAO(db.GetDB())
	prodInfo, err := prodInfoDao.CreateProdInfoIfNotExists(
		CreateProdInfoParams{
			BundleID: body.BundleID,
			ProdID:   body.ProdID,
			ProdName: body.ProdName,
			ProdDesc: body.ProdDesc,
			Price:    body.Price,
		},
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorMessage{
			Err: err.Error(),
		})

		return
	}

	log.Printf("DEBUG product info %v", prodInfo)

	c.JSON(http.StatusOK, struct {
		OK bool
	}{true})

}
