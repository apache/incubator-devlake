package helper

import (
	"time"

	"gorm.io/datatypes"
)

// Table structure for raw data storage
type RawData struct {
	ID        uint64 `gorm:"primaryKey"`
	Params    string `gorm:"type:varchar(255);index"`
	Data      datatypes.JSON
	Url       string
	CreatedAt time.Time
}
