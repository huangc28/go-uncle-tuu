package exporter

import (
	"huangc28/go-ios-iap-vendor/config"
	"huangc28/go-ios-iap-vendor/internal/middlewares"

	"github.com/gin-gonic/gin"

	"github.com/golobby/container/pkg/container"
)

func Routes(r *gin.RouterGroup, depCon container.Container) {
	g := r.Group(
		"/exporter",
		middlewares.JWTValidator(
			middlewares.JwtMiddlewareOptions{
				Secret: config.GetAppConf().APIJWTSecret,
			},
		),
	)

	g.POST("/report-received",
		func(c *gin.Context) {
			ReportReceived(c, depCon)
		},
	)

	g.POST("/report-not-received", func(c *gin.Context) {
		ReportUnreceived(c, depCon)
	})
}
