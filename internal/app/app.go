package app

import (
	"huangc28/go-ios-iap-vendor/internal/app/auth"
	"huangc28/go-ios-iap-vendor/internal/middlewares"

	"github.com/gin-gonic/gin"
)

func StartApp(e *gin.Engine) {
	e.Use(middlewares.ResponseLogger)
	rv1 := e.Group("/v1")

	rv1.POST(
		"/collect-product-info",
		collectProductInfoHandler,
	)

	rv1.GET(
		"/inventory",
		fetchInventoryHandler,
	)

	rv1.POST(
		"/add-item-to-inventory",
		addItemToInventory,
	)

	rv1.GET("/stock", GetAvailableStock)

	auth.Routes(rv1)
}
