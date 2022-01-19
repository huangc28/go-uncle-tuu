package inventory

import "github.com/gin-gonic/gin"

func Routes(r *gin.RouterGroup) {
	g := r.Group("/inventory")

	g.GET("/available-stock", GetAvailableStock)
}
