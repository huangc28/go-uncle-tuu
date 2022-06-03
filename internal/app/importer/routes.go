package importer

import (
	"github.com/gin-gonic/gin"
	"github.com/golobby/container/pkg/container"
)

func Routes(r *gin.RouterGroup, depCon container.Container) {
	g := r.Group("/importer")

	g.GET("/purchased-records", GetPurchasedRecordsHandler)

	g.POST("/upload-failed-list", UploadFailedList)

	g.POST("/upload-procurement", UploadProcurement)

	g.GET("/procurements", func(c *gin.Context) {
		GetProcurements(c, depCon)
	})
}
