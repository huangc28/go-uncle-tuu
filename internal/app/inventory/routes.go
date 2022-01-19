package inventory

import (
	"huangc28/go-ios-iap-vendor/config"
	"huangc28/go-ios-iap-vendor/internal/middlewares"

	"github.com/gin-gonic/gin"
)

func Routes(r *gin.RouterGroup) {
	g := r.Group(
		"/inventory",
		middlewares.JWTValidator(
			middlewares.JwtMiddlewareOptions{
				Secret: config.GetAppConf().APIJWTSecret,
			},
		),
	)

	g.GET(
		"/available-stock",
		GetAvailableStock,
	)
}
