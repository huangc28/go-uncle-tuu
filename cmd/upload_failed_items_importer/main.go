package main

import (
	"bufio"
	"encoding/json"
	"huangc28/go-ios-iap-vendor/config"
	"huangc28/go-ios-iap-vendor/db"
	"huangc28/go-ios-iap-vendor/internal/app/deps"
	"huangc28/go-ios-iap-vendor/internal/app/inventory"
	"io"
	"io/fs"
	"log"
	"os"
	"path"
	"strings"
	"time"
)

func init() {
	config.InitConfig()

	db.InitDB()
	if err := deps.Get().Run(); err != nil {
		log.Fatalf("failed to initialize dependency container %s", err.Error())
	}
}

func ReadFailedItemDir() ([]fs.DirEntry, error) {

	failedItemDirPath := path.Join(
		config.GetProjRootPath(),
		config.GetAppConf().ImportFailedFileDirPath,
	)

	// Read files from the directory path.
	f, err := os.Open(failedItemDirPath)

	if err != nil {
		return nil, err
	}

	defer f.Close()

	fileInfos, err := f.ReadDir(-1)

	if err != nil {
		return nil, err
	}

	return fileInfos, nil
}

var NewLineSeperator []byte = []byte{0xA, 0xA}

func NewlineSlicer(bytes []byte, needles []byte) int {
	var i int = -1

	if len(bytes) < len(needles) {
		return -1
	}
OuterLoop:
	for idx, b := range bytes {
		// Try to match first byte in the needles.
		if needles[0] != b {
			//i = -1
			continue
		} else {
			i = idx

		InnerLoop:
			for j := 1; j < len(needles); j++ {
				// What if bytes length is less then i+j? That means it's impossible to find the exact `
				// match of `needles` in `bytes` array. We simply break the loop and set the result
				// to be not found.
				if i+j >= len(bytes) {
					i = -1

					break InnerLoop
				}

				if bytes[i+j] != needles[j] {
					i = -1

					break InnerLoop
				}
			}

			// If the next three bytes all matches we simply break the outer loop.
			if i >= 0 {
				break OuterLoop
			}
		}
	}

	return i
}

func SplitByNewLine(data []byte, atEOF bool) (int, []byte, error) {
	// If we are at the end of file and we have no data in the current
	// data buffer, we request for more data.
	// This can happend if we are at the last segment of the data but coincidently the data buffer
	// does not have any data.
	if atEOF && len(data) == 0 {
		return 0, nil, nil
	}

	// If we found byte for new empty line (\r\n) in the data buffer [0xA, 0xA], we have a full event.
	// we can flush the data out of the buffer

	if i := NewlineSlicer(data, NewLineSeperator); i >= 0 {
		return i + 1, data[0:i], nil
	}

	// If we're at EOF, we have a final event
	if atEOF {
		return len(data), data, nil
	}

	//log.Printf("DEBUG raw byte %v", data)
	//log.Printf("DEBUG byte in string %v", string(data))

	return 0, nil, nil
}

type Stock struct {
	Receipt         string `json:"receipt"`
	TransactionID   string `json:"transaction_id"`
	ProdID          string `json:"prod_id"`
	TempReceipt     string `json:"temp_receipt"`
	TransactionDate string `json:"transaction_date"`
}

const (
	ISO8601Layout = "2006-01-02T15:04:05"
)

// This worker read upload failed items from import failed files, import them to inventory table.
// TODO
//  - use go routine to read file content.
//  - use loggers to instead of simply printing logs on the screen.
func main() {
	fileInfos, err := ReadFailedItemDir()

	if err != nil {
		log.Fatalf("failed to read failed item directory %v", err.Error())
	}

	inventoryDAO := inventory.NewInventoryDAO(db.GetDB())

	// Select product IDs from "product_info" to form a key-value pair hash like the following:
	//
	// {
	//    arktw_diamond_2: 2
	//    kgtw.cash.1: 4
	//    ...
	// }
	prodIDIDMap, err := inventoryDAO.GetProdInfoIDProdIDKeyValuePair()

	if err != nil {
		log.Printf("failed to get prod id to id map %v", prodIDIDMap)

		return
	}

	// sucessfulImportedFiles stores those files that has been imported to inventory successfully.
	sucessfulImportedFiles := make([]string, 0)
	gameItems := make([]*inventory.GameItem, 0)

	var closeAndResetFile = func(f *os.File) {
		f.Seek(0, 0)
		f.Close()
	}

	// Iterate through each failed items file.
FileLoop:
	for _, file := range fileInfos {
		// Read each file content and write data to a struct.
		// Read content from path and split it by empty line.
		filepath := path.Join(
			config.GetProjRootPath(),
			config.GetAppConf().ImportFailedFileDirPath,
			file.Name(),
		)

		f, err := os.Open(filepath)

		if err != nil {
			log.Printf("failed to open file %v %v", file.Name(), err)

			continue
		}

		scanner := bufio.NewScanner(f)
		scanner.Split(SplitByNewLine)

	ContentLoop:
		for scanner.Scan() {
			dataStr := strings.Trim(string(scanner.Bytes()), string([]byte{0xA}))

			r := strings.NewReader(dataStr)
			var stock Stock

			// If any panic error happened, we continue process the next file. Do not mark this file
			// as successfully imported.
			if err := json.NewDecoder(r).Decode(&stock); err != nil {
				if err == io.EOF {
					continue ContentLoop
				}

				log.Printf("failed to decode stock to struct for file, skipping to the next file %v %v", file.Name(), err)

				closeAndResetFile(f)
				continue FileLoop
			}

			json.NewDecoder(r).Decode(&stock)

			t, err := time.Parse(ISO8601Layout, stock.TransactionDate)

			if err != nil {
				log.Printf("failed to parse transaction date for file, skipping to the next file %v %v", file.Name(), err)

				closeAndResetFile(f)
				continue FileLoop
			}

			// Parse transaction_date
			gameItem := inventory.GameItem{
				ProdID:          stock.ProdID,
				Receipt:         stock.Receipt,
				TempReceipt:     stock.TempReceipt,
				TransactionID:   stock.TransactionID,
				TransactionDate: t,
			}

			gameItems = append(gameItems, &gameItem)
		}

		closeAndResetFile(f)
		sucessfulImportedFiles = append(sucessfulImportedFiles, filepath)
	}

	if err := inventoryDAO.BatchAddItemsToInventory(gameItems, prodIDIDMap); err != nil {
		log.Printf("failed to add items to inventory %v", err)

		return
	}

	log.Println("done importing game items, removing files...")
	for _, successfulImportedFile := range sucessfulImportedFiles {
		if err := os.Remove(successfulImportedFile); err != nil {
			log.Printf("failed to remove file %v, %v", successfulImportedFile, err)

			continue
		}
	}
}
