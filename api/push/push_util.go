package push

import (
	"github.com/merico-dev/lake/models"
)

func InsertThing(tableName string, thingToInsert map[string]interface{}) (int64, error) {
	tx := models.Db.Table(tableName).Create(thingToInsert)
	if tx.Error != nil {
		return 0, tx.Error
	}
	return tx.RowsAffected, nil
}
