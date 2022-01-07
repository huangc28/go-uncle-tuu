package auth

import "github.com/gin-gonic/gin"

func Routes(r *gin.RouterGroup) {
	g := r.Group("/auth")

	g.POST("/login", LoginHandler)
}
