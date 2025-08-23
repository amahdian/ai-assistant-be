package binding

import (
	"fmt"
	"strings"
)

const paramSeparator = " "

func extractTableAndColumnNames(param string) (tableName string, columnName string) {
	parts := strings.Split(param, paramSeparator)
	if len(parts) != 2 {
		panic(fmt.Sprintf("unsupported params format for 'exists' validator. the expected format is 'table_name,column_name'"))
	}

	tableName = parts[0]
	columnName = parts[1]
	return
}
