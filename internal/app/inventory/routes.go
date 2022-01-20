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

	g.GET("/reserved-stock", func(c *gin.Context) {
		GetReservedStock(c, depCon)
	})

	g.GET(
		"/available-stock",
		GetAvailableStock,
	)
}
