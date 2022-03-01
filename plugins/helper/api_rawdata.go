package helper

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
)

type RawData struct {
	Uuid      uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	Params    string    `gorm:"type:varchar(255);index"`
	Data      datatypes.JSON
	CreatedAt time.Time
}

func GetRawTableCreationSqls(table string) []string {
	return []string{
		fmt.Sprintf(`
			CREATE TABLE IF NOT EXISTS %s (
				uuid VARCHAR(50) NOT NULL,
				params varchar(255) NOT NULL,
				data json NULL,
				created_at DATETIME DEFAULT current_timestamp NULL,
				CONSTRAINT %s_PK PRIMARY KEY (uuid),
				INDEX %s_params_IDX (params)
			)
			`, table, table, table),
	}
}

func GetRawTableDeletionSql(table string) string {
	return fmt.Sprintf(`
	DELETE FROM %s WHERE params = ?
	`, table)
}
