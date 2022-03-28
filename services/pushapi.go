package services

func InsertRow(table string, rows []map[string]interface{}) (int64, error) {
	tx := db.Table(table).Create(rows)
	if tx.Error != nil {
		return 0, tx.Error
	}
	return tx.RowsAffected, nil
}
