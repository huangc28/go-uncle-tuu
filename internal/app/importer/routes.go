package importer

import "github.com/gin-gonic/gin"

func Routes(r *gin.RouterGroup) {
	g := r.Group("/importer")

	g.GET("/purchased-records", GetPurchasedRecordsHandler)
}
