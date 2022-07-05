package main

import "strings"

func getDataType(dataType string) string {
	starrocksDatatype := dataType
	if strings.HasPrefix(dataType, "varchar") {
		starrocksDatatype = "string"
	} else if strings.HasPrefix(dataType, "datetime") {
		starrocksDatatype = "datetime"
	} else if strings.HasPrefix(dataType, "bigint") {
		starrocksDatatype = "bigint"
	} else if dataType == "longtext" || dataType == "text" || dataType == "longblob" {
		starrocksDatatype = "string"
	} else if dataType == "tinyint(1)" {
		starrocksDatatype = "boolean"
	}
	return starrocksDatatype
}
