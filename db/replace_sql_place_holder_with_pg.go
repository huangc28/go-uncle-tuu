package db

import (
	"strconv"
	"strings"
)

func ReplaceSQLPlaceHolderWithPG(sqlStr, searchPattern string) string {
	tmpCount := strings.Count(sqlStr, searchPattern)
	for m := 1; m <= tmpCount; m++ {
		sqlStr = strings.Replace(sqlStr, searchPattern, "$"+strconv.Itoa(m), 1)
	}

	return sqlStr
}
