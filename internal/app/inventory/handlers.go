package inventory

import (
	"database/sql"
	"huangc28/go-ios-iap-vendor/db"
	"huangc28/go-ios-iap-vendor/internal/app/contracts"
	"huangc28/go-ios-iap-vendor/internal/apperrors"
	"huangc28/go-ios-iap-vendor/internal/pkg/requestbinder"
	"net/http"
	"time"

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

	// TODO: Only white listed user can access to all products even reserved products.
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

func GetAvailableStock(c *gin.Context, depCon container.Container) {
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

	var userDAO contracts.UserDAOer
	depCon.Make(&userDAO)

	user, err := userDAO.GetUserByUUID(c.GetString("user_uuid"), "id")

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

	dao := NewInventoryDAO(db.GetDB())

	stock, err := dao.GetAvailableStock(body.ProdID, user.ID)

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

type AddItemToInventoryBody struct {
	// ProdName example: "arktw_diamond_1".
	ProdID string `json:"prod_id" form:"prod_id" binding:"required"`

	// Receipt receipt string after successful transaction.
	Receipt         string    `json:"receipt" form:"receipt" binding:"required"`
	TempReceipt     string    `json:"temp_receipt" form:"temp_receipt" binding:"required"`
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

	dao := NewInventoryDAO(db.GetDB())

	// Add game item to inventory.
	if err := dao.AddItemToInventory(GameItem{
		ProdID:          body.ProdID,
		Receipt:         body.Receipt,
		TempReceipt:     body.TempReceipt,
		TransactionID:   body.TransactionID,
		TransactionDate: body.TransactionDate,
	}); err != nil {
		c.AbortWithError(
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

type AssignStocksParams struct {
	// Username indicates assignee's username
	Username string `form:"username" json:"username" binding:"required"`

	Stocks []StockParam `form:"stocks" json:"stocks" binding:"required"`
}

type StockParam struct {
	GameBundleID string `json:"game_bundle_id"`
	ProdID       string `json:"prod_id"`
	Quantity     int    `json:"quantity"`
}

// - Check if quantity is enough for each product
//
// Check the length the available inventory is greater than the number of requested products.
// TODO: add assignment_inventory table pivot table to trace the deliver status between assignmnts and stocks.
func assignStocks(c *gin.Context, depCon container.Container) {
	var stockParams AssignStocksParams

	if err := requestbinder.Bind(c, &stockParams); err != nil {
		c.AbortWithError(
			http.StatusInternalServerError,
			apperrors.NewErr(
				apperrors.FailedToBindAPIBody,
				err.Error(),
			),
		)

		return
	}

	var userDAO contracts.UserDAOer
	depCon.Make(&userDAO)
	assignee, err := userDAO.GetUserByUsername(stockParams.Username, "id")

	if err != nil {
		c.AbortWithError(
			http.StatusInternalServerError,
			apperrors.NewErr(
				apperrors.FailedToGetUserByUsername,
				err.Error(),
			),
		)

		return
	}

	// Retrieve all matched products in inventory with available=true.
	prodIDs := make([]string, len(stockParams.Stocks))
	for _, stock := range stockParams.Stocks {
		prodIDs = append(prodIDs, stock.ProdID)
	}
	invDAO := NewInventoryDAO(db.GetDB())
	avaiStocks, err := invDAO.GetAvailableStocksForProdIDs(prodIDs)

	if err != nil {
		c.AbortWithError(
			http.StatusInternalServerError,
			apperrors.NewErr(
				apperrors.FailedToGetAvailableStocksForProdIDs,
				err.Error(),
			),
		)

		return
	}

	// check if quantity in stock is enough for assigning.
	stockEnoughErr := IsQuantityInStockEnoughForAssigning(stockParams.Stocks, avaiStocks)

	if stockEnoughErr != nil {
		c.AbortWithError(
			http.StatusBadRequest,
			apperrors.NewErr(
				apperrors.NotEnoughStocks,
				stockEnoughErr.Error(),
			),
		)

		return
	}

	// Start assigning stock to user.
	avaiStockUUIDs := make([]string, len(avaiStocks))
	for _, avaiStock := range avaiStocks {
		avaiStockUUIDs = append(avaiStockUUIDs, avaiStock.UUID)
	}

	if err := invDAO.AssignStockToUser(int(assignee.ID), avaiStockUUIDs); err != nil {
		c.AbortWithError(
			http.StatusInternalServerError,
			apperrors.NewErr(
				apperrors.FailedToAssignStocksToUser,
			),
		)

		return
	}

	c.JSON(http.StatusOK, struct{}{})
}
