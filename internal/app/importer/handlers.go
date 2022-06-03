package importer

import (
	"context"
	"fmt"
	"huangc28/go-ios-iap-vendor/config"
	"huangc28/go-ios-iap-vendor/db"
	"huangc28/go-ios-iap-vendor/internal/app/contracts"
	"huangc28/go-ios-iap-vendor/internal/apperrors"
	"huangc28/go-ios-iap-vendor/internal/pkg/requestbinder"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"cloud.google.com/go/storage"
	"github.com/gin-gonic/gin"
	"github.com/golobby/container/pkg/container"
	gcsenhancer "github.com/huangc28/gcs_enhancer"
	"google.golang.org/api/option"
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

type ProcurementForm struct {
	File *multipart.FileHeader `form:"procurement" binding:"required"`
}

// TODO
//  - Create a table to log who uploads the procurement sheet, long with import status.
//  - Implement a worker to import procurement sheet into database.
func UploadProcurement(c *gin.Context) {
	form := &ProcurementForm{}

	if err := c.ShouldBind(form); err != nil {
		c.AbortWithError(
			http.StatusBadRequest,
			apperrors.NewErr(
				apperrors.FailedToBindAPIBody,
				err.Error(),
			),
		)

		return
	}

	ctx := context.Background()
	client, err := storage.NewClient(
		ctx,
		option.WithCredentialsFile(
			fmt.Sprintf("%s/%s",
				config.GetProjRootPath(),
				config.GetAppConf().GCSGoogleServiceAccountName,
			),
		),
	)

	if err != nil {
		c.AbortWithError(
			http.StatusInternalServerError,
			apperrors.NewErr(
				apperrors.FailedToInitGoogleStorageClient,
				err.Error(),
			),
		)
		return
	}

	f, err := form.File.Open()

	if err != nil {
		c.AbortWithError(
			http.StatusInternalServerError,
			apperrors.NewErr(
				apperrors.FailedToOpenUploadedFile,
				err.Error(),
			),
		)

		return
	}

	defer f.Close()
	enhancer := gcsenhancer.NewGCSEnhancer(client, config.GetAppConf().NewProcurementBucketName)
	uoInfo, err := enhancer.Upload(
		ctx,
		f,
		gcsenhancer.AppendUnixTimeStampToFilename(form.File.Filename),
		gcsenhancer.UploadOptions{},
	)

	if err != nil {
		c.AbortWithError(
			http.StatusInternalServerError,
			apperrors.NewErr(
				apperrors.FailedToUploadFileToGCS,
				err.Error(),
			),
		)

		return
	}

	procDAO := NewProcurementDAO(db.GetDB())
	procRec, err := procDAO.CreateProcurement(uoInfo.Filename)

	if err != nil {
		c.AbortWithError(
			http.StatusInternalServerError,
			apperrors.NewErr(
				apperrors.FailedToCreateProcurementRecord,
				err.Error(),
			),
		)

		return
	}

	c.JSON(http.StatusOK, struct {
		UploadedFilename string
		ImportStatus     string
		CreatedAt        time.Time
	}{
		uoInfo.Filename,
		string(procRec.Status),
		procRec.CreatedAt,
	})
}

func GetProcurements(c *gin.Context, depCon container.Container) {
	var procDAO contracts.ProcurementDAOer
	depCon.Make(&procDAO)

	procs, err := procDAO.GetProcurements()

	if err != nil {
		c.AbortWithError(
			http.StatusInternalServerError,
			apperrors.NewErr(
				apperrors.FailedToGetProcurements,
				err.Error(),
			),
		)

		return
	}

	c.JSON(http.StatusOK, TrfProcurements(procs))
}
