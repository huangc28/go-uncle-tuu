package importer

import (
	"fmt"
	"huangc28/go-ios-iap-vendor/config"
	"huangc28/go-ios-iap-vendor/db"
	"huangc28/go-ios-iap-vendor/internal/apperrors"
	"huangc28/go-ios-iap-vendor/internal/pkg/requestbinder"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type GetPurchasedRecordsBody struct {
	BundleID string `form:"bundle_id" json:"bundle_id" binding:"required,gt=0"`
	Page     int    `form:"page,default=0" json:"page,default=0"`
	PerPage  int    `form:"per_page,default=5" json:"per_page,default=5"`
}

func GetPurchasedRecordsHandler(c *gin.Context) {
	body := GetPurchasedRecordsBody{}

	if err := requestbinder.Bind(c, &body); err != nil {
		c.JSON(
			http.StatusBadRequest,
			apperrors.NewErr(
				apperrors.FailedToBindAPIBody,
				err.Error(),
			),
		)

		return
	}

	dao := NewImporterDAO(db.GetDB())

	recs, err := dao.GetPurchasedRecords(body.BundleID, body.PerPage, body.Page*body.PerPage)

	if err != nil {
		c.JSON(
			http.StatusInternalServerError,
			apperrors.NewErr(
				apperrors.FailedToGetPurchasedRecords,
				err.Error(),
			),
		)

		return
	}

	trfRecs := TtfPurchaseRecords(recs)

	c.JSON(http.StatusOK, struct {
		Data []TrfPurchaseRecord `json:"data"`
	}{trfRecs})
}

func fileNameWithoutExtension(fileName string) string {
	if pos := strings.LastIndexByte(fileName, '.'); pos != -1 {
		return fileName[:pos]
	}
	return fileName
}

func UploadFailedList(c *gin.Context) {
	// check if directory exists, if not make a directory.
	dirPath := filepath.Join(
		config.GetProjRootPath(),
		config.GetAppConf().ImportFailedFileDirPath,
	)

	if err := os.MkdirAll(dirPath, os.ModePerm); err != nil {
		c.AbortWithError(
			http.StatusInternalServerError,
			apperrors.NewErr(
				apperrors.FailedToCreateImportFailedFileDir,
				err.Error(),
			),
		)

		return
	}

	file, _ := c.FormFile("filename")

	log.Printf("file1 %v", file)
	dstName := filepath.Join(
		dirPath,
		fmt.Sprintf(
			"%s_%d%s",
			fileNameWithoutExtension(file.Filename),
			time.Now().Unix(),
			filepath.Ext(file.Filename),
		),
	)

	// Upload the file to specific dst.
	c.SaveUploadedFile(file, dstName)

	c.JSON(http.StatusOK, struct{}{})
}
