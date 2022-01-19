package app

import (
	"huangc28/go-ios-iap-vendor/internal/app/auth"
	"huangc28/go-ios-iap-vendor/internal/app/importer"
	"huangc28/go-ios-iap-vendor/internal/app/inventory"
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

	auth.Routes(rv1)

	importer.Routes(rv1)

	inventory.Routes(rv1)
}
