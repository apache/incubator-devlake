package db

import (
	"fmt"
)

func InsertThing(tableName string, thingToInsert map[string]interface{}) (int64, error) {
	keys, values := GetKeysAndValues(thingToInsert)

	rawSql := fmt.Sprintf("INSERT INTO %v (%v) VALUES(%v);", tableName, keys, values)

	res, err := Db.Exec(rawSql)
	if err != nil {
		return 0, err
	}
	return res.RowsAffected()
}

func GetKeysAndValues(thingToInsert map[string]interface{}) (keys string, values string) {
	count := 1
	endString := ", "
	for key, value := range thingToInsert {
		if count == len(thingToInsert) {
			endString = ""
		}
		keys += fmt.Sprintf("%v%v", key, endString)
		values += fmt.Sprintf("\"%v\"%v", value, endString)
		count++
	}
	return keys, values
}
