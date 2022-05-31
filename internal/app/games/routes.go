package games

import (
	"github.com/gin-gonic/gin"
)

func Routes(r *gin.RouterGroup) {
	g := r.Group("/games")

	g.GET("", GetGames)
	g.GET("/:game_bundle_id/products", GetProducts)
}
