package collector

import "github.com/gin-gonic/gin"

func Routes(r *gin.RouterGroup) {
	g := r.Group("/collector")

	g.POST("/collect-product-info", collectProductInfoHandler)
}
