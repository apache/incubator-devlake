package push

import (
	"fmt"

	"github.com/merico-dev/lake/models"
)

func InsertThing(tableName string, thingToInsert map[string]interface{}) (int64, error) {
	keys, values := GetKeysAndValues(thingToInsert)

	rawSql := fmt.Sprintf("INSERT INTO %v (%v) VALUES(%v);", tableName, keys, values)

	tx := models.Db.Exec(rawSql)
	if tx.Error != nil {
		return 0, tx.Error
	}
	return tx.RowsAffected, nil
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
