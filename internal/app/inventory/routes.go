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
		middlewares.JWTValidator(
			middlewares.JwtMiddlewareOptions{
				Secret: config.GetAppConf().APIJWTSecret,
			},
		),
	)

	// TODO: move this API to exporter
	// Check if there are enough quantity of the stock that he/she
	// wants to export.
	g.GET("/reserved-stock", func(c *gin.Context) {
		GetReservedStock(c, depCon)
	})

	g.GET("/available-stock", GetAvailableStock)

	// TODO: move this API to importer
	g.POST("/add-item-to-inventory", addItemToInventory)
}
