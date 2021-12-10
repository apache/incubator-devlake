package push

import (
	"github.com/merico-dev/lake/models"
)

func InsertRow(tableName string, rowToInsert map[string]interface{}) (int64, error) {
	tx := models.Db.Table(tableName).Create(rowToInsert)
	if tx.Error != nil {
		return 0, tx.Error
	}
	return tx.RowsAffected, nil
}
