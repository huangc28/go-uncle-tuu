package app

import (
	"database/sql"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"huangc28/ios-inapp-trade/db"
	"huangc28/ios-inapp-trade/internal/apperrors"
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

func collectProductInfoHandler(c *gin.Context) {
	body := CollectProductInfoBody{}

	// bind request handler
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

	// Store those information to DB.
	prodInfoDao := NewProdInfoDAO(db.GetDB())

	log.Printf("DEBUG %v %v", body.ProdID, body.BundleID)
	exists, err := prodInfoDao.IsProdInfoExists(body.ProdID, body.BundleID)

	if err != nil {
		c.JSON(
			http.StatusInternalServerError,
			apperrors.NewErr(apperrors.FailedToCheckProductExists, err.Error()),
		)

		return
	}

	if exists {
		c.JSON(
			http.StatusInternalServerError,
			apperrors.NewErr(apperrors.DuplicatedProductInfo),
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
		c.JSON(
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

type GetAvailableStockBody struct {
	// ProdID unique product id user intend to export from inventory.
	ProdID string `json:"prod_id" form:"prod_id" binding:"required"`
}

// TODO we need to add a column "delivered" to indicate if the product is delivered. This variable is set via client notification.
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
