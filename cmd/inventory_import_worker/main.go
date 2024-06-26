package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"huangc28/go-ios-iap-vendor/config"
	"huangc28/go-ios-iap-vendor/db"
	"huangc28/go-ios-iap-vendor/internal/app/contracts"
	"huangc28/go-ios-iap-vendor/internal/app/deps"
	"huangc28/go-ios-iap-vendor/internal/app/models"
	"huangc28/go-ios-iap-vendor/internal/apperrors"
	"io"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/lib/pq"
	log "github.com/sirupsen/logrus"

	"cloud.google.com/go/storage"
	"github.com/jmoiron/sqlx"
	"github.com/xuri/excelize/v2"
	"google.golang.org/api/option"
)

type ImportWorkerError struct {
	ErrCode         string `json:"err_code"`
	ProblematicFile string `json:"problematic_file"`
	ProblematicItem string `json:"problematic_item"`
	Message         string `json:"err"`
}

func (e *ImportWorkerError) Error() string {
	return e.Message
}

var (
	prodUUIDAndIDMap = map[string]int{}
	prodNameList     = []string{}
)

func getProdNameIDMap() error {
	query := `
SELECT
	id,
	prod_id
FROM
	product_info;
`

	rows, err := db.GetDB().Queryx(query)

	if err != nil {
		return err
	}

	for rows.Next() {
		var (
			ID     int
			prodID string
		)

		if err := rows.Scan(&ID, &prodID); err != nil {
			return err
		}

		prodUUIDAndIDMap[prodID] = ID
		prodNameList = append(prodNameList, prodID)
	}

	return nil
}

// This worker fetch unimported excel file from GCS and imported to DB.
// The worker will first read column `import_status` (pending, imported, import_failed) from database.
// Read all file names that are with `pending` status. Fetch GCS object from google cloud storage into a io.Reader.
//
// We will use `excelize` to read column data from this io.Reader.

// TODO
//   - What if we failed to init reader from GCS object? I think we will continue to read next file.
//     but record the failed reason in DB for the corresponding file.
//   - What if parseAndImportProcurementToDB failed? Jot down the filename and the failed reason. continue
//     parsing the next file.
func init() {
	config.InitConfig()
	db.InitDB()
	if err := deps.Get().Run(); err != nil {
		log.Fatalf("failed to initialize dependency container %s", err.Error())
	}
}

func main() {
	ticker := time.NewTicker(10 * time.Second)
	quit := make(chan struct{})
	ctx := context.Background()

	go func() {
		for {
			select {
			case <-ticker.C:
				StartImporting(ctx)
			case <-quit:
				ticker.Stop()
				return
			}

		}
	}()

	// Listen to system signal to stop the ticker.
	systemQuit := make(chan os.Signal)
	signal.Notify(systemQuit, syscall.SIGINT, syscall.SIGTERM)
	<-systemQuit

	log.Info("graceful shutdown...")

	// Close worker
	close(quit)

	log.Info("shutdown complete...")

}

func StartImporting(ctx context.Context) {

	// Fetch all `pending` procurement
	var procDAO contracts.ProcurementDAOer
	deps.Get().Container.Make(&procDAO)

	procs, err := procDAO.GetPendingProcurements()

	if err != nil {
		log.Errorf("failed to get pending procurements %v", err)

		return
	}

	if len(procs) == 0 {
		log.Info("no pending procurement")

		return
	}

	// For each procurement name, create new reader from google cloud storage.
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
		log.Errorf("failed to initialize GCS client %v", err)

		return
	}

	// Retrieve `prod_name` and `id` for every item in `product_info`. Structure a prod_name:id map.
	// We will use fuzzy search to find the cloest match and retrieve it's id using this map
	err = getProdNameIDMap()

	if err != nil {
		log.Errorf("failed to retrieve products info %v", err)

		return
	}

	successImportedProcs := make([]string, 0)
	failedProcs := make([]*ImportWorkerError, 0)

	for _, proc := range procs {
		procReader, err := client.
			Bucket(config.GetAppConf().NewProcurementBucketName).
			Object(proc.Filename).
			NewReader(ctx)

		if err != nil {
			ime := &ImportWorkerError{
				ErrCode:         apperrors.FailedToInitGCSReader,
				ProblematicFile: proc.Filename,
				Message:         fmt.Sprintf("failed to init GSC reader, 請聯絡 bryan  %s", err.Error()),
			}

			failedProcs = append(failedProcs, ime)

			log.Errorf("failed to init io.Reader for bucket object %s, error: %v", proc.Filename, err)

			continue
		}

		if err := parseAndImportProcurementToDB(procReader); err != nil {
			ime := err.(*ImportWorkerError)
			ime.ProblematicFile = proc.Filename

			failedProcs = append(failedProcs, ime)

			log.Errorf("failed to parse and import procurement excel %s, error: %v", proc.Filename, err)

			continue
		}

		successImportedProcs = append(successImportedProcs, proc.Filename)
		log.Infof("done importing file %s", proc.Filename)
	}

	if len(successImportedProcs) > 0 {
		if err := UpdateSuccessProcurementStatus(successImportedProcs); err != nil {
			log.Errorf("failed to update status of successfully uploaded procurements %v", err)
		}
	}

	if len(failedProcs) > 0 {
		if err := UpdateFailedProcurementStatus(failedProcs); err != nil {
			log.Errorf("failed to update status of failed procurements %v", err)
		}
	}

}

func UpdateFailedProcurementStatus(failedProcs []*ImportWorkerError) error {
	procUpdates := make([]string, 0)

	for _, failedProc := range failedProcs {
		failedProcByte, err := json.Marshal(&failedProc)

		if err != nil {
			log.Errorf("failed to marshal importer failed struct to string %s", err.Error())
			continue
		}

		procUpdates = append(
			procUpdates,
			fmt.Sprintf(
				"('%s', '%s'::import_status, '%s')",
				failedProc.ProblematicFile,
				models.ImportStatusFailed,
				string(failedProcByte),
			),
		)

	}

	statement := strings.TrimSuffix(strings.Join(procUpdates, ","), ",")

	query := fmt.Sprintf(`
UPDATE procurements AS t
SET
    status = c.status_v,
    failed_reason = c.failed_reason_v
FROM (
	values
	%s
) AS c(filename_v, status_v, failed_reason_v)
WHERE t.filename = c.filename_v;
`, statement)

	if _, err := db.GetDB().Exec(query); err != nil {
		return err
	}

	return nil
}

func UpdateSuccessProcurementStatus(procNames []string) error {
	query, args, err := sqlx.In(
		`
UPDATE
	procurements
SET
	status=?
WHERE filename IN (?)
`,
		models.ImportStatusImported,
		procNames,
	)

	if err != nil {
		return err
	}

	query = db.GetDB().Rebind(query)

	if _, err := db.GetDB().Exec(query, args...); err != nil {
		return err
	}

	return nil
}

const (
	GameName         string = "游戏名称"
	GameItemName            = "档位名称"
	GameItemPrice           = "档位价格"
	TransactionID           = "库存单号"
	ImportAt                = "入库时间"
	TempReceipt             = "临时客户端凭证"
	Receipt                 = "客户端凭证"
	ReceiptCreatedAt        = "凭证生成时间"
	GameItemUUID            = "档位代码"
)

var TitleSignature = map[string]bool{
	GameName:         true,
	GameItemName:     true,
	GameItemPrice:    true,
	TransactionID:    true,
	ImportAt:         true,
	TempReceipt:      true,
	Receipt:          true,
	ReceiptCreatedAt: true,
	GameItemUUID:     true,
}

// Things that we need to check before importing to DB
//   - The first row is a title row. All title should match preset string. ex 遊戲名稱
//   - No duplicated Transaction ID.
//   - Data is located at index where the title is at in the index row. For example "游戏名称" is located at index 1 in the title row, the data of "游戏名称" must be at index 1 in the data row.
//
// If I found "游戏名称" at index 0. I mark it in a hash map.
// [
//   游戏名称 => 0 ---> indicates that the data of 游戏名称 is located at 0 in the data row.
//   档位名称 => 1 ---> indicates that the data of 档位名称 is located at 1 in the data row.
// ]
//
// TODO
//   what if "Sheet1" is not found?
func parseAndImportProcurementToDB(procReader io.Reader) error {
	f, err := excelize.OpenReader(procReader)

	if err != nil {
		return err
	}

	rows, err := f.GetRows("Sheet1")

	if err != nil {
		return &ImportWorkerError{
			Message: fmt.Sprintf("採購單有找到，但是沒有找到表單 %s, %s", "Sheet1", err.Error()),
			ErrCode: apperrors.FailedToGetSheet,
		}
	}

	titleRow := rows[0]
	titleIndexMap := getTitlePositionAtTitleRow(titleRow)

	for title, pos := range titleIndexMap {
		if pos == -1 {
			return &ImportWorkerError{
				Message: fmt.Sprintf("採購單的標題有錯喔，再檢查一次 %s ", title),
				ErrCode: apperrors.TitleSigatureNotFound,
			}
		}
	}

	dataRows := rows[1:]
	invs, err := extractDataFromDataRows(dataRows, titleIndexMap)

	if err != nil {
		return err
	}

	if err := importInventoriesToDB(invs); err != nil {
		return err
	}

	return nil
}

func expandDataRowLengthToAlignTitleRow(dataRows [][]string) {
	// We can not use range, since range only copy it's value
	// instead of changing the property of underlying slice.
	for i := 0; i < len(dataRows); i++ {
		dataRows[i] = dataRows[i][:cap(dataRows[i])]
	}
}

func importInventoriesToDB(invs []*models.Inventory) error {
	query := `
INSERT INTO inventory (prod_id, transaction_id, receipt, temp_receipt, transaction_time)
VALUES (:prod_id, :transaction_id, :receipt, :temp_receipt, :transaction_time)
	`

	if _, err := db.GetDB().NamedExec(
		query,
		invs,
	); err != nil {
		pqErr := err.(*pq.Error)
		if pqErr.Code == "23505" {
			return &ImportWorkerError{
				ErrCode: apperrors.FailedToImportDueToDuplicateInventory,
				Message: fmt.Sprint("庫存有重複單號。你的採購單某些商品庫存有了"),
			}
		}

		return &ImportWorkerError{
			ErrCode: apperrors.FailedToImportInventory,
			Message: err.Error(),
		}
	}

	return nil
}

func extractDataFromDataRows(dataRows [][]string, titleIndexMap map[string]int) ([]*models.Inventory, error) {
	expandDataRowLengthToAlignTitleRow(dataRows)
	invs := make([]*models.Inventory, 0)
	for _, dataRow := range dataRows {
		gameName := dataRow[titleIndexMap[GameName]]
		gameItemName := dataRow[titleIndexMap[GameItemName]]
		gameItemUUID := dataRow[titleIndexMap[GameItemUUID]]
		transactionID := dataRow[titleIndexMap[TransactionID]]
		transactionTime := dataRow[titleIndexMap[ReceiptCreatedAt]]
		tempReceipt := dataRow[titleIndexMap[TempReceipt]]
		receipt := dataRow[titleIndexMap[Receipt]]

		// If gameItemUUID is not provided, we need to notify business buy to provide it in excel
		if len(gameItemUUID) == 0 {
			return nil, &ImportWorkerError{
				ErrCode:         apperrors.IOSGameItemUUIDNotProvided,
				ProblematicItem: gameItemName,
				Message:         fmt.Sprintf("採購單有商品沒有提供 \"档位代码\", 請再檢查一次 %s %s", gameName, gameItemName),
			}

		}
		// Retrieve item ID in our database via prodUUIDAndIDMap. If item not found, that means given item info has not been collected.
		gameItemIDInDB, exists := prodUUIDAndIDMap[gameItemUUID]

		if !exists {
			return nil, &ImportWorkerError{
				ErrCode:         apperrors.GameItemInfoHasNotBeenCollected,
				ProblematicItem: gameItemName,
				Message: fmt.Sprintf(
					"採購單上有些商品還沒有採集: %s %s。請確定採購單上的商品都採集後再上傳一次",
					gameItemName,
					gameName,
				),
			}
		}

		transactionTimeInt, err := strconv.ParseInt(transactionTime, 10, 64)

		if err != nil {
			return nil, &ImportWorkerError{
				ErrCode:         apperrors.FailedToParseTransactionTime,
				ProblematicItem: gameItemName,
				Message:         fmt.Sprintf("商品 %s, 憑證生成時間有錯 %v", transactionID, err),
			}
		}

		transactionTimeInst := time.Unix(transactionTimeInt/1000, 0)

		inv := &models.Inventory{
			ProdID: sql.NullInt32{
				Valid: true,
				Int32: int32(gameItemIDInDB),
			},
			TransactionID: sql.NullString{
				Valid:  true,
				String: transactionID,
			},
			Receipt: sql.NullString{
				Valid:  true,
				String: receipt,
			},
			TempReceipt: sql.NullString{
				Valid:  true,
				String: tempReceipt,
			},
			TransactionTime: transactionTimeInst,
		}

		invs = append(invs, inv)
	}

	return invs, nil
}

func getTitlePositionAtTitleRow(row []string) map[string]int {
	titleIndexMap := map[string]int{
		GameName:         -1,
		GameItemName:     -1,
		GameItemPrice:    -1,
		TransactionID:    -1,
		ImportAt:         -1,
		TempReceipt:      -1,
		Receipt:          -1,
		ReceiptCreatedAt: -1,
		GameItemUUID:     -1,
	}

	for idx, titleCell := range row {
		if _, exists := TitleSignature[titleCell]; !exists {
			continue
		}

		titleIndexMap[titleCell] = idx
	}

	return titleIndexMap
}
