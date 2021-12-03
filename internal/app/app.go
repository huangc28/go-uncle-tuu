package app

import (
	"huangc28/ios-inapp-trade/internal/middlewares"

	"github.com/gin-gonic/gin"
)

func StartApp(e *gin.Engine) {
	e.Use(middlewares.ResponseLogger)
	rv1 := e.Group("/v1")

	rv1.POST(
		"/collect-product-info",
		CollectProductInfoHandler,
	)

	rv1.GET(
		"/inventory",
		fetchInventoryHandler,
	)
}
