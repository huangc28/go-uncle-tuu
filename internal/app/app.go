package app

import (
	"huangc28/go-ios-iap-vendor/internal/app/auth"
	"huangc28/go-ios-iap-vendor/internal/app/collector"
	"huangc28/go-ios-iap-vendor/internal/app/deps"
	"huangc28/go-ios-iap-vendor/internal/app/exporter"
	"huangc28/go-ios-iap-vendor/internal/app/importer"
	"huangc28/go-ios-iap-vendor/internal/app/inventory"
	"huangc28/go-ios-iap-vendor/internal/apperrors"
	"huangc28/go-ios-iap-vendor/internal/middlewares"

	"github.com/gin-gonic/gin"
)

func StartApp(e *gin.Engine) {
	e.Use(middlewares.CORSMiddlewares())
	e.Use(middlewares.ResponseLogger)
	e.Use(apperrors.HandleError())
	rv1 := e.Group("/v1")

	rv1.GET(
		"/inventory",
		fetchInventoryHandler,
	)

	collector.Routes(rv1)

	auth.Routes(rv1)

	importer.Routes(rv1)

	inventory.Routes(rv1, deps.Get().Container)

	exporter.Routes(rv1, deps.Get().Container)
}
