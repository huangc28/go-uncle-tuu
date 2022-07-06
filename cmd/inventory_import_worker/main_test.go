package main

import (
	"bufio"
	"io"
	"log"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"github.com/xuri/excelize/v2"
)

// Test import stock procurement successfully.
//   - Make sure test DB connection is available
//   - Remove duplicated stock
//   - Read excel content
//   - Test when header column is a merge of 2 columns.

type InventoryImportWorkerTestSuite struct {
	suite.Suite
	GoodTestExcelName   string
	GoodTestExcelReader io.ReadSeekCloser
}

func (s *InventoryImportWorkerTestSuite) SetupTest() {
	var err error
	s.GoodTestExcelName = "./test_procurement.xlsx"
	s.GoodTestExcelReader, err = os.Open(s.GoodTestExcelName)
	checkErr(err)

}

func (s *InventoryImportWorkerTestSuite) TearDownAllSuite() {
	s.GoodTestExcelReader.Close()
}

func (s *InventoryImportWorkerTestSuite) TearDownTestSuite() {
	s.GoodTestExcelReader.Seek(0, 0)
}

func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

// Read title from GoodTestExcelName. Retrieve it's position. Makesure the positions are correct.
func (s *InventoryImportWorkerTestSuite) TestGetTitlePositionAtTitleRow() {
	assert := assert.New(s.T())
	fexcel, err := excelize.OpenReader(s.GoodTestExcelReader)
	checkErr(err)
	rows, err := fexcel.GetRows("Sheet1")
	checkErr(err)

	// locate the actual position of each title in title row so we can locate the
	// correct data of the specific title.
	titleIndexMap := getTitlePositionAtTitleRow(rows[0])

	// assert header is at the right position
	assert.Equal(0, titleIndexMap["游戏名称"])
	assert.Equal(1, titleIndexMap["档位名称"])
	assert.Equal(2, titleIndexMap["档位价格"])
	assert.Equal(3, titleIndexMap["库存单号"])
	assert.Equal(4, titleIndexMap["入库时间"])
	assert.Equal(6, titleIndexMap["临时客户端凭证"])
	assert.Equal(7, titleIndexMap["客户端凭证"])
	assert.Equal(8, titleIndexMap["凭证生成时间"])
	assert.Equal(11, titleIndexMap["档位代码"])
}

func (s *InventoryImportWorkerTestSuite) TestExtractDataFromRowsSuccess() {
	assert := assert.New(s.T())
	fexcel, err := excelize.OpenReader(s.GoodTestExcelReader)
	checkErr(err)
	dataRows, err := fexcel.GetRows("Sheet1")
	titleRow := dataRows[0]
	titleIndexMap := getTitlePositionAtTitleRow(titleRow)
	dataRows = dataRows[1:]
	data, err := extractDataFromDataRows(dataRows, titleIndexMap)
	if err != nil {
		assert.Equal("採購單有商品沒有提供 \"档位代码\", 請再檢查一次 天堂M 4000 鑽石", err.Error())
	}

	prodUUIDAndIDMap["com.ncsoft.lineage2mtw_a5"] = 71
	data, err = extractDataFromDataRows(dataRows, titleIndexMap)
	if err != nil {
		assert.Equal("採購單有商品沒有提供 \"档位代码\", 請再檢查一次 天堂M 4000 鑽石", err.Error())
	}
	log.Printf("DEBUG %v", data)
}

func (s *InventoryImportWorkerTestSuite) TestParseAndImportProcurementToDBSuccess() {
	testProcFile, err := os.Open(s.GoodTestExcelName)
	if err != nil {
		log.Fatal("failed to open test_procurement.xlsx")
	}
	defer testProcFile.Close()
	fbuf := bufio.NewReader(testProcFile)
	parseAndImportProcurementToDB(fbuf)
}

func TestInventoryImportWorkerTestSuite(t *testing.T) {
	suite.Run(t, new(InventoryImportWorkerTestSuite))
}
