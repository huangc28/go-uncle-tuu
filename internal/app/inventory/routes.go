package inventory

import (
	"huangc28/go-ios-iap-vendor/config"
	"huangc28/go-ios-iap-vendor/internal/middlewares"

	"github.com/gin-gonic/gin"
	"github.com/golobby/container/pkg/container"
)

func Routes(r *gin.RouterGroup, depCon container.Container) {
	g := r.Group(
		"/inventory",
	)

	// Deprecated Not used anywhere in frontend
	// TODO move this API to exporter
	// Check if there are enough quantity of the stock that he/she wants to export.
	g.GET(
		"/reserved-stock",
		middlewares.JWTValidator(
			middlewares.JwtMiddlewareOptions{
				Secret: config.GetAppConf().APIJWTSecret,
			},
		),
		func(c *gin.Context) {
			GetReservedStock(c, depCon)
		},
	)

	// Check if there are any stock that is eligible for exporting to user's game account.
	// If a reserved stock is found for the user, deliver the stock.
	g.GET(
		"/available-stock",
		middlewares.JWTValidator(middlewares.JwtMiddlewareOptions{
			Secret: config.GetAppConf().APIJWTSecret,
		}),
		func(c *gin.Context) {
			GetAvailableStock(c, depCon)
		},
	)

	// TODO: move this API to importer
	g.POST("/add-item-to-inventory", addItemToInventory)

	g.POST("/assign-stocks", func(c *gin.Context) {
		assignStocks(c, depCon)
	})
}
