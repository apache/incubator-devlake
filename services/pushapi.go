package services

func InsertRow(tableName string, rowToInsert map[string]interface{}) (int64, error) {
	tx := db.Table(tableName).Create(rowToInsert)
	if tx.Error != nil {
		return 0, tx.Error
	}
	return tx.RowsAffected, nil
}
