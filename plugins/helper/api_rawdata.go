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

// embedded fields for tool layer tables
type ExtractedRawData struct {
	// can be used for flushing outdated records from table
	RawDataParams string `gorm:"column:_raw_data_params,type:varchar(255);index"`
	RawDataTable  uint64 `gorm:"column:_raw_data_table"`
	// can be used for debugging
	RawDataId uint64 `gorm:"column:_raw_data_id"`
	// we can store record index into this field, which is helpful for debugging
	RawDataRemark string `gorm:"column:_raw_data_remark"`
}
